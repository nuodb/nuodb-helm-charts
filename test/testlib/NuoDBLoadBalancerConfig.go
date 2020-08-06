package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBLoadBalancerConfig struct {
	DbName         string `json:"dbName"`
	DefaultLbQuery string `json:"defaultLbQuery"`
	Prefilter      string `json:"prefilter"`
	IsGlobal       bool   `json:"isGlobal"`
}

func UnmarshalLoadBalancerConfigs(s string) (err error, lbConfigs []NuoDBLoadBalancerConfig) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBLoadBalancerConfig
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, lbConfigs
		}

		if err != nil {
			return
		}

		lbConfigs = append(lbConfigs, obj)
	}
}
