package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const issueJSON = `{"key":"PROJ-1","id":"10","fields":{"summary":"Fix login","status":{"name":"Done"},"labels":["a","b"],"components":[{"name":"api"}]}}`

func TestFilterFieldsNestedMissDoesNotLeak(t *testing.T) {
	got := string(filterFields(json.RawMessage(issueJSON), "key,fields.parent"))
	if strings.Contains(got, "summary") || strings.Contains(got, "components") {
		t.Fatalf("missing nested field exposed its siblings: %s", got)
	}
	if !strings.Contains(got, `"PROJ-1"`) {
		t.Fatalf("key was dropped: %s", got)
	}
}

func TestFilterFieldsTopLevelEnvelopeStillWorks(t *testing.T) {
	list := `{"isLast":true,"issues":[{"key":"PROJ-1","id":"10"},{"key":"PROJ-2","id":"11"}]}`
	got := string(filterFields(json.RawMessage(list), "key"))
	if !strings.Contains(got, `"PROJ-2"`) || strings.Contains(got, `"10"`) {
		t.Fatalf("envelope fallback broken: %s", got)
	}
}

func TestCheckSelectPaths(t *testing.T) {
	cases := []struct {
		name, data, sel string
		wantErr         string
	}{
		{"valid", issueJSON, "key,fields.summary,fields.status.name", ""},
		{"typo nested", issueJSON, "fields.summery", `Did you mean "fields.summary"?`},
		{"typo top", issueJSON, "kye", `Did you mean "key"?`},
		{"unknown top", issueJSON, "zzz", `no field "zzz"`},
		{"absent optional field", issueJSON, "fields.parent.key", ""},
		{"list element-wise", `[{"key":"A"},{"key":"B","fields":{"parent":{"key":"E"}}}]`, "fields.parent.key", ""},
		{"empty list skipped", `{"issues":[]}`, "anything", ""},
		{"camelCase as kebab", `{"issueType":"Bug"}`, "issue-type", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkSelectPaths(json.RawMessage(tc.data), tc.sel)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
			if !isCobraUsageError(err) && ExitCode(err) != 2 {
				t.Fatalf("bad --select should exit 2, got %d", ExitCode(err))
			}
		})
	}
}

func TestFlushOutputSummary(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	defer func() { outputSummary.limit, outputSummary.buf, outputSummary.selectErr = 0, nil, nil }()

	// Small output is printed as is.
	outputSummary.limit = 100
	outputSummary.buf = bytes.NewBufferString(`{"key":"PROJ-1"}`)
	var out bytes.Buffer
	if err := flushOutputSummary(&out); err != nil || out.String() != `{"key":"PROJ-1"}` {
		t.Fatalf("small output: %q %v", out.String(), err)
	}

	// Large output goes to a private file; stdout gets the summary.
	big := `{"results":{"issues":[` + strings.Repeat(`{"key":"PROJ-1"},`, 20) + `{"key":"PROJ-2"}]}}`
	outputSummary.buf = bytes.NewBufferString(big)
	out.Reset()
	if err := flushOutputSummary(&out); err != nil {
		t.Fatal(err)
	}
	var env struct {
		S struct {
			File    string   `json:"file"`
			Bytes   int      `json:"bytes"`
			Records int      `json:"records"`
			Keys    []string `json:"keys"`
		} `json:"output_summary"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("summary is not JSON: %q", out.String())
	}
	if env.S.Bytes != len(big) || env.S.Records != 21 || !strings.HasSuffix(env.S.File, ".json") {
		t.Fatalf("summary wrong: %+v", env.S)
	}
	saved, err := os.ReadFile(env.S.File)
	if err != nil || string(saved) != big {
		t.Fatalf("file does not hold the full output: %v", err)
	}
	if st, _ := os.Stat(env.S.File); st.Mode().Perm() != 0o600 {
		t.Fatalf("file should be readable by the owner only, is %v", st.Mode().Perm())
	}

	// A bad --select replaces the output with the error.
	outputSummary.buf = bytes.NewBufferString(`{"key":"PROJ-1"}`)
	outputSummary.selectErr = checkSelectPaths(json.RawMessage(issueJSON), "kye")
	out.Reset()
	if err := flushOutputSummary(&out); err == nil || out.Len() != 0 {
		t.Fatalf("bad --select should print nothing and fail: %q %v", out.String(), err)
	}
}

func TestSelectEnvelopePrefix(t *testing.T) {
	// keen prints a write under "data" and a read under "results"; a path
	// copied from that output works the same as the bare one.
	for _, sel := range []string{"data.key", "results.key", "Data.key"} {
		got := string(filterFields(json.RawMessage(issueJSON), sel))
		if got != `{"key":"PROJ-1"}` {
			t.Fatalf("%s: got %s", sel, got)
		}
	}
	// A response that really has a "data" field keeps it.
	if got := stripEnvelopePrefix(json.RawMessage(`{"data":{"id":1}}`), "data.id"); got != "data.id" {
		t.Fatalf("real data field was stripped: %s", got)
	}
}

func TestSelectMissAfterWrite(t *testing.T) {
	defer func() {
		outputSummary.limit, outputSummary.buf, outputSummary.selectErr, outputSummary.selectData = 0, nil, nil, nil
	}()

	// A write that went through prints the whole response and does not fail.
	outputSummary.buf = bytes.NewBufferString(`{"action":"post","data":{},"status":201,"success":true}`)
	filterFields(json.RawMessage(`{"id":"10001","self":"x"}`), "nope")
	var out bytes.Buffer
	if err := flushOutputSummary(&out); err != nil {
		t.Fatalf("successful write should exit 0: %v", err)
	}
	var env struct {
		Data map[string]string `json:"data"`
	}
	if json.Unmarshal(out.Bytes(), &env) != nil || env.Data["id"] != "10001" || env.Data["self"] != "x" {
		t.Fatalf("successful write should print the whole response: %q", out.String())
	}

	// A dry run (success false) and a read still fail and print nothing.
	for _, held := range []string{
		`{"action":"post","data":null,"dry_run":true,"status":0,"success":false}`,
		`{"meta":{"source":"live"},"results":{}}`,
	} {
		outputSummary.buf = bytes.NewBufferString(held)
		filterFields(json.RawMessage(`{"id":"10001"}`), "nope")
		out.Reset()
		if err := flushOutputSummary(&out); err == nil || out.Len() != 0 || ExitCode(err) != 2 {
			t.Fatalf("%s: should fail with exit 2 and print nothing: %q %v", held, out.String(), err)
		}
	}
}
