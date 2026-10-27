# Installation

## Binary Download

Prebuilt binaries for macOS (Apple Silicon) and Linux (x86_64, arm64) are on
[GitHub Releases](https://github.com/imogenlabs/keen/releases/latest). Each
archive holds both `keen` and `keen-mcp`:

```bash
curl -L https://github.com/imogenlabs/keen/releases/latest/download/keen_darwin_arm64.tar.gz | tar xz
sudo mv keen_darwin_arm64/keen keen_darwin_arm64/keen-mcp /usr/local/bin/
```

Swap `darwin_arm64` for `linux_amd64` or `linux_arm64` on Linux.

## From Source

Requires Go 1.26.10 or newer and Git.

```bash
git clone https://github.com/imogenlabs/keen.git
cd keen
go build -o keen ./cmd/keen/
```

Move the binary to your PATH:

```bash
# macOS/Linux
sudo mv keen /usr/local/bin/

# Or symlink
ln -s $(pwd)/keen ~/go/bin/keen
```

## MCP Server

```bash
go build -o keen-mcp ./cmd/keen-mcp/
```

The MCP server exposes every keen command (1,114 tools) over the Model Context Protocol (stdio transport). Add to your MCP config:

```json
{
  "mcpServers": {
    "keen": {
      "command": "/path/to/keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "you@example.com",
        "JIRA_CLOUD_PLATFORM_PASSWORD": "your-api-token",
        "JIRA_CLOUD_PLATFORM_URL": "https://your-site.atlassian.net"
      }
    }
  }
}
```

## Verify

```bash
keen doctor
keen --version
```
