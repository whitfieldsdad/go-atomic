package atomic

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/net"
)

type NetworkConnection struct {
	LocalHost  string `json:"local_host"`
	LocalPort  int64  `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int64  `json:"remote_port"`
	Pid        int64  `json:"pid"`
	Status     string `json:"status"`
}

func ListNetworkConnections() ([]*NetworkConnection, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list network connections")
	}
	connections := []*NetworkConnection{}
	for _, conn := range conns {
		connections = append(connections, &NetworkConnection{
			LocalHost:  conn.Laddr.IP,
			LocalPort:  int64(conn.Laddr.Port),
			RemoteHost: conn.Raddr.IP,
			RemotePort: int64(conn.Raddr.Port),
			Pid:        int64(conn.Pid),
			Status:     conn.Status,
		})
	}
	return connections, nil
}
