package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/ec_manager"
)

var ServerPrune = &cobra.Command{
	Use:   "server-prune",
	Short: "清理状态为ERROR的虚拟机",
	Run: func(cmd *cobra.Command, args []string) {
		wait, _ := cmd.Flags().GetBool("wait")
		yes, _ := cmd.Flags().GetBool("yes")
		ecManager := ec_manager.ECManager{}
		if err := ecManager.Init(); err != nil {
			logging.Error("%s", err)
			os.Exit(1)
		}
		ecManager.DelErrorServers(wait, yes)
	},
}

func init() {
	ServerPrune.Flags().BoolP("wait", "w", false, "等待虚拟删除完成")
	ServerPrune.Flags().BoolP("yes", "y", false, "所有问题自动回答'是'")
}
