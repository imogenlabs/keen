// Local customization: `keen release sync` records git release tags as Jira
// versions and adds them as Fix Versions on the issues their commits name.
// See .printing-press-patches.json.

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"keen/internal/client"
)

// releaseAPI is the slice of the Jira client release sync needs.
type releaseAPI interface {
	GetNoCache(ctx context.Context, path string, params map[string]string) (json.RawMessage, error)
	Post(ctx context.Context, path string, body any) (json.RawMessage, int, error)
	PutWithParams(ctx context.Context, path string, params map[string]string, body any) (json.RawMessage, int, error)
}

type jiraVersion struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate"`
	Description string `json:"description"`
}

type releaseTally struct {
	VersionsCreated []string `json:"versions_created"`
	VersionsUpdated []string `json:"versions_updated"`
	Added           int      `json:"fix_versions_added"`
	Already         int      `json:"fix_versions_already_set"`
	Missing         []string `json:"keys_not_in_project"`
	Refused         []string `json:"refused"`
	MovedToNext     []string `json:"moved_to_next,omitempty"`
}

type releaseSyncer struct {
	api      releaseAPI
	project  string
	nextName string
	urlFmt   string // printf format taking the tag, "" for no description
	dryRun   bool
	notify   bool
	out      io.Writer
	versions map[string]jiraVersion
	projID   string
	issues   map[string][]string // key -> fix version names
	tally    releaseTally
}

func (s *releaseSyncer) logf(format string, a ...any) { fmt.Fprintf(s.out, format+"\n", a...) }

func (s *releaseSyncer) loadVersions(ctx context.Context) error {
	pdata, err := s.api.GetNoCache(ctx, "/rest/api/3/project/"+url.PathEscape(s.project), nil)
	if err != nil {
		return fmt.Errorf("reading project %s: %w", s.project, err)
	}
	var proj struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(pdata, &proj); err != nil || proj.ID == "" {
		return fmt.Errorf("project %s has no id in Jira's answer", s.project)
	}
	s.projID = proj.ID
	data, err := s.api.GetNoCache(ctx, "/rest/api/3/project/"+url.PathEscape(s.project)+"/versions", nil)
	if err != nil {
		return fmt.Errorf("reading %s versions: %w", s.project, err)
	}
	var list []jiraVersion
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("parsing %s versions: %w", s.project, err)
	}
	s.versions = map[string]jiraVersion{}
	for _, v := range list {
		s.versions[v.Name] = v
	}
	return nil
}

func (s *releaseSyncer) loadIssues(ctx context.Context) error {
	s.issues = map[string][]string{}
	token := ""
	for {
		params := map[string]string{"jql": "project = " + s.project, "fields": "fixVersions", "maxResults": "100"}
		if token != "" {
			params["nextPageToken"] = token
		}
		data, err := s.api.GetNoCache(ctx, "/rest/api/3/search/jql", params)
		if err != nil {
			return fmt.Errorf("reading %s issues: %w", s.project, err)
		}
		var page struct {
			Issues []struct {
				Key    string `json:"key"`
				Fields struct {
					FixVersions []jiraVersion `json:"fixVersions"`
				} `json:"fields"`
			} `json:"issues"`
			NextPageToken string `json:"nextPageToken"`
			IsLast        bool   `json:"isLast"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return fmt.Errorf("parsing %s issues: %w", s.project, err)
		}
		for _, is := range page.Issues {
			names := []string{}
			for _, v := range is.Fields.FixVersions {
				names = append(names, v.Name)
			}
			s.issues[is.Key] = names
		}
		if page.NextPageToken == "" || page.IsLast {
			return nil
		}
		token = page.NextPageToken
	}
}

// ensureVersion creates the version or brings released/date/description in
// line with what the tag says. It never deletes a version.
func (s *releaseSyncer) ensureVersion(ctx context.Context, name string, released bool, date string) error {
	want := map[string]any{"released": released}
	if released {
		want["releaseDate"] = date
		if s.urlFmt != "" {
			want["description"] = fmt.Sprintf(s.urlFmt, name)
		}
	}
	v, ok := s.versions[name]
	if !ok {
		pid, _ := strconv.Atoi(s.projID)
		body := map[string]any{"name": name, "projectId": pid}
		for k, val := range want {
			body[k] = val
		}
		s.tally.VersionsCreated = append(s.tally.VersionsCreated, name)
		s.versions[name] = jiraVersion{Name: name, Released: released, ReleaseDate: date}
		if s.dryRun {
			s.logf("  would create version %s", name)
			return nil
		}
		data, _, err := s.api.Post(ctx, "/rest/api/3/version", body)
		if err != nil {
			return fmt.Errorf("creating version %s: %w", name, err)
		}
		var made jiraVersion
		_ = json.Unmarshal(data, &made)
		s.versions[name] = made
		s.logf("  created version %s", name)
		return nil
	}
	if v.Released == released && (!released || (v.ReleaseDate == date && (s.urlFmt == "" || v.Description == want["description"]))) {
		return nil
	}
	s.tally.VersionsUpdated = append(s.tally.VersionsUpdated, name)
	if s.dryRun {
		s.logf("  would update version %s", name)
		return nil
	}
	if _, _, err := s.api.PutWithParams(ctx, "/rest/api/3/version/"+v.ID, nil, want); err != nil {
		return fmt.Errorf("updating version %s: %w", name, err)
	}
	s.logf("  updated version %s", name)
	return nil
}

func (s *releaseSyncer) editFixVersion(ctx context.Context, key, op, name string) error {
	params := map[string]string{}
	if !s.notify {
		params["notifyUsers"] = "false"
	}
	body := map[string]any{"update": map[string]any{"fixVersions": []any{map[string]any{op: map[string]string{"name": name}}}}}
	_, _, err := s.api.PutWithParams(ctx, "/rest/api/3/issue/"+key, params, body)
	return err
}

// addFixVersion ADDS name to every key that exists in the project; it never
// replaces a Fix Version someone already set.
func (s *releaseSyncer) addFixVersion(ctx context.Context, keys []string, name string) {
	for _, key := range keys {
		have, ok := s.issues[key]
		switch {
		case !ok:
			s.tally.Missing = append(s.tally.Missing, key)
		case containsString(have, name):
			s.tally.Already++
		case s.dryRun:
			s.tally.Added++
			s.issues[key] = append(have, name)
		default:
			if err := s.editFixVersion(ctx, key, "add", name); err != nil {
				s.tally.Refused = append(s.tally.Refused, fmt.Sprintf("%s %s: %s", key, name, apiErrorSummary(err)))
				continue
			}
			s.tally.Added++
			s.issues[key] = append(have, name)
		}
	}
}

// release records one tag. With renameNext, work that waited on `next` and
// did not ship moves to a holding version first, then `next` is renamed to the
// tag, the tag's keys get it as a Fix Version, and the holding version becomes
// the fresh `next`. Every step can be repeated, so a run that stopped part-way
// is finished by running it again. A Fix Version a person set is never moved.
func (s *releaseSyncer) release(ctx context.Context, e releaseEntry, renameNext bool) error {
	renameNext = renameNext && s.nextName != ""
	pending := s.nextName + "-after-" + e.Tag
	if renameNext {
		if _, ok := s.versions[s.nextName]; ok {
			if _, tagged := s.versions[e.Tag]; !tagged {
				if err := s.holdUnshipped(ctx, e, pending); err != nil {
					return err
				}
				if err := s.renameVersion(ctx, s.nextName, e.Tag); err != nil {
					return err
				}
			}
		}
	}
	if err := s.ensureVersion(ctx, e.Tag, true, e.Date); err != nil {
		return err
	}
	s.addFixVersion(ctx, e.Keys, e.Tag)
	if !renameNext {
		return nil
	}
	if _, ok := s.versions[pending]; ok {
		if _, exists := s.versions[s.nextName]; !exists {
			return s.renameVersion(ctx, pending, s.nextName)
		}
	}
	return s.ensureVersion(ctx, s.nextName, false, "")
}

// holdUnshipped moves every issue on `next` that the tag does not contain to
// the holding version, so renaming `next` cannot release it.
func (s *releaseSyncer) holdUnshipped(ctx context.Context, e releaseEntry, pending string) error {
	inTag := map[string]bool{}
	for _, k := range e.Keys {
		inTag[k] = true
	}
	var hold []string
	for key, fv := range s.issues {
		if !inTag[key] && containsString(fv, s.nextName) {
			hold = append(hold, key)
		}
	}
	if len(hold) == 0 {
		return nil
	}
	if err := s.ensureVersion(ctx, pending, false, ""); err != nil {
		return err
	}
	for _, key := range sortedKeys(toSet(hold)) {
		s.tally.MovedToNext = append(s.tally.MovedToNext, key)
		if !s.dryRun {
			if err := s.editFixVersion(ctx, key, "add", pending); err != nil {
				return fmt.Errorf("holding %s back from %s: %s", key, e.Tag, apiErrorSummary(err))
			}
			if err := s.editFixVersion(ctx, key, "remove", s.nextName); err != nil {
				return fmt.Errorf("holding %s back from %s: %s", key, e.Tag, apiErrorSummary(err))
			}
		}
		fv := []string{}
		for _, n := range s.issues[key] {
			if n != s.nextName {
				fv = append(fv, n)
			}
		}
		s.issues[key] = append(fv, pending)
	}
	return nil
}

func (s *releaseSyncer) renameVersion(ctx context.Context, from, to string) error {
	v := s.versions[from]
	if s.dryRun {
		s.logf("  would rename %s to %s", from, to)
	} else if _, _, err := s.api.PutWithParams(ctx, "/rest/api/3/version/"+v.ID, nil, map[string]any{"name": to}); err != nil {
		return fmt.Errorf("renaming %s to %s: %w", from, to, err)
	} else {
		s.logf("  renamed %s to %s", from, to)
	}
	delete(s.versions, from)
	v.Name = to
	s.versions[to] = v
	for k, fv := range s.issues {
		for i, n := range fv {
			if n == from {
				s.issues[k][i] = to
			}
		}
	}
	return nil
}

func toSet(list []string) map[string]bool {
	set := make(map[string]bool, len(list))
	for _, v := range list {
		set[v] = true
	}
	return set
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func apiErrorSummary(err error) string {
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		body := strings.TrimSpace(apiErr.Body)
		if len(body) > 160 {
			body = body[:160] + "..."
		}
		return fmt.Sprintf("HTTP %d %s", apiErr.StatusCode, body)
	}
	return err.Error()
}

func newReleaseCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Record git release tags as Jira versions",
	}
	cmd.AddCommand(newReleaseSyncCmd(flags))
	return cmd
}

func newReleaseSyncCmd(flags *rootFlags) *cobra.Command {
	var repoPath, project, tag, head, tagPattern, nextName, githubRepo, report string
	var backfill, notify bool
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Create a Jira version per release tag and add it as Fix Version on the issues its commits name",
		Long: `For each release tag, create (or update) a released Jira version named as the
tag, dated with the tag's commit date and described by its GitHub release URL,
then ADD that version as a Fix Version on every issue of --project named in a
commit subject or body that the tag is the first release to contain.
Existing Fix Versions are never replaced. Keys that are not issues of the
project are skipped and counted.

An unreleased version (--next, default "next") holds the issues merged after
the newest tag. --tag renames it to the new tag, so a release picks up what
was already waiting, and then creates a fresh one.`,
		Example: `  keen release sync --repo-path . --project PROJ --backfill --head origin/dev --dry-run
  keen release sync --repo-path . --project PROJ --tag v1.4.0`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == "" || (tag == "") == !backfill {
				return usageErr(fmt.Errorf("need --project and exactly one of --tag or --backfill"))
			}
			repo := gitRepo{dir: repoPath}
			tags, err := repo.releaseTags(tagPattern)
			if err != nil {
				return err
			}
			if tag != "" && !containsString(tags, tag) {
				return fmt.Errorf("tag %q is not a release tag in %s (fetch tags, or check --tag-pattern)", tag, repoPath)
			}
			mapping, err := repo.releaseMapping(project, tags)
			if err != nil {
				return err
			}
			var nextEntry *releaseEntry
			if head != "" && nextName != "" && len(tags) > 0 {
				keys, err := repo.newKeys(project, tags[len(tags)-1], head)
				if err != nil {
					return err
				}
				nextEntry = &releaseEntry{Tag: nextName, Keys: keys}
			}
			if report != "" {
				all := mapping
				if nextEntry != nil {
					all = append(append([]releaseEntry{}, mapping...), *nextEntry)
				}
				data, _ := json.MarshalIndent(all, "", " ")
				if err := os.WriteFile(report, append(data, '\n'), 0o644); err != nil {
					return err
				}
			}
			if githubRepo == "" {
				if remote, err := repo.run("remote", "get-url", "origin"); err == nil {
					githubRepo = githubRepoFromRemote(remote)
				}
			}
			urlFmt := ""
			if githubRepo != "" {
				urlFmt = "https://github.com/" + githubRepo + "/releases/tag/%s"
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			c.DryRun = false // reads stay real; writes are skipped below
			s := &releaseSyncer{api: c, project: project, nextName: nextName, urlFmt: urlFmt,
				dryRun: flags.dryRun, notify: notify, out: cmd.ErrOrStderr()}
			ctx := cmd.Context()
			if err := s.loadVersions(ctx); err != nil {
				return err
			}
			if err := s.loadIssues(ctx); err != nil {
				return err
			}
			for _, e := range mapping {
				if backfill {
					if err := s.release(ctx, e, false); err != nil {
						return err
					}
				} else if e.Tag == tag {
					if err := s.release(ctx, e, true); err != nil {
						return err
					}
				}
			}
			if nextEntry != nil {
				if err := s.ensureVersion(ctx, nextName, false, ""); err != nil {
					return err
				}
				s.addFixVersion(ctx, nextEntry.Keys, nextName)
			}
			if err := flags.printJSON(cmd, s.tally); err != nil {
				return err
			}
			if len(s.tally.Refused) > 0 {
				return fmt.Errorf("%d fix version edit(s) refused", len(s.tally.Refused))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&repoPath, "repo-path", ".", "Git checkout holding the release tags")
	cmd.Flags().StringVar(&project, "project", "", "Jira project key whose issue keys to look for")
	cmd.Flags().StringVar(&tag, "tag", "", "Record this one release tag (run after publishing a release)")
	cmd.Flags().BoolVar(&backfill, "backfill", false, "Record every release tag")
	cmd.Flags().StringVar(&head, "head", "", "Ref holding merged-but-unreleased work for the next version (e.g. origin/main)")
	cmd.Flags().StringVar(&tagPattern, "tag-pattern", defaultReleaseTagPattern, "Regexp a release tag must match")
	cmd.Flags().StringVar(&nextName, "next", "next", "Name of the unreleased version (empty to skip it)")
	cmd.Flags().StringVar(&githubRepo, "github-repo", "", "owner/name for release URLs (default: read from the origin remote)")
	cmd.Flags().StringVar(&report, "report", "", "Write the tag -> date, keys mapping as JSON to this file")
	cmd.Flags().BoolVar(&notify, "notify", false, "Let Jira email watchers about each Fix Version edit")
	return cmd
}
