package integration

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/nuodb/nuodb-helm-charts/test/testlib"
)

const shebangPattern = "[#][!][ \t]*/[^/]+/(sh|bash)"

func TestLintScripts(t *testing.T) {
	rootPaths := [...]string{
		testlib.ADMIN_HELM_CHART_PATH,
		testlib.DATABASE_HELM_CHART_PATH,
		testlib.RESTORE_HELM_CHART_PATH,
		testlib.THP_HELM_CHART_PATH,
		testlib.YCSB_HELM_CHART_PATH,
	}

	shebang, err := regexp.Compile(shebangPattern)
	if err != nil {
		t.Error("Cannot compile patterh:", shebangPattern)
	}

	first20 := make([]byte, 20)

	var errorList []string

	defer func() {
		if len(errorList) > 0 {
			t.Error(strings.Join(errorList, "\nAND "))
		}
	}()

	for _, root := range rootPaths {
		rootPath := path.Join(root, "files")

		info, err := os.Stat(rootPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !info.IsDir() {
			continue
		}

		// open the path
		dir, err := os.Open(rootPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer dir.Close() // remember to close the File

		// get all files in the dir
		scripts, err := dir.Readdir(-1)
		if err != nil {
			log.Println("Error reading dir", rootPath, ":", err)
		}

		for _, scriptInfo := range scripts {
			scriptPath := path.Join(rootPath, scriptInfo.Name())

			script, err := os.Open(scriptPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer script.Close()

			_, err = script.Read(first20)
			if err != nil {
				fmt.Println(err)
				continue
			}

			groups := shebang.FindStringSubmatch(string(first20))
			if len(groups) == 0 {
				fmt.Println(scriptPath, " is not a shell script >>", string(first20), "<<")
				continue
			}

			fmt.Println("groups=", groups)
			cmd := exec.Command(groups[1], "-n", scriptPath)

			out, err := cmd.CombinedOutput()
			if err != nil {
				errorList = append(errorList, fmt.Sprintln(scriptPath, "LINT failed:", err, "\n", string(out)))
			}
		}
	}
}
