package openstack

import (
	"github.com/fjboy/ec-tools/pkg/openstack/identity"
	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

func TestNetQos(clientId string, serverId string) {
	client, err := identity.GetV3ClientFromEnv()
	if err != nil {
		logging.Error("获取客户端失败, %s", err)
		return
	}

	// TODO: 业务逻辑
	logging.Info("services %s", client.ServiceList())
	logging.Info("users %s", client.UserList())

}
