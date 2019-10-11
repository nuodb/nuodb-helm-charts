package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/consulting-helm/test/testlib"

	"gotest.tools/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func verifyLoadBalancer(t *testing.T, namespaceName string, kubectlOptions *k8s.KubectlOptions, balancerName string) {
	balancerService := k8s.GetService(t, kubectlOptions, balancerName)
	assert.Equal(t, balancerService.Name, balancerName)
}

func verifyLBPolicy(t *testing.T, namespaceName string, kubectlOptions *k8s.KubectlOptions, podName string) {
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

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	options := &helm.Options{
		SetValues: map[string]string{},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testadminsinglereplica-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 1)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, kubectlOptions, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, kubectlOptions, admin0) })
	t.Run("verifyPodKill", func(t *testing.T) { verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyAdminService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0) })
}
func TestKubernetesBasicAdminThreeReplicas(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-1", helmChartReleaseName)
	admin2 := fmt.Sprintf("%s-nuodb-2", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	options := &helm.Options{
		SetValues: map[string]string{"admin.replicas": "3"},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testadminthreereplicas-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 3)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin1, 100*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin2, 100*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)
	defer testlib.GetAppLog(t, namespaceName, admin1)
	defer testlib.GetAppLog(t, namespaceName, admin2)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, kubectlOptions, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, kubectlOptions, admin0) })
	t.Run("verifyPodKill", func(t *testing.T) { verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 3) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 3) })
	t.Run("verifyAdminService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0) })
}

func TestKubernetesUpgradeAdmin(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	options := &helm.Options{
		SetValues: map[string]string{"nuodb.image.tag": "4.0"},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testadminupgradeadmin-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 1)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)
	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")

	upgradedOptions := &helm.Options{
		SetValues: map[string]string{"nuodb.image.tag": "4.0.1"},
	}

	helm.Upgrade(t, upgradedOptions, helmChartPath, helmChartReleaseName)

	testlib.AwaitAdminPodUpgraded(t, namespaceName, admin0, "docker.io/nuodb/nuodb-ce:4.0.1", 300*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)
	defer testlib.GetAppLog(t, namespaceName, admin0)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, kubectlOptions, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, kubectlOptions, admin0) })
}

func TestKubernetesInvalidLicense(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	licenseString := "red-riding-hood"
	customFile := "customFile"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.configFiles.nuodb\\.lic":                 licenseString,
			fmt.Sprintf("admin.configFiles.%s", customFile): "TestKubernetesInvalidLicense"},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testadmininvalidlicense-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 1)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)

	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)

		// the license provided is not a valid PEM file
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, true)
	})
	t.Run("verifyLicenseFile", func(t *testing.T) {
		testlib.VerifyLicenseFile(t, namespaceName, admin0, licenseString)
	})

	t.Run("verifyCustomFileDoesNotGetMounted", func(t *testing.T) {
		testlib.VerifyCustomFileDoesNotGetMounted(t, namespaceName, admin0, customFile)
	})

}

func TestKubernetesBasicNameOverride(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
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

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, nonDefaultName, 1)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}
