package ec_manager

import (
	"os"
	"strconv"
	"strings"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/BytemanD/ec-tools/pkg/guest"
	"github.com/BytemanD/ec-tools/pkg/openstack/compute"
)

func (manager *ECManager) TestNetQos(serverId string, clientId string, pps bool) {
	var (
		clientVm, serverVm compute.Server
	)
	if clientId == "" {
		// TODO
		logging.Info("创建客户端虚拟机")
		clientVm, _ = manager.computeClient.ServerCreate(compute.ServerOpt{})
		if clientVm.Id == "" {
			logging.Fatal("创建客户端虚拟机失败")
		}
	} else {
		clientVm, _ = manager.computeClient.ServerShow(clientId)
		if clientVm.Id == "" {
			logging.Fatal("虚拟机 %s 不存在", clientId)
		}
	}
	if serverId == "" {
		// TODO
		logging.Info("创建服务端虚拟机")
		serverVm, _ = manager.computeClient.ServerCreate(compute.ServerOpt{})
		if clientVm.Id == "" {
			logging.Fatal("创建服务端虚拟机失败")
		}
	} else {
		serverVm, _ = manager.computeClient.ServerShow(serverId)
		if serverVm.Id == "" {
			logging.Fatal("虚拟机 %s 不存在", serverId)
			return
		}
	}
	if strings.ToUpper(serverVm.Status) != "ACTIVE" {
		logging.Error("期望虚拟机 %s 状态是 ACTIVE, 实际是 %s", serverVm.Id, serverVm.Status)
		return
	}
	if strings.ToUpper(clientVm.Status) != "ACTIVE" {
		logging.Error("期望虚拟机 %s 状态是 ACTIVE, 实际是 %s", clientVm.Id, clientVm.Status)
		return
	}

	inboundKb, outboundKb := PrintVmQosSetting(clientVm, serverVm)

	clientConn := guest.GuestConnection{Connection: clientVm.Host, Domain: clientVm.Id}
	serverConn := guest.GuestConnection{Connection: serverVm.Host, Domain: serverVm.Id}

	logging.Info("启动QGA测试")
	senderTotal, receiverTotal, err := guest.TestNetQos(clientConn, serverConn, pps)
	if err != nil {
		logging.Fatal("测试失败, %s", err)
	}
	if inboundKb != 0 {
		logging.Info("出向带宽误差: %v %%", (senderTotal-inboundKb)*100.0/inboundKb)
	}
	if outboundKb != 0 {
		logging.Info("入向带宽误差: %v %%", (receiverTotal-outboundKb)*100/outboundKb)
	}
}
func PrintVmQosSetting(clientServer compute.Server, serverServer compute.Server) (float64, float64) {
	tableWriter := table.NewWriter()
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	header1 := table.Row{
		"Server", "Bandwidth(KBytes/sec)", "Bandwidth(KBytes/sec)", "PPS", "PPS",
	}
	header2 := table.Row{"", "inbound", "outbound", "inbound", "outbound"}
	tableWriter.AppendHeader(header1, rowConfigAutoMerge)
	tableWriter.AppendHeader(header2, rowConfigAutoMerge)
	tableWriter.AppendRow(
		table.Row{
			serverServer.Id + "(Server)",
			serverServer.Flavor.ExtraSpecs["quota:vif_inbound_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_outbound_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_inbound_pps_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_outbound_pps_burst"],
		})
	tableWriter.AppendRow(
		table.Row{
			clientServer.Id + "(Client)",
			clientServer.Flavor.ExtraSpecs["quota:vif_inbound_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_outbound_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_inbound_pps_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_outbound_pps_burst"],
		})

	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.SetStyle(table.StyleLight)
	tableWriter.Style().Format.Header = text.FormatDefault
	tableWriter.Style().Options.SeparateRows = true
	tableWriter.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignHeader: text.AlignCenter, Align: text.AlignCenter},
		{Number: 2, AlignHeader: text.AlignCenter, Align: text.AlignRight},
		{Number: 3, AlignHeader: text.AlignCenter, Align: text.AlignRight},
		{Number: 4, AlignHeader: text.AlignCenter, Align: text.AlignRight},
	})
	logging.Info("虚拟机QoS信息:")
	tableWriter.Render()

	inboundKB, _ := strconv.Atoi(clientServer.Flavor.ExtraSpecs["quota:vif_inbound_burst"])
	outboundKB, _ := strconv.Atoi(clientServer.Flavor.ExtraSpecs["quota:vif_outbound_burst"])
	return parseFlavorBandwidthToKb(float64(inboundKB)), parseFlavorBandwidthToKb(float64(outboundKB))
}
func parseFlavorBandwidthToKb(kB float64) float64 {
	return kB * 8 / 1024 / 1024 * 1000 * 1000
}
