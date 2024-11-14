//go:build short
// +build short

package minikube

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func isMinikube(t *testing.T) bool {
	kubectlOptions := k8s.NewKubectlOptions("", "", "")
	ret, err := k8s.IsMinikubeE(t, kubectlOptions)
	require.NoError(t, err)
	return ret
}

func isDockerDesktop(t *testing.T) bool {
	kubectlOptions := k8s.NewKubectlOptions("", "", "")
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "config", "current-context")
	require.NoError(t, err)
	return strings.TrimSpace(output) == "docker-desktop"
}

func kubeHostSsh(t *testing.T, args ...string) string {
	var command string
	if isMinikube(t) {
		// Use `minikube ssh "sudo ..."`
		command = "minikube"
		args = append([]string{"sudo"}, args...)
		args = append([]string{"ssh"}, shellquote.Join(args...))
	} else {
		// Create privileged container sharing namespace with pid 1 on host
		command = "docker"
		args = append(
			[]string{"run", "--rm", "--privileged", "--pid=host", "alpine:edge", "nsenter", "-t", "1", "-m", "-u", "-n", "-i"},
			args...)
	}
	t.Logf("Running command on K8s host: %v", args)
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		var msg string
		if exitErr, ok := err.(*exec.ExitError); ok {
			msg = fmt.Sprintf("Command failed: error=%s, stderr=[%s], stdout=[%s]", exitErr.Error(), string(exitErr.Stderr), output)
		} else {
			msg = fmt.Sprintf("Command failed: error=%s, stdout=[%s]", err.Error(), output)
		}
		require.Fail(t, msg)
	}
	ret := string(output)
	t.Logf(ret)
	return ret
}

func rolloutRestart(t *testing.T, kubectlOptions *k8s.KubectlOptions, workload string) {
	k8s.RunKubectl(t, kubectlOptions, "rollout", "restart", workload)
	k8s.RunKubectl(t, kubectlOptions, "rollout", "status", workload, "--timeout=10s")
}

// prepareCsiDriver mounts a loopback filesystem as the data directory for the
// CSI hostpath driver so that when fsfreeze is invoked within the backup hooks
// container on the SM archive, it does not disrupt etcd and other Kubernetes
// services, which are normally using the same filesystem as the data directory.
func prepareCsiDriver(t *testing.T) {
	// Get pod for CSI hostpath driver
	kubectlOptions := k8s.NewKubectlOptions("", "", "")
	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=csi-hostpathplugin",
	}
	csiDriverPods := k8s.ListPods(t, kubectlOptions, listOptions)
	require.Len(t, csiDriverPods, 1)
	// Find workload for pod
	var workload string
	for _, ownerRef := range csiDriverPods[0].OwnerReferences {
		if ownerRef.Controller != nil && *ownerRef.Controller {
			workload = fmt.Sprintf("%s/%s", strings.ToLower(ownerRef.Kind), ownerRef.Name)
		}
	}
	require.NotEmpty(t, workload, "Did not find workload for pod %s", csiDriverPods[0])
	// Find hostpath directory for csi-data-dir
	var csiDataDir string
	for _, volume := range csiDriverPods[0].Spec.Volumes {
		if volume.Name == "csi-data-dir" && volume.HostPath != nil {
			csiDataDir = volume.HostPath.Path
		}
	}
	require.NotEmpty(t, csiDataDir, "Did not find csi-data-dir for pod %s", csiDriverPods[0])

	// Set namespace for kubectl options so that we can use it to restart
	// workload managing CSI driver
	kubectlOptions.Namespace = csiDriverPods[0].Namespace

	// Teardown functions are executed in reverse order of being registered.
	// Register restart of workload first so that it is executed after
	// umounting the loopback filesystem.
	testlib.AddTeardown(testlib.TEARDOWN_CSIDRIVER_FS, func() {
		rolloutRestart(t, kubectlOptions, workload)
	})

	// Create loopback device and filesystem
	kubeHostSsh(t, "dd", "if=/dev/zero", "of=/tmp/csi-vol.img", "bs=64M", "count=1")
	loopbackDevice := strings.TrimSpace(kubeHostSsh(t, "losetup", "--show", "-fP", "/tmp/csi-vol.img"))
	testlib.AddTeardown(testlib.TEARDOWN_CSIDRIVER_FS, func() {
		kubeHostSsh(t, "losetup", "-d", loopbackDevice)
	})
	kubeHostSsh(t, "mkfs.ext4", loopbackDevice)
	// Mount loopback filesystem on csi-data-dir
	kubeHostSsh(t, "mount", loopbackDevice, csiDataDir)
	testlib.AddTeardown(testlib.TEARDOWN_CSIDRIVER_FS, func() {
		kubeHostSsh(t, "umount", csiDataDir)
	})

	// Restart workload for CSI driver so that it uses loopback filesystem
	// as data directory for volumes
	rolloutRestart(t, kubectlOptions, workload)
}

func TestFsFreezeBackupHook(t *testing.T) {
	if !isMinikube(t) && !isDockerDesktop(t) {
		t.Skip("Can only run test on Minikube or Docker Desktop")
	}

	defer testlib.VerifyTeardown(t)

	// Prepare CSI driver to enable fsfreeze
	defer testlib.Teardown(testlib.TEARDOWN_CSIDRIVER_FS)
	prepareCsiDriver(t)

	// Create admin release
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	// Create database release with CSI driver storage class and backup hooks enabled
	admin := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := fmt.Sprintf("%s-0", admin)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	sourceDatabaseChartName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":                          testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.persistence.storageClass": testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.enabled":                  "true",
			"database.sm.noHotCopy.replicas":                             "1",
			"database.sm.hotCopy.replicas":                               "0",
			"database.backupHooks.enabled":                               "true",
			"database.backupHooks.freezeMode":                            "fsfreeze",
		},
	})

	// Invoke backup hooks on SM pod
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	smPod := fmt.Sprintf("sm-%s-nuodb-cluster0-demo-0", sourceDatabaseChartName)
	backupId := "123abc"

	// Collect backup hooks logs on test failure
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetAppLog(t, namespaceName, smPod, "_hooks", &corev1.PodLogOptions{Container: "backup-hooks"})
	})

	t.Run("verifyWritesFrozen", func(t *testing.T) {
		// Freeze writes to archive and defer invocation of unfreeze
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		defer testlib.InvokeBackupHook(t, namespaceName, smPod, "post-backup/"+backupId)

		// Try to create a file in archive directory in the background
		k8s.RunKubectl(t, kubectlOptions, "exec", smPod, "-c", "engine", "--",
			"sh", "-c", "nohup touch /var/opt/nuodb/archive/test-xyz &")

		// Wait some time and check that file has not been created
		time.Sleep(3 * time.Second)
		output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--",
			"ls", "/var/opt/nuodb/archive")
		require.NoError(t, err, output)
		require.NotContains(t, output, "test-xyz")
	})

	// Wait short amount of time for test file to be created
	testlib.Await(t, func() bool {
		output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", smPod, "-c", "engine", "--",
			"ls", "/var/opt/nuodb/archive")
		require.NoError(t, err, output)
		return strings.Contains(output, "test-xyz")
	}, 50*time.Millisecond)

	// Check that backup hook sidecar logged expected messages
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "logs", smPod, "-c", "backup-hooks")
	require.NoError(t, err, output)
	require.Contains(t, output, "Freezing writes to archive volume")
	require.Contains(t, output, "Unfreezing writes to archive volume")

	t.Run("negativeTests", func(t *testing.T) {
		// Negative test: use a bogus resource path
		response := testlib.GetBackupHookResponse(t, namespaceName, smPod, "liquidate-archive")
		require.False(t, response.Success)
		require.Equal(t, "No handler found for path /liquidate-archive", response.Message)

		// Negative test: use a resource path with too many components
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "pre-backup/123/456")
		require.False(t, response.Success)
		require.Equal(t, "2 parameter(s) expected but 3 supplied in request /pre-backup/123/456, including payload and query parameters", response.Message)

		// Negative test: use a resource path with not enough components
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "pre-backup")
		require.False(t, response.Success)
		require.Equal(t, "2 parameter(s) expected but 1 supplied in request /pre-backup, including payload and query parameters", response.Message)

		// Negative test: supply an empty backup ID
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "pre-backup/")
		require.False(t, response.Success)
		require.Equal(t, "Backup ID not specified", response.Message)

		// Negative test: invoke post-backup hook without backup in progress
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "post-backup/"+backupId)
		require.False(t, response.Success)
		require.Equal(t, "Unexpected backup ID: current=None, supplied="+backupId, response.Message)

		// Freeze writes to archive using direct invocation instead of HTTP
		k8s.RunKubectl(t, kubectlOptions, "exec", smPod, "-c", "backup-hooks", "--",
			"python", "/backup_hooks.py", "pre-hook", "--backup-id", backupId)
		// Defer unfreeze using direct invocation
		defer k8s.RunKubectl(t, kubectlOptions, "exec", smPod, "-c", "backup-hooks", "--",
			"python", "/backup_hooks.py", "post-hook", "--backup-id", backupId)

		// Negative test: invoke post-backup hook with incorrect backup ID
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "post-backup/bogus")
		require.False(t, response.Success)
		require.Equal(t, "Unexpected backup ID: current="+backupId+", supplied=bogus", response.Message)

		// Negative test: invoke pre-backup hook while backup is in progress
		response = testlib.GetBackupHookResponse(t, namespaceName, smPod, "pre-backup/"+backupId)
		require.False(t, response.Success)
		require.Equal(t, "Backup ID file /mnt/archive/nuodb/demo/backup.txt exists. Execute post-backup hook to complete current backup.", response.Message)
	})
}

func TestHotSnapBackupHook(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":                         testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory":                      testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.persistence.storageClass":                          testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.persistence.storageClass": testlib.SNAPSHOTABLE_STORAGE_CLASS,
			"database.sm.noHotCopy.journalPath.enabled":                  "true",
			"database.sm.noHotCopy.replicas":                             "1",
			"database.sm.hotCopy.replicas":                               "0",
			"database.backupHooks.enabled":                               "true",
			"database.backupHooks.freezeMode":                            "hotsnap",
		},
	}

	testlib.InjectTestValues(t, options)
	if options.SetValues["nuodb.image.tag"] == "" {
		// Specify the NuoDB version explicitly
		options.SetValues["nuodb.image.tag"] = "6.0.2"
	} else {
		// NuoDB version is injected by the test framework; skip the test if
		// HotSnap is not supported
		testlib.SkipTestOnNuoDBVersionCondition(t, "< 6.0.2")
	}

	// Create admin release
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &helm.Options{}, 1, "")

	// Create database release with CSI driver storage class and backup hooks enabled
	admin := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := fmt.Sprintf("%s-0", admin)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	sourceDatabaseChartName := testlib.StartDatabase(t, namespaceName, admin0, options)

	// Invoke backup hooks on SM pod
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	smPod := fmt.Sprintf("sm-%s-nuodb-cluster0-demo-0", sourceDatabaseChartName)
	backupId := "1234abc"

	// Collect backup hooks logs on test failure
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_DATABASE, t, func() {
		testlib.GetAppLog(t, namespaceName, smPod, "_hooks", &corev1.PodLogOptions{Container: "backup-hooks"})
	})

	// Collect backup hooks logs before it got restarted
	go testlib.GetAppLog(t, namespaceName, smPod, "_hooks_pre-restart",
		&corev1.PodLogOptions{Container: "backup-hooks", Follow: true})

	t.Run("verifyPauseArchiving", func(t *testing.T) {
		// Pause archiving and defer resume invocation
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		defer testlib.InvokeBackupHook(t, namespaceName, smPod, "post-backup/"+backupId)

		// Verify that NuoDB archiving is paused
		require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, smPod,
			fmt.Sprintf("pause archiving (pause id %s)", backupId), &corev1.PodLogOptions{}),
			"Expected NuoDB to pause archiving")
	})

	// Check that backup hook sidecar logged expected messages
	output, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "logs", smPod, "-c", "backup-hooks")
	require.NoError(t, err, output)
	require.Contains(t, output, "Pausing archiving for nuodb process")
	require.Contains(t, output, "Resuming archiving for nuodb process")
	require.Equal(t, 1, testlib.GetStringOccurrenceInLog(t, namespaceName, smPod,
		fmt.Sprintf("resume archiving (pause id %s)", backupId), &corev1.PodLogOptions{}),
		"Expected NuoDB to resume archiving")

	t.Run("verifyHooksTimeout", func(t *testing.T) {
		// Pause archiving with a smaller timeout to speed up the test but do
		// not invoke the resume hook
		backupId := strings.ReplaceAll(t.Name(), "/", "-")
		timeout := "5"
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId, "timeout", timeout)

		// Wait for the automatic backup canceling and metadata files removal
		testlib.Await(t, func() bool {
			return testlib.GetStringOccurrenceInLog(t, namespaceName, smPod,
				fmt.Sprintf("Canceling backup with ID %s. Timeout after %ss", backupId, timeout),
				&corev1.PodLogOptions{Container: "backup-hooks"}) > 0
		}, 10*time.Second)
		k8s.RunKubectl(t, kubectlOptions, "exec", smPod, "-c", "engine", "--",
			"sh", "-c", fmt.Sprintf(`
			i=0; 
			while [ -f "/var/opt/nuodb/archive/backup.txt" ]; do 
			if [ $i -ge %s ]; then 
				echo "ERROR: Backup metadata file not removed: Timeout after ${i}s"
				exit 1
			fi
			sleep 1
			((i=i+1))
			done`, timeout))

		// Verify that the archiving can be paused
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		defer testlib.InvokeBackupHook(t, namespaceName, smPod, "post-backup/"+backupId)
	})

	t.Run("verifyMetadataAfterEngineRestart", func(t *testing.T) {
		backupId := strings.ReplaceAll(t.Name(), "/", "")
		// Pause archiving and restart the SM
		pod := testlib.GetPod(t, namespaceName, smPod)
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		testlib.DeletePod(t, namespaceName, "pod/"+smPod)
		testlib.AwaitPodObjectRecreated(t, namespaceName, pod, 30*time.Second)
		testlib.AwaitPodUp(t, namespaceName, smPod, 120*time.Second)

		// Verify that resume archiving will fail as the previous pause
		// archiving was invalidated by the SM restart
		response := testlib.GetBackupHookResponse(t, namespaceName, smPod, "post-backup/"+backupId)
		require.False(t, response.Success)
		require.Contains(t, response.Message, "Archiving not paused")

		// Pause archiving and restart the SM again
		pod = testlib.GetPod(t, namespaceName, smPod)
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		testlib.DeletePod(t, namespaceName, "pod/"+smPod)
		testlib.AwaitPodObjectRecreated(t, namespaceName, pod, 30*time.Second)
		testlib.AwaitPodUp(t, namespaceName, smPod, 120*time.Second)

		// Verify that pause archiving will succeed as the previous operation
		// was invalidated by the SM restart
		testlib.InvokeBackupHook(t, namespaceName, smPod, "pre-backup/"+backupId)
		defer testlib.InvokeBackupHook(t, namespaceName, smPod, "post-backup/"+backupId)

		// Verify that SM restart was detected
		require.Equal(t, 1, testlib.GetRegexOccurrenceInLog(t, namespaceName, smPod,
			fmt.Sprintf("Unexpected start ID: .* SM process restarted while executing backup ID %s", backupId),
			&corev1.PodLogOptions{Container: "backup-hooks"}),
			"Expected to detect SM startId change")
	})
}
