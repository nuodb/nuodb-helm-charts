package integration

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func TestDatabaseBackupCronJobRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	testlib.SplitAndRenderCronJob(t, output, 2)
}

func TestDatabaseBackupCronJobDisabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.enableBackups": "false",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	assert.Contains(t, err.Error(), "could not find template")
}

func TestDatabaseBackupCronJobGarbage(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.enableBackups": "garbage",
		},
	}

	// Garbage value fails helm rendering
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid boolean value: garbage")
}

func TestDatabaseJournalBackupCronJobEnabled(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "true",
		},
	}

	// Verify that journal CronJob is rendered
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	testlib.SplitAndRenderCronJob(t, output, 3)

	// Verify that journal-hot-copy engine option is enabled
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
	journalFlagFound := false
	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		for _, arg := range obj.Spec.Template.Spec.Containers[0].Args {
			if strings.Contains(arg, "journal-hot-copy enable") {
				journalFlagFound = true
			}
		}
	}
	assert.True(t, journalFlagFound, "journal-hot-copy should be enabled")
}

func TestDatabaseJournalBackupCronJobGarbage(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "garbage",
		},
	}

	// Verify that helm rendering fails
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid boolean value: garbage")
}

func TestDatabaseJournalBackupCronJobDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Verify that journal CronJob is rendered
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	jobs := testlib.SplitAndRenderCronJob(t, output, 2)
	assert.Equal(t, 2, len(jobs))

	// Verify that journal-hot-copy engine option is enabled
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
	journalFlagFound := false
	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		for _, arg := range obj.Spec.Template.Spec.Containers[0].Args {
			if strings.Contains(arg, "journal-hot-copy enable") {
				journalFlagFound = true
			}
		}
	}
	assert.False(t, journalFlagFound, "journal-hot-copy should be missing")
}

func TestDatabaseJournalBackupCronJobFalse(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "false",
		},
	}

	// Verify that journal CronJob is rendered
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	jobs := testlib.SplitAndRenderCronJob(t, output, 2)
	assert.Equal(t, 2, len(jobs))

	// Verify that journal-hot-copy engine option is enabled
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
	journalFlagFound := false
	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		for _, arg := range obj.Spec.Template.Spec.Containers[0].Args {
			if strings.Contains(arg, "journal-hot-copy enable") {
				journalFlagFound = true
			}
		}
	}
	assert.False(t, journalFlagFound, "journal-hot-copy should be missing")
}

func TestDatabaseJournalBackupCronJobFalseFile(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		ValuesFiles: []string{"../files/database-journal-disabled.yaml"},
	}

	// Verify that journal CronJob is rendered
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	jobs := testlib.SplitAndRenderCronJob(t, output, 2)
	assert.Equal(t, 2, len(jobs))

	// Verify that journal-hot-copy engine option is enabled
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/statefulset.yaml"})
	journalFlagFound := false
	for _, obj := range testlib.SplitAndRenderStatefulSet(t, output, 1) {
		for _, arg := range obj.Spec.Template.Spec.Containers[0].Args {
			if strings.Contains(arg, "journal-hot-copy enable") {
				journalFlagFound = true
			}
		}
	}
	assert.False(t, journalFlagFound, "journal-hot-copy should be missing")
}

func TestDatabaseBackupCronJobRestartPolicyDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	for _, job := range testlib.SplitAndRenderCronJob(t, output, 2) {
		assert.Equal(t, job.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy, v1.RestartPolicyOnFailure)
	}
}

func TestDatabaseBackupCronJobRestartPolicyOverride(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.restartPolicy": "Never",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	for _, job := range testlib.SplitAndRenderCronJob(t, output, 2) {
		assert.Equal(t, job.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy, v1.RestartPolicyNever)
	}
}

func TestDatabaseBackupTimeout(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "true",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	// verify that default value for timeout is rendered
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 2) {
		require.NotEmpty(t, obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
		backupContainer := obj.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
		assert.Subset(t, backupContainer.Args, []string{
			"--timeout",
			"0",
		})
	}

	options = &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "true",
			"database.sm.hotCopy.timeout":               "1800",
			"database.sm.hotCopy.journalBackup.timeout": "950",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output = helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	// verify that configured value for timeout is rendered
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 2) {
		require.NotEmpty(t, obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
		backupContainer := obj.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
		if strings.Contains(obj.Name, "journal") {
			assert.Subset(t, backupContainer.Args, []string{
				"--timeout",
				"950",
			})
		} else {
			assert.Subset(t, backupContainer.Args, []string{
				"--timeout",
				"1800",
			})
		}
	}
}
