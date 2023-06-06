package compute

import (
	"encoding/json"
	"fmt"

	"github.com/fjboy/ec-tools/pkg/openstack/identity"
)

type Server struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	VmState    string `json:"OS-EXT-STS:vm_state"`
	PowerState string `json:"OS-EXT-STS:power_state"`
	Host       string `json:"OS-EXT-SRV-ATTR:host"`
	AZ         string `json:"OS-EXT-AZ:availability_zone"`
}

type ServerBody struct {
	Server Server `json:"server"`
}

type ComputeClientV2 struct {
	AuthClient identity.V3AuthClient
	Endpoint   string
}

func (computeClient *ComputeClientV2) getUrl(resource string, id string) string {
	url := fmt.Sprintf("%s/%s", computeClient.Endpoint, resource)
	if id != "" {
		url += "/" + id
	}
	return url
}

func (computeClient *ComputeClientV2) ServerList() string {
	return computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", ""), nil)
}

func (computeClient *ComputeClientV2) ServerShow(id string) Server {
	content := computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", id), nil)
	serverBody := ServerBody{}
	json.Unmarshal([]byte(content), &serverBody)
	return serverBody.Server
}
