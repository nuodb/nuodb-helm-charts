//go:build diagnostics
// +build diagnostics

// not the ideal plan, but diagnostics is currently 1/2 the time of other plans

package minikube

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
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
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
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

	testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		options := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":        "250m",
				"database.sm.resources.requests.memory":     testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":        "250m",
				"database.te.resources.requests.memory":     testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.noHotCopy.replicas":            "1",
				"database.sm.noHotCopy.journalPath.enabled": "false",
				"database.sm.hotCopy.journalPath.enabled":   "false",
			},
		}

		// Stage 1: start a database with 2 SMs and journalPath false

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

		opt := testlib.GetExtractedOptions(&options)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s-", databaseReleaseName, opt.ClusterName, opt.DbName)
		smPodNameTemplateHC := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy", databaseReleaseName, opt.ClusterName, opt.DbName)
		smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
		smPodNameHC0 := testlib.GetPodName(t, namespaceName, smPodNameTemplateHC)

		go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodNameHC0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

		statefulSets := testlib.FindAllStatefulSets(t, namespaceName)

		// Stage 2: Delete non-HC SM, upgrade to journalPath and restart

		testlib.ScaleStatefulSet(t, namespaceName, statefulSets.SmNonHCSet.Name, 0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		nonHCPvcName := fmt.Sprintf("archive-volume-sm-%s-nuodb-cluster0-demo-0", databaseReleaseName)
		nonHCPvc := testlib.GetPvc(t, namespaceName, nonHCPvcName)
		smHCPvcName := fmt.Sprintf("archive-volume-sm-%s-nuodb-cluster0-demo-hotcopy-0", databaseReleaseName)
		smHCPvc := testlib.GetPvc(t, namespaceName, smHCPvcName)

		testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
			kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
			k8s.RunKubectl(t, kubectlOptions, "get", "pvc")
			k8s.RunKubectl(t, kubectlOptions, "get", "pv")
			k8s.RunKubectlE(t, kubectlOptions, "describe", "pvc", nonHCPvc.Name)
			k8s.RunKubectlE(t, kubectlOptions, "describe", "pvc", smHCPvc.Name)
		})

		testlib.DeleteStatefulSet(t, namespaceName, statefulSets.SmNonHCSet.Name)

		// trigger Kubernetes-Aware-Admin to purge this archive from the
		// database; this is needed because currently nuodocker doesn't change
		// the archive object stored in the admin layer
		testlib.DeletePVC(t, namespaceName, nonHCPvc.Name)
		// await for the old archive storage to be recycled so that the newly
		// started SM will sync it from the running one
		testlib.AwaitPvDeleted(t, nonHCPvc.Spec.VolumeName, 120*time.Second)

		options.SetValues["database.sm.noHotCopy.journalPath.enabled"] = "true"

		helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

		// Stage 3: Delete HC SM, upgrade to journalPath and restart

		testlib.ScaleStatefulSet(t, namespaceName, statefulSets.SmHCSet.Name, 0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		testlib.DeleteStatefulSet(t, namespaceName, statefulSets.SmHCSet.Name)

		// trigger Kubernetes-Aware-Admin to purge this archive from the
		// database; this is needed because currently nuodocker doesn't change
		// the archive object stored in the admin layer
		testlib.DeletePVC(t, namespaceName, smHCPvc.Name)
		// await for the old archive storage to be recycled so that the newly
		// started SM will sync it from the running one
		testlib.AwaitPvDeleted(t, smHCPvc.Spec.VolumeName, 120*time.Second)

		options.SetValues["database.sm.hotCopy.journalPath.enabled"] = "true"

		helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)

	})
}
