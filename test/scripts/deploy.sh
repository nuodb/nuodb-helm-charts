#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
# set -o xtrace

# optionally set environment variable for
# 'google', 'azure', or 'amazon'.
: ${PLATFORM:="google"}
: ${TILLER_NAMESPACE:="nuodb"}
: ${DOMAIN_NAME:="cashews"}
: ${VALUE_CLASS:="tiny"}

# optionally set to 'nuodb' or nuodb-incubator'
# to test published artifacts in GCS-backed repos
: ${REPO_NAME:="stable"}

ME=`basename $0`
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

: ${SELF_ROOT:=${SCRIPT_DIR%test/scripts}}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. ${SELF_ROOT}/test/scripts/profile.sh

values_option="--values ${SELF_ROOT}/samples/${VALUE_CLASS}.yaml"

# permit qa to flexibly override values as code
: ${VALUES_OPTION_OVERRIDES_FILE:="${SELF_ROOT}/.helm.yaml"}

values_option_overrides=""
if [ -f ${VALUES_OPTION_OVERRIDES_FILE} ]; then
  values_option_overrides="--values ${VALUES_OPTION_OVERRIDES_FILE}"
fi

values_overrides="--set cloud.provider=${PLATFORM}"

if [ ${PLATFORM} == "amazon" ]; then
  values_overrides="${values_overrides} --set cloud.zones={us-east-2a,us-east-2b,us-east-2c}"
elif [ ${PLATFORM} == "azure" ]; then
  values_overrides="${values_overrides} --set cloud.zones={'0','1','2'}"
elif [ ${PLATFORM} == "google" ]; then
  values_overrides="${values_overrides} --set cloud.zones={us-central1-a,us-central1-b,us-central1-c}"
fi

log_error () {
  printf '\e[31mERROR: %s\n\e[39m' "$1" >&2
}

restart () {
  local term="${1?Specify term}"
  kubectl get pods --all-namespaces | grep ${term} | awk '{ print $5 }' &> /dev/null
  [ $? != 0 ];
}

status () {
  local term="${1?Specify term}"
  local expected="${2?Specify expected}"
  result=( $(kubectl get pods --all-namespaces | grep ${term} | awk '{ print $4 }') )
  [ $result == ${expected} ];
}

ready () {
  local term="${1?Specify term}"
  local count="${2?Specify count}"
  result=( $(kubectl get pods --all-namespaces | grep ${term} | awk '{ print $3 }') )
  [ $result == $count ];
}

check () {
  local query="${1?Specify query}"
  local ready="${2?Specify count}"
  local state="${3?Specify state}"
  if restart ${query} ; then log_error "${query} restarted > 0" ; fi
  if ! status ${query} ${state} ; then log_error "${query} status not ${state}" ; fi
  if ! ready ${query} ${ready} ; then log_error "${query} ready count not ${ready}" ; fi
}

check-pod () {
  local query="${1?Specify query}"
  local ready="${2?Specify count}"
  check $query $ready "Running"
}

check-job () {
  local query="${1?Specify query}"
  local ready="${2?Specify count}"
  check $query $ready "Completed"
}

helm install ${REPO_NAME}/transparent-hugepage --name transparent-hugepage \
  ${values_option} \
  ${values_option_overrides} \
  ${values_overrides}

sleep 30

check-pod "transparent-hugepage" "1/1"

helm install ${REPO_NAME}/admin --name admin \
  ${values_option} \
  ${values_option_overrides} \
  ${values_overrides} \
  --set admin.domain=${DOMAIN_NAME}

sleep 180

check-pod "admin-${DOMAIN_NAME}-0" "1/1"

# check and delete the lb policy jobs

check-job "job-lb-policy-nearest" "0/1"
kubectl delete jobs --all

kubectl scale sts admin-${DOMAIN_NAME} --replicas=3

sleep 180

check-pod "admin-${DOMAIN_NAME}-1" "1/1"
check-pod "admin-${DOMAIN_NAME}-2" "1/1"

CHARTS=( monitoring-influx monitoring-insights )
for CHART in "${CHARTS[@]}"
do
  helm install ${REPO_NAME}/$CHART -n $CHART \
    ${values_option} \
    ${values_option_overrides} \
    ${values_overrides} \
    --set admin.domain=${DOMAIN_NAME}
done

sleep 60 

check-pod "nuodb-dashboard-display" "1/1"
check-pod "nuodb-insights" "2/2"

helm install ${REPO_NAME}/database --name database \
  ${values_option} \
  ${values_option_overrides} \
  ${values_overrides} \
  --set admin.domain=${DOMAIN_NAME}

sleep 300

check-pod "sm-database-cashews-demo-0" "1/1"
check-pod "sm-database-cashews-demo-hotcopy-0" "1/1"
check-pod "te-database-cashews-demo" "1/1"

CHARTS=( backup )
for CHART in "${CHARTS[@]}"
do
  helm install ${REPO_NAME}/$CHART -n $CHART \
    ${values_option} \
    ${values_option_overrides} \
    ${values_overrides} \
    --set admin.domain=${DOMAIN_NAME}
done

sleep 30
kubectl delete job --all

kubectl scale rc ycsb-load --replicas=1

sleep 30

check-pod "ycsb-load" "1/1"

job_name=job-backup-full-$(python -c 'import random ; print "".join(map(lambda t: format(t, "02x"), [random.randrange(256) for x in range(6)]))')
kubectl create job $job_name \
  --from=cronjob/full-backup-demo-cronjob

sleep 60

check-job "job-backup-full" "0/1"
kubectl delete job --all

for num in {1..3}
do

  job_name=job-backup-incr-$(python -c 'import random ; print "".join(map(lambda t: format(t, "02x"), [random.randrange(256) for x in range(6)]))')
  kubectl create job $job_name \
    --from=cronjob/incremental-backup-demo-cronjob

  sleep 60

  check-job "job-backup-incr" "0/1"
  kubectl delete job --all

done
