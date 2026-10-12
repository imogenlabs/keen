package cli

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"keen/internal/client"
)

func TestExtractIssueKeys(t *testing.T) {
	got := sortedKeys(extractIssueKeys("PROJ", "PROJ-12: fix\n\nalso PROJ-7 and (PROJ-0009), PROJ-77/78\nPROJ-1 PROJ-2"))
	want := []string{"PROJ-1", "PROJ-2", "PROJ-7", "PROJ-9", "PROJ-12", "PROJ-77"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestExtractIssueKeysIgnoresLookalikes(t *testing.T) {
	got := extractIssueKeys("PROJ", "XPROJ-1 PROJ-2a PROJ-3x OTHER-4 PROJ- PROJ-0 SUB-PROJ-5 proj-6 PROJ_7")
	if len(got) != 0 {
		t.Fatalf("expected no keys, got %v", sortedKeys(got))
	}
}

func TestSortVersionTags(t *testing.T) {
	tags := []string{"v2.3.10", "v2.3.9", "2.3.2", "v10.0.0"}
	sortVersionTags(tags)
	if want := []string{"2.3.2", "v2.3.9", "v2.3.10", "v10.0.0"}; !reflect.DeepEqual(tags, want) {
		t.Fatalf("got %v, want %v", tags, want)
	}
}

func TestGithubRepoFromRemote(t *testing.T) {
	for in, want := range map[string]string{
		"git@github.com:acme/widget.git\n":   "acme/widget",
		"https://github.com/acme/widget":     "acme/widget",
		"https://github.com/acme/widget.git": "acme/widget",
		"https://example.com/acme/widget":    "",
	} {
		if got := githubRepoFromRemote(in); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}

// newTestRepo builds: PROJ-1 -> tag v1.0.9 -> PROJ-2 (body names PROJ-1 again)
// -> tags v1.0.10-rc1 and v1.0.10 -> PROJ-3 on the branch head only.
func newTestRepo(t *testing.T) gitRepo {
	t.Helper()
	dir := t.TempDir()
	env := append(os.Environ(), "GIT_AUTHOR_DATE=2026-09-01T12:00:00Z", "GIT_COMMITTER_DATE=2026-09-01T12:00:00Z")
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = env
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.name", "t")
	run("config", "user.email", "t@example.com")
	run("config", "tag.gpgSign", "false")
	run("config", "commit.gpgSign", "false")
	run("commit", "-q", "--allow-empty", "-m", "PROJ-1: first")
	run("tag", "v1.0.9")
	run("commit", "-q", "--allow-empty", "-m", "PROJ-2: second\n\nfollows PROJ-1")
	run("tag", "v1.0.10-rc1")
	run("tag", "v1.0.10")
	run("commit", "-q", "--allow-empty", "-m", "PROJ-3: unreleased")
	return gitRepo{dir: dir}
}

func TestReleaseMappingFirstContainingRelease(t *testing.T) {
	repo := newTestRepo(t)
	tags, err := repo.releaseTags(defaultReleaseTagPattern)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"v1.0.9", "v1.0.10"}; !reflect.DeepEqual(tags, want) {
		t.Fatalf("tags %v, want %v (pre-release tags must be left out)", tags, want)
	}
	m, err := repo.releaseMapping("PROJ", tags)
	if err != nil {
		t.Fatal(err)
	}
	want := []releaseEntry{
		{Tag: "v1.0.9", Date: "2026-09-01", Keys: []string{"PROJ-1"}},
		{Tag: "v1.0.10", Date: "2026-09-01", Keys: []string{"PROJ-2"}},
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("mapping %+v, want %+v", m, want)
	}
	next, err := repo.newKeys("PROJ", "v1.0.10", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(next, []string{"PROJ-3"}) {
		t.Fatalf("next %v, want [PROJ-3]", next)
	}
}

// fakeJira records writes and answers reads from fixed state.
type fakeJira struct {
	versions []jiraVersion
	issues   map[string][]string
	refuse   map[string]bool // issue keys whose edit is refused
	writes   []string
}

func (f *fakeJira) GetNoCache(_ context.Context, path string, _ map[string]string) (json.RawMessage, error) {
	switch {
	case strings.HasSuffix(path, "/versions"):
		return json.Marshal(f.versions)
	case strings.HasPrefix(path, "/rest/api/3/project/"):
		return json.RawMessage(`{"id":"100"}`), nil
	default:
		type fv struct {
			Name string `json:"name"`
		}
		var issues []map[string]any
		for k, names := range f.issues {
			var list []fv
			for _, n := range names {
				list = append(list, fv{n})
			}
			issues = append(issues, map[string]any{"key": k, "fields": map[string]any{"fixVersions": list}})
		}
		return json.Marshal(map[string]any{"issues": issues, "isLast": true})
	}
}

func (f *fakeJira) Post(_ context.Context, path string, body any) (json.RawMessage, int, error) {
	b, _ := json.Marshal(body)
	f.writes = append(f.writes, "POST "+path+" "+string(b))
	return json.RawMessage(`{"id":"9"}`), 201, nil
}

func (f *fakeJira) PutWithParams(_ context.Context, path string, _ map[string]string, body any) (json.RawMessage, int, error) {
	key := strings.TrimPrefix(path, "/rest/api/3/issue/")
	if f.refuse[key] {
		return nil, 400, &client.APIError{Method: "PUT", Path: path, StatusCode: 400, Body: "cannot be set"}
	}
	b, _ := json.Marshal(body)
	f.writes = append(f.writes, "PUT "+path+" "+string(b))
	return nil, 204, nil
}

func newTestSyncer(f *fakeJira) *releaseSyncer {
	s := &releaseSyncer{api: f, project: "PROJ", nextName: "next", urlFmt: "https://github.com/acme/widget/releases/tag/%s", out: io.Discard}
	_ = s.loadVersions(context.Background())
	_ = s.loadIssues(context.Background())
	return s
}

func TestReleaseRenamesNextAndAddsNeverReplaces(t *testing.T) {
	f := &fakeJira{
		versions: []jiraVersion{{ID: "7", Name: "next"}},
		issues: map[string][]string{
			"PROJ-1": {"next"},   // waited on next, is in the tag
			"PROJ-2": {"hotfix"}, // a person's Fix Version must stay
			"PROJ-4": {"next"},   // on next but not in the tag
			"PROJ-5": {},         // refused by the edit screen
			"PROJ-6": {"v2.0.0"}, // set by hand, never on next: untouched
		},
		refuse: map[string]bool{"PROJ-5": true},
	}
	s := newTestSyncer(f)
	err := s.release(context.Background(), releaseEntry{Tag: "v2.0.0", Date: "2026-10-01", Keys: []string{"PROJ-1", "PROJ-2", "PROJ-5", "PROJ-99"}}, true)
	if err != nil {
		t.Fatal(err)
	}
	all := strings.Join(f.writes, "\n")
	for _, want := range []string{
		`PUT /rest/api/3/version/7 {"name":"v2.0.0"}`,
		`PUT /rest/api/3/version/7 {"description":"https://github.com/acme/widget/releases/tag/v2.0.0","releaseDate":"2026-10-01","released":true}`,
		`PUT /rest/api/3/issue/PROJ-2 {"update":{"fixVersions":[{"add":{"name":"v2.0.0"}}]}}`,
		`POST /rest/api/3/version {"name":"next","projectId":100,"released":false}`,
		`PUT /rest/api/3/issue/PROJ-4 {"update":{"fixVersions":[{"add":{"name":"next"}}]}}`,
		`PUT /rest/api/3/issue/PROJ-4 {"update":{"fixVersions":[{"remove":{"name":"v2.0.0"}}]}}`,
	} {
		if !strings.Contains(all, want) {
			t.Errorf("missing write %s\nwrites:\n%s", want, all)
		}
	}
	if strings.Contains(all, `"set"`) || strings.Contains(all, "PROJ-1 ") || strings.Contains(all, "PROJ-6") {
		t.Errorf("a Fix Version was replaced or PROJ-1 re-edited:\n%s", all)
	}
	if s.tally.Added != 1 || s.tally.Already != 1 || !reflect.DeepEqual(s.tally.Missing, []string{"PROJ-99"}) || len(s.tally.Refused) != 1 {
		t.Errorf("tally %+v", s.tally)
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	f := &fakeJira{issues: map[string][]string{"PROJ-1": {}}}
	s := newTestSyncer(f)
	s.dryRun = true
	if err := s.release(context.Background(), releaseEntry{Tag: "v1.0.0", Date: "2026-01-01", Keys: []string{"PROJ-1"}}, false); err != nil {
		t.Fatal(err)
	}
	if len(f.writes) != 0 || s.tally.Added != 1 || len(s.tally.VersionsCreated) != 1 {
		t.Fatalf("writes %v tally %+v", f.writes, s.tally)
	}
}

func TestReleaseRerunTouchesNothing(t *testing.T) {
	f := &fakeJira{
		versions: []jiraVersion{
			{ID: "7", Name: "v2.0.0", Released: true, ReleaseDate: "2026-10-01", Description: "https://github.com/acme/widget/releases/tag/v2.0.0"},
			{ID: "8", Name: "next"},
		},
		issues: map[string][]string{"PROJ-1": {"v2.0.0"}, "PROJ-6": {"v2.0.0"}},
	}
	s := newTestSyncer(f)
	if err := s.release(context.Background(), releaseEntry{Tag: "v2.0.0", Date: "2026-10-01", Keys: []string{"PROJ-1"}}, true); err != nil {
		t.Fatal(err)
	}
	if len(f.writes) != 0 {
		t.Fatalf("a re-run must write nothing, wrote:\n%s", strings.Join(f.writes, "\n"))
	}
}
