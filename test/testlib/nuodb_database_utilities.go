package testlib

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/core/v1"

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

type ExtractedOptions struct {
	NrTePods          int
	NrSmHotCopyPods   int
	NrSmNoHotCopyPods int
	NrSmPods          int
	DbName            string
	ClusterName       string
}

func GetExtractedOptions(options *helm.Options) (opt ExtractedOptions) {
	var err error

	opt.NrTePods, err = strconv.Atoi(options.SetValues["database.te.replicas"])
	if err != nil {
		opt.NrTePods = 1
	}

	opt.NrSmHotCopyPods, err = strconv.Atoi(options.SetValues["database.sm.hotCopy.replicas"])
	if err != nil {
		opt.NrSmHotCopyPods = 1
	}

	opt.NrSmNoHotCopyPods, err = strconv.Atoi(options.SetValues["database.sm.noHotCopy.replicas"])
	if err != nil {
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

	return
}

func EnsureDatabaseNotRunning(t *testing.T, adminPod string, opt ExtractedOptions, kubectlOptions *k8s.KubectlOptions) {
	// invoke shutdown database; this may fail if the database is already NOT_RUNNING, which is okay
	k8s.RunKubectlE(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "shutdown", "database", "--db-name", opt.DbName)
	// wait for all database processes to exit
	k8s.RunKubectl(t, kubectlOptions, "exec", adminPod, "--", "nuocmd", "check", "database", "--db-name", opt.DbName, "--num-processes", "0", "--timeout", "30")
}

type DatabaseInstallationStep func(t *testing.T, options *helm.Options, helmChartReleaseName string)

func StartDatabaseTemplate(t *testing.T, namespaceName string, adminPod string, options *helm.Options, installationStep DatabaseInstallationStep, awaitDatabase bool) (helmChartReleaseName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

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
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)
	smPodName := fmt.Sprintf("sm-%s-nuodb-%s-%s", helmChartReleaseName, opt.ClusterName, opt.DbName)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	// with Async actions which do not return a cleanup method, create the teardown(s) first
	AddTeardown(TEARDOWN_DATABASE, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
		AwaitNoPods(t, namespaceName, helmChartReleaseName)
		EnsureDatabaseNotRunning(t, adminPod, opt, kubectlOptions)
		DeleteDatabase(t, namespaceName, opt.DbName, adminPod)
	})

	installationStep(t, options, helmChartReleaseName)

	if awaitDatabase {
		AwaitNrReplicasScheduled(t, namespaceName, tePodNameTemplate, opt.NrTePods)
		AwaitNrReplicasScheduled(t, namespaceName, smPodName, opt.NrSmPods)

		// NOTE: the Teardown logic will pick a TE/SM that is running during teardown time. Not the TE/SM that was running originally
		// this is relevant for any tests that restart TEs/SMs

		tePodName := GetPodName(t, namespaceName, tePodNameTemplate)

		AddTeardown(TEARDOWN_DATABASE, func() {
			go GetAppLog(t, namespaceName, GetPodName(t, namespaceName, tePodNameTemplate), "", &v12.PodLogOptions{Follow: true})
		})
		AwaitPodUp(t, namespaceName, tePodName, 180*time.Second)

		smPodName0 := GetPodName(t, namespaceName, smPodName)
		AddTeardown(TEARDOWN_DATABASE, func() {
			go GetAppLog(t, namespaceName, GetPodName(t, namespaceName, smPodName), "", &v12.PodLogOptions{Follow: true})
		})
		AwaitPodUp(t, namespaceName, smPodName0, 240*time.Second)

		AwaitDatabaseUp(t, namespaceName, adminPod, opt.DbName, opt.NrSmPods+opt.NrTePods)
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
	for key := range options.SetValues {
		if value, ok := databaseOptions.SetValues[key]; ok {
			options.SetValues[key] = value
		}
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	restore := func() {
		// Remove restore job if exist as it's not unique for a restore chart release
		k8s.RunKubectlE(t, kubectlOptions, "delete", "job", "restore-"+options.SetValues["database.name"])
		InjectTestValues(t, options)
		helm.Install(t, options, RESTORE_HELM_CHART_PATH, restName)
		AddTeardown(TEARDOWN_RESTORE, func() { helm.Delete(t, options, restName, true) })
		AwaitPodPhase(t, namespaceName, "restore-demo-", corev1.PodSucceeded, 120*time.Second)
	}

	if databaseOptions.SetValues["restore.autoRestart"] == "true" {
		AwaitDatabaseRestart(t, namespaceName, podName, "demo", databaseOptions, restore)
	} else {
		restore()
	}
}

func BackupDatabase(t *testing.T, namespaceName string, podName string,
	databaseName string, backupType string, backupGroup string) string {
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, opts,
		"exec", podName, "--",
		"nuobackup", "--type", backupType, "--db-name", databaseName,
		"--group", backupGroup, "--backup-root", "/var/opt/nuodb/backup",
	)
	require.NoError(t, err, "Error creating backup")
	require.True(t, strings.Contains(output, "completed"), "Error nuobackup: %s", output)
	return GetLatestBackup(t, namespaceName, podName, databaseName, backupGroup)
}

func GetLatestBackup(t *testing.T, namespaceName string, podName string,
	databaseName string, backupGroup string) string {
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	backupset, err := k8s.RunKubectlAndGetOutputE(t, opts,
		"exec", podName, "--", "bash", "-c",
		"nuobackup --type report-latest --db-name "+databaseName+
			" --group "+backupGroup+" --backup-root /var/opt/nuodb/backup 2>/dev/null",
	)
	require.NoError(t, err, "Error while reporting latest backupset")
	require.True(t, backupset != "")
	return backupset
}
