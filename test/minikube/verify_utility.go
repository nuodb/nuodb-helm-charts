package minikube

import (
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"testing"
	"time"
	"gotest.tools/assert"
)

func verifyLoadBalancer(t *testing.T, namespaceName string, balancerName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	balancerService := k8s.GetService(t, kubectlOptions, balancerName)
	assert.Equal(t, balancerService.Name, balancerName)
}

func verifyLBPolicy(t *testing.T, namespaceName string, podName string) {
	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")
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

func verifyAdminService(t *testing.T, namespaceName string, podName string) {
	serviceName := "nuodb"

	adminService := testlib.GetService(t, namespaceName, serviceName)
	assert.Equal(t, adminService.Name, serviceName)

	testlib.PingService(t, namespaceName, serviceName, podName)
}