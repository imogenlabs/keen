# Installation

## Prerequisites

- Go 1.26.3 or newer
- Git

## From Source

```bash
git clone https://github.com/smashingtags/keen.git
cd keen
go build -o atlas-cli ./cmd/keen/
```

Move the binary to your PATH:

```bash
# macOS/Linux
sudo mv atlas-cli /usr/local/bin/

# Or symlink
ln -s $(pwd)/atlas-cli ~/go/bin/atlas-cli
```

## MCP Server

```bash
go build -o atlas-mcp ./cmd/keen-mcp/
```

The MCP server exposes all 1,091 tools over the Model Context Protocol (stdio transport). Add to your MCP config:

```json
{
  "mcpServers": {
    "atlas": {
      "command": "/path/to/atlas-mcp",
      "args": []
    }
  }
}
```

## Verify

```bash
atlas-cli doctor
atlas-cli --version
```

## Batch Install (Printing Press Library CLIs)

To install all 32 CLIs on a machine:

```bash
# macOS/Linux
for cli in docker-hub posthog supabase firecrawl youtube spotify hackernews substack wikipedia slack notion chrome-history figma zoom sendgrid resend stripe mercury google-search-console google-ads producthunt beehiiv ahrefs digitalocean x-twitter pushover shopify gumroad open-meteo weather-goat icloud; do
  go install "github.com/mvanhorn/printing-press-library/library/*/$(echo $cli)/cmd/${cli}-pp-cli@latest" 2>/dev/null || echo "  $cli: skipped"
done
```
