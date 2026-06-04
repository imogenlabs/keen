---
name: pp-jira-pp-cli
description: "Printing Press CLI for Jira Pp Cli. Combined CLI for multiple API services"
author: "smashingtags"
license: "Apache-2.0"
argument-hint: "<command> [args] | install cli|mcp"
allowed-tools: "Read Bash"
metadata:
  openclaw:
    requires:
      bins:
        - keen
---

# Jira Pp Cli — Printing Press CLI

## Prerequisites: Install the CLI

This skill drives the `keen` binary. **You must verify the CLI is installed before invoking any command from this skill.** If it is missing, install it first:

1. Install via the Printing Press installer:
   ```bash
   npx -y @mvanhorn/printing-press-library install jira-pp-cli --cli-only
   ```
2. Verify: `keen --version`
3. Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is on `$PATH`.

If the `npx` install fails before this CLI has a public-library category, install Node or use the category-specific Go fallback after publish.

If `--version` reports "command not found" after install, the install step did not put the binary on `$PATH`. Do not proceed with skill commands until verification succeeds.

Combined CLI for multiple API services

## Command Reference

**addon-categories** — Manage addon categories

- `keen addon-categories get` — Get a specific category.
- `keen addon-categories get-addoncategories` — Get a list of categories associated with the parent application.

**addons** — Manage addons

- `keen addons create` — Create a new app. The associated vendor account must exist prior to the app's creation.
- `keen addons get` — Get a specific app.
- `keen addons get-archived` — Get a list of archived apps by a specific vendor matching the specified parameters.
- `keen addons get-vendor` — Get a list of apps for the specified vendor.
- `keen addons list` — Get a list of apps matching the specified parameters.
- `keen addons list-archived` — Get a list of archived apps matching the specified parameters.
- `keen addons list-listings` — Get a list of app banners matching the specified parameters.
- `keen addons list-search` — Get a list of apps with names matching a search term.
- `keen addons update` — Update a specific app. The request body must be a valid JSON Patch document.

**agile** — Manage agile

- `keen agile create-board` — Creates a new board. Board name, type and filter ID is required. * `name` - Must be less than 255 characters.
- `keen agile create-sprint` — Creates a future sprint. Sprint name and origin board id are required. Start date, end date, and goal are optional.
- `keen agile delete-board` — Deletes the board. Admin without the view permission can still remove the board.
- `keen agile delete-board-property` — Removes the property from the board identified by the id.
- `keen agile delete-property` — Removes the property from the sprint identified by the id.
- `keen agile delete-sprint` — Deletes a sprint. Once a sprint is deleted, all open issues in the sprint will be moved to the backlog.
- `keen agile estimate-issue-for-board` — Updates the estimation of the issue. boardId param is required.
- `keen agile get-all-boards` — Returns all boards. This only includes boards that the user has permission to view.
- `keen agile get-all-quick-filters` — Returns all quick filters from a board, for a given board ID.
- `keen agile get-all-sprints` — Returns all sprints from a board, for a given board ID. This only includes sprints that the user has permission to view.
- `keen agile get-all-versions` — Returns all versions from a board, for a given board ID.
- `keen agile get-board` — Returns the board for the given board ID. This board will only be returned if the user has permission to view it.
- `keen agile get-board-by-filter-id` — Returns any boards which use the provided filter id.
- `keen agile get-board-issues-for-epic` — Returns all issues that belong to an epic on the board, for the given epic ID and the board ID.
- `keen agile get-board-issues-for-sprint` — Get all issues you have access to that belong to the sprint from the board.
- `keen agile get-board-property` — Returns the value of the property with a given key from the board identified by the provided id.
- `keen agile get-board-property-keys` — Returns the keys of all properties for the board identified by the id.
- `keen agile get-configuration` — Get the board configuration. The response contains the following fields: * `id` - ID of the board.
- `keen agile get-epic` — Returns the epic for a given epic ID. This epic will only be returned if the user has permission to view it.
- `keen agile get-epics` — Returns all epics from the board, for the given board ID. This only includes epics that the user has permission to view.
- `keen agile get-features-for-board` — Get features for board
- `keen agile get-issue` — Returns a single issue, for a given issue ID or issue key.
- `keen agile get-issue-estimation-for-board` — Returns the estimation of the issue and a fieldId of the field that is used for it. `boardId` param is required.
- `keen agile get-issues-for-backlog` — Returns all issues from the board's backlog, for the given board ID.
- `keen agile get-issues-for-board` — Returns all issues from a board, for a given board ID. This only includes issues that the user has permission to view.
- `keen agile get-issues-for-epic` — Returns all issues that belong to the epic, for the given epic ID.
- `keen agile get-issues-for-sprint` — Returns all issues in a sprint, for a given sprint ID. This only includes issues that the user has permission to view.
- `keen agile get-issues-without-epic` — Returns all issues that do not belong to any epic. This only includes issues that the user has permission to view.
- `keen agile get-issues-without-epic-for-board` — Returns all issues that do not belong to any epic on a board, for a given board ID.
- `keen agile get-projects` — Returns all projects that are associated with the board, for the given board ID.
- `keen agile get-projects-full` — Returns all projects that are statically associated with the board, for the given board ID.
- `keen agile get-properties-keys` — Returns the keys of all properties for the sprint identified by the id.
- `keen agile get-property` — Returns the value of the property with a given key from the sprint identified by the provided id.
- `keen agile get-quick-filter` — Returns the quick filter for a given quick filter ID.
- `keen agile get-reports-for-board` — Get reports for board
- `keen agile get-sprint` — Returns the sprint for a given sprint ID.
- `keen agile move-issues-to-backlog` — Move issues to the backlog. This operation is equivalent to remove future and active sprints from a given set of issues.
- `keen agile move-issues-to-backlog-for-board` — Move issues to the backlog of a particular board (if they are already on that board).
- `keen agile move-issues-to-board` — Move issues from the backog to the board (if they are already in the backlog of that board).
- `keen agile move-issues-to-epic` — Moves issues to an epic, for a given epic id. Issues can be only in a single epic at the same time.
- `keen agile move-issues-to-sprint-and-rank` — Moves issues to a sprint, for a given sprint ID. Issues can only be moved to open or active sprints.
- `keen agile partially-update-epic` — Performs a partial update of the epic.
- `keen agile partially-update-sprint` — Performs a partial update of a sprint.
- `keen agile rank-epics` — Moves (ranks) an epic before or after a given epic.
- `keen agile rank-issues` — Moves (ranks) issues before or after a given issue. At most 50 issues may be ranked at once.
- `keen agile remove-issues-from-epic` — Removes issues from epics.
- `keen agile set-board-property` — Sets the value of the specified board's property.
- `keen agile set-property` — Sets the value of the specified sprint's property.
- `keen agile swap-sprint` — Swap the position of the sprint with the second sprint.
- `keen agile toggle-features` — Toggle features
- `keen agile update-sprint` — Performs a full update of a sprint. A full update means that the result will be exactly the same as the request body.

**announcement-banner** — This resource represents an announcement banner. Use it to retrieve and update banner configuration.

- `keen announcement-banner get-banner` — Returns the current announcement banner configuration.
- `keen announcement-banner set-banner` — Updates the announcement banner configuration.

**app** — Manage app

- `keen app get-custom-field-configuration` — Returns a [paginated](#pagination) list of configurations for a custom field of a [type](https://developer.atlassian.
- `keen app get-custom-fields-configurations` — Returns a [paginated](#pagination) list of configurations for list of custom fields of a [type](https://developer.
- `keen app update-custom-field-configuration` — Update the configuration for contexts of a custom field of a [type](https://developer.atlassian.
- `keen app update-custom-field-value` — Updates the value of a custom field on one or more issues.
- `keen app update-multiple-custom-field-values` — Updates the value of one or more custom fields on one or more issues.

**app-access-settings** — App Access Settings APIs

- `keen app-access-settings create-domain` — Registers a new approved-domain configuration for your organization, including the initial per-product configuration.
- `keen app-access-settings get-domain` — Returns the app access settings configuration for a specific domain in your organization.
- `keen app-access-settings list-domains` — Returns a paginated list of domain configurations for your organization's app access settings.
- `keen app-access-settings replace-domain` — Replaces the full per-product configuration for a domain.
- `keen app-access-settings upsert-product` — Adds or replaces a single product's access configuration on a domain.

**application-properties** — Manage application properties

- `keen application-properties get-advanced-settings` — Returns the application properties that are accessible on the *Advanced Settings* page.
- `keen application-properties get-application-property` — Returns all application properties or an application property.
- `keen application-properties set-application-property` — Changes the value of an application property. For example, you can change the value of the `jira.clone.

**applicationrole** — Manage applicationrole

- `keen applicationrole get-all-application-roles` — Returns all application roles.
- `keen applicationrole get-application-role` — Returns an application role.

**applications** — Information about applications

- `keen applications get` — Get a specific application.
- `keen applications list` — Get a list of applications.

**assets** — For uploading app/image files

- `keen assets create` — Use this resource to upload an app artifact (`.jar` or `.obr`) to Marketplace, so it can be used by an app version.
- `keen assets create-artifact` — Use this resource to copy and validate an app artifact (`.jar` or `.
- `keen assets create-image` — Use this resource to upload an image (`.jpg`, `.png`, `.gif`, `.
- `keen assets delete` — Removes a file from Marketplace storage which had been previously uploaded. This is only allowed if: 1.
- `keen assets get` — Returns details about a previously uploaded file that Marketplace is now storing.
- `keen assets list` — Get links to permitted asset resources.

**atlassian-access-events** — Events APIs

- `keen atlassian-access-events get` — Returns a filtered list of audit log events for an organization.
- `keen atlassian-access-events get-by-id` — Returns information about a single event by ID. #### Scopes **[OAuth 2.

**atlassian-access-groups** — Orgs Groups APIs

- `keen atlassian-access-groups <orgId>` — **This API is deprecated and will no longer work after June 30, 2026.

**atlassian-access-users** — Orgs Users APIs

- `keen atlassian-access-users create` — **The API is presently accessible exclusively to customers who hold a paid subscription.
- `keen atlassian-access-users create-orgs` — Invite people to your organization. When you invite someone: - they’re given app roles according to your invitation.
- `keen atlassian-access-users get` — Returns a list of managed accounts in an organization.
- `keen atlassian-access-users search` — **This API is deprecated and will no longer work after June 30, 2026.

**atlassian-connect** — Manage atlassian connect

- `keen atlassian-connect addon-properties-resource-delete-addon-property-delete` — Deletes an app's property.
- `keen atlassian-connect addon-properties-resource-get-addon-properties-get` — Gets all the properties of an app.
- `keen atlassian-connect addon-properties-resource-get-addon-property-get` — Returns the key and value of an app's property.
- `keen atlassian-connect addon-properties-resource-put-addon-property-put` — Sets the value of an app's property. Use this resource to store custom data for your app.
- `keen atlassian-connect app-issue-field-value-update-resource-update-issue-fields-put` — Updates the value of a custom field added by Connect apps on one or more issues.
- `keen atlassian-connect connect-to-forge-migration-fetch-task-resource-fetch-migration-task-get` — Returns the details of a Connect issue field's migration to Forge.
- `keen atlassian-connect connect-to-forge-migration-task-submission-resource-submit-task-post` — Submits a request to trigger migration of connect issue field to its Forge custom field counterpart.
- `keen atlassian-connect dynamic-modules-resource-get-modules-get` — Returns all modules registered dynamically by the calling app.
- `keen atlassian-connect dynamic-modules-resource-register-modules-post` — Registers a list of modules. **[Permissions](#permissions) required:** Only Connect apps can make this request.
- `keen atlassian-connect dynamic-modules-resource-remove-modules-delete` — Remove all or a list of modules registered by the calling app.
- `keen atlassian-connect migration-resource-update-entity-properties-value-put` — Updates the values of multiple entity properties for an object, up to 50 updates per request.
- `keen atlassian-connect migration-resource-workflow-rule-search-post` — Returns configurations for workflow transition rules migrated from server to cloud and owned by the calling Connect app.
- `keen atlassian-connect service-registry-resource-services-get` — Retrieve the attributes of given service registries.

**attachment** — Manage attachment

- `keen attachment get` — Returns the metadata for an attachment. Note that the attachment itself is not returned.
- `keen attachment get-content` — Returns the contents of an attachment.
- `keen attachment get-meta` — Returns the attachment settings, that is, whether attachments are enabled and the maximum attachment size allowed.
- `keen attachment get-thumbnail` — Returns the thumbnail of an attachment.
- `keen attachment remove` — Deletes an attachment from an issue. This operation can be accessed anonymously.

**auditing** — Manage auditing

- `keen auditing` — Returns a list of audit records.

**avatar** — This resource represents system and custom avatars. Use it to obtain the details of system or custom avatars, add and remove avatars from a project, issue type or priority and obtain avatar images.


**builds** — APIs related to integrating build data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraBuildInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/build/. They are also related to integrating Jira Software Cloud with on-premises tools using OAuth 2.0 credentials. See https://developer.atlassian.com/cloud/jira/software/integrate-jsw-cloud-with-onpremises-tools/.

These APIs are available to Forge apps with the `devops:buildInfoProvider` module in the Forge app's manifest. See https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-software-build-info/.

- `keen builds delete-by-key` — Delete the build data currently stored for the given `pipelineId` and `buildNumber` combination.
- `keen builds delete-by-property` — Bulk delete all builds data that match the given request.
- `keen builds get-by-key` — Retrieve the currently stored build data for the given `pipelineId` and `buildNumber` combination.
- `keen builds submit` — Update / insert builds data.

**changelog** — Manage changelog

- `keen changelog` — Bulk fetch changelogs for multiple issues and filter by fields Returns a paginated list of all changelogs for given

**classification-levels** — This resource represents classification levels.

- `keen classification-levels` — Returns all classification levels. **[Permissions](#permissions) required:** None.

**comment** — Manage comment

- `keen comment` — Returns a [paginated](#pagination) list of comments specified by a list of comment IDs.

**component** — Manage component

- `keen component create` — Creates a component. Use components to provide containers for issues within a project.
- `keen component delete` — Deletes a component. This operation can be accessed anonymously.
- `keen component find-for-projects` — Returns a [paginated](#pagination) list of all components in a project, including global (Compass)
- `keen component get` — Returns a component. This operation can be accessed anonymously.
- `keen component update` — Updates a component. Any fields included in the request are overwritten.

**config** — Manage config

- `keen config associate-projects-to-field-association-schemes` — Associate projects to field association schemes.
- `keen config clone-field-association-scheme` — Endpoint for cloning an existing field association scheme into a new one.
- `keen config create-field-association-scheme` — Endpoint for creating a new field association scheme.
- `keen config delete-field-association-scheme` — Delete a specified field association scheme **[Permissions](#permissions) required
- `keen config get-field-association-scheme-by-id` — Endpoint for fetching a field association scheme by its ID **[Permissions](#permissions) required
- `keen config get-field-association-scheme-item-parameters` — Retrieve field association parameters on a field association scheme **[Permissions](#permissions) required
- `keen config get-field-association-schemes` — REST endpoint for retrieving a paginated list of field association schemes with optional filtering.
- `keen config get-projects-with-field-schemes` — Get projects with field association schemes.
- `keen config remove-field-association-scheme-item-parameters` — Remove field association parameters overrides for work types.
- `keen config remove-fields-associated-with-schemes` — Remove fields associated with field association schemes.
- `keen config search-field-association-scheme-fields` — Search for fields belonging to a given field association scheme.
- `keen config search-field-association-scheme-projects` — REST Endpoint for searching for projects belonging to a given field association scheme **[Permissions](#permissions)
- `keen config update-field-association-scheme` — Endpoint for updating an existing field association scheme.
- `keen config update-field-association-scheme-item-parameters` — Update field association item parameters in field association schemes.
- `keen config update-fields-associated-with-schemes` — Update fields associated with field association schemes.

**configuration** — Manage configuration

- `keen configuration get` — Returns the [global settings](https://confluence.atlassian.com/x/qYXKM) in Jira.
- `keen configuration get-available-time-tracking-implementations` — Returns all time tracking providers.
- `keen configuration get-selected-time-tracking-implementation` — Returns the time tracking provider that is currently selected.
- `keen configuration get-shared-time-tracking` — Returns the time tracking settings. This includes settings such as the time format, default time unit, and others.
- `keen configuration select-time-tracking-implementation` — Selects a time tracking provider.
- `keen configuration set-shared-time-tracking` — Sets the time tracking settings.

**custom-field-option** — Manage custom field option

- `keen custom-field-option <id>` — Returns a custom field option. For example, an option in a select list.

**dashboard** — This resource represents dashboards. Use it to obtain the details of dashboards as well as get, create, update, or remove item properties and gadgets from dashboards.

- `keen dashboard bulk-edit` — Bulk edit dashboards. Maximum number of dashboards to be edited at the same time is 100.
- `keen dashboard create` — Creates a dashboard. **[Permissions](#permissions) required:** None.
- `keen dashboard delete` — Deletes a dashboard.
- `keen dashboard get` — Returns a dashboard. This operation can be accessed anonymously. **[Permissions](#permissions) required:** None.
- `keen dashboard get-all` — Returns a list of dashboards owned by or shared with the user.
- `keen dashboard get-all-available-gadgets` — Gets a list of all available gadgets that can be added to all dashboards.
- `keen dashboard get-paginated` — Returns a [paginated](#pagination) list of dashboards.
- `keen dashboard update` — Updates a dashboard, replacing all the dashboard details with those provided.

**data-policy** — Manage data policy

- `keen data-policy get-policies` — Returns data policies for the projects specified in the request.
- `keen data-policy get-policy` — Returns data policy for the workspace.

**deployments** — APIs related to integrating deployment data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraDeploymentInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/deployment/. They are also related to integrating Jira Software Cloud with on-premises tools using OAuth 2.0 credentials. See https://developer.atlassian.com/cloud/jira/software/integrate-jsw-cloud-with-onpremises-tools/.

These APIs are available to Forge apps with the `devops:deploymentInfoProvider` module in the Forge app's manifest. See https://developer.atlassian.com/platform/forge/manifest-reference/modules/jira-software-deployment-info/.

- `keen deployments delete-by-key` — Delete the currently stored deployment data for the given `pipelineId`
- `keen deployments delete-by-property` — Bulk delete all deployments that match the given request.
- `keen deployments get-by-key` — Retrieve the currently stored deployment data for the given `pipelineId`
- `keen deployments get-gating-status-by-key` — Retrieve the Deployment gating status for the given `pipelineId + environmentId + deploymentSequenceNumber` combination.
- `keen deployments submit` — Update / insert deployment data.

**devinfo** — Manage devinfo

- `keen devinfo delete-by-properties` — Deletes development information entities which have all the provided properties.
- `keen devinfo delete-entity` — Deletes particular development information entity. Deletion is performed asynchronously.
- `keen devinfo delete-repository` — Deletes the repository data stored by the given ID and all related development information entities.
- `keen devinfo exists-by-properties` — Checks if repositories which have all the provided properties exists.
- `keen devinfo get-repository` — For the specified repository ID, retrieves the repository and the most recent 400 development information entities.
- `keen devinfo store-development-information` — Stores development information provided in the request to make it available when viewing issues in Jira.

**devopscomponents** — Manage devopscomponents

- `keen devopscomponents delete-component-by-id` — Delete the Component data currently stored for the given ID. Deletion is performed asynchronously.
- `keen devopscomponents delete-components-by-property` — Bulk delete all Components that match the given request.
- `keen devopscomponents get-component-by-id` — Retrieve the currently stored Component data for the given ID.
- `keen devopscomponents submit-components` — Update / insert DevOps Component data.

**directories** — Manage directories

- `keen directories <orgId>` — Returns a page of directories in an organization that match the supplied parameters. #### Scopes **[OAuth 2.

**directory** — Org Directory APIs

- `keen directory assign-role-to-group` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory create` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory create-orgs` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory create-orgs-2` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory create-orgs-3` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory delete` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory delete-orgs` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory delete-orgs-2` — **This API is deprecated and will no longer work after June 30, 2026.
- `keen directory get` — **Additional response parameters of the API (for e.g.
- `keen directory revoke-role-to-group` — **This API is deprecated and will no longer work after June 30, 2026.

**domains** — Domain APIs

- `keen domains get` — Returns a list of domains in an organization one page at a time. #### Scopes **[OAuth 2.
- `keen domains get-by-id` — Returns information about a single verified domain by ID. #### Scopes **[OAuth 2.

**event-actions** — Manage event actions

- `keen event-actions <orgId>` — Returns information localized event actions #### Scopes **[OAuth 2.

**events** — Manage events

- `keen events` — Returns all issue events.

**events-stream** — Manage events stream

- `keen events-stream <orgId>` — Returns a paginated list of audit logs events for an organization.

**expression** — Manage expression

- `keen expression analyse` — Analyses and validates Jira expressions.
- `keen expression evaluate-jira` — Endpoint is currently being removed. [More details](https://developer.atlassian.
- `keen expression evaluate-jsisjira` — Evaluates a Jira expression and returns its value.

**featureflags** — Manage featureflags

- `keen featureflags delete-feature-flag-by-id` — Delete the Feature Flag data currently stored for the given ID. Deletion is performed asynchronously.
- `keen featureflags delete-feature-flags-by-property` — Bulk delete all Feature Flags that match the given request.
- `keen featureflags get-feature-flag-by-id` — Retrieve the currently stored Feature Flag data for the given ID.
- `keen featureflags submit-feature-flags` — Update / insert Feature Flag data.

**field** — Manage field

- `keen field create-associations` — Associates fields with projects. Fields will be associated with each issue type on the requested projects.
- `keen field create-custom` — Creates a custom field.
- `keen field delete-custom` — Deletes a custom field. The custom field is deleted whether it is in the trash or not.
- `keen field get` — Returns system and custom issue fields according to the following rules
- `keen field get-paginated` — Returns a [paginated](#pagination) list of fields for Classic Jira projects.
- `keen field get-trashed-paginated` — Returns a [paginated](#pagination) list of fields in the trash.
- `keen field remove-associations` — Unassociates a set of fields with a project and issue type context.
- `keen field update-custom` — Updates a custom field.

**fieldconfiguration** — Manage fieldconfiguration

- `keen fieldconfiguration create-field-configuration` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfiguration delete-field-configuration` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfiguration get-all-field-configurations` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfiguration update-field-configuration` — Deprecated, use [ Field schemes](https://developer.atlassian.

**fieldconfigurationscheme** — Manage fieldconfigurationscheme

- `keen fieldconfigurationscheme assign-field-configuration-scheme-to-project` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme create-field-configuration-scheme` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme delete-field-configuration-scheme` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme get-all-field-configuration-schemes` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme get-field-configuration-scheme-mappings` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme get-field-configuration-scheme-project-mapping` — Deprecated, use [ Field schemes](https://developer.atlassian.
- `keen fieldconfigurationscheme update-field-configuration-scheme` — Deprecated, use [ Field schemes](https://developer.atlassian.

**filter** — This resource represents [filters](https://confluence.atlassian.com/x/eQiiLQ). Use it to get, create, update, or delete filters. Also use it to configure the columns for a filter and set favorite filters.

- `keen filter create` — Creates a filter. The filter is shared according to the [default share scope](#api-rest-api-2-filter-post).
- `keen filter delete` — Delete a filter.
- `keen filter get` — Returns a filter. This operation can be accessed anonymously.
- `keen filter get-default-share-scope` — Returns the default sharing settings for new filters and dashboards for a user.
- `keen filter get-favourite` — Returns the visible favorite filters of the user. This operation can be accessed anonymously.
- `keen filter get-my` — Returns the filters owned by the user.
- `keen filter get-paginated` — Returns a [paginated](#pagination) list of filters.
- `keen filter set-default-share-scope` — Sets the default sharing for new filters and dashboards for a user.
- `keen filter update` — Updates a filter. Use this operation to update a filter's name, description, JQL, or sharing.

**forge** — Manage forge

- `keen forge bulk-pin-unpin-projects-async` — Bulk pin or unpin an issue panel (added by a Forge app) to or from multiple projects. The operation runs asynchronously.
- `keen forge delete-app-property` — Deletes a Forge app's property. **[Permissions](#permissions) required:** Only Forge apps can make this request.
- `keen forge get-app-property` — Returns the value of a Forge app's property.
- `keen forge get-app-property-keys` — Returns all property keys for the Forge app.
- `keen forge put-app-property` — Sets the value of a Forge app's property.

**group** — This resource represents groups of users. Use it to get, create, find, and delete groups as well as add and remove users from groups. (\[WARNING\] The standard Atlassian group names are default names only and can be edited or deleted. For example, an admin or Atlassian support could delete the default group jira-software-users or rename it to jsw-users at any point. See https://support.atlassian.com/user-management/docs/create-and-update-groups/ for details.)

- `keen group add-user-to` — Adds a user to a group.
- `keen group bulk-get` — Returns a [paginated](#pagination) list of groups.
- `keen group create` — Creates a group.
- `keen group get` — This operation is deprecated, use [`group/member`](#api-rest-api-2-group-member-get). Returns all users in a group.
- `keen group get-users-from` — Returns a [paginated](#pagination) list of all users in a group.
- `keen group remove` — Deletes a group.
- `keen group remove-user-from` — Removes a user from a group.

**groups** — This resource represents groups of users. Use it to get, create, find, and delete groups as well as add and remove users from groups. (\[WARNING\] The standard Atlassian group names are default names only and can be edited or deleted. For example, an admin or Atlassian support could delete the default group jira-software-users or rename it to jsw-users at any point. See https://support.atlassian.com/user-management/docs/create-and-update-groups/ for details.)

- `keen groups` — Returns a list of groups whose names contain a query string.

**groupuserpicker** — Manage groupuserpicker

- `keen groupuserpicker` — Returns a list of users and groups matching a string.

**instance** — Manage instance

- `keen instance` — Returns licensing information about the Jira instance. **[Permissions](#permissions) required:** None.

**internal** — Manage internal

- `keen internal` — Returns worklog details for a list of issue ID and worklog ID pairs.

**issue** — This resource represents Jira issues. Use it to:

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

- `keen issue archive` — Enables admins to archive up to 1000 issues in a single request using issue ID/key, returning details of the issue(s)
- `keen issue archive-async` — Enables admins to archive up to 100,000 issues in a single request using JQL
- `keen issue bulk-delete-property` — Deletes a property value from multiple issues. The issues to be updated can be specified by filter criteria.
- `keen issue bulk-fetch` — Returns the details for a set of requested issues. You can request up to 100 issues.
- `keen issue bulk-set-properties-by` — Sets or updates entity property values on issues.
- `keen issue bulk-set-properties-list` — Sets or updates a list of entity property values on issues.
- `keen issue bulk-set-property` — Sets a property value on multiple issues.
- `keen issue create` — Creates an issue or, where the option to create subtasks is enabled in Jira, a subtask.
- `keen issue create-bulk` — Creates upto **50** issues and, where the option to create subtasks is enabled in Jira, subtasks.
- `keen issue delete` — Deletes an issue. An issue cannot be deleted if it has one or more subtasks.
- `keen issue edit` — Edits an issue. Issue properties may be updated as part of the edit.
- `keen issue get` — Returns the details for an issue.
- `keen issue get-create-meta` — Returns details of projects, issue types within projects, and, when requested
- `keen issue get-create-meta-type-id` — Returns a page of field metadata for a specified project and issuetype id.
- `keen issue get-create-meta-types` — Returns a page of issue type metadata for a specified project.
- `keen issue get-is-watching-bulk` — Returns, for the user, details of the watched status of issues from a list.
- `keen issue get-limit-report` — Returns all issues breaching and approaching per-issue limits.
- `keen issue get-picker-resource` — Returns lists of issues matching a query string.
- `keen issue unarchive` — Enables admins to unarchive up to 1000 issues in a single request using issue ID/key, returning details of the issue(s)

**issue-link** — This resource represents links between issues. Use it to get, create, and delete links between issues.

To use it, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

- `keen issue-link delete` — Deletes an issue link. This operation can be accessed anonymously.
- `keen issue-link get` — Returns an issue link. This operation can be accessed anonymously.
- `keen issue-link link-issues` — Creates a link between two issues.

**issue-link-type** — This resource represents [issue link](#api-group-Issue-links) types. Use it to get, create, update, and delete link issue types as well as get lists of all link issue types.

To use it, the site must have [issue linking](https://confluence.atlassian.com/x/yoXKM) enabled.

- `keen issue-link-type create` — Creates an issue link type. Use this operation to create descriptions of the reasons why issues are linked.
- `keen issue-link-type delete` — Deletes an issue link type. To use this operation, the site must have [issue linking](https://confluence.atlassian.
- `keen issue-link-type get` — Returns a list of all issue link types. To use this operation, the site must have [issue linking](https://confluence.
- `keen issue-link-type get-issuelinktype` — Returns an issue link type. To use this operation, the site must have [issue linking](https://confluence.atlassian.
- `keen issue-link-type update` — Updates an issue link type. To use this operation, the site must have [issue linking](https://confluence.atlassian.

**issues** — This resource represents Jira issues. Use it to:

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

- `keen issues` — Enables admins to retrieve details of all archived issues.

**issuesecurityschemes** — Manage issuesecurityschemes

- `keen issuesecurityschemes associate-schemes-to-projects` — Associates an issue security scheme with a project and remaps security levels of issues to the new levels, if provided.
- `keen issuesecurityschemes create-issue-security-scheme` — Creates a security scheme with security scheme levels and levels' members.
- `keen issuesecurityschemes delete-security-scheme` — Deletes an issue security scheme.
- `keen issuesecurityschemes get-issue-security-scheme` — Returns an issue security scheme along with its security levels.
- `keen issuesecurityschemes get-issue-security-schemes` — Returns all [issue security schemes](https://confluence.atlassian.com/x/J4lKLg).
- `keen issuesecurityschemes get-security-level-members` — Returns a [paginated](#pagination) list of issue security level members.
- `keen issuesecurityschemes get-security-levels` — Returns a [paginated](#pagination) list of issue security levels.
- `keen issuesecurityschemes search-projects-using-security-schemes` — Returns a [paginated](#pagination) mapping of projects that are using security schemes.
- `keen issuesecurityschemes search-security-schemes` — Returns a [paginated](#pagination) list of issue security schemes.
- `keen issuesecurityschemes set-default-levels` — Sets default issue security levels for schemes.
- `keen issuesecurityschemes update-issue-security-scheme` — Updates the issue security scheme.

**issuetype** — Manage issuetype

- `keen issuetype create-issue-type` — Creates an issue type.
- `keen issuetype delete-issue-type` — Deletes the issue type.
- `keen issuetype get-issue-all-types` — Returns all issue types. This operation can be accessed anonymously.
- `keen issuetype get-issue-type` — Returns an issue type. This operation can be accessed anonymously.
- `keen issuetype get-issue-types-for-project` — Returns issue types for a project. This operation can be accessed anonymously.
- `keen issuetype update-issue-type` — Updates the issue type.

**issuetypescheme** — Manage issuetypescheme

- `keen issuetypescheme assign-issue-type-scheme-to-project` — Assigns an issue type scheme to a project.
- `keen issuetypescheme create-issue-type-scheme` — Creates an issue type scheme.
- `keen issuetypescheme delete-issue-type-scheme` — Deletes an issue type scheme. Only issue type schemes used in classic projects can be deleted.
- `keen issuetypescheme get-all-issue-type-schemes` — Returns a [paginated](#pagination) list of issue type schemes.
- `keen issuetypescheme get-issue-type-scheme-for-projects` — Returns a [paginated](#pagination) list of issue type schemes and, for each issue type scheme
- `keen issuetypescheme get-issue-type-schemes-mapping` — Returns a [paginated](#pagination) list of issue type scheme items.
- `keen issuetypescheme update-issue-type-scheme` — Updates an issue type scheme.

**issuetypescreenscheme** — Manage issuetypescreenscheme

- `keen issuetypescreenscheme assign-issue-type-screen-scheme-to-project` — Assigns an issue type screen scheme to a project. Issue type screen schemes can only be assigned to classic projects.
- `keen issuetypescreenscheme create-issue-type-screen-scheme` — Creates an issue type screen scheme.
- `keen issuetypescreenscheme delete-issue-type-screen-scheme` — Deletes an issue type screen scheme.
- `keen issuetypescreenscheme get-issue-type-screen-scheme-mappings` — Returns a [paginated](#pagination) list of issue type screen scheme items.
- `keen issuetypescreenscheme get-issue-type-screen-scheme-project-associations` — Returns a [paginated](#pagination) list of issue type screen schemes and, for each issue type screen scheme
- `keen issuetypescreenscheme get-issue-type-screen-schemes` — Returns a [paginated](#pagination) list of issue type screen schemes.
- `keen issuetypescreenscheme update-issue-type-screen-scheme` — Updates an issue type screen scheme.

**jira-cloud-platform-search** — Manage jira cloud platform search

- `keen jira-cloud-platform-search and-reconsile-issues-using-jql` — Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ).
- `keen jira-cloud-platform-search and-reconsile-issues-using-jql-post` — Searches for issues using [JQL](https://confluence.atlassian.com/x/egORLQ).
- `keen jira-cloud-platform-search count-issues` — Provide an estimated count of the issues that match the [JQL](https://confluence.atlassian.com/x/egORLQ).
- `keen jira-cloud-platform-search for-issues-using-jql` — Endpoint is currently being removed. [More details](https://developer.atlassian.
- `keen jira-cloud-platform-search for-issues-using-jql-post` — Endpoint is currently being removed. [More details](https://developer.atlassian.

**jira-cloud-platform-version** — Manage jira cloud platform version

- `keen jira-cloud-platform-version create` — Creates a project version. This operation can be accessed anonymously.
- `keen jira-cloud-platform-version delete` — Deletes a project version.
- `keen jira-cloud-platform-version get` — Returns a project version. This operation can be accessed anonymously.
- `keen jira-cloud-platform-version update` — Updates a project version. This operation can be accessed anonymously.

**jira-cloud-platform-workflow** — Manage jira cloud platform workflow

- `keen jira-cloud-platform-workflow create-transition-property` — This will be removed on [June 1, 2026](https://developer.atlassian.
- `keen jira-cloud-platform-workflow delete-inactive` — Deletes a workflow. The workflow cannot be deleted if it is: * an active workflow. * a system workflow.
- `keen jira-cloud-platform-workflow delete-transition-property` — This will be removed on [June 1, 2026](https://developer.atlassian.
- `keen jira-cloud-platform-workflow delete-transition-rule-configurations` — Deletes workflow transition rules from one or more workflows.
- `keen jira-cloud-platform-workflow get-paginated` — This will be removed on [June 1, 2026](https://developer.atlassian.
- `keen jira-cloud-platform-workflow get-transition-properties` — This will be removed on [June 1, 2026](https://developer.atlassian.
- `keen jira-cloud-platform-workflow get-transition-rule-configurations` — Returns a [paginated](#pagination) list of workflows with transition rules.
- `keen jira-cloud-platform-workflow list-history` — Returns a list of workflow history entries for a specified workflow id.
- `keen jira-cloud-platform-workflow read-from-history` — Returns a workflow and related statuses for a specified workflow id and version number.
- `keen jira-cloud-platform-workflow update-transition-property` — This will be removed on [June 1, 2026](https://developer.atlassian.
- `keen jira-cloud-platform-workflow update-transition-rule-configurations` — Updates configuration of workflow transition rules.

**jql** — This resource represents JQL search auto-complete details. Use it to obtain JQL search auto-complete data and suggestions for use in programmatic construction of queries or custom query builders. It also provides operations to:

 *  convert one or more JQL queries with user identifiers (username or user key) to equivalent JQL queries with account IDs.
 *  convert readable details in one or more JQL queries to IDs where a user doesn't have permission to view the entity whose details are readable.

- `keen jql get-auto-complete` — Returns reference data for JQL searches.
- `keen jql get-auto-complete-post` — Returns reference data for JQL searches.
- `keen jql get-field-auto-complete-for-query-string` — Returns the JQL search auto complete suggestions for a field.
- `keen jql get-precomputations` — Returns the list of a function's precomputations along with information about when they were created, updated
- `keen jql get-precomputations-by-id` — Returns function precomputations by IDs, along with information about when they were created, updated, and last used.
- `keen jql match-issues` — Checks whether one or more issues would be returned by one or more JQL queries.
- `keen jql migrate-queries` — Converts one or more JQL queries with user identifiers (username or user key)
- `keen jql parse-queries` — Parses and validates JQL queries. Validation is performed in context of the current user.
- `keen jql sanitise-queries` — Sanitizes one or more JQL queries by converting readable details into IDs where a user doesn't have permission to view
- `keen jql update-precomputations` — Update the precomputation value of a function created by a Forge/Connect app.

**label** — This resource represents available labels. Use it to get available labels for the global label field.

- `keen label` — Returns a [paginated](#pagination) list of labels.

**license** — Manage license

- `keen license get-approximate-application-count` — Returns the total approximate number of user accounts for a single Jira license.
- `keen license get-approximate-count` — Returns the approximate number of user accounts across all Jira licenses.

**license-types** — Information about license types

- `keen license-types get` — Get a specific license type.
- `keen license-types list` — Get a list of license types.

**mypermissions** — Manage mypermissions

- `keen mypermissions` — Returns a list of permissions indicating which permissions the user has.

**mypreferences** — Manage mypreferences

- `keen mypreferences get-locale` — Returns the locale for the user.
- `keen mypreferences get-preference` — Returns the value of a preference of the current user. Note that these keys are deprecated: * *jira.user.
- `keen mypreferences remove-preference` — Deletes a preference of the user, which restores the default value of system defined settings.
- `keen mypreferences set-locale` — Deprecated, use [ Update a user profile](https://developer.atlassian.
- `keen mypreferences set-preference` — Creates a preference for the user or updates a preference's value by sending a plain text string. For example, `false`.

**myself** — This resource represents information about the current user, such as basic details, group membership, application roles, preferences, and locale. Use it to get, create, update, and delete (restore default) values of the user's preferences and locale.

- `keen myself` — Returns details for the current user. **[Permissions](#permissions) required:** Permission to access Jira.

**notificationscheme** — Manage notificationscheme

- `keen notificationscheme create-notification-scheme` — Creates a notification scheme with notifications. You can create up to 1000 notifications per request.
- `keen notificationscheme delete-notification-scheme` — Deletes a notification scheme.
- `keen notificationscheme get-notification-scheme` — Returns a [notification scheme](https://confluence.atlassian.
- `keen notificationscheme get-notification-scheme-to-project-mappings` — Returns a [paginated](#pagination) mapping of project that have notification scheme assigned.
- `keen notificationscheme get-notification-schemes` — Returns a [paginated](#pagination) list of [notification schemes](https://confluence.atlassian.
- `keen notificationscheme update-notification-scheme` — Updates a notification scheme.

**operations** — APIs related to integrating Incident and Post-Incident Review (PIR) data with Jira Software.

These APIs are available to Atlassian Connect apps. To use these APIs you must have the `jiraOperationsInfoProvider` module in your Connect app's descriptor. See https://developer.atlassian.com/cloud/jira/software/modules/operations-information/.

- `keen operations delete-entity-by-property` — Bulk delete all Entties that match the given request.
- `keen operations delete-incident-by-id` — Delete the Incident data currently stored for the given ID. Deletion is performed asynchronously.
- `keen operations delete-review-by-id` — Delete the Review data currently stored for the given ID. Deletion is performed asynchronously.
- `keen operations delete-workspaces` — Bulk delete all Operations Workspaces that match the given request.
- `keen operations get-incident-by-id` — Retrieve the currently stored Incident data for the given ID.
- `keen operations get-review-by-id` — Retrieve the currently stored Review data for the given ID.
- `keen operations get-workspaces` — Retrieve the either all Operations Workspace IDs associated with the Jira site or a specific Operations Workspace ID
- `keen operations submit-entity` — Update / insert Incident or Review data.
- `keen operations submit-workspaces` — Insert Operations Workspace IDs to establish a relationship between them and the Jira site the app is installed in.

**orgs** — Orgs APIs

- `keen orgs get` — Returns a list of your organizations (based on your API key).
- `keen orgs get-by-id` — Returns information about a single organization by ID

**permissions** — This resource represents permissions. Use it to obtain details of all permissions and determine whether the user has certain permissions.

- `keen permissions get-all` — Returns all permissions, including: * global permissions. * project permissions. * global permissions added by plugins.
- `keen permissions get-bulk` — Returns: * for a list of global permissions, the global permissions granted to a user.
- `keen permissions get-permitted-projects` — Returns all the projects where the user is granted a list of project permissions.

**permissionscheme** — Manage permissionscheme

- `keen permissionscheme create-permission-scheme` — Creates a new permission scheme. You can create a permission scheme with or without defining a set of permission grants.
- `keen permissionscheme delete-permission-scheme` — Deletes a permission scheme.
- `keen permissionscheme get-all-permission-schemes` — Returns all permission schemes.
- `keen permissionscheme get-permission-scheme` — Returns a permission scheme. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen permissionscheme update-permission-scheme` — Updates a permission scheme.

**plans** — This resource represents plans. Use it to get, create, duplicate, update, trash and archive plans.

- `keen plans add-atlassian-team` — Adds an existing Atlassian team to a plan and configures their plannning settings.
- `keen plans archive` — Archives a plan. **[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.
- `keen plans create` — Creates a plan. **[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.
- `keen plans create-only-team` — Creates a plan-only team and configures their planning settings.
- `keen plans delete-only-team` — Deletes a plan-only team and their planning settings.
- `keen plans duplicate` — Duplicates a plan. **[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.
- `keen plans get` — Returns a [paginated](#pagination) list of plans.
- `keen plans get-atlassian-team` — Returns planning settings for an Atlassian team in a plan.
- `keen plans get-only-team` — Returns planning settings for a plan-only team.
- `keen plans get-plan` — Returns a plan. **[Permissions](#permissions) required:** *Administer Jira* [global permission](https://confluence.
- `keen plans get-teams` — Returns a [paginated](#pagination) list of plan-only and Atlassian teams in a plan.
- `keen plans remove-atlassian-team` — Removes an Atlassian team from a plan and deletes their planning settings.
- `keen plans trash` — Moves a plan to trash.
- `keen plans update` — Updates any of the following details of a plan using [JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902).
- `keen plans update-atlassian-team` — Updates any of the following planning settings of an Atlassian team in a plan using [JSON Patch](https://datatracker.
- `keen plans update-only-team` — Updates any of the following planning settings of a plan-only team using [JSON Patch](https://datatracker.ietf.

**policies** — Policies APIs

- `keen policies create-policy` — Create a policy for an org
- `keen policies delete-policy` — Delete a policy for an org
- `keen policies get` — Returns information about org policies
- `keen policies get-policy-by-id` — Returns information about a single policy by ID
- `keen policies update-policy` — Update a policy for an org

**priority** — Manage priority

- `keen priority create` — Creates an issue priority. Deprecation applies to iconUrl param in request body which will be sunset on 16th Mar 2025.
- `keen priority delete` — Deletes an issue priority. This operation is [asynchronous](#async).
- `keen priority get` — Returns an issue priority. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen priority get-priorities` — Returns the list of all issue priorities. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen priority move-priorities` — Changes the order of issue priorities.
- `keen priority search-priorities` — Returns a [paginated](#pagination) list of priorities.
- `keen priority set-default` — Sets default issue priority.
- `keen priority update` — Updates an issue priority. At least one request body parameter must be defined.

**priorityscheme** — Manage priorityscheme

- `keen priorityscheme create-priority-scheme` — Creates a new priority scheme.
- `keen priorityscheme delete-priority-scheme` — Deletes a priority scheme. This operation is only available for priority schemes without any associated projects.
- `keen priorityscheme get-available-priorities-by-priority-scheme` — Returns a [paginated](#pagination) list of priorities available for adding to a priority scheme.
- `keen priorityscheme get-priority-schemes` — Returns a [paginated](#pagination) list of priority schemes.
- `keen priorityscheme suggested-priorities-for-mappings` — Returns a [paginated](#pagination) list of priorities that would require mapping
- `keen priorityscheme update-priority-scheme` — Updates a priority scheme.

**products** — Information about products

- `keen products get` — Get a specific product.
- `keen products get-key` — Get a list of versions for the specified product.
- `keen products get-key-2` — Get the latest version of the specified product.
- `keen products get-key-3` — Get details about a specific version, matching the specified build number, of the specified product.
- `keen products get-key-4` — Get details about a specific version, matching the specified name, of the specified product.
- `keen products list` — Get a list of products matching the specified parameters.

**project** — This resource represents projects. Use it to get, create, update, and delete projects. Also get statuses available to a project, a project's notification schemes, and update a project's type.

- `keen project create` — Creates a project based on a project type template, as shown in the following table
- `keen project delete` — Deletes a project. You can't delete a project if it's archived.
- `keen project get` — Returns the [project details](https://confluence.atlassian.com/x/ahLpNw) for a project.
- `keen project get-accessible-type-by-key` — Returns a [project type](https://confluence.atlassian.com/x/Var1Nw) if it is accessible to the user.
- `keen project get-all` — Returns all projects visible to the user.
- `keen project get-all-accessible-types` — Returns all [project types](https://confluence.atlassian.com/x/Var1Nw) with a valid license.
- `keen project get-all-types` — Returns all [project types](https://confluence.atlassian.
- `keen project get-recent` — Returns a list of up to 20 projects recently viewed by the user that are still visible to the user.
- `keen project get-type-by-key` — Returns a [project type](https://confluence.atlassian.com/x/Var1Nw). This operation can be accessed anonymously.
- `keen project search` — Returns a [paginated](#pagination) list of projects visible to the user. This operation can be accessed anonymously.
- `keen project update` — Updates the [project details](https://confluence.atlassian.com/x/ahLpNw) of a project.

**project-category** — Manage project category

- `keen project-category create` — Creates a project category.
- `keen project-category get-all-project-categories` — Returns all project categories. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen project-category get-by-id` — Returns a project category. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen project-category remove` — Deletes a project category.
- `keen project-category update` — Updates a project category.

**project-template** — This resource represents project templates. Use it to create a new project from a custom template.

- `keen project-template create-project-with-custom-template` — Creates a project based on a custom template provided in the request.
- `keen project-template edit-template` — Edit custom template This API endpoint allows you to edit an existing customised template.
- `keen project-template live-template` — Get custom template This API endpoint allows you to get a live custom project template details by either templateKey or
- `keen project-template remove-template` — Remove custom template This API endpoint allows you to remove a specified customised template ***Note
- `keen project-template save-template` — Save custom template This API endpoint allows you to save a customised template ***Note

**projects** — This resource represents projects. Use it to get, create, update, and delete projects. Also get statuses available to a project, a project's notification schemes, and update a project's type.

- `keen projects` — Returns a [paginated](#pagination) list of fields for the requested projects and work types.

**projectvalidate** — Manage projectvalidate

- `keen projectvalidate get-valid-project-key` — Validates a project key and, if the key is invalid or in use, generates a valid random string for the project key.
- `keen projectvalidate get-valid-project-name` — Checks that a project name isn't in use. If the name isn't in use, the passed string is returned.
- `keen projectvalidate validate-project-key` — Validates a project key by confirming the key is a valid string and not in use.

**redact** — Manage redact

- `keen redact get-redaction-status` — Retrieves the current status of a redaction job ID.
- `keen redact redact` — Submit a job to redact issue field data.

**remotelinks** — Manage remotelinks

- `keen remotelinks delete-remote-link-by-id` — Delete the Remote Link data currently stored for the given ID. Deletion is performed asynchronously.
- `keen remotelinks delete-remote-links-by-property` — Bulk delete all Remote Links data that match the given request.
- `keen remotelinks get-remote-link-by-id` — Retrieve the currently stored Remote Link data for the given ID.
- `keen remotelinks submit-remote-links` — Update / insert Remote Link data.

**resolution** — Manage resolution

- `keen resolution create` — Creates an issue resolution.
- `keen resolution delete` — Deletes an issue resolution. This operation is [asynchronous](#async).
- `keen resolution get` — Returns a list of all issue resolution values. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen resolution get-id` — Returns an issue resolution value. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen resolution move` — Changes the order of issue resolutions.
- `keen resolution search` — Returns a [paginated](#pagination) list of resolutions.
- `keen resolution set-default` — Sets default issue resolution.
- `keen resolution update` — Updates an issue resolution.

**role** — Manage role

- `keen role create-project` — Creates a new project role with no [default actors](#api-rest-api-2-resolution-get).
- `keen role delete-project` — Deletes a project role. You must specify a replacement project role if you wish to delete a project role that is in use.
- `keen role fully-update-project` — Updates the project role's name and description. You must include both a name and a description in the request.
- `keen role get-all-project` — Gets a list of all project roles, complete with project role details and default actors.
- `keen role get-project-by-id` — Gets the project role details and the default actors associated with the role.
- `keen role partial-update-project` — Updates either the project role's name or its description.

**screens** — This resource represents the screens used to record issue details. Use it to:

 *  get details of all screens.
 *  get details of all the fields available for use on screens.
 *  create screens.
 *  delete screens.
 *  update screens.
 *  add a field to the default screen.

- `keen screens add-field-to-default` — Adds a field to the default tab of the default screen.
- `keen screens create` — Creates a screen with a default field tab.
- `keen screens delete` — Deletes a screen. A screen cannot be deleted if it is used in a screen scheme, workflow, or workflow draft.
- `keen screens get` — Returns a [paginated](#pagination) list of all screens or those specified by one or more screen IDs.
- `keen screens get-bulk-tabs` — Returns the list of tabs for a bulk of screens.
- `keen screens update` — Updates a screen. Only screens used in classic projects can be updated.

**screenscheme** — Manage screenscheme

- `keen screenscheme create-screen-scheme` — Creates a screen scheme.
- `keen screenscheme delete-screen-scheme` — Deletes a screen scheme. A screen scheme cannot be deleted if it is used in an issue type screen scheme.
- `keen screenscheme get-screen-schemes` — Returns a [paginated](#pagination) list of screen schemes. Only screen schemes used in classic projects are returned.
- `keen screenscheme update-screen-scheme` — Updates a screen scheme. Only screen schemes used in classic projects can be updated.

**security** — Manage security

- `keen security delete-linked-workspaces` — Bulk delete all linked Security Workspaces that match the given request.
- `keen security delete-vulnerabilities-by-property` — Bulk delete all Vulnerabilities that match the given request.
- `keen security delete-vulnerability-by-id` — Delete the Vulnerability data currently stored for the given ID. Deletion is performed asynchronously.
- `keen security get-linked-workspace-by-id` — Retrieve a specific Security Workspace linked to the Jira site for the given workspace ID.
- `keen security get-linked-workspaces` — Retrieve all Security Workspaces linked with the Jira site.
- `keen security get-vulnerability-by-id` — Retrieve the currently stored Vulnerability data for the given ID.
- `keen security submit-vulnerabilities` — Update / Insert Vulnerability data.
- `keen security submit-workspaces` — Insert Security Workspace IDs to establish a relationship between them and the Jira site the app is installed on.

**securitylevel** — Manage securitylevel

- `keen securitylevel <id>` — Returns details of an issue security level.

**server-info** — This resource provides information about the Jira instance.

- `keen server-info` — Returns information about the Jira instance. This operation can be accessed anonymously.

**servicedeskapi** — Manage servicedeskapi

- `keen servicedeskapi add-customers` — Adds one or more customers to a service desk.
- `keen servicedeskapi add-organization` — This method adds an organization to a service desk.
- `keen servicedeskapi add-request-participants` — This method adds participants to a customer request.
- `keen servicedeskapi add-users-to-organization` — This method adds users to an organization.
- `keen servicedeskapi answer-approval` — This method enables a user to **Approve** or **Decline** an approval on a customer request.
- `keen servicedeskapi attach-temporary-file` — This method adds one or more temporary attachments to a service desk
- `keen servicedeskapi check-request-type-permissions` — Returns: * a list of request type IDs where the given user has permission to administer.
- `keen servicedeskapi create-comment-with-attachment` — This method creates a comment on a customer request using one or more attachment files (uploaded using
- `keen servicedeskapi create-customer` — This method adds a customer to the Jira Service Management instance by passing a JSON file including an email address
- `keen servicedeskapi create-customer-request` — This method creates a customer request in a service desk.
- `keen servicedeskapi create-organization` — This method creates an organization by passing the name of the organization.
- `keen servicedeskapi create-request-comment` — This method creates a public or private (internal) comment on a customer request
- `keen servicedeskapi create-request-type` — This method enables a customer request type to be added to a service desk based on an issue type.
- `keen servicedeskapi delete-feedback` — This method deletes the feedback of request using it's `requestKey` or `requestId` **[Permissions](#permissions)
- `keen servicedeskapi delete-organization` — This method deletes an organization. Note that the organization is deleted regardless of other associations it may have.
- `keen servicedeskapi delete-property` — Removes an organization property.
- `keen servicedeskapi delete-property-servicedesk` — Removes a property from a request type.
- `keen servicedeskapi delete-request-type` — This method deletes a customer request type from a service desk, and removes it from all customer requests.
- `keen servicedeskapi get-all-request-types` — This method returns all customer request types used in the Jira Service Management instance
- `keen servicedeskapi get-approval-by-id` — This method returns an approval. Use this method to determine the status of an approval and the list of approvers.
- `keen servicedeskapi get-approvals` — This method returns all approvals on a customer request.
- `keen servicedeskapi get-articles` — Returns articles which match the given query string across all service desks.
- `keen servicedeskapi get-articles-servicedesk` — Returns articles which match the given query and belong to the knowledge base linked to the service desk.
- `keen servicedeskapi get-assets-workspaces` — Returns a list of Assets workspace IDs.
- `keen servicedeskapi get-attachment-content` — Returns the contents of an attachment.
- `keen servicedeskapi get-attachment-thumbnail` — Returns the thumbnail of an attachment.
- `keen servicedeskapi get-attachments-for-request` — This method returns all the attachments for a customer requests.
- `keen servicedeskapi get-comment-attachments` — This method returns the attachments referenced in a comment.
- `keen servicedeskapi get-customer-request-by-id-or-key` — This method returns a customer request.
- `keen servicedeskapi get-customer-request-status` — This method returns a list of all the statuses a customer Request has achieved.
- `keen servicedeskapi get-customer-requests` — This method returns all customer requests for the user executing the query.
- `keen servicedeskapi get-customer-transitions` — This method returns a list of transitions
- `keen servicedeskapi get-customers` — This method returns a list of the customers on a service desk.
- `keen servicedeskapi get-feedback` — This method retrieves a feedback of a request using it's `requestKey` or `requestId` **[Permissions](#permissions)
- `keen servicedeskapi get-info` — This method retrieves information about the Jira Service Management instance such as software version, builds
- `keen servicedeskapi get-insight-workspaces` — This endpoint is deprecated, please use /assets/workspace/.
- `keen servicedeskapi get-issues-in-queue` — This method returns the customer requests in a queue. Only fields that the queue is configured to show are returned.
- `keen servicedeskapi get-organization` — This method returns details of an organization.
- `keen servicedeskapi get-organizations` — This method returns a list of organizations in the Jira Service Management instance.
- `keen servicedeskapi get-organizations-servicedesk` — This method returns a list of all organizations associated with a service desk.
- `keen servicedeskapi get-properties-keys` — Returns the keys of all organization properties.
- `keen servicedeskapi get-properties-keys-servicedesk` — Returns the keys of all properties for a request type.
- `keen servicedeskapi get-property` — Returns the value of an organization property.
- `keen servicedeskapi get-property-servicedesk` — Returns the value of the property from a request type.
- `keen servicedeskapi get-queue` — This method returns a specific queues in a service desk.
- `keen servicedeskapi get-queues` — This method returns the queues in a service desk.
- `keen servicedeskapi get-request-comment-by-id` — This method returns details of a customer request's comment.
- `keen servicedeskapi get-request-comments` — This method returns all comments on a customer request.
- `keen servicedeskapi get-request-participants` — This method returns a list of all the participants on a customer request.
- `keen servicedeskapi get-request-type-by-id` — This method returns a customer request type from a service desk. This operation can be accessed anonymously.
- `keen servicedeskapi get-request-type-fields` — This method returns the fields for a service desk's customer request type.
- `keen servicedeskapi get-request-type-groups` — This method returns a service desk's customer request type groups.
- `keen servicedeskapi get-request-types` — This method returns all customer request types from a service desk.
- `keen servicedeskapi get-service-desk-by-id` — This method returns a service desk.
- `keen servicedeskapi get-service-desks` — This method returns all the service desks in the Jira Service Management instance that the user has permission to
- `keen servicedeskapi get-sla-information` — This method returns all the SLA records on a customer request. A customer request can have zero or more SLAs.
- `keen servicedeskapi get-sla-information-by-id` — This method returns the details for an SLA on a customer request.
- `keen servicedeskapi get-subscription-status` — This method returns the notification subscription status of the user making the request.
- `keen servicedeskapi get-users-in-organization` — This method returns all the users associated with an organization.
- `keen servicedeskapi invite-customer` — This method invites a customer to a specified service desk by sending them an email invitation
- `keen servicedeskapi perform-customer-transition` — This method performs a customer transition for a given request and transition.
- `keen servicedeskapi post-feedback` — This method adds a feedback on an request using it's `requestKey` or `requestId` **[Permissions](#permissions)
- `keen servicedeskapi remove-customers` — This method removes one or more customers from a service desk. The service desk must have closed access.
- `keen servicedeskapi remove-organization` — This method removes an organization from a service desk.
- `keen servicedeskapi remove-request-participants` — This method removes participants from a customer request.
- `keen servicedeskapi remove-users-from-organization` — This method removes users from an organization.
- `keen servicedeskapi revoke-portal-only-access-for-user` — This method revokes portal-only access for a particular user
- `keen servicedeskapi set-property` — Sets the value of an organization property. Use this resource to store custom data against an organization.
- `keen servicedeskapi set-property-servicedesk` — Sets the value of a request type property. Use this resource to store custom data against a request type.
- `keen servicedeskapi subscribe` — This method subscribes the user to receiving notifications from a customer request.
- `keen servicedeskapi unsubscribe` — This method unsubscribes the user from notifications from a customer request.

**settings** — Manage settings

- `keen settings get-issue-navigator-default-columns` — Returns the default issue navigator columns.
- `keen settings set-issue-navigator-default-columns` — Sets the default issue navigator columns.

**software** — Manage software

- `keen software get-approximate-issue-count-for-backlog` — Returns the approximate count of all issues from the board's backlog, for the given board ID.
- `keen software get-approximate-issue-count-for-board` — Returns the approximate count of all issues from a board, for a given board ID.
- `keen software get-board-issues-for-epic-jsis` — Returns all issues that belong to an epic on the board, for the given epic ID and the board ID.
- `keen software get-board-issues-for-sprint-jsis` — Get all issues you have access to that belong to the sprint from the board.
- `keen software get-issues-for-backlog-jsis` — Returns all issues from the board's backlog, for the given board ID.
- `keen software get-issues-for-board-jsis` — Returns all issues from a board, for a given board ID.
- `keen software get-issues-for-epic-jsis` — Returns all issues that belong to the epic, for the given epic ID.
- `keen software get-issues-for-sprint-jsis` — Returns all issues in a sprint, for a given sprint ID.
- `keen software get-issues-without-epic-for-board-jsis` — Returns all issues that do not belong to any epic on a board, for a given board ID.
- `keen software get-issues-without-epic-jsis` — Returns all issues that do not belong to any epic.

**status** — This resource represents statuses. Use it to search, get, create, delete, and change statuses.

- `keen status get` — Returns a status. The status must be associated with an active workflow to be returned.
- `keen status get-statuses` — Returns a list of all statuses associated with active workflows. This operation can be accessed anonymously.

**statuscategory** — Manage statuscategory

- `keen statuscategory get-status-categories` — Returns a list of all status categories. **[Permissions](#permissions) required:** Permission to access Jira.
- `keen statuscategory get-status-category` — Returns a status category.

**statuses** — Manage statuses

- `keen statuses create` — Creates statuses for a global or project scope.
- `keen statuses delete-by-id` — Deletes statuses by ID. **[Permissions](#permissions) required:** * *Administer projects* [project permission.
- `keen statuses get-by-id` — Returns a list of the statuses specified by one or more status IDs.
- `keen statuses get-by-name` — Returns a list of the statuses specified by one or more status names.
- `keen statuses search` — Returns a [paginated](https://developer.atlassian.
- `keen statuses update` — Updates statuses by ID. **[Permissions](#permissions) required:** * *Administer projects* [project permission.

**task** — This resource represents a [long-running asynchronous tasks](#async-operations). Use it to obtain details about the progress of a long-running task or cancel a long-running task.

- `keen task <taskId>` — Returns the status of a [long-running asynchronous task](#async).

**ui-modifications** — Manage ui modifications

- `keen ui-modifications create` — Creates a UI modification. UI modification can only be created by Forge apps.
- `keen ui-modifications delete` — Deletes a UI modification. All the contexts that belong to the UI modification are deleted too.
- `keen ui-modifications get` — Gets UI modifications. UI modifications can only be retrieved by Forge apps.
- `keen ui-modifications update` — Updates a UI modification. UI modification can only be updated by Forge apps.

**universal-avatar** — Manage universal avatar

- `keen universal-avatar delete-avatar` — Deletes an avatar from a project, issue type or priority.
- `keen universal-avatar get-avatar-image-by-id` — Returns a project, issue type or priority avatar image by ID. This operation can be accessed anonymously.
- `keen universal-avatar get-avatar-image-by-owner` — Returns the avatar image for a project, issue type or priority. This operation can be accessed anonymously.
- `keen universal-avatar get-avatar-image-by-type` — Returns the default project, issue type or priority avatar image. This operation can be accessed anonymously.
- `keen universal-avatar get-avatars` — Returns the system and custom avatars for a project, issue type or priority. This operation can be accessed anonymously.
- `keen universal-avatar store-avatar` — Loads a custom avatar for a project, issue type or priority.

**user** — This resource represent users. Use it to:

 *  get, get a list of, create, and delete users.
 *  get, set, and reset a user's default issue table columns.
 *  get a list of the groups the user belongs to.
 *  get a list of user account IDs for a list of usernames or user keys.

- `keen user bulk-get` — Returns a [paginated](#pagination) list of the users specified by one or more account IDs.
- `keen user bulk-get-migration` — Returns the account IDs for the users specified in the `key` or `username` parameters.
- `keen user create` — Creates a user. This resource is retained for legacy compatibility.
- `keen user delete-property` — Deletes a property from a user. Note: This operation does not access the [user properties](https://confluence.atlassian.
- `keen user find` — Returns a list of active users that match the search string and property.
- `keen user find-assignable` — Returns a list of users that can be assigned to an issue.
- `keen user find-bulk-assignable` — Returns a list of users who can be assigned issues in one or more projects.
- `keen user find-by-query` — Finds users with a structured query and returns a [paginated](#pagination) list of user details.
- `keen user find-for-picker` — Returns a list of users whose attributes match the query term.
- `keen user find-keys-by-query` — Finds users with a structured query and returns a [paginated](#pagination) list of user keys.
- `keen user find-with-all-permissions` — Returns a list of users who fulfill these criteria: * their user attributes match a search string.
- `keen user find-with-browse-permission` — Returns a list of users who fulfill these criteria: * their user attributes match a search string.
- `keen user get` — Returns a user. Privacy controls are applied to the response based on the user's preferences.
- `keen user get-default-columns` — Returns the default [issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user.
- `keen user get-email` — Returns a user's email address regardless of the user's profile visibility settings.
- `keen user get-email-bulk` — Returns a user's email address regardless of the user's profile visibility settings.
- `keen user get-groups` — Returns the groups to which a user belongs.
- `keen user get-property` — Returns the value of a user's property.
- `keen user get-property-keys` — Returns the keys of all properties for a user.
- `keen user remove` — Deletes a user. If the operation completes successfully then the user is removed from Jira's user base.
- `keen user reset-columns` — Resets the default [ issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user to the system default.
- `keen user set-columns` — Sets the default [ issue table columns](https://confluence.atlassian.com/x/XYdKLg) for the user.
- `keen user set-property` — Sets the value of a user's property. Use this resource to store custom data against a user.

**users** — This resource represent users. Use it to:

 *  get, get a list of, create, and delete users.
 *  get, set, and reset a user's default issue table columns.
 *  get a list of the groups the user belongs to.
 *  get a list of user account IDs for a list of usernames or user keys.

- `keen users get-all` — Returns a list of all users, including active users
- `keen users get-all-default` — Returns a list of all users, including active users

**vendors** — Information about vendors


**webhook** — This resource represents webhooks. Webhooks are calls sent to a URL when an event occurs in Jira for issues specified by a JQL query. Only Connect and OAuth 2.0 apps can register and manage webhooks. For more information, see [Webhooks](https://developer.atlassian.com/cloud/jira/platform/webhooks/#registering-a-webhook-via-the-jira-rest-api-for-connect-apps).

- `keen webhook delete-by-id` — Removes webhooks by ID. Only webhooks registered by the calling app are removed.
- `keen webhook get-dynamic-for-app` — Returns a [paginated](#pagination) list of the webhooks registered by the calling app.
- `keen webhook get-failed` — Returns webhooks that have recently failed to be delivered to the requesting app after the maximum number of retries.
- `keen webhook refresh` — Extends the life of webhook. Webhooks registered through the REST API expire after 30 days.
- `keen webhook register-dynamic` — Registers webhooks.

**wiki** — Manage wiki

- `keen wiki add-content-watcher` — Adds a user as a watcher to a piece of content.
- `keen wiki add-custom-content-permissions` — Adds new custom content permission to space.
- `keen wiki add-label-watcher` — Adds a user as a watcher to a label.
- `keen wiki add-labels-to-content` — Adds labels to a piece of content. Does not modify the existing labels.
- `keen wiki add-labels-to-space` — Adds labels to a piece of content. Does not modify the existing labels.
- `keen wiki add-permission-to-space` — Adds new permission to space.
- `keen wiki add-restrictions` — Adds restrictions to a piece of content. Note, this does not change any existing restrictions on the content.
- `keen wiki add-user-to-group-by-group-id` — Adds a user as a member in a group represented by its groupId **[Permissions](https://confluence.atlassian.
- `keen wiki archive-pages` — Archives a list of pages. The pages to be archived are specified as a list of content IDs.
- `keen wiki async-convert-content-body-request` — Converts a content body from one format to another format asynchronously. Returns the asyncId for the asynchronous task.
- `keen wiki async-convert-content-body-response` — Returns the content body for the corresponding `asyncId` of a completed conversion task.
- `keen wiki check-content-permission` — Check if a user or a group can perform an operation to the specified content. The `operation` to check must be provided.
- `keen wiki copy-page` — Copies a single page and its associated properties, permissions, attachments, and custom contents.
- `keen wiki copy-page-hierarchy` — Copy page hierarchy allows the copying of an entire hierarchy of pages and their associated properties
- `keen wiki create-attachment` — Adds an attachment to a piece of content. This method only adds a new attachment.
- `keen wiki create-audit-record` — Creates a record in the audit log. **[Permissions](https://confluence.atlassian.
- `keen wiki create-content-template` — Creates a new content template. Note, blueprint templates cannot be created via the REST API.
- `keen wiki create-group` — Creates a new user group. **[Permissions](https://confluence.atlassian.
- `keen wiki create-or-update-attachments` — Adds an attachment to a piece of content.
- `keen wiki create-private-space` — Creates a new space that is only visible to the creator.
- `keen wiki create-space` — Creates a new space. Note, currently you cannot set space labels when creating a space.
- `keen wiki delete-content-version` — Delete a historical version.
- `keen wiki delete-label-from-space` — Remove label from a space
- `keen wiki delete-page-tree` — Moves a pagetree rooted at a page to the space's trash: - If the content's type is `page` and its status is `current`
- `keen wiki delete-restrictions` — Removes all restrictions (read and update) on a piece of content. **[Permissions](https://confluence.atlassian.
- `keen wiki delete-space` — Permanently deletes a space without sending it to the trash. Note, the space will be deleted in a long running task.
- `keen wiki export-audit-records` — Exports audit records as a CSV file or ZIP file. **[Permissions](https://confluence.atlassian.
- `keen wiki get-all-label-content` — Returns label information and a list of contents associated with the label. **[Permissions](https://confluence.
- `keen wiki get-anonymous-user` — Returns information about how anonymous users are represented, like the profile picture and display name.
- `keen wiki get-audit-records` — Returns all records in the audit log, optionally for a certain date range.
- `keen wiki get-audit-records-for-time-period` — Returns records from the audit log, for a time period back from the current date.
- `keen wiki get-available-content-states` — Gets content states that are available for the content to be set as. Will return all enabled Space Content States.
- `keen wiki get-blueprint-templates` — Returns all templates provided by blueprints.
- `keen wiki get-bulk-user-lookup` — Returns user details for the ids provided in the request. Currently this API returns a maximum of 100 results.
- `keen wiki get-content-descendants` — Returns a map of the descendants of a piece of content.
- `keen wiki get-content-state` — Gets the current content state of the draft or current version of content.
- `keen wiki get-content-state-settings` — Get object describing whether content states are allowed at all
- `keen wiki get-content-template` — Returns a content template.
- `keen wiki get-content-templates` — Returns all content templates.
- `keen wiki get-content-watch-status` — Returns whether a user is watching a piece of content.
- `keen wiki get-contents-with-state` — Returns all content that has the provided content state in a space.
- `keen wiki get-current-user` — Returns the currently logged-in user.
- `keen wiki get-custom-content-states` — Get custom content states that authenticated user has created. **[Permissions](https://confluence.atlassian.
- `keen wiki get-descendants-of-type` — Returns all descendants of a given type, for a piece of content.
- `keen wiki get-global-theme` — Returns the globally assigned theme. **[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: None
- `keen wiki get-group-by-group-id` — Returns a user group for a given group id. **[Permissions](https://confluence.atlassian.
- `keen wiki get-group-members-by-group-id` — Returns the users that are members of a group. Use updated Get group API **[Permissions](https://confluence.atlassian.
- `keen wiki get-group-memberships-for-user` — Returns the groups that a user is a member of. **[Permissions](https://confluence.atlassian.
- `keen wiki get-groups` — Returns all user groups. The returned groups are ordered alphabetically in ascending order by group name.
- `keen wiki get-labels-for-space` — Returns a list of labels associated with a space.
- `keen wiki get-look-and-feel-settings` — Returns the look and feel settings for the site or a single space.
- `keen wiki get-privacy-unsafe-user-email` — Returns a user's email address regardless of the user’s profile visibility settings.
- `keen wiki get-privacy-unsafe-user-email-bulk` — Returns a user's email address regardless of the user’s profile visibility settings.
- `keen wiki get-restrictions` — Returns the restrictions on a piece of content. **[Permissions](https://confluence.atlassian.
- `keen wiki get-restrictions-by-operation` — Returns restrictions on a piece of content by operation.
- `keen wiki get-retention-period` — Returns the retention period for records in the audit log.
- `keen wiki get-space-content-states` — Get content states that are suggested in the space. **[Permissions](https://confluence.atlassian.
- `keen wiki get-space-settings` — Returns the settings of a space. Currently only the `routeOverrideEnabled` setting can be returned.
- `keen wiki get-space-theme` — Returns the theme selected for a space, if one is set.
- `keen wiki get-task` — Returns information about an active long-running task (e.g.
- `keen wiki get-tasks` — Returns information about all active long-running tasks (e.g.
- `keen wiki get-theme` — Returns a theme. This includes information about the theme name, description, and icon.
- `keen wiki get-themes` — Returns all themes, not including the default theme. **[Permissions](https://confluence.atlassian.
- `keen wiki get-user` — Returns a user.
- `keen wiki get-user-properties` — Returns the properties for a user as list of property keys.
- `keen wiki get-viewers` — Get the total number of distinct viewers a piece of content has.
- `keen wiki get-views` — Get the total number of views a piece of content has.
- `keen wiki get-watchers-for-space` — Returns a list of watchers of a space
- `keen wiki get-watches-for-page` — Returns the watches for a page. A user that watches a page will receive receive notifications when the page is updated.
- `keen wiki get-watches-for-space` — Returns all space watches for the space that the content is in.
- `keen wiki is-watching-label` — Returns whether a user is watching a label.
- `keen wiki is-watching-space` — Returns whether a user is watching a space.
- `keen wiki publish-legacy-draft` — Publishes a legacy draft of a page created from a blueprint.
- `keen wiki publish-shared-draft` — Publishes a shared draft of a page created from a blueprint. By default, the following objects are expanded: `body.
- `keen wiki remove-content-state` — Removes the content state of the content specified and creates a new version (publishes the content without changing
- `keen wiki remove-content-watcher` — Removes a user as a watcher from a piece of content.
- `keen wiki remove-group-by-id` — Delete user group. **[Permissions](https://confluence.atlassian.com/x/_AozKw) required**: User must be a site admin.
- `keen wiki remove-label-from-content` — Removes a label from a piece of content. Labels can't be deleted from archived content.
- `keen wiki remove-label-from-content-using-query-parameter` — Removes a label from a piece of content. Labels can't be deleted from archived content.
- `keen wiki remove-label-watcher` — Removes a user as a watcher from a label.
- `keen wiki remove-member-from-group-by-group-id` — Remove user as a member from a group. **[Permissions](https://confluence.atlassian.
- `keen wiki remove-permission` — Removes a space permission.
- `keen wiki remove-space-watch` — Removes a user as a watcher from a space.
- `keen wiki remove-template` — Deletes a template.
- `keen wiki reset-look-and-feel-settings` — Resets the custom look and feel settings for the site or a single space.
- `keen wiki reset-space-theme` — Resets the space theme.
- `keen wiki restore-content-version` — Restores a historical version to be the latest version.
- `keen wiki search-by-cql` — Searches for content using the [Confluence Query Language (CQL)](https://developer.atlassian.
- `keen wiki search-content-by-cql` — Returns the list of content that matches a Confluence Query Language (CQL) query.
- `keen wiki search-groups` — Get search results of groups by partial query provided.
- `keen wiki search-user` — Searches for users using user-specific queries from the [Confluence Query Language (CQL)](https://developer.atlassian.
- `keen wiki set-content-state` — Sets the content state of the content specified and creates a new version (publishes the content without changing the
- `keen wiki set-retention-period` — Sets the retention period for records in the audit log. The retention period can be set to a maximum of 1 year.
- `keen wiki set-space-theme` — Sets the theme for a space.
- `keen wiki update-content-template` — Updates a content template. Note, blueprint templates cannot be updated via the REST API.
- `keen wiki update-look-and-feel` — Sets the look and feel settings to the default (global) settings, the custom settings
- `keen wiki update-look-and-feel-settings` — Updates the look and feel settings for the site or for a single space. If custom settings exist, they are updated.
- `keen wiki update-restrictions` — Updates restrictions for a piece of content.
- `keen wiki update-space` — Updates the name, description, or homepage of a space.
- `keen wiki update-space-settings` — Updates the settings for a space. Currently only the `routeOverrideEnabled` setting can be updated.

**workflows** — This resource represents workflows. Use it to:

 *  Get workflows
 *  Create workflows
 *  Update workflows
 *  Delete inactive workflows
 *  Get workflow capabilities

- `keen workflows capabilities` — Get the list of workflow capabilities for a specific workflow using either the workflow ID
- `keen workflows create` — Create workflows and related statuses.
- `keen workflows get-default-editor` — Get the user's default workflow editor. This can be either the new editor or the legacy editor.
- `keen workflows read` — Returns a list of workflows and related statuses by providing workflow names, workflow IDs, or project and issue types.
- `keen workflows read-previews` — Returns a requested workflow within a given project.
- `keen workflows search` — Returns a [paginated](#pagination) list of global and project workflows.
- `keen workflows update` — Update workflows and related statuses.
- `keen workflows validate-create` — Validate the payload for bulk create workflows.
- `keen workflows validate-update` — Validate the payload for bulk update workflows.

**workflowscheme** — Manage workflowscheme

- `keen workflowscheme assign-scheme-to-project` — Assigns a workflow scheme to a project. This operation is performed only when there are no issues in the project.
- `keen workflowscheme create-workflow-scheme` — Creates a workflow scheme.
- `keen workflowscheme delete-workflow-scheme` — Deletes a workflow scheme.
- `keen workflowscheme get-all-workflow-schemes` — Returns a [paginated](#pagination) list of all workflow schemes, not including draft workflow schemes.
- `keen workflowscheme get-required-workflow-scheme-mappings` — Gets the required status mappings for the desired changes to a workflow scheme.
- `keen workflowscheme get-workflow-scheme` — Returns a workflow scheme.
- `keen workflowscheme get-workflow-scheme-project-associations` — Returns a list of the workflow schemes associated with a list of projects.
- `keen workflowscheme read-workflow-schemes` — Returns a list of workflow schemes by providing workflow scheme IDs or project IDs.
- `keen workflowscheme switch-workflow-scheme-for-project` — Switches a workflow scheme for a project. Workflow schemes can only be assigned to classic projects.
- `keen workflowscheme update-schemes` — Updates company-managed and team-managed project workflow schemes.
- `keen workflowscheme update-workflow-scheme` — Updates a company-manged project workflow scheme, including the name, default workflow, issue type to project mappings

**worklog** — Manage worklog

- `keen worklog get-for-ids` — Returns worklog details for a list of worklog IDs. The returned list of worklogs is limited to 1000 items.
- `keen worklog get-ids-of-deleted-since` — Returns a list of IDs and delete timestamps for worklogs deleted after a date and time.
- `keen worklog get-ids-of-modified-since` — Returns a list of IDs and update timestamps for worklogs updated after a date and time.

**workspaces** — Org Workspaces APIs

- `keen workspaces <orgId>` — A workspace refers to a specific instance of an Atlassian product that is accessed through a unique URL.


### Finding the right command

When you know what you want to do but not which command does it, ask the CLI directly:

```bash
keen which "<capability in your own words>"
```

`which` resolves a natural-language capability query to the best matching command from this CLI's curated feature index. Exit code `0` means at least one match; exit code `2` means no confident match — fall back to `--help` or use a narrower query.

## Auth Setup
Run `keen auth setup` to print the URL and steps for getting a key (add `--launch` to open the URL). Then set:

```bash
export JIRA_CLOUD_PLATFORM_USERNAME="<your-key>"
```

Or persist it in `~/.config/keen/config.toml`.

Run `keen doctor` to verify setup.

## Agent Mode

Add `--agent` to any command. Expands to: `--json --compact --no-input --no-color --yes`.

- **Pipeable** — JSON on stdout, errors on stderr
- **Filterable** — `--select` keeps a subset of fields. Dotted paths descend into nested structures; arrays traverse element-wise. Critical for keeping context small on verbose APIs:

  ```bash
  keen addon-categories get mock-value --agent --select id,name,status
  ```
- **Previewable** — `--dry-run` shows the request without sending
- **Offline-friendly** — sync/search commands can use the local SQLite store when available
- **Non-interactive** — never prompts, every input is a flag
- **Explicit retries** — use `--idempotent` only when an already-existing create should count as success, and `--ignore-missing` only when a missing delete target should count as success

### Response envelope

Commands that read from the local store or the API wrap output in a provenance envelope:

```json
{
  "meta": {"source": "live" | "local", "synced_at": "...", "reason": "..."},
  "results": <data>
}
```

Parse `.results` for data and `.meta.source` to know whether it's live or local. A human-readable `N results (live)` summary is printed to stderr only when stdout is a terminal AND no machine-format flag (`--json`, `--csv`, `--compact`, `--quiet`, `--plain`, `--select`) is set — piped/agent consumers and explicit-format runs get pure JSON on stdout.

## Agent Feedback

When you (or the agent) notice something off about this CLI, record it:

```
keen feedback "the --since flag is inclusive but docs say exclusive"
keen feedback --stdin < notes.txt
keen feedback list --json --limit 10
```

Entries are stored locally at `~/.local/share/keen/feedback.jsonl`. They are never POSTed unless `JIRA_PP_CLI_FEEDBACK_ENDPOINT` is set AND either `--send` is passed or `JIRA_PP_CLI_FEEDBACK_AUTO_SEND=true`. Default behavior is local-only.

Write what *surprised* you, not a bug report. Short, specific, one line: that is the part that compounds.

## Output Delivery

Every command accepts `--deliver <sink>`. The output goes to the named sink in addition to (or instead of) stdout, so agents can route command results without hand-piping. Three sinks are supported:

| Sink | Effect |
|------|--------|
| `stdout` | Default; write to stdout only |
| `file:<path>` | Atomically write output to `<path>` (tmp + rename) |
| `webhook:<url>` | POST the output body to the URL (`application/json` or `application/x-ndjson` when `--compact`) |

Unknown schemes are refused with a structured error naming the supported set. Webhook failures return non-zero and log the URL + HTTP status on stderr.

## Named Profiles

A profile is a saved set of flag values, reused across invocations. Use it when a scheduled agent calls the same command every run with the same configuration - HeyGen's "Beacon" pattern.

```
keen profile save briefing --json
keen --profile briefing addon-categories get mock-value
keen profile list --json
keen profile show briefing
keen profile delete briefing --yes
```

Explicit flags always win over profile values; profile values win over defaults. `agent-context` lists all available profiles under `available_profiles` so introspecting agents discover them at runtime.

## Async Jobs

For endpoints that submit long-running work, the generator detects the submit-then-poll pattern (a `job_id`/`task_id`/`operation_id` field in the response plus a sibling status endpoint) and wires up three extra flags on the submitting command:

| Flag | Purpose |
|------|---------|
| `--wait` | Block until the job reaches a terminal status instead of returning the job ID immediately |
| `--wait-timeout` | Maximum wait duration (default 10m, 0 means no timeout) |
| `--wait-interval` | Initial poll interval (default 2s; grows with exponential backoff up to 30s) |

Use async submission without `--wait` when you want to fire-and-forget; use `--wait` when you want one command to return the finished artifact.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage error (wrong arguments) |
| 3 | Resource not found |
| 4 | Authentication required |
| 5 | API error (upstream issue) |
| 7 | Rate limited (wait and retry) |
| 10 | Config error |

## Argument Parsing

Parse `$ARGUMENTS`:

1. **Empty, `help`, or `--help`** → show `keen --help` output
2. **Starts with `install`** → ends with `mcp` → MCP installation; otherwise → see Prerequisites above
3. **Anything else** → Direct Use (execute as CLI command with `--agent`)

## MCP Server Installation

Install the MCP binary from this CLI's published public-library entry or pre-built release, then register it:

```bash
claude mcp add keen-mcp -- keen-mcp
```

Verify: `claude mcp list`

## Direct Use

1. Check if installed: `which keen`
   If not found, offer to install (see Prerequisites at the top of this skill).
2. Match the user query to the best command from the Unique Capabilities and Command Reference above.
3. Execute with the `--agent` flag:
   ```bash
   keen <command> [subcommand] [args] --agent
   ```
4. If ambiguous, drill into subcommand help: `keen <command> --help`.
