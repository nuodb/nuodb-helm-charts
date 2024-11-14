package testlib

import (
	"github.com/ghodss/yaml"
)

type Registry struct {
	Nuodb struct {
		Image struct {
			Registry   string
			Repository string
			Tag        string
		}
	}
}

// UnmarshalImageYAML is used to unmarshal into map[string]string
func UnmarshalImageYAML(s string) (Registry, error) {
	registry := Registry{}
	err := yaml.Unmarshal([]byte(s), &registry)
	return registry, err
}
