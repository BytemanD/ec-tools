package guest

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

type GuestConnection struct {
	Connection string
	Domain     string
}

func getGuestConnection(guestAddr string) GuestConnection {
	addrList := strings.Split(guestAddr, ":")
	if len(addrList) == 2 {
		return GuestConnection{
			Connection: addrList[0],
			Domain:     addrList[1],
		}
	} else {
		return GuestConnection{
			Domain: addrList[0],
		}
	}
}

func TestNetQos(client string, server string) {
	clientConn := getGuestConnection(client)
	serverConn := getGuestConnection(server)

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
	logging.Info("获取客户端虚拟机IP地址")
	clientAddresses := clientGuest.GetIpaddrs()
	logging.Info("客户端虚拟机IP地址: %s", clientAddresses)

	logging.Info("获取服务端虚拟机IP地址")
	serverAddresses := serverGuest.GetIpaddrs()
	logging.Info("获取服务端虚拟机IP地址 %s", serverAddresses)

	fomatTime := time.Now().Format(time.RFC3339)

	for i := 0; i < len(serverAddresses); i++ {
		logfile := fmt.Sprintf("/tmp/iperf3_s_%s_%s", fomatTime, serverAddresses[i])
		logging.Info("启动服务端: %s", serverAddresses[i])
		serverGuest.RunIperfServer(serverAddresses[i], logfile)
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
	sender := regexp.MustCompile(" +([0-9.]+ +[a-zA-Z]+/sec) .* +sender")
	receiver := regexp.MustCompile(" +([0-9.]+ +[a-zA-Z]+/sec) .* +receiver")
	for i := 0; i < len(jobs); i++ {
		clientGuest.getExecStatusOutput(jobs[i].Pid)
		execResult := clientGuest.Exec(fmt.Sprintf("cat %s", jobs[i].Logfile), true)
		jobs[i].Sender = sender.FindString(execResult.OutData)
		jobs[i].Receiver = receiver.FindString(execResult.OutData)
	}

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"client -> server", "Sender", "Receiver"})
	splitSpace := regexp.MustCompile(" +|\n+")

	for i := 0; i < len(jobs); i++ {
		tableWriter.AppendRow(
			[]interface{}{
				fmt.Sprintf("%s -> %s", jobs[i].SourceIp, jobs[i].DestIp),
				strings.Join(splitSpace.Split(jobs[i].Sender, -1)[1:3], " "),
				strings.Join(splitSpace.Split(jobs[i].Receiver, -1)[1:3], " "),
			})
	}
	tableWriter.Render()

}
