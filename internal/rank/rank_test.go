// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package rank

import "testing"

func TestPlainCommandBeatsLongerOne(t *testing.T) {
	plain, _ := NewItem("issue transitions do", "Performs an issue transition.", false).Score("transition an issue")
	other, _ := NewItem("issue edit", "Edits an issue. issue transition is not supported.", false).Score("transition an issue")
	if plain <= other {
		t.Fatalf("issue transitions do scored %d, issue edit %d", plain, other)
	}
}

func TestDescriptionOnlyIsNotNamed(t *testing.T) {
	if _, named := NewItem("bitbucket-addon delete-app", "Close the app.", false).Score("close a sprint"); named {
		t.Fatal("a description-only hit must not count as a named match")
	}
}

func TestKeywordsCountWithoutCost(t *testing.T) {
	it := NewItem("transition", "Transition a Jira issue", false).WithKeywords("Transition a Jira issue")
	score, named := it.Score("transition issue")
	if !named || score != 10 {
		t.Fatalf("got score %d named %v, want 10 true", score, named)
	}
}

func TestExactNameWins(t *testing.T) {
	if s, _ := NewItem("issue_transitions_do", "", true).Score("issue_transitions_do"); s != Exact {
		t.Fatalf("exact name scored %d", s)
	}
}

func TestDeprecatedRanksLower(t *testing.T) {
	a, _ := NewItem("issue search", "", false).Score("issue search")
	b, _ := NewItem("issue search-post", "", true).Score("issue search")
	if b >= a {
		t.Fatalf("deprecated scored %d, current %d", b, a)
	}
}
