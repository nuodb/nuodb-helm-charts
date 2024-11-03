//go:build short
// +build short

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func scheduleDefault(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		SetValues: map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.Await(t, func() bool {
		daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
		return daemonSet.Status.DesiredNumberScheduled == 1
	}, 30*time.Second)
}

func scheduleLabel(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		ValuesFiles: []string{"../files/thp-affinity.yaml"},
		SetValues:   map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	testlib.LabelNodes(t, namespaceName, "test.nuodb.com/zone", randomSuffix)

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.Await(t, func() bool {
		daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
		return daemonSet.Status.DesiredNumberScheduled == 1
	}, 30*time.Second)
}

func scheduleLabelMismatch(t *testing.T, helmChartPath string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())
	helmChartReleaseName := fmt.Sprintf("thp-%s", randomSuffix)
	daemonName := fmt.Sprintf("%s-%s", "transparent-hugepage", randomSuffix)

	options := &helm.Options{
		ValuesFiles: []string{"../files/thp-affinity.yaml"},
		SetValues:   map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions

	testlib.LabelNodes(t, namespaceName, "test.nuodb.com/zone", "")

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	testlib.Await(t, func() bool {
		daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
		return daemonSet.Status.DesiredNumberScheduled == 0
	}, 30*time.Second)
}

func TestKubernetesDefaultMinikubeTHP(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN) // some namespace cleanup

	namespaceName := fmt.Sprintf("%sthp-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
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
