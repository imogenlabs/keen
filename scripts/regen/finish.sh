#!/usr/bin/env bash
# After resolving regen.sh merge conflicts in $REGEN_WORK/ours, bring the
# result into this checkout.   Usage: REGEN_WORK=<dir> scripts/regen/finish.sh
set -euo pipefail
REPO=$(git rev-parse --show-toplevel)
: "${REGEN_WORK:?set REGEN_WORK to the directory regen.sh printed}"
M="$REGEN_WORK/merge"
if grep -rl --include='*.go' --include='*.json' --include='*.md' '^<<<<<<< ' "$REGEN_WORK/ours" >/dev/null 2>&1; then
  echo "finish: conflict markers remain:" >&2
  grep -rl --include='*.go' --include='*.json' --include='*.md' '^<<<<<<< ' "$REGEN_WORK/ours" >&2
  exit 1
fi
git --git-dir="$M/.git" --work-tree="$REGEN_WORK/ours" add -A -f .
git --git-dir="$M/.git" --work-tree="$REGEN_WORK/ours" -c user.name=regen -c user.email=regen@localhost \
  commit -q -m "merge generator output"
git -C "$REPO" fetch -q "$M" theirs:refs/regen/theirs ours:refs/regen/merged
git -C "$REPO" read-tree -u --reset refs/regen/merged
python3 "$REPO/scripts/regen/postfix.py" "$REPO"
echo "finish: merged tree is in the checkout. Run scripts/regen/check.sh"
