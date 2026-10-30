// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Some Atlassian APIs put site ids in the URL: Assets is
// api.atlassian.com/ex/jira/{cloudId}/jsm/assets/workspace/{workspaceId}/v1.
// The refresh writes a marker in their place (scripts/regen/specfix.py,
// TEMPLATE_VARS); fillSiteVars swaps each marker for its value, taken from the
// env var of the same name, then auth.json, then a lookup on the site.

type siteVar struct {
	marker string
	config func(c *Client) string
	lookup func(ctx context.Context, c *Client) (string, error)
}

var siteVars = []siteVar{
	{
		marker: "KEEN_CLOUD_ID",
		config: func(c *Client) string { return c.Config.CloudID },
		lookup: func(ctx context.Context, c *Client) (string, error) {
			var v struct {
				CloudID string `json:"cloudId"`
			}
			if err := c.siteGet(ctx, "/_edge/tenant_info", &v); err != nil {
				return "", err
			}
			return v.CloudID, nil
		},
	},
	{
		marker: "KEEN_ASSETS_WORKSPACE_ID",
		config: func(c *Client) string { return c.Config.AssetsWorkspaceID },
		lookup: func(ctx context.Context, c *Client) (string, error) {
			var v struct {
				Values []struct {
					WorkspaceID string `json:"workspaceId"`
				} `json:"values"`
			}
			if err := c.siteGet(ctx, "/rest/servicedeskapi/assets/workspace", &v); err != nil {
				return "", err
			}
			if len(v.Values) == 0 {
				return "", nil
			}
			return v.Values[0].WorkspaceID, nil
		},
	},
}

type siteVarCache struct {
	mu     sync.Mutex
	values map[string]string
}

// fillSiteVars replaces every marker in rawURL with its value.
func (c *Client) fillSiteVars(ctx context.Context, rawURL string) (string, error) {
	for _, v := range siteVars {
		if !strings.Contains(rawURL, v.marker) {
			continue
		}
		value, err := c.siteVar(ctx, v)
		if err != nil {
			return "", err
		}
		rawURL = strings.ReplaceAll(rawURL, v.marker, value)
	}
	return rawURL, nil
}

func (c *Client) siteVar(ctx context.Context, v siteVar) (string, error) {
	if value := strings.TrimSpace(os.Getenv(v.marker)); value != "" {
		return value, nil
	}
	if c.Config != nil {
		if value := strings.TrimSpace(v.config(c)); value != "" {
			return value, nil
		}
	}
	c.siteVarCache.mu.Lock()
	defer c.siteVarCache.mu.Unlock()
	if value, ok := c.siteVarCache.values[v.marker]; ok {
		return value, nil
	}
	value, err := v.lookup(ctx, c)
	if err == nil && value == "" {
		err = fmt.Errorf("the site returned none")
	}
	if err != nil {
		return "", fmt.Errorf("finding %s from %s: %w\nhint: set it with: export %s=<value>", v.marker, c.BaseURL, err, v.marker)
	}
	if c.siteVarCache.values == nil {
		c.siteVarCache.values = map[string]string{}
	}
	c.siteVarCache.values[v.marker] = value
	return value, nil
}

func (c *Client) siteGet(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	auth, err := c.authHeader(ctx)
	if err != nil {
		return err
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s returned HTTP %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
