package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nuodb/nuodb-helm-charts/v3/test/testlib"
)

const shebangPattern = "[#][!][ \t]*.+[/ ](sh|bash)"

func TestLintShellScripts(t *testing.T) {
	rootPaths := [...]string{
		testlib.ADMIN_HELM_CHART_PATH,
		testlib.DATABASE_HELM_CHART_PATH,
		testlib.RESTORE_HELM_CHART_PATH,
		testlib.THP_HELM_CHART_PATH,
		testlib.YCSB_HELM_CHART_PATH,
	}

	shebang, err := regexp.Compile(shebangPattern)
	require.NoError(t, err, "Cannot compile pattern '%s'", shebangPattern)

	first20 := make([]byte, 20)

	for _, root := range rootPaths {
		rootPath := path.Join(root, "files")

		info, err := os.Stat(rootPath)
		require.NoError(t, err)

		if !info.IsDir() {
			continue
		}

		// open the path
		dir, err := os.Open(rootPath)
		require.NoError(t, err)

		defer dir.Close() // remember to close the File

		// get all files in the dir
		scripts, err := dir.Readdir(-1)
		require.NoError(t, err, "Error reading dir %s: %s", rootPath, err)

		for _, scriptInfo := range scripts {
			scriptPath := path.Join(rootPath, scriptInfo.Name())

			script, err := os.Open(scriptPath)
			require.NoError(t, err)

			defer script.Close()

			_, err = script.Read(first20)
			require.NoError(t, err)

			groups := shebang.FindStringSubmatch(string(first20))
			if len(groups) == 0 {
				t.Logf("Script %s is not a shell script >>%s<<", scriptPath, string(first20))
				continue
			}

			t.Logf("Running `%s -n %s`", groups[1], scriptPath)
			cmd := exec.Command(groups[1], "-n", scriptPath)

			out, err := cmd.CombinedOutput()
			assert.NoError(t, err, fmt.Sprintln(scriptPath, "LINT failed:", err, "\n", string(out)))
		}
	}
}
