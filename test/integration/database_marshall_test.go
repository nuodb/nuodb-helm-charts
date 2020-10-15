package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestDatabaseUnmarshal(t *testing.T) {
	s := (
		`{
  "incarnation": {
    "major": 1, 
    "minor": 2
  }, 
  "name": "demo", 
  "state": "NOT_RUNNING", 
  "uri": "https://localhost:8888/api/1/databases/demo"
}`)

	err, objects := testlib.UnmarshalDatabase(s)

	require.NoError(t, err)
	require.Equal(t, len(objects), 1)

	obj := objects[0]

	require.True(t, obj.Name == "demo")
	require.True(t, obj.State == "NOT_RUNNING")
	require.True(t, obj.Incarnation.Major == 1)
	require.True(t, obj.Incarnation.Minor == 2)

}
