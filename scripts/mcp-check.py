#!/usr/bin/env python3
"""Test keen-mcp the way an AI app uses it, against a live Atlassian site.

For each task it asks keen_search in plain words, checks the right tool comes
back near the top, runs it through keen_run and compares the answer with the
CLI's. With --claude it also gives Claude Code keen-mcp as its only tools and
checks it answers real questions without any tool call failing.

    scripts/mcp-check.py <dir with keen and keen-mcp>            # search, run, compare
    scripts/mcp-check.py <dir> --claude                          # and a real AI session

It reads the site from keen's normal config and picks its own test data (the
most recently updated issue, the first board, a Confluence page), so it runs
on any site. Tasks that need setup the site does not have (an admin key,
Assets) are skipped, not failed. check.sh runs it when KEEN_LIVE_SMOKE is set.
Exit 1 when anything fails.
"""
import json
import os
import shutil
import subprocess
import sys
import tempfile

TOP = 3  # the right tool must be in the first TOP results


class MCP:
    def __init__(self, binary):
        self.p = subprocess.Popen([binary], stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                                  stderr=subprocess.DEVNULL, text=True)
        self.n = 0
        self.rpc("initialize", {"protocolVersion": "2024-11-05", "capabilities": {},
                                "clientInfo": {"name": "mcp-check", "version": "1"}})
        self.send({"jsonrpc": "2.0", "method": "notifications/initialized"})

    def send(self, msg):
        self.p.stdin.write(json.dumps(msg) + "\n")
        self.p.stdin.flush()

    def rpc(self, method, params):
        self.n += 1
        self.send({"jsonrpc": "2.0", "id": self.n, "method": method, "params": params})
        while True:
            line = self.p.stdout.readline()
            if not line:
                raise RuntimeError("keen-mcp exited")
            msg = json.loads(line)
            if msg.get("id") == self.n:
                return msg

    def call(self, tool, args):
        """Returns (is_error, text)."""
        msg = self.rpc("tools/call", {"name": tool, "arguments": args})
        if "error" in msg:
            return True, json.dumps(msg["error"])
        res = msg["result"]
        return bool(res.get("isError")), res["content"][0]["text"]

    def close(self):
        self.p.terminate()


def cli(keen, *args):
    """The CLI's answer as JSON, without keen's envelope; None when it fails."""
    out = subprocess.run([keen, *args, "--json"], capture_output=True, text=True)
    if out.returncode != 0:
        return None
    try:
        data = json.loads(out.stdout)
    except ValueError:
        return None
    if isinstance(data, dict) and "meta" in data and "results" in data:
        data = data["results"]
    if isinstance(data, dict) and "success" in data and "data" in data:
        data = data["data"]
    return data


def pick(keen):
    """Test data from the site itself."""
    found = cli(keen, "find-issues", "--jql", "updated >= -3650d ORDER BY updated DESC", "--max-results", "1")
    issue_id = found["issues"][0]["id"] if found and found.get("issues") else None
    issue = cli(keen, "issue", "get", issue_id) if issue_id else None
    boards = cli(keen, "agile", "get-all-boards")
    pages = cli(keen, "pages", "get", "--limit", "1")
    return {
        "issue": issue["key"] if issue else None,
        "issue_id": issue_id,
        "board": str(boards["values"][0]["id"]) if boards and boards.get("values") else None,
        "page": pages["results"][0]["id"] if pages and pages.get("results") else None,
    }


def cases(d):
    """(task in words, tool keen_search should find, its arguments, CLI command, fields to compare or None for all)."""
    return [
        ("get issue by key", "issue_get", {"issueIdOrKey": d["issue"]}, ["issue", "get", d["issue"]], ("key", "id")),
        ("search issues with jql", "jira-cloud-platform-search_and-reconsile-issues-using-jql",
         {"jql": f"key = {d['issue']}"}, ["find-issues", "--jql", f"key = {d['issue']}"], ("issues",)),
        ("get comments on an issue", "issue_comment_get", {"issueIdOrKey": d["issue"]},
         ["issue", "comment", "get", d["issue"]], ("total",)),
        ("list sprints on a board", "agile_get-all-sprints", {"boardId": d["board"]},
         ["agile", "get-all-sprints", d["board"]], ("values",)),
        ("get current user", "myself_get-current-user", {}, ["myself", "get-current-user"], ("accountId",)),
        ("confluence page by id", "pages_get-by-id", {"id": d["page"]}, ["pages", "get-by-id", d["page"]], ("id", "title")),
        # Confluence hands out a fresh "next page" cursor on every call, so compare the spaces only.
        ("list confluence spaces", "spaces_get", {}, ["spaces", "get"], ("results",)),
        ("org admin organizations", "orgs_get", {}, ["orgs", "get"], None),
        ("assets usage", "assets-usage_get-tenant-info", {}, ["assets-usage"], None),
        ("bitbucket repositories in workspace", "bitbucket-repositories_list-in-workspace", {"workspace": "atlassian"},
         ["bitbucket-repositories", "list-in-workspace", "atlassian"], ("pagelen",)),
    ]


def same(mcp_text, cli_data, fields):
    try:
        got = json.loads(mcp_text)
    except ValueError:
        return False
    if fields is None:
        return got == cli_data
    return all(got.get(f) == cli_data.get(f) for f in fields)


def gone(mcp, tool, args):
    """True when the tool's API answers 410 Gone, or keen says Atlassian switched it off."""
    err, out = mcp.call("keen_run", {"tool": tool, "arguments": args})
    return err and ("410" in out or "switched this API off" in out)


def tool_check(mcp, keen, data):
    failures = 0
    listed = [t["name"] for t in mcp.rpc("tools/list", {})["result"]["tools"]]
    if sorted(listed) != ["keen_run", "keen_search"]:
        print(f"FAIL  tools/list offers {listed}, want keen_search and keen_run")
        failures += 1
    for task, tool, args, cmd, fields in cases(data):
        if any(v is None for v in args.values()) or None in cmd:
            print(f"SKIP  {task}: the site has no test data for it")
            continue
        expected = cli(keen, *cmd)
        if expected is None:
            print(f"SKIP  {task}: the CLI cannot do it on this site either (keen {' '.join(cmd)})")
            continue
        err, text = mcp.call("keen_search", {"query": task, "limit": 8})
        found = [] if err else json.loads(text)["tools"]
        names = [t["name"] for t in found]
        rank = names.index(tool) + 1 if tool in names else None
        # An AI tries the tools in the order search gives them. One ranked above
        # the right tool that Atlassian has switched off sends it into an error.
        dead = [t["name"] for t in found[:(rank or TOP) - 1] if t.get("readOnly") and gone(mcp, t["name"], args)]
        err, out = mcp.call("keen_run", {"tool": tool, "arguments": args})
        ok = rank is not None and rank <= TOP and not dead and not err and same(out, expected, fields)
        failures += not ok
        print(f"{'PASS' if ok else 'FAIL'}  {task}: search rank {rank}, run {'error' if err else 'ok'}, "
              f"{'same as' if not err and same(out, expected, fields) else 'differs from'} the CLI"
              + (f", ranked above it and switched off: {dead}" if dead else ""))
        if not ok:
            print("      found:", names[:TOP], "| answer:", out[:200].replace("\n", " "))
    return failures


def claude_check(mcp_bin, keen, data):
    """Give Claude Code only keen-mcp and check real questions come back right, with no failed tool call."""
    claude = shutil.which("claude") or os.path.expanduser("~/.local/bin/claude")
    if not os.path.exists(claude):
        print("FAIL  --claude: the claude CLI is not installed")
        return 1
    issue = cli(keen, "issue", "get", data["issue"])
    questions = [
        (f"What is the summary of Jira issue {data['issue']}? Reply with the summary only.",
         issue["fields"]["summary"]),
        (f"Search Jira with the JQL `key = {data['issue']}` and reply with the numeric id of the issue it finds, nothing else.",
         data["issue_id"]),
    ]
    failures = 0
    with tempfile.TemporaryDirectory() as tmp:
        config = os.path.join(tmp, "mcp.json")
        with open(config, "w") as f:
            json.dump({"mcpServers": {"keen": {"command": os.path.abspath(mcp_bin)}}}, f)
        for question, want in questions:
            run = subprocess.run(
                # --tools "" takes away Claude Code's own tools (it would otherwise
                # answer through the keen CLI in a shell); keen-mcp is all it has.
                [claude, "-p", question, "--mcp-config", config, "--strict-mcp-config", "--tools", "",
                 "--allowedTools", "mcp__keen__keen_search", "mcp__keen__keen_run",
                 "--output-format", "stream-json", "--verbose"],
                capture_output=True, text=True, cwd=tmp, timeout=600)
            answer, tool_errors, runs = "", [], 0
            for line in run.stdout.splitlines():
                try:
                    msg = json.loads(line)
                except ValueError:
                    continue
                if msg.get("type") == "result":
                    answer = msg.get("result") or ""
                for block in (msg.get("message") or {}).get("content") or []:
                    if isinstance(block, dict) and block.get("type") == "tool_use":
                        runs += block.get("name") == "mcp__keen__keen_run"
                    if isinstance(block, dict) and block.get("type") == "tool_result" and block.get("is_error"):
                        tool_errors.append(str(block.get("content"))[:160])
            ok = want.lower() in answer.lower() and not tool_errors and runs > 0
            failures += not ok
            print(f"{'PASS' if ok else 'FAIL'}  claude: {question[:70]}")
            if not ok:
                print(f"      want {want!r}, got {answer.strip()[:200]!r}; keen_run calls: {runs}; failed tool calls: {tool_errors}")
    return failures


def main():
    if len(sys.argv) < 2 or sys.argv[1].startswith("-"):
        sys.exit(__doc__)
    keen = os.path.join(sys.argv[1], "keen")
    mcp_bin = os.path.join(sys.argv[1], "keen-mcp")
    data = pick(keen)
    if not data["issue"]:
        sys.exit("mcp-check: could not find an issue on the site; is keen configured? (keen doctor)")
    mcp = MCP(mcp_bin)
    try:
        failures = tool_check(mcp, keen, data)
    finally:
        mcp.close()
    if "--claude" in sys.argv[2:]:
        failures += claude_check(mcp_bin, keen, data)
    print(f"mcp-check: {failures} failed")
    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
