#!/usr/bin/env bash

# optionally set environment variable for
# 'google', 'azure', or 'amazon'.
: ${PLATFORM:="google"}
: ${TILLER_NAMESPACE:="nuodb"}
: ${DOMAIN_NAME:="cashews"}
: ${VALUE_CLASS:="tiny"}

ME=`basename $0`
SCRIPT_DIR=`dirname $0`

: ${SELF_ROOT:=`dirname ${SCRIPT_DIR}`}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

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
  values_overrides="${values_overrides} --set cloud.zones={'0','1'}"
elif [ ${PLATFORM} == "google" ]; then
  values_overrides="${values_overrides} --set cloud.zones={us-central1-a,us-central1-b,us-central1-c}"
fi

. $SELF_ROOT/test/scripts/functions.sh

exit_code=0

(
  cd $SELF_ROOT
  for dir in `find ${SELF_ROOT}/stable -mindepth 1 -maxdepth 1 -type d`; do
      if test ! -f "$dir/Chart.yaml"; then
          continue
      fi
      if ! helm lint ${values_option} ${values_option_overrides} ${values_overrides} --set admin.domain=${DOMAIN_NAME} ${dir}; then
        log_error "Problem linting charts. Failures in '$dir'."
        exit_code=1
      fi
  done
)

exit ${exit_code}
