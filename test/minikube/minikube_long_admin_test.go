//go:build long
// +build long

package minikube

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/core/v1"
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

func TestDatabaseAdminAffinityLabels(t *testing.T) {
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 6.0.3")
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	options := helm.Options{}
	options.SetValues = map[string]string{
		"admin.adminLabels.host": "host1",
		"admin.adminLabels.zone": "us-east-1",
	}
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions
	admin := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := fmt.Sprintf("%s-0", admin)

	testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)

	testlib.VerifyAdminLabels(t, namespaceName, admin0,
		map[string]string{
			"host": "host1",
			"zone": "us-east-1",
		})

	dbName := "db"
	options.SetValues = map[string]string{
		"admin.affinityLabels":    "host zone",
		"database.te.labels.host": "host2",
		"database.te.labels.zone": "us-east-1",
		"database.sm.labels.host": "host2",
		"database.sm.labels.zone": "us-east-1",

		"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
		"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
		"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		"database.sm.noHotCopy.replicas":        "1",
		"database.sm.hotCopy.replicas":          "1",
		"database.name":                         dbName,
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	testlib.StartDatabase(t, namespaceName, admin0, &options)

	processes, err := testlib.GetDatabaseProcessesE(t, namespaceName, admin0, dbName)
	require.NoError(t, err)
	for _, process := range processes {
		require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, process.Hostname,
			"Looking for admin with labels matching: host zone", &v12.PodLogOptions{}))

		expectedAffinityLog := fmt.Sprintf("Preferring APs %s due to matching label zone=us-east-1", admin0)
		require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, process.Hostname,
			expectedAffinityLog, &v12.PodLogOptions{}), "Did not find expected log message %s", expectedAffinityLog)

	}
}
