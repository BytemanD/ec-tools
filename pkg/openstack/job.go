package openstack

import (
	"strings"

	"github.com/fjboy/magic-pocket/pkg/global/logging"

	"github.com/fjboy/ec-tools/pkg/guest"
	"github.com/fjboy/ec-tools/pkg/openstack/compute"
	"github.com/fjboy/ec-tools/pkg/openstack/identity"
)

func TestNetQos(clientId string, serverId string) {
	authClient, err := identity.GetV3ClientFromEnv()
	if err != nil {
		logging.Error("获取认证客户端失败, %s", err)
		return
	}
	computeClient, err := compute.GetComputeClientV2(authClient)
	if err != nil {
		logging.Error("获取计算客户端失败, %s", err)
		return
	}
	logging.Info("查询客户端和服务端虚拟机信息")
	clientVm := computeClient.ServerShow(clientId)
	serverVm := computeClient.ServerShow(serverId)
	if strings.ToUpper(clientVm.Status) != "ACTIVE" {
		logging.Error("期望虚拟机 %s 状态是 ACTIVE, 实际是 %s", clientVm.Id, clientVm.Status)
		return
	}
	if strings.ToUpper(serverVm.Status) != "ACTIVE" {
		logging.Error("期望虚拟机 %s 状态是 ACTIVE, 实际是 %s", serverVm.Id, serverVm.Status)
		return
	}

	clientConn := guest.GuestConnection{Connection: clientVm.Host, Domain: clientVm.Id}
	serverConn := guest.GuestConnection{Connection: serverVm.Host, Domain: serverVm.Id}

	logging.Info("开始通过 QGA 测试")
	guest.TestNetQos(clientConn, serverConn)
}
