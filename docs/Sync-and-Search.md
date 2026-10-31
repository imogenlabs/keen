# Sync and Search

keen can copy your Atlassian data to a local SQLite database and search it with full-text search. Once it is synced, searches, reports and checks run on your machine in well under a second and make no calls to Atlassian, so an AI agent can ask hundreds of questions without running into Atlassian's rate limits.

## Sync

```bash
keen sync                          # everything keen knows how to sync
keen sync --latest-only            # just refresh the newest page of each
```

Sync is resumable: run it again and it picks up where it stopped. Data lives in `~/.local/share/keen/data.db`.

## Search

```bash
keen search "payment failed" --data-source local
keen search "sprint bug" --type issue --data-source local
```

`--data-source` decides where an answer comes from:

| Value | Answers from |
|-------|--------------|
| `auto` (default) | Atlassian, falling back to the synced copy when Atlassian cannot be reached |
| `live` | Atlassian only |
| `local` | The synced copy only, with no network calls |

## Reports on synced data

```bash
keen analytics --type issue                        # how many issues are synced
keen stale --days 30                               # items nobody has touched in 30 days
keen orphans                                       # items missing important fields
keen load                                          # how much each person has assigned
```

## Fewer calls to Atlassian

keen keeps Atlassian traffic low in three ways:

- **Synced data.** Searches and reports on the local copy never call Atlassian.
- **A short response cache.** A read repeated within 5 minutes, from the command line or from the MCP server, is answered from keen's cache. Any successful write clears it, so a read after a write is always fresh. Use `--no-cache` on the command line, or `KEEN_MCP_CACHE=off` for the MCP server, to read live every time.
- **Backing off.** When Atlassian answers "too many requests", keen slows down and retries instead of failing.

See Atlassian's [rate limiting guide](https://developer.atlassian.com/cloud/jira/platform/rate-limiting/) for how the limits work.
