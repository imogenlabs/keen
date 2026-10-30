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
	name     string
	words    map[string]bool // from the name
	text     string          // lower-case name and description
	summary  map[string]any
	readOnly bool
}

type toolIndex []indexedTool

func buildToolIndex(all *server.MCPServer) toolIndex {
	var idx toolIndex
	for name, t := range all.ListTools() {
		desc := t.Tool.Description
		readOnly := t.Tool.Annotations.ReadOnlyHint != nil && *t.Tool.Annotations.ReadOnlyHint
		short := desc
		if len(short) > 300 {
			short = short[:300] + "..."
		}
		words := map[string]bool{}
		for _, w := range splitWords(name) {
			if !nameNoise[w] {
				words[w] = true
			}
		}
		idx = append(idx, indexedTool{
			name:     name,
			words:    words,
			text:     strings.ToLower(name + " " + desc),
			readOnly: readOnly,
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

// search ranks tools by how many query words appear in their name (weighted)
// and description, minus the name words the query did not ask for, so the
// plain tool beats a longer one that merely shares a word. An exact tool name
// comes first.
func (idx toolIndex) search(query string, limit int) []map[string]any {
	q := strings.TrimSpace(query)
	words := splitWords(q)
	type scored struct {
		t     indexedTool
		score int
	}
	var hits []scored
	for _, t := range idx {
		score := 0
		if strings.EqualFold(t.name, q) {
			score += 1000
		}
		matched := map[string]bool{}
		for _, w := range words {
			switch {
			case t.words[w]:
				score += 5
				matched[w] = true
			case strings.Contains(t.text, w):
				score++
			}
			for _, alt := range actionWords[w] {
				if t.words[alt] {
					score += 3
					matched[alt] = true
				}
			}
		}
		if len(matched) > 0 {
			// Name words the query did not ask for cost a little; a tool whose
			// name matched at all stays in the results.
			score = max(score-2*(len(t.words)-len(matched)), 1)
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

var stopWords = map[string]bool{"a": true, "an": true, "the": true, "to": true, "of": true, "in": true,
	"for": true, "by": true, "on": true, "my": true, "and": true, "or": true, "with": true,
	"is": true, "are": true, "who": true, "what": true, "which": true, "how": true, "me": true, "i": true,
	"do": true, "does": true, "can": true, "show": true}

// Words the generator puts in tool names that say nothing about the action.
var nameNoise = map[string]bool{"jira": true, "cloud": true, "platform": true, "using": true, "api": true}

// Words people use for an action, mapped to the words keen's tool names use.
var actionWords = map[string][]string{
	"list":   {"get", "all", "list"},
	"show":   {"get"},
	"find":   {"search", "get"},
	"move":   {"transition"},
	"change": {"update", "edit"},
	"remove": {"delete"},
	"make":   {"create"},
	"new":    {"create"},
	"done":   {"transition"},
	"close":  {"transition"},
	"status": {"transition"},
}

// splitWords lower-cases s, splits it on anything but letters and digits, drops
// filler words and a plural "s", so "list sprints" matches "sprint".
func splitWords(s string) []string {
	var out []string
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	}) {
		if stopWords[w] {
			continue
		}
		if len(w) > 3 && strings.HasSuffix(w, "s") && !strings.HasSuffix(w, "ss") {
			w = strings.TrimSuffix(w, "s")
		}
		out = append(out, w)
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
