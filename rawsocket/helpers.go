package rawsocket

import (
	"fmt"
	"net"
	"unsafe"
)

func Broadcast(addr net.Addr) (net.Addr, error) {
	macAddr, ok := addr.(*macIfaceAddr)
	if !ok {
		return nil, fmt.Errorf("invalid type %#v", addr)
	}

	return &macIfaceAddr{
		ifindex: macAddr.ifindex,
		name:    macAddr.name,
		mac:     [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}, nil
}

func MacAddr(stuff any) ([6]byte, error) {
	switch typed := any(stuff).(type) {
	case *socket:
		return typed.localAddr.mac, nil
	case *macIfaceAddr:
		return typed.mac, nil
	}

	return [6]byte{}, fmt.Errorf("invalid type %#v", stuff)
}

func Addr(stuff any, mac [6]byte) (net.Addr, error) {
	switch typed := any(stuff).(type) {
	case *socket:
		return &macIfaceAddr{
			ifindex: typed.localAddr.ifindex,
			name:    typed.localAddr.name,
			mac:     mac,
		}, nil

	case *macIfaceAddr:
		return &macIfaceAddr{
			ifindex: typed.ifindex,
			name:    typed.name,
			mac:     mac,
		}, nil
	}

	return nil, fmt.Errorf("invalid type %#v", stuff)
}

func htons(x uint16) uint16 {
	t := [2]byte{}
	*(*uint16)(unsafe.Pointer(&t[0])) = uint16(0x1234)

	if t[0] == 0x34 {
		x1 := byte(x)
		x2 := byte(x >> 8)
		return uint16(x1)<<8 | uint16(x2)
	}

	return x
}
