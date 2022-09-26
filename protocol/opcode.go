package protocol

import "fmt"

type (
	Opcode uint32
)

const (
	ReqHello Opcode = 0x88990100
	RepHello Opcode = 0x88990180
	ReqGet   Opcode = 0x88990101
	RepGet   Opcode = 0x88990181
	ReqSet   Opcode = 0x88990102
)

const (
	maskRpq   Opcode = 0xffffff00
	RpqLoop   Opcode = 0x88990300
	RpqEcho   Opcode = 0x88990200
	RpqLoop25 Opcode = 0x88992500
	RpqLoop23 Opcode = 0x88992300
)

func (o Opcode) IsValid() bool {
	switch o {
	case ReqHello, RepHello, ReqGet, RepGet, ReqSet:
		return true
	}

	switch o & maskRpq {
	case RpqLoop, RpqEcho, RpqLoop23, RpqLoop25:
		return true
	}

	return false
}

func (o Opcode) Header(dst, src [6]byte) RRCPacket {
	return o.Data(&HeaderPkt{
		DstMac: dst,
		SrcMac: src,
		Opcode: o,
	})
}

func (o Opcode) Data(header *HeaderPkt) RRCPacket {
	switch o {
	case ReqHello:
		return &ReqHelloGetSetPkt{Header: *header}
	case RepHello:
		return &RepHelloPkt{Header: *header}
	case ReqGet:
		return &ReqHelloGetSetPkt{Header: *header}
	case RepGet:
		return &RepGetPkt{Header: *header}
	case ReqSet:
		return &ReqHelloGetSetPkt{Header: *header}
	}

	switch o & maskRpq {
	case RpqLoop:
		return &RpqUnknownPkt{Header: *header}
	case RpqEcho:
		return &RpqUnknownPkt{Header: *header}
	case RpqLoop23:
		return &RpqUnknownPkt{Header: *header}
	case RpqLoop25:
		return &RpqUnknownPkt{Header: *header}
	}

	return &RpqUnknownPkt{Header: *header}
}

func (o Opcode) String() string {
	switch o {
	case ReqHello:
		return "Req Hello"
	case RepHello:
		return "Rep Hello"
	case ReqGet:
		return "Req Get"
	case RepGet:
		return "Rep Get"
	case ReqSet:
		return "Req Set"
	}

	switch o & maskRpq {
	case RpqLoop:
		return "Rpq Loop"
	case RpqEcho:
		return "Rpq Echo"
	case RpqLoop23:
		return "Rpq Loop23"
	case RpqLoop25:
		return "Rpq Loop25"
	}

	return fmt.Sprintf("unknown 0x%0x", uint32(o))
}
