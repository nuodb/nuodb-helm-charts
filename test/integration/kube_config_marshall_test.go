package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/require"

	"io/ioutil"
	"testing"
)


func TestKubeConfigUnmarshall(t *testing.T) {
	content, err := ioutil.ReadFile("../files/nuodb-dump.json")
	require.NoError(t, err)

	err, objects := testlib.UnmarshalNuoDBKubeConfig(string(content))

	require.NoError(t, err)
	require.Equal(t, len(objects), 1)

	config := objects[0]

	require.True(t, len(config.Deployments) == 1)
	require.True(t, len(config.StatefulSets) == 3)
	require.True(t, len(config.Pods) == 5)
	require.True(t, len(config.Volumes) == 3)

	// StatefulSets
	require.True(t, func() bool {_, ok := config.StatefulSets["admin-u7mxhy-nuodb-cluster0"]; return ok}())
	require.True(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo"]; return ok}())
	require.True(t, func() bool {_, ok := config.StatefulSets["sm-database-vslldk-nuodb-cluster0-demo-hotcopy"]; return ok}())

	// Deployments
	require.True(t, func() bool {_, ok := config.Deployments["te-database-vslldk-nuodb-cluster0-demo"]; return ok}())

	// Admin Volumes
	require.True(t, func() bool {_, ok := config.Volumes["raftlog-admin-u7mxhy-nuodb-cluster0-0"]; return ok}())

	// DB Volumes
	require.True(t, func() bool {_, ok := config.Volumes["archive-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	require.True(t, func() bool {_, ok := config.Volumes["backup-volume-sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())

	// Admin Pods
	require.True(t, func() bool {_, ok := config.Pods["admin-u7mxhy-nuodb-cluster0-0"]; return ok}())
	require.True(t, func() bool {_, ok := config.Pods["job-lb-policy-nearest-zs9jl"]; return ok}())

	// DB Pods
	require.True(t, func() bool {_, ok := config.Pods["sm-database-vslldk-nuodb-cluster0-demo-hotcopy-0"]; return ok}())
	require.True(t, func() bool {_, ok := config.Pods["hotcopy-demo-job-initial-549rc"]; return ok}())
	require.True(t, func() bool {_, ok := config.Pods["te-database-vslldk-nuodb-cluster0-demo-65c4cdf487-wbzj9"]; return ok}())

}
