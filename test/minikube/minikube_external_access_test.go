//go:build long
// +build long

package minikube

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"

	"github.com/Masterminds/semver"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/random"

	corev1 "k8s.io/api/core/v1"
)

func findServicePortE(service *corev1.Service, portName string) (*corev1.ServicePort, error) {
	for _, port := range service.Spec.Ports {
		if port.Name == portName {
			return &port, nil
		}
	}
	return nil, errors.New(
		fmt.Sprintf("Unable to find Port with name %s for service %s", portName, service.Name))
}

func findDomainExternalInfo(t *testing.T, namespaceName string, serviceName string, portName string) (string, int32) {
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
		servicePort, err := findServicePortE(domainService, portName)
		require.NoError(t, err)
		domainPort = servicePort.NodePort
	} else {
		domainAddress = domainService.Status.LoadBalancer.Ingress[0].IP
		servicePort, err := findServicePortE(domainService, portName)
		require.NoError(t, err)
		domainPort = servicePort.Port
	}
	return domainAddress, domainPort
}

func getNuoSQLVersion(t *testing.T) *semver.Version {
	pathToNuoSQL, ok := os.LookupEnv("NUOSQL_PATH")
	if !ok {
		pathToNuoSQL = "nuosql"
	}
	out, err := exec.Command(pathToNuoSQL, "--version").Output()
	if exiterr, ok := err.(*exec.ExitError); ok {
		// nuosql has exited with an non zero exit code; generate better error
		// message by including the command stderr
		require.NoError(t, err, string(exiterr.Stderr))
	} else {
		require.NoError(t, err)
	}
	match := regexp.MustCompile("NuoDB Client build (.*)").FindStringSubmatch(string(out))
	require.NotNil(t, match, out)
	version, err := semver.NewVersion(match[1])
	require.NoError(t, err)
	return version
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

func verifyDomainProcesses(t *testing.T, api_server string, databaseName string, keysPath string, expectedNrProcesses int) []testlib.NuoDBProcess {
	pathToNuoCmd, ok := os.LookupEnv("NUOCMD_PATH")
	if !ok {
		pathToNuoCmd = "nuocmd"
	}
	cmd := fmt.Sprintf("%s --api-server %s", pathToNuoCmd, api_server)
	if keysPath != "" {
		cmd += fmt.Sprintf(" --client-key %s --verify-server %s",
			path.Join(keysPath, testlib.NUOCMD_FILE), path.Join(keysPath, testlib.CA_CERT_FILE))
	}
	cmd += fmt.Sprintf(" --show-json get processes --db-name %s", databaseName)

	var out []byte
	getProcesses := func() error {
		var err error
		t.Logf("Running command: <%s>", cmd)
		out, err = exec.Command("sh", "-c", cmd).Output()
		if exiterr, ok := err.(*exec.ExitError); ok {
			// nuocmd has exited with an non zero exit code; generate better error
			// message by including the command stderr
			return fmt.Errorf(string(exiterr.Stderr))
		} else {
			return err
		}
	}

	// retry the connections that go via Ingress as HAProxy processes may need
	// to be restarted due to reconfiguration request and nuocmd client will
	// receive "Connection aborted." error
	err := testlib.Retry(t, getProcesses, 4, 5*time.Second)
	require.NoError(t, err)

	err, processes := testlib.Unmarshal(string(out))
	require.NoError(t, err)
	require.Equal(t, expectedNrProcesses, len(processes), "Unexpected number of domain processes")
	for _, process := range processes {
		require.Equal(t, "MONITORED", process.DState)
		require.Equal(t, "RUNNING", process.State)
	}
	return processes
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
				servicePort, err := findServicePortE(s, "48006-tcp")
				require.NoError(t, err)
				val, ok := process.Labels["external-port"]
				require.True(t, ok)
				require.Equal(t, strconv.Itoa(int(servicePort.NodePort)), val)
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

		verifyService(t, namespaceName, admin0, "nuodb-clusterip", "ClusterIP", false)
		verifyService(t, namespaceName, admin0, fmt.Sprintf("nuodb-%s", serviceSuffix), serviceType, false)

		defer testlib.Teardown(testlib.TEARDOWN_DATABASE) // ensure resources allocated in called functions are released when this function exits

		databaseOptions := helm.Options{
			SetValues: map[string]string{
				"database.sm.resources.requests.cpu":              "0.25",
				"database.sm.resources.requests.memory":           testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
				"database.te.resources.requests.cpu":              "0.25",
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
		databaseOptions.SetValues["database.sm.hotCopy.enablePod"] = "false"
		databaseOptions.SetValues["database.sm.noHotCopy.enablePod"] = "false"
		databaseOptions.SetValues["database.te.labels.tx-type"] = "TAP"
		group2ReleaseName := testlib.StartDatabase(t, namespaceName, admin0, &databaseOptions)
		opt := testlib.GetExtractedOptions(&databaseOptions)

		testlib.AwaitDatabaseUp(t, namespaceName, admin0, opt.DbName, 3)

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

		address, port := findDomainExternalInfo(t, namespaceName, fmt.Sprintf("nuodb-%s", serviceSuffix), "48004-tcp")

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

func TestKubernetesIngress(t *testing.T) {
	// this requires support for "external-address" and "external-port" labels
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 4.2.3")
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetesingress-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	defer testlib.Teardown(testlib.TEARDOWN_SECRETS)
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	defer testlib.Teardown(testlib.TEARDOWN_HAPROXY)
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)

	// install HAProxy Ingress controller
	haProxyOptions := helm.Options{}
	haProxyReleaseName := testlib.StartHAProxyIngress(t, &haProxyOptions, namespaceName)
	ingressClassName := haProxyReleaseName
	_, ingressPort := findDomainExternalInfo(t, namespaceName, fmt.Sprintf("%s-kubernetes-ingress", haProxyReleaseName), "https")

	options := helm.Options{
		SetValues: map[string]string{
			"admin.ingress.enabled":                 "true",
			"admin.ingress.api.hostname":            testlib.ADMIN_API_INGRESS_HOSTNAME,
			"admin.ingress.api.className":           ingressClassName,
			"admin.ingress.sql.hostname":            testlib.ADMIN_SQL_INGRESS_HOSTNAME,
			"admin.ingress.sql.className":           ingressClassName,
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.ingress.enabled":           "true",
			"database.te.ingress.hostname":          testlib.DATABASE_TE_INGRESS_HOSTNAME,
			"database.te.ingress.className":         ingressClassName,
			// The product doesn't support Ingress resource lookup, so define
			// the port manually; most of the time in production a service of
			// type=LoadBalancer will be provisioned for the Ingress Controller
			// which will keep the default HTTPs port to 443; helm charts will
			// configure automatically port 443 if nothing else is configured
			"database.te.labels.external-port": strconv.Itoa(int(ingressPort)),
		},
	}

	_, keysLocation := testlib.GenerateAndSetTLSKeys(t, &options, namespaceName)
	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	// install the database Helm release
	testlib.StartDatabase(t, namespaceName, admin0, &options)
	opt := testlib.GetExtractedOptions(&options)

	// verify that we can connect to REST API from the test machine
	admin_url := fmt.Sprintf("https://%s:%d", testlib.ADMIN_API_INGRESS_HOSTNAME, ingressPort)
	processes := verifyDomainProcesses(t, admin_url, opt.DbName, keysLocation, opt.NrSmPods+opt.NrTePods)
	// verify that external-address and external-port labels are set correctly
	for _, process := range processes {
		if process.Type == "TE" {
			val, ok := process.Labels["external-address"]
			require.True(t, ok)
			require.Equal(t, testlib.DATABASE_TE_INGRESS_HOSTNAME, val)
			val, ok = process.Labels["external-port"]
			require.True(t, ok)
			require.Equal(t, strconv.Itoa(int(ingressPort)), val)
		}
	}

	// this requires CDriver to support Server Name Indication (SNI)
	t.Run("sqlIngressConnectivity", func(t *testing.T) {
		constraint, err := semver.NewConstraint(">=5.0.4")
		require.NoError(t, err)
		nuosqlVersion := getNuoSQLVersion(t)
		if !constraint.Check(nuosqlVersion) {
			t.Skip("Skipping test because nuosql version %s does not support SNI", nuosqlVersion.String())
		}

		var expectedNodeIds []int32
		for _, process := range processes {
			if process.Type == "TE" {
				expectedNodeIds = append(expectedNodeIds, process.NodeId)
			}
		}

		verifyNuoSQLEngine(t, testlib.ADMIN_SQL_INGRESS_HOSTNAME, ingressPort,
			opt.DbName, map[string]string{
				"trustStore":       path.Join(keysLocation, testlib.CA_CERT_FILE),
				"verifyHostname":   "false",
				"allowSRPFallback": "false",
			},
			expectedNodeIds)
	})
}
