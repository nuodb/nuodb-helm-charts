#!/usr/bin/env bash

: ${CLUSTER_NAME:="helmtest"}

ME=`basename $0`
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

: ${SELF_ROOT:=${SCRIPT_DIR%test/scripts}}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. ${SELF_ROOT}/test/scripts/profile.sh

if [ "${PLATFORM}" == "google" ]; then
    gcloud container clusters get-credentials ${CLUSTER_NAME} --zone ${ZONE} --project cit-team
    kubectl config set-context $(kubectl config current-context) --namespace=${TILLER_NAMESPACE}
    kubectl config view | grep -A10 "name: $(kubectl config current-context)" | awk '$1=="access-token:"{print $2}'
elif [ "${PLATFORM}" == "azure" ]; then
    az login
    az aks get-credentials --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}
fi
