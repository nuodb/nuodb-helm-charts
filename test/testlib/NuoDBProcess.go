package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBProcess struct {
	Address  string            `json:"address"`
	DbName   string            `json:"dbName"`
	Type     string            `json:"type"`
	Host     string            `json:"host"`
	Hostname string            `json:"hostname"`
	Labels   map[string]string `json:"labels"`
}

func Unmarshal(s string) (err error, processes []NuoDBProcess) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBProcess
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, processes
		}

		if err != nil {
			return
		}

		processes = append(processes, obj)
	}
}
