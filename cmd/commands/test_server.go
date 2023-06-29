package commands

import (
	"github.com/BytemanD/ec-tools/pkg/openstack"
	"github.com/spf13/cobra"
)

var (
	times int
)
var TestServer = &cobra.Command{
	Use:   "test-server",
	Short: "云主机测试",
	Long:  "测试云主机创建/删除等操作",
	Run: func(cmd *cobra.Command, args []string) {
		openstack.TestServer(times)
	},
}

func init() {
	TestNetQos.Flags().IntVar(&times, "times", 0, "测试次数")
}
