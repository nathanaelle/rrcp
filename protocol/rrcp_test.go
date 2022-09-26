package protocol

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"

	"github.com/nathanaelle/rrcp/rawsocket"
)

const ifname = "eth0"

func TestDecode(t *testing.T) {
	sock, err := rawsocket.Socket(ifname)
	if err != nil {
		t.Errorf("getSocket() : %v\n", err)
		return
	}
	if sock == nil {
		t.Errorf("sock is nil !!\n")
		return
	}

	defer sock.Close()

	buff := make([]byte, 1<<7)

	/*
		if sock.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
			t.Errorf("deadline: %v", err)
		}
	*/
	size, _, err := sock.ReadFrom(buff)
	if err != nil {
		t.Errorf("sock.Read() : len %d %v\n", size, err)
		return
	}
	t.Logf("yeah ! read %d bytes\n", size)

	header := &HeaderPkt{}
	err = binary.Read(bytes.NewReader(buff[:size]), binary.BigEndian, header)
	if err != nil {
		t.Errorf("binary.Read() : len %d %v\n", size, err)
		return
	}

	t.Logf("%#v\n", header)
}

func TestRRCP(t *testing.T) {
	sock, err := rawsocket.Socket(ifname)
	if err != nil {
		t.Errorf("getSocket() : %v\n", err)
		return
	}
	if sock == nil {
		t.Errorf("sock is nil !!\n")
		return
	}

	defer sock.Close()

	hwaddr, err := rawsocket.MacAddr(sock)
	if err != nil {
		t.Errorf("MacAddr: %v", err)
		return
	}

	//hello := ReqHello.Header([6]byte{0x28, 0x87, 0xba, 0x55, 0x68, 0x8e}, hwaddr)
	hello := ReqHello.Header(macbroadcast, hwaddr)

	t.Logf("%s", hello)

	size, err := hello.Sendto(sock)
	if err != nil || size != defaultPktLen {
		t.Errorf("sock.Write() : len %d %v\n", size, err)
		return
	}

	t.Logf("yeah ! send %d bytes\n", size)
	readHello(t, sock)
	readHello(t, sock)
	readHello(t, sock)
}

func readHello(t *testing.T, sock net.PacketConn) {
	rbuff := make([]byte, 1<<7)

	size, _, err := sock.ReadFrom(rbuff)
	if err != nil {
		t.Errorf("sock.Read() : len %d %v\n", size, err)
		return
	}
	t.Logf("\nyeah ! read %d bytes\n", size)

	header := &HeaderPkt{}

	size, err = header.Write(rbuff)
	if err != nil {
		t.Errorf("header.Write() : header len %d %v\n", size, err)
		return
	}

	data := header.Data()

	err = binary.Read(bytes.NewReader(rbuff), binary.BigEndian, data)
	if err != nil {
		t.Errorf("binary.Read() : %v\n", err)
		return
	}
	t.Logf("%s\n", data)
}
