// +build short

package minikube

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
)

func TestKubernetesBasicAdminSingleReplica(t *testing.T) {
	testlib.AwaitTillerUp(t)
	defer testlib.VerifyTeardown(t)

	options := helm.Options{
		SetValues: map[string]string{},
	}

	defer testlib.Teardown(testlib.TEARDOWN_ADMIN)

	helmChartReleaseName, namespaceName := testlib.StartAdmin(t, &options, 1, "")

	admin0 := fmt.Sprintf("%s-nuodb-cluster0-0", helmChartReleaseName)
	headlessServiceName := fmt.Sprintf("nuodb")
	clusterServiceName := fmt.Sprintf("nuodb-clusterip")

	t.Run("verifyAdminState", func(t *testing.T) { testlib.VerifyAdminState(t, namespaceName, admin0) })
	t.Run("verifyOrderedLicensing", func(t *testing.T) {
		if os.Getenv("NUODB_LICENSE") == "ENTERPRISE" {
			t.Skip("Cannot test licensing in Enterprise Edition")
		}
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, false) // no error
	})
	t.Run("verifyAdminKvSetAndGet", func(t *testing.T) {
		testlib.VerifyAdminKvSetAndGet(t, admin0, namespaceName)
	})
	t.Run("verifyAdminHeadlessService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0, headlessServiceName, true) })
	t.Run("verifyAdminClusterService", func(t *testing.T) { verifyAdminService(t, namespaceName, admin0, clusterServiceName, false) })
	t.Run("verifyLBPolicy", func(t *testing.T) { verifyLBPolicy(t, namespaceName, admin0) })
	t.Run("verifyProcessKill", func(t *testing.T) { verifyKillProcess(t, namespaceName, admin0, helmChartReleaseName, 1) })
	t.Run("verifyPodKill", func(t *testing.T) {
		t.Skip("verifyPodKill is flaky")
		verifyPodKill(t, namespaceName, admin0, helmChartReleaseName, 1)
	})
}

func TestKubernetesInvalidLicense(t *testing.T) {
	testlib.AwaitTillerUp(t)
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
		testlib.VerifyLicenseIsCommunity(t, namespaceName, admin0)

		// the license provided is not a valid PEM file
		testlib.VerifyLicensingErrorsInLog(t, namespaceName, admin0, true)
	})
	t.Run("verifyLicenseFile", func(t *testing.T) {
		testlib.VerifyLicenseFile(t, namespaceName, admin0, licenseString)
	})

}

func TestKubernetesBasicNameOverride(t *testing.T) {
	testlib.AwaitTillerUp(t)
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
	testlib.AwaitTillerUp(t)
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
