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

func TestPromotedCommandWithPositionalArgUnchanged(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	for _, c := range root.Commands() {
		if c.Name() == "task" && c.Args != nil {
			t.Fatal("task <taskId> must keep accepting its positional argument")
		}
	}
}
