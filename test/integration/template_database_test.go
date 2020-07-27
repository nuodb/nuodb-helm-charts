package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
)

func TestDatabaseSecretsDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/secret.yaml"})

	for _, obj := range SplitAndRenderSecret(t, output, 1) {
		assert.Contains(t, obj.StringData, "database-name")
		assert.Contains(t, obj.StringData, "database-password")
		assert.Contains(t, obj.StringData, "database-username")
	}

}

func TestDatabaseConfigMaps(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	configs := make(map[string]bool)

	for _, obj := range SplitAndRenderConfigMap(t, output, 3) {
		for k := range obj.Data {
			configs[k] = true
		}
	}

	assert.Contains(t, configs, "nuosm")
	assert.Contains(t, configs, "nuote")
	assert.Contains(t, configs, "readinessprobe")
}

func TestDatabaseDaemonSetDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

	// helm3 wont render an empty template
	assert.Error(t, err)
}

func TestDatabaseDaemonSetEnabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

	SplitAndRenderDaemonSet(t, output, 2)
}

func TestDatabaseClusterServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service-clusterip.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "demo-clusterip", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP,  obj.Spec.Type)
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.Empty(t, obj.Spec.ClusterIP)
	}
}

func TestDatabaseHeadlessServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service-headless.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "demo", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP,  obj.Spec.Type)
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.Equal(t, "None",  obj.Spec.ClusterIP)
	}
}

func TestDatabaseServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.provider":                        "amazon",
			"database.te.externalAccess.enabled":    "true",
			"database.te.externalAccess.internalIP": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "demo-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer,  obj.Spec.Type)
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.Contains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
	}

}

func TestDatabaseStatefulSetDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	_, err  := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	// helm3 wont render an empty template
	assert.Error(t, err)
}

func TestDatabaseStatefulSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
		assert.Equal(t, "sm", obj.Spec.Selector.MatchLabels["component"])
		assert.Equal(t, "sm", obj.Spec.Template.ObjectMeta.Labels["component"])
	}
}

func TestDatabaseStatefulSetVolumes(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.sm.logPersistence.enabled": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
		if strings.Contains(obj.Name, "-hotcopy") {
			assert.True(t, strings.Contains(obj.Spec.VolumeClaimTemplates[0].ObjectMeta.Name, "archive-volume"))
			assert.True(t, strings.Contains(obj.Spec.VolumeClaimTemplates[1].ObjectMeta.Name, "backup-volume"))
			assert.True(t, strings.Contains(obj.Spec.VolumeClaimTemplates[2].ObjectMeta.Name, "log-volume"))
		} else {
			assert.True(t, strings.Contains(obj.Spec.VolumeClaimTemplates[0].ObjectMeta.Name, "archive-volume"))
			assert.True(t, strings.Contains(obj.Spec.VolumeClaimTemplates[1].ObjectMeta.Name, "log-volume"))
		}
	}

}

func TestDatabaseDeploymentRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

	for _, obj := range SplitAndRenderDeployment(t, output, 1) {
		assert.Equal(t, "te", obj.Spec.Selector.MatchLabels["component"])
		assert.Equal(t, "te", obj.Spec.Template.ObjectMeta.Labels["component"])
	}

}

func TestDatabaseOtherOptions(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.otherOptions.keystore": "/etc/nuodb/keys/nuoadmin.p12",
			"database.sm.otherOptions.keystore": "/etc/nuodb/keys/nuoadmin.p12",
			"admin.tlsKeyStore.secret":          "nuodb-keystore",
			"admin.tlsKeyStore.key":             "nuoadmin.p12",
			"admin.tlsKeyStore.password":        "changeIt",
		},
	}

	basicArgChecks := func(args []string) {
		assert.True(t, ArgContains(args, "--keystore"))
		assert.True(t, ArgContains(args, "/etc/nuodb/keys/nuoadmin.p12"))
	}

	basicEnvChecks := func(args []v1.EnvVar) {
		assert.True(t, EnvContains(args, "NUODOCKER_KEYSTORE_PASSWORD", "changeIt"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range SplitAndRenderDaemonSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})
}

func TestDatabaseCustomEnv(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues:   map[string]string{},
		ValuesFiles: []string{"../files/database-env.yaml"},
	}

	basicEnvChecks := func(args []v1.EnvVar) {
		expectedAltAddress := v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		}
		assert.True(t, EnvContainsValueFrom(args, "NUODB_ALT_ADDRESS", &expectedAltAddress))
		assert.True(t, EnvContains(args, "CUSTOM_ENV_VAR", "CUSTOM_ENV_VAR_VALUE"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range SplitAndRenderDaemonSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})
}

func TestDatabaseVPNRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.securityContext.capabilities[0]": "NET_ADMIN",
			"database.envFrom.configMapRef[0]":         "test-config",
		},
	}

	basicChecks := func(args []v1.Container) {
		assert.Contains(t, args[0].SecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))
		assert.True(t, EnvFromSourceContains(args[0].EnvFrom, "test-config"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers)
		}
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range SplitAndRenderDaemonSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers)
		}
	})

}

func TestDatabaseLabeling(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.labels.cloud":  "minikube",
			"database.te.labels.region": "local",
			"database.te.labels.zone":   "local-b",
			"database.sm.labels.cloud":  "minikube",
			"database.sm.labels.region": "local",
			"database.sm.labels.zone":   "local-b",
		},
	}

	basicChecks := func(args []string) {
		assert.True(t, ArgContains(args, "cloud minikube"))
		assert.True(t, ArgContains(args, "region local"))
		assert.True(t, ArgContains(args, "zone local-b"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers[0].Args)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers[0].Args)

			if testlib.IsStatefulSetHotCopyEnabled(&obj) {
				assert.True(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup cluster0"))
			}
		}
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range SplitAndRenderDaemonSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers[0].Args)

			if testlib.IsDaemonSetHotCopyEnabled(&obj) {
				assert.True(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup cluster0"))
			}
		}
	})
}

func TestReadinessProbe(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	basicChecks := func(spec v1.PodSpec) {
		container := spec.Containers[0]
		assert.True(t, container.ReadinessProbe != nil)
		assert.True(t, MountContains(container.VolumeMounts, "readinessprobe"))
		assert.True(t, VolumesContains(spec.Volumes, "readinessprobe"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec)
		}
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, "release-name", []string{"templates/daemonset.yaml"})

		for _, obj := range SplitAndRenderDaemonSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec)
		}
	})
}

func TestDatabaseConfigDoesNotContainEmptyBlocks(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.configFiles": "null",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	assert.NotContains(t, output, "---\n---")
}