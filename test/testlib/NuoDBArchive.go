package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBArchive struct {
	Id     int    `json:"id"`
	DbName string `json:"dbName"`
	Path   string `json:"path"`
	State  string `json:"state"`
}

func UnmarshalArchives(s string) (err error, archives []NuoDBArchive) {
	dec := json.NewDecoder(strings.NewReader(s))

	for {
		var obj NuoDBArchive
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, archives
		}

		if err != nil {
			return
		}

		archives = append(archives, obj)
	}
}
