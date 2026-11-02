// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"fmt"
	"strings"
)

// Local customization; see .printing-press-patches.json (analytics-group-by).
// Jira keeps status, assignee and type under "fields", so grouping by a
// top-level field reported every issue as <nil>. groupByValue follows dotted
// paths, knows the everyday issue names, and reads a nested object's name.

var issueGroupFields = map[string]string{
	"status":     "fields.status.name",
	"assignee":   "fields.assignee.displayName",
	"reporter":   "fields.reporter.displayName",
	"type":       "fields.issuetype.name",
	"priority":   "fields.priority.name",
	"project":    "fields.project.key",
	"resolution": "fields.resolution.name",
}

// groupByValue returns obj's value for field as a group label.
func groupByValue(obj map[string]any, field string) string {
	v, ok := lookupPath(obj, field)
	if !ok {
		if alias, known := issueGroupFields[field]; known {
			v, ok = lookupPath(obj, alias)
		}
	}
	if !ok {
		if f, isMap := obj["fields"].(map[string]any); isMap {
			v, ok = lookupPath(f, field)
		}
	}
	if !ok || v == nil {
		return "(none)"
	}
	if m, isMap := v.(map[string]any); isMap {
		for _, k := range []string{"name", "displayName", "value", "key"} {
			if s, isStr := m[k].(string); isStr {
				return s
			}
		}
	}
	return fmt.Sprintf("%v", v)
}

func lookupPath(obj map[string]any, path string) (any, bool) {
	var cur any = obj
	for _, part := range strings.Split(path, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = m[part]; !ok {
			return nil, false
		}
	}
	return cur, true
}
