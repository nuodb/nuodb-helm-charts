package testlib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/helm"
	coreosv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ArgContains(args []string, x string) bool {
	for _, n := range args {
		if strings.Contains(n, x) {
			return true
		}
	}
	return false
}

func EnvGet(envs []corev1.EnvVar, key string) (string, bool) {
	for _, n := range envs {
		if n.Name == key {
			return n.Value, true
		}
	}
	return "", false
}

func EnvGetValueFrom(envs []corev1.EnvVar, key string) (*corev1.EnvVarSource, bool) {
	for _, n := range envs {
		if n.Name == key {
			return n.ValueFrom, true
		}
	}
	return nil, false
}

func AssertEnvNotContains(t *testing.T, envs []corev1.EnvVar, key string) {
	actual, ok := EnvGet(envs, key)
	assert.False(t, ok, "Unexpected environment variable %s=%s", key, actual)
}

func AssertEnvContains(t *testing.T, envs []corev1.EnvVar, key, expected string) {
	actual, ok := EnvGet(envs, key)
	assert.True(t, ok, "Environment variable %s not set", key)
	assert.Equal(t, expected, actual)
}

func EnvContains(envs []corev1.EnvVar, key string, expected string) bool {
	if actual, ok := EnvGet(envs, key); ok {
		return actual == expected
	}
	return false
}

func EnvContainsValueFrom(envs []corev1.EnvVar, key string, valueFrom *corev1.EnvVarSource) bool {
	for _, n := range envs {
		if n.Name == key && cmp.Equal(n.ValueFrom, valueFrom) {
			return true
		}
	}
	return false
}

func AssertEnvContainsValueFrom(t *testing.T, envs []corev1.EnvVar, key string, expected corev1.EnvVarSource) {
	actual, ok := EnvGetValueFrom(envs, key)
	assert.True(t, ok, "Environment variable %s not set", key)
	assert.NotNil(t, actual, "Environment variable %s valueFrom not set")
	if actual != nil {
		assert.Equal(t, expected, *actual)
	}
}

func EnvFromSourceContains(envs []corev1.EnvFromSource, value string) bool {
	for _, n := range envs {
		if n.ConfigMapRef.Name == value {
			return true
		}
	}
	return false
}

func MountContains(mounts []corev1.VolumeMount, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func GetMount(mounts []corev1.VolumeMount, expectedName string) (*corev1.VolumeMount, bool) {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return &mount, true
		}
	}
	return nil, false
}

func VolumesContains(mounts []corev1.Volume, expectedName string) bool {
	for _, mount := range mounts {
		if mount.Name == expectedName {
			return true
		}
	}
	return false
}

func MapContains(actual map[string]string, expected map[string]string) (string, bool) {
	if actual == nil || expected == nil {
		return fmt.Sprintf("Map not initialized: actual=%#v, expected=%#v", actual, expected), false
	}
	for k, v := range expected {
		if val, ok := actual[k]; ok {
			if !assert.ObjectsAreEqual(v, val) {
				return fmt.Sprintf("%#v does not contain %#v: values for key '%s' does not match", actual, expected, k), false
			}
		} else {
			return fmt.Sprintf("%#v does not contain %#v: key '%s' missing", actual, expected, k), false
		}
	}
	return "", true
}

func GetVolume(volumes []corev1.Volume, expectedName string) (*corev1.Volume, bool) {
	for _, volume := range volumes {
		if volume.Name == expectedName {
			return &volume, true
		}
	}
	return nil, false
}

func GetVolumeClaim(vcp []corev1.PersistentVolumeClaim, expectedName string) (*corev1.PersistentVolumeClaim, bool) {
	for _, volume := range vcp {
		if volume.Name == expectedName {
			return &volume, true
		}
	}
	return nil, false
}

func SplitAndRender[T any](t *testing.T, output string, expectedNrObjects int, kind string) []T {
	objects := make([]T, 0)
	parts := strings.Split(output, "---\n")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, fmt.Sprintf("kind: %s", kind)) {
			var obj T
			helm.UnmarshalK8SYaml(t, part, &obj)

			objects = append(objects, obj)
		}
	}

	require.GreaterOrEqual(t, len(objects), expectedNrObjects)
	return objects
}

func SplitAndRenderPersistentVolumeClaim(t *testing.T, output string, expectedNrObjects int) []corev1.PersistentVolumeClaim {
	return SplitAndRender[corev1.PersistentVolumeClaim](t, output, expectedNrObjects, "PersistentVolumeClaim")
}

func SplitAndRenderConfigMap(t *testing.T, output string, expectedNrObjects int) []corev1.ConfigMap {
	return SplitAndRender[corev1.ConfigMap](t, output, expectedNrObjects, "ConfigMap")
}

func SplitAndRenderCronJob(t *testing.T, output string, expectedNrObjects int) []batchv1beta1.CronJob {
	return SplitAndRender[batchv1beta1.CronJob](t, output, expectedNrObjects, "CronJob")
}

func SplitAndRenderDaemonSet(t *testing.T, output string, expectedNrObjects int) []appsv1.DaemonSet {
	return SplitAndRender[appsv1.DaemonSet](t, output, expectedNrObjects, "DaemonSet")
}

func SplitAndRenderJob(t *testing.T, output string, expectedNrObjects int) []batchv1.Job {
	return SplitAndRender[batchv1.Job](t, output, expectedNrObjects, "Job")
}

func SplitAndRenderDeployment(t *testing.T, output string, expectedNrObjects int) []appsv1.Deployment {
	return SplitAndRender[appsv1.Deployment](t, output, expectedNrObjects, "Deployment")
}

func SplitAndRenderReplicationController(t *testing.T, output string, expectedNrObjects int) []corev1.ReplicationController {
	return SplitAndRender[corev1.ReplicationController](t, output, expectedNrObjects, "ReplicationController")
}

func SplitAndRenderSecret(t *testing.T, output string, expectedNrObjects int) []corev1.Secret {
	return SplitAndRender[corev1.Secret](t, output, expectedNrObjects, "Secret")
}

func SplitAndRenderService(t *testing.T, output string, expectedNrObjects int) []corev1.Service {
	return SplitAndRender[corev1.Service](t, output, expectedNrObjects, "Service")
}

func SplitAndRenderStatefulSet(t *testing.T, output string, expectedNrObjects int) []appsv1.StatefulSet {
	return SplitAndRender[appsv1.StatefulSet](t, output, expectedNrObjects, "StatefulSet")
}

func SplitAndRenderStorageClass(t *testing.T, output string, expectedNrObjects int) []storagev1.StorageClass {
	return SplitAndRender[storagev1.StorageClass](t, output, expectedNrObjects, "StorageClass")
}

func SplitAndRenderRole(t *testing.T, output string, expectedNrObjects int) []rbacv1.Role {
	return SplitAndRender[rbacv1.Role](t, output, expectedNrObjects, "Role")
}

func SplitAndRenderClusterRole(t *testing.T, output string, expectedNrObjects int) []rbacv1.ClusterRole {
	return SplitAndRender[rbacv1.ClusterRole](t, output, expectedNrObjects, "ClusterRole")
}

func SplitAndRenderClusterClusterRoleBinding(t *testing.T, output string, expectedNrObjects int) []rbacv1.ClusterRoleBinding {
	return SplitAndRender[rbacv1.ClusterRoleBinding](t, output, expectedNrObjects, "ClusterRoleBinding")
}

func SplitAndRenderServiceAccount(t *testing.T, output string, expectedNrObjects int) []corev1.ServiceAccount {
	return SplitAndRender[corev1.ServiceAccount](t, output, expectedNrObjects, "ServiceAccount")
}

func SplitAndRenderIngress(t *testing.T, output string, expectedNrObjects int) []networkingv1.Ingress {
	return SplitAndRender[networkingv1.Ingress](t, output, expectedNrObjects, "Ingress")
}

func SplitAndRenderPodMonitor(t *testing.T, output string, expectedNrObjects int) []coreosv1.PodMonitor {
	return SplitAndRender[coreosv1.PodMonitor](t, output, expectedNrObjects, "PodMonitor")
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
