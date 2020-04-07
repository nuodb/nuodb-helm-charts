// +build enterprise

package minikube

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"gotest.tools/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func verifyProcessLabels(t *testing.T, namespaceName string, adminPod string) (archiveVolumeClaims map[string]string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "processes", "--db-name", "demo")
	assert.NilError(t, err, output)

	err, objects := testlib.Unmarshal(output)
	assert.NilError(t, err, output)

	archiveVolumeClaims = make(map[string]string)
	for _, obj := range objects {
		podName, ok := obj.Labels["pod-name"]
		assert.Check(t, ok)
		// check that Pod exists
		pod := k8s.GetPod(t, options, podName)

		containerId, ok := obj.Labels["container-id"]
		assert.Check(t, ok)
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
			assert.NilError(t, err)
			archiveVolumeClaims[claimName] = strings.TrimSpace(output)
		} else {
			assert.Equal(t, "TE", obj.Type, "archive-pvc label should only be absent for TEs")
		}
	}
	return archiveVolumeClaims
}

type archive struct {
	ArchiveId string `json:"archiveId"`
	DbName    string `json:"dbName"`
}

func unmarshalArchives(s string) (err error, processes []archive) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj archive
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, processes
		}

		if err != nil {
			return
		}

		processes = append(processes, obj)
	}
}

func checkArchives(t *testing.T, namespaceName string, adminPod string, numExpected int, numExpectedRemoved int) (archives []archive, removedArchives []archive) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	// check archives
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", "demo")
	assert.NilError(t, err, output)

	err, archives = unmarshalArchives(output)
	assert.NilError(t, err)
	assert.Equal(t, numExpected, len(archives), output)

	// check removed archives
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", adminPod, "--",
		"nuocmd", "--show-json", "get", "archives", "--db-name", "demo", "--removed")
	assert.NilError(t, err, output)

	err, removedArchives = unmarshalArchives(output)
	assert.NilError(t, err)
	assert.Equal(t, numExpectedRemoved, len(removedArchives), output)
	return
}

func TestDomainResync(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			// TODO: delete image overrides before pushing
			"nuodb.image.repository":   "nuodb",
			"nuodb.image.tag":          "latest",
			"nuodb.image.pullPolicy":   "IfNotPresent",
			"busybox.image.pullPolicy": "IfNotPresent",
		},
	}, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

	testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    "0.25",
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "0.25",
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			// TODO: delete image overrides before pushing
			"nuodb.image.repository":   "nuodb",
			"nuodb.image.tag":          "latest",
			"nuodb.image.pullPolicy":   "IfNotPresent",
			"busybox.image.pullPolicy": "IfNotPresent",
		},
	})

	originalArchiveVolumeClaims := verifyProcessLabels(t, namespaceName, admin0)
	assert.Equal(t, 1, len(originalArchiveVolumeClaims))
	originalArchiveId := ""
	for _, archiveId := range originalArchiveVolumeClaims {
		originalArchiveId = archiveId
	}
	assert.Assert(t, originalArchiveId != "")

	// update replica count
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "statefulset", "-o", "jsonpath={.items[*].metadata.name}")
	assert.NilError(t, err, output)

	statefulSets := strings.Split(output, " ")
	assert.Equal(t, 3, len(statefulSets), "Expected 3 StatefulSets: Admin, SM, and hotcopy SM")

	// by default the hotcopy SM replica count is 1 and regular SM count is 0
	// scale regular SM replica count up to 1
	smStatefulSet := ""
	for _, statefulSet := range statefulSets {
		if strings.HasPrefix(statefulSet, "sm-") && !strings.Contains(statefulSet, "hotcopy") {
			k8s.RunKubectl(t, options, "scale", "statefulset", statefulSet, "--replicas=1")
			smStatefulSet = statefulSet
		}
	}
	assert.Assert(t, smStatefulSet != "")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 3)
	checkArchives(t, namespaceName, admin0, 2, 0)

	// scale hotcopy SM replica count down to 0
	hotCopySmStatefulSet := ""
	for _, statefulSet := range statefulSets {
		if strings.Contains(statefulSet, "hotcopy") {
			k8s.RunKubectl(t, options, "scale", "statefulset", statefulSet, "--replicas=0")
			hotCopySmStatefulSet = statefulSet
		}
	}
	assert.Assert(t, hotCopySmStatefulSet != "")
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
	// check that archive ID generated by hotcopy SM was removed
	_, removedArchives := checkArchives(t, namespaceName, admin0, 1, 1)
	assert.Equal(t, originalArchiveId, removedArchives[0].ArchiveId)

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
