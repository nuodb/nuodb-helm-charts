package minikube

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const TLS_SECRET_PASSWORD_YAML_TEMPLATE = `---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
apiVersion: v1
data:
  %s: %s
  password: %s
`

const TLS_SECRET_NO_PASSWORD_YAML_TEMPLATE = `---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
apiVersion: v1
data:
  %s: %s
`

func verifySecretFields(t *testing.T, namespaceName string, secretName string, fields ...string) {
	secret := testlib.GetSecret(t, namespaceName, secretName)
	for _, field := range fields {
		_, ok := secret.Data[field]
		assert.Check(t, ok)
	}
}

func readAll(path string) ([]byte, error) {
	file, ferr := os.Open(path)
	if ferr != nil {
		return nil, ferr
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content, rerr := ioutil.ReadAll(reader)
	if rerr != nil {
		return nil, rerr
	}
	return content, nil
}

func readAsBase64(path string) (string, error) {
	content, err := readAll(path)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

func createSecretDecl(path string, namespace string, name string, key string) (string, error) {
	base64, err := readAsBase64(path)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf(TLS_SECRET_NO_PASSWORD_YAML_TEMPLATE,
		name, namespace, key, base64)
	return text, nil
}

func createSecretPassDecl(path string, namespace string, name string, key string, password string) (string, error) {
	base64, err := readAsBase64(path)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf(TLS_SECRET_PASSWORD_YAML_TEMPLATE,
		name, namespace, key, base64, password)
	return text, nil
}

func verifyKeystore(t *testing.T, namespace string, podName string, keystore string, password string, matches string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "certificate", "--keystore", keystore, "--store-password", password)

	t.Log(output)
	t.Log(matches)

	assert.NilError(t, err)
	assert.Assert(t, strings.Compare(output, matches) != 0)
}

func createSecret(t *testing.T, namespaceName string, certName string, secretName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	keyDir := filepath.Join("..", "..", "keys")
	keyFile := filepath.Join(keyDir, certName)

	secretString, err := createSecretDecl(keyFile, namespaceName, secretName, certName)
	assert.NilError(t, err)

	k8s.KubectlApplyFromString(t, kubectlOptions, secretString)
	testlib.AddTeardown(testlib.TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromString(t, kubectlOptions, secretString) })

	fields := []string{certName}
	verifySecretFields(t, namespaceName, secretName, fields...)
}

func createSecretWithPassword(t *testing.T, namespaceName string, certName string, secretName string, password string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	keyDir := filepath.Join("..", "..", "keys")
	keyFile := filepath.Join(keyDir, certName)

	secretString, err := createSecretPassDecl(keyFile, namespaceName, secretName, certName, password)
	assert.NilError(t, err)

	k8s.KubectlApplyFromString(t, kubectlOptions, secretString)
	testlib.AddTeardown(testlib.TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromString(t, kubectlOptions, secretString) })

	fields := []string{certName, "password"}
	verifySecretFields(t, namespaceName, secretName, fields...)
}

func TestKubernetesTLS(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	namespaceName := fmt.Sprintf("test-admin-tls-%s", randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "")
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	// create the certs...
	createSecret(t, namespaceName, testlib.CA_CERT_FILE, testlib.CA_CERT_SECRET)
	createSecret(t, namespaceName, testlib.NUOCMD_FILE, testlib.NUOCMD_SECRET)
	createSecretWithPassword(t, namespaceName, testlib.KEYSTORE_FILE, testlib.KEYSTORE_SECRET, testlib.SECRET_PASSWORD)
	createSecretWithPassword(t, namespaceName, testlib.TRUSTSTORE_FILE, testlib.TRUSTSTORE_SECRET, testlib.SECRET_PASSWORD)

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

	helmChartReleaseName, namespaceName := startAdmin(t, &options, 3, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)


	t.Run("verifyKeystore", func(t *testing.T) {
		content, err := readAll("../../keys/default.certificate")
		assert.NilError(t, err)
		verifyKeystore(t, namespaceName, admin0, testlib.KEYSTORE_FILE, testlib.SECRET_PASSWORD, string(content))
	})

	t.Run("testDatabaseNoDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := options
		localOptions.SetValues["database.sm.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.sm.resources.requests.memory"] = "1Gi"
		localOptions.SetValues["database.te.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.te.resources.requests.memory"] = "1Gi"

		defer testlib.Teardown("database")

		startDatabase(t, namespaceName, admin0, &localOptions)
	})

	t.Run("testDatabaseDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := options
		localOptions.SetValues["database.sm.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.sm.resources.requests.memory"] = "1Gi"
		localOptions.SetValues["database.te.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.te.resources.requests.memory"] = "1Gi"

		localOptions.SetValues["database.te.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"
		localOptions.SetValues["database.sm.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"

		defer testlib.Teardown("database")

		startDatabase(t, namespaceName, admin0, &localOptions)
	})
}