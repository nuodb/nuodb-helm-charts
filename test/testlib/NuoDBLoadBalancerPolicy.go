package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBLoadBalancerPolicy struct {
	LbQuery    string `json:"lbQuery"`
	PolicyName string `json:"policyName"`
}

func UnmarshalLoadBalancerPolicies(s string) (err error, policies map[string]NuoDBLoadBalancerPolicy) {
	dec := json.NewDecoder(strings.NewReader(s))
	policies = make(map[string]NuoDBLoadBalancerPolicy)

	for {
		var obj NuoDBLoadBalancerPolicy
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, policies
		}

		if err != nil {
			return
		}

		policies[obj.PolicyName] = obj
	}
}
