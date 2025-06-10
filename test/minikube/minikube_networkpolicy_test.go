//go:build long
// +build long

package minikube

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestConnectivityWithNetworkPolicy(t *testing.T) {
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 8.0.0")

	defer testlib.VerifyTeardown(t)
	namespace := testlib.CreateNamespaceForTest(t, true)

	// Create network policy that limits connectivity to group=nuodb
	kubeOptions := k8s.NewKubectlOptions("", "", namespace)
	client, err := k8s.GetKubernetesClientFromOptionsE(t, kubeOptions)
	require.NoError(t, err)
	ctx := context.Background()
	networkPolicy := getNetworkPolicy(namespace)
	networkPolicyClient := client.NetworkingV1().NetworkPolicies(namespace)
	defer func() {
		_ = networkPolicyClient.Delete(ctx, networkPolicy.GetObjectMeta().GetName(), metav1.DeleteOptions{})
	}()
	_, err = networkPolicyClient.Create(ctx, networkPolicy, metav1.CreateOptions{})
	require.NoError(t, err)

	// Create nuodb/admin release with two APs and wait for them to become ready
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	adminReleaseName, _ := testlib.StartAdmin(t, &helm.Options{
		SetValues: map[string]string{
			"admin.replicas": "2",
		},
	}, 2, namespace)
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", adminReleaseName)

	// Create nuodb/database release and wait for database to become ready
	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	testlib.StartDatabase(t, namespace, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	// Verify that connectivity is enabled between AP pods
	err = k8s.RunKubectlE(t, kubeOptions, "exec", admin0, "-c", "admin", "--",
		"nuocmd", "check", "servers", "--wait-for-acks", "--timeout", "10")
	require.NoError(t, err)

	// Verify that connectivity (egress) is disabled out of domain
	err = k8s.RunKubectlE(t, kubeOptions, "exec", admin0, "-c", "admin", "--",
		"curl", "--silent", "--show-error", "--connect-timeout", "1", "https://google.com")
	errWithOutput, ok := err.(*shell.ErrWithCmdOutput)
	require.Truef(t, ok, "Expected shell.ErrWithCmdOutput, got: %+v", err)
	require.Contains(t, errWithOutput.Output.Stderr(), "Failed to connect to google.com port 443: Connection timed out")
}

func getNetworkPolicy(namespace string) *v1.NetworkPolicy {
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"group": "nuodb",
		},
	}
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nuodb-pol",
			Namespace: namespace,
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: labelSelector,
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
				v1.PolicyTypeEgress,
			},
			Ingress: []v1.NetworkPolicyIngressRule{{
				From: []v1.NetworkPolicyPeer{{
					PodSelector: &labelSelector,
				}},
			}},
			Egress: []v1.NetworkPolicyEgressRule{{
				To: []v1.NetworkPolicyPeer{{
					PodSelector: &labelSelector,
				}},
			}, {
				// Enable egress to DNS to enable hostname resolution
				To: []v1.NetworkPolicyPeer{{
					NamespaceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"kubernetes.io/metadata.name": "kube-system",
						},
					},
				}},
				Ports: []v1.NetworkPolicyPort{{
					Port:     ptr.To(intstr.FromInt(53)),
					Protocol: ptr.To(corev1.ProtocolUDP),
				}, {
					Port:     ptr.To(intstr.FromInt(53)),
					Protocol: ptr.To(corev1.ProtocolTCP),
				}},
			}},
		},
	}
}
