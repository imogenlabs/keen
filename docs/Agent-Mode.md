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

The `keen-mcp` binary gives agents every keen command (1,939 tools) through two MCP tools: `keen_search` (describe the task, get matching tools and their parameters) and `keen_run` (run one by name). The full list would not fit in most models' context; set `KEEN_MCP_TOOLS=all` if you want it anyway. `keen_search` ranks tools for APIs Atlassian has deprecated below the alternatives and leaves out ones it has switched off. A read repeated within 5 minutes is answered from keen's cache, so agents stay under Atlassian's rate limits; any write clears it, and `KEEN_MCP_CACHE=off` makes every read live. Each release is tested this way before it ships: `scripts/mcp-check.py` searches for everyday tasks, runs the tools and checks the answers against the CLI, and a Claude Code session with only keen-mcp answers real questions. Add to your MCP config:

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

## Claude Desktop, and the MCP registry

On a Mac, download `keen-mcp.mcpb` from the [latest release](https://github.com/imogenlabs/keen/releases/latest) and open it: Claude Desktop installs keen and asks for your Atlassian site, email and API token. keen is also listed in the official [MCP registry](https://registry.modelcontextprotocol.io) as `io.github.imogenlabs/keen`, updated with every release.

## Command Discovery

Agents should use `which` to find commands. It ranks the same way as keen-mcp's `keen_search`, covers keen's own commands (`find-issues`, `transition`, `sprints`, `stale`), and answers "no match" (exit 2) rather than guessing:

```bash
keen which "transition an issue"
keen which "list sprints on a board"
keen which "create page"
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

## Starting on an Issue

`keen context PROJ-42` gives an agent what it needs to start on an issue in one call: the issue,
its parent, comments, linked issues, sub-tasks, GitHub pull requests and commits (as `keen github
links` finds them) and linked Confluence pages.

```bash
keen context PROJ-42 --org my-org --agent
keen context PROJ-42 --budget 4000 --agent --select issue,comments,budget
```

The bundle stays under `--budget` tokens (default 8000, about four characters each). The issue
always goes in; its description gets at most half the budget. Then come the parent, GitHub work,
linked issues and sub-tasks, the three newest comments, Confluence pages and older comments.
Anything that does not fit is listed under `budget.omitted` with its size and the reason. A part
that cannot be read (no GitHub token, a page you cannot see) is a warning; the rest still comes back.

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

`--select` checks the fields you ask for. A likely typo fails with exit code 2 and a suggestion (`Did you mean "fields.summary"?`), instead of printing the whole response. A field Jira simply left out of this record, such as `fields.parent` on an issue without a parent, prints a warning and is skipped. Paths start inside the response, but a path copied from the output works too: `data.id` after a write is the same as `id`, and `results.key` after a read is the same as `key`. After a write that went through, a field that is not there never fails the command: keen warns, prints the whole response and exits 0, so you can see the write happened and do not make it again.

`--compact` (part of `--agent`) keeps only the usual key fields when you don't pass `--select`.
