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
	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

const LABEL_CLOUD = "minikube"
const LABEL_REGION = "local"
const LABEL_ZONE = "local-b"

func verifyKubernetesAccess(t *testing.T, namespaceName string, podName string, containerName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	serviceAccountDir := "/var/run/secrets/kubernetes.io/serviceaccount"

	// check namespace matches service account directory
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "cat", serviceAccountDir+"/namespace")
	require.NoError(t, err, output)
	require.Equal(t, namespaceName, output)

	// get authorization token
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "cat", serviceAccountDir+"/token")
	require.NoError(t, err, output)

	curlCmdPrefix := fmt.Sprintf("curl -s --cacert %s -H 'Authorization: Bearer %s' https://kubernetes.default.svc", serviceAccountDir+"/ca.crt", output)

	// check that we can access Pods
	url := "/api/v1/namespaces/" + namespaceName + "/pods"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.NoError(t, err, output)
	require.True(t, strings.Contains(output, "\"kind\": \"PodList\""), output)

	// check that we can access this Pod
	url = "/api/v1/namespaces/" + namespaceName + "/pods/" + podName
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.NoError(t, err, output)
	require.True(t, strings.Contains(output, "\"kind\": \"Pod\""), output)

	// check that we can access PersistentVolumeClaims
	url = "/api/v1/namespaces/" + namespaceName + "/persistentvolumeclaims"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"PersistentVolumeClaimList\""), output)

	// check that we can access Deployments
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/deployments"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"DeploymentList\""), output)

	// check that we can access StatefulSets
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/statefulsets"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"StatefulSetList\""), output)

	// check that we can create Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases"
	// when request data is specified without an explicit request method, POST is assumed
	leaseName := strings.ToLower(random.UniqueId())
	extraArgs := fmt.Sprintf(" -H 'Content-Type: application/json' -d  '{\"metadata\": {\"name\": \"%s\"}}'", leaseName)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can update Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases/" + leaseName
	// use create response as request payload, which contains the correct resourceVersion (update fails if the resourceVersion does not match)
	extraArgs = fmt.Sprintf(" -X PUT -H 'Content-Type: application/json' -d '%s'", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can get Leases
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", containerName, "--", "bash", "-c", curlCmdPrefix+url)
	require.True(t, strings.Contains(output, "\"kind\": \"Lease\""), output)
}

func verifyNuoSQL(t *testing.T, namespaceName string, adminPod string, databaseName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "-c", "admin", "-c", "admin", "--", "bash", "-c",
		fmt.Sprintf("echo \"select * from system.nodes;\" | nuosql %s@localhost --user dba --password secret", databaseName))

	require.NoError(t, err, output)

	require.True(t, strings.Contains(output, "Storage"))
	require.True(t, strings.Contains(output, "Transaction"))
}

func verifySecret(t *testing.T, namespaceName string) {
	secret := testlib.GetSecret(t, namespaceName, "nuodb-demo")

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
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
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
	clusterServiceName := fmt.Sprintf("%s-nuodb-cluster0-demo-clusterip", databaseReleaseName)

	t.Run("verifySecret", func(t *testing.T) { verifySecret(t, namespaceName) })
	t.Run("verifyDBClusterService", func(t *testing.T) { verifyDBService(t, namespaceName, admin0, clusterServiceName, false) })
	t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0, "demo") })
	t.Run("verifyPodLabeling", func(t *testing.T) { verifyPodLabeling(t, namespaceName, admin0) })
}

func TestSmVolumePermissionChange(t *testing.T) {

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	options := helm.Options{}
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.hotCopy.journalPath.enabled": "true",
		},
	}
	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	opt := testlib.GetExtractedOptions(&databaseOptions)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)

	// simulate the creation of a real filesystem by creating a lost+found
	// directory in all mounted volumes (archive, journal, and backup); with
	// certain PV provisioners, the directory is owned by root and has
	// permissions 700 (not group-writable)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	for _, dir := range []string{"/var/opt/nuodb/archive", "/var/opt/nuodb/journal", "/var/opt/nuodb/backup"} {
		k8s.RunKubectl(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--", "mkdir", dir+"/lost+found")
		k8s.RunKubectl(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--", "chmod", "700", dir+"/lost+found")

		output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--", "stat", "-c", "%a", dir+"/lost+found")
		require.NoError(t, err, output)
		require.Equal(t, "700", strings.TrimSpace(output))
	}

	// delete SM pod to cause init container to be invoked
	testlib.AwaitDatabaseRestart(t, namespaceName, admin0, opt.DbName, &databaseOptions, func() {
		k8s.RunKubectl(t, kubectlOptions, "delete", "pod", smPodName)
	})

	// check that directories group-writable after restart
	for _, dir := range []string{"/var/opt/nuodb/archive", "/var/opt/nuodb/journal", "/var/opt/nuodb/backup"} {
		output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--", "stat", "-c", "%a", dir+"/lost+found")
		require.NoError(t, err, output)
		require.Equal(t, "770", strings.TrimSpace(output))
	}
}

func TestKubernetesAccessWithinPods(t *testing.T) {

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
			verifyKubernetesAccess(t, namespaceName, admin0, "admin")

			// verify that SM and TE Pods can invoke K8s REST APIs
			opt := testlib.GetExtractedOptions(&databaseOptions)
			tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
			smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
			tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
			smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
			verifyKubernetesAccess(t, namespaceName, tePodName, "engine")
			verifyKubernetesAccess(t, namespaceName, smPodName, "engine")
		})
	})
}

func TestArchivePvcRecreated(t *testing.T) {
	options := helm.Options{}
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// Start database and wait for it to become running
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

	// Suspend admin process so that it cannot handle deletion of the archive PVC
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.RunKubectl(t, kubectlOptions, "exec", admin0, "-c", "admin", "--", "bash", "-c",
		"pgrep -x java | xargs -r kill -STOP")

	// Delete archive PVC and SM pod
	opt := testlib.GetExtractedOptions(&databaseOptions)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
	archivePvcName := "archive-volume-" + smPodName

	// Save archive PVC UID so that we can detect re-creation
	archivePvcUid, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "pvc/"+archivePvcName, "-o", "jsonpath={.metadata.uid}")
	require.NoError(t, err)
	k8s.RunKubectl(t, kubectlOptions, "delete", "pvc/"+archivePvcName, "pod/"+smPodName)

	// Wait for archive PVC to be re-created, which should happen
	// automatically due to the SM statefulset
	testlib.Await(t, func() bool {
		newArchivePvcUid, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "pvc/"+archivePvcName, "-o", "jsonpath={.metadata.uid}")
		require.NoError(t, err)
		return newArchivePvcUid != "" && newArchivePvcUid != archivePvcUid
	}, 120*time.Second)

	// Resume admin process
	k8s.RunKubectl(t, kubectlOptions, "exec", admin0, "-c", "admin", "--", "bash", "-c",
		"pgrep -x java | xargs -r kill -CONT")

	// If `nuodocker start sm` detects PVC re-creation, expect archive for
	// deleted PVC to be deleted, allowing database restart.
	minVersion := "6.0.2"
	t.Run("testOldArchiveDeleted", func(t *testing.T) {
		testlib.SkipTestOnNuoDBVersionCondition(t, "<"+minVersion)
		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, opt.NrSmPods+opt.NrTePods)

		// Check that there is one archive
		archives, removedArchives := testlib.GetArchives(t, namespaceName, admin0, opt.DbName)
		require.NoError(t, err)
		require.Equal(t, 1, len(archives))
		require.Equal(t, 0, len(removedArchives))

		// Look for expected SM logging
		out, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "logs", smPodName, "-c", "engine")
		require.NoError(t, err)
		require.Contains(t, out, "Deleting archive ID 0 associated with PVC "+archivePvcName)
	})
	// Otherwise, database restart will be blocked in
	// AWAITING_ARCHIVE_HISTORIES_MSG state waiting for the archive for the
	// deleted PVC to be started.
	t.Run("testOldArchiveNotDeleted", func(t *testing.T) {
		testlib.SkipTestOnNuoDBVersionCondition(t, ">="+minVersion)
		testlib.Await(t, func() bool {
			out, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "-c", "admin", "--",
				"nuocmd", "show", "domain", "--process-format", "{hostname} {durable_state}:{engine_state}")
			require.NoError(t, err)
			return strings.Contains(out, opt.DbName+" [state = AWAITING_ARCHIVE_HISTORIES_MSG]") &&
				strings.Contains(out, smPodName+" TRACKED:STARTING_UP")
		}, 120*time.Second)

		// Check that there are two archives
		archives, removedArchives := testlib.GetArchives(t, namespaceName, admin0, opt.DbName)
		require.NoError(t, err)
		require.Equal(t, 2, len(archives))
		require.Equal(t, 0, len(removedArchives))
	})
}

func TestHooksCustomHandlers(t *testing.T) {
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	adminReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)

	// Start database and wait for it to become running
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseOptions := helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                      testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                      testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.hotCopy.enablePod":                           "false",
			"database.sm.noHotCopy.replicas":                          "1",
			"database.backupHooks.enabled":                            "true",
			"database.backupHooks.customHandlers[0].method":           "GET",
			"database.backupHooks.customHandlers[0].path":             "/pid/{command}",
			"database.backupHooks.customHandlers[0].script":           `pgrep -x "${command}"`,
			"database.backupHooks.customHandlers[0].statusMappings.1": "404",
			"database.backupHooks.customHandlers[1].method":           "POST",
			"database.backupHooks.customHandlers[1].path":             "operation/{command}/execute",
			"database.backupHooks.customHandlers[1].script":           `"${command}" ${payload}`,
			"database.backupHooks.customHandlers[1].statusMappings.1": "400",
			"database.backupHooks.customHandlers[1].statusMappings.*": "500",
		},
	}
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Execute custom hooks from SM container
	opt := testlib.GetExtractedOptions(&databaseOptions)
	smPodNameTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseReleaseName, opt.ClusterName, opt.DbName)
	smPodName := testlib.GetPodName(t, namespaceName, smPodNameTemplate)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	pidFromHook, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "http://localhost:8000/pid/nuodb")
	require.NoError(t, err)

	// Get pid for nuodb process and make sure it matches pid from hook
	pid, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"pgrep", "-x", "nuodb")
	require.NoError(t, err)
	require.Equal(t, pid, pidFromHook)

	// Check HTTP code from hook response on success
	httpCode, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "-w", "%{http_code}", "http://localhost:8000/pid/nuodb", "-o", "/dev/null")
	require.NoError(t, err)
	require.Equal(t, "200", strings.TrimSpace(httpCode))

	// Check HTTP code from hook response on failure and content-type
	out, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "-w", "--- %{http_code} %{content_type} ---", "http://localhost:8000/pid/not_found")
	require.NoError(t, err)
	require.Contains(t, out, "--- 404 text/plain ---")

	// Execute arbitrary command using 'POST operation/{command}/execute'
	jsonData := `{"one": 1}`
	out, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "http://localhost:8000/operation/echo/execute", "-d", jsonData)
	require.NoError(t, err)
	require.Equal(t, jsonData, strings.TrimSpace(out))

	// Check content-type for JSON data
	contentType, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "-w", "%{content_type}", "http://localhost:8000/operation/echo/execute", "-d", jsonData, "-o", "/dev/null")
	require.NoError(t, err)
	require.Equal(t, "application/json", strings.TrimSpace(contentType))

	// Check HTTP code from hook response on exit code 1
	httpCode, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "-w", "%{http_code}", "http://localhost:8000/operation/exit/execute", "-d", "1")
	require.NoError(t, err)
	require.Equal(t, "400", strings.TrimSpace(httpCode))

	// Check HTTP code from hook response on arbitrary non-0 exit code
	httpCode, err = k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPodName, "-c", "engine", "--",
		"curl", "-s", "-w", "%{http_code}", "http://localhost:8000/operation/exit/execute", "-d", "123")
	require.NoError(t, err)
	require.Equal(t, "500", strings.TrimSpace(httpCode))
}

func TestKubernetesAltAddress(t *testing.T) {
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

	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_ADMIN, t, func() {
		k8s.RunKubectl(t, kubectlOptions, "get", "pods", "-o", "wide")
		testlib.DescribePods(t, namespaceName, admin)
	})

	// scale down the APs to 2 and delete PVC for scaled down AP; KAA will automatically shrink membership to 2
	k8s.RunKubectl(t, kubectlOptions, "scale", "statefulset", admin, "--replicas=2")
	admin2 := admin + "-2"
	testlib.AwaitServerState(
		t, namespaceName, admin0, admin2, "Disconnected", 60*time.Second)
	k8s.RunKubectl(t, kubectlOptions, "delete", "pvc", "raftlog-"+admin2)

	testlib.AwaitNrReplicasReady(t, namespaceName, admin, 2)
	// restart the current leader to bounce it
	leader := testlib.AwaitDomainLeader(t, namespaceName, admin0, 60*time.Second)
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
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

		options := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":      testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":   testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.sm.hotCopy.journalPath.enabled": "true",
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &options)
		// check that archives are created with external journal directory
		archives, _ := testlib.CheckArchives(t, namespaceName, admin0, "demo", 1, 0)
		for _, archive := range archives {
			require.Equal(t, "/var/opt/nuodb/journal/nuodb/demo", archive.JournalPath)
		}
	})
}

func TestKubernetesRestrictedDatabase(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"admin.initContainers.runInitDisk":                 "false",
			"admin.securityContext.enabledOnContainer":         "true",
			"admin.securityContext.readOnlyRootFilesystem":     "true",
			"admin.securityContext.capabilities.drop[0]":       "CAP_NET_RAW",
			"admin.securityContext.capabilities.drop[1]":       "ALL",
			"admin.securityContext.runAsNonRootGroup":          "true",
			"admin.resources.requests.cpu":                     "256m",
			"admin.resources.requests.memory":                  "512Mi",
			"database.initContainers.runInitDisk":              "false",
			"database.securityContext.enabledOnContainer":      "true",
			"database.securityContext.readOnlyRootFilesystem":  "true",
			"database.securityContext.capabilities.drop[0]":    "CAP_NET_RAW",
			"database.securityContext.capabilities.drop[1]":    "ALL",
			"database.securityContext.runAsNonRootGroup":       "true",
			"database.sm.resources.requests.cpu":               "256m",
			"database.sm.resources.requests.memory":            testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":               "256m",
			"database.te.resources.requests.memory":            testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.sm.hotCopy.jobResources.requests.cpu":    "100m",
			"database.sm.hotCopy.jobResources.requests.memory": "128Mi",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")
	adminStatefulSet := helmChartReleaseName + "-nuodb-cluster0"
	admin0 := adminStatefulSet + "-0"

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
	smPodNameTemplate := fmt.Sprintf("sm-%s", databaseReleaseName)
	tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)

	// Check that root filesystem is read-only for all pods
	for _, pod := range testlib.FindAllPodsInSchema(t, namespaceName) {
		var container string
		if strings.Contains(pod.Name, adminStatefulSet) {
			container = "admin"
		} else if strings.Contains(pod.Name, smPodNameTemplate) || strings.Contains(pod.Name, tePodNameTemplate) {
			container = "engine"
		}

		// Try to remove nuodb executable and check that this fails
		k8sOptions := &k8s.KubectlOptions{Namespace: namespaceName}
		output, err := k8s.RunKubectlAndGetOutputE(t, k8sOptions, "exec", pod.Name, "-c", container, "--", "rm", "-f", "/opt/nuodb/bin/nuodb")
		require.Error(t, err, "Expected error, but found output:\n---\n%s", output)
		require.Contains(t, err.Error(), "rm: cannot remove '/opt/nuodb/bin/nuodb': Read-only file system")
	}

	t.Run("verifyNuoSQL", func(t *testing.T) { verifyNuoSQL(t, namespaceName, admin0, "demo") })
}
