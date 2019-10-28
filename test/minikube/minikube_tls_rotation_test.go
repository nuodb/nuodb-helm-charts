package minikube

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestKubernetesTLSRotation(t *testing.T) {
	testlib.AwaitTillerUp(t)

	randomSuffix := strings.ToLower(random.UniqueId())

	namespaceName := fmt.Sprintf("test-tls-rotation-%s", randomSuffix)
	kubectlOptions := k8s.NewKubectlOptions("", "")
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	kubectlOptions.Namespace = namespaceName

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)

	initialTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
	}

	options := helm.Options{
		SetValues: map[string]string{
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
			"database.sm.resources.requests.memory": "500m",
			"database.te.resources.requests.cpu":    "500m",
			"database.te.resources.requests.memory": "500m",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	newTLSCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"nuocmd create keypair --keystore ca.p12 --store-password \"$DEFAULT_PASSWORD\" --dname \"CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=123456\" --validity 36500 --ca",
		"nuocmd create keypair --keystore " + testlib.KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" --dname \"CN=nuoadmin.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890\"",
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

	adminReleaseName, _ := testlib.RotateTLSCertificates(t, &options, &upgradedOptions, kubectlOptions, initialTLSCommands, newTLSCommands)
	admin0 := fmt.Sprintf("%s-nuodb-0", adminReleaseName)

	output, _ := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "exec", admin0, "--", "nuocmd", "--show-json", "get", "certificate-info")
	fmt.Println(output)
}
