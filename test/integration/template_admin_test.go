package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

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
			"admin.securityContext.capabilities":    "[ NET_ADMIN ]",
			"admin.envFrom.configMapRef[0]":         "test-config",
			"admin.options.leaderAssignmentTimeout": "30000",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderService(t, output, 1) {
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

	for _, obj := range SplitAndRenderService(t, output, 1) {
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

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer, obj.Spec.Type)
		assert.Empty(t, obj.Spec.ClusterIP)
		assert.Contains(t, obj.Annotations, "service.beta.kubernetes.io/aws-load-balancer-internal")
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
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

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
		assert.NotContains(t, obj.Annotations, "nuodb.com/sync-load-balancer-config")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-prefilter")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-default")
		assert.NotContains(t, obj.Annotations, "nuodb.com/load-balancer-policy.zone1")
	}
}

func TestAnnotationsRender(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.ADMIN_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.podAnnotations.key1": "value1",
			"admin.podAnnotations.key2": "value2",
			"admin.podAnnotations.key3\\.key3": "value3",
			"admin.podAnnotations.key4\\.key4/key4": "value4",
			"admin.podAnnotations.key5\\.key5/key5": "value5/value5",
			"admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject": `"true"`,
			"admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert": "|\n"+
			"{{- with secret \"nuodb.com/TLS\" -}}\n"+
			"  {{ .Data.data.tlsCACert }}\n"+
			"{{- end }}",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})

	for _, obj := range SplitAndRenderStatefulSet(t, output, 1) {
		assert.Equal(t, options.SetValues["admin.podAnnotations.key1"], obj.Annotations["key1"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key2"], obj.Annotations["key2"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key3\\.key3"], obj.Annotations["key3.key3"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key4\\.key4/key4"], obj.Annotations["key4.key4/key4"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.key5\\.key5/key5"], obj.Annotations["key5.key5/key5"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject"], obj.Annotations["vault.hashicorp.com/agent-inject"])
		assert.Equal(t, options.SetValues["admin.podAnnotations.vault\\.hashicorp\\.com/agent-inject-template-ca\\.cert"], obj.Annotations["vault.hashicorp.com/agent-inject-template-ca.cert"])
	}
}
