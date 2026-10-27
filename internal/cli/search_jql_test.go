package cli

import "testing"

func TestTextSearchJQLEscapes(t *testing.T) {
	got := textSearchJQL(`say "hi" \ now`)
	want := `text ~ "say \"hi\" \\ now" ORDER BY updated DESC`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
