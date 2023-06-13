package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/gitutils"
	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/cmd/commands"
)

var (
	debug   bool
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
		Short:   "常用工具合集",
		Long:    "Golang 实现的EC工具合集",
		Version: getVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := logging.INFO
			if debug {
				level = logging.DEBUG
			}
			logging.BasicConfig(logging.LogConfig{Level: level})
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "显示Debug信息")

	rootCmd.AddCommand(commands.QGACommand)
	rootCmd.AddCommand(commands.TestNetQos)
	rootCmd.AddCommand(commands.DelErrorServers)
	rootCmd.AddCommand(commands.InitConf)
	rootCmd.Execute()
}
