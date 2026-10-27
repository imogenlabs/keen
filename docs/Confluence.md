# Confluence

## Search

```bash
# Search by CQL
keen wiki search-content-by-cql --cql "type=page AND space=DOCS" --limit 10 --agent

# Full-text search
keen wiki search-by-cql --cql "text ~ 'deployment'" --agent
```

## Pages

```bash
# Child pages of a page
keen wiki get-content-descendants PAGE_ID --agent

# Copy a page
keen wiki copy-page PAGE_ID --stdin <<< '{"copyAttachments":true,"copyPermissions":true,"destination":{"type":"parent_page","value":"67890"}}'
```

Creating, reading and updating pages and blog posts are not in keen yet.

## Labels

```bash
# Add labels
keen wiki add-labels-to-content PAGE_ID --stdin <<< '[{"prefix":"global","name":"sprint-review"}]'

# Get labels for space
keen wiki get-labels-for-space DOCS --agent
```
