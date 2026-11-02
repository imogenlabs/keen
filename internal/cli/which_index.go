// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Local customization; see .printing-press-patches.json (shared-ranking).
// The generated whichIndex lists only API operations, so keen's own commands
// (find-issues, transition, stale, sprints, ...) could never be suggested.
// whichIndexFor adds every visible command missing from it, drops commands
// Atlassian has switched off (hidden by markDeprecatedCommands), and carries
// the "Deprecated by Atlassian." label so the ranker can rank those lower.
func whichIndexFor(root *cobra.Command) []whichEntry {
	cmds := map[string]*cobra.Command{}
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			walk(sub)
		}
		if c != root {
			cmds[strings.TrimPrefix(c.CommandPath(), root.Name()+" ")] = c
		}
	}
	walk(root)

	seen := map[string]bool{}
	out := make([]whichEntry, 0, len(whichIndex)+64)
	for _, e := range whichIndex {
		c := cmds[e.Command]
		if c != nil && c.Hidden {
			continue
		}
		if c != nil && strings.HasPrefix(c.Short, "Deprecated by Atlassian") {
			e.Description = "Deprecated by Atlassian. " + e.Description
		}
		seen[e.Command] = true
		out = append(out, e)
	}
	paths := make([]string, 0, len(cmds))
	for path := range cmds {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		c := cmds[path]
		if seen[path] || c.Hidden || !c.Runnable() || c.HasSubCommands() || c.Name() == "help" || strings.HasPrefix(path, "completion") {
			continue
		}
		out = append(out, whichEntry{Command: path, Description: c.Short, Group: whichNativeGroup})
	}
	return out
}

// whichNativeGroup marks keen's own commands; their summary counts as keywords.
const whichNativeGroup = "keen"

var (
	mdLink   = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	mdMarks  = regexp.MustCompile("[*`#_]+")
	htmlTag  = regexp.MustCompile(`<[^>]+>`)
	sentence = regexp.MustCompile(`^(.+?[.!?])(\s|$)`)
)

// whichOneLine turns an API description (often several paragraphs of
// markdown) into its first sentence, plain text, at most width characters.
func whichOneLine(desc string, width int) string {
	s := strings.TrimSpace(desc)
	if i := strings.Index(s, "\n\n"); i >= 0 {
		s = s[:i]
	}
	s = mdLink.ReplaceAllString(s, "$1")
	s = htmlTag.ReplaceAllString(s, "")
	s = mdMarks.ReplaceAllString(s, "")
	s = strings.Join(strings.Fields(s), " ")
	if m := sentence.FindStringSubmatch(s); m != nil {
		s = m[1]
	}
	if width < 20 {
		width = 20
	}
	if r := []rune(s); len(r) > width {
		s = strings.TrimRight(string(r[:width-1]), " ") + "…"
	}
	return s
}

// terminalWidth is the width of the terminal w writes to, from $COLUMNS, or
// 100 when it is not known.
func terminalWidth(w io.Writer) int {
	if f, ok := w.(*os.File); ok && f == os.Stdout {
		if n, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && n > 0 {
			return n
		}
	}
	return 100
}
