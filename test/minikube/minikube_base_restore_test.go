// +build short

package minikube

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func verifyPacketFetch(t *testing.T, namespaceName string, admin0 string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	// verify the container can actually download the file from the internet
	start := time.Now()
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions,
		"exec", admin0, "--",
		"bash", "-c",
		fmt.Sprintf("curl -k %s | tar tzf - ", testlib.IMPORT_ARCHIVE_URL),
	)
	require.NoError(t, err, "Could not fetch archive")
	elapsed := time.Since(start)
	t.Logf("Fetching package (%s) took %f seconds", testlib.IMPORT_ARCHIVE_URL, elapsed.Seconds())
	t.Log("tar contents: ", output)
}

func verifyBackup(t *testing.T, namespaceName string, podName string, databaseName string, options *helm.Options) {
	// verify that the backup has been documented by the Admin layer
	backupset, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions,
		"exec", podName, "--",
		"nuodocker", "get", "current-backup", "--db-name", databaseName,
	)

	require.NoError(t, err, "Error running: nuodocker get current-backup  ")
	require.True(t, backupset != "")
}

func TestKubernetesBackupDatabase(t *testing.T) {
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
				// Configure more frequent incremental schedule so that
				// a full backup is created as a prerequisite.
				"database.sm.hotCopy.incrementalSchedule": "?/1 * * * *",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		// Generate diagnose in case this test fails
		testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
			podName := testlib.GetPodName(t, namespaceName, "incremental-hotcopy-demo-cronjob")
			testlib.AwaitPodPhase(t, namespaceName, podName, corev1.PodFailed, 20*time.Second)
			testlib.GetAppLog(t, namespaceName, podName, "", &corev1.PodLogOptions{})
		})

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		testlib.AwaitJobSucceeded(t, namespaceName, "incremental-hotcopy-demo-cronjob", 120*time.Second)
		verifyBackup(t, namespaceName, admin0, "demo", &databaseOptions)
	})
}

func TestKubernetesRestoreDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"backup.persistence.enabled":            "true",
			"backup.persistence.size":               "1Gi",
			"database.te.logPersistence.enabled":    "true",
		},
	}

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("kubernetesrestoredatabase-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)
	// NuoDB 4.2 doesn't ship SSL certificates which will disable TLS in case
	// certificates are not generated; this is needed because NuoDB 4.0.8 image
	// will be used for the restore chart
	testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", "demo-log-te-volume")
	})

	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	opt := testlib.GetExtractedOptions(&options)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

	// Execute initial backup
	testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)

	t.Run("restoreDatabaseSameVersion", func(t *testing.T) {
		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		go testlib.GetAppLog(t, namespaceName, tePodName, "_same_pre-restart", &corev1.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodName0, "_same_pre-restart", &corev1.PodLogOptions{Follow: true})

		// restore database
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		testlib.RestoreDatabase(t, namespaceName, admin0, &options)

		go testlib.GetAppLog(t, namespaceName, smPodName0, "_same_post-restart", &corev1.PodLogOptions{Follow: true})

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})

	t.Run("restoreDatabaseBackwardsCompatibility", func(t *testing.T) {
		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		go testlib.GetAppLog(t, namespaceName, tePodName, "_compatible_pre-restart", &corev1.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodName0, "_compatible_pre-restart", &corev1.PodLogOptions{Follow: true})

		// restore database using pre 4.2 version of NuoDB image for the restore
		// chart; this should place "legacy" restore request
		options.SetValues["nuodb.image.registry"] = "docker.io"
		options.SetValues["nuodb.image.repository"] = "nuodb/nuodb-ce"
		options.SetValues["nuodb.image.tag"] = "4.0.8"
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		testlib.RestoreDatabase(t, namespaceName, admin0, &options)

		go testlib.GetAppLog(t, namespaceName, smPodName0, "_compatible_post-restart", &corev1.PodLogOptions{Follow: true})

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})
}

func TestKubernetesImportDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	verifyPacketFetch(t, namespaceName, admin0)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.autoImport.source":            "sftp://sftp.nuodb.com/incoming/restore.bak.tz",
				"database.autoImport.credentials":       "nuodb:wrongPass",
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		}

		// Install database and expect archive download failure during auto import
		databaseReleaseName := testlib.StartDatabaseNoWait(t, namespaceName, admin0, databaseOptions)

		opt := testlib.GetExtractedOptions(databaseOptions)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
		testlib.AwaitNrReplicasScheduled(t, namespaceName, smPodNameTemplate, opt.NrSmPods)
		smPodName0 := fmt.Sprintf("%s-hotcopy-0", smPodNameTemplate)
		testlib.AwaitPodLog(t, namespaceName, smPodName0, "_invalid-credentials")
		testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, smPodName0, 0, 120*time.Second)
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
			"Restore: unable to download/unpack backup", &corev1.PodLogOptions{Previous: true}), 1)

		// Use the correct IMPORT URL without credentials
		databaseOptions.SetValues["database.autoImport.source"] = testlib.IMPORT_ARCHIVE_URL
		databaseOptions.SetValues["database.autoImport.credentials"] = ""
		helm.Upgrade(t, databaseOptions, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		smPod0 := testlib.GetPod(t, namespaceName, smPodName0)
		testlib.DeletePod(t, namespaceName, "pod/"+smPodName0)
		testlib.AwaitPodObjectRecreated(t, namespaceName, smPod0, 30*time.Second)
		testlib.AwaitPodUp(t, namespaceName, smPodName0, 120*time.Second)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "HOCKEY"))
	})
}

func TestKubernetesAutoRestore(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test autoRestore without the Enterprise Edition")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	testlib.ApplyNuoDBLicense(t, namespaceName, admin0)

	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    "250m",
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "250m",
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.autoRestore.source":           ":latest",
			"database.sm.noHotCopy.replicas":        "1",
			"database.te.logPersistence.enabled":    "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", "demo-log-te-volume")
	})

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodNameTemplate := fmt.Sprintf("%s-hotcopy", smPodNameTemplate)
	smPodName0 := fmt.Sprintf("%s-0", smPodNameTemplate)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)
	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)

	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	testlib.CreateQuickstartSchema(t, namespaceName, admin0)
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)

	moveArchiveData := func(podName string) {
		// Move the archive data which will cause the SM to ASSERT when an atom needs to be loaded
		k8s.RunKubectl(t, kubectlOptions, "exec", podName, "--",
			"mv", "-f", "/var/opt/nuodb/archive/nuodb/demo", "/var/opt/nuodb/archive/nuodb/demo-moved")
		testlib.RunSQL(t, namespaceName, admin0, "demo", "select * from system.nodes")
		testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, podName, 0, 30*time.Second)
		testlib.AwaitPodLog(t, namespaceName, podName, "_post-restart")
	}

	t.Run("restartHotCopySM", func(t *testing.T) {
		moveArchiveData(hcSmPodName0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		// HC SM should restore the archive from the latest backup
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			fmt.Sprintf("Finished restoring /var/opt/nuodb/backup/%s to /var/opt/nuodb/archive/nuodb/demo", backupset),
			&corev1.PodLogOptions{}), 1)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})

	t.Run("restartNonHotCopySM", func(t *testing.T) {
		moveArchiveData(smPodName0)
		// nonHC SM should remove the archive metadata and SYNC the data from other SM
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})
}
