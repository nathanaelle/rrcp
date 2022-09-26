package rawsocket

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

type (
	socket struct {
		fd        int
		localAddr macIfaceAddr
	}
)

var realtekEtherType = htons(0x8899)

var _ net.PacketConn = &socket{}

func Socket(ifname string) (net.PacketConn, error) {
	sock := &socket{}

	gaddr, err := ByIfname(ifname)
	if err != nil {
		return nil, fmt.Errorf("getSocket: ByIfname [darwin] %w", err)
	}
	addr := gaddr.(*macIfaceAddr)
	sock.localAddr = *addr

	sock.fd, err = rawSocket()
	if err != nil {
		return nil, fmt.Errorf("getSocket: Socket() [darwin] %w", err)
	}

	if err := pktbind(sock); err != nil {
		return nil, err
	}

	return getSocket(sock)
}

func (s *socket) LocalAddr() net.Addr {
	r := s.localAddr
	return &r
}

func (s *socket) Close() error {
	return syscall.Close(s.fd)
}

func (*socket) SetDeadline(t time.Time) error {
	return nil
}

func (*socket) SetReadDeadline(t time.Time) error {
	return nil
}
func (*socket) SetWriteDeadline(t time.Time) error {
	return nil
}

func (s *socket) WriteTo(data []byte, remoteAddr net.Addr) (int, error) {
	return sendto(s, data, 0, remoteAddr)
}

func (s *socket) ReadFrom(data []byte) (n int, addr net.Addr, err error) {
	return recvfrom(s, data, syscall.MSG_TRUNC)
}
