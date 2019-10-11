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

func labelMinikubeNode(t *testing.T, namespace string, labelName string, labelValue string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	var labelString string

	if labelValue != "" {
		labelString = fmt.Sprintf("%s=%s", labelName, labelValue)
	} else {
		labelString = fmt.Sprintf("%s-", labelName)
	}

	k8s.RunKubectl(t, options, "label", "nodes", "minikube", labelString, "--overwrite")
}

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
		SetFiles:  map[string]string{"thp.affinity": "../files/thp-affinity.yaml"},
		SetValues: map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	labelMinikubeNode(t, namespaceName, "failure-domain.beta.kubernetes.io/zone", randomSuffix)

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
		SetFiles:  map[string]string{"thp.affinity": "../files/thp-affinity.yaml"},
		SetValues: map[string]string{"thp.fullnameOverride": daemonName},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	labelMinikubeNode(t, namespaceName, "failure-domain.beta.kubernetes.io/zone", "")

	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	defer helm.Delete(t, options, helmChartReleaseName, true)

	daemonSet := testlib.GetDaemonSet(t, namespaceName, daemonName)
	assert.Assert(t, daemonSet.Status.DesiredNumberScheduled == 0)
}

func TestKubernetesDefaultMinikubeTHP(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := "../../stable/transparent-hugepage"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	namespaceName := fmt.Sprintf("testthp-%s", randomSuffix)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	options.KubectlOptions.Namespace = namespaceName

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	/*
		These tests do not verify that THP can be turned off via DaemonSet.
		Minikube already has TPH off by default.
		They only verify that scheduling works as expected and that the charts work
	*/

	t.Run("scheduleDefault", func(t *testing.T) { scheduleDefault(t, helmChartPath, namespaceName) })
	t.Run("scheduleLabel", func(t *testing.T) { scheduleLabel(t, helmChartPath, namespaceName) })
	t.Run("scheduleLabelMismatch", func(t *testing.T) { scheduleLabelMismatch(t, helmChartPath, namespaceName) })
}
