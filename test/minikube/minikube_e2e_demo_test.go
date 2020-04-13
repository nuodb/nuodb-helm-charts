// +build long

package minikube

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
)

const YCSB_CONTROLLER_NAME = "ycsb-load"

func scaleYCSB(t *testing.T, namespaceName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	k8s.RunKubectl(t, kubectlOptions, "scale", "replicationcontroller", YCSB_CONTROLLER_NAME, "--replicas=1")

	testlib.AwaitNrReplicasScheduled(t, namespaceName, YCSB_CONTROLLER_NAME, 1)
	testlib.AwaitNrReplicasReady(t, namespaceName, YCSB_CONTROLLER_NAME, 1)

}

func TestKubernetesYCSB(t *testing.T) {
	testlib.AwaitTillerUp(t)
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

	scaleYCSB(t, namespaceName)

	ycsbPodName := testlib.GetPodName(t, namespaceName, YCSB_CONTROLLER_NAME)

	// let YCSB run for a couple of seconds
	time.Sleep(5 * time.Second)

	testlib.GetAppLog(t, namespaceName, ycsbPodName, "-ycsb")
}
