package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fjboy/ec-tools/pkg/guest"
	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

var (
	connection string
	uuid       bool
)
var QGACommand = &cobra.Command{
	Use:   "qga-exec <domain> <command>",
	Short: "QGA 命令执行工具",
	Long:  "执行 Libvirt QGA(qemu-guest-agent) 命令",
	Args:  cobra.ExactValidArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		command := args[1]
		domainGuest := guest.Guest{
			Connection: connection,
			Domain:     domainName,
			ByUUID:     uuid,
		}
		err := domainGuest.Connect()
		if err != nil {
			logging.Error("连接domain失败 %s", err)
			return
		}
		outData, errData := domainGuest.Exec(command)
		if outData != "" {
			fmt.Println(outData)
		}
		if errData != "" {
			fmt.Println(errData)
		}
	},
}

func init() {
	QGACommand.Flags().StringVarP(&connection, "connection", "c", "localhost", "连接地址")
	QGACommand.Flags().BoolVarP(&uuid, "uuid", "u", false, "通过 UUID 查找")
}
