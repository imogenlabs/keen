# Authentication

## Jira Cloud API Token

Generate at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

### Environment Variables

```bash
export JIRA_CLOUD_PLATFORM_USERNAME="you@example.com"
export JIRA_CLOUD_PLATFORM_PASSWORD="your-api-token"
```

### Config File

`~/.config/keen/config.toml`:

```toml
base_url = 'https://your-domain.atlassian.net'
username = 'you@example.com'
password = 'your-api-token'
```

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
