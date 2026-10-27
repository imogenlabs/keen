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

## Admin API

Admin endpoints use a different base URL (`api.atlassian.com`). The same token works if your account has org admin permissions.
