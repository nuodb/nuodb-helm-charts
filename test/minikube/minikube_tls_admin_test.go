//go:build long
// +build long

package minikube

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

const ENGINE_CERTIFICATE_LOG_TEMPLATE = `Engine Certificate: Certificate #%d CN %s`

func verifyKeystore(t *testing.T, namespace string, podName string, keystore string, password string, matches string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", "admin", "--", "nuocmd", "show", "certificate", "--keystore", keystore, "--store-password", password)
	output = testlib.RemoveEmptyLines(output)
	matches = testlib.RemoveEmptyLines(matches)

	t.Log("<" + output + ">")
	t.Log("<" + matches + ">")

	require.NoError(t, err)
	require.Equal(t, matches, output)
}

func TestKubernetesTLS(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetestls-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
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

	// sometimes the test fails because SMs doesn't go ready due to probe
	// timeout
	testlib.OverrideReadinessProbesTimeout(t, &options, "10")

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

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &localOptions)

		tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)
		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		go testlib.GetAppLog(t, namespaceName, tePodName, "", &corev1.PodLogOptions{Follow: true})

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

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &localOptions)

		tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)
		tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
		go testlib.GetAppLog(t, namespaceName, tePodName, "", &corev1.PodLogOptions{Follow: true})

		// TE certificate is not signed by the admin and the DN entry is the generic admin name
		// this is the 3rd pod name because: #0 and #1 are trusted certs, #2 is CA, #3 is admin (and engine)
		expectedLogLine := fmt.Sprintf(ENGINE_CERTIFICATE_LOG_TEMPLATE, 3, "nuoadmin.nuodb.com")
		testlib.VerifyCertificateInLog(t, namespaceName, tePodName, expectedLogLine)
	})
}
