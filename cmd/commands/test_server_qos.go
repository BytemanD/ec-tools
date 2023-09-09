package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/BytemanD/ec-tools/ec_manager"
)

var TestNetQos = &cobra.Command{
	Use:     "server-iperf3-test [server] --client [client]",
	Short:   "测试云主机网络QOS",
	Long:    "基于iperf3工具测试两个虚拟机的网络QOS",
	Example: strings.TrimRight(SERVER_IPERF3_TEST_EXAMPLE, "\n"),
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := cmd.Flags().GetString("client")
		pps, _ := cmd.Flags().GetBool("pps")

		ecManager := ec_manager.ECManager{}
		ecManager.Init()
		ecManager.TestNetQos(args[0], client, pps)
	},
}

func init() {
	TestNetQos.Flags().String("client", "", "客户端虚拟机UUID")
	TestNetQos.Flags().Bool("pps", false, "测试PPS")
	TestNetQos.MarkFlagRequired("client")
}
