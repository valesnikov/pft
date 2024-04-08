package main

import (
	"fmt"
	"net"
	"time"
)

func connectHost(port string) (net.Conn, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v:%s\n", getLocalIPs(), port)

	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectClient(addr string, port string) (net.Conn, error) {
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
