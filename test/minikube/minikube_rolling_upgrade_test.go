// +build upgrade

package minikube

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/require"

	v12 "k8s.io/api/core/v1"
)

const OLD_RELEASE = "4.0.4"

func verifyAllProcessesRunning(t *testing.T, namespaceName string, adminPod string, expectedNrProcesses int) {
	testlib.Await(t, func() bool {
		options := k8s.NewKubectlOptions("", "", namespaceName)

		output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--", "nuocmd", "show", "domain")
		require.NoError(t, err, "verifyAllProcessesRunning: running show domain failed")

		return strings.Count(output, "MONITORED:RUNNING") == expectedNrProcesses
	}, 30*time.Second)
}

func TestAdminReadinessProbe(t *testing.T) {
	// TODO: remove this whenever the image tested in nuodb-helm-charts CI
	// supports 'nuocmd check server' (singular)
	if os.Getenv("NUODB_DEV") != "true" {
		t.Skip("'nuocmd check server' is not supported in released versions")
	}

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	// create a two-server domain and induce a failure that makes it
	// impossible to elect a leader, causing 'nuocmd check server
	// --check-converged' to fail
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"admin.replicas": "2",
		},
	}, 2, "")
	adminStatefulSet := helmChartReleaseName + "-nuodb-cluster0"
	admin0 := adminStatefulSet + "-0"
	admin1 := adminStatefulSet + "-1"

	// make sure readinessprobe succeeds on both Admin processes
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "--", "readinessprobe")
	require.NoError(t, err, "readinessprobe failed: %s", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "--", "readinessprobe")
	require.NoError(t, err, "readinessprobe failed: %s", output)

	leader, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "--", "bash", "-c",
		"nuocmd show domain --server-format '{role} {id}' | sed -n 's/ *LEADER //p'")
	require.NoError(t, err)
	// delete Raft log on admin-0 and kill admin-0 so that it bootstraps a
	// disjoint domain when it is restarted and refuses messages from
	// admin-1
	k8s.RunKubectl(t, options, "exec", admin0, "--", "rm", "-f", "/var/opt/nuodb/raftlog")
	k8s.RunKubectl(t, options, "delete", "pod", admin0)
	// also kill admin-1 if it was the leader so that a new leader has to be
	// elected when it is restarted
	if admin1 == strings.TrimSpace(leader) {
		k8s.RunKubectl(t, options, "delete", "pod", admin1)
	}

	// make sure readinessprobe on admin-1 eventually fails, because there
	// is no leader for it to converge with even if the backwards-compatible
	// 'nuocmd check servers' (plural) succeeds; this verifies that the
	// fallback behavior does not apply to actual failures of the readiness
	// check
	testlib.Await(t, func() bool {
		// make sure 'nuocmd check servers' (plural) is successful, since it
		// only checks that the --api-server is ACTIVE
		output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "--", "nuocmd", "check", "servers")
		if err != nil {
			return false
		}

		// make sure readinessprobe fails
		output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "--", "readinessprobe")
		// make sure exit code is 1 to indicate non-parse error
		code, err := shell.GetExitCodeForRunCommandError(err)
		require.NoError(t, err)
		require.NotEqual(t, 2, code, "Exit code 2 should be reserved for parse errors")
		return code == 1
	}, 60*time.Second)

	// make sure 'nuocmd check servers' (plural) is successful, since it
	// only checks that the --api-server is ACTIVE
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "--", "nuocmd", "check", "servers")
	require.NoError(t, err, "'nuocmd check servers' failed: %s", output)
}

func TestAdminReadinessProbeFallback(t *testing.T) {
	// 'nuomcd check server' (singular) is unsupported for versions <=4.1.1;
	// this test verifies the fallback behavior of the readiness probe

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	// this test uses OLD_RELEASE, but any version <=4.1.1 will do
	helmOptions := helm.Options{
		SetValues: map[string]string{
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb-ce",
			"nuodb.image.tag":        OLD_RELEASE,
			"admin.bootstrapServers": "0",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helmOptions, 1, "")
	admin := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// make sure 'nuocmd check server' fails
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin, "--", "nuocmd", "check", "server")
	require.Error(t, err, "Expected 'nuocmd check server' to fail on %s: %s", OLD_RELEASE, output)
	// make sure exit code is 2 to indicate parse error
	code, err := shell.GetExitCodeForRunCommandError(err)
	require.NoError(t, err)
	require.Equal(t, 2, code)

	// make sure readinessprobe is successful
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin, "--", "readinessprobe")
	require.NoError(t, err, "readinessprobe failed on %s: %s", OLD_RELEASE, output)
}

func TestKubernetesUpgradeAdminMinorVersion(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb-ce",
			"nuodb.image.tag":        OLD_RELEASE,
			"admin.bootstrapServers": "0",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")

	expectedNewVersion := testlib.GetUpgradedReleaseVersion(t, &options)

	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, fmt.Sprintf(expectedNewVersion), 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}

func TestKubernetesUpgradeFullDatabaseMinorVersion(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb-ce",
			"nuodb.image.tag":        OLD_RELEASE,
			"admin.bootstrapServers": "0",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	adminHelmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminHelmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"nuodb.image.registry":                  "docker.io",
			"nuodb.image.repository":                "nuodb/nuodb-ce",
			"nuodb.image.tag":                       OLD_RELEASE,
		},
	}

	databaseHelmChartReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")
	testlib.DeletePod(t, namespaceName, "jobs/hotcopy-demo-job-initial")

	expectedNewVersion := testlib.GetUpgradedReleaseVersion(t, &options)

	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, adminHelmChartReleaseName)

	testlib.AwaitPodHasVersion(t, namespaceName, admin0, expectedNewVersion, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })

	opt := testlib.GetExtractedOptions(&databaseOptions)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	t.Run("expectAllEnginesReconnect", func(t *testing.T) {
		expectedNumberReconnects := 2

		testlib.Await(t, func() bool {
			return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0, "Reconnected with process with connectKey", &v12.PodLogOptions{}) == expectedNumberReconnects
		}, 30*time.Second)

	})

	t.Run("verifyAllProcessesRunning", func(t *testing.T) {
		verifyAllProcessesRunning(t, namespaceName, admin0, 2)
	})

	t.Run("upgradeDatabaseHelm", func(t *testing.T) {
		expectedNewDatabaseVersion := testlib.GetUpgradedReleaseVersion(t, &databaseOptions)

		helm.Upgrade(t, &databaseOptions, testlib.DATABASE_HELM_CHART_PATH, databaseHelmChartReleaseName)

		// make sure that we only have 1 TE and not 2
		testlib.SetDeploymentUpgradeStrategyToRecreate(t, namespaceName, fmt.Sprintf("te-%s-nuodb-cluster0-demo", databaseHelmChartReleaseName))

		testlib.AwaitPodTemplateHasVersion(t, namespaceName, "sm-database", expectedNewDatabaseVersion, 300*time.Second)
		testlib.AwaitPodTemplateHasVersion(t, namespaceName, "te-database", expectedNewDatabaseVersion, 300*time.Second)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

		verifyAllProcessesRunning(t, namespaceName, admin0, 2)
	})
}

func TestKubernetesRollingUpgradeAdminMinorVersion(t *testing.T) {
	t.Skip("4.0.7 Admin is not rolling upgradeable from pre-4.0.7")

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas":         "3",
			"admin.bootstrapServers": "0",
			"nuodb.image.registry":   "docker.io",
			"nuodb.image.repository": "nuodb/nuodb-ce",
			"nuodb.image.tag":        OLD_RELEASE,
		},
	}

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("rollingupgradeadminminorversion-%s", randomSuffix)
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

	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, admin1, "-previous", &v12.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, admin2, "-previous", &v12.PodLogOptions{Follow: true})

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
