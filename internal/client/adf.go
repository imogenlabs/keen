// Local customization: wrap plain-text rich-text fields into ADF for Jira v3.
// See .printing-press-patches.json (adf-plain-text).

package client

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Jira Platform v3 rejects a plain string where it expects an Atlassian
// Document Format doc ("Comment body is not valid!"). These are the
// write endpoints whose generated flags take such fields as strings.
var (
	commentPath = regexp.MustCompile(`^/rest/api/3/issue/[^/]+/comment(/[^/]+)?$`)
	worklogPath = regexp.MustCompile(`^/rest/api/3/issue/[^/]+/worklog(/[^/]+)?$`)
	issuePath   = regexp.MustCompile(`^/rest/api/3/issue(/[^/]+)?$`)
	bulkPath    = regexp.MustCompile(`^/rest/api/3/issue/bulk$`)
)

var issueRichText = []string{"description", "environment"}

// wrapRichText returns body with any plain-string rich-text field replaced
// by an ADF doc. Values that are already objects pass through untouched.
func wrapRichText(method, path string, body any) any {
	if method != "POST" && method != "PUT" {
		return body
	}
	path, _, _ = strings.Cut(path, "?")
	// issuePath also matches /issue/bulk.
	if !commentPath.MatchString(path) && !worklogPath.MatchString(path) && !issuePath.MatchString(path) {
		return body
	}
	m, ok := asMap(body)
	if !ok {
		return body
	}
	switch {
	case commentPath.MatchString(path):
		wrapKeys(m, "body")
	case worklogPath.MatchString(path):
		wrapKeys(m, "comment")
	case bulkPath.MatchString(path):
		// create-bulk: {"issueUpdates": [{"fields": {...}}, ...]}
		updates, _ := m["issueUpdates"].([]any)
		for _, u := range updates {
			if um, ok := u.(map[string]any); ok {
				if fields, ok := um["fields"].(map[string]any); ok {
					wrapKeys(fields, issueRichText...)
				}
			}
		}
	case issuePath.MatchString(path):
		if fields, ok := m["fields"].(map[string]any); ok {
			wrapKeys(fields, issueRichText...)
		}
	default:
		return body
	}
	return m
}

func wrapKeys(m map[string]any, keys ...string) {
	for _, k := range keys {
		if s, ok := m[k].(string); ok {
			m[k] = textToADF(s)
		}
	}
}

func asMap(body any) (map[string]any, bool) {
	switch b := body.(type) {
	case map[string]any:
		return b, true
	case json.RawMessage, []byte:
		var raw []byte
		if r, ok := b.(json.RawMessage); ok {
			raw = r
		} else {
			raw = b.([]byte)
		}
		var m map[string]any
		if json.Unmarshal(raw, &m) != nil {
			return nil, false
		}
		return m, true
	}
	return nil, false
}

// textToADF makes one paragraph per line; blank lines become empty paragraphs.
func textToADF(s string) map[string]any {
	var content []any
	for _, line := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		p := map[string]any{"type": "paragraph"}
		if line != "" {
			p["content"] = []any{map[string]any{"type": "text", "text": line}}
		}
		content = append(content, p)
	}
	return map[string]any{"type": "doc", "version": 1, "content": content}
}
