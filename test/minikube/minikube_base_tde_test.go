// +build long

package minikube

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func applyStoragePasswordSecret(t *testing.T, namespaceName string, name string, passwords []string) {
	assert.Greater(t, len(passwords), 0)
	opts := k8s.NewKubectlOptions("", "", namespaceName)
	var kubectlArgs []string
	kubectlArgs = append(kubectlArgs, "create", "secret", "generic", name, "--dry-run=client", "-o", "yaml")
	kubectlArgs = append(kubectlArgs, "--from-literal", fmt.Sprintf("target=%s", passwords[0]))
	for i, password := range passwords[1:] {
		kubectlArgs = append(kubectlArgs, "--from-literal", fmt.Sprintf("historical-%d=%s", i, password))
	}
	secret, err := k8s.RunKubectlAndGetOutputE(t, opts, kubectlArgs...)
	require.NoError(t, err)
	tmpfile, err := ioutil.TempFile("", "tde_secret")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.WriteString(secret)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)
	k8s.RunKubectl(t, opts, "apply", "-f", tmpfile.Name())
	testlib.AddTeardown(testlib.TEARDOWN_ADMIN, func() {
		k8s.RunKubectlE(t, opts, "delete", "secret", name)
	})
}

func performStoragePasswordsRotation(t *testing.T, namespaceName string, adminPod string,
	dbName string, secretName string, passwords []string) {
	assert.Greater(t, len(passwords), 0)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	applyStoragePasswordSecret(t, namespaceName, secretName, passwords)
	testlib.Await(t, func() bool {
		err := k8s.RunKubectlE(t, kubectlOptions,
			"exec", adminPod, "-c", "admin", "--",
			"nuocmd", "check", "data-encryption",
			"--db-name", dbName,
			"--password", passwords[0],
		)
		return err == nil
	}, 90*time.Second)
}

func TestAdminColdStartWithTDE(t *testing.T) {
	// TODO: remove this whenever the image tested in nuodb-helm-charts CI
	// supports 'tde_monitor' service, i.e. whenever the version
	// is bumped to >4.1.1
	if os.Getenv("NUODB_DEV") != "true" {
		t.Skip("'tde_monitor' service is not supported in released versions")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.tde.secrets.demo":                "demo-tde-secret",
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"admin.tde.storagePasswordsDir":         "/etc/nuodb/encryption",
		},
	}
	opt := testlib.GetExtractedOptions(&options)
	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%s-admin-cold-start-with-tde-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	password := strings.ToLower(random.UniqueId())
	applyStoragePasswordSecret(t, namespaceName, options.SetValues["admin.tde.secrets.demo"], []string{password})

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Successfully updated storage passwords for dbName=demo", &corev1.PodLogOptions{}) >= 1
	}, 30*time.Second)
	k8s.RunKubectl(t, kubectlOptions,
		"exec", admin0, "-c", "admin", "--",
		"nuocmd", "check", "data-encryption",
		"--db-name", opt.DbName,
		"--password", password,
	)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &options)
	// Enable TDE on database level
	output, _ := testlib.RunSQL(t, namespaceName, admin0, "demo", "alter database change encryption type AES128")
	require.False(t, strings.Contains(output, "Error"), "Failed to enable TDE: %s", output)

	// Restarting the only admin pod will mean that all storage passwords
	// will be "forgoten" by the admin
	adminPod := testlib.GetPod(t, namespaceName, admin0)
	testlib.DeletePod(t, namespaceName, "pod/"+admin0)
	testlib.AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)
	testlib.AwaitAdminFullyConnected(t, namespaceName, admin0, 1)
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Successfully updated storage passwords for dbName=demo", &corev1.PodLogOptions{}) >= 1
	}, 30*time.Second)

	// verify that restarting an SM pod will succeed to start encrypted archive
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
	testlib.DeletePod(t, namespaceName, "pod/"+smPodName0)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, 2)
}

func TestRestoreInPlaceWithTDE(t *testing.T) {
	// TODO: remove this whenever the image tested in nuodb-helm-charts CI
	// supports 'tde_monitor' service, i.e. whenever the version
	// is bumped to >4.1.1
	if os.Getenv("NUODB_DEV") != "true" {
		t.Skip("'tde_monitor' service is not supported in released versions")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.tde.secrets.demo":                "demo-tde-secret",
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.logPersistence.enabled":    "true",
		},
	}
	opt := testlib.GetExtractedOptions(&options)
	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%s-restore-in-place-with-tde-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	password := strings.ToLower(random.UniqueId())
	applyStoragePasswordSecret(t, namespaceName, options.SetValues["admin.tde.secrets.demo"], []string{password})

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Successfully updated storage passwords for dbName=demo", &corev1.PodLogOptions{}) >= 1
	}, 30*time.Second)
	k8s.RunKubectl(t, kubectlOptions,
		"exec", admin0, "-c", "admin", "--",
		"nuocmd", "check", "data-encryption",
		"--db-name", opt.DbName,
		"--password", password,
	)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		opt := testlib.GetExtractedOptions(&options)
		pvcName := fmt.Sprintf("%s-nuodb-%s-%s-log-te-volume", databaseChartName, opt.ClusterName, opt.DbName)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", pvcName)
	})

	// Enable TDE on database level
	output, _ := testlib.RunSQL(t, namespaceName, admin0, "demo", "alter database change encryption type AES128")
	require.False(t, strings.Contains(output, "Error"), "Failed to enable TDE: %s", output)

	// populate some data
	testlib.CreateQuickstartSchema(t, namespaceName, admin0)

	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	// Take encrypted backup
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)
	options.SetValues["restore.source"] = backupset

	// Drop USER.HOCKEY table before restore
	testlib.RunSQL(t, namespaceName, admin0, "demo", "drop table USER.HOCKEY")
	// restore database
	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
	testlib.RestoreDatabase(t, namespaceName, admin0, &options)
	// verify that the database contains USER.HOCKEY table AFTER the restore
	tables, _ := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	require.True(t, strings.Contains(tables, "HOCKEY"), "Show schema returned: %s", tables)

	// Perform TDE password rotation
	passwordNew := strings.ToLower(random.UniqueId())
	performStoragePasswordsRotation(t, namespaceName, admin0, opt.DbName,
		options.SetValues["admin.tde.secrets.demo"], []string{passwordNew, password})

	// Drop USER.HOCKEY table before restore
	testlib.RunSQL(t, namespaceName, admin0, "demo", "drop table USER.HOCKEY")
	// Restore backup encrypted with old password
	testlib.RestoreDatabase(t, namespaceName, admin0, &options)
	// verify that the database contains USER.HOCKEY table AFTER the restore
	tables, _ = testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	require.True(t, strings.Contains(tables, "HOCKEY"), "Show schema returned: ", tables)
}
