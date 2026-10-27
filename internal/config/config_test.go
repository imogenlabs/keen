package config

import (
	"os"
	"path/filepath"
	"testing"
)

func cleanEnv(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, k := range []string{"JIRA_PP_CLI_CONFIG", "JIRA_PP_CLI_BASE_URL", "JIRA_CLOUD_PLATFORM_URL",
		"JIRA_CLOUD_PLATFORM_USERNAME", "JIRA_CLOUD_PLATFORM_PASSWORD"} {
		t.Setenv(k, "")
	}
	return home
}

func TestLoadEmptyPathReadsAuthJSON(t *testing.T) {
	home := cleanEnv(t)
	dir := filepath.Join(home, ".config", "keen")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	body := `{"base_url":"https://example.atlassian.net","username":"u@example.com","password":"tok"}`
	if err := os.WriteFile(filepath.Join(dir, "auth.json"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://example.atlassian.net" || cfg.JiraCloudPlatformUsername != "u@example.com" || cfg.JiraCloudPlatformPassword != "tok" {
		t.Fatalf("auth.json not applied: %+v", cfg)
	}
}

func TestLoadBaseURLFromEnv(t *testing.T) {
	cleanEnv(t)
	t.Setenv("JIRA_CLOUD_PLATFORM_URL", "https://alias.atlassian.net")
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://alias.atlassian.net" {
		t.Fatalf("JIRA_CLOUD_PLATFORM_URL ignored: %q", cfg.BaseURL)
	}

	t.Setenv("JIRA_PP_CLI_BASE_URL", "https://wins.atlassian.net")
	cfg, err = Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://wins.atlassian.net" {
		t.Fatalf("JIRA_PP_CLI_BASE_URL should win: %q", cfg.BaseURL)
	}
}
