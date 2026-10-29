# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.1] - 2026-10-28

### Fixed
- `keen issue create-bulk` accepts plain-text `description` and `environment` fields, as `issue create` and `issue edit` already did. keen wraps them into Atlassian Document Format; before, Jira refused them with "must be ADF"
- Help examples show a value that fits each argument: `PROJ-42` for an issue key, `PROJ` for a project, `TEAM` for a space, `com.example.my-app` for an app. Before, 277 commands showed `your-token-here` there. The setup hint now asks for your email as the username, not a token

## [1.4.0] - 2026-10-28

### Added
- `--output-summary[=N]`: output over N bytes (default 4096) goes to a private temp file, and keen prints a one-line JSON summary with the file path, size, record count and top-level keys. Smaller output prints as usual. Works on every command

### Fixed
- `--select` with a mistyped field now fails (exit code 2) with a "did you mean" suggestion and the fields that exist there. Before, it could print the whole response
- `--select` on a nested field a record doesn't have (for example `fields.parent` on an issue without a parent) no longer prints every other field of that object. It prints a warning and leaves the field out

## [1.3.0] - 2026-10-27

### Added
- Regenerated from Atlassian's current OpenAPI specs (October 2026): 60 new commands and MCP tools. Among them: Jira bulk operations (`keen bulk submit-edit`, `submit-move`, `submit-transition`, `submit-delete`, watch/unwatch, progress), Jira policies (`keen policies list`, `apply`, `set-project`), `keen workflows copy`, Confluence content restrictions, user properties, page moves and macros, four Service Desk customer and knowledge base commands, and directory search in Atlassian Admin
- `scripts/regen/`: one command refreshes keen from Atlassian's current specs and merges the result onto keen's own changes; `check.sh` gates the result

### Fixed
- Commands that start a background job and wait for it (Jira bulk operations, workflow updates, Confluence body conversions) now poll Atlassian's real progress endpoint. Before, several polled an unrelated endpoint or a URL with a literal `{taskId}` in it
- `scripts/regen/regen.sh` can run repeatedly; a run on unchanged specs changes nothing

### Changed
- The org policy commands now live under `keen atlassian-access-policies`, because Jira added its own policies API. The old `keen policies create-policy`, `get`, `resources ...` names and MCP tools still work
- Commands for 21 APIs Atlassian retired (the Admin v1 directory API, app access settings, workflow transition properties) still exist under their old names and MCP tool names. They now say which command replaced them instead of failing with HTTP 401

## [1.2.2] - 2026-10-27

### Fixed
- The forms `keen which` suggests for shortcut commands, such as `keen myself get-current-user` and `keen server-info get`, work again. 1.2.0 and 1.2.1 refused the operation name as a stray word

## [1.2.1] - 2026-10-27

### Fixed
- `keen get KEY`, `keen sprints BOARD` and `keen transition KEY` work again. 1.2.0 refused their argument with "unknown command"

## [1.2.0] - 2026-10-27

### Added
- `keen wiki create-page`, `keen wiki get-page` and `keen wiki update-page` for Confluence pages. `update-page` reads the current version and sends the next one

### Fixed
- `keen search` works against Jira again. It called `/rest/api/3/search`, which Atlassian removed (HTTP 410); it now runs a JQL text search
- `--all` fetches every page from Jira, Agile and Service Desk lists that page with `startAt`; before, it stopped after the first page (for example 50 of 209 sprints)
- Plain text in `issue comment add --body`, worklog comments, and an issue's `description` or `environment` is wrapped into Atlassian Document Format. Jira v3 rejected it before with "Comment body is not valid!"
- `keen projects create` and other shortcut commands that take no arguments now refuse stray words instead of silently running a different request. The error points at the right command (`keen project create`)
- `keen-mcp` reads `~/.config/keen/auth.json` like the CLI does; before, it only looked at `config.toml` and environment variables
- `JIRA_CLOUD_PLATFORM_URL` sets the site URL, matching the username and password variables. `JIRA_PP_CLI_BASE_URL` still works and wins if both are set

### Changed
- Docs, `AGENTS.md`, `SKILL.md` and the tool manifest use the keen name throughout. The old install steps and generator names are gone, and the docs show `JIRA_CLOUD_PLATFORM_URL`
- The GitHub wiki is turned off. It was an outdated copy of `docs/`, still written for the old name
- SQLite driver updated (modernc.org/sqlite 1.61.0). The local database now really runs in WAL mode with a busy timeout, as keen always asked for; the old driver ignored those settings. Opening a fresh database from several processes at once retries instead of failing with "database is locked"

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
