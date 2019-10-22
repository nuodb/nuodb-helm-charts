package minikube

import (
	"fmt"
	"path/filepath"
	"runtime"
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

func getFunctionCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	nameFull := runtime.FuncForPC(pc).Name()    // main.foo
	nameEnd := filepath.Ext(nameFull)           // .foo
	name := strings.TrimPrefix(nameEnd, ".")    // foo

	return name
}

func startAdmin(t *testing.T, options *helm.Options, replicaCount int, namespace string) (helmChartReleaseName string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH
	helmChartReleaseName = fmt.Sprintf("admin-%s", randomSuffix)

	adminNames := make([]string, replicaCount)

	for i := 0; i < replicaCount; i++ {
		adminNames[i] = fmt.Sprintf("%s-nuodb-%d", helmChartReleaseName, i)
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	if namespace == "" {
		callerName := getFunctionCallerName()
		namespaceName = fmt.Sprintf("%s-%s", strings.ToLower(callerName), randomSuffix)
		k8s.CreateNamespace(t, kubectlOptions, namespaceName)
		testlib.AddTeardown(testlib.TEARDOWN_ADMIN, func() { k8s.DeleteNamespace(t, kubectlOptions, namespaceName) })
	} else {
		namespaceName = namespace
	}

	options.KubectlOptions.Namespace = namespaceName

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	testlib.AddTeardown("admin", func() { helm.Delete(t, options, helmChartReleaseName, true) })

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, replicaCount)

	for i := 0; i < replicaCount; i++ {
		adminName := adminNames[i] // array will be out of scope for defer

		// first await could be pulling the image from the repo
		testlib.AwaitAdminPodUp(t, namespaceName, adminName, 300*time.Second)
		testlib.AddTeardown("admin", func() { testlib.GetAppLog(t, namespaceName, adminName) })
	}

	for i := 0; i < replicaCount; i++ {
		testlib.AwaitAdminFullyConnected(t, namespaceName, adminNames[i], replicaCount)
	}

	return
}

func TestKubernetesBasicAdminSingleReplica(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := startAdmin(t, &options, 1, "")

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
func TestKubernetesBasicAdminThreeReplicas(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{
		SetValues: map[string]string{"admin.replicas": "3"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := startAdmin(t, &options, 3, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
	t.Run("verifyPodKill", func(t *testing.T) { verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 3) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 3) })
	t.Run("verifyAdminService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0) })
}

func TestKubernetesUpgradeAdmin(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{
		SetValues: map[string]string{"nuodb.image.tag": "4.0"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := startAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	lbName := fmt.Sprintf("%s-nuodb-balancer", helmChartReleaseName)

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")

	upgradedOptions := &helm.Options{
		SetValues: map[string]string{"nuodb.image.tag": "4.0.1"},
	}

	helm.Upgrade(t, upgradedOptions, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	testlib.AwaitAdminPodUpgraded(t, namespaceName, admin0, "docker.io/nuodb/nuodb-ce:4.0.1", 300*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
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

	helmChartReleaseName, namespaceName := startAdmin(t, options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

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
