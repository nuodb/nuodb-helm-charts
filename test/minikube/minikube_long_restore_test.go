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

func TestKubernetesRestoreMultipleSMs(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
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

	// Execute initial backup
	backupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", opt.ClusterName)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)

	t.Run("autoRestart", func(t *testing.T) {
		testlib.CreateQuickstartSchema(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = backupset
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

func TestKubernetesRestoreWithStorageGroups(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
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
			"database.sm.hotCopy.replicas":          "2",
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
	})

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodNameTemplate := fmt.Sprintf("%s-hotcopy", smPodNameTemplate)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)

	// Create 2 storage groups, each served by only one archive
	k8s.RunKubectl(t, kubectlOptions, "exec", admin0, "--",
		"nuocmd", "add", "storage-group", "--db-name", opt.DbName, "--sg-name", "sg0", "--archive-id", "0")
	k8s.RunKubectl(t, kubectlOptions, "exec", admin0, "--",
		"nuocmd", "add", "storage-group", "--db-name", opt.DbName, "--sg-name", "sg1", "--archive-id", "1")

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

	// Execute backup
	backupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", opt.ClusterName)

	// Insert more rows
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg0', '1001')")
	testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
		"INSERT INTO codes VALUES ('sg1', '2001')")

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)

	// restore database
	databaseOptions.SetValues["restore.source"] = backupset
	testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions)
	testlib.AwaitPodLog(t, namespaceName, hcSmPodName0, "_post-restart")

	// verify that the database does NOT contain the data from AFTER the backup

	output, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "select sg, count(*) from codes group by sg")
	require.NoError(t, err, "error running SQL: select sg, count(*) from codes group by sg")
	require.True(t, regexp.MustCompile(`sg0\s+1`).MatchString(output), "Unexpected data in sg0: ", output)
	require.True(t, regexp.MustCompile(`sg1\s+1`).MatchString(output), "Unexpected data in sg1: ", output)
	testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	testlib.CheckRestoreRequests(t, namespaceName, admin0, opt.DbName, "", "")
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
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)

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
