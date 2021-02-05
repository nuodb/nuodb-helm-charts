// +build short

package minikube

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const LABEL_CLOUD = "minikube"
const LABEL_REGION = "local"
const LABEL_ZONE = "local-b"

func populateCreateDBData(t *testing.T, namespaceName string, adminPod string) {
	// populate some data
	opts := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.RunKubectl(t, opts,
		"exec", adminPod, "--",
		"/opt/nuodb/bin/nuosql",
		"--user", "dba",
		"--password", "secret",
		"demo",
		"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql",
	)

	// verify that the database contains the populated data
	tables, err := testlib.RunSQL(t, namespaceName, adminPod, "demo", "show schema User")
	require.NoError(t, err, "error running SQL: show schema User")
	require.True(t, strings.Contains(tables, "HOCKEY"), "tables returned: ", tables)
}

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
	options := k8s.NewKubectlOptions("", "", namespaceName)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", "demo")

	require.NoError(t, err, output)

	err, objects := testlib.Unmarshal(output)

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

func verifyPacketFetch(t *testing.T, namespaceName string, admin0 string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	// verify the container can actually download the file from the internet
	start := time.Now()
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions,
		"exec", admin0, "--",
		"bash", "-c",
		fmt.Sprintf("curl -k %s | tar tzf - ", testlib.IMPORT_ARCHIVE_URL),
	)
	require.NoError(t, err, "Could not fetch archive")
	elapsed := time.Since(start)
	t.Logf("Fetching package (%s) took %f seconds", testlib.IMPORT_ARCHIVE_URL, elapsed.Seconds())
	t.Log("tar contents: ", output)
}

func verifyEngineAltAddress(t *testing.T, namespaceName string, admin0 string, expectedNrEngines int) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	podNames, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--",
		"bash", "-c",
		"nuocmd get processes | grep -o \"(address=[^/]*\"| cut -f2 -d'='")
	require.NoError(t, err)
	podNamesSlice := strings.Split(strings.TrimSuffix(podNames, "\n"), "\n")

	actualAltAddresses, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--",
		"bash", "-c",
		"nuocmd get processes | grep -o \"alt-address: [^}]*\" | cut -f2 -d' '")
	require.NoError(t, err)
	actualAltAddressesSlice := strings.Split(strings.TrimSuffix(actualAltAddresses, "\n"), "\n")

	require.True(t, len(podNamesSlice) == expectedNrEngines, "Expected number of process names don't match")
	require.True(t, len(actualAltAddressesSlice) == expectedNrEngines, "Expected number of process addresses don't match")

	for index, podName := range podNamesSlice {
		pod := k8s.GetPod(t, kubectlOptions, podName)
		require.True(t, pod.Status.PodIP == actualAltAddressesSlice[index], "Expected alt-address doesn't match")
	}
}

func verifyBackup(t *testing.T, namespaceName string, podName string, databaseName string, options *helm.Options) {
	// verify that the backup has been documented by the Admin layer
	backupset, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions,
		"exec", podName, "--",
		"nuodocker", "get", "current-backup", "--db-name", databaseName,
	)

	require.NoError(t, err, "Error running: nuodocker get current-backup  ")
	require.True(t, backupset != "")
}

func isRestoreRequestSupported(t *testing.T, namespaceName string, podName string) bool {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	err := k8s.RunKubectlE(t, kubectlOptions, "exec", podName, "--",
		"bash", "-c", "nuodocker request restore -h > /dev/null")
	return err == nil
}

func checkRestoreRequests(t *testing.T, namespaceName string, podName string, databaseName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	restoreRequest, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "--",
		"nuocmd", "get", "value", "--key", fmt.Sprintf("/nuodb/nuosm/database/%s/restore", databaseName))
	require.NoError(t, err)
	require.Empty(t, restoreRequest, "Legacy restore request should be cleared")
	if isRestoreRequestSupported(t, namespaceName, podName) {
		restoreRequest, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", podName, "--",
			"nuodocker", "get", "restore-requests", "--db-name", databaseName)
		require.NoError(t, err)
		require.Empty(t, restoreRequest, "Database restore requests should be cleared")
	}
}

func awaitPodLog(t *testing.T, namespaceName string, podName string, fileNameSuffix string) {
	testlib.AwaitNrReplicasScheduled(t, namespaceName, podName, 1)
	testlib.AwaitPodPhase(t, namespaceName, podName, corev1.PodRunning, 30*time.Second)
	go testlib.GetAppLog(t, namespaceName, podName, fileNameSuffix, &corev1.PodLogOptions{Follow: true})
}

func awaitContainerLog(t *testing.T, namespaceName string, podName string, containerName string, fileNameSuffix string) {
	testlib.AwaitNrReplicasScheduled(t, namespaceName, podName, 1)
	testlib.AwaitContainerStarted(t, namespaceName, podName, containerName, 30*time.Second)
	go testlib.GetAppLog(t, namespaceName, podName, fileNameSuffix+"-"+containerName,
		&corev1.PodLogOptions{Container: containerName, Follow: true})
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

func TestKubernetesBackupDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		podName := testlib.GetPodName(t, namespaceName, "incremental-hotcopy-demo-cronjob")
		testlib.GetAppLog(t, namespaceName, podName, "", &corev1.PodLogOptions{})
	})

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
				// Configure more frequent incremental schedule so that
				// a full backup is created as a prerequisite.
				"database.sm.hotCopy.incrementalSchedule": "?/1 * * * *",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		testlib.AwaitJobSucceeded(t, namespaceName, "incremental-hotcopy-demo-cronjob", 120*time.Second)
		verifyBackup(t, namespaceName, admin0, "demo", &databaseOptions)
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
				"database.enableDaemonSet":              "true",
				// prevent non-backup SM from scheduling
				"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
				// Configure more frequent incremental schedule so that
				// a full backup is created as a prerequisite.
				"database.sm.hotCopy.incrementalSchedule": "?/1 * * * *",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		testlib.AwaitJobSucceeded(t, namespaceName, "incremental-hotcopy-demo-cronjob", 120*time.Second)
		verifyBackup(t, namespaceName, admin0, "demo", &databaseOptions)
	})
}

func TestKubernetesRestoreDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"backup.persistence.enabled":            "true",
			"backup.persistence.size":               "1Gi",
			"database.te.logPersistence.enabled":    "true",
		},
	}

	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", "demo-log-te-volume")
	})

	opts := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = opts

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	smPodName0 := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

	// Execute initial backup
	testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)

	populateCreateDBData(t, namespaceName, admin0)

	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	// Dump restore container log in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetAppLog(t, namespaceName, smPodName0, "_restore", &corev1.PodLogOptions{Container: "restore"})
	})

	// restore database
	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
	testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions, true)

	// verify that the database does NOT contain the data from AFTER the backup
	tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	require.NoError(t, err, "error running SQL: show schema User")
	require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)

	checkRestoreRequests(t, namespaceName, admin0, opt.DbName)

}

func TestKubernetesImportDatabase(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	adminOptions := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &adminOptions, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	verifyPacketFetch(t, namespaceName, admin0)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		databaseOptions := &helm.Options{
			SetValues: map[string]string{
				"database.autoImport.source":            "sftp://sftp.nuodb.com/incoming/restore.bak.tz",
				"database.autoImport.credentials":       "nuodb:wrongPass",
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		}

		// Install database and expect archive download failure during auto import
		databaseReleaseName := testlib.StartDatabaseNoWait(t, namespaceName, admin0, databaseOptions)

		opt := testlib.GetExtractedOptions(databaseOptions)
		smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
		testlib.AwaitNrReplicasScheduled(t, namespaceName, smPodNameTemplate, opt.NrSmPods)
		smPodName0 := fmt.Sprintf("%s-hotcopy-0", smPodNameTemplate)
		awaitPodLog(t, namespaceName, smPodName0, "_invalid-credentials")
		testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, smPodName0, 0, 120*time.Second)
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
			"Restore: unable to download/unpack backup", &corev1.PodLogOptions{Previous: true}), 1)

		// Use the correct IMPORT URL without credentials
		databaseOptions.SetValues["database.autoImport.source"] = testlib.IMPORT_ARCHIVE_URL
		databaseOptions.SetValues["database.autoImport.credentials"] = ""
		helm.Upgrade(t, databaseOptions, testlib.DATABASE_HELM_CHART_PATH, databaseReleaseName)

		smPod0 := testlib.GetPod(t, namespaceName, smPodName0)
		testlib.DeletePod(t, namespaceName, "pod/"+smPodName0)
		testlib.AwaitPodObjectRecreated(t, namespaceName, smPod0, 30*time.Second)
		testlib.AwaitPodUp(t, namespaceName, smPodName0, 120*time.Second)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "HOCKEY"))
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		testlib.StartDatabase(t, namespaceName, admin0,
			&helm.Options{
				SetValues: map[string]string{
					"database.autoImport.source":            testlib.IMPORT_ARCHIVE_URL,
					"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
					"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
					"database.enableDaemonSet":              "true",
					// prevent non-backup SM from scheduling
					"database.sm.nodeSelectorNoHotCopyDS.inexistantTag": "required",
				},
			},
		)

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "HOCKEY"))
	})
}

func TestKubernetesRestoreMultipleSMs(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" {
		t.Skip("Cannot test multiple SMs without the Enterprise Edition")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.noHotCopy.replicas":        "1",
			"database.te.logPersistence.enabled":    "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", "demo-log-te-volume")
	})

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodNameTemplate := fmt.Sprintf("%s-hotcopy", smPodNameTemplate)
	smPodName0 := fmt.Sprintf("%s-0", smPodNameTemplate)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)

	// Execute initial backup
	backupset := testlib.BackupDatabase(t, namespaceName, hcSmPodName0, opt.DbName, "full", opt.ClusterName)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)

	t.Run("restartAllAtOnce", func(t *testing.T) {
		populateCreateDBData(t, namespaceName, admin0)
		// restore database
		databaseOptions.SetValues["restore.source"] = backupset
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions, true)
		awaitPodLog(t, namespaceName, smPodName0, "_auto_post-restart")
		awaitContainerLog(t, namespaceName, smPodName0, "restore", "_auto_post-restart")
		awaitPodLog(t, namespaceName, hcSmPodName0, "_auto_post-restart")
		awaitContainerLog(t, namespaceName, hcSmPodName0, "restore", "_auto_post-restart")

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		checkRestoreRequests(t, namespaceName, admin0, opt.DbName)
	})

	t.Run("restartInOrder", func(t *testing.T) {
		testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
			testlib.GetAppLog(t, namespaceName, smPodName0, "_manual_restore", &corev1.PodLogOptions{Container: "restore"})
		})
		populateCreateDBData(t, namespaceName, admin0)
		// restore database
		testlib.RestoreDatabase(t, namespaceName, admin0, &databaseOptions, false)

		// Manually scale down all SMs
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=0")
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", hcSmPodNameTemplate, "--replicas=0")
		testlib.AwaitNoPods(t, namespaceName, smPodNameTemplate)

		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=1")
		awaitPodLog(t, namespaceName, smPodName0, "_manual_post-restart")
		awaitContainerLog(t, namespaceName, smPodName0, "restore", "_manual_post-restart")

		if isRestoreRequestSupported(t, namespaceName, admin0) {
			// If nonHC SM is started first it should wait for the database restore to complete
			testlib.Await(t, func() bool {
				return testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
					"Waiting for database restore to complete",
					&corev1.PodLogOptions{}) == 1
			}, 60*time.Second)
		} else {
			// If nonHC SM is started first it should fail as the restore source won't be available
			testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, smPodName0, 1, 120*time.Second)
			require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, smPodName0,
				fmt.Sprintf("Backupset %s cannot be found in /var/opt/nuodb/backup", backupset),
				&corev1.PodLogOptions{Previous: true}), 1)
		}

		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=0")
		testlib.AwaitNoPods(t, namespaceName, smPodNameTemplate)

		// Restart SM pods in the correct order so that HC SM performs the restore
		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", hcSmPodNameTemplate, "--replicas=1")
		awaitPodLog(t, namespaceName, hcSmPodName0, "_manual_post-restart")
		awaitContainerLog(t, namespaceName, hcSmPodName0, "restore", "_manual_post-restart")
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmHotCopyPods+opt.NrTePods)

		k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", smPodNameTemplate, "--replicas=1")
		awaitPodLog(t, namespaceName, smPodName0, "_manual_post-restore")
		awaitContainerLog(t, namespaceName, smPodName0, "restore", "_manual_post-restore")
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// verify that the database does NOT contain the data from AFTER the backup
		tables, err := testlib.RunSQL(t, namespaceName, admin0, opt.DbName, "show schema User")
		require.NoError(t, err, "error running SQL: show schema User")
		require.True(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
		checkRestoreRequests(t, namespaceName, admin0, opt.DbName)
	})
}

func TestKubernetesAutoRestore(t *testing.T) {
	if os.Getenv("NUODB_LICENSE") != "ENTERPRISE" {
		t.Skip("Cannot test autoRestore without the Enterprise Edition")
	}
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.name":                         "demo",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.autoRestore.source":           ":latest",
			"database.sm.noHotCopy.replicas":        "1",
			"database.te.logPersistence.enabled":    "true",
			"database.env[0].name":                  "NUODB_DEBUG",
			"database.env[0].value":                 "debug",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0)
		testlib.RecoverCoresFromEngine(t, namespaceName, "te", "demo-log-te-volume")
	})

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	databaseOptions.KubectlOptions = kubectlOptions

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	hcSmPodNameTemplate := fmt.Sprintf("%s-hotcopy", smPodNameTemplate)
	smPodName0 := fmt.Sprintf("%s-0", smPodNameTemplate)
	hcSmPodName0 := fmt.Sprintf("%s-0", hcSmPodNameTemplate)
	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)

	go testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, hcSmPodName0, "_pre-restart", &corev1.PodLogOptions{Follow: true})

	populateCreateDBData(t, namespaceName, admin0)
	backupset := testlib.BackupDatabase(t, namespaceName, smPodName0, opt.DbName, "full", opt.ClusterName)

	removeArchiveData := func(podName string) {
		// Remove archive data and restart the pod
		k8s.RunKubectl(t, kubectlOptions, "exec", podName, "--", "rm", "-rf", "/var/opt/nuodb/archive/nuodb/demo")
		// Engine should fail with ASSERT and pod will be restarted
		testlib.RunSQL(t, namespaceName, admin0, "demo", "select * from system.nodes")
		testlib.AwaitPodRestartCountGreaterThan(t, namespaceName, podName, 0, 30*time.Second)
		awaitPodLog(t, namespaceName, podName, "_post-restart")
	}

	t.Run("restartHotCopySM", func(t *testing.T) {
		removeArchiveData(hcSmPodName0)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		// HC SM should restore the archive from the latest backup
		require.GreaterOrEqual(t, testlib.GetStringOccurrenceInLog(t, namespaceName, hcSmPodName0,
			fmt.Sprintf("Finished restoring /var/opt/nuodb/backup/%s to /var/opt/nuodb/archive/nuodb/demo", backupset),
			&corev1.PodLogOptions{}), 1)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})

	t.Run("restartNonHotCopySM", func(t *testing.T) {
		removeArchiveData(hcSmPodName0)
		// nonHC SM should remove the archive metadata and SYNC the data from other SM
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)
		testlib.CheckArchives(t, namespaceName, admin0, opt.DbName, 2, 0)
	})
}
