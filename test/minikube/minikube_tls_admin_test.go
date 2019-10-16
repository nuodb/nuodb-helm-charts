package minikube

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"gotest.tools/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

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

func TestKubernetesBasicAdminThreeReplicasTLS(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	helmChartReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	admin0 := fmt.Sprintf("%s-nuodb-0", helmChartReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-1", helmChartReleaseName)
	admin2 := fmt.Sprintf("%s-nuodb-2", helmChartReleaseName)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.replicas":               "3",
			"admin.tlsCACert.secret":       "nuodb-ca-cert",
			"admin.tlsCACert.key":          "ca.cert",
			"admin.tlsKeyStore.secret":     "nuodb-keystore",
			"admin.tlsKeyStore.key":        "nuoadmin.p12",
			"admin.tlsKeyStore.password":   "changeIt",
			"admin.tlsTrustStore.secret":   "nuodb-truststore",
			"admin.tlsTrustStore.key":      "nuoadmin-truststore.p12",
			"admin.tlsTrustStore.password": "changeIt",
			"admin.tlsClientPEM.secret":    "nuodb-client-pem",
			"admin.tlsClientPEM.key":       "nuocmd.pem",
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("test-admin-tls-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	// create the certs...

	tlsCaCertString, err := createSecretDecl("../../keys/ca.cert",
		namespaceName, "nuodb-ca-cert", "ca.cert")
	assert.NilError(t, err)
	tlsCaCertOptions := k8s.NewKubectlOptions("", "")
	tlsCaCertOptions.Namespace = namespaceName
	defer k8s.KubectlDeleteFromString(t, tlsCaCertOptions, tlsCaCertString)
	k8s.KubectlApplyFromString(t, tlsCaCertOptions, tlsCaCertString)
	tlsCaCertFields := []string{"ca.cert"}
	verifySecretFields(t, namespaceName, "nuodb-ca-cert", tlsCaCertFields...)

	tlsClientPemString, err := createSecretDecl("../../keys/nuocmd.pem",
		namespaceName, "nuodb-client-pem", "nuocmd.pem")
	assert.NilError(t, err)
	tlsClientPemOptions := k8s.NewKubectlOptions("", "")
	tlsClientPemOptions.Namespace = namespaceName
	defer k8s.KubectlDeleteFromString(t, tlsClientPemOptions, tlsClientPemString)
	k8s.KubectlApplyFromString(t, tlsClientPemOptions, tlsClientPemString)
	tlsClientPemFields := []string{"nuocmd.pem"}
	verifySecretFields(t, namespaceName, "nuodb-client-pem", tlsClientPemFields...)

	tlsKeyStoreString, err := createSecretPassDecl("../../keys/nuoadmin.p12",
		namespaceName, "nuodb-keystore", "nuoadmin.p12", "changeIt")
	assert.NilError(t, err)
	tlsKeyStoreOptions := k8s.NewKubectlOptions("", "")
	tlsKeyStoreOptions.Namespace = namespaceName
	defer k8s.KubectlDeleteFromString(t, tlsKeyStoreOptions, tlsKeyStoreString)
	k8s.KubectlApplyFromString(t, tlsKeyStoreOptions, tlsKeyStoreString)
	tlsKeyStoreFields := []string{"nuoadmin.p12", "password"}
	verifySecretFields(t, namespaceName, "nuodb-keystore", tlsKeyStoreFields...)

	tlsTrustStoreString, err := createSecretPassDecl("../../keys/nuoadmin-truststore.p12",
		namespaceName, "nuodb-truststore", "nuoadmin-truststore.p12", "changeIt")
	assert.NilError(t, err)
	tlsTrustStoreOptions := k8s.NewKubectlOptions("", "")
	tlsTrustStoreOptions.Namespace = namespaceName
	defer k8s.KubectlDeleteFromString(t, tlsTrustStoreOptions, tlsTrustStoreString)
	k8s.KubectlApplyFromString(t, tlsTrustStoreOptions, tlsTrustStoreString)
	tlsTrustStoreFields := []string{"nuoadmin-truststore.p12", "password"}
	verifySecretFields(t, namespaceName, "nuodb-truststore", tlsTrustStoreFields...)

	// install and check admin

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 3)

	// first await could be pulling the image from the repo
	testlib.AwaitAdminPodUp(t, namespaceName, admin0, 300*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin1, 100*time.Second)
	testlib.AwaitAdminPodUp(t, namespaceName, admin2, 100*time.Second)

	defer testlib.GetAppLog(t, namespaceName, admin0)
	defer testlib.GetAppLog(t, namespaceName, admin1)
	defer testlib.GetAppLog(t, namespaceName, admin2)

	t.Run("verifyKeystore", func(t *testing.T) {
		content, err := readAll("../../keys/default.certificate")
		assert.NilError(t, err)
		verifyKeystore(t, namespaceName, admin0, "nuoadmin.p12", "changeIt", string(content))
	})

	t.Run("testDatabaseNoDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.sm.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.sm.resources.requests.memory"] = "1Gi"
		localOptions.SetValues["database.te.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.te.resources.requests.memory"] = "1Gi"

		defer testlib.Teardown("database")

		startDatabase(t, namespaceName, admin0, localOptions)
	})

	t.Run("testDatabaseDirectEngineKeys", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.sm.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.sm.resources.requests.memory"] = "1Gi"
		localOptions.SetValues["database.te.resources.requests.cpu"] =    "500m"
		localOptions.SetValues["database.te.resources.requests.memory"] = "1Gi"

		localOptions.SetValues["database.te.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"
		localOptions.SetValues["database.sm.otherOptions.keystore"] = "/etc/nuodb/keys/nuoadmin.p12"

		defer testlib.Teardown("database")

		startDatabase(t, namespaceName, admin0, localOptions)
	})
}

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
