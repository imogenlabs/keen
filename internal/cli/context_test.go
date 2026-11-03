package cli

import (
	"context"
	"strings"
	"testing"
)

func contextJira() fakeJiraReader {
	long := strings.Repeat("Retry uploads that fail with a 5xx. ", 40)
	return fakeJiraReader{
		"/rest/api/3/issue/PROJ-42": `{"id":"1001","key":"PROJ-42","fields":{
			"summary":"Retry failed uploads","issuetype":{"name":"Story"},"status":{"name":"In Progress"},
			"priority":{"name":"High"},"assignee":{"displayName":"Sam"},"labels":["api"],
			"description":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"` + long + `"}]}]},
			"parent":{"key":"PROJ-1"},
			"issuelinks":[{"type":{"inward":"is blocked by","outward":"blocks"},"outwardIssue":{"key":"PROJ-50","fields":{"summary":"Ship uploads","status":{"name":"To Do"}}}}],
			"subtasks":[{"key":"PROJ-43","fields":{"summary":"Write the retry","status":{"name":"Done"}}}]}}`,
		"/rest/api/3/issue/PROJ-1": `{"id":"1000","key":"PROJ-1","fields":{"summary":"Uploads epic","issuetype":{"name":"Epic"},"status":{"name":"In Progress"}}}`,
		"/rest/api/3/issue/PROJ-42/comment": `{"total":4,"comments":[
			{"author":{"displayName":"Ana"},"created":"2026-11-01T10:00:00.000+0000","body":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"first"}]}]}},
			{"author":{"displayName":"Ana"},"created":"2026-11-02T10:00:00.000+0000","body":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"second"}]}]}},
			{"author":{"displayName":"Sam"},"created":"2026-11-03T10:00:00.000+0000","body":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"third"}]}]}},
			{"author":{"displayName":"Sam"},"created":"2026-11-04T10:00:00.000+0000","body":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"fourth"}]}]}}]}`,
		"/rest/api/3/issue/PROJ-42/remotelink": `[
			{"application":{"type":"com.atlassian.confluence"},"object":{"url":"https://x.atlassian.net/wiki/pages/viewpage.action?pageId=777","title":"Upload design"}},
			{"object":{"url":"https://github.com/acme/api/pull/7","title":"PR 7"}}]`,
		"/wiki/api/v2/pages/777":                 `{"title":"Upload design","body":{"storage":{"value":"<h1>Design</h1><p>Retry three times &amp; back off.</p>"}}}`,
		"/rest/dev-status/latest/issue/summary": emptyDevPanel,
		"/rest/api/3/issue/PROJ-7":              `{"id":"7","key":"PROJ-7","fields":{"summary":"Lonely","issuetype":{"name":"Task"},"status":{"name":"To Do"}}}`,
		"/rest/api/3/issue/PROJ-7/comment":      `{"total":0,"comments":[]}`,
		"/rest/api/3/issue/PROJ-7/remotelink":   `[]`,
	}
}

func TestContextBundlesEverything(t *testing.T) {
	jira := contextJira()
	gl := &githubLinker{jira: jira} // no token: GitHub search is a warning, web links still count
	b, err := buildContext(context.Background(), jira, gl, "PROJ-42", 8000)
	if err != nil {
		t.Fatal(err)
	}
	if b.Parent == nil || b.Parent.Key != "PROJ-1" {
		t.Fatalf("parent: %+v", b.Parent)
	}
	if len(b.Comments) != 4 || b.Comments[0].Body != "first" || b.Comments[3].Body != "fourth" {
		t.Fatalf("comments should all be there, oldest first: %+v", b.Comments)
	}
	if len(b.Linked) != 1 || b.Linked[0].Relation != "blocks" || len(b.Subtasks) != 1 {
		t.Fatalf("links: %+v %+v", b.Linked, b.Subtasks)
	}
	if len(b.GitHub) != 1 || b.GitHub[0].Number != 7 {
		t.Fatalf("github: %+v", b.GitHub)
	}
	if len(b.Confluence) != 1 || b.Confluence[0].Text != "Design\n\nRetry three times & back off." {
		t.Fatalf("confluence: %+v", b.Confluence)
	}
	if len(b.Budget.Omitted) != 0 || b.Budget.Used > 8000 {
		t.Fatalf("budget: %+v", b.Budget)
	}
	if !strings.Contains(strings.Join(b.Warnings, ";"), "no GitHub token") {
		t.Fatalf("expected the GitHub warning: %v", b.Warnings)
	}
}

func TestContextBudgetDropsLowestRankAndSaysSo(t *testing.T) {
	jira := contextJira()
	b, err := buildContext(context.Background(), jira, nil, "PROJ-42", 500)
	if err != nil {
		t.Fatal(err)
	}
	if b.Budget.Used > 500 {
		t.Fatalf("over budget: %d", b.Budget.Used)
	}
	if b.Issue.Key != "PROJ-42" || !strings.HasSuffix(b.Issue.Description, " …") {
		t.Fatalf("issue must stay, description cut: %q", b.Issue.Description[len(b.Issue.Description)-10:])
	}
	omitted := map[string]bool{}
	for _, o := range b.Budget.Omitted {
		if o.Tokens <= 0 || o.Reason == "" {
			t.Fatalf("omission without size or reason: %+v", o)
		}
		omitted[strings.Fields(o.Item)[0]] = true
	}
	if !omitted["issue"] {
		t.Fatalf("the description cut should be reported: %+v", b.Budget.Omitted)
	}
	// What did fit is the highest-ranked: the parent goes in before any comment.
	if len(b.Comments) > 0 && b.Parent == nil {
		t.Fatalf("a comment went in while the parent was left out: %+v", b.Budget.Omitted)
	}
	// Everything is either in the bundle or listed as left out.
	got := len(b.Comments) + len(b.Linked) + len(b.Subtasks) + len(b.Confluence) + len(b.GitHub)
	if b.Parent != nil {
		got++
	}
	if len(b.Budget.Omitted) < 3 {
		t.Fatalf("a 500-token budget should leave several items out: %+v", b.Budget.Omitted)
	}
	if got+len(b.Budget.Omitted)-1 != 8 { // parent, 4 comments, link, subtask, page (-1: the description cut)
		t.Fatalf("items lost: %d kept + %d omitted", got, len(b.Budget.Omitted))
	}
}

func TestContextIssueWithNothingAround(t *testing.T) {
	jira := contextJira()
	b, err := buildContext(context.Background(), jira, nil, "PROJ-7", 8000)
	if err != nil {
		t.Fatal(err)
	}
	if b.Issue.Summary != "Lonely" || b.Parent != nil || len(b.Comments)+len(b.Linked)+len(b.GitHub)+len(b.Confluence) != 0 || len(b.Warnings) != 0 {
		t.Fatalf("lonely issue: %+v", b)
	}
	if _, err := buildContext(context.Background(), jira, nil, "PROJ-404", 8000); err == nil {
		t.Fatal("a missing issue must be an error")
	}
}

func TestADFText(t *testing.T) {
	adf := `{"type":"doc","content":[
		{"type":"heading","content":[{"type":"text","text":"Steps"}]},
		{"type":"bulletList","content":[{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"one"}]}]}]},
		{"type":"paragraph","content":[{"type":"text","text":"see "},{"type":"inlineCard","attrs":{"url":"https://x/1"}},{"type":"hardBreak"},{"type":"mention","attrs":{"text":"@Sam"}}]}]}`
	if got := adfText([]byte(adf)); got != "Steps\none\n\nsee https://x/1\n@Sam" {
		t.Fatalf("adfText: %q", got)
	}
	if adfText([]byte(`"plain"`)) != "plain" || adfText(nil) != "" || adfText([]byte("null")) != "" {
		t.Fatal("plain/empty")
	}
}
