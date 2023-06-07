package commands

import (
	"github.com/fjboy/ec-tools/pkg/openstack"
	"github.com/spf13/cobra"
)

var DelErrorServers = &cobra.Command{
	Use:   "delete-error-servers",
	Short: "删除错误的虚拟机",
	Long:  "删除状态为ERROR的虚拟机",
	Run: func(cmd *cobra.Command, args []string) {
		openstack.DelErrorServers()
	},
}

func init() {
	TestNetQos.Flags().BoolVar(&directly, "watch", false, "等待虚拟删除完成")
}
