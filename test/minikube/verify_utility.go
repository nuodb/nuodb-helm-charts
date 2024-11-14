package minikube

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func verifyAdminService(t *testing.T, namespaceName string, podName string, serviceName string, ping bool) {

	adminService := testlib.GetService(t, namespaceName, serviceName)
	require.Equal(t, adminService.Name, serviceName)

	if ping {
		testlib.PingService(t, namespaceName, serviceName, podName)
	}
}

func verifyLBPolicy(t *testing.T, namespaceName string, podName string) {
	testlib.VerifyPolicyInstalled(t, namespaceName, podName)
}

func verifyPodKill(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillAdminPod(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitPodUp(t, namespaceName, podName, 100*time.Second)
}

func verifyKillProcess(t *testing.T, namespaceName string, podName string, helmChartReleaseName string, nrReplicasExpected int) {
	testlib.KillProcess(t, namespaceName, podName)
	testlib.AwaitNrReplicasScheduled(t, namespaceName, helmChartReleaseName, nrReplicasExpected)
	testlib.AwaitPodUp(t, namespaceName, podName, 100*time.Second)
}

func verifyLoadBalancer(t *testing.T, namespaceName string, adminPod string, deploymentOptions map[string]string) {
	actualLoadBalancerConfigurations, err := testlib.GetLoadBalancerConfigE(t, namespaceName, adminPod)
	require.NoError(t, err)
	actualLoadBalancerPolicies, err := testlib.GetLoadBalancerPoliciesE(t, namespaceName, adminPod)
	require.NoError(t, err)
	actualGlobalConfig, err := testlib.GetGlobalLoadBalancerConfigE(t, actualLoadBalancerConfigurations)
	require.NoError(t, err)
	actualDatabaseConfig, err := testlib.GetDatabaseLoadBalancerConfigE(t, "demo", actualLoadBalancerConfigurations)
	require.NoError(t, err)

	configuredPolicies := len(deploymentOptions)
	for opt, val := range deploymentOptions {
		t.Logf("requireing deployment option %s with value %s", opt, val)
		if strings.HasPrefix(opt, "admin.lbConfig.policies.") {
			// Verify that named policies are configured properly
			policyName := opt[strings.LastIndex(opt, ".")+1:]
			actualPolicy, ok := actualLoadBalancerPolicies[policyName]
			require.Truef(t, ok, "Unable to find named policy=%s", policyName)
			require.Equal(t, val, actualPolicy.LbQuery)
		} else if opt == "admin.lbConfig.prefilter" {
			if actualGlobalConfig != nil {
				require.Equal(t, val, actualGlobalConfig.Prefilter)
			}
		} else if opt == "admin.lbConfig.default" {
			if actualGlobalConfig != nil {
				require.Equal(t, val, actualGlobalConfig.DefaultLbQuery)
			}
		} else if opt == "database.lbConfig.prefilter" {
			if actualDatabaseConfig != nil {
				require.Equal(t, val, actualDatabaseConfig.Prefilter)
			}
		} else if opt == "database.lbConfig.default" {
			if actualDatabaseConfig != nil {
				require.Equal(t, val, actualDatabaseConfig.DefaultLbQuery)
			}
		} else {
			t.Logf("Deployment option %s skipped", opt)
			configuredPolicies--
		}
	}

	if deploymentOptions["admin.lbConfig.fullSync"] == "true" {
		// Verify that named policies match configured number of policies
		t.Logf("requireing load-balancer policies count is equal to configured policies via Helm")
		require.Equal(t, configuredPolicies, len(actualLoadBalancerPolicies))
	}
}
