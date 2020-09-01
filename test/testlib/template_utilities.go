package testlib

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/test/integration"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
)

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

	statefulSet := integration.SplitAndRenderStatefulSet(t, output, 1)[0]

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