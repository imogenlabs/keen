#!/usr/bin/env bash
# After resolving regen.sh merge conflicts in $REGEN_WORK/ours, bring the
# result into this checkout.   Usage: REGEN_WORK=<dir> scripts/regen/finish.sh
set -euo pipefail
# shellcheck source=/dev/null
source "$(git rev-parse --show-toplevel)/scripts/regen/config.sh"
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
git -C "$REPO" fetch -q "$M" +theirs:refs/regen/theirs +ours:refs/regen/merged
git -C "$REPO" read-tree -u --reset refs/regen/merged
python3 "$REPO/scripts/regen/postfix.py" "$REPO"
python3 "$REPO/scripts/regen/graphql.py" --write
echo "finish: merged tree is in the checkout. Run scripts/regen/check.sh"
# Today's generator output, stacked on the base so moving the base is a fast-forward.
next=$(git -C "$REPO" -c user.name=regen -c user.email=regen@localhost commit-tree "refs/regen/theirs^{tree}" \
  -p "$(git -C "$REPO" rev-parse "$BASE_REF")" -m "generator output $(date +%Y-%m-%d), renamed to keen")
git -C "$REPO" update-ref refs/regen/next-base "$next"
echo "finish: after this lands on main, move the base forward:"
echo "  git push origin refs/regen/next-base:refs/regen/base"
