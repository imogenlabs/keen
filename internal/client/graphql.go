// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// GraphQLPath is Atlassian's GraphQL gateway on the configured site.
const GraphQLPath = "/gateway/api/graphql"

// ErrGraphQL marks errors the gateway reported in the response body. GraphQL
// answers a failed query with HTTP 200, so the status alone would call it a
// success.
var ErrGraphQL = errors.New("graphql error")

// GraphQL runs one named query on the site's gateway and returns its data.
// Errors in the response fail the call even when the status is 200; the
// partial data that came with them is dropped.
func (c *Client) GraphQL(ctx context.Context, operation, query string, variables map[string]any, headers map[string]string) (json.RawMessage, error) {
	body := map[string]any{"operationName": operation, "query": query, "variables": variables}
	raw, _, err := c.PostWithHeaders(ctx, GraphQLPath, body, headers)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("reading the GraphQL response: %w", err)
	}
	if len(resp.Errors) > 0 {
		msgs := make([]string, 0, len(resp.Errors))
		for _, e := range resp.Errors {
			msgs = append(msgs, e.Message)
		}
		return nil, fmt.Errorf("%w: %s", ErrGraphQL, strings.Join(msgs, "; "))
	}
	return resp.Data, nil
}

// SiteCloudID returns the site's cloud id: KEEN_CLOUD_ID, cloud_id in
// auth.json, or a lookup on the site (cached).
func (c *Client) SiteCloudID(ctx context.Context) (string, error) {
	return c.siteVar(ctx, siteVars[0])
}
