// +build enterprise

package minikube

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	v12 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/ghodss/yaml"

	"github.com/gruntwork-io/terratest/modules/helm"
	rbacv1 "k8s.io/api/rbac/v1"
)

func modifyKubeInspectorRoleInPlace(t *testing.T, modificationFunc func(role *rbacv1.Role)) {
	inspectorRoleFile := filepath.Join(testlib.ADMIN_HELM_CHART_PATH, "templates", "role.yaml")
	originalData, err := ioutil.ReadFile(inspectorRoleFile)
	require.NoError(t, err)
	testlib.AddTeardown(testlib.TEARDOWN_ADMIN, func() { ioutil.WriteFile(inspectorRoleFile, originalData, 0644) })

	output := helm.RenderTemplate(t, &helm.Options{}, testlib.ADMIN_HELM_CHART_PATH, "release-name", []string{"templates/role.yaml"})
	roles := testlib.SplitAndRenderRole(t, output, 1)
	modificationFunc(&roles[0])
	roleBytes, err := yaml.Marshal(&roles[0])
	require.NoError(t, err)
	err = ioutil.WriteFile(inspectorRoleFile, roleBytes, 0644)
	require.NoError(t, err)
	out, _ := ioutil.ReadFile(inspectorRoleFile)
	t.Log("Modified roles file:\n" + string(out))
}

func TestKaaLimitedPermissions(t *testing.T) {
	// This test requires NuoDB 4.1.1+
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	modifyKubeInspectorRoleInPlace(t, func(role *rbacv1.Role) {
		// Remove "daemonsets" from the resources list
		role.Rules[1].Resources = []string{"statefulsets", "deployments"}
		// Remove "list" verb for pods and PVCs
		// The KAA resync logic should not remove processes or archives due to missing permissions
		role.Rules[0].Verbs = []string{"get", "watch"}
	})

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	databaseHelmChartReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	// Verify that KAA will not register informer for daemonsets
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Informer for resource 'daemonsets' not registered", &v12.PodLogOptions{}) == 1
	}, 30*time.Second)
	// Verify that KAA will not register informer for pods
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Informer for resource 'pods' not registered", &v12.PodLogOptions{}) == 1
	}, 30*time.Second)
	// Verify that KAA will not register informer for PVCs
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Informer for resource 'persistentvolumeclaims' not registered", &v12.PodLogOptions{}) == 1
	}, 30*time.Second)

	// Verify that resources that KAA have permissions for are available
	adminStatefulSet := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	teDeployment := fmt.Sprintf("te-%s-nuodb-cluster0-demo", databaseHelmChartReleaseName)
	config := testlib.GetNuoDBK8sConfigDump(t, namespaceName, admin0)
	require.True(t, func() bool { _, ok := config.StatefulSets[adminStatefulSet]; return ok }())
	require.True(t, func() bool { _, ok := config.Deployments[teDeployment]; return ok }())
	require.True(t, len(config.Volumes) == 0)
	require.True(t, len(config.Pods) == 0)
}

func TestKaaRolebindingDisabled(t *testing.T) {
	// This test requires NuoDB 4.2+ or 4.1.2+
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"nuodb.addRoleBinding": "false",
		},
	}, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	// Verify that KAA won't start due to limited mandatory permissions
	testlib.Await(t, func() bool {
		return testlib.GetStringOccurrenceInLog(t, namespaceName, admin0,
			"Not registering event listeners: service account unauthorized for resource 'leases'", &v12.PodLogOptions{}) == 1
	}, 30*time.Second)

	// Verify that no resource are avaialble via KKA
	config := testlib.GetNuoDBK8sConfigDump(t, namespaceName, admin0)
	require.True(t, len(config.StatefulSets) == 0)
	require.True(t, len(config.Deployments) == 0)
	require.True(t, len(config.Volumes) == 0)
	require.True(t, len(config.Pods) == 0)
}