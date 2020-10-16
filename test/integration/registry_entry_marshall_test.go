package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/assert"

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

	assert.NoError(t, err)

	assert.True(t, object.Nuodb.Image.Registry == "local")
	assert.True(t, object.Nuodb.Image.Repository == "master")
	assert.True(t, object.Nuodb.Image.Tag == "latest")
}
