//go:build upgrade
// +build upgrade

package minikube

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func upgradeAdminTest(t *testing.T, fromHelmVersion string, upgradeOptions *testlib.UpgradeOptions) {
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
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})

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

func upgradeDatabaseTest(t *testing.T, fromHelmVersion string, enableCron bool, upgradeOptions *testlib.UpgradeOptions) {
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.readinessTimeoutSeconds":         "5",
			"database.sm.resources.requests.cpu":    "250m",
			"database.sm.resources.requests.memory": "250Mi",
			"database.te.resources.requests.cpu":    "250m",
			"database.te.resources.requests.memory": "250Mi",
			"nuodb.image.pullPolicy":                "IfNotPresent",
			"database.sm.hotCopy.enableBackups":     strconv.FormatBool(enableCron),
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
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})

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
			&corev1.PodLogOptions{
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

		go testlib.GetAppLog(t, namespaceName, tePodName, "-pre-kill", &corev1.PodLogOptions{Follow: true})
		go testlib.GetAppLog(t, namespaceName, smPodName, "-pre-kill", &corev1.PodLogOptions{Follow: true})

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

	t.Run("NuoDB_From370_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.7.0", &testlib.UpgradeOptions{})
	})

	t.Run("NuoDB_From382_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.8.2", &testlib.UpgradeOptions{})
	})

	t.Run("NuoDB_From390_ToLocal", func(t *testing.T) {
		upgradeAdminTest(t, "3.9.0", &testlib.UpgradeOptions{})
	})
}

func TestUpgradeHelmFullDB(t *testing.T) {

	t.Run("NuoDB_From370_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.7.0", false, &testlib.UpgradeOptions{})
	})

	t.Run("NuoDB_From382_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.8.2", true, &testlib.UpgradeOptions{})
	})

	t.Run("NuoDB_From390_ToLocal", func(t *testing.T) {
		upgradeDatabaseTest(t, "3.9.0", true, &testlib.UpgradeOptions{})
	})
}

func TestCredentialImport(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	fromHelmVersion := "3.9.0"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.readinessTimeoutSeconds":         "5",
			"database.sm.resources.requests.cpu":    "250m",
			"database.sm.resources.requests.memory": "250Mi",
			"database.te.resources.requests.cpu":    "250m",
			"database.te.resources.requests.memory": "250Mi",
			"nuodb.image.pullPolicy":                "IfNotPresent",
			"database.generatePassword.enabled":     "true",
			"database.rootUser":                     "someUser",
		},
		Version: fromHelmVersion,
	}
	testlib.InferVersionFromTemplate(t, options)
	opt := testlib.GetExtractedOptions(options)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%supgradedatabasetest-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)

	oldUser, oldPass := testlib.GetDatabaseCredentials(t, namespaceName, opt.DomainName, opt.DbName)

	// unset the version and use local
	options.Version = ""
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = opts
	options.SetValues["database.rootUser"] = ""

	helm.Upgrade(t, options, testlib.ADMIN_HELM_CHART_PATH, helmChartReleaseName)

	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	// make sure the environment is stable before proceeding
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	testlib.UpgradeDatabase(t, namespaceName, databaseReleaseName, admin0, options, &testlib.UpgradeOptions{})

	newUser, newPass := testlib.GetDatabaseCredentials(t, namespaceName, opt.DomainName, opt.DbName)

	require.Equal(t, oldUser, newUser)
	require.Equal(t, oldPass, newPass)

	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

	// explicitly change the password
	rotatedPassword := "SomethingNew"
	rotatedUser := "NewUser"
	delete(options.SetValues, "database.generatePassword.enabled")
	options.SetValues["database.rootPassword"] = rotatedPassword
	options.SetValues["database.rootUser"] = rotatedUser

	testlib.UpgradeDatabase(t, namespaceName, databaseReleaseName, admin0, options, &testlib.UpgradeOptions{})

	newUser, newPass = testlib.GetDatabaseCredentials(t, namespaceName, opt.DomainName, opt.DbName)

	require.Equal(t, rotatedPassword, newPass)
	require.Equal(t, rotatedUser, newUser)
}
