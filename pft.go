package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path"
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

func checkHeaders(header [HEADER_SIZE]byte, conn net.Conn) int {
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

func sendFiles(files []string, conn net.Conn) int {
	defer conn.Close()
	if checkHeaders(SND_HEADER, conn) != 0 {
		return 1
	}

	sendBuf := make([]byte, BUFSIZE)

	for _, filepath := range files {

		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fStat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileName := fStat.Name()
		nameSize := len(fileName)
		fileSize := fStat.Size()

		sizeNameBuf := make([]byte, 16+nameSize)
		binary.BigEndian.PutUint64(sizeNameBuf[0:8], uint64(nameSize))
		binary.BigEndian.PutUint64(sizeNameBuf[8:16], uint64(fileSize))
		copy(sizeNameBuf[16:], []byte(fileName))

		_, err = conn.Write(sizeNameBuf)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		//fmt.Printf("Sending: %v\n", fileName)

		remaining := fileSize
		percentage := int64(-1)

		for remaining > 0 {
			var msg_size int = BUFSIZE
			if remaining < int64(BUFSIZE) {
				msg_size = int(remaining)
			}

			n, err := file.Read(sendBuf[:msg_size])
			if err != nil || n != msg_size {
				fmt.Println(err)
				return 1
			}

			n, err = conn.Write(sendBuf[:msg_size])
			if err != nil || n != msg_size {
				fmt.Println(err)
				return 1
			}

			remaining -= int64(msg_size)
			if 100-(remaining*100)/fileSize != percentage {
				percentage = 100 - (remaining*100)/fileSize
				fmt.Print("\033[2K\r");
				fmt.Printf("%v%% - %v", percentage, fileName);
			}
		}
		fmt.Print("\n")
		//fmt.Println("\nDone:", fileName)
	}

	_, err := conn.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func getFiles(destDir string, conn net.Conn) int {
	defer conn.Close()
	if checkHeaders(RCV_HEADER, conn) != 0 {
		return 1
	}
	recvBuf := make([]byte, BUFSIZE)
	fmt.Println("recv start")

	for {
		sizesBuf := [16]byte{}
		_, err := conn.Read(sizesBuf[:])

		if err != nil {
			fmt.Println(err)
			return 1
		}

		nameSize := binary.BigEndian.Uint64(sizesBuf[0:8])
		fileSize := binary.BigEndian.Uint64(sizesBuf[8:16])

		if nameSize == 0 {
			break
		}

		nameBuf := make([]byte, nameSize)
		_, err = conn.Read(nameBuf)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		fullName := string(nameBuf)
		fileName := path.Join(destDir, path.Base(fullName))
		tmpName := fileName + ".pft_tmp"

		file, err := os.Create(tmpName)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		defer os.Remove(tmpName)
		defer file.Close()

		//fmt.Printf("Getting: %v\n", fullName)

		remaining := fileSize
		percentage := int64(-1)

		for remaining > 0 {
			var msg_size int = BUFSIZE
			if remaining < uint64(BUFSIZE) {
				msg_size = int(remaining)
			}

			nRead, err := conn.Read(recvBuf[:msg_size])
			if err != nil {
				fmt.Println(err)
				return 1
			}

			nWrite, err := file.Write(recvBuf[:nRead])
			if err != nil {
				fmt.Println(err)
				return 1
			}

			if nWrite != nRead {
				fmt.Println("File write error", nWrite, nRead)
			}

			if int64(100-(remaining*100)/fileSize) != percentage {
				percentage = int64(100 - (remaining*100)/fileSize)
				fmt.Print("\033[2K\r");
				fmt.Printf("%v%% - %v", percentage, fileName);
			}

			remaining -= uint64(nWrite)
		}

		err = os.Rename(tmpName, fileName)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("\n")
		//fmt.Printf("Done: %v\n", fullName)
	}
	fmt.Println("recv done")
	return 0
}

func host_send() {
	ln, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", os.Args[2])

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
	} else {
		sendFiles(os.Args[3:], conn)
	}
}

func host_receive() {
	ln, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", os.Args[2])

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
	} else {
		getFiles(os.Args[3], conn)
	}
}

func client_send() {
	conn, err := net.Dial("tcp", os.Args[2]+":"+os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendFiles(os.Args[4:], conn)
}

func client_receive() {
	conn, err := net.Dial("tcp", os.Args[2]+":"+os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	getFiles(os.Args[4], conn)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("pft hs <port> [files]")
		fmt.Println("pft hr <port> <destdir>")
		fmt.Println("pft cs <addr> <port> [files]")
		fmt.Println("pft cr <addr> <port> <destdir>")
		os.Exit(1)
	}
	if os.Args[1] == "hs" {
		if len(os.Args) < 4 {
			fmt.Println("usage:")
			fmt.Println("pft hs <port> [files]")
			os.Exit(1)
		}
		host_send()
	} else if os.Args[1] == "hr" {
		if len(os.Args) != 4 {
			fmt.Println("usage:")
			fmt.Println("pft hr <port> <destdir>")
			os.Exit(1)
		}
		host_receive()
	} else if os.Args[1] == "cs" {
		if len(os.Args) < 5 {
			fmt.Println("usage:")
			fmt.Println("pft cs <addr> <port> [files]")
			os.Exit(1)
		}
		client_send()
	} else if os.Args[1] == "cr" {
		if len(os.Args) != 5 {
			fmt.Println("usage:")
			fmt.Println("pft cr <addr> <port> <destdir>")
			os.Exit(1)
		}
		client_receive()
	} else {
		fmt.Println("usage:")
		fmt.Println("pft hs <port> [files]")
		fmt.Println("pft hr <port> <destdir>")
		fmt.Println("pft cs <addr> <port> [files]")
		fmt.Println("pft cr <addr> <port> <destdir>")
		os.Exit(1)
	}
}
