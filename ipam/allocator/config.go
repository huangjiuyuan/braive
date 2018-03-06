package allocator

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
)

type IPAMConfig struct {
	*Range
	Name       string
	Type       string         `json:"type"`
	Routes     []*types.Route `json:"routes"`
	ResolvConf string         `json:"resolvConf"`
}

type Range struct {
	RangeStart net.IP    `json:"rangeStart,omitempty"`
	RangeEnd   net.IP    `json:"rangeEnd,omitempty"`
	Subnet     net.IPNet `json:"subnet"`
	Gateway    net.IP    `json:"gateway,omitempty"`
}

var (
	defaultRangeStart = "10.22.0.2"
	defaultRangeEnd   = "10.22.255.255"
	defaultSubnet     = "10.22.0.0/16"
	defaultGateway    = "10.22.0.1"
)

func LoadIPAllocator() (*IPAMConfig, error) {
	rangeStart := net.ParseIP(defaultRangeStart)
	if rangeStart == nil {
		return nil, fmt.Errorf("invalid value for rangeStart")
	}

	rangeEnd := net.ParseIP(defaultRangeEnd)
	if rangeEnd == nil {
		return nil, fmt.Errorf("invalid value for rangeEnd")
	}

	subnet, err := parseCIDR(defaultSubnet)
	if err != nil {
		return nil, fmt.Errorf("error on parsing subnet: %s", err)
	}

	gateway := net.ParseIP(defaultGateway)
	if gateway == nil {
		return nil, fmt.Errorf("invalid value for gateway")
	}

	config := &IPAMConfig{
		Range: &Range{
			RangeStart: rangeStart,
			RangeEnd:   rangeEnd,
			Subnet:     *subnet,
			Gateway:    gateway,
		},
		Name:       "braive",
		Type:       "braive-ipam",
		ResolvConf: "",
	}

	return config, nil
}

func parseCIDR(s string) (*net.IPNet, error) {
	ip, ipn, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	ipn.IP = ip
	return ipn, nil
}
