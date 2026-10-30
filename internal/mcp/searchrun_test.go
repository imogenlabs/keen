// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func fakeTools() *server.MCPServer {
	all := server.NewMCPServer("t", "")
	add := func(name, desc string) {
		all.AddTool(mcplib.NewTool(name, mcplib.WithDescription(desc), mcplib.WithString("issueIdOrKey")),
			func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
				b, _ := json.Marshal(map[string]any{"tool": name, "args": req.GetArguments()})
				return mcplib.NewToolResultText(string(b)), nil
			})
	}
	add("issue_get", "Returns the details for an issue by key.")
	add("issue_properties_get-issue-property-keys", "Returns the keys of all properties for an issue.")
	add("issue_transitions_do", "Performs an issue transition and, if the transition has a screen, updates the fields.")
	add("agile_move-issues-to-epic", "Moves issues to an epic.")
	add("agile_get-all-sprints", "Returns all sprints from a board.")
	add("agile_get-board-issues-for-sprint", "Get all issues on a board in a sprint.")
	add("teamwork_graph_user_manager", "Retrieves the manager for a specific user.")
	return all
}

func call(t *testing.T, s *server.MCPServer, tool string, args map[string]any) *mcplib.CallToolResult {
	t.Helper()
	req := mcplib.CallToolRequest{}
	req.Params.Name = tool
	req.Params.Arguments = args
	res, err := s.GetTool(tool).Handler(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func firstTool(t *testing.T, res *mcplib.CallToolResult) string {
	t.Helper()
	var out struct {
		Tools []struct{ Name string } `json:"tools"`
	}
	if err := json.Unmarshal([]byte(res.Content[0].(mcplib.TextContent).Text), &out); err != nil || len(out.Tools) == 0 {
		t.Fatalf("no tools in %v (%v)", res.Content, err)
	}
	return out.Tools[0].Name
}

func TestSearchRanksThePlainTool(t *testing.T) {
	s := server.NewMCPServer("t", "")
	RegisterSearchRun(s, fakeTools())
	for query, want := range map[string]string{
		"get issue by key":        "issue_get",
		"move issue to done":      "issue_transitions_do",
		"list sprints on a board": "agile_get-all-sprints",
		"who is my manager":       "teamwork_graph_user_manager",
		"issue_get":               "issue_get",
	} {
		if got := firstTool(t, call(t, s, "keen_search", map[string]any{"query": query})); got != want {
			t.Errorf("search %q: first = %s, want %s", query, got, want)
		}
	}
}

func TestRunCallsTheToolWithItsArguments(t *testing.T) {
	s := server.NewMCPServer("t", "")
	RegisterSearchRun(s, fakeTools())
	res := call(t, s, "keen_run", map[string]any{"tool": "issue_get", "arguments": map[string]any{"issueIdOrKey": "PROJ-1"}})
	if got := res.Content[0].(mcplib.TextContent).Text; got != `{"args":{"issueIdOrKey":"PROJ-1"},"tool":"issue_get"}` {
		t.Errorf("run result = %s", got)
	}
	if res := call(t, s, "keen_run", map[string]any{"tool": "nope"}); !res.IsError {
		t.Error("an unknown tool should be an error result")
	}
}

func TestRegisterModes(t *testing.T) {
	t.Setenv(ToolsModeEnv, "")
	s := server.NewMCPServer("t", "")
	Register(s)
	if n := len(s.ListTools()); n != 2 {
		t.Errorf("default mode lists %d tools, want 2 (keen_search, keen_run)", n)
	}
	t.Setenv(ToolsModeEnv, "all")
	s = server.NewMCPServer("t", "")
	Register(s)
	if n := len(s.ListTools()); n < 1000 {
		t.Errorf("KEEN_MCP_TOOLS=all lists %d tools, want every tool", n)
	}
}
