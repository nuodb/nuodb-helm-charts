package testlib

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"gotest.tools/assert"
	"os/exec"
	"testing"
)

var isOpenShift *bool

func isOpenShiftEnvironment(t *testing.T) bool {
	if isOpenShift == nil {
		output, err := exec.Command("oc", "status").Output()

		var isOs = (err == nil)
		isOpenShift = &isOs

		t.Logf("Running in OpenShift:\n%s", string(output))
	}

	return *isOpenShift
}

func createOpenShiftProject(t *testing.T, namespaceName string) {
	output, err := exec.Command("oc", "new-project", namespaceName).Output()
	assert.NilError(t, err, output)
	output, err = exec.Command("oc", "policy", "add-role-to-user", "edit", "system:serviceaccount:tiller:tiller").Output()
	assert.NilError(t, err, output)
	output, err = exec.Command("oc", "policy", "add-role-to-user", "hostaccess", "system:serviceaccount:tiller:tiller").Output()
	assert.NilError(t, err, output)
	output, err = exec.Command("oc", "policy", "add-role-to-user", "privileged", "system:serviceaccount:tiller:tiller").Output()
	assert.NilError(t, err, output)
	output, err = exec.Command("oc", "policy", "add-role-to-user", "anyuid", "system:serviceaccount:tiller:tiller").Output()
	assert.NilError(t, err, output)
	output, err = exec.Command("oc", "adm", "policy", "add-scc-to-user", "privileged", fmt.Sprintf("system:serviceaccount:%s:default", namespaceName)).Output()
	assert.NilError(t, err, output)
}

func InjectOpenShiftValues(t *testing.T, options *helm.Options) {
	if options.SetValues == nil {
		options.SetValues = make(map[string]string)
	}

	if isOpenShiftEnvironment(t) {
		options.SetValues["openshift.enabled"] = "true"
	}
}