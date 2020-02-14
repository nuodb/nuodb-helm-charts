package helm

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
)

// Install will install the selected helm chart with the provided options under the given release name. This will fail
// the test if there is an error.
func Install(t *testing.T, options *Options, chart string, releaseName string) {
	require.NoError(t, InstallE(t, options, chart, releaseName))
}

// InstallE will install the selected helm chart with the provided options under the given release name.
func InstallE(t *testing.T, options *Options, chart string, releaseName string) error {
	// If the chart refers to a path, convert to absolute path. Otherwise, pass straight through as it may be a remote
	// chart.
	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}
		chart = absChartDir
	}

	// Now call out to helm install to install the charts with the provided options
	// Declare err here so that we can update args later
	var err error
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return err
	}
	args = append(args, "-n", releaseName, chart)
	_, err = RunHelmCommandAndGetOutputE(t, options, "install", args...)
	return err
}
