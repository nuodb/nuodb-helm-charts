package testlib

import (
	"encoding/json"
	"io"
	"strings"
)

type NuoDBServerState struct {
	State              string `json:"state"`
	Latency            int64  `json:"latency"`
	LastAckDeltaMillis int64  `json:"lastAckDeltaMillis"`
}

type NuoDBServerTermIndexInfo struct {
	CommitIndex  int64 `json:"commitIndex"`
	CurrentTerm  int64 `json:"currentTerm"`
	LogLastIndex int64 `json:"logLastIndex"`
	LogLastTerm  int64 `json:"logLastTerm"`
	Valid        bool  `json:"valid"`
}

type NuoDBServerRoleInfo struct {
	LeaderServerId         string                   `json:"leaderServerId"`
	LocalPeerTermIndexInfo NuoDBServerTermIndexInfo `json:"localPeerTermIndexInfo"`
	Role                   string                   `json:"role"`
}

type NuoDBServer struct {
	Address         string              `json:"address"`
	ConnectedState  NuoDBServerState    `json:"connectedState"`
	Id              string              `json:"id"`
	IsEvicted       bool                `json:"isEvicted"`
	IsLocal         bool                `json:"isLocal"`
	LocalRoleInfo   NuoDBServerRoleInfo `json:"localRoleInfo"`
	PeerMemberState string              `json:"peerMemberState"`
	PeerState       string              `json:"peerState"`
	Version         string              `json:"version"`
	Labels          map[string]string   `json:"labels"`
}

func UnmarshalDomainServers(s string) (err error, servers map[string]NuoDBServer) {
	dec := json.NewDecoder(strings.NewReader(s))
	servers = make(map[string]NuoDBServer)

	for {
		var obj NuoDBServer
		err = dec.Decode(&obj)
		if err == io.EOF {
			// all done
			return nil, servers
		}

		if err != nil {
			return
		}

		servers[obj.Id] = obj
	}
}
