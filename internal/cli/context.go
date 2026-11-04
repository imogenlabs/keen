// Local customization: `keen context KEY` bundles what an agent needs to start
// on an issue in one call, within a token budget. See
// .printing-press-patches.json (keen-context).
//
// One hop from the issue: its parent, comments, linked issues, sub-tasks, the
// GitHub work `keen github links` finds, and Confluence pages linked to it.
// Items are ranked; when the bundle is over --budget the lowest-ranked are
// left out and listed under budget.omitted, never dropped silently.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type ctxIssue struct {
	Key         string   `json:"key"`
	Summary     string   `json:"summary"`
	Type        string   `json:"type,omitempty"`
	Status      string   `json:"status,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
	Reporter    string   `json:"reporter,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Created     string   `json:"created,omitempty"`
	Updated     string   `json:"updated,omitempty"`
	Description string   `json:"description,omitempty"`
}

type ctxRelated struct {
	Key      string `json:"key"`
	Relation string `json:"relation,omitempty"`
	Summary  string `json:"summary"`
	Status   string `json:"status,omitempty"`
}

type ctxComment struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Body    string `json:"body"`
}

type ctxPage struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text,omitempty"`
}

type ctxOmitted struct {
	Item   string `json:"item"`
	Tokens int    `json:"tokens"`
	Reason string `json:"reason"`
}

type ctxBundle struct {
	Issue      ctxIssue     `json:"issue"`
	Parent     *ctxIssue    `json:"parent,omitempty"`
	Comments   []ctxComment `json:"comments"`
	Linked     []ctxRelated `json:"linked_issues"`
	Subtasks   []ctxRelated `json:"subtasks"`
	GitHub     []githubLink `json:"github"`
	Confluence []ctxPage    `json:"confluence"`
	Budget     struct {
		Limit   int          `json:"limit"`
		Used    int          `json:"used"`
		Omitted []ctxOmitted `json:"omitted"`
	} `json:"budget"`
	Warnings []string `json:"warnings,omitempty"`
}

// approxTokens is the usual four-characters-per-token estimate.
func approxTokens(v any) int {
	b, _ := json.Marshal(v)
	return (len(b) + 3) / 4
}

var confluencePageID = regexp.MustCompile(`(?:pageId=|/pages/)([0-9]+)`)

// contextOutputFlags turns --compact off for the bundle. --budget already
// sized it, and compact's rule for big API objects (drop top-level comments
// and description) would throw away the comments under --agent.
func contextOutputFlags(flags *rootFlags) *rootFlags {
	f := *flags
	f.compact = false
	return &f
}

func newContextCmd(flags *rootFlags) *cobra.Command {
	var budget int
	var repos []string
	var org string
	cmd := &cobra.Command{
		Use:   "context <issue-key>",
		Short: "Start work on an issue: the issue with its parent epic, comments, linked issues, pull requests and Confluence pages, in one call within a token budget",
		Long: `Bundle one issue with what sits one hop from it: its parent, comments,
linked issues, sub-tasks, its GitHub pull requests and commits (as keen github
links finds them) and the Confluence pages linked to it.

The bundle is kept under --budget tokens (about four characters each). The
issue itself always comes first; then the parent, GitHub work, linked issues
and sub-tasks, the newest comments, Confluence pages and older comments. Items
that do not fit are listed under budget.omitted with their size, never dropped
silently. A part that cannot be read is a warning; the rest is still returned.`,
		Example: `  keen context PROJ-42
  keen context PROJ-42 --budget 4000 --agent
  keen context PROJ-42 --org my-org --agent --select issue,comments,budget`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToUpper(strings.TrimSpace(args[0]))
			if !issueKeyArg.MatchString(key) {
				return usageErr(fmt.Errorf("%q is not an issue key (expected PROJ-42)", args[0]))
			}
			if budget < 500 {
				return usageErr(fmt.Errorf("--budget must be at least 500 tokens"))
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			gl := &githubLinker{jira: c, githubAPI: "https://api.github.com",
				token: githubToken(), scopes: githubScopes(repos, org),
				http: &http.Client{Timeout: 30 * time.Second}}
			b, err := buildContext(cmd.Context(), c, gl, key, budget)
			if err != nil {
				return err
			}
			for _, w := range b.Warnings {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning:", w)
			}
			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), b, contextOutputFlags(flags))
			}
			return printContext(cmd.OutOrStdout(), b)
		},
	}
	cmd.Flags().IntVar(&budget, "budget", 8000, "Most tokens the bundle may use (about four characters each)")
	cmd.Flags().StringArrayVar(&repos, "repo", nil, "GitHub repo to search for the issue's work, owner/name (repeatable)")
	cmd.Flags().StringVar(&org, "org", "", "GitHub organization or user to search for the issue's work")
	return cmd
}

type jiraIssueJSON struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary     string                 `json:"summary"`
		Description json.RawMessage        `json:"description"`
		Labels      []string               `json:"labels"`
		Created     string                 `json:"created"`
		Updated     string                 `json:"updated"`
		IssueType   struct{ Name string }  `json:"issuetype"`
		Status      struct{ Name string }  `json:"status"`
		Priority    *struct{ Name string } `json:"priority"`
		Assignee    *struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Reporter *struct {
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Parent *struct {
			Key string `json:"key"`
		} `json:"parent"`
		Epic       string `json:"customfield_10014"`
		IssueLinks []struct {
			Type struct {
				Inward  string `json:"inward"`
				Outward string `json:"outward"`
			} `json:"type"`
			InwardIssue  *relatedJSON `json:"inwardIssue"`
			OutwardIssue *relatedJSON `json:"outwardIssue"`
		} `json:"issuelinks"`
		Subtasks []relatedJSON `json:"subtasks"`
	} `json:"fields"`
}

type relatedJSON struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string                `json:"summary"`
		Status  struct{ Name string } `json:"status"`
	} `json:"fields"`
}

func (j jiraIssueJSON) toIssue() ctxIssue {
	f := j.Fields
	i := ctxIssue{Key: j.Key, Summary: f.Summary, Type: f.IssueType.Name, Status: f.Status.Name,
		Labels: f.Labels, Created: f.Created, Updated: f.Updated, Description: adfText(f.Description)}
	if f.Priority != nil {
		i.Priority = f.Priority.Name
	}
	if f.Assignee != nil {
		i.Assignee = f.Assignee.DisplayName
	}
	if f.Reporter != nil {
		i.Reporter = f.Reporter.DisplayName
	}
	return i
}

const ctxIssueFields = "summary,description,labels,created,updated,issuetype,status,priority,assignee,reporter,parent,customfield_10014,issuelinks,subtasks"

func getIssueJSON(ctx context.Context, jira jiraReader, key string) (jiraIssueJSON, error) {
	var j jiraIssueJSON
	raw, err := jira.GetNoCache(ctx, "/rest/api/3/issue/"+url.PathEscape(key), map[string]string{"fields": ctxIssueFields})
	if err != nil {
		return j, err
	}
	return j, json.Unmarshal(raw, &j)
}

func buildContext(ctx context.Context, jira jiraReader, gl *githubLinker, key string, budget int) (*ctxBundle, error) {
	j, err := getIssueJSON(ctx, jira, key)
	if err != nil {
		return nil, err
	}
	b := &ctxBundle{Issue: j.toIssue(), Comments: []ctxComment{}, Linked: []ctxRelated{},
		Subtasks: []ctxRelated{}, GitHub: []githubLink{}, Confluence: []ctxPage{}}
	b.Budget.Limit = budget
	b.Budget.Omitted = []ctxOmitted{}
	warn := func(part string, err error) {
		if err != nil {
			b.Warnings = append(b.Warnings, part+": "+err.Error())
		}
	}

	// The parent is the parent field, or the older Epic Link field.
	parentKey := j.Fields.Epic
	if j.Fields.Parent != nil {
		parentKey = j.Fields.Parent.Key
	}
	var parent *ctxIssue
	if parentKey != "" {
		if pj, err := getIssueJSON(ctx, jira, parentKey); err != nil {
			warn("parent "+parentKey, err)
		} else {
			p := pj.toIssue()
			parent = &p
		}
	}
	var linked []ctxRelated
	for _, l := range j.Fields.IssueLinks {
		if l.OutwardIssue != nil {
			linked = append(linked, ctxRelated{Key: l.OutwardIssue.Key, Relation: l.Type.Outward,
				Summary: l.OutwardIssue.Fields.Summary, Status: l.OutwardIssue.Fields.Status.Name})
		}
		if l.InwardIssue != nil {
			linked = append(linked, ctxRelated{Key: l.InwardIssue.Key, Relation: l.Type.Inward,
				Summary: l.InwardIssue.Fields.Summary, Status: l.InwardIssue.Fields.Status.Name})
		}
	}
	var subtasks []ctxRelated
	for _, s := range j.Fields.Subtasks {
		subtasks = append(subtasks, ctxRelated{Key: s.Key, Summary: s.Fields.Summary, Status: s.Fields.Status.Name})
	}
	comments, err := getComments(ctx, jira, key)
	warn("comments", err)
	var github []githubLink
	if gl != nil {
		if err := gl.run(ctx, key); err != nil {
			warn("github", err)
		} else {
			github = gl.result.Links
			b.Warnings = append(b.Warnings, gl.result.Warnings...)
		}
	}
	pages, err := getConfluencePages(ctx, jira, key)
	warn("confluence", err)

	// Rank: the issue always; then parent, GitHub work, linked issues,
	// sub-tasks, the three newest comments, Confluence pages, older comments.
	type item struct {
		name string
		add  func()
		size int
	}
	var items []item
	if parent != nil {
		p := parent
		items = append(items, item{"parent " + p.Key, func() { b.Parent = p }, approxTokens(p)})
	}
	for _, g := range github {
		g := g
		name := "github " + g.Repo + " " + g.SHA
		if len(g.SHA) > 7 {
			name = "github " + g.Repo + " " + g.SHA[:7]
		}
		if g.Kind == "pull_request" {
			name = fmt.Sprintf("github %s#%d", g.Repo, g.Number)
		}
		items = append(items, item{name, func() { b.GitHub = append(b.GitHub, g) }, approxTokens(g)})
	}
	for _, r := range linked {
		r := r
		items = append(items, item{"linked " + r.Key, func() { b.Linked = append(b.Linked, r) }, approxTokens(r)})
	}
	for _, r := range subtasks {
		r := r
		items = append(items, item{"subtask " + r.Key, func() { b.Subtasks = append(b.Subtasks, r) }, approxTokens(r)})
	}
	// Newest first; the order is restored to oldest-first in the output.
	sort.SliceStable(comments, func(a, z int) bool { return comments[a].Created > comments[z].Created })
	commentItem := func(c ctxComment) item {
		return item{"comment " + c.Created + " by " + c.Author, func() { b.Comments = append(b.Comments, c) }, approxTokens(c)}
	}
	var older []item
	for i, c := range comments {
		if i < 3 {
			items = append(items, commentItem(c))
		} else {
			older = append(older, commentItem(c))
		}
	}
	for _, p := range pages {
		p := p
		items = append(items, item{"confluence " + p.Title, func() { b.Confluence = append(b.Confluence, p) }, approxTokens(p)})
	}
	items = append(items, older...)

	// The issue always goes in. Its description gets at most half the budget;
	// the cut is reported like any other omission.
	if keep := budget * 2; len(b.Issue.Description) > keep { // half the budget, at ~4 chars a token
		dropped := len(b.Issue.Description) - keep
		b.Issue.Description = strings.ToValidUTF8(b.Issue.Description[:keep], "") + " …"
		b.Budget.Omitted = append(b.Budget.Omitted, ctxOmitted{Item: "issue description (end)",
			Tokens: (dropped + 3) / 4, Reason: "description is over half the budget"})
	}
	// Fill greedily in rank order. The list of what was left out takes room
	// too, so measure the finished bundle and, while it is over, fill again
	// against a limit lowered by the overshoot.
	base := *b
	fill := func(limit int) {
		*b = base
		b.Budget.Omitted = append([]ctxOmitted{}, base.Budget.Omitted...)
		used := approxTokens(b)
		for _, it := range items {
			if used+it.size > limit {
				o := ctxOmitted{Item: it.name, Tokens: it.size, Reason: "over budget"}
				b.Budget.Omitted = append(b.Budget.Omitted, o)
				used += approxTokens(o)
				continue
			}
			it.add()
			used += it.size
		}
		sort.SliceStable(b.Comments, func(a, z int) bool { return b.Comments[a].Created < b.Comments[z].Created })
		b.Budget.Used = approxTokens(b)
	}
	for limit := budget; ; {
		fill(limit)
		if b.Budget.Used <= budget || limit <= 0 {
			break
		}
		limit -= b.Budget.Used - budget
	}
	if b.Budget.Used > budget {
		b.Warnings = append(b.Warnings, fmt.Sprintf(
			"the issue and the list of what was left out need %d tokens, over --budget %d", b.Budget.Used, budget))
	}
	return b, nil
}

func getComments(ctx context.Context, jira jiraReader, key string) ([]ctxComment, error) {
	var out []ctxComment
	for start := 0; ; {
		raw, err := jira.GetNoCache(ctx, "/rest/api/3/issue/"+url.PathEscape(key)+"/comment",
			map[string]string{"startAt": fmt.Sprint(start), "maxResults": "100", "orderBy": "created"})
		if err != nil {
			return out, err
		}
		var page struct {
			Total    int `json:"total"`
			Comments []struct {
				Author struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Created string          `json:"created"`
				Body    json.RawMessage `json:"body"`
			} `json:"comments"`
		}
		if err := json.Unmarshal(raw, &page); err != nil {
			return out, err
		}
		for _, c := range page.Comments {
			out = append(out, ctxComment{Author: c.Author.DisplayName, Created: c.Created, Body: adfText(c.Body)})
		}
		start += len(page.Comments)
		if len(page.Comments) == 0 || start >= page.Total {
			return out, nil
		}
	}
}

// getConfluencePages reads the Confluence pages linked to the issue: web links
// whose URL names a page id, read back from this same site.
func getConfluencePages(ctx context.Context, jira jiraReader, key string) ([]ctxPage, error) {
	raw, err := jira.GetNoCache(ctx, "/rest/api/3/issue/"+url.PathEscape(key)+"/remotelink", nil)
	if err != nil {
		return nil, err
	}
	var links []struct {
		Application struct {
			Type string `json:"type"`
		} `json:"application"`
		Object struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"object"`
	}
	if err := json.Unmarshal(raw, &links); err != nil {
		return nil, err
	}
	var pages []ctxPage
	seen := map[string]bool{}
	for _, l := range links {
		if l.Application.Type != "com.atlassian.confluence" && !strings.Contains(l.Object.URL, "/wiki/") {
			continue
		}
		m := confluencePageID.FindStringSubmatch(l.Object.URL)
		if m == nil || seen[m[1]] {
			continue
		}
		seen[m[1]] = true
		p := ctxPage{ID: m[1], Title: l.Object.Title, URL: l.Object.URL}
		if body, err := jira.GetNoCache(ctx, "/wiki/api/v2/pages/"+m[1], map[string]string{"body-format": "storage"}); err == nil {
			var page struct {
				Title string `json:"title"`
				Body  struct {
					Storage struct {
						Value string `json:"value"`
					} `json:"storage"`
				} `json:"body"`
			}
			if json.Unmarshal(body, &page) == nil {
				if page.Title != "" {
					p.Title = page.Title
				}
				p.Text = storageText(page.Body.Storage.Value)
			}
		}
		pages = append(pages, p)
	}
	return pages, nil
}

// adfText flattens an Atlassian Document Format value to plain text: one line
// per paragraph, heading or list item. A plain string is returned as is.
func adfText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var node any
	if json.Unmarshal(raw, &node) != nil {
		return ""
	}
	var sb strings.Builder
	var walk func(n any)
	walk = func(n any) {
		m, ok := n.(map[string]any)
		if !ok {
			return
		}
		switch m["type"] {
		case "text":
			t, _ := m["text"].(string)
			sb.WriteString(t)
		case "hardBreak":
			sb.WriteString("\n")
		case "mention", "emoji", "inlineCard":
			if a, ok := m["attrs"].(map[string]any); ok {
				for _, k := range []string{"text", "shortName", "url"} {
					if v, ok := a[k].(string); ok && v != "" {
						sb.WriteString(v)
						break
					}
				}
			}
		}
		if kids, ok := m["content"].([]any); ok {
			for _, k := range kids {
				walk(k)
			}
		}
		switch m["type"] {
		case "paragraph", "heading", "listItem", "codeBlock", "blockquote", "rule":
			sb.WriteString("\n")
		}
	}
	walk(node)
	return strings.TrimSpace(collapseBlankLines(sb.String()))
}

var (
	storageBlock = regexp.MustCompile(`(?i)</?(p|h[1-6]|li|tr|br|div|pre|blockquote)[^>]*>`)
	storageTag   = regexp.MustCompile(`<[^>]+>`)
)

// storageText turns Confluence storage format (XHTML) into plain text.
func storageText(s string) string {
	s = storageBlock.ReplaceAllString(s, "\n")
	s = storageTag.ReplaceAllString(s, "")
	return strings.TrimSpace(collapseBlankLines(html.UnescapeString(s)))
}

func collapseBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	out := lines[:0]
	blank := false
	for _, l := range lines {
		l = strings.TrimRight(l, " \t")
		if l == "" {
			if blank {
				continue
			}
			blank = true
		} else {
			blank = false
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

func printContext(w io.Writer, b *ctxBundle) error {
	p := func(format string, a ...any) { fmt.Fprintf(w, format, a...) }
	i := b.Issue
	p("%s  %s\n%s · %s", i.Key, i.Summary, i.Type, i.Status)
	if i.Assignee != "" {
		p(" · %s", i.Assignee)
	}
	p("\n")
	if i.Description != "" {
		p("\n%s\n", i.Description)
	}
	if b.Parent != nil {
		p("\nParent: %s  %s (%s)\n", b.Parent.Key, b.Parent.Summary, b.Parent.Status)
	}
	if len(b.GitHub) > 0 {
		p("\nGitHub:\n")
		for _, g := range b.GitHub {
			ref := fmt.Sprintf("#%d", g.Number)
			if g.Kind == "commit" && len(g.SHA) >= 7 {
				ref = g.SHA[:7]
			}
			p("  %s %s  %s\n", g.Repo, ref, g.Title)
		}
	}
	for _, sec := range []struct {
		title string
		items []ctxRelated
	}{{"Linked issues", b.Linked}, {"Sub-tasks", b.Subtasks}} {
		if len(sec.items) > 0 {
			p("\n%s:\n", sec.title)
			for _, r := range sec.items {
				p("  %s %s  %s (%s)\n", r.Relation, r.Key, r.Summary, r.Status)
			}
		}
	}
	for _, pg := range b.Confluence {
		p("\nConfluence: %s\n%s\n", pg.Title, pg.Text)
	}
	if len(b.Comments) > 0 {
		p("\nComments:\n")
		for _, c := range b.Comments {
			p("\n%s, %s:\n%s\n", c.Author, c.Created, c.Body)
		}
	}
	p("\n%d of %d tokens", b.Budget.Used, b.Budget.Limit)
	if n := len(b.Budget.Omitted); n > 0 {
		p("; left out %d item(s): ", n)
		var names []string
		for _, o := range b.Budget.Omitted {
			names = append(names, o.Item)
		}
		p("%s", strings.Join(names, ", "))
	}
	p("\n")
	return nil
}
