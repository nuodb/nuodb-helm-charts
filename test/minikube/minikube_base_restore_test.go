//go:build short
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

	"github.com/Masterminds/semver"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

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
			podName := testlib.GetPodName(t, namespaceName, "incremental-hotcopy-nuodb-demo-cluster0-0")
			testlib.AwaitPodPhase(t, namespaceName, podName, corev1.PodFailed, 20*time.Second)
			testlib.GetAppLog(t, namespaceName, podName, "", &corev1.PodLogOptions{})
		})

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		testlib.AwaitJobSucceeded(t, namespaceName, "incremental-hotcopy-nuodb-demo-cluster0-0", 120*time.Second)
		verifyBackup(t, namespaceName, admin0, "demo", &databaseOptions)
	})
}

func TestKubernetesJournalBackupSuspended(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 4.3")
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	options := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":        "250m",
			"database.sm.resources.requests.memory":     testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":        "250m",
			"database.te.resources.requests.memory":     testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.replicas":                      "0",
			"database.sm.hotCopy.replicas":              "2",
			"database.sm.hotCopy.journalBackup.enabled": "true",
		},
	}

	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	opt := testlib.GetExtractedOptions(&options)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	teDeployment := fmt.Sprintf("te-%s-%s-%s-%s", databaseReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
	smPodName0 := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy-0", databaseReleaseName, opt.ClusterName, opt.DbName)
	smPodName1 := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy-1", databaseReleaseName, opt.ClusterName, opt.DbName)

	// suspend all backup jobs
	backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
	backupGroup1 := fmt.Sprintf("%s-1", opt.ClusterName)
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup0)
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup1)

	// execute initial backup for backup group 1 which should fail as the
	// database is not initialized yet
	err := testlib.BackupDatabaseE(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Database not fully initialized by a Transaction Engine")

	// start the TE
	testlib.ScaleDeployment(t, namespaceName, teDeployment, 1)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, teDeployment, 1)
	tePodName := testlib.GetPodName(t, namespaceName, teDeployment)
	testlib.AwaitPodUp(t, namespaceName, tePodName, 120*time.Second)

	testlib.CreateQuickstartSchema(t, namespaceName, admin0)

	// Execute initial backup for backup group 0
	testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup0)

	restartSmAndExecuteJournalBackup := func(name string, backupGroup string) string {
		// restarting an SM will disable journal backup temporary until full or
		// incremental are requested and complete
		pod := testlib.GetPod(t, namespaceName, name)
		testlib.DeletePod(t, namespaceName, "pod/"+name)
		testlib.AwaitPodObjectRecreated(t, namespaceName, pod, 30*time.Second)
		testlib.AwaitPodUp(t, namespaceName, name, 120*time.Second)
		cronJobName := fmt.Sprintf("journal-hotcopy-%s-%s-%s", opt.DomainName, opt.DbName, backupGroup)
		testlib.DeleteJobPods(t, namespaceName, cronJobName)
		// trigger on-demand journal backup
		jobName := fmt.Sprintf("journal-backup-%s", strings.ToLower(random.UniqueId()))
		k8s.RunKubectl(t, kubectlOptions, "create", "job", "--from=cronjob/"+cronJobName, jobName)

		// Get logs from journal backup job in case it fails
		testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
			podName := testlib.GetPodName(t, namespaceName, jobName)
			testlib.GetAppLog(t, namespaceName, podName, "", &corev1.PodLogOptions{})
		})

		testlib.AwaitJobSucceeded(t, namespaceName, jobName, 120*time.Second)
		return testlib.GetPodName(t, namespaceName, jobName)
	}

	backupPodName := restartSmAndExecuteJournalBackup(smPodName0, backupGroup0)
	// verify that the journal backup fails and it's retried after requesting
	// incremental
	require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, backupPodName,
		"Executing incremental hot copy as journal hot copy is temporarily suspended", &corev1.PodLogOptions{}),
		"Incremental hot copy not requested to enable journal after sync")

	backupPodName = restartSmAndExecuteJournalBackup(smPodName1, backupGroup1)
	// verify that the journal backup fails and another full backup is requested
	// because the last full hot copy failed

	require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, backupPodName,
		"Executing incremental hot copy as journal hot copy is temporarily suspended", &corev1.PodLogOptions{}),
		"Incremental hot copy not requested to enable journal after sync")
	require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, backupPodName,
		"Executing full hotcopy as a prerequisite for incremental hotcopy", &corev1.PodLogOptions{}),
		"Full hot copy should be requested as previous full have failed")

	verifyBackup(t, namespaceName, admin0, "demo", &options)
}

func restoreDatabaseByArchiveType(t *testing.T, options helm.Options, namespaceName string, admin0 string, archiveType string) {
	isLsaType := archiveType == "lsa"
	name := "restoreFileArchive"
	if isLsaType {
		name = "restoreLsaArchive"
	}

	t.Run(name, func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		if isLsaType {
			options.SetValues["database.archiveType"] = "lsa"
		}

		databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &options)

		if isLsaType {
			delete(options.SetValues, "database.archiveType")
		}

		opts := k8s.NewKubectlOptions("", "", namespaceName)
		options.KubectlOptions = opts

		opt := testlib.GetExtractedOptions(&options)

		// Generate diagnose in case this test fails
		testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
			testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
			testlib.RecoverCoresFromEngine(t, namespaceName, "te",
				fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName))
		})

		tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

		// Execute initial backup
		backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
		testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup0)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		go testlib.GetAppLog(t, namespaceName, tePodName, "_same_pre-restart", &corev1.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodName0, "_same_pre-restart", &corev1.PodLogOptions{Follow: true})

		// restore database
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		testlib.RestoreDatabase(t, namespaceName, admin0, &options)

		go testlib.GetAppLog(t, namespaceName, smPodName0, "_same_post-restart", &corev1.PodLogOptions{Follow: true})

		if archiveType == "lsa" {
			require.False(t, testlib.HasFile(t, namespaceName, smPodName0, "/var/opt/nuodb/archive/nuodb/demo/1.atm"))
			require.True(t, testlib.HasDirectory(t, namespaceName, smPodName0, "/var/opt/nuodb/archive/nuodb/demo/1.cat"))
		} else {
			require.True(t, testlib.HasFile(t, namespaceName, smPodName0, "/var/opt/nuodb/archive/nuodb/demo/1.atm"))
		}

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
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
	namespaceName := fmt.Sprintf("%skubernetesrestoredatabase-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)
	// NuoDB 4.2 doesn't ship SSL certificates which will disable TLS in case
	// certificates are not generated; this is needed because NuoDB 4.0.8 image
	// will be used for the restore chart
	testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	restoreDatabaseByArchiveType(t, options, namespaceName, admin0, "")
	testlib.RunOnNuoDBVersionCondition(t, ">=6.0.0", func(version *semver.Version) {
		restoreDatabaseByArchiveType(t, options, namespaceName, admin0, "lsa")
	})
}

func TestKubernetesImportDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		defer testlib.Teardown(testlib.TEARDOWN_NGINX)

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
		remoteUrl := testlib.ServeFileViaHTTP(t, namespaceName, testlib.IMPORT_ARCHIVE_FILE)
		databaseOptions.SetValues["database.autoImport.source"] = remoteUrl
		databaseOptions.SetValues["database.autoImport.credentials"] = ""
		helm.Upgrade(t, databaseOptions, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		smPod0 := testlib.GetPod(t, namespaceName, smPodName0)
		testlib.DeletePod(t, namespaceName, "pod/"+smPodName0)
		testlib.AwaitPodObjectRecreated(t, namespaceName, smPod0, 30*time.Second)
		testlib.AwaitPodLog(t, namespaceName, smPodName0, "_no-credentials")
		testlib.AwaitPodUp(t, namespaceName, smPodName0, 120*time.Second)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema HOCKEY")
		require.NoError(t, err, "error running SQL: show schema HOCKEY")
		require.True(t, strings.Contains(tables, "PLAYERS"))
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

	testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)

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
		opt := testlib.GetExtractedOptions(&databaseOptions)
		pvcName := fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName)
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", pvcName)
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
	backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup0)

	moveArchiveData := func(podName string) {
		// Move the archive data which will cause the SM to ASSERT when an atom needs to be loaded
		k8s.RunKubectl(t, kubectlOptions, "exec", podName, "-c", "engine", "--",
			"mv", "-f", "/var/opt/nuodb/archive/nuodb/demo", "/var/opt/nuodb/archive/nuodb/demo-moved")
		testlib.RunSQL(t, namespaceName, admin0, "demo", "delete from hockey.hockey")
		// in cases where core file is generated, wait until it is dumped and processed
		testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, podName, 0, 120*time.Second)
		testlib.AwaitPodUp(t, namespaceName, podName, 90*time.Second)
		testlib.AwaitPodLog(t, namespaceName, podName, "_post-restart")
	}

	t.Run("restartHotCopySM", func(t *testing.T) {
		moveArchiveData(hcSmPodName0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		// HC SM should restore the archive from the latest backup

		expectedLine := fmt.Sprintf("Finished restoring /var/opt/nuodb/backup/%s to /var/opt/nuodb/archive/nuodb/demo", backupset)

		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0, expectedLine, &corev1.PodLogOptions{}),
			1, fmt.Sprintf("Expected line [%s] not found in logs from pod %s", expectedLine, hcSmPodName0))
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})

	t.Run("restartNonHotCopySM", func(t *testing.T) {
		moveArchiveData(smPodName0)
		// nonHC SM should remove the archive metadata and SYNC the data from other SM
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})
}

// Basic test for creating a database off of a VolumeSnapshot
func TestKubernetesSnapshotRestore(t *testing.T) {
	// TODO: Once it exists, use proper database snapshot function to freeze database, set backup id, and create volume snapshots
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	admin := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := fmt.Sprintf("%s-0", admin)

	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_ADMIN, t, func() {
		k8s.RunKubectl(t, kubectlOptions, "get", "pods", "-o", "wide")
		testlib.DescribePods(t, namespaceName, admin)
	})
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	sourceDatabaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":                          testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.persistence.storageClass": testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.enabled":                  "true",
			"database.sm.noHotCopy.replicas":                             "1",
			"database.sm.hotCopy.replicas":                               "0",
		},
	})

	output, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "CREATE TABLE testtbl (id INT); INSERT INTO testtbl (id) values (123)")

	require.NoError(t, err, output)

	smPod := fmt.Sprintf("sm-%s-nuodb-cluster0-demo-0", sourceDatabaseChartName)
	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		"echo \"123abc\" > /var/opt/nuodb/archive/nuodb/demo/backup.txt")
	require.NoError(t, err, output)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		"echo \"123abc\" > /var/opt/nuodb/journal/nuodb/demo/backup.txt")
	require.NoError(t, err, output)

	defer testlib.Teardown(testlib.TEARDOWN_SNAPSHOT)

	achiveVolumeName := "archive-volume-" + smPod
	archiveSnapshotName := "archive-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, archiveSnapshotName)

	journalVolumeName := "journal-volume-" + smPod
	journalSnapshotName := "journal-snapshot"
	testlib.SnapshotVolume(t, namespaceName, journalVolumeName, journalSnapshotName)

	helm.DeleteE(t, &helm.Options{KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName)}, sourceDatabaseChartName, true)

	restoredDb := "db-clone"
	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                                   testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                                testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                                   testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                                testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.noHotCopy.journalPath.persistence.storageClass":           testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.persistence.storageClass":                                    testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.name":                                                        restoredDb,
			"database.persistence.dataSourceRef.kind":                              "VolumeSnapshot",
			"database.persistence.dataSourceRef.name":                              archiveSnapshotName,
			"database.persistence.dataSourceRef.apiGroup":                          "snapshot.storage.k8s.io",
			"database.sm.noHotCopy.journalPath.persistence.dataSourceRef.kind":     "VolumeSnapshot",
			"database.sm.noHotCopy.journalPath.persistence.dataSourceRef.name":     journalSnapshotName,
			"database.sm.noHotCopy.journalPath.persistence.dataSourceRef.apiGroup": "snapshot.storage.k8s.io",
			"database.autoImport.backupId":                                         "123abc",
			"database.sm.noHotCopy.journalPath.enabled":                            "true",
			"database.sm.noHotCopy.replicas":                                       "1",
			"database.sm.hotCopy.replicas":                                         "0",
		},
	})

	output, err = testlib.RunSQL(t, namespaceName, admin0, restoredDb, "SELECT id FROM testtbl")
	require.NoError(t, err, output)

	require.True(t, strings.Contains(output, "123"))
}
