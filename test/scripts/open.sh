#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ME=`basename $0`
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

: ${SELF_ROOT:=${SCRIPT_DIR%test/scripts}}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. ${SELF_ROOT}/test/scripts/profile.sh

if [ "${PLATFORM}" == "google" ]; then

  echo "note bearer token to log into dashboard:"
  kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}')
  echo ""

  dashboard_url="http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/overview?namespace=default"
  echo "connect to:"
  echo "$dashboard_url"
  echo ""

elif [ "${PLATFORM}" == "azure" ]; then

  dashboard_url="http://localhost:8001/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy/"
  echo "connect to:"
  echo "$dashboard_url"
  echo ""

elif [ "${PLATFORM}" == "amazon" ]; then

  echo "note bearer token to log into dashboard:"
  kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')
  echo ""

  dashboard_url="http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login"
  echo "connect to:"
  echo "$dashboard_url"
  echo ""
  
fi

# opens a browser to the site...
open "$dashboard_url"

echo "starting proxy..."
kubectl proxy
echo ""