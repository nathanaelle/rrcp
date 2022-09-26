# RRCP

Go implementation of Realtek Remote Control Protocol (0x8899)

## Context

Several network switches use a realtek chipset.
These chipset are manageable through RRCP.

## Example 

### Direct use of RRCP protocol

in real situation, a lot of errors may occur, so you must handle all the errors.
this code just expose the workflow.

```
	// create a raw socket
	sock, _ := rawsocket.Socket("eth0")

	// defer the socket close
	defer sock.Close()

	// get the mac adress of the hardware
	hwaddr, _ := rawsocket.MacAddr(sock)

	// prepare a broadcast of a Hello discovery
	hello := ReqHello.Header(macbroadcast, hwaddr)

	// send the packet on the socket
	hello.Sendto(sock)

	// read a packet and copy it in a buffer
	rbuff := make([]byte, 1<<7)
	sock.ReadFrom(rbuff)

	// write the buffer in a header packet
	header := &HeaderPkt{}
	header.Write(rbuff)

	// generate a data structure for this header
	data := header.Data()

	// unserialize the packet in the ad-hoc data structure
	binary.Read(bytes.NewReader(rbuff), binary.BigEndian, data)

	// now data contains the decoded packet
```

## Status

alpha

| Implemented | tested | features                                                  |
| :---------: | :----: | :-------------------------------------------------------- |
|      ✓      |   ✓    | rawsocket                                                 |
|      ✓      |   ✓    | — linux litte endian (amd64)                              |
|      ✓      |   ·    | — linux big endian                                        |
|      ✓      |   ✗    | — macos/darwin little endian (amd64, M1)                  |
|      ·      |   ·    | rawsocket compatibility with the internal go network loop |
|      ·      |   ·    | — CLOEXEC                                                 |
|      ·      |   ·    | — NONBLOCK                                                |
|      ·      |   ·    | — Deadline() / WriteDeadline() / ReadDeadline()           |
|      ·      |   ·    | — syscall.Conn                                            |
|      ✓      |   ✓    | RRCP protocol                                             |
|      ✓      |   ✓    | — HELLO                                                   |
|      ✓      |   ✓    | — GET                                                     |
|      ✓      |   ✓    | — SET                                                     |
|      ✓      |   ✓    | — LOOP                                                    |
|      ✓      |   ✓    | — ECHO                                                    |
|      ·      |   ·    | chipsets specifications                                   |
|      ·      |   ·    | documentation                                             |
|      ·      |   ·    | tests                                                     |
|      ·      |   ·    | high level API                                            |

## Compatible chipsets

| Supported | Tested | Chipset     |
| :-------: | :----: | :---------- |
|     ?     |   ✗    | RTL 8305 SC |
|     ?     |   ✗    | RTL 8316 B  |
|     ?     |   ✗    | RTL 8316 BP |
|     ?     |   ✗    | RTL 8318 P  |
|     ?     |   ✗    | RTL 8324    |
|     ?     |   ✗    | RTL 8324 P  |
|     ?     |   ✗    | RTL 8324 BP |
|     ?     |   ✗    | RTL 8326    |
|     ?     |   ✗    | RTL 8326 S  |

## Compatible manufacturers

| Supported | Manufacturer | Model | Version | Year  | Chipset |
| :-------: | -----------: | :---- | :-----: | :---: | :------ |
|     ?     |            ? | ?     |    ?    |   ?   | ?       |

