// Local customization: `keen release sync` reads release tags and the issue keys
// named in their commits from a git checkout. See .printing-press-patches.json.

package cli

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// defaultReleaseTagPattern matches plain release tags (v1.2.3 or 1.2.3) and
// leaves out pre-releases such as v1.2.3-rc1, -dev.4 or -beta.
const defaultReleaseTagPattern = `^v?\d+\.\d+\.\d+$`

// releaseEntry is one Jira version's worth of work.
type releaseEntry struct {
	Tag  string   `json:"tag"`
	Date string   `json:"date,omitempty"`
	Keys []string `json:"keys"`
}

type gitRepo struct{ dir string }

func (g gitRepo) run(args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", g.dir}, args...)...).Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(ee.Stderr)))
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
}

// issueKeyPattern finds PROJECT-123 keys that are not part of a longer token
// (XPROJ-1, SUB-PROJ-1 and PROJ-1a are not PROJ keys).
func issueKeyPattern(project string) *regexp.Regexp {
	return regexp.MustCompile(`(?:^|[^A-Za-z0-9-])` + regexp.QuoteMeta(project) + `-(\d+)(?:[^0-9A-Za-z]|$)`)
}

// extractIssueKeys returns the distinct keys of one project named in text,
// normalised (PROJ-0042 becomes PROJ-42).
func extractIssueKeys(project, text string) map[string]bool {
	keys := map[string]bool{}
	re := issueKeyPattern(project)
	// The trailing context of one match can be the leading context of the
	// next ("PROJ-1 PROJ-2"), so scan line by line with overlapping retries.
	for i := 0; i < len(text); {
		loc := re.FindStringSubmatchIndex(text[i:])
		if loc == nil {
			break
		}
		if n, err := strconv.Atoi(text[i+loc[2] : i+loc[3]]); err == nil && n > 0 {
			keys[fmt.Sprintf("%s-%d", project, n)] = true
		}
		i += loc[3]
	}
	return keys
}

func versionParts(tag string) []int {
	parts := strings.Split(strings.TrimPrefix(tag, "v"), ".")
	out := make([]int, len(parts))
	for i, p := range parts {
		out[i], _ = strconv.Atoi(p)
	}
	return out
}

func sortVersionTags(tags []string) {
	sort.SliceStable(tags, func(i, j int) bool {
		a, b := versionParts(tags[i]), versionParts(tags[j])
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
}

func (g gitRepo) releaseTags(pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("--tag-pattern: %w", err)
	}
	out, err := g.run("tag", "-l")
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, t := range strings.Fields(out) {
		if re.MatchString(t) {
			tags = append(tags, t)
		}
	}
	sortVersionTags(tags)
	return tags, nil
}

func (g gitRepo) keysIn(project, revRange string) (map[string]bool, error) {
	out, err := g.run("log", "--format=%B", revRange)
	if err != nil {
		return nil, err
	}
	return extractIssueKeys(project, out), nil
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, _ := strconv.Atoi(keys[i][strings.LastIndex(keys[i], "-")+1:])
		b, _ := strconv.Atoi(keys[j][strings.LastIndex(keys[j], "-")+1:])
		return a < b
	})
	return keys
}

// newKeys is the set of keys first shipped by `to`: named in from..to and in
// nothing reachable from `from`. With from == "" it is everything up to `to`.
func (g gitRepo) newKeys(project, from, to string) ([]string, error) {
	if from == "" {
		set, err := g.keysIn(project, to)
		return sortedKeys(set), err
	}
	in, err := g.keysIn(project, from+".."+to)
	if err != nil {
		return nil, err
	}
	before, err := g.keysIn(project, from)
	if err != nil {
		return nil, err
	}
	for k := range before {
		delete(in, k)
	}
	return sortedKeys(in), nil
}

// releaseMapping returns one entry per tag (the first release containing each
// key), in version order.
func (g gitRepo) releaseMapping(project string, tags []string) ([]releaseEntry, error) {
	out := make([]releaseEntry, 0, len(tags))
	prev := ""
	for _, tag := range tags {
		keys, err := g.newKeys(project, prev, tag)
		if err != nil {
			return nil, err
		}
		date, err := g.run("log", "-1", "--format=%cs", tag)
		if err != nil {
			return nil, err
		}
		out = append(out, releaseEntry{Tag: tag, Date: strings.TrimSpace(date), Keys: keys})
		prev = tag
	}
	return out, nil
}

var githubRemoteRe = regexp.MustCompile(`github\.com[:/]([^/]+)/([^/]+?)(?:\.git)?/?$`)

// githubRepoFromRemote turns git@github.com:o/r.git or https://github.com/o/r
// into "o/r".
func githubRepoFromRemote(remote string) string {
	m := githubRemoteRe.FindStringSubmatch(strings.TrimSpace(remote))
	if m == nil {
		return ""
	}
	return m[1] + "/" + m[2]
}
