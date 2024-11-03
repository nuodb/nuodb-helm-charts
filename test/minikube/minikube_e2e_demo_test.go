//go:build long
// +build long

package minikube

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	corev1 "k8s.io/api/core/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestKubernetesYCSB(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	defer testlib.Teardown(testlib.TEARDOWN_YCSB)

	testlib.StartYCSBWorkload(t, namespaceName, &helm.Options{})
	testlib.ScaleYCSB(t, namespaceName, 1)

	ycsbPodName := testlib.GetPodName(t, namespaceName, testlib.YCSB_CONTROLLER_NAME)
	go testlib.GetAppLog(t, namespaceName, ycsbPodName, "-ycsb", &corev1.PodLogOptions{Follow: true})

	// let YCSB run for a couple of seconds
	time.Sleep(5 * time.Second)
}
