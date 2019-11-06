package testlib

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/otiai10/copy"
	"gotest.tools/assert"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
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

func createTLSGeneratorPod(t *testing.T, namespaceName string, image string, timeout time.Duration) string {
	if image == "" {
		image = "docker.io/nuodb/nuodb-ce:latest"
	}

	podName := "tls-generator"
	podTemplateString := fmt.Sprintf(TLS_GENERATOR_POD_TEMPLATE,
		podName, namespaceName, image)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	k8s.KubectlApplyFromString(t, kubectlOptions, podTemplateString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromStringE(t, kubectlOptions, podTemplateString) })

	AwaitNrReplicasScheduled(t, namespaceName, podName, 1)
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

func GenerateCustomCertificates(t *testing.T, namespaceName string, podName string, commands []string) {
	prependCommands := []string{
		"[ -d " + CERTIFICATES_GENERATION_PATH + " ] && rm -rf " + CERTIFICATES_BACKUP_PATH + " && mv " + CERTIFICATES_GENERATION_PATH + " " + CERTIFICATES_BACKUP_PATH,
		"rm -rf " + CERTIFICATES_GENERATION_PATH,
		"mkdir -p " + CERTIFICATES_GENERATION_PATH,
		"cd " + CERTIFICATES_GENERATION_PATH,
	}
	finalCommands := append(prependCommands, commands...)
	// Execute certificate generation commands
	ExecuteCommandsInPod(t, namespaceName, podName, finalCommands)
}

func CopyCertificatesToControlHost(t *testing.T, namespaceName string, podName string) string {
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
	AddTeardown(TEARDOWN_SECRETS, func() { BackupCertificateFilesOnTestFailure(t, namespaceName, realTargetDirectory) })
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

func BackupCertificateFilesOnTestFailure(t *testing.T, namespaceName string, srcDirectory string) {
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

func GenerateTLSConfiguration(t *testing.T, namespaceName string, commands []string, image string) (string, string) {
	podName := createTLSGeneratorPod(t, namespaceName, image, 30*time.Second)
	GenerateCustomCertificates(t, namespaceName, podName, commands)
	keysLocation := CopyCertificatesToControlHost(t, namespaceName, podName)

	CreateSecret(t, namespaceName, CA_CERT_FILE, CA_CERT_SECRET, keysLocation)
	CreateSecret(t, namespaceName, NUOCMD_FILE, NUOCMD_SECRET, keysLocation)
	CreateSecretWithPassword(t, namespaceName, KEYSTORE_FILE, KEYSTORE_SECRET, SECRET_PASSWORD, keysLocation)
	CreateSecretWithPassword(t, namespaceName, TRUSTSTORE_FILE, TRUSTSTORE_SECRET, SECRET_PASSWORD, keysLocation)

	return podName, keysLocation
}

func AddTrustedCertificate(t *testing.T, namespaceName string, podName string, tlsKeysLocation string, certFileName string) {
	options := k8s.NewKubectlOptions("", "")
	options.Namespace = namespaceName

	k8s.RunKubectl(t, options, "cp", filepath.Join(tlsKeysLocation, certFileName), podName+":/tmp")
	err := k8s.RunKubectlE(t, options, "exec", podName, "--", "nuocmd", "add", "trusted-certificate",
		"--alias", "ca_prime", "--cert", "/tmp/"+certFileName, "--timeout", "60")
	assert.NilError(t, err, "add trusted-certificate failed")

}