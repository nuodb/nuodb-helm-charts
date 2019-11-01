// +build long

package minikube

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"encoding/json"
	"strconv"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"gotest.tools/assert"
)

func unmarshalCertificateInfo(t *testing.T, certificateInfoJSON string) map[string]interface{} {
	var results map[string]interface{}
	json.Unmarshal([]byte(certificateInfoJSON), &results)
	return results
}

func verifyAdminCeritificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := unmarshalCertificateInfo(t, certificateInfoJSON)
	for _, value := range certificateInfo["serverCertificates"].(map[string]interface{}) {
		certSubjectName := value.(map[string]interface{})["subjectName"].(string)
		assert.Assert(t, strings.Contains(certSubjectName, expectedDN),
		"`%s` not found in:\n %s", expectedDN, certSubjectName)
	}
}

func verifyEngineCeritificates(t *testing.T, certificateInfoJSON string, expectedDN string) {
	certificateInfo := unmarshalCertificateInfo(t, certificateInfoJSON)
	for _, value := range certificateInfo["processCertificates"].(map[string]interface{}) {
		certIssuerName := value.(map[string]interface{})["issuerName"].(string)
		assert.Assert(t, strings.Contains(certIssuerName, expectedDN),
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

func TestKubernetesTLSRotation(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	namespaceName := fmt.Sprintf("test-tls-rotation-%s", randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "")
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	kubectlOptions.Namespace = namespaceName

	initialTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}

	// As nuodocker/nuoadmin wrapper is using peer insead of initialMembership, 
	//   we need to use persistence for admin Raft logs during the rolling upgrade.
	options := helm.Options{
		SetValues: map[string]string{
			"admin.persistence.enabled":             "true",
			"admin.replicas":                        "3",
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
			"database.sm.resources.requests.cpu":    "500m",
			"database.sm.resources.requests.memory": "500Mi",
			"database.te.resources.requests.cpu":    "500m",
			"database.te.resources.requests.memory": "500Mi",
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

	upgradedOptions := options
	upgradedOptions.SetValues = make(map[string]string, len(options.SetValues))
	for k, v := range options.SetValues {
		upgradedOptions.SetValues[k] = v
	}
	upgradedOptions.SetValues["admin.tlsCACert.secret"] = testlib.CA_CERT_SECRET_NEW
	upgradedOptions.SetValues["admin.tlsKeyStore.secret"] = testlib.KEYSTORE_SECRET_NEW

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	certGeneratorPodName, adminReleaseName, databaseReleaseName := startDomainWithTLSCertificates(t, &options, namespaceName, initialTLSCommands)
	
	// create the new certs...
	testlib.GenerateCustomCertificates(t, certGeneratorPodName, namespaceName, newTLSCommands)
	newTLSKeysLocation := testlib.CopyCertificatesToControlHost(t, certGeneratorPodName, namespaceName)
	testlib.CreateSecret(t, namespaceName, testlib.CA_CERT_FILE, testlib.CA_CERT_SECRET_NEW, newTLSKeysLocation)
	testlib.CreateSecretWithPassword(t, namespaceName, testlib.KEYSTORE_FILE, testlib.KEYSTORE_SECRET_NEW, testlib.SECRET_PASSWORD, newTLSKeysLocation)
	
	testlib.RotateTLSCertificates(t, &upgradedOptions, kubectlOptions, adminReleaseName, databaseReleaseName, newTLSKeysLocation)
	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)

	certificateInfo, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "--show-json", "get", "certificate-info")
	assert.NilError(t, err)
	
	t.Run("verifyAdminCeritificates", func(t *testing.T) {
		verifyAdminCeritificates(t, certificateInfo, expectedAdminDN)
	})

	t.Run("verifyEngineCeritificates", func(t *testing.T) {
		verifyEngineCeritificates(t, certificateInfo, expectedAdminDN)
	})
}
