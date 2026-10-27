#!/usr/bin/env python3
"""Re-apply keen's standing changes to code the generator adds fresh.

The merge keeps keen's edits to code that already existed. Code that is new in
this generator run (a new table, a new test) arrives in the generator's form,
so the same rule has to be applied to it here.

  nullable-fk-tests: keen makes typed-table foreign-key columns nullable (May
  2026), so a typed insert without the parent id succeeds. The generator's
  tests for new tables still expect it to fail; rewrite them to expect success.

  nullable-fk-columns: the same rule for the tables themselves. A new typed
  table arrives with "<parent>_id" TEXT NOT NULL; keen's columns are nullable.

Usage: postfix.py <repo-dir>
"""
import os
import re
import sys

NULLABLE_TEST = re.compile(
    r'if typed != 0 \{\n(\t+)t\.Fatalf\("[^"]* count = %d, want 0 \(typed insert violated NOT NULL on %q\)", typed, "[^"]+"\)\n',
)


PARENT_ID = re.compile(r'^(\s+"[a-z0-9_]+_id") TEXT NOT NULL,$', re.M)


def nullable_fk_columns(root):
    path = os.path.join(root, "internal/store/store.go")
    if not os.path.exists(path):
        return 0
    text = open(path, encoding="utf-8").read()
    new, n = PARENT_ID.subn(r"\1 TEXT,", text)
    if n:
        open(path, "w", encoding="utf-8").write(new)
    return n


def nullable_fk_tests(root):
    path = os.path.join(root, "internal/store/upsert_batch_test.go")
    if not os.path.exists(path):
        return 0
    text = open(path, encoding="utf-8").read()
    new, n = NULLABLE_TEST.subn(
        lambda m: "if typed != len(items) {\n" + m.group(1)
        + 't.Fatalf("typed count = %d, want %d (typed insert should succeed with nullable columns)", typed, len(items))\n',
        text,
    )
    if n:
        open(path, "w", encoding="utf-8").write(new)
    return n


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit(__doc__)
    print(f"postfix: nullable-fk-columns made {nullable_fk_columns(sys.argv[1])} columns nullable")
    print(f"postfix: nullable-fk-tests rewrote {nullable_fk_tests(sys.argv[1])} assertions")
