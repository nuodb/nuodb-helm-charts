package integration

import (
	"strings"
	"testing"

	"gotest.tools/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func ArgContains(args []string, x string) bool {
	for _, n := range args {
		if strings.Contains(n, x) {
			return true
		}
	}
	return false
}

func EnvContains(envs []v1.EnvVar, key string, value string) bool {
	for _, n := range envs {
		if (n.Name == key && n.Value == value) {
			return true
		}
	}
	return false
}

func TestDatabaseSecretsDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/secret.yaml"})

	var object v1.Secret
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, len(object.StringData) == 5)

	_, ok := object.StringData["database-name"]
	assert.Check(t, ok)

	_, ok = object.StringData["database-password"]
	assert.Check(t, ok)

	_, ok = object.StringData["database-username"]
	assert.Check(t, ok)
}

func TestDatabaseConfigMaps(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})
}

func TestDatabaseDaemonSetDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		assert.Check(t, !strings.Contains(part, "kind: DaemonSet"))
	}
}

func TestDatabaseDaemonSetEnabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	var cnt int

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, "kind: DaemonSet") {
			cnt++
		}
	}

	assert.Check(t, cnt == 2)
}

func TestDatabaseDeploymentConfigDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deploymentconfig.yaml"})

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		assert.Check(t, !strings.Contains(part, "kind: DeploymentConfig"))
	}
}

func TestDatabaseDeploymentConfigRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"openshift.enabled": "true", "openshift.enableDeploymentConfigs": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deploymentconfig.yaml"})

	assert.Check(t, strings.Contains(output, "kind: DeploymentConfig"))
}

func TestDatabaseServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	value, exists := object.Spec.Selector["component"]

	assert.Assert(t, exists)
	assert.Assert(t, value == "te")
}

func TestDatabaseServiceDaemonSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service.yaml"})

	assert.Check(t, strings.Contains(output, "kind: Service"))
}

func TestDatabaseStatefulSetDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.enableDaemonSet": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	assert.Check(t, !strings.Contains(output, "kind: StatefulSet"))
}

func TestDatabaseStatefulSet(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	var cnt int

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, "kind: StatefulSet") {
			cnt++

			var ss appsv1.StatefulSet
			helm.UnmarshalK8SYaml(t, part, &ss)

			skind, ok := ss.Spec.Selector.MatchLabels["component"]
			assert.Check(t, ok)
			assert.Check(t, skind == "sm")

			okind, ok := ss.Spec.Template.ObjectMeta.Labels["component"]
			assert.Check(t, ok)
			assert.Check(t, okind == "sm")
		}
	}

	assert.Check(t, cnt == 2)
}

func TestDatabaseDeploymentDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"openshift.enabled": "true", "openshift.enableDeploymentConfigs": "true"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

	assert.Check(t, !strings.Contains(output, "kind: Deployment"))
}

func TestDatabaseDeploymentRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

	assert.Check(t, strings.Contains(output, "kind: Deployment"))

	var dep appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &dep)

	skind, ok := dep.Spec.Selector.MatchLabels["component"]
	assert.Check(t, ok)
	assert.Check(t, skind == "te")

	okind, ok := dep.Spec.Template.ObjectMeta.Labels["component"]
	assert.Check(t, ok)
	assert.Check(t, okind == "te")
}

func TestDatabaseOtherOptions(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.otherOptions.keystore": "/etc/nuodb/keys/nuoadmin.p12",
			"database.sm.otherOptions.keystore": "/etc/nuodb/keys/nuoadmin.p12",
			"admin.tlsKeyStore.secret":     "nuodb-keystore",
			"admin.tlsKeyStore.key":        "nuoadmin.p12",
			"admin.tlsKeyStore.password":   "changeIt",
		},
	}

	basicArgChecks := func(args []string) {
		assert.Check(t, ArgContains(args, "--keystore"))
		assert.Check(t, ArgContains(args, "/etc/nuodb/keys/nuoadmin.p12"))
	}

	basicEnvChecks := func(args []v1.EnvVar) {
		assert.Check(t, EnvContains(args, "NUODOCKER_KEYSTORE_PASSWORD", "changeIt"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

		assert.Check(t, strings.Contains(output, "kind: Deployment"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

		basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
		basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
	})

	t.Run("testDeploymentConfig", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["openshift.enabled"] = "true"
		localOptions.SetValues["openshift.enableDeploymentConfigs"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, []string{"templates/deploymentconfig.yaml"})

		assert.Check(t, strings.Contains(output, "kind: DeploymentConfig"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

		basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
		basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

		var cnt int

		parts := strings.Split(output, "---")
		for _, part := range parts {
			if len(part) == 0 {
				continue
			}

			if strings.Contains(part, "kind: StatefulSet") {
				cnt++

				var obj appsv1.StatefulSet
				helm.UnmarshalK8SYaml(t, part, &obj)

				basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
				basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
			}
		}

		assert.Check(t, cnt == 2)
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, []string{"templates/daemonset.yaml"})

		var cnt int

		parts := strings.Split(output, "---")
		for _, part := range parts {
			if len(part) == 0 {
				continue
			}

			if strings.Contains(part, "kind: DaemonSet") {
				cnt++

				var obj appsv1.DaemonSet
				helm.UnmarshalK8SYaml(t, part, &obj)

				basicArgChecks(obj.Spec.Template.Spec.Containers[0].Args)
				basicEnvChecks(obj.Spec.Template.Spec.Containers[0].Env)
			}
		}

		assert.Check(t, cnt == 2)
	})
}

func TestDatabaseStandardVPNRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.securityContext.capabilities": "[ NET_ADMIN ]",
			"database.envFrom":                      "[ configMapRef: { name: test-config } ]",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	objs := make([]interface{}, 10)
	helm.UnmarshalK8SYaml(t, output, &objs)

	for _, k8obj := range objs {

		switch k8obj := k8obj.(type) {
		case appsv1.StatefulSet:
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities.Add[0], "NET_ADMIN")
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].EnvFrom[0].ConfigMapRef.Name, "test-config")

		case appsv1.Deployment:
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities.Add[0], "NET_ADMIN")
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].EnvFrom[0].ConfigMapRef.Name, "test-config")
		}
	}
}

func TestDatabaseDaemonSetVPNRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.securityContext.capabilities": "[ NET_ADMIN ]",
			"database.envFrom":                      "[ configMapRef: { name: test-config } ]",
			"database.enableDaemonSet":              "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	objs := make([]interface{}, 10)
	helm.UnmarshalK8SYaml(t, output, &objs)

	for _, k8obj := range objs {

		switch k8obj := k8obj.(type) {
		case appsv1.DaemonSet:
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities.Add[0], "NET_ADMIN")
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].EnvFrom[0].ConfigMapRef.Name, "test-config")

		case appsv1.Deployment:
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities.Add[0], "NET_ADMIN")
			assert.Check(t, k8obj.Spec.Template.Spec.Containers[0].EnvFrom[0].ConfigMapRef.Name, "test-config")
		}
	}
}

func TestDatabaseDeploymentConfigVPNRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.securityContext.capabilities": "[ NET_ADMIN ]",
			"database.envFrom":                      "[ configMapRef: { name: test-config } ]",
			"openshift.enabled":                     "true",
			"openshift.enableDeploymentConfigs":     "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

	// NTJ - Sadly, doesn't work. YAML.v3 still creates map[interface{}]interface{}  :(
	// m := make(map[string]interface{})

	parts := strings.Split(output, "---")

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DeploymentConfig") {
			continue
		}

		var object appsv1.StatefulSet
		helm.UnmarshalK8SYaml(t, part, &object)

		adminContainer := object.Spec.Template.Spec.Containers[0]
		assert.Check(t, adminContainer.SecurityContext.Capabilities.Add[0], "NET_ADMIN")
		assert.Check(t, adminContainer.EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "test-config")
	}
}

func TestDatabaseLabeling(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.te.labels.cloud": "minikube",
			"database.te.labels.region": "local",
			"database.te.labels.zone": "local-b",
			"database.sm.labels.cloud": "minikube",
			"database.sm.labels.region": "local",
			"database.sm.labels.zone": "local-b",
			},
	}

	basicChecks := func(args []string) {
		assert.Check(t, ArgContains(args, "cloud minikube"))
		assert.Check(t, ArgContains(args, "region local"))
		assert.Check(t, ArgContains(args, "zone local-b"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

		assert.Check(t, strings.Contains(output, "kind: Deployment"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

		basicChecks(obj.Spec.Template.Spec.Containers[0].Args)
	})

	t.Run("testDeploymentConfig", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["openshift.enabled"] = "true"
		localOptions.SetValues["openshift.enableDeploymentConfigs"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, []string{"templates/deploymentconfig.yaml"})

		assert.Check(t, strings.Contains(output, "kind: DeploymentConfig"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

		basicChecks(obj.Spec.Template.Spec.Containers[0].Args)
	})

	t.Run("testStatefulSet", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/statefulset.yaml"})

		var cnt int

		parts := strings.Split(output, "---")
		for _, part := range parts {
			if len(part) == 0 {
				continue
			}

			if strings.Contains(part, "kind: StatefulSet") {
				cnt++

				var obj appsv1.StatefulSet
				helm.UnmarshalK8SYaml(t, part, &obj)

				basicChecks(obj.Spec.Template.Spec.Containers[0].Args)

				if isStatefulSetHotCopyEnabled(&obj) {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup enabled"))
				} else {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup disabled"))
				}
			}
		}

		assert.Check(t, cnt == 2)
	})

	t.Run("testDaemonSet", func(t *testing.T) {
		// make a copy
		localOptions := *options
		localOptions.SetValues["database.enableDaemonSet"] = "true"

		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, &localOptions, helmChartPath, []string{"templates/daemonset.yaml"})

		var cnt int

		parts := strings.Split(output, "---")
		for _, part := range parts {
			if len(part) == 0 {
				continue
			}

			if strings.Contains(part, "kind: DaemonSet") {
				cnt++

				var obj appsv1.DaemonSet
				helm.UnmarshalK8SYaml(t, part, &obj)

				basicChecks(obj.Spec.Template.Spec.Containers[0].Args)

				if isDaemonSetHotCopyEnabled(&obj) {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup enabled"))
				} else {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup disabled"))
				}
			}
		}

		assert.Check(t, cnt == 2)
	})
}