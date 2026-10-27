// Local customization: top-level shortcut aliases for common agent operations.
// See .printing-press-patches.json for rationale.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func addAliasCommands(rootCmd *cobra.Command, flags *rootFlags) {
	// keen find-issues → keen jira-cloud-platform-search and-reconsile-issues-using-jql-post
	// Wraps with default fields so --agent returns key/summary/status, not just id
	findCmd := newJiraCloudPlatformSearchAndReconsileIssuesUsingJqlPostCmd(flags)
	origRunE := findCmd.RunE
	findCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("fields") && flags.compact {
			_ = cmd.Flags().Set("fields", `["key","summary","status","priority","issuetype","project"]`)
		}
		return origRunE(cmd, args)
	}
	findCmd.Use = "find-issues"
	findCmd.Aliases = []string{"fi"}
	findCmd.Short = "Search Jira issues by JQL (includes key/summary/status by default)"
	findCmd.Example = "  keen find-issues --jql \"project = INFRA AND status = 'In Progress'\" --agent"
	rootCmd.AddCommand(findCmd)

	// keen transition → keen issue transitions do
	transitionCmd := newIssueTransitionsDoCmd(flags)
	transitionCmd.Use = "transition [issueKey]"
	transitionCmd.Short = "Transition a Jira issue"
	transitionCmd.Example = "  keen transition INFRA-46 --transition-name QA --agent"
	// The Jira API requires transition.id (a numeric string); passing a name
	// fails with "'transition' identifier must be an integer". Resolve
	// --transition-name to its id here so agents can use human-readable names.
	transitionOrigRunE := transitionCmd.RunE
	transitionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("transition-name")
		id, _ := cmd.Flags().GetString("transition-id")
		if name != "" && id == "" && len(args) > 0 {
			resolved, err := resolveTransitionID(cmd.Context(), flags, args[0], name)
			if err != nil {
				return err
			}
			_ = cmd.Flags().Set("transition-id", resolved)
			// Clear the name so the body carries only transition.id.
			_ = cmd.Flags().Set("transition-name", "")
		}
		return transitionOrigRunE(cmd, args)
	}
	rootCmd.AddCommand(transitionCmd)

	// keen sprint → keen agile get-all-sprints
	sprintListCmd := newAgileGetAllSprintsCmd(flags)
	sprintListCmd.Use = "sprints [boardId]"
	sprintListCmd.Short = "List sprints for a board"
	sprintListCmd.Example = "  keen sprints 222 --agent"
	rootCmd.AddCommand(sprintListCmd)

	// keen get → keen issue get
	getIssueCmd := newIssueGetCmd(flags)
	getIssueCmd.Use = "get [issueKey]"
	getIssueCmd.Short = "Get a Jira issue"
	getIssueCmd.Example = "  keen get INFRA-46 --agent"
	rootCmd.AddCommand(getIssueCmd)
}

// resolveTransitionID looks up the numeric transition id for a transition name
// on a given issue, matching case-insensitively. Returns a helpful error
// listing the available transitions when the name does not match.
func resolveTransitionID(ctx context.Context, flags *rootFlags, issueKey, name string) (string, error) {
	c, err := flags.newClient()
	if err != nil {
		return "", err
	}
	data, err := c.Get(ctx, "/rest/api/3/issue/"+issueKey+"/transitions", nil)
	if err != nil {
		return "", classifyAPIError(err, flags)
	}
	var resp struct {
		Transitions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"transitions"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("parsing transitions for %s: %w", issueKey, err)
	}
	names := make([]string, 0, len(resp.Transitions))
	for _, t := range resp.Transitions {
		if strings.EqualFold(t.Name, name) {
			return t.ID, nil
		}
		names = append(names, t.Name)
	}
	return "", fmt.Errorf("no transition named %q available for %s (available: %s)", name, issueKey, strings.Join(names, ", "))
}

// rejectStrayArgs makes promoted endpoint commands that take no positional
// arguments refuse them. Without it `keen projects create` silently ran the
// projects/fields read and ignored "create". When a sibling group has that
// subcommand (`keen project create`), the error names it.
func rejectStrayArgs(rootCmd *cobra.Command) {
	for _, cmd := range rootCmd.Commands() {
		// Any word after the name ("get [issueKey]", "task <taskId>") declares an argument.
		if cmd.Annotations["pp:endpoint"] == "" || cmd.Args != nil || cmd.HasSubCommands() || len(strings.Fields(cmd.Use)) > 1 {
			continue
		}
		// `keen which` lists these as "<name> <operation>" (myself get-current-user),
		// so the operation's own name is accepted and ignored.
		_, op, _ := strings.Cut(cmd.Annotations["pp:endpoint"], ".")
		cmd.Args = func(c *cobra.Command, args []string) error {
			if len(args) == 0 || (len(args) == 1 && args[0] == op) {
				return nil
			}
			msg := fmt.Sprintf("unknown command %q for %q", args[0], c.CommandPath())
			for _, group := range []string{strings.TrimSuffix(c.Name(), "s"), c.Name() + "s"} {
				if group == c.Name() {
					continue
				}
				if g, _, err := rootCmd.Find([]string{group, args[0]}); err == nil && g.Name() == args[0] {
					msg += fmt.Sprintf("; did you mean %q?", g.CommandPath())
					break
				}
			}
			return fmt.Errorf("%s", msg)
		}
	}
}
