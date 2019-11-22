package testlib

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

func IsStatefulSetHotCopyEnabled(ss *appsv1.StatefulSet) bool {
	return strings.Contains(ss.Name, "hotcopy")
}

func IsDaemonSetHotCopyEnabled(ss *appsv1.DaemonSet) bool {
	return strings.Contains(ss.Name, "hotcopy")
}
