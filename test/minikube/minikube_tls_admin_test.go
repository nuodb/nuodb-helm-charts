// +build long

package minikube

import (
	"fmt"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const ENGINE_CERTIFICATE_LOG_TEMPLATE = `Engine Certificate: Certificate #%d CN %s`

func verifyKeystore(t *testing.T, namespace string, podName string, keystore string, password string, matches string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "certificate", "--keystore", keystore, "--store-password", password)
	output = testlib.RemoveEmptyLines(output)
	matches = testlib.RemoveEmptyLines(matches)

	t.Log("<" + output + ">")
	t.Log("<" + matches + ">")

	require.NoError(t, err)
	require.Equal(t, matches, output)
}

func TestKubernetesTLS(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("testkubernetestls-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)

	// create the certs and secrets...
	tlsCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
		"nuocmd show certificate --keystore " + testlib.KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" > nuoadmin.cert",
	}
	_, keysLocation := testlib.GenerateTLSConfiguration(t, namespaceName, tlsCommands)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas":               "3",
			"admin.tlsCACert.secret":       testlib.CA_CERT_SECRET,
			"admin.tlsCACert.key":          testlib.CA_CERT_FILE,
			"admin.tlsKeyStore.secret":     testlib.KEYSTORE_SECRET,
			"admin.tlsKeyStore.key":        testlib.KEYSTORE_FILE,
			"admin.tlsKeyStore.password":   testlib.SECRET_PASSWORD,
			"admin.tlsTrustStore.secret":   testlib.TRUSTSTORE_SECRET,
			"admin.tlsTrustStore.key":      testlib.TRUSTSTORE_FILE,
			"admin.tlsTrustStore.password": testlib.SECRET_PASSWORD,
			"admin.tlsClientPEM.secret":    testlib.NUOCMD_SECRET,
			"admin.tlsClientPEM.key":       testlib.NUOCMD_FILE,
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 3, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("verifyKeystore", func(t *testing.T) {
		content, err := testlib.ReadAll(filepath.Join(keysLocation, "nuoadmin.cert"))
		require.NoError(t, err)
		verifyKeystore(t, namespaceName, admin0, testlib.KEYSTORE_FILE, testlib.SECRET_PASSWORD, string(content))
	})

	t.Run("testDatabaseNoDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := options
		localOptions.SetValues["database.sm.resources.requests.cpu"] = testlib.MINIMAL_VIABLE_ENGINE_CPU
		localOptions.SetValues["database.sm.resources.requests.memory"] = testlib.MINIMAL_VIABLE_ENGINE_MEMORY
		localOptions.SetValues["database.te.resources.requests.cpu"] = testlib.MINIMAL_VIABLE_ENGINE_CPU
		localOptions.SetValues["database.te.resources.requests.memory"] = testlib.MINIMAL_VIABLE_ENGINE_MEMORY

		defer testlib.Teardown("database")

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &localOptions)

		tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)
		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		go testlib.GetAppLog(t, namespaceName, tePodName, "", &v12.PodLogOptions{Follow: true})

		// TE certificate is signed by the admin and the DN entry is the pod name
		// this is the 4th pod name because: #0 and #1 are trusted certs, #2 is CA, #3 is admin, #4 is engine
		expectedLogLine := fmt.Sprintf(ENGINE_CERTIFICATE_LOG_TEMPLATE, 4, tePodName)
		testlib.VerifyCertificateInLog(t, namespaceName, tePodName, expectedLogLine)
	})

	t.Run("testDatabaseDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := options
		localOptions.SetValues["database.sm.resources.requests.cpu"] = testlib.MINIMAL_VIABLE_ENGINE_CPU
		localOptions.SetValues["database.sm.resources.requests.memory"] = testlib.MINIMAL_VIABLE_ENGINE_MEMORY
		localOptions.SetValues["database.te.resources.requests.cpu"] = testlib.MINIMAL_VIABLE_ENGINE_CPU
		localOptions.SetValues["database.te.resources.requests.memory"] = testlib.MINIMAL_VIABLE_ENGINE_MEMORY

		localOptions.SetValues["database.te.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"
		localOptions.SetValues["database.sm.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"

		defer testlib.Teardown("database")

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &localOptions)

		tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)
		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		go testlib.GetAppLog(t, namespaceName, tePodName, "", &v12.PodLogOptions{Follow: true})

		// TE certificate is not signed by the admin and the DN entry is the generic admin name
		// this is the 3rd pod name because: #0 and #1 are trusted certs, #2 is CA, #3 is admin (and engine)
		expectedLogLine := fmt.Sprintf(ENGINE_CERTIFICATE_LOG_TEMPLATE, 3, "nuoadmin.nuodb.com")
		testlib.VerifyCertificateInLog(t, namespaceName, tePodName, expectedLogLine)
	})
}

func TestHashiCorpVault(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	namespaceName := fmt.Sprintf("testvault-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_VAULT)

	vaultOptions := helm.Options{}

	helmChartReleaseName := testlib.StartVault(t, &vaultOptions, namespaceName)
	vaultName := fmt.Sprintf("%s-vault-0", helmChartReleaseName)

	testlib.CreateVault(t, namespaceName, vaultName)
	testlib.EnableVaultKubernetesIntegration(t, namespaceName, vaultName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	testlib.CreateSecretsInVault(t, namespaceName, vaultName)

	adminOptions := helm.Options{
		ValuesFiles: []string{"../files/vault-annotations-admin.yaml"},
		SetValues: map[string]string {
			"vault.hashicorp.com/log-level": "trace", // increase debug level for testing
			"admin.replicas": "3",
			"admin.options.truststore-password": "$(</etc/nuodb/keys/nuoadmin-truststore.password)",
			"admin.options.keystore-password": "$(</etc/nuodb/keys/nuoadmin.password)",
			"admin.options.ssl": "true",
		},
	}

	adminHelmChartReleaseName, _ := testlib.StartAdmin(t, &adminOptions, 3, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminHelmChartReleaseName)

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })

	t.Run("testDatabaseHashiCorpVault", func(t *testing.T) {
		t.Skip("CircleCI K8s nodes do not have enough CPU to run Vault+admin+engines")
		engineOptions := helm.Options{
			ValuesFiles: []string{"../files/vault-annotations-database.yaml"},
			SetValues: map[string]string{
				"vault.hashicorp.com/log-level":         "trace", // increase debug level for testing
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			},
		}

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		testlib.StartDatabase(t, namespaceName, admin0, &engineOptions)
	})
}