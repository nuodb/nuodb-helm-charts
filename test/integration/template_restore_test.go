package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func verifyRestoreResourceLabels(t *testing.T, releaseName string, options *helm.Options, obj metav1.Object) {
	opt := testlib.GetExtractedOptions(options)
	labels := obj.GetLabels()
	app := fmt.Sprintf("%s-%s-%s-%s-restore", releaseName, opt.DomainName, opt.ClusterName, opt.DbName)
	msg, ok := testlib.MapContains(labels, map[string]string{
		"app":      app,
		"group":    "nuodb",
		"subgroup": "restore",
		"domain":   opt.DomainName,
		"database": opt.DbName,
		"chart":    "restore",
		"release":  releaseName,
	})
	require.Truef(t, ok, "Mandatory labels missing from resource %s: %s", obj.GetName(), msg)

	resourceLabels := make(map[string]string)
	for k, v := range options.SetValues {
		if strings.HasPrefix(k, "restore.resourceLabels.") {
			labelKey := strings.TrimPrefix(k, "restore.resourceLabels.")
			resourceLabels[labelKey] = v
		}
	}
	if len(resourceLabels) > 0 {
		msg, ok := testlib.MapContains(labels, resourceLabels)
		require.Truef(t, ok, "User supplied labels missing from resource %s: %s", obj.GetName(), msg)
	}
}

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

func TestRestoreWithProcessFilter(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.RESTORE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			// Select all HCSMs in cluster0
			"restore.processFilter": "and(label(backup cluster0) labels(role hotcopy))",
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
			"--process-filter",
			options.SetValues["restore.processFilter"],
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

func TestRestoreResourceLabels(t *testing.T) {
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
		verifyRestoreResourceLabels(t, "release-name", options, &obj)
		verifyRestoreResourceLabels(t, "release-name", options, &obj.Spec.Template)
	}
}
