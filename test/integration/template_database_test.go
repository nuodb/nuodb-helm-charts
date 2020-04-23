package integration

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
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
		if n.Name == key && n.Value == value {
			return true
		}
	}
	return false
}

func EnvContainsValueFrom(envs []v1.EnvVar, key string, valueFrom *v1.EnvVarSource) bool {
	for _, n := range envs {
		if n.Name == key && cmp.Equal(n.ValueFrom, valueFrom) {
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

	// check for the minimum 3 secret values: database-name, database-password, database-username
	assert.Check(t, len(object.StringData) >= 3)

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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

	configs := make(map[string]bool)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, "kind: ConfigMap") {

			var cm v1.ConfigMap
			helm.UnmarshalK8SYaml(t, part, &cm)

			for k := range cm.Data {
				configs[k] = true
			}
		}
	}

	assert.Check(t, configs["nuosm"])
	assert.Check(t, configs["nuote"])
	assert.Check(t, configs["readinessprobe"])
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

func TestDatabaseClusterServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service-clusterip.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "name: demo-clusterip"))
	assert.Check(t, strings.Contains(output, "type: ClusterIP"))
	assert.Check(t, !strings.Contains(output, "clusterIP: None"))
	assert.Check(t, strings.Contains(output, "component: te"))
}

func TestDatabaseHeadlessServiceRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service-headless.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "name: demo"))
	assert.Check(t, strings.Contains(output, "type: ClusterIP"))
	assert.Check(t, strings.Contains(output, "clusterIP: None"))
	assert.Check(t, strings.Contains(output, "component: te"))
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/service.yaml"})

	var object v1.Service
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, strings.Contains(output, "type: LoadBalancer"))
	assert.Check(t, strings.Contains(output, "name: demo-balancer"))
	assert.Check(t, strings.Contains(output, "kind: Service"))
	assert.Check(t, strings.Contains(output, "aws-load-balancer-internal"))
	assert.Check(t, strings.Contains(output, "component: te"))
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

func TestDatabaseStatefulSetVolumes(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{"database.sm.logPersistence.enabled": "true"},
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

			if strings.Contains(part, "-hotcopy") {
				assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[0].ObjectMeta.Name, "archive-volume"))
				assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[1].ObjectMeta.Name, "backup-volume"))
				assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[2].ObjectMeta.Name, "log-volume"))
			} else {
				assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[0].ObjectMeta.Name, "archive-volume"))
				assert.Check(t, strings.Contains(ss.Spec.VolumeClaimTemplates[1].ObjectMeta.Name, "log-volume"))
			}
		}
	}

	assert.Check(t, cnt == 2)
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
			"admin.tlsKeyStore.secret":          "nuodb-keystore",
			"admin.tlsKeyStore.key":             "nuoadmin.p12",
			"admin.tlsKeyStore.password":        "changeIt",
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
		assert.Check(t, EnvContainsValueFrom(args, "NUODB_ALT_ADDRESS", &expectedAltAddress))
		assert.Check(t, EnvContains(args, "CUSTOM_ENV_VAR", "CUSTOM_ENV_VAR_VALUE"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

		assert.Check(t, strings.Contains(output, "kind: Deployment"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

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
			"database.securityContext.capabilities[0]": "NET_ADMIN",
			"database.envFrom.configMapRef[0]":         "test-config",
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
			"database.securityContext.capabilities[0]": "NET_ADMIN",
			"database.envFrom.configMapRef[0]":         "test-config",
			"database.enableDaemonSet":                 "true",
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

				if testlib.IsStatefulSetHotCopyEnabled(&obj) {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup cluster0"))
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

				if testlib.IsDaemonSetHotCopyEnabled(&obj) {
					assert.Check(t, ArgContains(obj.Spec.Template.Spec.Containers[0].Args, "backup cluster0"))
				}
			}
		}

		assert.Check(t, cnt == 2)
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
		assert.Check(t, container.ReadinessProbe != nil)
		assert.Check(t, mountContains(container.VolumeMounts, "readinessprobe"))
		assert.Check(t, volumesContain(spec.Volumes, "readinessprobe"))
	}

	t.Run("testDeployment", func(t *testing.T) {
		// Run RenderTemplate to render the template and capture the output.
		output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

		assert.Check(t, strings.Contains(output, "kind: Deployment"))

		var obj appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &obj)

		basicChecks(obj.Spec.Template.Spec)
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

				basicChecks(obj.Spec.Template.Spec)
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

				basicChecks(obj.Spec.Template.Spec)
			}
		}

		assert.Check(t, cnt == 2)
	})
}

func mountContains(mounts []v1.VolumeMount, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func volumesContain(mounts []v1.Volume, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
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
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

	assert.Assert(t, !strings.Contains(output, "---\n---"))
}