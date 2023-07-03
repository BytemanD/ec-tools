package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/BytemanD/ec-tools/ec_manager"
	"github.com/BytemanD/ec-tools/pkg/guest"
)

func getGuestConnection(guestAddr string) guest.GuestConnection {
	addrList := strings.Split(guestAddr, ":")
	if len(addrList) == 2 {
		return guest.GuestConnection{
			Connection: addrList[0],
			Domain:     addrList[1],
		}
	} else {
		return guest.GuestConnection{
			Domain: addrList[0],
		}
	}
}

var (
	client   string
	server   string
	directly bool
)
var TestNetQos = &cobra.Command{
	Use:   "test-net-qos",
	Short: "测试网络QOS",
	Long:  "基于iperf3 工具测试两个虚拟机的网络QOS",
	Run: func(cmd *cobra.Command, args []string) {
		if directly {
			clientConn := getGuestConnection(client)
			serverConn := getGuestConnection(server)
			guest.TestNetQos(clientConn, serverConn)
		} else {
			ecManager := ec_manager.ECManager{}
			ecManager.Init()
			ecManager.TestNetQos(client, server)
		}
	},
}

func init() {
	TestNetQos.Flags().StringVarP(&client, "client", "c", "", "客户端虚拟机UUID")
	TestNetQos.Flags().StringVarP(&server, "server", "s", "", "服务端虚拟机UUID")
	TestNetQos.Flags().BoolVar(&directly, "directly", false, "直接通过 QGA 开始测试(不调用计算服务API)")
}
