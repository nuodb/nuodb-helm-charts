//go:build long
// +build long

package minikube

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestKubernetesBasicAdminThreeReplicas(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{"admin.replicas": "3"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 3, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-cluster0-1", helmChartReleaseName)
	admin2 := fmt.Sprintf("%s-nuodb-cluster0-2", helmChartReleaseName)
	clusterServiceName := fmt.Sprintf("nuodb-clusterip")

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		if os.Getenv("NUODB_LICENSE") == "ENTERPRISE" {
			t.Skip("Cannot test licensing in Enterprise Edition")
		}
		if os.Getenv("NUODB_LIMITED_LICENSE_CONTENT") == "" {
			t.Skip("Cannot test licensing without Limited License")
		}
		testlib.VerifyLicense(t, namespaceName, admin0, testlib.LIMITED)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyAdminClusterService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0, clusterServiceName, false) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })

	t.Run("verifyProcessKill", func(t *testing.T) {
		verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 3)
		verifyKillProcess(t, namespaceName, admin1, helmChartReleaseName, 3)
		verifyKillProcess(t, namespaceName, admin2, helmChartReleaseName, 3)
	})

	t.Run("verifyPodKill", func(t *testing.T) {
		t.Skip("verifyPodKill is flaky")
		verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 3)
		verifyPodKill(t, namespaceName, admin1, helmChartReleaseName, 3)
		verifyPodKill(t, namespaceName, admin2, helmChartReleaseName, 3)
	})
}
