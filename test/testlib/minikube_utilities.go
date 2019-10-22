package testlib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	v1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var RESULT_DIR = "results"

var teardownLists = make(map[string][]func())

/**
 * add a teardown function to the named list - for later execution
 */
func AddTeardown(name string, teardownFunc func()) {
	teardownLists[name] = append(teardownLists[name], teardownFunc)
}

/**
 * Call the stored teardown functions in the named list, in the correct order (last-in-first-out)
 */
func Teardown(name string) {
	list := teardownLists[name]
	delete(teardownLists, name)

	for x := len(list) - 1; x >= 0; x-- {
		list[x]()
	}
}

/**
* Verify all teardownLists have been executed already; and throw an ASSERT if not.
* Can be used to verify correct codng of a test that uses teardown.
*
* NOTE: while the funcs are called in the correct order for each list, there can be
* NO guarantee that the lists are iterated in the correct order.
*
* This function MUST NOT be used as a replacement for calling teardown() at the correct point in the code.
 */
func VerifyTeardown(t *testing.T) {
	remaining := len(teardownLists)

	// make a "best-effort" at releasing all remaining resources
	for _, list := range teardownLists {
		for x := len(list) - 1; x >= 0; x-- {
			list[x]()
		}
	}

	// release all funcs in all lists
	teardownLists = make(map[string][]func())

	assert.Check(t, remaining == 0, "Error - %d teardownLists were left uncleared", remaining)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func arePodConditionsMet(pod *corev1.Pod, condition corev1.PodConditionType,
	status corev1.ConditionStatus) bool {
	for _, cnd := range pod.Status.Conditions {
		if cnd.Type == condition && cnd.Status == status {
			return true
		}
	}

	return false
}

func findAllPodsInSchema(t *testing.T, namespace string) []corev1.Pod {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace
	filter := metav1.ListOptions{}
	return k8s.ListPods(t, options, filter)
}

func Await(t *testing.T, lmbd func() bool, timeout time.Duration) {
	for timeExpired := time.After(timeout); ; {
		select {
		case <-timeExpired:
			t.Log(string(debug.Stack()))
			t.Fatal("function call timed out")
		default:
			if lmbd() {
				return
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func AwaitTillerUp(t *testing.T) {
	Await(t, func() bool {
		for _, pod := range findAllPodsInSchema(t, "kube-system") {
			if strings.Contains(pod.Name, "tiller-deploy") {
				if arePodConditionsMet(&pod, corev1.PodReady, corev1.ConditionTrue) {
					return true
				}
			}
		}
		return false
	}, 30*time.Second)
}

func AwaitNrReplicasScheduled(t *testing.T, namespace string, expectedName string, nrReplicas int) {
	Await(t, func() bool {
		var cnt int
		for _, pod := range findAllPodsInSchema(t, namespace) {
			if strings.Contains(pod.Name, expectedName) {
				if arePodConditionsMet(&pod, corev1.PodScheduled, corev1.ConditionTrue) {
					cnt++
				}
			}
		}

		t.Logf("%d pods SCHEDULED for name '%s'\n", cnt, expectedName)

		return cnt == nrReplicas
	}, 30*time.Second)
}

func AwaitNoPods(t *testing.T, namespace string, expectedName string) {
	Await(t, func() bool {
		var cnt int
		for _, pod := range findAllPodsInSchema(t, namespace) {
			if strings.Contains(pod.Name, expectedName) {
				cnt++
			}
		}
		t.Logf("%d pods still RUNNING for name '%s'\n", cnt, expectedName)
		return cnt == 0
	}, 120*time.Second)
}

func findPod(t *testing.T, namespace string, expectedName string) (*corev1.Pod, error) {
	for _, pod := range findAllPodsInSchema(t, namespace) {
		if strings.Contains(pod.Name, expectedName) {
			return &pod, nil
		}
	}

	return nil, errors.New("did not find any pod matching name")
}

func GetPodName(t *testing.T, namespaceName string, expectedName string) string {
	tePod, err := findPod(t, namespaceName, expectedName)
	assert.NilError(t, err, "No pod found with name ", expectedName)

	return tePod.Name
}

func AwaitPodStatus(t *testing.T, namespace string, podName string, condition corev1.PodConditionType,
	status corev1.ConditionStatus, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	Await(t, func() bool {
		pod := k8s.GetPod(t, options, podName)
		return arePodConditionsMet(pod, condition, status)
	}, timeout)
}

func AwaitPodPhase(t *testing.T, namespace string, podName string, phase corev1.PodPhase, timeout time.Duration) {
	Await(t, func() bool {
		pod, err := findPod(t, namespace, podName)
		assert.NilError(t, err, "awaitPodPhase: could not find pod with name matching ", podName)

		return pod.Status.Phase == phase
	}, timeout)
}

func AwaitAdminPodUp(t *testing.T, namespace string, adminPodName string, timeout time.Duration) {
	AwaitPodStatus(t, namespace, adminPodName, corev1.PodReady, corev1.ConditionTrue, timeout)
}

func AwaitAdminPodUpgraded(t *testing.T, namespace string, adminPodName string, expectedVersion string, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	Await(t, func() bool {
		pod := k8s.GetPod(t, options, adminPodName)

		for _, container := range pod.Spec.Containers {
			t.Logf("Found container (%s) with image: %s", container.Name, container.Image)
			if container.Image == expectedVersion {
				return true
			}
		}

		return false
	}, timeout)
}

func AwaitBalancerTerminated(t *testing.T, namespace string, expectedName string) {
	Await(t, func() bool {
		for _, pod := range findAllPodsInSchema(t, namespace) {
			if strings.Contains(pod.Name, expectedName) {
				if pod.Status.Phase == "Succeeded" {
					t.Logf("Pod (%s) TERMINATED\n", expectedName)
					return true
				}
			}
		}
		return false
	}, 60*time.Second)
}

func VerifyAdminState(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "domain")

	assert.NilError(t, err, "verifyAdminState: running show domain failed")
	assert.Assert(t, strings.Contains(output, "ACTIVE"))
}

func AwaitAdminFullyConnected(t *testing.T, namespace string, podName string, numServers int) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "check", "servers",
		"--check-active", "--check-connected", "--check-leader",
		"--num-servers", strconv.Itoa(numServers),
		"--timeout", "300")
}

func AwaitDatabaseUp(t *testing.T, namespace string, podName string, databaseName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "check", "database",
		"--db-name", databaseName, "--check-running", "--check-liveness", "20",
		"--num-processes", "2",
		"--timeout", "300")
}

func VerifyPolicyInstalled(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "get", "load-balancers")
	assert.NilError(t, err, "VerifyPolicyInstalled: ", podName)
	assert.Assert(t, strings.Contains(output, "LoadBalancerPolicy"))
}

func VerifyLicenseFile(t *testing.T, namespace string, podName string, expectedLicense string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", "/etc/nuodb/nuodb.lic")
	assert.NilError(t, err, "verifyLicenseFile: exec cat nuodb.lic")
	assert.Equal(t, output, expectedLicense)
}

func VerifyCustomFileDoesNotGetMounted(t *testing.T, namespace string, podName string, unexpectedFile string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "ls", "/etc/nuodb/")
	assert.NilError(t, err)
	assert.Assert(t, !strings.Contains(output, unexpectedFile), output)
}

func VerifyLicenseIsCommunity(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "domain")
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(output, "server license: Community"), output)
}

func VerifyLicensingErrorsInLog(t *testing.T, namespace string, podName string, expectError bool) {
	buf, err := ioutil.ReadAll(getAppLogStream(t, namespace, podName))
	assert.NilError(t, err)

	fullLog := string(buf)

	assert.Equal(t, expectError, strings.Contains(fullLog, "Unable to verify configured license"), fullLog)
}

func VerifyCertificateInLog(t *testing.T, namespace string, podName string, expectedLogLine string) {
	buf, err := ioutil.ReadAll(getAppLogStream(t, namespace, podName))
	assert.NilError(t, err)

	fullLog := string(buf)

	assert.Assert(t, strings.Contains(standardizeSpaces(fullLog), expectedLogLine),
		"`%s` not found in:\n %s", expectedLogLine, fullLog)
}

func KillAdminPod(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "delete", "pod", podName)
	assert.NilError(t, err, "killAdminPod: delete pod returned an error")
	assert.Assert(t, strings.Contains(output, "deleted"), "`deleted` not found in %s", output)
}

func KillAdminProcess(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "ps")
	assert.NilError(t, err, "killAdminProcess: exec ps")
	parts := strings.Split(output, "\n")

	var pid string
	for _, part := range parts {
		if strings.Contains(part, "java") {
			pid = strings.Fields(part)[0]
		}
	}
	assert.Assert(t, pid != "", "pid not found in :%s\n", output)

	t.Logf("Killing pid %s in pod %s\n", pid, podName)

	k8s.RunKubectl(t, options, "exec", podName, "--", "kill", pid)
}

func GetService(t *testing.T, namespaceName string, serviceName string) *corev1.Service {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	return k8s.GetService(t, options, serviceName)
}

func PingService(t *testing.T, namespace string, serviceName string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	fullServiceName := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--",
		"ping", fullServiceName, "-c", "1",
	)
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(output, "1 received"))
}

func shouldPrintToStdout() bool {
	_, exists := os.LookupEnv("NUODB_PRINT_TO_STDOUT")
	return exists
}

func GetAppLog(t *testing.T, namespace string, podName string) {
	dirPath := filepath.Join("..", "..", RESULT_DIR, namespace)
	filePath := filepath.Join(dirPath, podName)

	_ = os.MkdirAll(dirPath, 0700)

	f, err := os.Create(filePath)
	assert.NilError(t, err)
	defer f.Close()

	// it is hard to recover this in Travis from the filesystem, without access to a AWS
	// print it to stdout instead
	var multiWriter io.Writer
	if t.Failed() && shouldPrintToStdout() {
		multiWriter = io.MultiWriter(f, os.Stdout)
	} else {
		multiWriter = io.MultiWriter(f)
	}

	_, err = io.Copy(multiWriter, getAppLogStream(t, namespace, podName))
	assert.NilError(t, err)
}

func getAppLogStream(t *testing.T, namespace string, podName string) io.ReadCloser {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	client, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	require.NoError(t, err)

	podLogOpts := corev1.PodLogOptions{}

	reader, err := client.CoreV1().Pods(options.Namespace).GetLogs(podName, &podLogOpts).Stream()
	assert.NilError(t, err)

	return reader
}

func GetSecret(t *testing.T, namespace string, secretName string) *corev1.Secret {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	return k8s.GetSecret(t, options, secretName)
}

func GetDaemonSet(t *testing.T, namespace string, daemonSetName string) *v1.DaemonSet {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "daemonset", daemonSetName,
		"-o", "yaml")
	assert.NilError(t, err, "getDaemonSet: kubectl get daemonSet")

	var object v1.DaemonSet
	helm.UnmarshalK8SYaml(t, output, &object)

	return &object
}

func decodeBase64(t *testing.T, input []byte) string {
	str, err := base64.StdEncoding.DecodeString(string(input))
	assert.NilError(t, err, "decodeBase64: base64.decodeString")
	return string(str)
}

func DeleteDatabase(t *testing.T, namespace string, dbName string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "delete", "database", "--db-name", dbName)
}

func DeletePod(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	k8s.RunKubectl(t, options, "delete", podName)
}

func RunSQL(t *testing.T, namespace string, podName string, databaseName string, sql string) (result string, err error) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespace

	// secrets := getSecret(t, namespace, databaseName)

	result, err = k8s.RunKubectlAndGetOutputE(t, options,
		"exec", podName, "--",
		"bash", "-c",
		fmt.Sprintf("echo \"%s;\" | /opt/nuodb/bin/nuosql --user dba --password secret %s", sql, databaseName),
	)

	assert.NilError(t, err, "runSQL: error trying to run ", sql)

	return result, err
}
