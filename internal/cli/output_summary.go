// Local customization: --output-summary and --select path checking.
// See .printing-press-patches.json (agent-output-summary).

package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// defaultSummaryLimit is the output size, in bytes, above which
// --output-summary writes the result to a file. About 1,000 tokens.
const defaultSummaryLimit = 4096

// outputSummary holds what one run of keen holds back: its output, the
// --output-summary limit, and the first --select path that matched nothing.
var outputSummary struct {
	limit     int
	buf       *bytes.Buffer
	selectErr error
}

// noteSelectPaths is called by filterFields. Generated commands apply
// --select in about a thousand places, so this is where a bad path is caught.
func noteSelectPaths(data json.RawMessage, fields string) {
	if outputSummary.selectErr == nil {
		outputSummary.selectErr = checkSelectPaths(data, fields)
	}
}

func addOutputSummaryFlag(rootCmd *cobra.Command, flags *rootFlags) {
	f := rootCmd.PersistentFlags()
	f.IntVar(&outputSummary.limit, "output-summary", 0,
		"Write output larger than N bytes to a temp file and print a short JSON summary instead (bare flag: "+
			strconv.Itoa(defaultSummaryLimit)+")")
	f.Lookup("output-summary").NoOptDefVal = strconv.Itoa(defaultSummaryLimit)

	pre := rootCmd.PersistentPreRunE
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if pre != nil {
			if err := pre(cmd, args); err != nil {
				return err
			}
		}
		// Hold output back for --output-summary, and for --select so a bad
		// path can be reported instead of printing an unfiltered response.
		if outputSummary.limit <= 0 && flags.selectFields == "" {
			return nil
		}
		outputSummary.buf = &bytes.Buffer{}
		if flags.deliverBuf != nil {
			cmd.SetOut(io.MultiWriter(outputSummary.buf, flags.deliverBuf))
		} else {
			cmd.SetOut(outputSummary.buf)
		}
		return nil
	}
}

// flushOutputSummary writes what --output-summary held back: the output
// itself when it is small, otherwise a file plus a summary line. It runs
// after the command, also when the command failed, so nothing is lost.
func flushOutputSummary(w io.Writer) error {
	buf := outputSummary.buf
	outputSummary.buf = nil
	if err := outputSummary.selectErr; err != nil {
		outputSummary.selectErr = nil
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	if buf == nil || buf.Len() == 0 {
		return nil
	}
	data := buf.Bytes()
	if outputSummary.limit <= 0 || len(data) <= outputSummary.limit {
		_, err := w.Write(data)
		return err
	}

	isJSON := json.Valid(bytes.TrimSpace(data))
	ext := ".txt"
	if isJSON {
		ext = ".json"
	}
	// os.CreateTemp honours TMPDIR and creates the file readable by the owner only.
	f, err := os.CreateTemp("", "keen-output-*"+ext)
	if err != nil {
		_, _ = w.Write(data) // never lose the result
		return fmt.Errorf("--output-summary: could not create the output file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		_, _ = w.Write(data)
		return fmt.Errorf("--output-summary: could not write %s: %w", f.Name(), err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("--output-summary: could not write %s: %w", f.Name(), err)
	}

	summary := map[string]any{"file": f.Name(), "bytes": len(data)}
	if isJSON {
		var v any
		if json.Unmarshal(data, &v) == nil {
			summary["records"] = countRecords(v)
			if keys := topKeys(v); len(keys) > 0 {
				summary["keys"] = keys
			}
		}
	}
	out, _ := json.Marshal(map[string]any{"output_summary": summary})
	_, err = fmt.Fprintln(w, string(out))
	return err
}

// countRecords counts the items in the main list of a response: a bare array,
// or the first array found at the top level or one level down (keen's
// {"results": {"issues": [...]}} shape). A single object counts as 1.
func countRecords(v any) int {
	switch t := v.(type) {
	case []any:
		return len(t)
	case map[string]any:
		for _, inner := range []any{t["results"], t["data"], t} {
			m, ok := inner.(map[string]any)
			if !ok {
				if arr, ok := inner.([]any); ok {
					return len(arr)
				}
				continue
			}
			for _, k := range mapKeys(m) {
				if arr, ok := m[k].([]any); ok {
					return len(arr)
				}
			}
		}
		return 1
	}
	return 1
}

func topKeys(v any) []string {
	m, ok := v.(map[string]any)
	if !ok {
		if arr, ok := v.([]any); ok && len(arr) > 0 {
			m, _ = arr[0].(map[string]any)
		}
	}
	keys := mapKeys(m)
	if len(keys) > 12 {
		keys = keys[:12]
	}
	return keys
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// checkSelectPaths reports the first --select path that matches nothing in
// data, with the field names available where the match stopped. An empty
// result (no records) is not checked: there is nothing to match against.
func checkSelectPaths(data json.RawMessage, fields string) error {
	var v any
	if json.Unmarshal(data, &v) != nil || isEmptyResult(v) {
		return nil
	}
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		parts := strings.Split(strings.ToLower(f), ".")
		if ok, at, avail := selectPathExists(v, parts, 0); !ok {
			guess := closestKey(parts[at], avail)
			// Jira leaves out fields that have no value (an issue without a
			// parent has no "parent"). A nested field with no near-miss is
			// treated as absent here, not as a mistake.
			if guess == "" && at > 0 {
				fmt.Fprintf(os.Stderr, "warning: --select: %q is not in this response; left out\n", f)
				continue
			}
			where := "at the top level"
			if at > 0 {
				where = "under " + strings.Join(strings.Split(f, ".")[:at], ".")
			}
			hint := ""
			if guess != "" {
				prefix := strings.Join(strings.Split(f, ".")[:at], ".")
				if prefix != "" {
					guess = prefix + "." + guess
				}
				hint = fmt.Sprintf(" Did you mean %q?", guess)
			}
			shown := avail
			if len(shown) > 15 {
				shown = append(append([]string{}, shown[:15]...), "…")
			}
			return usageErr(fmt.Errorf("--select: no field %q in this response (%s).%s Fields there: %s",
				f, where, hint, strings.Join(shown, ", ")))
		}
	}
	return nil
}

func isEmptyResult(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case []any:
		return len(t) == 0
	case map[string]any:
		if len(t) == 0 {
			return true
		}
		// A list envelope whose lists are all empty.
		found := false
		for _, x := range t {
			if arr, ok := x.([]any); ok {
				found = true
				if len(arr) > 0 {
					return false
				}
			}
		}
		return found
	}
	return false
}

// selectPathExists mirrors filterFields: case-insensitive keys, camelCase
// matching kebab-case, arrays searched element-wise, and list envelopes
// searched inside their arrays. It returns where the match stopped and the
// keys available there.
func selectPathExists(v any, parts []string, depth int) (bool, int, []string) {
	if len(parts) == 0 {
		return true, depth, nil
	}
	switch t := v.(type) {
	case []any:
		var avail []string
		best := depth
		for _, el := range t {
			ok, at, keys := selectPathExists(el, parts, depth)
			if ok {
				return true, at, nil
			}
			if at >= best {
				best, avail = at, keys
			}
		}
		return false, best, avail
	case map[string]any:
		for k, child := range t {
			if strings.ToLower(k) == parts[0] || camelToKebab(k) == parts[0] {
				return selectPathExists(child, parts[1:], depth+1)
			}
		}
		// List envelope fallback, as in filterFieldsRec.
		for _, child := range t {
			if arr, ok := child.([]any); ok && len(arr) > 0 {
				if ok, at, keys := selectPathExists(arr, parts, depth); ok {
					return true, at, keys
				}
			}
		}
		return false, depth, mapKeys(t)
	}
	return false, depth, nil
}

// closestKey returns the available key nearest to want (case-insensitive),
// if it is within a couple of typos.
func closestKey(want string, avail []string) string {
	best, bestDist := "", 3
	for _, k := range avail {
		if d := editDistance(want, strings.ToLower(k)); d < bestDist {
			best, bestDist = k, d
		}
	}
	return best
}

func editDistance(a, b string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}
