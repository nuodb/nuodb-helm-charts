// +build short

package minikube

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestUnconditionalDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", true, func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 3, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that an unconditionally skipped DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestSkippedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", false, func() {
		assert.FailNow(t, "This diagnostic teardown should not have been run")
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 2, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestUnconditionalFuncDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", func() bool { return true }, func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 3, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestSkippedFuncDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", func() bool { return false }, func() {
		assert.FailNow(t, "This diagnostic teardown should not have been run")
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 2, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *always* executed *before* any other teardown for the same name if the passed testing.T has Failed */
func TestTFailedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	tt := new(testing.T)

	testlib.AddDiagnosticTeardown("name", tt, func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	tt.Fail()

	testlib.Teardown("name")

	assert.Equal(t, 3, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *not* executed if the passed in T has *not* Failed */
func TestTNotFailedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", t, func() {
		assert.FailNow(t, "This diagnostic teardown should not have been called since T has not failed")
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 2, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *not* executed if the conditional is nil */
func TestNilDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", nil, func() {
		assert.FailNow(t, "Diagnostic teardown should not be run if the conditional is nil")
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.Teardown("name")

	assert.Equal(t, 2, tdcounter)

	testlib.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *always* executed if the ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS is true */
func TestUnconditionalEnvVarDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddDiagnosticTeardown("name", false, func() {
		tdcounter++
		assert.Equal(t, 1, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 3, tdcounter)
	})

	testlib.AddTeardown("name", func() {
		tdcounter++
		assert.Equal(t, 2, tdcounter)
	})

	testlib.AlwaysRunDiagnosticTeardowns = true
	defer func() { testlib.AlwaysRunDiagnosticTeardowns = false }()

	testlib.Teardown("name")

	assert.Equal(t, 3, tdcounter)

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
				"cloud.cluster.name":             "cluster1",
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
	testlib.InjectTestValues(t, &options)
}
