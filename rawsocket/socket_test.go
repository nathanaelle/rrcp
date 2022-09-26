package rawsocket

import "testing"

const ifname = "eth0"

func TestGetSocket(t *testing.T) {
	sock, err := Socket(ifname)
	if err != nil {
		t.Errorf("getRecvSocket() : %v", err)
		return
	}

	if err := sock.Close(); err != nil {
		t.Errorf("Close() : %v", err)
	}

	var data [1 << 7]byte

	n, from, err := sock.ReadFrom(data[:])
	if err != nil {
		t.Errorf("ReadFrom %v", err)
		return
	}
	if from == nil || n == 0 {
		t.Errorf("ReadFrom size %d from %v", n, from)
	}
}
