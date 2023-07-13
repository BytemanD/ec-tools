package commands

import (
	"github.com/BytemanD/ec-tools/ec_manager"
	"github.com/spf13/cobra"
)

var (
	times int
)
var TestServer = &cobra.Command{
	Use:   "server-test",
	Short: "云主机测试",
	Long:  "测试云主机创建/删除等操作",
	Run: func(cmd *cobra.Command, args []string) {
		ecManager := ec_manager.ECManager{}
		ecManager.Init()
		ecManager.TestServer(times)
	},
}

func init() {
	TestServer.Flags().IntVar(&times, "times", 0, "测试次数")
}
