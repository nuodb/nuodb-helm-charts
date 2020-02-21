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

func CreateNamespace(t *testing.T, namespaceName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")

	if IsOpenShiftEnvironment(t) {
		createOpenShiftProject(t, namespaceName)
	} else {
		k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	}

	AddTeardown(TEARDOWN_ADMIN, func() {
		GetK8sEventLog(t, namespaceName)
		k8s.DeleteNamespace(t, kubectlOptions, namespaceName)
	})
}

func StartAdmin(t *testing.T, options *helm.Options, replicaCount int, namespace string) (helmChartReleaseName string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	// Path to the helm chart we will test
	helmChartPath := ADMIN_HELM_CHART_PATH
	helmChartReleaseName = fmt.Sprintf("admin-%s", randomSuffix)

	if namespace == "" {
		callerName := getFunctionCallerName()
		namespaceName = fmt.Sprintf("%s-%s", strings.ToLower(callerName), randomSuffix)

		CreateNamespace(t, namespaceName)
	} else {
		namespaceName = namespace
	}

	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

    InjectTestVersion(t, options)
	InjectOpenShiftValues(t, options)
	helm.Install(t, options, helmChartPath, helmChartReleaseName)

	AddTeardown(TEARDOWN_ADMIN, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})


	adminNames := make([]string, replicaCount)

	for i := 0; i < replicaCount; i++ {
		if options.SetValues["admin.fullnameOverride"] != "" {
			adminNames[i] = fmt.Sprintf("%s-%d", options.SetValues["admin.fullnameOverride"], i)
		} else if  options.SetValues["admin.nameOverride"] != "" {
			adminNames[i] = fmt.Sprintf("%s-nuodb-cluster0-%s-%d", helmChartReleaseName, options.SetValues["admin.nameOverride"], i)
		} else {
			adminNames[i] = fmt.Sprintf("%s-nuodb-cluster0-%d", helmChartReleaseName, i)
		}
	}

	if options.SetValues["admin.fullnameOverride"] != "" {
		AwaitNrReplicasScheduled(t, namespaceName, options.SetValues["admin.fullnameOverride"], replicaCount)
	} else {
		AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, replicaCount)
	}

	for i := 0; i < replicaCount; i++ {
		adminName := adminNames[i] // array will be out of scope for defer

		defer func() {
			if (t.Failed()) {
				options := k8s.NewKubectlOptions("", "")
				options.Namespace = namespace
				// ignore any errors. This is already failed
				_ = k8s.RunKubectlE(t, options, "describe", "pod", adminName)
			}
		}()

		// first await could be pulling the image from the repo
		AwaitAdminPodUp(t, namespaceName, adminName, 300*time.Second)
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
