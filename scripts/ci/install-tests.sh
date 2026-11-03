#!/usr/bin/env bash
# Installs the snapshot's .deb and .rpm the way users do, and asks the image's
# MCP server for its tools. Run after `goreleaser release --snapshot` (CI).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

docker run --rm -v "$PWD/dist:/dist:ro" debian:stable bash -c \
  'apt-get update -qq && apt-get install -y -qq /dist/keen_linux_amd64.deb >/dev/null && keen version && command -v keen-mcp'
docker run --rm -v "$PWD/dist:/dist:ro" fedora:latest bash -c \
  'dnf install -y -q /dist/keen_linux_amd64.rpm >/dev/null && keen version && command -v keen-mcp'

img=$(docker image ls ghcr.io/imogenlabs/keen --format '{{.Repository}}:{{.Tag}}' | grep -- '-amd64$' | grep -v latest | head -1)
docker run --rm "$img" keen version
# The image's default command is keen-mcp: send it the MCP handshake and a tools/list.
reply=$( (printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"ci","version":"0"}}}' \
  '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; sleep 5) | docker run --rm -i "$img")
grep '"id":2' <<<"$reply" | grep -q '"name":"keen_run"' || { echo "keen-mcp in $img did not list keen_run" >&2; exit 1; }
echo "image $img: keen-mcp lists keen_search and keen_run"
