// Package compat lists keen command names that changed when Atlassian changed
// its API specs, so both the CLI and the MCP server keep the old names working.
// Local customization; see .printing-press-patches.json (command-compat).
package compat

// Retired is a command whose Atlassian API is gone. Path is the CLI path;
// the MCP tool name is the path joined with "_".
type Retired struct{ Path, Note string }

const (
	appAccess   = "Atlassian removed the app access settings API from its specification in 2026; it has no replacement."
	june30      = "Atlassian retired this API on 30 June 2026. Use: keen "
	transitions = "Atlassian removed this API on 1 June 2026. "
)

var RetiredCommands = []Retired{
	{"app-access-settings create-domain", appAccess},
	{"app-access-settings get-domain", appAccess},
	{"app-access-settings list-domains", appAccess},
	{"app-access-settings replace-domain", appAccess},
	{"app-access-settings upsert-product", appAccess},
	{"atlassian-access-groups search", june30 + "directories groups get"},
	{"atlassian-access-users search", june30 + "directories users get-directory"},
	{"atlassian-access-users create-orgs", "Atlassian retired this API in 2026. Use: keen atlassian-access-users create"},
	{"directory assign-role-to-group", june30 + "directories groups create-orgs-2"},
	{"directory create", june30 + "directories groups create"},
	{"directory create-orgs", june30 + "directories groups create-orgs"},
	{"directory create-orgs-2", june30 + "directories users create"},
	{"directory create-orgs-3", june30 + "directories users create-orgs"},
	{"directory delete", june30 + "directories groups delete"},
	{"directory delete-orgs", june30 + "directories users delete"},
	{"directory delete-orgs-2", june30 + "directories groups delete-orgs"},
	{"directory revoke-role-to-group", june30 + "directories groups create-orgs-3"},
	{"jira-cloud-platform-workflow create-transition-property", transitions + "Set transition properties with: keen workflows update"},
	{"jira-cloud-platform-workflow delete-transition-property", transitions + "Set transition properties with: keen workflows update"},
	{"jira-cloud-platform-workflow get-transition-properties", transitions + "Read transition properties with: keen workflows search"},
	{"jira-cloud-platform-workflow update-transition-property", transitions + "Set transition properties with: keen workflows update"},
}

// Renamed maps an old command group to the group that now holds its
// subcommands. Org policies moved when Jira added its own policies API.
var Renamed = map[string]string{
	"policies": "atlassian-access-policies",
}
