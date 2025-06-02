//go:build upgrade
// +build upgrade

package minikube

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

const OLD_RELEASE = "5.0"

func verifyAllProcessesRunning(t *testing.T, namespaceName string, adminPod string, expectedNrProcesses int) {
	testlib.Await(t, func() bool {
		options := k8s.NewKubectlOptions("", "", namespaceName)

		output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--", "nuocmd", "show", "domain")
		require.NoError(t, err, "verifyAllProcessesRunning: running show domain failed")

		return strings.Count(output, "MONITORED:RUNNING") == expectedNrProcesses
	}, 30*time.Second)
}

func TestAdminProbes(t *testing.T) {
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	// create a two-server domain and induce a failure that makes it
	// impossible to elect a leader, causing 'nuocmd check server
	// --check-converged' to fail
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"admin.livenessProbe.initialDelaySeconds": "30",
			"admin.livenessProbe.periodSeconds":       "15",
			"admin.livenessProbe.failureThreshold":    "1",
			"admin.readinessProbe.failureThreshold":   "1",
			"admin.replicas":                          "2",
		},
	}, 2, "")
	adminStatefulSet := helmChartReleaseName + "-nuodb-cluster0"
	admin0 := adminStatefulSet + "-0"
	admin1 := adminStatefulSet + "-1"

	// make sure both Admin Pods become Ready
	testlib.AwaitPodUp(t, namespaceName, admin0, 120*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin1, 120*time.Second)

	// make sure direct invocation of readinessprobe script succeeds on both
	// Admin processes
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "-c", "admin", "--", "readinessprobe")
	require.NoError(t, err, "readinessprobe failed: %s", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "-c", "admin", "--", "readinessprobe")
	require.NoError(t, err, "readinessprobe failed: %s", output)

	// make sure direct invocation of livenessprobe script succeeds on both
	// Admin processes
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "-c", "admin", "--", "livenessprobe")
	require.NoError(t, err, "livenessprobe failed: %s", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "-c", "admin", "--", "livenessprobe")
	require.NoError(t, err, "livenessprobe failed: %s", output)

	// get restart count for admin-1
	pod, err := testlib.FindPod(t, namespaceName, admin1)
	require.NoError(t, err)
	var restartCountBefore int32 = -1
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == "admin" {
			restartCountBefore = containerStatus.RestartCount
		}
	}
	require.NotEqual(t, -1, restartCountBefore, "Unable to get restart count for admin container")

	// delete raftlog PVC and pod for admin-0 so that it bootstraps a
	// disjoint domain when it is restarted and refuses messages from
	// admin-1
	k8s.RunKubectl(t, options, "delete", "pvc/raftlog-"+admin0, "pod/"+admin0)

	// check status for admin-1 pod
	testlib.Await(t, func() bool {
		pod, err := testlib.FindPod(t, namespaceName, admin1)
		require.NoError(t, err)
		var foundReadyCondition, foundRestartCount bool
		// check that Ready status condition is false
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != corev1.ConditionFalse {
					return false
				}
				foundReadyCondition = true
			}
		}
		// check that restart count for admin container has increased
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "admin" {
				if restartCountBefore >= containerStatus.RestartCount {
					return false
				}
				foundRestartCount = true
			}
		}
		return foundReadyCondition && foundRestartCount
	}, 120*time.Second)
}

func TestKubernetesUpgradeAdminMinorVersion(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb",
			"nuodb.image.tag":        OLD_RELEASE,
			"admin.bootstrapServers": "0",
		},
	}
	testlib.OverrideUpgradeContainerImage(t, &options)

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})

	expectedNewVersion := testlib.GetUpgradedReleaseVersion(t, &options)

	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}

func TestKubernetesUpgradeFullDatabase(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"nuodb.image.registry":                  "docker.io",
			"nuodb.image.repository":                "nuodb/nuodb",
			"nuodb.image.tag":                       OLD_RELEASE,
			"admin.bootstrapServers":                "0",
			"database.sm.resources.requests.cpu":    "0.25",
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "0.25",
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	}
	testlib.OverrideUpgradeContainerImage(t, &options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%supgradefulldatabase-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because NuoDB 4.2+ doesn't include
	// pre-generated keys and the engines can't reconnect to the upgraded admin tier
	testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	adminHelmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminHelmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	databaseHelmChartReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	expectedNewVersion := testlib.GetUpgradedReleaseVersion(t, &options)

	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, adminHelmChartReleaseName)

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)
	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })

	opt := testlib.GetExtractedOptions(&options)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	t.Run("expectAllEnginesReconnect", func(t *testing.T) {
		expectedNumberReconnects := 2

		testlib.Await(t, func() bool {
			return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0, "Reconnected with process with connectKey", &corev1.PodLogOptions{}) == expectedNumberReconnects
		}, 30*time.Second)

	})

	t.Run("verifyAllProcessesRunning", func(t *testing.T) {
		verifyAllProcessesRunning(t, namespaceName, admin0, 2)
	})

	t.Run("upgradeDatabaseHelm", func(t *testing.T) {
		expectedNewDatabaseVersion := testlib.GetUpgradedReleaseVersion(t, &options)

		helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseHelmChartReleaseName)

		// make sure that we only have 1 TE and not 2
		testlib.SetDeploymentUpgradeStrategyToRecreate(t, namespaceName, fmt.Sprintf("te-%s-nuodb-cluster0-demo", databaseHelmChartReleaseName))

		testlib.AwaitPodTemplateHasVersion(t, namespaceName, "sm-database", expectedNewDatabaseVersion, 300*time.Second)
		testlib.AwaitPodTemplateHasVersion(t, namespaceName, "te-database", expectedNewDatabaseVersion, 300*time.Second)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		verifyAllProcessesRunning(t, namespaceName, admin0, 2)
	})

	// In order to have a protocol upgrade, there has to be a protocol change between
	// OLD_RELEASE and the version under test.
	testlib.RunOnNuoDBVersionCondition(t, ">6.0.0", func(version *semver.Version) {
		// check that KAA will upgrade database protocol version and restart TE
		// automatically
		t.Run("verifyProtocolVersion", func(t *testing.T) {
			// start two TEs so that we can supply TE preference query
			options.SetValues["database.te.replicas"] = "2"
			helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseHelmChartReleaseName)

			testlib.Await(t, func() bool {
				// It is possible that the database will (briefly) have 3 processes while the upgrade is in progress
				// so we need to check that both the database is in the expected state and that there are not
				// leftover pods.
				err := k8s.RunKubectlE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "check", "database",
					"--db-name", opt.DbName, "--check-running", "--check-liveness", "20",
					"--num-processes", strconv.Itoa(3))
				return err == nil && len(testlib.GetPodNames(t, namespaceName, databaseHelmChartReleaseName)) == 3
			}, 300*time.Second)

			// find the startId of the second TE
			var toShutdownStartId int64 = -1
			toShutdownPod := ""
			processes, err := testlib.GetDatabaseProcessesE(t, namespaceName, admin0, opt.DbName)
			require.NoError(t, err)
			for _, process := range processes {
				if process.Type == "TE" {
					startId, err := strconv.ParseInt(process.StartId, 10, 64)
					require.NoError(t, err)
					if startId > toShutdownStartId {
						toShutdownStartId = startId
						toShutdownPod = process.Hostname
					}
				}
			}
			require.NotEqual(t, -1, toShutdownStartId, "Unable to find TE to shutdown")
			// enable automatic protocol upgrade
			options.SetValues["database.automaticProtocolUpgrade.enabled"] = "true"
			// supply LB query which will select the second TE to be shutdown
			// after protocol upgrade
			options.SetValues["database.automaticProtocolUpgrade.tePreferenceQuery"] =
				fmt.Sprintf("random(start_id(%d))", toShutdownStartId)

			helm.Upgrade(t, &options, testlib.DATABASE_HELM_CHART_PATH, databaseHelmChartReleaseName)

			// protocol upgrade is async task
			testlib.Await(t, func() bool {
				return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
					fmt.Sprintf("Upgrading protocol version for database dbName=%s", opt.DbName),
					&corev1.PodLogOptions{}) >= 1
			}, 60*time.Second)

			// fetch database effective version
			output, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions, "exec", admin0, "--",
				"nuocmd", "show", "database-versions", "--db-name", opt.DbName)
			require.NoError(t, err, "running show database-versions failed")
			pattern := regexp.MustCompile("effective version ID: ([0-9]+), effective version: (.+),")
			match := pattern.FindStringSubmatch(output)
			require.NotNil(t, match, "Unable to get database effective version from output")

			// effective version is string like this one 4.2|4.2.1|4.2.2; verify
			// that the OLD_RELEASE major.minor is not in the effective version
			// string
			parts := strings.Split(OLD_RELEASE, ".")
			require.GreaterOrEqual(t, len(parts), 2, "unable to get major.minor from OLD_RELEASE version")
			require.NotContains(t, match[2], fmt.Sprintf("%s.%s", parts[0], parts[1]))
			// verify that the new major.minor is in the effective version
			// string
			require.Contains(t, match[2], fmt.Sprintf("%d.%d", version.Major(), version.Minor()))

			// verify that a TE has been restarted and SQL layer version is
			// upgrade is performed
			testlib.Await(t, func() bool {
				return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
					fmt.Sprintf("Shutting down startId=%d to finalize the database protocol upgrade", toShutdownStartId),
					&corev1.PodLogOptions{}) >= 1
			}, 60*time.Second)
			testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, toShutdownPod, 0, 30*time.Second)
			testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, 3)
			output, err = testlib.RunSQL(t, namespaceName, admin0, opt.DbName,
				"select version from system.versions"+
					" where property = 'SYSTEM_TABLES_VERSION'"+
					" and version = GETEFFECTIVEPLATFORMVERSION()")
			require.NoError(t, err, "error cheking if SQL layer version is upgraded")
			require.Contains(t, output, match[1])

			// verify that the container image ID is stored in KV store so that
			// database version is not checked again for this container image
			// ID; this happens on the next resync once database versions
			// response is inspected
			testlib.Await(t, func() bool {
				adminPod := testlib.GetPod(t, namespaceName, admin0)
				actualImageId, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions, "exec", admin0, "-c", "admin", "--",
					"nuocmd", "get", "value", "--key", fmt.Sprintf("upgradeDatabaseLastObservedImage/%s", opt.DbName))
				require.NoError(t, err, "error getting last observed image ID")
				return adminPod.Status.ContainerStatuses[0].ImageID == actualImageId
			}, 60*time.Second)
		})
	})

}

func TestKubernetesRollingUpgradeAdminMinorVersion(t *testing.T) {
	t.Skip("4.0.7+ Admin is not rolling upgradeable from pre-4.0.7")

	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas":         "3",
			"admin.bootstrapServers": "0",
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb",
			"nuodb.image.tag":        OLD_RELEASE,
		},
	}
	testlib.OverrideUpgradeContainerImage(t, &options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%srollingupgradeadminminorversion-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because the older versions of helm charts have
	// hardcodded instances of "https://" in LB policy job and NuoDB 4.2+ image
	// doesn't contain pregenerated keys
	testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 3, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-cluster0-1", helmChartReleaseName)
	admin2 := fmt.Sprintf("%s-nuodb-cluster0-2", helmChartReleaseName)

	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, admin1, "-previous", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, admin2, "-previous", &corev1.PodLogOptions{Follow: true})

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")

	expectedNewVersion := testlib.GetUpgradedReleaseVersion(t, &options)

	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	// the rolling upgrade is done in reverse order
	testlib.AwaitPodHasVersion(t, namespaceName, admin2, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin2, 300*time.Second)

	testlib.AwaitPodHasVersion(t, namespaceName, admin1, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin1, 300*time.Second)

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}
