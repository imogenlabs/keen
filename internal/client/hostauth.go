// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

// Atlassian's admin APIs on api.atlassian.com do not take the site's email and
// API token. The organization APIs (/admin/..., /users/...) take an admin API
// key as a bearer token, and SCIM user provisioning (/scim/...) takes the key
// of one identity-provider directory. Bitbucket (api.bitbucket.org) takes an
// Atlassian API token with Bitbucket scopes, which the site token usually is
// not; without one keen sends no credential, so public repositories still
// work. Everything else uses the site login.

// ErrMissingKey is returned before a request is sent when the API it calls
// needs a key that is not configured.
var ErrMissingKey = errors.New("missing API key")

type hostCredential int

const (
	siteCredential hostCredential = iota
	adminCredential
	scimCredential
	bitbucketCredential
)

func credentialFor(rawURL string) hostCredential {
	u, err := url.Parse(rawURL)
	if err == nil && strings.EqualFold(u.Host, "api.bitbucket.org") {
		return bitbucketCredential
	}
	if err != nil || !strings.EqualFold(u.Host, "api.atlassian.com") {
		return siteCredential
	}
	switch {
	case strings.HasPrefix(u.Path, "/scim/"):
		return scimCredential
	case strings.HasPrefix(u.Path, "/admin/"), strings.HasPrefix(u.Path, "/users/"):
		return adminCredential
	}
	return siteCredential
}

// authHeaderFor returns the Authorization header for a request to rawURL.
func (c *Client) authHeaderFor(ctx context.Context, rawURL string) (string, error) {
	if c.Config == nil {
		return "", nil
	}
	switch credentialFor(rawURL) {
	case adminCredential:
		if c.Config.AdminAPIKey == "" {
			return "", errors.Join(ErrMissingKey, errors.New(adminKeyHint))
		}
		return "Bearer " + c.Config.AdminAPIKey, nil
	case scimCredential:
		if c.Config.ScimAPIKey == "" {
			return "", errors.Join(ErrMissingKey, errors.New(scimKeyHint))
		}
		return "Bearer " + c.Config.ScimAPIKey, nil
	case bitbucketCredential:
		if c.Config.BitbucketToken == "" {
			return "", nil
		}
		email := c.Config.BitbucketEmail
		if email == "" {
			email = c.Config.JiraCloudPlatformUsername
		}
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(email+":"+c.Config.BitbucketToken)), nil
	}
	return c.authHeader(ctx)
}

const (
	adminKeyHint = "this Atlassian admin API needs an organization admin API key (admin.atlassian.com > Settings > API keys)." +
		"\n      Set it with: export KEEN_ADMIN_API_KEY=<key>, or admin_api_key in ~/.config/keen/auth.json"
	scimKeyHint = "this user provisioning (SCIM) API needs the API key of your identity-provider directory" +
		" (admin.atlassian.com > Security > Identity providers)." +
		"\n      Set it with: export KEEN_SCIM_API_KEY=<key>, or scim_api_key in ~/.config/keen/auth.json"
	siteKeyHint = "Set your Atlassian email and API token:" +
		" export JIRA_CLOUD_PLATFORM_USERNAME=<email> JIRA_CLOUD_PLATFORM_PASSWORD=<api-token>"
)

// CredentialHint names the credential an API error message's URL needs, for
// the hint printed after a 400, 401 or 403.
func CredentialHint(msg string) string {
	for _, field := range strings.Fields(msg) {
		if strings.HasPrefix(field, "https://") {
			switch credentialFor(field) {
			case adminCredential:
				return "Set it with: export KEEN_ADMIN_API_KEY=<organization admin API key>"
			case scimCredential:
				return "Set it with: export KEEN_SCIM_API_KEY=<directory API key>"
			case bitbucketCredential:
				return "Set it with: export KEEN_BITBUCKET_TOKEN=<Atlassian API token with Bitbucket scopes>"
			}
			return siteKeyHint
		}
	}
	return siteKeyHint
}
