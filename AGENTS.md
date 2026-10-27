# keen — agent guide

How to use keen from an agent. For install and credentials, see `README.md` and `docs/`.

## Start here

Ask keen for the current truth instead of trusting a copied command list:

```bash
keen doctor --agent           # credentials and connectivity
keen agent-context --pretty   # what this build can do
keen which "<capability>"     # find the command for a task
keen <command> --help
```

Add `--agent` to every call: JSON, compact output, no prompts, no color.

The everyday shortcuts:

```bash
keen get PROJ-42 --agent
keen find-issues --jql "project = PROJ AND status != Done" --agent
keen transition PROJ-42 --transition-name "In Progress" --agent
keen issue comment add PROJ-42 --body "Plain text is fine" --agent
keen wiki get-page PAGE_ID --agent
```

## Before you change anything remote

Read the command's help and dry-run it first:

```bash
keen <command> --dry-run --agent
```

Use `--yes --no-input` only once the target and side effects are clear.

## Changing keen itself

Most of `internal/cli/` and `internal/mcp/` is generated (files marked `DO NOT EDIT`). keen is maintained here: fix it in this repository, keep edits to generated files narrow, and record each one in `.printing-press-patches.json` (id, summary, reason, files, validated outcome) so a future regeneration can re-apply it. Prefer a new hand-written file over editing a generated one. Do not open issues or pull requests on the generator's repository.
