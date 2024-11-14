package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

func TestRegistryEntryUnmarshal(t *testing.T) {
	s := (`nuodb:
  image:
    registry: local
    repository: master
    tag: latest
`)

	object, err := testlib.UnmarshalImageYAML(s)
	assert.NoError(t, err)

	assert.True(t, object.Nuodb.Image.Registry == "local")
	assert.True(t, object.Nuodb.Image.Repository == "master")
	assert.True(t, object.Nuodb.Image.Tag == "latest")
}
