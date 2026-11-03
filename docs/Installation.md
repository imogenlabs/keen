# Installation

## Homebrew

On macOS or Linux with [Homebrew](https://brew.sh):

```bash
brew install imogenlabs/tap/keen
```

This installs both `keen` and `keen-mcp`. `brew upgrade keen` picks up new releases.

## Debian, Ubuntu, Fedora, RHEL

Each release has `.deb` and `.rpm` packages for x86_64 (`amd64`) and arm64. They install `keen` and `keen-mcp` in `/usr/bin`:

```bash
# Debian, Ubuntu
curl -LO https://github.com/imogenlabs/keen/releases/latest/download/keen_linux_amd64.deb
sudo apt install ./keen_linux_amd64.deb

# Fedora, RHEL
curl -LO https://github.com/imogenlabs/keen/releases/latest/download/keen_linux_amd64.rpm
sudo dnf install ./keen_linux_amd64.rpm
```

On arm64, use `keen_linux_arm64.deb` or `keen_linux_arm64.rpm`.

## Docker

`ghcr.io/imogenlabs/keen` (also `smashingtags/keen` on Docker Hub), for amd64 and arm64. It runs `keen-mcp` by default; put `keen` and its arguments after the image name to run the CLI:

```bash
docker run --rm \
  -e JIRA_CLOUD_PLATFORM_URL=https://yoursite.atlassian.net \
  -e JIRA_CLOUD_PLATFORM_USERNAME=you@example.com \
  -e JIRA_CLOUD_PLATFORM_PASSWORD=your-api-token \
  ghcr.io/imogenlabs/keen keen get PROJ-42 --agent
```

Add `-v keen-data:/home/nonroot` to keep `keen sync`'s local copy between runs. For an MCP host, use `docker` as the command with `run -i --rm -e JIRA_CLOUD_PLATFORM_URL -e JIRA_CLOUD_PLATFORM_USERNAME -e JIRA_CLOUD_PLATFORM_PASSWORD ghcr.io/imogenlabs/keen` as its arguments, and set those three variables in the host's `env` block.

## Binary Download

Prebuilt binaries for macOS (Apple Silicon and Intel) and Linux (x86_64, arm64) are on
[GitHub Releases](https://github.com/imogenlabs/keen/releases/latest). Each
archive holds both `keen` and `keen-mcp`:

```bash
curl -L https://github.com/imogenlabs/keen/releases/latest/download/keen_darwin_arm64.tar.gz | tar xz
sudo mv keen_darwin_arm64/keen keen_darwin_arm64/keen-mcp /usr/local/bin/
```

Swap `darwin_arm64` for `darwin_amd64` on an Intel Mac, or `linux_amd64` or `linux_arm64` on Linux.

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

The MCP server gives agents every keen command (1,939 tools) over the Model Context Protocol (stdio transport), through two tools: `keen_search` and `keen_run`. Set `KEEN_MCP_TOOLS=all` to list every tool instead. Add to your MCP config:

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
