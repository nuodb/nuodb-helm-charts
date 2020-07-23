package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"testing"
)


func TestKubeConfigUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("../files/nuodb-dump.json")
	assert.NoError(t, err)

	err, objects := testlib.UnmarshalNuoDBKubeConfig(string(content))

	assert.NoError(t, err)
	assert.Equal(t, len(objects), 1)

	config := objects[0]

	assert.True(t, len(config.Deployments) == 1)
	assert.True(t, len(config.StatefulSets) == 3)
	assert.True(t, len(config.Pods) == 5)
	assert.True(t, len(config.Volumes) == 3)

	// StatefulSets
	assert.True(t, func() bool {_, ok := config.StatefulSets["admin-u7mxhy-nuodb-cluster0"]; return ok}())
	assert.True(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo"]; return ok}())
	assert.True(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo-hotcopy"]; return ok}())

	// Deployments
	assert.True(t, func() bool {_, ok := config.Deployments["te-database-vslldk-nuodb-cluster0-demo"]; return ok}())

	// Admin Volumes
	assert.True(t, func() bool {_, ok := config.Volumes["raftlog-admin-u7mxhy-nuodb-cluster0-0"]; return ok}())

	// DB Volumes
	assert.True(t, func() bool {_, ok := config.Volumes["archive-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	assert.True(t, func() bool {_, ok := config.Volumes["backup-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())

	// Admin Pods
	assert.True(t, func() bool {_, ok := config.Pods["admin-u7mxhy-nuodb-cluster0-0"]; return ok}())
	assert.True(t, func() bool {_, ok := config.Pods["job-lb-policy-nearest-zs9jl"]; return ok}())

	// DB Pods
	assert.True(t, func() bool {_, ok := config.Pods["sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	assert.True(t, func() bool {_, ok := config.Pods["hotcopy-demo-job-initial-549rc"]; return ok}())
	assert.True(t, func() bool {_, ok := config.Pods["te-database-vslldk-nuodb-cluster0-demo-65c4cdf487-wbzj9"]; return ok}())

}
