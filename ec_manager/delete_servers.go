package ec_manager

import (
	"fmt"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

func (manager *ECManager) DelErrorServers(waitDeleted bool, yes bool) {
	query := map[string]string{"status": "error"}
	logging.Info("查询虚拟机")
	servers := manager.computeClient.ServerList(query)
	logging.Info("状态为ERROR的虚拟机数量: %d\n", len(servers))
	if len(servers) == 0 {
		return
	}
	if !yes {
		var confirm string
		fmt.Println("即将删除虚拟机:")
		for _, server := range servers {
			fmt.Printf("%s (%s)\n", server.Id, server.Name)
		}
		for {
			fmt.Printf("是否删除 [yes/no]: ")
			fmt.Scanf("%s %d %f", &confirm)
			if confirm == "yes" || confirm == "y" {
				break
			} else if confirm == "no" || confirm == "n" {
				return
			} else {
				fmt.Printf("输入错误, 请重新输入!")
			}
		}
	}

	logging.Info("开始删除虚拟机")
	for _, server := range servers {
		logging.Info("删除虚拟机 %s(%s)", server.Id, server.Name)
		manager.computeClient.ServerDelete(server.Id)
		if waitDeleted {
			manager.computeClient.WaitServerDeleted(server.Id)
		}
	}
}
