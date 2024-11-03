package testlib

import (
	"encoding/json"
	"io"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type NuoDBKubeConfig struct {
	Pods         map[string]corev1.Pod         `json:"pods"`
	Deployments  map[string]appsv1.Deployment  `json:"deployments"`
	StatefulSets map[string]appsv1.StatefulSet `json:"statefulsets"`
	Volumes      map[string]corev1.Volume      `json:"volumes"`
	DaemonSets   map[string]appsv1.DaemonSet   `json:"daemonSets"`
}

func UnmarshalNuoDBKubeConfig(s string) (err error, kubeConfigs []NuoDBKubeConfig) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBKubeConfig
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, kubeConfigs
		}

		if err != nil {
			return
		}

		kubeConfigs = append(kubeConfigs, obj)
	}
}
