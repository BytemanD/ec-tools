package compute

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BytemanD/ec-tools/common"
)

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

func (computeClient *ComputeClientV2) ServerCreate(options ServerCreate) Server {
	if options.Flavor == "" {
		options.Flavor = common.CONF.Ec.Flavor
	}
	if options.Image == "" {
		options.Image = common.CONF.Ec.Image
	}
	if options.Name == "" {
		options.Name = fmt.Sprintf("ecTools-vm-%s", time.Now().Format("2006-01-02-15:04:05"))
	}
	if options.Networks == "" {
		options.Networks = "none"
	}

	serverCreateBody := ServeCreaterBody{
		Server: options,
	}

	body, _ := json.Marshal(serverCreateBody)
	content := computeClient.AuthClient.Request(
		"POST", computeClient.getUrl("servers", "", nil),
		body, nil, computeClient.BaseHeaders)
	serverBody := ServerBody{}
	json.Unmarshal([]byte(content), &serverBody)
	return serverBody.Server
}
