# keen CLI

Combined CLI for multiple API services

Printed by [@smashingtags](https://github.com/smashingtags) (smashingtags).

## Install

The recommended path installs both the `keen` binary and the `pp-keen` agent skill (Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, and other agents supported by the upstream [`skills`](https://github.com/vercel-labs/skills) CLI) in one shot:

```bash
npx -y @mvanhorn/printing-press-library install keen
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press-library install keen --cli-only
```

For skill only — installs the skill into the same agents as the default command above, but skips the CLI binary (use this to update or reinstall just the skill):

```bash
npx -y @mvanhorn/printing-press-library install keen --skill-only
```

To constrain the skill install to one or more specific agents (repeatable — agent names match the [`skills`](https://github.com/vercel-labs/skills) CLI):

```bash
npx -y @mvanhorn/printing-press-library install keen --agent claude-code
npx -y @mvanhorn/printing-press-library install keen --agent claude-code --agent codex
```

### Without Node

The generated install path is category-agnostic until this CLI is published. If `npx` is not available before publish, install Node or use the category-specific Go fallback from the public-library entry after publish.

### Pre-built binary

Download a pre-built binary for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/keen-current). On macOS, clear the Gatekeeper quarantine: `xattr -d com.apple.quarantine <binary>`. On Unix, mark it executable: `chmod +x <binary>`.

<!-- pp-hermes-install-anchor -->
## Install for Hermes

From the Hermes CLI:

```bash
hermes skills install mvanhorn/printing-press-library/cli-skills/pp-keen --force
```

Inside a Hermes chat session:

```bash
/skills install mvanhorn/printing-press-library/cli-skills/pp-keen --force
```

## Install for OpenClaw

Tell your OpenClaw agent (copy this):

```
Install the pp-keen skill from https://github.com/mvanhorn/printing-press-library/tree/main/cli-skills/pp-keen. The skill defines how its required CLI can be installed.
```

## Use with Claude Desktop

This CLI ships an [MCPB](https://github.com/modelcontextprotocol/mcpb) bundle — Claude Desktop's standard format for one-click MCP extension installs (no JSON config required).

To install:

1. Download the `.mcpb` for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/keen-current).
2. Double-click the `.mcpb` file. Claude Desktop opens and walks you through the install.
3. Fill in `JIRA_CLOUD_PLATFORM_USERNAME` when Claude Desktop prompts you.

Requires Claude Desktop 1.0.0 or later. Pre-built bundles ship for macOS Apple Silicon (`darwin-arm64`) and Windows (`amd64`, `arm64`); for other platforms, use the manual config below.

<details>
<summary>Manual JSON config (advanced)</summary>

If you can't use the MCPB bundle (older Claude Desktop, unsupported platform), install the MCP binary and configure it manually.


Install the MCP binary from this CLI's published public-library entry or pre-built release.

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "keen": {
      "command": "keen-mcp",
      "env": {
        "JIRA_CLOUD_PLATFORM_USERNAME": "<your-key>"
      }
    }
  }
}
```

</details>

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your API key from your API provider's developer portal. The key typically looks like a long alphanumeric string.

```bash
export JIRA_CLOUD_PLATFORM_USERNAME="<paste-your-key>"
```

You can also persist this in your config file at `~/.config/keen/config.toml`.

### 3. Verify Setup

```bash
keen doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
keen addon-categories get mock-value
```

## Usage

Run `keen --help` for the full command reference and flag list.

## Commands

### addon

The addon resource is intended to use used by Bitbucket Cloud Connect
Apps, and only supports JWT authentication.

- **`keen addon delete-app`** - Deletes the application for the user.

This endpoint is intended to be used by Bitbucket Connect apps
and only supports JWT authentication -- that is how Bitbucket
identifies the particular installation of the app. Developers
with applications registered in the "Develop Apps" section
of Bitbucket Marketplace need not use this endpoint as
updates for those applications can be sent out via the
UI of that section.

```
$ curl -X DELETE https://api.bitbucket.org/2.0/addon \
  -H "Authorization: JWT <JWT Token>"
```
- **`keen addon update-installed-app`** - Updates the application installation for the user.

This endpoint is intended to be used by Bitbucket Connect apps
and only supports JWT authentication -- that is how Bitbucket
identifies the particular installation of the app. Developers
with applications registered in the "Develop Apps" section
of Bitbucket need not use this endpoint as updates for those
applications can be sent out via the UI of that section.

Passing an empty body will update the installation using the
existing descriptor URL.

```
$ curl -X PUT https://api.bitbucket.org/2.0/addon \
  -H "Authorization: JWT <JWT Token>" \
  --header "Content-Type: application/json" \
  --data '{}'
```

The new `descriptor` for the installation can be also provided
in the body directly.

```
$ curl -X PUT https://api.bitbucket.org/2.0/addon \
  -H "Authorization: JWT <JWT Token>" \
  --header "Content-Type: application/json" \
  --data '{"descriptor": $NEW_DESCRIPTOR}'
```

In both these modes the URL of the descriptor cannot be changed. To
change the descriptor location and upgrade an installation
the request must be made exclusively with a `descriptor_url`.

 ```
$ curl -X PUT https://api.bitbucket.org/2.0/addon \
  -H "Authorization: JWT <JWT Token>" \
  --header "Content-Type: application/json" \
  --data '{"descriptor_url": $NEW_URL}'
```

The `descriptor_url` must exactly match the marketplace registration
that Atlassian has for the application. Contact your Atlassian
developer advocate to update this registration. Once the registration
has been updated you may call this resource for each installation.

Note that the scopes of the application cannot be increased
in the new descriptor nor reduced to none.

### addon-categories

Manage addon categories

- **`keen addon-categories get`** - Get a specific category.
- **`keen addon-categories get-addoncategories`** - Get a list of categories associated with the parent application. 

For non-UPM clients, this resource returns all categories regardless of the specified application.

### addons

Manage addons

- **`keen addons create`** - Create a new app.

The associated vendor account must exist prior to the app's creation. 

Note that the app's first version should be included in the request body (in `_embedded.version`).

This resource requires authentication.
- **`keen addons get`** - Get a specific app.
- **`keen addons get-archived`** - Get a list of archived apps by a specific vendor matching the specified parameters.
- **`keen addons get-vendor`** - Get a list of apps for the specified vendor.
- **`keen addons list`** - Get a list of apps matching the specified parameters.
- **`keen addons list-archived`** - Get a list of archived apps matching the specified parameters.
- **`keen addons list-listings`** - Get a list of app banners matching the specified parameters.
- **`keen addons list-search`** - Get a list of apps with names matching a search term. This resource's responses are limited to just a single page of results. Use the app collection resource (`/rest/2/addons`) for more results.
- **`keen addons update`** - Update a specific app.

The request body must be a valid <a href="https://tools.ietf.org/html/rfc6902">JSON Patch</a> document. The properties which can be referenced in the PATCH are the same ones returned by a GET on this URI.

This resource requires authentication.

### admin

Manage admin

- **`keen admin delete`** - Delete a user in our SCIM DB with a Atlassian Account ID (AAID). This will apply to all directories in your organization matching that AAID and only works for managed users. 
<br/>
<br/>
You will have to completely reprovision the user to their respective groups after deletion. 
<br/>
<br/>
Explore more about [updating managed SCIM email addresses](../../email-change/).
- **`keen admin get-scim-links`** - Get SCIM Links for a Atlassian Account ID (AAID).
- **`keen admin get-scim-links-by-email`** - Get SCIM Links for an email address in an organization.
- **`keen admin unlink-scim-user`** - Unlinks a SCIM user from their Atlassian account without deleting the user.

### admin-control-admin

Manage admin control admin

- **`keen admin-control-admin add-users-to-policy`** - Add users to an authentication policy to address the security of different user sets. [Understand how to add users to a policy and check the status](https://developer.atlassian.com/cloud/admin/auth-policy-cookbook/)

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin attach-detach-resources-v2`** - Add or remove resources to a policy

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin bulk-fetch-auth-policy`** - Get authentication policy information for a given list of managed users. This is a bulk action.
- **`keen admin-control-admin create-policy`** - Create a policy aligned with your organization's standards.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin create-policy-v2`** - Create a policy aligned with your organization's standards.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin create-resource`** - Add a new resource to a policy

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin delete-policy`** - Delete a policy with a policyId

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `delete:policies:admin`
- **`keen admin-control-admin delete-resource`** - Delete one resource from a policy

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin delete-resources`** - Remove all resources from a policy.
- **`keen admin-control-admin delete-resources-v2`** - Remove all resources from a policy.
- **`keen admin-control-admin get-policies`** - Returns comprehensive details on organizational policies, including both rules and resources.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-policies-v2`** - Returns comprehensive details on organizational policies, including both rules and status.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-policy`** - Returns information about a policy by policyId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-policy-v2`** - Returns information about a policy by policyId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-resources`** - #### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-resources-v2`** - #### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`
- **`keen admin-control-admin get-task-status`** - Verify that users are assigned to the intended policy and report errors, if any.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin publish-draft-policies`** - Publishes policies by a bulk request for a specific ruleName. This is the only way to create or modify published policies.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin update-policy`** - Update a policy with a policyId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin update-policy-v2`** - Update a policy with a policyId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin update-resource`** - Delete one resource from a policy

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:policies:admin`
- **`keen admin-control-admin validate-policy`** - Validate a policy to view potential issues in your policy

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:policies:admin`

### admin-key

The [Admin Key feature](https://support.atlassian.com/confluence-cloud/docs/bypass-access-restrictions-on-a-page-with-admin-key/) is only offered with Confluence Cloud Premium and Enterprise to organization or site admins.

**NOTE:** Users can only make calls to our public REST APIs using a valid Admin Key by adding the following HTTP header on their requests: `Atl-Confluence-With-Admin-Key: true`. Visit the [Confluence Cloud Changelog](https://developer.atlassian.com/cloud/confluence/changelog/#CHANGE-2480) for more details.

- **`keen admin-key disable`** - Disables admin key access for the calling user within the site.

**[Permissions](https://support.atlassian.com/user-management/docs/give-users-admin-permissions/#Centralized-user-management-content) required**:
User must be an organization or site admin.
- **`keen admin-key enable`** - Enables admin key access for the calling user within the site. If an admin key already exists for the user, a new one will be issued with an updated expiration time.

**Note:** The `durationInMinutes` field within the request body is optional. If the request body is empty or if the `durationInMinutes` is set to 0 minutes, a new admin key will be issued to the calling user with a default duration of 10 minutes.

**[Permissions](https://support.atlassian.com/user-management/docs/give-users-admin-permissions/#Centralized-user-management-content) required**:
User must be an organization or site admin.
- **`keen admin-key get`** - Returns information about the admin key if one is currently enabled for the calling user within the site.

**[Permissions](https://support.atlassian.com/user-management/docs/give-users-admin-permissions/#Centralized-user-management-content) required**:
User must be an organization or site admin.

### agents

Use the APIs to invoke agents in bulk for multiple issues. Bulk invocation responses report session-start outcomes; the agent's work continues asynchronously.

- **`keen agents start-batch`** - Starts agent sessions for a batch of issues, with an agent and optional prompt specified per invocation. You can submit between 1 and 50 invocations in one request, including multiple invocations for the same issue. Supply numeric issue-ID strings and Rovo agent IDs. Session-start outcomes are returned inline; the agent's work continues asynchronously.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen agents start-uniform`** - Starts the same agent on between 1 and 50 issues, with an optional prompt shared by every invocation. Supply numeric issue-ID strings, and a Rovo agent ID. Session-start outcomes are returned inline; the agent's work continues asynchronously.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.

### agile

Manage agile

- **`keen agile create-board`** - Creates a new board. Board name, type and filter ID is required.

 *  `name` \- Must be less than 255 characters.
 *  `type` \- Valid values: scrum, kanban
 *  `filterId` \- ID of a filter that the user has permissions to view. Note, if the user does not have the 'Create shared objects' permission and tries to create a shared board, a private board will be created instead (remember that board sharing depends on the filter sharing).
 *  `location` \- The container that the board will be located in. `location` must include the `type` property (Valid values: project, user). If choosing 'project', then a project must be specified by a `projectKeyOrId` property in `location`. If choosing 'user', the current user is chosen by default. The `projectKeyOrId` property should not be provided.

Note:

 *  If you want to create a new project with an associated board, use the [Jira platform REST API](https://docs.atlassian.com/jira/REST/latest). For more information, see the [Create project](#api-rest-api-3-project-post) method. The `projectTypeKey` for software boards must be 'software' and the `projectTemplateKey` must be either `com.pyxis.greenhopper.jira:gh-kanban-template` or `com.pyxis.greenhopper.jira:gh-scrum-template`.
 *  You can create a filter using the [Jira REST API](https://docs.atlassian.com/jira/REST/latest). For more information, see the [Create filter](#api-rest-api-3-filter-post) method.
 *  If you do not ORDER BY the Rank field for the filter of your board, you will not be able to reorder issues on the board.
- **`keen agile create-sprint`** - Creates a future sprint. Sprint name and origin board id are required. Start date, end date, and goal are optional.

Note that the sprint name is trimmed. Also, when starting sprints from the UI, the "endDate" set through this call is ignored and instead the last sprint's duration is used to fill the form.
- **`keen agile delete-board`** - Deletes the board. Admin without the view permission can still remove the board.
- **`keen agile delete-board-property`** - Removes the property from the board identified by the id. Ths user removing the property is required to have permissions to modify the board.
- **`keen agile delete-property`** - Removes the property from the sprint identified by the id. Ths user removing the property is required to have permissions to modify the sprint.
- **`keen agile delete-sprint`** - Deletes a sprint. Once a sprint is deleted, all open issues in the sprint will be moved to the backlog.
- **`keen agile estimate-issue-for-board`** - Updates the estimation of the issue. boardId param is required. This param determines which field will be updated on a issue.

Note that this resource changes the estimation field of the issue regardless of appearance the field on the screen.

Original time tracking estimation field accepts estimation in formats like "1w", "2d", "3h", "20m" or number which represent number of minutes. However, internally the field stores and returns the estimation as a number of seconds.

The field used for estimation on the given board can be obtained from [board configuration resource](#agile/1.0/board-getConfiguration). More information about the field are returned by [edit meta resource](#api-rest-api-3-issue-issueIdOrKey-editmeta-get) or [field resource](#api-rest-api-3-field-get).
- **`keen agile get-all-boards`** - Returns all boards. This only includes boards that the user has permission to view.

**Deprecation notice:** The required OAuth 2.0 scopes will be updated on February 15, 2024.

 *  `read:board-scope:jira-software`, `read:project:jira`
- **`keen agile get-all-quick-filters`** - Returns all quick filters from a board, for a given board ID.
- **`keen agile get-all-sprints`** - Returns all sprints from a board, for a given board ID. This only includes sprints that the user has permission to view.
- **`keen agile get-all-versions`** - Returns all versions from a board, for a given board ID. This only includes versions that the user has permission to view. Note, if the user does not have permission to view the board, no versions will be returned at all. Returned versions are ordered by the name of the project from which they belong and then by sequence defined by user.
- **`keen agile get-board`** - Returns the board for the given board ID. This board will only be returned if the user has permission to view it. Admins without the view permission will see the board as a private one, so will see only a subset of the board's data (board location for instance).
- **`keen agile get-board-by-filter-id`** - Returns any boards which use the provided filter id. This method can be executed by users without a valid software license in order to find which boards are using a particular filter.
- **`keen agile get-board-issues-for-epic`** - Returns all issues that belong to an epic on the board, for the given epic ID and the board ID. This only includes issues that the user has permission to view. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen agile get-board-issues-for-sprint`** - Get all issues you have access to that belong to the sprint from the board. Issue returned from this resource contains additional fields like: sprint, closedSprints, flagged and epic. Issues are returned ordered by rank. JQL order has higher priority than default rank.
- **`keen agile get-board-property`** - Returns the value of the property with a given key from the board identified by the provided id. The user who retrieves the property is required to have permissions to view the board.
- **`keen agile get-board-property-keys`** - Returns the keys of all properties for the board identified by the id. The user who retrieves the property keys is required to have permissions to view the board.
- **`keen agile get-configuration`** - Get the board configuration. The response contains the following fields:

 *  `id` \- ID of the board.
 *  `name` \- Name of the board.
 *  `filter` \- Reference to the filter used by the given board.
 *  `location` \- Reference to the container that the board is located in. Includes the container type (Valid values: project, user).
 *  `subQuery` (Kanban only) - JQL subquery used by the given board.
 *  `columnConfig` \- The column configuration lists the columns for the board, in the order defined in the column configuration. For each column, it shows the issue status mapping as well as the constraint type (Valid values: none, issueCount, issueCountExclSubs) for the min/max number of issues. Note, the last column with statuses mapped to it is treated as the "Done" column, which means that issues in that column will be marked as already completed.
 *  `estimation` (Scrum only) - Contains information about type of estimation used for the board. Valid values: none, issueCount, field. If the estimation type is "field", the ID and display name of the field used for estimation is also returned. Note, estimates for an issue can be updated by a PUT /rest/api/3/issue/\{issueIdOrKey\} request, however the fields must be on the screen. "timeoriginalestimate" field will never be on the screen, so in order to update it "originalEstimate" in "timetracking" field should be updated.
 *  `ranking` \- Contains information about custom field used for ranking in the given board.
- **`keen agile get-epic`** - Returns the epic for a given epic ID. This epic will only be returned if the user has permission to view it. **Note:** This operation does not work for epics in next-gen projects.
- **`keen agile get-epics`** - Returns all epics from the board, for the given board ID. This only includes epics that the user has permission to view. Note, if the user does not have permission to view the board, no epics will be returned at all.
- **`keen agile get-features-for-board`** - Get features for board
- **`keen agile get-issue`** - Returns a single issue, for a given issue ID or issue key. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic.
- **`keen agile get-issue-estimation-for-board`** - Returns the estimation of the issue and a fieldId of the field that is used for it. `boardId` param is required. This param determines which field will be updated on a issue.

Original time internally stores and returns the estimation as a number of seconds.

The field used for estimation on the given board can be obtained from [board configuration resource](#agile/1.0/board-getConfiguration). More information about the field are returned by [edit meta resource](#api-rest-api-3-issue-getEditIssueMeta) or [field resource](#api-rest-api-3-field-get).
- **`keen agile get-issues-for-backlog`** - Returns all issues from the board's backlog, for the given board ID. This only includes issues that the user has permission to view. The backlog contains incomplete issues that are not assigned to any future or active sprint. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen agile get-issues-for-board`** - Returns all issues from a board, for a given board ID. This only includes issues that the user has permission to view. An issue belongs to the board if its status is mapped to the board's column. Epic issues do not belongs to the scrum boards. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen agile get-issues-for-epic`** - Returns all issues that belong to the epic, for the given epic ID. This only includes issues that the user has permission to view. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank. **Note:** If you are querying a next-gen project, do not use this operation. Instead, search for issues that belong to an epic by using the [Search for issues using JQL](https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-search-get) operation in the Jira platform REST API. Build your JQL query using the `parent` clause. For more information on the `parent` JQL field, see [Advanced searching](https://confluence.atlassian.com/x/dAiiLQ#Advancedsearching-fieldsreference-Parent).
- **`keen agile get-issues-for-sprint`** - Returns all issues in a sprint, for a given sprint ID. This only includes issues that the user has permission to view. By default, the returned issues are ordered by rank.
- **`keen agile get-issues-without-epic`** - Returns all issues that do not belong to any epic. This only includes issues that the user has permission to view. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank. **Note:** If you are querying a next-gen project, do not use this operation. Instead, search for issues that don't belong to an epic by using the [Search for issues using JQL](https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-search-get) operation in the Jira platform REST API. Build your JQL query using the `parent is empty` clause. For more information on the `parent` JQL field, see [Advanced searching](https://confluence.atlassian.com/x/dAiiLQ#Advancedsearching-fieldsreference-Parent).
- **`keen agile get-issues-without-epic-for-board`** - Returns all issues that do not belong to any epic on a board, for a given board ID. This only includes issues that the user has permission to view. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen agile get-projects`** - Returns all projects that are associated with the board, for the given board ID. If the user does not have permission to view the board, no projects will be returned at all. Returned projects are ordered by the name.

A project is associated with a board if the board filter contains reference the project or there is an issue from the project that belongs to the board.

The board filter contains reference the project only if JQL query guarantees that returned issues will be returned from the project set defined in JQL. For instance the query `project in (ABC, BCD) AND reporter = admin` have reference to ABC and BCD projects but query `project in (ABC, BCD) OR reporter = admin` doesn't have reference to any project.

An issue belongs to the board if its status is mapped to the board's column. Epic issues do not belongs to the scrum boards.
- **`keen agile get-projects-full`** - Returns all projects that are statically associated with the board, for the given board ID. Returned projects are ordered by the name.

A project is associated with a board if the board filter contains reference the project.

The board filter contains reference the project only if JQL query guarantees that returned issues will be returned from the project set defined in JQL. For instance the query `project in (ABC, BCD) AND reporter = admin` have reference to ABC and BCD projects but query `project in (ABC, BCD) OR reporter = admin` doesn't have reference to any project.
- **`keen agile get-properties-keys`** - Returns the keys of all properties for the sprint identified by the id. The user who retrieves the property keys is required to have permissions to view the sprint.
- **`keen agile get-property`** - Returns the value of the property with a given key from the sprint identified by the provided id. The user who retrieves the property is required to have permissions to view the sprint.
- **`keen agile get-quick-filter`** - Returns the quick filter for a given quick filter ID. The quick filter will only be returned if the user can view the board that the quick filter belongs to.
- **`keen agile get-reports-for-board`** - Get reports for board
- **`keen agile get-sprint`** - Returns the sprint for a given sprint ID. The sprint will only be returned if the user can view the board that the sprint was created on, or view at least one of the issues in the sprint.
- **`keen agile move-issues-to-backlog`** - Move issues to the backlog.  
This operation is equivalent to remove future and active sprints from a given set of issues. At most 50 issues may be moved at once.
- **`keen agile move-issues-to-backlog-for-board`** - Move issues to the backlog of a particular board (if they are already on that board).  
This operation is equivalent to remove future and active sprints from a given set of issues if the board has sprints If the board does not have sprints this will put the issues back into the backlog from the board. At most 50 issues may be moved at once.
- **`keen agile move-issues-to-board`** - Move issues from the backog to the board (if they are already in the backlog of that board).  
This operation either moves an issue(s) onto a board from the backlog (by adding it to the issueList for the board) Or transitions the issue(s) to the first column for a kanban board with backlog. At most 50 issues may be moved at once.
- **`keen agile move-issues-to-epic`** - Moves issues to an epic, for a given epic id. Issues can be only in a single epic at the same time. That means that already assigned issues to an epic, will not be assigned to the previous epic anymore. The user needs to have the edit issue permission for all issue they want to move and to the epic. The maximum number of issues that can be moved in one operation is 50. **Note:** This operation does not work for epics in next-gen projects.
- **`keen agile move-issues-to-sprint-and-rank`** - Moves issues to a sprint, for a given sprint ID. Issues can only be moved to open or active sprints. The maximum number of issues that can be moved in one operation is 50.
- **`keen agile partially-update-epic`** - Performs a partial update of the epic. A partial update means that fields not present in the request JSON will not be updated. Valid values for color are `color_1` to `color_9`. **Note:** This operation does not work for epics in next-gen projects.
- **`keen agile partially-update-sprint`** - Performs a partial update of a sprint. A partial update means that fields not present in the request JSON will not be updated.

Notes:

 *  For closed sprints, only the name and goal can be updated; changes to other fields will be ignored.
 *  A sprint can be started by updating the state to 'active'. This requires the sprint to be in the 'future' state and have a startDate and endDate set.
 *  A sprint can be completed by updating the state to 'closed'. This action requires the sprint to be in the 'active' state. This sets the completeDate to the time of the request.
 *  Other changes to state are not allowed.
 *  The completeDate field cannot be updated manually.
- **`keen agile rank-epics`** - Moves (ranks) an epic before or after a given epic.

If rankCustomFieldId is not defined, the default rank field will be used.

**Note:** This operation does not work for epics in next-gen projects.
- **`keen agile rank-issues`** - Moves (ranks) issues before or after a given issue. At most 50 issues may be ranked at once.

This operation may fail for some issues, although this will be rare. In that case the 207 status code is returned for the whole response and detailed information regarding each issue is available in the response body.

If rankCustomFieldId is not defined, the default rank field will be used.
- **`keen agile remove-issues-from-epic`** - Removes issues from epics. The user needs to have the edit issue permission for all issue they want to remove from epics. The maximum number of issues that can be moved in one operation is 50. **Note:** This operation does not work for epics in next-gen projects. Instead, update the issue using `\{ fields: \{ parent: \{\} \} \}`
- **`keen agile set-board-property`** - Sets the value of the specified board's property.

You can use this resource to store a custom data against the board identified by the id. The user who stores the data is required to have permissions to modify the board.
- **`keen agile set-property`** - Sets the value of the specified sprint's property.

You can use this resource to store a custom data against the sprint identified by the id. The user who stores the data is required to have permissions to modify the sprint.
- **`keen agile swap-sprint`** - Swap the position of the sprint with the second sprint.
- **`keen agile toggle-features`** - Toggle features
- **`keen agile update-sprint`** - Performs a full update of a sprint. A full update means that the result will be exactly the same as the request body. Any fields not present in the request JSON will be set to null.

Notes:

 *  For closed sprints, only the name and goal can be updated; changes to other fields will be ignored.
 *  A sprint can be started by updating the state to 'active'. This requires the sprint to be in the 'future' state and have a startDate and endDate set.
 *  A sprint can be completed by updating the state to 'closed'. This action requires the sprint to be in the 'active' state. This sets the completeDate to the time of the request.
 *  Other changes to state are not allowed.
 *  The completeDate field cannot be updated manually.

### announcement-banner

This resource represents an announcement banner. Use it to retrieve and update banner configuration.

- **`keen announcement-banner get-banner`** - Returns the current announcement banner configuration.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen announcement-banner set-banner`** - Updates the announcement banner configuration.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### api-config

Manage api config

- **`keen api-config reference-type-create`** - Update a reference type
- **`keen api-config reference-type-list`** - Get reference type
- **`keen api-config status-create`** - Create a new status
- **`keen api-config status-delete`** - Delete an existing status
- **`keen api-config status-find`** - Find a status by id
- **`keen api-config status-list`** - Find all status
- **`keen api-config status-update`** - Update an existing status

### api-import

Manage api import

- **`keen api-import <id>`** - Start configured imports. To see an ongoing import see the Progress resource

### api-keys

API Key APIs

- **`keen api-keys get-all-by-org-id`** - Gets all user API keys in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:keys:admin`
- **`keen api-keys get-count-by-org-id`** - Gets count of user API keys in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:keys:admin`
- **`keen api-keys revoke`** - Revokes an existing API key

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `delete:keys:admin`

### api-tokens

API Token APIs

- **`keen api-tokens bulk-revoke`** - Revokes all managed user API tokens in an organization by orgID.
#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `delete:tokens:admin`
- **`keen api-tokens get-all-by-org-id`** - Gets all [user API tokens](/cloud/admin/api-access/rest/intro/#API%20Tokens) in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:tokens:admin`
- **`keen api-tokens get-count-by-org-id`** - Gets count of [user API tokens](/cloud/admin/api-access/rest/intro/#API%20Tokens) in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:tokens:admin`

### app

Manage app

- **`keen app get-custom-field-configuration`** - Returns a [paginated](#pagination) list of configurations for a custom field of a [type](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field-type/) created by a [Forge app](https://developer.atlassian.com/platform/forge/).

The result can be filtered by one of these criteria:

 *  `id`.
 *  `fieldContextId`.
 *  `issueId`.
 *  `projectKeyOrId` and `issueTypeId`.

Otherwise, all configurations are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg). Jira permissions are not required for the Forge app that provided the custom field type.
- **`keen app get-custom-fields-configurations`** - Returns a [paginated](#pagination) list of configurations for list of custom fields of a [type](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field-type/) created by a [Forge app](https://developer.atlassian.com/platform/forge/).

The result can be filtered by one of these criteria:

 *  `id`.
 *  `fieldContextId`.
 *  `issueId`.
 *  `projectKeyOrId` and `issueTypeId`.

Otherwise, all configurations for the provided list of custom fields are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg). Jira permissions are not required for the Forge app that provided the custom field type.
- **`keen app update-custom-field-configuration`** - Update the configuration for contexts of a custom field of a [type](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field-type/) created by a [Forge app](https://developer.atlassian.com/platform/forge/).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg). Jira permissions are not required for the Forge app that created the custom field type.
- **`keen app update-custom-field-value`** - Updates the value of a custom field on one or more issues.

Apps can only perform this operation on [custom fields](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field/) and [custom field types](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field-type/) declared in their own manifests.

**[Permissions](#permissions) required:** Only the app that owns the custom field or custom field type can update its values with this operation.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen app update-multiple-custom-field-values`** - Updates the value of one or more custom fields on one or more issues. Combinations of custom field and issue should be unique within the request.

Apps can only perform this operation on [custom fields](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field/) and [custom field types](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field-type/) declared in their own manifests.

**[Permissions](#permissions) required:** Only the app that owns the custom field or custom field type can update its values with this operation.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.

### application-properties

Manage application properties

- **`keen application-properties get-advanced-settings`** - Returns the application properties that are accessible on the *Advanced Settings* page. To navigate to the *Advanced Settings* page in Jira, choose the Jira icon > **Jira settings** > **System**, **General Configuration** and then click **Advanced Settings** (in the upper right).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen application-properties get-application-property`** - Returns all application properties or an application property.

If you specify a value for the `key` parameter, then an application property is returned as an object (not in an array). Otherwise, an array of all editable application properties is returned. See [Set application property](#api-rest-api-3-application-properties-id-put) for descriptions of editable properties.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen application-properties set-application-property`** - Changes the value of an application property. For example, you can change the value of the `jira.clone.prefix` from its default value of *CLONE -* to *Clone -* if you prefer sentence case capitalization. Editable properties are described below along with their default values.

#### Advanced settings ####

The advanced settings below are also accessible in [Jira](https://confluence.atlassian.com/x/vYXKM).

| Key | Description | Default value |  
| -- | -- | -- |  
| `jira.clone.prefix` | The string of text prefixed to the title of a cloned issue. | `CLONE -` |  
| `jira.date.picker.java.format` | The date format for the Java (server-side) generated dates. This must be the same as the `jira.date.picker.javascript.format` format setting. | `d/MMM/yy` |  
| `jira.date.picker.javascript.format` | The date format for the JavaScript (client-side) generated dates. This must be the same as the `jira.date.picker.java.format` format setting. | `%e/%b/%y` |  
| `jira.date.time.picker.java.format` | The date format for the Java (server-side) generated date times. This must be the same as the `jira.date.time.picker.javascript.format` format setting. | `dd/MMM/yy h:mm a` |  
| `jira.date.time.picker.javascript.format` | The date format for the JavaScript (client-side) generated date times. This must be the same as the `jira.date.time.picker.java.format` format setting. | `%e/%b/%y %I:%M %p` |  
| `jira.issue.actions.order` | The default order of actions (such as *Comments* or *Change history*) displayed on the issue view. | `asc` |  
| `jira.view.issue.links.sort.order` | The sort order of the list of issue links on the issue view. | `type, status, priority` |  
| `jira.comment.collapsing.minimum.hidden` | The minimum number of comments required for comment collapsing to occur. A value of `0` disables comment collapsing. | `4` |  
| `jira.newsletter.tip.delay.days` | The number of days before a prompt to sign up to the Jira Insiders newsletter is shown. A value of `-1` disables this feature. | `7` |  


#### Look and feel ####

The settings listed below adjust the [look and feel](https://confluence.atlassian.com/x/VwCLLg).

| Key | Description | Default value |  
| -- | -- | -- |  
| `jira.lf.date.time` | The [ time format](https://docs.oracle.com/javase/6/docs/api/index.html?java/text/SimpleDateFormat.html). | `h:mm a` |  
| `jira.lf.date.day` | The [ day format](https://docs.oracle.com/javase/6/docs/api/index.html?java/text/SimpleDateFormat.html). | `EEEE h:mm a` |  
| `jira.lf.date.complete` | The [ date and time format](https://docs.oracle.com/javase/6/docs/api/index.html?java/text/SimpleDateFormat.html). | `dd/MMM/yy h:mm a` |  
| `jira.lf.date.dmy` | The [ date format](https://docs.oracle.com/javase/6/docs/api/index.html?java/text/SimpleDateFormat.html). | `dd/MMM/yy` |  
| `jira.date.time.picker.use.iso8061` | When enabled, sets Monday as the first day of the week in the date picker, as specified by the ISO8601 standard. | `false` |  
| `jira.lf.logo.url` | The URL of the logo image file. | `/images/icon-jira-logo.png` |  
| `jira.lf.logo.show.application.title` | Controls the visibility of the application title on the sidebar. | `false` |  
| `jira.lf.favicon.url` | The URL of the favicon. | `/favicon.ico` |  
| `jira.lf.favicon.hires.url` | The URL of the high-resolution favicon. | `/images/64jira.png` |  
| `jira.lf.navigation.bgcolour` | The background color of the sidebar. | `#0747A6` |  
| `jira.lf.navigation.highlightcolour` | The color of the text and logo of the sidebar. | `#DEEBFF` |  
| `jira.lf.hero.button.base.bg.colour` | The background color of the hero button. | `#3b7fc4` |  
| `jira.title` | The text for the application title. The application title can also be set in *General settings*. | `Jira` |  
| `jira.option.globalsharing` | Whether filters and dashboards can be shared with anyone signed into Jira. | `true` |  
| `xflow.product.suggestions.enabled` | Whether to expose product suggestions for other Atlassian products within Jira. | `true` |  


#### Other settings ####

| Key | Description | Default value |  
| -- | -- | -- |  
| `jira.issuenav.criteria.autoupdate` | Whether instant updates to search criteria is active. | `true` |  


*Note: Be careful when changing [application properties and advanced settings](https://confluence.atlassian.com/x/vYXKM).*

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### applicationrole

Manage applicationrole

- **`keen applicationrole get-all-application-roles`** - Returns all application roles. In Jira, application roles are managed using the [Application access configuration](https://confluence.atlassian.com/x/3YxjL) page.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen applicationrole get-application-role`** - Returns an application role.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### applications

Information about applications

- **`keen applications get`** - Get a specific application.
- **`keen applications list`** - Get a list of applications.

### aql

Manage aql

- **`keen aql`** - This endpoint is deprecated and will be removed on <b>18th of September 2024</b>. Please use POST `/object/aql` instead. Find objects based on Assets Query Language (AQL)

### assets

For uploading app/image files

- **`keen assets create`** - Use this resource to upload an app artifact (`.jar` or `.obr`) to Marketplace, so it can be used by an app version.

The request body is the file data.

If successful, the resource will respond with a resource link in `_links.self` that can be used as the `artifact` link when creating an app or app version. The response may include additional properties; other than `fileInfo`, these are used internally and should be ignored.

This resource requires authentication.

Example: This curl command uploads the artifact present in a file ('my-file.jar').
``` bash
curl --request POST \
  --url https://marketplace.atlassian.com/rest/2/assets/artifact \
  --user 'email@example.com:<api_token>' \
  --header 'Accept: application/json' \
  --header 'Content-Type: multipart/form-data' \
  --form file=@"my-file.jar"
```
- **`keen assets create-artifact`** - Use this resource to copy and validate an app artifact (`.jar` or `.obr`) that is available from an external URL, so it can be used by an app version.

If successful, you will get a resource link in `_links.self` that can be used as the `artifact` link when creating an app or app version.

Note that this step can be skipped by just setting the `artifact` link to the external URL where the artifact is hosted when creating an app version; Marketplace will then download the artifact automatically.

This resource requires authentication.
- **`keen assets create-image`** - Use this resource to upload an image (`.jpg`, `.png`, `.gif`, `.bmp`) to Marketplace, so you can use it for an app or vendor.

The request body is the file data.

If successful, you will get a resource link in _links.self that you can use as the appropriate image link in another resource, for example the `logo` link if you are creating an app.

Note that if the image is available from an external URL, you can skip this step by just setting the image link tothat external URL when you are creating the resource that uses it.
- **`keen assets delete`** - Removes a file from Marketplace storage which had been previously uploaded.

This is only allowed if:

1. you, rather than another user, uploaded the file, and
2. the asset is not currently being used by any other resource (for example, it cannot be the logo image for an existing app).

This resource requires authentication.
- **`keen assets get`** - Returns details about a previously uploaded file that Marketplace is now storing.

The asset can either be an image file or an app artifact. If it is an app artifact, the response may include additional properties which are used internally and should be ignored.
- **`keen assets list`** - Get links to permitted asset resources.

### atlassian-access-events

Events APIs

- **`keen atlassian-access-events get`** - Returns a filtered list of audit log events for an organization.
Use this endpoint for more granular and detailed querying. 

If you simply need to paginate through all events, consider using the [/events-stream](https://developer.atlassian.com/cloud/admin/organization/rest/api-group-events/#api-v1-orgs-orgid-events-stream-get) endpoint.  

These rate limits for this endpoint be lowered effective end of May 2025 as follows:
 - *Rate limit per user*: *10* requests per minute
 - *Rate limit per API path*: *10* requests per minute

 Please migrate to the polling API to guarantee uninterrupted service for use cases involving a high request rate.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:events:admin`
- **`keen atlassian-access-events get-by-id`** - **This API is deprecated. Requests to it fail starting December 21, 2026.** Use the [Get events endpoint](https://developer.atlassian.com/cloud/admin/organization/rest/api-group-events/#api-v1-orgs-orgid-events-get) or the [Poll events endpoint](https://developer.atlassian.com/cloud/admin/organization/rest/api-group-events/#api-v1-orgs-orgid-events-stream-get) instead; both return the same event fields. See [Migrate from the single-event lookup](/cloud/admin/organization/event-by-id-deprecation/) for guidance.

Returns information about a single event by ID.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:events:admin`

### atlassian-access-policies

Policies APIs

- **`keen atlassian-access-policies create-policy`** - Create a policy for an org
- **`keen atlassian-access-policies delete-policy`** - Delete a policy for an org
- **`keen atlassian-access-policies get`** - Returns information about org policies
- **`keen atlassian-access-policies get-policy-by-id`** - Returns information about a single policy by ID
- **`keen atlassian-access-policies update-policy`** - Update a policy for an org

### atlassian-access-users

Orgs Users APIs

- **`keen atlassian-access-users create`** - Invite people to your organization. When you invite someone:
- they’re given app roles according to your invitation.
- they’re added to directories based on apps in your invitation.
- they’re added to groups according to your invitation.
- they receive an email invitation if the `sendNotification` field is set to `true` and the `notificationText` field contains a message to include in the email invitation.

**This API is only available to customers who have at least one paid subscription in their organization.**
- **`keen atlassian-access-users get`** - Returns a list of managed accounts in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:accounts:admin`

### atlassian-connect

Manage atlassian connect

- **`keen atlassian-connect addon-properties-resource-delete-addon-property-delete`** - Deletes an app's property.

**[Permissions](#permissions) required:** Only a Connect app whose key matches `addonKey` can make this request.
Additionally, Forge apps can access Connect app properties (stored against the same `app.connect.key`).
- **`keen atlassian-connect addon-properties-resource-get-addon-properties-get`** - Gets all the properties of an app. The reserved key `connect_client_key_019cdff3-8bfb-71fe-9628-875b700aebb8` is not returned.

**[Permissions](#permissions) required:** Only a Connect app whose key matches `addonKey` can make this request.
Additionally, Forge apps can access Connect app properties (stored against the same `app.connect.key`).
- **`keen atlassian-connect addon-properties-resource-get-addon-property-get`** - Returns the key and value of an app's property. The property key `connect_client_key_019cdff3-8bfb-71fe-9628-875b700aebb8`
is reserved. It returns a synthetic, read-only property containing the Connect `clientKey` for the requested tenant.
This is intended for Forge apps with `app.connect.key` to retrieve the Connect client key during migration.

**[Permissions](#permissions) required:** Only a Connect app whose key matches `addonKey` can make this request.
Additionally, Forge apps can access Connect app properties (stored against the same `app.connect.key`).
- **`keen atlassian-connect addon-properties-resource-put-addon-property-put`** - Sets the value of an app's property. Use this resource to store custom data for your app.

The value of the request body must be a [valid](http://tools.ietf.org/html/rfc4627), non-empty JSON blob. The maximum length is 32768 characters.

**[Permissions](#permissions) required:** Only a Connect app whose key matches `addonKey` can make this request.
Additionally, Forge apps can access Connect app properties (stored against the same `app.connect.key`).
- **`keen atlassian-connect app-issue-field-value-update-resource-update-issue-fields-put`** - Updates the value of a custom field added by Connect apps on one or more issues.
The values of up to 200 custom fields can be updated.

**[Permissions](#permissions) required:** Only Connect apps can make this request
- **`keen atlassian-connect connect-to-forge-migration-fetch-task-resource-fetch-migration-task-get`** - Returns the details of a Connect issue field's migration to Forge.

When migrating a Connect app to Forge, [Issue Field](https://developer.atlassian.com/cloud/jira/software/modules/issue-field/) modules
must be converted to [Custom field](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field/). When the
Forge version of the app is installed, Forge creates a
[background task](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-tasks/#api-group-tasks) to track the
migration of field data across. This endpoint returns the status and other details of that background task.

For more details, see
[Jira modules > Jira Custom Fields](https://developer.atlassian.com/platform/adopting-forge-from-connect/migrate-jira-custom-fields/).

**[Permissions](#permissions) required:** Only Connect and Forge apps can make this request.
- **`keen atlassian-connect connect-to-forge-migration-task-submission-resource-submit-task-post`** - Submits a request to trigger migration of connect issue field to its Forge custom field counterpart.

When migrating a Connect app to Forge, [Issue Field](https://developer.atlassian.com/cloud/jira/software/modules/issue-field/) modules
must be converted to [Custom field](https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-custom-field/) modules.
This endpoint triggers the background migration of field data. Use the GET endpoint to retrieve
the status and progress of the task.

For more details, see
[Jira modules > Jira Custom Fields](https://developer.atlassian.com/platform/adopting-forge-from-connect/migrate-jira-custom-fields/).

**[Permissions](#permissions) required:** Only Connect and Forge apps can make this request.
- **`keen atlassian-connect dynamic-modules-resource-get-modules-get`** - Returns all modules registered dynamically by the calling app.

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen atlassian-connect dynamic-modules-resource-register-modules-post`** - Registers a list of modules.

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen atlassian-connect dynamic-modules-resource-remove-modules-delete`** - Remove all or a list of modules registered by the calling app.

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen atlassian-connect migration-resource-update-entity-properties-value-put`** - Updates the values of multiple entity properties for an object, up to 50 updates per request. This operation is for use by Connect apps during app migration.
- **`keen atlassian-connect migration-resource-workflow-rule-search-post`** - Returns configurations for workflow transition rules migrated from server to cloud and owned by the calling Connect app.
- **`keen atlassian-connect service-registry-resource-services-get`** - Retrieve the attributes of given service registries.

**[Permissions](#permissions) required:** Only Connect apps can make this request and the servicesIds belong to the tenant you are requesting

### attachment

Manage attachment

- **`keen attachment get`** - Returns the metadata for an attachment. Note that the attachment itself is not returned.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  If attachments are added in private comments, the comment-level restriction will be applied.
- **`keen attachment get-content`** - Returns the contents of an attachment. A `Range` header can be set to define a range of bytes within the attachment to download. See the [HTTP Range header standard](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range) for details.

To return a thumbnail of the attachment, use [Get attachment thumbnail](#api-rest-api-3-attachment-thumbnail-id-get).

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** For the issue containing the attachment:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  If attachments are added in private comments, the comment-level restriction will be applied.
- **`keen attachment get-meta`** - Returns the attachment settings, that is, whether attachments are enabled and the maximum attachment size allowed.

Note that there are also [project permissions](https://confluence.atlassian.com/x/yodKLg) that restrict whether users can create and delete attachments.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen attachment get-thumbnail`** - Returns the thumbnail of an attachment.

To return the attachment contents, use [Get attachment content](#api-rest-api-3-attachment-content-id-get).

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** For the issue containing the attachment:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  If attachments are added in private comments, the comment-level restriction will be applied.
- **`keen attachment remove`** - Deletes an attachment from an issue.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** For the project holding the issue containing the attachment:

 *  *Delete own attachments* [project permission](https://confluence.atlassian.com/x/yodKLg) to delete an attachment created by the calling user.
 *  *Delete all attachments* [project permission](https://confluence.atlassian.com/x/yodKLg) to delete an attachment created by any user.

### attachments

Manage attachments

- **`keen attachments delete`** - Delete an attachment by id.

Deleting an attachment moves the attachment to the trash, where it can be restored later. To permanently delete an attachment (or "purge" it),
the endpoint must be called on a **trashed** attachment with the following param `purge=true`.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
Permission to view the container of the attachment.
Permission to delete attachments in the space.
[`manage/content`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space (if attempting to purge).

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen attachments get`** - Returns all attachments. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the container of the attachment.
- **`keen attachments get-by-id`** - Returns a specific attachment.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the attachment's container.

### auditing

Manage auditing

- **`keen auditing`** - Returns a list of audit records. The list can be filtered to include items:

 *  where each item in `filter` has at least one match in any of these fields:
    
     *  `summary`
     *  `category`
     *  `eventSource`
     *  `objectItem.name` If the object is a user, account ID is available to filter.
     *  `objectItem.parentName`
     *  `objectItem.typeName`
     *  `changedValues.changedFrom`
     *  `changedValues.changedTo`
     *  `remoteAddress`
    
    For example, if `filter` contains *man ed*, an audit record containing `summary": "User added to group"` and `"category": "group management"` is returned.
 *  created on or after a date and time.
 *  created or or before a date and time.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### avatar

This resource represents system and custom avatars. Use it to obtain the details of system or custom avatars, add and remove avatars from a project, issue type or priority and obtain avatar images.


### bitbucket-user

The users resource allows you to access public information
associated with a user account. Most resources in the users
endpoint have been deprecated in favor of workspaces.

- **`keen bitbucket-user get-current`** - Returns the currently logged in user.
- **`keen bitbucket-user get-email-address-for-current`** - Returns details about a specific one of the authenticated user's
email addresses.

Details describe whether the address has been confirmed by the user and
whether it is the user's primary address or not.
- **`keen bitbucket-user get-permission-on-workspace`** - Returns the caller's effective role; as in, the highest level of privilege
the caller has for the workspace.
If the calling user is a member of multiple groups with distinct roles, only the
highest level is returned.

Permissions can be:

* `owner`
* `create-project`
* `collaborator` (deprecated; see this
[deprecation announcement](/cloud/bitbucket/deprecation-notice-collaborator-role/) for more details)
* `member`
- **`keen bitbucket-user list-email-addresses-for-current`** - Returns all the authenticated user's email addresses. Both
confirmed and unconfirmed.
- **`keen bitbucket-user list-repository-permissions-in-workspace-for`** - Returns an object for each repository the caller has explicit access to in the
specified workspace and their effective permission — the highest level of
permission the caller has. This does not return public repositories that the
user was not granted any specific permission in, and does not distinguish between
explicit and implicit privileges.

Permissions can be:

* `admin`
* `write`
* `read`

Results may be further [filtered or sorted](/cloud/bitbucket/rest/intro/#filtering) by
repository or permission by adding the following query string
parameters:

* `q=repository.name="bits"` or `q=permission>"read"`
* `sort=repository.name`

Note that the query parameter values need to be URL escaped so that `=`
would become `%3D`.
- **`keen bitbucket-user list-workspaces-for-current`** - Returns an object for each workspace accessible to the caller. This object
also contains details on whether the caller has admin permissions on the workspace
(`"administrator" = true`) or not (`"administrator" = false`).

Queries support filtering based on administrator permissions,
[sorting](/cloud/bitbucket/rest/intro/#sorting-query-results) or
[filtering](/cloud/bitbucket/rest/intro/#filtering) by `slug`. Results can
be [paginated](/cloud/bitbucket/rest/intro/#pagination).

### bitbucket-users

The users resource allows you to access public information
associated with a user account. Most resources in the users
endpoint have been deprecated in favor of workspaces.

- **`keen bitbucket-users <selected_user>`** - Gets the public information associated with a user account.

If the user's profile is private, `location`, `website` and
`created_on` elements are omitted.

Note that the user object returned by this operation is changing significantly, due to privacy changes.
See the [announcement](https://developer.atlassian.com/cloud/bitbucket/bitbucket-api-changes-gdpr/#changes-to-bitbucket-user-objects) for details.

### bitbucket-workspaces

A workspace is where you create repositories, collaborate on
your code, and organize different streams of work in your Bitbucket
Cloud account. Workspaces replace the use of teams and users in API
calls.

- **`keen bitbucket-workspaces <workspace>`** - Returns the requested workspace.

### blogposts

Manage blogposts

- **`keen blogposts create-blog-post`** - Creates a new blog post in the space specified by the spaceId.

By default this will create the blog post as a non-draft, unless the status is specified as draft.
If creating a non-draft, the title must not be empty.

Currently only supports the storage representation specified in the body.representation enums below
- **`keen blogposts delete-blog-post`** - Delete a blog post by id.

By default this will delete blog posts that are non-drafts. To delete a blog post that is a draft, the endpoint must be called on a 
draft with the following param `draft=true`. Discarded drafts are not sent to the trash and are permanently deleted.

Deleting a blog post that is not a draft moves the blog post to the trash, where it can be restored later.
To permanently delete a blog post (or "purge" it), the endpoint must be called on a **trashed** blog post with the following param `purge=true`.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
Permission to view the blog post and its corresponding space.
Permission to delete blog posts in the space.
[`manage/content`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space (if attempting to purge).

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen blogposts get-blog-post-by-id`** - Returns a specific blog post.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the blog post and its corresponding space.
- **`keen blogposts get-blog-posts`** - Returns all blog posts. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only blog posts that the user has permission to view will be returned.
- **`keen blogposts update-blog-post`** - Update a blog post by id.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the blog post and its corresponding space. Permission to update blog posts in the space.

### builds

APIs related to integrating build data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraBuildInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/build/. They are also related to integrating Jira Software Cloud with on-premises tools using OAuth 2.0 credentials. See https://developer.atlassian.com/cloud/jira/software/integrate-jsw-cloud-with-onpremises-tools/.

These APIs are available to Forge apps with the `devops:buildInfoProvider` module in the Forge app's manifest. See https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-software-build-info/.

- **`keen builds delete-by-key`** - Delete the build data currently stored for the given `pipelineId` and `buildNumber` combination.

Deletion is performed asynchronously. The `getBuildByKey` operation can be used to confirm that data has been
deleted successfully (if needed).
- **`keen builds delete-by-property`** - Bulk delete all builds data that match the given request.

One or more query params must be supplied to specify Properties to delete by.
Optional param `_updateSequenceNumber` is no longer supported.
If more than one Property is provided, data will be deleted that matches ALL of the
Properties (e.g. treated as an AND).

See the documentation for the `submitBuilds` operation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&repoId=repo-345

Deletion is performed asynchronously. The `getBuildByKey` operation can be used to confirm that data has been
deleted successfully (if needed).
- **`keen builds get-by-key`** - Retrieve the currently stored build data for the given `pipelineId` and `buildNumber` combination.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen builds submit`** - Update / insert builds data.

Builds are identified by the combination of `pipelineId` and `buildNumber`, and existing build data for the same
build will be replaced if it exists and the `updateSequenceNumber` of the existing data is less than the
incoming data.

Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are
available within a short period of time, but may take some time during peak load and/or maintenance times.
The `getBuildByKey` operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple builds being submitted in one request, each is validated individually prior to
submission. Details of which build failed submission (if any) are available in the response object.

### bulk

Manage bulk

- **`keen bulk get-available-transitions`** - Use this API to retrieve a list of transitions available for the specified issues that can be used or bulk transition operations. You can submit either single or multiple issues in the query to obtain the available transitions.

The response will provide the available transitions for issues, organized by their respective workflows. **Only the transitions that are common among the issues within that workflow and do not involve any additional field updates will be included.** For bulk transitions that require additional field updates, please utilise the Jira Cloud UI.

You can request available transitions for up to 1,000 issues in a single operation. This API uses pagination to return responses, delivering 50 workflows at a time.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Transition [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/permissions-for-company-managed-projects/#Transition-issues/) in all projects that contain the selected issues.
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk get-editable-fields`** - Use this API to get a list of fields visible to the user to perform bulk edit operations. You can pass single or multiple issues in the query to get eligible editable fields. This API uses pagination to return responses, delivering 50 fields at a time.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  Depending on the field, any field-specific permissions required to edit it.
- **`keen bulk get-operation-progress`** - Use this to get the progress state for the specified bulk operation `taskId`.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).

If the task is running, this resource will return:

    {"taskId":"10779","status":"RUNNING","progressPercent":65,"submittedBy":{"accountId":"5b10a2844c20165700ede21g"},"created":1690180055963,"started":1690180056206,"updated":169018005829}

If the task has completed, then this resource will return:

    {"processedAccessibleIssues":[10001,10002],"created":1709189449954,"progressPercent":100,"started":1709189450154,"status":"COMPLETE","submittedBy":{"accountId":"5b10a2844c20165700ede21g"},"invalidOrInaccessibleIssueCount":0,"taskId":"10000","totalIssueCount":2,"updated":1709189450354}

**Note:** You can view task progress for up to 14 days from creation.
- **`keen bulk submit-delete`** - Use this API to submit a bulk delete request. You can delete up to 1,000 issues in a single operation.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Delete [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/permissions-for-company-managed-projects/#Delete-issues/) in all projects that contain the selected issues.
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk submit-edit`** - Use this API to submit a bulk edit request and simultaneously edit multiple issues. There are limits applied to the number of issues and fields that can be edited. A single request can accommodate a maximum of 1000 issues (including subtasks) and 200 fields.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  Edit [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk submit-move`** - Use this API to submit a bulk issue move request. You can move multiple issues from multiple projects in a single request, but they must all be moved to a single project, issue type, and parent. You can't move more than 1000 issues (including subtasks) at once.

#### Scenarios: ####

This is an early version of the API and it doesn't have full feature parity with the Bulk Move UI experience.

 *  Moving issue of type A to issue of type B in the same project or a different project: `SUPPORTED`
 *  Moving multiple issues of type A in one or more projects to multiple issues of type B in one of the source projects or a different project: `SUPPORTED`
 *  Moving issues of multiple issue types in one or more projects to issues of a single issue type in one of the source project or a different project: **`SUPPORTED`**  
    E.g. Moving issues of story and task issue types in project 1 and project 2 to issues of task issue type in project 3
 *  Moving a standard parent issue of type A with its multiple subtask issue types in one project to standard issue of type B and multiple subtask issue types in the same project or a different project: `SUPPORTED`
 *  Moving standard issues with their subtasks to a parent issue in the same project or a different project without losing their relation: `SUPPORTED`
 *  Moving an epic issue with its child issues to a different project without losing their relation: `SUPPORTED`  
    This usecase is **supported using multiple requests**. Move the epic in one request and then move the children in a separate request with target parent set to the epic issue id  
      
    (Alternatively, move them individually and stitch the relationship back with the Bulk Edit API)

#### Limits applied to bulk issue moves: ####

When using the bulk move, keep in mind that there are limits on the number of issues and fields you can include.

 *  You can move up to 1,000 issues in a single operation, including any subtasks.
 *  The total combined number of fields across all issues must not exceed 1,500,000. For example, if each issue includes 15,000 fields, then the maximum number of issues that can be moved is 100.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Move [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in source projects.
 *  Create [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in destination projects.
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in destination projects, if moving subtasks only.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk submit-transition`** - Use this API to submit a bulk issue status transition request. You can transition multiple issues, alongside with their valid transition Ids. You can transition up to 1,000 issues in a single operation.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Transition [issues permission](https://support.atlassian.com/jira-cloud-administration/docs/permissions-for-company-managed-projects/#Transition-issues/) in all projects that contain the selected issues.
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk submit-unwatch`** - Use this API to submit a bulk unwatch request. You can unwatch up to 1,000 issues in a single operation.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen bulk submit-watch`** - Use this API to submit a bulk watch request. You can watch up to 1,000 issues in a single operation.

**[Permissions](#permissions) required:**

 *  Global bulk change [permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-global-permissions/).
 *  Browse [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) in all projects that contain the selected issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.

### changelog

Manage changelog

- **`keen changelog`** - Bulk fetch changelogs for multiple issues and filter by fields

Returns a paginated list of all changelogs for given issues sorted by changelog date and issue IDs, starting from the oldest changelog and smallest issue ID.

Issues are identified by their ID or key, and optionally changelogs can be filtered by their field IDs. You can request the changelogs of up to 1000 issues and can filter them by up to 10 field IDs.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the projects that the issues are in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issues.

### classification-levels

This resource represents classification levels.

- **`keen classification-levels`** - Returns all classification levels.

**[Permissions](#permissions) required:** None.

### comment

Manage comment

- **`keen comment`** - Returns a [paginated](#pagination) list of comments specified by a list of comment IDs.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Comments are returned where the user:

 *  has *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the comment.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  If the comment has visibility restrictions, belongs to the group or has the role visibility is restricted to.

### comments

Manage comments


### compass

Manage compass

- **`keen compass create-event`** - Sends an event into an event source. This API is rate limited. Only 100 requests per user per minute are allowed.
- **`keen compass delete-apispec`** - Delete the API spec file for a component.
- **`keen compass delete-attachment`** - Delete file attached to a component. If you need to have attachment upload enabled for your Compass Forge app, contact Atlassian support.
- **`keen compass delete-lock-file`** - Delete all the package dependencies associated with a particular source.The `sourceId` parameter needs to be Base64 encoded.
- **`keen compass evaluate-entitlements`** - Get entitlement results for a given user and site
- **`keen compass get-attachment`** - Get a file attached to a component. If you need to have attachment upload enabled for your Compass Forge app, contact Atlassian support.
- **`keen compass handle-webhook-invocation`** - Third parties can call these inbound webhooks to dothings within compass, like create deployment events for webhooks coming from github
- **`keen compass insert-metric-value`** - Sends a metric value into a metric source for a [component](/cloud/compass/components/what-is-a-component/). This API is rate limited. Only 100 requests per user per minute are allowed.
- **`keen compass upload-apispec`** - Upload api specs for a component.

The request body should be a file with the api-spec content.

Example: This curl command uploads a spec file (swagger.yaml) for the component.
```bash
curl --request PUT \
--url 'https://your-domain.atlassian.net/gateway/api/compass/v1/component/{componentId}/api_specs' \
--user 'email@example.com:<api_token>' \
--header 'Content-Type: multipart/form-data' \
--form file=@/path/to/file/swagger.yaml
```
- **`keen compass upload-attachment`** - Attach a file to a component. If you need to have attachment upload enabled for your Compass Forge app, contact Atlassian support.

The request body should be a file.

Example: This curl command attaches a file to a component.
```bash
curl --request PUT \
--url 'https://your-domain.atlassian.net/gateway/api/compass/v1/component/{componentId}/app/{forgeAppId}/attachment/{key}' \
--user 'email@example.com:<api_token>' \
--header 'Content-Type: multipart/form-data' \
--form file=@/path/to/attachment.txt
```
- **`keen compass upload-lock-file`** - Upload package dependencies for a component.

The file with package dependency information. Currently, we only accept package-lock.json files that are below 2MB. \
\
The baseSourceUrl should be a url link to your lock file, ideally with "#lines-" appended at the end, if from Bitbucket, or "#L" appended if a GitHub link. This will let link to specific lines in your lock file where dependencies were found. \
\
The sourceId should be a string that uniquely identifies the given lock file. We recommend using the filepath to the lock file, with the "/" characters replaced with "_" characters.

Example: This curl command uploads a lock file to populate package dependencies for the component.
```bash
curl --request PUT \
--url 'https://your-domain.atlassian.net/gateway/api/compass/v1/package_dependencies/lock_file' \
--user 'email@example.com:<api_token>' \
--header 'Content-Type: multipart/form-data' \
--form file=@/path/to/file/package-lock.json;type=application/json
--form componentId=<componentId>
--form baseSourceUrl=<url to your lock file>
--form sourceId=<a unique identifier for your file>
```

### component

Manage component

- **`keen component create`** - Creates a component. Use components to provide containers for issues within a project. Use components to provide containers for issues within a project.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project in which the component is created or *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen component delete`** - Deletes a component.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the component or *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen component find-for-projects`** - Returns a [paginated](#pagination) list of all components in a project, including global (Compass) components when applicable.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
- **`keen component get`** - Returns a component.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for project containing the component.
- **`keen component update`** - Updates a component. Any fields included in the request are overwritten. If `leadAccountId` is an empty string ("") the component lead is removed.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the component or *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### config

Manage config

- **`keen config associate-projects-to-field-association-schemes`** - Associate projects to field association schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config clone-field-association-scheme`** - Endpoint for cloning an existing field association scheme into a new one.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config create-field-association-scheme`** - Endpoint for creating a new field association scheme.

A new scheme is **not** copied from, or based on, any existing field association scheme. Instead, it is initialised with a minimal default set of critical fields sourced from the instance's own *system* and *product* fields (the fields returned by the product's field API), rather than from a scheme you specify.

To create a scheme that is based on an existing one, use the *Clone field scheme* endpoint instead.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config delete-field-association-scheme`** - Delete a specified field association scheme

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config get-field-association-scheme-by-id`** - Endpoint for fetching a field association scheme by its ID

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config get-field-association-scheme-item-parameters`** - Retrieve field association parameters on a field association scheme

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config get-field-association-schemes`** - REST endpoint for retrieving a paginated list of field association schemes with optional filtering.

This endpoint allows clients to fetch field association schemes with optional filtering by project IDs and text queries. The response includes scheme details with navigation links and filter metadata when applicable.

Filtering Behavior:

 *  When projectId or query parameters are provided, the response includes matchedFilters metadata showing which filters were applied.
 *  When no filters are applied, matchedFilters is omitted from individual scheme objects

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config get-projects-with-field-schemes`** - Get projects with field association schemes. This will be a temporary API but useful when transitioning from the legacy field configuration APIs to the new ones.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config remove-field-association-scheme-item-parameters`** - Remove field association parameters overrides for work types.

**Request body limits:** a single request can remove at most **100 parameters**. The number of removals is counted as the sum, across every removal entry, of the number of work types multiplied by the number of parameter names in that entry. Exceeding this limit returns a 400 response with the error "Request exceeds maximum number of parameters that can be removed 100".

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config remove-fields-associated-with-schemes`** - Remove fields associated with field association schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config search-field-association-scheme-fields`** - Search for fields belonging to a given field association scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config search-field-association-scheme-projects`** - REST Endpoint for searching for projects belonging to a given field association scheme

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config update-field-association-scheme`** - Endpoint for updating an existing field association scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config update-field-association-scheme-item-parameters`** - Update field association item parameters in field association schemes.

**Request body limits:** a single request can update parameters for at most **100 fields** (the top-level field ID keys of the request body). In addition, each field may include at most 50 update items, and each update item may reference at most 50 field schemes. Exceeding any of these limits returns a 400 response with an error such as "The fieldIds must contain no more than 100 items."

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen config update-fields-associated-with-schemes`** - Update fields associated with field association schemes.

**Request body limits:** a single request can update associations for at most **100 fields** (the top-level field ID keys of the request body), and each field association item may reference at most 50 field schemes. Exceeding a limit returns a 400 response with an error such as "The fieldIds must contain no more than 100 items."

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### configuration

Manage configuration

- **`keen configuration get`** - Returns the [global settings](https://confluence.atlassian.com/x/qYXKM) in Jira. These settings determine whether optional features (for example, subtasks, time tracking, and others) are enabled. If time tracking is enabled, this operation also returns the time tracking configuration.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen configuration get-available-time-tracking-implementations`** - Returns all time tracking providers. By default, Jira only has one time tracking provider: *JIRA provided time tracking*. However, you can install other time tracking providers via apps from the Atlassian Marketplace. For more information on time tracking providers, see the documentation for the [ Time Tracking Provider](https://developer.atlassian.com/cloud/jira/platform/modules/time-tracking-provider/) module.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen configuration get-selected-time-tracking-implementation`** - Returns the time tracking provider that is currently selected. Note that if time tracking is disabled, then a successful but empty response is returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen configuration get-shared-time-tracking`** - Returns the time tracking settings. This includes settings such as the time format, default time unit, and others. For more information, see [Configuring time tracking](https://confluence.atlassian.com/x/qoXKM).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen configuration select-time-tracking-implementation`** - Selects a time tracking provider.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen configuration set-shared-time-tracking`** - Sets the time tracking settings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### confluence-cloud-app

Manage confluence cloud app

- **`keen confluence-cloud-app delete-forge-property`** - Deletes a Forge app property. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestconfluence/#method-signature)** requests from Forge.
- **`keen confluence-cloud-app get-forge-properties`** - Gets Forge app properties. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestconfluence/#method-signature)** requests from Forge.
- **`keen confluence-cloud-app get-forge-property`** - Gets a Forge app property by property key. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestconfluence/#method-signature)** requests from Forge.
- **`keen confluence-cloud-app put-forge-property`** - Creates or updates a Forge app property. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestconfluence/#method-signature)** requests from Forge.

### confluence-cloud-classification-levels

Manage confluence cloud classification levels

- **`keen confluence-cloud-classification-levels`** - Returns a list of [classification levels](https://developer.atlassian.com/cloud/admin/dlp/rest/intro/#Classification%20level) 
available.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Permission to access the Confluence site ('Can use' global permission).

### confluence-cloud-user

Manage confluence cloud user

- **`keen confluence-cloud-user check-access-by-email`** - Returns the list of emails from the input list that do not have access to site.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen confluence-cloud-user invite-by-email`** - Invite a list of emails to the site.

Ignores all invalid emails and no action is taken for the emails that already have access to the site.

<b>NOTE:</b> This API is asynchronous and may take some time to complete.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).

### content

Manage content

- **`keen content`** - Converts a list of content ids into their associated content types. This is useful for users migrating from v1 to v2
who may have stored just content ids without their associated type. This will return types as they should be used in v2.
Notably, this will return `inline-comment` for inline comments and `footer-comment` for footer comments, which is distinct from them
both being represented by `comment` in v1.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the requested content. Any content that the user does not have permission to view or does not exist will map to `null` in the response.

### custom-content

Manage custom content

- **`keen custom-content create`** - Creates a new custom content in the given space, page, blogpost or other custom content.

Only one of `spaceId`, `pageId`, `blogPostId`, or `customContentId` is required in the request body.
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to create custom content in the space.
- **`keen custom-content delete`** - Delete a custom content by id.

Deleting a custom content will either move it to the trash or permanently delete it (purge it), depending on the apiSupport.
To permanently delete a **trashed** custom content, the endpoint must be called with the following param `purge=true`.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
Permission to view the content of the page or blogpost and its corresponding space.
Permission to delete custom content in the space.
[`manage/content`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space (if attempting to purge).

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen custom-content get-by-id`** - Returns a specific piece of custom content. 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the custom content, the container of the custom content, and the corresponding space (if different from the container).
- **`keen custom-content get-by-type`** - Returns all custom content for a given type. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the custom content, the container of the custom content, and the corresponding space (if different from the container).
- **`keen custom-content update`** - Update a custom content by id.
At most one of `spaceId`, `pageId`, `blogPostId`, or `customContentId` is allowed in the request body.
Note that if `spaceId` is specified, it must be the same as the `spaceId` used for creating the custom content
as moving custom content to a different space is not supported.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to update custom content in the space.

### custom-field-option

Manage custom field option

- **`keen custom-field-option <id>`** - Returns a custom field option. For example, an option in a select list.

Note that this operation **only works for issue field select list options created in Jira or using operations from the [Issue custom field options](#api-group-Issue-custom-field-options) resource**, it cannot be used with issue field select list options created by Connect apps.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** The custom field option is returned as follows:

 *  if the user has the *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
 *  if the user has the *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for at least one project the custom field is used in, and the field is visible in at least one layout the user has permission to view.

### dashboard

This resource represents dashboards. Use it to obtain the details of dashboards as well as get, create, update, or remove item properties and gadgets from dashboards.

- **`keen dashboard bulk-edit`** - Bulk edit dashboards. Maximum number of dashboards to be edited at the same time is 100.

**[Permissions](#permissions) required:** None

The dashboards to be updated must be owned by the user, or the user must be an administrator.
- **`keen dashboard create`** - Creates a dashboard.

**[Permissions](#permissions) required:** None.
- **`keen dashboard delete`** - Deletes a dashboard.

**[Permissions](#permissions) required:** None

The dashboard to be deleted must be owned by the user.
- **`keen dashboard get`** - Returns a dashboard.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.

However, to get a dashboard, the dashboard must be shared with the user or the user must own it. Note, users with the *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg) are considered owners of the System dashboard. The System dashboard is considered to be shared with all other users.
- **`keen dashboard get-all`** - Returns a list of dashboards owned by or shared with the user. The list may be filtered to include only favorite or owned dashboards.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen dashboard get-all-available-gadgets`** - Gets a list of all available gadgets that can be added to all dashboards.

**[Permissions](#permissions) required:** None.
- **`keen dashboard get-paginated`** - Returns a [paginated](#pagination) list of dashboards. This operation is similar to [Get dashboards](#api-rest-api-3-dashboard-get) except that the results can be refined to include dashboards that have specific attributes. For example, dashboards with a particular name. When multiple attributes are specified only filters matching all attributes are returned.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** The following dashboards that match the query parameters are returned:

 *  Dashboards owned by the user. Not returned for anonymous users.
 *  Dashboards shared with a group that the user is a member of. Not returned for anonymous users.
 *  Dashboards shared with a private project that the user can browse. Not returned for anonymous users.
 *  Dashboards shared with a public project.
 *  Dashboards shared with the public.
- **`keen dashboard update`** - Updates a dashboard, replacing all the dashboard details with those provided.

**[Permissions](#permissions) required:** None

The dashboard to be updated must be owned by the user.

### data-loss-prevention-classification-levels

Classification Level APIs"

- **`keen data-loss-prevention-classification-levels archive-level`** - -
Archives a [classification level](/cloud/admin/dlp/rest/intro/#classification%20level) with the supplied levelId (batch not currently supported).
When you archive a published classification level:
- Users won’t be able to classify pages at this level
- Any page or issues classified at this level will become unclassified
- Pages and issues will retain history

In the event that the classification level is restored and published, those pages and issues will regain this classification level.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels create-level`** - Creates a draft [classification level](cloud/admin/dlp/rest/intro/#classification%20level) for an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels edit-level`** - Edits a [classification level](/cloud/admin/dlp/rest/intro/#classification%20level) with the supplied levelId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels get-level`** - Gets a [classification level](/cloud/admin/dlp/rest/intro/#classification%20level) with the supplied levelId.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels get-level-list`** - Gets all [classification levels](/cloud/admin/dlp/rest/intro/#classification%20level) in an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels publish-level`** - Publishes one or more [classification level](/cloud/admin/dlp/rest/intro/#classification%20level) with the supplied levelIds.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels reorder`** - Changes the order of [classification levels](/cloud/admin/dlp/rest/intro/#classification%20level). The most sensitive classification level should be ranked 1.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`
- **`keen data-loss-prevention-classification-levels restore-level`** - Restores a [classification level](/cloud/admin/dlp/rest/intro/#classification%20level) with supplied levelId. When you restore an archived classification level, it’s restored as a draft. 

When you’re ready for your users to start classifying Confluence pages and Jira issues, you can publish the classification level.

If the classification level had been published before being archived, restored, and published again, any pages and issues classified at this level before it was archived will regain this classification level.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `write:classification-levels:admin`

### data-policies

Manage data policies

- **`keen data-policies get-data-policy-metadata`** - Returns data policy metadata for the workspace.

**[Permissions](#permissions) required:**
Only apps can make this request.
Permission to access the Confluence site ('Can use' global permission).
- **`keen data-policies get-data-policy-spaces`** - Returns all spaces. The results will be sorted by id ascending. The number of results is limited by the `limit` parameter and
additional results (if available) will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Only apps can make this request.
Permission to access the Confluence site ('Can use' global permission).
Only spaces that the app has permission to view will be returned.

### data-policy

Manage data policy

- **`keen data-policy get-policies`** - Returns data policies for the projects specified in the request.
- **`keen data-policy get-policy`** - Returns data policy for the workspace.

### databases

Manage databases

- **`keen databases create`** - Creates a database in the space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the corresponding space. Permission to create a database in the space.
- **`keen databases delete`** - Delete a database by id.

Deleting a database moves the database to the trash, where it can be restored later

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the database and its corresponding space.
Permission to delete databases in the space.
- **`keen databases get-by-id`** - Returns a specific database.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the database and its corresponding space.

### deployments

APIs related to integrating deployment data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraDeploymentInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/deployment/. They are also related to integrating Jira Software Cloud with on-premises tools using OAuth 2.0 credentials. See https://developer.atlassian.com/cloud/jira/software/integrate-jsw-cloud-with-onpremises-tools/.

These APIs are available to Forge apps with the `devops:deploymentInfoProvider` module in the Forge app's manifest. See https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-software-deployment-info/.

- **`keen deployments delete-by-key`** - Delete the currently stored deployment data for the given `pipelineId`, `environmentId` and `deploymentSequenceNumber` combination.

Deletion is performed asynchronously. The `getDeploymentByKey` operation can be used to confirm that data has been deleted successfully (if needed).
- **`keen deployments delete-by-property`** - Bulk delete all deployments that match the given request.

One or more query params must be supplied to specify the Properties to delete by. Optional param `_updateSequenceNumber` is no longer supported.
If more than one Property is provided, data will be deleted that matches ALL of the Properties (i.e. treated as AND).
See the documentation for the `submitDeployments` operation for more details.

Example operation: DELETE /bulkByProperties?accountId=account-123&createdBy=user-456

Deletion is performed asynchronously. The `getDeploymentByKey` operation can be used to confirm that data has been deleted successfully (if needed).
- **`keen deployments get-by-key`** - Retrieve the currently stored deployment data for the given `pipelineId`, `environmentId` and `deploymentSequenceNumber` combination.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen deployments get-gating-status-by-key`** - Retrieve the  Deployment gating status for the given `pipelineId + environmentId + deploymentSequenceNumber` combination.
Only apps that define the `jiraDeploymentInfoProvider` module can access this resource. This resource requires the 'READ' scope.
- **`keen deployments submit`** - Update / insert deployment data.

Deployments are identified by the combination of `pipelineId`, `environmentId` and `deploymentSequenceNumber`, and existing deployment data for the same deployment will be replaced if it exists and the `updateSequenceNumber` of existing data is less than the incoming data.

Submissions are processed asynchronously. Submitted data will eventually be available in Jira. Most updates are available within a short period of time, but may take some time during peak load and/or maintenance times. The `getDeploymentByKey` operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple deployments being submitted in one request, each is validated individually prior to submission. Details of which deployments failed submission (if any) are available in the response object.

### devinfo

Manage devinfo

- **`keen devinfo delete-by-properties`** - Deletes development information entities which have all the provided properties. Repositories which have properties that match ALL of the properties (i.e. treated as an AND), and all their related development information (such as commits, branches and pull requests), will be deleted. For example if request is `DELETE bulk?accountId=123&projectId=ABC` entities which have properties `accountId=123` and `projectId=ABC` will be deleted. Optional param `_updateSequenceId` is no longer supported. Deletion is performed asynchronously: specified entities will eventually be removed from Jira.
- **`keen devinfo delete-entity`** - Deletes particular development information entity. Deletion is performed asynchronously.
- **`keen devinfo delete-repository`** - Deletes the repository data stored by the given ID and all related development information entities. Deletion is performed asynchronously.
- **`keen devinfo exists-by-properties`** - Checks if repositories which have all the provided properties exists. For example, if request is `GET existsByProperties?accountId=123&projectId=ABC` then result will be positive only if there is at least one repository with both properties `accountId=123` and `projectId=ABC`. Special property `_updateSequenceId` can be used to filter all entities with updateSequenceId less or equal than the value specified. In addition to the optional `_updateSequenceId`, one or more query params must be supplied to specify properties to search by.
- **`keen devinfo get-repository`** - For the specified repository ID, retrieves the repository and the most recent 400 development information entities. The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen devinfo store-development-information`** - Stores development information provided in the request to make it available when viewing issues in Jira. Existing repository and entity data for the same ID will be replaced if the updateSequenceId of existing data is less than the incoming data. Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are available within a short period of time, but may take some time during peak load and/or maintenance times.

### devopscomponents

Manage devopscomponents

- **`keen devopscomponents delete-component-by-id`** - Delete the Component data currently stored for the given ID.

Deletion is performed asynchronously. The getComponentById operation can be used to confirm that data has been deleted successfully (if needed).

Only Connect apps that define the `jiraDevOpsComponentProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.
- **`keen devopscomponents delete-components-by-property`** - Bulk delete all Components that match the given request.

One or more query params must be supplied to specify Properties to delete by.
If more than one Property is provided, data will be deleted that matches ALL of the Properties (e.g. treated as an AND).
See the documentation for the submitComponents operation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&createdBy=user-456

Deletion is performed asynchronously. The getComponentById operation can be used to confirm that data has been deleted successfully (if needed).

Only Connect apps that define the `jiraDevOpsComponentProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.
- **`keen devopscomponents get-component-by-id`** - Retrieve the currently stored Component data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.

Only Connect apps that define the `jiraDevOpsComponentProvider` module can access this resource.
This resource requires the 'READ' scope for Connect apps.
- **`keen devopscomponents submit-components`** - Update / insert DevOps Component data.

Components are identified by their ID, and existing Component data for the same ID will be replaced if it exists and the updateSequenceNumber of existing data is less than the incoming data.

Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are available within a short period of time, but may take some time during peak load and/or maintenance times. The getComponentById operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple Components being submitted in one request, each is validated individually prior to submission. Details of which Components failed submission (if any) are available in the response object.

A maximum of 1000 components can be submitted in one request.

Only Connect apps that define the `jiraDevOpsComponentProvider` module can access this resource.
This resource requires the 'WRITE' scope for Connect apps.

### directories

Manage directories

- **`keen directories <orgId>`** - Returns a page of directories in an organization that match the supplied parameters.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:directories:admin`

### directory

Org Directory APIs

- **`keen directory <orgId> <accountId>`** - **Additional response parameters of the API (for e.g., `added_to_org`) are available only to customers using the new user management experience.** Learn more about the [new user management experience](https://community.atlassian.com/t5/Atlassian-Access-articles/User-management-for-cloud-admins-just-got-easier/ba-p/1576592).

Specifications:
- Return a user’s last active date for each product listed in Atlassian Administration.
- Active is defined as viewing a product's page for a minimum of 2 seconds. 
- The data for the last activity may be delayed by up to 24 hours.
- If the user has not accessed a product, the `product_access` response field will be empty. 

Learn the fastest way to call the API with a detailed [tutorial](https://developer.atlassian.com/cloud/admin/organization/user-last-active-dates/).

### domains

Domain APIs

- **`keen domains get`** - Returns a list of domains in an organization one page at a time.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:domains:admin`
- **`keen domains get-by-id`** - Returns information about a single verified domain by ID.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:domains:admin`

### embeds

Manage embeds

- **`keen embeds create-smart-link`** - Creates a Smart Link in the content tree in the space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the corresponding space. Permission to create a Smart Link in the content tree in the space.
- **`keen embeds delete-smart-link`** - Delete a Smart Link in the content tree by id.

Deleting a Smart Link in the content tree moves the Smart Link to the trash, where it can be restored later

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the Smart Link in the content tree and its corresponding space.
Permission to delete Smart Links in the content tree in the space.
- **`keen embeds get-smart-link-by-id`** - Returns a specific Smart Link in the content tree.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the Smart Link in the content tree and its corresponding space.

### event-actions

Manage event actions

- **`keen event-actions <orgId>`** - Returns information localized event actions

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:events:admin`

### events

Manage events

- **`keen events`** - Returns all issue events.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### events-stream

Manage events stream

- **`keen events-stream <orgId>`** - Returns a paginated list of audit logs events for an organization. Use this endpoint if you want to retrieve events in a simple, paginated manner with time-based filtering.

If you need more advanced filtering, refer to the [/events](https://developer.atlassian.com/cloud/admin/organization/rest/api-group-events/#api-v1-orgs-orgid-events-get) endpoint.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:events:admin`

### expression

Manage expression

- **`keen expression analyse`** - Analyses and validates Jira expressions.

As an experimental feature, this operation can also attempt to type-check the expressions.

Learn more about Jira expressions in the [documentation](https://developer.atlassian.com/cloud/jira/platform/jira-expressions/).

**[Permissions](#permissions) required**: None.
- **`keen expression evaluate-jira`** - Endpoint is currently being removed. [More details](https://developer.atlassian.com/changelog/#CHANGE-2046)

Evaluates a Jira expression and returns its value.

This resource can be used to test Jira expressions that you plan to use elsewhere, or to fetch data in a flexible way. Consult the [Jira expressions documentation](https://developer.atlassian.com/cloud/jira/platform/jira-expressions/) for more details.

#### Context variables ####

The following context variables are available to Jira expressions evaluated by this resource. Their presence depends on various factors; usually you need to manually request them in the context object sent in the payload, but some of them are added automatically under certain conditions.

 *  `user` ([User](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#user)): The current user. Always available and equal to `null` if the request is anonymous.
 *  `app` ([App](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#app)): The [Connect app](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) that made the request. Available only for authenticated requests made by Connect Apps (read more here: [Authentication for Connect apps](https://developer.atlassian.com/cloud/jira/platform/security-for-connect-apps/)).
 *  `issue` ([Issue](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue)): The current issue. Available only when the issue is provided in the request context object.
 *  `issues` ([List](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#list) of [Issues](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue)): A collection of issues matching a JQL query. Available only when JQL is provided in the request context object.
 *  `project` ([Project](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#project)): The current project. Available only when the project is provided in the request context object.
 *  `sprint` ([Sprint](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#sprint)): The current sprint. Available only when the sprint is provided in the request context object.
 *  `board` ([Board](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#board)): The current board. Available only when the board is provided in the request context object.
 *  `serviceDesk` ([ServiceDesk](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#servicedesk)): The current service desk. Available only when the service desk is provided in the request context object.
 *  `customerRequest` ([CustomerRequest](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#customerrequest)): The current customer request. Available only when the customer request is provided in the request context object.

Also, custom context variables can be passed in the request with their types. Those variables can be accessed by key in the Jira expression. These variable types are available for use in a custom context:

 *  `user`: A [user](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#user) specified as an Atlassian account ID.
 *  `issue`: An [issue](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue) specified by ID or key. All the fields of the issue object are available in the Jira expression.
 *  `json`: A JSON object containing custom content.
 *  `list`: A JSON list of `user`, `issue`, or `json` variable types.

This operation can be accessed anonymously.

**[Permissions](#permissions) required**: None. However, an expression may return different results for different users depending on their permissions. For example, different users may see different comments on the same issue.  
Permission to access Jira Software is required to access Jira Software context variables (`board` and `sprint`) or fields (for example, `issue.sprint`).
- **`keen expression evaluate-jsisjira`** - Evaluates a Jira expression and returns its value. The difference between this and `eval` is that this endpoint uses the enhanced search API when evaluating JQL queries. This API is eventually consistent, unlike the strongly consistent `eval` API. This allows for better performance and scalability. In addition, this API's response for JQL evaluation is based on a scrolling view (backed by a `nextPageToken`) instead of a paginated view (backed by `startAt` and `totalCount`).

This resource can be used to test Jira expressions that you plan to use elsewhere, or to fetch data in a flexible way. Consult the [Jira expressions documentation](https://developer.atlassian.com/cloud/jira/platform/jira-expressions/) for more details.

#### Context variables ####

The following context variables are available to Jira expressions evaluated by this resource. Their presence depends on various factors; usually you need to manually request them in the context object sent in the payload, but some of them are added automatically under certain conditions.

 *  `user` ([User](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#user)): The current user. Always available and equal to `null` if the request is anonymous.
 *  `app` ([App](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#app)): The [Connect app](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) that made the request. Available only for authenticated requests made by Connect apps (read more here: [Authentication for Connect apps](https://developer.atlassian.com/cloud/jira/platform/security-for-connect-apps/)).
 *  `issue` ([Issue](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue)): The current issue. Available only when the issue is provided in the request context object.
 *  `issues` ([List](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#list) of [Issues](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue)): A collection of issues matching a JQL query. Available only when JQL is provided in the request context object.
 *  `project` ([Project](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#project)): The current project. Available only when the project is provided in the request context object.
 *  `sprint` ([Sprint](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#sprint)): The current sprint. Available only when the sprint is provided in the request context object.
 *  `board` ([Board](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#board)): The current board. Available only when the board is provided in the request context object.
 *  `serviceDesk` ([ServiceDesk](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#servicedesk)): The current service desk. Available only when the service desk is provided in the request context object.
 *  `customerRequest` ([CustomerRequest](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#customerrequest)): The current customer request. Available only when the customer request is provided in the request context object.

In addition, you can pass custom context variables along with their types. You can then access them from the Jira expression by key. You can use the following variables in a custom context:

 *  `user`: A [user](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#user) specified as an Atlassian account ID.
 *  `issue`: An [issue](https://developer.atlassian.com/cloud/jira/platform/jira-expressions-type-reference#issue) specified by ID or key. All the fields of the issue object are available in the Jira expression.
 *  `json`: A JSON object containing custom content.
 *  `list`: A JSON list of `user`, `issue`, or `json` variable types.

This operation can be accessed anonymously.

**[Permissions](#permissions) required**: None. However, an expression may return different results for different users depending on their permissions. For example, different users may see different comments on the same issue.  
Permission to access Jira Software is required to access Jira Software context variables (`board` and `sprint`) or fields (for example, `issue.sprint`).

### featureflags

Manage featureflags

- **`keen featureflags delete-feature-flag-by-id`** - Delete the Feature Flag data currently stored for the given ID.

Deletion is performed asynchronously. The getFeatureFlagById operation can be used to confirm that data has been deleted successfully (if needed).
- **`keen featureflags delete-feature-flags-by-property`** - Bulk delete all Feature Flags that match the given request.

One or more query params must be supplied to specify Properties to delete by. Optional param `_updateSequenceId` is no longer supported.
If more than one Property is provided, data will be deleted that matches ALL of the Properties (e.g. treated as an AND).
See the documentation for the submitFeatureFlags operation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&createdBy=user-456

Deletion is performed asynchronously. The getFeatureFlagById operation can be used to confirm that data has been deleted successfully (if needed).
- **`keen featureflags get-feature-flag-by-id`** - Retrieve the currently stored Feature Flag data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen featureflags submit-feature-flags`** - Update / insert Feature Flag data.

Feature Flags are identified by their ID, and existing Feature Flag data for the same ID will be replaced if it exists and the updateSequenceId of existing data is less than the incoming data.

Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are available within a short period of time, but may take some time during peak load and/or maintenance times. The getFeatureFlagById operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple Feature Flags being submitted in one request, each is validated individually prior to submission. Details of which Feature Flags failed submission (if any) are available in the response object.

### field

Manage field

- **`keen field create-associations`** - Associates fields with projects.

Fields will be associated with each issue type on the requested projects.

Fields will be associated with all projects that share the same field configuration which the provided projects are using. This means that while the field will be associated with the requested projects, it will also be associated with any other projects that share the same field configuration.

If a success response is returned it means that the field association has been created in any applicable contexts where it wasn't already present.

Up to 50 fields and up to 100 projects can be associated in a single request. If more fields or projects are provided a 400 response will be returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen field create-custom`** - Creates a custom field.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen field delete-custom`** - Deletes a custom field. The custom field is deleted whether it is in the trash or not. See [Edit or delete a custom field](https://confluence.atlassian.com/x/Z44fOw) for more information on trashing and deleting custom fields.

This operation is [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen field get`** - Returns system and custom issue fields according to the following rules:

 *  Fields that cannot be added to the issue navigator are always returned.
 *  Fields that cannot be placed on an issue screen are always returned.
 *  Fields that depend on global Jira settings are only returned if the setting is enabled. That is, timetracking fields, subtasks, votes, and watches.
 *  Fields that are not associated to any used field configurations or screens are not returned.
 *  For all other fields, this operation only returns the fields that the user has permission to view (that is, the field is used in at least one project that the user has *Browse Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for.)

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen field get-paginated`** - Returns a [paginated](#pagination) list of fields for Classic Jira projects. The list can include:

 *  all fields
 *  specific fields, by defining `id`
 *  fields that contain a string in the field name or description, by defining `query`
 *  specific fields that contain a string in the field name or description, by defining `id` and `query`

Use `type` must be set to `custom` to show custom fields only.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen field get-trashed-paginated`** - Returns a [paginated](#pagination) list of fields in the trash. The list may be restricted to fields whose field name or description partially match a string.

Only custom fields can be queried, `type` must be set to `custom`.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen field remove-associations`** - Unassociates a set of fields with a project and issue type context.

Fields will be unassociated with all projects/issue types that share the same field configuration which the provided project and issue types are using. This means that while the field will be unassociated with the provided project and issue types, it will also be unassociated with any other projects and issue types that share the same field configuration.

If a success response is returned it means that the field association has been removed in any applicable contexts where it was present.

Up to 50 fields and up to 100 projects and issue types can be unassociated in a single request. If more fields or projects are provided a 400 response will be returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen field update-custom`** - Updates a custom field.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### fieldconfiguration

Manage fieldconfiguration

- **`keen fieldconfiguration create-field-configuration`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Creates a field configuration. The field configuration is created with the same field properties as the default configuration, with all the fields being optional.

This operation can only create configurations for use in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfiguration delete-field-configuration`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Deletes a field configuration.

This operation can only delete configurations used in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfiguration get-all-field-configurations`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Returns a [paginated](#pagination) list of field configurations. The list can be for all field configurations or a subset determined by any combination of these criteria:

 *  a list of field configuration item IDs.
 *  whether the field configuration is a default.
 *  whether the field configuration name or description contains a query string.

Only field configurations used in company-managed (classic) projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfiguration update-field-configuration`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Updates a field configuration. The name and the description provided in the request override the existing values.

This operation can only update configurations used in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### fieldconfigurationscheme

Manage fieldconfigurationscheme

- **`keen fieldconfigurationscheme assign-field-configuration-scheme-to-project`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Assigns a field configuration scheme to a project. If the field configuration scheme ID is `null`, the operation assigns the default field configuration scheme.

Field configuration schemes can only be assigned to classic projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme create-field-configuration-scheme`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Creates a field configuration scheme.

This operation can only create field configuration schemes used in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme delete-field-configuration-scheme`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Deletes a field configuration scheme.

This operation can only delete field configuration schemes used in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme get-all-field-configuration-schemes`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Returns a [paginated](#pagination) list of field configuration schemes.

Only field configuration schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme get-field-configuration-scheme-mappings`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Returns a [paginated](#pagination) list of field configuration issue type items.

Only items used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme get-field-configuration-scheme-project-mapping`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Returns a [paginated](#pagination) list of field configuration schemes and, for each scheme, a list of the projects that use it.

The list is sorted by field configuration scheme ID. The first item contains the list of project IDs assigned to the default field configuration scheme.

Only field configuration schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen fieldconfigurationscheme update-field-configuration-scheme`** - Deprecated, use [ Field schemes](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-field-schemes/#api-group-field-schemes) which supports field association schemes.

Updates a field configuration scheme.

This operation can only update field configuration schemes used in company-managed (classic) projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### filter

This resource represents [filters](https://confluence.atlassian.com/x/eQiiLQ). Use it to get, create, update, or delete filters. Also use it to configure the columns for a filter and set favorite filters.

- **`keen filter create`** - Creates a filter. The filter is shared according to the [default share scope](#api-rest-api-3-filter-post). The filter is not selected as a favorite.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen filter delete`** - Delete a filter.

**[Permissions](#permissions) required:** Permission to access Jira, however filters can only be deleted by the creator of the filter or a user with *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen filter get`** - Returns a filter.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None, however, the filter is only returned where it is:

 *  owned by the user.
 *  shared with a group that the user is a member of.
 *  shared with a private project that the user has *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for.
 *  shared with a public project.
 *  shared with the public.
- **`keen filter get-default-share-scope`** - Returns the default sharing settings for new filters and dashboards for a user.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen filter get-favourite`** - Returns the visible favorite filters of the user.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** A favorite filter is only visible to the user where the filter is:

 *  owned by the user.
 *  shared with a group that the user is a member of.
 *  shared with a private project that the user has *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for.
 *  shared with a public project.
 *  shared with the public.

For example, if the user favorites a public filter that is subsequently made private that filter is not returned by this operation.
- **`keen filter get-my`** - Returns the filters owned by the user. If `includeFavourites` is `true`, the user's visible favorite filters are also returned.

**[Permissions](#permissions) required:** Permission to access Jira, however, a favorite filters is only visible to the user where the filter is:

 *  owned by the user.
 *  shared with a group that the user is a member of.
 *  shared with a private project that the user has *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for.
 *  shared with a public project.
 *  shared with the public.

For example, if the user favorites a public filter that is subsequently made private that filter is not returned by this operation.
- **`keen filter get-paginated`** - Returns a [paginated](#pagination) list of filters. Use this operation to get:

 *  specific filters, by defining `id` only.
 *  filters that match all of the specified attributes. For example, all filters for a user with a particular word in their name. When multiple attributes are specified only filters matching all attributes are returned.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None, however, only the following filters that match the query parameters are returned:

 *  filters owned by the user.
 *  filters shared with a group that the user is a member of.
 *  filters shared with a private project that the user has *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for.
 *  filters shared with a public project.
 *  filters shared with the public.
- **`keen filter set-default-share-scope`** - Sets the default sharing for new filters and dashboards for a user.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen filter update`** - Updates a filter. Use this operation to update a filter's name, description, JQL, or sharing.

**[Permissions](#permissions) required:** Permission to access Jira, however the user must own the filter.

### folders

Manage folders

- **`keen folders create`** - Creates a folder in the space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the corresponding space. Permission to create a folder in the space.
- **`keen folders delete`** - Delete a folder by id.

Deleting a folder moves the folder to the trash, where it can be restored later

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the folder and its corresponding space.
Permission to delete folders in the space.
- **`keen folders get-by-id`** - Returns a specific folder.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the folder and its corresponding space.

### footer-comments

Manage footer comments

- **`keen footer-comments create`** - Create a footer comment.

The footer comment can be made against several locations: 
- at the top level (specifying pageId or blogPostId in the request body)
- as a reply (specifying parentCommentId in the request body)
- against an attachment (note: this is different than the comments added via the attachment properties page on the UI, which are referred to as version comments)
- against a custom content

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to create comments in the space.
- **`keen footer-comments delete`** - Deletes a footer comment. This is a permanent deletion and cannot be reverted.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to delete comments in the space.
- **`keen footer-comments get`** - Returns all footer comments. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the container and its corresponding space.
- **`keen footer-comments get-by-id`** - Retrieves a footer comment by id

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the container and its corresponding space.
- **`keen footer-comments update`** - Update a footer comment. This can be used to update the body text of a comment.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to create comments in the space.

### forge

Manage forge

- **`keen forge bulk-pin-unpin-projects-async`** - Bulk pin or unpin an issue panel (added by a Forge app) to or from multiple projects.

The operation runs asynchronously. The response includes a task ID - use the [Get task](#api-rest-api-3-task-taskId-get) endpoint to check progress.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen forge delete-app-property`** - Deletes a Forge app's property.

**[Permissions](#permissions) required:** Only Forge apps can make this request. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestjira/#method-signature)** requests from Forge.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen forge get-app-property`** - Returns the value of a Forge app's property.

**[Permissions](#permissions) required:** Only Forge apps can make this request. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestjira/#method-signature)** requests from Forge.
- **`keen forge get-app-property-keys`** - Returns all property keys for the Forge app.

**[Permissions](#permissions) required:** Only Forge apps can make this request. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestjira/#method-signature)** requests from Forge.
- **`keen forge get-bulk-pin-status`** - Get the pin status of an issue panel (added by a Forge app) for multiple projects.

The operation is read-only and runs synchronously. Projects that do not exist, or that you do not have permission to access, are returned in the response with the panel reported as not pinned and the reason in the `error` field; the request itself still succeeds.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen forge put-app-property`** - Sets the value of a Forge app's property.
These values can be retrieved in [Jira expressions](/cloud/jira/platform/jira-expressions/)
through the `app` [context variable](/cloud/jira/platform/jira-expressions/#context-variables).
They are also available in [entity property display conditions](/platform/forge/manifest-reference/display-conditions/entity-property-conditions/).

For other use cases, use the [Storage API](/platform/forge/runtime-reference/storage-api/).

The value of the request body must be a [valid](http://tools.ietf.org/html/rfc4627), non-empty JSON blob. The maximum length is 32768 characters.

**[Permissions](#permissions) required:** Only Forge apps can make this request. This API can only be accessed using **[asApp()](https://developer.atlassian.com/platform/forge/apis-reference/fetch-api-product.requestjira/#method-signature)** requests from Forge.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.

### global

Manage global

- **`keen global <id>`** - Update general configuration for object schema

### group

This resource represents groups of users. Use it to get, create, find, and delete groups as well as add and remove users from groups. (\[WARNING\] The standard Atlassian group names are default names only and can be edited or deleted. For example, an admin or Atlassian support could delete the default group jira-software-users or rename it to jsw-users at any point. See https://support.atlassian.com/user-management/docs/create-and-update-groups/ for details.)

- **`keen group add-user-to`** - Adds a user to a group.

**[Permissions](#permissions) required:** Site administration (that is, member of the *site-admin* [group](https://confluence.atlassian.com/x/24xjL)).
- **`keen group bulk-get`** - Returns a [paginated](#pagination) list of groups.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen group create`** - Creates a group.

**[Permissions](#permissions) required:** Site administration (that is, member of the *site-admin* [group](https://confluence.atlassian.com/x/24xjL)).
- **`keen group get`** - This operation is deprecated, use [`group/member`](#api-rest-api-3-group-member-get).

Returns all users in a group.

**[Permissions](#permissions) required:** either of:

 *  *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen group get-users-from`** - Returns a [paginated](#pagination) list of all users in a group.

Note that users are ordered by username, however the username is not returned in the results due to privacy reasons.

**[Permissions](#permissions) required:** either of:

 *  *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen group remove`** - Deletes a group.

**[Permissions](#permissions) required:** Site administration (that is, member of the *site-admin* strategic [group](https://confluence.atlassian.com/x/24xjL)).
- **`keen group remove-user-from`** - Removes a user from a group.

**[Permissions](#permissions) required:** Site administration (that is, member of the *site-admin* [group](https://confluence.atlassian.com/x/24xjL)).

### groups

This resource represents groups of users. Use it to get, create, find, and delete groups as well as add and remove users from groups. (\[WARNING\] The standard Atlassian group names are default names only and can be edited or deleted. For example, an admin or Atlassian support could delete the default group jira-software-users or rename it to jsw-users at any point. See https://support.atlassian.com/user-management/docs/create-and-update-groups/ for details.)

- **`keen groups`** - Returns a list of groups whose names contain a query string. A list of group names can be provided to exclude groups from the results.

The primary use case for this resource is to populate a group picker suggestions list. To this end, the returned object includes the `html` field where the matched query term is highlighted in the group name with the HTML strong tag. Also, the groups list is wrapped in a response object that contains a header for use in the picker, specifically *Showing X of Y matching groups*.

The list returns with the groups sorted. If no groups match the list criteria, an empty list is returned.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg). Anonymous calls and calls by users without the required permission return an empty list.

*Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg). Without this permission, calls where query is not an exact match to an existing group will return an empty list.

### groupuserpicker

Manage groupuserpicker

- **`keen groupuserpicker`** - Returns a list of users and groups matching a string. The string is used:

 *  for users, to find a case-insensitive match with display name and e-mail address. Note that if a user has hidden their email address in their user profile, partial matches of the email address will not find the user. An exact match is required.
 *  for groups, to find a case-sensitive match with group name.

For example, if the string *tin* is used, records with the display name *Tina*, email address *sarah@tinplatetraining.com*, and the group *accounting* would be returned.

Optionally, the search can be refined to:

 *  the projects and issue types associated with a custom field, such as a user picker. The search can then be further refined to return only users and groups that have permission to view specific:
    
     *  projects.
     *  issue types.
    
    If multiple projects or issue types are specified, they must be a subset of those enabled for the custom field or no results are returned. For example, if a field is enabled for projects A, B, and C then the search could be limited to projects B and C. However, if the search is limited to projects B and D, nothing is returned.
 *  not return Connect app users and groups.
 *  return groups that have a case-insensitive match with the query.

The primary use case for this resource is to populate a picker field suggestion list with users or groups. To this end, the returned object includes an `html` field for each list. This field highlights the matched query term in the item name with the HTML strong tag. Also, each list is wrapped in a response object that contains a header for use in a picker, specifically *Showing X of Y matching groups*.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/yodKLg).

### hook-events

Manage hook events

- **`keen hook-events list`** - Returns the webhook resource or subject types on which webhooks can be registered.
Each resource/subject type contains an `events` link that returns the paginated list of specific events each individual subject type can emit.
This endpoint is publicly accessible and does not require authentication or scopes.
- **`keen hook-events list-for-resource`** - Returns a paginated list of all valid webhook events for the specified entity.
**The team and user webhooks are deprecated, and you should use workspace instead. For more information, see [the announcement](https://developer.atlassian.com/cloud/bitbucket/bitbucket-api-teams-deprecation/).**
This is public data that does not require any scopes or authentication.
NOTE: The example response is a truncated response object for the `workspace` `subject_type`. We return the same structure for the other `subject_type` objects.

### icon

Manage icon

- **`keen icon find`** - Load a single icon by id
- **`keen icon find-global`** - Return all global icons i.e. icons not associated with a particular object schema

### importsource

Manage importsource

- **`keen importsource <id>`** - Retrieves a specific import source configuration by its ID. If scheduled imports are enabled, the response includes scheduling information.

### inline-comments

Manage inline comments

- **`keen inline-comments create`** - Create an inline comment. This can be at the top level (specifying pageId or blogPostId in the request body)
or as a reply (specifying parentCommentId in the request body). Note the inlineCommentProperties object in the
request body is used to select the text the inline comment should be tied to. This is what determines the text 
highlighting when viewing a page in Confluence.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to create comments in the space.
- **`keen inline-comments delete`** - Deletes an inline comment. This is a permanent deletion and cannot be reverted.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to delete comments in the space.
- **`keen inline-comments get`** - Returns all inline comments. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page and its corresponding space.
- **`keen inline-comments get-by-id`** - Retrieves an inline comment by id

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space.
- **`keen inline-comments update`** - Update an inline comment. This can be used to update the body text of a comment and/or to resolve the comment

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content of the page or blogpost and its corresponding space. Permission to create comments in the space.

### instance

Manage instance

- **`keen instance`** - Returns licensing information about the Jira instance.

**[Permissions](#permissions) required:** None.

### internal

Manage internal

- **`keen internal`** - Returns worklog details for a list of issue ID and worklog ID pairs.

This is an internal API for bulk fetching worklogs by their issue and worklog IDs. Worklogs that don't exist will be filtered out from the response.

The returned list of worklogs is limited to 1000 items.

**[Permissions](#permissions) required:** This is an internal service-to-service API that requires ASAP authentication. No user permission checks are performed as this bypasses normal user context.

### invalid-jqls-access-log

Manage invalid jqls access log

- **`keen invalid-jqls-access-log`** - Returns JQL searches that need updating because Jira can no longer support them. Use the entries to find affected saved filters and JQL used by apps, automations, scripts, or REST integrations.

Entries are ordered by creation time, newest first, and retained for three days. Pass the `nextPageToken` returned in a previous response to retrieve the next page. Results may be restricted to a user by external account ID.

**[Permissions](#permissions) required:** *Administer Jira* global permission.

### iql

Manage iql

- **`keen iql`** - This endpoint is deprecated and will be removed on <b>18th of September 2024</b>. Please use POST `/object/aql` instead.

### issue

This resource represents Jira issues. Use it to:

 *  create or edit issues, individually or in bulk.
 *  retrieve metadata about the options for creating or editing issues.
 *  delete an issue.
 *  assign a user to an issue.
 *  get issue changelogs.
 *  send notifications about an issue.
 *  get details of the transitions available for an issue.
 *  transition an issue.
 *  Archive issues.
 *  Unarchive issues.
 *  Export archived issues.

- **`keen issue archive`** - Enables admins to archive up to 1000 issues in a single request using issue ID/key, returning details of the issue(s) archived in the process and the errors encountered, if any.

**Note that:**

 *  you can't archive subtasks directly, only through their parent issues
 *  you can only archive issues from software, service management, and business projects

**[Permissions](#permissions) required:** Jira admin or site admin: [global permission](https://confluence.atlassian.com/x/x4dKLg)

**License required:** Premium or Enterprise

**Signed-in users only:** This API can't be accessed anonymously.
- **`keen issue archive-async`** - Enables admins to archive up to 100,000 issues in a single request using JQL, returning the URL to check the status of the submitted request.

You can use the [get task](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-tasks/#api-rest-api-3-task-taskid-get) and [cancel task](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-tasks/#api-rest-api-3-task-taskid-cancel-post) APIs to manage the request.

**Note that:**

 *  you can't archive subtasks directly, only through their parent issues
 *  you can only archive issues from software, service management, and business projects

**[Permissions](#permissions) required:** Jira admin or site admin: [global permission](https://confluence.atlassian.com/x/x4dKLg)

**License required:** Premium or Enterprise

**Signed-in users only:** This API can't be accessed anonymously.

**Rate limiting:** Only a single request per jira instance can be active at any given time.
- **`keen issue bulk-delete-property`** - Deletes a property value from multiple issues. The issues to be updated can be specified by filter criteria.

The criteria the filter used to identify eligible issues are:

 *  `entityIds` Only issues from this list are eligible.
 *  `currentValue` Only issues with the property set to this value are eligible.

If both criteria is specified, they are joined with the logical *AND*: only issues that satisfy both criteria are considered eligible.

If no filter criteria are specified, all the issues visible to the user and where the user has the EDIT\_ISSUES permission for the issue are considered eligible.

This operation is:

 *  transactional, either the property is deleted from all eligible issues or, when errors occur, no properties are deleted.
 *  [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:**

 *  *Browse projects* [ project permission](https://confluence.atlassian.com/x/yodKLg) for each project containing issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  *Edit issues* [project permission](https://confluence.atlassian.com/x/yodKLg) for each issue.
- **`keen issue bulk-fetch`** - Returns the details for a set of requested issues.

By default you can request up to 100 issues in a single call. You can request up to 1000 issues in a single call when the request is shaped so that it can be served efficiently, that is, when *all* of the following are true:

 *  the `fields` parameter explicitly names at least one field to include — a request that contains only exclusions is **not** eligible, and neither are the `*all` and `*navigable` wildcards or the default navigable field set, because the number of resolved fields depends on the site's configuration;
 *  no more than 100 fields are explicitly included;
 *  none of the included fields returns multiple values (for example `comment`, `worklog`, or `attachment`); and
 *  the `expand` parameter does not include `changelog`, `editmeta`, `operations`, `renderedFields`, `transitions`, or `versionedRepresentations`.

Requests that do not meet all of these conditions can include at most 100 issues; larger requests are rejected with a 400 error.

Each issue is identified by its ID or key, however, if the identifier doesn't match an issue, a case-insensitive search and check for moved issues is performed. If a matching issue is found its details are returned, a 302 or other redirect is **not** returned.

Issues will be returned in ascending `id` order. If there are errors, Jira will return a list of issues which couldn't be fetched along with error messages.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue bulk-set-properties-by`** - Sets or updates entity property values on issues. Up to 10 entity properties can be specified for each issue and up to 100 issues included in the request.

The value of the request body must be a [valid](http://tools.ietf.org/html/rfc4627), non-empty JSON.

This operation is:

 *  [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.
 *  non-transactional. Updating some entities may fail. Such information will available in the task result.

**[Permissions](#permissions) required:**

 *  *Browse projects* and *Edit issues* [project permissions](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue bulk-set-properties-list`** - Sets or updates a list of entity property values on issues. A list of up to 10 entity properties can be specified along with up to 10,000 issues on which to set or update that list of entity properties.

The value of the request body must be a [valid](http://tools.ietf.org/html/rfc4627), non-empty JSON. The maximum length of single issue property value is 32768 characters. This operation can be accessed anonymously.

This operation is:

 *  transactional, either all properties are updated in all eligible issues or, when errors occur, no properties are updated.
 *  [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:**

 *  *Browse projects* and *Edit issues* [project permissions](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue bulk-set-property`** - Sets a property value on multiple issues.

The value set can be a constant or determined by a [Jira expression](https://developer.atlassian.com/cloud/jira/platform/jira-expressions/). Expressions must be computable with constant complexity when applied to a set of issues. Expressions must also comply with the [restrictions](https://developer.atlassian.com/cloud/jira/platform/jira-expressions/#restrictions) that apply to all Jira expressions.

The issues to be updated can be specified by a filter.

The filter identifies issues eligible for update using these criteria:

 *  `entityIds` Only issues from this list are eligible.
 *  `currentValue` Only issues with the property set to this value are eligible.
 *  `hasProperty`:
    
     *  If *true*, only issues with the property are eligible.
     *  If *false*, only issues without the property are eligible.

If more than one criteria is specified, they are joined with the logical *AND*: only issues that satisfy all criteria are eligible.

If an invalid combination of criteria is provided, an error is returned. For example, specifying a `currentValue` and `hasProperty` as *false* would not match any issues (because without the property the property cannot have a value).

The filter is optional. Without the filter all the issues visible to the user and where the user has the EDIT\_ISSUES permission for the issue are considered eligible.

This operation is:

 *  transactional, either all eligible issues are updated or, when errors occur, none are updated.
 *  [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for each project containing issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  *Edit issues* [project permission](https://confluence.atlassian.com/x/yodKLg) for each issue.
- **`keen issue create`** - Creates an issue or, where the option to create subtasks is enabled in Jira, a subtask. A transition may be applied, to move the issue or subtask to a workflow step other than the default start step, and issue properties set.

The content of the issue or subtask is defined using `update` and `fields`. The fields that can be set in the issue or subtask are determined using the [ Get create issue metadata](#api-rest-api-3-issue-createmeta-get). These are the same fields that appear on the issue's create screen. Note that the `description`, `environment`, and any `textarea` type custom fields (multi-line text fields) take Atlassian Document Format content. Single line custom fields (`textfield`) accept a string and don't handle Atlassian Document Format content.

Creating a subtask differs from creating an issue as follows:

 *  `issueType` must be set to a subtask issue type (use [ Get create issue metadata](#api-rest-api-3-issue-createmeta-get) to find subtask issue types).
 *  `parent` must contain the ID or key of the parent issue.

In a next-gen project any issue may be made a child providing that the parent and child are members of the same project.

**[Permissions](#permissions) required:** *Browse projects* and *Create issues* [project permissions](https://confluence.atlassian.com/x/yodKLg) for the project in which the issue or subtask is created.
- **`keen issue create-bulk`** - Creates upto **50** issues and, where the option to create subtasks is enabled in Jira, subtasks. Transitions may be applied, to move the issues or subtasks to a workflow step other than the default start step, and issue properties set.

The content of each issue or subtask is defined using `update` and `fields`. The fields that can be set in the issue or subtask are determined using the [ Get create issue metadata](#api-rest-api-3-issue-createmeta-get). These are the same fields that appear on the issues' create screens. Note that the `description`, `environment`, and any `textarea` type custom fields (multi-line text fields) take Atlassian Document Format content. Single line custom fields (`textfield`) accept a string and don't handle Atlassian Document Format content.

Creating a subtask differs from creating an issue as follows:

 *  `issueType` must be set to a subtask issue type (use [ Get create issue metadata](#api-rest-api-3-issue-createmeta-get) to find subtask issue types).
 *  `parent` the must contain the ID or key of the parent issue.

**[Permissions](#permissions) required:** *Browse projects* and *Create issues* [project permissions](https://confluence.atlassian.com/x/yodKLg) for the project in which each issue or subtask is created.
- **`keen issue delete`** - Deletes an issue.

An issue cannot be deleted if it has one or more subtasks. To delete an issue with subtasks, set `deleteSubtasks`. This causes the issue's subtasks to be deleted with the issue.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse projects* and *Delete issues* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue edit`** - Edits an issue. Issue properties may be updated as part of the edit. Please note that issue transition is not supported and is ignored here. To transition an issue, please use [Transition issue](#api-rest-api-3-issue-issueIdOrKey-transitions-post).

The edits to the issue's fields are defined using `update` and `fields`. The fields that can be edited are determined using [ Get edit issue metadata](#api-rest-api-3-issue-issueIdOrKey-editmeta-get).

**Note:** This endpoint doesn't check screen configurations to determine if a field is editable. For more context, see the [Deprecation of override screen security](https://community.developer.atlassian.com/t/deprecation-of-override-screen-security/97153) announcement.

The parent field may be set by key or ID. For standard issue types, the parent may be removed by setting `update.parent.set.none` to *true*. Note that the `description`, `environment`, and any `textarea` type custom fields (multi-line text fields) take Atlassian Document Format content. Single line custom fields (`textfield`) accept a string and don't handle Atlassian Document Format content.

Connect apps having an app user with *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), and Forge apps acting on behalf of users with *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), can override the screen security configuration using `overrideScreenSecurity` and `overrideEditableFlag`.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse projects* and *Edit issues* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue get`** - Returns the details for an issue.

The issue is identified by its ID or key, however, if the identifier doesn't match an issue, a case-insensitive search and check for moved issues is performed. If a matching issue is found its details are returned, a 302 or other redirect is **not** returned. The issue key returned in the response is the key of the issue found.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue get-adf-limit-report`** - Returns all issues whose ADF (rich text) field data breaches the universal ADF size limit.

Unlike the issue limit report, which reports issues breaching per-issue entity *count* limits, this endpoint reports issues whose ADF field *byte size* exceeds that limit. The reported ADF field types are `comment_adf`, `worklog_adf`, `customfield_adf`, `description_adf` and `environment_adf`. The reported value for each issue is the number of breaching entities for that field (always 1 for the single-value description and environment fields).

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) is required for the project the issues are in. Results may be incomplete otherwise
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issue get-create-meta`** - Returns details of projects, issue types within projects, and, when requested, the create screen fields for each issue type for the user. Use the information to populate the requests in [ Create issue](#api-rest-api-3-issue-post) and [Create issues](#api-rest-api-3-issue-bulk-post).

Deprecated, see [Create Issue Meta Endpoint Deprecation Notice](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-1304).

The request can be restricted to specific projects or issue types using the query parameters. The response will contain information for the valid projects, issue types, or project and issue type combinations requested. Note that invalid project, issue type, or project and issue type combinations do not generate errors.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Create issues* [project permission](https://confluence.atlassian.com/x/yodKLg) in the requested projects.
- **`keen issue get-create-meta-type-id`** - Returns a page of field metadata for a specified project and issuetype id. Use the information to populate the requests in [ Create issue](#api-rest-api-3-issue-post) and [Create issues](#api-rest-api-3-issue-bulk-post).

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Create issues* [project permission](https://confluence.atlassian.com/x/yodKLg) in the requested projects.
- **`keen issue get-create-meta-types`** - Returns a page of issue type metadata for a specified project. Use the information to populate the requests in [ Create issue](#api-rest-api-3-issue-post) and [Create issues](#api-rest-api-3-issue-bulk-post).

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Create issues* [project permission](https://confluence.atlassian.com/x/yodKLg) in the requested projects.
- **`keen issue get-is-watching-bulk`** - Returns, for the user, details of the watched status of issues from a list. If an issue ID is invalid, the returned watched status is `false`.

This operation requires the **Allow users to watch issues** option to be *ON*. This option is set in General configuration for Jira. See [Configuring Jira application options](https://confluence.atlassian.com/x/uYXKM) for details.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen issue get-limit-report`** - Returns all issues breaching and approaching per-issue limits.

**[Permissions](#permissions) required:**

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) is required for the project the issues are in. Results may be incomplete otherwise
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issue get-picker-resource`** - Returns lists of issues matching a query string. Use this resource to provide auto-completion suggestions when the user is looking for an issue using a word or string.

This operation returns two lists:

 *  `History Search` which includes issues from the user's history of created, edited, or viewed issues that contain the string in the `query` parameter.
 *  `Current Search` which includes issues that match the JQL expression in `currentJQL` and contain the string in the `query` parameter.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen issue unarchive`** - Enables admins to unarchive up to 1000 issues in a single request using issue ID/key, returning details of the issue(s) unarchived in the process and the errors encountered, if any.

**Note that:**

 *  you can't unarchive subtasks directly, only through their parent issues
 *  you can only unarchive issues from software, service management, and business projects

**[Permissions](#permissions) required:** Jira admin or site admin: [global permission](https://confluence.atlassian.com/x/x4dKLg)

**License required:** Premium or Enterprise

**Signed-in users only:** This API can't be accessed anonymously.

### issue-link

This resource represents links between issues. Use it to get, create, and delete links between issues.

To use it, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

- **`keen issue-link delete`** - Deletes an issue link.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  Browse project [project permission](https://confluence.atlassian.com/x/yodKLg) for all the projects containing the issues in the link.
 *  *Link issues* [project permission](https://confluence.atlassian.com/x/yodKLg) for at least one of the projects containing issues in the link.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, permission to view both of the issues.
- **`keen issue-link get`** - Returns an issue link.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse project* [project permission](https://confluence.atlassian.com/x/yodKLg) for all the projects containing the linked issues.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, permission to view both of the issues.
- **`keen issue-link link-issues`** - Creates a link between two issues. Use this operation to indicate a relationship between two issues and optionally add a comment to the from (outward) issue. To use this resource the site must have [Issue Linking](https://confluence.atlassian.com/x/yoXKM) enabled.

This resource returns nothing on the creation of an issue link. To obtain the ID of the issue link, use `https://your-domain.atlassian.net/rest/api/3/issue/[linked issue key]?fields=issuelinks`.

If the link request duplicates a link, the response indicates that the issue link was created. If the request included a comment, the comment is added.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Browse project* [project permission](https://confluence.atlassian.com/x/yodKLg) for all the projects containing the issues to be linked,
 *  *Link issues* [project permission](https://confluence.atlassian.com/x/yodKLg) on the project containing the from (outward) issue,
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
 *  If the comment has visibility restrictions, belongs to the group or has the role visibility is restricted to.

### issue-link-type

This resource represents [issue link](#api-group-Issue-links) types. Use it to get, create, update, and delete link issue types as well as get lists of all link issue types.

To use it, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

- **`keen issue-link-type create`** - Creates an issue link type. Use this operation to create descriptions of the reasons why issues are linked. The issue link type consists of a name and descriptions for a link's inward and outward relationships.

To use this operation, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issue-link-type delete`** - Deletes an issue link type.

To use this operation, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issue-link-type get`** - Returns a list of all issue link types.

To use this operation, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for a project in the site.
- **`keen issue-link-type get-issuelinktype`** - Returns an issue link type.

To use this operation, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for a project in the site.
- **`keen issue-link-type update`** - Updates an issue link type.

To use this operation, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### issues

This resource represents Jira issues. Use it to:

 *  create or edit issues, individually or in bulk.
 *  retrieve metadata about the options for creating or editing issues.
 *  delete an issue.
 *  assign a user to an issue.
 *  get issue changelogs.
 *  send notifications about an issue.
 *  get details of the transitions available for an issue.
 *  transition an issue.
 *  Archive issues.
 *  Unarchive issues.
 *  Export archived issues.

- **`keen issues`** - Enables admins to retrieve details of all archived issues. Upon a successful request, the admin who submitted it will receive an email with a link to download a CSV file with the issue details.

Note that this API only exports the values of system fields and archival-specific fields (`ArchivedBy` and `ArchivedDate`). Custom fields aren't supported.

**[Permissions](#permissions) required:** Jira admin or site admin: [global permission](https://confluence.atlassian.com/x/x4dKLg)

**License required:** Premium or Enterprise

**Signed-in users only:** This API can't be accessed anonymously.

**Rate limiting:** Only a single request can be active at any given time.

### issuesecurityschemes

Manage issuesecurityschemes

- **`keen issuesecurityschemes associate-schemes-to-projects`** - Associates an issue security scheme with a project and remaps security levels of issues to the new levels, if provided.

This operation is [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes create-issue-security-scheme`** - Creates a security scheme with security scheme levels and levels' members. You can create up to 100 security scheme levels and security scheme levels' members per request.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes delete-security-scheme`** - Deletes an issue security scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes get-issue-security-scheme`** - Returns an issue security scheme along with its security levels.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
 *  *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for a project that uses the requested issue security scheme.
- **`keen issuesecurityschemes get-issue-security-schemes`** - Returns all [issue security schemes](https://confluence.atlassian.com/x/J4lKLg).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes get-security-level-members`** - Returns a [paginated](#pagination) list of issue security level members.

Only issue security level members in the context of classic projects are returned.

Filtering using parameters is inclusive: if you specify both security scheme IDs and level IDs, the result will include all issue security level members from the specified schemes and levels.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes get-security-levels`** - Returns a [paginated](#pagination) list of issue security levels.

Only issue security levels in the context of classic projects are returned.

Filtering using IDs is inclusive: if you specify both security scheme IDs and level IDs, the result will include both specified issue security levels and all issue security levels from the specified schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes search-projects-using-security-schemes`** - Returns a [paginated](#pagination) mapping of projects that are using security schemes. You can provide either one or multiple security scheme IDs or project IDs to filter by. If you don't provide any, this will return a list of all mappings. Only issue security schemes in the context of classic projects are supported. **[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes search-security-schemes`** - Returns a [paginated](#pagination) list of issue security schemes.  
If you specify the project ID parameter, the result will contain issue security schemes and related project IDs you filter by. Use \{@link IssueSecuritySchemeResource\#searchProjectsUsingSecuritySchemes(String, String, Set, Set)\} to obtain all projects related to scheme.

Only issue security schemes in the context of classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes set-default-levels`** - Sets default issue security levels for schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuesecurityschemes update-issue-security-scheme`** - Updates the issue security scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### issuetype

Manage issuetype

- **`keen issuetype create-issue-type`** - Creates an issue type.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetype delete-issue-type`** - Deletes the issue type. If the issue type is in use, all uses are updated with the alternative issue type (`alternativeIssueTypeId`). A list of alternative issue types are obtained from the [Get alternative issue types](#api-rest-api-3-issuetype-id-alternatives-get) resource.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetype get-issue-all-types`** - Returns all issue types.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issue types are only returned as follows:

 *  if the user has the *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), all issue types are returned.
 *  if the user has the *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for one or more projects, the issue types associated with the projects the user has permission to browse are returned.
 *  if the user is anonymous then they will be able to access projects with the *Browse projects* for anonymous users
 *  if the user authentication is incorrect they will fall back to anonymous
- **`keen issuetype get-issue-type`** - Returns an issue type.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) in a project the issue type is associated with or *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetype get-issue-types-for-project`** - Returns issue types for a project.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) in the relevant project or *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetype update-issue-type`** - Updates the issue type.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### issuetypescheme

Manage issuetypescheme

- **`keen issuetypescheme assign-issue-type-scheme-to-project`** - Assigns an issue type scheme to a project.

If any issues in the project are assigned issue types not present in the new scheme, the operation will fail. To complete the assignment those issues must be updated to use issue types in the new scheme.

Issue type schemes can only be assigned to classic projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme create-issue-type-scheme`** - Creates an issue type scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme delete-issue-type-scheme`** - Deletes an issue type scheme.

Only issue type schemes used in classic projects can be deleted. Only issue type schemes not associated with a project can be deleted

A validation error will be returned if the specified scheme is associated with one or more projects. Use [Get issue type scheme API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-type-schemes/#api-rest-api-3-issuetypescheme-get) (with the projects expand, and id query parameter) to get a list of projects. Then, use [Assign issue type scheme to project API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-type-schemes/#api-rest-api-3-issuetypescheme-project-put) to associate all projects to another scheme before deleting.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme get-all-issue-type-schemes`** - Returns a [paginated](#pagination) list of issue type schemes.

Only issue type schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme get-issue-type-scheme-for-projects`** - Returns a [paginated](#pagination) list of issue type schemes and, for each issue type scheme, a list of the projects that use it.

Only issue type schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme get-issue-type-schemes-mapping`** - Returns a [paginated](#pagination) list of issue type scheme items.

Only issue type scheme items used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescheme update-issue-type-scheme`** - Updates an issue type scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### issuetypescreenscheme

Manage issuetypescreenscheme

- **`keen issuetypescreenscheme assign-issue-type-screen-scheme-to-project`** - Assigns an issue type screen scheme to a project.

Issue type screen schemes can only be assigned to classic projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme create-issue-type-screen-scheme`** - Creates an issue type screen scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme delete-issue-type-screen-scheme`** - Deletes an issue type screen scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme get-issue-type-screen-scheme-mappings`** - Returns a [paginated](#pagination) list of issue type screen scheme items.

Only issue type screen schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme get-issue-type-screen-scheme-project-associations`** - Returns a [paginated](#pagination) list of issue type screen schemes and, for each issue type screen scheme, a list of the projects that use it.

Only issue type screen schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme get-issue-type-screen-schemes`** - Returns a [paginated](#pagination) list of issue type screen schemes.

Only issue type screen schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen issuetypescreenscheme update-issue-type-screen-scheme`** - Updates an issue type screen scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### jira-cloud-platform-search

Manage jira cloud platform search

- **`keen jira-cloud-platform-search and-reconsile-issues-using-jql`** - Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ). Recent updates might not be immediately visible in the returned search results. If you need [read-after-write](https://developer.atlassian.com/cloud/jira/platform/search-and-reconcile/) consistency, you can utilize the `reconcileIssues` parameter to ensure stronger consistency assurances. This operation can be accessed anonymously.

If the JQL query expression is too large to be encoded as a query parameter, use the [POST](#api-rest-api-3-search-post) version of this resource.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen jira-cloud-platform-search and-reconsile-issues-using-jql-post`** - Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ). Recent updates might not be immediately visible in the returned search results. If you need [read-after-write](https://developer.atlassian.com/cloud/jira/platform/search-and-reconcile/) consistency, you can utilize the `reconcileIssues` parameter to ensure stronger consistency assurances. This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen jira-cloud-platform-search count-issues`** - Provide an estimated count of the issues that match the [JQL](https://confluence.atlassian.com/x/egORLQ). Recent updates might not be immediately visible in the returned output. This endpoint requires JQL to be bounded.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen jira-cloud-platform-search for-issues-using-jql`** - Endpoint is currently being removed. [More details](https://developer.atlassian.com/changelog/#CHANGE-2046)

Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ).

If the JQL query expression is too large to be encoded as a query parameter, use the [POST](#api-rest-api-3-search-post) version of this resource.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen jira-cloud-platform-search for-issues-using-jql-post`** - Endpoint is currently being removed. [More details](https://developer.atlassian.com/changelog/#CHANGE-2046)

Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ).

There is a [GET](#api-rest-api-3-search-get) version of this resource that can be used for smaller JQL query expressions.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Issues are included in the response where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the issue.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.

### jira-cloud-platform-version

Manage jira cloud platform version

- **`keen jira-cloud-platform-version create`** - Creates a project version.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg) or *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project the version is added to.
- **`keen jira-cloud-platform-version delete`** - Deletes a project version.

Deprecated, use [ Delete and replace version](#api-rest-api-3-version-id-removeAndSwap-post) that supports swapping version values in custom fields, in addition to the swapping for `fixVersion` and `affectedVersion` provided in this resource.

Alternative versions can be provided to update issues that use the deleted version in `fixVersion` or `affectedVersion`. If alternatives are not provided, occurrences of `fixVersion` and `affectedVersion` that contain the deleted version are cleared.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg) or *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that contains the version.
- **`keen jira-cloud-platform-version get`** - Returns a project version.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project containing the version.
- **`keen jira-cloud-platform-version update`** - Updates a project version.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg) or *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that contains the version.

### jira-cloud-platform-workflow

Manage jira cloud platform workflow

- **`keen jira-cloud-platform-workflow delete-inactive`** - Deletes a workflow.

The workflow cannot be deleted if it is:

 *  an active workflow.
 *  a system workflow.
 *  associated with any workflow scheme.
 *  associated with any draft workflow scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen jira-cloud-platform-workflow delete-transition-rule-configurations`** - Deletes workflow transition rules from one or more workflows. These rule types are supported:

 *  [post functions](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-post-function/)
 *  [conditions](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-condition/)
 *  [validators](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-validator/)

Only rules created by the calling Connect app can be deleted.

**Note:** The `draft` parameter in the request body WorkflowId is deprecated and will be removed from this API on [November 2, 2026](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-3147).

**[Permissions](#permissions) required:** Only Connect apps can use this operation.
- **`keen jira-cloud-platform-workflow get-paginated`** - This will be removed on [June 1, 2026](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-2569); use [Search workflows](#api-rest-api-3-workflows-search-get) instead.

Returns a [paginated](#pagination) list of published classic workflows. When workflow names are specified, details of those workflows are returned. Otherwise, all published classic workflows are returned.

This operation does not return next-gen workflows.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen jira-cloud-platform-workflow get-transition-rule-configurations`** - Returns a [paginated](#pagination) list of workflows with transition rules. The workflows can be filtered to return only those containing workflow transition rules:

 *  of one or more transition rule types, such as [workflow post functions](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-post-function/).
 *  matching one or more transition rule keys.

Only workflows containing transition rules created by the calling [Connect](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) or [Forge](https://developer.atlassian.com/cloud/jira/platform/index/#forge-apps) app are returned.

Due to server-side optimizations, workflows with an empty list of rules may be returned; these workflows can be ignored.

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) or [Forge](https://developer.atlassian.com/cloud/jira/platform/index/#forge-apps) apps can use this operation.
- **`keen jira-cloud-platform-workflow list-history`** - Returns a list of workflow history entries for a specified workflow id.

**Note:** Stored workflow data expires after 60 days. Additionally, no data from before the 30th of October 2025 is available.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to access all, including project-scoped, workflows
 *  At least one of the *Administer projects* and *View (read-only) workflow* project permissions to access project-scoped workflows
- **`keen jira-cloud-platform-workflow read-from-history`** - Returns a workflow and related statuses for a specified workflow id and version number.

**Note:** Stored workflow data expires after 60 days. Additionally, no data from before the 30th of October 2025 is available.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to access all, including project-scoped, workflows
 *  At least one of the *Administer projects* and *View (read-only) workflow* project permissions to access project-scoped workflows
- **`keen jira-cloud-platform-workflow update-transition-rule-configurations`** - Updates configuration of workflow transition rules. The following rule types are supported:

 *  [post functions](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-post-function/)
 *  [conditions](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-condition/)
 *  [validators](https://developer.atlassian.com/cloud/jira/platform/modules/workflow-validator/)

Only rules created by the calling [Connect](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) or [Forge](https://developer.atlassian.com/cloud/jira/platform/index/#forge-apps) app can be updated.

To assist with app migration, this operation can be used to:

 *  Disable a rule.
 *  Add a `tag`. Use this to filter rules in the [Get workflow transition rule configurations](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-workflow-transition-rules/#api-rest-api-3-workflow-rule-config-get).

Rules are enabled if the `disabled` parameter is not provided.

**Note:** The `draft` parameter in the request body WorkflowId is deprecated and will be removed from this API on [November 2, 2026](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-3147).

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) or [Forge](https://developer.atlassian.com/cloud/jira/platform/index/#forge-apps) apps can use this operation.

### jql

This resource represents JQL search auto-complete details. Use it to obtain JQL search auto-complete data and suggestions for use in programmatic construction of queries or custom query builders. It also provides operations to:

 *  convert one or more JQL queries with user identifiers (username or user key) to equivalent JQL queries with account IDs.
 *  convert readable details in one or more JQL queries to IDs where a user doesn't have permission to view the entity whose details are readable.

- **`keen jql get-auto-complete`** - Returns reference data for JQL searches. This is a downloadable version of the documentation provided in [Advanced searching - fields reference](https://confluence.atlassian.com/x/gwORLQ) and [Advanced searching - functions reference](https://confluence.atlassian.com/x/hgORLQ), along with a list of JQL-reserved words. Use this information to assist with the programmatic creation of JQL queries or the validation of queries built in a custom query builder.

To filter visible field details by project or collapse non-unique fields by field type then [Get field reference data (POST)](#api-rest-api-3-jql-autocompletedata-post) can be used.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen jql get-auto-complete-post`** - Returns reference data for JQL searches. This is a downloadable version of the documentation provided in [Advanced searching - fields reference](https://confluence.atlassian.com/x/gwORLQ) and [Advanced searching - functions reference](https://confluence.atlassian.com/x/hgORLQ), along with a list of JQL-reserved words. Use this information to assist with the programmatic creation of JQL queries or the validation of queries built in a custom query builder.

This operation can filter the custom fields returned by project. Invalid project IDs in `projectIds` are ignored. System fields are always returned.

It can also return the collapsed field for custom fields. Collapsed fields enable searches to be performed across all fields with the same name and of the same field type. For example, the collapsed field `Component - Component[Dropdown]` enables dropdown fields `Component - cf[10061]` and `Component - cf[10062]` to be searched simultaneously.

**[Permissions](#permissions) required:** None.
- **`keen jql get-field-auto-complete-for-query-string`** - Returns the JQL search auto complete suggestions for a field.

Suggestions can be obtained by providing:

 *  `fieldName` to get a list of all values for the field.
 *  `fieldName` and `fieldValue` to get a list of values containing the text in `fieldValue`.
 *  `fieldName` and `predicateName` to get a list of all predicate values for the field.
 *  `fieldName`, `predicateName`, and `predicateValue` to get a list of predicate values containing the text in `predicateValue`.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen jql get-precomputations`** - Returns the list of a function's precomputations along with information about when they were created, updated, and last used. Each precomputation has a `value` \- the JQL fragment to replace the custom function clause with.

**[Permissions](#permissions) required:** This API is only accessible to apps and apps can only inspect their own functions.

The new `read:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen jql get-precomputations-by-id`** - Returns function precomputations by IDs, along with information about when they were created, updated, and last used. Each precomputation has a `value` \- the JQL fragment to replace the custom function clause with.

**[Permissions](#permissions) required:** This API is only accessible to apps and apps can only inspect their own functions.

The new `read:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen jql match-issues`** - Checks whether one or more issues would be returned by one or more JQL queries. Up to 10 JQL queries can be specified and up to 50 issue IDs included in the request.

**[Permissions](#permissions) required:** None, however, issues are only matched against JQL queries where the user has:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen jql migrate-queries`** - Converts one or more JQL queries with user identifiers (username or user key) to equivalent JQL queries with account IDs.

You may wish to use this operation if your system stores JQL queries and you want to make them GDPR-compliant. For more information about GDPR-related changes, see the [migration guide](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-user-privacy-api-migration-guide/).

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen jql parse-queries`** - Parses and validates JQL queries.

Validation is performed in context of the current user.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen jql sanitise-queries`** - Sanitizes one or more JQL queries by converting readable details into IDs where a user doesn't have permission to view the entity.

For example, if the query contains the clause *project = 'Secret project'*, and a user does not have browse permission for the project "Secret project", the sanitized query replaces the clause with *project = 12345"* (where 12345 is the ID of the project). If a user has the required permission, the clause is not sanitized. If the account ID is null, sanitizing is performed for an anonymous user.

Note that sanitization doesn't make the queries GDPR-compliant, because it doesn't remove user identifiers (username or user key). If you need to make queries GDPR-compliant, use [Convert user identifiers to account IDs in JQL queries](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-jql/#api-rest-api-3-jql-sanitize-post).

Before sanitization each JQL query is parsed. The queries are returned in the same order that they were passed.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen jql update-precomputations`** - Update the precomputation value of a function created by a Forge/Connect app.

**[Permissions](#permissions) required:** An API for apps to update their own precomputations.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.

### label

This resource represents available labels. Use it to get available labels for the global label field.

- **`keen label`** - Returns a [paginated](#pagination) list of labels.

### labels

Manage labels

- **`keen labels`** - Returns all labels. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only labels that the user has permission to view will be returned.

### license

Manage license

- **`keen license get-approximate-application-count`** - Returns the total approximate number of user accounts for a single Jira license. Note that this information is cached with a 7-day lifecycle and could be stale at the time of call.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen license get-approximate-count`** - Returns the approximate number of user accounts across all Jira licenses. Note that this information is cached with a 7-day lifecycle and could be stale at the time of call.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### license-types

Information about license types

- **`keen license-types get`** - Get a specific license type.
- **`keen license-types list`** - Get a list of license types.

### manage

Manage APIs

- **`keen manage activate-user`** - Activates the specified user account. The permission to make use of this resource is exposed by the
`lifecycle.enablement` privilege.

User accounts that were deactivated due to US export controls cannot be reactivated using this API. If you believe
the account was incorrectly blocked, please contact [Atlassian Support](https://support.atlassian.com/contact).

User accounts that have been deleted need the deletion to be canceled before reactivating.
- **`keen manage cancel-delete-account`** - This API will:
 - Cancel the scheduled deletion of the specified managed account.
 - Restore and activate the user’s account.
 
 Specifications:
 - You can cancel the deletion within the 14-day grace period of deleting a managed account. After that the account is permanently deleted.
 
 The permission to make use of this resource is exposed by the `lifecycle.delete` privilege. Learn more about [Get user management permissions API](https://developer.atlassian.com/cloud/admin/user-management/rest/api-group-manage/#api-users-account-id-manage-get) to manage the specified user.
- **`keen manage deactivate-user`** - Deactivate (block) the specified user account from logging into Atlassian. The permission to make use of this resource is exposed by the `lifecycle.enablement` privilege.
You can optionally set a message associated with the block. If none is supplied, a default message will be used.
- **`keen manage delete-account`** - This API will:
- Delete a managed account from Atlassian Administration.
- Withdraw complete access to all products and services listed in Atlassian Administration.
- Remove reference to the account from all lists under Directory in Atlassian Administration.

Specifications:
- Deleting an account is permanent. If you think you’ll need the account again, we recommend you [deactivate](https://support.atlassian.com/user-management/docs/deactivate-a-managed-account/) it instead.
- Before you permanently delete the account, you’ll have a 14-day grace period, during which the account will appear as temporarily deactivated.

Learn more about [deleting a managed account](https://support.atlassian.com/user-management/docs/delete-a-managed-account/).

Learn the fastest way to get the paramaters and delete account with a detailed [tutorial](https://developer.atlassian.com/cloud/admin/user-management/delete-managed-account/#delete-account). 

The permission to make use of this resource is exposed by the `lifecycle.delete` privilege. Learn more about [Get user management permissions API](https://developer.atlassian.com/cloud/admin/user-management/rest/api-group-manage/#api-users-account-id-manage-get) to manage the specified user.
- **`keen manage delete-api-token`** - Deletes a specifid API token by ID.
- **`keen manage get-api-tokens`** - Gets the API tokens owned by the specified user.
- **`keen manage get-profile`** - Returns information about a single Atlassian account by ID
- **`keen manage get-user-management-permissions`** - Returns the set of permissions you have for managing the specified Atlassian account
- **`keen manage set-email`** - Sets the specified user's email address. Before using this endpoint, you must [verify the target domain](https://confluence.atlassian.com/x/gjcWN) as the new email address will be considered verified.
The permission to make use of this resource is exposed by the `email.set` privilege.
This call invalidates all active sessions.
- **`keen manage update-profile`** - Updates fields in a user account. The `profile.write` privilege details which fields you can change.

### mypermissions

Manage mypermissions

- **`keen mypermissions`** - Returns a list of permissions indicating which permissions the user has. Details of the user's permissions can be obtained in a global, project, issue or comment context.

The user is reported as having a project permission:

 *  in the global context, if the user has the project permission in any project.
 *  for a project, where the project permission is determined using issue data, if the user meets the permission's criteria for any issue in the project. Otherwise, if the user has the project permission in the project.
 *  for an issue, where a project permission is determined using issue data, if the user has the permission in the issue. Otherwise, if the user has the project permission in the project containing the issue.
 *  for a comment, where the user has both the permission to browse the comment and the project permission for the comment's parent issue. Only the BROWSE\_PROJECTS permission is supported. If a `commentId` is provided whose `permissions` does not equal BROWSE\_PROJECTS, a 400 error will be returned.

This means that users may be shown as having an issue permission (such as EDIT\_ISSUES) in the global context or a project context but may not have the permission for any or all issues. For example, if Reporters have the EDIT\_ISSUES permission a user would be shown as having this permission in the global context or the context of a project, because any user can be a reporter. However, if they are not the user who reported the issue queried they would not have EDIT\_ISSUES permission for that issue.

For [Jira Service Management project permissions](https://support.atlassian.com/jira-cloud-administration/docs/customize-jira-service-management-permissions/), this will be evaluated similarly to a user in the customer portal. For example, if the BROWSE\_PROJECTS permission is granted to Service Project Customer - Portal Access, any users with access to the customer portal will have the BROWSE\_PROJECTS permission.

Global permissions are unaffected by context.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.

### mypreferences

Manage mypreferences

- **`keen mypreferences get-locale`** - Returns the locale for the user.

If the user has no language preference set (which is the default setting) or this resource is accessed anonymous, the browser locale detected by Jira is returned. Jira detects the browser locale using the *Accept-Language* header in the request. However, if this doesn't match a locale available Jira, the site default locale is returned.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen mypreferences get-preference`** - Returns the value of a preference of the current user.

Note that these keys are deprecated:

 *  *jira.user.locale* The locale of the user. By default this is not set and the user takes the locale of the instance.
 *  *jira.user.timezone* The time zone of the user. By default this is not set and the user takes the timezone of the instance.

These system preferences keys will be deprecated by 15/07/2024. You can still retrieve these keys, but it will not have any impact on Notification behaviour.

 *  *user.notifications.watcher* Whether the user gets notified when they are watcher.
 *  *user.notifications.assignee* Whether the user gets notified when they are assignee.
 *  *user.notifications.reporter* Whether the user gets notified when they are reporter.
 *  *user.notifications.mentions* Whether the user gets notified when they are mentions.

Use [ Update a user profile](https://developer.atlassian.com/cloud/admin/user-management/rest/#api-users-account-id-manage-profile-patch) from the user management REST API to manage timezone and locale instead.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen mypreferences remove-preference`** - Deletes a preference of the user, which restores the default value of system defined settings.

Note that these keys are deprecated:

 *  *jira.user.locale* The locale of the user. By default, not set. The user takes the instance locale.
 *  *jira.user.timezone* The time zone of the user. By default, not set. The user takes the instance timezone.

Use [ Update a user profile](https://developer.atlassian.com/cloud/admin/user-management/rest/#api-users-account-id-manage-profile-patch) from the user management REST API to manage timezone and locale instead.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen mypreferences set-locale`** - Deprecated, use [ Update a user profile](https://developer.atlassian.com/cloud/admin/user-management/rest/#api-users-account-id-manage-profile-patch) from the user management REST API instead.

Sets the locale of the user. The locale must be one supported by the instance of Jira.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen mypreferences set-preference`** - Creates a preference for the user or updates a preference's value by sending a plain text string. For example, `false`. An arbitrary preference can be created with the value containing up to 255 characters. In addition, the following keys define system preferences that can be set or created:

 *  *user.notifications.mimetype* The mime type used in notifications sent to the user. Defaults to `html`.
 *  *user.default.share.private* Whether new [ filters](https://confluence.atlassian.com/x/eQiiLQ) are set to private. Defaults to `true`.
 *  *user.keyboard.shortcuts.disabled* Whether keyboard shortcuts are disabled. Defaults to `false`.
 *  *user.autowatch.disabled* Whether the user automatically watches issues they create or add a comment to. By default, not set: the user takes the instance autowatch setting.
 *  *user.notifiy.own.changes* Whether the user gets notified of their own changes.

Note that these keys are deprecated:

 *  *jira.user.locale* The locale of the user. By default, not set. The user takes the instance locale.
 *  *jira.user.timezone* The time zone of the user. By default, not set. The user takes the instance timezone.

These system preferences keys will be deprecated by 15/07/2024. You can still use these keys to create arbitrary preferences, but it will not have any impact on Notification behaviour.

 *  *user.notifications.watcher* Whether the user gets notified when they are watcher.
 *  *user.notifications.assignee* Whether the user gets notified when they are assignee.
 *  *user.notifications.reporter* Whether the user gets notified when they are reporter.
 *  *user.notifications.mentions* Whether the user gets notified when they are mentions.

Use [ Update a user profile](https://developer.atlassian.com/cloud/admin/user-management/rest/#api-users-account-id-manage-profile-patch) from the user management REST API to manage timezone and locale instead.

**[Permissions](#permissions) required:** Permission to access Jira.

### myself

This resource represents information about the current user, such as basic details, group membership, application roles, preferences, and locale. Use it to get, create, update, and delete (restore default) values of the user's preferences and locale.

- **`keen myself`** - Returns details for the current user.

**[Permissions](#permissions) required:** Permission to access Jira.

### notificationscheme

Manage notificationscheme

- **`keen notificationscheme create-notification-scheme`** - Creates a notification scheme with notifications. You can create up to 1000 notifications per request.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen notificationscheme delete-notification-scheme`** - Deletes a notification scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen notificationscheme get-notification-scheme`** - Returns a [notification scheme](https://confluence.atlassian.com/x/8YdKLg), including the list of events and the recipients who will receive notifications for those events.

**[Permissions](#permissions) required:** Permission to access Jira, however, the user must have permission to administer at least one project associated with the notification scheme.
- **`keen notificationscheme get-notification-scheme-to-project-mappings`** - Returns a [paginated](#pagination) mapping of project that have notification scheme assigned. You can provide either one or multiple notification scheme IDs or project IDs to filter by. If you don't provide any, this will return a list of all mappings. Note that only company-managed (classic) projects are supported. This is because team-managed projects don't have a concept of a default notification scheme. The mappings are ordered by projectId.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen notificationscheme get-notification-schemes`** - Returns a [paginated](#pagination) list of [notification schemes](https://confluence.atlassian.com/x/8YdKLg) ordered by the display name.

*Note that you should allow for events without recipients to appear in responses.*

**[Permissions](#permissions) required:** Permission to access Jira, however, the user must have permission to administer at least one project associated with a notification scheme for it to be returned.
- **`keen notificationscheme update-notification-scheme`** - Updates a notification scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### oauth-clients

OAuth Client Management APIs

- **`keen oauth-clients count`** - Gets the count of OAuth clients in the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen oauth-clients create`** - Creates a new OAuth client for a service account in the specified organization.
Set `tokenEndpointAuthMethod` to `private_key_jwt` and supply `publicKeys` to
authenticate using a signed JWT. Private Key JWT must be enabled for your organization.
The response contains `clientId` and, when available, the supplied `keyId`;
`clientSecret` is null for `private_key_jwt` clients.

An active admin session with org management permissions is supported. API access requires OAuth client creation to be enabled for your organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `write:service-accounts-tokens:admin`
- **`keen oauth-clients delete`** - Deletes an existing OAuth client from the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `delete:service-accounts-tokens:admin`
- **`keen oauth-clients get`** - Retrieves a specific OAuth client by ID for the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen oauth-clients list`** - Retrieves all OAuth clients for the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen oauth-clients update`** - Updates the name, token endpoint authentication method, or public keys of an
OAuth client in the specified organization. Omitted fields remain unchanged.
Supply `publicKeys` to replace the configured key set, including when rotating keys.
When switching to `private_key_jwt`, supply the public keys to use for authentication.
Private Key JWT must be enabled for your organization.

The response includes the supplied `keyId` when available. When explicitly switching
from `private_key_jwt` to a client-secret method (`basic`, `post`, `basic_pkce`, or
`post_pkce`), the response can include `clientSecret`. Save it securely at this point;
GET responses do not return it. Other updates return a null `clientSecret`.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `write:service-accounts-tokens:admin`

### object

Manage object

- **`keen object by-aql`** - Fetch Objects by AQL
- **`keen object count-by-aql`** - This API provides the total count of objects that match a specified AQL query. Please note that this operation may incur performance latency.
- **`keen object create`** - Create a new object in Assets
- **`keen object delete`** - Delete the referenced object
- **`keen object find`** - Load one object
- **`keen object navigator-list`** - Retrieve a list of objects based on an AQL. Deprecated from <b>30 September 2024</b>. Please use POST /object/aql instead. For more information please see https://developer.atlassian.com/changelog/#CHANGE-1661.
- **`keen object navigator-list-deprecated`** - Deprecated. Use `/object/navlist/aql` instead.
- **`keen object update`** - Update an existing object in Assets

### objectconnectedtickets

Manage objectconnectedtickets


### objectschema

Manage objectschema

- **`keen objectschema schema-create`** - Create a new object schema
- **`keen objectschema schema-delete`** - Delete a schema
- **`keen objectschema schema-find`** - Find a schema by id
- **`keen objectschema schema-list`** - Resource to find object schemas in Assets
- **`keen objectschema schema-update`** - Update an object schema

### objecttype

Manage objecttype

- **`keen objecttype object-type-create`** - Create a new object type
- **`keen objecttype object-type-delete`** - Delete an object type
- **`keen objecttype object-type-find`** - Find an object type by id
- **`keen objecttype object-type-update`** - Update an existing object type

### objecttypeattribute

Manage objecttypeattribute

- **`keen objecttypeattribute object-type-attribute-create`** - Create a new attribute on the given object type
- **`keen objecttypeattribute object-type-attribute-delete`** - Delete an existing object type attribute
- **`keen objecttypeattribute object-type-attribute-update`** - Update an existing object type attribute

### operations

APIs related to integrating Incident and Post-Incident Review (PIR) data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraOperationsInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/operations-information/.

- **`keen operations delete-entity-by-property`** - Bulk delete all Entties that match the given request.

One or more query params must be supplied to specify Properties to delete by.
If more than one Property is provided, data will be deleted that matches ALL of the Properties (e.g. treated as an AND).
See the documentation for the submitEntity operation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&createdBy=user-456

Deletion is performed asynchronously. The getIncidentById operation can be used to confirm that data has been deleted successfully (if needed).

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.
- **`keen operations delete-incident-by-id`** - Delete the Incident data currently stored for the given ID.

Deletion is performed asynchronously. The getIncidentById operation can be used to confirm that data has been deleted successfully (if needed).

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.
- **`keen operations delete-review-by-id`** - Delete the Review data currently stored for the given ID.

Deletion is performed asynchronously. The getReviewById operation can be used to confirm that data has been deleted successfully (if needed).

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.
- **`keen operations delete-workspaces`** - Bulk delete all Operations Workspaces that match the given request.

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'DELETE' scope for Connect apps.

e.g. DELETE /bulk?workspaceIds=111-222-333,444-555-666
- **`keen operations get-incident-by-id`** - Retrieve the currently stored Incident data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'READ' scope for Connect apps.
- **`keen operations get-review-by-id`** - Retrieve the currently stored Review data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'READ' scope for Connect apps.
- **`keen operations get-workspaces`** - Retrieve the either all Operations Workspace IDs associated with the Jira site or a specific Operations Workspace ID for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.

e.g. GET /workspace?workspaceId=111-222-333

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'READ' scope for Connect apps.
- **`keen operations submit-entity`** - Update / insert Incident or Review data.

Incidents and reviews are identified by their ID, and existing Incident and Review data for the same ID will be replaced if it exists and the updateSequenceNumber of existing data is less than the incoming data.

Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are available within a short period of time, but may take some time during peak load and/or maintenance times. The getIncidentById or getReviewById operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple Incidents and Reviews being submitted in one request, each is validated individually prior to submission. Details of which entities failed submission (if any) are available in the response object.

A maximum of 1000 incidents can be submitted in one request.

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'WRITE' scope for Connect apps.
- **`keen operations submit-workspaces`** - Insert Operations Workspace IDs to establish a relationship between them and the Jira site the app is installed in. If a relationship between the Workspace ID and Jira already exists then the workspace ID will be ignored and Jira will process the rest of the entries.

Only Connect apps that define the `jiraOperationsInfoProvider` module can access this resource.
This resource requires the 'WRITE' scope for Connect apps.

### orgs

Orgs APIs

- **`keen orgs get`** - Returns a list of your organizations (based on your API key).
- **`keen orgs get-by-id`** - Returns information about a single organization by ID

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:orgs:admin`

### pages

Manage pages

- **`keen pages create`** - Creates a page in the space.

Pages are created as published by default unless specified as a draft in the status field. If creating a published page, the title must be specified.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the corresponding space. Permission to create a page in the space.
- **`keen pages delete`** - Delete a page by id.

By default this will delete pages that are non-drafts. To delete a page that is a draft, the endpoint must be called on a 
draft with the following param `draft=true`. Discarded drafts are not sent to the trash and are permanently deleted.

Deleting a page moves the page to the trash, where it can be restored later. To permanently delete a page (or "purge" it),
the endpoint must be called on a **trashed** page with the following param `purge=true`.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
Permission to view the page and its corresponding space.
Permission to delete pages in the space.
[`manage/content`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space (if attempting to purge).

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen pages get`** - Returns all pages. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only pages that the user has permission to view will be returned.
- **`keen pages get-by-id`** - Returns a specific page.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the page and its corresponding space.
- **`keen pages update`** - Update a page by id.

When the "current" version is updated, the provided body content is considered as the latest version. This latest body content
will be attempted to be merged into the draft version through a content reconciliation algorithm. If two versions are significantly diverged, 
the latest provided content may entirely override what was previously in the draft. 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the page and its corresponding space. Permission to update pages in the space.

### permissions

This resource represents permissions. Use it to obtain details of all permissions and determine whether the user has certain permissions.

- **`keen permissions get-all`** - Returns all permissions, including:

 *  global permissions.
 *  project permissions.
 *  global permissions added by plugins.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen permissions get-bulk`** - Returns:

 *  for a list of global permissions, the global permissions granted to a user.
 *  for a list of project permissions and lists of projects and issues, for each project permission a list of the projects and issues a user can access or manipulate.

If no account ID is provided, the operation returns details for the logged in user.

Note that:

 *  Invalid project and issue IDs are ignored.
 *  A maximum of 1000 projects and 1000 issues can be checked.
 *  Null values in `globalPermissions`, `projectPermissions`, `projectPermissions.projects`, and `projectPermissions.issues` are ignored.
 *  Empty strings in `projectPermissions.permissions` are ignored.

**Deprecation notice:** The required OAuth 2.0 scopes will be updated on June 15, 2024.

 *  **Classic**: `read:jira-work`
 *  **Granular**: `read:permission:jira`

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg) to check the permissions for other users, otherwise none. However, Connect apps can make a call from the app server to the product to obtain permission details for any user, without admin permission. This Connect app ability doesn't apply to calls made using AP.request() in a browser.
- **`keen permissions get-permitted-projects`** - Returns all the projects where the user is granted a list of project permissions.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.

### permissionscheme

Manage permissionscheme

- **`keen permissionscheme create-permission-scheme`** - Creates a new permission scheme. You can create a permission scheme with or without defining a set of permission grants.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen permissionscheme delete-permission-scheme`** - Deletes a permission scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen permissionscheme get-all-permission-schemes`** - Returns all permission schemes.

### About permission schemes and grants ###

A permission scheme is a collection of permission grants. A permission grant consists of a `holder` and a `permission`.

#### Holder object ####

The `holder` object contains information about the user or group being granted the permission. For example, the *Administer projects* permission is granted to a group named *Teams in space administrators*. In this case, the type is `"type": "group"`, and the parameter is the group name, `"parameter": "Teams in space administrators"` and the value is group ID, `"value": "ca85fac0-d974-40ca-a615-7af99c48d24f"`.

The `holder` object is defined by the following properties:

 *  `type` Identifies the user or group (see the list of types below).
 *  `parameter` As a group's name can change, use of `value` is recommended. The value of this property depends on the `type`. For example, if the `type` is a group, then you need to specify the group name.
 *  `value` The value of this property depends on the `type`. If the `type` is a group, then you need to specify the group ID. For other `type` it has the same value as `parameter`

The following `types` are available. The expected values for `parameter` and `value` are given in parentheses (some types may not have a `parameter` or `value`):

 *  `anyone` Grant for anonymous users.
 *  `applicationRole` Grant for users with access to the specified application (application name, application name). See [Update product access settings](https://confluence.atlassian.com/x/3YxjL) for more information.
 *  `assignee` Grant for the user currently assigned to an issue.
 *  `group` Grant for the specified group (`parameter` : group name, `value` : group ID).
 *  `groupCustomField` Grant for a user in the group selected in the specified custom field (`parameter` : custom field ID, `value` : custom field ID).
 *  `projectLead` Grant for a project lead.
 *  `projectRole` Grant for the specified project role (`parameter` :project role ID, `value` : project role ID).
 *  `reporter` Grant for the user who reported the issue.
 *  `sd.customer.portal.only` Jira Service Desk only. Grants customers permission to access the customer portal but not Jira. See [Customizing Jira Service Desk permissions](https://confluence.atlassian.com/x/24dKLg) for more information.
 *  `user` Grant for the specified user (`parameter` : user ID - historically this was the userkey but that is deprecated and the account ID should be used, `value` : user ID).
 *  `userCustomField` Grant for a user selected in the specified custom field (`parameter` : custom field ID, `value` : custom field ID).

#### Built-in permissions ####

The [built-in Jira permissions](https://confluence.atlassian.com/x/yodKLg) are listed below. Apps can also define custom permissions. See the [project permission](https://developer.atlassian.com/cloud/jira/platform/modules/project-permission/) and [global permission](https://developer.atlassian.com/cloud/jira/platform/modules/global-permission/) module documentation for more information.

**Administration permissions**

 *  `ADMINISTER_PROJECTS`
 *  `EDIT_WORKFLOW`
 *  `EDIT_ISSUE_LAYOUT`

**Project permissions**

 *  `BROWSE_PROJECTS`
 *  `MANAGE_SPRINTS_PERMISSION` (Jira Software only)
 *  `SERVICEDESK_AGENT` (Jira Service Desk only)
 *  `VIEW_DEV_TOOLS` (Jira Software only)
 *  `VIEW_READONLY_WORKFLOW`

**Issue permissions**

 *  `ASSIGNABLE_USER`
 *  `ASSIGN_ISSUES`
 *  `CLOSE_ISSUES`
 *  `CREATE_ISSUES`
 *  `DELETE_ISSUES`
 *  `EDIT_ISSUES`
 *  `LINK_ISSUES`
 *  `MODIFY_REPORTER`
 *  `MOVE_ISSUES`
 *  `RESOLVE_ISSUES`
 *  `SCHEDULE_ISSUES`
 *  `SET_ISSUE_SECURITY`
 *  `TRANSITION_ISSUES`

**Voters and watchers permissions**

 *  `MANAGE_WATCHERS`
 *  `VIEW_VOTERS_AND_WATCHERS`

**Comments permissions**

 *  `ADD_COMMENTS`
 *  `DELETE_ALL_COMMENTS`
 *  `DELETE_OWN_COMMENTS`
 *  `EDIT_ALL_COMMENTS`
 *  `EDIT_OWN_COMMENTS`

**Attachments permissions**

 *  `CREATE_ATTACHMENTS`
 *  `DELETE_ALL_ATTACHMENTS`
 *  `DELETE_OWN_ATTACHMENTS`

**Time tracking permissions**

 *  `DELETE_ALL_WORKLOGS`
 *  `DELETE_OWN_WORKLOGS`
 *  `EDIT_ALL_WORKLOGS`
 *  `EDIT_OWN_WORKLOGS`
 *  `WORK_ON_ISSUES`

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen permissionscheme get-permission-scheme`** - Returns a permission scheme.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen permissionscheme update-permission-scheme`** - Updates a permission scheme. Below are some important things to note when using this resource:

 *  If a permissions list is present in the request, then it is set in the permission scheme, overwriting *all existing* grants.
 *  If you want to update only the name and description, then do not send a permissions list in the request.
 *  Sending an empty list will remove all permission grants from the permission scheme.

If you want to add or delete a permission grant instead of updating the whole list, see [Create permission grant](#api-rest-api-3-permissionscheme-schemeId-permission-post) or [Delete permission scheme entity](#api-rest-api-3-permissionscheme-schemeId-permission-permissionId-delete).

See [About permission schemes and grants](../api-group-permission-schemes/#about-permission-schemes-and-grants) for more details.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### plans

This resource represents plans. Use it to get, create, duplicate, update, trash and archive plans.

- **`keen plans add-atlassian-team`** - Adds an existing Atlassian team to a plan and configures their plannning settings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans archive`** - Archives a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans create`** - Creates a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans create-only-team`** - Creates a plan-only team and configures their planning settings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans delete-only-team`** - Deletes a plan-only team and their planning settings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans duplicate`** - Duplicates a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans get`** - Returns a [paginated](#pagination) list of plans.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans get-atlassian-team`** - Returns planning settings for an Atlassian team in a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans get-only-team`** - Returns planning settings for a plan-only team.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans get-plan`** - Returns a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans get-teams`** - Returns a [paginated](#pagination) list of plan-only and Atlassian teams in a plan.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans remove-atlassian-team`** - Removes an Atlassian team from a plan and deletes their planning settings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans trash`** - Moves a plan to trash.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen plans update`** - Updates any of the following details of a plan using [JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902).

 *  name
 *  leadAccountId
 *  scheduling
    
     *  estimation with StoryPoints, Days or Hours as possible values
     *  startDate
        
         *  type with DueDate, TargetStartDate, TargetEndDate or DateCustomField as possible values
         *  dateCustomFieldId
     *  endDate
        
         *  type with DueDate, TargetStartDate, TargetEndDate or DateCustomField as possible values
         *  dateCustomFieldId
     *  inferredDates with None, SprintDates or ReleaseDates as possible values
     *  dependencies with Sequential or Concurrent as possible values
 *  issueSources
    
     *  type with Board, Project or Filter as possible values
     *  value
 *  exclusionRules
    
     *  numberOfDaysToShowCompletedIssues
     *  issueIds
     *  workStatusIds
     *  workStatusCategoryIds
     *  issueTypeIds
     *  releaseIds
 *  crossProjectReleases
    
     *  name
     *  releaseIds
 *  customFields
    
     *  customFieldId
     *  filter
 *  permissions
    
     *  type with View or Edit as possible values
     *  holder
        
         *  type with Group or AccountId as possible values
         *  value

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

*Note that "add" operations do not respect array indexes in target locations. Call the "Get plan" endpoint to find out the order of array elements.*
- **`keen plans update-atlassian-team`** - Updates any of the following planning settings of an Atlassian team in a plan using [JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902).

 *  planningStyle
 *  issueSourceId
 *  sprintLength
 *  capacity

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

*Note that "add" operations do not respect array indexes in target locations. Call the "Get Atlassian team in plan" endpoint to find out the order of array elements.*
- **`keen plans update-only-team`** - Updates any of the following planning settings of a plan-only team using [JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902).

 *  name
 *  planningStyle
 *  issueSourceId
 *  sprintLength
 *  capacity
 *  memberAccountIds

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

*Note that "add" operations do not respect array indexes in target locations. Call the "Get plan-only team" endpoint to find out the order of array elements.*

### policies

This resource represents policies. Use it to get, create, update, and delete policies, to apply them to and remove them from projects, and to inspect the policies and rules that apply to a project.

- **`keen policies get-policy-state`** - Returns the policies applied to a project.
- **`keen policies get-project-rules`** - Returns the configuration of the queried rule type IDs for a given project.

Use this to check whether specific rules are active on a project, and if so, how they are configured.

Only rules that are actually active on the project are returned. Rules that are not active are omitted from the response.
- **`keen policies list`** - Returns a paginated summary of all known policies, including their names, descriptions, and which projects each policy is assigned to.

When \{@code tanuki-358-cursor-policy-listing\} is enabled, omit cursor for the first page and pass the returned \{@code nextPageCursor\} for subsequent pages. The legacy \{@code startAt\} and \{@code total\} fields are null because an offset and total are not available in cursor mode. When disabled, the offset contract is unchanged.
- **`keen policies set-project`** - Sets the complete set of policies on a project (full replacement).

A project may have at most one policy assigned. Jira Admin only.
- **`keen policies validate-policy-application`** - Validates whether a set of policies can be applied to a project without actually applying them.

Runs the compliance checks for each specified policy against the project's current configuration. This can be called before applying policies to preview potential issues. Jira Admin only.

### priority

Manage priority

- **`keen priority create`** - Creates an issue priority.

**Deprecation notice:** The `iconUrl` parameter was sunset on 16th Mar 2025, and replaced with `avatarId`. See [CHANGE-1525](https://developer.atlassian.com/changelog/#CHANGE-1525).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priority delete`** - Deletes an issue priority.

This operation is [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priority get`** - Returns an issue priority. To fetch multiple priorities at once, use [Search priorities](#api-rest-api-3-priority-search-get) instead.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen priority get-priorities`** - Returns the list of all issue priorities.

**Deprecated:** Use [Search priorities](#api-rest-api-3-priority-search-get) instead. **[Permissions](#permissions) required:** Permission to access Jira.
- **`keen priority move-priorities`** - Changes the order of issue priorities.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priority search-priorities`** - Returns a [paginated](#pagination) list of priorities. The list can contain all priorities or a subset determined by any combination of these criteria:

 *  a list of priority IDs. Any invalid priority IDs are ignored.
 *  a list of project IDs. Only priorities that are available in these projects will be returned. Any invalid project IDs are ignored.
 *  whether the field configuration is a default. This returns priorities from company-managed (classic) projects only, as there is no concept of default priorities in team-managed projects.

**Deprecation notice:** The `onlyDefault` parameter is deprecated and will be removed at a later date. See [CHANGE-1655](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-1655).

**Deprecation notice:** The `isDefault` property of priorities is deprecated and will be removed at a later date. See [CHANGE-1655](https://developer.atlassian.com/cloud/jira/platform/changelog/#CHANGE-1655).

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen priority set-default`** - Sets default issue priority.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priority update`** - Updates an issue priority.

At least one request body parameter must be defined.

**Deprecation notice:** The `iconUrl` parameter was sunset on 16th Mar 2025, and replaced with `avatarId`. See [CHANGE-1525](https://developer.atlassian.com/changelog/#CHANGE-1525).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### priorityscheme

Manage priorityscheme

- **`keen priorityscheme create-priority-scheme`** - Creates a new priority scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priorityscheme delete-priority-scheme`** - Deletes a priority scheme.

This operation is only available for priority schemes without any associated projects. Any associated projects must be removed from the priority scheme before this operation can be performed.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priorityscheme get-available-priorities-by-priority-scheme`** - Returns a [paginated](#pagination) list of priorities available for adding to a priority scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priorityscheme get-priority-schemes`** - Returns a [paginated](#pagination) list of priority schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priorityscheme suggested-priorities-for-mappings`** - Returns a [paginated](#pagination) list of priorities that would require mapping, given a change in priorities or projects associated with a priority scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen priorityscheme update-priority-scheme`** - Updates a priority scheme. This includes its details, the lists of priorities and projects in it

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### products

Information about products

- **`keen products get`** - Get a specific product.
- **`keen products get-key`** - Get a list of versions for the specified product.
- **`keen products get-key-2`** - Get the latest version of the specified product.
- **`keen products get-key-3`** - Get details about a specific version, matching the specified build number, of the specified product.
- **`keen products get-key-4`** - Get details about a specific version, matching the specified name, of the specified product.
- **`keen products list`** - Get a list of products matching the specified parameters.

### progress

Manage progress

- **`keen progress <id>`** - Show ongoing import process

### project

This resource represents projects. Use it to get, create, update, and delete projects. Also get statuses available to a project, a project's notification schemes, and update a project's type.

- **`keen project create`** - Creates a project based on a project type template, as shown in the following table:

| Project Type Key | Project Template Key |  
|--|--|  
| `business` | `com.atlassian.jira-core-project-templates:jira-core-simplified-content-management`, `com.atlassian.jira-core-project-templates:jira-core-simplified-document-approval`, `com.atlassian.jira-core-project-templates:jira-core-simplified-lead-tracking`, `com.atlassian.jira-core-project-templates:jira-core-simplified-process-control`, `com.atlassian.jira-core-project-templates:jira-core-simplified-procurement`, `com.atlassian.jira-core-project-templates:jira-core-simplified-project-management`, `com.atlassian.jira-core-project-templates:jira-core-simplified-recruitment`, `com.atlassian.jira-core-project-templates:jira-core-simplified-task-tracking` |  
| `service_desk` | `com.atlassian.servicedesk:simplified-it-service-management`, `com.atlassian.servicedesk:simplified-external-service-desk`, `com.atlassian.servicedesk:simplified-hr-service-desk`, `com.atlassian.servicedesk:simplified-facilities-service-desk`, `com.atlassian.servicedesk:simplified-legal-service-desk`, `com.atlassian.servicedesk:simplified-analytics-service-desk`, `com.atlassian.servicedesk:simplified-marketing-service-desk`, `com.atlassian.servicedesk:simplified-design-service-desk`, `com.atlassian.servicedesk:simplified-sales-service-desk`, `com.atlassian.servicedesk:simplified-finance-service-desk`, `com.atlassian.servicedesk:company-managed-blank-service-project`, `com.atlassian.servicedesk:company-managed-general-service-project`, `com.atlassian.servicedesk:team-managed-general-service-project`, `com.atlassian.servicedesk:next-gen-it-service-desk`, `com.atlassian.servicedesk:next-gen-hr-service-desk`, `com.atlassian.servicedesk:next-gen-legal-service-desk`, `com.atlassian.servicedesk:next-gen-marketing-service-desk`, `com.atlassian.servicedesk:next-gen-facilities-service-desk`, `com.atlassian.servicedesk:next-gen-analytics-service-desk`, `com.atlassian.servicedesk:next-gen-finance-service-desk`, `com.atlassian.servicedesk:next-gen-design-service-desk`, `com.atlassian.servicedesk:next-gen-sales-service-desk` |  
| `software` | `com.pyxis.greenhopper.jira:gh-simplified-agility-kanban`, `com.pyxis.greenhopper.jira:gh-simplified-agility-scrum`, `com.pyxis.greenhopper.jira:gh-simplified-basic`, `com.pyxis.greenhopper.jira:gh-simplified-kanban-classic`, `com.pyxis.greenhopper.jira:gh-simplified-scrum-classic` |  
| `customer_service` | `com.atlassian.jcs:customer-service-management` |  
The project types are available according to the installed Jira features as follows:

 *  Jira Core, the default, enables `business` projects.
 *  Jira Service Management enables `service_desk` projects.
 *  Jira Software enables `software` projects.
 *  Jira Customer Service enables `customer_service` projects.

To determine which features are installed, go to **Jira settings** > **Apps** > **Manage apps** and review the System Apps list. To add Jira Software or Jira Service Management into a JIRA instance, use **Jira settings** > **Apps** > **Finding new apps**. For more information, see [ Managing add-ons](https://confluence.atlassian.com/x/S31NLg).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project delete`** - Deletes a project.

You can't delete a project if it's archived. To delete an archived project, restore the project and then delete it. To restore a project, use the Jira UI.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project get`** - Returns the [project details](https://confluence.atlassian.com/x/ahLpNw) for a project.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
- **`keen project get-accessible-type-by-key`** - Returns a [project type](https://confluence.atlassian.com/x/Var1Nw) if it is accessible to the user.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen project get-all`** - Returns all projects visible to the user. Deprecated, use [ Get projects paginated](#api-rest-api-3-project-search-get) that supports search and pagination.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Projects are returned only where the user has *Browse Projects* or *Administer projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
- **`keen project get-all-accessible-types`** - Returns all [project types](https://confluence.atlassian.com/x/Var1Nw) with a valid license.
- **`keen project get-all-types`** - Returns all [project types](https://confluence.atlassian.com/x/Var1Nw), whether or not the instance has a valid license for each type.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen project get-recent`** - Returns a list of up to 20 projects recently viewed by the user that are still visible to the user.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Projects are returned only where the user has one of:

 *  *Browse Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
 *  *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project get-type-by-key`** - Returns a [project type](https://confluence.atlassian.com/x/Var1Nw).

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen project search`** - Returns a [paginated](#pagination) list of projects visible to the user.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** Projects are returned only where the user has one of:

 *  *Browse Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
 *  *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project.
 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project update`** - Updates the [project details](https://confluence.atlassian.com/x/ahLpNw) of a project.

All parameters are optional in the body of the request. Schemes will only be updated if they are included in the request, any omitted schemes will be left unchanged.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg). is only needed when changing the schemes or project key. Otherwise you will only need *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg)

### project-category

Manage project category

- **`keen project-category create`** - Creates a project category.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project-category get-all-project-categories`** - Returns all project categories.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen project-category get-by-id`** - Returns a project category.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen project-category remove`** - Deletes a project category.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project-category update`** - Updates a project category.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### project-template

This resource represents project templates. Use it to create a new project from a custom template.

- **`keen project-template create-project-with-custom-template`** - Creates a project based on a custom template provided in the request.

The request body should contain the project details and the capabilities that comprise the project:

 *  `details` \- represents the project details settings
 *  `template` \- represents a list of capabilities responsible for creating specific parts of a project

A capability is defined as a unit of configuration for the project you want to create.

This operation is:

 *  [asynchronous](#async). Follow the `Location` link in the response header to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

***Note: This API is only supported for Jira Enterprise edition.***

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen project-template edit-template`** - Edit custom template

This API endpoint allows you to edit an existing customised template.

***Note: Custom Templates are only supported for Jira Enterprise edition.***
- **`keen project-template live-template`** - Get custom template

This API endpoint allows you to get a live custom project template details by either templateKey or projectId

***Note: Custom Templates are only supported for Jira Enterprise edition.***
- **`keen project-template remove-template`** - Remove custom template

This API endpoint allows you to remove a specified customised template

***Note: Custom Templates are only supported for Jira Enterprise edition.***
- **`keen project-template save-template`** - Save custom template

This API endpoint allows you to save a customised template

***Note: Custom Templates are only supported for Jira Enterprise edition.***

### projects

This resource represents projects. Use it to get, create, update, and delete projects. Also get statuses available to a project, a project's notification schemes, and update a project's type.

- **`keen projects`** - Returns a [paginated](#pagination) list of fields for the requested projects and work types.

Only fields that are available for the specified combination of projects and work types are returned. This endpoint allows filtering to specific fields if field IDs are provided.

**[Permissions](#permissions) required:** Permission to access Jira.

### projectvalidate

Manage projectvalidate

- **`keen projectvalidate get-valid-project-key`** - Validates a project key and, if the key is invalid or in use, generates a valid random string for the project key.

**[Permissions](#permissions) required:** None.
- **`keen projectvalidate get-valid-project-name`** - Checks that a project name isn't in use. If the name isn't in use, the passed string is returned. If the name is in use, this operation attempts to generate a valid project name based on the one supplied, usually by adding a sequence number. If a valid project name cannot be generated, a 404 response is returned.

**[Permissions](#permissions) required:** None.
- **`keen projectvalidate validate-project-key`** - Validates a project key by confirming the key is a valid string and not in use.

**[Permissions](#permissions) required:** None.

### redact

Manage redact

- **`keen redact get-redaction-status`** - Retrieves the current status of a redaction job ID.

The jobStatus will be one of the following:

 *  IN\_PROGRESS - The redaction job is currently in progress
 *  COMPLETED - The redaction job has completed successfully.
 *  PENDING - The redaction job has not started yet
- **`keen redact redact`** - Submit a job to redact issue field data. This will trigger the redaction of the data in the specified fields asynchronously.

The redaction status can be polled using the job id.

### remotelinks

Manage remotelinks

- **`keen remotelinks delete-remote-link-by-id`** - Delete the Remote Link data currently stored for the given ID.

Deletion is performed asynchronously. The `getRemoteLinkById` operation can be used to confirm that data has been
deleted successfully (if needed).
- **`keen remotelinks delete-remote-links-by-property`** - Bulk delete all Remote Links data that match the given request.

One or more query params must be supplied to specify Properties to delete by.
Optional param `_updateSequenceNumber` is no longer supported. If more than one Property is provided,
data will be deleted that matches ALL of the Properties (e.g. treated as an AND).

See the documentation for the `submitRemoteLinks` operation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&repoId=repo-345

Deletion is performed asynchronously. The `getRemoteLinkById` operation can be used to confirm that data has been
deleted successfully (if needed).
- **`keen remotelinks get-remote-link-by-id`** - Retrieve the currently stored Remote Link data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen remotelinks submit-remote-links`** - Update / insert Remote Link data.

Remote Links are identified by their ID, existing Remote Link data for the same ID will be replaced if it
exists and the updateSequenceId of existing data is less than the incoming data.

Submissions are performed asynchronously. Submitted data will eventually be available in Jira; most updates are
available within a short period of time, but may take some time during peak load and/or maintenance times.
The `getRemoteLinkById` operation can be used to confirm that data has been stored successfully (if needed).

In the case of multiple Remote Links being submitted in one request, each is validated individually prior to
submission. Details of which Remote LInk failed submission (if any) are available in the response object.

### repositories

A Git repository is a virtual storage of your project. It
allows you to save versions of your code, which you can access
when needed. The repo resource allows you to access public repos,
or repos that belong to a specific workspace.

- **`keen repositories create-repository`** - Creates a new repository.

Note: In order to set the project for the newly created repository,
pass in either the project key or the project UUID as part of the
request body as shown in the examples below:

```
$ curl -X POST -H "Content-Type: application/json" -d '{
    "scm": "git",
    "project": {
        "key": "MARS"
    }
}' https://api.bitbucket.org/2.0/repositories/teamsinspace/hablanding
```

or

```
$ curl -X POST -H "Content-Type: application/json" -d '{
    "scm": "git",
    "project": {
        "key": "{ba516952-992a-4c2d-acbd-17d502922f96}"
    }
}' https://api.bitbucket.org/2.0/repositories/teamsinspace/hablanding
```

The project must be assigned for all repositories. If the project is not provided,
the repository is automatically assigned to the oldest project in the workspace.

Note: In the examples above, the workspace ID `teamsinspace`,
and/or the repository name `hablanding` can be replaced by UUIDs.
- **`keen repositories delete-repository`** - Deletes the repository. This is an irreversible operation.

This does not affect its forks.
- **`keen repositories get-repository`** - Returns the object describing this repository.
- **`keen repositories list-in-workspace`** - Returns a paginated list of all repositories owned by the specified
workspace.

The result can be narrowed down based on the authenticated user's role.

E.g. with `?role=contributor`, only those repositories that the
authenticated user has write access to are returned (this includes any
repo the user is an admin on, as that implies write access).

This endpoint also supports filtering and sorting of the results. See
[filtering and sorting](/cloud/bitbucket/rest/intro/#filtering) for more details.
- **`keen repositories update-repository`** - Since this endpoint can be used to both update and to create a
repository, the request body depends on the intent.

#### Creation

See the POST documentation for the repository endpoint for an example
of the request body.

#### Update

Note: Changing the `name` of the repository will cause the location to
be changed. This is because the URL of the repo is derived from the
name (a process called slugification). In such a scenario, it is
possible for the request to fail if the newly created slug conflicts
with an existing repository's slug. But if there is no conflict,
the new location will be returned in the `Location` header of the
response.

### resolution

Manage resolution

- **`keen resolution create`** - Creates an issue resolution.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen resolution delete`** - Deletes an issue resolution.

This operation is [asynchronous](#async). Follow the `location` link in the response to determine the status of the task and use [Get task](#api-rest-api-3-task-taskId-get) to obtain subsequent updates.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen resolution get`** - Returns a list of all issue resolution values.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen resolution get-id`** - Returns an issue resolution value.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen resolution move`** - Changes the order of issue resolutions.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen resolution search`** - Returns a [paginated](#pagination) list of resolutions. The list can contain all resolutions or a subset determined by any combination of these criteria:

 *  a list of resolutions IDs.
 *  whether the field configuration is a default. This returns resolutions from company-managed (classic) projects only, as there is no concept of default resolutions in team-managed projects.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen resolution set-default`** - Sets default issue resolution.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen resolution update`** - Updates an issue resolution.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### role

Manage role

- **`keen role create-project`** - Creates a new project role with no [default actors](#api-rest-api-3-resolution-get). You can use the [Add default actors to project role](#api-rest-api-3-role-id-actors-post) operation to add default actors to the project role after creating it.

*Note that although a new project role is available to all projects upon creation, any default actors that are associated with the project role are not added to projects that existed prior to the role being created.*<

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen role delete-project`** - Deletes a project role. You must specify a replacement project role if you wish to delete a project role that is in use.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen role fully-update-project`** - Updates the project role's name and description. You must include both a name and a description in the request.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen role get-all-project`** - Gets a list of all project roles, complete with project role details and default actors.

### About project roles ###

[Project roles](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-roles/) are a flexible way to to associate users and groups with projects. In Jira Cloud, the list of project roles is shared globally with all projects, but each project can have a different set of actors associated with it (unlike groups, which have the same membership throughout all Jira applications).

Project roles are used in [permission schemes](#api-rest-api-3-permissionscheme-get), [email notification schemes](#api-rest-api-3-notificationscheme-get), [issue security levels](#api-rest-api-3-issuesecurityschemes-get), [comment visibility](#api-rest-api-3-comment-list-post), and workflow conditions.

#### Members and actors ####

In the Jira REST API, a member of a project role is called an *actor*. An *actor* is a group or user associated with a project role.

Actors may be set as [default members](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-roles/#Specifying-'default-members'-for-a-project-role) of the project role or set at the project level:

 *  Default actors: Users and groups that are assigned to the project role for all newly created projects. The default actors can be removed at the project level later if desired.
 *  Actors: Users and groups that are associated with a project role for a project, which may differ from the default actors. This enables you to assign a user to different roles in different projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen role get-project-by-id`** - Gets the project role details and the default actors associated with the role. The list of default actors is sorted by display name.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen role partial-update-project`** - Updates either the project role's name or its description.

You cannot update both the name and description at the same time using this operation. If you send a request with a name and a description only the name is updated.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### scim

Manage scim

- **`keen scim create-agroup-in-active-directory`** - Creates a read-only group in the organization's directory. You can only edit groups from your identity provider. 

**Note:** An attempt to create a group with an existing name will fail with a 409 (Conflict) error.
- **`keen scim create-auser-in-an-active-directory`** - Creates a user in the directory.
**Note:** An attempt to create an existing user will fail with a 409 (Conflict) error.

Use this API to manage accounts outside your organization when assigning these users to SCIM groups.

If there's already a managed Atlassian account associated with the specified email address on the Atlassian platform, the user in your identity provider will be connected or linked to the user in your Atlassian organization.
- **`keen scim delete-agroup`** - Deletes a group to remove the group from the organization's directory.
 
 **Note**: An attempt to delete a non-existent group will fail with a 404 (Resource Not found) error.

 **Note**: Deleting a synced group from your identity provider will delete the group from your organization's directory and associated sites. 
 1. If this group is used for allocating product license (granting role in a product), then members of this group may lose access to corresponding product after group deletion. 
 2. If this group is used to grant permissions in product, then members of this group may lose their permissions in the corresponding product.
- **`keen scim delete-auser-from-an-active-directory`** - Deleting a user via the SCIM APIs will unlink the user from your identity provider and deactivate the user within Atlassian if they are managed by your organization.

The deleted user is not available for future requests until created with a new `userId`. If the user is deactivated they can be activated again via [Atlassian Administration](https://admin.atlassian.com/).

**Note:** Executing this API call will result in the deletion of the SCIM record, and there is no method to reverse these changes except by creating a new SCIM record with [Create a user API](https://developer.atlassian.com/cloud/admin/user-provisioning/rest/api-group-users/#api-scim-directory-directoryid-users-post).
- **`keen scim get`** - Gets the details of a group based on the id.
- **`keen scim get-all-groups-from-an-active-directory`** - Get groups from the directory. Filter the groups by name supported with a single exact match (`eq`) against the `displayName` attribute. 

**Note**: While this API enables pagination, sorting functionality is not supported.
- **`keen scim get-auser-from-active-directory`** - Retrieves a user from the directory based on their `userId`.
- **`keen scim get-config`** - Get metadata about the supported SCIM features. This is a service provider configuration  endpoint providing supported SCIM features. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-extension-user-schemas`** - Get the user enterprise extension schemas from the SCIM provider. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-group-resource-type`** - Retrieves group resource type of this SCIM service provider. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-group-schemas`** - Get the group schemas from the SCIM provider. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-resource-types`** - Get different types of resources available on a SCIM service provider (e.g., Users and Groups).  
**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-schemas`** - Get all SCIM features metadata of your organization. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-user-resource-type`** - Retrieves user resource types from the SCIM service provider. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-user-schemas`** - Get the user schemas from the SCIM provider. 

**Note:** This API does not support filtering, pagination, or sorting.
- **`keen scim get-users-from-an-active-directory`** - Get users from the specified directory. Filtering is supported with a single exact match  (`eq`) against the `userName` and `externalId` attributes.

 **Note**: While this API enables pagination, sorting functionality is not supported.
- **`keen scim patch`** - Updates a group's information in the directory and manages group membership.

**Note:** Renaming groups after they've synced to your Atlassian organization isn't supported in this 
release of User Provisioning API. To rename a group, create a new group with the desired 
name, update membership, and then delete the old group.

#### Example

Some HTTP headers omitted and JSON payloads formatted for readability.

```
# Request
PATCH /scim/directory/2fb21891-7bee-4c2d-a61a-ade3834c8b2b/Groups/50202593-bc47-45df-8fa0-3f63343aa3c1 HTTP/1.1
Accept: application/scim+json
Accept-Charset: utf-8
Content-Type: application/scim+json; charset=utf-8
Authorization: Bearer 0j6lDgrjU7HmGagocgLe
Host: api.atlassian.com

{
   "schemas":[
      "urn:ietf:params:scim:api:messages:2.0:PatchOp"
   ],
   "Operations":[
      {
         "op":"add",
         "path":"members",
         "value":[
            {
               "value":"c6993c94-dbda-40f1-b6f0-18c855522ade",
               "display":"dave.meyer@demotime.authteam.com"
            },
            {
               "value":"f0ae48f7-1466-445e-85ea-e83ef754aefd",
               "display":"lingbo.lu@demotime.authteam.com"
            },
            {
               "value":"432d6f10-2e28-454e-be99-0f8c732a046f",
               "display":"joanna@demotime.authteam.com"
            }
         ]
      }
   ]
}

# Response
HTTP/1.1 200
Content-Type: application/scim+json

{
   "schemas":[
      "urn:ietf:params:scim:schemas:core:2.0:Group"
   ],
   "id":"50202593-bc47-45df-8fa0-3f63343aa3c1",
   "displayName":"demotime-confluence-users",
   "members":[
      {
         "type":"User",
         "value":"f0ae48f7-1466-445e-85ea-e83ef754aefd",
         "display":"lingbo.lu@demotime.authteam.com",
         "$ref":"https://api.atlassian.com/scim/directory/2fb21891-7bee-4c2d-a61a-ade3834c8b2b/Users/f0ae48f7-1466-445e-85ea-e83ef754aefd"
      },
      {
         "type":"User",
         "value":"c6993c94-dbda-40f1-b6f0-18c855522ade",
         "display":"dave.meyer@demotime.authteam.com",
         "$ref":"https://api.atlassian.com/scim/directory/2fb21891-7bee-4c2d-a61a-ade3834c8b2b/Users/c6993c94-dbda-40f1-b6f0-18c855522ade"
      },
      {
         "type":"User",
         "value":"432d6f10-2e28-454e-be99-0f8c732a046f",
         "display":"joanna@demotime.authteam.com",
         "$ref":"https://api.atlassian.com/scim/directory/2fb21891-7bee-4c2d-a61a-ade3834c8b2b/Users/432d6f10-2e28-454e-be99-0f8c732a046f"
      }
   ],
   "meta":{
      "resourceType":"Group",
      "location":"https://api.atlassian.com/scim/directory/2fb21891-7bee-4c2d-a61a-ade3834c8b2b/Groups/50202593-bc47-45df-8fa0-3f63343aa3c1",
      "lastModified":"2018-09-26T17:49:09.420654Z",
      "created":"2018-09-26T17:41:35.49073Z"
   }
}
```
- **`keen scim patch-user-information-in-an-active-directory`** - Updates a user's information in the directory based on their `userId` via `PATCH`. Refer to  [Service Provider Configuration APIs](https://developer.atlassian.com/cloud/admin/user-provisioning/rest/api-group-service-provider-configuration/#api-group-service-provider-configuration) for details on supported operations.
- **`keen scim put`** - Updates the details of a group with its unique ID.
- **`keen scim update-user-information-in-an-active-directory`** - Update the directory-based user information using the user attributes associated with their `userId`. User information  is replaced attribute-by-attribute, with the exception of immutable and read-only  attributes. Existing values of unspecified attributes are cleaned.

### screens

This resource represents the screens used to record issue details. Use it to:

 *  get details of all screens.
 *  get details of all the fields available for use on screens.
 *  create screens.
 *  delete screens.
 *  update screens.
 *  add a field to the default screen.

- **`keen screens add-field-to-default`** - Adds a field to the default tab of the default screen.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screens create`** - Creates a screen with a default field tab.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screens delete`** - Deletes a screen. A screen cannot be deleted if it is used in a screen scheme, workflow, or workflow draft.

Only screens used in classic projects can be deleted.
- **`keen screens get`** - Returns a [paginated](#pagination) list of all screens or those specified by one or more screen IDs.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screens get-bulk-tabs`** - Returns the list of tabs for a bulk of screens.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screens update`** - Updates a screen. Only screens used in classic projects can be updated.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### screenscheme

Manage screenscheme

- **`keen screenscheme create-screen-scheme`** - Creates a screen scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screenscheme delete-screen-scheme`** - Deletes a screen scheme. A screen scheme cannot be deleted if it is used in an issue type screen scheme.

Only screens schemes used in classic projects can be deleted.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screenscheme get-screen-schemes`** - Returns a [paginated](#pagination) list of screen schemes.

Only screen schemes used in classic projects are returned.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen screenscheme update-screen-scheme`** - Updates a screen scheme. Only screen schemes used in classic projects can be updated.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### security

Manage security

- **`keen security delete-linked-workspaces`** - Bulk delete all linked Security Workspaces that match the given request.

e.g. DELETE /bulk?workspaceIds=111-222-333,444-555-666
- **`keen security delete-vulnerabilities-by-property`** - Bulk delete all Vulnerabilities that match the given request.

One or more query params must be supplied to specify Properties to delete by.
If more than one Property is provided, data will be deleted that matches ALL of the Properties (e.g. treated as an AND).
Read the POST bulk endpoint documentation for more details.

e.g. DELETE /bulkByProperties?accountId=account-123&createdBy=user-456

Deletion is performed asynchronously. The GET vulnerability endpoint can be used to confirm that data has been deleted successfully (if needed).
- **`keen security delete-vulnerability-by-id`** - Delete the Vulnerability data currently stored for the given ID.

Deletion is performed asynchronously. The GET vulnerability endpoint can be used to confirm that data has been deleted successfully (if needed).
- **`keen security get-linked-workspace-by-id`** - Retrieve a specific Security Workspace linked to the Jira site for the given workspace ID.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen security get-linked-workspaces`** - Retrieve all Security Workspaces linked with the Jira site.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen security get-vulnerability-by-id`** - Retrieve the currently stored Vulnerability data for the given ID.

The result will be what is currently stored, ignoring any pending updates or deletes.
- **`keen security submit-vulnerabilities`** - Update / Insert Vulnerability data.

Vulnerabilities are identified by their ID, any existing Vulnerability data with the same ID will be replaced if it exists and the updateSequenceNumber of the existing data is less than the incoming data.

Submissions are performed asynchronously. Most updates are available within a short period of time but may take some time during peak load and/or maintenance times. The GET vulnerability endpoint can be used to confirm that data has been stored successfully (if needed).

In the case of multiple Vulnerabilities being submitted in one request, each is validated individually prior to submission. Details of Vulnerabilities that failed submission (if any) are available in the response object.

A maximum of 1000 vulnerabilities can be submitted in one request.
- **`keen security submit-workspaces`** - Insert Security Workspace IDs to establish a relationship between them and the Jira site the app is installed on. If a relationship between the workspace ID and Jira already exists then the workspace ID will be ignored and Jira will process the rest of the entries.

### securitylevel

Manage securitylevel

- **`keen securitylevel <id>`** - Returns details of an issue security level.

Use [Get issue security scheme](#api-rest-api-3-issuesecurityschemes-id-get) to obtain the IDs of issue security levels associated with the issue security scheme.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.

### server-info

This resource provides information about the Jira instance.

- **`keen server-info`** - Returns information about the Jira instance.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.

### service-accounts

Service Account lifecycle and credential APIs

- **`keen service-accounts count-api-tokens`** - Gets count of API tokens for specified service accounts within an organization.

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen service-accounts create`** - Creates a new service account for the specified organization. The service account is created with the given display name and optional description, then invited with the specified permission rules and optional additional group memberships.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `write:service-accounts:admin`
- **`keen service-accounts delete`** - Deletes a service account from the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `write:service-accounts:admin`
- **`keen service-accounts list`** - Retrieves all service accounts for the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts:admin`
- **`keen service-accounts list-api-tokens`** - Retrieves API tokens belonging to service accounts in the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen service-accounts list-oauth-clients`** - Retrieves OAuth clients belonging to service accounts in the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `read:service-accounts-tokens:admin`
- **`keen service-accounts update`** - Updates a service account's display name and/or description within the specified organization.

#### Scopes
**[OAuth 2.0 scopes](/cloud/admin/scopes/) required:** `write:service-accounts:admin`

### servicedeskapi

Manage servicedeskapi

- **`keen servicedeskapi add-customers`** - Adds one or more customers to a service desk. If any of the passed customers are associated with the service desk, no changes will be made for those customers and the resource returns a 204 success code.

**[Permissions](#permissions) required**: Service desk administrator
- **`keen servicedeskapi add-customers-servicedesk`** - Adds one or more customers to a service desk on behalf of jsd-nutmeg.

This endpoint is restricted to jsd-nutmeg via ASAP authentication. It provides the same capability as the public \{@code POST /servicedeskapi/servicedesk/\{serviceDeskId\}/customer\} endpoint, but does not require a User Context Token (UCT) or Connect app user \\u2014 authorization is enforced entirely via the ASAP token.

No user permission checks are performed; \{@code null\} is passed as the acting user to bypass the permission check in the underlying service.

If any of the passed customers are already associated with the service desk, no changes will be made for those customers and the resource returns a 204 success code.
- **`keen servicedeskapi add-organization`** - This method adds an organization to a service desk. If the organization ID is already associated with the service desk, no change is made and the resource returns a 204 success code.

**[Permissions](#permissions) required**: Service desk's agent.
- **`keen servicedeskapi add-request-participants`** - This method adds participants to a customer request.

**[Permissions](#permissions) required**: Permission to manage participants on the customer request.

Note, participants can be added when creating a customer request using the [request](#api-request-post) resource, by defining the participants in the `requestParticipants` field.
- **`keen servicedeskapi add-users-to-organization`** - This method adds users to an organization.

**[Permissions](#permissions) required**: Service desk administrator or agent. Note: Permission to add users to an organization can be switched to users with the Jira administrator permission, using the **[Organization management](https://confluence.atlassian.com/servicedeskcloud/setting-up-service-desk-users-732528877.html#Settingupservicedeskusers-manageorgsManageorganizations)** feature.
- **`keen servicedeskapi answer-approval`** - This method enables a user to **Approve** or **Decline** an approval on a customer request. The approval is assumed to be owned by the user making the call.

**[Permissions](#permissions) required**: User is assigned to the approval request.
- **`keen servicedeskapi attach-temporary-file`** - This method adds one or more temporary attachments to a service desk, which can then be permanently attached to a customer request using [servicedeskapi/request/\{issueIdOrKey\}/attachment](#api-request-issueIdOrKey-attachment-post).

**Note**: It is possible for a service desk administrator to turn off the ability to add attachments to a service desk.

This method expects a multipart request. The media-type multipart/form-data is defined in RFC 1867. Most client libraries have classes that make dealing with multipart posts simple. For instance, in Java the Apache HTTP Components library provides [MultiPartEntity](http://hc.apache.org/httpcomponents-client-ga/httpmime/apidocs/org/apache/http/entity/mime/MultipartEntity.html).

Because this method accepts multipart/form-data, it has XSRF protection on it. This means you must submit a header of X-Atlassian-Token: no-check with the request or it will be blocked.

The name of the multipart/form-data parameter that contains the attachments must be `file`.

For example, to upload a file called `myfile.txt` in the Service Desk with ID 10001 use

    curl -D- -u customer:customer -X POST -H "X-ExperimentalApi: opt-in" -H "X-Atlassian-Token: no-check" -F "file=@myfile.txt" https://your-domain.atlassian.net/rest/servicedeskapi/servicedesk/10001/attachTemporaryFile

**[Permissions](#permissions) required**: Permission to add attachments in this Service Desk.
- **`keen servicedeskapi check-request-type-permissions`** - Returns:

 *  a list of request type IDs where the given user has permission to administer.
 *  a list of request type IDs where the given user has permission to submit the request.

If no account ID is provided, the operation returns details for the logged in user.

Note that:

 *  invalid request type IDs are ignored.
 *  a maximum of 50 request types can be checked.

**[Permissions](#permissions) required:**

 *  *Administer Jira* or *Project Administrator* to check the permissions for other users.

However, Connect apps can make a call from the app server to the product to obtain permission details for any user, without admin permission. This Connect app ability doesn't apply to calls made using AP.request() in a browser.
- **`keen servicedeskapi create-comment-with-attachment`** - This method creates a comment on a customer request using one or more attachment files (uploaded using [servicedeskapi/servicedesk/\{serviceDeskId\}/attachTemporaryFile](https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-servicedesk/#api-rest-servicedeskapi-servicedesk-servicedeskid-attachtemporaryfile-post)), with the visibility set by `public`. See

 *  GET [servicedeskapi/request/\{issueIdOrKey\}/attachment](./#api-rest-servicedeskapi-request-issueidorkey-attachment-get)
 *  GET [servicedeskapi/request/\{issueIdOrKey\}/comment/\{commentId\}/attachment](./#api-rest-servicedeskapi-request-issueidorkey-comment-commentid-attachment-get)

**[Permissions](#permissions) required**: Permission to add an attachment.

**Request limitations**: Customers can set public visibility only.
- **`keen servicedeskapi create-customer`** - This method adds a customer to the Jira Service Management instance by passing a JSON file including an email address and display name. The display name does not need to be unique. The record's identifiers, `name` and `key`, are automatically generated from the request details.

**[Permissions](#permissions) required**: Jira Administrator Global permission
- **`keen servicedeskapi create-customer-customer`** - Creates a customer account on behalf of jsd-nutmeg.

This endpoint is restricted to jsd-nutmeg via ASAP authentication. It provides the same capability as the public \{@code POST /servicedeskapi/customer\} endpoint, but does not require a User Context Token (UCT) or Connect app user \\u2014 authorization is enforced entirely via the ASAP token.

No user permission checks are performed; \{@code null\} is passed as the acting user to bypass the permission check in the underlying service.
- **`keen servicedeskapi create-customer-request`** - This method creates a customer request in a service desk.

The JSON request must include the service desk and customer request type, as well as any fields that are required for the request type. A list of the fields required by a customer request type can be obtained using [servicedesk/\{serviceDeskId\}/requesttype/\{requestTypeId\}/field](#api-servicedesk-serviceDeskId-requesttype-requestTypeId-field-get).

The fields required for a customer request type depend on the user's permissions:

 *  `raiseOnBehalfOf` is not available to Users who have the customer permission only.
 *  `requestParticipants` is not available to Users who have the customer permission only or if the feature is turned off for customers.

`requestFieldValues` is a map of Jira field IDs and their values. See [Field input formats](#fieldformats), for details of each field's JSON semantics and the values they can take.

**[Permissions](#permissions) required**: Permission to create requests in the specified service desk.
- **`keen servicedeskapi create-organization`** - This method creates an organization by passing the name of the organization.

**[Permissions](#permissions) required**: Service desk administrator or agent. Note: Permission to create organizations can be switched to users with the Jira administrator permission, using the **[Organization management](https://confluence.atlassian.com/servicedeskcloud/setting-up-service-desk-users-732528877.html#Settingupservicedeskusers-manageorgsManageorganizations)** feature.
- **`keen servicedeskapi create-request-comment`** - This method creates a public or private (internal) comment on a customer request, with the comment visibility set by `public`. The user recorded as the author of the comment.

**[Permissions](#permissions) required**: User has Add Comments permission.

**Request limitations**: Customers can set comments to public visibility only.
- **`keen servicedeskapi create-request-type`** - This method enables a customer request type to be added to a service desk based on an issue type. Note that not all customer request type fields can be specified in the request and these fields are given the following default values:

 *  Request type icon is given the headset icon.
 *  Request type groups is left empty, which means this customer request type will not be visible on the [customer portal](https://confluence.atlassian.com/servicedeskcloud/configuring-the-customer-portal-732528918.html).
 *  Request type status mapping is left empty, so the request type has no custom status mapping but inherits the status map from the issue type upon which it is based.
 *  Request type field mapping is set to show the required fields as specified by the issue type used to create the customer request type.

  
These fields can be updated by a service desk administrator using the **Request types** option in **Project settings**.  
Request Types are created in next-gen projects by creating Issue Types. Please use the Jira Cloud Platform Create issue type endpoint instead.

**[Permissions](#permissions) required**: Service desk's administrator
- **`keen servicedeskapi delete-feedback`** - This method deletes the feedback of request using it's `requestKey` or `requestId`

**[Permissions](#permissions) required**: User must be the reporter or an Atlassian Connect app.
- **`keen servicedeskapi delete-organization`** - This method deletes an organization. Note that the organization is deleted regardless of other associations it may have. For example, associations with service desks.

**[Permissions](#permissions) required**: Jira administrator.
- **`keen servicedeskapi delete-property`** - Removes an organization property. Organization properties are a type of entity property which are available to the API only, and not shown in Jira Service Management. [Learn more](https://developer.atlassian.com/cloud/jira/platform/jira-entity-properties/).

For operations relating to organization detail field values which are visible in Jira Service Management, see the [Customer Service Management REST API](https://developer.atlassian.com/cloud/customer-service-management/rest/v1/api-group-organization/#api-group-organization).

**[Permissions](#permissions) required**: Service Desk Administrator or Agent.

Note: Permission to manage organizations can be switched to users with the Jira administrator permission, using the **[Organization management](https://confluence.atlassian.com/servicedeskcloud/setting-up-service-desk-users-732528877.html#Settingupservicedeskusers-manageorgsManageorganizations)** feature.
- **`keen servicedeskapi delete-property-servicedesk`** - Removes a property from a request type.

Properties for a Request Type in next-gen are stored as Issue Type properties and therefore can also be deleted by calling the Jira Cloud Platform [Delete issue type property](https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-issuetype-issueTypeId-properties-propertyKey-delete) endpoint.

**[Permissions](#permissions) required**: Jira project administrator with a Jira Service Management agent license.
- **`keen servicedeskapi delete-request-type`** - This method deletes a customer request type from a service desk, and removes it from all customer requests.  
This only supports classic projects.

**[Permissions](#permissions) required**: Service desk administrator.
- **`keen servicedeskapi get-all-request-types`** - This method returns all customer request types used in the Jira Service Management instance, optionally filtered by a query string.

Use [servicedeskapi/servicedesk/\{serviceDeskId\}/requesttype](#api-servicedesk-serviceDeskId-requesttype-get) to find the customer request types supported by a specific service desk.

The returned list of customer request types can be filtered using the `searchQuery` parameter. The parameter is matched against the customer request types' `name` or `description`. For example, searching for "Install", "Inst", "Equi", or "Equipment" will match a customer request type with the *name* "Equipment Installation Request".

**Note:** This API by default will filter out request types hidden in the portal (i.e. request types without groups and request types where a user doesn't have permission) when `searchQuery` is provided, unless `includeHiddenRequestTypesInSearch` is set to true. Restricted request types will not be returned for those who aren't admins.

**[Permissions](#permissions) required**: Any
- **`keen servicedeskapi get-approval-by-id`** - This method returns an approval. Use this method to determine the status of an approval and the list of approvers.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-approvals`** - This method returns all approvals on a customer request.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-articles`** - Returns articles which match the given query string across all service desks.

**[Permissions](#permissions) required**: Permission to access the [customer portal](https://confluence.atlassian.com/servicedeskcloud/configuring-the-customer-portal-732528918.html).
- **`keen servicedeskapi get-articles-servicedesk`** - Returns articles which match the given query and belong to the knowledge base linked to the service desk.

**[Permissions](#permissions) required**: Permission to access the service desk.
- **`keen servicedeskapi get-assets-workspaces`** - Returns a list of Assets workspace IDs. Include a workspace ID in the path to access the [Assets REST APIs](https://developer.atlassian.com/cloud/assets/rest).

**[Permissions](#permissions) required**: Any
- **`keen servicedeskapi get-attachment-content`** - Returns the contents of an attachment.

To return a thumbnail of the attachment, use [servicedeskapi/request/\{issueIdOrKey\}/attachment/\{attachmentId\}/thumbnail](./#api-rest-servicedeskapi-request-issueidorkey-attachment-attachmentid-thumbnail-get).

**[Permissions](#permissions) required:** For the issue containing the attachment:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen servicedeskapi get-attachment-thumbnail`** - Returns the thumbnail of an attachment.

To return the attachment contents, use [servicedeskapi/request/\{issueIdOrKey\}/attachment/\{attachmentId\}](#api-rest-servicedeskapi-request-issueidorkey-attachment-attachmentid-get).

**[Permissions](#permissions) required:** For the issue containing the attachment:

 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project that the issue is in.
 *  If [issue-level security](https://confluence.atlassian.com/x/J4lKLg) is configured, issue-level security permission to view the issue.
- **`keen servicedeskapi get-attachments-for-request`** - This method returns all the attachments for a customer requests.

**[Permissions](#permissions) required**: Permission to view the customer request.

**Response limitations**: Customers will only get a list of public attachments.
- **`keen servicedeskapi get-comment-attachments`** - This method returns the attachments referenced in a comment.

**[Permissions](#permissions) required**: Permission to view the customer request.

**Response limitations**: Customers can only view public comments, and retrieve their attachments, on requests where they are the reporter or a participant whereas agents can see both internal and public comments.
- **`keen servicedeskapi get-customer-request-by-id-or-key`** - This method returns a customer request.

**[Permissions](#permissions) required**: Permission to access the specified service desk.

**Response limitations**: For customers, only a request they created, was created on their behalf, or they are participating in will be returned.

**Note:** `requestFieldValues` does not include hidden fields. To get a list of request type fields that includes hidden fields, see [/rest/servicedeskapi/servicedesk/\{serviceDeskId\}/requesttype/\{requestTypeId\}/field](https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-servicedesk/#api-rest-servicedeskapi-servicedesk-servicedeskid-requesttype-requesttypeid-field-get)
- **`keen servicedeskapi get-customer-request-status`** - This method returns a list of all the statuses a customer Request has achieved. A status represents the state of an issue in its workflow. An issue can have one active status only. The list returns the status history in chronological order, most recent (current) status first.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-customer-requests`** - This method returns all customer requests for the user executing the query.

The returned customer requests are ordered chronologically by the latest activity on each request. For example, the latest status transition or comment.

**[Permissions](#permissions) required**: Permission to access the specified service desk.

**Response limitations**: For customers, the list returned will include request they created (or were created on their behalf) or are participating in only.
- **`keen servicedeskapi get-customer-transitions`** - This method returns a list of transitions, the workflow processes that moves a customer request from one status to another, that the user can perform on a request. Use this method to provide a user with a list if the actions they can take on a customer request.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-customers`** - This method returns a list of the customers on a service desk.

The returned list of customers can be filtered using the `query` parameter. The parameter is matched against customers' `displayName`, `name`, or `email`. For example, searching for "John", "Jo", "Smi", or "Smith" will match a user with display name "John Smith".

**[Permissions](#permissions) required**: Permission to view this Service Desk's customers.
- **`keen servicedeskapi get-feedback`** - This method retrieves a feedback of a request using it's `requestKey` or `requestId`

**[Permissions](#permissions) required**: User has view request permissions.
- **`keen servicedeskapi get-info`** - This method retrieves information about the Jira Service Management instance such as software version, builds, and related links.

**[Permissions](#permissions) required**: None, the user does not need to be logged in.
- **`keen servicedeskapi get-insight-workspaces`** - This endpoint is deprecated, please use /assets/workspace/.
- **`keen servicedeskapi get-issues-in-queue`** - This method returns the customer requests in a queue. Only fields that the queue is configured to show are returned. For example, if a queue is configured to show description and due date, then only those two fields are returned for each customer request in the queue.

**[Permissions](#permissions) required**: Service desk's agent.
- **`keen servicedeskapi get-organization`** - This method returns details of an organization. Use this method to get organization details whenever your application component is passed an organization ID but needs to display other organization details.

To get organization detail field values which are visible in Jira Service Management, see the [Customer Service Management REST API](https://developer.atlassian.com/cloud/customer-service-management/rest/v1/api-group-organization/#api-group-organization).

**[Permissions](#permissions) required**: Any

**Response limitations**: Customers can only retrieve organization of which they are members.
- **`keen servicedeskapi get-organizations`** - This method returns a list of organizations in the Jira Service Management instance. Use this method when you want to present a list of organizations or want to locate an organization by name.

**[Permissions](#permissions) required**: Any. However, to fetch organizations based on `accountId` the user must have a Service Desk agent license.

**Response limitations**: If the user is a customer, only those organizations of which the customer is a member are listed.
- **`keen servicedeskapi get-organizations-servicedesk`** - This method returns a list of all organizations associated with a service desk.

**[Permissions](#permissions) required**: Service desk's agent.
- **`keen servicedeskapi get-properties-keys`** - Returns the keys of all organization properties. Organization properties are a type of entity property which are available to the API only, and not shown in Jira Service Management. [Learn more](https://developer.atlassian.com/cloud/jira/platform/jira-entity-properties/).

To get organization detail field values which are visible in Jira Service Management, see the [Customer Service Management REST API](https://developer.atlassian.com/cloud/customer-service-management/rest/v1/api-group-organization/#api-group-organization).

**[Permissions](#permissions) required**: Any

**Response limitations**: Customers can only access properties of organizations of which they are members.
- **`keen servicedeskapi get-properties-keys-servicedesk`** - Returns the keys of all properties for a request type.

Properties for a Request Type in next-gen are stored as Issue Type properties and therefore the keys of all properties for a request type are also available by calling the Jira Cloud Platform [Get issue type property keys](https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-issuetype-issueTypeId-properties-get) endpoint.

**[Permissions](#permissions) required**: The user must have permission to view the request type.
- **`keen servicedeskapi get-property`** - Returns the value of an organization property. Use this method to obtain the JSON content for an organization's property. Organization properties are a type of entity property which are available to the API only, and not shown in Jira Service Management. [Learn more](https://developer.atlassian.com/cloud/jira/platform/jira-entity-properties/).

To get organization detail field values which are visible in Jira Service Management, see the [Customer Service Management REST API](https://developer.atlassian.com/cloud/customer-service-management/rest/v1/api-group-organization/#api-group-organization).

**[Permissions](#permissions) required**: Any

**Response limitations**: Customers can only access properties of organizations of which they are members.
- **`keen servicedeskapi get-property-servicedesk`** - Returns the value of the property from a request type.

Properties for a Request Type in next-gen are stored as Issue Type properties and therefore also available by calling the Jira Cloud Platform [Get issue type property](https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-issuetype-issueTypeId-properties-propertyKey-get) endpoint.

**[Permissions](#permissions) required**: User must have permission to view the request type.
- **`keen servicedeskapi get-queue`** - This method returns a specific queues in a service desk. To include a customer request count for the queue (in the `issueCount` field) in the response, set the query parameter `includeCount` to true (its default is false).

**[Permissions](#permissions) required**: service desk's Agent.
- **`keen servicedeskapi get-queues`** - This method returns the queues in a service desk. To include a customer request count for each queue (in the `issueCount` field) in the response, set the query parameter `includeCount` to true (its default is false).

**[Permissions](#permissions) required**: service desk's Agent.
- **`keen servicedeskapi get-request-comment-by-id`** - This method returns details of a customer request's comment.

**[Permissions](#permissions) required**: Permission to view the customer request.

**Response limitations**: Customers can only view public comments on requests where they are the reporter or a participant whereas agents can see both internal and public comments.
- **`keen servicedeskapi get-request-comments`** - This method returns all comments on a customer request. No permissions error is provided if, for example, the user doesn't have access to the service desk or request, the method simply returns an empty response.

**[Permissions](#permissions) required**: Permission to view the customer request.

**Response limitations**: Customers are returned public comments only.
- **`keen servicedeskapi get-request-participants`** - This method returns a list of all the participants on a customer request.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-request-type-by-id`** - This method returns a customer request type from a service desk.

This operation can be accessed anonymously.

**[Permissions](#permissions) required**: Permission to access the service desk.
- **`keen servicedeskapi get-request-type-fields`** - This method returns the fields for a service desk's customer request type.

Also, the following information about the user's permissions for the request type is returned:

 *  `canRaiseOnBehalfOf` returns `true` if the user has permission to raise customer requests on behalf of other customers. Otherwise, returns `false`.
 *  `canAddRequestParticipants` returns `true` if the user can add customer request participants. Otherwise, returns `false`.

**[Permissions](#permissions) required**: Permission to view the Service Desk. However, hidden fields would be visible to only Service desk's Administrator.
- **`keen servicedeskapi get-request-type-groups`** - This method returns a service desk's customer request type groups. Jira Service Management administrators can arrange the customer request type groups in an arbitrary order for display on the customer portal; the groups are returned in this order.

**[Permissions](#permissions) required**: Permission to view the service desk.
- **`keen servicedeskapi get-request-types`** - This method returns all customer request types from a service desk. There are two parameters for filtering the returned list:

 *  `groupId` which filters the results to items in the customer request type group.
 *  `searchQuery` which is matched against request types' `name` or `description`. For example, the strings "Install", "Inst", "Equi", or "Equipment" will match a request type with the *name* "Equipment Installation Request".

**Note:** This API by default will filter out request types hidden in the portal (i.e. request types without groups and request types where a user doesn't have permission) when `searchQuery` is provided, unless `includeHiddenRequestTypesInSearch` is set to true. Restricted request types will not be returned for those who aren't admins.

**[Permissions](#permissions) required**: Permission to access the service desk.
- **`keen servicedeskapi get-service-desk-by-id`** - This method returns a service desk. Use this method to get service desk details whenever your application component is passed a service desk ID but needs to display other service desk details.

**[Permissions](#permissions) required**: Permission to access the Service Desk. For example, being the Service Desk's Administrator or one of its Agents or Users.
- **`keen servicedeskapi get-service-desks`** - This method returns all the service desks in the Jira Service Management instance that the user has permission to access. Use this method where you need a list of service desks or need to locate a service desk by name or keyword.

**Note:** This method will be slow if the instance has hundreds of service desks. If you want to fetch a single service desk by its ID, use [/rest/servicedeskapi/servicedesk/\{serviceDeskId\}](./#api-rest-servicedeskapi-servicedesk-servicedeskid-get) instead.

**[Permissions](#permissions) required**: Any
- **`keen servicedeskapi get-sla-information`** - This method returns all the SLA records on a customer request. A customer request can have zero or more SLAs. Each SLA can have recordings for zero or more "completed cycles" and zero or 1 "ongoing cycle". Each cycle includes information on when it started and stopped, and whether it breached the SLA goal.

**[Permissions](#permissions) required**:

 *  Agent for the Service Desk containing the queried customer request, AND
 *  Browse Projects permission on the project containing the customer request, including any restrictions imposed by issue security schemes or custom permission schemes on the specific issue.
- **`keen servicedeskapi get-sla-information-by-id`** - This method returns the details for an SLA on a customer request.

**[Permissions](#permissions) required**:

 *  Agent for the Service Desk containing the queried customer request, AND
 *  Browse Projects permission on the project containing the customer request, including any restrictions imposed by issue security schemes or custom permission schemes on the specific issue.
- **`keen servicedeskapi get-subscription-status`** - This method returns the notification subscription status of the user making the request. Use this method to determine if the user is subscribed to a customer request's notifications.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi get-users-in-organization`** - This method returns all the users associated with an organization. Use this method where you want to provide a list of users for an organization or determine if a user is associated with an organization.

**[Permissions](#permissions) required**: Service desk administrator or agent.
- **`keen servicedeskapi invite-customer`** - This method invites a customer to a specified service desk by sending them an email invitation, creating a new customer account if one does not already exist. The display name does not need to be unique. The record's identifiers, `name` and `key`, are automatically generated from the request details.

**[Permissions](#permissions) required**: Jira Administrator Global permission & Service desk administrator
- **`keen servicedeskapi perform-customer-transition`** - This method performs a customer transition for a given request and transition. An optional comment can be included to provide a reason for the transition.

**[Permissions](#permissions) required**: The user must be able to view the request and have the Transition Issues permission. If a comment is passed the user must have the Add Comments permission.
- **`keen servicedeskapi post-feedback`** - This method adds a feedback on an request using it's `requestKey` or `requestId`

**[Permissions](#permissions) required**: User must be the reporter or an Atlassian Connect app.
- **`keen servicedeskapi remove-customers`** - This method removes one or more customers from a service desk. The service desk must have closed access. If any of the passed customers are not associated with the service desk, no changes will be made for those customers and the resource returns a 204 success code.

**[Permissions](#permissions) required**: Services desk administrator
- **`keen servicedeskapi remove-organization`** - This method removes an organization from a service desk. If the organization ID does not match an organization associated with the service desk, no change is made and the resource returns a 204 success code.

**[Permissions](#permissions) required**: Service desk's agent.
- **`keen servicedeskapi remove-request-participants`** - This method removes participants from a customer request.

**[Permissions](#permissions) required**: Permission to manage participants on the customer request.
- **`keen servicedeskapi remove-users-from-organization`** - This method removes users from an organization.

**[Permissions](#permissions) required**: Service desk administrator or agent. Note: Permission to delete users from an organization can be switched to users with the Jira administrator permission, using the **[Organization management](https://confluence.atlassian.com/servicedeskcloud/setting-up-service-desk-users-732528877.html#Settingupservicedeskusers-manageorgsManageorganizations)** feature.
- **`keen servicedeskapi revoke-portal-only-access-for-user`** - This method revokes portal-only access for a particular user, removing their ability to log in to the Jira Service Management customer portal as a portal-only user. After revocation, the user cannot submit or view requests through the portal.

**[Permissions](#permissions) required:** Site administration (that is, member of the *site-admin* [group](https://confluence.atlassian.com/x/24xjL)).
- **`keen servicedeskapi set-property`** - Sets the value of an organization property. Use this resource to store custom data against an organization. Organization properties are a type of entity property which are available to the API only, and not shown in Jira Service Management. [Learn more](https://developer.atlassian.com/cloud/jira/platform/jira-entity-properties/).

To store organization detail field values which are visible in Jira Service Management, see the [Customer Service Management REST API](https://developer.atlassian.com/cloud/customer-service-management/rest/v1/api-group-organization/#api-group-organization).

**[Permissions](#permissions) required**: Service Desk Administrator or Agent.

Note: Permission to manage organizations can be switched to users with the Jira administrator permission, using the **[Organization management](https://confluence.atlassian.com/servicedeskcloud/setting-up-service-desk-users-732528877.html#Settingupservicedeskusers-manageorgsManageorganizations)** feature.
- **`keen servicedeskapi set-property-servicedesk`** - Sets the value of a request type property. Use this resource to store custom data against a request type.

Properties for a Request Type in next-gen are stored as Issue Type properties and therefore can also be set by calling the Jira Cloud Platform [Set issue type property](https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-issuetype-issueTypeId-properties-propertyKey-put) endpoint.

**[Permissions](#permissions) required**: Jira project administrator with a Jira Service Management agent license.
- **`keen servicedeskapi subscribe`** - This method subscribes the user to receiving notifications from a customer request.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi unsubscribe`** - This method unsubscribes the user from notifications from a customer request.

**[Permissions](#permissions) required**: Permission to view the customer request.
- **`keen servicedeskapi validate-customer-request`** - Validates a customer request payload without creating (persisting) a request.

This endpoint runs exactly the same structural and semantic validations as [Create customer request](#api-request-post) \\u2014 including ProForma form validation \\u2014 but performs **no mutation**: no issue is created and no side effects (attachments, comments, analytics) run.

The response is intentionally verbose and structured so that it can be consumed by automated agents (for example an LLM repairing an invalid payload): every failure carries a machine-readable location (field id / form entity) and a human-readable reason. A valid payload returns HTTP 200 with \{@code valid: true\}; an invalid payload returns HTTP 400 with \{@code valid: false\} together with the field, form and general validation errors.

**[Permissions](#permissions) required**: Permission to create requests in the specified service desk.
- **`keen servicedeskapi view-article`** - View knowledge base article

### settings

Manage settings

- **`keen settings get-issue-navigator-default-columns`** - Returns the default issue navigator columns.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen settings set-issue-navigator-default-columns`** - Sets the default issue navigator columns.

The `columns` parameter accepts a navigable field value and is expressed as HTML form data. To specify multiple columns, pass multiple `columns` parameters. For example, in curl:

`curl -X PUT -d columns=summary -d columns=description https://your-domain.atlassian.net/rest/api/3/settings/columns`

If no column details are sent, then all default columns are removed.

A navigable field is one that can be used as a column on the issue navigator. Find details of navigable issue columns using [Get fields](#api-rest-api-3-field-get).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### snippets

Snippets allow you share code segments or files with yourself, members of
your workspace, or the world.

Like pull requests, repositories and workspaces, the full set of snippets
is defined by what the current user has access to. This includes all
snippets owned by any of the workspaces the user is a member of, or
snippets by other users that the current user is either watching or has
 collaborated on (for instance by commenting on it).

- **`keen snippets create`** - Creates a new snippet under the authenticated user's account.

Snippets can contain multiple files. Both text and binary files are
supported.

The simplest way to create a new snippet from a local file:

    $ curl -u username:password -X POST https://api.bitbucket.org/2.0/snippets               -F file=@image.png

Creating snippets through curl has a few limitations and so let's look
at a more complicated scenario.

Snippets are created with a multipart POST. Both `multipart/form-data`
and `multipart/related` are supported. Both allow the creation of
snippets with both meta data (title, etc), as well as multiple text
and binary files.

The main difference is that `multipart/related` can use rich encoding
for the meta data (currently JSON).


multipart/related (RFC-2387)
----------------------------

This is the most advanced and efficient way to create a paste.

    POST /2.0/snippets/evzijst HTTP/1.1
    Content-Length: 1188
    Content-Type: multipart/related; start="snippet"; boundary="===============1438169132528273974=="
    MIME-Version: 1.0

    --===============1438169132528273974==
    Content-Type: application/json; charset="utf-8"
    MIME-Version: 1.0
    Content-ID: snippet

    {
      "title": "My snippet",
      "is_private": true,
      "scm": "git",
      "files": {
          "foo.txt": {},
          "image.png": {}
        }
    }

    --===============1438169132528273974==
    Content-Type: text/plain; charset="us-ascii"
    MIME-Version: 1.0
    Content-Transfer-Encoding: 7bit
    Content-ID: "foo.txt"
    Content-Disposition: attachment; filename="foo.txt"

    foo

    --===============1438169132528273974==
    Content-Type: image/png
    MIME-Version: 1.0
    Content-Transfer-Encoding: base64
    Content-ID: "image.png"
    Content-Disposition: attachment; filename="image.png"

    iVBORw0KGgoAAAANSUhEUgAAABQAAAAoCAYAAAD+MdrbAAABD0lEQVR4Ae3VMUoDQRTG8ccUaW2m
    TKONFxArJYJamCvkCnZTaa+VnQdJSBFl2SMsLFrEWNjZBZs0JgiL/+KrhhVmJRbCLPx4O+/DT2TB
    cbblJxf+UWFVVRNsEGAtgvJxnLm2H+A5RQ93uIl+3632PZyl/skjfOn9Gvdwmlcw5aPUwimG+NT5
    EnNN036IaZePUuIcK533NVfal7/5yjWeot2z9ta1cAczHEf7I+3J0ws9Cgx0fsOFpmlfwKcWPuBQ
    73Oc4FHzBaZ8llq4q1mr5B2mOUCt815qYR8eB1hG2VJ7j35q4RofaH7IG+Xrf/PfJhfmwtfFYoIN
    AqxFUD6OMxcvkO+UfKfkOyXfKdsv/AYCHMLVkHAFWgAAAABJRU5ErkJggg==
    --===============1438169132528273974==--

The request contains multiple parts and is structured as follows.

The first part is the JSON document that describes the snippet's
properties or meta data. It either has to be the first part, or the
request's `Content-Type` header must contain the `start` parameter to
point to it.

The remaining parts are the files of which there can be zero or more.
Each file part should contain the `Content-ID` MIME header through
which the JSON meta data's `files` element addresses it. The value
should be the name of the file.

`Content-Disposition` is an optional MIME header. The header's
optional `filename` parameter can be used to specify the file name
that Bitbucket should use when writing the file to disk. When present,
`filename` takes precedence over the value of `Content-ID`.

When the JSON body omits the `files` element, the remaining parts are
not ignored. Instead, each file is added to the new snippet as if its
name was explicitly linked (the use of the `files` elements is
mandatory for some operations like deleting or renaming files).


multipart/form-data
-------------------

The use of JSON for the snippet's meta data is optional. Meta data can
also be supplied as regular form fields in a more conventional
`multipart/form-data` request:

    $ curl -X POST -u credentials https://api.bitbucket.org/2.0/snippets               -F title="My snippet"               -F file=@foo.txt -F file=@image.png

    POST /2.0/snippets HTTP/1.1
    Content-Length: 951
    Content-Type: multipart/form-data; boundary=----------------------------63a4b224c59f

    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="file"; filename="foo.txt"
    Content-Type: text/plain

    foo

    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="file"; filename="image.png"
    Content-Type: application/octet-stream

    ?PNG

    IHDR?1??I.....
    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="title"

    My snippet
    ------------------------------63a4b224c59f--

Here the meta data properties are included as flat, top-level form
fields. The file attachments use the `file` field name. To attach
multiple files, simply repeat the field.

The advantage of `multipart/form-data` over `multipart/related` is
that it can be easier to build clients.

Essentially all properties are optional, `title` and `files` included.


Sharing and Visibility
----------------------

Snippets can be either public (visible to anyone on Bitbucket, as well
as anonymous users), or private (visible only to members of the workspace).
This is controlled through the snippet's `is_private` element:

* **is_private=false** -- everyone, including anonymous users can view
  the snippet
* **is_private=true** -- only workspace members can view the snippet

To create the snippet under a workspace, just append the workspace ID
to the URL. See [`/2.0/snippets/{workspace}`](/cloud/bitbucket/rest/api-group-snippets/#api-snippets-workspace-post).
- **`keen snippets create-for-workspace`** - Identical to [`/snippets`](/cloud/bitbucket/rest/api-group-snippets/#api-snippets-post), except that the new snippet will be
created under the workspace specified in the path parameter
`{workspace}`.
- **`keen snippets delete`** - Deletes a snippet and returns an empty response.
- **`keen snippets delete-previous-revision-of`** - Deletes the snippet.

Note that this only works for versioned URLs that point to the latest
commit of the snippet. Pointing to an older commit results in a 405
status code.

To delete a snippet, regardless of whether or not concurrent changes
are being made to it, use `DELETE /snippets/{encoded_id}` instead.
- **`keen snippets get`** - Retrieves a single snippet.

Snippets support multiple content types:

* application/json
* multipart/related
* multipart/form-data


application/json
----------------

The default content type of the response is `application/json`.
Since JSON is always `utf-8`, it cannot reliably contain file contents
for files that are not text. Therefore, JSON snippet documents only
contain the filename and links to the file contents.

This means that in order to retrieve all parts of a snippet, N+1
requests need to be made (where N is the number of files in the
snippet).


multipart/related
-----------------

To retrieve an entire snippet in a single response, use the
`Accept: multipart/related` HTTP request header.

    $ curl -H "Accept: multipart/related" https://api.bitbucket.org/2.0/snippets/evzijst/1

Response:

    HTTP/1.1 200 OK
    Content-Length: 2214
    Content-Type: multipart/related; start="snippet"; boundary="===============1438169132528273974=="
    MIME-Version: 1.0

    --===============1438169132528273974==
    Content-Type: application/json; charset="utf-8"
    MIME-Version: 1.0
    Content-ID: snippet

    {
      "links": {
        "self": {
          "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj"
        },
        "html": {
          "href": "https://bitbucket.org/snippets/evzijst/kypj"
        },
        "comments": {
          "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj/comments"
        },
        "watchers": {
          "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj/watchers"
        },
        "commits": {
          "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj/commits"
        }
      },
      "id": kypj,
      "title": "My snippet",
      "created_on": "2014-12-29T22:22:04.790331+00:00",
      "updated_on": "2014-12-29T22:22:04.790331+00:00",
      "is_private": false,
      "files": {
        "foo.txt": {
          "links": {
            "self": {
              "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj/files/367ab19/foo.txt"
            },
            "html": {
              "href": "https://bitbucket.org/snippets/evzijst/kypj#file-foo.txt"
            }
          }
        },
        "image.png": {
          "links": {
            "self": {
              "href": "https://api.bitbucket.org/2.0/snippets/evzijst/kypj/files/367ab19/image.png"
            },
            "html": {
              "href": "https://bitbucket.org/snippets/evzijst/kypj#file-image.png"
            }
          }
        }
      ],
      "owner": {
        "username": "evzijst",
        "nickname": "evzijst",
        "display_name": "Erik van Zijst",
        "uuid": "{d301aafa-d676-4ee0-88be-962be7417567}",
        "links": {
          "self": {
            "href": "https://api.bitbucket.org/2.0/users/evzijst"
          },
          "html": {
            "href": "https://bitbucket.org/evzijst"
          },
          "avatar": {
            "href": "https://bitbucket-staging-assetroot.s3.amazonaws.com/c/photos/2013/Jul/31/erik-avatar-725122544-0_avatar.png"
          }
        }
      },
      "creator": {
        "username": "evzijst",
        "nickname": "evzijst",
        "display_name": "Erik van Zijst",
        "uuid": "{d301aafa-d676-4ee0-88be-962be7417567}",
        "links": {
          "self": {
            "href": "https://api.bitbucket.org/2.0/users/evzijst"
          },
          "html": {
            "href": "https://bitbucket.org/evzijst"
          },
          "avatar": {
            "href": "https://bitbucket-staging-assetroot.s3.amazonaws.com/c/photos/2013/Jul/31/erik-avatar-725122544-0_avatar.png"
          }
        }
      }
    }

    --===============1438169132528273974==
    Content-Type: text/plain; charset="us-ascii"
    MIME-Version: 1.0
    Content-Transfer-Encoding: 7bit
    Content-ID: "foo.txt"
    Content-Disposition: attachment; filename="foo.txt"

    foo

    --===============1438169132528273974==
    Content-Type: image/png
    MIME-Version: 1.0
    Content-Transfer-Encoding: base64
    Content-ID: "image.png"
    Content-Disposition: attachment; filename="image.png"

    iVBORw0KGgoAAAANSUhEUgAAABQAAAAoCAYAAAD+MdrbAAABD0lEQVR4Ae3VMUoDQRTG8ccUaW2m
    TKONFxArJYJamCvkCnZTaa+VnQdJSBFl2SMsLFrEWNjZBZs0JgiL/+KrhhVmJRbCLPx4O+/DT2TB
    cbblJxf+UWFVVRNsEGAtgvJxnLm2H+A5RQ93uIl+3632PZyl/skjfOn9Gvdwmlcw5aPUwimG+NT5
    EnNN036IaZePUuIcK533NVfal7/5yjWeot2z9ta1cAczHEf7I+3J0ws9Cgx0fsOFpmlfwKcWPuBQ
    73Oc4FHzBaZ8llq4q1mr5B2mOUCt815qYR8eB1hG2VJ7j35q4RofaH7IG+Xrf/PfJhfmwtfFYoIN
    AqxFUD6OMxcvkO+UfKfkOyXfKdsv/AYCHMLVkHAFWgAAAABJRU5ErkJggg==
    --===============1438169132528273974==--

multipart/form-data
-------------------

As with creating new snippets, `multipart/form-data` can be used as an
alternative to `multipart/related`. However, the inherently flat
structure of form-data means that only basic, root-level properties
can be returned, while nested elements like `links` are omitted:

    $ curl -H "Accept: multipart/form-data" https://api.bitbucket.org/2.0/snippets/evzijst/kypj

Response:

    HTTP/1.1 200 OK
    Content-Length: 951
    Content-Type: multipart/form-data; boundary=----------------------------63a4b224c59f

    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="title"
    Content-Type: text/plain; charset="utf-8"

    My snippet
    ------------------------------63a4b224c59f--
    Content-Disposition: attachment; name="file"; filename="foo.txt"
    Content-Type: text/plain

    foo

    ------------------------------63a4b224c59f
    Content-Disposition: attachment; name="file"; filename="image.png"
    Content-Transfer-Encoding: base64
    Content-Type: application/octet-stream

    iVBORw0KGgoAAAANSUhEUgAAABQAAAAoCAYAAAD+MdrbAAABD0lEQVR4Ae3VMUoDQRTG8ccUaW2m
    TKONFxArJYJamCvkCnZTaa+VnQdJSBFl2SMsLFrEWNjZBZs0JgiL/+KrhhVmJRbCLPx4O+/DT2TB
    cbblJxf+UWFVVRNsEGAtgvJxnLm2H+A5RQ93uIl+3632PZyl/skjfOn9Gvdwmlcw5aPUwimG+NT5
    EnNN036IaZePUuIcK533NVfal7/5yjWeot2z9ta1cAczHEf7I+3J0ws9Cgx0fsOFpmlfwKcWPuBQ
    73Oc4FHzBaZ8llq4q1mr5B2mOUCt815qYR8eB1hG2VJ7j35q4RofaH7IG+Xrf/PfJhfmwtfFYoIN
    AqxFUD6OMxcvkO+UfKfkOyXfKdsv/AYCHMLVkHAFWgAAAABJRU5ErkJggg==
    ------------------------------5957323a6b76--
- **`keen snippets get-previous-revision-of`** - Identical to `GET /snippets/encoded_id`, except that this endpoint
can be used to retrieve the contents of the snippet as it was at an
older revision, while `/snippets/encoded_id` always returns the
snippet's current revision.

Note that only the snippet's file contents are versioned, not its
meta data properties like the title.

Other than that, the two endpoints are identical in behavior.
- **`keen snippets list-in-workspace`** - Returns a paginated list of snippets owned by `{workspace}`.

To limit the set of returned snippets, apply the
`?role=[owner|contributor|member]` query parameter where the roles are
defined as follows:

* `owner`: snippets owned by `{workspace}` that also belong to the current user
    (only returns results when `{workspace}` is the current user's personal workspace)
* `contributor`: snippets owned by `{workspace}` that the current user is watching,
    plus any owned by `{workspace}` and the current user
* `member`: all snippets owned by `{workspace}` if the current user is a member,
    otherwise only those the current user is watching

When no role is specified, all snippets owned by `{workspace}` are returned.

If the current user is not a member of `{workspace}`, only public snippets are
returned regardless of role.

The returned response is a normal paginated JSON list. This endpoint
only supports `application/json` responses and no
`multipart/form-data` or `multipart/related`. As a result, it is not
possible to include the file contents.
- **`keen snippets update`** - Used to update a snippet. Use this to add and delete files and to
change a snippet's title.

To update a snippet, one can either PUT a full snapshot, or only the
parts that need to be changed.

The contract for PUT on this API is that properties missing from the
request remain untouched so that snippets can be efficiently
manipulated with differential payloads.

To delete a property (e.g. the title, or a file), include its name in
the request, but omit its value (use `null`).

As in Git, explicit renaming of files is not supported. Instead, to
rename a file, delete it and add it again under another name. This can
be done atomically in a single request. Rename detection is left to
the SCM.

PUT supports three different content types for both request and
response bodies:

* `application/json`
* `multipart/related`
* `multipart/form-data`

The content type used for the request body can be different than that
used for the response. Content types are specified using standard HTTP
headers.

Use the `Content-Type` and `Accept` headers to select the desired
request and response format.


application/json
----------------

As with creation and retrieval, the content type determines what
properties can be manipulated. `application/json` does not support
file contents and is therefore limited to a snippet's meta data.

To update the title, without changing any of its files:

    $ curl -X POST -H "Content-Type: application/json" https://api.bitbucket.org/2.0/snippets/evzijst/kypj             -d '{"title": "Updated title"}'


To delete the title:

    $ curl -X POST -H "Content-Type: application/json" https://api.bitbucket.org/2.0/snippets/evzijst/kypj             -d '{"title": null}'

Not all parts of a snippet can be manipulated. The owner and creator
for instance are immutable.


multipart/related
-----------------

`multipart/related` can be used to manipulate all of a snippet's
properties. The body is identical to a POST. properties omitted from
the request are left unchanged. Since the `start` part contains JSON,
the mechanism for manipulating the snippet's meta data is identical
to `application/json` requests.

To update one of a snippet's file contents, while also changing its
title:

    PUT /2.0/snippets/evzijst/kypj HTTP/1.1
    Content-Length: 288
    Content-Type: multipart/related; start="snippet"; boundary="===============1438169132528273974=="
    MIME-Version: 1.0

    --===============1438169132528273974==
    Content-Type: application/json; charset="utf-8"
    MIME-Version: 1.0
    Content-ID: snippet

    {
      "title": "My updated snippet",
      "files": {
          "foo.txt": {}
        }
    }

    --===============1438169132528273974==
    Content-Type: text/plain; charset="us-ascii"
    MIME-Version: 1.0
    Content-Transfer-Encoding: 7bit
    Content-ID: "foo.txt"
    Content-Disposition: attachment; filename="foo.txt"

    Updated file contents.

    --===============1438169132528273974==--

Here only the parts that are changed are included in the body. The
other files remain untouched.

Note the use of the `files` list in the JSON part. This list contains
the files that are being manipulated. This list should have
corresponding multiparts in the request that contain the new contents
of these files.

If a filename in the `files` list does not have a corresponding part,
it will be deleted from the snippet, as shown below:

    PUT /2.0/snippets/evzijst/kypj HTTP/1.1
    Content-Length: 188
    Content-Type: multipart/related; start="snippet"; boundary="===============1438169132528273974=="
    MIME-Version: 1.0

    --===============1438169132528273974==
    Content-Type: application/json; charset="utf-8"
    MIME-Version: 1.0
    Content-ID: snippet

    {
      "files": {
        "image.png": {}
      }
    }

    --===============1438169132528273974==--

To simulate a rename, delete a file and add the same file under
another name:

    PUT /2.0/snippets/evzijst/kypj HTTP/1.1
    Content-Length: 212
    Content-Type: multipart/related; start="snippet"; boundary="===============1438169132528273974=="
    MIME-Version: 1.0

    --===============1438169132528273974==
    Content-Type: application/json; charset="utf-8"
    MIME-Version: 1.0
    Content-ID: snippet

    {
        "files": {
          "foo.txt": {},
          "bar.txt": {}
        }
    }

    --===============1438169132528273974==
    Content-Type: text/plain; charset="us-ascii"
    MIME-Version: 1.0
    Content-Transfer-Encoding: 7bit
    Content-ID: "bar.txt"
    Content-Disposition: attachment; filename="bar.txt"

    foo

    --===============1438169132528273974==--


multipart/form-data
-----------------

Again, one can also use `multipart/form-data` to manipulate file
contents and meta data atomically.

    $ curl -X PUT http://localhost:12345/2.0/snippets/evzijst/kypj             -F title="My updated snippet" -F file=@foo.txt

    PUT /2.0/snippets/evzijst/kypj HTTP/1.1
    Content-Length: 351
    Content-Type: multipart/form-data; boundary=----------------------------63a4b224c59f

    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="file"; filename="foo.txt"
    Content-Type: text/plain

    foo

    ------------------------------63a4b224c59f
    Content-Disposition: form-data; name="title"

    My updated snippet
    ------------------------------63a4b224c59f

To delete a file, omit its contents while including its name in the
`files` field:

    $ curl -X PUT https://api.bitbucket.org/2.0/snippets/evzijst/kypj -F files=image.png

    PUT /2.0/snippets/evzijst/kypj HTTP/1.1
    Content-Length: 149
    Content-Type: multipart/form-data; boundary=----------------------------ef8871065a86

    ------------------------------ef8871065a86
    Content-Disposition: form-data; name="files"

    image.png
    ------------------------------ef8871065a86--

The explicit use of the `files` element in `multipart/related` and
`multipart/form-data` is only required when deleting files.
The default mode of operation is for file parts to be processed,
regardless of whether or not they are listed in `files`, as a
convenience to the client.
- **`keen snippets update-previous-revision-of`** - Identical to `UPDATE /snippets/encoded_id`, except that this endpoint
takes an explicit commit revision. Only the snippet's "HEAD"/"tip"
(most recent) version can be updated and requests on all other,
older revisions fail by returning a 405 status.

Usage of this endpoint over the unrestricted `/snippets/encoded_id`
could be desired if the caller wants to be sure no concurrent
modifications have taken place between the moment of the UPDATE
request and the original GET.

This can be considered a so-called "Compare And Swap", or CAS
operation.

Other than that, the two endpoints are identical in behavior.

### software

Manage software

- **`keen software get-approximate-issue-count-for-backlog`** - Returns the approximate count of all issues from the board's backlog, for the given board ID. This is equivalent to counting the issues on all pages returned by [Get issues for backlog enhanced](https://developer.atlassian.com/cloud/jira/software/rest/api-group-board/#api-rest-software-1-0-board-boardid-backlog-get). Recent updates might not be immediately visible in the returned output. This only includes issues that the user has permission to view.
- **`keen software get-approximate-issue-count-for-board`** - Returns the approximate count of all issues from a board, for a given board ID. This is equivalent to counting the issues on all pages returned by [Get issues for board enhanced](https://developer.atlassian.com/cloud/jira/software/rest/api-group-board/#api-rest-software-1-0-board-boardid-issue-get). Recent updates might not be immediately visible in the returned output. This only includes issues that the user has permission to view.
- **`keen software get-board-issues-for-epic-jsis`** - Returns all issues that belong to an epic on the board, for the given epic ID and the board ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen software get-board-issues-for-sprint-jsis`** - Get all issues you have access to that belong to the sprint from the board. Result pagination is token-based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource contains additional fields like: sprint, closedSprints, flagged, and epic. Issues are returned ordered by rank. JQL order has higher priority than default rank.
- **`keen software get-issues-for-backlog-jsis`** - Returns all issues from the board's backlog, for the given board ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Note, if the user does not have permission to view the board, no issues will be returned at all. The backlog contains incomplete issues that are not assigned to any future or active sprint. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen software get-issues-for-board-jsis`** - Returns all issues from a board, for a given board ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Note, if the user does not have permission to view the board, no issues will be returned at all. An issue belongs to the board if its status is mapped to the board's column. Epic issues do not belong to scrum boards. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen software get-issues-for-epic-jsis`** - Returns all issues that belong to the epic, for the given epic ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank. **Note:** If you are querying a Team Managed project, do not use this operation. Instead, search for issues that belong to an epic by using the [Search for issues using JQL enhanced search](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-search/#api-rest-api-3-search-jql-get) operation in the Jira platform REST API. Build your JQL query using the `parent` clause. For more information on the `parent` JQL field, see [Advanced searching](https://confluence.atlassian.com/x/dAiiLQ#Advancedsearching-fieldsreference-Parent).
- **`keen software get-issues-for-sprint-jsis`** - Returns all issues in a sprint, for a given sprint ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen software get-issues-without-epic-for-board-jsis`** - Returns all issues that do not belong to any epic on a board, for a given board ID. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
- **`keen software get-issues-without-epic-jsis`** - Returns all issues that do not belong to any epic. Result pagination is token based, using `nextPageToken` and `maxResults`. This only includes issues that the user has permission to view. Issues returned from this resource include Software project fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank. **Note:** If you are querying a Team Managed project, do not use this operation. Instead, search for issues that don't belong to an epic by using the [Search for issues using JQL enhanced search](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-search/#api-rest-api-3-search-jql-get) operation in the Jira platform REST API. Build your JQL query using the `parent is empty` clause. For more information on the `parent` JQL field, see [Advanced searching](https://confluence.atlassian.com/x/dAiiLQ#Advancedsearching-fieldsreference-Parent).

### space-permissions

Manage space permissions

- **`keen space-permissions bulk-assign-roles`** - Bulk assigns roles for one or more permission combination IDs obtained from the space permission
combinations. Supports targeting all spaces, specific spaces, or excluding specific spaces.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a Confluence administrator.
- **`keen space-permissions bulk-remove-access`** - Bulk removes access for one or more permission combination IDs obtained from the space permission
combinations. This removes all space permissions for the specified combinations across
the targeted spaces.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a Confluence administrator.
- **`keen space-permissions generate-combinations`** - Submits a task to refresh the space permission combinations in the database, which identifies
all unique permission combinations across the site. This provides permission combination IDs
that can be used with the assign-roles and remove-access endpoints.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a Confluence administrator.
- **`keen space-permissions get-available`** - Retrieves the available space permissions.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site.
- **`keen space-permissions get-transition-task-status`** - Retrieves the status of an async space permission transition task. Use the taskId returned
from the generate-combinations, assign-roles, or remove-access endpoints to poll for
progress and completion.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a Confluence administrator.
- **`keen space-permissions list-combinations`** - Lists the unique unassigned space permission combinations currently present on the tenant.
Combinations that already map to a space role are filtered out server-side. Each row carries
the decoded set of space permissions and the principal types that currently hold the
combination — these inform which `principalType` values are valid to include in the matching
bulk role-assignments request.

Results are always sorted by `principalCount` descending. Sort field and sort order are not
configurable; page size is controlled by the `limit` query parameter (default 25, min 1,
max 250). Use the `cursor` field to page through additional results. The `generatedAt` field
reflects the last audit run that populated the combinations table — call the
generate-combinations endpoint to refresh stale data.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a Confluence administrator.

### space-role-mode

Manage space role mode

- **`keen space-role-mode`** - Retrieves the space role mode.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).

### space-roles

Manage space roles

- **`keen space-roles create`** - Create a space role.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be an organization or site admin. Connect and Forge app users are not authorized to access this resource.
- **`keen space-roles delete`** - Delete a space role

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be an organization or site admin. Connect and Forge app users are not authorized to access this resource.
- **`keen space-roles get-available`** - Retrieves the available space roles.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site; if requesting a certain space's roles, permission to view the space.
- **`keen space-roles get-by-id`** - Retrieves the space role by ID.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site.
- **`keen space-roles update`** - Update a space role.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be an organization or site admin. Connect and Forge app users are not authorized to access this resource.

### spaces

Manage spaces

- **`keen spaces create`** - Creates a Space as specified in the payload.

Available on tenants with [Role-Based Access Control](https://support.atlassian.com/confluence-cloud/docs/manage-user-roles/). 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to create spaces.
- **`keen spaces get`** - Returns all spaces. The results will be sorted by id ascending. The number of results is limited by the `limit` parameter and
additional results (if available) will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only spaces that the user has permission to view will be returned.
- **`keen spaces get-by-id`** - Returns a specific space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the space.

### status

This resource represents statuses. Use it to search, get, create, delete, and change statuses.

- **`keen status get`** - Returns a status. The status must be associated with an active workflow to be returned.

If a name is used on more than one status, only the status found first is returned. Therefore, identifying the status by its ID may be preferable.

This operation can be accessed anonymously.

[Permissions](#permissions) required: *Browse projects* [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) for the project.
- **`keen status get-statuses`** - Returns a list of all statuses associated with active workflows.

This operation can be accessed anonymously.

[Permissions](#permissions) required: *Browse projects* [project permission](https://support.atlassian.com/jira-cloud-administration/docs/manage-project-permissions/) for the project.

### statuscategory

Manage statuscategory

- **`keen statuscategory get-status-categories`** - Returns a list of all status categories.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen statuscategory get-status-category`** - Returns a status category. Status categories provided a mechanism for categorizing [statuses](#api-rest-api-3-status-idOrName-get).

**[Permissions](#permissions) required:** Permission to access Jira.

### statuses

Manage statuses

- **`keen statuses create`** - Creates statuses for a global or project scope.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)
- **`keen statuses delete-by-id`** - Deletes statuses by ID.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)
- **`keen statuses get-by-id`** - Returns a list of the statuses specified by one or more status IDs.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)
- **`keen statuses get-by-name`** - Returns a list of the statuses specified by one or more status names.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Browse projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
- **`keen statuses search`** - Returns a [paginated](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#pagination) list of statuses that match a search on name or project.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)
- **`keen statuses update`** - Updates statuses by ID.

**[Permissions](#permissions) required:**

 *  *Administer projects* [project permission.](https://confluence.atlassian.com/x/yodKLg)
 *  *Administer Jira* [project permission.](https://confluence.atlassian.com/x/yodKLg)

### task

This resource represents a [long-running asynchronous tasks](#async-operations). Use it to obtain details about the progress of a long-running task or cancel a long-running task.

- **`keen task <taskId>`** - Returns the status of a [long-running asynchronous task](#async).

When a task has finished, this operation returns the JSON blob applicable to the task. See the documentation of the operation that created the task for details. Task details are not permanently retained. As of September 2019, details are retained for 14 days although this period may change without notice.

**Deprecation notice:** The required OAuth 2.0 scopes will be updated on June 15, 2024.

 *  `read:jira-work`

**[Permissions](#permissions) required:** either of:

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
 *  Creator of the task.

### tasks

Manage tasks

- **`keen tasks get`** - Returns all tasks. The number of results is limited by the `limit` parameter and additional results (if available)
will be available through the `next` URL present in the `Link` response header.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only tasks that the user has permission to view will be returned.
- **`keen tasks get-by-id`** - Returns a specific task. 

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the containing page or blog post and its corresponding space.
- **`keen tasks update`** - Update a task by id. This endpoint currently only supports updating task status.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the containing page or blog post and view its corresponding space.

### teams

The teams resource has been deprecated, and the workspaces
endpoint should be used instead.

The teams resource returns all the teams that the authenticated
user is associated with.


### ui-modifications

Manage ui modifications

- **`keen ui-modifications create`** - Creates a UI modification. UI modification can only be created by Forge apps.

Each app can define up to 3000 UI modifications. Each UI modification can define up to 1000 contexts. The same context can be assigned to maximum 100 UI modifications.

**Context types:**

 *  **Jira contexts:** For Jira view types, use `projectId` and `issueTypeId`. One field can act as a wildcard. Supported Jira views:
    
     *  `GIC` \- Jira global issue create
     *  `IssueView` \- Jira issue view
     *  `IssueTransition` \- Jira issue transition
 *  **Jira Service Management contexts:** For Jira Service Management view types, use `portalId` and `requestTypeId`. Wildcards are not supported. Supported JSM views:
    
     *  `JSMRequestCreate` \- Jira Service Management request create portal view
 *  **Agent view contexts:** For Agent view types, use `projectId` and `issueTypeId` like Jira contexts, and optionally set `requestTypeId`. `portalId` must not be set. One of `projectId`, `issueTypeId`, or `viewType` can act as a wildcard. Supported Agent views:
    
     *  `GICAgentView` \- Agent view variant of Jira global issue create
     *  `IssueViewAgentView` \- Agent view variant of Jira issue view
     *  `IssueTransitionAgentView` \- Agent view variant of Jira issue transition

**[Permissions](#permissions) required:**

 *  *None* if the UI modification is created without contexts.
 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for one or more projects, if the UI modification is created with contexts.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen ui-modifications delete`** - Deletes a UI modification. All the contexts that belong to the UI modification are deleted too. UI modification can only be deleted by Forge apps.

**[Permissions](#permissions) required:** None.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen ui-modifications get`** - Gets UI modifications. UI modifications can only be retrieved by Forge apps.

**[Permissions](#permissions) required:** None.

The new `read:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.
- **`keen ui-modifications update`** - Updates a UI modification. UI modification can only be updated by Forge apps.

Each UI modification can define up to 1000 contexts. The same context can be assigned to maximum 100 UI modifications.

**Context types:**

 *  **Jira contexts:** For Jira view types, use `projectId` and `issueTypeId`. One field can act as a wildcard. Supported Jira views:
    
     *  `GIC` \- Jira global issue create
     *  `IssueView` \- Jira issue view
     *  `IssueTransition` \- Jira issue transition
 *  **Jira Service Management contexts:** For Jira Service Management view types, use `portalId` and `requestTypeId`. Wildcards are not supported. Supported JSM views:
    
     *  `JSMRequestCreate` \- Jira Service Management request create portal view
 *  **Agent view contexts:** For Agent view types, use `projectId` and `issueTypeId` like Jira contexts, and optionally set `requestTypeId`. `portalId` must not be set. One of `projectId`, `issueTypeId`, or `viewType` can act as a wildcard. Supported Agent views:
    
     *  `GICAgentView` \- Agent view variant of Jira global issue create
     *  `IssueViewAgentView` \- Agent view variant of Jira issue view
     *  `IssueTransitionAgentView` \- Agent view variant of Jira issue transition

**[Permissions](#permissions) required:**

 *  *None* if the UI modification is created without contexts.
 *  *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for one or more projects, if the UI modification is created with contexts.

The new `write:app-data:jira` OAuth scope is 100% optional now, and not using it won't break your app. However, we recommend adding it to your app's scope list because we will eventually make it mandatory.

### universal-avatar

Manage universal avatar

- **`keen universal-avatar delete-avatar`** - Deletes an avatar from a project, issue type or priority.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen universal-avatar get-avatar-image-by-id`** - Returns a project, issue type or priority avatar image by ID.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  For system avatars, none.
 *  For custom project avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project the avatar belongs to.
 *  For custom issue type avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for at least one project the issue type is used in.
 *  For priority avatars, none.
- **`keen universal-avatar get-avatar-image-by-owner`** - Returns the avatar image for a project, issue type or priority.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  For system avatars, none.
 *  For custom project avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project the avatar belongs to.
 *  For custom issue type avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for at least one project the issue type is used in.
 *  For priority avatars, none.
- **`keen universal-avatar get-avatar-image-by-type`** - Returns the default project, issue type or priority avatar image.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** None.
- **`keen universal-avatar get-avatars`** - Returns the system and custom avatars for a project, issue type or priority.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  for custom project avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for the project the avatar belongs to.
 *  for custom issue type avatars, *Browse projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for at least one project the issue type is used in.
 *  for system avatars, none.
 *  for priority avatars, none.
- **`keen universal-avatar store-avatar`** - Loads a custom avatar for a project, issue type or priority.

Specify the avatar's local file location in the body of the request. Also, include the following headers:

 *  `X-Atlassian-Token: no-check` To prevent XSRF protection blocking the request, for more information see [Special Headers](#special-request-headers).
 *  `Content-Type: image/image type` Valid image types are JPEG, GIF, or PNG.

For example:  
`curl --request POST `

`--user email@example.com:<api_token> `

`--header 'X-Atlassian-Token: no-check' `

`--header 'Content-Type: image/< image_type>' `

`--data-binary "<@/path/to/file/with/your/avatar>" `

`--url 'https://your-domain.atlassian.net/rest/api/3/universal_avatar/type/{type}/owner/{entityId}'`

The avatar is cropped to a square. If no crop parameters are specified, the square originates at the top left of the image. The length of the square's sides is set to the smaller of the height or width of the image.

The cropped image is then used to create avatars of 16x16, 24x24, 32x32, and 48x48 in size.

After creating the avatar use:

 *  [Update issue type](#api-rest-api-3-issuetype-id-put) to set it as the issue type's displayed avatar.
 *  [Set project avatar](#api-rest-api-3-project-projectIdOrKey-avatar-put) to set it as the project's displayed avatar.
 *  [Update priority](#api-rest-api-3-priority-id-put) to set it as the priority's displayed avatar.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### usage

Manage usage

- **`keen usage`** - Retrieves comprehensive usage statistics for the current tenant including total object counts and a per-schema breakdown for billing and analytics.

### user

This resource represent users. Use it to:

 *  get, get a list of, create, and delete users.
 *  get, set, and reset a user's default issue table columns.
 *  get a list of the groups the user belongs to.
 *  get a list of user account IDs for a list of usernames or user keys.

- **`keen user bulk-get`** - Returns a [paginated](#pagination) list of the users specified by one or more account IDs.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen user bulk-get-migration`** - Returns the account IDs for the users specified in the `key` or `username` parameters. Note that multiple `key` or `username` parameters can be specified.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen user create`** - Creates a user. This resource is retained for legacy compatibility. As soon as a more suitable alternative is available this resource will be deprecated.

**Note:** This API does not support Forge apps.

If the user exists and has access to Jira, the operation returns a 201 status. If the user exists but does not have access to Jira & no new jira-products are requested, the operation returns a 400 status.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg). The caller has to be an **organization admin**.
- **`keen user delete-property`** - Deletes a property from a user.

Note: This operation does not access the [user properties](https://confluence.atlassian.com/x/8YxjL) created and maintained in Jira.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to delete a property from any user.
 *  Access to Jira, to delete a property from the calling user's record.
- **`keen user find`** - Returns a list of active users that match the search string and property.

This operation first applies a filter to match the search string and property, and then takes the filtered users in the range defined by `startAt` and `maxResults`, up to the thousandth user. To get all the users who match the search string and property, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

This operation can be accessed anonymously.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg). Anonymous calls or calls by users without the required permission return empty search results.
- **`keen user find-assignable`** - Returns a list of users that can be assigned to an issue. Use this operation to find the list of users who can be assigned to:

 *  a new issue, by providing the `projectKeyOrId`.
 *  an updated issue, by providing the `issueKey` or `issueId`.
 *  to an issue during a transition (workflow action), by providing the `issueKey` or `issueId` and the transition id in `actionDescriptorId`. You can obtain the IDs of an issue's valid transitions using the `transitions` option in the `expand` parameter of [ Get issue](#api-rest-api-3-issue-issueIdOrKey-get).

In all these cases, you can pass an account ID to determine if a user can be assigned to an issue. The user is returned in the response if they can be assigned to the issue or issue transition.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that can be assigned the issue. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who can be assigned the issue, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg) or *Assign issues* [project permission](https://confluence.atlassian.com/x/yodKLg)
- **`keen user find-bulk-assignable`** - Returns a list of users who can be assigned issues in one or more projects. The list may be restricted to users whose attributes match a string.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that can be assigned issues in the projects. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who can be assigned issues in the projects, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for each project specified in `projectKeys`.
- **`keen user find-by-query`** - Finds users with a structured query and returns a [paginated](#pagination) list of user details.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that match the structured query. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who match the structured query, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).

The query statements are:

 *  `is assignee of PROJ` Returns the users that are assignees of at least one issue in project *PROJ*.
 *  `is assignee of (PROJ-1, PROJ-2)` Returns users that are assignees on the issues *PROJ-1* or *PROJ-2*.
 *  `is reporter of (PROJ-1, PROJ-2)` Returns users that are reporters on the issues *PROJ-1* or *PROJ-2*.
 *  `is watcher of (PROJ-1, PROJ-2)` Returns users that are watchers on the issues *PROJ-1* or *PROJ-2*.
 *  `is voter of (PROJ-1, PROJ-2)` Returns users that are voters on the issues *PROJ-1* or *PROJ-2*.
 *  `is commenter of (PROJ-1, PROJ-2)` Returns users that have posted a comment on the issues *PROJ-1* or *PROJ-2*.
 *  `is transitioner of (PROJ-1, PROJ-2)` Returns users that have performed a transition on issues *PROJ-1* or *PROJ-2*.
 *  `[propertyKey].entity.property.path is "property value"` Returns users with the entity property value. For example, if user property `location` is set to value `{"office": {"country": "AU", "city": "Sydney"}}`, then it's possible to use `[location].office.city is "Sydney"` to match the user.

The list of issues can be extended as needed, as in *(PROJ-1, PROJ-2, ... PROJ-n)*. Statements can be combined using the `AND` and `OR` operators to form more complex queries. For example:

`is assignee of PROJ AND [propertyKey].entity.property.path is "property value"`
- **`keen user find-for-picker`** - Returns a list of users whose attributes match the query term. The returned object includes the `html` field where the matched query term is highlighted with the HTML strong tag. A list of account IDs can be provided to exclude users from the results.

This operation takes the users in the range defined by `maxResults`, up to the thousandth user, and then returns only the users from that range that match the query term. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who match the query term, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg). Anonymous calls and calls by users without the required permission return search results for an exact name match only.
- **`keen user find-keys-by-query`** - Finds users with a structured query and returns a [paginated](#pagination) list of user keys.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that match the structured query. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who match the structured query, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).

The query statements are:

 *  `is assignee of PROJ` Returns the users that are assignees of at least one issue in project *PROJ*.
 *  `is assignee of (PROJ-1, PROJ-2)` Returns users that are assignees on the issues *PROJ-1* or *PROJ-2*.
 *  `is reporter of (PROJ-1, PROJ-2)` Returns users that are reporters on the issues *PROJ-1* or *PROJ-2*.
 *  `is watcher of (PROJ-1, PROJ-2)` Returns users that are watchers on the issues *PROJ-1* or *PROJ-2*.
 *  `is voter of (PROJ-1, PROJ-2)` Returns users that are voters on the issues *PROJ-1* or *PROJ-2*.
 *  `is commenter of (PROJ-1, PROJ-2)` Returns users that have posted a comment on the issues *PROJ-1* or *PROJ-2*.
 *  `is transitioner of (PROJ-1, PROJ-2)` Returns users that have performed a transition on issues *PROJ-1* or *PROJ-2*.
 *  `[propertyKey].entity.property.path is "property value"` Returns users with the entity property value. For example, if user property `location` is set to value `{"office": {"country": "AU", "city": "Sydney"}}`, then it's possible to use `[location].office.city is "Sydney"` to match the user.

The list of issues can be extended as needed, as in *(PROJ-1, PROJ-2, ... PROJ-n)*. Statements can be combined using the `AND` and `OR` operators to form more complex queries. For example:

`is assignee of PROJ AND [propertyKey].entity.property.path is "property value"`
- **`keen user find-with-all-permissions`** - Returns a list of users who fulfill these criteria:

 *  their user attributes match a search string.
 *  they have a set of permissions for a project or issue.

If no search string is provided, a list of all users with the permissions is returned.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that match the search string and have permission for the project or issue. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who match the search string and have permission for the project or issue, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to get users for any project.
 *  *Administer Projects* [project permission](https://confluence.atlassian.com/x/yodKLg) for a project, to get users for that project.
- **`keen user find-with-browse-permission`** - Returns a list of users who fulfill these criteria:

 *  their user attributes match a search string.
 *  they have permission to browse issues.

Use this resource to find users who can browse:

 *  an issue, by providing the `issueKey`.
 *  any issue in a project, by providing the `projectKey`.

This operation takes the users in the range defined by `startAt` and `maxResults`, up to the thousandth user, and then returns only the users from that range that match the search string and have permission to browse issues. This means the operation usually returns fewer users than specified in `maxResults`. To get all the users who match the search string and have permission to browse issues, use [Get all users](#api-rest-api-3-users-search-get) and filter the records in your code.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

This operation can be accessed anonymously.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg). Anonymous calls and calls by users without the required permission return empty search results.
- **`keen user get`** - Returns a user.

Privacy controls are applied to the response based on the user's preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen user get-default-columns`** - Returns the default [issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user. If `accountId` is not passed in the request, the calling user's details are returned.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLgl), to get the column details for any user.
 *  Permission to access Jira, to get the calling user's column details.
- **`keen user get-email`** - Returns a user's email address regardless of the user's profile visibility settings. For Connect apps, this API is only available to apps approved by Atlassian, according to these [guidelines](https://community.developer.atlassian.com/t/guidelines-for-requesting-access-to-email-address/27603). For Forge apps, this API only supports access via asApp() requests.
- **`keen user get-email-bulk`** - Returns a user's email address regardless of the user's profile visibility settings. For Connect apps, this API is only available to apps approved by Atlassian, according to these [guidelines](https://community.developer.atlassian.com/t/guidelines-for-requesting-access-to-email-address/27603). For Forge apps, this API only supports access via asApp() requests.
- **`keen user get-groups`** - Returns the groups to which a user belongs.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen user get-property`** - Returns the value of a user's property. If no property key is provided [Get user property keys](#api-rest-api-3-user-properties-get) is called.

Note: This operation does not access the [user properties](https://confluence.atlassian.com/x/8YxjL) created and maintained in Jira.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to get a property from any user.
 *  Access to Jira, to get a property from the calling user's record.
- **`keen user get-property-keys`** - Returns the keys of all properties for a user.

Note: This operation does not access the [user properties](https://confluence.atlassian.com/x/8YxjL) created and maintained in Jira.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to access the property keys on any user.
 *  Access to Jira, to access the calling user's property keys.
- **`keen user remove`** - Deletes a user. If the operation completes successfully then the user is removed from Jira's user base. This operation does not delete the user's Atlassian account.

**[Permissions](#permissions) required:** Site administration (that is, membership of the *site-admin* [group](https://confluence.atlassian.com/x/24xjL)).
- **`keen user reset-columns`** - Resets the default [ issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user to the system default. If `accountId` is not passed, the calling user's default columns are reset.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to set the columns on any user.
 *  Permission to access Jira, to set the calling user's columns.
- **`keen user set-columns`** - Sets the default [ issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user. If an account ID is not passed, the calling user's default columns are set. If no column details are sent, then all default columns are removed.

The parameters for this resource are expressed as HTML form data. For example, in curl:

`curl -X PUT -d columns=summary -d columns=description https://your-domain.atlassian.net/rest/api/3/user/columns?accountId=5b10ac8d82e05b22cc7d4ef5'`

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to set the columns on any user.
 *  Permission to access Jira, to set the calling user's columns.
- **`keen user set-property`** - Sets the value of a user's property. Use this resource to store custom data against a user.

Note: This operation does not access the [user properties](https://confluence.atlassian.com/x/8YxjL) created and maintained in Jira.

**[Permissions](#permissions) required:**

 *  *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg), to set a property on any user.
 *  Access to Jira, to set a property on the calling user's record.

### users

This resource represent users. Use it to:

 *  get, get a list of, create, and delete users.
 *  get, set, and reset a user's default issue table columns.
 *  get a list of the groups the user belongs to.
 *  get a list of user account IDs for a list of usernames or user keys.

- **`keen users get-all`** - Returns a list of all users, including active users, inactive users and previously deleted users that have an Atlassian account.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen users get-all-default`** - Returns a list of all users, including active users, inactive users and previously deleted users that have an Atlassian account.

Privacy controls are applied to the response based on the users' preferences. This could mean, for example, that the user's email address is hidden. See the [Profile visibility overview](https://developer.atlassian.com/cloud/jira/platform/profile-visibility/) for more details.

**[Permissions](#permissions) required:** *Browse users and groups* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### users-bulk

Manage users bulk

- **`keen users-bulk`** - Returns user details for the ids provided in the request body.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
The user must be able to view user profiles in the Confluence site.

### vendors

Information about vendors


### webhook

This resource represents webhooks. Webhooks are calls sent to a URL when an event occurs in Jira for issues specified by a JQL query. Only Connect and OAuth 2.0 apps can register and manage webhooks. For more information, see [Webhooks](https://developer.atlassian.com/cloud/jira/platform/webhooks/#registering-a-webhook-via-the-jira-rest-api-for-connect-apps).

- **`keen webhook delete-by-id`** - Removes webhooks by ID. Only webhooks registered by the calling app are removed. If webhooks created by other apps are specified, they are ignored.

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/#connect-apps) and [OAuth 2.0](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps) apps can use this operation.
- **`keen webhook get-dynamic-for-app`** - Returns a [paginated](#pagination) list of the webhooks registered by the calling app.

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/#connect-apps) and [OAuth 2.0](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps) apps can use this operation.
- **`keen webhook get-failed`** - Returns webhooks that have recently failed to be delivered to the requesting app after the maximum number of retries.

After 72 hours the failure may no longer be returned by this operation.

The oldest failure is returned first.

This method uses a cursor-based pagination. To request the next page use the failure time of the last webhook on the list as the `failedAfter` value or use the URL provided in `next`.

**[Permissions](#permissions) required:** Only [Connect apps](https://developer.atlassian.com/cloud/jira/platform/index/#connect-apps) can use this operation.
- **`keen webhook refresh`** - Extends the life of webhook. Webhooks registered through the REST API expire after 30 days. Call this operation to keep them alive.

Unrecognized webhook IDs (those that are not found or belong to other apps) are ignored.

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/#connect-apps) and [OAuth 2.0](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps) apps can use this operation.
- **`keen webhook register-dynamic`** - Registers webhooks.

**NOTE:** for non-public OAuth apps, webhooks are delivered only if there is a match between the app owner and the user who registered a dynamic webhook.

**[Permissions](#permissions) required:** Only [Connect](https://developer.atlassian.com/cloud/jira/platform/#connect-apps) and [OAuth 2.0](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps) apps can use this operation.

### whiteboards

Manage whiteboards

- **`keen whiteboards create`** - Creates a whiteboard in the space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the corresponding space. Permission to create a whiteboard in the space.
- **`keen whiteboards delete`** - Delete a whiteboard by id.

Deleting a whiteboard moves the whiteboard to the trash, where it can be restored later

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the whiteboard and its corresponding space.
Permission to delete whiteboards in the space.
- **`keen whiteboards get-by-id`** - Returns a specific whiteboard.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the whiteboard and its corresponding space.

### wiki

Manage wiki

- **`keen wiki add-content-watcher`** - Adds a user as a watcher to a piece of content. Choose the user by doing
one of the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

Note, you must add the `X-Atlassian-Token: no-check` header when making a
request, as this operation has XSRF protection.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki add-custom-content-permissions`** - Adds new custom content permission to space.

If the permission to be added is a group permission, the group can be identified
by its group name or group id.

Note: Only apps can access this REST resource and only make changes to the respective app permissions.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki add-group-to-content-restriction-by-group-id`** - Adds a group to a content restriction by Group Id. That is, grant read or update
permission to the group for a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki add-label-watcher`** - Adds a user as a watcher to a label. Choose the user by doing one of the
following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

Note, you must add the `X-Atlassian-Token: no-check` header when making a
request, as this operation has XSRF protection.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki add-labels-to-content`** - Adds labels to a piece of content. Does not modify the existing labels.

Notes:

- Labels can also be added when creating content ([Create content](#api-content-post)).
- Labels can be updated when updating content ([Update content](#api-content-id-put)).
This will delete the existing labels and replace them with the labels in
the request.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki add-labels-to-space`** - Adds labels to a piece of content. Does not modify the existing labels.

Notes:

- Labels can also be added when creating content ([Create content](#api-content-post)).
- Labels can be updated when updating content ([Update content](#api-content-id-put)).
This will delete the existing labels and replace them with the labels in
the request.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki add-permission-to-space`** - Adds new permission to space.

If the permission to be added is a group permission, the group can be identified
by its group name or group id.

Note: Apps cannot access this REST resource - including when utilizing user impersonation.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki add-restrictions`** - Adds the supplied users and groups to the content's restrictions, preserving
existing restrictions. Use `read` for viewing restrictions and `update` for
editing restrictions. Omitted or empty user and group lists add no restrictions.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki add-space-watcher`** - Adds a user as a watcher to a space. Choose the user by doing one of the
following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

Note, you must add the `X-Atlassian-Token: no-check` header when making a
request, as this operation has XSRF protection.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki add-user-to-content-restriction`** - Adds a user to a content restriction. That is, grant read or update
permission to the user for a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki add-user-to-group-by-group-id`** - Adds a user as a member in a group represented by its groupId

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a site admin.
- **`keen wiki archive-pages`** - Archives a list of pages. The pages to be archived are specified as a list of content IDs.
This API accepts the archival request and returns a task ID.
The archival process happens asynchronously.
Use the /longtask/<taskId> REST API to get the copy task status.

Each content ID needs to resolve to page objects that are not already in an archived state.
The content IDs need not belong to the same space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Archive' permission for each of the pages in the corresponding space it belongs to.
- **`keen wiki async-convert-content-body-request`** - Converts a content body from one format to another format asynchronously.
Returns the asyncId for the asynchronous task.

Supported conversions:

- atlas_doc_format: editor, export_view, storage, styled_view, view
- storage: atlas_doc_format, editor, export_view, styled_view, view
- editor: storage

No other conversions are supported at the moment.
Once a conversion is completed, it will be available for 5 minutes at the result endpoint.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
If request specifies 'contentIdContext', 'View' permission for the space, and permission to view the content.
- **`keen wiki async-convert-content-body-response`** - Returns the content body for the corresponding `asyncId` of a completed conversion task. If
the task is not completed, the task status is returned instead.

Once a conversion task is completed, the result can be obtained for up to 5 minutes, or
until an identical conversion request is made again with the `allowCache` parameter set to
false.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
If request specifies 'contentIdContext', 'View' permission for the space, and permission to view the content.
- **`keen wiki bulk-async-convert-content-body-request`** - Asynchronously converts content bodies from one format to another format in bulk. Use the Content body
REST API to get the status of conversion tasks. Note that there is a maximum limit of 10 conversions per
request to this endpoint.

Supported conversions:

- storage: editor, export_view, styled_view, view
- editor: storage

Once a conversion task is completed, it is available for polling for up to 5 minutes.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space, and permission to view the content if the `spaceKeyContext` or
`contentIdContext` are present.
- **`keen wiki bulk-async-convert-content-body-response`** - Returns the content body for the corresponding `asyncId` of a completed conversion task. If
the task is not completed, the task status is returned instead.

Once a conversion task is completed, the result can be obtained for up to 5 minutes, or
until an identical conversion request is made again with the `allowCache` parameter set to
false.

Note that there is a maximum limit of 50 task results per request to this endpoint.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki check-content-permission`** - Check if a user or a group can perform an operation to the specified content. The `operation` to check
must be provided. The user’s account ID or the ID of the group can be provided in the `subject` to check
permissions against a specified user or group. The following permission checks are done to make sure that the
user or group has the proper access:

- site permissions
- space permissions
- content restrictions

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission) if checking permission for self,
otherwise 'Confluence Administrator' global permission is required.
- **`keen wiki copy-page`** - Copies a single page and its associated properties, permissions, attachments, and custom contents.
 The `id` path parameter refers to the content ID of the page to copy. The target of the page to be copied
 is defined using the `destination` in the request body and can be one of the following types.

  - `space`: page will be copied to the specified space as a root page on the space
  - `parent_page`: page will be copied as a child of the specified parent page
  - `parent_content`: page will be copied as a child of the specified parent content
  - `existing_page`: page will be copied and replace the specified page

By default, the following objects are expanded: `space`, `history`, `version`.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: 'Add' permission for the space that the content will be copied in and permission to update the content if copying to an `existing_page`.
- **`keen wiki copy-page-hierarchy`** - Copy page hierarchy allows the copying of an entire hierarchy of pages and their associated properties, permissions and attachments.
 The id path parameter refers to the content id of the page to copy, and the new parent of this copied page is defined using the destinationPageId in the request body.
 The titleOptions object defines the rules of renaming page titles during the copy;
 for example, search and replace can be used in conjunction to rewrite the copied page titles.

 Response example:
 <pre><code>
 {
      "id" : "1180606",
      "links" : {
           "status" : "/rest/api/longtask/1180606"
      }
 }
 </code></pre>
 Use the /longtask/<taskId> REST API to get the copy task status.
- **`keen wiki create-attachment`** - Adds an attachment to a piece of content. This method only adds a new
attachment. If you want to update an existing attachment, use
[Create or update attachments](#api-content-id-child-attachment-put).

Note, you must set a `X-Atlassian-Token: nocheck` header on the request
for this method, otherwise it will be blocked. This protects against XSRF
attacks, which is necessary as this method accepts multipart/form-data.

The media type 'multipart/form-data' is defined in [RFC 7578](https://www.ietf.org/rfc/rfc7578.txt).
Most client libraries have classes that make it easier to implement
multipart posts, like the [MultipartEntityBuilder](https://hc.apache.org/httpcomponents-client-5.1.x/current/httpclient5/apidocs/)
Java class provided by Apache HTTP Components.

Note, according to [RFC 7578](https://tools.ietf.org/html/rfc7578#section-4.5),
in the case where the form data is text,
the charset parameter for the "text/plain" Content-Type may be used to
indicate the character encoding used in that part. In the case of this
API endpoint, the `comment` body parameter should be sent with `type=text/plain`
and `charset=utf-8` values. This will force the charset to be UTF-8.

Example: This curl command attaches a file ('example.txt') to a container
(id='123') with a comment and `minorEdits`=true.

``` bash
curl -D- \
  -u admin:admin \
  -X POST \
  -H 'X-Atlassian-Token: nocheck' \
  -F 'file=@"example.txt"' \
  -F 'minorEdit="true"' \
  -F 'comment="Example attachment comment"; type=text/plain; charset=utf-8' \
  https://myhost/wiki/rest/api/content/123/child/attachment
```
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki create-audit-record`** - Creates a record in the audit log.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki create-content-template`** - Creates a new content template. Note, blueprint templates cannot be created via the REST API.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space to create a space template or 'Confluence Administrator'
global permission to create a global template.
- **`keen wiki create-group`** - Creates a new user group.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a site admin.
- **`keen wiki create-or-update-attachments`** - Adds an attachment to a piece of content. If the attachment already exists
for the content, then the attachment is updated (i.e. a new version of the
attachment is created).

Note, you must set a `X-Atlassian-Token: nocheck` header on the request
for this method, otherwise it will be blocked. This protects against XSRF
attacks, which is necessary as this method accepts multipart/form-data.

The media type 'multipart/form-data' is defined in [RFC 7578](https://www.ietf.org/rfc/rfc7578.txt).
Most client libraries have classes that make it easier to implement
multipart posts, like the [MultipartEntityBuilder](https://hc.apache.org/httpcomponents-client-5.1.x/current/httpclient5/apidocs/)
Java class provided by Apache HTTP Components.

Note, according to [RFC 7578](https://tools.ietf.org/html/rfc7578#section-4.5),
in the case where the form data is text,
the charset parameter for the "text/plain" Content-Type may be used to
indicate the character encoding used in that part. In the case of this
API endpoint, the `comment` body parameter should be sent with `type=text/plain`
and `charset=utf-8` values. This will force the charset to be UTF-8.

Example: This curl command attaches a file ('example.txt') to a piece of
content (id='123') with a comment and `minorEdits`=true. If the 'example.txt'
file already exists, it will update it with a new version of the attachment.

``` bash
curl -D- \
  -u admin:admin \
  -X PUT \
  -H 'X-Atlassian-Token: nocheck' \
  -F 'file=@"example.txt"' \
  -F 'minorEdit="true"' \
  -F 'comment="Example attachment comment"; type=text/plain; charset=utf-8' \
  http://myhost/rest/api/content/123/child/attachment
```
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki create-private-space`** - Creates a new space that is only visible to the creator. This method is
the same as the [Create space](#api-space-post) method with permissions
set to the current user only. Note, currently you cannot set space
labels when creating a space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Create Space(s)' global permission.
- **`keen wiki create-relationship`** - Creates a relationship between two entities (user, space, content). The
'favourite' relationship is supported by default, but you can use this method
to create any type of relationship between two entities.

For example, the following method creates a 'sibling' relationship between
two pieces of content:
`PUT /wiki/rest/api/relation/sibling/from/content/123/to/content/456`

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki create-space`** - Creates a new space. Note, currently you cannot set space labels when
creating a space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Create Space(s)' global permission.
- **`keen wiki create-user-property`** - Creates a property for a user. For more information  about user properties, see [Confluence entity properties]
(https://developer.atlassian.com/cloud/confluence/confluence-entity-properties/).
`Note`, these properties stored against a user are on a Confluence site level and not space/content level.

`Note:` the number of properties which could be created per app in a tenant for each user might be
restricted by fixed system limits.
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki delete-content-version`** - Delete a historical version. This does not delete the changes made to the
content in that version, rather the changes for the deleted version are
rolled up into the next version. Note, you cannot delete the current version.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki delete-label-from-space`** - Remove label from a space
- **`keen wiki delete-page-tree`** - Moves a pagetree rooted at a page to the space's trash:

- If the content's type is `page` and its status is `current`, it will be trashed including
all its descendants.
- For every other combination of content type and status, this API is not supported.

This API accepts the pageTree delete request and returns a task ID.
The delete process happens asynchronously.

 Response example:
 <pre><code>
 {
      "id" : "1180606",
      "links" : {
           "status" : "/rest/api/longtask/1180606"
      }
 }
 </code></pre>
 Use the `/longtask/<taskId>` REST API to get the copy task status.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Delete' permission for the space that the content is in.
- **`keen wiki delete-relationship`** - Deletes a relationship between two entities (user, space, content).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
For favourite relationships, the current user can only delete their own
favourite relationships. A space administrator can delete favourite
relationships for any user.
- **`keen wiki delete-restrictions`** - Removes all restrictions (read and update) on a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki delete-space`** - Permanently deletes a space without sending it to the trash. Note, the space will be deleted in a long running task.
Therefore, the space may not be deleted yet when this method has
returned. Clients should poll the status link that is returned in the
response until the task completes.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki delete-user-property`** - Deletes a property for the given user.
For more information about user properties, see
[Confluence entity properties](https://developer.atlassian.com/cloud/confluence/confluence-entity-properties/).
`Note`, these properties stored against a user are on a Confluence site level and not space/content level.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki download-attatchment`** - Redirects the client to a URL that serves an attachment's binary data.
- **`keen wiki export-audit-records`** - Exports audit records as a CSV file or ZIP file.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki find-sources-for-target`** - Returns all target entities that have a particular relationship to the
source entity. Note, relationships are one way.

For example, the following method finds all users that have a 'collaborator'
relationship to a piece of content with an ID of '1234':
`GET /wiki/rest/api/relation/collaborator/to/content/1234/from/user`
Note, 'collaborator' is an example custom relationship type.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view both the target entity and source entity.
- **`keen wiki find-target-from-source`** - Returns all target entities that have a particular relationship to the
source entity. Note, relationships are one way.

For example, the following method finds all content that the current user
has an 'ignore' relationship with:
`GET /wiki/rest/api/relation/ignore/from/user/current/to/content`
Note, 'ignore' is an example custom relationship type.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view both the target entity and source entity.
- **`keen wiki get-all-label-content`** - Returns label information and a list of contents associated with the label.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission). Only contents
that the user is permitted to view is returned.
- **`keen wiki get-and-async-convert-macro-body-by-macro-id`** - Returns Async Id of the conversion task which will convert the macro into a content body of the desired format.
The result will be available for 5 minutes after completion of the conversion.

About the macro ID: When a macro is created in a new version of content,
Confluence will generate a random ID for it, unless an ID is specified
(by an app). The macro ID will look similar to this: '884bd9-0cb8-41d5-98be-f80943c14f96'.
The ID is then persisted as new versions of content are created, and is
only modified by Confluence if there are conflicting IDs.

For Forge macros, the value for macro ID is the "local ID" of that particular ADF node.
This value can be retrieved either client-side by calling view.getContext() and accessing "localId"
on the resulting object, or server-side by examining the "local-id" parameter node inside the "parameters" node.

Note that there are other attributes named "local-id", but only this particular one is used to store the macro ID.

Example:
<ac:adf-node type="extension">
  <ac:adf-attribute key="extension-type">com.atlassian.ecosystem</ac:adf-attribute>
  <ac:adf-attribute key="parameters">
      <ac:adf-parameter key="local-id">e9c4aa10-73fa-417c-888d-48c719ae4165</ac:adf-parameter>
  </ac:adf-parameter>
</ac:adf-node>

Note, to preserve backwards compatibility this resource will also match on
the hash of the macro body, even if a macro ID is found. This check will
eventually become redundant, as macro IDs are generated for pages and
transparently propagate out to all instances.

This backwards compatibility logic does not apply to Forge macros; those
can only be retrieved by their ID.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content that the macro is in.
- **`keen wiki get-and-convert-macro-body-by-macro-id`** - Returns the body of a macro in format specified in path, for the given macro ID.
This includes information like the name of the macro, the body of the macro,
and any macro parameters.

About the macro ID: When a macro is created in a new version of content,
Confluence will generate a random ID for it, unless an ID is specified
(by an app). The macro ID will look similar to this: '50884bd9-0cb8-41d5-98be-f80943c14f96'.
The ID is then persisted as new versions of content are created, and is
only modified by Confluence if there are conflicting IDs.

For Forge macros, the value for macro ID is the "local ID" of that particular ADF node.
This value can be retrieved either client-side by calling view.getContext() and accessing "localId"
on the resulting object, or server-side by examining the "local-id" parameter node inside the "parameters" node.

Note that there are other attributes named "local-id", but only this particular one is used to store the macro ID.

Example:
<ac:adf-node type="extension">
  <ac:adf-attribute key="extension-type">com.atlassian.ecosystem</ac:adf-attribute>
  <ac:adf-attribute key="parameters">
      <ac:adf-parameter key="local-id">e9c4aa10-73fa-417c-888d-48c719ae4165</ac:adf-parameter>
  </ac:adf-parameter>
</ac:adf-node>

Note, to preserve backwards compatibility this resource will also match on
the hash of the macro body, even if a macro ID is found. This check will
eventually become redundant, as macro IDs are generated for pages and
transparently propagate out to all instances.

This backwards compatibility logic does not apply to Forge macros; those
can only be retrieved by their ID.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content that the macro is in.
- **`keen wiki get-anonymous-user`** - Returns information about how anonymous users are represented, like the
profile picture and display name.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-audit-records`** - Returns all records in the audit log, optionally for a certain date range.
This contains information about events like space exports, group membership
changes, app installations, etc. For more information, see
[Audit log](https://confluence.atlassian.com/confcloud/audit-log-802164269.html)
in the Confluence administrator's guide.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki get-audit-records-for-time-period`** - Returns records from the audit log, for a time period back from the current
date. For example, you can use this method to get the last 3 months of records.

This contains information about events like space exports, group membership
changes, app installations, etc. For more information, see
[Audit log](https://confluence.atlassian.com/confcloud/audit-log-802164269.html)
in the Confluence administrator's guide.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki get-available-content-states`** - Gets content states that are available for the content to be set as.
Will return all enabled Space Content States.
Will only return most the 3 most recently published custom content states to match UI editor list.
To get all custom content states, use the /content-states endpoint.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki get-blueprint-templates`** - Returns all templates provided by blueprints. Use this method to retrieve
all global blueprint templates or all blueprint templates in a space.

Note, all global blueprints are inherited by each space. Space blueprints
can be customised without affecting the global blueprints.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space to view blueprints for the space and permission
to access the Confluence site ('Can use' global permission) to view global blueprints.
- **`keen wiki get-bulk-user-lookup`** - Returns user details for the ids provided in the request.
Currently this API returns a maximum of 100 results.
If more than 100 account ids are passed in, then the first 100 will be returned.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-content-descendants`** - Returns a map of the descendants of a piece of content. This is similar
to [Get content children](#api-content-id-child-get), except that this
method returns child pages at all levels, rather than just the direct
child pages.

A piece of content has different types of descendants, depending on its type:

- `page`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `whiteboard`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `database`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `embed`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `folder`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `blogpost`: descendant is `comment`, `attachment`
- `attachment`: descendant is `comment`
- `comment`: descendant is `attachment`

The map will always include all descendant types that are valid for the content.
However, if the content has no instances of a descendant type, the map will
contain an empty array for that descendant type.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space, and permission to view the content if it
is a page.
- **`keen wiki get-content-restriction-status-for-user`** - Returns whether the specified content restriction applies to a user.
For example, if a page with `id=123` has a `read` restriction for a user with an account ID of
`384093:32b4d9w0-f6a5-3535-11a3-9c8c88d10192`, the following request will return `true`:

`/wiki/rest/api/content/123/restriction/byOperation/read/user?accountId=384093:32b4d9w0-f6a5-3535-11a3-9c8c88d10192`

Note that a response of `true` does not guarantee that the user can view the page, as it does not account for
account-inherited restrictions, space permissions, or even product access. For more
information, see [Confluence permissions](https://confluence.atlassian.com/x/_AozKw).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-content-state`** - Gets the current content state of the draft or current version of content. To specify the draft version, set
the parameter status to draft, otherwise archived or current will get the relevant published state.
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-content-state-settings`** - Get object describing whether content states are allowed at all, if custom content states or space content states
are restricted, and a list of space content states allowed for the space if they are not restricted.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki get-content-template`** - Returns a content template. This includes information about template,
like the name, the space or blueprint that the template is in, the body
of the template, and more.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space to view space templates and permission to
access the Confluence site ('Can use' global permission) to view global templates.
- **`keen wiki get-content-templates`** - Returns all content templates. Use this method to retrieve all global
content templates or all content templates in a space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space to view space templates and permission to
access the Confluence site ('Can use' global permission) to view global templates.
- **`keen wiki get-content-watch-status`** - Returns whether a user is watching a piece of content. Choose the user by
doing one of the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-contents-with-state`** - Returns all content that has the provided content state in a space.

If the expand query parameter is used with the `body.export_view` and/or `body.styled_view` properties, then the query limit parameter will be restricted to a maximum value of 25.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space.
- **`keen wiki get-current-user`** - Returns the currently logged-in user. This includes information about
the user, like the display name, userKey, account ID, profile picture,
and more.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-custom-content-states`** - Get custom content states that authenticated user has created.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**
Must have user authentication.
- **`keen wiki get-descendants-of-type`** - Returns all descendants of a given type, for a piece of content. This is
similar to [Get content children by type](#api-content-id-child-type-get),
except that this method returns child pages at all levels, rather than just
the direct child pages.

A piece of content has different types of descendants, depending on its type:

- `page`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `whiteboard`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `database`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `embed`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `folder`: descendant is `page`, `whiteboard`, `database`, `embed`, `folder`, `comment`, `attachment`
- `blogpost`: descendant is `comment`, `attachment`
- `attachment`: descendant is `comment`
- `comment`: descendant is `attachment`

Custom content types that are provided by apps can also be returned.

If the expand query parameter is used with the `body.export_view` and/or `body.styled_view` properties, then the query limit parameter will be restricted to a maximum value of 25.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space, and permission to view the content if it
is a page.
- **`keen wiki get-global-theme`** - Returns the globally assigned theme.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: None
- **`keen wiki get-group-by-group-id`** - Returns a user group for a given group id.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-group-members-by-group-id`** - Returns the users that are members of a group.

Use updated Get group API

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-group-memberships-for-user`** - Returns the groups that a user is a member of.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-groups`** - Returns all user groups. The returned groups are ordered alphabetically in
ascending order by group name.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-individual-group-restriction-status-by-group-id`** - Returns whether the specified content restriction applies to a group.
For example, if a page with `id=123` has a `read` restriction for the `123456` group id,
the following request will return `true`:

`/wiki/rest/api/content/123/restriction/byOperation/read/byGroupId/123456`

Note that a response of `true` does not guarantee that the group can view the page, as it does not account for
account-inherited restrictions, space permissions, or even product access. For more
information, see [Confluence permissions](https://confluence.atlassian.com/x/_AozKw).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-labels-for-space`** - Returns a list of labels associated with a space. Can provide a prefix as well as other filters to
select different types of labels.
- **`keen wiki get-look-and-feel-settings`** - Returns the look and feel settings for the site or a single space. This
includes attributes such as the color scheme, padding, and border radius.

The look and feel settings for a space can be inherited from the global
look and feel settings or provided by a theme.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
None
- **`keen wiki get-macro-body-by-macro-id`** - Returns the body of a macro in storage format, for the given macro ID.
This includes information like the name of the macro, the body of the macro,
and any macro parameters. This method is mainly used by Cloud apps.

About the macro ID: When a macro is created in a new version of content,
Confluence will generate a random ID for it, unless an ID is specified
(by an app). The macro ID will look similar to this: '50884bd9-0cb8-41d5-98be-f80943c14f96'.
The ID is then persisted as new versions of content are created, and is
only modified by Confluence if there are conflicting IDs.

For Forge macros, the value for macro ID is the "local ID" of that particular ADF node.
This value can be retrieved either client-side by calling view.getContext() and accessing "localId"
on the resulting object, or server-side by examining the "local-id" parameter node inside the "parameters" node.

Note that there are other attributes named "local-id", but only this particular one is used to store the macro ID.

Example:
<ac:adf-node type="extension">
  <ac:adf-attribute key="extension-type">com.atlassian.ecosystem</ac:adf-attribute>
  <ac:adf-attribute key="parameters">
      <ac:adf-parameter key="local-id">e9c4aa10-73fa-417c-888d-48c719ae4165</ac:adf-parameter>
  </ac:adf-parameter>
</ac:adf-node>

Note, to preserve backwards compatibility this resource will also match on
the hash of the macro body, even if a macro ID is found. This check will
eventually become redundant, as macro IDs are generated for pages and
transparently propagate out to all instances.

This backwards compatibility logic does not apply to Forge macros; those
can only be retrieved by their ID.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content that the macro is in.
- **`keen wiki get-modules`** - Returns all modules registered dynamically by the calling app.

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen wiki get-privacy-unsafe-user-email`** - Returns a user's email address regardless of the user’s profile visibility settings. For Connect apps, this API is only available to apps approved by
Atlassian, according to these [guidelines](https://community.developer.atlassian.com/t/guidelines-for-requesting-access-to-email-address/27603).
For Forge apps, this API only supports access via asApp() requests.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-privacy-unsafe-user-email-bulk`** - Returns a user's email address regardless of the user’s profile visibility settings. For Connect apps, this API is only available to apps approved by
Atlassian, according to these [guidelines](https://community.developer.atlassian.com/t/guidelines-for-requesting-access-to-email-address/27603).
For Forge apps, this API only supports access via asApp() requests.

Any accounts which are not available will not be included in the result.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-relationship`** - Find whether a particular type of relationship exists from a source
entity to a target entity. Note, relationships are one way.

For example, you can use this method to find whether the current user has
selected a particular page as a favorite (i.e. 'save for later'):
`GET /wiki/rest/api/relation/favourite/from/user/current/to/content/123`

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view both the target entity and source entity.
- **`keen wiki get-restrictions`** - Returns the restrictions on a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-restrictions-by-operation`** - Returns restrictions on a piece of content by operation. This method is
similar to [Get restrictions](#api-content-id-restriction-get) except that
the operations are properties of the return object, rather than items in
a results array.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-restrictions-for-operation`** - Returns the restictions on a piece of content for a given operation (read
or update).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the content.
- **`keen wiki get-retention-period`** - Returns the retention period for records in the audit log. The retention
period is how long an audit record is kept for, from creation date until
it is deleted.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki get-space-content-states`** - Get content states that are suggested in the space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space.
- **`keen wiki get-space-settings`** - Returns the settings of a space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'View' permission for the space.
- **`keen wiki get-space-theme`** - Returns the theme selected for a space, if one is set. If no space
theme is set, this means that the space is inheriting the global look
and feel settings.

**[Permissions required](https://confluence.atlassian.com/x/_AozKw)**: ‘View’ permission for the space.
- **`keen wiki get-system-info`** - Returns the system information for the Confluence Cloud tenant. This
information is used by Atlassian.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-task`** - Returns information about an active long-running task (e.g. space export),
such as how long it has been running and the percentage of the task that
has completed.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-tasks`** - Returns information about all active long-running tasks (e.g. space export),
such as how long each task has been running and the percentage of each task
that has completed.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-theme`** - Returns a theme. This includes information about the theme name,
description, and icon.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: None
- **`keen wiki get-themes`** - Returns all themes, not including the default theme.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: None
- **`keen wiki get-user`** - Returns a user. This includes information about the user, such as the
display name, account ID, profile picture, and more. The information returned may be
restricted by the user's profile visibility settings.

**Note:** to add, edit, or delete users in your organization, see the
[user management REST API](/cloud/admin/user-management/about/).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-user-properties`** - Returns the properties for a user as list of property keys. For more information
about user properties, see [Confluence entity properties](https://developer.atlassian.com/cloud/confluence/confluence-entity-properties/).
`Note`, these properties stored against a user are on a Confluence site level and not space/content level.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-user-property`** - Returns the property corresponding to `key` for a user. For more information
about user properties, see [Confluence entity properties](https://developer.atlassian.com/cloud/confluence/confluence-entity-properties/).
`Note`, these properties stored against a user are on a Confluence site level and not space/content level.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-viewers`** - Get the total number of distinct viewers a piece of content has.
- **`keen wiki get-views`** - Get the total number of views a piece of content has.
- **`keen wiki get-watchers-for-space`** - Returns a list of watchers of a space
- **`keen wiki get-watches-for-page`** - Returns the watches for a page. A user that watches a page will receive
receive notifications when the page is updated.

If you want to manage watches for a page, use the following `user` methods:

- [Get content watch status for user](#api-user-watch-content-contentId-get)
- [Add content watch](#api-user-watch-content-contentId-post)
- [Remove content watch](#api-user-watch-content-contentId-delete)

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki get-watches-for-space`** - Returns all space watches for the space that the content is in. A user that
watches a space will receive receive notifications when any content in the
space is updated.

If you want to manage watches for a space, use the following `user` methods:

- [Get space watch status for user](#api-user-watch-space-spaceKey-get)
- [Add space watch](#api-user-watch-space-spaceKey-post)
- [Remove space watch](#api-user-watch-space-spaceKey-delete)

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
- **`keen wiki is-watching-label`** - Returns whether a user is watching a label. Choose the user by doing one
of the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki is-watching-space`** - Returns whether a user is watching a space. Choose the user by
doing one of the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki move-page`** - Move a page to a new location relative to a target page:

* `before` - move the page under the same parent as the target, before the target in the list of children
* `after` - move the page under the same parent as the target, after the target in the list of children
* `append` - move the page to be a child of the target

Caution: This API can move pages to the top level of a space. Top-level pages are difficult to find in the UI
because they do not show up in the page tree display. To avoid this, never use `before` or `after` positions
when the `targetId` is a top-level page.
- **`keen wiki publish-legacy-draft`** - Publishes a legacy draft of a page created from a blueprint. Legacy drafts
will eventually be removed in favor of shared drafts. For now, this method
works the same as [Publish shared draft](#api-content-blueprint-instance-draftId-put).

By default, the following objects are expanded: `body.storage`, `history`, `space`, `version`, `ancestors`.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the draft and 'Add' permission for the space that
the content will be created in.
- **`keen wiki publish-shared-draft`** - Publishes a shared draft of a page created from a blueprint.

By default, the following objects are expanded: `body.storage`, `history`, `space`, `version`, `ancestors`.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the draft and 'Add' permission for the space that
the content will be created in.
- **`keen wiki register-modules`** - Registers a list of modules. For the list of modules that support dynamic registration, see [Dynamic modules](https://developer.atlassian.com/cloud/confluence/dynamic-modules/).

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen wiki remove-content-state`** - Removes the content state of the content specified and creates a new version
(publishes the content without changing the body) of the content with the new status.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki remove-content-watcher`** - Removes a user as a watcher from a piece of content. Choose the user by
doing one of the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki remove-group-by-id`** - Delete user group.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a site admin.
- **`keen wiki remove-group-from-content-restriction`** - Removes a group from a content restriction. That is, remove read or update
permission for the group for a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki remove-label-from-content`** - Removes a label from a piece of content. Labels can't be deleted from archived content.
This is similar to [Remove label from content using query parameter](#api-content-id-label-delete)
except that the label name is specified via a path parameter.

Use this method if the label name does not have "/" characters, as the path
parameter does not accept "/" characters for security reasons. Otherwise,
use [Remove label from content using query parameter](#api-content-id-label-delete).

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki remove-label-from-content-using-query-parameter`** - Removes a label from a piece of content. Labels can't be deleted from archived content.
This is similar to [Remove label from content](#api-content-id-label-label-delete)
except that the label name is specified via a query parameter.

Use this method if the label name has "/" characters, as
[Remove label from content using query parameter](#api-content-id-label-delete)
does not accept "/" characters for the label name.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki remove-label-watcher`** - Removes a user as a watcher from a label. Choose the user by doing one of
the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki remove-member-from-group-by-group-id`** - Remove user as a member from a group.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
User must be a site admin.
- **`keen wiki remove-modules`** - Remove all or a list of modules registered by the calling app.

**[Permissions](#permissions) required:** Only Connect apps can make this request.
- **`keen wiki remove-permission`** - Removes a space permission. Note that removing Read Space permission for a user or group will remove all
the space permissions for that user or group.

Note: Apps cannot access this REST resource - including when utilizing user impersonation.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki remove-space-watch`** - Removes a user as a watcher from a space. Choose the user by doing one of
the following:

- Specify a user via a query parameter: Use the `accountId` to identify the user.
- Do not specify a user: The currently logged-in user will be used.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission or 'Space Administrator' permission for the relevant space if specifying a user, otherwise
permission to access the Confluence site ('Can use' global permission).
- **`keen wiki remove-template`** - Deletes a template. This results in different actions depending on the
type of template:

- If the template is a content template, it is deleted.
- If the template is a modified space-level blueprint template, it reverts
to the template inherited from the global-level blueprint template.
- If the template is a modified global-level blueprint template, it reverts
to the default global-level blueprint template.

 Note, unmodified blueprint templates cannot be deleted.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
        'Admin' permission for the space to delete a space template or 'Confluence Administrator'
        global permission to delete a global template.
- **`keen wiki remove-user-from-content-restriction`** - Removes a user from a content restriction. That is, remove read or update
permission for the user for a piece of content.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki reset-look-and-feel-settings`** - Resets the custom look and feel settings for the site or a single space.
This changes the values of the custom settings to be the same as the
default settings. It does not change which settings (default or custom)
are selected. Note, the default space settings are inherited from the
current global settings.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:

- If `spaceKey` is specified, [`manage/look-and-feel`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space.
- If `spaceKey` is omitted, 'Confluence Administrator' global permission.

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen wiki reset-space-theme`** - Resets the space theme. This means that the space will inherit the
global look and feel settings

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki restore-content-version`** - Restores a historical version to be the latest version. That is, a new version
is created with the content of the historical version.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki search-by-cql`** - Searches for content using the
[Confluence Query Language (CQL)](https://developer.atlassian.com/cloud/confluence/advanced-searching-using-cql/).

**Note that CQL input queries submitted through the `/wiki/rest/api/search` endpoint no longer support user-specific fields like `user`, `user.fullname`, `user.accountid`, and `user.userkey`.** 
See this [deprecation notice](https://developer.atlassian.com/cloud/confluence/deprecation-notice-search-api/) for more details.

Example initial call:
```
/wiki/rest/api/search?cql=type=page&limit=25
```

Example response:
```
{
  "results": [
    { ... },
    { ... },
    ...
    { ... }
  ],
  "limit": 25,
  "size": 25,
  ...
  "_links": {
    "base": "<url>",
    "context": "<url>",
    "next": "/rest/api/search?cql=type=page&limit=25&cursor=raNDoMsTRiNg",
    "self": "<url>"
  }
}
```

When additional results are available, returns `next` and `prev` URLs to retrieve them in subsequent calls. The URLs each contain a cursor that points to the appropriate set of results. Use `limit` to specify the number of results returned in each call.

Example subsequent call (taken from example response):
```
/wiki/rest/api/search?cql=type=page&limit=25&cursor=raNDoMsTRiNg
```
The response to this will have a `prev` URL similar to the `next` in the example response.

If the expand query parameter is used with the `body.export_view` and/or `body.styled_view` properties, then the query limit parameter will be restricted to a maximum value of 25.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to view the entities. Note, only entities that the user has
permission to view will be returned.
- **`keen wiki search-content-by-cql`** - Returns the list of content that matches a Confluence Query Language
(CQL) query. For information on CQL, see:
[Advanced searching using CQL](https://developer.atlassian.com/cloud/confluence/advanced-searching-using-cql/).

Example initial call:
```
/wiki/rest/api/content/search?cql=type=page&limit=25
```

Example response:
```
{
  "results": [
    { ... },
    { ... },
    ...
    { ... }
  ],
  "limit": 25,
  "size": 25,
  ...
  "_links": {
    "base": "<url>",
    "context": "<url>",
    "next": "/rest/api/content/search?cql=type=page&limit=25&cursor=raNDoMsTRiNg",
    "self": "<url>"
  }
}
```

When additional results are available, returns `next` and `prev` URLs to retrieve them in subsequent calls. The URLs each contain a cursor that points to the appropriate set of results. Use `limit` to specify the number of results returned in each call.
Example subsequent call (taken from example response):
```
/wiki/rest/api/content/search?cql=type=page&limit=25&cursor=raNDoMsTRiNg
```
The response to this will have a `prev` URL similar to the `next` in the example response.

If the expand query parameter is used with the `body.export_view` and/or `body.styled_view` properties, then the query limit parameter will be restricted to a maximum value of 25.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).
Only content that the user has permission to view will be returned.
- **`keen wiki search-groups`** - Get search results of groups by partial query provided.
- **`keen wiki search-user`** - Searches for users using user-specific queries from the
[Confluence Query Language (CQL)](https://developer.atlassian.com/cloud/confluence/advanced-searching-using-cql/).

Note that CQL input queries submitted through the `/wiki/rest/api/search/user` endpoint only support user-specific fields like `user`, `user.fullname`, `user.accountid`, and `user.userkey`.

Note that some user fields may be set to null depending on the user's privacy settings.
These are: email, profilePicture, displayName, and timeZone.
- **`keen wiki set-content-state`** - Sets the content state of the content specified and creates a new version
(publishes the content without changing the body) of the content with the new state.

You may pass in either an id of a state, or the name and color of a desired new state.
If all 3 are passed in, id will be used.
If the name and color passed in already exist under the current user's existing custom states, the existing state will be reused.
If custom states are disabled in the space of the content (which can be determined by getting the content state space settings of the content's space)
then this set will fail.

You may not remove a content state via this PUT request. You must use the DELETE method. A specified state is required in the body of this request.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki set-retention-period`** - Sets the retention period for records in the audit log. The retention period
can be set to a maximum of 1 year.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Confluence Administrator' global permission.
- **`keen wiki set-space-theme`** - Sets the theme for a space. Note, if you want to reset the space theme to
the default Confluence theme, use the 'Reset space theme' method instead
of this method.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki update-attachment-data`** - Updates the binary data of an attachment, given the attachment ID, and
optionally the comment and the minor edit field.

This method is essentially the same as [Create or update attachments](#api-content-id-child-attachment-put),
except that it matches the attachment ID rather than the name.

Note, you must set a `X-Atlassian-Token: nocheck` header on the request
for this method, otherwise it will be blocked. This protects against XSRF
attacks, which is necessary as this method accepts multipart/form-data.

The media type 'multipart/form-data' is defined in [RFC 7578](https://www.ietf.org/rfc/rfc7578.txt).
Most client libraries have classes that make it easier to implement
multipart posts, like the [MultipartEntityBuilder](https://hc.apache.org/httpcomponents-client-5.1.x/current/httpclient5/apidocs/)
Java class provided by Apache HTTP Components.

Note, according to [RFC 7578](https://tools.ietf.org/html/rfc7578#section-4.5),
in the case where the form data is text,
the charset parameter for the "text/plain" Content-Type may be used to
indicate the character encoding used in that part. In the case of this
API endpoint, the `comment` body parameter should be sent with `type=text/plain`
and `charset=utf-8` values. This will force the charset to be UTF-8.

Example: This curl command updates an attachment (id='att456') that is attached
to a piece of content (id='123') with a comment and `minorEdits`=true.

``` bash
curl -D- \
  -u admin:admin \
  -X POST \
  -H 'X-Atlassian-Token: nocheck' \
  -F 'file=@"example.txt"' \
  -F 'minorEdit="true"' \
  -F 'comment="Example attachment comment"; type=text/plain; charset=utf-8' \
  http://myhost/rest/api/content/123/child/attachment/att456/data
```
**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki update-attachment-properties`** - Updates the attachment properties, i.e. the non-binary data of an attachment
like the filename, media-type, comment, and parent container.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to update the content.
- **`keen wiki update-content-template`** - Updates a content template. Note, blueprint templates cannot be updated
via the REST API.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space to update a space template or 'Confluence Administrator'
global permission to update a global template.
- **`keen wiki update-look-and-feel`** - Sets the look and feel settings to the default (global) settings, the
custom settings, or the current theme's settings for a space.
The custom and theme settings can only be selected if there is already
a theme set for a space. Note, the default space settings are inherited
from the current global settings.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
[`manage/look-and-feel`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space.

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen wiki update-look-and-feel-settings`** - Updates the look and feel settings for the site or for a single space.
If custom settings exist, they are updated. If no custom settings exist,
then a set of custom settings is created.

Note, if a theme is selected for a space, the space look and feel settings
are provided by the theme and cannot be overridden.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:

- If `spaceKey` is specified, [`manage/look-and-feel`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space.
- If `spaceKey` is omitted, 'Confluence Administrator' global permission.

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen wiki update-restrictions`** - Updates restrictions for a piece of content. For each operation in the request,
the supplied `user` and `group` lists replace the corresponding existing lists.
Omitted operations and omitted user or group lists are preserved. An explicitly
empty list (`{"results": []}` or `[]`) clears that list.

Use `read` for viewing restrictions and `update` for editing restrictions.
While content remains view restricted, an editor must also be able to view it.
Removing an editor's `read` restriction while retaining their `update` restriction
returns `400`. Remove their `update` restriction too, or clear all `read`
restrictions to remove direct viewing restrictions. Restrictions inherited
from ancestors still apply.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to edit the content.
- **`keen wiki update-space`** - Updates the name, description, or homepage of a space.

-   For security reasons, permissions cannot be updated via the API and
must be changed via the user interface instead.
-   Currently you cannot set space labels when updating a space.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
'Admin' permission for the space.
- **`keen wiki update-space-settings`** - Updates the settings for a space.

**[Permissions](https://support.atlassian.com/confluence-cloud/docs/what-are-confluences-roles/) required**:
[`manage/space`](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) permission for the space.

**Note:** To find the display name for each permission ID, call the [Get available space permissions](https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space-permissions/#api-space-permissions-get) API.
- **`keen wiki update-user-property`** - Updates a property for the given user. Note, you cannot update the key of a user property, only the value.
For more information about user properties, see
[Confluence entity properties](https://developer.atlassian.com/cloud/confluence/confluence-entity-properties/).
`Note`, these properties stored against a user are on a Confluence site level and not space/content level.

**[Permissions](https://confluence.atlassian.com/x/_AozKw) required**:
Permission to access the Confluence site ('Can use' global permission).

### workflows

This resource represents workflows. Use it to:

 *  Get workflows
 *  Create workflows
 *  Update workflows
 *  Delete inactive workflows
 *  Get workflow capabilities

- **`keen workflows capabilities`** - Get the list of workflow capabilities for a specific workflow using either the workflow ID, or the project and issue type ID pair. The response includes the scope of the workflow, defined as global/project-based, and a list of project types that the workflow is scoped to. It also includes all rules organised into their broad categories (conditions, validators, actions, triggers, screens) as well as the source location (Atlassian-provided, Connect, Forge).

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to access all, including global-scoped, workflows
 *  *Administer projects* project permissions to access project-scoped workflows

The current list of Atlassian-provided rules:

#### Validators ####

A validator rule that checks if a user has the required permissions to execute the transition in the workflow.

##### Permission validator #####

A validator rule that checks if a user has the required permissions to execute the transition in the workflow.

    {
       "ruleKey": "system:check-permission-validator",
       "parameters": {
         "permissionKey": "ADMINISTER_PROJECTS"
       }
     }

Parameters:

 *  `permissionKey` The permission required to perform the transition. Allowed values: [built-in Jira permissions](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-permission-schemes/#built-in-permissions).

##### Parent or child blocking validator #####

A validator to block the child issue's transition depending on the parent issue's status.

    {
       "ruleKey" : "system:parent-or-child-blocking-validator"
       "parameters" : {
         "blocker" : "PARENT"
         "statusIds" : "1,2,3"
       }
     }

Parameters:

 *  `blocker` currently only supports `PARENT`.
 *  `statusIds` a comma-separated list of status IDs.

##### Previous status validator #####

A validator that checks if an issue has transitioned through specified previous status(es) before allowing the current transition to occur.

    {
       "ruleKey": "system:previous-status-validator",
       "parameters": {
         "previousStatusIds": "10014",
         "mostRecentStatusOnly": "true"
       }
     }

Parameters:

 *  `previousStatusIds` a comma-separated list of status IDs, currently only support one ID.
 *  `mostRecentStatusOnly` when `true` only considers the most recent status for the condition evaluation. Allowed values: `true`, `false`.

##### Validate a field value #####

A validation that ensures a specific field's value meets the defined criteria before allowing an issue to transition in the workflow.

Depending on the rule type, the result will vary:

###### Field required ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "fieldRequired",
         "fieldsRequired": "assignee",
         "ignoreContext": "true",
         "errorMessage": "An assignee must be set!"
       }
     }

Parameters:

 *  `fieldsRequired` the ID of the field that is required. For a custom field, it would look like `customfield_123`.
 *  `ignoreContext` controls the impact of context settings on field validation. When set to `true`, the validator doesn't check a required field if its context isn't configured for the current issue. When set to `false`, the validator requires a field even if its context is invalid. Allowed values: `true`, `false`.
 *  `errorMessage` is the error message to display if the user does not provide a value during the transition. A default error message will be shown if you don't provide one (Optional).

###### Field changed ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "fieldChanged",
         "groupsExemptFromValidation": "6862ac20-8672-4f68-896d-4854f5efb79e",
         "fieldKey": "versions",
         "errorMessage": "Affect versions must be modified before transition"
       }
     }

Parameters:

 *  `groupsExemptFromValidation` a comma-separated list of group IDs to be exempt from the validation.
 *  `fieldKey` the ID of the field that has changed. For a custom field, it would look like `customfield_123`.
 *  `errorMessage` the error message to display if the user does not provide a value during the transition. A default error message will be shown if you don't provide one (Optional).

###### Field has a single value ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "fieldHasSingleValue",
         "fieldKey": "created",
         "excludeSubtasks": "true"
       }
     }

Parameters:

 *  `fieldKey` the ID of the field to validate. For a custom field, it would look like `customfield_123`.
 *  `excludeSubtasks` Option to exclude values copied from sub-tasks. Allowed values: `true`, `false`.

###### Field matches regular expression ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "fieldMatchesRegularExpression",
         "regexp": "[0-9]{4}",
         "fieldKey": "description",
         "errorMessage": "Description must contain a 4-digit year"
       }
     }

Parameters:

 *  `regexp` the regular expression used to validate the field's content.
 *  `fieldKey` the ID of the field to validate. For a custom field, it would look like `customfield_123`.
 *  `errorMessage` the error message to display if the field value does not match the regular expression. A default error message will be shown if you don't provide one (Optional).

###### Date field comparison ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "dateFieldComparison",
         "date1FieldKey": "duedate",
         "date2FieldKey": "customfield_10054",
         "includeTime": "true",
         "conditionSelected": ">="
       }
     }

Parameters:

 *  `date1FieldKey` the ID of the first field to compare. For a custom field, it would look like `customfield_123`.
 *  `date2FieldKey` the ID of the second field to compare. For a custom field, it would look like `customfield_123`.
 *  `includeTime` if `true`, compares both date and time. Allowed values: `true`, `false`.
 *  `conditionSelected` the condition to compare with. Allowed values: `>`, `>=`, `=`, `<=`, `<`, `!=`.

###### Date range comparison ######

    {
       "ruleKey": "system:validate-field-value",
       "parameters": {
         "ruleType": "windowDateComparison",
         "date1FieldKey": "customfield_10009",
         "date2FieldKey": "customfield_10054",
         "numberOfDays": "3"
       }
     }

Parameters:

 *  `date1FieldKey` the ID of the first field to compare. For a custom field, it would look like `customfield_123`.
 *  `date2FieldKey` the ID of the second field to compare. For a custom field, it would look like `customfield_123`.
 *  `numberOfDays` maximum number of days past the reference date (`date2FieldKey`) to pass validation.

This rule is composed by aggregating the following legacy rules:

 *  FieldRequiredValidator
 *  FieldChangedValidator
 *  FieldHasSingleValueValidator
 *  RegexpFieldValidator
 *  DateFieldValidator
 *  WindowsDateValidator

##### Pro forma: Forms attached validator #####

Validates that one or more forms are attached to the issue.

    {
       "ruleKey" : "system:proforma-forms-attached"
       "parameters" : {}
     }

##### Proforma: Forms submitted validator #####

Validates that all forms attached to the issue have been submitted.

    {
       "ruleKey" : "system:proforma-forms-submitted"
       "parameters" : {}
     }

#### Conditions ####

Conditions enable workflow rules that govern whether a transition can execute.

##### Check field value #####

A condition rule evaluates as true if a specific field's value meets the defined criteria. This rule ensures that an issue can only transition to the next step in the workflow if the field's value matches the desired condition.

    {
       "ruleKey": "system:check-field-value",
       "parameters": {
         "fieldId": "description",
         "fieldValue": "[\"Done\"]",
         "comparator": "=",
         "comparisonType": "STRING"
       }
     }

Parameters:

 *  `fieldId` The ID of the field to check the value of. For non-system fields, it will look like `customfield_123`. Note: `fieldId` is used interchangeably with the idea of `fieldKey` here, they refer to the same field.
 *  `fieldValue` the list of values to check against the field\\u2019s value.
 *  `comparator` The comparison logic. Allowed values: `>`, `>=`, `=`, `<=`, `<`, `!=`.
 *  `comparisonType` The type of data being compared. Allowed values: `STRING`, `NUMBER`, `DATE`, `DATE_WITHOUT_TIME`, `OPTIONID`.

##### Restrict issue transition #####

This rule ensures that issue transitions are restricted based on user accounts, roles, group memberships, and permissions, maintaining control over who can transition an issue. This condition evaluates as `true` if any of the following criteria is met.

    {
       "ruleKey": "system:restrict-issue-transition",
       "parameters": {
         "accountIds": "allow-reporter,5e68ac137d64450d01a77fa0",
         "roleIds": "10002,10004",
         "groupIds": "703ff44a-7dc8-4f4b-9aa6-a65bf3574fa4",
         "permissionKeys": "ADMINISTER_PROJECTS",
         "groupCustomFields": "customfield_10028",
         "allowUserCustomFields": "customfield_10072,customfield_10144,customfield_10007",
         "denyUserCustomFields": "customfield_10107"
       }
     }

Parameters:

 *  `accountIds` a comma-separated list of the user account IDs. It also allows generic values like: `allow-assignee`, `allow-reporter`, and `accountIds` Note: This is only supported in team-managed projects
 *  `roleIds` a comma-separated list of role IDs.
 *  `groupIds` a comma-separated list of group IDs.
 *  `permissionKeys` a comma-separated list of permission keys. Allowed values: [built-in Jira permissions](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-permission-schemes/#built-in-permissions).
 *  `groupCustomFields` a comma-separated list of group custom field IDs.
 *  `allowUserCustomFields` a comma-separated list of user custom field IDs to allow for issue transition.
 *  `denyUserCustomFields` a comma-separated list of user custom field IDs to deny for issue transition.

This rule is composed by aggregating the following legacy rules:

 *  AllowOnlyAssignee
 *  AllowOnlyReporter
 *  InAnyProjectRoleCondition
 *  InProjectRoleCondition
 *  UserInAnyGroupCondition
 *  UserInGroupCondition
 *  PermissionCondtion
 *  InGroupCFCondition
 *  UserIsInCustomFieldCondition

##### Previous status condition #####

A condition that evaluates based on an issue's previous status(es) and specific criteria.

    {
       "ruleKey" : "system:previous-status-condition"
       "parameters" : {
         "previousStatusIds" : "10004",
         "not": "true",
         "mostRecentStatusOnly" : "true",
         "includeCurrentStatus": "true",
         "ignoreLoopTransitions": "true"
       }
     }

Parameters:

 *  `previousStatusIds` a comma-separated list of status IDs, current only support one ID.
 *  `not` indicates if the condition should be reversed. When `true` it checks that the issue has not been in the selected statuses. Allowed values: `true`, `false`.
 *  `mostRecentStatusOnly` when true only considers the most recent status for the condition evaluation. Allowed values: `true`, `false`.
 *  `includeCurrentStatus` includes the current status when evaluating if the issue has been through the selected statuses. Allowed values: `true`, `false`.
 *  `ignoreLoopTransitions` ignore loop transitions. Allowed values: `true`, `false`.

##### Parent or child blocking condition #####

A condition to block the parent\\u2019s issue transition depending on the child\\u2019s issue status.

    {
       "ruleKey" : "system:parent-or-child-blocking-condition"
       "parameters" : {
         "blocker" : "CHILD",
         "statusIds" : "1,2,3"
       }
     }

Parameters:

 *  `blocker` currently only supports `CHILD`.
 *  `statusIds` a comma-separated list of status IDs.

##### Separation of duties #####

A condition preventing the user from performing, if the user has already performed a transition on the issue.

    {
       "ruleKey": "system:separation-of-duties",
       "parameters": {
         "fromStatusId": "10161",
         "toStatusId": "10160"
       }
     }

Parameters:

 *  `fromStatusId` represents the status ID from which the issue is transitioning. It ensures that the user performing the current transition has not performed any actions when the issue was in the specified status.
 *  `toStatusId` represents the status ID to which the issue is transitioning. It ensures that the user performing the current transition is not the same user who has previously transitioned the issue.

##### Restrict transitions #####

A condition preventing all users from transitioning the issue can also optionally include APIs as well.

    {
       "ruleKey": "system:restrict-from-all-users",
       "parameters": {
         "restrictMode": "users"
       }
     }

Parameters:

 *  `restrictMode` restricts the issue transition including/excluding APIs. Allowed values: `"users"`, `"usersAndAPI"`.

##### Jira Service Management block until approved #####

Block an issue transition until approval. Note: This is only supported in team-managed projects.

    {
       "ruleKey": "system:jsd-approvals-block-until-approved",
       "parameters": {
         "approvalConfigurationJson": "{"statusExternalUuid...}"
       }
     }

Parameters:

 *  `approvalConfigurationJson` a stringified JSON holding the Jira Service Management approval configuration.

##### Jira Service Management block until rejected #####

Block an issue transition until rejected. Note: This is only supported in team-managed projects.

    {
       "ruleKey": "system:jsd-approvals-block-until-rejected",
       "parameters": {
         "approvalConfigurationJson": "{"statusExternalUuid...}"
       }
     }

Parameters:

 *  `approvalConfigurationJson` a stringified JSON holding the Jira Service Management approval configuration.

##### Block in progress approval #####

Condition to block issue transition if there is pending approval. Note: This is only supported in company-managed projects.

    {
       "ruleKey": "system:block-in-progress-approval",
       "parameters": {}
     }

#### Post functions ####

Post functions carry out any additional processing required after a workflow transition is executed.

##### Change assignee #####

A post function rule that changes the assignee of an issue after a transition.

    {
       "ruleKey": "system:change-assignee",
       "parameters": {
         "type": "to-selected-user",
         "accountId": "example-account-id"
       }
     }

Parameters:

 *  `type` the parameter used to determine the new assignee. Allowed values: `to-selected-user`, `to-unassigned`, `to-current-user`, `to-current-user`, `to-default-user`, `to-default-user`
 *  `accountId` the account ID of the user to assign the issue to. This parameter is required only when the type is `"to-selected-user"`.

##### Copy field value #####

A post function that automates the process of copying values between fields during a specific transition, ensuring data consistency and reducing manual effort.

    {
       "ruleKey": "system:copy-value-from-other-field",
       "parameters": {
         "sourceFieldKey": "description",
         "targetFieldKey": "components",
         "issueSource": "SAME"
       }
     }

Parameters:

 *  `sourceFieldKey` the field key to copy from. For a custom field, it would look like `customfield_123`
 *  `targetFieldKey` the field key to copy to. For a custom field, it would look like `customfield_123`
 *  `issueSource` `SAME` or `PARENT`. Defaults to `SAME` if no value is provided.

##### Update field #####

A post function that updates or appends a specific field with the given value.

    {
       "ruleKey": "system:update-field",
       "parameters": {
         "field": "customfield_10056",
         "value": "asdf",
         "mode": "append"
       }
     }

Parameters:

 *  `field` the ID of the field to update. For a custom field, it would look like `customfield_123`
 *  `value` the value to update the field with.
 *  `mode` `append` or `replace`. Determines if a value will be appended to the current value, or if the current value will be replaced.

##### Trigger webhook #####

A post function that automatically triggers a predefined webhook when a transition occurs in the workflow.

    {
       "ruleKey": "system:trigger-webhook",
       "parameters": {
         "webhookId": "1"
       }
     }

Parameters:

 *  `webhookId` the ID of the webhook.

##### Trigger agent #####

A post function that triggers a Jira AI agent for the issue after the transition runs, using the configured agent and an optional prompt. The agent run is asynchronous: it is scheduled after the transition finishes and does not block the transition.

    {
       "ruleKey": "system:trigger-agent",
       "parameters": {
         "agentId": "712020:3c6d3f05-331a-4488-932e-37c34b704720",
         "promptValue": ""
       }
     }

Parameters:

 *  `agentId` the identifier of the agent to trigger (the Atlassian account ID of the agent).
 *  `promptValue` optional text passed to the agent as a user prompt after the transition runs; use an empty string if no extra prompt is needed.

#### Screen ####

##### Remind people to update fields #####

A screen rule that prompts users to update a specific field when they interact with an issue screen during a transition. This rule is useful for ensuring that users provide or modify necessary information before moving an issue to the next step in the workflow.

    {
       "ruleKey": "system:remind-people-to-update-fields",
       "params": {
         "remindingFieldIds": "assignee,customfield_10025",
         "remindingMessage": "The message",
         "remindingAlwaysAsk": "true"
       }
     }

Parameters:

 *  `remindingFieldIds` a comma-separated list of field IDs. Note: `fieldId` is used interchangeably with the idea of `fieldKey` here, they refer to the same field.
 *  `remindingMessage` the message to display when prompting the users to update the fields.
 *  `remindingAlwaysAsk` always remind to update fields. Allowed values: `true`, `false`.

##### Shared transition screen #####

A common screen that is shared between transitions in a workflow.

    {
       "ruleKey": "system:transition-screen",
       "params": {
         "screenId": "3"
       }
     }

Parameters:

 *  `screenId` the ID of the screen.

#### Connect & Forge ####

##### Connect rules #####

Validator/Condition/Post function for Connect app.

    {
       "ruleKey": "connect:expression-validator",
       "parameters": {
         "appKey": "com.atlassian.app",
         "config": "",
         "id": "90ce590f-e90c-4cd3-8281-165ce41f2ac3",
         "disabled": "false",
         "tag": ""
       }
     }

Parameters:

 *  `ruleKey` Validator: `connect:expression-validator`, Condition: `connect:expression-condition`, and Post function: `connect:remote-workflow-function`
 *  `appKey` the reference to the Connect app
 *  `config` a JSON payload string describing the configuration
 *  `id` the ID of the rule
 *  `disabled` determine if the Connect app is disabled. Allowed values: `true`, `false`.
 *  `tag` additional tags for the Connect app

##### Forge rules #####

Validator/Condition/Post function for Forge app.

    {
       "ruleKey": "forge:expression-validator",
       "parameters": {
         "key": "ari:cloud:ecosystem::extension/{appId}/{environmentId}/static/{moduleKey}",
         "config": "{"searchString":"workflow validator"}",
         "id": "a865ddf6-bb3f-4a7b-9540-c2f8b3f9f6c2",
         "disabled": "false",
         "tag": ""
       }
     }

Parameters:

 *  `ruleKey` Validator: `forge:expression-validator`, Condition: `forge:expression-condition`, and Post function: `forge:workflow-post-function`
 *  `key` the identifier for the Forge app
 *  `config` the persistent stringified JSON configuration for the Forge rule
 *  `id` the ID of the Forge rule
 *  `disabled` determine if the Forge app is disabled. Allowed values: `true`, `false`.
 *  `tag` additional tags for the Forge app
- **`keen workflows copy`** - Copies an existing workflow, and the statuses it uses, into a new workflow with the given name. The copy is created in the same scope as the workflow it is copied from. If no description is provided, the copy is created with an empty description.

Copying a workflow requires permission both to read the workflow being copied and to create the copy, which is created in the same scope as its source.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to copy all, including project-scoped, workflows
 *  To copy a project-scoped workflow, either the *Edit workflows* project permission, or both the *View (read-only) workflow* and *Administer projects* project permissions
- **`keen workflows create`** - Create workflows and related statuses.

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to create all, including global-scoped, workflows
 *  *Administer projects* project permissions to create project-scoped workflows
- **`keen workflows get-default-editor`** - Get the user's default workflow editor. This can be either the new editor or the legacy editor.
- **`keen workflows read`** - Returns a list of workflows and related statuses by providing workflow names, workflow IDs, or project and issue types.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to access all, including project-scoped, workflows
 *  At least one of the *Administer projects* and *View (read-only) workflow* project permissions to access project-scoped workflows
- **`keen workflows read-previews`** - Returns a requested workflow within a given project. The response provides a read-only preview of the workflow, omitting full configuration details.

**[Permissions](#permissions) required:**

 *  At least one of the *Administer projects* and *View (read-only) workflow* project permissions
- **`keen workflows search`** - Returns a [paginated](#pagination) list of global and project workflows. If workflow names are specified in the query string, details of those workflows are returned. Otherwise, all workflows are returned.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to access all, including project-scoped, workflows
 *  At least one of the *Administer projects* and *View (read-only) workflow* project permissions to access project-scoped workflows
- **`keen workflows update`** - Update workflows and related statuses.

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to create all, including global-scoped, workflows
 *  *Administer projects* project permissions to create project-scoped workflows
- **`keen workflows validate-create`** - Validate the payload for bulk create workflows.

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to create all, including global-scoped, workflows
 *  *Administer projects* project permissions to create project-scoped workflows
- **`keen workflows validate-update`** - Validate the payload for bulk update workflows.

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to create all, including global-scoped, workflows
 *  *Administer projects* project permissions to create project-scoped workflows

### workflowscheme

Manage workflowscheme

- **`keen workflowscheme assign-scheme-to-project`** - Assigns a workflow scheme to a project. This operation is performed only when there are no issues in the project.

Workflow schemes can only be assigned to classic projects.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme create-workflow-scheme`** - Creates a workflow scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme delete-workflow-scheme`** - Deletes a workflow scheme. Note that a workflow scheme cannot be deleted if it is active (that is, being used by at least one project).

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme get-all-workflow-schemes`** - Returns a [paginated](#pagination) list of all workflow schemes, not including draft workflow schemes.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme get-required-workflow-scheme-mappings`** - Gets the required status mappings for the desired changes to a workflow scheme. The results are provided per issue type and workflow. When updating a workflow scheme, status mappings can be provided per issue type, per workflow, or both.

**[Permissions](#permissions) required:**

 *  *Administer Jira* permission to update all, including global-scoped, workflow schemes.
 *  *Administer projects* project permission to update project-scoped workflow schemes.
- **`keen workflowscheme get-workflow-scheme`** - Returns a workflow scheme.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme get-workflow-scheme-project-associations`** - Returns a list of the workflow schemes associated with a list of projects. Each returned workflow scheme includes a list of the requested projects associated with it. Any team-managed or non-existent projects in the request are ignored and no errors are returned.

If the project is associated with the `Default Workflow Scheme` no ID is returned. This is because the way the `Default Workflow Scheme` is stored means it has no ID.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme read-workflow-schemes`** - Returns a list of workflow schemes by providing workflow scheme IDs or project IDs.

**[Permissions](#permissions) required:**

 *  *Administer Jira* global permission to access all, including project-scoped, workflow schemes
 *  *Administer projects* project permissions to access project-scoped workflow schemes
- **`keen workflowscheme switch-workflow-scheme-for-project`** - Switches a workflow scheme for a project.

Workflow schemes can only be assigned to classic projects.

**Calculating required mappings:** If statuses from the current workflow scheme won't exist in the target workflow scheme, you must provide `mappingsByIssueTypeOverride` to specify how issues with those statuses should be migrated. Use [the required workflow scheme mappings API](#api-rest-api-3-workflowscheme-update-mappings-post) to determine which statuses and issue types require mappings.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).
- **`keen workflowscheme update-schemes`** - Updates company-managed and team-managed project workflow schemes. This API doesn't have a concept of draft, so any changes made to a workflow scheme are immediately available. When changing the available statuses for issue types, an [asynchronous task](#async) migrates the issues as defined in the provided mappings.

**[Permissions](#permissions) required:**

 *  *Administer Jira* project permission to update all, including global-scoped, workflow schemes.
 *  *Administer projects* project permission to update project-scoped workflow schemes.
- **`keen workflowscheme update-workflow-scheme`** - Updates a company-manged project workflow scheme, including the name, default workflow, issue type to project mappings, and more. If the workflow scheme is active (that is, being used by at least one project), then a draft workflow scheme is created or updated instead, provided that `updateDraftIfNeeded` is set to `true`.

**[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.atlassian.com/x/x4dKLg).

### worklog

Manage worklog

- **`keen worklog get-for-ids`** - Returns worklog details for a list of worklog IDs.

The returned list of worklogs is limited to 1000 items.

**[Permissions](#permissions) required:** Permission to access Jira, however, worklogs are only returned where either of the following is true:

 *  the worklog is set as *Viewable by All Users*.
 *  the user is a member of a project role or group with permission to view the worklog.
- **`keen worklog get-ids-of-deleted-since`** - Returns a list of IDs and delete timestamps for worklogs deleted after a date and time.

This resource is paginated, with a limit of 1000 worklogs per page. Each page lists worklogs from oldest to youngest. If the number of items in the date range exceeds 1000, `until` indicates the timestamp of the youngest item on the page. Also, `nextPage` provides the URL for the next page of worklogs. The `lastPage` parameter is set to true on the last page of worklogs.

This resource does not return worklogs deleted during the minute preceding the request.

**[Permissions](#permissions) required:** Permission to access Jira.
- **`keen worklog get-ids-of-modified-since`** - Returns a list of IDs and update timestamps for worklogs updated after a date and time.

This resource is paginated, with a limit of 1000 worklogs per page. Each page lists worklogs from oldest to youngest. If the number of items in the date range exceeds 1000, `until` indicates the timestamp of the youngest item on the page. Also, `nextPage` provides the URL for the next page of worklogs. The `lastPage` parameter is set to true on the last page of worklogs.

This resource does not return worklogs updated during the minute preceding the request.

**[Permissions](#permissions) required:** Permission to access Jira, however, worklogs are only returned where either of the following is true:

 *  the worklog is set as *Viewable by All Users*.
 *  the user is a member of a project role or group with permission to view the worklog.

### workspaces

Org Workspaces APIs

- **`keen workspaces <orgId>`** - A workspace refers to a specific instance of an Atlassian product that is accessed through a unique URL. Whenever a user initiates or adds a new product instance, it results in the creation of a distinct workspace.

This API will:
- Return a paginated list of workspaces in a given org
- Return more details about an organization's products (including product URL).

#### Scopes
**[Authorization scopes](/cloud/admin/scopes/) required:** `read:workspaces:admin`


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
keen addon-categories get mock-value

# JSON for scripting and agents
keen addon-categories get mock-value --json

# Filter to specific fields
keen addon-categories get mock-value --json --select id,name,status

# Dry run — show the request without sending
keen addon-categories get mock-value --dry-run

# Agent mode — JSON + compact + no prompts in one flag
keen addon-categories get mock-value --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Explicit retries** - add `--idempotent` to create retries and `--ignore-missing` to delete retries when a no-op success is acceptable
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - write commands can accept structured input when their help lists `--stdin`
- **Offline-friendly** - sync/search commands can use the local SQLite store when available
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
keen doctor
```

Verifies configuration, credentials, and connectivity to the API.

## Configuration

Config file: `~/.config/keen/config.toml`

Static request headers can be configured under `headers`; per-command header overrides take precedence.

Environment variables:

| Name | Kind | Required | Description |
| --- | --- | --- | --- |
| `JIRA_CLOUD_PLATFORM_USERNAME` | per_call | Yes |  |
| `JIRA_CLOUD_PLATFORM_PASSWORD` | per_call | Yes | Set to your API credential. |

## Troubleshooting
**Authentication errors (exit code 4)**
- Run `keen doctor` to check credentials
- Verify the environment variable is set: `echo $JIRA_CLOUD_PLATFORM_USERNAME`
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
