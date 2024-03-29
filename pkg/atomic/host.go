package atomic

import (
	"os"
	"runtime"

	"github.com/NextronSystems/go-osversion"
	"github.com/charmbracelet/log"
	"github.com/denisbrodbeck/machineid"
)

const (
	appId = "a86148b4-19e2-4533-a16f-3f3e96e92848"
)

type Host struct {
	Artifact `json:",inline"`
	Hostname string `json:"hostname"`
	OS       OS     `json:"os"`
}

func GetHost() Host {
	id := GetHostId()
	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Failed to get hostname: %s", hostname)
	}
	return Host{
		Artifact: Artifact{
			Id:   id,
		},
		Hostname: hostname,
		OS:       GetOS(),
	}
}

var _hostId string

func GetHostId() string {
	if _hostId == "" {
		id, err := machineid.ID()
		if err != nil {
			log.Fatalf("Failed to get host ID")
		}
		_hostId = NewUUID5(appId, []byte(id))
	}
	return _hostId
}

type OS struct {
	Type     string `json:"type"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
}

func GetOS() OS {
	version, err := osversion.Get()
	if err != nil {
		log.Fatalf("Failed to get OS version: %s", err)
	}
	return OS{
		Type:     runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  version,
	}
}
