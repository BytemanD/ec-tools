package compute

import (
	"encoding/json"
	"fmt"

	"github.com/fjboy/ec-tools/pkg/openstack/identity"
)

type VersionBody struct {
	Version ServerVersion `json:"version"`
}
type ServerVersion struct {
	MinVersion string `json:"min_version"`
	Version    string `json:"version"`
}

type ComputeClientV2 struct {
	AuthClient    identity.V3AuthClient
	Endpoint      string
	ServerVersion ServerVersion
	BaseHeaders   map[string]string
}

func (computeClient *ComputeClientV2) getUrl(resource string, id string, query map[string]string) string {
	url := fmt.Sprintf("%s/%s", computeClient.Endpoint, resource)
	if id != "" {
		url += "/" + id
	}
	var queryString string
	for key, value := range query {
		queryString += fmt.Sprintf("%s=%s", key, value)
	}

	if queryString != "" {
		return url + "?" + queryString
	} else {
		return url
	}
}

// X-OpenStack-Nova-API-Version
func (computeClient *ComputeClientV2) UpdateVersion() {
	version := computeClient.AuthClient.Request("GET", computeClient.Endpoint, nil, nil, nil)
	versionBody := VersionBody{}
	json.Unmarshal([]byte(version), &versionBody)
	computeClient.BaseHeaders = map[string]string{}
	computeClient.ServerVersion = versionBody.Version
	computeClient.BaseHeaders["X-OpenStack-Nova-API-Version"] = computeClient.ServerVersion.Version
}

func (cmpCli *ComputeClientV2) ServerList(query map[string]string) []Server {
	serversBody := ServersBody{}

	content := cmpCli.AuthClient.Request(
		"GET", cmpCli.getUrl("servers", "", query),
		nil, query, cmpCli.BaseHeaders)
	json.Unmarshal([]byte(content), &serversBody)
	return serversBody.Servers
}

func (computeClient *ComputeClientV2) ServerShow(id string) Server {
	content := computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
	serverBody := ServerBody{}
	json.Unmarshal([]byte(content), &serverBody)
	return serverBody.Server
}

func (computeClient *ComputeClientV2) ServerDelete(id string) {
	computeClient.AuthClient.Request(
		"DELETE", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
}
