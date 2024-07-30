package multicluster

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestMultiClusterAdminLabelAffinity(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	// For multi-cluster deployment to work correctly, there are two prerequisites:
	// - the same namespace name should be used in all clusters
	// - the same admin Helm release should be used in all cluster (probably something to fix)
	randomSuffix := strings.ToLower(random.UniqueId())
	adminReleaseName := fmt.Sprintf("admin-%s", randomSuffix)
	namespaceName := fmt.Sprintf("%skubernetesbasicmulticluster-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)

	options := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	context := context.Background()
	cluster1Context := testlib.NewClusterDeploymentContext(context,
		&options, testlib.MULTI_CLUSTER_1, testlib.MULTI_CLUSTER_1)
	cluster2Context := testlib.NewClusterDeploymentContext(context,
		&options, testlib.MULTI_CLUSTER_2, testlib.MULTI_CLUSTER_1)

	// Add each DNS server as an upstream resolver for the other
	testlib.UpdateDnsConfig(t, cluster1Context, cluster2Context)
	testlib.UpdateDnsConfig(t, cluster2Context, cluster1Context)

	defer testlib.Teardown(testlib.TEARDOWN_MULTICLUSTER)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	testlib.DeployWithContext(t,
		cluster1Context,
		func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
			testlib.CreateNamespace(t, namespaceName)
			testlib.StartAdminCustomRelease(t, options, 1, namespaceName, adminReleaseName)
			admin0 := fmt.Sprintf("%s-nuodb-%s-0", adminReleaseName, context.ThisCluster.Name)
			databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)
			// Store deployment details in the context
			context.AdminReleaseName = adminReleaseName
			context.DatabaseReleaseName = databaseReleaseName
			context.Namespace = namespaceName
		},
	)

	testlib.DeployWithContext(t,
		cluster2Context,
		func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
			testlib.CreateNamespace(t, namespaceName)
			testlib.StartAdminCustomRelease(t, options, 1, namespaceName, adminReleaseName)
			admin0 := fmt.Sprintf("%s-nuodb-%s-0", adminReleaseName, context.ThisCluster.Name)
			testlib.ApplyLicense(t, namespaceName, admin0, testlib.ENTERPRISE)
			databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)
			// Store deployment details in the context
			context.AdminReleaseName = adminReleaseName
			context.DatabaseReleaseName = databaseReleaseName
			context.Namespace = namespaceName
		},
	)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
			admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
			testlib.GetDiagnoseOnTestFailure(t, context.Namespace, admin0)
		})
	})

	t.Run("verifyDomain", func(t *testing.T) {
		testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
			admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
			testlib.AwaitAdminFullyConnected(t, context.Namespace, admin0, 2)
			testlib.AwaitDatabaseUp(t, context.Namespace, admin0, "demo", 4)
		})
	})
}
