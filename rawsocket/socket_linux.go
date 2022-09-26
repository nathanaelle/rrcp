package rawsocket

import (
	"fmt"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func getSocket(sock *socket) (recv net.PacketConn, err error) {
	pmreq := &unix.PacketMreq{
		Ifindex: int32(sock.localAddr.ifindex),
		Type:    syscall.PACKET_MR_PROMISC,
	}
	copy(pmreq.Address[:], sock.localAddr.mac[:])

	err = unix.SetsockoptPacketMreq(sock.fd, syscall.SOL_PACKET, syscall.PACKET_ADD_MEMBERSHIP, pmreq)
	if err := solError(sock.fd, "packet_add_membership", err); err != nil {
		return nil, fmt.Errorf("getSocket: [linux] %w", err)
	}

	return sock, nil
}
