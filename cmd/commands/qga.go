package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/pkg/guest"
)

var (
	connection string
	uuid       bool
)
var QGACommand = &cobra.Command{
	Use:   "qga-exec <domain> <command>",
	Short: "QGA 命令执行工具",
	Long:  "执行 Libvirt QGA(qemu-guest-agent) 命令",
	Args:  cobra.MatchAll(cobra.ExactArgs(2)),
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
		execResult := domainGuest.Exec(command, true)
		if execResult.OutData != "" {
			fmt.Println(execResult.OutData)
		}
		if execResult.ErrData != "" {
			fmt.Println(execResult.ErrData)
		}
	},
}

var CopyCommand = &cobra.Command{
	Use:   "qga-copy <domain> <file> <remote path>",
	Short: "QGA 拷贝文件",
	Long:  "使用 Libvirt QGA 拷贝小文件",
	Args:  cobra.MatchAll(cobra.ExactArgs(3)),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		filePath := args[1]
		guestPath := args[2]
		connection, _ := cmd.Flags().GetString("connection")
		domainGuest := guest.Guest{
			Connection: connection, Domain: domainName}
		err := domainGuest.Connect()
		if err != nil {
			logging.Fatal("连接domain失败, %s", err)
		}
		guestFile, err := domainGuest.CopyFile(filePath, guestPath)
		if err != nil {
			logging.Fatal("copy file failed %s", err)
		} else {
			logging.Info("remote file is %s", guestFile)
		}
	},
}

func init() {
	QGACommand.Flags().StringVar(&connection, "connection", "localhost", "连接地址")
	QGACommand.Flags().BoolVarP(&uuid, "uuid", "u", false, "通过 UUID 查找")

	CopyCommand.Flags().String("connection", "localhost", "连接地址")
}
