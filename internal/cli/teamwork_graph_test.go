// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTeamworkGraphQuery(t *testing.T) {
	q := graphQuery{Field: "teamworkGraph_userTeams", OptIn: "TeamworkGraphContextAPIs",
		Args: []graphArg{{Name: "userId", Type: "ID!"}, {Name: "first", Type: "Int"}, {Name: "after", Type: "String"}}}
	got := teamworkGraphQuery(q, map[string]any{"userId": "ari:cloud:identity::user/1", "first": 5})
	for _, want := range []string{
		"query KeenTeamworkGraph($userId: ID!, $first: Int)",
		"teamworkGraph_userTeams(userId: $userId, first: $first)",
		`@optIn(to: "TeamworkGraphContextAPIs")`,
		"pageInfo { hasNextPage endCursor }",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("query is missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "$after") {
		t.Errorf("query declares an argument that was not given:\n%s", got)
	}
}

func TestEveryTeamworkGraphQueryHasACommand(t *testing.T) {
	root := RootCmd()
	for _, q := range teamworkGraphQueries {
		cmd, _, err := root.Find([]string{"teamwork-graph", q.Command})
		if err != nil || cmd.Name() != q.Command {
			t.Errorf("keen teamwork-graph %s: not found (%v)", q.Command, err)
		}
		if strings.TrimSuffix(q.Returns, "!") != cypherConnection {
			t.Errorf("%s returns %s; keen only reads %s", q.Field, q.Returns, cypherConnection)
		}
	}
	if len(teamworkGraphQueries) == 0 {
		t.Fatal("no Teamwork Graph queries in the table")
	}
}

func TestCypherRows(t *testing.T) {
	raw := `{"edges":[{"node":{"columns":[
		{"key":"entity","value":{"__typename":"GraphStoreCypherQueryV2AriNode","id":"ari:x","data":{"__typename":"JiraIssue"}}},
		{"key":"teams","value":{"__typename":"GraphStoreCypherQueryV2NodeList","nodes":[{"id":"ari:t","data":{"__typename":"TeamV2"}}]}},
		{"key":"count","value":{"__typename":"GraphStoreCypherQueryV2IntObject","int":3}},
		{"key":"manager","value":null}]}}],
		"pageInfo":{"hasNextPage":true,"endCursor":"c1"}}`
	var r cypherResult
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(r.rows())
	want := `{"pageInfo":{"endCursor":"c1","hasNextPage":true},"rows":[{"count":3,"entity":{"id":"ari:x","type":"JiraIssue"},"manager":null,"teams":[{"id":"ari:t","type":"TeamV2"}]}]}`
	if string(out) != want {
		t.Errorf("rows =\n%s\nwant\n%s", out, want)
	}
}

func TestARIPrefixAndNames(t *testing.T) {
	if got := withARIPrefix("userId", "70121:abc"); got != "ari:cloud:identity::user/70121:abc" {
		t.Errorf("user id = %s", got)
	}
	if got := withARIPrefix("userId", "ari:cloud:identity::user/x"); got != "ari:cloud:identity::user/x" {
		t.Errorf("an ARI was changed: %s", got)
	}
	if got := argFlagName("viewedEntityTypes"); got != "viewed-entity-types" {
		t.Errorf("flag = %s", got)
	}
	if got := operationName("query KeenX($a: ID!) { x }"); got != "KeenX" {
		t.Errorf("operation = %q", got)
	}
}
