package testlib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/otiai10/copy"
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

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.KubectlApplyFromString(t, kubectlOptions, podTemplateString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromStringE(t, kubectlOptions, podTemplateString) })

	AwaitPodUp(t, namespaceName, podName, timeout)

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
	assert.NoError(t, err)

	set := make(map[string]bool)
	for _, file := range files {
		set[file.Name()] = true
		t.Logf("Found generated certificate file: %s", file.Name())
	}
	for _, expectedFile := range expectedFiles {
		assert.True(t, set[expectedFile] == true, "Unable to find certificate file %s in path %s", expectedFile, directory)
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
	ExecuteCommandsInPod(t, namespaceName, podName, finalCommands)
}

func CopyCertificatesToControlHost(t *testing.T, podName string, namespaceName string) string {
	options := k8s.NewKubectlOptions("", "", namespaceName)

	prefix := "tls-keys"
	targetDirectory, err := ioutil.TempDir("", prefix)
	assert.NoError(t, err, "Unable to create TMP directory with prefix ", prefix)
	AddTeardown(TEARDOWN_SECRETS, func() { os.RemoveAll(targetDirectory) })

	realTargetDirectory, err := filepath.EvalSymlinks(targetDirectory)
	assert.NoError(t, err)

	k8s.RunKubectl(t, options, "cp", podName+":"+CERTIFICATES_GENERATION_PATH, realTargetDirectory)
	t.Logf("Certificate files location: %s", realTargetDirectory)
	AddTeardown(TEARDOWN_SECRETS, func() { BackupCerificateFilesOnTestFailure(t, namespaceName, realTargetDirectory) })
	verifyCertificateFiles(t, realTargetDirectory)

	return realTargetDirectory
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

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	assert.NoError(t, err, "Unable to find/convert admin.replicas value")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)

	k8s.RunKubectl(t, kubectlOptions, "cp", filepath.Join(tlsKeysLocation, CA_CERT_FILE_NEW), admin0+":/tmp")
	err = k8s.RunKubectlE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "add", "trusted-certificate",
		"--alias", "ca_prime", "--cert", "/tmp/"+CA_CERT_FILE_NEW, "--timeout", "60")
	assert.NoError(t, err, "add trusted-certificate failed")

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
		// Faster approach will be to restart all PODs. A prerequisite for this
		// is to have the same secrets update before hand.
		adminPod := GetPod(t, namespaceName, admin0)
		k8s.RunKubectl(t, kubectlOptions, "delete", "pod", "--selector=domain=nuodb")
		AwaitPodObjectRecreated(t, namespaceName, adminPod, 30*time.Second)
		AwaitPodUp(t, namespaceName, admin0, 300*time.Second)
		AwaitAdminFullyConnected(t, namespaceName, admin0, adminReplicaCount)
	}
	AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
}
