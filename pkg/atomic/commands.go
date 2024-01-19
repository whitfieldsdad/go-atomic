package atomic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hairyhenderson/go-which"
	"github.com/pkg/errors"
)

var (
	WindowsPowerShell = "powershell"
	PowerShellCore    = "pwsh"
	PowerShell        = getPowerShellCommandType()
	CommandPrompt     = "command_prompt"
	Sh                = "sh"
	Bash              = "bash"
)

var (
	commandShims = map[string][]string{
		WindowsPowerShell: {"powershell", "-ExecutionPolicy", "Bypass", "-Command"},
		PowerShellCore:    {"pwsh", "-Command"},
		CommandPrompt:     {"cmd", "/c"},
		Sh:                {"sh", "-c"},
		Bash:              {"bash", "-c"},
	}
)

func getDefaultCommandType() string {
	if runtime.GOOS == "windows" {
		return CommandPrompt
	}
	return Bash
}

func getPowerShellCommandType() string {
	if runtime.GOOS == "windows" {
		return WindowsPowerShell
	}
	return PowerShellCore
}

func wrapShellCommand(command, commandType string) ([]string, error) {
	if commandType == "" {
		commandType = getDefaultCommandType()
	}
	argv := commandShims[commandType]
	if argv == nil {
		return nil, errors.Errorf("invalid command type: %s", commandType)
	}
	argv = append(argv, command)
	return argv, nil
}

type ShellCommand struct {
	Command        string                 `json:"command"`
	CommandType    string                 `json:"command_type"`
	InputArguments map[string]interface{} `json:"input_arguments,omitempty"`
}

func NewShellCommand(command, commandType string) ShellCommand {
	return NewShellCommandWithArgs(command, commandType, nil)
}

func NewShellCommandWithArgs(command, commandType string, args map[string]interface{}) ShellCommand {
	if commandType == "" {
		commandType = getDefaultCommandType()
	}
	return ShellCommand{
		Command:        command,
		CommandType:    commandType,
		InputArguments: args,
	}
}

func (c *ShellCommand) GetArgv() ([]string, error) {
	command := InterpolateCommandArgs(c.Command, c.InputArguments)
	argv, err := wrapShellCommand(command, c.CommandType)
	if err != nil {
		return nil, err
	}
	return argv, nil
}

func (c *ShellCommand) Exec(ctx context.Context) (*Process, error) {
	argv, err := c.GetArgv()
	if err != nil {
		return nil, err
	}
	return ExecArgv(ctx, argv)
}

func MergeInputArgs(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}
	return m
}

func InterpolateCommandArgs(command string, args map[string]interface{}) string {
	for k, v := range args {
		placeholder := fmt.Sprintf("#{%s}", k)
		if strings.Contains(command, placeholder) {
			command = strings.ReplaceAll(command, placeholder, fmt.Sprintf("%v", v))
		}
	}
	return command
}

type ExecutedCommand struct {
	StartTime  time.Time    `json:"start_time"`
	EndTime    time.Time    `json:"end_time"`
	Command    ShellCommand `json:"command"`
	User       User         `json:"user"`
	Subprocess *Process     `json:"subprocess,omitempty"`
	Error      error        `json:"error,omitempty"`
}

func ExecArgv(ctx context.Context, argv []string) (*Process, error) {
	command := strings.Join(argv, " ")
	log.Infof("Executing command: %s", command)

	path := which.Which(argv[0])
	if path == "" {
		return nil, errors.New("cannot execute command - executable not found")
	}
	cmd := exec.Command(path, argv[1:]...)
	cmd.SysProcAttr = getSysProcAttrs()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}
	log.Infof("Command started %s (PID: %d)", command, cmd.Process.Pid)

	// Lookup information about the process
	pid := cmd.Process.Pid
	process, err := GetProcess(pid)
	if err != nil {
		process, _ := NewProcess(pid)
		process.PPID = os.Getpid()
		process.Argv = argv
		process.Argc = len(argv)
		process.Executable, _ = GetFile(path)
	}

	log.Infof("Waiting for command to exit (PID: %d, PPID: %d)", process.PID, process.PPID)

	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for command to exit")
	}

	process.Stdout = stdout.String()
	process.Stderr = stderr.String()
	exitCode := cmd.ProcessState.ExitCode()
	process.ExitCode = &exitCode

	log.Infof("Command exited with code %d (PID: %d, PPID: %d)", exitCode, process.PID, process.PPID)
	return process, nil
}
