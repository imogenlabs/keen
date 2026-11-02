// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Local customization; see .printing-press-patches.json (project-delete).
// The generator made "project delete" a group for its one subcommand
// (project-asynchronously), which left DELETE /rest/api/3/project/{key}
// without a command although keen-mcp has project_delete. The group now runs
// that delete when given a project, and still holds the subcommand.
func withProjectDelete(cmd *cobra.Command, flags *rootFlags) *cobra.Command {
	var enableUndo bool
	cmd.Use = "delete [projectIdOrKey]"
	cmd.Short = "Deletes a project (or: delete project-asynchronously <key>)"
	cmd.Long = "Deletes a project. You can't delete an archived project; restore it first. " +
		"Needs the Administer Jira global permission. With --enable-undo the project goes to the Jira recycle bin."
	cmd.Example = "  keen project delete PROJ --dry-run\n  keen project delete PROJ --enable-undo\n  keen project delete project-asynchronously PROJ"
	cmd.Annotations = map[string]string{"pp:method": "DELETE", "pp:path": "/rest/api/3/project/{projectIdOrKey}"}
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().BoolVar(&enableUndo, "enable-undo", false, "Put the project in the Jira recycle bin, where it can be restored")
	listSubs := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listSubs(cmd, args)
		}
		c, err := flags.newClient()
		if err != nil {
			return err
		}
		path := replacePathParam("/rest/api/3/project/{projectIdOrKey}", "projectIdOrKey", args[0])
		params := map[string]string{}
		if cmd.Flags().Changed("enable-undo") {
			params["enableUndo"] = fmt.Sprintf("%v", enableUndo)
		}
		_, statusCode, err := c.DeleteWithParams(cmd.Context(), path, params)
		if err != nil {
			return classifyDeleteError(err, flags)
		}
		if flags.dryRun {
			return nil
		}
		if flags.asJSON || !isTerminal(cmd.OutOrStdout()) {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{
				"action": "delete", "resource": "project", "path": path, "status": statusCode, "success": true,
			})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted project %s\n", args[0])
		return nil
	}
	return cmd
}
