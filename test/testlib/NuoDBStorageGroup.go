package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBStorageGroup struct {
	Id               int               `json:"sgId"`
	Name             string            `json:"sgName"`
	DbName           string            `json:"dbName"`
	State            string            `json:"state"`
	ArchiveStates    map[string]string `json:"archiveStates"`
	ProcessStates    map[string]string `json:"processStates"`
	LeaderCandidates []string          `json:"leaderCandidates"`
}

func UnmarshalStorageGroups(s string) (err error, sgs []NuoDBStorageGroup) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBStorageGroup
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, sgs
		}

		if err != nil {
			return
		}

		sgs = append(sgs, obj)
	}
}
