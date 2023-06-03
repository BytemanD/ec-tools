package commands

import (
	"github.com/spf13/cobra"

	"github.com/fjboy/ec-tools/pkg/guest"
	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

var (
	client string
	server string
)
var TestNetQos = &cobra.Command{
	Use:   "test-net-qos",
	Short: "测试网络QOS",
	Long:  "基于iperf3 工具测试两个虚拟机的网络QOS",
	Run: func(cmd *cobra.Command, args []string) {
		if client == "" || server == "" {
			logging.Error("非法参数, client 和 server 不能为空")
			return
		}
		guest.TestNetQos(client, server)
	},
}

func init() {
	TestNetQos.Flags().StringVarP(&client, "client", "c", "", "客户端虚拟机UUID")
	TestNetQos.Flags().StringVarP(&server, "server", "s", "", "服务端虚拟机UUID")
}