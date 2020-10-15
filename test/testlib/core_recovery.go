package testlib

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"testing"
	"time"
)

const DEBUG_POD = `
apiVersion: v1
kind: Pod
metadata:
  name: %s
  labels:
    app: nuodb
    group: nuodb
spec:
  containers:
    - name: wait
      image: docker.io/busybox:latest
      imagePullPolicy: IfNotPresent
      command: [ "/bin/sh", "-c", "--" ]
      args: [ "while true; do sleep 30; done;" ]
      volumeMounts:
      - name: log-volume
        mountPath: /var/log/nuodb
  volumes:
  - name: log-volume
    persistentVolumeClaim:
      claimName: %s
`

func RecoverCoresFromEngine(t *testing.T, namespaceName string, engineType string, pvcName string) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	targetDirPath := filepath.Join(pwd, RESULT_DIR, namespaceName, "cores", engineType)
	_ = os.MkdirAll(targetDirPath, 0700)

	debugPodName := fmt.Sprintf("%s-debug-pod", engineType)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.KubectlApplyFromString(t, kubectlOptions, fmt.Sprintf(DEBUG_POD, debugPodName, pvcName))

	AwaitPodUp(t, namespaceName, debugPodName, 30* time.Second)

	k8s.RunKubectl(t, kubectlOptions, "exec", debugPodName, "--", "ls", "-lah", "/var/log/nuodb/")

	k8s.RunKubectl(t, kubectlOptions, "cp", debugPodName+":/var/log/nuodb/", targetDirPath)
}
