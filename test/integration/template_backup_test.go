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
