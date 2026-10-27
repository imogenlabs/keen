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
# Create a page (body is Confluence storage format; --representation wiki for wiki markup)
keen wiki create-page --space DOCS --title "Release notes" --parent-id 12345 --stdin < notes.xhtml

# Read a page: title, version, URL and storage body
keen wiki get-page PAGE_ID --agent

# Replace a page's body; keen reads the current version and sends the next one
keen wiki update-page PAGE_ID --stdin --message "Add v1.2.0" < notes.xhtml

# Child pages of a page
keen wiki get-content-descendants PAGE_ID --agent

# Copy a page
keen wiki copy-page PAGE_ID --stdin <<< '{"copyAttachments":true,"copyPermissions":true,"destination":{"type":"parent_page","value":"67890"}}'
```

Markdown is not converted. Blog posts are not covered by these commands.

## Labels

```bash
# Add labels
keen wiki add-labels-to-content PAGE_ID --stdin <<< '[{"prefix":"global","name":"sprint-review"}]'

# Get labels for space
keen wiki get-labels-for-space DOCS --agent
```
