#!/usr/bin/env python3
"""Rename CLI Printing Press output from its default names to keen's.

The generator is run with --name jira-pp-cli, so its output uses the same
names as keen's first commit (module jira-pp-cli-pp-cli). This applies the
same rename keen made by hand, so a three-way merge only sees real changes.

Usage: normalize.py <generated-dir>   (edits the directory in place)
"""
import os
import re
import sys

from placeholders import fix_placeholders

DIRS = {
    "cmd/jira-pp-cli-pp-cli": "cmd/keen",
    "cmd/jira-pp-cli-pp-mcp": "cmd/keen-mcp",
}

# Order matters: the long names first, then the bare API name.
SUBS = [
    (re.compile(r"jira-pp-cli-pp-cli"), "keen"),
    (re.compile(r"jira-pp-cli-pp-mcp"), "keen-mcp"),
    (re.compile(r"'%s-pp-cli\b"), "'%s"),
    (re.compile(r"jira-pp-cli"), "keen"),
    (re.compile(r"Jira Pp Cli"), "keen"),
]

# Async commands poll a status endpoint. The generator guesses it, differently on
# each run and often wrongly, and jobs.go only fills in {id} / {job_id}. Pin the
# endpoint Atlassian documents for each, so the output is correct and repeatable.
JOB_STATUS = {
    **{f"internal/cli/bulk_{op}.go": ("bulk", "get-operation-progress", "/rest/api/3/bulk/queue/{id}")
       for op in ("submit-delete", "submit-edit", "submit-move", "submit-transition",
                  "submit-watch", "submit-unwatch", "get-operation-progress")},
    "internal/cli/workflows_update.go": ("task", "get", "/rest/api/3/task/{id}"),
    "internal/cli/workflowscheme_read-workflow-schemes.go": ("task", "get", "/rest/api/3/task/{id}"),
    "internal/cli/wiki_async-convert-content-body-request.go":
        ("wiki", "async-convert-content-body-response", "/wiki/rest/api/contentbody/convert/async/{id}"),
    "internal/cli/wiki_get-and-async-convert-macro-body-by-macro-id.go":
        ("wiki", "async-convert-content-body-response", "/wiki/rest/api/contentbody/convert/async/{id}"),
}
STATUS_RESOURCE = re.compile(r'(StatusResource: )"[^"]*"')
STATUS_ENDPOINT = re.compile(r'(StatusEndpoint: )"[^"]*"')
WAIT_PATH = re.compile(r'(WaitForJob\(ctx, c, )"[^"]*"')


def pin_job_status(root):
    pinned = 0
    for rel, (res, ep, path) in JOB_STATUS.items():
        f = os.path.join(root, rel)
        if not os.path.exists(f):
            continue
        text = open(f, encoding="utf-8").read()
        new = STATUS_RESOURCE.sub(lambda m: f'{m.group(1)}"{res}"', text)
        new = STATUS_ENDPOINT.sub(lambda m: f'{m.group(1)}"{ep}"', new)
        new = WAIT_PATH.sub(lambda m: f'{m.group(1)}"{path}"', new)
        if new != text:
            open(f, "w", encoding="utf-8").write(new)
            pinned += 1
    return pinned


# Clear names for top-level command groups the generator names badly (a name
# clash, or a generic Assets resource). Only the names people type and read
# change: the command, help examples, `keen which`, MCP tool names, SKILL.md and
# the tool manifest. Store tables, sync resource names (keen sync --resources),
# files and Go identifiers keep the generator's name. A group is only ever renamed before it ships; after that, add the old
# name to internal/compat instead of changing an entry here.
GROUP_RENAMES = {
    "admin-control-admin": "admin-control",
    "manage": "user-management",
    "admin": "user-provisioning",
    "aql": "assets-aql",
    "iql": "assets-iql",
    "icon": "assets-icon",
    "object": "assets-object",
    "objectschema": "assets-objectschema",
    "objecttype": "assets-objecttype",
    "objecttypeattribute": "assets-objecttypeattribute",
    "objectconnectedtickets": "assets-object-tickets",
    "importsource": "assets-import-source",
    "api-config": "assets-config",
    "api-import": "assets-import",
    "global": "assets-global",
    "progress": "assets-progress",
    "usage": "assets-usage",
}


def group_patterns(old, new):
    o = re.escape(old)
    return [
        (re.compile(r'(\{Command: ")' + o + r'([ "])'), r"\g<1>" + new + r"\2"),
        (re.compile(r'(Group: ")' + o + r'(")'), r"\g<1>" + new + r"\2"),
        (re.compile(r'(\bkeen )' + o + r'([ "`\n])'), r"\g<1>" + new + r"\2"),
        (re.compile(r'(NewTool\(")' + o + r'([_"])'), r"\g<1>" + new + r"\2"),
        (re.compile(r'("name": +")' + o + r'(_)'), r"\g<1>" + new + r"\2"),
        (re.compile(r'(\*\*)' + o + r'(\*\* — Manage )[^\n]*'), r"\g<1>" + new + r"\2" + new.replace("-", " ")),
    ]


def rename_groups(root):
    """Apply GROUP_RENAMES; returns the number of files changed."""
    subs = [p for old, new in GROUP_RENAMES.items() for p in group_patterns(old, new)]
    changed = 0
    targets = [os.path.join(root, f) for f in ("SKILL.md", "tools-manifest.json")]
    for sub in ("internal/cli", "internal/mcp"):
        for dirpath, _, files in os.walk(os.path.join(root, sub)):
            targets += [os.path.join(dirpath, f) for f in files if f.endswith(".go")]
    for path in targets:
        if not os.path.exists(path):
            continue
        text = open(path, encoding="utf-8").read()
        new = text
        for pat, rep in subs:
            new = pat.sub(rep, new)
        # The group command itself: Use and Short in its own file.
        name = os.path.basename(path)
        for old, repl in GROUP_RENAMES.items():
            if name in (old + ".go", "promoted_" + old + ".go"):
                new = re.sub(r'(Use:\s+")' + re.escape(old) + '([ "])', r"\g<1>" + repl + r"\2", new)
                new = re.sub(r'(Short:\s+")Manage [^"]*"', r"\g<1>Manage " + repl.replace("-", " ") + '"', new)
        if new != text:
            open(path, "w", encoding="utf-8").write(new)
            changed += 1
    return changed


TEXT_EXT = {".go", ".mod", ".md", ".json", ".yaml", ".yml", ".toml", ".txt", ""}


def normalize(root):
    for old, new in DIRS.items():
        src, dst = os.path.join(root, old), os.path.join(root, new)
        if os.path.isdir(src):
            os.makedirs(os.path.dirname(dst), exist_ok=True)
            os.rename(src, dst)
    changed = 0
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d not in (".git", "specs")]
        for name in filenames:
            path = os.path.join(dirpath, name)
            if os.path.splitext(name)[1] not in TEXT_EXT or os.path.islink(path):
                continue
            try:
                text = open(path, encoding="utf-8").read()
            except (UnicodeDecodeError, OSError):
                continue
            new = text
            for pat, rep in SUBS:
                new = pat.sub(rep, new)
            if new != text:
                open(path, "w", encoding="utf-8").write(new)
                changed += 1
    return changed


if __name__ == "__main__":
    if len(sys.argv) != 2 or not os.path.isdir(sys.argv[1]):
        sys.exit(__doc__)
    print(f"normalized {normalize(sys.argv[1])} files, pinned {pin_job_status(sys.argv[1])} job status endpoints, "
          f"fixed example placeholders in {len(fix_placeholders(sys.argv[1]))} files, "
          f"renamed command groups in {rename_groups(sys.argv[1])} files")
