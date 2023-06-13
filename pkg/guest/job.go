package guest

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type GuestConnection struct {
	Connection string
	Domain     string
}

type Job struct {
	SourceIp string
	DestIp   string
	Pid      int
	Logfile  string
	Output   string
	Sender   string
	Receiver string
}

// 使用 iperf3 工具测试网络限速
//
// 参数为客户端和服务端虚拟机的连接消息，格式: "连接地址:虚拟机 UUID"。例如：
//
//	192.168.192.168:a6ee919a-4026-4f0b-8d7e-404950a91eb3
func TestNetQos(clientConn GuestConnection, serverConn GuestConnection) {
	clientGuest := Guest{
		Connection: clientConn.Connection,
		Domain:     clientConn.Domain,
		QGATimeout: 60,
		ByUUID:     true}
	serverGuest := Guest{Connection: serverConn.Connection, Domain: serverConn.Domain, ByUUID: true}
	err := clientGuest.Connect()
	logging.Info("连接客户端虚拟机: %s", clientGuest.Domain)
	if err != nil {
		logging.Error("连接客户端虚拟机失败, %s", err)
		return
	}
	logging.Info("连接服务端虚拟机: %s", serverGuest.Domain)
	err = serverGuest.Connect()
	if err != nil {
		logging.Error("连接服务端虚拟机失败, %s", err)
		return
	}
	logging.Info("获取客户端和服务端虚拟机IP地址")
	clientAddresses := clientGuest.GetIpaddrs()
	serverAddresses := serverGuest.GetIpaddrs()

	logging.Info("客户端虚拟机IP地址: %s", clientAddresses)
	logging.Info("服务端虚拟机IP地址: %s", serverAddresses)

	if len(clientAddresses) == 0 || len(serverAddresses) == 0 {
		logging.Fatal("客户端和服务端虚拟机必须至少有一张启用的网卡")
	}

	fomatTime := time.Now().Format(time.RFC3339)
	serverPids := []int{}
	for _, serverAddress := range serverAddresses {
		logfile := fmt.Sprintf("/tmp/iperf3_s_%s_%s", fomatTime, serverAddresses)
		logging.Info("启动服务端: %s", serverAddress)
		execResult := serverGuest.RunIperfServer(serverAddress, logfile)
		if execResult.Failed {
			return
		}
		serverPids = append(serverPids, execResult.Pid)
	}
	if len(serverPids) > 0 {
		defer serverGuest.Kill(9, serverPids)
	}

	jobs := []Job{}

	for i := 0; i < len(clientAddresses) && i < len(serverAddresses); i++ {
		logfile := fmt.Sprintf("/tmp/iperf3_c_%s_%s", fomatTime, serverAddresses[i])
		logging.Info("启动客户端: %s -> %s", clientAddresses[i], serverAddresses[i])
		execResult := clientGuest.RunIperfClient(clientAddresses[i], serverAddresses[i], logfile)
		jobs = append(jobs, Job{
			SourceIp: clientAddresses[i],
			DestIp:   serverAddresses[i],
			Pid:      execResult.Pid,
			Logfile:  logfile,
		})
	}

	logging.Info("等待测试结果")
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	tableWriter := table.NewWriter()

	tableWriter.AppendHeader(
		table.Row{"Client -> Server", "Bandwidth(KBytes/sec)", "Bandwidth(KBytes/sec)"},
		rowConfigAutoMerge,
	)
	tableWriter.AppendHeader(
		table.Row{"", "Sender", "Receiver"},
		rowConfigAutoMerge, rowConfigAutoMerge, rowConfigAutoMerge,
	)

	senderReg := regexp.MustCompile(" +([0-9.]+) KBytes/sec.* +sender")
	receiverReg := regexp.MustCompile(" +([0-9.]+) KBytes/sec .* +receiver")
	var (
		senderTotal   int
		receiverTotal int
	)
	for _, job := range jobs {
		// 等待命令执行结束
		clientGuest.getExecStatusOutput(job.Pid)
		// 获取 sender 和 receiver
		execResult := clientGuest.Cat(job.Logfile)
		matchedSender := senderReg.FindStringSubmatch(execResult.OutData)
		matchedReceiver := receiverReg.FindStringSubmatch(execResult.OutData)
		if len(matchedSender) >= 2 {
			job.Sender = matchedSender[1]
			number, _ := strconv.Atoi(job.Sender)
			senderTotal += number
		}
		if len(matchedReceiver) >= 2 {
			job.Receiver = matchedReceiver[1]
			number, _ := strconv.Atoi(job.Receiver)
			receiverTotal += number
		}
		tableWriter.AppendRow(
			table.Row{
				fmt.Sprintf("%s -> %s", job.SourceIp, job.DestIp),
				job.Sender, job.Receiver,
			})
	}
	tableWriter.AppendFooter(table.Row{"Total", senderTotal, receiverTotal})

	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.SetAutoIndex(true)
	tableWriter.SetStyle(table.StyleLight)
	tableWriter.Style().Format.Header = text.FormatDefault
	tableWriter.Style().Options.SeparateRows = true
	tableWriter.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, AlignHeader: text.AlignCenter, Align: text.AlignCenter, AlignFooter: text.AlignCenter},
		{Number: 2, AutoMerge: true, AlignHeader: text.AlignCenter, Align: text.AlignRight, AlignFooter: text.AlignRight},
		{Number: 3, AutoMerge: true, AlignHeader: text.AlignCenter, Align: text.AlignRight, AlignFooter: text.AlignRight},
	})
	tableWriter.Render()
}
