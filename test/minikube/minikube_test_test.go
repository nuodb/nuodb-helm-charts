package minikube

import (
	"testing"
	"time"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"

	// NTJ - commented until we can get TestAwaitFailure() working - then uncomment; or remove if the solution does not use this package
	// verify "github.com/stretchr/testify/assert"

	"gotest.tools/assert"
)

/**
 * A set of tests thattest the test infrastructure
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
		assert.Check(t, tdcounter == 3)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Check(t, tdcounter == 2)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Check(t, tdcounter == 1)
	})

	testlib.Teardown("")

	assert.Check(t, tdcounter == 3)

	testlib.VerifyTeardown(t)
}

func TestNamedTeardown(t *testing.T) {
	tdcounter := 0

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Check(t, tdcounter == 3)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Check(t, tdcounter == 2)
	})

	testlib.AddTeardown("", func() {
		tdcounter++
		assert.Check(t, tdcounter == 1)
	})

	testlib.AddTeardown("other", func() {
		tdcounter++
		assert.Check(t, tdcounter == 5)
	})

	testlib.AddTeardown("other", func() {
		tdcounter++
		assert.Check(t, tdcounter == 4)
	})

	testlib.Teardown("")
	assert.Check(t, tdcounter == 3)

	testlib.Teardown("other")
	assert.Check(t, tdcounter == 5)

	testlib.VerifyTeardown(t)
}
