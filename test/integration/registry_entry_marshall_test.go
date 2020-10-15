package integration

import (
	"github.com/nuodb/nuodb-helm-charts/test/testlib"
	"github.com/stretchr/testify/require"

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

	require.NoError(t, err)

	require.True(t, object.Nuodb.Image.Registry == "local")
	require.True(t, object.Nuodb.Image.Repository == "master")
	require.True(t, object.Nuodb.Image.Tag == "latest")
}
