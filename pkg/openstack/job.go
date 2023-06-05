package openstack

import (
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

func TestNetQos(clientId string, serverId string) {
	authOpts, _ := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		logging.Error("初始化认证客户端失败, %s", err)
		return
	}
	ecClient, err := openstack.NewComputeV2(
		provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	if err != nil {
		logging.Error("初始化 compute v2 客户端失败, %s", err)
		return
	}
	pager := servers.List(ecClient, servers.ListOpts{})
	pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, _ := servers.ExtractServers(page)
		for _, s := range serverList {
			logging.Info("111111 %s", s.Name)
		}
		return true, nil
	})
	// server, id := servers.Get(ecClient, clientId).Extract()
	// logging.Info("xxx %s", server.Name)
	// logging.Info("%s", id)
}
