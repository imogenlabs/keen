package cli

import (
	"slices"
	"strings"
	"testing"
)

func TestSyncSkipsRemovedEndpoints(t *testing.T) {
	got := withoutRemoved(defaultSyncResources())
	if slices.Contains(got, "jira-cloud-platform-search") {
		t.Error("default sync still includes jira-cloud-platform-search (GET /rest/api/3/search, switched off)")
	}
	if !slices.Contains(got, "jira-cloud-platform-search-jql") {
		t.Error("default sync lost jira-cloud-platform-search-jql")
	}
}

func TestRemovedAndDeprecatedCommands(t *testing.T) {
	root := RootCmd()
	removed, _, err := root.Find([]string{"jira-cloud-platform-search", "for-issues-using-jql"})
	if err != nil || !removed.Hidden || !strings.Contains(removed.Short, "find-issues") {
		t.Errorf("removed search command: hidden=%v short=%q err=%v", removed.Hidden, removed.Short, err)
	}
	if err := removed.RunE(removed, nil); err == nil || !strings.Contains(err.Error(), "410 Gone") {
		t.Errorf("running the removed search command: %v", err)
	}
}
