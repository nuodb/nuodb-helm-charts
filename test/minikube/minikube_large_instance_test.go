// +build large

package minikube

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/random"
)

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