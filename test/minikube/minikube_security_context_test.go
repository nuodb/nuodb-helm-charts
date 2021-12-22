// +build long

package minikube

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func checkUser(t *testing.T, namespaceName string, podName string, container string, expectedUid int, expectedGid int, expectedSupplementaryGid int) {
	// check uid
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", container, "--", "id", "-u")
	require.NoError(t, err)
	uid, err := strconv.Atoi(strings.TrimSpace(output))
	require.Equal(t, expectedUid, uid)

	// check gid
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", container, "--", "id", "-g")
	require.NoError(t, err)
	gid, err := strconv.Atoi(strings.TrimSpace(output))
	require.Equal(t, expectedGid, gid)

	// check supplementary gids
	output, err = k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", container, "--", "id", "-G")
	require.NoError(t, err)
	var found bool
	for _, token := range strings.Split(strings.TrimSpace(output), " ") {
		gid, err = strconv.Atoi(token)
		if err == nil && gid == expectedSupplementaryGid {
			found = true
		}
	}
	require.True(t, found, "gid %d not found: %s", expectedSupplementaryGid, output)
}

func checkOwnerGid(t *testing.T, namespaceName string, podName string, container string, filename string, expectedGid int) {
	// check that the specified file is owned by the supplied gid
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "exec", podName, "-c", container, "--", "stat", "-c", "%g", filename)
	require.NoError(t, err)
	gid, err := strconv.Atoi(strings.TrimSpace(output))
	require.Equal(t, expectedGid, gid)
}

func securityContextTest(t *testing.T, adminUid int, adminGid int, adminFsGroup int, databaseUid int, databaseGid, databaseFsGroup int, optionOverrides *map[string]string) {
	// verify that user can run as the supplied uid and that fsGroup becomes
	// owner gid for volumes; this test specifically checks secrets because
	// the hostpath storage-class used by Minikube does not support fsGroup

	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetessecuritycontext-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)

	// create the certs and secrets...
	tlsCommands := []string{
		"export DEFAULT_PASSWORD='" + testlib.SECRET_PASSWORD + "'",
		"setup-keys.sh",
		"nuocmd show certificate --keystore " + testlib.KEYSTORE_FILE + " --store-password \"$DEFAULT_PASSWORD\" > nuoadmin.cert",
	}
	testlib.GenerateTLSConfiguration(t, namespaceName, tlsCommands)

	// create a database with TLS enabled
	options := helm.Options{
		SetValues: map[string]string{
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
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"admin.securityContext.runAsUser":       strconv.Itoa(adminUid),
			"admin.securityContext.fsGroup":         strconv.Itoa(adminFsGroup),
			"database.securityContext.runAsUser":    strconv.Itoa(databaseUid),
			"database.securityContext.fsGroup":      strconv.Itoa(databaseFsGroup),
		},
	}

	// set overrides
	for k, v := range *optionOverrides {
		options.SetValues[k] = v
	}

	// sometimes the test fails because SMs doesn't go ready due to probe
	// timeout
	testlib.OverrideReadinessProbesTimeout(t, &options, "10")

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	adminStatefulSet := fmt.Sprintf("%s-nuodb-cluster0", helmChartReleaseName)
	admin0 := adminStatefulSet+"-0"

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

	smPodNameTemplate := fmt.Sprintf("sm-%s", databaseReleaseName)
	tePodNameTemplate := fmt.Sprintf("te-%s", databaseReleaseName)

	// for all Pods, check that the user has the expected uid:gid and
	// supplementary gid, and check that TLS secrets mounted into the
	// container have the expected ownership (owned by gid fsGroup)
	for _, pod := range testlib.FindAllPodsInSchema(t, namespaceName) {
		if strings.Contains(pod.Name, adminStatefulSet) {
			checkUser(t, namespaceName, pod.Name, "admin", adminUid, adminGid, adminFsGroup)
			checkOwnerGid(t, namespaceName, pod.Name, "admin", "/etc/nuodb/keys/nuocmd.pem", adminFsGroup)
		} else if strings.Contains(pod.Name, smPodNameTemplate) || strings.Contains(pod.Name, tePodNameTemplate) {
			checkUser(t, namespaceName, pod.Name, "engine", databaseUid, databaseGid, databaseFsGroup)
			checkOwnerGid(t, namespaceName, pod.Name, "engine", "/etc/nuodb/keys/nuocmd.pem", databaseFsGroup)
		}
	}
}

func TestSecurityContextEnabled(t *testing.T) {
	// use arbitary uid and fsGroup, which can be different for both charts
	adminUid := 1234
	adminGid := 0
	adminFsGroup := 2000
	databaseUid := 5678
	databaseGid := 0
	databaseFsGroup := 4000
	optionOverrides := map[string]string {
		"admin.securityContext.enabled":    "true",
		"database.securityContext.enabled": "true",
	}
	securityContextTest(t, adminUid, adminGid, adminFsGroup, databaseUid, databaseGid, databaseFsGroup, &optionOverrides)
}

func TestSecurityContextRunAsNonRootGroup(t *testing.T) {
	// users with non-0 gid is supported in 4.3
	testlib.SkipTestOnNuoDBVersionCondition(t, "<4.3")

	// runAsNonRootGroup only supports 1000:1000
	adminUid := 1000
	adminGid := 1000
	adminFsGroup := 1000
	databaseUid := 1000
	databaseGid := 1000
	databaseFsGroup := 1000
	optionOverrides := map[string]string {
		"admin.securityContext.runAsNonRootGroup":    "true",
		"database.securityContext.runAsNonRootGroup": "true",
	}
	securityContextTest(t, adminUid, adminGid, adminFsGroup, databaseUid, databaseGid, databaseFsGroup, &optionOverrides)
}

func TestSecurityContextFsGroupOnly(t *testing.T) {
	// fsGroup omits runAsUser and runAsGroup, but the we expect the image
	// to have a default user of 1000:0
	adminUid := 1000
	adminGid := 0
	adminFsGroup := 2000
	databaseUid := 1000
	databaseGid := 0
	databaseFsGroup := 4000
	optionOverrides := map[string]string {
		"admin.securityContext.fsGroupOnly":    "true",
		"database.securityContext.fsGroupOnly": "true",
	}
	securityContextTest(t, adminUid, adminGid, adminFsGroup, databaseUid, databaseGid, databaseFsGroup, &optionOverrides)
}
