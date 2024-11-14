package multicluster

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func verifyNuoSQL(t *testing.T, namespaceName string, adminPod string, databaseName string) {
	output, err := testlib.RunSQL(t, namespaceName, adminPod, databaseName, "select * from system.nodes")
	require.NoError(t, err, output)
	require.True(t, strings.Contains(output, "Storage"))
	require.True(t, strings.Contains(output, "Transaction"))
}

func TestKubernetesRemoveOrphanNamespaces(t *testing.T) {
	// Golang 1.4+ provides TestMain() function hook which can be used to
	// execute setup/teardown tasks, however, it receives a *testing.M instance
	// and all methods in our test framework rely on passing *testing.T
	// instance. Execute this cleanup tasks as separate test case for simplicity
	// as it's only needed in multi-cluster infrastructure.
	defer testlib.VerifyTeardown(t)

	context := context.Background()
	// Create contexts for all available clusters
	testlib.NewClusterDeploymentContext(context,
		&helm.Options{}, testlib.MULTI_CLUSTER_1, testlib.MULTI_CLUSTER_1)
	testlib.NewClusterDeploymentContext(context,
		&helm.Options{}, testlib.MULTI_CLUSTER_2, testlib.MULTI_CLUSTER_1)

	defer testlib.Teardown(testlib.TEARDOWN_MULTICLUSTER)

	// Sometimes when previous test job is canceled in CircleCI the testing
	// teardown logic is not executed which leaves orphan test pods in the
	// shared Kubernetes cluster. Remove any long-running test namespaces from
	// all clusters in multi-cluster setup.
	testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
		t.Logf("Removing orphan namespaces in cluster name=%s, context=%s",
			context.ThisCluster.Name, context.ThisCluster.Context)
		testlib.RemoveOrphanNamespaces(t)
	})
}

func TestKubernetesBasicMultiCluster(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" && os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
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
	t.Run("verifyNuoSQL", func(t *testing.T) {
		testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
			admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
			verifyNuoSQL(t, context.Namespace, admin0, "demo")
		})
	})

	testlib.RunOnNuoDBVersionCondition(t, ">=5.0.3", func(version *semver.Version) {
		t.Run("testDatabaseResync", func(t *testing.T) {
			testlib.DeployWithContext(t,
				cluster1Context,
				func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
					// delete database controllers and PVCs
					helm.Delete(t, options, context.DatabaseReleaseName, true)
					testlib.AwaitNoPods(t, namespaceName, context.DatabaseReleaseName)
					smHCPvcName := fmt.Sprintf("archive-volume-sm-%s-nuodb-%s-demo-hotcopy-0",
						context.DatabaseReleaseName, context.ThisCluster.Name)
					k8s.RunKubectl(t, kubectlOptions, "delete", "pvc", smHCPvcName)
				},
			)
			// check that the database is not deleted
			testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
				admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
				testlib.AwaitDatabaseUp(t, context.Namespace, admin0, "demo", 2)
				testlib.CheckArchives(t, namespaceName, admin0, "demo", 1, 0)
				db, err := testlib.GetDatabaseE(t, namespaceName, admin0, "demo")
				require.NoError(t, err)
				require.NotEqual(t, "TOMBSTONE", db.State)
			})

			testlib.DeployWithContext(t,
				cluster2Context,
				func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
					admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
					// delete the database controllers
					helm.Delete(t, options, context.DatabaseReleaseName, true)
					testlib.AwaitNoPods(t, namespaceName, context.DatabaseReleaseName)

					// check that database is NOT_RUNNING
					testlib.Await(t, func() bool {
						db, err := testlib.GetDatabaseE(t, namespaceName, admin0, "demo")
						require.NoError(t, err)
						return db.State == "NOT_RUNNING"
					}, 30*time.Second)
				},
			)
			// check that the database is not deleted
			testlib.ExecuteInAllClusters(t, func(context *testlib.ClusterDeploymentContext) {
				admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
				testlib.CheckArchives(t, namespaceName, admin0, "demo", 1, 0)
				db, err := testlib.GetDatabaseE(t, namespaceName, admin0, "demo")
				require.NoError(t, err)
				require.NotEqual(t, "TOMBSTONE", db.State)
			})

			testlib.DeployWithContext(t,
				cluster2Context,
				func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
					admin0 := fmt.Sprintf("%s-nuodb-%s-0", context.AdminReleaseName, context.ThisCluster.Name)
					// delete all database PVCs
					smHCPvcName := fmt.Sprintf("archive-volume-sm-%s-nuodb-%s-demo-hotcopy-0",
						context.DatabaseReleaseName, context.ThisCluster.Name)
					k8s.RunKubectl(t, kubectlOptions, "delete", "pvc", smHCPvcName)

					// check that database is deleted
					testlib.Await(t, func() bool {
						db, err := testlib.GetDatabaseE(t, namespaceName, admin0, "demo")
						require.NoError(t, err)
						return db.State == "TOMBSTONE"
					}, 30*time.Second)
					testlib.CheckArchives(t, namespaceName, admin0, "demo", 0, 0)
				},
			)
		})
	})
}
