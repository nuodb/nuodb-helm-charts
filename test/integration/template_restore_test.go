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
			"--manual",
			"false",
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
			"--manual",
			"false",
		})
	}
}

func TestRestoreSpecificArchives(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.archiveIds[0]": "1",
			"restore.archiveIds[1]": "2",
			"restore.archiveIds[2]": "3",
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
			"true",
			"--manual",
			"false",
			"--archive-ids",
			"1 2 3",
		})
	}
}

func TestRestoreSpecificLabels(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.labels.labelA": "labelA_value",
			"restore.labels.labelB": "labelB_value",
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
			"true",
			"--manual",
			"false",
			"--labels",
			"labelA labelA_value labelB labelB_value",
		})
	}
}

func TestRestoreSpecificLabelsAsString(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.labels": "labelA labelA_value labelB labelB_value",
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
			"true",
			"--manual",
			"false",
			"--labels",
			"labelA labelA_value labelB labelB_value",
		})
	}
}

func TestManualRestore(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.manual": "true",
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
			"true",
			"--manual",
			"true",
		})
	}
}

func TestRestoreRequestStripLevels(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.stripLevels": "2",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/job.yaml"})

	for _, obj := range testlib.SplitAndRenderJob(t, output, 1) {
		require.NotEmpty(t, obj.Spec.Template.Spec.Containers)
		restoreContainer := obj.Spec.Template.Spec.Containers[0]
		assert.True(t, testlib.EnvContains(restoreContainer.Env, "NUODB_RESTORE_REQUEST_STRIP_LEVELS", "2"))
	}
}

func TestRestoreRequestSource(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"restore.source": ":garbage",
		},
	}

	// There are reserved restore sources starting with ":" which are ":latest"
	// and ":group-latest"; special sources that are not in allowed fails helm
	// rendering; we can't really validate other sources as URLs and backup set
	// names are also allowed
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/job.yaml"})
	require.Error(t, err)
}
