// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

// keyStatus reports an optional API key without revealing it.
func keyStatus(key, env, usedFor string) string {
	if key != "" {
		return "configured"
	}
	return "not set (only needed for " + usedFor + "; " + env + ")"
}
