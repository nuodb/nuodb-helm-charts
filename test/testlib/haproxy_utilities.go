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
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func StartHAProxyIngress(t *testing.T, options *helm.Options, namespaceName string) string {
	randomSuffix := strings.ToLower(random.UniqueId())

	helmChartReleaseName := fmt.Sprintf("haproxy-%s", randomSuffix)

	defaultOptions := helm.Options{
		SetValues: map[string]string{
			"controller.image.tag":                 "1.10.11",
			"controller.replicaCount":              "1",
			"controller.service.type":              "NodePort",
			"controller.ingressClass":              helmChartReleaseName,
			"controller.ingressClassResource.name": helmChartReleaseName,
			"controller.resources.requests.cpu":    "150m",
		},
	}

	// set default values if not defined
	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}
	for k, v := range defaultOptions.SetValues {
		if _, ok := options.SetValues[k]; !ok {
			options.SetValues[k] = v
		}
	}
	// count expected pod number
	podCount, err := strconv.Atoi(options.SetValues["controller.replicaCount"])
	require.NoError(t, err)
	if options.SetValues["defaultBackend.enabled"] == "true" {
		c, err := strconv.Atoi(options.SetValues["defaultBackend.replicaCount"])
		require.NoError(t, err)
		podCount += c
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	helm.Install(t, options, "haproxytech/kubernetes-ingress", helmChartReleaseName)
	AddTeardown(TEARDOWN_HAPROXY, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})

	// there are two pods here, the haproxy, the default backend
	AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, podCount)
	haProxyPods := getHAProxyControllerPods(t, namespaceName, helmChartReleaseName)
	require.True(t, len(haProxyPods) > 0, "unable to find HAProxy controller pod")
	for _, haProxyPodName := range haProxyPods {
		AwaitPodUp(t, namespaceName, haProxyPodName, 300*time.Second)
	}

	AddTeardown(TEARDOWN_HAPROXY, func() {
		for _, haProxyPodName := range haProxyPods {
			if _, err := k8s.GetPodE(t, kubectlOptions, haProxyPodName); err != nil {
				t.Logf("HAProxy pod '%s' is not available and logs can not be retrieved", haProxyPodName)
			} else {
				go GetAppLog(t, namespaceName, haProxyPodName, "", &corev1.PodLogOptions{Follow: true})
			}
		}
	})

	return helmChartReleaseName
}

func getHAProxyControllerPods(t *testing.T, namespaceName string, helmChartReleaseName string) []string {
	haProxyNameTemplate := fmt.Sprintf("%s-kubernetes-ingress", helmChartReleaseName)
	pods, err := findPods(t, namespaceName, haProxyNameTemplate)
	require.NoError(t, err)
	var names []string
	for _, pod := range pods {
		// Filter out Kubernetes jobs
		if _, ok := pod.GetLabels()["job-name"]; !ok {
			names = append(names, pod.Name)
		}
	}
	return names
}
