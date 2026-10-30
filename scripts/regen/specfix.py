#!/usr/bin/env python3
"""Make Atlassian's server URLs usable by the generator before it runs.

Two spec shapes need it:

  site template: Confluence v2 gives its server as https://{your-domain}/wiki/api/v2.
  The generator cannot fill {your-domain} and writes https://no-default/... into
  every command. Point it at your-domain.atlassian.net, the server the Jira spec
  uses, and move /wiki/api/v2 onto each path, so the commands become paths on the
  configured base_url like every other site command.

  ids in the server path: Assets gives https://api.atlassian.com/ex/jira/{cloudId}/
  jsm/assets/workspace/{workspaceId}/v1. Move that path onto every operation and
  declare each {id} as a path parameter, so the ids become command arguments.

The downloaded copies in specs/ stay as Atlassian publishes them; this writes
the fixed copies the generator reads.

Usage: specfix.py <in.json> <out.json>
"""
import json
import re
import sys
from urllib.parse import urlsplit

VAR = re.compile(r"\{([^}]+)\}")
METHODS = ("get", "put", "post", "delete", "patch", "head", "options")


def fix(spec):
    servers = spec.get("servers") or []
    if not servers:
        return spec
    url = servers[0]["url"]
    host_part, _, path_part = url.partition("//")[2].partition("/")
    path_part = "/" + path_part.rstrip("/") if path_part else ""
    if not VAR.search(url):
        return spec  # nothing to fill; leave the spec exactly as published

    if VAR.search(host_part):
        host_part = "your-domain.atlassian.net"

    # The server path moves onto every operation, so the server is a bare host.
    ids = VAR.findall(path_part)
    spec["paths"] = {path_part + p: item for p, item in spec["paths"].items()}
    if ids:
        for item in spec["paths"].values():
            for method, op in item.items():
                if method not in METHODS:
                    continue
                params = op.setdefault("parameters", [])
                have = {p.get("name") for p in params if p.get("in") == "path"}
                missing = [n for n in ids if n not in have]
                params[:0] = [{"name": n, "in": "path", "required": True, "schema": {"type": "string"}}
                              for n in missing]

    scheme = urlsplit(url).scheme or "https"
    spec["servers"] = [{"url": f"{scheme}://{host_part}"}]
    return spec


if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit(__doc__)
    with open(sys.argv[1]) as f:
        spec = json.load(f)
    with open(sys.argv[2], "w") as f:
        json.dump(fix(spec), f)
