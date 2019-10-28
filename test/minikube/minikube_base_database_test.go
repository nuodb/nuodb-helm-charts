package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const LABEL_CLOUD = "minikube"
const LABEL_REGION = "local"
const LABEL_ZONE = "local-b"

func populateCreateDBData(t *testing.T, namespaceName string, adminPod string) {
	// populate some data
	opts := k8s.NewKubectlOptions("", "")
	opts.Namespace = namespaceName
	k8s.RunKubectl(t, opts,
		"exec", adminPod, "--",
		"/opt/nuodb/bin/nuosql",
		"--user", "dba",
		"--password", "secret",
		"demo",
		"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql",
	)
}

func verifyNuoSQL(t *testing.T, namespaceName string, adminPod string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--", "bash", "-c",
		"echo \"select * from system.nodes;\" | nuosql demo@localhost --user dba --password secret")

	assert.NilError(t, err, output)

	assert.Check(t, strings.Contains(output, "Storage"))
	assert.Check(t, strings.Contains(output, "Transaction"))
}

func verifySecret(t *testing.T, namespaceName string) {
	secret := testlib.GetSecret(t, namespaceName, "demo.nuodb.com")

	_, ok := secret.Data["database-name"]
	assert.Check(t, ok)

	_, ok = secret.Data["database-password"]
	assert.Check(t, ok)

	_, ok = secret.Data["database-username"]
	assert.Check(t, ok)
}

func verifyDBService(t *testing.T, namespaceName string, podName string) {
	serviceName := "demo"

	adminService := testlib.GetService(t, namespaceName, serviceName)
	assert.Equal(t, adminService.Name, serviceName)

	testlib.PingService(t, namespaceName, serviceName, podName)
}

func verifyPodLabeling(t *testing.T, namespaceName string, adminPod string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", "demo")

	assert.NilError(t, err, output)

	err, objects := testlib.Unmarshal(output)

	for _, obj := range objects {
		val, ok := obj.Labels["cloud"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_CLOUD)

		val, ok = obj.Labels["region"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_REGION)

		val, ok = obj.Labels["zone"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_ZONE)
	}

}

func verifyPacketFetch(t *testing.T, namespaceName string, admin0 string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	// verify the container can actually download the file from the internet
	start := time.Now()
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions,
		"exec", admin0, "--",
		"bash", "-c",
		fmt.Sprintf("curl -k %s | tar tzf - | head -n 10", testlib.IMPORT_ARCHIVE_URL),
	)
	assert.NilError(t, err, "Could not fetch archive")
	elapsed := time.Since(start)
	t.Logf("Fetching package (%s) took %f seconds", testlib.IMPORT_ARCHIVE_URL, elapsed.Seconds())
	t.Log("tar contents: ", output)
}

func backupDatabase(t *testing.T, namespaceName string, podName string, databaseName string, options *helm.Options) {
	randomSuffix := strings.ToLower(random.UniqueId())

	bakName := fmt.Sprintf("backup-full-%s", randomSuffix)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	testlib.AddTeardown(testlib.TEARDOWN_BACKUP, func() { helm.Delete(t, options, bakName, true) })
	helm.Install(t, options, testlib.BACKUP_HELM_CHART_PATH, bakName)

	// wait for the backup to both start _and_ complete successfully
	backupJob := fmt.Sprintf("backup-%s-job-full", databaseName)
	testlib.AwaitPodPhase(t, namespaceName, backupJob, corev1.PodSucceeded, 120*time.Second)

	// verify that the backup has been documented by the Admin layer
	backupName := fmt.Sprintf("%s/%s", testlib.LAST_BACKUP_PREFIX, databaseName)
	backupset, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions,
		"exec", podName, "--",
		"nuocmd", "get", "value",
		"--key", backupName,
	)

	assert.NilError(t, err, "Error running: nuocmd get value --key ", backupName)
	assert.Check(t, backupset != "")
}

func restoreDatabase(t *testing.T, namespaceName string) {
	// run the restore chart - which flags the database to restore on next startup
	restName := "restore-demo"
	options := &helm.Options{
		SetValues: map[string]string{
			"database.name":     "demo",
			"restore.target":    "demo",
			"restore.backupSet": ":latest",
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	helm.Install(t, options, testlib.RESTORE_HELM_CHART_PATH, restName)
	testlib.AddTeardown(testlib.TEARDOWN_RESTORE, func() { helm.Delete(t, options, restName, true) })

	testlib.AwaitPodPhase(t, namespaceName, restName, corev1.PodSucceeded, 120*time.Second)
}

func TestKubernetesBasicDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{"database.sm.resources.requests.cpu": "500m",
				"database.sm.resources.requests.memory": "1Gi",
				"database.te.resources.requests.cpu":    "500m",
				"database.te.resources.requests.memory": "1Gi",
				"database.te.labels.cloud":              LABEL_CLOUD,
				"database.te.labels.region":             LABEL_REGION,
				"database.te.labels.zone":               LABEL_ZONE,
				"database.sm.labels.cloud":              LABEL_CLOUD,
				"database.sm.labels.region":             LABEL_REGION,
				"database.sm.labels.zone":               LABEL_ZONE,
			},
		})

		t.Run("verifySecret", func(t *testing.T) { verifySecret(t, namespaceName) })
		t.Run("verifyDBService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0) })
		t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0) })
		t.Run("verifyPodLabeling", func(t *testing.T) { verifyPodLabeling(t, namespaceName, admin0) })
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		testlib.StartDatabase(t, namespaceName, admin0,
			&helm.Options{
				SetValues: map[string]string{"database.sm.resources.requests.cpu": "500m",
					"database.sm.resources.requests.memory": "1Gi",
					"database.te.resources.requests.cpu":    "500m",
					"database.te.resources.requests.memory": "1Gi",
					"database.te.labels.cloud":              LABEL_CLOUD,
					"database.te.labels.region":             LABEL_REGION,
					"database.te.labels.zone":               LABEL_ZONE,
					"database.sm.labels.cloud":              LABEL_CLOUD,
					"database.sm.labels.region":             LABEL_REGION,
					"database.sm.labels.zone":               LABEL_ZONE,
					"database.enableDaemonSet":              "true",
					// prevent non-backup SM from scheduling
					"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
				},
			},
		)

		t.Run("verifySecret", func(t *testing.T) { verifySecret(t, namespaceName) })
		t.Run("verifyDBService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0) })
		t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0) })
		t.Run("verifyPodLabeling", func(t *testing.T) { verifyPodLabeling(t, namespaceName, admin0) })
	})
}

func TestKubernetesBackupDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    "500m",
				"database.sm.resources.requests.memory": "1Gi",
				"database.te.resources.requests.cpu":    "500m",
				"database.te.resources.requests.memory": "1Gi",
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		backupDatabase(t, namespaceName, admin0, "demo", &databaseOptions)
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    "500m",
				"database.sm.resources.requests.memory": "1Gi",
				"database.te.resources.requests.cpu":    "500m",
				"database.te.resources.requests.memory": "1Gi",
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
				"database.enableDaemonSet":              "true",
				// prevent non-backup SM from scheduling
				"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		backupDatabase(t, namespaceName, admin0, "demo", &databaseOptions)
	})
}

func TestKubernetesRestoreDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.name":                         "demo",
				"database.sm.resources.requests.cpu":    "500m",
				"database.sm.resources.requests.memory": "1Gi",
				"database.te.resources.requests.cpu":    "500m",
				"database.te.resources.requests.memory": "1Gi",
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		opts := k8s.NewKubectlOptions("", "")
		opts.Namespace = namespaceName
		databaseOptions.KubectlOptions = opts

		// populate some data
		k8s.RunKubectl(t, opts,
			"exec", admin0, "--",
			"/opt/nuodb/bin/nuosql",
			"--user", "dba",
			"--password", "secret",
			"demo",
			"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql")

		// run a manual backup
		k8s.RunKubectl(t, opts,
			"exec", admin0, "--",
			"nuodocker", "backup", "database",
			"--db-name", "demo",
			"--type", "full",
			"--timeout", "120",
			"--backup-root", "/var/opt/nuodb/backup",
		)

		// populate some more data
		k8s.RunKubectl(t, opts,
			"exec", admin0, "--",
			"/opt/nuodb/bin/nuosql",
			"--user", "dba",
			"--password", "secret",
			"demo",
			"--file", "/opt/nuodb/samples/quickstart/sql/Teams.sql",
		)

		defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
		restoreDatabase(t, namespaceName)

		// and restart database - to trigger the restore
		k8s.RunKubectl(t, opts,
			"exec", admin0, "--",
			"nuocmd", "shutdown", "database",
			"--db-name", "demo",
		)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo")

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		assert.NilError(t, err, "error running SQL: show schema User")
		assert.Check(t, strings.Contains(tables, "HOCKEY"), "tables returned: ", tables)

		// verify that the database does NOT contain the data from AFTER the backup
		count, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "select 'count=' || count(*) from User.Teams")
		assert.NilError(t, err, "error running SQL: select count(*) from User.Teams")
		assert.Check(t, strings.Contains(count, "count=0"), "count returned: ", count)
	})
}

func TestKubernetesImportDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)

	verifyPacketFetch(t, namespaceName, admin0)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.import.url":                   testlib.IMPORT_ARCHIVE_URL,
				"database.sm.resources.requests.cpu":    "500m",
				"database.sm.resources.requests.memory": "1Gi",
				"database.te.resources.requests.cpu":    "500m",
				"database.te.resources.requests.memory": "1Gi",
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		})

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		assert.NilError(t, err, "error running SQL: show schema User")
		assert.Check(t, strings.Contains(tables, "HOCKEY"))
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		testlib.StartDatabase(t, namespaceName, admin0,
			&helm.Options{
				SetValues: map[string]string{
					"database.import.url":                   testlib.IMPORT_ARCHIVE_URL,
					"database.sm.resources.requests.cpu":    "500m",
					"database.sm.resources.requests.memory": "1Gi",
					"database.te.resources.requests.cpu":    "500m",
					"database.te.resources.requests.memory": "1Gi",
					"database.enableDaemonSet":              "true",
					// prevent non-backup SM from scheduling
					"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
				},
			},
		)

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		assert.NilError(t, err, "error running SQL: show schema User")
		assert.Check(t, strings.Contains(tables, "HOCKEY"))
	})
}
