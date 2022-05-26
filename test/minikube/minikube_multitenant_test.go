//go:build short
// +build short

package minikube

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/require"

	"github.com/Masterminds/semver"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestKubernetesMultiTenantDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
		k8s.RunKubectl(t, kubectlOptions, "get", "pods", "-o", "wide")
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
	})

	// provision two database "green" and "blue" managed by the same admin
	// domain
	for _, dbName := range []string{"green", "blue"} {
		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.name":                         dbName,
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			},
		})

		t.Run(fmt.Sprintf("verifyNuoSQL-%s", dbName), func(t *testing.T) {
			verifyNuoSQL(t, namespaceName, admin0, dbName)
		})
	}
}

func TestKubernetesNamespaceCoexistence(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%snamespacecoexistence-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// provision two domains "green" and "blue" in the same namespace; each
	// domain will manage a NuoDB database "demo" which is installed as a
	// primary database Helm release
	for i, project := range []string{"green", "blue"} {
		lbPolicyKey := fmt.Sprintf("admin.lbConfig.policies.%s", project)
		lbQuery := fmt.Sprintf("round_robin(first(label(project %s) any))", project)
		options := helm.Options{
			SetValues: map[string]string{
				"admin.domain":                          project,
				"admin.resourceLabels.project":          project,
				lbPolicyKey:                             lbQuery,
				"database.resourceLabels.project":       project,
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.lbConfig.default":             lbQuery,
			},
		}
		if i > 0 {
			// the RBAC needs to be provisioned before hand or enabled only with
			// the first helm release; it will be enabled only for the first
			// domain installation
			options.SetValues["nuodb.addRoleBinding"] = "false"
			options.SetValues["nuodb.addServiceAccount"] = "false"
		}
		adminReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)

		opt := testlib.GetExtractedOptions(&options)
		admin0 := fmt.Sprintf("%s-%s-%s-0", adminReleaseName, opt.DomainName, opt.ClusterName)

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

		t.Run(fmt.Sprintf("verifyNuoSQL-%s", project), func(t *testing.T) {
			verifyNuoSQL(t, namespaceName, admin0, "demo")
		})

		t.Run(fmt.Sprintf("verifyResourceLabels-%s", project), func(t *testing.T) {
			adminSts := fmt.Sprintf("%s-%s-%s", adminReleaseName, opt.DomainName, opt.ClusterName)
			hcsmSts := fmt.Sprintf("sm-%s-%s-%s-%s-hotcopy", databaseReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
			smSts := fmt.Sprintf("sm-%s-%s-%s-%s", databaseReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
			for _, stsName := range []string{adminSts, hcsmSts, smSts} {
				sts := testlib.GetStatefulSet(t, namespaceName, stsName)
				msg, ok := testlib.MapContains(sts.GetLabels(), map[string]string{"project": project})
				require.Truef(t, ok, "The 'project' label doesn't match for statefulset %s: %s", stsName, msg)
			}
		})

		testlib.RunOnNuoDBVersionCondition(t, ">=4.3.3", func(version *semver.Version) {
			t.Run(fmt.Sprintf("verifyLoadBalancerPolicy-%s", project), func(t *testing.T) {
				// there should be 3 load balancer configurations on each
				// domain: nearest, green|blue and __default/demo
				testlib.AwaitNrLoadBalancerPolicies(t, namespaceName, admin0, 3)
				verifyLoadBalancer(t, namespaceName, admin0, options.SetValues)
			})
		})
	}
}
