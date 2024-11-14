//go:build long
// +build long

package minikube

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	schedulingv1 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestPriorityClassNonexistent(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetespriorityclassnonexistent-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// create a domain and database with a custom priority class
	priorityClass := "test-priorityclass-" + randomSuffix
	options := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.noHotCopy.replicas":        "1",
			"admin.priorityClass":                   priorityClass,
			"database.priorityClasses.sm":           priorityClass,
			"database.priorityClasses.te":           priorityClass,
		},
	}

	// install admin helm chart without waiting for it to become running
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, _ := testlib.StartAdminNoWait(t, &options, 1, namespaceName)
	adminStatefulSet := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := adminStatefulSet + "-0"

	// install database helm chart without waiting for it to become running
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabaseNoWait(t, namespaceName, admin0, &options)
	smPodNameTemplate := fmt.Sprintf("sm-%s", databaseReleaseName)
	tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)

	// pods for the admin and database StatefulSets/Deployments should fail
	// to be created because the specified priority class does not exist
	testlib.Await(t, func() bool {
		filter := metav1.ListOptions{FieldSelector: "type=Warning,reason=FailedCreate"}
		events := testlib.GetEvents(t, namespaceName, filter)
		var expectedFailures []interface{}
		for _, event := range events {
			if !strings.Contains(event.Message, "no PriorityClass with name high-priority was found") {
				objectName := event.InvolvedObject.Name
				if strings.HasPrefix(objectName, adminStatefulSet) || strings.HasPrefix(objectName, smPodNameTemplate) || strings.HasPrefix(objectName, tePodNameTemplate) {
					expectedFailures = append(expectedFailures, event)
				}
			}
		}
		return len(expectedFailures) == 4
	}, time.Second*10)
}

func TestPriorityClass(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetespriorityclass-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	// create the priority class and make sure it is deleted when the test
	// finishes (priority classes are cluster-scoped resources, so they do
	// not go away when the namespace is deleted)
	priorityClass := "test-priorityclass-" + randomSuffix
	defer deletePriorityClass(t, priorityClass)
	var priorityClassValue int32 = int32(random.Random(1, 1000))
	createPriorityClass(t, priorityClass, priorityClassValue, "Priority class for testing")

	// create a domain and database with a custom priority class
	options := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"admin.priorityClass":                   priorityClass,
			"database.priorityClasses.sm":           priorityClass,
			"database.priorityClasses.te":           priorityClass,
		},
	}

	// sometimes the test fails because SMs doesn't go ready due to probe
	// timeout
	testlib.OverrideReadinessProbesTimeout(t, &options, "10")

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	adminStatefulSet := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := adminStatefulSet + "-0"

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
	smPodNameTemplate := fmt.Sprintf("sm-%s", databaseReleaseName)
	tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)

	// for all Pods, check the priorityClassName and priority
	for _, pod := range testlib.FindAllPodsInSchema(t, namespaceName) {
		if strings.HasPrefix(pod.Name, adminStatefulSet) || strings.HasPrefix(pod.Name, smPodNameTemplate) || strings.HasPrefix(pod.Name, tePodNameTemplate) {
			require.Equal(t, priorityClass, pod.Spec.PriorityClassName)
			require.Equal(t, priorityClassValue, *pod.Spec.Priority)
		}
	}
}

func createPriorityClass(t *testing.T, name string, value int32, description string) {
	options := k8s.NewKubectlOptions("", "", "")
	client, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	require.NoError(t, err)

	var priorityClass schedulingv1.PriorityClass
	priorityClass.Name = name
	priorityClass.Value = value
	priorityClass.Description = description
	_, err = client.SchedulingV1().PriorityClasses().Create(context.TODO(), &priorityClass, metav1.CreateOptions{})
	require.NoError(t, err)
}

func deletePriorityClass(t *testing.T, name string) {
	options := k8s.NewKubectlOptions("", "", "")
	client, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	require.NoError(t, err)

	err = client.SchedulingV1().PriorityClasses().Delete(context.TODO(), name, metav1.DeleteOptions{})
	require.NoError(t, err)
}
