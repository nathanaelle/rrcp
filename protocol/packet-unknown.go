package protocol

import (
	"fmt"
	"io"
	"net"

	"github.com/nathanaelle/rrcp/rawsocket"
)

type (
	RpqUnknownPkt struct {
		Header  HeaderPkt
		Unknown [44]byte
	}
)

func (rup *RpqUnknownPkt) String() string {
	return fmt.Sprintf("%s\n%#v", &rup.Header, rup.Unknown)
}

func (rup *RpqUnknownPkt) Read(data []byte) (int, error) {
	if len(data) < defaultPktLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", defaultPktLen, len(data))
	}

	if _, err := rup.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}

	copy(data[headerLen:], rup.Unknown[:])

	return defaultPktLen, nil
}

func (rup *RpqUnknownPkt) WriteTo(w io.Writer) (int64, error) {
	var data [defaultPktLen]byte

	if _, err := rup.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}
	copy(data[headerLen:], rup.Unknown[:])

	n, err := w.Write(data[:])
	return int64(n), err
}

func (rup *RpqUnknownPkt) Sendto(s net.PacketConn) (int, error) {
	var data [defaultPktLen]byte

	if _, err := rup.Header.Read(data[0:headerLen]); err != nil {
		return 0, fmt.Errorf("Header : %w", err)
	}
	copy(data[headerLen:], rup.Unknown[:])

	remote, err := rawsocket.Addr(s, rup.Header.DstMac)
	if err != nil {
		return 0, fmt.Errorf("Addr : %w", err)
	}
	return s.WriteTo(data[:], remote)
}
