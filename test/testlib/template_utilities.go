package testlib

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"strings"
	"testing"
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

func EnvFromSourceContains(envs []v1.EnvFromSource, value string) bool {
	for _, n := range envs {
		if n.ConfigMapRef.Name == value {
			return true
		}
	}
	return false
}

func MountContains(mounts []v1.VolumeMount, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func VolumesContains(mounts []v1.Volume, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func SplitAndRenderConfigMap(t *testing.T, output string, expectedNrObjects int) []v1.ConfigMap {
	objects := make([]v1.ConfigMap, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "ConfigMap")) {
			var obj v1.ConfigMap
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderDaemonSet(t *testing.T, output string, expectedNrObjects int) []appsv1.DaemonSet {
	objects := make([]appsv1.DaemonSet, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "DaemonSet")) {
			var obj appsv1.DaemonSet
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderDeployment(t *testing.T, output string, expectedNrObjects int) []appsv1.Deployment {
	objects := make([]appsv1.Deployment, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "Deployment")) {
			var obj appsv1.Deployment
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderReplicationController(t *testing.T, output string, expectedNrObjects int) []v1.ReplicationController {
	objects := make([]v1.ReplicationController, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "ReplicationController")) {
			var obj v1.ReplicationController
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderSecret(t *testing.T, output string, expectedNrObjects int) []v1.Secret {
	objects := make([]v1.Secret, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "Secret")) {
			var obj v1.Secret
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderService(t *testing.T, output string, expectedNrObjects int) []v1.Service {
	objects := make([]v1.Service, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "Service")) {
			var obj v1.Service
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderStatefulSet(t *testing.T, output string, expectedNrObjects int) []appsv1.StatefulSet {
	objects := make([]appsv1.StatefulSet, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "StatefulSet")) {
			var obj appsv1.StatefulSet
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderStorageClass(t *testing.T, output string, expectedNrObjects int) []storagev1.StorageClass {
	objects := make([]storagev1.StorageClass, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "StorageClass")) {
			var obj storagev1.StorageClass
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}


func IsStatefulSetHotCopyEnabled(ss *appsv1.StatefulSet) bool {
	return strings.Contains(ss.Name, "hotcopy")
}

func IsDaemonSetHotCopyEnabled(ss *appsv1.DaemonSet) bool {
	return strings.Contains(ss.Name, "hotcopy")
}

func InferVersionFromTemplate(t *testing.T, options *helm.Options) {
	// prefer injected values
	InjectTestValues(t, options)

	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}

	// inject already specified these
	if options.SetValues["nuodb.image.registry"] != "" ||
		options.SetValues["nuodb.image.repository"] != "" ||
		options.SetValues["nuodb.image.tag"] != "" {
		return
	}

	// pick the version that is in the current charts

	output := helm.RenderTemplate(t, options, ADMIN_HELM_CHART_PATH, "admin-tmp", []string{"templates/statefulset.yaml"})

	statefulSet := SplitAndRenderStatefulSet(t, output, 1)[0]

	t.Logf("Using NuoDB image: %s", statefulSet.Spec.Template.Spec.Containers[0].Image)

	parts := strings.Split(statefulSet.Spec.Template.Spec.Containers[0].Image, "/")
	registry := parts[0]
	afterRegistry := strings.Join(parts[1:], "/")
	part2 := strings.Split(afterRegistry, ":")
	repository := part2[0]
	tag := part2[1]

	options.SetValues["nuodb.image.registry"] = registry
	options.SetValues["nuodb.image.repository"] = repository
	options.SetValues["nuodb.image.tag"] = tag
}