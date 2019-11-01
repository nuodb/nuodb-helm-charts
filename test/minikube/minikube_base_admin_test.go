// +build short

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"gotest.tools/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
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

func TestKubernetesBasicAdminSingleReplica(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
	t.Run("verifyPodKill", func(t *testing.T) { verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyAdminService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0) })
}

func TestKubernetesInvalidLicense(t *testing.T) {
	testlib.AwaitTillerUp(t)

	licenseString := "red-riding-hood"
	customFile := "customFile"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.configFiles.nuodb\\.lic":                 licenseString,
			fmt.Sprintf("admin.configFiles.%s", customFile): "TestKubernetesInvalidLicense"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)

		// the license provided is not a valid PEM file
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, true)
	})
	t.Run("verifyLicenseFile", func(t *testing.T) {
		testlib.VerifyLicenseFile(t, namespaceName, admin0, licenseString)
	})

}

func TestKubernetesBasicNameOverride(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	nonDefaultName := "nondefault-adminname"
	admin0 := fmt.Sprintf("%s-0", nonDefaultName)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.nameOverride":     "aws-a",
			"admin.fullnameOverride": nonDefaultName,
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testadminsinglereplica-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	helm.Install(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, nonDefaultName, 1)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}
