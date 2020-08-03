// +build enterprise

package minikube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"



	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getStatefulSets(t *testing.T, namespaceName string) *appsv1.StatefulSetList {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	assert.NoError(t, err)

	statefulSets, err := clientset.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	assert.NoError(t, err)

	return statefulSets
}

func verifyProcessLabels(t *testing.T, namespaceName string, adminPod string) (archiveVolumeClaims map[string]int) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", "demo")
	assert.NoError(t, err, output)

	err, objects := testlib.Unmarshal(output)
	assert.NoError(t, err, output)

	archiveVolumeClaims = make(map[string]int)
	for _, obj := range objects {
		podName, ok := obj.Labels["pod-name"]
		assert.True(t, ok)
		// check that Pod exists
		pod := k8s.GetPod(t, options, podName)

		containerId, ok := obj.Labels["container-id"]
		assert.True(t, ok)
		// check that Pod has container ID
		for _, containerStatus := range pod.Status.ContainerStatuses {
			assert.Equal(t, "docker://"+containerId, containerStatus.ContainerID)
		}

		claimName, ok := obj.Labels["archive-pvc"]
		if ok {
			assert.Equal(t, "SM", obj.Type, "archive-pvc label should only be present for SMs")
			// check that PVC exists
			k8s.RunKubectl(t, options, "get", "pvc", claimName)
			// add mapping of PVC to archive ID
			output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
				"nuocmd", "get", "value", "--key", "archiveVolumeClaims/"+claimName)
			assert.NoError(t, err)

			archiveId, err := strconv.Atoi(strings.TrimSpace(output))
			assert.NoError(t, err)
			archiveVolumeClaims[claimName] = archiveId
		} else {
			assert.Equal(t, "TE", obj.Type, "archive-pvc label should only be absent for TEs")
		}
	}
	return archiveVolumeClaims
}

func checkArchives(t *testing.T, namespaceName string, adminPod string, numExpected int, numExpectedRemoved int) (archives []testlib.NuoDBArchive, removedArchives []testlib.NuoDBArchive) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	// check archives
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", "demo")
	assert.NoError(t, err, output)

	err, archives = testlib.UnmarshalArchives(output)
	assert.NoError(t, err)
	assert.Equal(t, numExpected, len(archives), output)

	// check removed archives
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", "demo", "--removed")
	assert.NoError(t, err, output)

	err, removedArchives = testlib.UnmarshalArchives(output)
	assert.NoError(t, err)
	assert.Equal(t, numExpectedRemoved, len(removedArchives), output)
	return
}

func checkInitialMembership(t assert.TestingT, configJson string, expectedSize int) {
	type initialMembershipEntry struct {
		Transport string `json:"transport"`
		Version   string `json:"version"`
	}
	var adminConfig struct {
		InitialMembership map[string]initialMembershipEntry `json:"initialMembership"`
	}
	dec := json.NewDecoder(strings.NewReader(configJson))
	err := dec.Decode(&adminConfig)
	if err != io.EOF {
		assert.NoError(t, err, "Unable to deserialize admin config")
	}
	assert.Equal(t, expectedSize, len(adminConfig.InitialMembership))
}

func TestReprovisionAdmin0(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"admin.replicas":         "2",
			"admin.bootstrapServers": "2",
		},
	}, 2, "")

	adminStatefulSet := helmChartReleaseName + "-nuodb-cluster0"
	admin0 := adminStatefulSet + "-0"
	admin1 := adminStatefulSet + "-1"

	// check initial membership on admin-0
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "--",
		"nuocmd", "--show-json", "get", "server-config", "--this-server")
	assert.NoError(t, err, output)
	checkInitialMembership(t, output, 2)

	// check initial membership on admin-1
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", admin1, "--",
		"nuocmd", "--show-json", "get", "server-config", "--this-server")
	assert.NoError(t, err, output)
	checkInitialMembership(t, output, 2)

	// store a value in the KV store via admin-0
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "set", "value", "--key", "testKey", "--value", "0", "--unconditional")

	// save the original Pod object
	originalPod := k8s.GetPod(t, options, admin0)

	// delete Raft data and Pod for admin-0
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"bash", "-c", "rm $NUODB_VARDIR/raftlog")
	k8s.RunKubectl(t, options, "delete", "pod", admin0)

	// wait until the Pod is rescheduled
	testlib.AwaitPodObjectRecreated(t, namespaceName, originalPod, 300*time.Second)
	testlib.AwaitPodUp(t, namespaceName, admin0, 300*time.Second)

	// make sure admin0 rejoins
	k8s.RunKubectl(t, options, "exec", admin1, "--",
		"nuocmd", "check", "servers", "--check-connected", "--num-servers", "2", "--check-leader", "--timeout", "300")
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "check", "servers", "--check-connected", "--num-servers", "2", "--check-leader", "--timeout", "300")

	// conditionally update value in the KV store via admin-0; if admin-0
	// rejoined with admin-1 rather than bootstrapping a new domain, then it
	// should have the current value
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "set", "value", "--key", "testKey", "--value", "1", "--expected-value", "0")

	// conditionally update value in the KV store via admin-1
	k8s.RunKubectl(t, options, "exec", admin1, "--",
		"nuocmd", "set", "value", "--key", "testKey", "--value", "2", "--expected-value", "1")
}

func TestAdminScaleDown(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"admin.replicas": "2",
		},
	}, 2, "")

	adminStatefulSet := helmChartReleaseName + "-nuodb-cluster0"
	admin0 := adminStatefulSet + "-0"
	admin1 := adminStatefulSet + "-1"

	// scale down Admin StatefulSet
	options := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.RunKubectl(t, options, "scale", "statefulset", adminStatefulSet, "--replicas=1")

	// wait for scaled-down Admin to show as "Disconnected"
	testlib.Await(t, func() bool {
		output, _ := k8s.RunKubectlAndGetOutputE(t, options, "exec", admin0, "--",
			"nuocmd", "show", "domain", "--server-format", "{id} {connected_state}")
		return strings.Contains(output, admin1+" Disconnected")
	}, 300*time.Second)

	// wait for scaled-down Admin Pod to be deleted
	testlib.AwaitNoPods(t, namespaceName, admin1)

	// commit a Raft command to confirm that remaining Admin has consensus
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "set", "value", "--key", "testKey", "--value", "testValue", "--unconditional")

	// admin1 is still in membership, though it is excluded from consensus;
	// delete PVC to cause it to be completely removed from the membership;
	// this should allow the Admin health-check to succeed
	k8s.RunKubectl(t, options, "delete", "pvc", "raftlog-"+admin1)
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "check", "servers", "--check-connected", "--num-servers", "1", "--check-leader", "--timeout", "300")

	// scale up Admin StatefulSet and make sure admin1 rejoins
	k8s.RunKubectl(t, options, "scale", "statefulset", adminStatefulSet, "--replicas=2")
	k8s.RunKubectl(t, options, "exec", admin0, "--",
		"nuocmd", "check", "servers", "--check-connected", "--num-servers", "2", "--check-leader", "--timeout", "300")
	k8s.RunKubectl(t, options, "exec", admin1, "--",
		"nuocmd", "check", "servers", "--check-connected", "--num-servers", "2", "--check-leader", "--timeout", "300")
}

func TestDomainResync(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    "0.25",
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "0.25",
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	originalArchiveVolumeClaims := verifyProcessLabels(t, namespaceName, admin0)
	assert.Equal(t, 1, len(originalArchiveVolumeClaims))
	originalArchiveId := -1
	for _, archiveId := range originalArchiveVolumeClaims {
		originalArchiveId = archiveId
	}
	assert.True(t, originalArchiveId != -1)

	// update replica count
	options := k8s.NewKubectlOptions("", "", namespaceName)

	statefulSets := getStatefulSets(t, namespaceName).Items
	assert.Equal(t, 3, len(statefulSets), "Expected 3 StatefulSets: Admin, SM, and hotcopy SM")

	// by default the hotcopy SM replica count is 1 and regular SM count is 0
	// scale regular SM replica count up to 1
	smStatefulSet := ""
	for _, statefulSet := range statefulSets {
		name := statefulSet.Name
		if strings.HasPrefix(name, "sm-") && !strings.Contains(name, "hotcopy") {
			k8s.RunKubectl(t, options, "scale", "statefulset", name, "--replicas=1")
			smStatefulSet = name
		}
	}
	assert.True(t, smStatefulSet != "")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)
	checkArchives(t, namespaceName, admin0, 2, 0)

	// scale hotcopy SM replica count down to 0
	hotCopySmStatefulSet := ""
	for _, statefulSet := range statefulSets {
		name := statefulSet.Name
		if strings.Contains(name, "hotcopy") {
			k8s.RunKubectl(t, options, "scale", "statefulset", name, "--replicas=0")
			hotCopySmStatefulSet = name
		}
	}
	assert.True(t, hotCopySmStatefulSet != "")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
	// check that archive ID generated by hotcopy SM was removed
	_, removedArchives := checkArchives(t, namespaceName, admin0, 1, 1)
	assert.Equal(t, originalArchiveId, removedArchives[0].Id)

	// scale hotcopy SM replica count back up to 1; the removed archive ID should be resurrected
	k8s.RunKubectl(t, options, "scale", "statefulset", hotCopySmStatefulSet, "--replicas=1")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)
	checkArchives(t, namespaceName, admin0, 2, 0)

	// scale hotcopy SM replica count back down to 0
	k8s.RunKubectl(t, options, "scale", "statefulset", hotCopySmStatefulSet, "--replicas=0")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
	checkArchives(t, namespaceName, admin0, 1, 1)

	// explicitly delete the scaled-down PVC and make sure the archive ID is purged
	for claimName, _ := range originalArchiveVolumeClaims {
		k8s.RunKubectl(t, options, "delete", "pvc", claimName)
	}
	testlib.Await(t, func() bool {
		checkArchives(t, namespaceName, admin0, 1, 0)
		return true
	}, 300*time.Second)
}

func TestNuoDBKubeDiagnostics(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	databaseHelmChartReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	config := testlib.GetNuoDBK8sConfigDump(t, namespaceName, admin0)

	assert.True(t, func() bool { _, ok := config.Pods[admin0]; return ok }())

	tePodNameTemplate := fmt.Sprintf("te-%s-nuodb-%s-%s", databaseHelmChartReleaseName, "cluster0", "demo")
	tePodName := testlib.GetPodName(t, namespaceName, tePodNameTemplate)
	assert.True(t, func() bool { _, ok := config.Pods[tePodName]; return ok }())

	smPodTemplate := fmt.Sprintf("sm-%s-nuodb-%s-%s", databaseHelmChartReleaseName, "cluster0", "demo")
	smPodName := testlib.GetPodName(t, namespaceName, smPodTemplate)
	assert.True(t, func() bool { _, ok := config.Pods[smPodName]; return ok }())

}
