package protocol

import (
	"fmt"
	"io"
	"net"
)

type (
	RRCPacket interface {
		io.WriterTo
		io.Reader
		fmt.Stringer

		Sendto(net.PacketConn) (int, error)
	}

	// rpqLoop25Pkt = 480 bits = 60 bytes
	rpqLoop25Pkt struct {
		Header    HeaderPkt
		Res00idx1 [2]byte
		SelfMac1  [6]byte
		Res00idx2 [2]byte
		Res80idx1 byte
		Res00idx3 [2]byte
		SelfMac2  [6]byte
		Res80idx2 byte
		Res00idx4 [24]byte
	}
)

const defaultPktLen int = 60
const defaultAuthKey uint16 = 0x2379
const defaultCRC uint32 = 0xffff0000

var macbroadcast = [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
