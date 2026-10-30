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

Markdown is not converted. For blog posts, use the `blogposts` commands below.

## Labels

```bash
# Add labels
keen wiki add-labels-to-content PAGE_ID --stdin <<< '[{"prefix":"global","name":"sprint-review"}]'

# Get labels for space
keen wiki get-labels-for-space DOCS --agent
```

## Confluence v2 API

Atlassian is moving Confluence to its v2 API. keen has every v2 operation as a command grouped by what it works on: `pages`, `blogposts`, `spaces`, `folders`, `databases`, `whiteboards`, `custom-content`, `attachments`, `footer-comments`, `inline-comments`, `labels`, `tasks` and more. Run `keen <group> --help` to see them. The `wiki` commands above keep working.

```bash
# List spaces, then the pages in one
keen spaces get --agent --select results.id,results.key,results.name
keen pages get --space-id 65986 --limit 10 --agent

# One page with its body
keen pages get-by-id 66038 --body-format storage --agent

# Blog posts and folders
keen blogposts get --space-id 65986 --agent
keen folders get-by-id 123456 --agent
```
