package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestDatabaseUnmarshal(t *testing.T) {
	s := (`{
  "incarnation": {
    "major": 1, 
    "minor": 2
  }, 
  "name": "demo", 
  "state": "NOT_RUNNING", 
  "uri": "https://localhost:8888/api/1/databases/demo"
}`)

	err, objects := testlib.UnmarshalDatabase(s)

	assert.NoError(t, err)
	assert.Equal(t, len(objects), 1)

	obj := objects[0]

	assert.True(t, obj.Name == "demo")
	assert.True(t, obj.State == "NOT_RUNNING")
	assert.True(t, obj.Incarnation.Major == 1)
	assert.True(t, obj.Incarnation.Minor == 2)

}
