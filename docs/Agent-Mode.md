# Agent Mode

The CLI is designed for AI agents first. The `--agent` flag sets `--json --compact --no-input --no-color` in one shot.

## In Claude Code

Reference the CLI in CLAUDE.md:

```markdown
**Keen:** `~/repos/keen/keen` — 1,091-endpoint CLI for all Atlassian APIs.
Use this instead of raw curl for ALL Atlassian API calls.
Auth: `export JIRA_CLOUD_PLATFORM_USERNAME=... JIRA_CLOUD_PLATFORM_PASSWORD=...`
```

## MCP Server

The `keen-mcp` binary exposes all 1,091 tools as MCP tools. Add to your MCP config:

```json
{
  "mcpServers": {
    "keen": {
      "command": "/path/to/keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "you@example.com",
        "JIRA_CLOUD_PLATFORM_PASSWORD": "your-token"
      }
    }
  }
}
```

## Command Discovery

Agents should use `which` to find commands:

```bash
keen which "create sprint"
keen which "search page"
keen which "addon version"
```

## Output Parsing

Agent mode returns JSON envelopes:

```json
{
  "action": "post",
  "data": { ... },
  "path": "/rest/api/3/...",
  "resource": "...",
  "status": 200,
  "success": true
}
```

The actual API response is in `data`. Use `--compact` to minimize token usage. Use `--select "key,summary,status"` to pick specific fields.
