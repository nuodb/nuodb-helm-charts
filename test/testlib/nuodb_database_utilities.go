package testlib

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	corev1 "k8s.io/api/core/v1"
)

func StartDatabase(t *testing.T, namespaceName string, adminPod string, options *helm.Options) (helmChartReleaseName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	helmChartReleaseName = fmt.Sprintf("database-%s", randomSuffix)
	tePodNameTemplate := fmt.Sprintf("te-%s", helmChartReleaseName)
	smPodName := fmt.Sprintf("sm-%s-nuodb-demo", helmChartReleaseName)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	// with Async actions which do not return a cleanup method, create the teardown(s) first
	AddTeardown(TEARDOWN_DATABASE, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
		AwaitNoPods(t, namespaceName, "database")
		DeleteDatabase(t, namespaceName, "demo", adminPod)
	})

	helm.Install(t, options, DATABASE_HELM_CHART_PATH, helmChartReleaseName)

	nrTePods, err := strconv.Atoi(options.SetValues["database.te.replicas"])
	if err != nil {
		nrTePods = 1
	}

	nrSmHotCopyPods, err := strconv.Atoi(options.SetValues["database.sm.hotCopy.replicas"])
	if err != nil {
		nrSmHotCopyPods = 1
	}

	nrSmNoHotCopyPods, err := strconv.Atoi(options.SetValues["database.sm.noHotCopy.replicas"])
	if err != nil {
		nrSmNoHotCopyPods = 0
	}

	nrSmPods := nrSmNoHotCopyPods + nrSmHotCopyPods

	AwaitNrReplicasScheduled(t, namespaceName, tePodNameTemplate, nrTePods)
	AwaitNrReplicasScheduled(t, namespaceName, smPodName, nrSmPods)

	tePodName := GetPodName(t, namespaceName, tePodNameTemplate)
	AwaitPodStatus(t, namespaceName, tePodName, corev1.PodReady, corev1.ConditionTrue, 120*time.Second)

	smPodName0 := GetPodName(t, namespaceName, smPodName)
	AwaitPodStatus(t, namespaceName, smPodName0, corev1.PodReady, corev1.ConditionTrue, 120*time.Second)

	AwaitDatabaseUp(t, namespaceName, adminPod, "demo", nrSmPods + nrTePods)

	return
}
