// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"keen/internal/config"
)

func TestGraphQLErrorsWithStatus200Fail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GraphQLPath || r.Header.Get("X-Query-Context") != "ctx" {
			http.Error(w, "wrong request", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"data":null,"errors":[{"message":"field removed"}]}`))
	}))
	defer srv.Close()
	c := New(&config.Config{BaseURL: srv.URL, JiraCloudPlatformUsername: "u", JiraCloudPlatformPassword: "p"}, 0, 0)
	c.NoCache = true
	_, err := c.GraphQL(context.Background(), "Q", "query Q { x }", nil, map[string]string{"X-Query-Context": "ctx"})
	if !errors.Is(err, ErrGraphQL) {
		t.Fatalf("err = %v; want ErrGraphQL for errors in a 200 response", err)
	}
}
