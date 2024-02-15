package atomic

import "github.com/pkg/errors"

type NetworkLocation struct {
	Hostname      string   `json:"hostname"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	MACAddress    string   `json:"mac_address"`
}

func GetNetworkLocation() (*NetworkLocation, error) {
	hostname, _ := GetHostname()
	pnic, err := GetPrimaryNetworkInterface()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get primary network interface")
	}
	loc, err := GetNetworkLocationFromNetworkInterface(pnic.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get network location from primary network interface")
	}
	loc.Hostname = hostname
	return loc, nil
}
