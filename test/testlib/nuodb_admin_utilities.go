package testlib

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func getFunctionCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	nameFull := runtime.FuncForPC(pc).Name() // main.foo
	nameEnd := filepath.Ext(nameFull)        // .foo
	name := strings.TrimPrefix(nameEnd, ".") // foo

	return name
}

func StartAdmin(t *testing.T, options *helm.Options, replicaCount int, namespace string) (helmChartReleaseName string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := ADMIN_HELM_CHART_PATH
	helmChartReleaseName = fmt.Sprintf("admin-%s", randomSuffix)

	adminNames := make([]string, replicaCount)

	for i := 0; i < replicaCount; i++ {
		adminNames[i] = fmt.Sprintf("%s-nuodb-cluster0-%d", helmChartReleaseName, i)
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions

	if namespace == "" {
		callerName := getFunctionCallerName()
		namespaceName = fmt.Sprintf("%s-%s", strings.ToLower(callerName), randomSuffix)
		k8s.CreateNamespace(t, kubectlOptions, namespaceName)
		AddTeardown(TEARDOWN_ADMIN, func() {
			GetK8sEventLog(t, namespaceName)
			k8s.DeleteNamespace(t, kubectlOptions, namespaceName)
		})
	} else {
		namespaceName = namespace
	}

	options.KubectlOptions.Namespace = namespaceName

	InjectTestVersion(t, options)
	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	AddTeardown("admin", func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})

	AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, replicaCount)

	for i := 0; i < replicaCount; i++ {
		adminName := adminNames[i] // array will be out of scope for defer

		defer func() {
			if t.Failed() {
				options := k8s.NewKubectlOptions("", "")
				options.Namespace = namespace
				// ignore any errors. This is already failed
				_ = k8s.RunKubectlE(t, options, "describe", "pod", adminName)
			}
		}()

		// first await could be pulling the image from the repo
		AwaitPodUp(t, namespaceName, adminName, 300*time.Second)
		AddTeardown("admin", func() {
			GetAppLog(t, namespaceName, adminName, "")
			GetAdminEventLog(t, namespaceName, adminName)
		})
	}

	for i := 0; i < replicaCount; i++ {
		AwaitAdminFullyConnected(t, namespaceName, adminNames[i], replicaCount)
	}

	return
}
