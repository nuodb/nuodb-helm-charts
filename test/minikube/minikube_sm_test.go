// +build long

/*
 * Tests to verify operation of SM startup.
 * Separate so that changes to SM startup in general - and the SM startup scipt(s) in particular
 * can be tested in isolation.
 */

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestKubernetesStartSM(t *testing.T) {
	// skip this test until the code to detect and recover from database startup failure is robust
	t.Skip()

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
				"database.env[0].name":                  "NUODB_DEBUG",
				"database.env[0].value":                 "debug",
			},
		}

		databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		fmt.Println("startDatabase returned...")

		opt := testlib.GetExtractedOptions(&databaseOptions)
		smPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
		sm0 := testlib.GetPodName(t, namespaceName, smPodTemplate)

		fmt.Println("extracted databaseOptions; sm name is ", sm0)

		kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

		testlib.AddDiagnosticTeardown(
			testlib.TEARDOWN_DATABASE,
			t,
			func() { testlib.GetFile(t, namespaceName, sm0, "/var/log/nuodb", "nuosm.log") },
		)

		defer func() {
			status, _ := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "show", "database", "--db-name", opt.DbName)
			fmt.Println("database status", status)
		}()

		t.Run("verifyRestore", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
		t.Run("verifyNonexstentRestore", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", "urgle-furgle") })

		// have to wait for the SM to be running again...
		t.Run("fixArchive1", func(t *testing.T) { healDatabase(t, namespaceName, sm0, admin0, &databaseOptions) })
		t.Run("verifynonexistentRestoreURL", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", "urgle://furgle/gurgle") })
		t.Run("fixArchive2", func(t *testing.T) { healDatabase(t, namespaceName, sm0, admin0, &databaseOptions) })
		// t.Run("verifyEmptyRestore", func(t *testing.T) {
		// 	verifyRestoreFails(t, namespaceName, sm0, "archive", testlib.RESTORE_EMPTYARCHIVE_URL)
		// })
		// t.Run("fixArchive3", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
		// t.Run("verifyWrongRestore", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", testlib.RESTORE_ARCHIVE2_URL) })
		// t.Run("fixArchive4", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
	})

	// t.Run("startDatabaseDaemonSet", func(t *testing.T) {
	// 	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	// 	databaseOptions := helm.Options{
	// 		SetValues: map[string]string{
	// 			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
	// 			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
	// 			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
	// 			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
	// 			"backup.persistence.enabled":            "true",
	// 			"backup.persistence.size":               "1Gi",
	// 			"database.enableDaemonSet":              "true",
	// 			// prevent non-backup SM from scheduling
	// 			"database.sm.nodeSelectorNoHotCopyDS.nonexistantTag": "required",
	// 		},
	// 	}

	// 	testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// 	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// 	opt := testlib.GetExtractedOptions(&databaseOptions)
	// 	smPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	// 	sm0 := testlib.GetPodName(t, namespaceName, smPodTemplate)

	// 	t.Run("verifyRestore", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
	// 	t.Run("verifyNonexstentRestore", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", "urgle-furgle") })
	// 	t.Run("fixArchive1", func(t *testing.T) { healDatabase(t, namespaceName, sm0, admin0, &databaseOptions) })
	// 	t.Run("verifynonexistentRestoreURL", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", "urgle://furgle/gurgle") })
	// 	t.Run("fixArchive2", func(t *testing.T) { healDatabase(t, namespaceName, sm0, admin0, &databaseOptions) })
	// 	// t.Run("verifyEmptyRestore", func(t *testing.T) {
	// 	// 	verifyRestoreFails(t, namespaceName, sm0, "archive", testlib.RESTORE_EMPTYARCHIVE_URL)
	// 	// })
	// 	// t.Run("fixArchive3", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
	// 	// t.Run("verifyWrongRestore", func(t *testing.T) { verifyRestoreFails(t, namespaceName, sm0, "archive", testlib.RESTORE_ARCHIVE2_URL) })
	// 	// t.Run("fixArchive4", func(t *testing.T) { verifyRestorePasses(t, namespaceName, sm0, "archive", ":latest") })
	// })

}

func TestKubernetesRestartSM(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
				"database.autoRestore.source":           ":latest",
				"database.env[0].name":                  "NUODB_DEBUG",
				"database.env[0].value":                 "debug",
			},
		}

		databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		fmt.Println("startDatabase returned...")

		opt := testlib.GetExtractedOptions(&databaseOptions)
		smPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
		sm0 := testlib.GetPodName(t, namespaceName, smPodTemplate)

		kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

		testlib.AddDiagnosticTeardown(
			testlib.TEARDOWN_DATABASE,
			t,
			func() { testlib.GetFile(t, namespaceName, sm0, "/var/log/nuodb", "nuosm.log") },
		)

		restart := func() {
			k8s.RunKubectl(t, kubectlOptions, "exec", sm0, "--", "rm", "-rf", "/var/opt/nuodb/archive/nuodb/demo")

			result, _ := testlib.RunSQL(t, namespaceName, admin0, "demo", "create table user.newtable (id bigint)")
			t.Log(result)
		}

		defer func() {
			status, _ := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "show", "database", "--db-name", opt.DbName)
			fmt.Println("database status on EXIT", status)
		}()

		t.Run("verifyRestartOnLostArchive", func(t *testing.T) {
			testlib.AwaitDatabaseRestart(t, namespaceName, admin0, "demo", &databaseOptions, restart)

			status, _ := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", sm0, "--", "ls", "-l", "/var/opt/nuodb/archive/nuodb/demo")
			t.Log("after archive delete: ", status)
		})
	})

	/*
		t.Run("startDatabaseDaemonSet", func(t *testing.T) {
			defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
			databaseOptions := helm.Options{
				SetValues: map[string]string{
					"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"backup.persistence.enabled":            "true",
					"backup.persistence.size":               "1Gi",
					"database.enableDaemonSet":              "true",
					// prevent non-backup SM from scheduling
					"database.sm.nodeSelectorNoHotCopyDS.nonexistantTag": "required",
					"database.autoRestore.source":                        ":latest",
				},
			}

			testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

			databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

			opt := testlib.GetExtractedOptions(&databaseOptions)
			smPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
			sm0 := testlib.GetPodName(t, namespaceName, smPodTemplate)

			restart := func() {
				kubectlOptions := k8s.NewKubectlOptions("", "")
				kubectlOptions.Namespace = namespaceName

				k8s.RunKubectl(t, kubectlOptions, "exec", sm0, "--", "rm", "-rf", "/var/opt/nuodb/archive/nuodb/demo")
				testlib.RunSQL(t, namespaceName, admin0, "demo", "select 1 from dual")
			}

			t.Run("verifyRestartOnLostArchive", func(t *testing.T) {
				testlib.AwaitDatabaseRestart(t, namespaceName, admin0, "demo", &databaseOptions, restart)
			})
		})
	*/

}

func healDatabase(t *testing.T, namespaceName string, podName string, adminName string, databaseOptions *helm.Options) {
	opts := testlib.GetExtractedOptions(databaseOptions)

	restoreArchive(t, namespaceName, podName, "archive", ":latest")
	testlib.AwaitDatabaseUp(t, namespaceName, adminName, opts.DbName, opts.NrTePods+opts.NrSmPods)
}

func verifyRestoreFails(t *testing.T, namespaceName string, podName string, backupType string, backupName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	restoreArchive(t, namespaceName, podName, backupType, backupName)
	testlib.VerifyProcessRestartFails(t, namespaceName, podName, func() { k8s.RunKubectl(t, kubectlOptions, "exec", podName, "--", "kill", "1") })
}

func verifyRestorePasses(t *testing.T, namespaceName string, podName string, backupType string, backupName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	restoreArchive(t, namespaceName, podName, backupType, backupName)
	testlib.AwaitProcessRestart(t, namespaceName, podName, func() { k8s.RunKubectl(t, kubectlOptions, "exec", podName, "--", "kill", "1") })
}

func restoreArchive(t *testing.T, namespaceName string, podName string, backupType string, backupName string) {
	// run the restore chart - and then manually restart the SM
	randomSuffix := strings.ToLower(random.UniqueId())

	restoreName := fmt.Sprintf("restore-demo-%s", randomSuffix)
	options := &helm.Options{
		SetValues: map[string]string{
			"database.name":       "demo",
			"restore.target":      "demo",
			"restore.source":      backupName,
			"restore.type":        backupType,
			"restore.autoRestart": "false",
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	testlib.InjectTestVersion(t, options)

	helm.Install(t, options, testlib.RESTORE_HELM_CHART_PATH, restoreName)
	defer helm.Delete(t, options, restoreName, true)
	defer k8s.RunKubectl(t, kubectlOptions, "delete", "job", "restore-demo")

	testlib.AwaitPodPhase(t, namespaceName, "restore-demo-", corev1.PodSucceeded, 120*time.Second)
	restorePod := testlib.GetPodName(t, namespaceName, "restore-demo")
	req, _ := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "--", "nuocmd", "get", "value", "--key", "$NUODB_RESTORE_REQUEST_PREFIX/demo/restore")
	t.Log("restore pod=", restorePod, "; restore request=", req)

	testlib.GetFile(t, namespaceName, "restore-demo", "/var/log/nuodb", "nuorestore.log")
}
