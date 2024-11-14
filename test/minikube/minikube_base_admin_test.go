//go:build short
// +build short

package minikube

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestKubernetesBasicAdminSingleReplica(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	clusterServiceName := fmt.Sprintf("nuodb-clusterip")

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		if os.Getenv("NUODB_LICENSE") == "ENTERPRISE" {
			t.Skip("Cannot test licensing in Enterprise Edition")
		}
		if os.Getenv("NUODB_LIMITED_LICENSE_CONTENT") == "" {
			t.Skip("Cannot test licensing without Limited License")
		}
		testlib.VerifyLicense(t, namespaceName, admin0, testlib.LIMITED)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyAdminKvSetAndGet", func(t *testing.T) {
		testlib.VerifyAdminKvSetAndGet(t, admin0, namespaceName)
	})
	t.Run("verifyAdminClusterService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0, clusterServiceName, false) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyPodKill", func(t *testing.T) {
		t.Skip("verifyPodKill is flaky")
		verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 1)
	})
}

func TestKubernetesAdminLicenseSecret(t *testing.T) {
	testlib.SkipTestOnNuoDBVersionCondition(t, "< 5.1.1")
	if os.Getenv("NUODB_LICENSE_CONTENT") == "" {
		t.Skip("Cannot run this test without a valid license")
	}
	defer testlib.VerifyTeardown(t)

	randomSuffix := strings.ToLower(random.UniqueId())
	namespaceName := fmt.Sprintf("%skubernetesadminlicensesecret-%s", testlib.NAMESPACE_NAME_PREFIX, randomSuffix)
	testlib.CreateNamespace(t, namespaceName)

	ctx := context.Background()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nuodb-license",
			Namespace: namespaceName,
		},
		StringData: map[string]string{
			"nuodb.lic": "garbage",
		},
	}

	options := helm.Options{
		SetValues: map[string]string{
			"admin.license.secret": "nuodb-license",
		},
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, kubectlOptions)
	require.NoError(t, err)
	clientset.CoreV1().Secrets(namespaceName).Create(ctx, secret, metav1.CreateOptions{})

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, _ := testlib.StartAdmin(t, &options, 1, namespaceName)

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	testlib.VerifyLicense(t, namespaceName, admin0, testlib.UNLICENSED)
	testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, true)

	// Update the secret with the NuoDB EE license
	licenseContentBytes, err := base64.StdEncoding.DecodeString(os.Getenv("NUODB_LICENSE_CONTENT"))
	require.NoError(t, err)
	secret, err = clientset.CoreV1().Secrets(namespaceName).Get(ctx, secret.Name, metav1.GetOptions{})
	require.NoError(t, err)
	secret.Data["nuodb.lic"] = licenseContentBytes
	secret, err = clientset.CoreV1().Secrets(namespaceName).Update(ctx, secret, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("Waiting for license update")
	testlib.Await(t, func() bool {
		if err := testlib.VerifyLicenseE(t, namespaceName, admin0, testlib.ENTERPRISE); err != nil {
			return false
		}
		return true
	}, 180*time.Second)
}

func TestKubernetesInvalidLicense(t *testing.T) {
	testlib.SkipTestOnNuoDBVersionCondition(t, ">= 6.0.0")
	defer testlib.VerifyTeardown(t)

	licenseString := "red-riding-hood"
	customFile := "customFile"

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.configFiles.nuodb\\.lic":                 licenseString,
			fmt.Sprintf("admin.configFiles.%s", customFile): "TestKubernetesInvalidLicense"},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)

	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		if os.Getenv("NUODB_LICENSE") == "ENTERPRISE" {
			t.Skip("Cannot test licensing in Enterprise Edition")
		}
		testlib.VerifyLicense(t, namespaceName, admin0, testlib.UNLICENSED)

		// the license provided is not a valid PEM file
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, true)
	})
	t.Run("verifyLicenseFile", func(t *testing.T) {
		testlib.VerifyLicenseFile(t, namespaceName, admin0, licenseString)
	})

}

func TestKubernetesBasicNameOverride(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.nameOverride": "aws-a",
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, options, 1, "")
	admin0 := fmt.Sprintf("%s-nuodb-cluster0-%s-0", helmChartReleaseName, "aws-a")

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}

func TestKubernetesFullNameOverride(t *testing.T) {
	defer testlib.VerifyTeardown(t)

	nonDefaultName := "nondefault-adminname"
	admin0 := fmt.Sprintf("%s-0", nonDefaultName)

	options := &helm.Options{
		SetValues: map[string]string{
			"admin.fullnameOverride": nonDefaultName,
		},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	_, namespaceName := testlib.StartAdmin(t, options, 1, "")

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
}
