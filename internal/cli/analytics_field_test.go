// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"testing"
)

func TestGroupByValueReadsIssueFields(t *testing.T) {
	var issue map[string]any
	_ = json.Unmarshal([]byte(`{"key":"PROJ-1","fields":{"status":{"name":"In Progress"},"assignee":null,
		"issuetype":{"name":"Bug"},"priority":{"name":"High"},"labels":["cli"]}}`), &issue)
	cases := map[string]string{
		"status":             "In Progress",
		"fields.status.name": "In Progress",
		"type":               "Bug",
		"priority":           "High",
		"assignee":           "(none)",
		"key":                "PROJ-1",
		"issuetype":          "Bug",
		"missing":            "(none)",
	}
	for field, want := range cases {
		if got := groupByValue(issue, field); got != want {
			t.Errorf("groupByValue(%q) = %q, want %q", field, got, want)
		}
	}
}

func TestProjectDeleteTakesAProject(t *testing.T) {
	root := newRootCmd(&rootFlags{})
	c, rest, err := root.Find([]string{"project", "delete", "PROJ"})
	if err != nil || c.Name() != "delete" || len(rest) != 1 {
		t.Fatalf("project delete PROJ resolved to %v %v %v", c.Name(), rest, err)
	}
	if err := c.Args(c, rest); err != nil {
		t.Fatalf("project delete PROJ rejected: %v", err)
	}
	if c.Annotations["pp:method"] != "DELETE" {
		t.Errorf("annotations = %v", c.Annotations)
	}
	if sub, _, err := root.Find([]string{"project", "delete", "project-asynchronously", "PROJ"}); err != nil || sub.Name() != "project-asynchronously" {
		t.Errorf("subcommand lost: %v %v", sub, err)
	}
}
