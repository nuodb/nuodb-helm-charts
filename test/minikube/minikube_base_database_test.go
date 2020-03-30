// +build short

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"gotest.tools/assert"
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
	opts := k8s.NewKubectlOptions("", "")
	opts.Namespace = namespaceName
	k8s.RunKubectl(t, opts,
		"exec", adminPod, "--",
		"/opt/nuodb/bin/nuosql",
		"--user", "dba",
		"--password", "secret",
		"demo",
		"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql",
	)
}

func verifyKubernetesAccess(t *testing.T, namespaceName string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	serviceAccountDir := "/var/run/secrets/kubernetes.io/serviceaccount"

	// check namespace matches service account directory
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", serviceAccountDir+"/namespace")
	assert.NilError(t, err, output)
	assert.Equal(t, namespaceName, output)

	// get authorization token
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", serviceAccountDir+"/token")
	assert.NilError(t, err, output)

	curlCmdPrefix := fmt.Sprintf("curl -s --cacert %s -H 'Authorization: Bearer %s' https://kubernetes.default.svc", serviceAccountDir+"/ca.crt", output)

	// check that we can access Pods
	url := "/api/v1/namespaces/" + namespaceName + "/pods"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.NilError(t, err, output)
	assert.Check(t, strings.Contains(output, "\"kind\": \"PodList\""), output)

	// check that we can access this Pod
	url = "/api/v1/namespaces/" + namespaceName + "/pods/" + podName
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.NilError(t, err, output)
	assert.Check(t, strings.Contains(output, "\"kind\": \"Pod\""), output)

	// check that we can access PersistentVolumeClaims
	url = "/api/v1/namespaces/" + namespaceName + "/persistentvolumeclaims"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.Check(t, strings.Contains(output, "\"kind\": \"PersistentVolumeClaimList\""), output)

	// check that we can access Deployments
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/deployments"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.Check(t, strings.Contains(output, "\"kind\": \"DeploymentList\""), output)

	// check that we can access StatefulSets
	url = "/apis/apps/v1/namespaces/" + namespaceName + "/statefulsets"
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.Check(t, strings.Contains(output, "\"kind\": \"StatefulSetList\""), output)

	// check that we can create Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases"
	// when request data is specified without an explicit request method, POST is assumed
	leaseName := strings.ToLower(random.UniqueId())
	extraArgs := fmt.Sprintf(" -H 'Content-Type: application/json' -d  '{\"metadata\": {\"name\": \"%s\"}}'", leaseName)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	assert.Check(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can update Leases
	url = "/apis/coordination.k8s.io/v1/namespaces/" + namespaceName + "/leases/" + leaseName
	// use create response as request payload, which contains the correct resourceVersion (update fails if the resourceVersion does not match)
	extraArgs = fmt.Sprintf(" -X PUT -H 'Content-Type: application/json' -d '%s'", output)
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url+extraArgs)
	assert.Check(t, strings.Contains(output, "\"kind\": \"Lease\""), output)

	// check that we can get Leases
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", curlCmdPrefix+url)
	assert.Check(t, strings.Contains(output, "\"kind\": \"Lease\""), output)
}

func verifyNuoSQL(t *testing.T, namespaceName string, adminPod string, databaseName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--", "bash", "-c",
		fmt.Sprintf("echo \"select * from system.nodes;\" | nuosql %s@localhost --user dba --password secret", databaseName))

	assert.NilError(t, err, output)

	assert.Check(t, strings.Contains(output, "Storage"))
	assert.Check(t, strings.Contains(output, "Transaction"))
}

func verifySecret(t *testing.T, namespaceName string) {
	secret := testlib.GetSecret(t, namespaceName, "demo.nuodb.com")

	_, ok := secret.Data["database-name"]
	assert.Check(t, ok)

	_, ok = secret.Data["database-password"]
	assert.Check(t, ok)

	_, ok = secret.Data["database-username"]
	assert.Check(t, ok)
}

func verifyDBService(t *testing.T, namespaceName string, podName string, serviceName string, ping bool) {

	dBService := testlib.GetService(t, namespaceName, serviceName)
	assert.Equal(t, dBService.Name, serviceName)

	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
}

func verifyPodLabeling(t *testing.T, namespaceName string, adminPod string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", "demo")

	assert.NilError(t, err, output)

	err, objects := testlib.Unmarshal(output)

	for _, obj := range objects {
		val, ok := obj.Labels["cloud"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_CLOUD)

		val, ok = obj.Labels["region"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_REGION)

		val, ok = obj.Labels["zone"]
		assert.Check(t, ok)
		assert.Check(t, val == LABEL_ZONE)
	}

}

func verifyPacketFetch(t *testing.T, namespaceName string, admin0 string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	// verify the container can actually download the file from the internet
	start := time.Now()
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions,
		"exec", admin0, "--",
		"bash", "-c",
		fmt.Sprintf("curl -k %s | tar tzf - ", testlib.IMPORT_ARCHIVE_URL),
	)
	assert.NilError(t, err, "Could not fetch archive")
	elapsed := time.Since(start)
	t.Logf("Fetching package (%s) took %f seconds", testlib.IMPORT_ARCHIVE_URL, elapsed.Seconds())
	t.Log("tar contents: ", output)
}

func verifyEngineAltAddress(t *testing.T, namespaceName string, admin0 string, expectedNrEngines int) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	podNames, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--",
		"bash", "-c",
		"nuocmd get processes | grep -o \"(address=[^/]*\"| cut -f2 -d'='")
	assert.NilError(t, err)
	podNamesSlice := strings.Split(strings.TrimSuffix(podNames, "\n"), "\n")

	actualAltAddresses, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--",
		"bash", "-c",
		"nuocmd get processes | grep -o \"alt-address: [^}]*\" | cut -f2 -d' '")
	assert.NilError(t, err)
	actualAltAddressesSlice := strings.Split(strings.TrimSuffix(actualAltAddresses, "\n"), "\n")

	assert.Assert(t, len(podNamesSlice) == expectedNrEngines, "Expected number of process names don't match")
	assert.Assert(t, len(actualAltAddressesSlice) == expectedNrEngines, "Expected number of process addresses don't match")

	for index, podName := range podNamesSlice {
		pod := k8s.GetPod(t, kubectlOptions, podName)
		assert.Assert(t, pod.Status.PodIP == actualAltAddressesSlice[index], "Expected alt-address doesn't match")
	}
}

func backupDatabase(t *testing.T, namespaceName string, podName string, databaseName string, options *helm.Options) {

	// wait for the backup to both start _and_ complete successfully
	backupJob := fmt.Sprintf("hotcopy-%s-job-initial", databaseName)
	testlib.AwaitPodPhase(t, namespaceName, backupJob, corev1.PodSucceeded, 120*time.Second)

	// verify that the backup has been documented by the Admin layer
	backupset, err := k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions,
		"exec", podName, "--",
		"nuodocker", "get", "current-backup", "--db-name", databaseName,
	)

	assert.NilError(t, err, "Error running: nuodocker get current-backup  ")
	assert.Check(t, backupset != "")
}

func restoreDatabase(t *testing.T, namespaceName string, podName string, databaseOptions *helm.Options) {
	// run the restore chart - which flags the database to restore on next startup
	randomSuffix := strings.ToLower(random.UniqueId())

	restName := fmt.Sprintf("restore-demo-%s", randomSuffix)
	options := &helm.Options{
		SetValues: map[string]string{
			"database.name":       "demo",
			"restore.target":      "demo",
			"restore.source":      ":latest",
			"restore.autoRestart": "true",
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "")
	options.KubectlOptions = kubectlOptions
	options.KubectlOptions.Namespace = namespaceName

	restore := func() {
		testlib.InjectTestVersion(t, options)
		helm.Install(t, options, testlib.RESTORE_HELM_CHART_PATH, restName)
		testlib.AddTeardown(testlib.TEARDOWN_RESTORE, func() { helm.Delete(t, options, restName, true) })

		testlib.AwaitPodPhase(t, namespaceName, "restore-demo-", corev1.PodSucceeded, 120*time.Second)
	}

	testlib.AwaitDatabaseRestart(t, namespaceName, podName, "demo", databaseOptions, restore)
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
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		backupDatabase(t, namespaceName, admin0, "demo", &databaseOptions)
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
			},
		}

		testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

		populateCreateDBData(t, namespaceName, admin0)

		defer testlib.Teardown(testlib.TEARDOWN_BACKUP)
		backupDatabase(t, namespaceName, admin0, "demo", &databaseOptions)
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
		},
	}

	databaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)

	// Generate diagnose in case this test fails
	testlib.AddTeardown(testlib.TEARDOWN_DATABASE, func() { testlib.GetDiagnoseOnTestFailure(t, namespaceName, admin0) })

	opts := k8s.NewKubectlOptions("", "")
	opts.Namespace = namespaceName
	databaseOptions.KubectlOptions = opts

	// wait for the initial backup to complete
	databaseName := "demo"
	backupJob := fmt.Sprintf("hotcopy-%s-job-initial-", databaseName)
	testlib.AwaitPodPhase(t, namespaceName, backupJob, corev1.PodSucceeded, 120*time.Second)

	// populate some data
	k8s.RunKubectl(t, opts,
		"exec", admin0, "--",
		"/opt/nuodb/bin/nuosql",
		"--user", "dba",
		"--password", "secret",
		"demo",
		"--file", "/opt/nuodb/samples/quickstart/sql/create-db.sql",
	)

	// verify that the database contains the populated data
	tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	assert.NilError(t, err, "error running SQL: show schema User")
	assert.Check(t, strings.Contains(tables, "HOCKEY"), "tables returned: ", tables)

	opt := testlib.GetExtractedOptions(&databaseOptions)
	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)
	smPodName := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseChartName, opt.ClusterName, opt.DbName)

	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	testlib.GetAppLog(t, namespaceName, tePodName, "_pre-restart")

	smPodName0 := testlib.GetPodName(t, namespaceName, smPodName)
	testlib.GetAppLog(t, namespaceName, smPodName0, "_pre-restart")

	// restore database
	defer testlib.Teardown(testlib.TEARDOWN_RESTORE)
	restoreDatabase(t, namespaceName, admin0, &databaseOptions)

	// verify that the database does NOT contain the data from AFTER the backup
	tables, err = testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
	assert.NilError(t, err, "error running SQL: show schema User")
	assert.Check(t, strings.Contains(tables, "No tables found in schema "), "Show schema returned: ", tables)
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

		testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"database.autoImport.source":            testlib.IMPORT_ARCHIVE_URL,
				"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"backup.persistence.enabled":            "true",
				"backup.persistence.size":               "1Gi",
			},
		})

		// verify that the database contains the restored data
		tables, err := testlib.RunSQL(t, namespaceName, admin0, "demo", "show schema User")
		assert.NilError(t, err, "error running SQL: show schema User")
		assert.Check(t, strings.Contains(tables, "HOCKEY"))
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
		assert.NilError(t, err, "error running SQL: show schema User")
		assert.Check(t, strings.Contains(tables, "HOCKEY"))
	})
}
