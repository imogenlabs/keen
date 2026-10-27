# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- `keen-mcp` reads `~/.config/keen/auth.json` like the CLI does; before, it only looked at `config.toml` and environment variables
- `JIRA_CLOUD_PLATFORM_URL` sets the site URL, matching the username and password variables. `JIRA_PP_CLI_BASE_URL` still works and wins if both are set

## [1.1.0] - 2026-10-27

### Added
- `keen release sync`: record git release tags as Jira versions (released, dated, linked to the GitHub release) and add each as a Fix Version on the issues its commits name, with an unreleased `next` version that the next release renames
- Top-level shortcuts: `keen find-issues`, `keen transition`, `keen sprints`, `keen get`
- `--transition-name` on `keen transition` resolves the name to the transition id; an unknown name lists the available transitions
- Credentials can live in `~/.config/keen/auth.json` instead of environment variables
- Compact output (`--agent`) keeps the issue `summary` and flattens Jira's nested `fields`

### Changed
- Built with Go 1.26.10. Cobra, pflag, go-toml and mcp-go (v1.4.0) are updated. SQLite stays on the driver 1.0.0 shipped with
- `keen-mcp` identifies itself to MCP clients as `keen` with its release version, not "Jira Pp Cli" 1.0.0
- Builds made without release flags report version `dev`
- Release archives are named `keen_<os>_<arch>.tar.gz` with no version, so `releases/latest/download/` links stay stable
- The rename is complete: module, commands, help text and config all say `keen`
- The repository moved to the imogenlabs GitHub organization
- CI and release workflows pin every action by commit SHA; CI's token is read-only

### Fixed
- 14 Go standard library vulnerabilities reachable from keen (net/http, crypto/tls, crypto/x509 and others)
- The download links in the README returned 404
- A missing or placeholder `base_url` now fails with a clear error instead of calling a nonexistent host
- `keen sync` no longer warns about skipped records for pagination envelopes

### Removed
- Homebrew formula publishing. No formula was ever published; install from the release archives

## [1.0.0] - 2026-05-27

### Added
- 1,091 Atlassian API endpoints (Jira Platform v3, Jira Agile, Confluence, Marketplace, Service Desk, Admin Orgs)
- SQLite-backed local sync with FTS5 full-text search
- MCP server (`keen-mcp`) for AI agent integration
- Jira API v3 support (migrated from deprecated v2)
- Cross-platform binaries: darwin/arm64, linux/amd64, linux/arm64
- Agent mode (`--agent`) for structured JSON output
- Command discovery via `keen which "<keyword>"`
- Local analytics, stale item detection, workload analysis

### Fixed
- Nullable FK columns in typed tables for partial sync data
- v1-to-v2 schema migration for databases missing resource_type column
- FTS5 ambiguous rank column in search queries
