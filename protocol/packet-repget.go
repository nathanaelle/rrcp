package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/nathanaelle/rrcp/rawsocket"
)

type (
	RepGetPkt struct {
		Header       HeaderPkt
		AuthKey      uint16
		Register     uint16
		RegisterData uint32
		Reserved     [8]byte
		ResZero      [28]byte
	}
)

func (rgp *RepGetPkt) String() string {
	return fmt.Sprintf("%s\n[%0x] 0x%0x => %d(0x%0x)", &rgp.Header, rgp.AuthKey, rgp.Register, rgp.RegisterData, rgp.RegisterData)
}

func (rgp *RepGetPkt) Read(data []byte) (int, error) {
	if len(data) < defaultPktLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", defaultPktLen, len(data))
	}

	if _, err := rgp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rgp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rgp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rgp.RegisterData)
	binary.BigEndian.PutUint64(data[headerLen+8:], 0)
	binary.BigEndian.PutUint32(data[headerLen+16:], 0)
	binary.BigEndian.PutUint64(data[headerLen+20:], 0)
	binary.BigEndian.PutUint64(data[headerLen+28:], 0)
	binary.BigEndian.PutUint64(data[headerLen+36:], 0)
	binary.BigEndian.PutUint64(data[headerLen+44:], 0)
	binary.BigEndian.PutUint64(data[headerLen+52:], 0)

	return defaultPktLen, io.EOF
}

func (rgp *RepGetPkt) WriteTo(w io.Writer) (int64, error) {
	var data [defaultPktLen]byte

	if _, err := rgp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rgp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rgp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rgp.RegisterData)

	n, err := w.Write(data[:])
	return int64(n), err
}

func (rgp *RepGetPkt) Sendto(s net.PacketConn) (int, error) {
	var data [defaultPktLen]byte

	if _, err := rgp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rgp.AuthKey)
	binary.BigEndian.PutUint16(data[headerLen+2:], rgp.Register)
	binary.BigEndian.PutUint32(data[headerLen+4:], rgp.RegisterData)

	remote, err := rawsocket.Addr(s, rgp.Header.DstMac)
	if err != nil {
		return 0, fmt.Errorf("Addr : %w", err)
	}
	return s.WriteTo(data[:], remote)
}
