#!/usr/bin/env python3
"""Package keen-mcp for the official MCP registry.

    scripts/mcp-registry.py <keen-mcp binary> <version> <out dir>

Writes <out dir>/keen-mcp.mcpb, an MCP bundle (manifest.json plus the binary;
Claude Desktop installs it in one click), and <out dir>/server.json, the
registry entry pointing at that bundle on the GitHub release with its
checksum. The release workflow builds the binary as a universal Mac binary,
runs this, uploads the bundle, then publishes server.json with mcp-publisher.
"""
import hashlib
import json
import os
import pathlib
import sys
import zipfile

REPO = "imogenlabs/keen"
NAME = "io.github.imogenlabs/keen"
DESCRIPTION = "Atlassian Cloud for AI apps: Jira, Confluence, Bitbucket, Service Management, Admin and more."
SITE = "JIRA_CLOUD_PLATFORM_URL"
USER = "JIRA_CLOUD_PLATFORM_USERNAME"
TOKEN = "JIRA_CLOUD_PLATFORM_PASSWORD"


def manifest(version):
    return {
        "manifest_version": "0.3",
        "name": "keen",
        "display_name": "keen for Atlassian",
        "version": version,
        "description": DESCRIPTION,
        "long_description": "Jira, Jira Software, Service Management, Confluence, Assets, Compass, Bitbucket, "
                            "Marketplace and Atlassian Admin, plus the org chart. Ask for a task in plain words "
                            "with keen_search, then run the tool it finds with keen_run.",
        "author": {"name": "Imogen Labs", "url": f"https://github.com/{REPO}"},
        "homepage": f"https://github.com/{REPO}",
        "repository": {"type": "git", "url": f"https://github.com/{REPO}"},
        "license": "Apache-2.0",
        "server": {
            "type": "binary",
            "entry_point": "server/keen-mcp",
            "mcp_config": {
                "command": "${__dirname}/server/keen-mcp",
                "args": [],
                "env": {SITE: "${user_config.site}", USER: "${user_config.email}", TOKEN: "${user_config.token}"},
            },
        },
        "tools": [
            {"name": "keen_search", "description": "Find keen tools for a task, with their parameters."},
            {"name": "keen_run", "description": "Run a keen tool found with keen_search."},
        ],
        "user_config": {
            "site": {"type": "string", "title": "Atlassian site", "required": True,
                     "description": "Your site's address, e.g. https://yourteam.atlassian.net"},
            "email": {"type": "string", "title": "Email", "required": True,
                      "description": "The email you sign in to Atlassian with"},
            "token": {"type": "string", "title": "API token", "required": True, "sensitive": True,
                      "description": "Create one at https://id.atlassian.com/manage-profile/security/api-tokens"},
        },
        "compatibility": {"platforms": ["darwin"]},
    }


def server_json(version, sha256):
    return {
        "$schema": "https://static.modelcontextprotocol.io/schemas/2025-12-11/server.schema.json",
        "name": NAME,
        "title": "keen",
        "description": DESCRIPTION,
        "version": version,
        "repository": {"url": f"https://github.com/{REPO}", "source": "github"},
        "websiteUrl": f"https://github.com/{REPO}",
        "packages": [{
            "registryType": "mcpb",
            "identifier": f"https://github.com/{REPO}/releases/download/v{version}/keen-mcp.mcpb",
            "fileSha256": sha256,
            "transport": {"type": "stdio"},
            "environmentVariables": [
                {"name": SITE, "description": "Your Atlassian site, e.g. https://yourteam.atlassian.net", "isRequired": True},
                {"name": USER, "description": "The email you sign in to Atlassian with", "isRequired": True},
                {"name": TOKEN, "description": "An Atlassian API token", "isRequired": True, "isSecret": True},
            ],
        }],
    }


def main():
    if len(sys.argv) != 4:
        sys.exit(__doc__)
    binary, version, out = pathlib.Path(sys.argv[1]), sys.argv[2].removeprefix("v"), pathlib.Path(sys.argv[3])
    out.mkdir(parents=True, exist_ok=True)
    bundle = out / "keen-mcp.mcpb"
    with zipfile.ZipFile(bundle, "w", zipfile.ZIP_DEFLATED) as z:
        z.writestr("manifest.json", json.dumps(manifest(version), indent=2) + "\n")
        info = zipfile.ZipInfo("server/keen-mcp")
        info.external_attr = 0o755 << 16
        info.compress_type = zipfile.ZIP_DEFLATED
        z.writestr(info, binary.read_bytes())
    sha = hashlib.sha256(bundle.read_bytes()).hexdigest()
    (out / "server.json").write_text(json.dumps(server_json(version, sha), indent=2) + "\n")
    print(f"{bundle} sha256 {sha}")


if __name__ == "__main__":
    main()
