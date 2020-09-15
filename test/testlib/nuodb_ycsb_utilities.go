package testlib

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func StartYCSBWorkload(t *testing.T, namespaceName string, options *helm.Options) (helmChartReleaseName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	InjectTestValues(t, options)

	helmChartReleaseName = fmt.Sprintf("ycsb-%s", randomSuffix)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	// with Async actions which do not return a cleanup method, create the teardown(s) first
	AddTeardown(TEARDOWN_YCSB, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})

	helm.Install(t, options, YCSB_HELM_CHART_PATH, helmChartReleaseName)

	Await(t, func() bool {
		return GetReplicationController(t, namespaceName, helmChartReleaseName) != nil
	}, 30*time.Second)

	return
}
