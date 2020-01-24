// +build short

package minikube

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"

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

func TestGetExtractedOptions(t *testing.T) {

	t.Run("emptyOptions", func(t *testing.T) {
		opt := testlib.GetExtractedOptions(&helm.Options{
			SetValues: map[string]string{},
		})

		assert.Check(t, opt.DbName == "demo")
		assert.Check(t, opt.NrTePods == 1)
		assert.Check(t, opt.NrSmPods == 1)
		assert.Check(t, opt.ClusterName == "cluster0")
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

		assert.Check(t, opt.DbName == "green")
		assert.Check(t, opt.NrTePods == 2)
		assert.Check(t, opt.NrSmHotCopyPods == 2)
		assert.Check(t, opt.NrSmNoHotCopyPods == 2)
		assert.Check(t, opt.NrSmPods == 4)
		assert.Check(t, opt.ClusterName == "cluster1")
	})

}

func TestParseIncarnations(t *testing.T) {

	incarnation := `   incarnation:(3, 0)
	[SM] sm-database-j7uzoi-nuodb-cluster0-demo-hotcopy-0/172.17.0.11:48006 [start_id = 3] [server_id = admin-uz7aiw-nuodb-cluster0-0] [pid = 152] [node_id = 1] [last_ack =  1.09] MONITORED:RUNNING
	[TE] te-database-j7uzoi-nuodb-cluster0-demo-7789b6c57f-99xnz/172.17.0.4:48006 [start_id = 4] [server_id = admin-uz7aiw-nuodb-cluster0-0] [pid = 38] [node_id = 2] [last_ack =  0.56] MONITORED:RUNNING`

	result := testlib.ParseDatabaseIncarnation(t, incarnation)

	assert.Check(t, result[0] == 3)
	assert.Check(t, result[1] == 0)
}
