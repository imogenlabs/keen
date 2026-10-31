package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

// fakeSite counts requests by method and serves a value that changes after
// each write, so a stale read is visible.
func fakeSite(t *testing.T) (gets, writes *int32) {
	t.Helper()
	gets, writes = new(int32), new(int32)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			atomic.AddInt32(gets, 1)
			w.Header().Set("Content-Type", "application/json")
			body := `{"summary":"before"}`
			if atomic.LoadInt32(writes) > 0 {
				body = `{"summary":"after"}`
			}
			_, _ = w.Write([]byte(body))
			return
		}
		atomic.AddInt32(writes, 1)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("JIRA_CLOUD_PLATFORM_URL", srv.URL)
	t.Setenv("JIRA_CLOUD_PLATFORM_USERNAME", "agent@example.com")
	t.Setenv("JIRA_CLOUD_PLATFORM_PASSWORD", "test-token")
	return gets, writes
}

func callHandler(t *testing.T, method string) string {
	t.Helper()
	h := makeAPIHandler(method, "/rest/api/3/issue/PROJ-1", method == "GET", false, nil, nil, nil)
	res, err := h(context.Background(), mcplib.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("%s returned an error: %v", method, res.Content)
	}
	return res.Content[0].(mcplib.TextContent).Text
}

func TestMCPReadsAreCachedAndWritesClearTheCache(t *testing.T) {
	gets, _ := fakeSite(t)
	t.Setenv("KEEN_MCP_CACHE", "")

	callHandler(t, "GET")
	callHandler(t, "GET")
	if n := atomic.LoadInt32(gets); n != 1 {
		t.Fatalf("two identical reads sent %d requests to Atlassian, want 1", n)
	}

	callHandler(t, "PUT")
	if got := callHandler(t, "GET"); got != `{"summary":"after"}` {
		t.Fatalf("read after a write = %s, want the new value", got)
	}
	if n := atomic.LoadInt32(gets); n != 2 {
		t.Fatalf("read after a write was served from cache (%d requests), want a fresh one", n)
	}
}

func TestMCPCacheOff(t *testing.T) {
	gets, _ := fakeSite(t)
	t.Setenv("KEEN_MCP_CACHE", "off")

	callHandler(t, "GET")
	callHandler(t, "GET")
	if n := atomic.LoadInt32(gets); n != 2 {
		t.Fatalf("KEEN_MCP_CACHE=off sent %d requests for two reads, want 2", n)
	}
}
