#!/usr/bin/env python3
"""Give generated help examples placeholders that fit the argument.

The generator fills every argument whose name ends in "Key" with
"your-token-here", so `keen issue comment get-issue --help` shows a token where
the issue key goes. This replaces it by argument (or flag) name. Arguments that
really are tokens keep it.

Usage: placeholders.py <dir>           (edits <dir>/internal/cli in place)
       placeholders.py --check <dir>   (lists files it would change; exit 1 if any)
"""
import glob
import os
import re
import sys

TOKEN = "your-token-here"

# Argument or flag name, lower-cased with dashes removed -> placeholder.
BY_NAME = {
    "issueidorkey": "PROJ-42",
    "issuekey": "PROJ-42",
    "idorkey": "PROJ-42",
    "epicidorkey": "PROJ-7",
    "requestidorkey": "HELP-12",
    "issueidsorkeys": "PROJ-42,PROJ-43",
    "projectidorkey": "PROJ",
    "projectkey": "PROJ",
    "projectkeys": "PROJ",
    "spacekey": "TEAM",
    "addonkey": "com.example.my-app",
    "connectkey": "com.example.my-app",
    "propertykey": "my-property",
    "fieldkey": "customfield_10000",
    "fieldidorkey": "customfield_10000",
    "jiraissuefieldskey": "customfield_10000",
    "applicationkey": "jira-software",
    "productkey": "jira-software",
    "licensetypekey": "jira-software",
    "sourcekey": "PROJ",
    "targetkey": "PROJ",
    "themekey": "my-theme",
    "templatekey": "com.atlassian.jira-core-project-templates:jira-core-simplified-process-control",
    "modulekey": "my-module",
    "featurekey": "jsw.classic.roadmap",
    "operationkey": "my-operation",
}

# A bare "key" means something different on each command.
BARE_KEY = {
    "addons": "com.example.my-app",
    "atlassian-connect": "com.example.my-app",
    "project": "PROJ",
    "projects": "PROJ",
    "mypreferences": "my-preference",
    "configuration": "JIRA",
    "applicationrole": "jira-software",
    "wiki": "my-property",
}

USE = re.compile(r'Use:\s+"([^"]*)"')
EXAMPLE = re.compile(r'(Example:\s+")([^"]*)(")')


def placeholder(name, group):
    n = name.lower().replace("-", "")
    if "token" in n:
        return TOKEN
    if n == "key":
        return BARE_KEY.get(group, "example-key")
    return BY_NAME.get(n, "example-key")


def fix_example(example, use):
    words = use.split()
    if not words:
        return example
    args = [re.sub(r"[<>\[\].]", "", w) for w in words[1:] if w[:1] in "<["]
    tokens = example.split(" ")
    # The command group is the word after "keen".
    group = ""
    nonempty = [t for t in tokens if t]
    if len(nonempty) > 1 and nonempty[0] == "keen":
        group = nonempty[1]
    # Positional values follow the command's own name (its last mention
    # before the flags: in "keen issue notify notify X" it is the second).
    prefix = next((k for k, t in enumerate(tokens) if t.startswith("--")), len(tokens))
    hits = [k for k in range(prefix) if tokens[k] == words[0]]
    start = hits[-1] + 1 if hits else len(tokens)
    i, a = start, 0
    while i < len(tokens) and not tokens[i].startswith("--") and a < len(args):
        if tokens[i] == TOKEN:
            tokens[i] = placeholder(args[a], group)
        i, a = i + 1, a + 1
    # Flag values: --name your-token-here
    for j in range(len(tokens) - 1):
        if tokens[j].startswith("--") and tokens[j + 1] == TOKEN:
            tokens[j + 1] = placeholder(tokens[j][2:], group)
    return " ".join(tokens)


def fix_placeholders(root, write=True):
    changed = []
    for path in glob.glob(os.path.join(root, "internal/cli/*.go")):
        text = open(path, encoding="utf-8").read()
        if TOKEN not in text:
            continue
        use = USE.search(text)
        new = EXAMPLE.sub(
            lambda m: m.group(1) + fix_example(m.group(2), use.group(1) if use else "") + m.group(3), text)
        # The env var hint asks for a username, not a token.
        for wrong in ("your-token-here", "<your-token>"):
            new = new.replace(f'JIRA_CLOUD_PLATFORM_USERNAME=\\"{wrong}\\"',
                              'JIRA_CLOUD_PLATFORM_USERNAME=\\"you@example.com\\"')
        if new != text:
            if write:
                open(path, "w", encoding="utf-8").write(new)
            changed.append(os.path.relpath(path, root))
    return changed


if __name__ == "__main__":
    args = sys.argv[1:]
    check = args[:1] == ["--check"]
    if check:
        args = args[1:]
    if len(args) != 1 or not os.path.isdir(args[0]):
        sys.exit(__doc__)
    changed = fix_placeholders(args[0], write=not check)
    if check:
        for f in changed:
            print(f"placeholders: wrong example placeholder in {f}")
        sys.exit(1 if changed else 0)
    print(f"placeholders: fixed examples in {len(changed)} files")
