package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/nathanaelle/rrcp/rawsocket"
)

type (
	RepHelloPkt struct {
		Header       HeaderPkt
		AuthKey      uint16
		DownlinkPort byte
		UplinkPort   byte
		UplinkMac    [6]byte
		ChipId       uint16
		VendorId     uint32
		ResZero      [28]byte
	}
)

func (rhp *RepHelloPkt) String() string {
	var args [13]any
	args[0] = rhp.Header.String()
	args[1] = rhp.AuthKey
	args[2] = rhp.VendorId
	args[3] = rhp.ChipId
	args[4] = rhp.ChipId
	args[5] = rhp.DownlinkPort
	args[6] = rhp.UplinkPort
	for i, r := range rhp.UplinkMac {
		args[7+i] = r
	}

	return fmt.Sprintf("%s\n[%04x] vendor %0d/%0d[0x%0x] dwl %0d upl %0d %02x:%02x:%02x:%02x:%02x:%02x", args[:]...)
}

func (rhp *RepHelloPkt) Read(data []byte) (int, error) {
	if len(data) < defaultPktLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", defaultPktLen, len(data))
	}

	if _, err := rhp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhp.AuthKey)
	data[headerLen+2] = rhp.DownlinkPort
	data[headerLen+3] = rhp.UplinkPort
	copy(data[headerLen+4:headerLen+10], rhp.UplinkMac[:])
	binary.BigEndian.PutUint16(data[headerLen+10:], rhp.ChipId)
	binary.BigEndian.PutUint32(data[headerLen+12:], rhp.VendorId)
	binary.BigEndian.PutUint32(data[headerLen+16:], 0)
	binary.BigEndian.PutUint64(data[headerLen+20:], 0)
	binary.BigEndian.PutUint64(data[headerLen+28:], 0)
	binary.BigEndian.PutUint64(data[headerLen+36:], 0)
	binary.BigEndian.PutUint64(data[headerLen+44:], 0)
	binary.BigEndian.PutUint64(data[headerLen+52:], 0)

	return defaultPktLen, nil
}

func (rhp *RepHelloPkt) WriteTo(w io.Writer) (int64, error) {
	var data [defaultPktLen]byte

	if _, err := rhp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhp.AuthKey)
	data[headerLen+2] = rhp.DownlinkPort
	data[headerLen+3] = rhp.UplinkPort
	copy(data[headerLen+4:headerLen+10], rhp.UplinkMac[:])
	binary.BigEndian.PutUint16(data[headerLen+10:], rhp.ChipId)
	binary.BigEndian.PutUint32(data[headerLen+12:], rhp.VendorId)

	n, err := w.Write(data[:])
	return int64(n), err
}

func (rhp *RepHelloPkt) Sendto(s net.PacketConn) (int, error) {
	var data [defaultPktLen]byte

	if _, err := rhp.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	binary.BigEndian.PutUint16(data[headerLen:], rhp.AuthKey)
	data[headerLen+2] = rhp.DownlinkPort
	data[headerLen+3] = rhp.UplinkPort
	copy(data[headerLen+4:headerLen+10], rhp.UplinkMac[:])
	binary.BigEndian.PutUint16(data[headerLen+10:], rhp.ChipId)
	binary.BigEndian.PutUint32(data[headerLen+12:], rhp.VendorId)

	remote, err := rawsocket.Addr(s, rhp.Header.DstMac)
	if err != nil {
		return 0, fmt.Errorf("Addr : %w", err)
	}
	return s.WriteTo(data[:], remote)
}
