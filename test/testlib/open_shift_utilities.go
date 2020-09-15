package testlib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var isOpenShift *bool

func IsOpenShiftEnvironment(t *testing.T) bool {
	if isOpenShift == nil {
		// Use OpenShift if user is logged-in and session token is available
		output, err := exec.Command("oc", "whoami", "-t").Output()

		var isOs = (err == nil)
		isOpenShift = &isOs

		if isOs {
			t.Logf("Running in OpenShift:\n%s", string(output))
		}
	}

	return *isOpenShift
}

func createOpenShiftProject(t *testing.T, namespaceName string) {
	output, err := exec.Command("oc", "new-project", namespaceName).Output()
	assert.NoError(t, err, output)

	pwd, err := os.Getwd()
	assert.NoError(t, err)

	sccFilePath := filepath.Join(pwd, "..", "..", "deploy", "nuodb-scc.yaml")

	output, err = exec.Command("oc", "apply", "-n", namespaceName, "-f", sccFilePath).Output()
	assert.NoError(t, err, output)

	output, err = exec.Command("oc", "adm", "policy", "add-scc-to-user", "nuodb-scc", fmt.Sprintf("system:serviceaccount:%s:default", namespaceName)).Output()
	assert.NoError(t, err, output)

	output, err = exec.Command("oc", "adm", "policy", "add-scc-to-user", "nuodb-scc", fmt.Sprintf("system:serviceaccount:%s:nuodb", namespaceName)).Output()
	assert.NoError(t, err, output)
}
