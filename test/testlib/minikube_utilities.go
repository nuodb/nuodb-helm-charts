package testlib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"

	v1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/** Lists of the teardown and diagnostic teardown funcs */
var teardownLists = make(map[string][]func())
var diagnosticTeardownLists = make(map[string][]func())

/** Exported var - initialised from the EnvVar, but can be reset in code if desired */
var AlwaysRunDiagnosticTeardowns = strings.EqualFold(os.Getenv(ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS), "true")

/**
 * add a teardown function to the named list - for deferred execution.
 *
 * The teardown functions are called in reverse order of insertion, by a call to Teardown(name).
 *
 * The typical idiom is:
 * <pre>
 *   testlib.AddTeardown("DATABASE", func() { ...})
 *   // possibly more testlib.AddTeardown("DATABASE", func() { ... })
 *   defer testlib.Teardown("DATABASE")
 * <pre>
 */
func AddTeardown(name string, teardownFunc func()) {
	teardownLists[name] = append(teardownLists[name], teardownFunc)
}

/**
 * add a diagnostic teardown func to be called before any other teardowns in the named list - to aid diagnostics/debugging.
 * This allows a diagnostic teardown to do such things as:
 * <ul>
 *   <li>Generate logging and debug information immediately prior to resurce teardown
 *   <li>call time.Sleep() to allow inspection and/or debugging of the exit state before teardown.
 * <ul>
 *
 * NOTE: it is generally undesirable to add multiple diagnostic teardowns that sleep - so it would usually be best to
 * add any Sleep() debug teardown to the innermost teardown list.
 * Nonetheless, there are use-cases where multiple Sleep() teardowns are useful - to allow inspecting different
 * intermediate states.
 */
func AddDiagnosticTeardown(name string, condition interface{}, teardownFunc func()) {
	tdfunc := func() {
		shouldIdoIt := AlwaysRunDiagnosticTeardowns

		if !shouldIdoIt {
			switch c := condition.(type) {
			case *testing.T:
				shouldIdoIt = c.Failed()

			case func() bool:
				shouldIdoIt = c()

			case bool:
				shouldIdoIt = c

			default:
				shouldIdoIt = c != nil
			}
		}

		if shouldIdoIt {
			teardownFunc()
		}
	}

	diagnosticTeardownLists[name] = append(diagnosticTeardownLists[name], tdfunc)
}

/**
 * Call the stored teardown functions in the named list, in the correct order (last-in-first-out)
 *
 * NOTE: Any DIAGNOSTIC teardowns - those added with AddDiagnosticTeardown() for this name - are called BEFORE any other teardowns for this name.
 *
 * The typical use of Teardown is with a deferred call:
 * defer testlib.Teardown("SOME NAME")
 * See: testlib.AddTeardown(); testlib.AddDiagnosticTeardown()
 */
func Teardown(name string) {
	// ensure both list and diagnostic list are removed.
	defer func() { delete(diagnosticTeardownLists, name) }()
	defer func() { delete(teardownLists, name) }()

	list := teardownLists[name]
	list = append(list, diagnosticTeardownLists[name]...) // append any diagnostic funcs - so they are called FIRST

	for x := len(list) - 1; x >= 0; x-- {
		list[x]()
	}
}

/**
* Verify all teardownLists have been executed already; and throw an ASSERT if not.
* Can be used to verify correct coding of a test that uses teardown - and to ensure eventual release of resources.
*
* NOTE: while the funcs are called in the correct order for each list,
* there can be NO guarantee that the lists are iterated in the correct order.
*
* This function MUST NOT be used as a replacement for calling teardown() at the correct point in the code.
 */
func VerifyTeardown(t *testing.T) {

	// ensure all funcs in all lists are released
	defer func() { teardownLists = make(map[string][]func()) }()
	defer func() { diagnosticTeardownLists = make(map[string][]func()) }()

	// append each diagnostic list to the corresponding (possibly empty) teardown list
	for name, list := range diagnosticTeardownLists {
		teardownLists[name] = append(teardownLists[name], list...)
	}

	// release all remaining resources - this is a "best effort" as the order of iterating the map is arbitrary
	uncleared := make([]string, 0)

	// make a "best-effort" at releasing all remaining resources
	for name, list := range teardownLists {
		uncleared = append(uncleared, name)

		for x := len(list) - 1; x >= 0; x-- {
			list[x]()
		}
	}

	assert.Equal(t, 0, len(uncleared), "Error - %d teardownLists were left uncleared: %s", len(uncleared), uncleared)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func RemoveEmptyLines(s string) string {
	regex, err := regexp.Compile("(\r|\r\n|\n){2,}")
	if err != nil {
		return s
	}
	s = regex.ReplaceAllString(s, "\n")
	s = strings.TrimRight(s, "\n")

	return s
}

func InjectOpenShiftOverrides(t *testing.T, options *helm.Options) {
	if !IsOpenShiftEnvironment(t) ||
		options.SetValues["admin.readinessTimeoutSeconds"] != "" {
		return
	}

	t.Log("Using OpenShift specific injects")

	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}

	// OpenShift and CodeReadyContainers readiness probes are slower
	options.SetValues["admin.readinessTimeoutSeconds"] = "5"
}

func InjectTestValuesFile(t *testing.T, options *helm.Options) {
	dat, err := ioutil.ReadFile(INJECT_VALUES_FILE)
	if err != nil {
		return
	}
	t.Logf("Using injected values file=%s with content:%s\n", INJECT_VALUES_FILE, string(dat))
	options.ValuesFiles = append(options.ValuesFiles, INJECT_VALUES_FILE)
}

func InjectTestVersion(t *testing.T, options *helm.Options) {
	dat, err := ioutil.ReadFile(INJECT_FILE)
	if err != nil {
		return
	}

	// do not inject anything if the test overrides these
	// access to nil map yields the default
	if options.SetValues["nuodb.image.registry"] != "" ||
		options.SetValues["nuodb.image.repository"] != "" ||
		options.SetValues["nuodb.image.tag"] != "" {

		return
	}

	t.Log("Using injected values:\n", string(dat))

	err, image := UnmarshalImageYAML(string(dat))
	assert.NoError(t, err)

	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}

	options.SetValues["nuodb.image.registry"] = image.Nuodb.Image.Registry
	options.SetValues["nuodb.image.repository"] = image.Nuodb.Image.Repository
	options.SetValues["nuodb.image.tag"] = image.Nuodb.Image.Tag
}

func InjectTestValues(t *testing.T, options *helm.Options) {
	InjectTestValuesFile(t, options)
	InjectOpenShiftOverrides(t, options)
	InjectTestVersion(t, options)
}

func GetUpgradedReleaseVersion(t *testing.T, options *helm.Options) string {
	// reset all image tags
	delete(options.SetValues, "nuodb.image.registry")
	delete(options.SetValues, "nuodb.image.repository")
	delete(options.SetValues, "nuodb.image.tag")

	InferVersionFromTemplate(t, options)

	return fmt.Sprintf("%s/%s:%s", options.SetValues["nuodb.image.registry"],
		options.SetValues["nuodb.image.repository"],
		options.SetValues["nuodb.image.tag"])

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
	options := k8s.NewKubectlOptions("", "", namespace)
	filter := metav1.ListOptions{}
	return k8s.ListPods(t, options, filter)
}

func Await(t *testing.T, lmbd func() bool, timeout time.Duration) {
	now := time.Now()
	for timeExpired := time.After(timeout); ; {
		select {
		case <-timeExpired:
			t.Logf("Function %s timed out", runtime.FuncForPC(reflect.ValueOf(lmbd).Pointer()).Name())
			t.Logf("Full stack trace of caller:\n%s", string(debug.Stack()))
			t.Fatalf("function call timed out after %f seconds. Start of await was '%s'", timeout.Seconds(), now)
		default:
			if lmbd() {
				return
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func AwaitTillerUp(t *testing.T) {
	version, err := helm.RunHelmCommandAndGetOutputE(t, &helm.Options{}, "version", "--short")
	assert.NoError(t, err)

	t.Logf("Using Helm %s", version)

	if strings.Contains(version, "v3.") {
		return
	}

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
	timeout := 30 * time.Second
	if nrReplicas > 1 {
		timeout *= time.Duration(nrReplicas)
	}
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
	}, timeout)
}

func AwaitNrReplicasReady(t *testing.T, namespace string, expectedName string, nrReplicas int) {
	Await(t, func() bool {
		var cnt int
		for _, pod := range findAllPodsInSchema(t, namespace) {
			if strings.Contains(pod.Name, expectedName) {
				if arePodConditionsMet(&pod, corev1.PodReady, corev1.ConditionTrue) {
					cnt++
				}
			}
		}

		t.Logf("%d pods READY for name '%s'\n", cnt, expectedName)

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

	for _, pod := range findAllPodsInSchema(t, namespace) {
		t.Logf("Pods %s\n", pod.Name)
	}

	return nil, errors.New("did not find any pod matching name")
}

func GetPod(t *testing.T, namespace string, podName string) *corev1.Pod {
	options := k8s.NewKubectlOptions("", "", namespace)

	return k8s.GetPod(t, options, podName)
}

func GetPodName(t *testing.T, namespaceName string, expectedName string) string {
	tePod, err := findPod(t, namespaceName, expectedName)
	assert.NoError(t, err, "No pod found with name ", expectedName)

	return tePod.Name
}

func AwaitPodStatus(t *testing.T, namespace string, podName string, condition corev1.PodConditionType,
	status corev1.ConditionStatus, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "", namespace)

	Await(t, func() bool {
		pod := k8s.GetPod(t, options, podName)
		return arePodConditionsMet(pod, condition, status)
	}, timeout)
}

func AwaitPodPhase(t *testing.T, namespace string, podName string, phase corev1.PodPhase, timeout time.Duration) {
	Await(t, func() bool {
		pod, err := findPod(t, namespace, podName)
		assert.NoError(t, err, "awaitPodPhase: could not find pod with name matching ", podName)

		return pod.Status.Phase == phase
	}, timeout)
}

func AwaitPodUp(t *testing.T, namespace string, adminPodName string, timeout time.Duration) {
	AwaitPodStatus(t, namespace, adminPodName, corev1.PodReady, corev1.ConditionTrue, timeout)
}

func AwaitPodObjectRecreated(t *testing.T, namespace string, pod *corev1.Pod, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "", namespace)

	Await(t, func() bool {
		currentPod, err := k8s.GetPodE(t, options, pod.Name)

		if err != nil {
			return false
		}

		return currentPod.UID != pod.UID
	}, timeout)
}

func AwaitPodTemplateHasVersion(t *testing.T, namespace string, podNameTemplate string, expectedVersion string, timeout time.Duration) {
	Await(t, func() bool {
		pod, err := findPod(t, namespace, podNameTemplate)

		if err != nil {
			t.Logf("No pod found with name %s", podNameTemplate)
			return false
		}

		for _, container := range pod.Spec.Containers {
			t.Logf("Found container (%s) with image: %s", container.Name, container.Image)
			if container.Image == expectedVersion {
				return true
			}
		}

		return false
	}, timeout)
}

func AwaitPodHasVersion(t *testing.T, namespace string, podName string, expectedVersion string, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "", namespace)

	Await(t, func() bool {
		pod, err := k8s.GetPodE(t, options, podName)

		if err != nil {
			t.Logf("No pod found with name %s", podName)
			return false
		}

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
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "domain")

	assert.NoError(t, err, "verifyAdminState: running show domain failed")
	assert.True(t, strings.Contains(output, "ACTIVE"))

}

func AwaitAdminFullyConnected(t *testing.T, namespace string, podName string, numServers int) {
	options := k8s.NewKubectlOptions("", "", namespace)

	k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "check", "servers",
		"--check-active", "--check-connected", "--check-leader",
		"--num-servers", strconv.Itoa(numServers),
		"--timeout", "300")
}

func AwaitDatabaseUp(t *testing.T, namespace string, podName string, databaseName string, numProcesses int) {
	options := k8s.NewKubectlOptions("", "", namespace)

	err := k8s.RunKubectlE(t, options, "exec", podName, "--", "nuocmd", "check", "database",
		"--db-name", databaseName, "--check-running", "--check-liveness", "20",
		"--num-processes", strconv.Itoa(numProcesses),
		"--timeout", "300")

	if err != nil {
		_ = k8s.RunKubectlE(t, options, "exec", podName, "--", "nuocmd", "show", "domain")
	}

	assert.NoError(t, err, "Check database failed. DB not ready after 300s")
}

func GetDiagnoseOnTestFailure(t *testing.T, namespace string, podName string) {
	if t.Failed() && shouldGetDiagnose() {
		options := k8s.NewKubectlOptions("", "", namespace)

		pwd, err := os.Getwd()
		assert.NoError(t, err)

		targetDirPath := filepath.Join(pwd, RESULT_DIR, namespace, "diagnose")
		_ = os.MkdirAll(targetDirPath, 0700)

		// Get cores
		// Once DB-29847 is implemented, we can set a --timeout or --wait-forever flags
		// So that core dump streams doesn't timeout in minikube environments
		t.Log("Generating diagnose archive...")
		k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "get", "diagnose-info",
			"--include-cores", "--output-dir", "/tmp")
		diagnoseFile, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "bash", "-c", "ls -1 /tmp | grep diagnose-")
		assert.NoError(t, err, "Can not find diagnose archive")

		k8s.RunKubectl(t, options, "cp", podName+":/tmp/"+diagnoseFile, filepath.Join(targetDirPath, diagnoseFile))
	}
}

func GetDatabaseIncarnation(t *testing.T, namespace string, podName string, databaseName string) *DBVersion {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "--show-json", "get", "databases")
	assert.NoError(t, err)

	err, databases := UnmarshalDatabase(output)
	assert.NoError(t, err)

	for _, db := range databases {
		if db.Name == databaseName {
			return &db.Incarnation
		}
	}

	t.Logf("GetDatabaseIncarnation did not find DB name: %s", databaseName)
	t.FailNow()
	return nil
}

func AwaitDatabaseRestart(t *testing.T, namespace string, podName string, databaseName string, databaseOptions *helm.Options, restart func()) {
	incarnation := GetDatabaseIncarnation(t, namespace, podName, databaseName)

	restart()

	Await(t, func() bool {
		return GetDatabaseIncarnation(t, namespace, podName, databaseName).Major > incarnation.Major
	}, 300*time.Second)

	opts := GetExtractedOptions(databaseOptions)
	AwaitDatabaseUp(t, namespace, podName, databaseName, opts.NrTePods+opts.NrSmPods)
}

func GetPodRestartCount(t *testing.T, namespace string, podName string) int32 {
	options := k8s.NewKubectlOptions("", "", namespace)

	pod := k8s.GetPod(t, options, podName)

	var restartCount int32
	for _, status := range pod.Status.ContainerStatuses {
		restartCount += status.RestartCount
	}

	return restartCount
}

func AwaitPodRestartCountGreaterThan(t *testing.T, namespace string, podName string, expectedRestartCount int32,
	timeout time.Duration) {
	Await(t, func() bool {
		return GetPodRestartCount(t, namespace, podName) > expectedRestartCount
	}, timeout)
}

func VerifyPolicyInstalled(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "get", "load-balancers")
	assert.NoError(t, err, "VerifyPolicyInstalled: ", podName)
	assert.True(t, strings.Contains(output, "LoadBalancerPolicy"))
}

func VerifyLicenseFile(t *testing.T, namespace string, podName string, expectedLicense string) {
	options := k8s.NewKubectlOptions("", "", namespace)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "cat", "/etc/nuodb/nuodb.lic")
	assert.NoError(t, err, "verifyLicenseFile: exec cat nuodb.lic")
	assert.Equal(t, output, expectedLicense)
}

func VerifyLicenseIsCommunity(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--", "nuocmd", "show", "domain")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(output, "server license: Community"), output)
}

func VerifyLicensingErrorsInLog(t *testing.T, namespace string, podName string, expectError bool) {
	buf, err := ioutil.ReadAll(getAppLogStream(t, namespace, podName, &corev1.PodLogOptions{}))
	assert.NoError(t, err)

	fullLog := string(buf)

	assert.Equal(t, expectError, strings.Contains(fullLog, "Unable to verify license"), fullLog)
}

func GetStringOccurrenceInLog(t *testing.T, namespace string, podName string, expectedLogLine string, podLogOptions *corev1.PodLogOptions) int {
	buf, err := ioutil.ReadAll(getAppLogStream(t, namespace, podName, podLogOptions))
	assert.NoError(t, err)

	fullLog := string(buf)

	return strings.Count(fullLog, expectedLogLine)

}

func VerifyCertificateInLog(t *testing.T, namespace string, podName string, expectedLogLine string) {
	buf, err := ioutil.ReadAll(getAppLogStream(t, namespace, podName, &corev1.PodLogOptions{}))
	assert.NoError(t, err)

	fullLog := string(buf)

	assert.True(t, strings.Contains(standardizeSpaces(fullLog), expectedLogLine),
		"`%s` not found in:\n %s", expectedLogLine, fullLog)
}

func KillAdminPod(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "delete", "pod", podName)
	assert.NoError(t, err, "killAdminPod: delete pod returned an error")
	assert.True(t, strings.Contains(output, "deleted"), "`deleted` not found in %s", output)
}

func KillProcess(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	t.Logf("Killing pid 1 in pod %s\n", podName)
	k8s.RunKubectl(t, options, "exec", podName, "--", "kill", "1")

	AwaitPodRestartCountGreaterThan(t, namespace, podName, 0, 30*time.Second)
}

func GetService(t *testing.T, namespaceName string, serviceName string) *corev1.Service {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	return k8s.GetService(t, options, serviceName)
}

func PingService(t *testing.T, namespace string, serviceName string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	fullServiceName := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)

	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "--",
		"ping", fullServiceName, "-c", "1",
	)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(output, "1 received"))
}

func shouldGetDiagnose() bool {
	_, exists := os.LookupEnv("NUODB_GET_DIAGNOSE")
	return exists
}

func GetK8sEventLog(t *testing.T, namespace string) {
	dirPath := filepath.Join(RESULT_DIR, namespace)
	filePath := filepath.Join(dirPath, K8S_EVENT_LOG_FILE)

	_ = os.MkdirAll(dirPath, 0700)

	f, err := os.Create(filePath)
	assert.NoError(t, err)
	defer f.Close()

	options := k8s.NewKubectlOptions("", "", namespace)

	client, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	assert.NoError(t, err)

	var opts metav1.ListOptions

	events, err := client.CoreV1().Events(namespace).Watch(context.TODO(), opts)
	assert.NoError(t, err)

	writer := io.Writer(f)

	for event := range events.ResultChan() {
		_, err = fmt.Fprintln(writer, event)
		assert.NoError(t, err)
	}

	t.Log("Fully consumed k8s event log")
}

func GetAppLog(t *testing.T, namespace string, podName string, fileNameSuffix string, podLogOptions *corev1.PodLogOptions) {
	dirPath := filepath.Join(RESULT_DIR, namespace)
	filePath := filepath.Join(dirPath, podName+fileNameSuffix+".log")

	_ = os.MkdirAll(dirPath, 0700)

	f, err := os.Create(filePath)
	assert.NoError(t, err)
	defer f.Close()

	writer := io.Writer(f)

	reader := getAppLogStream(t, namespace, podName, podLogOptions)
	assert.NotNil(t, reader)
	_, err = io.Copy(writer, reader)
	assert.NoError(t, err)

	t.Logf("Finished reading log file %s", filePath)
}

func getAppLogStream(t *testing.T, namespace string, podName string, podLogOptions *corev1.PodLogOptions) io.ReadCloser {
	options := k8s.NewKubectlOptions("", "", namespace)

	client, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	assert.NoError(t, err)

	reader, err := client.CoreV1().Pods(options.Namespace).GetLogs(podName, podLogOptions).Stream(context.TODO())
	assert.NoError(t, err)

	return reader
}

func GetAdminEventLog(t *testing.T, namespace string, podName string) {
	pwd, err := os.Getwd()
	assert.NoError(t, err)

	dirPath := filepath.Join(pwd, RESULT_DIR, namespace)
	filePath := filepath.Join(dirPath, podName+"_nuoadmin_event.log")

	_ = os.MkdirAll(dirPath, 0700)

	f, err := os.Create(filePath)
	assert.NoError(t, err)
	defer f.Close()

	options := k8s.NewKubectlOptions("", "", namespace)

	// ignore errors
	_ = k8s.RunKubectlE(t, options,
		"cp",
		fmt.Sprintf("%s/%s:%s", namespace, podName, "/var/log/nuodb/nuoadmin_event.log"),
		filePath,
	)
}

func GetSecret(t *testing.T, namespace string, secretName string) *corev1.Secret {
	options := k8s.NewKubectlOptions("", "", namespace)

	return k8s.GetSecret(t, options, secretName)
}

func GetDaemonSet(t *testing.T, namespace string, daemonSetName string) *v1.DaemonSet {
	options := k8s.NewKubectlOptions("", "", namespace)

	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	assert.NoError(t, err)

	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})

	return daemonSet
}

func GetReplicationController(t *testing.T, namespace string, replicationControllerName string) *corev1.ReplicationController {
	options := k8s.NewKubectlOptions("", "", namespace)

	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, options)
	assert.NoError(t, err)

	controller, err := clientset.CoreV1().ReplicationControllers(namespace).Get(context.TODO(), replicationControllerName, metav1.GetOptions{})

	return controller
}

func DeleteDatabase(t *testing.T, namespace string, dbName string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	k8s.RunKubectl(t, options, "exec", podName, "--", "nuocmd", "delete", "database", "--db-name", dbName)
}

func DeletePod(t *testing.T, namespace string, podName string) {
	options := k8s.NewKubectlOptions("", "", namespace)

	k8s.RunKubectl(t, options, "delete", podName)
}

func RunSQL(t *testing.T, namespace string, podName string, databaseName string, sql string) (result string, err error) {
	options := k8s.NewKubectlOptions("", "", namespace)

	// secrets := getSecret(t, namespace, databaseName)

	result, err = k8s.RunKubectlAndGetOutputE(t, options,
		"exec", podName, "--",
		"bash", "-c",
		fmt.Sprintf("echo \"%s;\" | /opt/nuodb/bin/nuosql --user dba --password secret %s", sql, databaseName),
	)

	assert.NoError(t, err, "runSQL: error trying to run ", sql)

	return result, err
}

func GetNuoDBK8sConfigDump(t *testing.T, namespace string, podName string) NuoDBKubeConfig {
	dumpFileName := "nuodb-dump.json"

	options := k8s.NewKubectlOptions("", "", namespace)

	pwd, err := os.Getwd()
	assert.NoError(t, err)

	targetDirPath := filepath.Join(pwd, RESULT_DIR, namespace, "k8s-dump")
	_ = os.MkdirAll(targetDirPath, 0700)

	targetFile := filepath.Join(targetDirPath, dumpFileName)

	k8s.RunKubectl(t, options,
		"exec", podName, "--",
		"bash", "-c",
		"nuocmd --show-json get kubernetes-config > /tmp/nuodb-dump.json",
	)

	k8s.RunKubectl(t, options, "cp", podName+":/tmp/nuodb-dump.json", targetFile)

	content, err := ioutil.ReadFile(targetFile)
	assert.NoError(t, err)
	err, unmarshalledDump := UnmarshalNuoDBKubeConfig(string(content))
	assert.NoError(t, err)
	assert.Equal(t, len(unmarshalledDump), 1)
	return unmarshalledDump[0]
}

func ExecuteCommandsInPod(t *testing.T, namespaceName string, podName string, commands []string) {
	tmpfile, err := ioutil.TempFile("", "script")
	if err != nil {
		assert.NoError(t, err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.WriteString("set -ev" + "\n"); err != nil {
		assert.NoError(t, err)
	}

	for _, item := range commands {
		if _, err := tmpfile.WriteString(item + "\n"); err != nil {
			assert.NoError(t, err)
		}
	}
	if err := tmpfile.Close(); err != nil {
		assert.NoError(t, err)
	}

	options := k8s.NewKubectlOptions("", "", namespaceName)

	// Transfer the TEMP script to POD and execute it
	k8s.RunKubectl(t, options, "cp", tmpfile.Name(), podName+":/tmp")
	k8s.RunKubectl(t, options, "exec", podName, "--", "chmod", "a+x", "/tmp/"+filepath.Base(tmpfile.Name()))
	err = k8s.RunKubectlE(t, options, "exec", podName, "--", "sh", "/tmp/"+filepath.Base(tmpfile.Name()))
	assert.NoError(t, err, "executeCommandsInPod: Script returned error.")
}

func UnmarshalJSONObject(t *testing.T, stringJSON string) map[string]interface{} {
	var results map[string]interface{}
	err := json.Unmarshal([]byte(stringJSON), &results)
	assert.NoError(t, err)
	return results
}

func VerifyAdminKvSetAndGet(t *testing.T, podName string, namespaceName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	// verify the KV store can write and read in a reasonable time
	start := time.Now()
	output, err := k8s.RunKubectlAndGetOutputE(t, options,
		"exec", podName, "--", "nuocmd", "set", "value", "--key", "test/minikube", "--value", "testVal", "--unconditional",
	)
	assert.NoError(t, err, "Could not set KV value")
	elapsed := time.Since(start)

	assert.True(t, elapsed.Seconds() < 2.0, fmt.Sprintf("KV set took longer than 2s: %s", elapsed))

	start = time.Now()
	output, err = k8s.RunKubectlAndGetOutputE(t, options,
		"exec", podName, "--", "nuocmd", "get", "value", "--key", "test/minikube",
	)
	assert.NoError(t, err, "Could not get KV value")
	elapsed = time.Since(start)

	assert.True(t, elapsed.Seconds() < 2.0, fmt.Sprintf("KV get took longer than 2s: %s", elapsed))

	assert.True(t, output == "testVal", fmt.Sprintf("KV get returned the wrong value: %s", output))
}

func LabelNodes(t *testing.T, namespaceName string, labelName string, labelValue string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	var labelString string

	if labelValue != "" {
		labelString = fmt.Sprintf("%s=%s", labelName, labelValue)
	} else {
		labelString = fmt.Sprintf("%s-", labelName)
	}

	nodes := k8s.GetNodes(t, options)

	assert.True(t, len(nodes) > 0)

	for _, node := range nodes {
		err := k8s.RunKubectlE(t, options, "label", "node", node.Name, labelString, "--overwrite")
		assert.NoError(t, err, "Labeling node %s with '%s' failed", node.Name, labelString)
	}
}
