// +build long

package minikube

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
	"encoding/json"
	"strconv"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"gotest.tools/assert"
)

func unmarshalCertificateInfo(t *testing.T, certificateInfoJSON string) map[string]interface{} {
	var results map[string]interface{}
	err := json.Unmarshal([]byte(certificateInfoJSON), &results)
	assert.NilError(t, err)
	return results
}

func verifyAdminCertificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := unmarshalCertificateInfo(t, certificateInfoJSON)
	for _, value := range certificateInfo["serverCertificates"].(map[string]interface{}) {
		certSubjectName := value.(map[string]interface{})["subjectName"].(string)
		assert.Check(t, strings.Contains(certSubjectName, expectedDN),
		"`%s` not found in:\n %s", expectedDN, certSubjectName)
	}
}

func verifyEngineCertificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := unmarshalCertificateInfo(t, certificateInfoJSON)
	for _, value := range certificateInfo["processCertificates"].(map[string]interface{}) {
		certIssuerName := value.(map[string]interface{})["issuerName"].(string)
		assert.Check(t, strings.Contains(certIssuerName, expectedDN),
		"`%s` not found in:\n %s", expectedDN, certIssuerName)
	}
}

func startDomainWithTLSCertificates(t *testing.T, options *helm.Options, namespaceName string, tlsCommands []string) (string, string, string) {
	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	assert.NilError(t, err, "Unable to find/convert admin.replicas value")

	// create initial certs...
	certGeneratorPodName, _ := testlib.GenerateTLSConfiguration(t, namespaceName, tlsCommands, "")

	adminReleaseName, namespaceName := testlib.StartAdmin(t, options, adminReplicaCount, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)

	return certGeneratorPodName, adminReleaseName, databaseReleaseName
}

func TestKubernetesTLSCARotation(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	namespaceName := fmt.Sprintf("test-tls-ca-rotation-%s", randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "")
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	kubectlOptions.Namespace = namespaceName

	initialTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas":                        "2",
			"admin.tlsCACert.secret":                testlib.CA_CERT_SECRET,
			"admin.tlsCACert.key":                   testlib.CA_CERT_FILE,
			"admin.tlsKeyStore.secret":              testlib.KEYSTORE_SECRET,
			"admin.tlsKeyStore.key":                 testlib.KEYSTORE_FILE,
			"admin.tlsKeyStore.password":            testlib.SECRET_PASSWORD,
			"admin.tlsTrustStore.secret":            testlib.TRUSTSTORE_SECRET,
			"admin.tlsTrustStore.key":               testlib.TRUSTSTORE_FILE,
			"admin.tlsTrustStore.password":          testlib.SECRET_PASSWORD,
			"admin.tlsClientPEM.secret":             testlib.NUOCMD_SECRET,
			"admin.tlsClientPEM.key":                testlib.NUOCMD_FILE,
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "250m", // during upgrade we will be running 2 of these
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	}

	expectedCaDN := "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=123456"
	expectedAdminDN := "CN=nuoadmin.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"

	newTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"nuocmd create keypair --keystore ca.p12 --store-password \"$DEFAULT_PASSWORD\" --dname \"" + expectedCaDN + "\" --validity 36500 --ca",
		"nuocmd create keypair --keystore " + testlib.KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" --dname \"" + expectedAdminDN + "\"",
		"nuocmd sign certificate --keystore " + testlib.KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" --ca-keystore ca.p12 --ca-store-password \"$DEFAULT_PASSWORD\" --validity 36500 --ca --update",
		"nuocmd show certificate --keystore ca.p12 --store-password \"$DEFAULT_PASSWORD\" --cert-only > " + testlib.CA_CERT_FILE_NEW,
		"cp " + filepath.Join(testlib.CERTIFICATES_BACKUP_PATH, testlib.TRUSTSTORE_FILE) + " " + testlib.CERTIFICATES_GENERATION_PATH,
		"cp " + filepath.Join(testlib.CERTIFICATES_BACKUP_PATH, testlib.NUOCMD_FILE) + " " + testlib.CERTIFICATES_GENERATION_PATH,
		"cat " + filepath.Join(testlib.CERTIFICATES_BACKUP_PATH, testlib.CA_CERT_FILE) + " > ca.cert",
		"cat " + testlib.CA_CERT_FILE_NEW + " >> ca.cert",
	}

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	certGeneratorPodName, adminReleaseName, _ := startDomainWithTLSCertificates(t, &options, namespaceName, initialTLSCommands)

	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)

	// create the new certs and replace the existing secrets with new content
	testlib.GenerateCustomCertificates(t, namespaceName, certGeneratorPodName, newTLSCommands)
	newTLSKeysLocation := testlib.CopyCertificatesToControlHost(t, namespaceName, certGeneratorPodName)
	testlib.CreateSecret(t, namespaceName, testlib.CA_CERT_FILE, testlib.CA_CERT_SECRET, newTLSKeysLocation)
	testlib.CreateSecretWithPassword(t, namespaceName, testlib.KEYSTORE_FILE, testlib.KEYSTORE_SECRET, testlib.SECRET_PASSWORD, newTLSKeysLocation)

	testlib.AddTrustedCertificate(t, namespaceName, admin0, newTLSKeysLocation, testlib.CA_CERT_FILE_NEW)

	// Rolling upgrade could take a lot of time due to readiness probes.
	// Faster approach will be to restart all PODs. A prerequisite for this
	// is to have the same secrets update as before.
	testlib.DeleteAllPodsInDomain(t, namespaceName, "nuodb")
	testlib.AwaitPodPhase(t, namespaceName, admin0, v1.PodRunning, 60*time.Second)
	testlib.AwaitAdminFullyConnected(t, namespaceName, admin0, 2)
	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

	certificateInfo, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "--show-json", "get", "certificate-info")
	assert.NilError(t, err)
	
	t.Run("verifyAdminCertificates", func(t *testing.T) {
		verifyAdminCertificates(t, certificateInfo, expectedAdminDN)
	})

	t.Run("verifyEngineCertificates", func(t *testing.T) {
		verifyEngineCertificates(t, certificateInfo, expectedAdminDN)
	})
}

func TestKubernetesTLSAdminRotation(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	namespaceName := fmt.Sprintf("test-tls-admin-rotation-%s", randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "")
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	kubectlOptions.Namespace = namespaceName

	initialTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}

	options := helm.Options{
		SetValues: map[string]string{
			"admin.replicas":                        "2",
			"admin.tlsCACert.secret":                testlib.CA_CERT_SECRET,
			"admin.tlsCACert.key":                   testlib.CA_CERT_FILE,
			"admin.tlsKeyStore.secret":              testlib.KEYSTORE_SECRET,
			"admin.tlsKeyStore.key":                 testlib.KEYSTORE_FILE,
			"admin.tlsKeyStore.password":            testlib.SECRET_PASSWORD,
			"admin.tlsTrustStore.secret":            testlib.TRUSTSTORE_SECRET,
			"admin.tlsTrustStore.key":               testlib.TRUSTSTORE_FILE,
			"admin.tlsTrustStore.password":          testlib.SECRET_PASSWORD,
			"admin.tlsClientPEM.secret":             testlib.NUOCMD_SECRET,
			"admin.tlsClientPEM.key":                testlib.NUOCMD_FILE,
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    "250m", // during upgrade we will be running 2 of these
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	}

	expectedAdminDN := "CN=nuoadmin.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	certGeneratorPodName, adminReleaseName, _ := startDomainWithTLSCertificates(t, &options, namespaceName, initialTLSCommands)

	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-1", adminReleaseName)

	// generate new admin key using the same CA and store it in a new PKCS12 file
	newTLSCommands := []string{
		"cd " + testlib.CERTIFICATES_GENERATION_PATH,
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"nuocmd create keypair --keystore " + testlib.NEW_KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" --dname \"" + expectedAdminDN + "\"",
		"nuocmd sign certificate --keystore " + testlib.NEW_KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" --ca-keystore ca.p12 --ca-store-password \"$DEFAULT_PASSWORD\" --validity 36500 --ca --update",
	}

	// create new admin cert
	testlib.ExecuteCommandsInPod(t, namespaceName, certGeneratorPodName, newTLSCommands)

	newTLSKeysLocation := testlib.CopyCertificatesToControlHost(t, namespaceName, certGeneratorPodName)
	testlib.CreateSecretWithPassword(t, namespaceName, testlib.NEW_KEYSTORE_FILE, testlib.NEW_KEYSTORE_SECRET, testlib.SECRET_PASSWORD, newTLSKeysLocation)

	// update options to use the new key saved in the new secret
	options.SetValues["admin.tlsKeyStore.key"] =   	testlib.NEW_KEYSTORE_FILE
	options.SetValues["admin.tlsKeyStore.secret"] = testlib.NEW_KEYSTORE_SECRET

	// Upgrade admin release
	testlib.DeletePod(t, namespaceName, "jobs/job-lb-policy-nearest")
	helm.Upgrade(t, &options, testlib.ADMIN_HELM_CHART_PATH, adminReleaseName)

	// wait for both to go through the restart and cycle the readiness
	testlib.AwaitPodStatus(t, namespaceName, admin1, v1.PodReady, v1.ConditionFalse, 300*time.Second) // NOT READY
	testlib.AwaitPodStatus(t, namespaceName, admin1, v1.PodReady, v1.ConditionTrue, 300*time.Second)
	testlib.AwaitPodStatus(t, namespaceName, admin0, v1.PodReady, v1.ConditionFalse, 300*time.Second) // NOT READY
	testlib.AwaitPodStatus(t, namespaceName, admin0, v1.PodReady, v1.ConditionTrue, 300*time.Second)

	testlib.AwaitAdminFullyConnected(t, namespaceName, admin0, 2)

	testlib.AwaitDatabaseUp(t, namespaceName, admin0, "demo", 2)

	t.Run("verifyAdminState", func(t *testing.T) {
		testlib.VerifyAdminState(t, namespaceName, admin0)
	})

	t.Run("verifyAdminCertificates", func(t *testing.T) {
		certificateInfo, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "--show-json", "get", "certificate-info")
		assert.NilError(t, err)
		verifyAdminCertificates(t, certificateInfo, expectedAdminDN)
	})

	t.Run("expectAllEnginesReconnect", func(t *testing.T) {
		expectedNumberReconnects := 2

		testlib.Await(t, func() bool {
			foundReconnects := 0
			admins := []string{admin0, admin1}
			for _, admin := range admins {
				foundReconnects += testlib.GetStringOccurenceInLog(t, namespaceName, admin, "Reconnected with process with connectKey")
			}

			return foundReconnects == expectedNumberReconnects
		},30*time.Second )
	})

	t.Run("verifyAllProcessesRunning", func(t *testing.T) {
		verifyAllProcessesRunning(t, namespaceName, admin0, 2)
	})

}