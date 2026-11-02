// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package client

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"keen/internal/config"
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

// End to end: a dry run with real-looking credentials prints no four
// characters of the password or of the encoded header, wherever it prints.
func TestDryRunOutputShowsNoPartOfTheSecret(t *testing.T) {
	password := "ATATT3xFfGF0zyxwvutsrqponmlkjihg"
	cfg := &config.Config{BaseURL: "https://example.atlassian.net", JiraCloudPlatformUsername: "someone@example.com", JiraCloudPlatformPassword: password}
	c := New(cfg, 5*time.Second, 0)
	c.DryRun = true

	r, w, _ := os.Pipe()
	stderr := os.Stderr
	os.Stderr = w
	_, err := c.Get(context.Background(), "/rest/api/3/myself", map[string]string{"expand": "groups"})
	w.Close()
	os.Stderr = stderr
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)
	encoded := base64.StdEncoding.EncodeToString([]byte(cfg.JiraCloudPlatformUsername + ":" + password))
	for _, secret := range []string{password, encoded} {
		for i := 0; i+4 <= len(secret); i++ {
			if strings.Contains(string(out), secret[i:i+4]) {
				t.Fatalf("dry-run output leaks %q:\n%s", secret[i:i+4], out)
			}
		}
	}
	if !strings.Contains(string(out), "Authorization: Basic ****") {
		t.Errorf("expected the masked header in:\n%s", out)
	}
}
