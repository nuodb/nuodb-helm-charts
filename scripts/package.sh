#!/bin/sh

fail() {
    printf "$1\n" >&2
    exit 1
}

set -e

# Change to root directory
cd "$(dirname "$0")/.."

# Create directory for packaged Helm charts
dest="package/v$(./get-version.sh)"
mkdir -p "$dest"

# Package all Helm charts
charts="$(find stable -maxdepth 2 -name Chart.yaml -exec dirname {}  \;)"
for dir in $charts; do
    helm package "$dir" -d "$dest"
done
