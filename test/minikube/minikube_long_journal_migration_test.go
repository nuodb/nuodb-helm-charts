// +build diagnostics
// not the ideal plan, but diagnostics is currently 1/2 the time of other plans

package minikube

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestChangingJournalLocationFails(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		options := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "false",
			},
		}

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

		options.SetValues["database.sm.hotCopy.journalPath.enabled"] = "true"

		err := helm.UpgradeE(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)
		require.Error(t, err)
	})
}


func TestChangingJournalLocationWithMultipleSMs(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test autoRestore without the Enterprise Edition")
	}

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	testlib.ApplyNuoDBLicense(t, namespaceName, admin0)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		options := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.noHotCopy.replicas":        "1",
				"database.sm.noHotCopy.journalPath.enabled": "false",
				"database.sm.hotCopy.journalPath.enabled": "false",
			},
		}

		// Stage 1: start a database with 2 SMs and journalPath false

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

		statefulSets := testlib.FindAllStatefulSets(t, namespaceName)

		// Stage 2: Delete non-HC SM, upgrade to journalPath and restart

		testlib.ScaleStatefulSet(t, namespaceName, statefulSets.SmNonHCSet.Name, 0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		nonHCPVC := fmt.Sprintf("archive-volume-sm-%s-nuodb-cluster0-demo-0", databaseReleaseName)

		// trigger Kubernetes-Aware-Admin to purge this archive from the database
		testlib.DeletePVC(t, namespaceName, nonHCPVC)

		testlib.DeleteStatefulSet(t, namespaceName, statefulSets.SmNonHCSet.Name)

		options.SetValues["database.sm.noHotCopy.journalPath.enabled"] = "true"

		helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

		// Stage 3: Delete HC SM, upgrade to journalPath and restart

		testlib.ScaleStatefulSet(t, namespaceName, statefulSets.SmHCSet.Name, 0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		smHCPVC := fmt.Sprintf("archive-volume-sm-%s-nuodb-cluster0-demo-hotcopy-0", databaseReleaseName)

		// trigger Kubernetes-Aware-Admin to purge this archive from the database
		testlib.DeletePVC(t, namespaceName, smHCPVC)

		testlib.DeleteStatefulSet(t, namespaceName, statefulSets.SmHCSet.Name)

		options.SetValues["database.sm.hotCopy.journalPath.enabled"] = "true"

		helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

	})
}