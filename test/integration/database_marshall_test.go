package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"gotest.tools/assert"
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

	assert.NilError(t, err)
	assert.Equal(t, len(objects), 1)

	obj := objects[0]

	assert.Check(t, obj.Name == "demo")
	assert.Check(t, obj.State == "NOT_RUNNING")
	assert.Check(t, obj.Incarnation.Major == 1)
	assert.Check(t, obj.Incarnation.Minor == 2)

}
