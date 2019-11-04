package testlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
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

func createTLSGeneratorPod(t *testing.T, namespaceName string, image string, timeout time.Duration) string {
	if image == "" {
		image = "docker.io/nuodb/nuodb-ce:latest"
	}

	podName := "tls-generator-" + strings.ToLower(random.UniqueId())
	podTemplateString := fmt.Sprintf(TLS_GENERATOR_POD_TEMPLATE,
		podName, namespaceName, image)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	k8s.KubectlApplyFromString(t, kubectlOptions, podTemplateString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromString(t, kubectlOptions, podTemplateString) })

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
	verifyCertificateFiles(t, realTargetDirectory)

	return realTargetDirectory
}

func GenerateTLSConfiguration(t *testing.T, namespaceName string, commands []string, image string) (string, string) {
	podName := createTLSGeneratorPod(t, namespaceName, image, 30*time.Second)
	GenerateCustomCertificates(t, podName, namespaceName, commands)
	keysLocation := CopyCertificatesToControlHost(t, podName, namespaceName)

	CreateSecret(t, namespaceName, CA_CERT_FILE, CA_CERT_SECRET, keysLocation)
	CreateSecret(t, namespaceName, NUOCMD_FILE, NUOCMD_SECRET, keysLocation)
	CreateSecretWithPassword(t, namespaceName, KEYSTORE_FILE, KEYSTORE_SECRET, SECRET_PASSWORD, keysLocation)
	CreateSecretWithPassword(t, namespaceName, TRUSTSTORE_FILE, TRUSTSTORE_SECRET, SECRET_PASSWORD, keysLocation)

	return podName, keysLocation
}

func RotateTLSCertificates(t *testing.T, options *helm.Options,
	kubectlOptions *k8s.KubectlOptions, adminReleaseName string, databaseReleaseName string, tlsKeysLocation string) {

	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	assert.NilError(t, err, "Unable to find/convert admin.replicas value")
	namespaceName := kubectlOptions.Namespace
	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)

	k8s.RunKubectl(t, kubectlOptions, "cp", filepath.Join(tlsKeysLocation, CA_CERT_FILE_NEW), admin0+":/tmp")
	err = k8s.RunKubectlE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "add", "trusted-certificate",
		"--alias", "ca_prime", "--cert", "/tmp/"+CA_CERT_FILE_NEW, "--timeout", "60")
	assert.NilError(t, err, "add trusted-certificate failed")

	// Upgrade admin release
	DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")
	helm.Upgrade(t, options, ADMIN_HELM_CHART_PATH, adminReleaseName)

	adminStatefulSet := fmt.Sprintf("%s-nuodb", adminReleaseName)
	k8s.RunKubectl(t, kubectlOptions, "rollout", "status", "sts/"+adminStatefulSet, "--timeout", "300s")
	AwaitAdminFullyConnected(t, namespaceName, admin0, adminReplicaCount)

	// Upgrade database release
	helm.Upgrade(t, options, DATABASE_HELM_CHART_PATH, databaseReleaseName)

	AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)
}
