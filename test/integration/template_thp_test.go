package integration

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"gotest.tools/assert"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
	"testing"
)

func TestThpOpenShift(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/transparent-hugepage"

	options := &helm.Options{
		SetValues: map[string]string{
			"openshift.enabled": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/daemonset.yaml"})

	cnt := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: DaemonSet") {
			continue
		}

		cnt += 1

		var object appsv1.DaemonSet
		helm.UnmarshalK8SYaml(t, part, &object)

		assert.Assert(t, object.Spec.Template.Spec.InitContainers[0].SecurityContext.Privileged != nil)
		assert.Assert(t, object.Spec.Template.Spec.Containers[0].SecurityContext.Privileged != nil)

		assert.Check(t, *object.Spec.Template.Spec.InitContainers[0].SecurityContext.Privileged == true)
		assert.Check(t, *object.Spec.Template.Spec.Containers[0].SecurityContext.Privileged == true)
	}

	assert.Check(t, cnt == 1)
}