# CLAUDE.md — keen

`keen` is a generated Atlassian CLI (Jira, Confluence, Agile, Marketplace, JSM, Admin — ~1,091 endpoints), produced by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press). This file is for working **on** the keen codebase. For using keen as an agent (commands, discovery, `--agent` mode), see `AGENTS.md`.

## Public repo — do not leak

This repository is **public**. Never commit secrets, credentials, API tokens, or the Imogen Labs operational canon (cloud IDs, account IDs, the repo→Jira map, tenant facts). Those live in private repos and in Bitwarden. A canon block was once leaked here and had to be force-removed — do not reintroduce it.

## Build, test, run

```bash
go build ./...                 # compile everything
go build -o bin/keen ./cmd/keen        # the CLI
go build -o bin/keen-mcp ./cmd/keen-mcp # the MCP server
go test ./...                  # tests
bin/keen doctor --agent        # verify auth + connectivity
```

Config lives at `~/.config/keen/auth.json`:
```json
{ "base_url": "https://your-site.atlassian.net", "username": "you@example.com", "password": "<API_TOKEN>" }
```
`newClient()` fails loudly if `base_url` is empty or the `your-domain.atlassian.net` placeholder.

## Generated vs hand-written

- Most files under `internal/cli/` are **generated** ("DO NOT EDIT" header) by Printing Press. Fix systemic issues **upstream** in the generator, not in the generated tree.
- Narrow, intentional local edits are allowed but must be recorded in `.printing-press-patches.json` with a reason, so a regenerate can re-apply them.
- Hand-written customization (top-level aliases like `keen transition`/`find-issues`/`get`, and name→id resolution) lives in `internal/cli/aliases.go`.
- The SQLite local store/sync lives in `internal/store/`; the HTTP client in `internal/client/`.

## Regenerating

The generate command is in `README.md`. Use `--name keen` (not the old `jira-pp-cli`) so the binary keeps the `keen` identity. The published spec source on the Printing Press library is still named `jira-pp-cli`.

## Conventions

- One ticket = one branch = one PR (Jira project `INFRA`). Commit `KEY: imperative subject`.
- Do not add `Co-Authored-By` or any agent/AI attribution to commits or PRs.
