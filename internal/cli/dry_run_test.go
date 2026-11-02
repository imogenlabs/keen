// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"testing"
)

func TestDryRunEnvelopeIsNotLive(t *testing.T) {
	out, err := wrapWithProvenance(json.RawMessage(`{"dry_run": true}`), DataProvenance{Source: dryRunSource})
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		Meta map[string]any `json:"meta"`
	}
	if err := json.Unmarshal(out, &env); err != nil {
		t.Fatal(err)
	}
	if env.Meta["source"] == "live" || env.Meta["dry_run"] != true {
		t.Errorf("dry-run meta = %v", env.Meta)
	}
}
