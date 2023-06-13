package openstack

import (
	"fmt"
	"os"
	"strings"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/BytemanD/ec-tools/pkg/guest"
	"github.com/BytemanD/ec-tools/pkg/openstack/compute"
	"github.com/BytemanD/ec-tools/pkg/openstack/identity"

	"github.com/BytemanD/ec-tools/common"
)

func getAuthedClient() (compute.ComputeClientV2, error) {
	authClient, err := identity.GetV3ClientFromEnv()
	if err != nil {
		logging.Error("获取认证客户端失败, %s", err)
		return compute.ComputeClientV2{}, fmt.Errorf("获取计算客户端失败")
	}
	computeClient, err := compute.GetComputeClientV2(authClient)
	if err != nil {
		logging.Error("获取计算客户端失败, %s", err)
		return compute.ComputeClientV2{}, fmt.Errorf("获取计算客户端失败")
	}
	computeClient.UpdateVersion()
	return computeClient, nil
}

func PrintVmQosSetting(clientServer compute.Server, serverServer compute.Server) {
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
			clientServer.Id + "(Client)",
			clientServer.Flavor.ExtraSpecs["quota:vif_inbound_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_outbound_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_inbound_pps_burst"],
			clientServer.Flavor.ExtraSpecs["quota:vif_outbound_pps_burst"],
		})
	tableWriter.AppendRow(
		table.Row{
			serverServer.Id + "(Server)",
			serverServer.Flavor.ExtraSpecs["quota:vif_inbound_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_outbound_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_inbound_pps_burst"],
			serverServer.Flavor.ExtraSpecs["quota:vif_outbound_pps_burst"],
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
}

func loadOpenrc() error {
	if common.CONF.Ec.AuthOpenrc == "" {
		return fmt.Errorf("authOpenrc is null")
	}
	lines, err := common.ReadLines(common.CONF.Ec.AuthOpenrc)
	if err != nil {
		return err
	}
	for _, line := range lines {
		cols := strings.Split(line, " ")
		if len(cols) != 2 {
			continue
		}
		envValue := strings.Split(cols[1], "=")
		if len(envValue) != 2 {
			continue
		}
		os.Setenv(envValue[0], envValue[1])
	}
	return nil
}

func initConfig() {
	if err := loadOpenrc(); err != nil {
		logging.Fatal("导入环境变量失败, %s", err)
	}
}

func TestNetQos(clientId string, serverId string) {
	initConfig()
	computeClient, err := getAuthedClient()
	if err != nil {
		logging.Fatal("获取认证客户端失败, %s", err)
	}
	var (
		clientVm, serverVm compute.Server
	)
	if clientId == "" {
		// TODO
		logging.Info("创建客户端虚拟机")
		clientVm = computeClient.ServerCreate(compute.ServerCreate{})
		if clientVm.Id == "" {
			logging.Fatal("创建客户端虚拟机失败")
		}
	} else {
		clientVm = computeClient.ServerShow(clientId)
		if clientVm.Id == "" {
			logging.Fatal("虚拟机 %s 不存在", clientId)
		}
	}
	if serverId == "" {
		// TODO
		logging.Info("创建服务端虚拟机")
		serverVm = computeClient.ServerCreate(compute.ServerCreate{})
		if clientVm.Id == "" {
			logging.Fatal("创建服务端虚拟机失败")
		}
	} else {
		serverVm = computeClient.ServerShow(serverId)
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
	if strings.ToUpper(serverVm.Status) != "ACTIVE" {
		logging.Error("期望虚拟机 %s 状态是 ACTIVE, 实际是 %s", serverVm.Id, serverVm.Status)
		return
	}

	PrintVmQosSetting(clientVm, serverVm)

	clientConn := guest.GuestConnection{Connection: clientVm.Host, Domain: clientVm.Id}
	serverConn := guest.GuestConnection{Connection: serverVm.Host, Domain: serverVm.Id}

	logging.Info("开始通过 QGA 测试")
	guest.TestNetQos(clientConn, serverConn)
}

func DelErrorServers() {
	initConfig()
	computeClient, err := getAuthedClient()
	if err != nil {
		return
	}
	query := map[string]string{}
	query["status"] = "error"
	logging.Info("查询虚拟机")
	servers := computeClient.ServerList(query)
	if len(servers) == 0 {
		logging.Warning("无状态为ERROR的虚拟机")
		return
	}
	logging.Info("开始删除虚拟机")
	for _, server := range servers {
		logging.Info("删除虚拟机 %s(%s)", server.Id, server.Name)
		computeClient.ServerDelete(server.Id)
	}
}
