package atomic

import (
	"os"
	"github.com/charmbracelet/log"
)

type Host struct {
	Hostname string `json:"hostname"`
	OS       OS     `json:"os"`
}

func GetHost() Host {
	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Failed to get hostname: %s", hostname)
	}
	return Host{
		Hostname: hostname,
		OS:       GetOS(),
	}
}
