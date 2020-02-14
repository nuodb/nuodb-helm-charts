package testlib

import (
	"gopkg.in/yaml.v2"
)

type Registry struct {
	Nuodb struct {
		Image struct {
			Registry string
			Repository string
			Tag string
		}
	}
}

// UnmarshalImageYAML is used to unmarshal into map[string]string
func UnmarshalImageYAML(s string) (err error, registry Registry){
	registry = Registry{}

	err = yaml.Unmarshal([]byte(s), &registry)

	return
}