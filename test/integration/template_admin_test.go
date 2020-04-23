package integration

import (
	"strings"
	"testing"

	"gotest.tools/assert"
	appsv1 "k8s.io/api/apps/v1"
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

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

	assert.Assert(t, !found, "no matching config map was found")
}

func TestAdminLicenseCanBeSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"
	licenseString := "red-riding-hood"

	options := &helm.Options{
		SetValues: map[string]string{"admin.configFiles.nuodb\\.lic": licenseString},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

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

			assert.Assert(t, ok, "license not properly set")
			assert.Equal(t, val, licenseString)
		}

	}

	assert.Assert(t, found, "no matching config map was found")
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: StatefulSet") {
			continue
		}

		var object appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &object)

		adminContainer := object.Spec.Template.Spec.Containers[0]
		assert.Check(t, adminContainer.SecurityContext.Capabilities.Add[0] == "NET_ADMIN")
		assert.Check(t, adminContainer.EnvFrom[0].ConfigMapRef.LocalObjectReference.Name == "test-config")
		assert.Check(t, adminContainer.Args[0] == "nuoadmin")
		assert.Check(t, adminContainer.Args[1] == "--")

		// make sure all expected admin option overrides appear in command-line
		adminOptions := make(map[string]bool)
		for _, option := range adminContainer.Args[2:] {
			adminOptions[option] = true
		}
		assert.Check(t, adminOptions["pendingReconnectTimeout=60000"])
		assert.Check(t, adminOptions["processLivenessCheckSec=30"])
		assert.Check(t, adminOptions["leaderAssignmentTimeout=30000"])
	}
}

func TestAdminStatefulSetComponentLabel(t *testing.T) {
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

		skind, ok := ss.Spec.Selector.MatchLabels["component"]
		assert.Check(t, ok)
		assert.Check(t, skind == "admin")

		_, ok = ss.ObjectMeta.Labels["chart"]
		assert.Check(t, ok)
		_, ok = ss.ObjectMeta.Labels["release"]
		assert.Check(t, ok)

		okind, ok := ss.Spec.Template.ObjectMeta.Labels["component"]
		assert.Check(t, ok)
		assert.Check(t, okind == "admin")

		_, ok = ss.Spec.Template.ObjectMeta.Labels["chart"]
		assert.Check(t, ok)
		_, ok = ss.Spec.Template.ObjectMeta.Labels["release"]
		assert.Check(t, ok)
	}
}

func TestAdminClusterServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service-clusterip.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "name: nuodb-clusterip"))
	assert.Check(t, strings.Contains(output, "type: ClusterIP"))
	assert.Check(t, !strings.Contains(output, "clusterIP: None"))
}

func TestAdminHeadlessServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service-headless.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "name: nuodb"))
	assert.Check(t, strings.Contains(output, "type: ClusterIP"))
	assert.Check(t, strings.Contains(output, "clusterIP: None"))
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "name: nuodb-balancer"))
	assert.Check(t, strings.Contains(output, "type: LoadBalancer"))
	assert.Check(t, strings.Contains(output, "aws-load-balancer-internal"))
}

func TestAdminStatefulSetVolumes(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/admin"

	options := &helm.Options{
		SetValues: map[string]string{"admin.logPersistence.enabled": "true"},
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

		assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[0].ObjectMeta.Name, "raftlog"))
		assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[1].ObjectMeta.Name, "log-volume"))
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

		environmentals := make(map[string]string)

		for _, val := range ss.Spec.Template.Spec.Containers[0].Env {
			environmentals[val.Name] = val.Value
		}

		assert.Check(t, strings.EqualFold(environmentals["NUODB_DOMAIN_ENTRYPOINT"], "RELEASE-NAME-nuodb-cluster-1-admin-0.nuodb.$(NAMESPACE).svc.cluster1.local"))
		assert.Check(t, strings.EqualFold(environmentals["NUODB_ALT_ADDRESS"], "$(POD_NAME).nuodb.$(NAMESPACE).svc.cluster2.local"))

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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

	assert.Assert(t, !strings.Contains(output, "---\n---"))
}