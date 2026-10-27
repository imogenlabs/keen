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
)

// wrapRichText returns body with any plain-string rich-text field replaced
// by an ADF doc. Values that are already objects pass through untouched.
func wrapRichText(method, path string, body any) any {
	if method != "POST" && method != "PUT" {
		return body
	}
	path, _, _ = strings.Cut(path, "?")
	var keys []string
	nested := false
	switch {
	case commentPath.MatchString(path):
		keys = []string{"body"}
	case worklogPath.MatchString(path):
		keys = []string{"comment"}
	case issuePath.MatchString(path):
		keys, nested = []string{"description", "environment"}, true
	default:
		return body
	}

	m, ok := asMap(body)
	if !ok {
		return body
	}
	target := m
	if nested {
		if target, ok = m["fields"].(map[string]any); !ok {
			return body
		}
	}
	for _, k := range keys {
		if s, ok := target[k].(string); ok {
			target[k] = textToADF(s)
		}
	}
	return m
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
