package integration

import (
	"strconv"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"gotest.tools/assert"
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)
		assert.Check(t, (*containers)[0].Resources.Limits == nil)
		assert.Check(t, (*containers)[0].Resources.Requests == nil)
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	partCounter := 0
	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		partCounter++

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)
		assert.Check(t, (*containers)[0].Resources.Limits == nil)
		assert.Assert(t, (*containers)[0].Resources.Requests != nil)

		assert.Check(t, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0) == 1)
		assert.Check(t, (*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga) == 4,
			(*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga))

	}

	assert.Equal(t, partCounter, 1)
}

func TestResourcesDatabaseDefaults(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	partCounter := 0
	foundBackupEnabled := false
	foundBackupDisabled := false

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)
		assert.Assert(t, (*containers)[0].Resources.Limits != nil)
		assert.Assert(t, (*containers)[0].Resources.Requests != nil)

		// the memory is confusing Gi with G. We are using power of two (1024). But scaled value is using 1000

		assert.Check(t, (*containers)[0].Resources.Limits.Cpu().ScaledValue(0) == 8)
		assert.Check(t, (*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga) == 18,
			(*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga))

		assert.Check(t, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0) == 4)
		assert.Check(t, (*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga) == 9,
			(*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga))

		// make sure the replica counts are correct
		if testlib.IsStatefulSetHotCopyEnabled(&ss) {
			assert.Check(t, *ss.Spec.Replicas == 1)
			foundBackupEnabled = true
		} else {
			assert.Check(t, *ss.Spec.Replicas == 0)
			foundBackupDisabled = true
		}
	}

	assert.Check(t, foundBackupEnabled)
	assert.Check(t, foundBackupDisabled)

	assert.Equal(t, partCounter, 2)
}

func TestResourcesDatabaseOverridden(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	noHotCopyReplicas := 2
	hotcopyReplicas := 1

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    "1",
			"database.sm.resources.requests.memory": "4G",
			"database.sm.noHotCopy.replicas":        strconv.Itoa(noHotCopyReplicas),
			"database.sm.hotCopy.replicas":          strconv.Itoa(hotcopyReplicas),
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	partCounter := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)
		assert.Assert(t, (*containers)[0].Resources.Limits != nil)
		assert.Assert(t, (*containers)[0].Resources.Requests != nil)

		// the memory is confusing Gi with G. We are using power of two (1024). But scaled value is using 1000

		assert.Check(t, (*containers)[0].Resources.Limits.Cpu().ScaledValue(0) == 8)
		assert.Check(t, (*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga) == 18,
			(*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga))

		assert.Check(t, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0) == 1)
		assert.Check(t, (*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga) == 4,
			(*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga))

		// make sure the replica counts are correct
		if testlib.IsStatefulSetHotCopyEnabled(&ss) {
			assert.Check(t, *ss.Spec.Replicas == int32(hotcopyReplicas))
		} else {
			assert.Check(t, *ss.Spec.Replicas == int32(noHotCopyReplicas))
		}
	}

	assert.Equal(t, partCounter, 2)
}

func TestPullSecretsRenderAllNuoDB(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{"nuodb.image.pullSecrets": "{fooBar}"},
	}

	helm.RenderTemplate(t, options, "../../stable/admin", []string{"templates/job.yaml"})
	helm.RenderTemplate(t, options, "../../stable/admin", []string{"templates/statefulset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/database", []string{"templates/statefulset.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", []string{"templates/deployment.yaml"})

	helm.RenderTemplate(t, options, "../../stable/transparent-hugepage", []string{"templates/daemonset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/monitoring-influx", []string{"templates/deployment.yaml"})

	helm.RenderTemplate(t, options, "../../stable/monitoring-insights", []string{"templates/pod.yaml"})

	helm.RenderTemplate(t, options, "../../stable/backup", []string{"templates/cronjob.yaml"})
	helm.RenderTemplate(t, options, "../../stable/backup", []string{"templates/job.yaml"})

	helm.RenderTemplate(t, options, "../../stable/restore", []string{"templates/job.yaml"})
}

func TestPullSecretsRenderAllGlobal(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{"global.imagePullSecrets": "{fooBar}"},
	}

	helm.RenderTemplate(t, options, "../../stable/admin", []string{"templates/job.yaml"})
	helm.RenderTemplate(t, options, "../../stable/admin", []string{"templates/statefulset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/database", []string{"templates/statefulset.yaml"})
	helm.RenderTemplate(t, options, "../../stable/database", []string{"templates/deployment.yaml"})

	helm.RenderTemplate(t, options, "../../stable/transparent-hugepage", []string{"templates/daemonset.yaml"})

	helm.RenderTemplate(t, options, "../../stable/monitoring-influx", []string{"templates/deployment.yaml"})

	helm.RenderTemplate(t, options, "../../stable/monitoring-insights", []string{"templates/pod.yaml"})

	helm.RenderTemplate(t, options, "../../stable/backup", []string{"templates/cronjob.yaml"})
	helm.RenderTemplate(t, options, "../../stable/backup", []string{"templates/job.yaml"})

	helm.RenderTemplate(t, options, "../../stable/restore", []string{"templates/job.yaml"})
}

func TestPingTimeoutSetStatefulSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)

		assert.Check(t, listContains((*containers)[0].Args, "ping-timeout"))
	}
}

func TestPingTimeoutSetDaemonSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		var ss appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)

		assert.Check(t, listContains((*containers)[0].Args, "ping-timeout"))
	}
}

func TestSpecificOptionsDeployment(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.te.engineOptions.verbose": "advanced-txn"},
	}

	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

	var dep appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &dep)

	containers := &dep.Spec.Template.Spec.Containers
	assert.Assert(t, len(*containers) >= 1)
	assert.Check(t, listContains((*containers)[0].Args, "advanced-txn"))
}

func TestSpecificOptionsStatefulSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.sm.engineOptions.verbose": "advanced-txn"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)

		assert.Check(t, listContains((*containers)[0].Args, "advanced-txn"))
	}
}

func TestSpecificOptionsDaemonSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true",
			"database.sm.engineOptions.verbose": "advanced-txn",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		var ss appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)

		assert.Check(t, listContains((*containers)[0].Args, "advanced-txn"))
	}
}

func TestResourcesDaemonSetsDefaults(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	partCounter := 0
	foundBackupEnabled := false
	foundBackupDisabled := false

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		containers := &ss.Spec.Template.Spec.Containers

		assert.Assert(t, len(*containers) >= 1)
		assert.Assert(t, (*containers)[0].Resources.Limits != nil)
		assert.Assert(t, (*containers)[0].Resources.Requests != nil)

		// the memory is confusing Gi with G. We are using power of two (1024). But scaled value is using 1000

		assert.Check(t, (*containers)[0].Resources.Limits.Cpu().ScaledValue(0) == 8)
		assert.Check(t, (*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga) == 18,
			(*containers)[0].Resources.Limits.Memory().ScaledValue(resource.Giga))

		assert.Check(t, (*containers)[0].Resources.Requests.Cpu().ScaledValue(0) == 4)
		assert.Check(t, (*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga) == 9,
			(*containers)[0].Resources.Requests.Memory().ScaledValue(resource.Giga))

		if testlib.IsDaemonSetHotCopyEnabled(&ss) {
			foundBackupEnabled = true
		} else {
			foundBackupDisabled = true
		}

	}

	assert.Check(t, foundBackupEnabled)
	assert.Check(t, foundBackupDisabled)

	assert.Equal(t, partCounter, 2)
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	partCounter := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		assert.Check(t, !testlib.IsStatefulSetHotCopyEnabled(&ss), "Found stateful set with backup enabled")
	}

	assert.Equal(t, partCounter, 1)
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	partCounter := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		assert.Check(t, testlib.IsStatefulSetHotCopyEnabled(&ss), "Found stateful set with backup enabled")
	}

	assert.Equal(t, partCounter, 1)
}

func TestDatabaseBackupDisabledDaemonSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.enableDaemonSet":      "true",
			"database.sm.hotCopy.enablePod": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	partCounter := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		assert.Check(t, !testlib.IsDaemonSetHotCopyEnabled(&ss), "Found daemon set with backup enabled")
	}

	// with daemonSet
	assert.Equal(t, partCounter, 1)
}

func TestDatabaseNoBackupDisabledDaemonSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.enableDaemonSet":        "true",
			"database.sm.noHotCopy.enablePod": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	partCounter := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {

		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		partCounter += 1

		var ss appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &ss)

		assert.Check(t, testlib.IsDaemonSetHotCopyEnabled(&ss), "Found daemon set with backup enabled")
	}

	// with daemonSet
	assert.Equal(t, partCounter, 1)
}
