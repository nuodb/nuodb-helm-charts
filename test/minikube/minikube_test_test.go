// +build short

package minikube

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"


)

/**
 * A set of tests that test the test infrastructure
 */

func TestAwaitSuccess(t *testing.T) {
	testlib.Await(t, func() bool { return true }, 2*time.Second)
}

// NTJ - supposed to work - test still fails - giving up until a solution can be found.
// Leaving here - but commented - until we either find a solution, or decide that this type of test is impossible in Go.
// func TestAwaitFailure(t *testing.T) {
// 	verify.Panics(t, func() { await(t, func() bool { return false }, 2*time.Second) }, "Await timeout did not panic")
// }

func TestTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.Teardown("")

	assert.Equal(t, 3, tdcounter)

	testlib.VerifyTeardown(t)
}

func TestNamedTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.AddTeardown("other", func() {
		tdcounter++
		assert.Equal(t, 5, tdcounter)
	})

	testlib.AddTeardown("other", func() {
		tdcounter++
		assert.Equal(t, 4, tdcounter)
	})

	testlib.Teardown("")
	assert.Equal(t, 3, tdcounter)

	testlib.Teardown("other")
	assert.Equal(t, 5, tdcounter)

	testlib.VerifyTeardown(t)
}

func TestGetExtractedOptions(t *testing.T) {

	t.Run("emptyOptions", func(t *testing.T) {
		opt := testlib.GetExtractedOptions(&helm.Options{
			SetValues: map[string]string{},
		})

		assert.Equal(t, "demo", opt.DbName)
		assert.Equal(t, 1, opt.NrTePods)
		assert.Equal(t, 1, opt.NrSmPods)
		assert.Equal(t, "cluster0", opt.ClusterName)
	})

	t.Run("overriddenOptions", func(t *testing.T) {
		opt := testlib.GetExtractedOptions(&helm.Options{
			SetValues: map[string]string{
				"database.name":                  "green",
				"database.te.replicas":           "2",
				"database.sm.hotCopy.replicas":   "2",
				"database.sm.noHotCopy.replicas": "2",
				"cloud.cluster.name":              "cluster1",
			},
		})

		assert.Equal(t, "green", opt.DbName)
		assert.Equal(t, 2, opt.NrTePods)
		assert.Equal(t, 2, opt.NrSmHotCopyPods)
		assert.Equal(t, 2, opt.NrSmNoHotCopyPods)
		assert.Equal(t, 4, opt.NrSmPods)
		assert.Equal(t, "cluster1", opt.ClusterName)
	})

}

func TestInjection(t *testing.T) {
	options := helm.Options{
		SetValues: map[string]string{},
	}

	testlib.InjectTestVersion(t, &options)
}