package rawsocket

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func rawSocket() (int, error) {
	return syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(realtekEtherType))
}

func pktbind(sock *socket) error {
	sa := &syscall.SockaddrLinklayer{
		Ifindex:  sock.localAddr.ifindex,
		Protocol: realtekEtherType,
		Pkttype:  syscall.PACKET_HOST,
		Hatype:   syscall.ARPHRD_NETROM,
	}

	return solError(sock.fd, "bind", syscall.Bind(sock.fd, sa))
}

func recvfrom(sock *socket, data []byte, flags int) (n int, addr net.Addr, err error) {
	var from syscall.Sockaddr

	n, from, err = syscall.Recvfrom(sock.fd, data, syscall.MSG_TRUNC)
	if err != nil {
		return
	}

	addr, err = fromSockaddr(sock, from)
	if err != nil {
		n = 0
	}

	return
}

func sendto(sock *socket, data []byte, flags int, remoteAddr net.Addr) (int, error) {
	saddr, err := toSockaddr(remoteAddr)
	if err != nil {
		return 0, fmt.Errorf("Sockaddr %w", err)
	}

	err = syscall.Sendto(sock.fd, data, flags, saddr)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func solError(fd int, syscallName string, oldErr error) error {
	if oldErr != nil {
		return os.NewSyscallError(syscallName, oldErr)
	}

	errno, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_ERROR)
	if err == nil && errno == 0 {
		return nil
	}

	if err != nil {
		return os.NewSyscallError(syscallName, err)
	}

	return os.NewSyscallError(syscallName, syscall.Errno(errno))
}

func toSockaddr(addr net.Addr) (syscall.Sockaddr, error) {
	maciface, ok := addr.(*macIfaceAddr)
	if !ok {
		return nil, fmt.Errorf("invalid type %#v", addr)
	}

	sall := &syscall.SockaddrLinklayer{
		Protocol: realtekEtherType,
		Ifindex:  maciface.ifindex,
		Halen:    uint8(len(maciface.mac)),
	}
	copy(sall.Addr[:], maciface.mac[:])

	return sall, nil
}

func fromSockaddr(sock *socket, saddr syscall.Sockaddr) (net.Addr, error) {
	scaddr, ok := saddr.(*syscall.SockaddrLinklayer)
	if !ok {
		return nil, fmt.Errorf("invalid type %#v", saddr)
	}

	var mac [6]byte
	copy(mac[:], scaddr.Addr[:])

	addr := &macIfaceAddr{
		ifindex: scaddr.Ifindex,
		mac:     mac,
		name:    sock.localAddr.name,
	}

	return addr, nil
}
