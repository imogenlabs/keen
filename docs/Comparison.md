# Comparison

How keen compares with the other ways to reach Atlassian from a terminal or an AI app.

| | keen | mcp-atlassian | Atlassian Rovo MCP | jira-cli |
|---|---|---|---|---|
| Products | Jira, Jira Software, Service Management, Confluence, Assets, Compass, Bitbucket, Marketplace, Atlassian Admin, org chart and Teamwork Graph | Jira, Confluence | Jira, Confluence and other Atlassian Cloud apps | Jira |
| Operations | 1,843 API endpoints and 18 graph queries | 98 tools | Curated tools | Jira commands |
| Command line | Yes | No | No | Yes |
| MCP server | Yes, through two tools (search, then run) | Yes | Yes (hosted by Atlassian) | No |
| Local copy and offline search | Yes | No | No | No |
| Atlassian Cloud | Yes | Yes | Yes | Yes |
| Server and Data Center | No | Yes | No | Yes |
| Install | Homebrew, one binary, Claude Desktop bundle | uvx, pip, Docker | Nothing to install | Homebrew, one binary |
| Price | Free and open source (Apache-2.0) | Free and open source | Included with Atlassian Cloud | Free and open source (MIT) |

## When to choose something else

- **Jira Server or Data Center:** keen covers Atlassian Cloud only. mcp-atlassian and jira-cli support self-hosted Jira.
- **Nothing to install:** Atlassian's own Rovo MCP server runs in Atlassian's cloud.

## Why the tool count does not overwhelm an AI app

Listing 1,939 tools one by one would take about half a million tokens. keen's MCP server offers two instead: `keen_search` finds the tools for a task in plain words, and `keen_run` runs one. See [[Agent Mode]].
