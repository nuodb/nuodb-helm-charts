package testlib

import (
	"github.com/gruntwork-io/terratest/modules/k8s"
	"gotest.tools/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const DEBUG_POD = `
apiVersion: v1
kind: Pod
metadata:
  name: debug-pod
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
      - name: log-te-volume
        mountPath: /var/log/nuodb
  volumes:
  - name: log-te-volume
    persistentVolumeClaim:
      claimName: log-te-volume
`

func RecoverCoresFromTEs(t *testing.T, namespaceName string) {
	pwd, err := os.Getwd()
	assert.NilError(t, err)

	targetDirPath := filepath.Join(pwd, RESULT_DIR, namespaceName, "cores")
	_ = os.MkdirAll(targetDirPath, 0700)

	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	k8s.KubectlApplyFromString(t, kubectlOptions, DEBUG_POD)

	AwaitPodUp(t, namespaceName, "debug-pod", 30* time.Second)

	k8s.RunKubectl(t, kubectlOptions, "exec", "debug-pod", "--", "ls", "-lah", "/var/log/nuodb/crash")

	k8s.RunKubectl(t, kubectlOptions, "cp", "debug-pod:/var/log/nuodb/crash", targetDirPath)
}
