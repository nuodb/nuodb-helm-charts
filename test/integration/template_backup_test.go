package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func verifyBackupResourceLabels(t *testing.T, options *helm.Options, obj metav1.Object) {
	msg, ok := testlib.MapContains(obj.GetLabels(), map[string]string{
		"subgroup":     "backup",
		"backup-group": "cluster0-0",
	})
	require.Truef(t, ok, "Backup resource labels do not match for resource %s: %s", obj.GetName(), msg)
}

func truncate(s string, max int) string {
	return s[:max]
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

func TestDatabaseCronJobNames(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH
	templateFn := func(options *helm.Options) string {
		backupGroupPrefix := options.SetValues["database.sm.hotCopy.backupGroupPrefix"]
		if backupGroupPrefix == "" {
			// use the cluster name by default
			backupGroupPrefix = options.SetValues["cloud.cluster.name"]
		}
		backupGroupTemplate := fmt.Sprintf("%s-%%d", backupGroupPrefix)
		var hasBackupGroups bool
		for k := range options.SetValues {
			if strings.HasPrefix(k, "database.sm.hotCopy.backupGroups.") {
				hasBackupGroups = true
				break
			}
		}
		if hasBackupGroups {
			// leave the caller to specify the full backup group name
			backupGroupTemplate = "%s"
		}
		return fmt.Sprintf("%%s-hotcopy-%s-%s-%s",
			options.SetValues["admin.domain"],
			options.SetValues["database.name"],
			backupGroupTemplate)
	}
	assertContainsPrefix := func(t *testing.T, list map[string]interface{}, prefix string) {
		for k := range list {
			if strings.HasPrefix(k, prefix) {
				return
			}
		}
		assert.Fail(t, "%v does not contain element with prefix %q", list, prefix)
	}
	t.Run("testDefault", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"cloud.cluster.name":                        "cluster0",
				"admin.domain":                              "nuodb",
				"database.name":                             "demo",
				"database.sm.hotCopy.replicas":              "2",
				"database.sm.hotCopy.journalBackup.enabled": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/cronjob.yaml"})
		actual := make(map[string]interface{})
		for _, obj := range testlib.SplitAndRenderCronJob(t, output, 2) {
			actual[obj.Name] = struct{}{}
		}
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "full", 0))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "full", 1))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "incremental", 0))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "incremental", 1))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "journal", 0))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "journal", 1))
	})

	t.Run("testBackupGroups", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"cloud.cluster.name": "cluster0",
				"admin.domain":       "nuodb",
				"database.name":      "demo",
				"database.sm.hotCopy.journalBackup.enabled":          "true",
				"database.sm.hotCopy.backupGroups.aws.labels":        "cloud aws",
				"database.sm.hotCopy.backupGroups.gcp.processFilter": "label(cloud gcp)",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/cronjob.yaml"})
		actual := make(map[string]interface{})
		for _, obj := range testlib.SplitAndRenderCronJob(t, output, 2) {
			actual[obj.Name] = struct{}{}
		}
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "full", "aws"))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "full", "gcp"))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "incremental", "aws"))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "incremental", "gcp"))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "journal", "aws"))
		assert.Contains(t, actual, fmt.Sprintf(templateFn(options), "journal", "gcp"))
	})

	t.Run("testLongDatabaseName", func(t *testing.T) {
		options := &helm.Options{
			SetValues: map[string]string{
				"cloud.cluster.name":                        "cluster0",
				"admin.domain":                              "nuodb",
				"database.name":                             "superlongdatabasename",
				"database.sm.hotCopy.replicas":              "2",
				"database.sm.hotCopy.journalBackup.enabled": "true",
			},
		}

		output := helm.RenderTemplate(t, options, helmChartPath,
			"release-name", []string{"templates/cronjob.yaml"})
		actual := make(map[string]interface{})
		for _, obj := range testlib.SplitAndRenderCronJob(t, output, 2) {
			actual[obj.Name] = struct{}{}
		}
		// verify that unique names are generated for all CronJobs
		assert.Equal(t, 6, len(actual))
		for k := range actual {
			// verify that all CronJob names are less than 52 characters
			assert.LessOrEqual(t, len(k), 52, k)
		}
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "full", 0), 44))
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "full", 1), 44))
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "incremental", 0), 44))
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "incremental", 1), 44))
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "journal", 0), 44))
		assertContainsPrefix(t, actual, truncate(fmt.Sprintf(templateFn(options), "journal", 1), 44))
	})
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
		assert.Equal(t, job.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy, corev1.RestartPolicyOnFailure)
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
		assert.Equal(t, job.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy, corev1.RestartPolicyNever)
	}
}

func TestDatabaseBackupCronJobEnvironment(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := testlib.DATABASE_HELM_CHART_PATH

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.hotCopy.journalBackup.enabled": "true",
			"database.env[0].name":                      "foo",
			"database.env[0].value":                     "bar",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})

	for _, job := range testlib.SplitAndRenderCronJob(t, output, 3) {
		assert.Equal(t, 1, len(job.Spec.JobTemplate.Spec.Template.Spec.Containers))
		container := job.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
		testlib.EnvContains(container.Env, "foo", "bar")
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
		expectedBackupLabels := fmt.Sprintf(
			"pod-name sm-database-nuodb-cluster0-demo-hotcopy-%s",
			string(backupGroup[len(backupGroup)-1]))
		backupLabels := obj.Spec.JobTemplate.ObjectMeta.Annotations["backup-group-labels"]
		assert.Equal(t, expectedBackupLabels, backupLabels)
		assert.Subset(t, backupContainer.Args, []string{
			"--labels",
			expectedBackupLabels,
		})
		assert.Empty(t, obj.Spec.JobTemplate.ObjectMeta.Annotations["backup-group-process-filter"])
		expectedOperation := "full-hotcopy"
		expectedSchedule := options.SetValues["database.sm.hotCopy.fullSchedule"]
		if strings.Contains(obj.Name, "incremental") {
			expectedOperation = "incremental-hotcopy"
			expectedSchedule = options.SetValues["database.sm.hotCopy.incrementalSchedule"]
		} else if strings.Contains(obj.Name, "journal") {
			expectedOperation = "journal-hotcopy"
			expectedSchedule = options.SetValues["database.sm.hotCopy.journalBackup.journalSchedule"]
		}
		assert.Equal(t, expectedOperation, obj.Spec.JobTemplate.ObjectMeta.Annotations["operation"])
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
			"database.sm.hotCopy.backupGroups.gcp.processFilter":       "label(cloud gcp)",
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
			assert.Equal(t, "cloud aws", obj.Spec.JobTemplate.ObjectMeta.Annotations["backup-group-labels"])
			assert.Subset(t, backupContainer.Args, []string{
				"--labels",
				"cloud aws",
			})
		} else {
			assert.Equal(t, "label(cloud gcp)", obj.Spec.JobTemplate.ObjectMeta.Annotations["backup-group-process-filter"])
			assert.Subset(t, backupContainer.Args, []string{
				"--process-filter",
				"label(cloud gcp)",
			})
			assert.Empty(t, obj.Spec.JobTemplate.ObjectMeta.Annotations["backup-group-labels"])
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
		verifyDatabaseResourceLabels(t, "release-name", options, &obj)
		verifyBackupResourceLabels(t, options, &obj)
		verifyDatabaseResourceLabels(t, "release-name", options, &obj.Spec.JobTemplate)
		verifyBackupResourceLabels(t, options, &obj.Spec.JobTemplate)
	}
}

func TestDatabaseBackupCronJobLongName(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../../stable/database"

	options := &helm.Options{
		SetValues: map[string]string{
			"database.name": "really-long-name-xxxxxxxxxxxxxxx",
			"database.sm.hotCopy.journalBackup.enabled": "true",
		},
	}

	// Verify that hotcopy CronJob name doesn't exceed 52 chars
	output := helm.RenderTemplate(t, options, helmChartPath, "release-name", []string{"templates/cronjob.yaml"})
	for _, obj := range testlib.SplitAndRenderCronJob(t, output, 3) {
		assert.LessOrEqual(t, len(obj.Name), 52)
	}
}
