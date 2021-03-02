package testlib

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func StartVault(t *testing.T, options *helm.Options, namespaceName string) string {
	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}

	options.SetValues["server.dev.enabled"] = "true"

	randomSuffix := strings.ToLower(random.UniqueId())

	helmChartReleaseName := fmt.Sprintf("hc-%s", randomSuffix)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	helm.Install(t, options, "hashicorp/vault ", helmChartReleaseName)
	AddTeardown(TEARDOWN_VAULT, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})

	// there are two pods here, the vault itself and an agent-injector
	AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, 2)

	vaultName := fmt.Sprintf("%s-vault-0", helmChartReleaseName)

	AwaitPodUp(t, namespaceName, vaultName, 300*time.Second)

	// enable audit logging
	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "audit", "enable", "file", "file_path=stdout")

	AddTeardown(TEARDOWN_VAULT, func() {
		_, err := k8s.GetPodE(t, kubectlOptions, vaultName)
		if err != nil {
			t.Logf("Vault pod '%s' is not available and logs can not be retrieved", vaultName)
		} else {
			go GetAppLog(t, namespaceName, vaultName, "", &v12.PodLogOptions{Follow: true})
		}
	})

	return helmChartReleaseName
}

func CreateSecretsInVault(t *testing.T, namespaceName string, vaultName string, tlsKeyLocation string) {
	createSecretFromFile(t, namespaceName, vaultName, filepath.Join(tlsKeyLocation, CA_CERT_FILE), "tlsCACert")
	patchSecretFromFile(t, namespaceName, vaultName, filepath.Join(tlsKeyLocation, NUOCMD_FILE), "tlsClientPEM")
	patcSecretFromString(t, namespaceName, vaultName, SECRET_PASSWORD, "tlsKeyStorePassword")
	patcSecretFromString(t, namespaceName, vaultName, SECRET_PASSWORD, "tlsTrustStorePassword")
	patchSecretFromFile(t, namespaceName, vaultName, filepath.Join(tlsKeyLocation, KEYSTORE_FILE), "tlsKeyStore")
	patchSecretFromFile(t, namespaceName, vaultName, filepath.Join(tlsKeyLocation, TRUSTSTORE_FILE), "tlsTrustStore")
}

func createSecretFromFile(t *testing.T, namespaceName string, vaultName string, keyFile string, secretName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	base64, err := readAsBase64(keyFile)
	require.NoError(t, err)

	secret := fmt.Sprintf("%s=%s", secretName, base64)

	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "kv", "put", "nuodb.com/TLS", secret)
}

func patchSecretFromFile(t *testing.T, namespaceName string, vaultName string, keyFile string, secretName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	base64, err := readAsBase64(keyFile)
	require.NoError(t, err)

	secret := fmt.Sprintf("%s=%s", secretName, base64)

	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "kv", "patch", "nuodb.com/TLS", secret)
}

func patcSecretFromString(t *testing.T, namespaceName string, vaultName string, keyText string, secretName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	secret := fmt.Sprintf("%s=%s", secretName, keyText)

	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "kv", "patch", "nuodb.com/TLS", secret)
}

func EnableVaultKubernetesIntegration(t *testing.T, namespaceName string, vaultName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	// install security policy
	policyLocation := "/home/vault/nuodb-policy.hcl"
	k8s.RunKubectl(t, kubectlOptions, "cp", "../files/nuodb-policy.hcl", vaultName+":"+policyLocation)
	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "policy", "write", "nuodb-policy", policyLocation)

	// enable Kubernetes Vault Auth
	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "auth", "enable", "kubernetes")

	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "sh", "-x", "-c",
		`vault write auth/kubernetes/config \
		token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
		kubernetes_host="https://${KUBERNETES_PORT_443_TCP_ADDR}:443" \
		kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`)

	// give permissions to test namespace
	namespaces := fmt.Sprintf("bound_service_account_namespaces=%s", namespaceName)
	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "write", "auth/kubernetes/role/nuodb",
		"bound_service_account_names=nuodb", namespaces, "policies=nuodb-policy", "ttl=1h")
}

func CreateVault(t *testing.T, namespaceName string, vaultName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.RunKubectl(t, kubectlOptions, "exec", vaultName, "--", "vault", "secrets", "enable", "-version=2", "-path=nuodb.com", "kv")
}
