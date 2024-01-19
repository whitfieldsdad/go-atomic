package atomic

import (
	"time"

	"github.com/google/uuid"
	ps "github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	Id         string     `json:"id"`
	Time       time.Time  `json:"time"`
	PID        int        `json:"pid"`
	PPID       int        `json:"ppid"`
	Name       string     `json:"name,omitempty"`
	CWD        string     `json:"cwd,omitempty"`
	Argv       []string   `json:"argv,omitempty"`
	Argc       int        `json:"argc,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	ExitCode   *int       `json:"exit_code,omitempty"`
	Stdout     string     `json:"stdout,omitempty"`
	Stderr     string     `json:"stderr,omitempty"`
	User       *User      `json:"user,omitempty"`
	UserIds    []int      `json:"user_ids,omitempty"`
	GroupIds   []int      `json:"group_ids,omitempty"`
	Executable *File      `json:"executable,omitempty"`
}

func NewProcess(pid int) (*Process, error) {
	result := &Process{
		Id:   uuid.NewString(),
		Time: time.Now(),
		PID:  pid,
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
		argc int
	)
	argv, err = p.CmdlineSlice()
	if len(argv) != 0 {
		argc = len(argv)
	}

	var uids []int
	int32uids, _ := p.Uids()
	if len(int32uids) > 0 {
		for _, uid := range int32uids {
			uids = append(uids, int(uid))
		}
	}

	var gids []int
	int32gids, _ := p.Groups()
	if len(int32gids) > 0 {
		for _, gid := range int32gids {
			gids = append(gids, int(gid))
		}
	}

	var createTime *time.Time
	createTimeMs, err := p.CreateTime()
	if err == nil {
		t := time.UnixMilli(createTimeMs)
		createTime = &t
	}

	process := Process{
		Id:         uuid.NewString(),
		Time:       time.Now(),
		PID:        int(p.Pid),
		PPID:       int(ppid),
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
