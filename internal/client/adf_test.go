package client

import (
	"encoding/json"
	"testing"
)

func TestWrapRichText(t *testing.T) {
	adfDoc := map[string]any{"type": "doc", "version": 1, "content": []any{}}
	cases := []struct {
		name, method, path string
		body               any
		want               string
	}{
		{"comment string", "POST", "/rest/api/3/issue/PROJ-1/comment", map[string]any{"body": "a\nb"},
			`{"body":{"content":[{"content":[{"text":"a","type":"text"}],"type":"paragraph"},{"content":[{"text":"b","type":"text"}],"type":"paragraph"}],"type":"doc","version":1}}`},
		{"comment ADF untouched", "POST", "/rest/api/3/issue/PROJ-1/comment", map[string]any{"body": adfDoc},
			`{"body":{"content":[],"type":"doc","version":1}}`},
		{"comment edit with query", "PUT", "/rest/api/3/issue/PROJ-1/comment/100?notifyUsers=false", map[string]any{"body": "x"},
			`{"body":{"content":[{"content":[{"text":"x","type":"text"}],"type":"paragraph"}],"type":"doc","version":1}}`},
		{"worklog comment", "POST", "/rest/api/3/issue/PROJ-1/worklog", map[string]any{"comment": "", "timeSpent": "1h"},
			`{"comment":{"content":[{"type":"paragraph"}],"type":"doc","version":1},"timeSpent":"1h"}`},
		{"issue create description", "POST", "/rest/api/3/issue", json.RawMessage(`{"fields":{"summary":"s","description":"d"}}`),
			`{"fields":{"description":{"content":[{"content":[{"text":"d","type":"text"}],"type":"paragraph"}],"type":"doc","version":1},"summary":"s"}}`},
		{"v2 path untouched", "POST", "/rest/api/2/issue/PROJ-1/comment", map[string]any{"body": "x"}, `{"body":"x"}`},
		{"GET untouched", "GET", "/rest/api/3/issue/PROJ-1/comment", map[string]any{"body": "x"}, `{"body":"x"}`},
		{"other path untouched", "POST", "/rest/api/3/issue/PROJ-1/transitions", map[string]any{"body": "x"}, `{"body":"x"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(wrapRichText(tc.method, tc.path, tc.body))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got  %s\nwant %s", got, tc.want)
			}
		})
	}
}
