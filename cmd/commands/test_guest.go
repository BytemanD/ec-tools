package commands

import (
	"strings"

	"github.com/BytemanD/ec-tools/pkg/guest"
	"github.com/spf13/cobra"
)

func getGuestConnection(guestAddr string) guest.GuestConnection {
	addrList := strings.Split(guestAddr, ":")
	if len(addrList) == 2 {
		return guest.GuestConnection{
			Connection: addrList[0],
			Domain:     addrList[1],
		}
	} else {
		return guest.GuestConnection{
			Domain: addrList[0],
		}
	}
}

var TestGuestBps = &cobra.Command{
	Use:     "test-guest-bps <server>",
	Short:   "测试实例BPS",
	Long:    "基于 iperf3 工具测试实例的网络BPS",
	Args:    cobra.ExactArgs(1),
	Example: "guest-bps-test <hostA>:<guest-uuid1> --client <hostB>:<guest-uuid2>",
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := cmd.Flags().GetString("client")

		serverConn := getGuestConnection(args[0])
		clientConn := getGuestConnection(client)
		guest.TestNetQos(clientConn, serverConn)
	},
}

func init() {
	TestGuestBps.Flags().String("client", "", "客户端实例UUID")
	TestGuestBps.MarkFlagRequired("client")
}
