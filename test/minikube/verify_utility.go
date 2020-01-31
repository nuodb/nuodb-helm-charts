package minikube

import (
        "fmt"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"gotest.tools/assert"
	"testing"
	"time"
)

func verifyAdminService(t *testing.T, namespaceName string, podName string, serviceName string, ping bool) {

	adminService := testlib.GetService(t, namespaceName, serviceName)
	assert.Equal(t, adminService.Name, serviceName)

	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
}

func verifyLBPolicy(t *testing.T, namespaceName string, podName string, domainName string) {
	jobLBPolicy := fmt.Sprintf("%s-job-lb-policy-nearest", domainName)
	testlib.AwaitBalancerTerminated(t, namespaceName, jobLBPolicy)
	testlib.VerifyPolicyInstalled(t, namespaceName, podName)
}

func verifyPodKill(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillAdminPod(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitAdminPodUp(t, namespaceName, podName, 100*time.Second)
}

func verifyKillProcess(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillAdminProcess(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitAdminPodUp(t, namespaceName, podName, 100*time.Second)
}
