package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestDatabaseSecretsDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/secret.yaml"})

	for _, obj := range testlib.SplitAndRenderSecret(t, output, 1) {
		assert.Contains(t, obj.StringData, "database-name")
		assert.Contains(t, obj.StringData, "database-password")
		assert.Contains(t, obj.StringData, "database-username")
	}

}

func TestDatabaseConfigMaps(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.tde.storagePasswordsDir": "/etc/nuodb/encryption",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	configs := make(map[string]string)

	for _, obj := range testlib.SplitAndRenderConfigMap(t, output, 3) {
		for k, v := range obj.Data {
			configs[k] = v
		}
	}

	assert.Contains(t, configs, "nuosm")
	assert.Contains(t, configs, "nuote")
	assert.Contains(t, configs, "readinessprobe")
	assert.Contains(t, configs, "NUODB_STORAGE_PASSWORDS_DIR")
	assert.Equal(t, configs["NUODB_STORAGE_PASSWORDS_DIR"], "/etc/nuodb/encryption")
}

func TestDatabaseClusterServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service-clusterip.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 2) {
		assert.Equal(t, v1.ServiceTypeClusterIP, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)

		if obj.Name == "demo-clusterip" {
			// This is the cluster IP targeting all database TEs
			assert.Equal(t, "te", obj.Spec.Selector["component"])
			assert.Equal(t, "nuodb", obj.Spec.Selector["domain"])
			assert.Equal(t, "demo", obj.Spec.Selector["database"])
		} else {
			// This is the cluster IP targeting only TEs in this TE group
			assert.Equal(t, "te", obj.Spec.Selector["component"])
			assert.Equal(t, "release-name-nuodb-cluster0-demo-database", obj.Spec.Selector["app"])
		}
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

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "demo", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP, obj.Spec.Type)
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.Equal(t, "None", obj.Spec.ClusterIP)
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

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database", obj.Spec.Selector["app"])
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.Contains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
	}

	// render external AWS NLB annotations
	options.SetValues["database.te.externalAccess.internalIP"] = "false"
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-type"], "external")
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"], "ip")
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-scheme"], "internet-facing")
	}

	// render custom annotations for the external service
	options.SetValues["database.te.externalAccess.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-name"] = "nuodb-demo-nlb"
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})
	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-name"], "nuodb-demo-nlb")
		assert.NotContains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-scheme")
	}

}

func TestDatabaseNodePortServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.provider":                        "amazon",
			"database.te.externalAccess.enabled":    "true",
			"database.te.externalAccess.type":       "NodePort",
			"database.te.externalAccess.internalIP": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database-nodeport", obj.Name)
		assert.Equal(t, v1.ServiceTypeNodePort, obj.Spec.Type)
		assert.Equal(t, "release-name-nuodb-cluster0-demo-database", obj.Spec.Selector["app"])
		assert.Equal(t, "te", obj.Spec.Selector["component"])
		assert.NotContains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
	}
}

func TestDatabaseStatefulSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
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

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
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

	for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
		assert.Equal(t, "te", obj.Spec.Selector.MatchLabels["component"])
		assert.Equal(t, "te", obj.Spec.Template.ObjectMeta.Labels["component"])
	}

}

func TestDatabaseOtherOptions(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.otherOptions.keystore":               "/etc/nuodb/keys/nuoadmin.p12",
			"database.sm.otherOptions.keystore":               "/etc/nuodb/keys/nuoadmin.p12",
			"database.te.otherOptions.enable-external-access": "false",
			"database.sm.otherOptions.enable-external-access": "false",
			"database.te.otherOptions.resolve-hostname":       "true",
			"database.sm.otherOptions.resolve-hostname":       "true",
			"database.te.otherOptions.some-other-flag":        "",
			"database.sm.otherOptions.some-other-flag":        "",
			"admin.tlsKeyStore.secret":                        "nuodb-keystore",
			"admin.tlsKeyStore.key":                           "nuoadmin.p12",
			"admin.tlsKeyStore.password":                      "changeIt",
		},
	}

	basicArgChecks := func(args []string) {
		t.Log(args)
		assert.True(t, testlib.ArgContains(args, "--keystore"))
		assert.True(t, testlib.ArgContains(args, "/etc/nuodb/keys/nuoadmin.p12"))
		assert.True(t, testlib.ArgContains(args, "--resolve-hostname"))
		assert.False(t, testlib.ArgContains(args, "--enable-external-access"))
		assert.False(t, testlib.ArgContains(args, "--some-other-flag"))
		assert.False(t, testlib.ArgContains(args, "true"))
		assert.NotContains(t, args, "")
	}

	basicEnvChecks := func(args []v1.EnvVar) {
		assert.True(t, testlib.EnvContains(args, "NUODOCKER_KEYSTORE_PASSWORD", "changeIt"))
	}

	basicInitContainerCommandChecks := func(commands []string) {
		assert.NotContains(t, commands, "/var/opt/nuodb/journal")
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)

			require.NotEmpty(t, obj.Spec.Template.Spec.InitContainers)
			basicInitContainerCommandChecks(obj.Spec.Template.Spec.InitContainers[0].Command)
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
		assert.True(t, testlib.EnvContainsValueFrom(args, "NUODB_ALT_ADDRESS", &expectedAltAddress))
		assert.True(t, testlib.EnvContains(args, "CUSTOM_ENV_VAR", "CUSTOM_ENV_VAR_VALUE"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
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
			"database.securityContext.enabledOnContainer": "true",
			"database.securityContext.capabilities[0]":    "NET_ADMIN",
			"database.envFrom.configMapRef[0]":            "test-config",
		},
	}

	basicChecks := func(args []v1.Container) {
		assert.Contains(t, args[0].SecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))
		assert.True(t, testlib.EnvFromSourceContains(args[0].EnvFrom, "test-config"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
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
		assert.True(t, testlib.ArgContains(args, "cloud minikube"))
		assert.True(t, testlib.ArgContains(args, "region local"))
		assert.True(t, testlib.ArgContains(args, "zone local-b"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers[0].Args)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec.Containers[0].Args)

			if testlib.IsStatefulSetHotCopyEnabled(&obj) {
				assert.True(t, testlib.ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup cluster0"))
				assert.True(t, testlib.ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "role hotcopy"))
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
		assert.True(t, testlib.MountContains(container.VolumeMounts, "readinessprobe"))
		assert.True(t, testlib.VolumesContains(spec.Volumes, "readinessprobe"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
			basicChecks(obj.Spec.Template.Spec)
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
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

func TestLoadBalancerConfigurationRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.lbConfig.prefilter": "not(label(zone DR))",
			"database.lbConfig.default":   "random(first(label(node ${NODE_NAME:-}) any))",
		},
	}

	assertLoadBalancerAnnotations := func(annotations map[string]string) {
		assert.Equal(t, options.SetValues["database.lbConfig.prefilter"], annotations["nuodb.com/load-balancer-prefilter"])
		assert.Equal(t, options.SetValues["database.lbConfig.default"], annotations["nuodb.com/load-balancer-default"])
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assertLoadBalancerAnnotations(obj.Annotations)
		}
	})
}

func TestDefaultLoadBalancerConfigurationRendersOnlyOnEntryPointCluster(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.cluster.name":          "aws0",
			"database.lbConfig.prefilter": "not(label(zone DR))",
			"database.lbConfig.default":   "random(first(label(node ${NODE_NAME:-}) any))",
		},
	}

	assertLoadBalancerAnnotations := func(annotations map[string]string) {
		assert.NotContains(t, annotations, "nuodb.com/load-balancer-prefilter")
		assert.NotContains(t, annotations, "nuodb.com/load-balancer-default")
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assertLoadBalancerAnnotations(obj.Annotations)
		}
	})
}

func TestDefaultLoadBalancerConfigurationNotRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	assertLoadBalancerAnnotations := func(annotations map[string]string) {
		assert.NotContains(t, annotations, "nuodb.com/load-balancer-prefilter")
		assert.NotContains(t, annotations, "nuodb.com/load-balancer-default")
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assertLoadBalancerAnnotations(obj.Annotations)
		}
	})
}

func TestAutomaticDatabaseProtocolUpgrade(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	t.Run("testAutomaticUpgradeOnEntrypointCluster", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.automaticProtocolUpgrade.enabled": "true",
			},
		}
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Equal(t, "true", obj.Annotations["nuodb.com/automatic-database-protocol-upgrade"])
		}
	})

	t.Run("testAutomaticUpgradeOnOtherCluster", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"cloud.cluster.name":                        "aws0",
				"database.automaticProtocolUpgrade.enabled": "true",
			},
		}
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.NotContains(t, obj.Annotations, "nuodb.com/automatic-database-protocol-upgrade")
		}
	})

	t.Run("testAutomaticUpgradeWithPreference", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.automaticProtocolUpgrade.enabled":           "true",
				"database.automaticProtocolUpgrade.tePreferenceQuery": "random(label(region tiebreaker))",
			},
		}
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Equal(t, "true", obj.Annotations["nuodb.com/automatic-database-protocol-upgrade"])
			assert.Equal(t,
				options.SetValues["database.automaticProtocolUpgrade.tePreferenceQuery"],
				obj.Annotations["nuodb.com/automatic-database-protocol-upgrade.te-preference-query"])
		}
	})
}

func TestDatabasePodAnnotationsRender(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"cc.podAnnotations.key1":                                       "value1",
			"database.podAnnotations.key2":                                 "value2",
			"database.podAnnotations.key3\\.key3":                          "value3",
			"database.podAnnotations.key4\\.key4/key4":                     "value4",
			"database.podAnnotations.key5\\.key5/key5":                     "value5/value5",
			"database.podAnnotations.vault\\.hashicorp\\.com/agent-inject": `"true"`,
			"database.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert": "|\n" +
				"{{- with secret \"nuodb.com/TLS\" -}}\n" +
				"  {{ .Data.data.tlsCACert }}\n" +
				"{{- end }}",
		},
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})

		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Equal(t, options.SetValues["database.podAnnotations.key1"], obj.Spec.Template.ObjectMeta.Annotations["key1"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key2"], obj.Spec.Template.ObjectMeta.Annotations["key2"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key3\\.key3"], obj.Spec.Template.ObjectMeta.Annotations["key3.key3"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key4\\.key4/key4"], obj.Spec.Template.ObjectMeta.Annotations["key4.key4/key4"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key5\\.key5/key5"], obj.Spec.Template.ObjectMeta.Annotations["key5.key5/key5"])
			assert.Equal(t, options.SetValues["database.podAnnotations.vault\\.hashicorp\\.com/agent-inject"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject"])
			assert.Equal(t, options.SetValues["database.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject-template-ca.cert"])
		}
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			assert.Equal(t, options.SetValues["database.podAnnotations.key1"], obj.Spec.Template.ObjectMeta.Annotations["key1"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key2"], obj.Spec.Template.ObjectMeta.Annotations["key2"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key3\\.key3"], obj.Spec.Template.ObjectMeta.Annotations["key3.key3"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key4\\.key4/key4"], obj.Spec.Template.ObjectMeta.Annotations["key4.key4/key4"])
			assert.Equal(t, options.SetValues["database.podAnnotations.key5\\.key5/key5"], obj.Spec.Template.ObjectMeta.Annotations["key5.key5/key5"])
			assert.Equal(t, options.SetValues["database.podAnnotations.vault\\.hashicorp\\.com/agent-inject"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject"])
			assert.Equal(t, options.SetValues["database.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject-template-ca.cert"])
		}
	})
}

func TestDatabaseStoragePasswordsRender(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.tde.secrets.demo":        "tde-secret",
			"admin.tde.storagePasswordsDir": "/etc/nuodb/encryption",
		},
	}

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			database := "demo"
			container := obj.Spec.Template.Spec.Containers[0]
			mount, ok := testlib.GetMount(container.VolumeMounts, "tde-volume-"+database)
			assert.True(t, ok, "mount tde-volume-%s not found", database)
			assert.True(t, mount.ReadOnly)
			assert.Equal(t, options.SetValues["admin.tde.storagePasswordsDir"]+"/"+database, mount.MountPath)
			volume, ok := testlib.GetVolume(obj.Spec.Template.Spec.Volumes, "tde-volume-"+database)
			assert.True(t, ok, "volume tde-volume-%s not found", database)
			assert.Equal(t, options.SetValues["admin.tde.secrets."+database], volume.VolumeSource.Secret.SecretName)
		}
	})
}

func TestDatabaseSeparateJournal(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	t.Run("testStatefulSetDefaults", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
			},
		}

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			container := obj.Spec.Template.Spec.Containers[0]

			assert.True(t, testlib.EnvContains(container.Env, "SEPARATE_JOURNAL", "true"))

			mount, ok := testlib.GetMount(container.VolumeMounts, "journal-volume")
			assert.True(t, ok, "mount journal-volume not found")
			assert.EqualValues(t, "/var/opt/nuodb/journal", mount.MountPath)

			claim, ok := testlib.GetVolumeClaim(obj.Spec.VolumeClaimTemplates, "journal-volume")
			assert.True(t, ok, "volume journal-volume not found")
			assert.Equal(t, v1.ReadWriteOnce, claim.Spec.AccessModes[0])
			assert.Nil(t, claim.Spec.StorageClassName)
		}
	})

	t.Run("testStatefulSetOverrides", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":                      "true",
				"database.sm.noHotCopy.journalPath.enabled":                    "true",
				"database.sm.hotCopy.journalPath.persistence.accessModes[0]":   "ReadWriteMany",
				"database.sm.noHotCopy.journalPath.persistence.accessModes[0]": "ReadWriteMany",
				"database.sm.hotCopy.journalPath.persistence.storageClass":     "non-default",
				"database.sm.noHotCopy.journalPath.persistence.storageClass":   "non-default",
			},
		}

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			container := obj.Spec.Template.Spec.Containers[0]

			assert.True(t, testlib.EnvContains(container.Env, "SEPARATE_JOURNAL", "true"))

			mount, ok := testlib.GetMount(container.VolumeMounts, "journal-volume")
			assert.True(t, ok, "mount journal-volume not found")
			assert.EqualValues(t, "/var/opt/nuodb/journal", mount.MountPath)

			claim, ok := testlib.GetVolumeClaim(obj.Spec.VolumeClaimTemplates, "journal-volume")
			assert.True(t, ok, "volume journal-volume not found")
			assert.Equal(t, v1.ReadWriteMany, claim.Spec.AccessModes[0])
			assert.EqualValues(t, "non-default", *claim.Spec.StorageClassName)
		}
	})

	t.Run("testStatefulDefaultFalse", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled": "",
			},
		}

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			container := obj.Spec.Template.Spec.Containers[0]

			assert.True(t, testlib.EnvContains(container.Env, "SEPARATE_JOURNAL", "false"))

			_, ok := testlib.GetMount(container.VolumeMounts, "journal-volume")
			assert.False(t, ok, "mount journal-volume not found")

			_, ok = testlib.GetVolumeClaim(obj.Spec.VolumeClaimTemplates, "journal-volume")
			assert.False(t, ok, "volume journal-volume not found")
		}
	})

}

func TestDatabaseSecurityContext(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	t.Run("testDefault", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.Nil(t, securityContext)
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.Nil(t, containerSecurityContext)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.Nil(t, securityContext)
			containerSecurityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			assert.Nil(t, containerSecurityContext)
		}
	})

	t.Run("testEnabled", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
				"database.securityContext.enabled":          "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1000), *securityContext.FSGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1000), *securityContext.FSGroup)
		}
	})

	t.Run("testRunAsNonRootGroup", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":    "true",
				"database.sm.noHotCopy.journalPath.enabled":  "true",
				"database.securityContext.runAsNonRootGroup": "true",
				"database.securityContext.runAsUser":         "5555",
				"database.securityContext.fsGroup":           "1234",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(1000), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(1000), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}
	})

	t.Run("testFsGroupOnly", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
				"database.securityContext.fsGroupOnly":      "true",
				"database.securityContext.fsGroup":          "1234",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// user and group should be absent
			assert.Nil(t, securityContext.RunAsUser)
			assert.Nil(t, securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// user and group should be absent
			assert.Nil(t, securityContext.RunAsUser)
			assert.Nil(t, securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}
	})

	t.Run("testEnabledPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":    "true",
				"database.sm.noHotCopy.journalPath.enabled":  "true",
				"database.securityContext.enabled":           "true",
				"database.securityContext.runAsNonRootGroup": "true",
				"database.securityContext.runAsUser":         "5555",
				"database.securityContext.fsGroup":           "1234",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(5555), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(5555), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}
	})

	t.Run("testRunAsNonRootGroupPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":    "true",
				"database.sm.noHotCopy.journalPath.enabled":  "true",
				"database.securityContext.runAsNonRootGroup": "true",
				"database.securityContext.fsGroupOnly":       "true",
				"database.securityContext.runAsUser":         "5555",
				"database.securityContext.fsGroup":           "1234",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(1000), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			securityContext := dep.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(1000), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}
	})

	t.Run("testContainerEnabled", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.securityContext.enabled":            "true",
				"database.securityContext.enabledOnContainer": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.False(t, *containerSecurityContext.Privileged)
			assert.False(t, *containerSecurityContext.AllowPrivilegeEscalation)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(0), *containerSecurityContext.RunAsGroup)
		}
		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.False(t, *containerSecurityContext.Privileged)
			assert.False(t, *containerSecurityContext.AllowPrivilegeEscalation)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(0), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerPrivileged", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.securityContext.enabledOnContainer":       "true",
				"database.securityContext.privileged":               "true",
				"database.securityContext.allowPrivilegeEscalation": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.True(t, *containerSecurityContext.Privileged)
			assert.True(t, *containerSecurityContext.AllowPrivilegeEscalation)
		}
		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.True(t, *containerSecurityContext.Privileged)
			assert.True(t, *containerSecurityContext.AllowPrivilegeEscalation)
		}
	})

	t.Run("testCapabilitiesAdd", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.securityContext.enabledOnContainer":  "true",
				"database.securityContext.capabilities.add[0]": "NET_ADMIN",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))
			assert.Nil(t, containerSecurityContext.Capabilities.Drop)
		}
		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))
			assert.Nil(t, containerSecurityContext.Capabilities.Drop)
		}
	})

	t.Run("testCapabilitiesDrop", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.securityContext.enabledOnContainer":   "true",
				"database.securityContext.capabilities.drop[0]": "CAP_NET_RAW",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Drop, v1.Capability("CAP_NET_RAW"))
			assert.Nil(t, containerSecurityContext.Capabilities.Add)
		}
		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Drop, v1.Capability("CAP_NET_RAW"))
			assert.Nil(t, containerSecurityContext.Capabilities.Add)
		}
	})

	t.Run("testContainerRunAsNonRootGroup", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":     "true",
				"database.sm.noHotCopy.journalPath.enabled":   "true",
				"database.securityContext.runAsNonRootGroup":  "true",
				"database.securityContext.runAsUser":          "5555",
				"database.securityContext.enabledOnContainer": "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerFsGroupOnly", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":     "true",
				"database.sm.noHotCopy.journalPath.enabled":   "true",
				"database.securityContext.fsGroupOnly":        "true",
				"database.securityContext.enabledOnContainer": "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// user and group should be absent
			assert.Nil(t, containerSecurityContext.RunAsUser)
			assert.Nil(t, containerSecurityContext.RunAsGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// user and group should be absent
			assert.Nil(t, containerSecurityContext.RunAsUser)
			assert.Nil(t, containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerEnabledPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":     "true",
				"database.sm.noHotCopy.journalPath.enabled":   "true",
				"database.securityContext.enabled":            "true",
				"database.securityContext.runAsNonRootGroup":  "true",
				"database.securityContext.runAsUser":          "5555",
				"database.securityContext.enabledOnContainer": "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Equal(t, int64(5555), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(0), *containerSecurityContext.RunAsGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Equal(t, int64(5555), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(0), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerRunAsNonRootGroupPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":     "true",
				"database.sm.noHotCopy.journalPath.enabled":   "true",
				"database.securityContext.runAsNonRootGroup":  "true",
				"database.securityContext.fsGroupOnly":        "true",
				"database.securityContext.runAsUser":          "5555",
				"database.securityContext.enabledOnContainer": "true",
			},
		}

		// check security context on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}

		// check security context on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			containerSecurityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}
	})
}

func TestDatabaseInitContainers(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	t.Run("testDefault", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
			},
		}

		// check init containers on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			// look for expected init-disk container
			initContainers := obj.Spec.Template.Spec.InitContainers
			assert.Equal(t, 1, len(initContainers))
			container, err := getContainerNamed(initContainers, "init-disk")
			assert.NoError(t, err)
			// check that security context for container specifies root user and group
			securityContext := container.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(0), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
		}

		// check init containers on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			// look for expected init-disk container
			initContainers := dep.Spec.Template.Spec.InitContainers
			assert.Equal(t, 1, len(initContainers))
			container, err := getContainerNamed(initContainers, "init-disk")
			assert.NoError(t, err)
			// check that security context for container specifies root user and group
			securityContext := container.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(0), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
		}
	})

	t.Run("testRunInitDiskAsNonRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
				"database.initContainers.runInitDisk":       "true",
				"database.initContainers.runInitDiskAsRoot": "false",
			},
		}

		// check init containers on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			// look for expected init-disk container
			initContainers := obj.Spec.Template.Spec.InitContainers
			assert.Equal(t, 1, len(initContainers))
			container, err := getContainerNamed(initContainers, "init-disk")
			assert.NoError(t, err)
			// check that security context is not defined
			securityContext := container.SecurityContext
			assert.Nil(t, securityContext)
		}

		// check init containers on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			// look for expected init-disk container
			initContainers := dep.Spec.Template.Spec.InitContainers
			assert.Equal(t, 1, len(initContainers))
			container, err := getContainerNamed(initContainers, "init-disk")
			assert.NoError(t, err)
			// check that security context is not defined
			securityContext := container.SecurityContext
			assert.Nil(t, securityContext)
		}
	})

	t.Run("testDisabled", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"database.sm.hotCopy.journalPath.enabled":   "true",
				"database.sm.noHotCopy.journalPath.enabled": "true",
				"database.initContainers.runInitDisk":       "false",
			},
		}

		// check init containers on SM StatefulSets
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			initContainers := obj.Spec.Template.Spec.InitContainers
			assert.Equal(t, 0, len(initContainers))
		}

		// check init containers on TE Deployment
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/deployment.yaml"})
		for _, dep := range testlib.SplitAndRenderDeployment(t, output, 1) {
			initContainers := dep.Spec.Template.Spec.InitContainers
			assert.Equal(t, 0, len(initContainers))
		}
	})
}

func TestDatabaseServiceAccount(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	t.Run("testUsage", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"nuodb.serviceAccount": "foo",
			},
		}

		// correct ServiceAccount is used
		output := helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/statefulset.yaml", "templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			assert.Equal(t, "foo", obj.Spec.Template.Spec.ServiceAccountName)
		}
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Equal(t, "foo", obj.Spec.Template.Spec.ServiceAccountName)
		}
	})

	t.Run("testDefaultServiceAccount", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"nuodb.serviceAccount": "",
			},
		}

		// the default ServiceAccount for the namespace will be used
		output := helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/statefulset.yaml", "templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}

		options = &helm.Options{
			SetValues: map[string]string{
				"nuodb.serviceAccount": "null",
			},
		}

		// the default ServiceAccount for the namespace will be used
		output = helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/statefulset.yaml", "templates/deployment.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 2) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}
		for _, obj := range testlib.SplitAndRenderDeployment(t, output, 1) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}
	})

}
