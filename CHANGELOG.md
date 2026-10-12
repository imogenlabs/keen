# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `keen release sync`: record git release tags as Jira versions (released, dated, linked to the GitHub release) and add each as a Fix Version on the issues its commits name, with an unreleased `next` version that the next release renames

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
