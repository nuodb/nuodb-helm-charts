package testlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
	// 4.2.2 has Java 11.0.12, which uses default algorithms for signatures
	// and encryption of private keys that are not available in Java 8, used
	// by NuoDB up to 4.0.7; to avoid breaking rolling upgrade tests, fix
	// version at 4.2.1
	image := "docker.io/nuodb/nuodb-ce:4.2.1"

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
	require.NoError(t, err)

	set := make(map[string]bool)
	for _, file := range files {
		set[file.Name()] = true
		t.Logf("Found generated certificate file: %s", file.Name())
	}
	for _, expectedFile := range expectedFiles {
		require.True(t, set[expectedFile] == true, "Unable to find certificate file %s in path %s", expectedFile, directory)
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
	require.NoError(t, err, "Unable to create TMP directory with prefix ", prefix)
	AddTeardown(TEARDOWN_SECRETS, func() { os.RemoveAll(targetDirectory) })

	realTargetDirectory, err := filepath.EvalSymlinks(targetDirectory)
	require.NoError(t, err)

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
	opt := GetExtractedOptions(options)

	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	require.NoError(t, err, "Unable to find/convert admin.replicas value")
	adminStatefulSet := fmt.Sprintf("%s-%s-%s", adminReleaseName, opt.DomainName, opt.ClusterName)
	admin0 := fmt.Sprintf("%s-0", adminStatefulSet)

	// Add the new CA certificate to all APs and engines
	k8s.RunKubectl(t, kubectlOptions, "cp", filepath.Join(tlsKeysLocation, CA_CERT_FILE_NEW), admin0+":/tmp")
	err = k8s.RunKubectlE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "add", "trusted-certificate",
		"--alias", "ca_prime", "--cert", "/tmp/"+CA_CERT_FILE_NEW, "--timeout", "60")
	require.NoError(t, err, "add trusted-certificate failed")

	if helmUpgrade == true {
		// Read the server certificate from the keystore file
		k8s.RunKubectl(t, kubectlOptions, "cp", filepath.Join(tlsKeysLocation, KEYSTORE_FILE), admin0+":/tmp")
		adminPem, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "-c", "admin", "--",
			"nuocmd", "show", "certificate", "--keystore", "/tmp/"+KEYSTORE_FILE,
			"--store-password", SECRET_PASSWORD, "--cert-only")
		require.NoError(t, err)

		// Upgrade admin release
		helm.Upgrade(t, options, ADMIN_HELM_CHART_PATH, adminReleaseName)
		// The total delay for the new server certificates to be picked up by
		// APs is kubelet syncFrequency (by default 60sec) + NuoAdmin
		// keystoreUpdateIntervalMs (by default 5sec)
		AwaitAdminKeystoreReload(t, namespaceName, admin0, adminPem, 120*time.Second)
		AwaitAdminFullyConnected(t, namespaceName, admin0, adminReplicaCount)

		// Upgrade database release
		tePodNameTemplate := fmt.Sprintf("te-%s-%s-%s-%s", databaseReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
		smPodNameTemplate := fmt.Sprintf("sm-%s-%s-%s-%s", databaseReleaseName, opt.DomainName, opt.ClusterName, opt.DbName)
		smPod := GetPod(t, namespaceName, GetPodName(t, namespaceName, smPodNameTemplate))
		tePod := GetPod(t, namespaceName, GetPodName(t, namespaceName, tePodNameTemplate))
		helm.Upgrade(t, options, DATABASE_HELM_CHART_PATH, databaseReleaseName)
		// Wait for engines to be restarted since the keystore has been changed;
		// the TE pod can be recreated with a different name
		AwaitPodObjectRecreated(t, namespaceName, smPod, 30*time.Second)
		Await(t, func() bool {
			pod, err := FindPod(t, namespaceName, tePodNameTemplate)
			if err != nil {
				t.Logf("%s: %s", err.Error(), tePodNameTemplate)
				return false
			}
			return tePod.UID != pod.UID
		}, 30*time.Second)
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

func AwaitAdminKeystoreReload(t *testing.T, namespace string, podName string, adminPem string, timeout time.Duration) {
	options := k8s.NewKubectlOptions("", "", namespace)
	adminPem = strings.TrimSpace(adminPem)
	Await(t, func() bool {
		data, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", "admin", "--",
			"nuocmd", "--show-json", "get", "certificate-info")
		require.NoError(t, err)
		err, certificateInfo := UnmarshalCertificateInfo(data)
		require.NoError(t, err)
		// Wait for all APs to reload their keystore files and return the
		// updated server certificate
		var notUpdated []string
		for admin, info := range certificateInfo.ServerCertificates {
			actual := strings.TrimSpace(info.CertificatePem)
			if actual != adminPem {
				notUpdated = append(notUpdated, admin)
			}
		}
		if len(notUpdated) > 0 {
			t.Logf("Server certificates not updated for APs %s", notUpdated)
			return false
		}
		return true
	}, timeout)
}

func GenerateAndSetTLSKeys(t *testing.T, options *helm.Options, namespaceName string) (string, string) {
	tlsCommands := []string{
		"export DEFAULT_PASSWORD='" + SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}
	podName, keysLocation := GenerateTLSConfiguration(t, namespaceName, tlsCommands)
	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}
	options.SetValues["admin.tlsCACert.secret"] = CA_CERT_SECRET
	options.SetValues["admin.tlsCACert.key"] = CA_CERT_FILE
	options.SetValues["admin.tlsKeyStore.secret"] = KEYSTORE_SECRET
	options.SetValues["admin.tlsKeyStore.key"] = KEYSTORE_FILE
	options.SetValues["admin.tlsKeyStore.password"] = SECRET_PASSWORD
	options.SetValues["admin.tlsTrustStore.secret"] = TRUSTSTORE_SECRET
	options.SetValues["admin.tlsTrustStore.key"] = TRUSTSTORE_FILE
	options.SetValues["admin.tlsTrustStore.password"] = SECRET_PASSWORD
	options.SetValues["admin.tlsClientPEM.secret"] = NUOCMD_SECRET
	options.SetValues["admin.tlsClientPEM.key"] = NUOCMD_FILE
	return podName, keysLocation
}
