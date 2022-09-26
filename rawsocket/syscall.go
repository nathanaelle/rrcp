package rawsocket

import (
	"unsafe"
)

type (
	Sockaddr interface {
		Sockaddr() (unsafe.Pointer, uint32, error)
	}
)
