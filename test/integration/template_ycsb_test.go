package integration

import (
	"strings"
	"testing"

	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestYcsbConfigMapRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../incubator/demo-ycsb"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

	configMapCount := 0

	parts := strings.Split(output, "---")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if !strings.Contains(part, "kind: ConfigMap") {
			continue
		}

		var object v1.ConfigMap
		helm.UnmarshalK8SYaml(t, part, &object)

		configMapCount += 1

	}

	assert.Equal(t, 3, configMapCount)
}

func TestYcsbRCRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../incubator/demo-ycsb"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/replicationcontroller.yaml"})

	var object v1.ReplicationController
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, object.Name == "ycsb-load")
	assert.Check(t, *object.Spec.Replicas == 0)
}

func TestYcsbRCReplicas(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../incubator/demo-ycsb"

	options := &helm.Options{
		SetValues: map[string]string{
			"ycsb.replicas": "1",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/replicationcontroller.yaml"})

	var object v1.ReplicationController
	helm.UnmarshalK8SYaml(t, output, &object)

	assert.Check(t, object.Name == "ycsb-load")
	assert.Check(t, *object.Spec.Replicas == 1)
}