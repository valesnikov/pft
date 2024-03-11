package main

import (
	"errors"
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

var ErrHeaders = errors.New("check headers: receive and send headers do not match")

func checkHeaders(header [HEADER_SIZE]byte, conn io.ReadWriter) error {
	var hdr = [HEADER_SIZE]byte{}

	if header == SND_HEADER {
		_, err := conn.Write(header[:]) //send header
		if err != nil {
			fmt.Println(err)
			return ErrHeaders
		}
		_, err = io.ReadFull(conn, hdr[:])
		if err != nil {
			fmt.Println(err)
			return ErrHeaders
		}
	} else if header == RCV_HEADER {
		_, err := io.ReadFull(conn, hdr[:])
		if err != nil {
			fmt.Println(err)
			return ErrHeaders
		}
		
		_, err = conn.Write(header[:]) //send header
		if err != nil {
			fmt.Println(err)
			return ErrHeaders
		}
	}

	if header == SND_HEADER && hdr == RCV_HEADER { //cmp headers
		return nil
	} else if header == RCV_HEADER && hdr == SND_HEADER { //cmp headers
		return nil
	} else {
		return ErrHeaders
	}
}
