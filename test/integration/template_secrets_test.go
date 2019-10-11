package integration

import (
	"testing"

	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestSecretsDatabaseDefault(t *testing.T) {
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
