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
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	v12 "k8s.io/api/core/v1"
)

func upgradeAdminTest(t *testing.T, fromHelmVersion string, upgradeOptions *testlib.UpgradeOptions) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.readinessTimeoutSeconds": "5",
		},
		Version: fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%supgradeadmintest-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because the older versions of helm charts have
	// hardcoded instances of "https://" in LB policy job and NuoDB 4.2+ image
	// doesn't contain pregenerated keys
	testlib.GenerateAndSetTLSKeys(t, options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	adminPod := testlib.GetPod(t, namespaceName, admin0)

	// unset the version and use local
	options.Version = ""
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	helm.Upgrade(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	if upgradeOptions.AdminPodShouldGetRecreated {
		testlib.AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
	}

	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)
}

func upgradeDatabaseTest(t *testing.T, fromHelmVersion string, upgradeOptions *testlib.UpgradeOptions) {
	if os.Getenv("NUODB_LICENSE") == "ENTERPRISE" {
		t.Skip("Can not test helm upgrade in this environment. See DB-33858")
	}

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.readinessTimeoutSeconds":         "5",
			"database.sm.resources.requests.cpu":    "250m",
			"database.sm.resources.requests.memory": "250Mi",
			"database.te.resources.requests.cpu":    "250m",
			"database.te.resources.requests.memory": "250Mi",
			"nuodb.image.pullPolicy":                "IfNotPresent",
		},
		Version: fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%supgradedatabasetest-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because the older versions of helm charts have
	// hardcoded instances of "https://" in LB policy job and NuoDB 4.2+ image
	// doesn't contain pregenerated keys
	testlib.GenerateAndSetTLSKeys(t, options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)

	adminPod := testlib.GetPod(t, namespaceName, admin0)

	// unset the version and use local
	options.Version = ""
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	helm.Upgrade(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	if upgradeOptions.AdminPodShouldGetRecreated {
		testlib.AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
	}

	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	opt := testlib.GetExtractedOptions(options)

	// make sure the DB is properly reconnected before restarting
	err := testlib.AwaitE(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Reconnected with process with connectKey",
			&v12.PodLogOptions{
				Container: "admin",
			}) == 2
	}, 120*time.Second)

	if err != nil {
		// in some environments, the engine does not manage to reconnect to an admin after the admin pod was restarted
		// the only way to proceed is to restart the engine
		// see https://github.com/nuodb/nuodb-helm-charts/issues/238
		t.Log(err)
		t.Logf("WARNING: engine did not reconnect with admin. Killing all engines to make upgrade proceed!")

		opt := testlib.GetExtractedOptions(options)
		tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)

		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

		go testlib.GetAppLog(t, namespaceName, tePodName, "-pre-kill", &v12.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodName, "-pre-kill", &v12.PodLogOptions{Follow: true})

		testlib.KillProcess(t, namespaceName, tePodName)
		testlib.KillProcess(t, namespaceName, smPodName)

		testlib.AwaitPodUp(t, namespaceName, tePodName, 300*time.Second)
		testlib.AwaitPodUp(t, namespaceName, smPodName, 300*time.Second)

	}
	// make sure the environment is stable before proceeding
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	testlib.UpgradeDatabase(t, namespaceName, databaseReleaseName, admin0, options, upgradeOptions)
}

func TestUpgradeHelm(t *testing.T) {
	t.Run("NuoDB_From310_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.1.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})

	t.Run("NuoDB_From320_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.2.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})

	t.Run("NuoDB_From330_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.3.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})
}

func TestUpgradeHelmFullDB(t *testing.T) {
	t.Run("NuoDB_From310_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.1.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})

	t.Run("NuoDB_From320_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.2.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})

	t.Run("NuoDB_From330_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.3.0", &testlib.UpgradeOptions{
			AdminPodShouldGetRecreated: true,
		})
	})
}
