package integration

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
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

	assert.ErrorContains(t, err, "could not find template")
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