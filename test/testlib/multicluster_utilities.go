package testlib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const CONTEXT_CLUSTER_KEY = CONTEXT_KEY("cluster")

var MULTI_CLUSTER_1 = K8sCluster{
	Name:    "cluster-nuodb-1",
	Domain:  "cluster1.local",
	Context: "yin",
}

var MULTI_CLUSTER_2 = K8sCluster{
	Name:    "cluster-nuodb-2",
	Domain:  "cluster2.local",
	Context: "yang",
}

type CONTEXT_KEY string

// Global variable that holds all multi-cluster NuoDB deployments
var clusterDeployments = make(map[string]*ClusterDeploymentContext)

type K8sCluster struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Context string `json:"context"`
}

type ClusterDeploymentContext struct {
	Options             *helm.Options
	ThisCluster         K8sCluster
	EntrypointCluster   K8sCluster
	AdminReleaseName    string
	DatabaseReleaseName string
	Namespace           string
}

func UnmarshalClusters(s string) (err error, clusters map[string]K8sCluster) {
	dec := json.NewDecoder(strings.NewReader(s))
	clusters = make(map[string]K8sCluster)

	for {
		var obj K8sCluster
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, clusters
		}

		if err != nil {
			return
		}

		clusters[obj.Name] = obj
	}
}

func CopyMap(m map[string]string) map[string]string {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

/**
 * Injects user-defined Kubernetes cluster information
 *
 * The information is expected to be in stored in
 * clustersInject.yaml file in JSON format as a list.
 * See all cluster names in MULTI_CLUSTER_* constants.
 *
 */
func InjectClusters(t *testing.T, cluster K8sCluster) K8sCluster {
	dat, err := os.ReadFile(INJECT_CLUSTERS_FILE)
	if err != nil {
		return cluster
	}

	t.Log("Using injected clusters:\n", string(dat))

	err, clustersMap := UnmarshalClusters(string(dat))
	require.NoError(t, err)

	if val, ok := clustersMap[cluster.Name]; ok {
		// The cluster information is overwritten by injected data
		return val
	}
	return cluster
}

/**
 * Creates cluster deployment context
 *
 * New cluster deployment context will be created and attached to
 * provided parent context. The new context will be stored in
 * clusterDeployments global map so that it can be used later if needed.
 *
 * The typical idiom is:
 * <pre>
 *   defer testlib.Teardown(testlib.TEARDOWN_MULTICLUSTER)
 *
 *   context := context.Background()
 *   deploymentContext := testlib.NewClusterDeploymentContext(context,
 *	   &helm.Options{},
 *     testlib.MULTI_CLUSTER_1,
 *     testlib.MULTI_CLUSTER_1)
 * <pre>
 */
func NewClusterDeploymentContext(parent context.Context, options *helm.Options, thisCluster K8sCluster, entrypointCluster K8sCluster) context.Context {
	deploymentContext := &ClusterDeploymentContext{
		ThisCluster:       thisCluster,
		EntrypointCluster: entrypointCluster,
		Options:           options,
	}
	AddTeardown(TEARDOWN_MULTICLUSTER, func() { delete(clusterDeployments, thisCluster.Name) })
	clusterDeployments[thisCluster.Name] = deploymentContext
	return context.WithValue(parent, CONTEXT_CLUSTER_KEY, deploymentContext)
}

/**
 * Change kubectl current context
 *
 * A change and revert functions are returned so that
 * the caller can use them.
 *
 */
func ChangeCluster(t *testing.T, cluster K8sCluster) (func(), func()) {
	clusterContext := cluster.Context
	kubectlOptions := k8s.NewKubectlOptions("", "", "")
	current, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "config", "current-context")
	require.NoError(t, err, "Unable to get current context")
	changeBackFunc := func() { k8s.RunKubectl(t, kubectlOptions, "config", "use-context", current) }
	changeFunc := func() { k8s.RunKubectl(t, kubectlOptions, "config", "use-context", clusterContext) }
	changeFunc()
	return changeFunc, changeBackFunc
}

func GetClusterDeploymentContext(t *testing.T, context context.Context) *ClusterDeploymentContext {
	value := context.Value(CONTEXT_CLUSTER_KEY)
	require.NotNil(t, value, "Unable to retrieve cluster deployment context")
	deploymentContext, ok := value.(*ClusterDeploymentContext)
	require.True(t, ok, "Unexpected type for cluster deployment context")
	return deploymentContext
}

/**
 * Deploy NuoDB Helm charts using deployment context
 *
 * This function is typically used when performing a multi-cluster deployment.
 *
 * The typical idiom is:
 * <pre>
 *	testlib.DeployWithContext(t,
 *		deploymentContext,
 *		func(context *testlib.ClusterDeploymentContext, options *helm.Options) {
 *			testlib.CreateNamespace(t, namespaceName)
 *			testlib.StartAdminCustomRelease(t, options, 1, namespaceName, adminReleaseName)
 *			admin0 := fmt.Sprintf("%s-nuodb-%s-0", adminReleaseName, context.ThisCluster.Name)
 *			databaseReleaseName := testlib.StartDatabase(t, namespaceName, admin0, options)
 *			// Store deployment details in the context
 *			context.AdminReleaseName = adminReleaseName
 *			context.DatabaseReleaseName = databaseReleaseName
 *			context.Namespace = namespaceName
 *		},
 *	)
 * <pre>
 */
func DeployWithContext(t *testing.T, context context.Context, deployFunc func(context *ClusterDeploymentContext, options *helm.Options)) {
	// Retrieve cluster deployment information from the context
	deploymentContext := GetClusterDeploymentContext(t, context)
	thisCluster := deploymentContext.ThisCluster
	entrypointCluster := deploymentContext.EntrypointCluster
	thisCluster = InjectClusters(t, thisCluster)
	entrypointCluster = InjectClusters(t, entrypointCluster)
	options := deploymentContext.Options

	// Set all variables needed by NuoDB Helm Charts in multi-cluster mode
	optionsCopy := &helm.Options{
		SetValues:      CopyMap(options.SetValues),
		KubectlOptions: options.KubectlOptions,
	}
	optionsCopy.SetValues["cloud.cluster.name"] = thisCluster.Name
	optionsCopy.SetValues["cloud.cluster.entrypointName"] = entrypointCluster.Name
	optionsCopy.SetValues["cloud.cluster.domain"] = thisCluster.Domain
	optionsCopy.SetValues["cloud.cluster.entrypointDomain"] = entrypointCluster.Domain

	// Change the current context to this cluster
	changeClusterFunc, changeClusterBackFunc := ChangeCluster(t, thisCluster)
	// Change the cluster before and after any additional teardown steps are executed
	AddGlobalTeardown(changeClusterBackFunc)
	AddGlobalDiagnosticTeardown(true, changeClusterBackFunc)
	defer AddGlobalTeardown(changeClusterFunc)
	defer AddGlobalDiagnosticTeardown(true, changeClusterFunc)

	// Execute provided actions in context of this cluster
	deployFunc(deploymentContext, optionsCopy)
	changeClusterBackFunc()
}

/**
 * Execute arbitrary actions in context of a cluster deployment
 *
 * The deployment context will be passed in the custom function as parameter
 * which can be used to obtain information about the deployment.
 *
 */
func WithClusterDeployment(t *testing.T, context *ClusterDeploymentContext, actionsFunc func(context *ClusterDeploymentContext)) {
	cluster := context.ThisCluster
	cluster = InjectClusters(t, cluster)
	// Change the current context to this cluster
	_, changeClusterBackFunc := ChangeCluster(t, cluster)
	actionsFunc(context)
	changeClusterBackFunc()
}

/**
 * Execute arbitrary actions on all known cluster deployments
 *
 * The cluster deployment is stored in clusterDeployments variable when
 * a deployment context for it is created.
 *
 */
func ExecuteInAllClusters(t *testing.T, actionFunc func(context *ClusterDeploymentContext)) {
	for _, deploymentContext := range clusterDeployments {
		WithClusterDeployment(t, deploymentContext, actionFunc)
	}
}

func AdjustPodTimeout(podName string, timeout time.Duration) time.Duration {
	if strings.Contains(podName, MULTI_CLUSTER_1.Name) || strings.Contains(podName, MULTI_CLUSTER_2.Name) {
		// in Azure sometimes it takes 3 to 5 min for the disks to be
		// provisioned and attached to the nodes
		return time.Duration(int64(1.5 * math.Max(float64(timeout), float64(300*time.Second))))
	}
	return timeout
}

const (
	DNS_CONFIG_MARKER = "# BEGIN BASE CONFIG"
	COREDNS_NS        = "kube-system"
	COREDNS_CM        = "coredns"
	COREFILE_KEY      = "Corefile"
	DNS_SERVICE       = "dns-service"
)

func getClusterDnsServer(t *testing.T, clusterContext context.Context) string {
	var ret string
	DeployWithContext(t,
		clusterContext,
		func(context *ClusterDeploymentContext, options *helm.Options) {
			// Copy kubectl options and update namespace
			require.NotNil(t, options.KubectlOptions)
			kubectlOptions := *options.KubectlOptions
			kubectlOptions.Namespace = COREDNS_NS
			// Get endpoint for DNS server
			out, err := k8s.RunKubectlAndGetOutputE(t, &kubectlOptions,
				"get", "service", DNS_SERVICE,
				"-o", "jsonpath={.status.loadBalancer.ingress[0].hostname}")
			require.NoError(t, err)
			ret = strings.TrimSpace(out)
		},
	)
	return ret
}

func getDnsServerIps(t *testing.T, clusterContext context.Context) []string {
	host := getClusterDnsServer(t, clusterContext)
	if host == "" {
		return nil
	}
	addresses, err := net.LookupHost(host)
	require.NoError(t, err)
	sort.Strings(addresses)
	return addresses
}

func getDnsConfigSnippet(t *testing.T, context context.Context) string {
	dnsServers := getDnsServerIps(t, context)
	if dnsServers == nil {
		return ""
	}
	deploymentContext := GetClusterDeploymentContext(t, context)
	return fmt.Sprintf(`%s:53 {
    errors
    cache 30
    forward . %s {
      force_tcp
    }
}`, deploymentContext.ThisCluster.Domain, strings.Join(dnsServers, " "))
}

func updateDnsConfig(t *testing.T, ctx context.Context, kubectlOptions *k8s.KubectlOptions, dnsConfigSnippet string) {
	// Create K8s client and get configmap for CoreDNS
	clientset, err := GetKubernetesClientFromOptionsE(t, kubectlOptions)
	require.NoError(t, err, "Unable to create K8s client")
	cm, err := clientset.CoreV1().ConfigMaps(COREDNS_NS).Get(ctx, COREDNS_CM, metav1.GetOptions{})
	require.NoError(t, err, "Unable to get CoreDNS configmap")
	config, ok := cm.Data[COREFILE_KEY]
	require.True(t, ok, "Did not find key %s in CoreDNS configmap", COREFILE_KEY)

	// Find marker string separating generated config from base config
	idx := strings.Index(config, DNS_CONFIG_MARKER)
	require.GreaterOrEqual(t, idx, 0, "Did not find marker string in CoreDNS config")

	// Update configmap if generated config changed
	updatedConfig := dnsConfigSnippet + "\n" + config[idx:]
	if updatedConfig == config {
		t.Logf("DNS configuration up to date")
		return
	}
	t.Logf("Adding DNS config snippet %s", dnsConfigSnippet)
	cm.Data[COREFILE_KEY] = updatedConfig
	_, err = clientset.CoreV1().ConfigMaps(COREDNS_NS).Update(ctx, cm, metav1.UpdateOptions{})
	require.NoError(t, err)
}

// UpdateDnsConfig adds the DNS server from one cluster as an upstream resolver
// for the other, replacing any stale IP addresses that may currently be appear
// in the CoreDNS configuration.
func UpdateDnsConfig(t *testing.T, fromCtx context.Context, toCtx context.Context) {
	dnsConfigSnippet := getDnsConfigSnippet(t, fromCtx)
	if dnsConfigSnippet != "" {
		DeployWithContext(t,
			toCtx,
			func(context *ClusterDeploymentContext, options *helm.Options) {
				updateDnsConfig(t, toCtx, options.KubectlOptions, dnsConfigSnippet)
			},
		)
	}
}
