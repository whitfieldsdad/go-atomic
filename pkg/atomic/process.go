package atomic

import (
	"time"

	ps "github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	Artifact   `json:",inline"`
	PID        int32      `json:"pid"`
	PPID       int32      `json:"ppid"`
	Name       string     `json:"name,omitempty"`
	CWD        string     `json:"cwd,omitempty"`
	Argv       []string   `json:"argv,omitempty"`
	Argc       int32      `json:"argc,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	ExitCode   *int       `json:"exit_code,omitempty"`
	Stdout     string     `json:"stdout,omitempty"`
	Stderr     string     `json:"stderr,omitempty"`
	User       *User      `json:"user,omitempty"`
	UserIds    []int32    `json:"user_ids,omitempty"`
	GroupIds   []int32    `json:"group_ids,omitempty"`
	Executable *File      `json:"executable,omitempty"`
}

func NewProcess(pid int32) (*Process, error) {
	result := &Process{
		Artifact: NewArtifact(),
		PID:      pid,
	}
	return result, nil
}

func GetProcess(pid int) (*Process, error) {
	p, err := ps.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	info := parseProcess(p)
	return &info, nil
}

func ListProcesses() ([]Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	var rows []Process
	for _, p := range processes {
		row := parseProcess(p)
		rows = append(rows, row)
	}
	return rows, nil
}

func parseProcess(p *ps.Process) Process {
	name, _ := p.Name()
	cwd, _ := p.Cwd()
	ppid, _ := p.Ppid()

	var executable *File
	exe, _ := p.Exe()
	if exe != "" {
		executable, _ = GetFile(exe)
	}

	var user *User
	username, err := p.Username()
	if err == nil {
		user, _ = GetUserByUsername(username)
	}

	var (
		argv []string
		argc int32
	)
	argv, _ = p.CmdlineSlice()
	if len(argv) != 0 {
		argc = int32(len(argv))
	}

	uids, _ := p.Uids()
	gids, _ := p.Groups()

	var createTime *time.Time
	createTimeMs, err := p.CreateTime()
	if err == nil {
		t := time.UnixMilli(createTimeMs)
		createTime = &t
	}

	process := Process{
		Artifact:   NewArtifact(),
		PID:        p.Pid,
		PPID:       ppid,
		Name:       name,
		CWD:        cwd,
		Argv:       argv,
		Argc:       argc,
		CreateTime: createTime,
		Executable: executable,
		User:       user,
		UserIds:    uids,
		GroupIds:   gids,
	}
	return process
}

func IsElevated() bool {
	return isElevated()
}
