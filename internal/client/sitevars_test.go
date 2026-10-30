// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"keen/internal/config"
)

const assetsURL = "https://api.atlassian.com/ex/jira/KEEN_CLOUD_ID/jsm/assets/workspace/KEEN_ASSETS_WORKSPACE_ID/v1/objectschema/list"

func siteVarServer(t *testing.T, workspaces string, calls *int) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls++
		switch r.URL.Path {
		case "/_edge/tenant_info":
			_, _ = w.Write([]byte(`{"cloudId":"cloud-1"}`))
		case "/rest/servicedeskapi/assets/workspace":
			_, _ = w.Write([]byte(`{"values":[` + workspaces + `]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	cfg := &config.Config{JiraCloudPlatformUsername: "me@example.com", JiraCloudPlatformPassword: "tok"}
	return &Client{BaseURL: srv.URL, Config: cfg, HTTPClient: srv.Client()}
}

func TestFillSiteVarsLooksUpOnce(t *testing.T) {
	t.Setenv("KEEN_CLOUD_ID", "")
	t.Setenv("KEEN_ASSETS_WORKSPACE_ID", "")
	calls := 0
	c := siteVarServer(t, `{"workspaceId":"ws-1"}`, &calls)
	for i := 0; i < 3; i++ {
		got, err := c.fillSiteVars(context.Background(), assetsURL)
		if err != nil {
			t.Fatal(err)
		}
		if want := "https://api.atlassian.com/ex/jira/cloud-1/jsm/assets/workspace/ws-1/v1/objectschema/list"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
	if calls != 2 {
		t.Errorf("site lookups = %d, want 2 (one per id, then cached)", calls)
	}
}

func TestFillSiteVarsPrefersEnvAndConfig(t *testing.T) {
	t.Setenv("KEEN_CLOUD_ID", "cloud-env")
	t.Setenv("KEEN_ASSETS_WORKSPACE_ID", "")
	calls := 0
	c := siteVarServer(t, `{"workspaceId":"ws-site"}`, &calls)
	c.Config.AssetsWorkspaceID = "ws-config"
	got, err := c.fillSiteVars(context.Background(), assetsURL)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "/cloud-env/") || !strings.Contains(got, "/ws-config/") || calls != 0 {
		t.Errorf("got %s with %d lookups; want env cloud id, config workspace, no lookups", got, calls)
	}
}

func TestFillSiteVarsNoWorkspace(t *testing.T) {
	t.Setenv("KEEN_CLOUD_ID", "cloud-env")
	t.Setenv("KEEN_ASSETS_WORKSPACE_ID", "")
	calls := 0
	c := siteVarServer(t, ``, &calls)
	_, err := c.fillSiteVars(context.Background(), assetsURL)
	if err == nil || !strings.Contains(err.Error(), "export KEEN_ASSETS_WORKSPACE_ID") {
		t.Fatalf("err = %v; want a hint naming KEEN_ASSETS_WORKSPACE_ID", err)
	}
}

func TestFillSiteVarsLeavesOtherURLs(t *testing.T) {
	calls := 0
	c := siteVarServer(t, ``, &calls)
	u := "https://example.atlassian.net/rest/api/3/myself"
	if got, err := c.fillSiteVars(context.Background(), u); err != nil || got != u || calls != 0 {
		t.Errorf("got %s, %v, %d lookups; want the URL unchanged and no lookups", got, err, calls)
	}
}
