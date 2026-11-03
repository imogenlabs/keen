// Local customization: `keen github links KEY` lists a Jira issue's GitHub pull
// requests and commits, looked up live. See .printing-press-patches.json
// (github-links).
//
// Three sources, each optional, merged so every PR and commit appears once:
//   - Jira's development panel (dev-status), filled when the GitHub for Jira
//     app is connected;
//   - GitHub search for the key in PR titles and bodies and in commit messages,
//     the way Atlassian links work (fix/KEY-12-slug, "fix: KEY-12 subject");
//   - web links on the issue that point at a GitHub PR.
// A source that fails is a warning, never a failed command.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type githubLink struct {
	Kind    string   `json:"kind"` // "pull_request" or "commit"
	URL     string   `json:"url"`
	Repo    string   `json:"repo,omitempty"`
	Number  int      `json:"number,omitempty"`
	SHA     string   `json:"sha,omitempty"`
	Title   string   `json:"title,omitempty"`
	State   string   `json:"state,omitempty"`
	Branch  string   `json:"branch,omitempty"`
	Sources []string `json:"sources"`
}

type githubLinksResult struct {
	Issue    string       `json:"issue"`
	Links    []githubLink `json:"links"`
	Warnings []string     `json:"warnings,omitempty"`
}

type jiraReader interface {
	GetNoCache(ctx context.Context, path string, params map[string]string) (json.RawMessage, error)
}

type githubLinker struct {
	jira      jiraReader
	githubAPI string // https://api.github.com, a test server in tests
	token     string
	scopes    []string // search qualifiers: "repo:o/r" or "org:o"
	http      *http.Client
	result    githubLinksResult
	byID      map[string]int
}

var (
	issueKeyArg   = regexp.MustCompile(`^[A-Z][A-Z0-9]+-[0-9]+$`)
	githubPullURL = regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/pull/([0-9]+)`)
	githubRepoURL = regexp.MustCompile(`/repos/([^/]+/[^/]+)$`)
)

func newGithubCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github",
		Short: "Find the GitHub work behind Jira issues",
	}
	cmd.AddCommand(newGithubLinksCmd(flags))
	return cmd
}

func newGithubLinksCmd(flags *rootFlags) *cobra.Command {
	var repos []string
	var org string
	cmd := &cobra.Command{
		Use:   "links <issue-key>",
		Short: "List a Jira issue's GitHub pull requests and commits",
		Long: `List the GitHub pull requests and commits for one Jira issue, from three
sources, each PR and commit once:

  jira-dev-panel    Jira's development panel (needs the GitHub for Jira app)
  github-search     PR titles and bodies and commit messages that name the key,
                    in --repo, --org, or the current checkout's GitHub origin
  jira-remote-link  web links on the issue that point at a GitHub PR

GitHub search uses KEEN_GITHUB_TOKEN, GH_TOKEN or GITHUB_TOKEN, or the gh CLI's
login. A source that is not set up, or fails, is reported as a warning and the
others still run.`,
		Example: `  keen github links PROJ-42
  keen github links PROJ-42 --org my-org --agent
  keen github links PROJ-42 --repo my-org/api --repo my-org/web`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToUpper(strings.TrimSpace(args[0]))
			if !issueKeyArg.MatchString(key) {
				return usageErr(fmt.Errorf("%q is not an issue key (expected PROJ-42)", args[0]))
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			l := &githubLinker{jira: c, githubAPI: "https://api.github.com",
				token: githubToken(), scopes: githubScopes(repos, org),
				http: &http.Client{Timeout: 30 * time.Second}}
			if err := l.run(cmd.Context(), key); err != nil {
				return err
			}
			for _, w := range l.result.Warnings {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning:", w)
			}
			if flags.asJSON {
				return flags.printJSON(cmd, l.result)
			}
			return printGithubLinks(cmd.OutOrStdout(), l.result)
		},
	}
	cmd.Flags().StringArrayVar(&repos, "repo", nil, "GitHub repo to search, owner/name (repeatable)")
	cmd.Flags().StringVar(&org, "org", "", "GitHub organization or user to search")
	return cmd
}

// githubToken finds a token without asking: keen's own variable, the two the
// gh CLI and Actions use, then the gh CLI's stored login.
func githubToken() string {
	for _, v := range []string{"KEEN_GITHUB_TOKEN", "GH_TOKEN", "GITHUB_TOKEN"} {
		if t := strings.TrimSpace(os.Getenv(v)); t != "" {
			return t
		}
	}
	if out, err := exec.Command("gh", "auth", "token").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

// githubScopes turns --repo and --org into search qualifiers. With neither, it
// searches the GitHub repo the current directory is a checkout of.
func githubScopes(repos []string, org string) []string {
	var scopes []string
	for _, r := range repos {
		scopes = append(scopes, "repo:"+strings.TrimSpace(r))
	}
	if org != "" {
		scopes = append(scopes, "org:"+strings.TrimSpace(org))
	}
	if len(scopes) == 0 {
		if remote, err := (gitRepo{dir: "."}).run("remote", "get-url", "origin"); err == nil {
			if r := githubRepoFromRemote(remote); r != "" {
				scopes = append(scopes, "repo:"+r)
			}
		}
	}
	return scopes
}

func (l *githubLinker) run(ctx context.Context, key string) error {
	l.result = githubLinksResult{Issue: key, Links: []githubLink{}}
	l.byID = map[string]int{}
	raw, err := l.jira.GetNoCache(ctx, "/rest/api/3/issue/"+url.PathEscape(key), map[string]string{"fields": "summary"})
	if err != nil {
		return err
	}
	var issue struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &issue); err != nil || issue.ID == "" {
		return fmt.Errorf("could not read issue %s", key)
	}
	l.warn("jira dev panel", l.fromDevPanel(ctx, issue.ID))
	l.warn("github search", l.fromGithubSearch(ctx, key))
	l.warn("jira web links", l.fromRemoteLinks(ctx, key))
	sort.SliceStable(l.result.Links, func(i, j int) bool {
		a, b := l.result.Links[i], l.result.Links[j]
		if a.Kind != b.Kind {
			return a.Kind == "pull_request"
		}
		if a.Repo != b.Repo {
			return a.Repo < b.Repo
		}
		return a.Number < b.Number
	})
	return nil
}

func (l *githubLinker) warn(source string, err error) {
	if err != nil {
		l.result.Warnings = append(l.result.Warnings, source+": "+err.Error())
	}
}

// add records a link, or adds the source to one already found.
func (l *githubLinker) add(g githubLink, source string) {
	id := g.SHA
	if g.Kind == "pull_request" {
		id = strings.ToLower(g.Repo) + "#" + strconv.Itoa(g.Number)
	}
	if i, ok := l.byID[id]; ok {
		have := &l.result.Links[i]
		for _, s := range have.Sources {
			if s == source {
				return
			}
		}
		have.Sources = append(have.Sources, source)
		if have.Title == "" {
			have.Title = g.Title
		}
		if have.State == "" {
			have.State = g.State
		}
		if have.Branch == "" {
			have.Branch = g.Branch
		}
		return
	}
	g.Sources = []string{source}
	l.byID[id] = len(l.result.Links)
	l.result.Links = append(l.result.Links, g)
}

func (l *githubLinker) fromDevPanel(ctx context.Context, issueID string) error {
	raw, err := l.jira.GetNoCache(ctx, "/rest/dev-status/latest/issue/summary", map[string]string{"issueId": issueID})
	if err != nil {
		return err
	}
	var summary struct {
		Summary map[string]struct {
			ByInstanceType map[string]json.RawMessage `json:"byInstanceType"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(raw, &summary); err != nil {
		return err
	}
	// The app's instance type (for example "GitHub") is read, not assumed.
	types := map[string]bool{}
	for _, data := range []string{"pullrequest", "repository"} {
		for t := range summary.Summary[data].ByInstanceType {
			types[t] = true
		}
	}
	for t := range types {
		if err := l.devPanelPulls(ctx, issueID, t); err != nil {
			return err
		}
		if err := l.devPanelCommits(ctx, issueID, t); err != nil {
			return err
		}
	}
	return nil
}

func (l *githubLinker) devPanelPulls(ctx context.Context, issueID, appType string) error {
	raw, err := l.jira.GetNoCache(ctx, "/rest/dev-status/latest/issue/detail",
		map[string]string{"issueId": issueID, "applicationType": appType, "dataType": "pullrequest"})
	if err != nil {
		return err
	}
	var detail struct {
		Detail []struct {
			PullRequests []struct {
				URL    string `json:"url"`
				Name   string `json:"name"`
				Status string `json:"status"`
				Source struct {
					Branch string `json:"branch"`
				} `json:"source"`
			} `json:"pullRequests"`
		} `json:"detail"`
	}
	if err := json.Unmarshal(raw, &detail); err != nil {
		return err
	}
	for _, d := range detail.Detail {
		for _, pr := range d.PullRequests {
			if m := githubPullURL.FindStringSubmatch(pr.URL); m != nil {
				n, _ := strconv.Atoi(m[2])
				l.add(githubLink{Kind: "pull_request", URL: pr.URL, Repo: m[1], Number: n,
					Title: pr.Name, State: strings.ToLower(pr.Status), Branch: pr.Source.Branch}, "jira-dev-panel")
			}
		}
	}
	return nil
}

func (l *githubLinker) devPanelCommits(ctx context.Context, issueID, appType string) error {
	raw, err := l.jira.GetNoCache(ctx, "/rest/dev-status/latest/issue/detail",
		map[string]string{"issueId": issueID, "applicationType": appType, "dataType": "repository"})
	if err != nil {
		return err
	}
	var detail struct {
		Detail []struct {
			Repositories []struct {
				Name    string `json:"name"`
				Commits []struct {
					ID      string `json:"id"`
					URL     string `json:"url"`
					Message string `json:"message"`
				} `json:"commits"`
			} `json:"repositories"`
		} `json:"detail"`
	}
	if err := json.Unmarshal(raw, &detail); err != nil {
		return err
	}
	for _, d := range detail.Detail {
		for _, r := range d.Repositories {
			for _, c := range r.Commits {
				l.add(githubLink{Kind: "commit", URL: c.URL, Repo: r.Name, SHA: c.ID,
					Title: firstLine(c.Message)}, "jira-dev-panel")
			}
		}
	}
	return nil
}

func (l *githubLinker) fromGithubSearch(ctx context.Context, key string) error {
	if l.token == "" {
		return fmt.Errorf("no GitHub token (set KEEN_GITHUB_TOKEN or GH_TOKEN, or run gh auth login); skipped")
	}
	if len(l.scopes) == 0 {
		return fmt.Errorf("no --repo or --org, and this directory is not a GitHub checkout; skipped")
	}
	// GitHub's search matches words loosely ("PROJ-4" finds "PROJ-42"), so
	// every hit is checked for the exact key before it counts.
	project := key[:strings.LastIndex(key, "-")]
	re := issueKeyPattern(project)
	names := func(text string) bool {
		for _, m := range re.FindAllStringSubmatch(text, -1) {
			if project+"-"+m[1] == key {
				return true
			}
		}
		return false
	}
	scope := strings.Join(l.scopes, " ")

	var prs struct {
		Items []struct {
			HTMLURL       string `json:"html_url"`
			Number        int    `json:"number"`
			Title         string `json:"title"`
			Body          string `json:"body"`
			State         string `json:"state"`
			RepositoryURL string `json:"repository_url"`
			PullRequest   struct {
				MergedAt *string `json:"merged_at"`
			} `json:"pull_request"`
		} `json:"items"`
	}
	if err := l.githubGet(ctx, "/search/issues", `"`+key+`" type:pr `+scope, &prs); err != nil {
		return err
	}
	for _, it := range prs.Items {
		if !names(it.Title + "\n" + it.Body) {
			continue
		}
		repo := ""
		if m := githubRepoURL.FindStringSubmatch(it.RepositoryURL); m != nil {
			repo = m[1]
		}
		state := it.State
		if it.PullRequest.MergedAt != nil {
			state = "merged"
		}
		l.add(githubLink{Kind: "pull_request", URL: it.HTMLURL, Repo: repo, Number: it.Number,
			Title: it.Title, State: state}, "github-search")
	}

	var commits struct {
		Items []struct {
			SHA        string `json:"sha"`
			HTMLURL    string `json:"html_url"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
			Commit struct {
				Message string `json:"message"`
			} `json:"commit"`
		} `json:"items"`
	}
	if err := l.githubGet(ctx, "/search/commits", `"`+key+`" `+scope, &commits); err != nil {
		return err
	}
	for _, it := range commits.Items {
		if !names(it.Commit.Message) {
			continue
		}
		l.add(githubLink{Kind: "commit", URL: it.HTMLURL, Repo: it.Repository.FullName, SHA: it.SHA,
			Title: firstLine(it.Commit.Message)}, "github-search")
	}
	return nil
}

func (l *githubLinker) githubGet(ctx context.Context, path, query string, out any) error {
	u := l.githubAPI + path + "?per_page=100&q=" + url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+l.token)
	resp, err := l.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(body, &e)
		return fmt.Errorf("GitHub %s returned HTTP %d: %s", path, resp.StatusCode, e.Message)
	}
	return json.Unmarshal(body, out)
}

func (l *githubLinker) fromRemoteLinks(ctx context.Context, key string) error {
	raw, err := l.jira.GetNoCache(ctx, "/rest/api/3/issue/"+url.PathEscape(key)+"/remotelink", nil)
	if err != nil {
		return err
	}
	var links []struct {
		Object struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"object"`
	}
	if err := json.Unmarshal(raw, &links); err != nil {
		return err
	}
	for _, rl := range links {
		if m := githubPullURL.FindStringSubmatch(rl.Object.URL); m != nil {
			n, _ := strconv.Atoi(m[2])
			l.add(githubLink{Kind: "pull_request", URL: "https://github.com/" + m[1] + "/pull/" + m[2],
				Repo: m[1], Number: n, Title: rl.Object.Title}, "jira-remote-link")
		}
	}
	return nil
}

func printGithubLinks(w io.Writer, r githubLinksResult) error {
	if len(r.Links) == 0 {
		_, err := fmt.Fprintf(w, "No GitHub pull requests or commits found for %s.\n", r.Issue)
		return err
	}
	for _, g := range r.Links {
		ref := "#" + strconv.Itoa(g.Number)
		if g.Kind == "commit" && len(g.SHA) >= 7 {
			ref = g.SHA[:7]
		}
		state := g.State
		if state == "" {
			state = "-"
		}
		if _, err := fmt.Fprintf(w, "%-30s %-9s %-7s %s  [%s]\n", g.Repo, ref, state, g.Title,
			strings.Join(g.Sources, ", ")); err != nil {
			return err
		}
	}
	return nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}
