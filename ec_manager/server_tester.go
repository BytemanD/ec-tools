package ec_manager

import (
	"os"
	"sync"

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
	var wg sync.WaitGroup
	workers := make(chan struct{}, 1)
	if common.CONF.TestServer.Workers != 0 {
		workers = make(chan struct{}, common.CONF.TestServer.Workers)
	}
	wg.Add(CONF.TestServer.Times)
	for i := 1; i <= CONF.TestServer.Times; i++ {
		workers <- struct{}{}
		go func() {
			serverTest(*manager)
			<-workers
			wg.Done()
		}()
	}
	wg.Wait()
	logging.Info("all test finished")
}

func serverTest(manager ECManager) {
	logging.Info("%d create server", os.Getegid())
	server, err := manager.computeClient.WaitServerCreate(compute.ServerOpt{})
	if err != nil {
		logging.Error("[server: %s] create server failed, %s", server.Id, err)
		if server.Id != "" {
			manager.computeClient.WaitServerDeleted(server.Id)
		}
		return
	}
	logging.Info("[server: %s] created", server.Id)
	logging.Info("[server: %s] deleting", server.Id)
	manager.computeClient.WaitServerDeleted(server.Id)
	logging.Info("[server: %s] deleted", server.Id)
}
