package guest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"libvirt.org/go/libvirt"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type GuestExecArguments struct {
	CaptureOutput bool     `json:"capture-output"`
	Path          string   `json:"path"`
	Arg           []string `json:"arg"`
}
type GuestExecStatusArguments struct {
	Pid int `json:"pid"`
}

type QemuAgentCommand struct {
	Execute   string             `json:"execute"`
	Arguments GuestExecArguments `json:"arguments"`
}
type QGACExecStatus struct {
	Execute   string                   `json:"execute"`
	Arguments GuestExecStatusArguments `json:"arguments"`
}
type QgaExecReturn struct {
	Pid int `json:"pid"`
}
type QgaExecStatusReturn struct {
	Exited  bool   `json:"exited"`
	OutData string `json:"out-data"`
	ErrData string `json:"err-data"`
}
type QgaExecResult struct {
	Return QgaExecReturn `json:"return"`
}
type QgaExecStatusResult struct {
	Return QgaExecStatusReturn `json:"return"`
}
type QgaGFileOpen struct {
	Execute   string                 `json:"execute"`
	Arguments GuestFileOpenArguments `json:"arguments"`
}
type GuestFileOpenArguments struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
}
type QgaFileOpenReturn struct {
	Return int `json:"return"`
}
type QgaGFileWrite struct {
	Execute   string                  `json:"execute"`
	Arguments GuestFileWriteArguments `json:"arguments"`
}
type GuestFileWriteArguments struct {
	Handle int    `json:"handle"`
	Bufb64 string `json:"buf-b64"`
}
type QgaWriteOpenReturn struct {
	Return QgaWriteOpenReturnArguments `json:"return"`
}
type QgaWriteOpenReturnArguments struct {
	Count int  `json:"count"`
	Eof   bool `json:"eof"`
}
type QgaGFileClose struct {
	Execute   string                  `json:"execute"`
	Arguments GuestFileCloseArguments `json:"arguments"`
}
type GuestFileCloseArguments struct {
	Handle int `json:"handle"`
}

func getGuestExecArguments(command string) GuestExecArguments {
	commandArgs := strings.Split(command, " ")
	return GuestExecArguments{
		CaptureOutput: true,
		Path:          commandArgs[0],
		Arg:           commandArgs[1:],
	}
}
func getGuestExecStatusArguments(pid int) GuestExecStatusArguments {
	return GuestExecStatusArguments{
		Pid: pid,
	}
}

type ExecResult struct {
	Pid     int
	OutData string
	ErrData string
	Failed  bool
}

func (guest Guest) Exec(command string, wait bool) ExecResult {
	qemuAgentCommand := QemuAgentCommand{
		Execute:   "guest-exec",
		Arguments: getGuestExecArguments(command),
	}
	jsonData, _ := json.Marshal(qemuAgentCommand)
	result, err := guest.runQemuAgentCommand(jsonData)
	if err != nil {
		return ExecResult{Failed: true, ErrData: fmt.Sprintf("%s", err)}
	}
	var qgaExecResult QgaExecResult
	json.Unmarshal([]byte(result), &qgaExecResult)
	if !wait {
		return ExecResult{Pid: qgaExecResult.Return.Pid}
	}
	outData, errData := guest.getExecStatusOutput(qgaExecResult.Return.Pid)
	logging.Debug("OutData: %s, ErrData: %s", outData, errData)

	return ExecResult{
		Pid:     qgaExecResult.Return.Pid,
		OutData: outData,
		ErrData: errData,
	}
}

func (guest Guest) runQemuAgentCommand(jsonData []byte) (string, error) {
	logging.Debug("QGA 命令: %s", fmt.Sprintf("%s", jsonData))
	result, err := guest.domain.QemuAgentCommand(
		fmt.Sprintf("%s", jsonData), libvirt.DOMAIN_QEMU_AGENT_COMMAND_MIN, 0)
	if err != nil {
		logging.Error("执行失败: %s", err)
		return "", err
	}
	logging.Debug("命令执行结果: %s", result)
	return result, nil
}

// guest-exec-status
func (guest Guest) getExecStatusOutput(pid int) (string, string) {
	qemuAgentCommand := QGACExecStatus{
		Execute:   "guest-exec-status",
		Arguments: getGuestExecStatusArguments(pid),
	}
	jsonData, _ := json.Marshal(qemuAgentCommand)
	var qgaExecResult QgaExecStatusResult
	startTime := time.Now()
	for {
		result, _ := guest.runQemuAgentCommand(jsonData)
		json.Unmarshal([]byte(result), &qgaExecResult)
		if qgaExecResult.Return.Exited {
			break
		}
		if guest.QGATimeout > 0 &&
			time.Since(startTime).Seconds() >= float64(time.Second)*float64(guest.QGATimeout) {
			break
		}
		time.Sleep(time.Second * 1)
	}
	outDecode, _ := base64.StdEncoding.DecodeString(qgaExecResult.Return.OutData)
	errDecode, _ := base64.StdEncoding.DecodeString(qgaExecResult.Return.ErrData)
	return string(outDecode), string(errDecode)
}

func (guest Guest) GetIpaddrs() []string {
	execResult := guest.Exec("ip a", true)
	reg := regexp.MustCompile("inet ([0-9.]+)")
	matchedIPAddresses := reg.FindAllStringSubmatch(execResult.OutData, -1)
	ipAddresses := []string{}
	for _, matchedIPAddress := range matchedIPAddresses {
		if len(matchedIPAddress) < 2 || matchedIPAddress[1] == "127.0.0.1" {
			continue
		}
		ipAddresses = append(ipAddresses, matchedIPAddress[1])
	}
	return ipAddresses
}

func (guest Guest) Cat(args ...string) ExecResult {
	return guest.Exec(fmt.Sprintf("cat %s", strings.Join(args, " ")), true)
}

func (guest Guest) Kill(single int, pids []int) ExecResult {
	pidString := []string{}
	for _, pid := range pids {
		pidString = append(pidString, fmt.Sprintf("%d", pid))
	}
	return guest.Exec(
		fmt.Sprintf("kill -%d %s", single, strings.Join(pidString, " ")),
		true)
}

// Return pid
func (guest Guest) RunIperf3(args ...string) ExecResult {
	return guest.Exec(fmt.Sprintf("iperf3 %s", strings.Join(args, " ")), false)
}

func (guest Guest) RunIperfServer(serverIp string, logfile string, options string) ExecResult {
	return guest.RunIperf3(
		"-s", "--bind", serverIp, "--logfile", logfile, options,
	)
}

func (guest Guest) RunIperfClient(clientIp string, serverIp string, logfile string, options string) ExecResult {
	return guest.RunIperf3(
		"-c", serverIp, "--bind", clientIp, "--logfile", logfile, options,
	)
}
func (guest Guest) RunIperfClientUdp(clientIp string, serverIp string, logfile string, options string) ExecResult {
	return guest.RunIperf3(
		"-u", "-c", serverIp, "--bind", clientIp, "--logfile", logfile, options,
	)
}
func (guest Guest) HasCommand(command string) bool {
	execResult := guest.Exec(fmt.Sprintf("whereis %s", command), true)
	if execResult.Failed {
		return false
	}
	stringList := strings.Split(execResult.OutData, " ")
	if len(stringList) < 2 || stringList[1] == "" {
		return false
	}
	return true
}

func (guest Guest) RpmInstall(packagePath string) error {
	logging.Info("[%s] 安装 %v", guest.Domain, packagePath)
	result := guest.Exec(fmt.Sprintf("rpm -ivh %s", packagePath), true)
	if result.Failed {
		return fmt.Errorf("%s install failed, %s", packagePath, result.ErrData)
	}
	return nil
}

func (guest Guest) FileWrite(filePath string, content string) error {
	// file open
	logging.Debug("%s file open", filePath)
	fileOpenCommand := QgaGFileOpen{
		Execute:   "guest-file-open",
		Arguments: GuestFileOpenArguments{Path: filePath, Mode: "w+"},
	}
	jsonData, _ := json.Marshal(fileOpenCommand)
	result, _ := guest.runQemuAgentCommand(jsonData)
	var qgaFileOpenResult QgaFileOpenReturn
	json.Unmarshal([]byte(result), &qgaFileOpenResult)
	logging.Debug("file %s open handle: %d", filePath, qgaFileOpenResult.Return)
	if qgaFileOpenResult.Return == 0 {
		return fmt.Errorf("open file failed, return: %s", result)
	}
	// file write
	logging.Debug("%s file write", filePath)
	fileWriteCommand := QgaGFileWrite{
		Execute: "guest-file-write",
		Arguments: GuestFileWriteArguments{
			Handle: qgaFileOpenResult.Return,
			Bufb64: base64.StdEncoding.EncodeToString([]byte(content))},
	}
	jsonData, _ = json.Marshal(fileWriteCommand)
	result, _ = guest.runQemuAgentCommand(jsonData)
	var qgaFIleWriteResult QgaWriteOpenReturn
	json.Unmarshal([]byte(result), &qgaFIleWriteResult)
	logging.Debug("file %s write result %v", filePath, qgaFIleWriteResult)
	// file close
	logging.Debug("%s file close", filePath)
	fileCloseCommand := QgaGFileClose{
		Execute:   "guest-file-close",
		Arguments: GuestFileCloseArguments{Handle: qgaFileOpenResult.Return},
	}
	jsonData, _ = json.Marshal(fileCloseCommand)
	guest.runQemuAgentCommand(jsonData)
	return nil
}

func (guest Guest) CopyFile(localFile string, remotePath string) (string, error) {
	// TODO: 限制文件大小 <= 160k
	f, err := os.OpenFile(localFile, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(f)
	remoteFile := remotePath + "/" + filepath.Base(localFile)
	logging.Info("[%s] 拷贝文件 %s --> %s", guest.Domain, localFile, remotePath)
	return remoteFile, guest.FileWrite(remoteFile, string(bytes))
}
