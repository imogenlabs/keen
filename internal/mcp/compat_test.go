package mcp

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"

	"keen/internal/compat"
)

func TestCompatToolsKeepOldNames(t *testing.T) {
	s := server.NewMCPServer("keen", "test", server.WithToolCapabilities(false))
	RegisterTools(s)
	RegisterCompatTools(s)
	for _, r := range compat.RetiredCommands {
		name := strings.ReplaceAll(r.Path, " ", "_")
		if s.GetTool(name) == nil {
			t.Errorf("retired tool %s is missing", name)
		}
	}
	for _, old := range []string{"policies_create-policy", "policies_get", "policies_resources_add-to-policy", "policies_validate_policy"} {
		if s.GetTool(old) == nil {
			t.Errorf("renamed tool %s lost its old name", old)
		}
	}
	if s.GetTool("policies_list") == nil {
		t.Error("new Jira tool policies_list is missing")
	}
}
