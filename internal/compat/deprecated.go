package compat

import "regexp"

// Atlassian's specs still list some operations it has already switched off:
// they answer 410 Gone. keen keeps their command and tool names, hides them, and
// says what to use instead. Deprecated operations that still work are listed in
// deprecated_gen.go and only labelled.

// Replacement names what to use instead of a removed operation.
type Replacement struct{ Command, Tool string }

// RemovedEndpoints maps a removed "METHOD path" (as in DeprecatedEndpoints) to
// its replacement. Add one when Atlassian turns off an operation it still lists.
var RemovedEndpoints = map[string]Replacement{
	"GET /rest/api/3/search": {`find-issues --jql "<JQL>"`, "jira-cloud-platform-search_and-reconsile-issues-using-jql"},
	"POST /rest/api/3/search": {`jira-cloud-platform-search and-reconsile-issues-using-jql-post`,
		"jira-cloud-platform-search_and-reconsile-issues-using-jql-post"},
}

var pathParam = regexp.MustCompile(`\{[^}]*\}`)

// EndpointKey is the "METHOD path" form DeprecatedEndpoints uses.
func EndpointKey(method, path string) string {
	return method + " " + pathParam.ReplaceAllString(path, "{}")
}

// RemovedNote is the message for a removed operation, for the CLI or a tool.
func RemovedNote(r Replacement, forTool bool) string {
	if forTool {
		return "Atlassian has switched this API off (it answers 410 Gone). Use the " + r.Tool + " tool instead."
	}
	return "Atlassian has switched this API off (it answers 410 Gone). Use: keen " + r.Command
}
