#!/bin/sh

fail() {
    printf "$1\n" >&2
    exit 1
}

# Change to root directory
cd "$(dirname "$0")"

# Scrape version from all Chart.yaml files at second level of directory hierarchy
VERSION="$(find stable -maxdepth 2 -name Chart.yaml | xargs -r sed -n 's/^version *: *//p' | sort | uniq)"

# Check that a unique version was found and return it
[ "$VERSION" = "$(echo "$VERSION" | head -n1)" ] || fail "Inconsistent versions in Chart.yaml:\n$VERSION"
echo "$VERSION"
