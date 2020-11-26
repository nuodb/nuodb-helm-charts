package testlib

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	v12 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

var AdminRolesRequirePatching = false

func getFunctionCallerName() string {
	pc, _, _, _ := runtime.Caller(3)
	nameFull := runtime.FuncForPC(pc).Name() // main.foo
	nameEnd := filepath.Ext(nameFull)        // .foo
	name := strings.TrimPrefix(nameEnd, ".") // foo

	return name
}

func CreateNamespace(t *testing.T, namespaceName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	if IsOpenShiftEnvironment(t) {
		createOpenShiftProject(t, namespaceName)
	} else {
		k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	}

	// this method is async
	go GetK8sEventLog(t, namespaceName)

	AddTeardown(TEARDOWN_ADMIN, func() {
		k8s.DeleteNamespace(t, kubectlOptions, namespaceName)
	})
}

type AdminInstallationStep func(t *testing.T, options *helm.Options, helmChartReleaseName string)

func StartAdminTemplate(t *testing.T, options *helm.Options, replicaCount int, namespace string, installStep AdminInstallationStep) (helmChartReleaseName string, namespaceName string) {
	randomSuffix := strings.ToLower(random.UniqueId())

	helmChartReleaseName = fmt.Sprintf("admin-%s", randomSuffix)

	if namespace == "" {
		callerName := getFunctionCallerName()
		namespaceName = fmt.Sprintf("%s-%s", strings.ToLower(callerName), randomSuffix)

		CreateNamespace(t, namespaceName)
	} else {
		namespaceName = namespace
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	InjectTestValues(t, options)
	installStep(t, options, helmChartReleaseName)

	AddTeardown(TEARDOWN_ADMIN, func() {
		helm.Delete(t, options, helmChartReleaseName, true)
	})

	adminNames := make([]string, replicaCount)
	var adminStatefulSet string

	for i := 0; i < replicaCount; i++ {
		if options.SetValues["admin.fullnameOverride"] != "" {
			adminStatefulSet = fmt.Sprintf("%s", options.SetValues["admin.fullnameOverride"])
			adminNames[i] = fmt.Sprintf("%s-%d", adminStatefulSet, i)
		} else if options.SetValues["admin.nameOverride"] != "" {
			adminStatefulSet = fmt.Sprintf("%s-nuodb-cluster0-%s", helmChartReleaseName, options.SetValues["admin.nameOverride"])
			adminNames[i] = fmt.Sprintf("%s-%d", adminStatefulSet, i)
		} else {
			adminStatefulSet = fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
			adminNames[i] = fmt.Sprintf("%s-%d", adminStatefulSet, i)
		}
	}

	defer func() {
		// collect some useful diagnostics
		if t.Failed() {
			options := k8s.NewKubectlOptions("", "", namespaceName)
			// ignore any errors. This is already failed
			_ = k8s.RunKubectlE(t, options, "describe", "statefulset", adminStatefulSet)
		}
	}()

	if options.SetValues["admin.fullnameOverride"] != "" {
		AwaitNrReplicasScheduled(t, namespaceName, options.SetValues["admin.fullnameOverride"], replicaCount)
	} else {
		AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, replicaCount)
	}

	if AdminRolesRequirePatching {
		// workaround for https://github.com/nuodb/nuodb-helm-charts/issues/140
		output := helm.RenderTemplate(t, options, ADMIN_HELM_CHART_PATH, helmChartReleaseName, []string{"templates/role.yaml"})
		k8s.RunKubectl(t, kubectlOptions, "patch", "role", "nuodb-kube-inspector", "-p", output)
	}

	for i := 0; i < replicaCount; i++ {
		adminName := adminNames[i] // array will be out of scope for defer

		defer func() {
			if t.Failed() {
				options := k8s.NewKubectlOptions("", "", namespaceName)
				// ignore any errors. This is already failed
				_ = k8s.RunKubectlE(t, options, "describe", "pod", adminName)
			}
		}()

		// first await could be pulling the image from the repo
		AwaitPodUp(t, namespaceName, adminName, 300*time.Second)

		AddTeardown(TEARDOWN_ADMIN, func() {
			_, err := k8s.GetPodE(t, kubectlOptions, adminName)
			if err != nil {
				t.Logf("Admin pod '%s' is not available and logs can not be retrieved", adminName)
			} else {
				go GetAppLog(t, namespaceName, adminName, "", &v12.PodLogOptions{Follow: true})
				GetAdminEventLog(t, namespaceName, adminName)
			}
		})
	}

	for i := 0; i < replicaCount; i++ {
		AwaitAdminFullyConnected(t, namespaceName, adminNames[i], replicaCount)
	}

	return
}

func StartAdmin(t *testing.T, options *helm.Options, replicaCount int, namespace string) (string, string) {
	return StartAdminTemplate(t, options, replicaCount, namespace, func(t *testing.T, options *helm.Options, helmChartReleaseName string) {
		if options.Version == "" {
			helm.Install(t, options, ADMIN_HELM_CHART_PATH, helmChartReleaseName)
		} else {
			helm.Install(t, options, "nuodb/admin ", helmChartReleaseName)
		}
	})
}

func GetLoadBalancerPoliciesE(t *testing.T, namespaceName string, adminPod string) (map[string]NuoDBLoadBalancerPolicy, error) {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "load-balancers")
	if err == nil {
		err, policiesMap := UnmarshalLoadBalancerPolicies(output)
		return policiesMap, err
	}
	return nil, err
}

func GetLoadBalancerConfigE(t *testing.T, namespaceName string, adminPod string) ([]NuoDBLoadBalancerConfig, error) {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "load-balancer-config")
	if err == nil {
		err, configs := UnmarshalLoadBalancerConfigs(output)
		return configs, err
	}
	return nil, err
}

func AwaitNrLoadBalancerPolicies(t *testing.T, namespace string, podName string, expectedNumber int) {
	Await(t, func() bool {
		policies, err := GetLoadBalancerPoliciesE(t, namespace, podName)
		return err == nil && len(policies) == expectedNumber
	}, 30*time.Second)
}
