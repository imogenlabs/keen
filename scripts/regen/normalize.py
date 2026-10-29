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
          f"fixed example placeholders in {len(fix_placeholders(sys.argv[1]))} files")
