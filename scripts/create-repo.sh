#!/bin/sh

cd "$(dirname "$0")"

check_command() {
    if ! command -v "$1" >/dev/null; then
        echo "Command not found: $1"
        exit 1
    fi
}

set -e

check_command kubectl
check_command helm

: ${KUSTOMIZE_VERSION:="5.3.0"}

if command -v go >/dev/null; then
    : ${OS:="$(go env GOOS)"}
    : ${ARCH:="$(go env GOARCH)"}
fi

# Download kustomize if not available
if ! command -v kustomize >/dev/null && [ ! -x bin/kustomize ]; then
    : ${OS?="OS must be specified to download kustomize"}
    : ${ARCH?="Architecture must be specified to download kustomize"}

    mkdir -p bin
    curl -s -L "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" | tar -x -C bin
fi

# Make sure kustomize is on system path
if ! command -v kustomize >/dev/null; then
    export PATH="$(pwd)/bin:$PATH"
fi

# Package Helm charts from current commit
rm -rf ../package/
rm -rf repo/charts/
mkdir -p repo/charts
./package.sh
find ../package/stable -name \*.tgz -exec cp {} repo/charts/ \;

# Download published Helm repo index and merge charts from current commit
: ${HELM_REPO_URL:="http://nuodb-helm-repo"}
curl -s https://nuodb.github.io/nuodb-helm-charts/index.yaml > repo/charts/index.yaml
helm repo index --merge repo/charts/index.yaml --url "$HELM_REPO_URL" repo/charts

# Delete current Helm repo if one exists
: ${HELM_REPO_LABEL:="nuodb-helm-repo"}
./delete-repo.sh

# Create Helm repo
cd repo
echo > kustomization.yaml
kustomize edit add resource deployment.yaml
kustomize edit add resource service.yaml
kustomize edit add configmap charts --disableNameSuffixHash --from-file=charts/*
kustomize edit add label "app.kubernetes.io/instance:$HELM_REPO_LABEL" --without-selector
kustomize edit add label "app.kubernetes.io/name:$HELM_REPO_LABEL" --without-selector
kustomize build | kubectl create -f -

echo "Helm repository created from commit $(git rev-parse --short HEAD)."
echo "To use repository in NuoDB Operator:"
echo "  helm install nuodb-cp-operator nuodb-cp/nuodb-cp-operator ... --set cpOperator.nuodbRepoOverride=\"$HELM_REPO_URL\""
