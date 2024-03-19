#!/bin/sh

: ${HELM_REPO_LABEL:="nuodb-helm-repo"}
kubectl get configmaps,services,deployments.apps -l "app.kubernetes.io/instance=$HELM_REPO_LABEL" -o name | xargs -r kubectl delete
