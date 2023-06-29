package compute

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/common"
)

func (cmpCli *ComputeClientV2) ServerList(query map[string]string) []Server {
	serversBody := ServersBody{}

	resp := cmpCli.AuthClient.Request(
		"GET", cmpCli.getUrl("servers", "", query),
		nil, query, cmpCli.BaseHeaders)
	json.Unmarshal(resp.Body, &serversBody)
	return serversBody.Servers
}

func (computeClient *ComputeClientV2) ServerShow(id string) (Server, error) {
	resp := computeClient.AuthClient.Request(
		"GET", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
	if err := resp.JudgeStatus(); err != nil {
		return Server{}, err
	}
	serverBody := ServerBody{}
	json.Unmarshal(resp.Body, &serverBody)
	return serverBody.Server, nil
}

func (computeClient *ComputeClientV2) ServerDelete(id string) {
	computeClient.AuthClient.Request(
		"DELETE", computeClient.getUrl("servers", id, nil), nil, nil, computeClient.BaseHeaders)
}

func (computeClient *ComputeClientV2) ServerCreate(options ServerOpt) Server {
	if options.Flavor == "" {
		options.Flavor = common.CONF.Ec.Flavor
	}
	if options.Image == "" {
		options.Image = common.CONF.Ec.Image
	}
	if options.Name == "" {
		options.Name = fmt.Sprintf("ecTools-server-%s", time.Now().Format("2006-01-02-15:04:05"))
	}
	if options.Networks == "" {
		options.Networks = "none"
	}

	body, _ := json.Marshal(ServeCreaterBody{Server: options})
	resp := computeClient.AuthClient.Request(
		"POST", computeClient.getUrl("servers", "", nil),
		body, nil, computeClient.BaseHeaders)
	serverBody := ServerBody{}
	json.Unmarshal(resp.Body, &serverBody)
	return serverBody.Server
}
func (client *ComputeClientV2) WaitServerCreate(options ServerOpt) Server {
	server := client.ServerCreate(options)
	return client.WaitServerStatusSecond(server.Id, "ACTIVE", 5)
}

func (client *ComputeClientV2) WaitServerStatusSecond(serverId string, status string, second int) Server {
	var server Server
	for {
		server, _ = client.ServerShow(serverId)
		logging.Debug("server stauts is %s", server.Status)
		if strings.ToUpper(server.Status) == strings.ToUpper(status) {
			break
		}
		time.Sleep(time.Second * time.Duration(second))
	}
	return server
}

func (client *ComputeClientV2) WaitServerStatus(serverId string, status string) Server {
	return client.WaitServerStatusSecond(serverId, status, 1)
}

func (client *ComputeClientV2) WaitServerDeleted(id string) {
	client.ServerDelete(id)
	for {
		server, err := client.ServerShow(id)
		logging.Debug("server status is %s", server.Status)
		if err != nil || server.Id != "" {
			break
		}
		time.Sleep(time.Second * 5)
	}
}
