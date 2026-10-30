#!/usr/bin/env python3
"""Compare Atlassian's current specs with the copies keen was built from.

    scripts/regen/specdiff.py            # download every spec in config.sh and compare
    scripts/regen/specdiff.py --specs D  # compare specs already downloaded into D

Prints, per spec, the operations Atlassian added, removed or changed, and exits
1 when anything differs (0 when keen is current). A spec listed in config.sh
with no copy under specs/ yet counts as new. Used by the weekly spec-check
workflow; run it by hand to see what a regen would pick up. It also compares the
GraphQL schema behind the Teamwork Graph commands (graphql.py --check).
"""
import hashlib
import json
import pathlib
import subprocess
import sys
import urllib.request

REPO = pathlib.Path(__file__).resolve().parents[2]
METHODS = ("get", "put", "post", "delete", "patch", "head", "options")


def spec_list():
    out = subprocess.run(
        ["bash", "-c", 'source scripts/regen/config.sh && printf "%s\\n" "${SPECS[@]}"'],
        cwd=REPO, capture_output=True, text=True, check=True).stdout
    return [line.split("=", 1) for line in out.splitlines() if line]


def operations(spec):
    ops = {}
    for path, item in spec.get("paths", {}).items():
        for method, op in item.items():
            if method in METHODS:
                text = json.dumps(op, sort_keys=True)
                ops[f"{method.upper()} {path}"] = hashlib.sha256(text.encode()).hexdigest()
    return ops


def main():
    specs_dir = None
    if len(sys.argv) == 3 and sys.argv[1] == "--specs":
        specs_dir = pathlib.Path(sys.argv[2])
    changed = False
    for name, url in spec_list():
        if specs_dir:
            new = json.loads((specs_dir / f"{name}.json").read_text())
        else:
            with urllib.request.urlopen(url, timeout=60) as r:
                new = json.load(r)
        if "paths" not in new:
            sys.exit(f"specdiff: {name} from {url} is not an OpenAPI document")
        old_file = REPO / "specs" / f"{name}.json"
        if not old_file.exists():
            print(f"## {name} (new, {len(operations(new))} operations)\n{url}\n")
            changed = True
            continue
        old = json.loads(old_file.read_text())
        old_ops, new_ops = operations(old), operations(new)
        added = sorted(new_ops.keys() - old_ops.keys())
        removed = sorted(old_ops.keys() - new_ops.keys())
        modified = sorted(k for k in old_ops.keys() & new_ops.keys() if old_ops[k] != new_ops[k])
        # Request and response bodies live in components; a change there changes flags.
        schemas = old.get("components") != new.get("components")
        if not (added or removed or modified or schemas):
            continue
        changed = True
        print(f"## {name}: {len(added)} added, {len(removed)} removed, {len(modified)} changed"
              + (", shared schemas changed" if schemas else ""))
        for label, items in (("added", added), ("removed", removed), ("changed", modified)):
            for op in items[:50]:
                print(f"- {label}: `{op}`")
            if len(items) > 50:
                print(f"- ...and {len(items) - 50} more {label}")
        print()
    # The Teamwork Graph commands come from the GraphQL schema, not an OpenAPI spec.
    if subprocess.run([sys.executable, str(REPO / "scripts" / "regen" / "graphql.py"), "--check"]).returncode:
        changed = True
    if not changed:
        print("specdiff: keen is built from Atlassian's current specs")
    sys.exit(1 if changed else 0)


if __name__ == "__main__":
    main()
