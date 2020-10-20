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
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	v12 "k8s.io/api/core/v1"
)

type UpdateOptions struct {
	adminPodShouldGetRecreated            bool
	adminRolesRequirePatching             bool
	adminBootstrapServersOverrideRequired bool
}

func upgradeAdminTest(t *testing.T, fromHelmVersion string, updateOptions *UpdateOptions) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{},
		Version:   fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)

	if updateOptions.adminBootstrapServersOverrideRequired {
		options.SetValues["admin.bootstrapServers"] = "0"
	}
	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("upgradeadmintest-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because the older versions of helm charts have
	// hardcodded instances of "https://" in LB policy job and NuoDB 4.2+ image
	// doesn't contain pregenerated keys
	testlib.GenerateAndSetTLSKeys(t, options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	if updateOptions.adminRolesRequirePatching {
		testlib.AdminRolesRequirePatching = true
		testlib.AddTeardown(testlib.TEARDOWN_ADMIN, func() {
			testlib.AdminRolesRequirePatching = false
		})
	}

	helmChartReleaseName, _ := testlib.StartAdmin(t, options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")

	adminPod := testlib.GetPod(t, namespaceName, admin0)

	// unset the version and use local
	options.Version = ""
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	helm.Upgrade(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	if updateOptions.adminPodShouldGetRecreated {
		testlib.AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
	}

	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)
}

func upgradeDatabaseTest(t *testing.T, fromHelmVersion string, updateOptions *UpdateOptions) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.bootstrapServers":                "0",
			"database.sm.resources.requests.cpu":    "250m",
			"database.sm.resources.requests.memory": "250Mi",
			"database.te.resources.requests.cpu":    "250m",
			"database.te.resources.requests.memory": "250Mi",
		},
		Version: fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)

	if updateOptions.adminBootstrapServersOverrideRequired {
		options.SetValues["admin.bootstrapServers"] = "0"
	}
	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("upgradedatabasetest-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// Enable TLS during upgrade because the older versions of helm charts have
	// hardcodded instances of "https://" in LB policy job and NuoDB 4.2+ image
	// doesn't contain pregenerated keys
	testlib.GenerateAndSetTLSKeys(t, options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	if updateOptions.adminRolesRequirePatching {
		testlib.AdminRolesRequirePatching = true
		testlib.AddTeardown(testlib.TEARDOWN_ADMIN, func() {
			testlib.AdminRolesRequirePatching = false
		})
	}

	helmChartReleaseName, _ := testlib.StartAdmin(t, options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)

	testlib.AwaitBalancerTerminated(t, namespaceName, "job-lb-policy")

	// all jobs need to be deleted before an upgrade can be performed
	// so far we have not found an automated way to delete them as part of a pre-upgrade hook
	// if we find it, this line can be removed and the test should still pass
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")
	testlib.DeletePod(t, namespaceName, "jobs/hotcopy-demo-job-initial")

	adminPod := testlib.GetPod(t, namespaceName, admin0)

	// unset the version and use local
	options.Version = ""
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts

	helm.Upgrade(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	if updateOptions.adminPodShouldGetRecreated {
		testlib.AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
	}

	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	opt := testlib.GetExtractedOptions(options)

	// make sure the DB is properly reconnected before restarting
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	helm.Upgrade(t, options, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
}

func TestUpgradeHelm(t *testing.T) {
	t.Run("NuoDB40X_From231_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "2.3.1", &UpdateOptions{
			adminPodShouldGetRecreated:            true,
			adminRolesRequirePatching:             true,
			adminBootstrapServersOverrideRequired: true,
		})
	})

	t.Run("NuoDB40X_From240_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "2.4.0", &UpdateOptions{
			adminRolesRequirePatching: true,
		})
	})

	t.Run("NuoDB40X_From241_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "2.4.1", &UpdateOptions{})
	})
}

func TestUpgradeHelmFullDB(t *testing.T) {
	t.Run("NuoDB40X_From231_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "2.3.1", &UpdateOptions{
			adminPodShouldGetRecreated:            true,
			adminRolesRequirePatching:             true,
			adminBootstrapServersOverrideRequired: true,
		})
	})

	t.Run("NuoDB40X_From240_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "2.4.0", &UpdateOptions{
			adminRolesRequirePatching:  true,
			adminPodShouldGetRecreated:	true,
		})
	})

	t.Run("NuoDB40X_From241_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "2.4.1", &UpdateOptions{
			adminPodShouldGetRecreated:	true,
		})
	})
}
