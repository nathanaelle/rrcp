package rawsocket

import (
	"net"
)

func getSocket(sock *socket) (recv net.PacketConn, err error) {
	return sock, nil
}
