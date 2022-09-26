package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	// headerPkt = 128 bits = 16 bytes
	HeaderPkt struct {
		DstMac [6]byte
		SrcMac [6]byte
		Opcode Opcode
	}
)

var _ io.Reader = &HeaderPkt{}
var _ io.Writer = &HeaderPkt{}

const headerLen = 16

func (h *HeaderPkt) Write(data []byte) (int, error) {
	if len(data) < headerLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", headerLen, len(data))
	}

	copy(h.DstMac[:], data[0:6])
	copy(h.SrcMac[:], data[6:12])
	h.Opcode = Opcode(binary.BigEndian.Uint32(data[12:16]))

	return headerLen, nil
}

func (h *HeaderPkt) Read(data []byte) (int, error) {
	if len(data) < headerLen {
		return 0, fmt.Errorf("need at least %d bytes got only %d", headerLen, len(data))
	}

	copy(data[0:6], h.DstMac[:])
	copy(data[6:12], h.SrcMac[:])
	binary.BigEndian.PutUint32(data[12:16], uint32(h.Opcode))

	return headerLen, nil
}

func (h *HeaderPkt) String() string {
	var args [13]any
	args[12] = h.Opcode
	for i, r := range h.DstMac {
		args[0+i] = r
	}
	for i, r := range h.SrcMac {
		args[6+i] = r
	}

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x <- %02x:%02x:%02x:%02x:%02x:%02x [%s]", args[:]...)
}

func (h *HeaderPkt) Data() RRCPacket {
	return h.Opcode.Data(h)
}
