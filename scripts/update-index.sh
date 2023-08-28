#!/bin/sh

fail() {
    printf "$1\n" >&2
    exit 1
}

check_package() {
    [ -d "package/$1" ] || fail "package/$1 directory does not exist... run package.sh first"
}

update_index() {
    # Merge new charts into current index
    : ${GH_RELEASES_URL:="https://github.com/nuodb/nuodb-helm-charts/releases/download"}
    helm repo index "package/$1" --merge "$1/index.yaml" --url "$GH_RELEASES_URL"

    # Stage update to index file
    mv "package/$1/index.yaml" "$1/index.yaml"
    git add "$1/index.yaml"
}

set -e

# Make sure there are no uncommitted changes
GIT_STATUS="$(git status --porcelain)"
[ "$GIT_STATUS" = "" ] || fail "Cannot publish charts with uncommitted changes:\n$GIT_STATUS"

# Change to root directory and make sure package directories exist
cd "$(dirname "$0")/.."
check_package stable
check_package incubator

# Checkout gh-pages and fast forward to origin
git checkout gh-pages
git merge --ff-only origin/gh-pages

# Update indexes with new Helm charts
update_index stable
update_index incubator

# Commit changes to indexes
git commit -m "Add $(ls package/stable | sed 's|/||') charts to indexes"

# Push changes if PUSH_UPDATE=true
if [ "$PUSH_UPDATE" = true ]; then
    git push
fi
