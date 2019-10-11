#!/usr/bin/env bash

# set -o errexit
# set -o nounset
# set -o pipefail
# set -o xtrace

# optionally set environment variable for
# 'google', 'azure', or 'amazon'.
: ${PLATFORM:="google"}
: ${TILLER_NAMESPACE:="nuodb"}
: ${DOMAIN_NAME:="cashews"}

ME=`basename $0`
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

: ${SELF_ROOT:=${SCRIPT_DIR%test/scripts}}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. ${SELF_ROOT}/test/scripts/profile.sh

RELEASES=( demo-ycsb restored-database backup database monitoring-influx monitoring-insights admin transparent-hugepage )
for RELEASE in "${RELEASES[@]}"
do
    helm delete --purge ${RELEASE}
    sleep 1
done

kubectl delete jobs --all
kubectl delete pvc --all
kubectl delete pv --all
