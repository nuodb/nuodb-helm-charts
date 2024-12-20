package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestAutoRestoreGarbage(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.autoRestore.type": "garbage",
		},
	}

	// Garbage value fails helm rendering
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})
	require.Error(t, err)
}

func TestAutoRestoreSource(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.autoRestore.source": ":garbage",
		},
	}

	// There are reserved restore sources starting with ":" which are ":latest"
	// and ":group-latest"; special sources that are not in allowed fails helm
	// rendering; we can't really validate other sources as URLs and backup set
	// names are also allowed
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})
	require.Error(t, err)
}

func TestAutoRestoreDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	var found = false

	for _, obj := range testlib.SplitAndRenderConfigMap(t, output, 1) {
		if obj.Name == "release-name-nuodb-cluster0-demo-database-restore" {
			found = true
			assert.EqualValues(t, "stream", obj.Data["NUODB_AUTO_RESTORE_TYPE"])
		}
	}

	require.True(t, found)
}

func TestAutoRestoreValidValueStream(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.autoRestore.type": "stream",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	var found = false

	for _, obj := range testlib.SplitAndRenderConfigMap(t, output, 1) {
		if obj.Name == "release-name-nuodb-cluster0-demo-database-restore" {
			found = true
			assert.Empty(t, obj.Data["NUODB_AUTO_RESTORE"])
			assert.EqualValues(t, "stream", obj.Data["NUODB_AUTO_RESTORE_TYPE"])
		}
	}

	require.True(t, found)
}

func TestAutoRestoreValidValueBackupset(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.autoRestore.type": "backupset",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	var found = false

	for _, obj := range testlib.SplitAndRenderConfigMap(t, output, 1) {
		if obj.Name == "release-name-nuodb-cluster0-demo-database-restore" {
			found = true
			assert.Empty(t, obj.Data["NUODB_AUTO_RESTORE"])
			assert.EqualValues(t, "backupset", obj.Data["NUODB_AUTO_RESTORE_TYPE"])
		}
	}

	require.True(t, found)
}

func TestRestoreConfigMapSecondaryRelease(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.primaryRelease":         "false",
			"database.autoImport.source":      "http://files.nuodb.com/backup.tar.gz",
			"database.autoImport.type":        "backupset",
			"database.autoImport.stripLevels": "2",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/configmap.yaml"})

	var found = false

	for _, obj := range testlib.SplitAndRenderConfigMap(t, output, 1) {
		if obj.Name == "release-name-nuodb-cluster0-demo-database-restore" {
			found = true
			assert.EqualValues(t, options.SetValues["database.autoImport.stripLevels"], obj.Data["NUODB_IMPORT_STRIP_LEVELS"])
			assert.EqualValues(t, options.SetValues["database.autoImport.source"], obj.Data["NUODB_AUTO_IMPORT"])
			assert.EqualValues(t, options.SetValues["database.autoImport.type"], obj.Data["NUODB_AUTO_IMPORT_TYPE"])
		}
	}

	require.True(t, found)
}
