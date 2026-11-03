package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeJiraReader answers by path; a missing path is an error, as a 404 would be.
type fakeJiraReader map[string]string

func (f fakeJiraReader) GetNoCache(_ context.Context, path string, params map[string]string) (json.RawMessage, error) {
	if params["dataType"] != "" {
		path += "?" + params["dataType"]
	}
	if body, ok := f[path]; ok {
		return json.RawMessage(body), nil
	}
	return nil, fmt.Errorf("HTTP 404 for %s", path)
}

const emptyDevPanel = `{"summary":{"pullrequest":{"byInstanceType":{}},"repository":{"byInstanceType":{}}}}`

func githubServer(t *testing.T, prs, commits string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/search/issues":
			fmt.Fprint(w, prs)
		case "/search/commits":
			fmt.Fprint(w, commits)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestGithubLinksMergesSources(t *testing.T) {
	jira := fakeJiraReader{
		"/rest/api/3/issue/PROJ-42":             `{"id":"1001"}`,
		"/rest/dev-status/latest/issue/summary": `{"summary":{"pullrequest":{"byInstanceType":{"GitHub":{}}}}}`,
		"/rest/dev-status/latest/issue/detail?pullrequest": `{"detail":[{"pullRequests":[
			{"url":"https://github.com/acme/api/pull/7","name":"fix: PROJ-42 retry","status":"MERGED","source":{"branch":"fix/PROJ-42-retry"}}]}]}`,
		"/rest/dev-status/latest/issue/detail?repository": `{"detail":[{"repositories":[{"name":"acme/api","commits":[
			{"id":"abc1234def","url":"https://github.com/acme/api/commit/abc1234def","message":"fix: PROJ-42 retry\n\nbody"}]}]}]}`,
		"/rest/api/3/issue/PROJ-42/remotelink": `[{"object":{"url":"https://github.com/acme/api/pull/7","title":"PR 7"}},
			{"object":{"url":"https://github.com/acme/web/pull/3","title":"web PR"}},
			{"object":{"url":"https://example.com/doc","title":"not GitHub"}}]`,
	}
	srv := githubServer(t,
		`{"items":[
			{"html_url":"https://github.com/acme/api/pull/7","number":7,"title":"fix: PROJ-42 retry","state":"closed","repository_url":"https://api.github.com/repos/acme/api","pull_request":{"merged_at":"2026-11-01T00:00:00Z"}},
			{"html_url":"https://github.com/acme/api/pull/9","number":9,"title":"feat: PROJ-420 other","state":"open","repository_url":"https://api.github.com/repos/acme/api","pull_request":{}}]}`,
		`{"items":[
			{"sha":"abc1234def","html_url":"https://github.com/acme/api/commit/abc1234def","repository":{"full_name":"acme/api"},"commit":{"message":"fix: PROJ-42 retry"}},
			{"sha":"fff0000aaa","html_url":"https://github.com/acme/api/commit/fff0000aaa","repository":{"full_name":"acme/api"},"commit":{"message":"chore: PROJ-4 unrelated"}}]}`)

	l := &githubLinker{jira: jira, githubAPI: srv.URL, token: "tok", scopes: []string{"org:acme"}, http: srv.Client()}
	if err := l.run(context.Background(), "PROJ-42"); err != nil {
		t.Fatal(err)
	}
	r := l.result
	if len(r.Warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", r.Warnings)
	}
	// PR 7 found three ways, once; web PR 3 from the web link; PR 9 (PROJ-420)
	// and the PROJ-4 commit are not this issue; one commit, found twice.
	if len(r.Links) != 3 {
		t.Fatalf("want 3 links, got %d: %+v", len(r.Links), r.Links)
	}
	pr7 := r.Links[0]
	if pr7.Repo != "acme/api" || pr7.Number != 7 || pr7.Branch != "fix/PROJ-42-retry" ||
		strings.Join(pr7.Sources, ",") != "jira-dev-panel,github-search,jira-remote-link" {
		t.Fatalf("PR 7 wrong: %+v", pr7)
	}
	if r.Links[1].Repo != "acme/web" || r.Links[1].Number != 3 {
		t.Fatalf("web PR wrong: %+v", r.Links[1])
	}
	c := r.Links[2]
	if c.Kind != "commit" || c.Title != "fix: PROJ-42 retry" || strings.Join(c.Sources, ",") != "jira-dev-panel,github-search" {
		t.Fatalf("commit wrong: %+v", c)
	}
}

func TestGithubLinksSourceFailuresAreWarnings(t *testing.T) {
	jira := fakeJiraReader{
		"/rest/api/3/issue/PROJ-42":             `{"id":"1001"}`,
		"/rest/dev-status/latest/issue/summary": emptyDevPanel,
		"/rest/api/3/issue/PROJ-42/remotelink":  `[{"object":{"url":"https://github.com/acme/api/pull/7"}}]`,
	}
	// No token: GitHub is skipped with a warning; the web link still counts.
	l := &githubLinker{jira: jira, scopes: []string{"org:acme"}}
	if err := l.run(context.Background(), "PROJ-42"); err != nil {
		t.Fatal(err)
	}
	if len(l.result.Links) != 1 || len(l.result.Warnings) != 1 || !strings.Contains(l.result.Warnings[0], "no GitHub token") {
		t.Fatalf("no token: %+v", l.result)
	}
	// GitHub refuses (403): a warning naming the status, same result.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message":"rate limited"}`)
	}))
	defer srv.Close()
	l = &githubLinker{jira: jira, githubAPI: srv.URL, token: "tok", scopes: []string{"org:acme"}, http: srv.Client()}
	if err := l.run(context.Background(), "PROJ-42"); err != nil {
		t.Fatal(err)
	}
	if len(l.result.Links) != 1 || !strings.Contains(strings.Join(l.result.Warnings, ";"), "HTTP 403: rate limited") {
		t.Fatalf("403: %+v", l.result)
	}
	// An issue that does not exist is an error, not an empty list.
	l = &githubLinker{jira: jira}
	if err := l.run(context.Background(), "PROJ-1"); err == nil {
		t.Fatal("missing issue should fail")
	}
}

func TestGithubScopesFromFlags(t *testing.T) {
	got := githubScopes([]string{"acme/api", " acme/web "}, "acme")
	if strings.Join(got, " ") != "repo:acme/api repo:acme/web org:acme" {
		t.Fatalf("scopes: %v", got)
	}
}
