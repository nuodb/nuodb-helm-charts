// +build long

package minikube

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FILE_PLUGIN_CONFIGMAP = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: nuocollector-insights-file
  labels:
    nuodb.com/nuocollector-plugin: insights
data:
  file.conf: |-
    [[outputs.file]]
    files = ["stdout"]
    data_format = "influx"
`

func createOutputFilePlugin(t *testing.T, namespaceName string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.KubectlApplyFromString(t, kubectlOptions, FILE_PLUGIN_CONFIGMAP)
	testlib.AddTeardown(testlib.TEARDOWN_INSIGHTS, func() { k8s.KubectlDeleteFromStringE(t, kubectlOptions, FILE_PLUGIN_CONFIGMAP) })
	// Wait a bit so that the insights-config sidecar refresh the Telegraf configuration
	time.Sleep(3 * time.Second)
}

func checkInsightsMetricsLine(t *testing.T, namespaceName string, podName string,
	expectedLine string, minOccurances int) bool {
	count := testlib.GetRegexOccurrenceInLog(t, namespaceName, podName, expectedLine, &v12.PodLogOptions{Container: "insights"})
	if count >= minOccurances {
		t.Logf("Found %d occurances of '%s' in pod %s log", count, expectedLine, podName)
		return true
	}
	return false
}

func verifyCollectionForAdmin(t *testing.T, namespaceName string, app string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	pods := k8s.ListPods(t, options, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s,component=admin", app),
	})
	for _, pod := range pods {
		// Verfy that logfile nuocollector input plugin is collecting and outputing data to stdout
		message := "logfile,db_tag=nuolog,host=" + pod.Name
		t.Logf("Searching string '%s' in pod %s logs", message, pod.Name)
		testlib.Await(t, func() bool {
			return checkInsightsMetricsLine(t, namespaceName, pod.Name, message, 1)
		}, 60*time.Second)
	}
}

func verifyCollectionForDatabase(t *testing.T, namespaceName string, app string, dbName string) {
	options := k8s.NewKubectlOptions("", "", namespaceName)
	pods := k8s.ListPods(t, options, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s,database=%s,component in (sm, te)", app, dbName),
	})
	for _, pod := range pods {
		// Verfy that monitor nuocollector input plugin is collecting and outputing data to stdout
		message := fmt.Sprintf("[^,]+,db=%s,db_tag=nuodb,host=%s", dbName, pod.Name)
		t.Logf("Searching string '%s' in pod %s logs", message, pod.Name)
		testlib.Await(t, func() bool {
			return checkInsightsMetricsLine(t, namespaceName, pod.Name, message, 1)
		}, 600*time.Second)
		// Verfy that nuodb_thread nuocollector input plugin is collecting and outputing data to stdout
		message = "nuodb_thread,db_tag=nuodb_internal,exe=[^,]+,host=" + pod.Name
		t.Logf("Searching string '%s' in pod %s logs", message, pod.Name)
		testlib.Await(t, func() bool {
			return checkInsightsMetricsLine(t, namespaceName, pod.Name, message, 1)
		}, 60*time.Second)
		// Verfy that nuodb_msgtrace nuocollector input plugin is collecting and outputing data to stdout
		message = fmt.Sprintf("nuodb_msgtrace,db_tag=nuodb_internal,dbname=%s,host=%s", dbName, pod.Name)
		t.Logf("Searching string '%s' in pod %s logs", message, pod.Name)
		testlib.Await(t, func() bool {
			return checkInsightsMetricsLine(t, namespaceName, pod.Name, message, 1)
		}, 60*time.Second)
		// Verfy that nuodb_synctrace nuocollector input plugin is collecting and outputing data to stdout
		message = fmt.Sprintf("nuodb_synctrace,db_tag=nuodb_internal,dbname=%s,host=%s", dbName, pod.Name)
		t.Logf("Searching string '%s' in pod %s logs", message, pod.Name)
		testlib.Await(t, func() bool {
			return checkInsightsMetricsLine(t, namespaceName, pod.Name, message, 1)
		}, 120*time.Second)
	}
}

func TestInsightsMetricsCollection(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{
			"insights.enabled":                      "true",
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			// Load custom user plugin for the admin
			"insights.plugins.admin.log": `
[[inputs.tail]]
  files = ["/var/log/nuodb/nuoadmin*.log"]
  from_beginning = true
  name_override= "logfile"
  data_format = "grok"
  grok_patterns= [ "%{CUSTOM_LOGLINE}" ]
  grok_custom_patterns = '''
  CUSTOM_LOGLINE %{TIMESTAMP_ISO8601:timestamp:ts-"2006-01-02T15:04:05.000-0700"}%{SPACE}(?:%{LOGLEVEL:loglevel:tag}%{SPACE}(?:%{NOTSPACE:logger:tag}%{SPACE})?)?%{GREEDYDATA:message}
  '''
  [inputs.tail.tags]
    db_tag = "nuolog"`,
		},
	}

	// DB-32319: Enable TLS because NuoDB 4.2+ image doesn't contain pregenerated keys
	// and it looks like nuocollector have problem with downgrading gracefully to HTTP
	// Remove this once DB-32319 is fixed
	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("insightsmetricscollection-%s", randomSuffix)
	testlib.CreateNamespace(t, namespaceName)
	testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	adminReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)

	t.Run("startDatabaseStatefulSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		databaseRealeaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)
		createOutputFilePlugin(t, namespaceName)
		defer testlib.Teardown(testlib.TEARDOWN_INSIGHTS)
		testlib.PopulateDBWithQuickstart(t, namespaceName, admin0, "demo")
		t.Run("verifyMetricsCollection", func(t *testing.T) {
			verifyCollectionForAdmin(t, namespaceName, fmt.Sprintf("%s-nuodb-cluster0", adminReleaseName))
			verifyCollectionForDatabase(t, namespaceName, fmt.Sprintf("%s-nuodb-%s-%s", databaseRealeaseName, "cluster0", "demo"), "demo")
		})
	})

	t.Run("startDatabaseDaemonSet", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits
		options.SetValues["database.enableDaemonSet"] = "true"
		// Start only hotcopy SM daemonset
		options.SetValues["database.sm.noHotCopy.enablePod"] = "false"
		databaseRealeaseName := testlib.StartDatabase(t, namespaceName, admin0, &options)

		createOutputFilePlugin(t, namespaceName)
		defer testlib.Teardown(testlib.TEARDOWN_INSIGHTS)
		testlib.PopulateDBWithQuickstart(t, namespaceName, admin0, "demo")
		t.Run("verifyMetricsCollection", func(t *testing.T) {
			verifyCollectionForAdmin(t, namespaceName, fmt.Sprintf("%s-nuodb-cluster0", adminReleaseName))
			verifyCollectionForDatabase(t, namespaceName, fmt.Sprintf("%s-nuodb-%s-%s", databaseRealeaseName, "cluster0", "demo"), "demo")
		})
	})

	t.Run("startDatabaseMultiTenant", func(t *testing.T) {
		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		greenDatabaseRealeaseName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"insights.enabled":                      "true",
				"database.name":                         "green",
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
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
			},
		})

		blueDatabaseRealeaseName := testlib.StartDatabase(t, namespaceName, admin0, &helm.Options{
			SetValues: map[string]string{
				"insights.enabled":                      "true",
				"database.name":                         "blue",
				"database.sm.resources.requests.cpu":    "250m",
				"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":    "250m",
				"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
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
			},
		})

		createOutputFilePlugin(t, namespaceName)
		defer testlib.Teardown(testlib.TEARDOWN_INSIGHTS)

		t.Run("verifyMetricsCollectionForAdmin", func(t *testing.T) {
			verifyCollectionForAdmin(t, namespaceName, fmt.Sprintf("%s-nuodb-cluster0", adminReleaseName))
		})
		testlib.PopulateDBWithQuickstart(t, namespaceName, admin0, "green")
		t.Run("verifyMetricsCollectionForGreen", func(t *testing.T) {
			verifyCollectionForDatabase(t, namespaceName, fmt.Sprintf("%s-nuodb-%s-%s", greenDatabaseRealeaseName, "cluster0", "green"), "green")
		})
		testlib.PopulateDBWithQuickstart(t, namespaceName, admin0, "blue")
		t.Run("verifyMetricsCollectionForBlue", func(t *testing.T) {
			verifyCollectionForDatabase(t, namespaceName, fmt.Sprintf("%s-nuodb-%s-%s", blueDatabaseRealeaseName, "cluster0", "blue"), "blue")
		})
	})

}