// +build short

package minikube

import (
	"fmt"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"gotest.tools/assert"
)


func scheduleDefault(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		SetValues: map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
	assert.Assert(t, daemonSet.Status.DesiredNumberScheduled == 1)
}

func scheduleLabel(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		ValuesFiles: []string{"../files/thp-affinity.yaml"},
		SetValues:   map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	testlib.LabelNodes(t, namespaceName, "test.nuodb.com/zone", randomSuffix)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
	assert.Assert(t, daemonSet.Status.DesiredNumberScheduled == 1)
}

func scheduleLabelMismatch(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		ValuesFiles: []string{"../files/thp-affinity.yaml"},
		SetValues:   map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	testlib.LabelNodes(t, namespaceName, "test.nuodb.com/zone", "")

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
	assert.Assert(t, daemonSet.Status.DesiredNumberScheduled == 0)
}

func TestKubernetesDefaultMinikubeTHP(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN) // some namespace cleanup

	namespaceName := fmt.Sprintf("testthp-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	/*
		These tests do not verify that THP can be turned off via DaemonSet.
		Minikube already has TPH off by default.
		They only verify that scheduling works as expected and that the charts work
	*/

	t.Run("scheduleDefault", func(t *testing.T) { scheduleDefault(t, testlib.THP_HELM_CHART_PATH, namespaceName) })
	t.Run("scheduleLabel", func(t *testing.T) { scheduleLabel(t, testlib.THP_HELM_CHART_PATH, namespaceName) })
	t.Run("scheduleLabelMismatch", func(t *testing.T) { scheduleLabelMismatch(t, testlib.THP_HELM_CHART_PATH, namespaceName) })
}
