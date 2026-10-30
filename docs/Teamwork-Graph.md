# Org Chart and Teamwork Graph

Atlassian's Teamwork Graph links people, teams, projects and work across your site. keen reads it with your normal login: no Rovo credits, no second tool.

```bash
# Who is my manager, who reports to me, my whole chain up to ten levels
keen teamwork-graph user-manager --agent
keen teamwork-graph user-direct-reports --agent
keen teamwork-graph user-report-chain --agent

# Someone else: an account id or a full ARI
keen teamwork-graph user-manager --user 5b10ac8d82e05b22cc7d4ef5 --agent

# Teams, and who is in one
keen teamwork-graph user-teams --agent
keen teamwork-graph team-users --team <team-id> --agent

# What I created, commented on, was mentioned in, and who mentions me most
keen teamwork-graph user-created --first 20 --agent
keen teamwork-graph user-top-recent-mentioners --agent

# What I viewed in a date range
keen teamwork-graph user-viewed --start-date 2026-10-01T00:00:00Z --end-date 2026-10-30T00:00:00Z \
  --viewed-entity-types JiraIssue,ConfluencePage --agent
```

Run `keen teamwork-graph --help` for all 18 queries. `--user` defaults to you.

## Output

Each answer is a list of rows plus paging:

```json
{"rows": [{"entity": {"id": "ari:cloud:jira:...:issue/22167", "type": "JiraIssue"}}],
 "pageInfo": {"hasNextPage": true, "endCursor": "IAM="}}
```

Pass `--after <endCursor>` for the next page. A graph item is `{id, type}`; read it with the matching keen command (for example `keen get` for a Jira issue). On a site with no org chart or teams, these queries return empty rows: that is the answer, not an error.

## Any other GraphQL query

`keen graphql` runs a named read query on your site's gateway:

```bash
keen graphql --query 'query Me { me { user { name } } }' --agent
keen graphql --file query.graphql --variables '{"id":"..."}'
```

Experimental fields need `@optIn(to: "...")`; the error message names the value. Mutations are refused. An error in the response fails the command, even though Atlassian answers it with HTTP 200.

## Keeping up with Atlassian

The 18 commands are built from Atlassian's GraphQL schema. The weekly spec check compares the live schema with `specs/graphql.json` and opens an issue when a query changes; a refresh (`scripts/regen/regen.sh`) rebuilds the commands from it.
