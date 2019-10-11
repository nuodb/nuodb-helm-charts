package integration

import (
	"encoding/json"
	"github.com/nuodb/consulting-helm/test/testlib"
	"gotest.tools/assert"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	s := (
		`{
  "address": "172.17.0.6", 
  "dbName": "demo", 
  "durableState": "MONITORED", 
  "host": "admin-kljkzo-nuodb-0", 
  "hostname": "te-database-rggfj1-nuodb-demo-57f984dcd5-s7nhr", 
  "ipAddress": "172.17.0.5", 
  "isExternalStartup": true, 
  "labels": {
    "cloud": "minikube", 
    "region": "local", 
    "zone": "local-b"
  }, 
  "lastHeardFrom": 3611, 
  "nodeId": 2, 
  "options": {
    "agent": "admin-kljkzo-nuodb-0.nuodb.testkubernetesbasicdatabase-kljkzo.svc:48004", 
    "commit": "safe", 
    "database": "demo", 
    "engine-type": "TE", 
    "ext-start": "true", 
    "geo-region": "0", 
    "log-over-conn": "enable", 
    "max-lost-archives": "0", 
    "mem": "8000000000", 
    "node-port": "48006", 
    "ping-timeout": "60", 
    "read-stdin": "true", 
    "region-name": "Default", 
    "user-data": "{\"incarnation\":{\"major\":1,\"minor\":0},\"startId\":\"1\",\"labels\":{\"cloud\":\"minikube\",\"zone\":\"local-b\",\"region\":\"local\"}}", 
    "verbose": "error,flush,warn"
  }, 
  "pid": 39, 
  "port": 48006, 
  "regionId": 0, 
  "regionName": "Default", 
  "startId": "1", 
  "state": "RUNNING", 
  "storageGroups": "https://localhost:8888/api/1/processes/1/storageGroups", 
  "type": "TE", 
  "version": "4.0.1-1"
}
`)

	err, objects := testlib.Unmarshal(s)

	assert.NilError(t, err)
	assert.Equal(t, len(objects), 1)

	obj := objects[0]

	assert.Check(t, obj.Hostname == "te-database-rggfj1-nuodb-demo-57f984dcd5-s7nhr")
	assert.Check(t, obj.Host == "admin-kljkzo-nuodb-0")
	assert.Check(t, obj.DbName == "demo")

	_, ok := obj.Labels["cloud"]
	assert.Check(t, ok)

	_, ok = obj.Labels["region"]
	assert.Check(t, ok)

	_, ok = obj.Labels["zone"]
	assert.Check(t, ok)

}

func TestUnmarshallMany(t *testing.T) {
	s := (
		`{
  "address": "172.17.0.7", 
  "archiveDir": "/var/opt/nuodb/archive/nuodb/demo", 
  "archiveId": 0, 
  "dbName": "demo", 
  "durableState": "MONITORED", 
  "host": "admin-kljkzo-nuodb-0", 
  "hostname": "sm-database-rggfj1-nuodb-demo-hotcopy-0", 
  "ipAddress": "172.17.0.5", 
  "isExternalStartup": true, 
  "labels": {
    "cloud": "minikube", 
    "region": "local", 
    "zone": "local-b"
  }
}
{
  "address": "172.17.0.6", 
  "dbName": "demo", 
  "durableState": "MONITORED", 
  "host": "admin-kljkzo-nuodb-0", 
  "hostname": "te-database-rggfj1-nuodb-demo-57f984dcd5-s7nhr", 
  "ipAddress": "172.17.0.5", 
  "isExternalStartup": true, 
  "labels": {
    "cloud": "minikube", 
    "region": "local", 
    "zone": "local-b"
  }
}

`)
	err, objects := testlib.Unmarshal(s)
	assert.NilError(t, err)

	assert.Check(t, len(objects) == 2)

	for _, obj := range objects {
		assert.Check(t, obj.Host == "admin-kljkzo-nuodb-0")
		assert.Check(t, obj.DbName == "demo")

		_, ok := obj.Labels["cloud"]
		assert.Check(t, ok)

		_, ok = obj.Labels["region"]
		assert.Check(t, ok)

		_, ok = obj.Labels["zone"]
		assert.Check(t, ok)
	}

}

func TestMarshall(t *testing.T) {

	labels := map[string]string {
		"cloud": "minikube",
		"region": "local",
		"zone": "local-b",
	}

	obj := testlib.NuoDBProcess {
		Address: "172.17.0.6",
		DbName: "demo",
		Host: "admin-kljkzo-nuodb-0",
		Hostname: "te-database-rggfj1-nuodb-demo-57f984dcd5-s7nhr",
		Labels: labels,
	}
	b, err := json.Marshal(&obj)

	assert.NilError(t, err)

	t.Log(string(b))
}