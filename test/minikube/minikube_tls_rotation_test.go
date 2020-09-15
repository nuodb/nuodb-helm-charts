// +build long

package minikube

import (
	"fmt"
	v12 "k8s.io/api/core/v1"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func verifyAdminCertificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := testlib.UnmarshalJSONObject(t, certificateInfoJSON)
	for _, value := range certificateInfo["serverCertificates"].(map[string]interface{}) {
		certSubjectName := value.(map[string]interface{})["subjectName"].(string)
		assert.Contains(t, certSubjectName, expectedDN,
			"`%s` not found in:\n %s", expectedDN, certSubjectName)
	}
}

func verifyEngineCertificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := testlib.UnmarshalJSONObject(t, certificateInfoJSON)
	for _, value := range certificateInfo["processCertificates"].(map[string]interface{}) {
		certIssuerName := value.(map[string]interface{})["issuerName"].(string)
		assert.Contains(t, certIssuerName, expectedDN,
			"`%s` not found in:\n %s", expectedDN, certIssuerName)
	}
}

func startDomainWithTLSCertificates(t *testing.T, options *helm.Options, namespaceName string, tlsCommands []string) (string, string, string) {
	adminReplicaCount, err := strconv.Atoi(options.SetValues["admin.replicas"])
	assert.NoError(t, err, "Unable to find/convert admin.replicas value")

	// create initial certs...
	certGeneratorPodName, _ := testlib.GenerateTLSConfiguration(t, namespaceName, tlsCommands)

	adminReleaseName, namespaceName := testlib.StartAdmin(t, options, adminReplicaCount, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)

	return certGeneratorPodName, adminReleaseName, databaseReleaseName
}

func TestKubernetesTLSRotation(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("testtlsrotation-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	initialTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}

	// As nuodocker/nuoadmin wrapper is using peer instead of initialMembership,
	// we need to use persistence for admin Raft logs during the rolling upgrade.
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

	certGeneratorPodName, adminReleaseName, databaseReleaseName := startDomainWithTLSCertificates(t, &options, namespaceName, initialTLSCommands)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)
	admin1 := fmt.Sprintf("%s-nuodb-cluster0-1", adminReleaseName)

	// get the OLD log
	go testlib.GetAppLog(t, namespaceName, admin0, "-previous", &v12.PodLogOptions{Follow: true})
	go testlib.GetAppLog(t, namespaceName, admin1, "-previous", &v12.PodLogOptions{Follow: true})

	// create the new certs...
	testlib.GenerateCustomCertificates(t, certGeneratorPodName, namespaceName, newTLSCommands)
	newTLSKeysLocation := testlib.CopyCertificatesToControlHost(t, certGeneratorPodName, namespaceName)
	testlib.CreateSecret(t, namespaceName, testlib.CA_CERT_FILE, testlib.CA_CERT_SECRET, newTLSKeysLocation)
	testlib.CreateSecretWithPassword(t, namespaceName, testlib.KEYSTORE_FILE, testlib.KEYSTORE_SECRET, testlib.SECRET_PASSWORD, newTLSKeysLocation)

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	testlib.AddDiagnosticTeardown(testlib.TEARDOWN_ADMIN, t, func() {
		k8s.RunKubectl(t, kubectlOptions, "get", "pods", "-o", "wide")
	})

	testlib.RotateTLSCertificates(t, &options, namespaceName, adminReleaseName, databaseReleaseName, newTLSKeysLocation, false)

	certificateInfo, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "--show-json", "get", "certificate-info")
	assert.NoError(t, err)

	t.Run("verifyAdminCertificates", func(t *testing.T) {
		verifyAdminCertificates(t, certificateInfo, expectedAdminDN)
	})

	t.Run("verifyEngineCertificates", func(t *testing.T) {
		verifyEngineCertificates(t, certificateInfo, expectedAdminDN)
	})
}
