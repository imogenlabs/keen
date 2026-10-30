# Authentication

## Jira Cloud API Token

Generate at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

### Config File

`~/.config/keen/auth.json` (used by `keen` and `keen-mcp`):

```json
{ "base_url": "https://your-site.atlassian.net", "username": "you@example.com", "password": "your-api-token" }
```

Keep it private: `chmod 600 ~/.config/keen/auth.json`.

Both also read `~/.config/keen/config.toml` when there is no `auth.json`. Note the key names:

```toml
base_url = 'https://your-site.atlassian.net'
cloud_platform_username = 'you@example.com'
cloud_platform_password = 'your-api-token'
```

### Environment Variables

Environment variables override either file and work for both `keen` and `keen-mcp`:

```bash
export JIRA_CLOUD_PLATFORM_USERNAME="you@example.com"
export JIRA_CLOUD_PLATFORM_PASSWORD="your-api-token"
export JIRA_CLOUD_PLATFORM_URL="https://your-site.atlassian.net"
```

`JIRA_PP_CLI_BASE_URL` also sets the site and wins if both are set.

### Verify

```bash
keen doctor
```

Expected output:
```
  OK Config: ok
  OK Auth: configured
  OK Credentials: ok
```

## Marketplace API

The same Jira API token works for Marketplace operations. The CLI automatically routes requests to `marketplace.atlassian.com` for addon/vendor endpoints.

## Admin APIs

Atlassian's admin APIs on `api.atlassian.com` do not take your email and API token. keen sends each one the key it needs and uses your normal login for everything else:

| Commands | Address | Key |
|----------|---------|-----|
| Organizations, policies, API tokens, service accounts, data loss prevention, user management (`orgs`, `admin-control`, `api-tokens`, `api-keys`, `service-accounts`, `user-management`, `user-provisioning`, ...) | `/admin/...`, `/users/...` | Organization admin API key: admin.atlassian.com > Settings > API keys |
| User provisioning (`scim`) | `/scim/...` | The API key of your identity-provider directory: admin.atlassian.com > Security > Identity providers |

```bash
export KEEN_ADMIN_API_KEY="your-org-admin-api-key"
export KEEN_SCIM_API_KEY="your-directory-api-key"   # only for scim commands
```

Or add `admin_api_key` and `scim_api_key` to `~/.config/keen/auth.json`. Without the key, an admin command stops before calling Atlassian and says which key to set. `keen doctor` shows whether each key is set.

## Assets

Assets commands (all named `assets-...`: `assets-objectschema`, `assets-object`, `assets-aql`, `assets-import-source`, ...) use your normal login. Their address needs your site's cloud id and Assets workspace id; keen looks both up from your site the first time a command needs them. To skip the lookup, set them yourself:

```bash
export KEEN_CLOUD_ID="your-cloud-id"
export KEEN_ASSETS_WORKSPACE_ID="your-workspace-id"
```

or `cloud_id` and `assets_workspace_id` in `~/.config/keen/auth.json`. A `--dry-run` shows the address with `KEEN_CLOUD_ID` and `KEEN_ASSETS_WORKSPACE_ID` in place of the ids.
