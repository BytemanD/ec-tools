package compute

import (
	"encoding/json"
	"fmt"

	"github.com/fjboy/ec-tools/pkg/openstack/identity"
)

type Flavor struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	OriginalName string            `json:"original_name"`
	Ram          int               `json:"ram"`
	Vcpus        int               `json:"vcpus"`
	ExtraSpecs   map[string]string `json:"extra_specs"`
}

type Server struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	VmState    string `json:"OS-EXT-STS:vm_state"`
	PowerState string `json:"OS-EXT-STS:power_state"`
	Host       string `json:"OS-EXT-SRV-ATTR:host"`
	AZ         string `json:"OS-EXT-AZ:availability_zone"`
	Flavor     Flavor `json:"flavor"`
}

type ServerBody struct {
	Server Server `json:"server"`
}
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

func (computeClient *ComputeClientV2) getUrl(resource string, id string) string {
	url := fmt.Sprintf("%s/%s", computeClient.Endpoint, resource)
	if id != "" {
		url += "/" + id
	}
	return url
}

// X-OpenStack-Nova-API-Version
func (computeClient *ComputeClientV2) UpdateVersion() {
	version := computeClient.AuthClient.Request("GET", computeClient.Endpoint, nil, nil)
	versionBody := VersionBody{}
	json.Unmarshal([]byte(version), &versionBody)
	computeClient.BaseHeaders = map[string]string{}
	computeClient.ServerVersion = versionBody.Version
	computeClient.BaseHeaders["X-OpenStack-Nova-API-Version"] = computeClient.ServerVersion.Version
}

func (computeClient *ComputeClientV2) ServerList() string {
	return computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", ""), nil, computeClient.BaseHeaders)
}

func (computeClient *ComputeClientV2) ServerShow(id string) Server {
	content := computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", id), nil, computeClient.BaseHeaders)
	serverBody := ServerBody{}
	json.Unmarshal([]byte(content), &serverBody)
	return serverBody.Server
}
