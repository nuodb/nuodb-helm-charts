#!/bin/sh

fail() {
    printf "$1\n" >&2
    exit 1
}

package() {
    # Create directory for packaged Helm charts
    dest="package/$1/v$(./get-version.sh)"
    mkdir -p "$dest"

    # Package all Helm charts under specified directory
    charts="$(find "$1" -maxdepth 2 -name Chart.yaml -exec dirname {}  \;)"
    for dir in $charts; do
        helm package "$dir" -d "$dest"
    done
}

set -e

# Change to root directory
cd "$(dirname "$0")/.."

# Package Helm charts
package stable
package incubator
