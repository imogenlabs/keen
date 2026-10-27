package cli

import (
	"strings"
	"testing"
)

func TestPromotedCommandRejectsStrayArgs(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	root.SetArgs([]string{"projects", "create", "--dry-run"})
	root.SilenceUsage, root.SilenceErrors = true, true
	err := root.Execute()
	if err == nil {
		t.Fatal("keen projects create ran instead of failing")
	}
	if !strings.Contains(err.Error(), `did you mean "keen project create"`) {
		t.Fatalf("error does not point at keen project create: %v", err)
	}
	if !isCobraUsageError(err) {
		t.Fatalf("stray-arg error should be a usage error (exit 2): %v", err)
	}
}

func TestCommandsThatDeclareArgsStillTakeThem(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	checked := 0
	for _, c := range root.Commands() {
		if len(strings.Fields(c.Use)) < 2 {
			continue
		}
		checked++
		if c.Args != nil {
			if err := c.Args(c, []string{"PROJ-1"}); err != nil {
				t.Errorf("%q declares an argument but rejects one: %v", c.Use, err)
			}
		}
	}
	for _, name := range []string{"get", "sprints", "transition", "task"} {
		c, _, err := root.Find([]string{name, "X-1"})
		if err != nil || c.Name() != name {
			t.Errorf("keen %s X-1 does not resolve to %s: %v", name, name, err)
		}
	}
	if checked < 4 {
		t.Fatalf("only %d top-level commands declare arguments; the check is not seeing them", checked)
	}
}
