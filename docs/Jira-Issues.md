# Jira Issues

## Search (JQL)

```bash
# Basic search
keen jira-cloud-platform-search and-reconsile-issues-using-jql-post \
  --jql "project = PROJ AND status = 'In Progress'" --agent

# With specific fields
keen jira-cloud-platform-search and-reconsile-issues-using-jql-post \
  --jql "project = PROJ ORDER BY rank" \
  --fields '["key","summary","status","priority","customfield_10016"]' \
  --max-results 50 --agent

# All issues in a sprint
keen jira-cloud-platform-search and-reconsile-issues-using-jql-post \
  --jql "sprint = 7" --agent
```

## Get Issue

```bash
keen issue get PROJ-42 --agent
```

## Edit Issue

```bash
# Change priority
keen issue edit PROJ-42 --fields '{"priority":{"name":"High"}}' --agent

# Set story points (the field id varies by site; customfield_10016 is common)
keen issue edit PROJ-42 --fields '{"customfield_10016":3}' --agent

# Set resolution
keen issue edit PROJ-42 --fields '{"resolution":{"name":"Done"}}' --agent
```

### Descriptions

Jira wants `description` and `environment` as Atlassian Document Format. Send plain
text and keen wraps it for you, one paragraph per line. The text is kept exactly as
written: `## User Story` stays the literal characters, not a heading. This works for
`issue create`, `issue edit` and `issue create-bulk`. A value that is already an ADF
object is sent unchanged.

```bash
keen issue edit PROJ-42 --stdin --agent <<< '{"fields":{"description":"First line\nSecond line"}}'
```

## Transitions

```bash
# Get available transitions
keen issue transitions get PROJ-42 --agent

# Transition to a status (by transition ID)
keen issue transitions do PROJ-42 --transition-id 41 --agent
```

Transition IDs differ from one workflow to the next. List them with
`keen issue transitions get`, or pass the name and let keen look the ID up:

```bash
keen transition PROJ-42 --transition-name "Done" --agent
```

## Comments

```bash
# Add a comment. Plain text is wrapped into ADF, one paragraph per line
keen issue comment add PROJ-42 --body "Deployed to prod."

# Or send ADF yourself
keen issue comment add PROJ-42 --stdin <<< '{
  "body": {
    "type": "doc",
    "version": 1,
    "content": [{"type":"paragraph","content":[{"type":"text","text":"Deployed to prod."}]}]
  }
}'
```

## Bulk Operations

```bash
# Create up to 50 issues in one call; descriptions can be plain text
keen issue create-bulk --stdin --agent <<< '{"issueUpdates":[
  {"fields":{"project":{"key":"PROJ"},"issuetype":{"name":"Task"},"summary":"First","description":"Plain text is fine"}},
  {"fields":{"project":{"key":"PROJ"},"issuetype":{"name":"Task"},"summary":"Second"}}
]}'

# Fetch up to 100 issues by key or id
keen issue bulk-fetch --stdin <<< '{"issueIdsOrKeys":["PROJ-42","PROJ-43"]}'
```

Bulk edit, move, status changes, delete and watch are under `keen bulk`
(`submit-edit`, `submit-move`, `submit-transition`, `submit-delete`, `submit-watch`,
`submit-unwatch`). They run as a Jira job; add `--wait` and keen waits for it to finish.
