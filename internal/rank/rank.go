// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

// Package rank scores a command or MCP tool against a few plain words
// ("transition an issue", "list sprints on a board"). keen which and keen-mcp's
// keen_search both use it, so the CLI and the MCP server cannot disagree about
// which command does a job. Local customization; see
// .printing-press-patches.json (shared-ranking).
package rank

import "strings"

// Item is something that can be found: a CLI command path ("issue transitions
// do") or an MCP tool name ("issue_transitions_do"), with its description.
type Item struct {
	words      map[string]bool // from the name
	keywords   map[string]bool // extra words that count like name words, without the cost
	name       string
	text       string // lower-case name and description
	Deprecated bool
}

// NewItem prepares name and description for scoring.
func NewItem(name, description string, deprecated bool) Item {
	words := map[string]bool{}
	for _, w := range Words(name) {
		if !nameNoise[w] {
			words[w] = true
		}
	}
	return Item{words: words, name: name, text: strings.ToLower(name + " " + description), Deprecated: deprecated}
}

// WithKeywords returns it with the words of text counting like name words, but
// without the cost of unmatched name words. keen's own commands have one-word
// names ("transition", "sprints"); their summary says what they act on.
func (it Item) WithKeywords(text string) Item {
	it.keywords = map[string]bool{}
	for _, w := range Words(text) {
		if !nameNoise[w] {
			it.keywords[w] = true
		}
	}
	return it
}

// Exact is the score of an exact name match; nothing else reaches it.
const Exact = 1000

// Score ranks it against query. Query words found in the name count 5, words
// people use for the action in the name ("move" for "transition") count 3, and
// words found only in the description count 1. Name words the query did not ask
// for cost 2 each, so the plain command beats a longer one sharing a word, and a
// deprecated item loses 5. Keywords (WithKeywords) count like name words.
// named reports whether any query word hit the name or a keyword; a hit on the
// description alone is a weak guess.
func (it Item) Score(query string) (score int, named bool) {
	q := strings.TrimSpace(query)
	if strings.EqualFold(it.name, q) {
		return Exact, true
	}
	matched := map[string]bool{} // name words the query hit
	hit := false
	for _, w := range Words(q) {
		switch {
		case it.words[w]:
			score += 5
			matched[w] = true
		case it.keywords[w]:
			score += 5
			hit = true
		case strings.Contains(it.text, w):
			score++
		}
		for _, alt := range actionWords[w] {
			switch {
			case it.words[alt]:
				score += 3
				matched[alt] = true
			case it.keywords[alt]:
				score += 3
				hit = true
			}
		}
	}
	hit = hit || len(matched) > 0
	if hit {
		score = max(score-2*(len(it.words)-len(matched)), 1)
	}
	if it.Deprecated && score > 0 {
		score = max(score-5, 1)
	}
	return score, hit
}

var stopWords = map[string]bool{"a": true, "an": true, "the": true, "to": true, "of": true, "in": true,
	"for": true, "by": true, "on": true, "my": true, "and": true, "or": true, "with": true,
	"is": true, "are": true, "who": true, "what": true, "which": true, "how": true, "me": true, "i": true,
	"do": true, "does": true, "can": true, "show": true}

// Words the generator puts in names that say nothing about the action.
var nameNoise = map[string]bool{"jira": true, "cloud": true, "platform": true, "using": true, "api": true}

// Words people use for an action, mapped to the words keen's names use.
var actionWords = map[string][]string{
	"list":   {"get", "all", "list"},
	"show":   {"get"},
	"find":   {"search"},
	"search": {"find"},
	"move":   {"transition"},
	"change": {"update", "edit"},
	"remove": {"delete"},
	"make":   {"create"},
	"new":    {"create"},
	"done":   {"transition"},
	"close":  {"transition", "update"},
	"status": {"transition"},
	"ticket": {"issue"},
	"jql":    {"find", "search"},
}

// Words lower-cases s, splits it on anything but letters and digits, drops
// filler words and a plural "s", so "list sprints" matches "sprint".
func Words(s string) []string {
	var out []string
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	}) {
		if stopWords[w] {
			continue
		}
		if len(w) > 3 && strings.HasSuffix(w, "s") && !strings.HasSuffix(w, "ss") {
			w = strings.TrimSuffix(w, "s")
		}
		out = append(out, w)
	}
	return out
}
