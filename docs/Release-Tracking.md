# Release Tracking

`keen release sync` records your git release tags in Jira: each tag becomes a Jira version, and the issues that release shipped get it as their Fix Version.

## Record every release so far

```bash
keen release sync --repo-path . --project PROJ --backfill --head origin/main --dry-run
```

`--backfill` walks every release tag. It also keeps an unreleased version called `next` for work merged since the last tag. Drop `--dry-run` to write.

## After each new release

```bash
keen release sync --repo-path . --project PROJ --tag v1.4.0 --head origin/main
```

The `next` version becomes the new tag, and a fresh `next` is created.

## What it writes

- A version named after the tag, released on the tag's commit date and linked to its GitHub release.
- On each issue whose key appears in that release's commits, and in no earlier release: the version as a Fix Version. Existing Fix Versions are added to, never replaced.

It is safe to run again: nothing is written twice.

## In GitHub Actions

Run it on every tag push, with your site URL, email and API token as repository secrets:

```yaml
on:
  push:
    tags: ["v*"]
jobs:
  record:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - run: curl -fsSL https://github.com/imogenlabs/keen/releases/latest/download/keen_linux_amd64.tar.gz | tar xz
      - run: keen_linux_amd64/keen release sync --repo-path . --project PROJ --tag "$GITHUB_REF_NAME" --head origin/main
        env:
          JIRA_CLOUD_PLATFORM_URL: ${{ secrets.JIRA_BASE_URL }}
          JIRA_CLOUD_PLATFORM_USERNAME: ${{ secrets.JIRA_USER }}
          JIRA_CLOUD_PLATFORM_PASSWORD: ${{ secrets.JIRA_API_TOKEN }}

## The GitHub work for one issue

`keen release sync` reads keys from commits, so keys belong in branch names, commit messages and
PR titles (`fix/PROJ-42-retry-uploads`, `fix: PROJ-42 retry failed uploads`). To see the pull
requests and commits for one issue:

```bash
keen github links PROJ-42 --org my-org
keen github links PROJ-42 --repo my-org/api --agent
```

It merges three sources, each PR and commit once: Jira's development panel (when the GitHub for
Jira app is connected), GitHub search for the key in PR titles, bodies and commit messages, and
web links on the issue that point at a PR. GitHub search uses `KEEN_GITHUB_TOKEN`, `GH_TOKEN`,
`GITHUB_TOKEN` or the `gh` CLI's login; without a token, or without `--org`/`--repo` outside a
checkout, that source is skipped with a warning and the others still answer.
