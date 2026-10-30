# Agent Mode

The CLI is designed for AI agents first. The `--agent` flag sets `--json --compact --no-input --no-color` in one shot.

## In Claude Code

Reference the CLI in CLAUDE.md:

```markdown
**Keen:** `keen` — 1,843-endpoint CLI for all Atlassian APIs.
Use this instead of raw curl for ALL Atlassian API calls.
Auth: `~/.config/keen/auth.json`, or `export JIRA_CLOUD_PLATFORM_USERNAME=... JIRA_CLOUD_PLATFORM_PASSWORD=...`
```

## MCP Server

The `keen-mcp` binary gives agents every keen command (1,939 tools) through two MCP tools: `keen_search` (describe the task, get matching tools and their parameters) and `keen_run` (run one by name). The full list would not fit in most models' context; set `KEEN_MCP_TOOLS=all` if you want it anyway. `keen_search` ranks tools for APIs Atlassian has deprecated below the alternatives and leaves out ones it has switched off. Add to your MCP config:

```json
{
  "mcpServers": {
    "keen": {
      "command": "/path/to/keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "you@example.com",
        "JIRA_CLOUD_PLATFORM_PASSWORD": "your-token",
        "JIRA_CLOUD_PLATFORM_URL": "https://your-site.atlassian.net"
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

The actual API response is in `data`.

## Keeping Output Small

Every byte keen prints lands in the agent's context. Three ways to cut it:

```bash
# Keep only the fields you need (dotted paths reach nested ones)
keen get PROJ-42 --agent --select key,fields.summary,fields.status.name

# Large results go to a file; stdout gets a one-line summary
keen find-issues --jql "project = PROJ" --agent --output-summary
# {"output_summary":{"bytes":41309,"file":"/tmp/keen-output-1234.json","keys":["action","data",...],"records":17}}

# Set your own limit in bytes (default 4096, about 1,000 tokens)
keen find-issues --jql "project = PROJ" --agent --output-summary=20000
```

`--output-summary` prints small results as usual. The file is created in your temp directory (`TMPDIR`), readable only by you; delete it when you're done.

`--select` checks the fields you ask for. A likely typo fails with exit code 2 and a suggestion (`Did you mean "fields.summary"?`), instead of printing the whole response. A field Jira simply left out of this record, such as `fields.parent` on an issue without a parent, prints a warning and is skipped.

`--compact` (part of `--agent`) keeps only the usual key fields when you don't pass `--select`.
