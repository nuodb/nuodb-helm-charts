package minikube

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"testing"
	"time"
	"gotest.tools/assert"
)

func verifyAdminService(t *testing.T, namespaceName string, podName string, serviceName string, ping bool) {

	adminService := testlib.GetService(t, namespaceName, serviceName)
	assert.Equal(t, adminService.Name, serviceName)

	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
}

func verifyLBPolicy(t *testing.T, namespaceName string, podName string) {
	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")
	testlib.VerifyPolicyInstalled(t, namespaceName, podName)
}

func verifyPodKill(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillAdminPod(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitPodUp(t, namespaceName, podName, 100*time.Second)
}

func verifyKillProcess(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillProcess(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitPodUp(t, namespaceName, podName, 100*time.Second)
}
