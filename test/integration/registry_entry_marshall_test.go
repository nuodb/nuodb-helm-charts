package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"gotest.tools/assert"
	"testing"
)

func TestRegistryEntryUnmarshal(t *testing.T) {
	s := (
`nuodb:
  image:
    registry: local
    repository: master
    tag: latest
`)

	err, object := testlib.UnmarshalImageYAML(s)

	assert.NilError(t, err)

	assert.Check(t, object.Nuodb.Image.Registry == "local")
	assert.Check(t, object.Nuodb.Image.Repository == "master")
	assert.Check(t, object.Nuodb.Image.Tag == "latest")
}
