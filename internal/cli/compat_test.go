package cli

import (
	"strings"
	"testing"

	"keen/internal/compat"
)

func TestRetiredCommandsExplainTheReplacement(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	for _, r := range compat.RetiredCommands {
		c, _, err := root.Find(strings.Fields(r.Path))
		if err != nil || c.CommandPath() != "keen "+r.Path {
			t.Errorf("keen %s does not resolve: %v", r.Path, err)
			continue
		}
		if c.Hidden && c.RunE != nil {
			err := c.RunE(c, nil)
			if err == nil || !strings.Contains(err.Error(), "Atlassian") {
				t.Errorf("keen %s should explain the retirement, got %v", r.Path, err)
			}
		}
	}
}

func TestOldPolicyCommandNamesStillWork(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	for _, path := range []string{
		"policies create-policy", "policies delete-policy", "policies get", "policies get-policy-by-id",
		"policies update-policy", "policies resources add-to-policy", "policies resources delete-policy",
		"policies resources update-policy", "policies validate policy",
	} {
		c, _, err := root.Find(strings.Fields(path))
		if err != nil || c.CommandPath() != "keen "+path {
			t.Errorf("keen %s no longer resolves: %v", path, err)
			continue
		}
		if c.Annotations["pp:endpoint"] == "" {
			t.Errorf("keen %s resolves to %q, which is not an API command", path, c.CommandPath())
		}
	}
	// The new Jira policies commands are still there.
	if c, _, err := root.Find([]string{"policies", "list"}); err != nil || c.Name() != "list" {
		t.Errorf("keen policies list missing: %v", err)
	}
}
