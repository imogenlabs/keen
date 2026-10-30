#!/usr/bin/env python3
"""Resolve conflict hunks in one file by taking one side.

Usage: resolve.py ours|theirs <file>...

Use it for files where one side is right for every hunk: `theirs` for pure
generator tables (which.go, sync.go maps, internal/mcp/tools.go), `ours` where
keen's hand fix must win (search.go). Mixed files are resolved by hand.
"""
import re
import sys

HUNK = re.compile(r"^<<<<<<< [^\n]*\n(.*?)^=======\n(.*?)^>>>>>>> [^\n]*\n", re.S | re.M)


def resolve(path, side):
    text = open(path, encoding="utf-8").read()
    new, n = HUNK.subn(lambda m: m.group(1) if side == "ours" else m.group(2), text)
    if "<<<<<<< " in new or "\n>>>>>>> " in new:
        raise SystemExit(f"{path}: unbalanced conflict markers")
    open(path, "w", encoding="utf-8").write(new)
    return n


if __name__ == "__main__":
    if len(sys.argv) < 3 or sys.argv[1] not in ("ours", "theirs"):
        sys.exit(__doc__)
    for p in sys.argv[2:]:
        print(f"{p}: {resolve(p, sys.argv[1])} hunks -> {sys.argv[1]}")
