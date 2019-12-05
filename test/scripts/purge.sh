#!/usr/bin/env bash

ME=`basename $0`
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

# optionally set to 'nuodb' or nuodb-incubator'
# to test published artifacts in GCS-backed repos
: ${REPO_NAME:="stable"}

: ${SELF_ROOT:=${SCRIPT_DIR%test/scripts}}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. ${SELF_ROOT}/test/scripts/profile.sh

# =========================================================

CHARTS=( admin database backup )
for CHART in "${CHARTS[@]}"
do
    helm del --purge $CHART
done

for j in $(kubectl get jobs -o custom-columns=:.metadata.name)
do
    kubectl delete jobs $j
done

kubectl delete pvc --all
kubectl delete pv --all
kubectl delete namespace nuodb

# =========================================================

# delete storage classes...

echo "deleting storage class chart..."
helm del --purge --tiller-namespace kube-system storage-class

# delete cluster scoped tiller...

echo "deleting system tiller..."

kubectl -n kube-system delete deployment tiller-deploy
kubectl delete clusterrolebinding tiller-system
kubectl -n kube-system delete serviceaccount tiller-system
kubectl delete service tiller-deploy -n kube-system

# delete dashboards, etc...
if [ ! "${PLATFORM}" == "azure" ]; then
    echo "deleting dashboards..."
    kubectl delete -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
    kubectl delete -f ${SELF_ROOT}/test/files/dashboard-adminuser.yaml
fi

if [ "${PLATFORM}" == "azure" ]; then
    kubectl delete clusterrolebinding kubernetes-dashboard
fi
