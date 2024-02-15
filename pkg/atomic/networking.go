package atomic

import (
	"fmt"
	"net"
	"os"
)

func GetHostname() (string, error) {
	return os.Hostname()
}

func GetPrimaryIPv4Address() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func IsLoopbackAddress(ip string) bool {
	v := net.ParseIP(ip)
	return v != nil && v.IsLoopback()
}

func GetHostnameFromIPAddress(ipAddress string) (string, error) {
	addresses, err := net.LookupAddr(ipAddress)
	if err != nil {
		return "", err
	}
	fmt.Println(addresses)
	if len(addresses) == 0 {
		return "", nil
	}
	return addresses[0], nil
}

func GetIPv4AddressFromHostname(hostname string) (string, error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	for _, address := range addresses {
		if address.To4() != nil {
			return address.String(), nil
		}
	}
	return "", nil
}

func GetIPv6AddressFromHostname(hostname string) (string, error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	for _, address := range addresses {
		if address.To16() != nil {
			return address.String(), nil
		}
	}
	return "", nil
}
