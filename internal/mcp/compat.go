// Local customization: keep MCP tool names working across Atlassian spec changes.
// See .printing-press-patches.json (command-compat).

package mcp

import (
	"context"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"keen/internal/compat"
)

// RegisterCompatTools adds the old names after RegisterTools: renamed tools
// answer under their old name too, retired and removed ones say what replaced them.
func RegisterCompatTools(s *server.MCPServer) {
	for oldGroup, newGroup := range compat.Renamed {
		for name, t := range s.ListTools() {
			if !strings.HasPrefix(name, newGroup+"_") {
				continue
			}
			oldName := oldGroup + strings.TrimPrefix(name, newGroup)
			if s.GetTool(oldName) != nil {
				continue // the old name is taken by a live tool
			}
			tool := t.Tool
			tool.Name = oldName
			tool.Description = "Old name of " + name + ". " + tool.Description
			s.AddTool(tool, t.Handler)
		}
	}
	for name, key := range compat.DeprecatedTools {
		r, ok := compat.RemovedEndpoints[key]
		t := s.GetTool(name)
		if !ok || t == nil {
			continue
		}
		note := compat.RemovedNote(r, true)
		tool := t.Tool
		tool.Description = "Removed by Atlassian. " + note
		s.AddTool(tool, func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
			return mcplib.NewToolResultError(note), nil
		})
	}
	for _, r := range compat.RetiredCommands {
		name := strings.ReplaceAll(r.Path, " ", "_")
		if s.GetTool(name) != nil {
			continue
		}
		note := r.Note
		s.AddTool(
			mcplib.NewTool(name, mcplib.WithDescription("Retired by Atlassian: "+note)),
			func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
				return mcplib.NewToolResultError(note), nil
			},
		)
	}
}
