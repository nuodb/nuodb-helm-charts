package integration

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