package rawsocket

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"
)

type (
	RAWsockaddrNDRV struct {
		SndLen    uint8
		SndFamily uint8
		SndName   [14]byte
	}

	sockaddrNDRV struct {
		Name string
		raw  RAWsockaddrNDRV
	}
)

func (sa *sockaddrNDRV) Sockaddr() (unsafe.Pointer, uint32, error) {
	sa.raw.SndLen = 16
	sa.raw.SndFamily = syscall.AF_NDRV
	sa.raw.SndName = [14]byte{}
	bytes := []byte(sa.Name)
	copy(sa.raw.SndName[:], bytes)

	return unsafe.Pointer(&sa.raw), 16, nil
}

func toSockaddr(addr net.Addr) (Sockaddr, error) {
	maciface, ok := addr.(*macIfaceAddr)
	if !ok {
		return nil, fmt.Errorf("invalid type %#v", addr)
	}

	ndrv := &sockaddrNDRV{
		Name: maciface.name,
	}

	return ndrv, nil
}

func fromSockaddr(pc net.PacketConn, saddr Sockaddr) (net.Addr, error) {
	_, ok := pc.(*socket)
	if !ok {
		return nil, fmt.Errorf("invalid type %#v", pc)
	}

	// need implementation
	// so produce an error that dump the packet
	return nil, fmt.Errorf("%#v", saddr)
}

func rawSocket() (int, error) {
	return syscall.Socket(syscall.AF_NDRV, syscall.SOCK_RAW, int(realtekEtherType))
}

func pktbind(sock *socket) error {
	sa := &sockaddrNDRV{
		Name: sock.localAddr.name,
	}

	ptr, ptrlen, err := sa.Sockaddr()
	if err != nil {
		return err
	}

	_, _, errno := syscall.Syscall(syscall.SYS_BIND, uintptr(sock.fd), uintptr(ptr), uintptr(ptrlen))
	if errno == 0 {
		return nil
	}
	return os.NewSyscallError("SYS_BIND", syscall.Errno(errno))
}

func sendto(sock *socket, data []byte, flags int, remoteAddr net.Addr) (int, error) {
	var ptrlen uint32
	var ptr unsafe.Pointer

	saddr, err := toSockaddr(remoteAddr)
	if err != nil {
		return 0, fmt.Errorf("Sockaddr %w", err)
	}

	ptr, ptrlen, err = saddr.Sockaddr()
	if err != nil {
		return 0, err
	}

	_, _, errno := syscall.Syscall6(syscall.SYS_SENDTO, uintptr(sock.fd), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(flags), uintptr(ptr), uintptr(ptrlen))
	if errno == 0 {
		return len(data), nil
	}

	return 0, os.NewSyscallError("SYS_SENDTO", syscall.Errno(errno))
}

func recvfrom(sock *socket, data []byte, flags int) (n int, addr net.Addr, err error) {
	from := &sockaddrNDRV{}
	ptr, ptrlen, _ := from.Sockaddr()

	r0, _, errno := syscall.Syscall6(syscall.SYS_RECVFROM, uintptr(sock.fd), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(flags), uintptr(ptr), uintptr(ptrlen))
	if errno != 0 {
		return int(r0), nil, os.NewSyscallError("SYS_RECVFROM", syscall.Errno(errno))
	}
	n = int(r0)

	err = fmt.Errorf("todo : decode %q", from.Name)
	return
	/*
		addr, err = fromSockaddr(sock, from)
		if err != nil {
			n = 0
		}

		return
	*/
}
