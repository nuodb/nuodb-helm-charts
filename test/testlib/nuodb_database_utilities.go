package testlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/core/v1"

	"github.com/Masterminds/semver"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

const UPGRADE_STRATEGY = `
spec:
  strategy:
    $retainKeys:
    - type
    type: Recreate
`

const SNAPSHOT_TEMPLATE = `
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: %s
spec:
  volumeSnapshotClassName: %s
  source:
    persistentVolumeClaimName: %s
`

type ExtractedOptions struct {
	NrTePods              int
	NrSmHotCopyPods       int
	NrSmNoHotCopyPods     int
	NrSmPods              int
	DbName                string
	ClusterName           string
	EntrypointClusterName string
	DbPrimaryRelease      bool
	DomainName            string
}

var nuoDBVersionMap map[string]string = make(map[string]string)

var NuoDBVersionRegex = regexp.MustCompile(`^([0-9]+\.[0-9]+(?:\.[0-9]+)?).*$`)

func GetExtractedOptions(options *helm.Options) (opt ExtractedOptions) {
	var err error

	opt.NrTePods, err = strconv.Atoi(options.SetValues["database.te.replicas"])
	if err != nil {
		opt.NrTePods = 1
	}

	if options.SetValues["database.te.enablePod"] == "false" {
		opt.NrTePods = 0
	}

	opt.NrSmHotCopyPods, err = strconv.Atoi(options.SetValues["database.sm.hotCopy.replicas"])
	if err != nil {
		opt.NrSmHotCopyPods = 1
	}

	if options.SetValues["database.sm.hotCopy.enablePod"] == "false" {
		opt.NrSmHotCopyPods = 0
	}

	opt.NrSmNoHotCopyPods, err = strconv.Atoi(options.SetValues["database.sm.noHotCopy.replicas"])
	if err != nil {
		opt.NrSmNoHotCopyPods = 0
	}

	if options.SetValues["database.sm.noHotCopy.enablePod"] == "false" {
		opt.NrSmNoHotCopyPods = 0
	}

	opt.NrSmPods = opt.NrSmNoHotCopyPods + opt.NrSmHotCopyPods

	opt.DbName = options.SetValues["database.name"]
	if len(opt.DbName) == 0 {
		opt.DbName = "demo"
	}

	opt.ClusterName = options.SetValues["cloud.cluster.name"]
	if len(opt.ClusterName) == 0 {
		opt.ClusterName = "cluster0"
	}

	opt.EntrypointClusterName = options.SetValues["cloud.cluster.entrypointName"]
	if len(opt.EntrypointClusterName) == 0 {
		opt.EntrypointClusterName = "cluster0"
	}

	if options.SetValues["database.primaryRelease"] == "false" {
		opt.DbPrimaryRelease = false
	} else {
		opt.DbPrimaryRelease = true
	}

	opt.DomainName = options.SetValues["admin.domain"]
	if opt.DomainName == "" {
		opt.DomainName = "nuodb"
	}

	return
}

func EnsureDatabaseNotRunning(t *testing.T, adminPod string, opt ExtractedOptions, kubectlOptions *k8s.KubectlOptions) {
	// invoke shutdown database; this may fail if the database is already NOT_RUNNING, which is okay
	k8s.RunKubectlE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "shutdown", "database", "--db-name", opt.DbName)

	// wait for all database processes to exit
	_, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "check", "database", "--db-name", opt.DbName, "--num-processes", "0", "--timeout", "30")
	if err != nil {
		// suppress error caused by database entering TOMBSTONE state,
		// which can happen due to KAA deleting the database object
		// after the database Helm release is uninstalled
		require.Contains(t, err.Error(), "Database "+opt.DbName+" in TOMBSTONE state")
	}
}

func RestartDatabasePods(t *testing.T, namespaceName string, helmChartReleaseName string, options *helm.Options) {
	opt := GetExtractedOptions(options)
	hcSmPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s-hotcopy", helmChartReleaseName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)
	var toDelete []string
	tes := GetPodNames(t, namespaceName, tePodNameTemplate)
	require.Equal(t, opt.NrTePods, len(tes), "Unexpected number of TE Pods")
	toDelete = append(toDelete, tes...)
	for i := 0; i < opt.NrSmHotCopyPods; i++ {
		toDelete = append(toDelete, fmt.Sprintf("%s-%d", hcSmPodNameTemplate, i))
	}
	for i := 0; i < opt.NrSmNoHotCopyPods; i++ {
		toDelete = append(toDelete, fmt.Sprintf("%s-%d", smPodNameTemplate, i))
	}
	for _, podName := range toDelete {
		DeletePod(t, namespaceName, "pod/"+podName)
	}
	AwaitNrReplicasScheduled(t, namespaceName, tePodNameTemplate, opt.NrTePods)
	AwaitNrReplicasScheduled(t, namespaceName, smPodNameTemplate, opt.NrSmPods)
}

type DatabaseInstallationStep func(t *testing.T, options *helm.Options, helmChartReleaseName string)

func StartDatabaseTemplate(t *testing.T, namespaceName string, adminPod string, options *helm.Options, installationStep DatabaseInstallationStep, awaitDatabase bool) (helmChartReleaseName string) {
	//Truncation done to reduce pod name length
	randomSuffix := strings.ToLower(random.UniqueId())[:5]

	InjectTestValues(t, options)
	opt := GetExtractedOptions(options)

	if IsOpenShiftEnvironment(t) {
		THPReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
		AddTeardown(TEARDOWN_DATABASE, func() {
			helm.Delete(t, options, THPReleaseName, true)
		})
		helm.Install(t, options, THP_HELM_CHART_PATH, THPReleaseName)

		AwaitNrReplicasReady(t, namespaceName, THPReleaseName, 1)
	}

	helmChartReleaseName = fmt.Sprintf("database-%s", randomSuffix)
	tePodNameTemplate := fmt.Sprintf("te-%s-%s-%s-%s", helmChartReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-%s-%s-%s", helmChartReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	// with Async actions which do not return a cleanup method, create the teardown(s) first
	AddTeardown(TEARDOWN_DATABASE, func() {
		// suppress error as the helm release might be already deleted
		helm.DeleteE(t, options, helmChartReleaseName, true)
		AwaitNoPods(t, namespaceName, helmChartReleaseName)
		// Delete database only if admin pod exists and tearing down the entrypoint cluster
		_, err := FindPod(t, namespaceName, adminPod)
		if err == nil && opt.ClusterName == opt.EntrypointClusterName && opt.DbPrimaryRelease {
			db, err := GetDatabaseE(t, namespaceName, adminPod, opt.DbName)
			// delete the database if it is not deleted already
			if err == nil && db.State != "TOMBSTONE" {
				EnsureDatabaseNotRunning(t, adminPod, opt, kubectlOptions)
				DeleteDatabase(t, namespaceName, opt.DbName, adminPod)
				// purge all database archives so that multiple database tests can
				// be executed with a single admin deployment
				PurgeDatabaseArchives(t, namespaceName, adminPod, opt.DbName)
			}
		}
	})

	installationStep(t, options, helmChartReleaseName)

	if awaitDatabase {
		AddDiagnosticTeardown(TEARDOWN_DATABASE, t, func() {
			DescribePods(t, namespaceName, tePodNameTemplate)
			DescribePods(t, namespaceName, smPodNameTemplate)
		})
		AwaitNrReplicasScheduled(t, namespaceName, tePodNameTemplate, opt.NrTePods)
		AwaitNrReplicasScheduled(t, namespaceName, smPodNameTemplate, opt.NrSmPods)

		// NOTE: the Teardown logic will pick a TE/SM that is running during teardown time. Not the TE/SM that was running originally
		// this is relevant for any tests that restart TEs/SMs
		AddTeardown(TEARDOWN_DATABASE, func() {
			pods, _ := findPods(t, namespaceName, tePodNameTemplate)
			for _, tePod := range pods {
				t.Logf("Getting log from TE pod: %s", tePod.Name)
				go GetAppLog(t, namespaceName, tePod.Name, "", &v12.PodLogOptions{Follow: true})
			}
		})
		// the TEs will become RUNNING after the SMs as they need an entry node
		// so use the same ready timeout for both
		readyTimeout := AdjustPodTimeout(tePodNameTemplate, 300*time.Second)
		if opt.NrTePods > 0 {
			tePodName := GetPodName(t, namespaceName, tePodNameTemplate)
			AwaitPodUp(t, namespaceName, tePodName, readyTimeout)
		}

		if opt.NrSmPods > 0 {
			smPodName0 := GetPodName(t, namespaceName, smPodNameTemplate)
			AddTeardown(TEARDOWN_DATABASE, func() {
				pods, _ := findPods(t, namespaceName, smPodNameTemplate)
				for _, smPod := range pods {
					t.Logf("Getting log from SM pod: %s", smPod.Name)
					go GetAppLog(t, namespaceName, smPod.Name, "", &v12.PodLogOptions{Follow: true})
				}
			})
			AwaitPodUp(t, namespaceName, smPodName0, readyTimeout)
		}

		// Await num of database processes only for single cluster deployments
		// and if the database release is primary; in multi-clusters the await
		// logic should be called once all clusters are installed with the
		// database chart
		if opt.ClusterName == opt.EntrypointClusterName && opt.DbPrimaryRelease {
			AwaitDatabaseUp(t, namespaceName, adminPod, opt.DbName, opt.NrSmPods+opt.NrTePods)
		}
	}

	return
}

func InstallDatabase(t *testing.T, options *helm.Options, helmChartReleaseName string) {
	if options.Version == "" {
		helm.Install(t, options, DATABASE_HELM_CHART_PATH, helmChartReleaseName)
	} else {
		helm.Install(t, options, "nuodb/database", helmChartReleaseName)
	}
}

func StartDatabase(t *testing.T, namespace string, adminPod string, options *helm.Options) string {
	return StartDatabaseTemplate(t, namespace, adminPod, options, InstallDatabase, true)
}

func StartDatabaseNoWait(t *testing.T, namespace string, adminPod string, options *helm.Options) string {
	return StartDatabaseTemplate(t, namespace, adminPod, options, InstallDatabase, false)
}

func UpgradeDatabase(t *testing.T, namespaceName string, helmChartReleaseName string, adminPod string, options *helm.Options, upgradeOptions *UpgradeOptions) {
	opt := GetExtractedOptions(options)

	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)
	smPodName := fmt.Sprintf("sm-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	// get the OLD log
	tePodName := GetPodName(t, namespaceName, tePodNameTemplate)
	tePod := GetPod(t, namespaceName, tePodName)
	go GetAppLog(t, namespaceName, tePodName, "-previous", &v12.PodLogOptions{Follow: true})

	smPodName0 := GetPodName(t, namespaceName, smPodName)
	smPod0 := GetPod(t, namespaceName, smPodName0)
	go GetAppLog(t, namespaceName, smPodName0, "-previous", &v12.PodLogOptions{Follow: true})

	AddDiagnosticTeardown(TEARDOWN_DATABASE, t, func() {
		_ = k8s.RunKubectlE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "show", "domain")
		_ = k8s.RunKubectlE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "show", "archives")
		_ = k8s.RunKubectlE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "show", "database", "--db-name", "demo", "--all-incarnations")
	})

	SetDeploymentUpgradeStrategyToRecreate(t, namespaceName, fmt.Sprintf("te-%s-nuodb-cluster0-demo", helmChartReleaseName))

	helm.Upgrade(t, options, DATABASE_HELM_CHART_PATH, helmChartReleaseName)

	if upgradeOptions.TePodShouldGetRecreated {
		AwaitPodObjectRecreated(t, namespaceName, tePod, 30*time.Second)
	}
	// the TE ReplicaSet can be recreated with a different name so fetch the pod
	// using a template continuously
	Await(t, func() bool {
		pod, err := FindPod(t, namespaceName, tePodNameTemplate)
		if err != nil {
			t.Logf("%s: %s", err.Error(), tePodNameTemplate)
			return false
		}
		return arePodConditionsMet(pod, corev1.PodReady, corev1.ConditionTrue)
	}, 180*time.Second)

	if upgradeOptions.SmPodShouldGetRecreated {
		AwaitPodObjectRecreated(t, namespaceName, smPod0, 30*time.Second)
	}
	AwaitPodUp(t, namespaceName, smPodName0, 300*time.Second)

	// Await num of database processes only for single cluster deployment;
	// in multi-clusters the await logic should be called once all clusters
	// are installed with the database chart
	if opt.ClusterName == opt.EntrypointClusterName {
		AwaitDatabaseUp(t, namespaceName, adminPod, opt.DbName, opt.NrSmPods+opt.NrTePods)
	}
}

func SetDeploymentUpgradeStrategyToRecreate(t *testing.T, namespaceName string, deploymentName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.RunKubectl(t, kubectlOptions, "patch", "deployment", deploymentName, "-p", UPGRADE_STRATEGY)
}

func RestoreDatabase(t *testing.T, namespaceName string, podName string, databaseOptions *helm.Options) {
	// run the restore chart - which flags the database to restore on next startup
	randomSuffix := strings.ToLower(random.UniqueId())

	restName := fmt.Sprintf("restore-demo-%s", randomSuffix)
	options := &helm.Options{
		SetValues: map[string]string{
			"database.name":       "demo",
			"restore.target":      "demo",
			"restore.source":      ":latest",
			"restore.autoRestart": "true",
		},
	}
	for key, value := range databaseOptions.SetValues {
		options.SetValues[key] = value
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	restore := func() {
		// Get restore pod logs and events on failure
		AddDiagnosticTeardown(TEARDOWN_RESTORE, t, func() {
			restorePodName := GetPodName(t, namespaceName, "restore-demo-")
			k8s.RunKubectl(t, kubectlOptions, "describe", "pod", restorePodName)
			GetAppLog(t, namespaceName, restorePodName, "restore-job", &v12.PodLogOptions{})
		})
		// Remove restore job if exist as it's not unique for a restore chart release
		k8s.RunKubectlE(t, kubectlOptions, "delete", "job", "restore-"+options.SetValues["database.name"])
		InjectTestValues(t, options)
		helm.Install(t, options, RESTORE_HELM_CHART_PATH, restName)
		AddTeardown(TEARDOWN_RESTORE, func() { helm.Delete(t, options, restName, true) })
		// Using a bit longer timeout here as we might be performing restore using
		// older image which needs to be pulled
		AwaitPodPhase(t, namespaceName, "restore-demo-", corev1.PodSucceeded, 300*time.Second)
	}

	if options.SetValues["restore.autoRestart"] == "true" {
		AwaitDatabaseRestart(t, namespaceName, podName, "demo", databaseOptions, restore)
	} else {
		restore()
	}
}

func BackupDatabase(
	t *testing.T,
	namespaceName string,
	podName string,
	databaseName string,
	backupType string,
	backupGroup string,
) string {
	err := BackupDatabaseE(t, namespaceName, podName, databaseName, backupType, backupGroup)
	require.NoError(t, err)
	return GetLatestBackup(t, namespaceName, podName, databaseName, backupGroup)
}

func BackupDatabaseE(
	t *testing.T,
	namespaceName string,
	podName string,
	databaseName string,
	backupType string,
	backupGroup string,
) error {
	cronJobName := fmt.Sprintf("%s-hotcopy-nuodb-%s-%s", backupType, databaseName, backupGroup)
	jobName := fmt.Sprintf("backup-database-%s", strings.ToLower(random.UniqueId()))
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	if err := k8s.RunKubectlE(t, opts, "create", "job", "--from=cronjob/"+cronJobName, jobName); err != nil {
		return err
	}
	defer func() {
		// Get backup logs and delete the job
		if pod, err := FindPod(t, namespaceName, jobName); err == nil {
			GetAppLog(t, namespaceName, pod.Name, "", &v12.PodLogOptions{Container: "nuodb"})
		}
		k8s.RunKubectl(t, opts, "delete", "job", jobName)
	}()
	var backupErr error
	if err := AwaitE(t, func() bool {
		pod, err := FindPod(t, namespaceName, jobName)
		if err != nil {
			return false
		}
		if pod.Status.Phase == corev1.PodSucceeded {
			return true
		}
		var restartCount int32
		for _, containerStatus := range pod.Status.ContainerStatuses {
			restartCount += containerStatus.RestartCount
		}
		if restartCount > 0 {
			// some of the containers failed and was restarted; stop waiting
			buf, err := io.ReadAll(getAppLogStream(
				t, namespaceName, pod.Name,
				&corev1.PodLogOptions{Previous: true, Container: "nuodb"}))
			if err != nil {
				backupErr = err
			} else {
				backupErr = fmt.Errorf("database backup failed: %s", string(buf))
			}
			return true
		}
		return false
	}, 120*time.Second); err != nil {
		return fmt.Errorf("database backup exceed 120 sec: %w", err)
	}
	return backupErr
}

func GetLatestBackup(t *testing.T, namespaceName string, podName string,
	databaseName string, backupGroup string) string {
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	backupset, err := k8s.RunKubectlAndGetOutputE(t, opts,
		"exec", podName, "-c", "engine", "--", "bash", "-c",
		"nuobackup --type report-latest --db-name "+databaseName+
			" --group "+backupGroup+" --backup-root /var/opt/nuodb/backup 2>/dev/null",
	)
	require.NoError(t, err, "Error while reporting latest backupset")
	require.True(t, backupset != "")
	return backupset
}

func GetLatestBackupGroup(t *testing.T, namespaceName string, podName string, databaseName string) string {
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	group, err := k8s.RunKubectlAndGetOutputE(t, opts,
		"exec", podName, "-c", "engine", "--", "bash", "-c",
		"nuobackup --type report-latest --db-name "+databaseName+
			" --backup-root /var/opt/nuodb/backup 2>/dev/null",
	)
	require.NoError(t, err, "Error while reporting latest backup group")
	require.True(t, group != "")
	return group
}

func GetArchives(t *testing.T, namespaceName string, adminPod string, dbName string) (archives []NuoDBArchive, removedArchives []NuoDBArchive) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	// get archives
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", dbName)
	require.NoError(t, err, output)
	err, archives = UnmarshalArchives(output)
	require.NoError(t, err)

	// get removed archives
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", dbName, "--removed")
	require.NoError(t, err, output)
	err, removedArchives = UnmarshalArchives(output)
	require.NoError(t, err)
	return
}

func CheckArchives(t *testing.T, namespaceName string, adminPod string, dbName string, numExpected int, numExpectedRemoved int) (archives []NuoDBArchive, removedArchives []NuoDBArchive) {
	archives, removedArchives = GetArchives(t, namespaceName, adminPod, dbName)
	require.Equal(t, numExpected, len(archives), archives)
	require.Equal(t, numExpectedRemoved, len(removedArchives), removedArchives)
	return
}

func AwaitStorageGroup(t *testing.T, namespaceName, adminPod, dbName, sgName string, timeout time.Duration) *NuoDBStorageGroup {
	var sg *NuoDBStorageGroup
	Await(t, func() bool {
		sg = GetStorageGroup(t, namespaceName, adminPod, dbName, sgName)
		if len(sg.ArchiveStates) == 0 {
			return false
		}
		for archiveId, state := range sg.ArchiveStates {
			if state != "ADDED" {
				t.Logf("Waiting for storage group sg=%s, archiveId=%s, state=%s",
					sgName, archiveId, state)
				return false
			}
		}
		for nodeId, state := range sg.ProcessStates {
			if state != "RUNNING" {
				t.Logf("Waiting for process in storage group sg=%s, nodeId=%s, state=%s",
					sgName, nodeId, state)
				return false
			}
		}
		return true
	}, timeout)
	return sg
}

func GetStorageGroup(t *testing.T, namespaceName, adminPod, dbName, sgName string) *NuoDBStorageGroup {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "storage-groups", "--db-name", dbName)
	require.NoError(t, err, output)
	err, sgs := UnmarshalStorageGroups(output)
	require.NoError(t, err)
	// the admin converts all storage group names to upper case
	expectedSgName := strings.ToUpper(sgName)
	var found *NuoDBStorageGroup
	for _, sg := range sgs {
		if sg.Name == expectedSgName {
			found = &sg
			break
		}
	}
	require.NotNil(t, found, "Storage group sgName=%s, dbName=%s not found", sgName, dbName)
	return found
}

func CreateQuickstartSchema(t *testing.T, namespaceName string, adminPod string) {
	opts := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.RunKubectl(t, opts,
		"exec", adminPod, "--",
		"/opt/nuodb/bin/nuosql",
		"--user", "dba",
		"--password", "secret",
		"demo",
		"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql",
	)

	// verify that the database contains the populated data
	tables, err := RunSQL(t, namespaceName, adminPod, "demo", "show schema User")
	require.NoError(t, err, "error running SQL: show schema User")
	require.True(t, strings.Contains(tables, "HOCKEY"), "tables returned: ", tables)
}

func IsRestoreRequestSupported(t *testing.T, namespaceName string, podName string) bool {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	err := k8s.RunKubectlE(t, kubectlOptions, "exec", podName, "--",
		"bash", "-c", "nuodocker request restore -h > /dev/null")
	return err == nil
}

func CheckRestoreRequests(t *testing.T, namespaceName string, podName string, databaseName string,
	expectedValue string, expectedLegacyValue string) (string, string) {

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	legacyRestoreRequest := ""
	restoreRequest := ""
	Await(t, func() bool {
		legacyRestoreRequest, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "-c", "admin", "--",
			"nuocmd", "get", "value", "--key", fmt.Sprintf("/nuodb/nuosm/database/%s/restore", databaseName))
		// Legacy restore request is cleared async
		return err == nil && legacyRestoreRequest == expectedLegacyValue
	}, 30*time.Second)
	if IsRestoreRequestSupported(t, namespaceName, podName) {
		restoreRequest, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "-c", "admin", "--",
			"nuodocker", "get", "restore-requests", "--db-name", databaseName)
		require.NoError(t, err)
		require.Equal(t, expectedValue, restoreRequest)
	}
	return restoreRequest, legacyRestoreRequest
}

func CreateNginxDeployment(t *testing.T, namespaceName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.RunKubectl(t, kubectlOptions, "create", "deployment", NGINX_DEPLOYMENT, "--image=nginx")
	k8s.RunKubectl(t, kubectlOptions, "create", "service", "clusterip", NGINX_DEPLOYMENT, "--tcp=80:80")
	podName := GetPodName(t, namespaceName, NGINX_DEPLOYMENT)
	AwaitPodUp(t, namespaceName, podName, 60*time.Second)

	AddDiagnosticTeardown(TEARDOWN_NGINX, t, func() {
		k8s.RunKubectl(t, kubectlOptions, "describe", "deployment", NGINX_DEPLOYMENT)
		podName := GetPodName(t, namespaceName, NGINX_DEPLOYMENT)
		GetAppLog(t, namespaceName, podName, "_nginx", &corev1.PodLogOptions{})
	})

	AddTeardown(TEARDOWN_NGINX, func() {
		k8s.RunKubectl(t, kubectlOptions, "delete", "service", NGINX_DEPLOYMENT)
		k8s.RunKubectl(t, kubectlOptions, "delete", "deployment", NGINX_DEPLOYMENT)
	})
}

func ServePodFileViaHTTP(t *testing.T, namespaceName string, srcPodName string, filePath string) string {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "deployment", NGINX_DEPLOYMENT)
	if !strings.Contains(output, NGINX_DEPLOYMENT) || err != nil {
		CreateNginxDeployment(t, namespaceName)
	}

	prefix := "nginx-html"
	tmpDir, err := os.MkdirTemp("", prefix)
	require.NoError(t, err, "Unable to create TMP directory with prefix ", prefix)
	defer os.RemoveAll(tmpDir)
	fileName := filepath.Base(filePath)
	localFilePath := path.Join(tmpDir, fileName)

	k8s.RunKubectl(t, kubectlOptions, "cp", srcPodName+":"+filePath, localFilePath)
	nginxPod := GetPodName(t, namespaceName, NGINX_DEPLOYMENT)
	k8s.RunKubectl(t, kubectlOptions, "cp", localFilePath, nginxPod+":/usr/share/nginx/html/"+fileName)
	return fmt.Sprintf("http://nginx.%s.svc.cluster.local/%s", namespaceName, fileName)
}

func ServeFileViaHTTP(t *testing.T, namespaceName string, localFilePath string) string {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "deployment", NGINX_DEPLOYMENT)
	if !strings.Contains(output, NGINX_DEPLOYMENT) || err != nil {
		CreateNginxDeployment(t, namespaceName)
	}
	fileName := filepath.Base(localFilePath)
	nginxPod := GetPodName(t, namespaceName, NGINX_DEPLOYMENT)
	k8s.RunKubectl(t, kubectlOptions, "cp", localFilePath, nginxPod+":/usr/share/nginx/html/"+fileName)
	return fmt.Sprintf("http://nginx.%s.svc.cluster.local/%s", namespaceName, fileName)
}

func GetNuoDBVersion(t *testing.T, namespaceName string, options *helm.Options) string {
	if version, ok := nuoDBVersionMap[options.SetValues["nuodb.image.tag"]]; ok {
		return version
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	randomSuffix := strings.ToLower(random.UniqueId())
	podName := fmt.Sprintf("nuodb-version-%s", randomSuffix)
	InferVersionFromTemplate(t, options)
	InjectTestValues(t, options)
	defer func() {
		if t.Failed() {
			DescribePods(t, namespaceName, podName)
			GetAppLog(t, namespaceName, podName, "", &corev1.PodLogOptions{})
		}
		// Delete the helper pod
		k8s.RunKubectlE(t, kubectlOptions, "delete", "pod", podName)
	}()
	nuodbImage := fmt.Sprintf(
		"%s/%s:%s",
		options.SetValues["nuodb.image.registry"],
		options.SetValues["nuodb.image.repository"],
		options.SetValues["nuodb.image.tag"])
	output, err := k8s.RunKubectlAndGetOutputE(
		t, kubectlOptions, "run", podName,
		"--image", nuodbImage, "--restart=Never", "--attach", "--command",
		"--image-pull-policy=IfNotPresent", "--",
		"nuodb", "--version")
	require.NoError(t, err, "Unable to check NuoDB version in helper pod")
	match := regexp.MustCompile("NuoDB Server build (.*)").FindStringSubmatch(output)
	require.NotNil(t, match, "Unable to match NuoDB version from output")
	nuoDBVersionMap[options.SetValues["nuodb.image.tag"]] = match[1]
	return match[1]
}

func RunOnNuoDBVersion(t *testing.T, versionCheckFunc func(*semver.Version) bool, actionFunc func(*semver.Version)) {
	RunOnNuoDBVersionFromOptions(t, &helm.Options{}, versionCheckFunc, actionFunc)
}

func RunOnNuoDBVersionFromOptions(t *testing.T, options *helm.Options, versionCheckFunc func(*semver.Version) bool, actionFunc func(*semver.Version)) {
	versionString, ok := options.SetValues["nuodb.image.tag"]
	if !ok || !NuoDBVersionRegex.MatchString(versionString) {
		versionString = GetNuoDBVersion(t, "default", options)
	}
	version, err := semver.NewVersion(NuoDBVersionRegex.ReplaceAllString(versionString, "${1}"))
	require.NoError(t, err, "Unable to parse NuoDB version", versionString)

	if versionCheckFunc(version) {
		actionFunc(version)
	} else {
		t.Logf("Conditional logic not executed for NuoDB version %s", versionString)
	}
}

func SkipTestOnNuoDBVersion(t *testing.T, versionCheckFunc func(*semver.Version) bool) {
	RunOnNuoDBVersion(t, versionCheckFunc, func(version *semver.Version) {
		t.Skip("Skipping test when using NuoDB version", version)
	})
}

/**
 * Skip test if NuoDB version matches the provided condition.
 *
 * For more information about supported condition strings, please
 * check https://github.com/Masterminds/semver.
 *
 */
func SkipTestOnNuoDBVersionCondition(t *testing.T, condition string) {
	SkipTestOnNuoDBVersion(t, func(version *semver.Version) bool {
		c, err := semver.NewConstraint(condition)
		require.NoError(t, err)
		return c.Check(version)
	})
}

func RunOnNuoDBVersionCondition(t *testing.T, condition string, actionFunc func(*semver.Version)) {
	RunOnNuoDBVersion(t, func(version *semver.Version) bool {
		c, err := semver.NewConstraint(condition)
		require.NoError(t, err)
		return c.Check(version)
	}, actionFunc)
}

func RunOnNuoDBVersionFromOptionCondition(t *testing.T, options *helm.Options, condition string, actionFunc func(*semver.Version)) {
	RunOnNuoDBVersionFromOptions(t, options, func(version *semver.Version) bool {
		c, err := semver.NewConstraint(condition)
		require.NoError(t, err)
		return c.Check(version)
	}, actionFunc)
}

func GetDatabaseProcessesE(t *testing.T, namespaceName string, adminPod string, dbName string) (processes []NuoDBProcess, err error) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", dbName)
	if err != nil {
		return nil, err
	}
	err, processes = Unmarshal(output)
	return
}

func PurgeDatabaseArchives(t *testing.T, namespaceName string, adminPod string, dbName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", dbName)
	require.NoError(t, err, output)
	err, archives := UnmarshalArchives(output)
	require.NoError(t, err, output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", dbName, "--removed")
	require.NoError(t, err, output)
	err, removed := UnmarshalArchives(output)
	require.NoError(t, err, output)
	for _, archive := range append(archives, removed...) {
		k8s.RunKubectl(t, options, "exec", adminPod, "-c", "admin", "--",
			"nuocmd", "delete", "archive", "--archive-id", strconv.Itoa(archive.Id), "--purge")
	}
}

func SuspendDatabaseBackupJobs(t *testing.T, namespaceName string, domain, dbName string, backupGroup string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	// the CronJob apiVersion may be different depending on the Kubernetes version
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "cronjobs", "-oname")
	require.NoError(t, err)
	cronJobs := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	expectedSuffix := fmt.Sprintf("-hotcopy-%s-%s-%s", domain, dbName, backupGroup)
	for _, cronJob := range cronJobs {
		cronJob = strings.TrimSpace(cronJob)
		if strings.HasSuffix(cronJob, expectedSuffix) {
			t.Logf("Suspending %s", cronJob)
			k8s.RunKubectl(t, kubectlOptions, "patch", cronJob,
				"-p", "{\"spec\" : {\"suspend\" : true }}")
			// Extract the cronjob name and delete any jobs or pods that are
			// already scheduled
			parts := strings.Split(cronJob, "/")
			require.Equal(t, 2, len(parts))
			cronJobName := parts[1]
			output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "jobs", "-o",
				fmt.Sprintf(`jsonpath={.items[?(@.metadata.ownerReferences[*].name=="%s")].metadata.name}`, cronJobName))
			require.NoError(t, err)
			output = strings.TrimSpace(output)
			if output != "" {
				for _, jobName := range strings.Split(output, " ") {
					t.Logf("Deleting job %s", jobName)
					k8s.RunKubectl(t, kubectlOptions, "delete", "job", jobName)
				}
			}
			DeleteJobPods(t, namespaceName, cronJobName)
		}
	}
}

// Create a snapshot of a pvc and wait for it to become ready
func SnapshotVolume(t *testing.T, namespaceName string, pvcName string, snapName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.KubectlApplyFromString(t, kubectlOptions, fmt.Sprintf(SNAPSHOT_TEMPLATE, snapName, VOLUME_SNAPSHOT_CLASS, pvcName))

	Await(t, func() bool {
		out, error := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "volumesnapshot", snapName, "-o", "jsonpath='{.status.readyToUse}'")
		return error == nil && strings.TrimSpace(out) == "'true'"
	}, 30*time.Second)
}

func getSnapshotName(t *testing.T, tmpl *template.Template, backupId, volumeType string) string {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &map[string]string{
		"backupId":   backupId,
		"volumeType": volumeType,
	})
	require.NoError(t, err)
	return buf.String()
}

func getArchiveSnapshotName(t *testing.T, tmpl *template.Template, backupId string) string {
	return getSnapshotName(t, tmpl, backupId, "archive")
}

func getJournalSnapshotName(t *testing.T, tmpl *template.Template, backupId string) string {
	return getSnapshotName(t, tmpl, backupId, "journal")
}

type BackupHookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func GetBackupHookResponse(t *testing.T, namespaceName, smPod, resourcePath string, payload ...string) *BackupHookResponse {
	cmdArgs := []string{
		"exec", smPod, "-c", "engine", "--",
		"curl", "-s", "-X", "POST", "http://0.0.0.0:8000/" + resourcePath,
	}
	// Prepare the payload as map
	if len(payload) > 0 {
		require.Equal(t, 0, len(payload)%2, "Backup hook payload must contain even number of entries")
		payloadMap := make(map[string]string, len(payload)/2)
		for i := 1; i < len(payload); i += 2 {
			payloadMap[payload[i-1]] = payload[i]
		}
		payloadJson, err := json.Marshal(payloadMap)
		require.NoError(t, err, payloadMap)
		cmdArgs = append(cmdArgs, "-d", string(payloadJson))
	}
	// Issue request to backup hook server. The nuodb image has curl, so use
	// it in the SM container to issue HTTP request against the backup hook
	// server running in a sidecar container on the same pod.
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, cmdArgs...)
	require.NoError(t, err, output)

	// Decode response and make sure request was successful
	buffer := bytes.NewBufferString(output)
	decoder := json.NewDecoder(buffer)
	var response BackupHookResponse
	err = decoder.Decode(&response)
	require.NoError(t, err, "Unable to decode output as JSON: %s", output)

	// Make sure output was fully decoded
	unread := strings.TrimSpace(buffer.String())
	require.Zero(t, len(unread), "Did not fully decode output: %s", output)
	return &response
}

func InvokeBackupHook(t *testing.T, namespaceName, smPod, resourcePath string, payload ...string) {
	response := GetBackupHookResponse(t, namespaceName, smPod, resourcePath, payload...)
	require.True(t, response.Success, "Unexpected response: %s", response.Message)
}

func SnapshotSm(t *testing.T, namespaceName, smPod, backupId, snapshotNameTemplate string, journalVolume bool) {
	// Invoke pre-backup hook
	InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)

	// Parse snapshot name template
	tmpl, err := template.New("").Parse(snapshotNameTemplate)
	require.NoError(t, err)

	// Snapshot archive
	achiveVolumeName := "archive-volume-" + smPod
	archiveSnapshotName := getArchiveSnapshotName(t, tmpl, backupId)
	SnapshotVolume(t, namespaceName, achiveVolumeName, archiveSnapshotName)

	// Snapshot journal if journal volume is separate
	if journalVolume {
		journalVolumeName := "journal-volume-" + smPod
		journalSnapshotName := getJournalSnapshotName(t, tmpl, backupId)
		SnapshotVolume(t, namespaceName, journalVolumeName, journalSnapshotName)
	}

	// Invoke post-backup hook
	InvokeBackupHook(t, namespaceName, smPod, "post-backup/"+backupId)
}
