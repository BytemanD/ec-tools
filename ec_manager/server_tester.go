package ec_manager

import (
	"os"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/common"
	"github.com/BytemanD/ec-tools/pkg/openstack/compute"
)

func (manager *ECManager) TestServer(times int) {
	CONF := common.CONF
	if CONF.Ec.Flavor == "" {
		logging.Fatal("flavor can not be null")
	}
	if CONF.Ec.Image == "" {
		logging.Fatal("image can not be null")
		os.Exit(1)
	}
	var testIimes int
	if times > 0 {
		testIimes = times
	} else {
		testIimes = CONF.TestServer.Times
	}
	if testIimes == 0 {
		logging.Fatal("test times must >= 1")
	}

	for i := 1; i <= CONF.TestServer.Times; i++ {
		logging.Info("create server %d", i)
		server, err := manager.computeClient.WaitServerCreate(compute.ServerOpt{})
		logging.Info("created server %s", server.Id)
		if err != nil {
			logging.Error("create server failed, %s", err)
			continue
		}
		logging.Info("delete server %s", server.Id)
		manager.computeClient.WaitServerDeleted(server.Id)
	}
}
