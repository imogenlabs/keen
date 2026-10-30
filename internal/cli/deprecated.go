// Local customization: label commands for operations Atlassian deprecated, and
// point the ones it has switched off at their replacement.
// See .printing-press-patches.json (deprecated-endpoints).

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"keen/internal/compat"
)

func markDeprecatedCommands(cmd *cobra.Command) {
	for _, c := range cmd.Commands() {
		markDeprecatedCommands(c)
	}
	method, path := cmd.Annotations["pp:method"], cmd.Annotations["pp:path"]
	if method == "" || path == "" {
		return
	}
	key := compat.EndpointKey(method, path)
	if r, ok := compat.RemovedEndpoints[key]; ok {
		note := compat.RemovedNote(r, false)
		cmd.Short = "Removed by Atlassian. " + note
		cmd.Hidden = true
		cmd.RunE = func(c *cobra.Command, args []string) error {
			return fmt.Errorf("keen %s: %s", c.CommandPath()[len(c.Root().Name())+1:], note)
		}
		return
	}
	if _, ok := compat.DeprecatedEndpoints[key]; ok {
		cmd.Short = "Deprecated by Atlassian. " + cmd.Short
	}
}

// withoutRemoved drops sync resources whose endpoint Atlassian has switched off.
func withoutRemoved(resources []string) []string {
	var out []string
	for _, r := range resources {
		path, err := syncResourcePath(r)
		if err != nil {
			out = append(out, r)
			continue
		}
		if _, gone := compat.RemovedEndpoints[compat.EndpointKey("GET", path)]; !gone {
			out = append(out, r)
		}
	}
	return out
}
