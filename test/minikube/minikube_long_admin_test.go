// +build long

package minikube

import (
	"fmt"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestKubernetesUpgradeAdmin(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{
		SetValues: map[string]string{"nuodb.image.tag": "4.0"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

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

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, "docker.io/nuodb/nuodb-ce:4.0.1", 300*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyLoadBalancer", func(t *testing.T) { verifyLoadBalancer(t, namespaceName, lbName) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
}