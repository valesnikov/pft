package main

import (
	"fmt"
	"io"
)

/*
len(HEADER) byte - header
while {filename size} != 0 {
	8 byte - filename size
	8 byte - file size
	{filename size} byte - filename
	{file size} byte - file
}
*/

const BUFSIZE int = 1024 * 1024 //1MiB

const HEADER_SIZE = 8

var SND_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x73, 0x30, 0x30, 0x31, 0x0a} //pfts001\n
var RCV_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x72, 0x30, 0x30, 0x31, 0x0a} //pftr001\n

func checkHeaders(header [HEADER_SIZE]byte, conn io.ReadWriteCloser) int {
	_, err := conn.Write(header[:]) //send header
	if err != nil {
		fmt.Println(err)
		return 1
	}

	hdr := [HEADER_SIZE]byte{} //receiver header
	_, err = conn.Read(hdr[:])
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if header == SND_HEADER && hdr == RCV_HEADER { //cmp headers
		return 0
	} else if header == RCV_HEADER && hdr == SND_HEADER { //cmp headers
		return 0
	} else {
		fmt.Println("The headings don't match")
		return 1
	}
}
