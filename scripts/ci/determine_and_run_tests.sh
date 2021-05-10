#!/usr/bin/env bash

# (C) Copyright NuoDB, Inc. 2019-2021  All Rights Reserved
# This file is licensed under the BSD 3-Clause License.
# See https://github.com/nuodb/nuodb-helm-charts/blob/master/LICENSE

echo "Running $TEST_SUITE"

if [[ $TEST_SUITE = "basic"  ]]; then
  ./test/test_suite.sh
  go test -timeout 10m -v ./test/integration
elif [[ $TEST_SUITE = "minikube-short"  ]]; then
  go test -timeout 50m -v ./test/minikube -tags=short
elif [[ $TEST_SUITE = "minikube-long"  ]]; then
  go test -timeout 50m -v ./test/minikube -tags=long
elif [[ $TEST_SUITE = "minikube-diagnostics"  ]]; then
  go test -timeout 50m -v ./test/minikube -tags=diagnostics
elif [[ $TEST_SUITE = "minikube-upgrade"  ]]; then
  go test -timeout 50m -v ./test/minikube -tags=upgrade
fi