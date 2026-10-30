// Copyright 2026 smashingtags. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"keen/internal/client"
)

// Org chart and Teamwork Graph: Atlassian's teamworkGraph_* queries exist only
// on the GraphQL gateway. Their arguments, opt-in names and headers come from
// teamworkGraphQueries (teamwork_graph_gen.go), which scripts/regen/graphql.py
// writes from the live schema, so a schema change is picked up by a refresh.

type graphArg struct {
	Name string // GraphQL argument name
	Type string // GraphQL type, e.g. "ID!" or "[TeamworkGraphUserViewedEntityType!]!"
}

type graphQuery struct {
	Command     string
	Field       string
	Summary     string
	OptIn       string // value for @optIn(to: ...) on experimental fields
	ForceDynamo bool   // the field's description asks for X-Force-Dynamo: true
	Returns     string
	Args        []graphArg
}

// cypherConnection is the type every teamworkGraph_* field returns today. Its
// rows are columns of key/value pairs; each value is one member of a union.
const cypherConnection = "GraphStoreCypherQueryV2Connection"

const cypherSelection = `edges { cursor node { columns { key value { __typename
  ... on GraphStoreCypherQueryV2AriNode { id data { __typename } }
  ... on GraphStoreCypherQueryV2StringObject { string: value }
  ... on GraphStoreCypherQueryV2IntObject { int: value }
  ... on GraphStoreCypherQueryV2FloatObject { float: value }
  ... on GraphStoreCypherQueryV2BooleanObject { bool: value }
  ... on GraphStoreCypherQueryV2TimestampObject { timestamp: value }
  ... on GraphStoreCypherQueryV2StringListObject { strings: values }
  ... on GraphStoreCypherQueryV2IntListObject { ints: values }
  ... on GraphStoreCypherQueryV2FloatListObject { floats: values }
  ... on GraphStoreCypherQueryV2BoolListObject { bools: values }
  ... on GraphStoreCypherQueryV2TimestampListObject { timestamps: values }
  ... on GraphStoreCypherQueryV2NodeList { nodes { id data { __typename } } }
  ... on GraphStoreCypherQueryV2Path { elements } } } } }
pageInfo { hasNextPage endCursor }`

// Ids the gateway wants as ARIs; a bare id gets the prefix.
var ariPrefix = map[string]string{
	"userId": "ari:cloud:identity::user/",
	"teamId": "ari:cloud:identity::team/",
}

func addTeamworkGraphCommands(rootCmd *cobra.Command, flags *rootFlags) {
	group := &cobra.Command{
		Use:   "teamwork-graph",
		Short: "Org chart and Teamwork Graph: managers, reports, teams, collaborators and a person's work",
		Long: "Atlassian's Teamwork Graph queries. Users default to you; give an account id or ARI with --user.\n" +
			"On a site without an org chart or teams the answers are empty, not errors.",
	}
	for _, q := range teamworkGraphQueries {
		group.AddCommand(newTeamworkGraphCmd(flags, q))
	}
	rootCmd.AddCommand(group, newGraphQLCmd(flags))
}

func argFlagName(arg string) string {
	name := strings.TrimSuffix(arg, "Id")
	return strings.ToLower(regexp.MustCompile(`([a-z0-9])([A-Z])`).ReplaceAllString(name, "$1-$2"))
}

func newTeamworkGraphCmd(flags *rootFlags, q graphQuery) *cobra.Command {
	strs := map[string]*string{}
	ints := map[string]*int{}
	lists := map[string]*[]string{}
	cmd := &cobra.Command{
		Use:         q.Command,
		Short:       q.Summary,
		Args:        cobra.NoArgs,
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSuffix(q.Returns, "!") != cypherConnection {
				return fmt.Errorf("keen cannot read %s results yet; run the query with keen graphql", q.Returns)
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			vars := map[string]any{}
			for _, a := range q.Args {
				fl := argFlagName(a.Name)
				switch {
				case ints[a.Name] != nil && cmd.Flags().Changed(fl):
					vars[a.Name] = *ints[a.Name]
				case lists[a.Name] != nil && (cmd.Flags().Changed(fl) || strings.HasSuffix(a.Type, "!")):
					vars[a.Name] = append([]string{}, *lists[a.Name]...)
				case strs[a.Name] != nil && *strs[a.Name] != "":
					vars[a.Name] = withARIPrefix(a.Name, *strs[a.Name])
				case a.Name == "userId":
					me, err := currentAccountID(cmd.Context(), c)
					if err != nil {
						return err
					}
					vars[a.Name] = ariPrefix["userId"] + me
				case strings.HasSuffix(a.Type, "!"):
					return usageErr(fmt.Errorf("--%s is required", fl))
				}
			}
			data, err := runTeamworkGraph(cmd.Context(), c, q, vars)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			return flags.printJSON(cmd, data)
		},
	}
	for _, a := range q.Args {
		fl := argFlagName(a.Name)
		base := strings.Trim(a.Type, "[]!")
		switch {
		case strings.HasPrefix(a.Type, "["):
			lists[a.Name] = cmd.Flags().StringSlice(fl, nil, base+" values, comma separated")
		case base == "Int":
			ints[a.Name] = cmd.Flags().Int(fl, 0, a.Name)
		default:
			help := a.Name
			if a.Name == "userId" {
				help = "account id or ARI (default: you)"
			} else if strings.HasSuffix(a.Type, "!") {
				help += " (required)"
			}
			if base == "DateTime" {
				help += ", e.g. 2026-10-01T00:00:00Z"
			}
			strs[a.Name] = cmd.Flags().String(fl, "", help)
		}
	}
	cmd.Example = "  keen teamwork-graph " + q.Command + exampleFlags(q)
	return cmd
}

func exampleFlags(q graphQuery) string {
	var b strings.Builder
	for _, a := range q.Args {
		if !strings.HasSuffix(a.Type, "!") || a.Name == "userId" {
			continue
		}
		switch strings.Trim(a.Type, "[]!") {
		case "DateTime":
			b.WriteString(" --" + argFlagName(a.Name) + " 2026-10-01T00:00:00Z")
		default:
			b.WriteString(" --" + argFlagName(a.Name) + " <" + argFlagName(a.Name) + ">")
		}
	}
	return b.String()
}

func withARIPrefix(arg, v string) string {
	if p, ok := ariPrefix[arg]; ok && !strings.HasPrefix(v, "ari:") {
		return p + v
	}
	return v
}

func currentAccountID(ctx context.Context, c *client.Client) (string, error) {
	raw, err := c.Get(ctx, "/rest/api/3/myself", nil)
	if err != nil {
		return "", err
	}
	var me struct {
		AccountID string `json:"accountId"`
	}
	if err := json.Unmarshal(raw, &me); err != nil {
		return "", fmt.Errorf("finding your account id: %w", err)
	}
	if me.AccountID == "" {
		return "", fmt.Errorf("finding your account id: /rest/api/3/myself returned none")
	}
	return me.AccountID, nil
}

// teamworkGraphQuery builds the named query for q with the given variables.
func teamworkGraphQuery(q graphQuery, vars map[string]any) string {
	var decls, args []string
	for _, a := range q.Args {
		if _, ok := vars[a.Name]; ok {
			decls = append(decls, "$"+a.Name+": "+a.Type)
			args = append(args, a.Name+": $"+a.Name)
		}
	}
	var b strings.Builder
	b.WriteString("query KeenTeamworkGraph")
	if len(decls) > 0 {
		b.WriteString("(" + strings.Join(decls, ", ") + ")")
	}
	b.WriteString(" { " + q.Field)
	if len(args) > 0 {
		b.WriteString("(" + strings.Join(args, ", ") + ")")
	}
	if q.OptIn != "" {
		b.WriteString(` @optIn(to: "` + q.OptIn + `")`)
	}
	b.WriteString(" { " + cypherSelection + " } }")
	return b.String()
}

func runTeamworkGraph(ctx context.Context, c *client.Client, q graphQuery, vars map[string]any) (any, error) {
	headers, err := graphHeaders(ctx, c)
	if err != nil {
		return nil, err
	}
	if q.ForceDynamo {
		headers["X-Force-Dynamo"] = "true"
	}
	data, err := c.GraphQL(ctx, "KeenTeamworkGraph", teamworkGraphQuery(q, vars), vars, headers)
	if err != nil {
		return nil, err
	}
	var out map[string]cypherResult
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("reading %s: %w", q.Field, err)
	}
	return out[q.Field].rows(), nil
}

// graphHeaders names the site as the query context; the gateway refuses
// Teamwork Graph queries without it.
func graphHeaders(ctx context.Context, c *client.Client) (map[string]string, error) {
	cloudID, err := c.SiteCloudID(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]string{"X-Query-Context": "ari:cloud:platform::site/" + cloudID}, nil
}

type cypherResult struct {
	Edges []struct {
		Cursor string `json:"cursor"`
		Node   struct {
			Columns []struct {
				Key   string                     `json:"key"`
				Value map[string]json.RawMessage `json:"value"`
			} `json:"columns"`
		} `json:"node"`
	} `json:"edges"`
	PageInfo struct {
		HasNextPage bool    `json:"hasNextPage"`
		EndCursor   *string `json:"endCursor"`
	} `json:"pageInfo"`
}

// rows flattens each edge to {column: value}. A graph node becomes {id, type},
// a list of nodes a list of those, and a scalar or list its value.
func (r cypherResult) rows() map[string]any {
	rows := make([]map[string]any, 0, len(r.Edges))
	for _, e := range r.Edges {
		row := map[string]any{}
		for _, col := range e.Node.Columns {
			row[col.Key] = columnValue(col.Value)
		}
		rows = append(rows, row)
	}
	page := map[string]any{"hasNextPage": r.PageInfo.HasNextPage}
	if r.PageInfo.EndCursor != nil {
		page["endCursor"] = *r.PageInfo.EndCursor
	}
	return map[string]any{"rows": rows, "pageInfo": page}
}

func columnValue(v map[string]json.RawMessage) any {
	if v == nil {
		return nil
	}
	if _, ok := v["id"]; ok {
		return graphNode(v)
	}
	if raw, ok := v["nodes"]; ok {
		var nodes []map[string]json.RawMessage
		if json.Unmarshal(raw, &nodes) == nil {
			out := make([]any, 0, len(nodes))
			for _, n := range nodes {
				out = append(out, graphNode(n))
			}
			return out
		}
	}
	for k, raw := range v {
		if k != "__typename" {
			return raw
		}
	}
	return v
}

func graphNode(v map[string]json.RawMessage) map[string]any {
	node := map[string]any{"id": json.RawMessage(v["id"])}
	var data struct {
		Typename string `json:"__typename"`
	}
	if raw, ok := v["data"]; ok && json.Unmarshal(raw, &data) == nil && data.Typename != "" {
		node["type"] = data.Typename
	}
	return node
}

func newGraphQLCmd(flags *rootFlags) *cobra.Command {
	var query, variables, file string
	var headers []string
	cmd := &cobra.Command{
		Use:   "graphql",
		Short: "Run a read query on Atlassian's GraphQL gateway",
		Long: "Runs a named GraphQL query on your site's gateway and prints its data. The site is sent as\n" +
			"the query context. Experimental fields need @optIn(to: \"...\") in the query; the error says which.\n" +
			"Errors in the response fail the command, even though the gateway answers them with HTTP 200.\n" +
			"Queries only: mutations are refused, so agents can use this as a read-only tool.",
		Example: "  keen graphql --query 'query Me { me { user { name } } }'\n" +
			"  keen graphql --file query.graphql --variables '{\"id\":\"...\"}' --header X-Force-Dynamo=true",
		Args:        cobra.NoArgs,
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if file != "" {
				b, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				query = string(b)
			}
			if regexp.MustCompile(`(?m)^\s*(mutation|subscription)\b`).MatchString(query) {
				return usageErr(fmt.Errorf("keen graphql runs queries only, not mutations or subscriptions"))
			}
			op := operationName(query)
			if op == "" {
				return usageErr(fmt.Errorf("give the query a name (query MyQuery { ... }); the gateway requires one"))
			}
			vars := map[string]any{}
			if variables != "" {
				if err := json.Unmarshal([]byte(variables), &vars); err != nil {
					return usageErr(fmt.Errorf("--variables is not a JSON object: %w", err))
				}
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			h, err := graphHeaders(cmd.Context(), c)
			if err != nil {
				return err
			}
			for _, kv := range headers {
				k, v, ok := strings.Cut(kv, "=")
				if !ok {
					return usageErr(fmt.Errorf("--header %q: want Name=value", kv))
				}
				h[k] = v
			}
			data, err := c.GraphQL(cmd.Context(), op, query, vars, h)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			return flags.printJSON(cmd, data)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "the GraphQL query, with an operation name")
	cmd.Flags().StringVar(&file, "file", "", "read the query from a file")
	cmd.Flags().StringVar(&variables, "variables", "", "variables as a JSON object")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "extra request header, Name=value (repeatable)")
	return cmd
}

var operationNamePattern = regexp.MustCompile(`(?m)^\s*query\s+([A-Za-z_][A-Za-z0-9_]*)`)

func operationName(query string) string {
	if m := operationNamePattern.FindStringSubmatch(query); m != nil {
		return m[1]
	}
	return ""
}
