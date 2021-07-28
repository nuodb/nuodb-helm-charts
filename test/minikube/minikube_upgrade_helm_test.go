// +build upgrade

package minikube

import (
	"fmt"
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
		SetValues: map[string]string{},
		Version:   fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("upgradeadmintest-%s", randomSuffix)
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

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
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
	namespaceName := fmt.Sprintf("upgradedatabasetest-%s", randomSuffix)
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

	// TODO: check that 'helm upgrade' did not trigger a restart of database
	// pods and that the database is still running; currently it is possible
	// that the restarted admin does not receive a reconnection from the
	// database processes either due to DNS resolution or due to the
	// database process not detecting the loss of the admin connection

	testlib.UpgradeDatabase(t, namespaceName, databaseReleaseName, admin0, options, upgradeOptions)

	opt := testlib.GetExtractedOptions(options)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
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
}
