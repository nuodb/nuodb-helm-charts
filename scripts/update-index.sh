#!/bin/sh

fail() {
    printf "$1\n" >&2
    exit 1
}

set -e

# Make sure there are no uncommitted changes
GIT_STATUS="$(git status --porcelain)"
[ "$GIT_STATUS" = "" ] || fail "Cannot publish charts with uncommitted changes:\n$GIT_STATUS"

# Change to root directory and make sure package directory exists
cd "$(dirname "$0")/.."
[ -d package ] || fail "package directory does not exist... run package.sh first"

# Checkout gh-pages and fast forward to origin
git checkout gh-pages
git merge --ff-only origin/gh-pages

# Update index with new Helm charts
: ${GH_RELEASES_URL:="https://github.com/nuodb/nuodb-helm-charts/releases/download"}
helm repo index package --merge index.yaml --url "$GH_RELEASES_URL"

# Commit and push change if PUSH_UPDATE=true
if [ "$PUSH_UPDATE" = true ]; then
    mv package/index.yaml .
    git add index.yaml
    git commit -m "Adding $(ls package | sed 's|/||') charts to index"
    git push
fi
