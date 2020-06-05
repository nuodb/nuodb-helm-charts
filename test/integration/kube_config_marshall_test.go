package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"gotest.tools/assert"
	"io/ioutil"
	"testing"
)


func TestKubeConfigUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("../files/nuodb-dump.json")
	assert.NilError(t, err)

	err, objects := testlib.UnmarshalNuoDBKubeConfig(string(content))

	assert.NilError(t, err)
	assert.Equal(t, len(objects), 1)

	config := objects[0]

	assert.Check(t, len(config.Deployments) == 1)
	assert.Check(t, len(config.StatefulSets) == 3)
	assert.Check(t, len(config.Pods) == 5)
	assert.Check(t, len(config.Volumes) == 3)

	// StatefulSets
	assert.Check(t, func() bool {_, ok := config.StatefulSets["admin-u7mxhy-nuodb-cluster0"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo-hotcopy"]; return ok}())

	// Deployments
	assert.Check(t, func() bool {_, ok := config.Deployments["te-database-vslldk-nuodb-cluster0-demo"]; return ok}())

	// Admin Volumes
	assert.Check(t, func() bool {_, ok := config.Volumes["raftlog-admin-u7mxhy-nuodb-cluster0-0"]; return ok}())

	// DB Volumes
	assert.Check(t, func() bool {_, ok := config.Volumes["archive-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.Volumes["backup-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())

	// Admin Pods
	assert.Check(t, func() bool {_, ok := config.Pods["admin-u7mxhy-nuodb-cluster0-0"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.Pods["job-lb-policy-nearest-zs9jl"]; return ok}())

	// DB Pods
	assert.Check(t, func() bool {_, ok := config.Pods["sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.Pods["hotcopy-demo-job-initial-549rc"]; return ok}())
	assert.Check(t, func() bool {_, ok := config.Pods["te-database-vslldk-nuodb-cluster0-demo-65c4cdf487-wbzj9"]; return ok}())

}
