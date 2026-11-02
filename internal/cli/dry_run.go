// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"context"
	"encoding/json"

	"keen/internal/client"
)

// Local customization; see .printing-press-patches.json (dry-run-output).
// A dry run sends nothing, so a read must not look like a live answer: no
// result count, no write to the local copy (which warned that it could find
// no ids), and no "source": "live" in the envelope.

// dryRunSource is the provenance of a read that was only printed.
const dryRunSource = "dry_run"

// dryRunRead prints the request the read would send and returns the client's
// dry-run marker with dry-run provenance.
func dryRunRead(ctx context.Context, c *client.Client, path string, params, headers map[string]string) (json.RawMessage, DataProvenance, error) {
	data, err := c.GetWithHeaders(ctx, path, params, headers)
	return data, DataProvenance{Source: dryRunSource}, err
}
