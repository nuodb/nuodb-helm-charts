package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/configmap.yaml"})

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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/statefulset.yaml"})


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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/statefulset.yaml"})

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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/service-clusterip.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-clusterip", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP,  obj.Spec.Type)
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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/service-headless.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb", obj.Name)
		assert.Equal(t, v1.ServiceTypeClusterIP,  obj.Spec.Type)
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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/service.yaml"})

	for _, obj := range SplitAndRenderService(t, output, 1) {
		assert.Equal(t, "nuodb-balancer", obj.Name)
		assert.Equal(t, v1.ServiceTypeLoadBalancer,  obj.Spec.Type)
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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/statefulset.yaml"})

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
			"cloud.cluster.name": "cluster-2",
			"cloud.cluster.entrypointName": "cluster-1",
			"cloud.cluster.domain": "cluster2.local",
			"cloud.cluster.entrypointDomain": "cluster1.local",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/statefulset.yaml"})

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
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name",  []string{"templates/configmap.yaml"})


	assert.NotContains(t, output, "---\n---")
}