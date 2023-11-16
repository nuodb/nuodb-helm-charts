//go:build long
// +build long

package minikube

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func verifyExternalJournal(t *testing.T, namespaceName string, adminPod string,
	databaseReleaseName string, databaseOptions *helm.Options) {
	opt := testlib.GetExtractedOptions(databaseOptions)
	if _, ok := databaseOptions.SetValues["database.autoImport.source"]; ok {
		// verify that the journal content is moved to the external journal dir
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
		smPodName0 := fmt.Sprintf("%s-hotcopy-0", smPodNameTemplate)
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
			fmt.Sprintf("Moving restored journal content to /var/opt/nuodb/journal/nuodb/%s", opt.DbName),
			&corev1.PodLogOptions{}), 1)
		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, adminPod, "demo", "show schema HOCKEY")
		require.NoError(t, err, "error running SQL: show schema HOCKEY")
		require.True(t, strings.Contains(tables, "PLAYERS"))
	}
	// check that archives are created with external journal directory
	archives, _ := testlib.CheckArchives(t, namespaceName, adminPod, opt.DbName, opt.NrSmPods, 0)
	for _, archive := range archives {
		require.Equal(t, fmt.Sprintf("/var/opt/nuodb/journal/nuodb/%s", opt.DbName), archive.JournalPath)
	}
}

func verifyBackupSet(t *testing.T, namespaceName string, backupSet string,
	expectedBackupGroup string, expectedIndex int, expectedPods string) {

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	pods := strings.Split(expectedPods, " ")
	// verify that the correct information is recorded in the KV store for this
	// backup set
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", pods[0], "-c", "engine", "--",
		"nuobackup", "--type", "report-backups", "--group", expectedBackupGroup)
	require.NoError(t, err, output)
	expectedLine := fmt.Sprintf("%s:%d %s %s", expectedBackupGroup, expectedIndex, backupSet, expectedPods)
	require.Contains(t, output, expectedLine)
	require.NotContains(t, output, expectedLine+" ")
	// verify that the backup set exist on all expected HCSMs
	for _, podName := range pods {
		output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "-c", "engine", "--",
			"ls", "-1", "/var/opt/nuodb/backup")
		require.Contains(t, output, backupSet)
	}
}

func TestKubernetesRestoreMultipleSMs(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
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
			"database.sm.noHotCopy.replicas":        "1",
			"database.te.logPersistence.enabled":    "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)
	opt := testlib.GetExtractedOptions(&databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		pvcName := fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName)
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", pvcName)
	})

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodNameTemplate := fmt.Sprintf("%s-hotcopy", smPodNameTemplate)
	smPodName0 := fmt.Sprintf("%s-0", smPodNameTemplate)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)

	// Execute initial backup
	backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
	backupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup0)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)

	t.Run("autoRestart", func(t *testing.T) {
		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = backupset
		databaseOptions.SetValues["restore.labels"] = "pod-name " + hcSmPodName0
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
		testlib.AwaitPodLog(t, namespaceName, smPodName0, "_auto_post-restart")
		testlib.AwaitPodLog(t, namespaceName, hcSmPodName0, "_auto_post-restart")

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})

	t.Run("manualRestart", func(t *testing.T) {
		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.autoRestart"] = "false"
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)

		// Manually scale down all SMs
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=0")
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", hcSmPodNameTemplate, "--replicas=0")
		testlib.AwaitNoPods(t, namespaceName, smPodNameTemplate)

		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=1")
		testlib.AwaitPodLog(t, namespaceName, smPodName0, "_manual_post-restart")

		if testlib.IsRestoreRequestSupported(t, namespaceName, admin0) {
			// If nonHC SM is started first it should wait for the database restore to complete
			testlib.Await(t, func() bool {
				return testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
					"Waiting for database restore to complete",
					&corev1.PodLogOptions{}) == 1
			}, 60*time.Second)
		} else {
			// If nonHC SM is started first it should fail as the restore source won't be available
			testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, smPodName0, 1, 120*time.Second)
			require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
				fmt.Sprintf("Backupset %s cannot be found in /var/opt/nuodb/backup", backupset),
				&corev1.PodLogOptions{Previous: true}), 1)
		}

		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=0")
		testlib.AwaitNoPods(t, namespaceName, smPodNameTemplate)

		// Restart SM pods in the correct order so that HC SM performs the restore
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", hcSmPodNameTemplate, "--replicas=1")
		testlib.AwaitPodLog(t, namespaceName, hcSmPodName0, "_manual_post-restart")
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmHotCopyPods+opt.NrTePods)
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=1")
		testlib.AwaitPodLog(t, namespaceName, smPodName0, "_manual_post-restore")
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})
}

func TestKubernetesRestoreMultipleBackupGroups(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
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
			"database.sm.hotCopy.replicas":          "2",
			"database.sm.noHotCopy.replicas":        "1",
			"database.te.logPersistence.enabled":    "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
			// multiple restore operations with autoRestart=true may cause
			// containers to be reported as "CrashLoopBackOff" although the
			// engines will exit with zero return code
			"restore.autoRestart": "false",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)
	opt := testlib.GetExtractedOptions(&databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		pvcName := fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName)
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", pvcName)
	})

	hcSmPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)
	hcSmPodName1 := fmt.Sprintf("%s-1", hcSmPodNameTemplate)
	backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
	backupGroup1 := fmt.Sprintf("%s-1", opt.ClusterName)

	// Suspend backup jobs for all backup groups
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup0)
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup1)

	t.Run("restoreToLatest", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		// Execute backup for backup group 1
		backupsetGroup1 := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup1)
		verifyBackupSet(t, namespaceName, backupsetGroup1, backupGroup1, 1, hcSmPodName1)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = ":latest"
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
		testlib.RestartDatabasePods(t, namespaceName, databaseChartName, &databaseOptions)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

		// HCSM with ordinal 0 should not be selected for restore
		require.Equal(t, 0, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			"Restoring ", &corev1.PodLogOptions{}))
		// verify that the correct backupset is used to restore the archive of
		// HCSM with ordinal 1
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName1,
			fmt.Sprintf("Restoring %s", backupsetGroup1), &corev1.PodLogOptions{}), 1)

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 3, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})

	t.Run("restoreToGroupLatest", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		// Execute backup for backup group 0
		backupsetGroup0 := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup0)
		verifyBackupSet(t, namespaceName, backupsetGroup0, backupGroup0, 1, hcSmPodName0)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = "cluster0-0:latest"
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
		testlib.RestartDatabasePods(t, namespaceName, databaseChartName, &databaseOptions)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

		// HCSM with ordinal 1 should not be selected for restore
		require.Equal(t, 0, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName1,
			"Restoring ", &corev1.PodLogOptions{}))
		// verify that the correct backupset is used to restore the archive of
		// HCSM with ordinal 0
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			fmt.Sprintf("Restoring %s", backupsetGroup0), &corev1.PodLogOptions{}), 1)
		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 3, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})

	t.Run("restoreToGroupSpecific", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		// Create another backup for backup group 0 (index 2)
		newBackupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup0)
		verifyBackupSet(t, namespaceName, newBackupset, backupGroup0, 2, hcSmPodName0)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = "cluster0-0:2"
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
		testlib.RestartDatabasePods(t, namespaceName, databaseChartName, &databaseOptions)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

		// HCSM with ordinal 1 should not be selected for restore
		require.Equal(t, 0, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName1,
			"Restoring ", &corev1.PodLogOptions{}))
		// verify that the correct backupset is used to restore the archive of
		// HCSM with ordinal 0
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			fmt.Sprintf("Restoring %s", newBackupset), &corev1.PodLogOptions{}), 1)
		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 3, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})
}

func TestKubernetesRestoreCustomBackupGroups(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)

	processFilterTemplate := "and(label(backup cluster0) label(pod-name *-hotcopy-%d))"
	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.name":                                         "demo",
			"database.sm.resources.requests.cpu":                    "250m",
			"database.sm.resources.requests.memory":                 testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                    "250m",
			"database.te.resources.requests.memory":                 testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.hotCopy.replicas":                          "2",
			"database.te.logPersistence.enabled":                    "true",
			"database.env[0].name":                                  "NUODB_DEBUG",
			"database.env[0].value":                                 "debug",
			"database.sm.hotCopy.backupGroups.group0.processFilter": fmt.Sprintf(processFilterTemplate, 0),
		},
	}

	opt := testlib.GetExtractedOptions(&databaseOptions)

	restoreOptions := helm.Options{
		SetValues: map[string]string{
			"restore.target": "demo",
			// multiple restore operations with autoRestart=true may cause
			// containers to be reported as "CrashLoopBackOff" although the
			// engines will exit with zero return code
			"restore.autoRestart": "false",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	hcSmPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)
	hcSmPodName1 := fmt.Sprintf("%s-1", hcSmPodNameTemplate)

	// Add another backup group using labels and upgrade the Helm release
	labelsTemplate := "pod-name %s-%d"
	databaseOptions.SetValues["database.sm.hotCopy.backupGroups.group1.labels"] = fmt.Sprintf(labelsTemplate, hcSmPodNameTemplate, 1)
	helm.Upgrade(t, &databaseOptions, testlib.DATABASE_HELM_CHART_PATH, databaseChartName)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		pvcName := fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName)
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", pvcName)
	})

	backupGroup0 := "group0"
	backupGroup1 := "group1"

	// Suspend backup jobs for all backup groups
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup0)
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, backupGroup1)

	t.Run("restoreWithProcessFilter", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		// Create another backup for backup group 0 (index 1)
		newBackupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup0)
		verifyBackupSet(t, namespaceName, newBackupset, backupGroup0, 1, hcSmPodName0)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)

		// restore database from specific backup set by selecting HCSM with index 0
		restoreOptions.SetValues["restore.source"] = newBackupset
		restoreOptions.SetValues["restore.processFilter"] = fmt.Sprintf(processFilterTemplate, 0)
		testlib.RestoreDatabase(t, namespaceName, admin0, &restoreOptions)
		testlib.RestartDatabasePods(t, namespaceName, databaseChartName, &databaseOptions)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

		// HCSM with ordinal 1 should not be selected for restore
		require.Equal(t, 0, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName1,
			"Restoring ", &corev1.PodLogOptions{}))
		// verify that the correct backupset is used to restore the archive of
		// HCSM with ordinal 0
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			fmt.Sprintf("Restoring %s", newBackupset), &corev1.PodLogOptions{}), 1)
		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})

	t.Run("restoreWithLabels", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		// Execute backup for backup group 1
		newBackupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", backupGroup1)
		verifyBackupSet(t, namespaceName, newBackupset, backupGroup1, 1, hcSmPodName1)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)

		// restore the database using labels which should take precedence over
		// restore.processFilter value
		restoreOptions.SetValues["restore.source"] = newBackupset
		restoreOptions.SetValues["restore.labels"] = fmt.Sprintf(labelsTemplate, hcSmPodNameTemplate, 1)
		testlib.RestoreDatabase(t, namespaceName, admin0, &restoreOptions)
		testlib.RestartDatabasePods(t, namespaceName, databaseChartName, &databaseOptions)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrTePods+opt.NrSmPods)

		// HCSM with ordinal 0 should not be selected for restore
		require.Equal(t, 0, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			"Restoring ", &corev1.PodLogOptions{}))
		// verify that the correct backupset is used to restore the archive of
		// HCSM with ordinal 1
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName1,
			fmt.Sprintf("Restoring %s", newBackupset), &corev1.PodLogOptions{}), 1)

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})
}

func TestKubernetesRestoreWithStorageGroups(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 5.0.3")
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
			"database.te.logPersistence.enabled":    "true",
			"database.sm.storageGroup.enabled":      "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
			// include both HCSMs as each of them is serving separate storage group
			"database.sm.hotCopy.backupGroups.all-sg.labels": "role hotcopy",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
	})

	// install database primary release which serves sg0
	databaseOptions.SetValues["database.sm.storageGroup.name"] = "sg0"
	sg0ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// install database secondary release which serves sg1
	databaseOptions.SetValues["database.primaryRelease"] = "false"
	databaseOptions.SetValues["database.te.enablePod"] = "false"
	databaseOptions.SetValues["database.sm.storageGroup.name"] = "sg1"
	sg1ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// set the total number of engines as in the test utilities they are
	// inferred from the values
	databaseOptions.SetValues["database.te.enablePod"] = "true"
	databaseOptions.SetValues["database.te.replicas"] = "1"
	databaseOptions.SetValues["database.sm.hotCopy.replicas"] = "2"

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", sg0ReleaseName, opt.ClusterName, opt.DbName)
	sg0HcSmPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", sg0ReleaseName, opt.ClusterName, opt.DbName)
	sg1HcSmPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", sg1ReleaseName, opt.ClusterName, opt.DbName)
	sg0HcSmPodName0 := fmt.Sprintf("%s-hotcopy-0", sg0HcSmPodTemplate)
	sg1HcSmPodName0 := fmt.Sprintf("%s-hotcopy-0", sg1HcSmPodTemplate)

	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	verifyStorageGroup := func(sgName string) {
		// wait for archives to be added to the storage group
		sg := testlib.AwaitStorageGroup(t, namespaceName, admin0, opt.DbName, sgName, 30*time.Second)
		require.Equal(t, sg.State, "Available")
		require.GreaterOrEqual(t, len(sg.ArchiveStates), 1)
		require.GreaterOrEqual(t, len(sg.ProcessStates), 1)
		require.GreaterOrEqual(t, len(sg.LeaderCandidates), 1)
	}
	verifyStorageGroup("sg0")
	verifyStorageGroup("sg1")

	// Populate some data partitioned and stored in each storage group
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"CREATE TABLE codes ( sg char(4), zip char(5) ) PARTITION BY RANGE (zip) ("+
			" PARTITION p_sg0 VALUES LESS THAN ('2000') STORE IN sg0"+
			" PARTITION p_sg1 VALUES LESS THAN ('3000') STORE IN sg1)")
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"CREATE TABLE codes ( sg char(4), zip char(5) ) PARTITION BY RANGE (zip) ("+
			" PARTITION p_sg0 VALUES LESS THAN ('2000') STORE IN sg0"+
			" PARTITION p_sg1 VALUES LESS THAN ('3000') STORE IN sg1)")
	// Insert 1 row of data in each storage group
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg0', '1001')")
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg1', '2001')")

	// Suspend the backup jobs and perform a backup
	testlib.SuspendDatabaseBackupJobs(t, namespaceName, opt.DomainName, opt.DbName, "all-sg")
	testlib.BackupDatabase(t, namespaceName, sg0HcSmPodName0, opt.DbName, "full", "all-sg")

	// Insert more rows
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg0', '1001')")
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg1', '2001')")

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, sg0HcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, sg1HcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)

	// Get SM pod logs if the test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetAppLog(t, namespaceName, sg0HcSmPodName0, "_post-restore", &corev1.PodLogOptions{})
		testlib.GetAppLog(t, namespaceName, sg1HcSmPodName0, "_post-restore", &corev1.PodLogOptions{})
	})

	// restore database to the latest backup
	testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
	testlib.AwaitPodLog(t, namespaceName, sg0HcSmPodName0, "_post-restart")
	testlib.AwaitPodLog(t, namespaceName, sg1HcSmPodName0, "_post-restart")

	// verify that the database does NOT contain the data from AFTER the backup
	output, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "select sg, count(*) from codes group by sg")
	require.NoError(t, err, "error running SQL: select sg, count(*) from codes group by sg")
	require.True(t, regexp.MustCompile(`sg0\s+1`).MatchString(output), "Unexpected data in sg0: ", output)
	require.True(t, regexp.MustCompile(`sg1\s+1`).MatchString(output), "Unexpected data in sg1: ", output)
	testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
}

func TestKubernetesImportWithStorageGroups(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 5.0.3")
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
			"database.te.logPersistence.enabled":    "true",
			"database.te.replicas":                  "0",
			"database.sm.storageGroup.enabled":      "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	defer testlib.Teardown(testlib.TEARDOWN_NGINX)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
	})

	// Upload backups taken from TestKubernetesRestoreWithStorageGroups test
	// database on HTTP server
	sg0BackupUrl := testlib.ServeFileViaHTTP(t, namespaceName, testlib.IMPORT_SG0_BACKUP_FILE)
	sg1BackupUrl := testlib.ServeFileViaHTTP(t, namespaceName, testlib.IMPORT_SG1_BACKUP_FILE)

	// Install database primary release which serves sg0
	databaseOptions.SetValues["database.sm.storageGroup.name"] = "sg0"
	databaseOptions.SetValues["database.autoImport.type"] = "backupset"
	databaseOptions.SetValues["database.autoImport.source"] = sg0BackupUrl
	sg0ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Install database secondary release which serves sg1
	databaseOptions.SetValues["database.primaryRelease"] = "false"
	databaseOptions.SetValues["database.te.enablePod"] = "false"
	databaseOptions.SetValues["database.sm.storageGroup.name"] = "sg1"
	databaseOptions.SetValues["database.autoImport.source"] = sg1BackupUrl
	sg1ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Set the total number of engines as in the test utilities they are
	// inferred from the values
	databaseOptions.SetValues["database.sm.hotCopy.replicas"] = "2"

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	sg0HcSmPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", sg0ReleaseName, opt.ClusterName, opt.DbName)
	sg1HcSmPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", sg1ReleaseName, opt.ClusterName, opt.DbName)
	sg0HcSmPodName0 := fmt.Sprintf("%s-hotcopy-0", sg0HcSmPodTemplate)
	sg1HcSmPodName0 := fmt.Sprintf("%s-hotcopy-0", sg1HcSmPodTemplate)

	// Get SM pod logs if the test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetAppLog(t, namespaceName, sg0HcSmPodName0, "_post-restore", &corev1.PodLogOptions{})
		testlib.GetAppLog(t, namespaceName, sg1HcSmPodName0, "_post-restore", &corev1.PodLogOptions{})
	})

	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	verifyStorageGroup := func(sgName string) {
		// wait for archives to be added to the storage group
		sg := testlib.AwaitStorageGroup(t, namespaceName, admin0, opt.DbName, sgName, 30*time.Second)
		require.Equal(t, sg.State, "Available")
		require.GreaterOrEqual(t, len(sg.ArchiveStates), 1)
		require.GreaterOrEqual(t, len(sg.ProcessStates), 1)
		require.GreaterOrEqual(t, len(sg.LeaderCandidates), 1)
	}
	verifyStorageGroup("sg0")
	verifyStorageGroup("sg1")

	// Start TEs after all storage groups are available
	teDeployment := fmt.Sprintf("te-%s-nuodb-%s-%s", sg0ReleaseName, opt.ClusterName, opt.DbName)
	testlib.ScaleDeployment(t, namespaceName, teDeployment, 1)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, 3)

	// Verify that the database contains expected data in both storage groups
	output, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "select sg, count(*) from codes group by sg")
	require.NoError(t, err, "error running SQL: select sg, count(*) from codes group by sg")
	require.True(t, regexp.MustCompile(`sg0\s+1`).MatchString(output), "Unexpected data in sg0: ", output)
	require.True(t, regexp.MustCompile(`sg1\s+1`).MatchString(output), "Unexpected data in sg1: ", output)
}

func TestKubernetesRestoreDatabaseWithURL(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
	})

	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	opt := testlib.GetExtractedOptions(&options)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

	// Execute initial backup
	backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup0)

	testlib.CreateQuickstartSchema(t, namespaceName, admin0)

	defer testlib.Teardown(testlib.TEARDOWN_NGINX)

	// prepare backupset tarball and upload it on HTTP server
	tarFilePath := fmt.Sprintf("/tmp/%s.tar.gz", backupset)
	k8s.RunKubectl(t, opts, "exec", smPodName0, "--", "bash", "-c",
		fmt.Sprintf("cd /var/opt/nuodb/ && tar -czvf %s backup/%s", tarFilePath, backupset))
	remoteUrl := testlib.ServePodFileViaHTTP(t, namespaceName, smPodName0, tarFilePath)

	// restore database and set stripLevels setting
	options.SetValues["restore.source"] = remoteUrl
	options.SetValues["restore.stripLevels"] = "2"
	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
	testlib.RestoreDatabase(t, namespaceName, admin0, &options)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)

	go testlib.GetAppLog(t, namespaceName, smPodName0, "_post-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, tePodName, "_post-restart", &corev1.PodLogOptions{Follow: true})

	// verify that the database does NOT contain the data from AFTER the backup
	tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	require.NoError(t, err, "error running SQL: show schema User")
	require.Contains(t, tables, "No tables found in schema ", "Show schema returned: ", tables)
	testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
}

func TestKubernetesImportDatabaseSeparateJournal(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_NGINX)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	remoteUrl := testlib.ServeFileViaHTTP(t, namespaceName, testlib.IMPORT_ARCHIVE_FILE)

	t.Run("autoImportStream", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.autoImport.source":              remoteUrl,
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "true",
			},
		}

		// Install database and check that the journal content of the cold
		// backup is moved to the target journal location
		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, databaseOptions)
		verifyExternalJournal(t, namespaceName, admin0, databaseReleaseName, databaseOptions)
	})

	t.Run("autoImportBackupset", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.autoImport.source":              remoteUrl,
				"database.autoImport.type":                "backupset",
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "true",
			},
		}

		// Regardless of the specified backup type the journal content should be
		// moved to the target journal location in case a cold backup is
		// downloaded
		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, databaseOptions)
		verifyExternalJournal(t, namespaceName, admin0, databaseReleaseName, databaseOptions)
	})
}

func TestKubernetesRestoreDatabaseSeparateJournal(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_NGINX)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	remoteUrl := testlib.ServeFileViaHTTP(t, namespaceName, testlib.IMPORT_ARCHIVE_FILE)

	t.Run("databaseRestoreStream", func(t *testing.T) {
		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.name":                           "demo",
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "true",
				"restore.source":                          remoteUrl,
			},
		}

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, databaseOptions)

		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		testlib.RestoreDatabase(t, namespaceName, admin0, databaseOptions)

		verifyExternalJournal(t, namespaceName, admin0, databaseReleaseName, databaseOptions)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, "demo", "", "")
	})

	t.Run("databaseRestoreBackupset", func(t *testing.T) {
		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.name":                           "demo",
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "true",
			},
		}

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, databaseOptions)

		opt := testlib.GetExtractedOptions(databaseOptions)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
		smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

		// Execute initial backup
		backupGroup0 := fmt.Sprintf("%s-0", opt.ClusterName)
		backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", backupGroup0)

		testlib.CreateQuickstartSchema(t, namespaceName, admin0)

		// set restore source to the initial backupset
		databaseOptions.SetValues["restore.source"] = backupset
		databaseOptions.SetValues["restore.labels"] = "role hotcopy"

		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		testlib.RestoreDatabase(t, namespaceName, admin0, databaseOptions)

		verifyExternalJournal(t, namespaceName, admin0, databaseReleaseName, databaseOptions)

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.Contains(t, tables, "No tables found in schema ", "Show schema returned: ", tables)
		testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
	})
}

func TestCornerCaseKubernetesSnapshotRestore(t *testing.T) {
	// Set up domain
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

	options.KubectlOptions = kubectlOptions

	testDb := func(dbName string, archiveSnapshotName string, journalSnapshotName string, backupId string, shouldStart bool) string {
		retVal := ""

		values := map[string]string{
			"database.sm.resources.requests.cpu":          testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":       testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":          testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":       testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":           testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.name":                               dbName,
			"database.persistence.dataSourceRef.kind":     "VolumeSnapshot",
			"database.persistence.dataSourceRef.name":     archiveSnapshotName,
			"database.persistence.dataSourceRef.apiGroup": "snapshot.storage.k8s.io",
		}
		if backupId != "" {
			values["database.autoImport.backup_id"] = backupId
		}

		if journalSnapshotName != "" {
			values["database.sm.hotCopy.journalPath.enabled"] = "true"
			values["database.sm.hotCopy.journalPath.persistence.dataSourceRef.kind"] = "VolumeSnapshot"
			values["database.sm.hotCopy.journalPath.persistence.dataSourceRef.name"] = journalSnapshotName
			values["database.sm.hotCopy.journalPath.persistence.dataSourceRef.apiGroup"] = "snapshot.storage.k8s.io"
			values["database.sm.hotCopy.journalPath.persistence.storageClass"] = testlib.SNAPSHOTABLE_STORAGE_CLASS
		}

		var databaseChartName string

		if shouldStart {
			databaseChartName = testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
				SetValues: values,
			})

			output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "-c", "admin", "-c", "admin", "--", "bash", "-c",
				fmt.Sprintf("echo \"SELECT id FROM testtbl;\" | nuosql %s --user dba --password secret", dbName))

			require.NoError(t, err, output)

			require.True(t, strings.Contains(output, "123"))
		} else {
			databaseChartName = testlib.StartDatabaseNoWait(t, namespaceName, admin0, &helm.Options{
				SetValues: values,
			})

			smPod := fmt.Sprintf("sm-%s-nuodb-cluster0-%s-hotcopy-0", databaseChartName, dbName)
			var pod *corev1.Pod
			testlib.Await(t, func() bool {
				var err error
				pod, err = k8s.GetPodE(t, kubectlOptions, smPod)
				return err == nil && len(pod.Status.ContainerStatuses) > 0 &&
					pod.Status.ContainerStatuses[0].State.Terminated != nil &&
					pod.Status.ContainerStatuses[0].State.Terminated.Reason == "Error"
			}, 60*time.Second)
			retVal = k8s.GetPodLogs(t, kubectlOptions, pod, "engine")
		}

		helm.DeleteE(t, &options, databaseChartName, true)
		return retVal
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	sourceDb := "src-noj"
	sourceDatabaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":     testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.name":                         sourceDb,
		},
	})

	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "-c", "admin", "--", "bash", "-c",
		fmt.Sprintf("echo \"CREATE TABLE testtbl (id INT); INSERT INTO testtbl (id) values (123);\" | nuosql %s --user dba --password secret", sourceDb))

	require.NoError(t, err, output)

	defer testlib.Teardown(testlib.TEARDOWN_SNAPSHOT)

	smPod := fmt.Sprintf("sm-%s-nuodb-cluster0-%s-hotcopy-0", sourceDatabaseChartName, sourceDb)
	achiveVolumeName := "archive-volume-" + smPod

	noJournalNoBidSnapshotName := "noj-nobid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, noJournalNoBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("echo \"{\\\"id\\\": \\\"123abc\\\"}\" > /var/opt/nuodb/archive/nuodb/%s/backup.json", sourceDb))
	require.NoError(t, err, output)

	noJournalBidSnapshotName := "noj-bid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, noJournalBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("echo \"{\\\"id\\\": \\\"123abd\\\"}\" > /var/opt/nuodb/archive/nuodb/%s/backup.json", sourceDb))
	require.NoError(t, err, output)

	noJournalBadBidSnapshotName := "noj-badbid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, noJournalBadBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		"mkdir /var/opt/nuodb/archive/nuodb/foo && touch /var/opt/nuodb/archive/nuodb/foo/info.json")
	require.NoError(t, err, output)

	noJournalExtraArchiveSnapshotName := "noj-dupe-arch-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, noJournalExtraArchiveSnapshotName)

	helm.DeleteE(t, &options, sourceDatabaseChartName, true)

	testDb("noj-nobid", noJournalNoBidSnapshotName, "", "", true)
	testDb("noj-gdbid", noJournalBidSnapshotName, "", "123abc", true)
	testDb("noj-ignbid", noJournalBadBidSnapshotName, "", "", true)

	output = testDb("noj-misbid", noJournalNoBidSnapshotName, "", "123abc", false)
	require.Contains(t, output, "Incorrect backup id in archive")

	output = testDb("noj-badbid", noJournalBadBidSnapshotName, "", "123abc", false)
	require.Contains(t, output, "Incorrect backup id in archive")

	output = testDb("noj-2arch", noJournalExtraArchiveSnapshotName, "", "123abc", false)
	require.Contains(t, output, "Did not find exactly 1 archive:")

	sourceDb = "src-journ"
	sourceDatabaseChartName = testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                       testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                    testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                       testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                    testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":                        testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.hotCopy.journalPath.persistence.storageClass": testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.name": sourceDb,
			"database.sm.hotCopy.journalPath.enabled": "true",
		},
	})

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "-c", "admin", "--", "bash", "-c",
		fmt.Sprintf("echo \"CREATE TABLE testtbl (id INT); INSERT INTO testtbl (id) values (123);\" | nuosql %s --user dba --password secret", sourceDb))

	require.NoError(t, err, output)

	smPod = fmt.Sprintf("sm-%s-nuodb-cluster0-%s-hotcopy-0", sourceDatabaseChartName, sourceDb)
	achiveVolumeName = "archive-volume-" + smPod
	journalVolumeName := "journal-volume-" + smPod

	journalNoBidSnapshotName := "jor-nobid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, journalVolumeName, journalNoBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("echo \"{\\\"id\\\": \\\"123abc\\\"}\" > /var/opt/nuodb/archive/nuodb/%s/backup.json", sourceDb))
	require.NoError(t, err, output)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("echo \"{\\\"id\\\": \\\"123abc\\\"}\" > /var/opt/nuodb/journal/nuodb/%s/backup.json", sourceDb))
	require.NoError(t, err, output)

	archiveNeedsJournalSnapshotName := "archive-jor-snapshot"
	testlib.SnapshotVolume(t, namespaceName, achiveVolumeName, archiveNeedsJournalSnapshotName)

	journalBidSnapshotName := "jor-bid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, journalVolumeName, journalBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("echo \"{\\\"id\\\": \\\"123abd\\\"}\" > /var/opt/nuodb/journal/nuodb/%s/backup.json", sourceDb))
	require.NoError(t, err, output)

	journalBadBidSnapshotName := "jor-badbid-snapshot"
	testlib.SnapshotVolume(t, namespaceName, journalVolumeName, journalBadBidSnapshotName)

	output, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--", "bash", "-c",
		fmt.Sprintf("mv /var/opt/nuodb/journal/nuodb/%s /var/opt/nuodb/journal/nuodb/wrong", sourceDb))
	require.NoError(t, err, output)

	journalWrongPathSnapshotName := "jor-moved-snapshot"
	testlib.SnapshotVolume(t, namespaceName, journalVolumeName, journalWrongPathSnapshotName)

	helm.DeleteE(t, &options, sourceDatabaseChartName, true)

	testDb("jor-nobid", archiveNeedsJournalSnapshotName, journalNoBidSnapshotName, "", true)
	testDb("jor-gdbid", archiveNeedsJournalSnapshotName, journalBidSnapshotName, "123abc", true)
	testDb("jor-ignbid", archiveNeedsJournalSnapshotName, journalBadBidSnapshotName, "", true)

	output = testDb("jor-misbid", archiveNeedsJournalSnapshotName, journalNoBidSnapshotName, "123abc", false)
	require.Contains(t, output, "Incorrect backup id in journal")

	output = testDb("jor-badbid", archiveNeedsJournalSnapshotName, journalBadBidSnapshotName, "123abc", false)
	require.Contains(t, output, "Incorrect backup id in journal")

	output = testDb("jor-mvjor", archiveNeedsJournalSnapshotName, journalWrongPathSnapshotName, "123abc", false)
	require.Contains(t, output, "Did not find a journal snapshot at '/var/opt/nuodb/journal/nuodb/src-journ'")
}
