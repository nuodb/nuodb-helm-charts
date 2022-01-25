package testlib

import (
	"fmt"
	"strings"
	"testing"

	"k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
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

func GetMount(mounts []v1.VolumeMount, expectedName string) (*v1.VolumeMount, bool) {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return &mount, true
		}
	}
	return nil, false
}

func VolumesContains(mounts []v1.Volume, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func GetVolume(volumes []v1.Volume, expectedName string) (*v1.Volume, bool) {
	for _, volume := range volumes {
		if volume.Name == expectedName {
			return &volume, true
		}
	}
	return nil, false
}

func GetVolumeClaim(vcp []v1.PersistentVolumeClaim, expectedName string) (*v1.PersistentVolumeClaim, bool) {
	for _, volume := range vcp {
		if volume.Name == expectedName {
			return &volume, true
		}
	}
	return nil, false
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

func SplitAndRenderCronJob(t *testing.T, output string, expectedNrObjects int) []v1beta1.CronJob {
	objects := make([]v1beta1.CronJob, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "CronJob")) {
			var obj v1beta1.CronJob
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

func SplitAndRenderJob(t *testing.T, output string, expectedNrObjects int) []batchv1.Job {
	objects := make([]batchv1.Job, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "Job")) {
			var obj batchv1.Job
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

func SplitAndRenderRole(t *testing.T, output string, expectedNrObjects int) []rbacv1.Role {
	objects := make([]rbacv1.Role, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "Role")) {
			var obj rbacv1.Role
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)

	return objects
}

func SplitAndRenderServiceAccount(t *testing.T, output string, expectedNrObjects int) []v1.ServiceAccount {
	objects := make([]v1.ServiceAccount, 0)

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", "ServiceAccount")) {
			var obj v1.ServiceAccount
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

func AssertResourceValue(t *testing.T, options *helm.Options, key string, actual *resource.Quantity) {
	if expected, ok := options.SetValues[key]; ok {
		require.Equal(t, 0, actual.Cmp(resource.MustParse(expected)),
			fmt.Sprintf("Resource mismatch key='%s', expected='%s', actual='%s'", key, expected, actual.String()))
	}
}
