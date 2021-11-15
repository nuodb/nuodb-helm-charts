// +build long

package minikube

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/random"

	corev1 "k8s.io/api/core/v1"
)

func findServiceNodePortE(service *corev1.Service, portName string) (int32, error) {
	for _, port := range service.Spec.Ports {
		if port.Name == portName {
			return port.NodePort, nil
		}
	}
	return 0, errors.New(
		fmt.Sprintf("Unable to find NodePort for service %s", service.Name))
}

func findDomainExternalInfo(t *testing.T, namespaceName string, serviceName string) (string, int32) {
	var domainAddress string
	var domainPort int32
	domainService := testlib.GetService(t, namespaceName, serviceName)
	if domainService.Spec.Type == corev1.ServiceTypeNodePort {
		// get the internal address of the first Kubernetes node
		for _, address := range testlib.GetNodesInternalAddresses(t) {
			domainAddress = address
			break
		}
		var err error
		// get the nodePort of the domain service
		domainPort, err = findServiceNodePortE(domainService, "48004-tcp")
		require.NoError(t, err)
	} else {
		domainAddress = domainService.Status.LoadBalancer.Ingress[0].IP
		domainPort = 48004
	}
	return domainAddress, domainPort
}

func verifyNuoSQLEngine(t *testing.T, address string, port int32, databaseName string,
	properties map[string]string, expectedNodeIds []int32) {

	pathToNuoSQL, ok := os.LookupEnv("NUOSQL_PATH")
	if !ok {
		pathToNuoSQL = "nuosql"
	}
	cmd := fmt.Sprintf("echo \"select GETNODEID() as NODEID from dual;\" | %s %s@%s:%d --user dba --password secret --vertical-output",
		pathToNuoSQL, databaseName, address, port)
	for key, value := range properties {
		cmd += fmt.Sprintf(" --connection-property %s='%s'", key, value)
	}
	t.Logf("Running command: <%s>", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if exiterr, ok := err.(*exec.ExitError); ok {
		// nuosql has exited with an non zero exit code; generate better error
		// message by including the command stderr
		require.NoError(t, err, string(exiterr.Stderr))
	} else {
		require.NoError(t, err)
	}

	match := regexp.MustCompile("NODEID: ([0-9]+)").FindStringSubmatch(string(out))
	require.NotNil(t, match, out)
	actualNodeId, err := strconv.Atoi(match[1])
	require.NoError(t, err)
	require.Contains(t, expectedNodeIds, int32(actualNodeId), "Unexpected TE nodeId")
}

func verifyService(t *testing.T, namespaceName string, podName string, serviceName string,
	expectedServiceType corev1.ServiceType, ping bool) *corev1.Service {

	service := testlib.GetService(t, namespaceName, serviceName)
	require.Equal(t, expectedServiceType, service.Spec.Type)
	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
	return service
}

func verifyProcessExternalAccessLabels(t *testing.T, namespaceName string, adminPod string,
	databaseName string, services map[string]*corev1.Service) {

	processes, err := testlib.GetDatabaseProcessesE(t, namespaceName, adminPod, databaseName)
	require.NoError(t, err)

	findServiceForProcessE := func(process *testlib.NuoDBProcess) (*corev1.Service, error) {
		for group, service := range services {
			if strings.Contains(process.Hostname, group) {
				return service, nil
			}
		}
		return nil, errors.New(
			fmt.Sprintf("Unable to find corresponding service for process %s", process.Hostname))
	}

	for _, process := range processes {
		if process.Type == "TE" {
			t.Logf("Validating process labels for TE hostname=%s", process.Hostname)
			s, err := findServiceForProcessE(&process)
			require.NoError(t, err)
			if s.Spec.Type == corev1.ServiceTypeLoadBalancer {
				// verify that external-address is set correctly
				val, ok := process.Labels["external-address"]
				require.True(t, ok)
				require.Equal(t, s.Status.LoadBalancer.Ingress[0].IP, val)
			} else {
				// verify that external-address is set
				_, ok := process.Labels["external-address"]
				require.True(t, ok)
				// verify that external-port is set correctly
				nodePort, err := findServiceNodePortE(s, "48006-tcp")
				require.NoError(t, err)
				val, ok := process.Labels["external-port"]
				require.True(t, ok)
				require.Equal(t, strconv.Itoa(int(nodePort)), val)
			}
		}
	}

}

func TestKubernetesMultipleTEGroups(t *testing.T) {
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 4.2.4")
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	deployMultipleTEGroups := func(t *testing.T, serviceType corev1.ServiceType) {
		adminOptions := helm.Options{
			SetValues: map[string]string{
				"admin.externalAccess.enabled": "true",
				"admin.externalAccess.type":    fmt.Sprintf("%s", serviceType),
			},
		}

		defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

		randomSuffix := strings.ToLower(random.UniqueId())
		namespaceName := fmt.Sprintf("%skubernetesmultipletegroups-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
		testlib.CreateNamespace(t, namespaceName)

		helmChartReleaseName, _ := testlib.StartAdmin(t, &adminOptions, 1, namespaceName)

		admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
		var serviceSuffix string
		if serviceType == corev1.ServiceTypeNodePort {
			serviceSuffix = "nodeport"
		} else {
			serviceSuffix = "balancer"
		}

		verifyService(t, namespaceName, admin0, "nuodb", "ClusterIP", true)
		verifyService(t, namespaceName, admin0, "nuodb-clusterip", "ClusterIP", false)
		verifyService(t, namespaceName, admin0, fmt.Sprintf("nuodb-%s", serviceSuffix), serviceType, false)

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":              testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.sm.resources.requests.memory":           testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":              testlib.MINIMAL_VIABLE_ENGINE_CPU,
				"database.te.resources.requests.memory":           testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.externalAccess.enabled":              "true",
				"database.te.externalAccess.type":                 fmt.Sprintf("%s", serviceType),
				"database.te.otherOptions.enable-external-access": "true",
			},
			// this is needed to define the NODE_IP environment variable
			ValuesFiles: []string{"../files/database-env.yaml"},
		}

		if serviceType == corev1.ServiceTypeNodePort {
			// set TE external-address to the Kubernetes node IP; in the real
			// scenario this will be a load balancer on top of all cluster nodes
			databaseOptions.SetValues["database.te.labels.external-address"] = "$(NODE_IP)"
		}

		// install the primary Helm release
		databaseOptions.SetValues["database.te.labels.tx-type"] = "OLTP"
		group1ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)
		// all other Helm releases will be secondary releases with TEs only
		databaseOptions.SetValues["database.primaryRelease"] = "false"
		databaseOptions.SetValues["database.sm.hotCopy.replicas"] = "0"
		databaseOptions.SetValues["database.sm.noHotCopy.replicas"] = "0"
		databaseOptions.SetValues["database.te.labels.tx-type"] = "TAP"
		group2ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)
		opt := testlib.GetExtractedOptions(&databaseOptions)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, 3)

		verifyService(t, namespaceName, admin0, "demo", "ClusterIP", true)
		verifyService(t, namespaceName, admin0, "demo-clusterip", "ClusterIP", false)

		// get balancer service per database Helm release
		balancerServices := make(map[string]*corev1.Service)
		for _, group := range []string{group1ReleaseName, group2ReleaseName} {
			// each group installs unique balancer that matches TEs in the group only
			balancerServiceName := fmt.Sprintf("%s-nuodb-%s-%s-%s", group, opt.ClusterName,
				opt.DbName, serviceSuffix)
			s := verifyService(t, namespaceName, admin0, balancerServiceName, serviceType, false)
			balancerServices[group] = s
		}

		// verify external-address and external-port labels for TE processes
		verifyProcessExternalAccessLabels(t, namespaceName, admin0, opt.DbName, balancerServices)

		address, port := findDomainExternalInfo(t, namespaceName, fmt.Sprintf("nuodb-%s", serviceSuffix))

		for txType, group := range map[string]string{"OLTP": group1ReleaseName, "TAP": group2ReleaseName} {
			var expectedNodeIds []int32
			processes, err := testlib.GetDatabaseProcessesE(t, namespaceName, admin0, opt.DbName)
			require.NoError(t, err)
			for _, process := range processes {
				if strings.Contains(process.Hostname, group) && process.Type == "TE" {
					expectedNodeIds = append(expectedNodeIds, process.NodeId)
				}
			}
			// verify that expected TE engine is returned by the admin
			verifyNuoSQLEngine(t, address, port, opt.DbName,
				map[string]string{"LBQuery": fmt.Sprintf("random(label(tx-type %s))", txType)},
				expectedNodeIds)
		}
	}

	t.Run("testServiceLoadBalancer", func(t *testing.T) { deployMultipleTEGroups(t, corev1.ServiceTypeLoadBalancer) })
	t.Run("testServiceNodePort", func(t *testing.T) { deployMultipleTEGroups(t, corev1.ServiceTypeNodePort) })
}
