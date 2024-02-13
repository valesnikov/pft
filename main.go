package main

import (
	"fmt"
	"net"
	"os"
)

func host_send() int {
	ln, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		fmt.Print(err)
		return 1
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", os.Args[2])

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return 1
	} else {
		return sendFiles(os.Args[3:], conn)
	}
}

func host_receive() int {
	ln, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		fmt.Print(err)
		return 1
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", os.Args[2])

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return 1
	} else {
		return getFiles(os.Args[3], conn)
	}
}

func client_send() int {
	conn, err := net.Dial("tcp", os.Args[2]+":"+os.Args[3])
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return sendFiles(os.Args[4:], conn)
}

func client_receive() int {
	conn, err := net.Dial("tcp", os.Args[2]+":"+os.Args[3])
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return getFiles(os.Args[4], conn)
}

func invalid_usage() {
	fmt.Print(
	"usage:\n"+
	"pft hs <port> [files]\n" +
	"pft hr <port> <destdir>\n" +
	"pft cs <addr> <port> [files]\n" +
	"pft cr <addr> <port> <destdir>\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		invalid_usage()
	}
	if os.Args[1] == "hs" {
		if len(os.Args) < 4 {
			invalid_usage()
		}
		os.Exit(host_send())
	} else if os.Args[1] == "hr" {
		if len(os.Args) != 4 {
			invalid_usage()
		}
		os.Exit(host_receive())
	} else if os.Args[1] == "cs" {
		if len(os.Args) < 5 {
			invalid_usage()
		}
		os.Exit(client_send())
	} else if os.Args[1] == "cr" {
		if len(os.Args) != 5 {
			invalid_usage()
		}
		os.Exit(client_receive())
	} else {
		invalid_usage()
	}
}
