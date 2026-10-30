#!/usr/bin/env python3
"""Make Atlassian's server URLs usable by the generator before it runs.

Two spec shapes need it:

  site template: Confluence v2 gives its server as https://{your-domain}/wiki/api/v2.
  The generator cannot fill {your-domain} and writes https://no-default/... into
  every command. Point it at your-domain.atlassian.net, the server the Jira spec
  uses. The server path stays, so the commands are grouped by resource (pages,
  spaces, ...); postfix.py then makes those URLs relative to the configured site.

  ids in the server path: Assets gives https://api.atlassian.com/ex/jira/{cloudId}/
  jsm/assets/workspace/{workspaceId}/v1. An id named in TEMPLATE_VARS (config.sh)
  is replaced by a plain marker (KEEN_CLOUD_ID), which keen swaps for the real
  value at run time (internal/client/sitevars.go). The generator's own server
  variables are dropped when it merges several specs, so a marker it leaves
  alone is used instead. Any other id moves onto every operation as a path
  parameter.

Command names come from operation ids, so two id shapes are fixed too:

  missing ids: the user management spec has none, and the generator falls back
  to names like "create-users-2". Each missing id is made from the summary
  ("Deactivate a user" -> deactivateUser).

  a shared prefix: admin control names every operation "APIs_...", which put
  "apis-" in front of every command. A prefix ending in "_" that every id in
  the spec shares is dropped.

  auth schemes: keen chooses each API's credential itself (internal/client/
  hostauth.go). A spec whose auth schemes would make the generator rebuild keen's
  login around them (Bitbucket's OAuth2 did) is listed in STRIP_AUTH and has its
  security schemes and requirements removed (--no-auth).

The id rules run only for specs listed in FIX_IDS in config.sh (--ids): applied
to a spec keen already ships, they would rename its commands. Specs that need
none of this are left exactly as published. The downloaded
copies in specs/ stay as Atlassian publishes them; this writes the fixed copies
the generator reads.

Usage: specfix.py [--ids] [--no-auth] <in.json> <out.json>
"""
import json
import os
import re
import sys
from urllib.parse import urlsplit

VAR = re.compile(r"\{([^}]+)\}")
METHODS = ("get", "put", "post", "delete", "patch", "head", "options")


def operations_of(spec):
    for item in spec.get("paths", {}).values():
        for method, op in item.items():
            if method in METHODS:
                yield op


def singular(word):
    if word.endswith("ies"):
        return word[:-3] + "y"
    return word[:-1] if word.endswith("s") else word


def camel(words):
    return "".join(w[:1].upper() + w[1:] for w in words)


def fix_operation_ids(spec):
    ops = list(operations_of(spec))
    made = {}  # id made from a summary -> [(path, op)]
    for path, item in spec.get("paths", {}).items():
        for method, op in item.items():
            if method in METHODS and not op.get("operationId") and op.get("summary"):
                summary = re.sub(r"['\u2019]s\b", "", op["summary"])  # "project's" -> "project"
                words = [w for w in re.findall(r"[A-Za-z0-9]+", summary) if w.lower() not in ("a", "an", "the")]
                op["operationId"] = words[0].lower() + camel(w.lower() for w in words[1:])
                made.setdefault(op["operationId"], []).append((path, op))
    # Two operations with the same summary: tell them apart by the area of the
    # API they sit in, or else by the id the narrower one takes.
    for same in made.values():
        if len(same) < 2:
            continue
        areas = [p.strip("/").split("/")[0] for p, _ in same]
        if len(set(areas)) == len(areas):
            for (p, op), area in zip(same, areas):
                op["operationId"] += "In" + camel(re.findall(r"[A-Za-z0-9]+", singular(area)))
            continue
        for p, op in same:
            params = re.findall(r"\{([^}]+)\}", p)
            if len(params) == max(len(re.findall(r"\{([^}]+)\}", q)) for q, _ in same) and len({len(re.findall(r"\{", q)) for q, _ in same}) > 1:
                op["operationId"] += "By" + camel(re.findall(r"[A-Za-z0-9]+", params[-1]))
    ids = [op.get("operationId", "") for op in ops]
    prefix = re.match(r"[A-Za-z0-9]+_", ids[0]) if ids else None
    if prefix and all(i.startswith(prefix.group()) for i in ids):
        for op in ops:
            op["operationId"] = op["operationId"][len(prefix.group()):]


def template_vars():
    """TEMPLATE_VARS from the environment: "cloudId=KEEN_CLOUD_ID workspaceId=..."."""
    pairs = (p.split("=", 1) for p in os.environ.get("TEMPLATE_VARS", "").split() if "=" in p)
    return dict(pairs)


def strip_auth(spec):
    spec.pop("security", None)
    spec.get("components", {}).pop("securitySchemes", None)
    for op in operations_of(spec):
        op.pop("security", None)


def fix(spec, ids=False, no_auth=False):
    if ids:
        fix_operation_ids(spec)
    if no_auth:
        strip_auth(spec)
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

    # Ids keen fills at run time become markers in the server URL; the rest
    # move onto every operation as path parameters.
    known = template_vars()
    for v in VAR.findall(path_part):
        if v in known:
            path_part = path_part.replace("{" + v + "}", known[v])
    ids = VAR.findall(path_part)
    if ids:
        spec["paths"] = {path_part + p: item for p, item in spec["paths"].items()}
        path_part = ""
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
    spec["servers"] = [{"url": f"{scheme}://{host_part}{path_part}"}]
    return spec


if __name__ == "__main__":
    args = sys.argv[1:]
    ids = "--ids" in args
    no_auth = "--no-auth" in args
    args = [a for a in args if a not in ("--ids", "--no-auth")]
    if len(args) != 2:
        sys.exit(__doc__)
    with open(args[0]) as f:
        spec = json.load(f)
    with open(args[1], "w") as f:
        json.dump(fix(spec, ids, no_auth), f)
