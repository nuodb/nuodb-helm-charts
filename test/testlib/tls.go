package testlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/otiai10/copy"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
)

const TLS_GENERATOR_POD_TEMPLATE = `---
apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  containers:
  - name: generate-tls-certs
    image: %s
    command: ['tail', '-f', '/dev/null']
`

func createTLSGeneratorPod(t *testing.T, namespaceName string, timeout time.Duration) string {
	image := "docker.io/nuodb/nuodb-ce:latest"

	podName := "tls-generator-" + strings.ToLower(random.UniqueId())
	podTemplateString := fmt.Sprintf(TLS_GENERATOR_POD_TEMPLATE,
		podName, namespaceName, image)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	k8s.KubectlApplyFromString(t, kubectlOptions, podTemplateString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromStringE(t, kubectlOptions, podTemplateString) })

	AwaitPodStatus(t, namespaceName, podName, corev1.PodReady, corev1.ConditionTrue, timeout)

	return podName
}

func verifyCertificateFiles(t *testing.T, directory string) {
	expectedFiles := []string{
		KEYSTORE_FILE,
		TRUSTSTORE_FILE,
		CA_CERT_FILE,
		NUOCMD_FILE,
	}

	files, err := ioutil.ReadDir(directory)
	assert.NilError(t, err)

	set := make(map[string]bool)
	for _, file := range files {
		set[file.Name()] = true
		t.Logf("Found generated certificate file: %s", file.Name())
	}
	for _, expectedFile := range expectedFiles {
		assert.Assert(t, set[expectedFile] == true, "Unable to find certificate file %s in path %s", expectedFile, directory)
	}
}

func GenerateCustomCertificates(t *testing.T, podName string, namespaceName string, commands []string) {
	prependCommands := []string{
		"[ -d " + CERTIFICATES_GENERATION_PATH + " ] && rm -rf " + CERTIFICATES_BACKUP_PATH + " && mv " + CERTIFICATES_GENERATION_PATH + " " + CERTIFICATES_BACKUP_PATH,
		"rm -rf " + CERTIFICATES_GENERATION_PATH,
		"mkdir -p " + CERTIFICATES_GENERATION_PATH,
		"cd " + CERTIFICATES_GENERATION_PATH,
	}
	finalCommands := append(prependCommands, commands...)
	// Execute certificate generation commands
	ExecuteCommandsInPod(t, podName, namespaceName, finalCommands)
}

func CopyCertificatesToControlHost(t *testing.T, podName string, namespaceName string) string {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	prefix := "tls-keys"
	targetDirectory, err := ioutil.TempDir("", prefix)
	assert.NilError(t, err, "Unable to create TMP directory with prefix ", prefix)
	AddTeardown(TEARDOWN_SECRETS, func() { os.RemoveAll(targetDirectory) })

	realTargetDirectory, err := filepath.EvalSymlinks(targetDirectory)
	assert.NilError(t, err)

	k8s.RunKubectl(t, options, "cp", podName+":"+CERTIFICATES_GENERATION_PATH, realTargetDirectory)
	t.Logf("Certificate files location: %s", realTargetDirectory)
	AddTeardown(TEARDOWN_SECRETS, func() { BackupCerificateFilesOnTestFailure(t, namespaceName, realTargetDirectory) })
	AddTeardown(TEARDOWN_SECRETS, func() { PrintCertificateFilesOnTestFailure(t, realTargetDirectory) })
	verifyCertificateFiles(t, realTargetDirectory)

	return realTargetDirectory
}

func PrintCertificateFilesOnTestFailure(t *testing.T, srcDirectory string) {
	if t.Failed() && shouldPrintToStdout() {
		files, _ := ioutil.ReadDir(srcDirectory)
		t.Logf("Printing certificate files in %s.", srcDirectory)
		for _, file := range files {
			output, _ := exec.Command("hexdump", "-ve", `16/1 "%02x " "\n"`, filepath.Join(srcDirectory, file.Name())).CombinedOutput()
			t.Logf("%s:\n%s", file.Name(), string(output))
		}
	}
}

func BackupCerificateFilesOnTestFailure(t *testing.T, namespaceName string, srcDirectory string) {
	targetDirPath := filepath.Join(RESULT_DIR, namespaceName, filepath.Base(srcDirectory))
	_ = os.MkdirAll(targetDirPath, 0700)
	if t.Failed() {
		err := copy.Copy(srcDirectory, targetDirPath)
		if err != nil {
			t.Logf("Unable to backup certificates in %s", srcDirectory)
		}
		t.Logf("Certificate files copied from %s to %s", srcDirectory, targetDirPath)
	}
}

func GenerateTLSConfiguration(t *testing.T, namespaceName string, commands []string) (string, string) {
	podName := createTLSGeneratorPod(t, namespaceName, 300*time.Second) // this might pull an image
	GenerateCustomCertificates(t, podName, namespaceName, commands)
	keysLocation := CopyCertificatesToControlHost(t, podName, namespaceName)

	CreateSecret(t, namespaceName, CA_CERT_FILE, CA_CERT_SECRET, keysLocation)
	CreateSecret(t, namespaceName, NUOCMD_FILE, NUOCMD_SECRET, keysLocation)
	CreateSecretWithPassword(t, namespaceName, KEYSTORE_FILE, KEYSTORE_SECRET, SECRET_PASSWORD, keysLocation)
	CreateSecretWithPassword(t, namespaceName, TRUSTSTORE_FILE, TRUSTSTORE_SECRET, SECRET_PASSWORD, keysLocation)

	return podName, keysLocation
}

func RotateTLSCertificates(t *testing.T, options *helm.Options, namespaceName string,
	adminReleaseName string, databaseReleaseName string, tlsKeysLocation string, helmUpgrade bool) {

	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	assert.NilError(t, err, "Unable to find/convert admin.replicas value")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)

	k8s.RunKubectl(t, kubectlOptions, "cp", filepath.Join(tlsKeysLocation, CA_CERT_FILE_NEW), admin0+":/tmp")
	err = k8s.RunKubectlE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "add", "trusted-certificate",
		"--alias", "ca_prime", "--cert", "/tmp/"+CA_CERT_FILE_NEW, "--timeout", "60")
	assert.NilError(t, err, "add trusted-certificate failed")

	if helmUpgrade == true {
		// Upgrade admin release
		DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")
		helm.Upgrade(t, options, ADMIN_HELM_CHART_PATH, adminReleaseName)

		adminStatefulSet := fmt.Sprintf("%s-nuodb", adminReleaseName)
		k8s.RunKubectl(t, kubectlOptions, "rollout", "status", "sts/"+adminStatefulSet, "--timeout", "300s")
		AwaitAdminFullyConnected(t, namespaceName, admin0, adminReplicaCount)

		// Upgrade database release
		helm.Upgrade(t, options, DATABASE_HELM_CHART_PATH, databaseReleaseName)
	} else {
		// Rolling upgrade could take a lot of time due to readiness probes.
		// Faster approach will be to restart all PODs. A prerequsite for this 
		// is to have the same secrets update before hand.
		k8s.RunKubectl(t, kubectlOptions, "delete", "pod", "--selector=domain=nuodb")
		AwaitPodPhase(t, namespaceName, admin0, corev1.PodRunning, 60*time.Second)
		AwaitAdminFullyConnected(t, namespaceName, admin0, adminReplicaCount)
	}
	AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
}
