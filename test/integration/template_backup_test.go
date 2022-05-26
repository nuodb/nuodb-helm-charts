package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func verifyBackupResourceLabels(t *testing.T, releaseName string, options *helm.Options, obj metav1.Object) {
	msg, ok := testlib.MapContains(obj.GetLabels(), map[string]string{
		"subgroup":     "backup",
		"backup-group": "cluster0-0",
	})
	require.Truef(t, ok, "Backup resource labels do not match for resource %s: %s", obj.GetName(), msg)
	verifyDatabaseResourceLabels(t, "release-name", options, obj)
}

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

func TestDatabaseBackupGroupsDefault(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.replicas":                      "2",
			"database.sm.hotCopy.journalBackup.enabled":         "true",
			"database.sm.hotCopy.fullSchedule":                  "35 22 * * 6",
			"database.sm.hotCopy.incrementalSchedule":           "35 22 * * 0-5",
			"database.sm.hotCopy.journalBackup.journalSchedule": "?/15 * * * *",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "database", []string{"templates/cronjob.yaml"})

	// one set of CronJobs per backup group is rendered
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 6) {
		require.NotEmpty(t, obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
		backupContainer := obj.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
		backupGroup := obj.ObjectMeta.Labels["backup-group"]
		assert.Contains(t, backupGroup, "cluster0")
		assert.NotEmpty(t, backupGroup, "Backup group label is empty")
		assert.Subset(t, backupContainer.Args, []string{
			"--labels",
			fmt.Sprintf("pod-name sm-database-nuodb-cluster0-demo-hotcopy-%s",
				string(backupGroup[len(backupGroup)-1])),
		})
		expectedSchedule := options.SetValues["database.sm.hotCopy.fullSchedule"]
		if strings.Contains(obj.Name, "incremental") {
			expectedSchedule = options.SetValues["database.sm.hotCopy.incrementalSchedule"]
		} else if strings.Contains(obj.Name, "journal") {
			expectedSchedule = options.SetValues["database.sm.hotCopy.journalBackup.journalSchedule"]
		}
		assert.Equal(t, expectedSchedule, obj.Spec.Schedule)
	}
}

func TestDatabaseBackupGroupsCustom(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.replicas":                             "2",
			"database.sm.hotCopy.journalBackup.enabled":                "true",
			"database.sm.hotCopy.fullSchedule":                         "35 22 * * 6",
			"database.sm.hotCopy.incrementalSchedule":                  "35 22 * * 0-5",
			"database.sm.hotCopy.journalBackup.journalSchedule":        "?/15 * * * *",
			"database.sm.hotCopy.backupGroups.aws.labels":              "cloud aws",
			"database.sm.hotCopy.backupGroups.aws.fullSchedule":        "35 22 * * 1",
			"database.sm.hotCopy.backupGroups.aws.incrementalSchedule": "35 22 * * 2-7",
			"database.sm.hotCopy.backupGroups.gcp.labels":              "cloud gcp",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "database", []string{"templates/cronjob.yaml"})

	// one set of CronJobs per backup group is rendered
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 6) {
		require.NotEmpty(t, obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
		backupContainer := obj.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
		backupGroup := obj.ObjectMeta.Labels["backup-group"]
		assert.NotEmpty(t, backupGroup, "Backup group label is empty")
		if backupGroup == "aws" {
			assert.Subset(t, backupContainer.Args, []string{
				"--labels",
				"cloud aws",
			})
		} else {
			assert.Subset(t, backupContainer.Args, []string{
				"--labels",
				"cloud gcp",
			})
		}
		expectedSchedule := options.SetValues["database.sm.hotCopy.fullSchedule"]
		if backupGroup == "aws" {
			expectedSchedule = options.SetValues["database.sm.hotCopy.backupGroups.aws.fullSchedule"]
		}
		if strings.Contains(obj.Name, "incremental") {
			expectedSchedule = options.SetValues["database.sm.hotCopy.incrementalSchedule"]
			if backupGroup == "aws" {
				expectedSchedule = options.SetValues["database.sm.hotCopy.backupGroups.aws.incrementalSchedule"]
			}
		} else if strings.Contains(obj.Name, "journal") {
			expectedSchedule = options.SetValues["database.sm.hotCopy.journalBackup.journalSchedule"]
		}
		assert.Equal(t, expectedSchedule, obj.Spec.Schedule)
	}
}

func TestDatabaseBackupResourceLabels(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "true",
			"database.resourceLabels.foo":               "foo",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	// verify that default value for timeout is rendered
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 3) {
		verifyBackupResourceLabels(t, "release-name", options, &obj)
		verifyBackupResourceLabels(t, "release-name", options, &obj.Spec.JobTemplate)
	}
}
