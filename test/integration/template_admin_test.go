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

		if strings.Contains(part, "nuodb-admin-configuration") {
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
			"admin.envFrom[0].configMapRef.name":    "test-config",
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
		assert.Check(t, adminContainer.Args[2] == "pendingReconnectTimeout=60000")
		assert.Check(t, adminContainer.Args[3] == "leaderAssignmentTimeout=30000")
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
			"cloud.provider": "amazon",
			"admin.externalAccess.enabled": "true",
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