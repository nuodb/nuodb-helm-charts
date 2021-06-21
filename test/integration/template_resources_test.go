package integration

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func listContains(arr []string, s string) bool {
	for _, ele := range arr {
		if strings.Contains(ele, s) {
			return true
		}
	}

	return false
}

func TestResourcesAdminDefaults(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		containers := &obj.Spec.Template.Spec.Containers

		assert.NotEmpty(t, containers)
		assert.Nil(t, (*containers)[0].Resources.Limits)
		assert.Nil(t, (*containers)[0].Resources.Requests)
	}
}

func TestResourcesAdminOverridden(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.resources.requests.cpu":    "1",
			"admin.resources.requests.memory": "4G",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		containers := &obj.Spec.Template.Spec.Containers

		assert.NotEmpty(t, containers)
		assert.Nil(t, (*containers)[0].Resources.Limits)
		assert.NotNil(t, (*containers)[0].Resources.Requests)

		assert.EqualValues(t, 1, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0))
		assert.EqualValues(t, 4, (*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga),
			(*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga))
	}
}

func TestResourcesDatabaseDefaults(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	foundBackupEnabled := false
	foundBackupDisabled := false

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
		containers := &obj.Spec.Template.Spec.Containers

		assert.NotEmpty(t, containers)
		assert.NotNil(t, (*containers)[0].Resources.Limits)
		assert.NotNil(t, (*containers)[0].Resources.Requests)

		// the memory is confusing Gi with G. We are using power of two (1024). But scaled value is using 1000

		assert.EqualValues(t, 4, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0))
		assert.EqualValues(t, 8*1024*1024*1024, (*containers)[0].Resources.Requests.Memory().ScaledValue(0))

		assert.EqualValues(t, 8, (*containers)[0].Resources.Limits.Cpu().ScaledValue(0))
		assert.EqualValues(t, 16*1024*1024*1024, (*containers)[0].Resources.Limits.Memory().ScaledValue(0))

		assert.True(t, testlib.ArgContains((*containers)[0].Args, "mem 8Gi"))

		// make sure the replica counts are correct
		if testlib.IsStatefulSetHotCopyEnabled(&obj) {
			assert.EqualValues(t, 1, *obj.Spec.Replicas)
			foundBackupEnabled = true
		} else {
			assert.Zero(t, *obj.Spec.Replicas)
			foundBackupDisabled = true
		}
	}

	assert.True(t, foundBackupEnabled)
	assert.True(t, foundBackupDisabled)
}

func TestResourcesDatabaseOverridden(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	noHotCopyReplicas := 2
	hotcopyReplicas := 1

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    "1",
			"database.sm.resources.requests.memory": "4Gi",
			"database.sm.noHotCopy.replicas":        strconv.Itoa(noHotCopyReplicas),
			"database.sm.hotCopy.replicas":          strconv.Itoa(hotcopyReplicas),
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	foundBackupEnabled := false
	foundBackupDisabled := false

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
		containers := &obj.Spec.Template.Spec.Containers

		assert.NotEmpty(t, containers)
		assert.NotNil(t, (*containers)[0].Resources.Limits)
		assert.NotNil(t, (*containers)[0].Resources.Requests)

		// the memory is confusing Gi with G. We are using power of two (1024). But scaled value is using 1000

		assert.EqualValues(t, 1, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0))
		assert.EqualValues(t, 4*1024*1024*1024, (*containers)[0].Resources.Requests.Memory().ScaledValue(0))

		assert.EqualValues(t, 8, (*containers)[0].Resources.Limits.Cpu().ScaledValue(0))
		assert.EqualValues(t, 16*1024*1024*1024, (*containers)[0].Resources.Limits.Memory().ScaledValue(0))

		assert.True(t, testlib.ArgContains((*containers)[0].Args, "mem 4Gi"))

		// make sure the replica counts are correct
		if testlib.IsStatefulSetHotCopyEnabled(&obj) {
			assert.EqualValues(t, hotcopyReplicas, *obj.Spec.Replicas)
			foundBackupEnabled = true
		} else {
			assert.EqualValues(t, noHotCopyReplicas, *obj.Spec.Replicas)
			foundBackupDisabled = true
		}
	}

	assert.True(t, foundBackupEnabled)
	assert.True(t, foundBackupDisabled)
}

func TestPullSecretsRenderAllNuoDB(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"nuodb.image.pullSecrets":              "{fooBar}",
			"admin.legacy.loadBalancerJob.enabled": "true",
		},
	}

	helm.RenderTemplate(t, options, "../../stable/admin", "release-name", []string{"templates/job.yaml"})
	helm.RenderTemplate(t, options, "../../stable/admin", "release-name", []string{"templates/statefulset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/statefulset.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/deployment.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/cronjob.yaml"})

	helm.RenderTemplate(t, options, "../../stable/transparent-hugepage", "release-name", []string{"templates/daemonset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/restore", "release-name", []string{"templates/job.yaml"})
}

func TestPullSecretsRenderAllGlobal(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"global.imagePullSecrets":              "{fooBar}",
			"admin.legacy.loadBalancerJob.enabled": "true",
		},
	}

	helm.RenderTemplate(t, options, "../../stable/admin", "release-name", []string{"templates/job.yaml"})
	helm.RenderTemplate(t, options, "../../stable/admin", "release-name", []string{"templates/statefulset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/statefulset.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/deployment.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", "release-name", []string{"templates/cronjob.yaml"})

	helm.RenderTemplate(t, options, "../../stable/transparent-hugepage", "release-name", []string{"templates/daemonset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/restore", "release-name", []string{"templates/job.yaml"})
}

func TestPingTimeout(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			assert.True(t, listContains(obj.Spec.Template.Spec.Containers[0].Args, "ping-timeout"))
		}
	})
}

func TestSpecificOptions(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.engineOptions.verbose": "advanced-txn",
			"database.sm.engineOptions.verbose": "advanced-txn",
		},
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			assert.True(t, listContains(obj.Spec.Template.Spec.Containers[0].Args, "advanced-txn"))
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			assert.True(t, listContains(obj.Spec.Template.Spec.Containers[0].Args, "advanced-txn"))
		}
	})
}

func TestDatabaseBackupDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.enablePod": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.False(t, testlib.IsStatefulSetHotCopyEnabled(&obj), "Found stateful set with backup enabled")
	}
}

func TestDatabaseNonBackupDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.noHotCopy.enablePod": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.True(t, testlib.IsStatefulSetHotCopyEnabled(&obj), "Found stateful set with backup disabled")
	}
}

func assertExpectedLines(t *testing.T, optionsMap *map[string]string, helmChartName string, templateNames []string, expectedLines *map[string]int) {
	options := &helm.Options{
		SetValues: *optionsMap,
	}

	output, err := helm.RenderTemplateE(t, options, "../../stable/"+helmChartName, "release-name", templateNames)

	if expectedLines == nil {
		assert.Error(t, err)
		return
	} else {
		assert.NoError(t, err)
	}

	actualLines := make(map[string]int)
	// iterate through all lines of rendered output, removing any trailing spaces
	for _, line := range regexp.MustCompile(" *\n").Split(output, -1) {
		if count, ok := (*expectedLines)[line]; ok {
			assert.True(t, count != 0, "Unexpected line: "+line)
			actualLines[line]++
		}
	}

	for line, cnt := range *expectedLines {
		assert.Equal(t, cnt, actualLines[line], "Unexpected number of occurrences of "+line)
	}
}

func TestAddRoleBindingEnabled(t *testing.T) {
	optionsMap := map[string]string{}
	templateNames := []string{
		"templates/role.yaml",
		"templates/rolebinding.yaml",
		"templates/serviceaccount.yaml",
		"templates/statefulset.yaml",
	}
	// expect Role, RoleBinding, and ServiceAccount to be created
	expectedLines := map[string]int{
		"kind: Role":                      1,
		"kind: RoleBinding":               1,
		"kind: ServiceAccount":            1,
		"      serviceAccountName: nuodb": 1,
	}
	assertExpectedLines(t, &optionsMap, "admin", templateNames, &expectedLines)
}

func TestAddRoleBindingDisabled(t *testing.T) {
	// disable creation of role and role binding
	optionsMap := map[string]string{
		"nuodb.serviceAccount": "default",
		"nuodb.addRoleBinding": "false",
	}
	templateNames := []string{
		"templates/role.yaml",
		"templates/rolebinding.yaml",
		"templates/serviceaccount.yaml",
		"templates/statefulset.yaml",
	}

	assertExpectedLines(t, &optionsMap, "admin", templateNames, nil)
}

func TestDeploymentServiceAccount(t *testing.T) {
	optionsMap := map[string]string{}
	templateNames := []string{
		"templates/deployment.yaml",
	}
	expectedLines := map[string]int{
		"      serviceAccountName: nuodb": 1,
	}
	assertExpectedLines(t, &optionsMap, "database", templateNames, &expectedLines)
}

func TestStatefulSetServiceAccount(t *testing.T) {
	optionsMap := map[string]string{}
	templateNames := []string{
		"templates/statefulset.yaml",
	}
	// there should be two serviceAccountName declarations, for SM and hotcopy SM
	expectedLines := map[string]int{
		"kind: StatefulSet":               2,
		"      serviceAccountName: nuodb": 2,
	}
	assertExpectedLines(t, &optionsMap, "database", templateNames, &expectedLines)
}