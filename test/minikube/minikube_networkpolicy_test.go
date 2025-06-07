/*
//go:build long
// +build long
*/

package minikube

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestDatabaseConnectivityWithNetworkPolicies(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{}
	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)
	helmChartReleaseName, namespace := testlib.StartAdmin(t, &options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	k8sOptions := k8s.NewKubectlOptions("", "", namespace)
	client, err := k8s.GetKubernetesClientFromOptionsE(t, k8sOptions)
	require.NoError(t, err)

	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "nuodb",
		},
	}
	networkPolicy := v1.NetworkPolicy{
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
			}},
		},
	}
	_, err = client.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), &networkPolicy, metav1.CreateOptions{})
	require.NoError(t, err)

	defer testlib.Teardown(testlib.TEARDOWN_DATABASE)
	testlib.StartDatabase(t, namespace, admin0, &helm.Options{
		SetValues: map[string]string{
			"database.sm.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.sm.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
			"database.te.resources.requests.cpu":    testlib.MINIMAL_VIABLE_ENGINE_CPU,
			"database.te.resources.requests.memory": testlib.MINIMAL_VIABLE_ENGINE_MEMORY,
		},
	})

	// Insert long sleep to enable interaction
	// TODO: remove
	time.Sleep(5 * time.Hour)
}
