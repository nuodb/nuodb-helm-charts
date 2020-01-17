package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBDatabase struct {
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