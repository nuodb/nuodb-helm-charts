package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func checkSidecarContainers(t *testing.T, containers []v1.Container, options *helm.Options, chartPath string) {
	require.NotEmpty(t, containers)
	found := 0
	securityContextEnabled := options.SetValues["admin.securityContext.enabledOnContainer"] == "true" ||
		options.SetValues["database.securityContext.enabledOnContainer"] == "true"

	for _, container := range containers {
		t.Logf("Inspecting container %s in chart %s", container.Name, chartPath)
		if container.Name == "nuocollector" {
			found++
			assert.Equal(t, container.Image, fmt.Sprintf("%s/%s:%s",
				options.SetValues["nuocollector.image.registry"],
				options.SetValues["nuocollector.image.repository"],
				options.SetValues["nuocollector.image.tag"]))
		} else if container.Name == "nuocollector-config" {
			found++
			assert.Equal(t, container.Image, fmt.Sprintf("%s/%s:%s",
				options.SetValues["nuocollector.watcher.registry"],
				options.SetValues["nuocollector.watcher.repository"],
				options.SetValues["nuocollector.watcher.tag"]))
			assert.Contains(t, container.Env, v1.EnvVar{
				Name:  "FOLDER",
				Value: "/etc/telegraf/telegraf.d/dynamic/",
			})
			assert.Contains(t, container.Env, v1.EnvVar{
				Name:  "REQ_URL",
				Value: "http://127.0.0.1:5000/reload",
			})
			if chartPath == testlib.ADMIN_HELM_CHART_PATH {
				assert.Contains(t, container.Env, v1.EnvVar{
					Name:  "LABEL",
					Value: "nuodb.com/nuocollector-plugin in (release-name-nuodb-cluster0-admin, insights)",
				})
			} else {
				assert.Contains(t, container.Env, v1.EnvVar{
					Name:  "LABEL",
					Value: "nuodb.com/nuocollector-plugin in (release-name-nuodb-cluster0-demo-database, insights)",
				})
			}
		} else {
			// This is probably the main container
			continue
		}
		assert.Contains(t, container.VolumeMounts, v1.VolumeMount{
			Name:      "nuocollector-config",
			MountPath: "/etc/telegraf/telegraf.d/dynamic/",
		})
		assert.Contains(t, container.VolumeMounts, v1.VolumeMount{
			Name:      "log-volume",
			MountPath: "/var/log/nuodb",
		})
		if securityContextEnabled {
			assert.NotNil(t, container.SecurityContext)
		} else {
			assert.Nil(t, container.SecurityContext)
		}
		testlib.AssertResourceValue(t, options, "nuocollector.resources.limits.cpu", container.Resources.Limits.Cpu())
		testlib.AssertResourceValue(t, options, "nuocollector.resources.limits.memory", container.Resources.Limits.Memory())
		testlib.AssertResourceValue(t, options, "nuocollector.resources.requests.cpu", container.Resources.Requests.Cpu())
		testlib.AssertResourceValue(t, options, "nuocollector.resources.requests.memory", container.Resources.Requests.Memory())
	}

	expectedContainersCount := 0
	if options.SetValues["nuocollector.enabled"] == "true" {
		expectedContainersCount = 2
	}
	assert.Equal(t, expectedContainersCount, found)
}

func checkSpecVolumes(t *testing.T, volumes []v1.Volume, options *helm.Options, chartPath string) {
	require.NotEmpty(t, volumes)
	found := false
	for _, volume := range volumes {
		if volume.Name == "nuocollector-config" {
			found = true
			// Check that empty dir is mounted
			assert.NotNil(t, volume.EmptyDir)
		}
	}
	if options.SetValues["nuocollector.enabled"] == "true" {
		assert.True(t, found, "nuocollector-config should be declared as volume")
	} else {
		assert.False(t, found, "nuocollector-config is declared as volume with nuocollector disabled")
	}
}

func checkPluginsRendered(t *testing.T, configMaps []v1.ConfigMap, options *helm.Options, chartPath string, expectedNrPlugins int) {
	found := 0
	for _, cm := range configMaps {
		if labelValue, ok := cm.Labels["nuodb.com/nuocollector-plugin"]; ok {
			found++
			if chartPath == testlib.ADMIN_HELM_CHART_PATH {
				assert.Equal(t, "release-name-nuodb-cluster0-admin", labelValue)
			} else {
				assert.Equal(t, "release-name-nuodb-cluster0-demo-database", labelValue)
			}
			parts := strings.Split(cm.Name, "-")
			assert.Greater(t, len(parts), 1)
			pluginName := parts[len(parts)-1]
			var expectedData string
			ok := false
			if chartPath == testlib.ADMIN_HELM_CHART_PATH {
				expectedData, ok = options.SetValues["nuocollector.plugins.admin."+pluginName]
			} else {
				expectedData, ok = options.SetValues["nuocollector.plugins.database."+pluginName]
			}
			if ok {
				// Check content only for plugins specified in options
				assert.NotEmpty(t, expectedData)
				assert.Equal(t, expectedData, cm.Data[pluginName+".conf"])
			}
		}
	}
	assert.Equal(t, expectedNrPlugins, found)
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

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			t.Logf("Inspecting database statefulset: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})

	t.Run("testDatabaseDeploymentSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database deployment
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			t.Logf("Inspecting database deployment: %s", obj.Name)
			checkSpecVolumes(t, obj.Spec.Template.Spec.Volumes, options, helmChartPath)
			checkSidecarContainers(t, obj.Spec.Template.Spec.Containers, options, helmChartPath)
		}
	})
}

func TestNuoDBCollectorSidecarsEnabled(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"nuocollector.enabled":            "true",
			"nuocollector.image.registry":     "docker.io",
			"nuocollector.image.repository":   "nuodb/nuocd",
			"nuocollector.image.tag":          "1.0.0",
			"nuocollector.watcher.registry":   "docker.io",
			"nuocollector.watcher.repository": "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":        "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestNuoDBCollectorSidecarsDisabled(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"nuocollector.enabled":            "false",
			"nuocollector.image.registry":     "docker.io",
			"nuocollector.image.repository":   "nuodb/nuocd",
			"nuocollector.image.tag":          "1.0.0",
			"nuocollector.watcher.registry":   "docker.io",
			"nuocollector.watcher.repository": "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":        "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestNuoDBCollectorPluginsRendered(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"nuocollector.enabled":               "true",
			"nuocollector.image.registry":        "docker.io",
			"nuocollector.image.repository":      "nuodb/nuocd",
			"nuocollector.image.tag":             "1.0.0",
			"nuocollector.watcher.registry":      "docker.io",
			"nuocollector.watcher.repository":    "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":           "latest",
			"nuocollector.plugins.admin.file":    "[[outputs.file]]\nfiles = [\"/var/log/nuodb/metrics.log\"]\ndata_format = \"json\"",
			"nuocollector.plugins.database.file": "[[outputs.file]]\nfiles = [\"/var/log/nuodb/metrics.log\"]\ndata_format = \"json\"",
		},
	}

	t.Run("testAdminPlugins", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect admin nuocollector configMaps
		helmChartPath := testlib.ADMIN_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/nuocollector-configmap.yaml"})
		configMaps := testlib.SplitAndRenderConfigMap(t, output, 1)
		// Check that default and custom plugins are rendered
		checkPluginsRendered(t, configMaps, options, helmChartPath, 1)
	})

	t.Run("testDatabaseStatefulsetSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database nuocollector configMaps
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/nuocollector-configmap.yaml"})
		configMaps := testlib.SplitAndRenderConfigMap(t, output, 1)
		// Check that default and custom plugins are rendered
		checkPluginsRendered(t, configMaps, options, helmChartPath, 1)

	})
}

func TestNuoDBCollectorSidecarsSecurityContext(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.securityContext.enabledOnContainer":    "true",
			"database.securityContext.enabledOnContainer": "true",
			"nuocollector.enabled":                        "true",
			"nuocollector.image.registry":                 "docker.io",
			"nuocollector.image.repository":               "nuodb/nuocd",
			"nuocollector.image.tag":                      "1.0.0",
			"nuocollector.watcher.registry":               "docker.io",
			"nuocollector.watcher.repository":             "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":                    "latest",
		},
	}
	executeSidecarTests(t, options)
}

func TestNuoDBCollectorResources(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"admin.securityContext.enabledOnContainer":    "true",
			"database.securityContext.enabledOnContainer": "true",
			"nuocollector.enabled":                        "true",
			"nuocollector.image.registry":                 "docker.io",
			"nuocollector.image.repository":               "nuodb/nuocd",
			"nuocollector.image.tag":                      "1.0.0",
			"nuocollector.watcher.registry":               "docker.io",
			"nuocollector.watcher.repository":             "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":                    "latest",
			"nuocollector.resources.limits.cpu":           "100m",
			"nuocollector.resources.limits.memory":        "128Mi",
			"nuocollector.resources.requests.cpu":         "100m",
			"nuocollector.resources.requests.memory":      "128Mi",
		},
	}
	executeSidecarTests(t, options)
}
