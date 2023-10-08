package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/gitutils"
	"github.com/BytemanD/easygo/pkg/global/logging"

	"github.com/BytemanD/ec-tools/cmd/commands"
	"github.com/BytemanD/ec-tools/common"
)

var (
	Version string
)

func getVersion() string {
	if Version == "" {
		return gitutils.GetVersion()
	}
	return fmt.Sprint(Version)
}

func main() {
	rootCmd := cobra.Command{
		Use:     "ec-tools",
		Short:   "EC常用工具集",
		Long:    "Golang 实现的EC工具集",
		Version: getVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")
			conf, _ := cmd.Flags().GetString("conf")
			level := logging.INFO
			if debug {
				level = logging.DEBUG
			}
			logging.BasicConfig(logging.LogConfig{Level: level})
			err := common.LoadConf(conf)
			if err != nil && cmd.Name() != commands.DumpConf.Name() {
				logging.Warning("加载配置文件失败, %s", err)
			}
			if !debug && common.CONF.Debug {
				logging.BasicConfig(logging.LogConfig{Level: logging.DEBUG})
			}
			common.LogLines()
		},
	}

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "显示Debug信息")
	rootCmd.PersistentFlags().StringP("conf", "c", "", "自定义配置文件")

	rootCmd.AddCommand(
		commands.QGACommand, commands.CopyCommand,
		commands.TestNetQos,
		commands.DumpConf, commands.TestServer,
		commands.GuestIperf3Test,
	)
	rootCmd.Execute()
}
