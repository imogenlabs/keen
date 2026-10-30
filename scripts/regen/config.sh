# Settings for scripts/regen. Sourced by regen.sh and check.sh.

# CLI Printing Press binary. keen is generated with v4.15 (built 2026-05-25).
GENERATOR=${GENERATOR:-$HOME/go/bin/cli-printing-press}

# The generator is run under its original name so the output lines up with the
# base; normalize.py renames it to keen afterwards.
GENERATOR_FLAGS=(--name jira-pp-cli --owner smashingtags --lenient
  --auth-preference basicAuth --max-resources 1500 --validate=false)

# Previous generator output, renamed. Kept on GitHub as the hidden ref refs/regen/base
# (not a branch), fetched by regen.sh, moved forward after each regen lands on main.
BASE_REF=${BASE_REF:-refs/regen/base}

# Every Atlassian API keen is built from, one per line: name=spec URL.
# To add an API, add its line at the END (the order is the generator's order, and
# appending keeps every existing command where it is), then run regen.sh.
# The weekly spec-check workflow watches every URL here.
SPECS=(
  "jira-platform=https://developer.atlassian.com/cloud/jira/platform/swagger-v3.v3.json"
  "jira-agile=https://developer.atlassian.com/cloud/jira/software/swagger.v3.json"
  "confluence=https://developer.atlassian.com/cloud/confluence/swagger.v3.json"
  "marketplace=https://developer.atlassian.com/platform/marketplace/swagger.v3.json"
  "jira-service-desk=https://developer.atlassian.com/cloud/jira/service-desk/swagger.v3.json"
  "admin-orgs=https://developer.atlassian.com/cloud/admin/organization/swagger.v3.json"
  "confluence-v2=https://developer.atlassian.com/cloud/confluence/openapi-v2.v3.json"
)

# Files keen writes by hand: the merge always keeps keen's copy.
KEEN_OWNED=(README.md AGENTS.md CHANGELOG.md CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md
  NOTICE LICENSE .printing-press-patches.json ".github/**" "docs/**" "scripts/**" .gitignore go.sum)

# Generator byproducts that are not part of keen.
GENERATOR_ONLY=(manifest.json spec.json .printing-press-patches.json)

# Old names that must never appear in keen. Env var names JIRA_PP_CLI_* are kept
# on purpose (existing setups use them) and are not matched here.
OLD_NAMES='jira-pp-cli|jira-pp-mcp|pp-cli|pp-mcp|Jira Pp Cli|Jira PP CLI|Printing Press CLI for|atlas-cli|Atlas CLI'
