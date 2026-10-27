// Local customization: keep command names working across Atlassian spec changes.
// See .printing-press-patches.json (command-compat).

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"keen/internal/compat"
)

// renamedGroups builds a fresh copy of the group that now holds a renamed
// group's subcommands (see compat.Renamed).
var renamedGroups = map[string]func(*rootFlags) *cobra.Command{
	"atlassian-access-policies": newAtlassianAccessPoliciesCmd,
}

func childNamed(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func addCompatCommands(rootCmd *cobra.Command, flags *rootFlags) {
	// The org policy commands moved from `keen policies ...` to
	// `keen atlassian-access-policies ...` when Jira gained its own policies API.
	// The old names keep working, hidden, alongside the new Jira ones.
	for oldGroup, newGroup := range compat.Renamed {
		build, ok := renamedGroups[newGroup]
		old := childNamed(rootCmd, oldGroup)
		if !ok || old == nil {
			continue
		}
		moved := build(flags)
		for _, c := range moved.Commands() {
			if childNamed(old, c.Name()) == nil {
				moved.RemoveCommand(c)
				c.Hidden = true
				old.AddCommand(c)
			}
		}
	}

	for _, r := range compat.RetiredCommands {
		words := strings.Fields(r.Path)
		parent := rootCmd
		for _, w := range words[:len(words)-1] {
			next := childNamed(parent, w)
			if next == nil {
				next = &cobra.Command{Use: w, Hidden: true, RunE: parentNoSubcommandRunE(flags)}
				parent.AddCommand(next)
			}
			parent = next
		}
		leaf := words[len(words)-1]
		if childNamed(parent, leaf) != nil {
			continue // never shadow a live command
		}
		note, path := r.Note, r.Path
		parent.AddCommand(&cobra.Command{
			Use:    leaf,
			Short:  "Retired by Atlassian: " + note,
			Hidden: true,
			Args:   cobra.ArbitraryArgs,
			// Old scripts pass the old flags; accept them so the message is what they see.
			FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("keen %s: %s", path, note)
			},
		})
	}
}
