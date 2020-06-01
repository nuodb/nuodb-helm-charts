package testlib

import (
	"encoding/json"
	"io"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type NuoDBKubeConfig struct {
	Pods map[string]corev1.Pod `json:"pods"`
	Deployments map[string]v1.Deployment `json:"deployments"`
	StatefulSets map[string]v1.StatefulSet `json:"statefulsets"`
	Volumes map[string]corev1.Volume `json:"volumes"`
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
