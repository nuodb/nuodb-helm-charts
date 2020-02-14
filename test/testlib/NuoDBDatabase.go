package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type DBVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type NuoDBDatabase struct {
	Incarnation DBVersion `json:"incarnation"`
	Name string `json:"name"`
	Processes string `json:"processes"`
	State string `json:"state"`
}

func UnmarshalDatabase(s string) (err error, databases []NuoDBDatabase) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBDatabase
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, databases
		}

		if err != nil {
			return
		}

		databases = append(databases, obj)
	}
}