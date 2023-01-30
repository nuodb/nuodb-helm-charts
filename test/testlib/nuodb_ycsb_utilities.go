package testlib

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	corev1 "k8s.io/api/core/v1"
)

const YCSB_CONTROLLER_NAME = "ycsb-load"

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

	AddDiagnosticTeardown(TEARDOWN_YCSB, t, func() {
		DescribePods(t, namespaceName, YCSB_CONTROLLER_NAME)
	})

	if options.Version == "" {
		helm.Install(t, options, YCSB_HELM_CHART_PATH, helmChartReleaseName)
	} else {
		helm.Install(t, options, "nuodb-incubator/demo-ycsb ", helmChartReleaseName)
	}

	Await(t, func() bool {
		return GetReplicationController(t, namespaceName, helmChartReleaseName) != nil
	}, 30*time.Second)

	return
}

func ScaleYCSB(t *testing.T, namespaceName string, replicas int) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.RunKubectl(t, kubectlOptions, "scale", "replicationcontroller", YCSB_CONTROLLER_NAME,
		fmt.Sprintf("--replicas=%d", replicas))

	if replicas > 0 {
		AwaitNrReplicasScheduled(t, namespaceName, YCSB_CONTROLLER_NAME, replicas)
		// wait longer for the the YCSB image to be pulled and then check that all
		// requested replicas are reported READY
		AwaitPodPhase(t, namespaceName, YCSB_CONTROLLER_NAME, corev1.PodRunning, 300*time.Second)
		AwaitNrReplicasReady(t, namespaceName, YCSB_CONTROLLER_NAME, replicas)
	} else {
		AwaitNoPods(t, namespaceName, YCSB_CONTROLLER_NAME)
	}
}
