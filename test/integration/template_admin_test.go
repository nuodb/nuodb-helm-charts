package integration

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestAdminDefaultLicense(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	found := false

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: ConfigMap") {
			continue
		}

		if strings.Contains(part, "nuodb-admin-configuration") {
			found = true

			var object v1.ConfigMap
			helm.UnmarshalK8SYaml(t, part, &object)

			assert.Equal(t, len(object.Data), 0)
		}

	}

	assert.True(t, !found, "no matching config map was found")
}

func TestAdminLicenseCanBeSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	licenseString := "red-riding-hood"

	options := &helm.Options{
		SetValues: map[string]string{"admin.configFiles.nuodb\\.lic": licenseString},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	found := false

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: ConfigMap") {
			continue
		}

		if strings.Contains(part, "nuodb-cluster0-admin-configuration") {
			found = true

			var object v1.ConfigMap
			helm.UnmarshalK8SYaml(t, part, &object)

			val, ok := object.Data["nuodb.lic"]

			assert.True(t, ok, "license not properly set")
			assert.Equal(t, val, licenseString)
		}

	}

	assert.True(t, found, "no matching config map was found")
}

func TestAdminStatefulSetVPNRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.securityContext.enabledOnContainer": "true",
			"admin.securityContext.capabilities[0]":    "NET_ADMIN",
			"admin.envFrom.configMapRef[0]":            "test-config",
			"admin.options.leaderAssignmentTimeout":    "30000",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		require.NotEmpty(t, obj.Spec.Template.Spec.Containers)

		adminContainer := obj.Spec.Template.Spec.Containers[0]

		assert.True(t, adminContainer.EnvFrom[0].ConfigMapRef.LocalObjectReference.Name == "test-config")
		assert.Contains(t, adminContainer.SecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))

		assert.Equal(t, "nuoadmin", adminContainer.Args[0])
		assert.Equal(t, "--", adminContainer.Args[1])

		assert.Contains(t, adminContainer.Args[2:], "pendingReconnectTimeout=60000")
		assert.Contains(t, adminContainer.Args[2:], "processLivenessCheckSec=30")
		assert.Contains(t, adminContainer.Args[2:], "leaderAssignmentTimeout=30000")
	}
}

func TestAdminStatefulSetComponentLabel(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, "admin", obj.Spec.Selector.MatchLabels["component"])

		assert.Contains(t, obj.ObjectMeta.Labels, "chart")
		assert.Contains(t, obj.ObjectMeta.Labels, "release")

		assert.Equal(t, "admin", obj.Spec.Template.ObjectMeta.Labels["component"])

		assert.Contains(t, obj.Spec.Template.ObjectMeta.Labels, "chart")
		assert.Contains(t, obj.Spec.Template.ObjectMeta.Labels, "release")
	}
}

func TestAdminClusterServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service-clusterip.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-clusterip", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
	}
}

func TestAdminHeadlessServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service-headless.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP, obj.Spec.Type)
		assert.Equal(t, "None", obj.Spec.ClusterIP)
	}
}

func TestAdminServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.provider":                  "amazon",
			"admin.externalAccess.enabled":    "true",
			"admin.externalAccess.internalIP": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
		assert.Contains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
		assert.Contains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-scheme")
	}

	// render external AWS NLB annotations
	options.SetValues["admin.externalAccess.internalIP"] = "false"
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-type"], "external")
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"], "ip")
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-scheme"], "internet-facing")
	}

	// render custom annotations for the external service
	options.SetValues["admin.externalAccess.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-name"] = "nuodb-admin-nlb"
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})
	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Equal(t, obj.Annotations["service.beta.kubernetes.io/aws-load-balancer-name"], "nuodb-admin-nlb")
		assert.NotContains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-scheme")
	}
}

func TestAdminNodePortServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.provider":                  "amazon",
			"admin.externalAccess.enabled":    "true",
			"admin.externalAccess.type":       "NodePort",
			"admin.externalAccess.internalIP": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/service.yaml"})

	for _, obj := range testlib.SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-nodeport", obj.Name)
		assert.Equal(t, v1.ServiceTypeNodePort, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
		assert.NotContains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
	}
}

func TestAdminStatefulSetVolumes(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{"admin.logPersistence.enabled": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		vcts := make(map[string]bool)

		for _, val := range obj.Spec.VolumeClaimTemplates {
			vcts[val.ObjectMeta.Name] = true
		}

		assert.Contains(t, vcts, "raftlog")
		assert.Contains(t, vcts, "log-volume")

	}
}

func TestAdminMultiClusterEnvVars(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.cluster.name":             "cluster-2",
			"cloud.cluster.entrypointName":   "cluster-1",
			"cloud.cluster.domain":           "cluster2.local",
			"cloud.cluster.entrypointDomain": "cluster1.local",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		environmentals := make(map[string]string)

		require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
		for _, val := range obj.Spec.Template.Spec.Containers[0].Env {
			environmentals[val.Name] = val.Value
		}

		assert.True(t, strings.EqualFold(environmentals["NUODB_DOMAIN_ENTRYPOINT"], "RELEASE-NAME-nuodb-cluster-1-admin-0.nuodb.$(NAMESPACE).svc.cluster1.local"))
		assert.True(t, strings.EqualFold(environmentals["NUODB_ALT_ADDRESS"], "$(POD_NAME).nuodb.$(NAMESPACE).svc.cluster2.local"))

	}
}

func TestConfigDoesNotContainEmptyBlocks(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.configFiles": "null",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	assert.NotContains(t, output, "---\n---")
}

func TestBootstrapServersRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.bootstrapServers": "5",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, options.SetValues["admin.bootstrapServers"], obj.Annotations["nuodb.com/bootstrap-servers"])
		assert.Equal(t, options.SetValues["admin.bootstrapServers"], obj.Labels["bootstrapServers"])
	}
}

func TestGlobalLoadBalancerConfigRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.lbConfig.prefilter":      "not(label(region tiebreaker))",
			"admin.lbConfig.default":        "random(first(label(node ${NODE_NAME:-}) any))",
			"admin.lbConfig.policies.zone1": "round_robin(first(label(zone zone1) any))",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, options.SetValues["admin.lbConfig.prefilter"], obj.Annotations["nuodb.com/load-balancer-prefilter"])
		assert.Equal(t, options.SetValues["admin.lbConfig.default"], obj.Annotations["nuodb.com/load-balancer-default"])
		assert.Equal(t, options.SetValues["admin.lbConfig.policies.zone1"], obj.Annotations["nuodb.com/load-balancer-policy.zone1"])
		// The nearest policy is rendered by default
		assert.Contains(t, obj.Annotations, "nuodb.com/load-balancer-policy.nearest")
		assert.NotContains(t, obj.Annotations, "nuodb.com/sync-load-balancer-config")
	}
}

func TestGlobalLoadBalancerConfigFullSyncRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.lbConfig.fullSync": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, "true", obj.Annotations["nuodb.com/sync-load-balancer-config"])
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-prefilter")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-default")
	}
}

func TestGlobalLoadBalancerConfigRendersOnlyOnEntryPointCluster(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"cloud.cluster.name":            "aws0",
			"admin.lbConfig.fullSync":       "true",
			"admin.lbConfig.prefilter":      "not(label(region tiebreaker))",
			"admin.lbConfig.default":        "random(first(label(node ${NODE_NAME:-}) any))",
			"admin.lbConfig.policies.zone1": "round_robin(first(label(zone zone1) any))",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.NotContains(t, obj.Annotations, "nuodb.com/sync-load-balancer-config")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-prefilter")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-default")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-policy.zone1")
	}
}

func TestAdminPodAnnotationsRender(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.podAnnotations.key1":                                 "value1",
			"admin.podAnnotations.key2":                                 "value2",
			"admin.podAnnotations.key3\\.key3":                          "value3",
			"admin.podAnnotations.key4\\.key4/key4":                     "value4",
			"admin.podAnnotations.key5\\.key5/key5":                     "value5/value5",
			"admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject": `"true"`,
			"admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert": "|\n" +
				"{{- with secret \"nuodb.com/TLS\" -}}\n" +
				"  {{ .Data.data.tlsCACert }}\n" +
				"{{- end }}",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, options.SetValues["admin.podAnnotations.key1"], obj.Spec.Template.ObjectMeta.Annotations["key1"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key2"], obj.Spec.Template.ObjectMeta.Annotations["key2"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key3\\.key3"], obj.Spec.Template.ObjectMeta.Annotations["key3.key3"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key4\\.key4/key4"], obj.Spec.Template.ObjectMeta.Annotations["key4.key4/key4"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key5\\.key5/key5"], obj.Spec.Template.ObjectMeta.Annotations["key5.key5/key5"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert"], obj.Spec.Template.ObjectMeta.Annotations["vault.hashicorp.com/agent-inject-template-ca.cert"])
	}
}

func TestAdminStoragePasswordsRender(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.tde.secrets.databaseA":   "tde-secret-db-a",
			"admin.tde.secrets.databaseB":   "tde-secret-db-b",
			"admin.tde.storagePasswordsDir": "/etc/nuodb/encryption",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		container := obj.Spec.Template.Spec.Containers[0]
		assert.Contains(t, container.Args, "tdeMonitor.storagePasswordsDir="+options.SetValues["admin.tde.storagePasswordsDir"])
		for _, database := range []string{"databaseA", "databaseB"} {
			mount, ok := testlib.GetMount(container.VolumeMounts, "tde-volume-"+database)
			assert.True(t, ok, "mount tde-volume-%s not found", database)
			assert.True(t, mount.ReadOnly)
			assert.Equal(t, options.SetValues["admin.tde.storagePasswordsDir"]+"/"+database, mount.MountPath)
			volume, ok := testlib.GetVolume(obj.Spec.Template.Spec.Volumes, "tde-volume-"+database)
			assert.True(t, ok, "volume tde-volume-%s not found", database)
			assert.Equal(t, options.SetValues["admin.tde.secrets."+database], volume.VolumeSource.Secret.SecretName)
		}
	}
}

func TestLoadBalancerLegacyRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.legacy.loadBalancerJob.enabled": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/job.yaml"})

	for _, obj := range testlib.SplitAndRenderJob(t, output, 1) {
		require.EqualValues(t, "job-lb-policy-nearest", obj.Name)
	}
}

func TestAdminEvictedServers(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.evicted.servers": "{nuoadmin-1,nuoadmin-2}",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		container := obj.Spec.Template.Spec.Containers[0]
		assert.Contains(t, container.Args, "--evicted-servers")
		assert.Contains(t, container.Args, "nuoadmin-1,nuoadmin-2")
	}
}

func TestAdminSecurityContext(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	t.Run("testDefault", func(t *testing.T) {
		options := &helm.Options{}
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.Nil(t, securityContext)
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.Nil(t, containerSecurityContext)
		}
	})

	t.Run("testEnabled", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabled": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1000), *securityContext.FSGroup)
		}
	})

	t.Run("testRunAsNonRootGroup", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup": "true",
				"admin.securityContext.runAsUser":         "5555",
				"admin.securityContext.fsGroup":           "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
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
				"admin.securityContext.fsGroupOnly": "true",
				"admin.securityContext.fsGroup":     "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
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
				"admin.securityContext.enabled":           "true",
				"admin.securityContext.runAsNonRootGroup": "true",
				"admin.securityContext.runAsUser":         "5555",
				"admin.securityContext.fsGroup":           "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(5555), *securityContext.RunAsUser)
			assert.Equal(t, int64(0), *securityContext.RunAsGroup)
			assert.Equal(t, int64(1234), *securityContext.FSGroup)
		}
	})

	t.Run("testRunAsNonRootGroupPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup": "true",
				"admin.securityContext.fsGroupOnly":       "true",
				"admin.securityContext.runAsUser":         "5555",
				"admin.securityContext.fsGroup":           "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
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
				"admin.securityContext.enabled":            "true",
				"admin.securityContext.enabledOnContainer": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
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
				"admin.securityContext.enabledOnContainer":       "true",
				"admin.securityContext.privileged":               "true",
				"admin.securityContext.allowPrivilegeEscalation": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.True(t, *containerSecurityContext.Privileged)
			assert.True(t, *containerSecurityContext.AllowPrivilegeEscalation)
		}
	})

	t.Run("testCapabilitiesAdd", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabledOnContainer":  "true",
				"admin.securityContext.capabilities.add[0]": "NET_ADMIN",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Add, v1.Capability("NET_ADMIN"))
			assert.Nil(t, containerSecurityContext.Capabilities.Drop)
		}
	})

	t.Run("testCapabilitiesDrop", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabledOnContainer":   "true",
				"admin.securityContext.capabilities.drop[0]": "CAP_NET_RAW",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Contains(t, containerSecurityContext.Capabilities.Drop, v1.Capability("CAP_NET_RAW"))
			assert.Nil(t, containerSecurityContext.Capabilities.Add)
		}
	})

	t.Run("testContainerRunAsNonRootGroup", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup":  "true",
				"admin.securityContext.enabledOnContainer": "true",
				"admin.securityContext.runAsUser":          "5555",
				"admin.securityContext.fsGroup":            "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerEnabledPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabled":            "true",
				"admin.securityContext.runAsNonRootGroup":  "true",
				"admin.securityContext.enabledOnContainer": "true",
				"admin.securityContext.runAsUser":          "5555",
				"admin.securityContext.fsGroup":            "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			assert.Equal(t, int64(5555), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(0), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testContainerRunAsNonRootGroupPrecedence", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup":  "true",
				"admin.securityContext.fsGroupOnly":        "true",
				"admin.securityContext.enabledOnContainer": "true",
				"admin.securityContext.runAsUser":          "5555",
				"admin.securityContext.fsGroup":            "1234",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			containerSecurityContext := obj.Spec.Template.Spec.Containers[0].SecurityContext
			assert.NotNil(t, containerSecurityContext)
			// runAsUser should be disregarded, since we can only use 1000:1000 or <uid>:0
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsUser)
			assert.Equal(t, int64(1000), *containerSecurityContext.RunAsGroup)
		}
	})

	t.Run("testEnabledRunInitAsRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabled": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Nil(t, securityContext.RunAsNonRoot)
		}
	})

	t.Run("testEnabledRunInitAsNonRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.enabled":          "true",
				"admin.initContainers.runInitDiskAsRoot": "false",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, true, *securityContext.RunAsNonRoot)
		}
	})

	t.Run("testRunAsNonRootGroupRunInitAsRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Nil(t, securityContext.RunAsNonRoot)
		}
	})

	t.Run("testRunAsNonRootRunInitAsNonRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.securityContext.runAsNonRootGroup": "true",
				"admin.initContainers.runInitDiskAsRoot":  "false",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			securityContext := obj.Spec.Template.Spec.SecurityContext
			assert.NotNil(t, securityContext)
			assert.Equal(t, int64(1000), *securityContext.RunAsUser)
			assert.Equal(t, true, *securityContext.RunAsNonRoot)
		}
	})
}

func getContainerNamed(containers []v1.Container, name string) (*v1.Container, error) {
	var containerNames string
	for _, container := range containers {
		if container.Name == name {
			return &container, nil
		}
		if containerNames != "" {
			containerNames += ", "
		}
		containerNames += container.Name
	}
	return nil, errors.New(fmt.Sprintf("No container named %s found in [%s]", name, containerNames))
}

func TestAdminInitContainers(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	t.Run("testDefault", func(t *testing.T) {
		options := &helm.Options{}
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
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
	})

	t.Run("testRunInitDiskAsNonRoot", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"admin.initContainers.runInitDisk":       "true",
				"admin.initContainers.runInitDiskAsRoot": "false",
			},
		}
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			// look for expected init-disk container
			initContainers := obj.Spec.Template.Spec.InitContainers
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
				"admin.initContainers.runInitDisk": "false",
			},
		}
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			initContainers := obj.Spec.Template.Spec.InitContainers
			assert.Equal(t, 0, len(initContainers))
		}
	})
}

func TestAdminServiceAccount(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	t.Run("testCreated", func(t *testing.T) {
		output := helm.RenderTemplate(t, &helm.Options{}, helmChartPath,
			"release-name", []string{"templates/serviceaccount.yaml", "templates/statefulset.yaml"})

		// Verify that nuodb ServiceAccount is created
		for _, obj := range testlib.SplitAndRenderServiceAccount(t, output, 1) {
			assert.Equal(t, "nuodb", obj.Name)
		}

		// Verify that the correct service account name is used by the admin Pods
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			assert.Equal(t, "nuodb", obj.Spec.Template.Spec.ServiceAccountName)
		}
	})

	t.Run("testNotCreated", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"nuodb.addServiceAccount": "false",
			},
		}

		// nuodb ServiceAccount is not created; a template file passed with
		// "--show-only" is checked by Helm after rendering so expect an error;
		// ref: https://github.com/helm/helm/issues/7295
		_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/serviceaccount.yaml"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "could not find template templates/serviceaccount.yaml in chart")

		// it is expected for the nuodb ServiceAccout to be pre-created
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			assert.Equal(t, "nuodb", obj.Spec.Template.Spec.ServiceAccountName)
		}
	})

	t.Run("testDefaultServiceAccount", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"nuodb.addServiceAccount": "false",
				"nuodb.serviceAccount":    "",
			},
		}

		// nuodb ServiceAccount is not created
		_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/serviceaccount.yaml"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "could not find template templates/serviceaccount.yaml in chart")

		// the default ServiceAccount for the namespace will be used
		output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}

		options = &helm.Options{
			SetValues: map[string]string{
				"nuodb.addServiceAccount": "false",
				"nuodb.serviceAccount":    "null",
			},
		}

		// nuodb ServiceAccount is not created
		_, err = helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/serviceaccount.yaml"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "could not find template templates/serviceaccount.yaml in chart")

		// the default ServiceAccount for the namespace will be used
		output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
		for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
			assert.Empty(t, obj.Spec.Template.Spec.ServiceAccountName)
		}

		options = &helm.Options{
			SetValues: map[string]string{
				"nuodb.addServiceAccount": "true",
				"nuodb.serviceAccount":    "",
			},
		}

		// nuodb ServiceAccount is not created
		_, err = helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/serviceaccount.yaml"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "could not find template templates/serviceaccount.yaml in chart")
	})

}
