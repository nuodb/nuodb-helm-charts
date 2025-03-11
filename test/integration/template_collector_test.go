package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func checkSidecarContainers(t *testing.T, containers []corev1.Container, options *helm.Options, chartPath string) {
	require.NotEmpty(t, containers)
	found := 0
	securityContextEnabled := options.SetValues["admin.securityContext.enabledOnContainer"] == "true" ||
		options.SetValues["database.securityContext.enabledOnContainer"] == "true"
	logPersistenceEnabled := options.SetValues["admin.logPersistence.enabled"] == "true" ||
		options.SetValues["database.sm.logPersistence.enabled"] == "true" ||
		options.SetValues["database.te.logPersistence.enabled"] == "true"

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
			assert.Contains(t, container.Env, corev1.EnvVar{
				Name:  "FOLDER",
				Value: "/etc/telegraf/telegraf.d/dynamic/",
			})
			assert.Contains(t, container.Env, corev1.EnvVar{
				Name:  "REQ_URL",
				Value: "http://127.0.0.1:5000/reload",
			})
			if chartPath == testlib.ADMIN_HELM_CHART_PATH {
				assert.Contains(t, container.Env, corev1.EnvVar{
					Name:  "LABEL",
					Value: "nuodb.com/nuocollector-plugin in (release-name-nuodb-cluster0-admin, insights)",
				})
			} else {
				assert.Contains(t, container.Env, corev1.EnvVar{
					Name:  "LABEL",
					Value: "nuodb.com/nuocollector-plugin in (release-name-nuodb-cluster0-demo-database, insights)",
				})
			}
			assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
				Name:      "eph-volume",
				MountPath: "/tmp",
				SubPath:   "tmp-watcher",
			})
		} else {
			// This is probably the main container
			continue
		}
		assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
			Name:      "eph-volume",
			MountPath: "/etc/telegraf/telegraf.d/dynamic/",
			SubPath:   "telegraf",
		})
		if logPersistenceEnabled {
			assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
				Name:      "log-volume",
				MountPath: "/var/log/nuodb",
			})
		} else {
			assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
				Name:      "eph-volume",
				MountPath: "/var/log/nuodb",
				SubPath:   "log",
			})
		}
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

func checkSpecVolumes(t *testing.T, volumes []corev1.Volume, options *helm.Options, chartPath string) {
	if options.SetValues["nuocollector.enabled"] == "false" {
		return
	}
	// Check that ephemeral volume is enabled. This is needed to share telegraf config.
	for _, volume := range volumes {
		if volume.Name == "eph-volume" {
			return
		}
	}
	assert.Fail(t, "eph-volume should be declared as volume")
}

func checkPluginsRendered(t *testing.T, configMaps []corev1.ConfigMap, options *helm.Options, chartPath string, expectedNrPlugins int) {
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

type assertFunction func(chartPath string, volumes []corev1.Volume, containers []corev1.Container)

func executeSidecarTestsWithAsserts(t *testing.T, options *helm.Options, assertFn assertFunction) {
	t.Run("testAdminSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect admin statefulset
		helmChartPath := testlib.ADMIN_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.ADMIN_HELM_CHART_PATH, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			t.Logf("Inspecting admin statefulset: %s", obj.Name)
			assertFn(helmChartPath, obj.Spec.Template.Spec.Volumes, obj.Spec.Template.Spec.Containers)
		}
	})

	t.Run("testDatabaseStatefulsetSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database statefulset
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			t.Logf("Inspecting database statefulset: %s", obj.Name)
			assertFn(helmChartPath, obj.Spec.Template.Spec.Volumes, obj.Spec.Template.Spec.Containers)
		}
	})

	t.Run("testDatabaseDeploymentSidecars", func(t *testing.T) {
		// Run RenderTemplate to render the template and inspect database deployment
		helmChartPath := testlib.DATABASE_HELM_CHART_PATH
		output := helm.RenderTemplate(t, options, testlib.DATABASE_HELM_CHART_PATH, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			t.Logf("Inspecting database deployment: %s", obj.Name)
			assertFn(helmChartPath, obj.Spec.Template.Spec.Volumes, obj.Spec.Template.Spec.Containers)
		}
	})
}

func executeSidecarTests(t *testing.T, options *helm.Options) {
	executeSidecarTestsWithAsserts(t, options, func(chartPath string, volumes []corev1.Volume, containers []corev1.Container) {
		checkSpecVolumes(t, volumes, options, chartPath)
		checkSidecarContainers(t, containers, options, chartPath)
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

func TestNuoDBCollectorSidecarsLogPersistence(t *testing.T) {

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.logPersistence.enabled":       "true",
			"database.sm.logPersistence.enabled": "true",
			"database.te.logPersistence.enabled": "true",
			"nuocollector.enabled":               "true",
			"nuocollector.image.registry":        "docker.io",
			"nuocollector.image.repository":      "nuodb/nuocd",
			"nuocollector.image.tag":             "1.0.0",
			"nuocollector.watcher.registry":      "docker.io",
			"nuocollector.watcher.repository":    "kiwigrid/k8s-sidecar",
			"nuocollector.watcher.tag":           "latest",
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

func TestNuoDBCollectorEnv(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"nuocollector.enabled":                            "true",
			"nuocollector.env[0].name":                        "foo",
			"nuocollector.env[0].value":                       "bar",
			"nuocollector.env[1].name":                        "baz",
			"nuocollector.env[1].valueFrom.secretKeyRef.name": "telegraf-creds",
			"nuocollector.env[1].valueFrom.secretKeyRef.key":  "password",
		},
	}
	executeSidecarTestsWithAsserts(t, options, func(chartPath string, volumes []corev1.Volume, containers []corev1.Container) {
		if chartPath != testlib.DATABASE_HELM_CHART_PATH {
			// environment variables are supported only for the database chart
			return
		}
		var nuocollector *corev1.Container
		for i := range containers {
			if containers[i].Name == "nuocollector" {
				nuocollector = &containers[i]
			}
		}
		assert.NotNil(t, nuocollector)
		testlib.AssertEnvContains(t, nuocollector.Env, "foo", "bar")
		testlib.AssertEnvContainsValueFrom(t, nuocollector.Env, "baz", corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "telegraf-creds",
				},
				Key: "password",
			},
		})
	})
}

func TestNuoDBCollectorPorts(t *testing.T) {
	options := &helm.Options{
		SetValues: map[string]string{
			"nuocollector.enabled":                "true",
			"nuocollector.ports[0].name":          "http-metrics",
			"nuocollector.ports[0].containerPort": "9001",
			"nuocollector.ports[0].protocol":      "TCP",
		},
	}
	executeSidecarTestsWithAsserts(t, options, func(chartPath string, volumes []corev1.Volume, containers []corev1.Container) {
		if chartPath != testlib.DATABASE_HELM_CHART_PATH {
			// ports are supported only for the database chart
			return
		}
		var nuocollector *corev1.Container
		for i := range containers {
			if containers[i].Name == "nuocollector" {
				nuocollector = &containers[i]
			}
		}
		assert.NotNil(t, nuocollector)
		assert.Contains(t, nuocollector.Ports, corev1.ContainerPort{
			Name:          "http-metrics",
			ContainerPort: 9001,
			Protocol:      corev1.ProtocolTCP,
		})
	})
}
