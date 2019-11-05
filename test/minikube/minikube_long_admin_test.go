// +build short

package minikube

import (
	"fmt"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestKubernetesBasicAdminThreeReplicas(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{
		SetValues: map[string]string{"admin.replicas": "3"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 3, "")

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