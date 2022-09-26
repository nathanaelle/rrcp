package rawsocket

import (
	"fmt"
	"net"
)

type (
	macIfaceAddr struct {
		mac     [6]byte
		name    string
		ifindex int
	}
)

var _ net.Addr = &macIfaceAddr{}

func ByIfname(ifname string) (net.Addr, error) {
	ifr, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	var mac [6]byte
	hwaddr := []byte(ifr.HardwareAddr)
	copy(mac[:], hwaddr)

	return &macIfaceAddr{
		ifindex: ifr.Index,
		name:    ifr.Name,
		mac:     mac,
	}, nil
}

func (*macIfaceAddr) Network() string {
	return "ether(0x8899)"
}

func (addr *macIfaceAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x@%s[%d]",
		addr.mac[0], addr.mac[1], addr.mac[2], addr.mac[3], addr.mac[4], addr.mac[5],
		addr.name, addr.ifindex,
	)
}
