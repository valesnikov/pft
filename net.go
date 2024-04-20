package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func connectHost(port string) (net.Conn, error) {
	wrap_err := func(err error) error { return fmt.Errorf("connect host with on \"%s\" port:\n%w", port, err) }

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, wrap_err(err)
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v:%s\n", getLocalIPs(), port)

	conn, err := ln.Accept()
	if err != nil {
		return nil, wrap_err(err)
	}
	return conn, nil
}

func connectClient(addr string, port string) (net.Conn, error) {
	_ = func(err error) error {
		return fmt.Errorf("connect client to \"%s\" on \"%s\" port:\n%w", addr, port, err)
	}

	fmt.Printf("Awaiting connection to %s:%s", addr, port)
	fmt.Println("")
RETRY:
	conn, err := net.Dial("tcp", addr+":"+port)
	if err != nil {
		cleanLine()
		fmt.Print(err)
		time.Sleep(250 * time.Millisecond)
		goto RETRY
	}
	return conn, nil
}

func getLocalIPs() []string {
	res := make([]string, 0, 1)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return res
	}
	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				res = append(res, ipnet.IP.String())
			}
		}
	}
	return res
}

const (
	ZstdComressionFlag uint32 = 1 << iota // Z
)

func exchangeFlags(flags uint32, conn io.ReadWriter) (uint32, error) {
	wrap_err := func(err error) error { return fmt.Errorf("exchange flags: \"%b\":\n%w", flags, err) }

	sfb := make([]byte, 4)
	rfb := make([]byte, 4)
	binary.BigEndian.PutUint32(sfb[0:4], flags)

	_, err := conn.Write(sfb)
	if err != nil {
		return 0, wrap_err(err)
	}
	_, err = io.ReadFull(conn, rfb)
	if err != nil {
		return 0, wrap_err(err)
	}

	return binary.BigEndian.Uint32(rfb[0:4]), nil
}
