#!/usr/bin/env bash

# optionally set environment variable for
# 'google', 'azure', or 'amazon'.
: ${PLATFORM:="google"}
: ${TILLER_NAMESPACE:="nuodb"}
: ${DOMAIN_NAME:="cashews"}

ME=`basename $0`
SCRIPT_DIR=`dirname $0`

: ${SELF_ROOT:=`dirname ${SCRIPT_DIR}`}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

. $SELF_ROOT/test/scripts/functions.sh

exit_code=0

run_tests () {
  cd $SELF_ROOT

  # run bash-based tests...
  for test in `find test -type f -name "*_test.sh"`; do
    echo "running $test..."
    if ! ${test}; then
      log_error "Problem running test: ${test}"
      exit_code=1
    fi
  done

}

run_tests

exit ${exit_code}
