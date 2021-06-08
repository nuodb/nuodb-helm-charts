// +build short

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const LABEL_CLOUD = "minikube"
const LABEL_REGION = "local"
const LABEL_ZONE = "local-b"

func verifyKubernetesAccess(t *testing.T, namespaceName string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	serviceAccountDir := "/var/run/secrets/kubernetes.io/serviceaccount"

	// check namespace matches service account directory
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", serviceAccountDir+"/namespace")
	require.NoError(t, err, output)
	require.Equal(t, namespaceName, output)

	// get authorization token
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", serviceAccountDir+"/token")
	require.NoError(t, err, output)

	curlCmdPrefix := fmt.Sprintf("curl -s --cacert %s -H 'Authorization: Bearer %s' https://kubernetes.default.svc", serviceAccountDir+"/ca.crt", output)

	// check that we can access Pods
	url := "/api/v1/namespaces/" + namespaceName + "/pods"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.NoError(t, err, output)
	require.True(t, strings.Contains(output, "\"kind\": \"PodList\""), output)

	// check that we can access this Pod
	url = "/api/v1/namespaces/" + namespaceName + "/pods/" + podName
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.NoError(t, err, output)
	require.True(t, strings.Contains(output, "\"kind\": \"Pod\""), output)

	// check that we can access PersistentVolumeClaims
	url = "/api/v1/namespaces/" + namespaceName + "/persistentvolumeclaims"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"PersistentVolumeClaimList\""), output)

	// check that we can access Deployments
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/deployments"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"DeploymentList\""), output)

	// check that we can access StatefulSets
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/statefulsets"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"StatefulSetList\""), output)

	// check that we can create Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases"
	// when request data is specified without an explicit request method, POST is assumed
	leaseName := strings.ToLower(random.UniqueId())
	extraArgs := fmt.Sprintf(" -H 'Content-Type: application/json' -d  '{\"metadata\": {\"name\": \"%s\"}}'", leaseName)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can update Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases/" + leaseName
	// use create response as request payload, which contains the correct resourceVersion (update fails if the resourceVersion does not match)
	extraArgs = fmt.Sprintf(" -X PUT -H 'Content-Type: application/json' -d '%s'", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can get Leases
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)
}

func verifyNuoSQL(t *testing.T, namespaceName string, adminPod string, databaseName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--", "bash", "-c",
		fmt.Sprintf("echo \"select * from system.nodes;\" | nuosql %s@localhost --user dba --password secret", databaseName))

	require.NoError(t, err, output)

	require.True(t, strings.Contains(output, "Storage"))
	require.True(t, strings.Contains(output, "Transaction"))
}

func verifySecret(t *testing.T, namespaceName string) {
	secret := testlib.GetSecret(t, namespaceName, "demo.nuodb.com")

	_, ok := secret.Data["database-name"]
	require.True(t, ok)

	_, ok = secret.Data["database-password"]
	require.True(t, ok)

	_, ok = secret.Data["database-username"]
	require.True(t, ok)
}

func verifyDBService(t *testing.T, namespaceName string, podName string, serviceName string, ping bool) {

	dBService := testlib.GetService(t, namespaceName, serviceName)
	require.Equal(t, dBService.Name, serviceName)

	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
}

func verifyPodLabeling(t *testing.T, namespaceName string, adminPod string) {
	objects, err := testlib.GetDatabaseProcessesE(t, namespaceName, adminPod, "demo")
	require.NoError(t, err)

	for _, obj := range objects {
		val, ok := obj.Labels["cloud"]
		require.True(t, ok)
		require.True(t, val == LABEL_CLOUD)

		val, ok = obj.Labels["region"]
		require.True(t, ok)
		require.True(t, val == LABEL_REGION)

		val, ok = obj.Labels["zone"]
		require.True(t, ok)
		require.True(t, val == LABEL_ZONE)
	}

}

func verifyEngineAltAddress(t *testing.T, namespaceName string, adminPod string, expectedNrEngines int) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	objects, err := testlib.GetDatabaseProcessesE(t, namespaceName, adminPod, "demo")
	require.NoError(t, err)
	require.EqualValues(t, len(objects), expectedNrEngines)

	for _, obj := range objects {
		pod := k8s.GetPod(t, kubectlOptions, obj.Hostname)
		require.EqualValues(t, pod.Status.PodIP, obj.Options["alt-address"])
	}
}

func TestKubernetesBasicDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	headlessServiceName := fmt.Sprintf("demo")
	clusterServiceName := fmt.Sprintf("demo-clusterip")

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.labels.cloud":              LABEL_CLOUD,
				"database.te.labels.region":             LABEL_REGION,
				"database.te.labels.zone":               LABEL_ZONE,
				"database.sm.labels.cloud":              LABEL_CLOUD,
				"database.sm.labels.region":             LABEL_REGION,
				"database.sm.labels.zone":               LABEL_ZONE,
			},
		})

		t.Run("verifySecret", func(t *testing.T) { verifySecret(t, namespaceName) })
		t.Run("verifyDBHeadlessService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0, headlessServiceName, true) })
		t.Run("verifyDBClusterService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0, clusterServiceName, false) })
		t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0, "demo") })
		t.Run("verifyPodLabeling", func(t *testing.T) { verifyPodLabeling(t, namespaceName, admin0) })
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		testlib.StartDatabase(t, namespaceName, admin0,
			&helm.Options{
				SetValues: map[string]string{
					"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"database.te.labels.cloud":              LABEL_CLOUD,
					"database.te.labels.region":             LABEL_REGION,
					"database.te.labels.zone":               LABEL_ZONE,
					"database.sm.labels.cloud":              LABEL_CLOUD,
					"database.sm.labels.region":             LABEL_REGION,
					"database.sm.labels.zone":               LABEL_ZONE,
					"database.enableDaemonSet":              "true",
					// prevent non-backup SM from scheduling
					"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
				},
			},
		)

		t.Run("verifySecret", func(t *testing.T) { verifySecret(t, namespaceName) })
		t.Run("verifyDBHeadlessService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0, headlessServiceName, true) })
		t.Run("verifyDBClusterService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0, clusterServiceName, false) })
		t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0, "demo") })
		t.Run("verifyPodLabeling", func(t *testing.T) { verifyPodLabeling(t, namespaceName, admin0) })
	})

	t.Run("startDatabaseStatefulSetMultiTenant", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.name":                         "green",
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			},
		})

		t.Run("verifyNuoSQL-green", func(t *testing.T) {
			verifyNuoSQL(t, namespaceName, admin0, "green")
		})

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.name":                         "blue",
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			},
		})

		t.Run("verifyNuoSQL-blue", func(t *testing.T) {
			verifyNuoSQL(t, namespaceName, admin0, "blue")
		})
	})
}

func TestKubernetesAccessWithinPods(t *testing.T) {
	testlib.AwaitTillerUp(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseVerifyKubeAccess", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			},
		}
		databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		t.Run("verifyKubernetesAccess", func(t *testing.T) {
			// verify that Admin Pod can invoke K8s REST APIs
			verifyKubernetesAccess(t, namespaceName, admin0)

			// verify that SM and TE Pods can invoke K8s REST APIs
			opt := testlib.GetExtractedOptions(&databaseOptions)
			tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
			smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
			tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
			smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
			verifyKubernetesAccess(t, namespaceName, tePodName)
			verifyKubernetesAccess(t, namespaceName, smPodName)
		})
	})
}

func TestKubernetesAltAddress(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSetWithAltAddress", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.engineOptions.alt-address": "$(NUODB_ALT_ADDRESS)",
				"database.te.engineOptions.alt-address": "$(NUODB_ALT_ADDRESS)",
			},
			ValuesFiles: []string{"../files/database-env.yaml"},
		})
		expectedNrEngines := 2
		t.Run("verifyEnginesAltAddress", func(t *testing.T) { verifyEngineAltAddress(t, namespaceName, admin0, expectedNrEngines) })
	})
}

func TestKubernetesStartDatabaseShrinkedAdmin(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas": "3",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 3, "")

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	admin := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := fmt.Sprintf("%s-0", admin)

	// scale down the APs to 2; KAA will automatically evict the shut down AP
	k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", admin, "--replicas=2")
	testlib.AwaitServerState(
		t, namespaceName, admin0, fmt.Sprintf("%s-2", admin), "Disconnected", 30*time.Second)
	testlib.AwaitNrReplicasReady(t, namespaceName, admin, 2)
	// restart the current leader to bounce it
	leader := testlib.AwaitDomainLeader(t, namespaceName, admin0, 20*time.Second)
	testlib.DeletePod(t, namespaceName, "pod/"+leader)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, leader, 1)
	testlib.AwaitPodUp(t, namespaceName, leader, 90*time.Second)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	// make sure that SM database processes can start
	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})
}

func TestKubernetesSeparateJournalLocation(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		options := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.journalPath.enabled":       "true",
			},
		}

		databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
		t.Logf(databaseReleaseName)

	})
}
