package main

import (
	"github.com/spf13/cobra"

	"github.com/fjboy/ec-tools/cmd/commands"
	"github.com/fjboy/magic-pocket/pkg/global/gitutils"
	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

var (
	debug bool
)

func main() {
	rootCmd := cobra.Command{
		Use:     "ec-tools",
		Short:   "常用工具合集",
		Long:    "Golang 实现的EC工具合集",
		Version: gitutils.GetVersion(),
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
	rootCmd.Execute()
}
