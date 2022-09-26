package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/nathanaelle/rrcp/rawsocket"
)

type (
	ReqHelloGetSetPkt struct {
		Header       HeaderPkt
		AuthKey      uint16
		Register     uint16
		RegisterData uint32
		Reserved     [8]byte
		ResZero      [28]byte
	}
)

func (rhgsp *ReqHelloGetSetPkt) String() string {
	switch rhgsp.Header.Opcode {
	case ReqHello:
		return fmt.Sprintf("%s\n[%0x] HELLO", &rhgsp.Header, rhgsp.AuthKey)

	case ReqGet:
		return fmt.Sprintf("%s\n[%0x] 0x%0x ->", &rhgsp.Header, rhgsp.AuthKey, rhgsp.Register)

	case ReqSet:
		return fmt.Sprintf("%s\n[%0x] 0x%0x <- %d(0x%0x)", &rhgsp.Header, rhgsp.AuthKey, rhgsp.Register, rhgsp.RegisterData, rhgsp.RegisterData)
	}

	return fmt.Sprintf("Reached the unreachable : %s\n[%0x] 0x%0x <- %d(0x%0x)", &rhgsp.Header, rhgsp.AuthKey, rhgsp.Register, rhgsp.RegisterData, rhgsp.RegisterData)
}

func (rhgsp *ReqHelloGetSetPkt) Read(data []byte) (int, error) {
	if len(data) < defaultPktLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", defaultPktLen, len(data))
	}

	if _, err := rhgsp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhgsp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rhgsp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rhgsp.RegisterData)
	binary.BigEndian.PutUint64(data[headerLen+8:], 0)
	binary.BigEndian.PutUint32(data[headerLen+16:], 0)
	binary.BigEndian.PutUint64(data[headerLen+20:], 0)
	binary.BigEndian.PutUint64(data[headerLen+28:], 0)
	binary.BigEndian.PutUint64(data[headerLen+36:], 0)
	binary.BigEndian.PutUint64(data[headerLen+44:], 0)
	binary.BigEndian.PutUint64(data[headerLen+52:], 0)

	return defaultPktLen, io.EOF
}

func (rhgsp *ReqHelloGetSetPkt) WriteTo(w io.Writer) (int64, error) {
	var data [defaultPktLen]byte

	if _, err := rhgsp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhgsp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rhgsp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rhgsp.RegisterData)

	n, err := w.Write(data[:])
	return int64(n), err
}

func (rhgsp *ReqHelloGetSetPkt) Sendto(s net.PacketConn) (int, error) {
	var data [defaultPktLen]byte

	if _, err := rhgsp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhgsp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rhgsp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rhgsp.RegisterData)

	remote, err := rawsocket.Addr(s, rhgsp.Header.DstMac)
	if err != nil {
		return 0, fmt.Errorf("Addr : %w", err)
	}
	return s.WriteTo(data[:], remote)
}
