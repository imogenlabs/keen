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
]

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
    print(f"normalized {normalize(sys.argv[1])} files")
