// Local customization: Confluence page create/get/update on the v2 API.
// Atlassian moved content CRUD out of the v1 spec keen is generated from,
// so these are hand-written. See .printing-press-patches.json (wiki-pages).

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"keen/internal/client"

	"github.com/spf13/cobra"
)

type wikiPage struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	SpaceID string `json:"spaceId"`
	Status  string `json:"status"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Value string `json:"value"`
		} `json:"storage"`
	} `json:"body"`
	Links struct {
		Base  string `json:"base"`
		WebUI string `json:"webui"`
	} `json:"_links"`
}

type wikiPageSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	URL     string `json:"url,omitempty"`
	Body    string `json:"body,omitempty"`
}

func (p wikiPage) summary(baseURL string, withBody bool) wikiPageSummary {
	s := wikiPageSummary{ID: p.ID, Title: p.Title, Version: p.Version.Number}
	if p.Links.WebUI != "" {
		base := p.Links.Base
		if base == "" {
			base = baseURL + "/wiki"
		}
		s.URL = base + p.Links.WebUI
	}
	if withBody {
		s.Body = p.Body.Storage.Value
	}
	return s
}

// readPageBody returns --body, or stdin when --stdin is set.
func readPageBody(body string, fromStdin bool) (string, error) {
	if !fromStdin {
		return body, nil
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	return string(b), nil
}

// readOnly runs a GET even under --dry-run, so a dry run can still show the
// exact write it would make.
func readOnly(c *client.Client, fn func() (json.RawMessage, error)) (json.RawMessage, error) {
	dry := c.DryRun
	c.DryRun = false
	defer func() { c.DryRun = dry }()
	return fn()
}

func getWikiPage(cmd *cobra.Command, c *client.Client, id string) (wikiPage, error) {
	var p wikiPage
	data, err := readOnly(c, func() (json.RawMessage, error) {
		return c.Get(cmd.Context(), "/wiki/api/v2/pages/"+id, map[string]string{"body-format": "storage"})
	})
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("parsing page %s: %w", id, err)
	}
	return p, nil
}

func addWikiPageCommands(rootCmd *cobra.Command, flags *rootFlags) {
	wiki, _, err := rootCmd.Find([]string{"wiki"})
	if err != nil || wiki.Name() != "wiki" {
		return
	}
	wiki.AddCommand(newWikiCreatePageCmd(flags), newWikiGetPageCmd(flags), newWikiUpdatePageCmd(flags))
}

func newWikiCreatePageCmd(flags *rootFlags) *cobra.Command {
	var space, title, parentID, body, representation string
	var fromStdin, draft bool
	cmd := &cobra.Command{
		Use:   "create-page",
		Short: "Create a Confluence page",
		Long: "Create a Confluence page in a space. The body is Confluence storage format (XHTML) by\n" +
			"default; --representation wiki takes wiki markup. Markdown is not converted.",
		Example: "  keen wiki create-page --space DOCS --title \"Release notes\" --body \"<p>Hello</p>\"\n" +
			"  keen wiki create-page --space DOCS --title Notes --parent-id 12345 --stdin < page.xhtml",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if space == "" || title == "" {
				return usageErr(fmt.Errorf("--space and --title are required"))
			}
			value, err := readPageBody(body, fromStdin)
			if err != nil {
				return err
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			data, err := readOnly(c, func() (json.RawMessage, error) {
				return c.Get(cmd.Context(), "/wiki/api/v2/spaces", map[string]string{"keys": space})
			})
			if err != nil {
				return classifyAPIError(err, flags)
			}
			var spaces struct {
				Results []struct {
					ID string `json:"id"`
				} `json:"results"`
			}
			if err := json.Unmarshal(data, &spaces); err != nil || len(spaces.Results) == 0 {
				return fmt.Errorf("space %q not found", space)
			}
			req := map[string]any{
				"spaceId": spaces.Results[0].ID,
				"status":  "current",
				"title":   title,
				"body":    map[string]any{"representation": representation, "value": value},
			}
			if draft {
				req["status"] = "draft"
			}
			if parentID != "" {
				req["parentId"] = parentID
			}
			data, _, err = c.Post(cmd.Context(), "/wiki/api/v2/pages", req)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			if flags.dryRun {
				return nil
			}
			var p wikiPage
			if err := json.Unmarshal(data, &p); err != nil {
				return fmt.Errorf("parsing created page: %w", err)
			}
			return flags.printJSON(cmd, p.summary(c.BaseURL, false))
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&title, "title", "", "Page title (required)")
	cmd.Flags().StringVar(&parentID, "parent-id", "", "ID of the parent page")
	cmd.Flags().StringVar(&body, "body", "", "Page body")
	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "Read the page body from stdin")
	cmd.Flags().StringVar(&representation, "representation", "storage", "Body format: storage or wiki")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create the page as a draft")
	return cmd
}

func newWikiGetPageCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "get-page <id>",
		Short:   "Get a Confluence page: title, version, URL and storage body",
		Example: "  keen wiki get-page 12345",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			p, err := getWikiPage(cmd, c, args[0])
			if err != nil {
				return classifyAPIError(err, flags)
			}
			return flags.printJSON(cmd, p.summary(c.BaseURL, true))
		},
	}
}

func newWikiUpdatePageCmd(flags *rootFlags) *cobra.Command {
	var title, body, representation, message string
	var fromStdin bool
	cmd := &cobra.Command{
		Use:   "update-page <id>",
		Short: "Replace a Confluence page's body; the version number is bumped for you",
		Long: "Replace a Confluence page's body. keen reads the current version and sends the next one,\n" +
			"so there is no version to look up. The title stays unless --title is given.",
		Example: "  keen wiki update-page 12345 --stdin < page.xhtml\n" +
			"  keen wiki update-page 12345 --body \"<p>New text</p>\" --message \"Fix typo\"",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value, err := readPageBody(body, fromStdin)
			if err != nil {
				return err
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			cur, err := getWikiPage(cmd, c, args[0])
			if err != nil {
				return classifyAPIError(err, flags)
			}
			if title == "" {
				title = cur.Title
			}
			status := cur.Status
			if status == "" {
				status = "current"
			}
			// Confluence keeps a draft at version 1; only published pages advance.
			next := cur.Version.Number + 1
			if status == "draft" {
				next = cur.Version.Number
			}
			version := map[string]any{"number": next}
			if message != "" {
				version["message"] = message
			}
			req := map[string]any{
				"id":      cur.ID,
				"status":  status,
				"title":   title,
				"body":    map[string]any{"representation": representation, "value": value},
				"version": version,
			}
			data, _, err := c.Put(cmd.Context(), "/wiki/api/v2/pages/"+cur.ID, req)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			if flags.dryRun {
				return nil
			}
			var p wikiPage
			if err := json.Unmarshal(data, &p); err != nil {
				return fmt.Errorf("parsing updated page: %w", err)
			}
			return flags.printJSON(cmd, p.summary(c.BaseURL, false))
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "New title (default: keep the current one)")
	cmd.Flags().StringVar(&body, "body", "", "New page body")
	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "Read the new body from stdin")
	cmd.Flags().StringVar(&representation, "representation", "storage", "Body format: storage or wiki")
	cmd.Flags().StringVar(&message, "message", "", "Version comment")
	return cmd
}
