#!/usr/bin/env bash
# Regenerate keen from the current Atlassian OpenAPI specs without losing
# keen's own changes.
#
# How it works: keen's generated code is a three-way merge.
#   base   = the generator's output last time (renamed to keen), kept on the
#            branch build/generator-output
#   ours   = keen as it is now (HEAD), hand changes included
#   theirs = the generator's output today, from today's specs, renamed the same way
# git merges base->theirs onto ours, so spec changes land and keen's edits stay.
# Files keen maintains by hand (README, docs, CHANGELOG, AGENTS.md, the patch list)
# always keep keen's version.
#
# Usage (from the repo root, on a fresh branch):
#   scripts/regen/regen.sh                 # download current specs, generate, merge
#   scripts/regen/regen.sh --specs DIR     # use specs already in DIR
# Afterwards: resolve any conflicts it lists, then run scripts/regen/check.sh.
set -euo pipefail

REPO=$(git rev-parse --show-toplevel)
HERE="$REPO/scripts/regen"
# shellcheck source=/dev/null
source "$HERE/config.sh"

SPECS_IN=""
if [[ "${1:-}" == "--specs" ]]; then SPECS_IN=$2; fi

[[ -z "$(git -C "$REPO" status --porcelain)" ]] || { echo "regen: commit or stash your changes first" >&2; exit 1; }
command -v "$GENERATOR" >/dev/null || { echo "regen: generator not found: $GENERATOR (see config.sh)" >&2; exit 1; }

WORK=${REGEN_WORK:-$(mktemp -d "${TMPDIR:-/tmp}/keen-regen-$(date +%Y%m%d-%H%M%S)-XXXX")}
echo "regen: working in $WORK"
mkdir -p "$WORK/specs" "$WORK/base" "$WORK/ours" "$WORK/theirs"

# 1. Specs
if [[ -n "$SPECS_IN" ]]; then
  cp "$SPECS_IN"/*.json "$WORK/specs/"
else
  for entry in "${SPECS[@]}"; do
    name=${entry%%=*}; url=${entry#*=}
    curl -fsSL "$url" -o "$WORK/specs/$name.json"
    python3 -c "import json,sys; d=json.load(open(sys.argv[1])); assert 'paths' in d" "$WORK/specs/$name.json" \
      || { echo "regen: $name spec from $url is not an OpenAPI document" >&2; exit 1; }
  done
fi

# 2. Generate today's output (named as the base was), then rename it to keen.
args=()
for name in "${SPEC_ORDER[@]}"; do args+=(--spec "$WORK/specs/$name.json"); done
(cd "$WORK" && "$GENERATOR" generate "${args[@]}" "${GENERATOR_FLAGS[@]}" --output "$WORK/theirs" > "$WORK/generate.log" 2>&1) \
  || { echo "regen: generator failed, see $WORK/generate.log" >&2; exit 1; }
python3 "$HERE/normalize.py" "$WORK/theirs"
# The generator records where it read the specs; keep that path inside the repo.
grep -rlF "$WORK/" "$WORK/theirs" 2>/dev/null | while IFS= read -r f; do
  python3 -c "import sys; p,w=sys.argv[1],sys.argv[2]; t=open(p).read(); open(p,'w').write(t.replace(w,''))" "$f" "$WORK/"
done
mkdir -p "$WORK/theirs/specs"
cp "$WORK/specs/"*.json "$WORK/theirs/specs/"
for f in "${GENERATOR_ONLY[@]}"; do [[ -e "$WORK/theirs/$f" ]] && mv "$WORK/theirs/$f" "$WORK/$(basename "$f").generated"; done

# 3. Base and ours
git -C "$REPO" archive "$BASE_REF" | tar -x -C "$WORK/base"
git -C "$REPO" archive HEAD | tar -x -C "$WORK/ours"

# 4. Three-way merge in a scratch repository.
M="$WORK/merge"
git init -q "$M"
snap() { # snap <dir> <message>
  git --git-dir="$M/.git" --work-tree="$1" read-tree --empty
  git --git-dir="$M/.git" --work-tree="$1" add -A -f .
  git --git-dir="$M/.git" --work-tree="$1" -c user.name=regen -c user.email=regen@localhost commit -q -m "$2"
}
snap "$WORK/base" base
git --git-dir="$M/.git" branch -q -M base
git --git-dir="$M/.git" branch -q theirs base
git --git-dir="$M/.git" branch -q ours base
git --git-dir="$M/.git" symbolic-ref HEAD refs/heads/theirs; snap "$WORK/theirs" theirs
git --git-dir="$M/.git" symbolic-ref HEAD refs/heads/ours;   snap "$WORK/ours" ours

for f in "${KEEN_OWNED[@]}"; do echo "$f merge=ours"; done > "$M/.git/info/attributes"
git --git-dir="$M/.git" config merge.ours.driver true
git --git-dir="$M/.git" --work-tree="$WORK/ours" checkout -q -f ours
set +e
# No rename detection: a retired command and a new one look alike and would be paired up.
git --git-dir="$M/.git" --work-tree="$WORK/ours" -c user.name=regen -c user.email=regen@localhost \
  -c merge.renames=false -c diff.renames=false \
  merge -q --no-edit -m "merge generator output" theirs > "$WORK/merge.log" 2>&1
merge_status=$?
set -e

# keen's own files: when the generator dropped or changed one, keep keen's copy.
if [[ $merge_status -ne 0 ]]; then
  git --git-dir="$M/.git" --work-tree="$WORK/ours" diff --name-only --diff-filter=U > "$WORK/conflicts.txt"
  while IFS= read -r f; do
    for pat in "${KEEN_OWNED[@]}"; do
      # shellcheck disable=SC2053
      if [[ "$f" == $pat ]]; then
        git --git-dir="$M/.git" --work-tree="$WORK/ours" checkout -q --ours -- "$f"
        git --git-dir="$M/.git" --work-tree="$WORK/ours" add -f -- "$f"
        echo "regen: kept keen's $f"
        break
      fi
    done
  done < "$WORK/conflicts.txt"
  if [[ -z "$(git --git-dir="$M/.git" --work-tree="$WORK/ours" diff --name-only --diff-filter=U)" ]]; then
    git --git-dir="$M/.git" --work-tree="$WORK/ours" -c user.name=regen -c user.email=regen@localhost \
      commit -q --no-edit -m "merge generator output"
    merge_status=0
  fi
fi

# 5. Bring the merged tree into this checkout (git updates and removes tracked files).
if [[ $merge_status -ne 0 ]]; then
  echo "regen: merge has conflicts. Resolve them in $WORK/ours (look for <<<<<<<), then run:"
  echo "  REGEN_WORK=$WORK scripts/regen/finish.sh"
  git --git-dir="$M/.git" --work-tree="$WORK/ours" diff --name-only --diff-filter=U
  exit 2
fi
git -C "$REPO" fetch -q "$M" +theirs:refs/regen/theirs +ours:refs/regen/merged
git -C "$REPO" read-tree -u --reset refs/regen/merged
python3 "$HERE/postfix.py" "$REPO"
echo "regen: merged. Review with git status / git diff, then run scripts/regen/check.sh"
# Today's generator output, stacked on the base so moving the base is a fast-forward.
next=$(git -C "$REPO" -c user.name=regen -c user.email=regen@localhost commit-tree "refs/regen/theirs^{tree}" \
  -p "$(git -C "$REPO" rev-parse "$BASE_REF")" -m "generator output $(date +%Y-%m-%d), renamed to keen")
git -C "$REPO" update-ref refs/regen/next-base "$next"
echo "regen: after this lands on main, move the base forward:"
echo "  git push origin refs/regen/next-base:refs/heads/build/generator-output"
