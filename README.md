# Keen

<p align="center"><strong>The complete Atlassian platform CLI.</strong></p>

<p align="center">
    <a href="https://github.com/imogenlabs/keen/releases/latest">
        <img src="https://img.shields.io/github/v/release/smashingtags/keen?label=Release&logo=github" alt="Release">
    </a>
    <a href="https://github.com/imogenlabs/keen/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="Apache 2.0 License">
    </a>
    <a href="https://go.dev">
        <img src="https://img.shields.io/github/go-mod/go-version/smashingtags/keen" alt="Go Version">
    </a>
    <a href="https://github.com/imogenlabs/keen/actions/workflows/ci.yml">
        <img src="https://github.com/imogenlabs/keen/actions/workflows/ci.yml/badge.svg" alt="CI">
    </a>
    <a href="https://goreportcard.com/report/github.com/imogenlabs/keen">
        <img src="https://goreportcard.com/badge/github.com/imogenlabs/keen" alt="Go Report Card">
    </a>
</p>

<p align="center">
    <a href="https://imogenlabs.ai">
        <img src="https://img.shields.io/badge/Imogen_Labs-AI-8B5CF6" alt="Imogen Labs">
    </a>
</p>

---

1,091 API endpoints across Jira, Confluence, Agile, Marketplace, Service Desk, and Admin — in a single Go binary. SQLite-backed local sync with full-text search. MCP server for AI agents. Built for automation first, humans second.

## What is Keen?

Every Atlassian workflow you automate is a wall of `curl` calls with hand-rolled JSON payloads and fragile `jq` parsing. Every AI agent that touches Jira hits rate limits after a few dozen operations.

Keen replaces all of that with typed commands, structured output, local SQLite caching with FTS5 full-text search, and an MCP server that gives AI agents fast, uncapped access to your entire Atlassian stack.

**Free and open source.** Apache-2.0 license. No account required. No telemetry.

## API Coverage

| API | Endpoints | What you get |
|-----|-----------|-------------|
| **Jira Platform v3** | 692 | Issues, transitions, search, bulk ops, fields, workflows, permissions |
| **Jira Agile** | 105 | Sprints, boards, backlogs, epics, ranking, velocity |
| **Confluence** | 100 | Pages, blog posts, spaces, search, templates, permissions |
| **Marketplace v2** | 149 | Addons, versions, vendors, distribution, reviews, reporting |
| **Service Desk** | 70 | Requests, queues, customers, SLAs, approvals |
| **Admin Orgs** | 61 | Organizations, directories, groups, policies, domains |

## Your First 5 Minutes

```bash
# 1. Install
brew install smashingtags/tap/keen

# 2. Set your credentials
export JIRA_CLOUD_PLATFORM_USERNAME="you@example.com"
export JIRA_CLOUD_PLATFORM_PASSWORD="your-api-token"
export JIRA_CLOUD_PLATFORM_URL="https://yoursite.atlassian.net"

# 3. Check connectivity
keen doctor

# 4. Sync your Jira data locally
keen sync

# 5. Search instantly (no API calls)
keen search "authentication bug" --data-source local
```

Generate an API token at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

## Install

### Homebrew (macOS / Linux)

```bash
brew install smashingtags/tap/keen
```

### Binary download

Download from [GitHub Releases](https://github.com/imogenlabs/keen/releases/latest):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/imogenlabs/keen/releases/latest/download/keen_darwin_arm64.tar.gz | tar xz
sudo mv keen_darwin_arm64/keen /usr/local/bin/
sudo mv keen_darwin_arm64/keen-mcp /usr/local/bin/

# Linux (x86_64)
curl -L https://github.com/imogenlabs/keen/releases/latest/download/keen_linux_amd64.tar.gz | tar xz
sudo mv keen_linux_amd64/keen /usr/local/bin/
sudo mv keen_linux_amd64/keen-mcp /usr/local/bin/
```

### From source

```bash
git clone https://github.com/imogenlabs/keen.git
cd keen
make build-all   # builds bin/keen and bin/keen-mcp
```

Requires Go 1.22+.

## MCP Server

Keen includes an MCP server (`keen-mcp`) that exposes all 1,091 endpoints as tools for AI agents. Works with Claude Code, Cursor, Windsurf, and any MCP-compatible host.

### Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "keen": {
      "command": "keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "you@example.com",
        "JIRA_CLOUD_PLATFORM_PASSWORD": "your-api-token",
        "JIRA_CLOUD_PLATFORM_URL": "https://yoursite.atlassian.net"
      }
    }
  }
}
```

### Cursor

Add to `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "keen": {
      "command": "keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "you@example.com",
        "JIRA_CLOUD_PLATFORM_PASSWORD": "your-api-token",
        "JIRA_CLOUD_PLATFORM_URL": "https://yoursite.atlassian.net"
      }
    }
  }
}
```

### Windsurf

Add to `~/.codeium/windsurf/mcp_config.json` with the same structure.

## Local Sync & Offline Search

Keen syncs your Atlassian data to a local SQLite database with FTS5 full-text search. Queries run in under 100ms against the local cache instead of hitting the Jira Cloud API.

```bash
keen sync                                          # hydrate local cache
keen search "sprint bug" --data-source local        # FTS search (sub-100ms)
keen analytics                                      # analytics on synced data
keen stale --days 30                                # find stale items
keen orphans                                        # find orphaned items
keen load                                           # workload distribution
```

This is especially useful for AI agents that need to make many queries without hitting [Atlassian's rate limits](https://developer.atlassian.com/cloud/jira/platform/rate-limiting/).

## Agent Mode

Every command supports `--agent` which sets JSON output, no interactive prompts, compact format, and no color in one flag. Designed for AI agent consumption.

```bash
keen agile get-all-boards --agent
keen issue transitions get EOS-42 --agent | jq '.results.transitions[] | {id, name}'
```

## Key Flags

| Flag | What it does |
|------|-------------|
| `--agent` | JSON + no prompts + compact (use this for all scripted calls) |
| `--compact` | Only key fields (id, name, status, timestamps) |
| `--json` | Force JSON output |
| `--csv` | CSV output |
| `--dry-run` | Show what would be sent without sending |
| `--stdin` | Read request body as JSON from stdin |
| `--select "id,name"` | Pick specific fields from response |
| `--quiet` | One value per line |
| `--data-source local` | Query local SQLite cache only |

## Command Discovery

```bash
keen which "sprint"          # find commands by keyword
keen which "create page"     # natural language search
keen api                     # browse all API interfaces
keen api wiki                # browse Confluence endpoints
keen api agile               # browse Agile endpoints
```

## Common Workflows

### Sprint Management

```bash
keen agile get-all-boards --agent
keen agile get-all-sprints BOARD_ID --state active,future --agent
keen agile create-sprint --stdin <<< '{"name":"Sprint 5","originBoardId":216}'
keen agile move-issues-to-sprint-and-rank SPRINT_ID \
  --stdin <<< '{"issues": ["EOS-42", "EOS-43"]}'
keen agile partially-update-sprint SPRINT_ID --stdin <<< '{"state": "active"}'
```

### Issue Lifecycle

```bash
keen jira-cloud-platform-search and-reconsile-issues-using-jql-post \
  --jql "project = EOS AND status != Done" --agent
keen issue transitions get EOS-42 --agent
keen issue transitions do EOS-42 --transition-id 41 --agent
keen issue comment create EOS-42 --stdin <<< '{"body":{"type":"doc","version":1,"content":[{"type":"paragraph","content":[{"type":"text","text":"Done."}]}]}}'
```

### Confluence

```bash
keen wiki get-spaces --agent
keen wiki search-content-by-cql --cql "type=page AND space=EOS" --agent
keen wiki create-or-update-content --stdin <<< '{
  "type": "page", "title": "ADR: API Gateway",
  "space": {"key": "EOS"},
  "body": {"storage": {"value": "<p>Content</p>", "representation": "storage"}}
}'
```

### Marketplace

```bash
keen addons get io.example.my-addon --agent
keen addons list-search --q "story points" --agent
keen addons versions get io.example.my-addon --agent
```

## Comparison

| Feature | Keen | Atlassian Rovo MCP | jira-cli | ACLI (Appfire) |
|---------|------|--------------------|----------|----------------|
| Endpoints | 1,091 | 46+ | Jira only | Many |
| Local cache | SQLite + FTS5 | No | No | No |
| MCP server | Yes | Yes (official) | No | No |
| Products | 6 (Jira, Confluence, Agile, Marketplace, Service Desk, Admin) | 5 | 1 | 4 |
| Price | Free (Apache-2.0) | Free with Cloud | Free (MIT) | $360-$6,195/yr |
| Binary | Single Go binary | Remote server | Single Go binary | Requires JVM |
| Offline search | Yes | No | No | No |

## Architecture

```
┌─────────────┐    ┌─────────────┐
│  keen (CLI)  │    │  keen-mcp   │
│              │    │ (MCP Server)│
└──────┬───────┘    └──────┬──────┘
       │                   │
       └─────────┬─────────┘
                 │
          ┌──────┴──────┐
          │ CLI Router  │
          │  (Cobra)    │
          └──────┬──────┘
                 │
       ┌─────────┴─────────┐
       │                   │
┌──────┴──────┐    ┌───────┴───────┐
│ HTTP Client │    │ SQLite Store  │
│ (API calls) │    │  (FTS5 sync)  │
└──────┬──────┘    └───────────────┘
       │
┌──────┴──────────────────────────┐
│        Atlassian Cloud          │
│ Jira · Confluence · Agile       │
│ Marketplace · Service Desk      │
│ Admin Orgs                      │
└─────────────────────────────────┘
```

## Regeneration

Generated from 6 official Atlassian OpenAPI specs using [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press). Source specs in `specs/`. When Atlassian updates their APIs, regenerate:

```bash
cli-printing-press generate \
  --spec specs/jira-platform.json \
  --spec specs/jira-agile.json \
  --spec specs/confluence.json \
  --spec specs/marketplace.json \
  --spec specs/jira-service-desk.json \
  --spec specs/admin-orgs.json \
  --name keen --owner smashingtags --force --lenient \
  --auth-preference basicAuth --max-resources 1500
```

Do not hand-edit generated Go files. They are regenerated from specs.

## Security

See [SECURITY.md](SECURITY.md) for vulnerability reporting. Keen stores API credentials in environment variables only — never on disk. The local SQLite database contains cached API responses and should be treated as sensitive.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug fixes welcome. Open an issue before submitting PRs for new features.

## License

Apache-2.0

---

<p align="center">Built by <a href="https://imogenlabs.ai">Imogen Labs</a></p>
