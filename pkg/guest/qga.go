package guest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fjboy/magic-pocket/pkg/global/logging"
	"libvirt.org/go/libvirt"
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
type QACExecStatus struct {
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
func (guest *Guest) Exec(command string) (string, string) {
	qemuAgentCommand := QemuAgentCommand{
		Execute:   "guest-exec",
		Arguments: getGuestExecArguments(command),
	}
	jsonData, _ := json.Marshal(qemuAgentCommand)
	result, _ := guest.runQemuAgentCommand(jsonData)
	var qgaExecResult QgaExecResult
	json.Unmarshal([]byte(result), &qgaExecResult)
	return guest.getExecStatusOutput(qgaExecResult.Return.Pid)
}

func (guest *Guest) runQemuAgentCommand(jsonData []byte) (string, error) {
	logging.Debug("QGA 命令: %s", fmt.Sprintf("%s", jsonData))
	result, err := guest.domain.QemuAgentCommand(
		fmt.Sprintf("%s", jsonData),
		libvirt.DOMAIN_QEMU_AGENT_COMMAND_MIN,
		0)
	if err != nil {
		logging.Error("执行失败: %s", err)
		return "", err
	}
	logging.Debug("命令执行结果: %s", result)
	return result, nil
}

// guest-exec-status
func (guest *Guest) getExecStatusOutput(pid int) (string, string) {
	qemuAgentCommand := QACExecStatus{
		Execute:   "guest-exec-status",
		Arguments: getGuestExecStatusArguments(pid),
	}
	jsonData, _ := json.Marshal(qemuAgentCommand)
	var qgaExecResult QgaExecStatusResult
	for {
		result, _ := guest.runQemuAgentCommand(jsonData)
		json.Unmarshal([]byte(result), &qgaExecResult)
		if qgaExecResult.Return.Exited {
			break
		}
		time.Sleep(time.Second)
	}
	outDecode, _ := base64.StdEncoding.DecodeString(qgaExecResult.Return.OutData)
	errDecode, _ := base64.StdEncoding.DecodeString(qgaExecResult.Return.ErrData)
	return string(outDecode), string(errDecode)
}
