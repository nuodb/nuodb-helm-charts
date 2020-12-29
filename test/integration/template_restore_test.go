package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestRestoreDefaults(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/job.yaml"})

	for _, obj := range testlib.SplitAndRenderJob(t, output, 1) {
		require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
		restoreContainer := obj.Spec.Template.Spec.Containers[0]
		assert.ElementsMatch(t, restoreContainer.Args, []string{
			"nuorestore",
			"--type",
			"database",
			"--db-name",
			"demo",
			"--source",
			":latest",
			"--auto",
			"true",
		})
	}
}

func TestRestoreNoDatabaseRestart(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.autoRestart": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/job.yaml"})

	for _, obj := range testlib.SplitAndRenderJob(t, output, 1) {
		require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
		restoreContainer := obj.Spec.Template.Spec.Containers[0]
		assert.ElementsMatch(t, restoreContainer.Args, []string{
			"nuorestore",
			"--type",
			"database",
			"--db-name",
			"demo",
			"--source",
			":latest",
			"--auto",
			"false",
		})
	}
}
