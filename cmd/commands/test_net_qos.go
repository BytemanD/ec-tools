package commands

import (
	"github.com/spf13/cobra"

	"github.com/BytemanD/ec-tools/ec_manager"
)

var TestNetQos = &cobra.Command{
	Use:   "test-net-qos [server]",
	Short: "测试网络QOS",
	Long:  "基于iperf3 工具测试两个虚拟机的网络QOS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := cmd.Flags().GetString("client")
		ecManager := ec_manager.ECManager{}
		ecManager.Init()
		ecManager.TestNetQos(client, args[0])
	},
}

func init() {
	TestNetQos.Flags().String("client", "", "客户端虚拟机UUID")
	TestNetQos.MarkFlagRequired("client")
}
