// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"encoding/base64"
	"strings"
	"testing"
)

// No four characters of a credential may appear in what keen prints for it.
func TestMaskTokenShowsNoPartOfTheSecret(t *testing.T) {
	secret := base64.StdEncoding.EncodeToString([]byte("someone@example.com:ATATT3xFfGF0abcdefghijklmnop"))
	for _, header := range []string{"Basic " + secret, "Bearer " + secret, secret} {
		got := maskToken(header)
		for i := 0; i+4 <= len(secret); i++ {
			if strings.Contains(got, secret[i:i+4]) {
				t.Fatalf("maskToken(%q...) = %q leaks %q", header[:10], got, secret[i:i+4])
			}
		}
	}
	if got := maskToken("Basic " + secret); got != "Basic ****" {
		t.Errorf("got %q, want the scheme kept", got)
	}
}
