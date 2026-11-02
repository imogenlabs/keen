// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package mcp

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"keen/internal/compat"
	"keen/internal/rank"
)

// keen has about 1,900 tools. Listed one by one, their names, descriptions and
// schemas come to about 2 MB, more than most models can hold, so an AI app
// that connects to keen-mcp would fill its memory before doing any work. By
// default keen-mcp offers two tools instead: keen_search finds the tools for a
// task, keen_run calls one by name through the same handler. KEEN_MCP_TOOLS=all
// lists every tool as before.

// ToolsModeEnv selects the tool list keen-mcp offers: "all" for every tool,
// anything else for keen_search and keen_run.
const ToolsModeEnv = "KEEN_MCP_TOOLS"

// Register adds keen's tools to s in the mode ToolsModeEnv selects.
func Register(s *server.MCPServer) {
	if strings.EqualFold(strings.TrimSpace(os.Getenv(ToolsModeEnv)), "all") {
		RegisterTools(s)
		RegisterCompatTools(s)
		return
	}
	all := server.NewMCPServer("keen-all", "", server.WithToolCapabilities(false))
	RegisterTools(all)
	RegisterCompatTools(all)
	RegisterSearchRun(s, all)
}

// RegisterSearchRun adds keen_search and keen_run to s, backed by the tools in all.
func RegisterSearchRun(s, all *server.MCPServer) {
	index := buildToolIndex(all)
	s.AddTool(mcplib.NewTool("keen_search",
		mcplib.WithDescription("Find keen tools for a task. keen covers Atlassian Cloud: Jira (issues, search, "+
			"projects, workflows), Jira Software (boards, sprints), Service Management, Confluence (pages, spaces, "+
			"blog posts), Assets, Compass, Bitbucket, Marketplace, Atlassian Admin, and the org chart / Teamwork "+
			"Graph. Describe what you want in a few words (\"transition an issue\", \"list sprints on a board\", "+
			"\"confluence page by id\"); you get the best matching tools with their parameters. Then call keen_run."),
		mcplib.WithString("query", mcplib.Required(), mcplib.Description("what you want to do, in a few words, or an exact tool name")),
		mcplib.WithNumber("limit", mcplib.Description("how many tools to return (default 8, at most 25)")),
		mcplib.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		args := req.GetArguments()
		query, _ := args["query"].(string)
		limit := 8
		if n, ok := args["limit"].(float64); ok && n > 0 {
			limit = min(int(n), 25)
		}
		matches := index.search(query, limit)
		if len(matches) == 0 {
			return mcplib.NewToolResultError("no keen tool matches " + query + "; try other words, e.g. the product and the action"), nil
		}
		return mcplib.NewToolResultJSON(map[string]any{"tools": matches})
	})
	s.AddTool(mcplib.NewTool("keen_run",
		mcplib.WithDescription("Run a keen tool found with keen_search. Pass its exact name and its arguments as an "+
			"object. Tools marked readOnly only read; others can create, change or delete data in Atlassian."),
		mcplib.WithString("tool", mcplib.Required(), mcplib.Description("exact tool name from keen_search")),
		mcplib.WithObject("arguments", mcplib.Description("the tool's arguments, as keen_search listed them")),
	), func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		args := req.GetArguments()
		name, _ := args["tool"].(string)
		tool := all.GetTool(name)
		if tool == nil {
			return mcplib.NewToolResultError("no keen tool named " + name + "; find it with keen_search"), nil
		}
		inner := mcplib.CallToolRequest{Header: req.Header}
		inner.Params.Name = name
		toolArgs, _ := args["arguments"].(map[string]any)
		if toolArgs == nil {
			toolArgs = map[string]any{}
		}
		inner.Params.Arguments = toolArgs
		return tool.Handler(ctx, inner)
	})
}

type indexedTool struct {
	name    string
	item    rank.Item
	summary map[string]any
	removed bool // Atlassian switched the API off; found only by exact name
}

type toolIndex []indexedTool

func buildToolIndex(all *server.MCPServer) toolIndex {
	var idx toolIndex
	for name, t := range all.ListTools() {
		desc := t.Tool.Description
		readOnly := t.Tool.Annotations.ReadOnlyHint != nil && *t.Tool.Annotations.ReadOnlyHint
		key, deprecated := compat.DeprecatedTools[name]
		_, removed := compat.RemovedEndpoints[key]
		if deprecated && !removed {
			desc = "Deprecated by Atlassian; prefer another tool if one fits. " + desc
		}
		short := desc
		if len(short) > 300 {
			short = short[:300] + "..."
		}
		idx = append(idx, indexedTool{
			name:    name,
			item:    rank.NewItem(name, desc, deprecated && !removed),
			removed: removed,
			summary: map[string]any{
				"name":        name,
				"description": short,
				"readOnly":    readOnly,
				"parameters":  json.RawMessage(mustJSON(t.Tool.InputSchema)),
			},
		})
	}
	sort.Slice(idx, func(i, j int) bool { return idx[i].name < idx[j].name })
	return idx
}

// search ranks tools with rank.Item.Score, the same ranking keen which uses.
// Removed tools show up only by exact name. An exact tool name comes first.
func (idx toolIndex) search(query string, limit int) []map[string]any {
	type scored struct {
		t     indexedTool
		score int
	}
	var hits []scored
	for _, t := range idx {
		score, _ := t.item.Score(query)
		if t.removed && score < rank.Exact {
			continue
		}
		if score > 0 {
			hits = append(hits, scored{t, score})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].score != hits[j].score {
			return hits[i].score > hits[j].score
		}
		return len(hits[i].t.name) < len(hits[j].t.name)
	})
	out := []map[string]any{}
	for i := 0; i < len(hits) && i < limit; i++ {
		out = append(out, hits[i].t.summary)
	}
	return out
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return b
}
