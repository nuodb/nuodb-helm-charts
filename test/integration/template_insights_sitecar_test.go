package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"
)

func checkSidecarContainers(t *testing.T, containers []v1.Container, options *helm.Options, chartPath string) {
	assert.NotEmpty(t, containers)
	adminPluginsLabel, adminPluginsLabelOK := options.SetValues["admin.insightsPluginsLabel"]
	databasePluginsLabel, databasePluginsLabelOK := options.SetValues["database.insightsPluginsLabel"]
	found := 0

	for _, container := range containers {
		t.Logf("Inspecting container %s in chart %s", container.Name, chartPath)
		if container.Name == "insights" {
			found++
			assert.Equal(t, container.Image, fmt.Sprintf("%s/%s:%s",
				options.SetValues["insights.image.registry"],
				options.SetValues["insights.image.repository"],
				options.SetValues["insights.image.tag"]))
		} else if container.Name == "insights-config" {
			found++
			assert.Equal(t, container.Image, fmt.Sprintf("%s/%s:%s",
				options.SetValues["insights.watcher.registry"],
				options.SetValues["insights.watcher.repository"],
				options.SetValues["insights.watcher.tag"]))
			assert.Contains(t, container.Env, v1.EnvVar{
				Name:  "FOLDER",
				Value: "/etc/telegraf/telegraf.d/",
			})
			assert.Contains(t, container.Env, v1.EnvVar{
				Name:  "REQ_URL",
				Value: "http://127.0.0.1:5000/reload",
			})
			if chartPath == testlib.ADMIN_HELM_CHART_PATH {
				assert.Contains(t, container.Env, v1.EnvVar{
					Name:  "LABEL",
					Value: adminPluginsLabel,
				})
			} else {
				assert.Contains(t, container.Env, v1.EnvVar{
					Name:  "LABEL",
					Value: databasePluginsLabel,
				})
			}
		} else {
			// This is probably the main container
			continue
		}
		assert.Contains(t, container.VolumeMounts, v1.VolumeMount{
			Name:      "insights-config",
			MountPath: "/etc/telegraf/telegraf.d/",
		})
		assert.Contains(t, container.VolumeMounts, v1.VolumeMount{
			Name:      "log-volume",
			MountPath: "/var/log/nuodb",
		})
	}

	expectedContainersCount := 0
	if options.SetValues["insights.enabled"] == "true" && chartPath == testlib.ADMIN_HELM_CHART_PATH {
		if adminPluginsLabelOK && adminPluginsLabel != "" {
			expectedContainersCount = 2
		}
	} else if options.SetValues["insights.enabled"] == "true" && chartPath == testlib.DATABASE_HELM_CHART_PATH {
		if databasePluginsLabelOK && databasePluginsLabel != "" {
			expectedContainersCount = 2
		}
	}
	assert.Equal(t, expectedContainersCount, found)
}

func checkSpecVolumes(t *testing.T, volumes []v1.Volume, options *helm.Options, chartPath string) {
	assert.NotEmpty(t, volumes)
	adminPluginsLabel, adminPluginsLabelOK := options.SetValues["admin.insightsPluginsLabel"]
	databasePluginsLabel, databasePluginsLabelOK := options.SetValues["database.insightsPluginsLabel"]
	found := false
	for _, volume := range volumes {
		if volume.Name == "insights-config" {
			found = true
			// Check that empty dir is mounted
			assert.NotNil(t, volume.EmptyDir)
		}
	}
	if options.SetValues["insights.enabled"] == "true" &&
		((adminPluginsLabelOK && chartPath == testlib.ADMIN_HELM_CHART_PATH) || (databasePluginsLabelOK && chartPath == testlib.DATABASE_HELM_CHART_PATH)) &&
		((adminPluginsLabel != "" && chartPath == testlib.ADMIN_HELM_CHART_PATH) || (databasePluginsLabel != "" && chartPath == testlib.DATABASE_HELM_CHART_PATH)) {
		assert.True(t, found)
	} else {
		assert.False(t, found)
	}
}

func executeSidecarTests(t *testing.T, options *helm.Options) {
	t.Run("testAdminSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect admin statefulset
		helmChartPath := testlib.ADMIN_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.ADMIN_HELM_CHART_PATH, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			t.Logf("Inspecting admin statefulset: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})

	t.Run("testDatabaseStatefulsetSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database statefulset
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			t.Logf("Inspecting database statefulset: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})

	t.Run("testDatabaseDeploymentSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database statefulset
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			t.Logf("Inspecting database deployment: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})

	t.Run("testDatabaseDaemonsetSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database daemonset
		options.SetValues["database.enableDaemonSet"] = "true"
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range testlib.SplitAndRenderDaemonSet(t, output, 1) {
			t.Logf("Inspecting database daemonset: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})
}

func TestInsightsSidecarsEnabled(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"insights.enabled":              "true",
			"admin.insightsPluginsLabel":    "admin-plugin-label",
			"database.insightsPluginsLabel": "database-plugin-label",
			"insights.image.registry":       "docker.io",
			"insights.image.repository":     "nuodb/nuocd",
			"insights.image.tag":            "1.0.0",
			"insights.watcher.registry":     "docker.io",
			"insights.watcher.repository":   "kiwigrid/k8s-sidecar",
			"insights.watcher.tag":          "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestInsightsSidecarsDisabled(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"insights.enabled":              "false",
			"admin.insightsPluginsLabel":    "admin-plugin-label",
			"database.insightsPluginsLabel": "database-plugin-label",
			"insights.image.registry":       "docker.io",
			"insights.image.repository":     "nuodb/nuocd",
			"insights.image.tag":            "1.0.0",
			"insights.watcher.registry":     "docker.io",
			"insights.watcher.repository":   "kiwigrid/k8s-sidecar",
			"insights.watcher.tag":          "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestInsightsSidecarsAdminPluginLabelUnset(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"insights.enabled":              "true",
			"admin.insightsPluginsLabel":    "",
			"database.insightsPluginsLabel": "database-plugin-label",
			"insights.image.registry":       "docker.io",
			"insights.image.repository":     "nuodb/nuocd",
			"insights.image.tag":            "1.0.0",
			"insights.watcher.registry":     "docker.io",
			"insights.watcher.repository":   "kiwigrid/k8s-sidecar",
			"insights.watcher.tag":          "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestInsightsSidecarsDatabasePluginLabelUnset(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"insights.enabled":              "true",
			"admin.insightsPluginsLabel":    "admin-plugin-label",
			"database.insightsPluginsLabel": "",
			"insights.image.registry":       "docker.io",
			"insights.image.repository":     "nuodb/nuocd",
			"insights.image.tag":            "1.0.0",
			"insights.watcher.registry":     "docker.io",
			"insights.watcher.repository":   "kiwigrid/k8s-sidecar",
			"insights.watcher.tag":          "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestInsightsSidecarsBothPluginLabelsUnset(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"insights.enabled":              "true",
			"admin.insightsPluginsLabel":    "",
			"database.insightsPluginsLabel": "",
			"insights.image.registry":       "docker.io",
			"insights.image.repository":     "nuodb/nuocd",
			"insights.image.tag":            "1.0.0",
			"insights.watcher.registry":     "docker.io",
			"insights.watcher.repository":   "kiwigrid/k8s-sidecar",
			"insights.watcher.tag":          "latest",
		},
	}
	executeSidecarTests(t, options)
}
