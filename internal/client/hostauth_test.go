// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"errors"
	"strings"
	"testing"

	"keen/internal/config"
)

func TestCredentialFor(t *testing.T) {
	cases := map[string]hostCredential{
		"https://example.atlassian.net/rest/api/3/myself":                   siteCredential,
		"https://example.atlassian.net/wiki/api/v2/pages":                   siteCredential,
		"https://api.atlassian.com/admin/v1/orgs":                           adminCredential,
		"https://api.atlassian.com/admin/control/v1/orgs/o/policies":        adminCredential,
		"https://api.atlassian.com/users/abc/manage/profile":                adminCredential,
		"https://api.atlassian.com/scim/directory/d/Users":                  scimCredential,
		"https://api.atlassian.com/ex/jira/c/jsm/assets/workspace/w/v1/aql": siteCredential,
		"https://marketplace.atlassian.com/rest/2/addons":                   siteCredential,
		"https://evil.example/admin/v1/orgs":                                siteCredential,
		"https://api.atlassian.com.evil.example/admin/v1/orgs":              siteCredential,
	}
	for u, want := range cases {
		if got := credentialFor(u); got != want {
			t.Errorf("credentialFor(%q) = %v, want %v", u, got, want)
		}
	}
}

func TestAuthHeaderFor(t *testing.T) {
	cfg := &config.Config{JiraCloudPlatformUsername: "me@example.com", JiraCloudPlatformPassword: "tok"}
	c := &Client{Config: cfg}
	ctx := context.Background()

	site, err := c.authHeaderFor(ctx, "https://example.atlassian.net/rest/api/3/myself")
	if err != nil || !strings.HasPrefix(site, "Basic ") {
		t.Fatalf("site header = %q, %v; want Basic", site, err)
	}

	_, err = c.authHeaderFor(ctx, "https://api.atlassian.com/admin/v1/orgs")
	if !errors.Is(err, ErrMissingKey) || !strings.Contains(err.Error(), "KEEN_ADMIN_API_KEY") {
		t.Fatalf("admin without key: err = %v; want ErrMissingKey naming KEEN_ADMIN_API_KEY", err)
	}
	_, err = c.authHeaderFor(ctx, "https://api.atlassian.com/scim/directory/d/Users")
	if !errors.Is(err, ErrMissingKey) || !strings.Contains(err.Error(), "KEEN_SCIM_API_KEY") {
		t.Fatalf("scim without key: err = %v; want ErrMissingKey naming KEEN_SCIM_API_KEY", err)
	}

	cfg.AdminAPIKey, cfg.ScimAPIKey = "admin-k", "scim-k"
	if h, _ := c.authHeaderFor(ctx, "https://api.atlassian.com/admin/v1/orgs"); h != "Bearer admin-k" {
		t.Errorf("admin header = %q, want Bearer admin-k", h)
	}
	if h, _ := c.authHeaderFor(ctx, "https://api.atlassian.com/scim/directory/d/Users"); h != "Bearer scim-k" {
		t.Errorf("scim header = %q, want Bearer scim-k", h)
	}
	if h, _ := c.authHeaderFor(ctx, "https://example.atlassian.net/rest/api/3/myself"); h != site {
		t.Errorf("site header changed when admin keys are set: %q", h)
	}
}

func TestCredentialHint(t *testing.T) {
	cases := map[string]string{
		"GET https://api.atlassian.com/admin/v1/orgs returned HTTP 401: {}":         "KEEN_ADMIN_API_KEY",
		"GET https://api.atlassian.com/scim/directory/d/Users returned HTTP 401":    "KEEN_SCIM_API_KEY",
		"GET https://example.atlassian.net/rest/api/3/myself returned HTTP 401: {}": "JIRA_CLOUD_PLATFORM_PASSWORD=<api-token>",
		"no url here": "JIRA_CLOUD_PLATFORM_USERNAME=<email>",
	}
	for msg, want := range cases {
		if got := CredentialHint(msg); !strings.Contains(got, want) {
			t.Errorf("CredentialHint(%q) = %q, want it to mention %q", msg, got, want)
		}
	}
}
