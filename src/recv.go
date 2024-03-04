package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path"
)

func getFiles(destDir string, conn net.Conn) int {
	defer conn.Close()
	if checkHeaders(RCV_HEADER, conn) != 0 {
		return 1
	}
	recvBuf := make([]byte, BUFSIZE)
	fmt.Println("Started receiving")

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

			remaining -= uint64(nWrite)
			if int64(100-(remaining*100)/fileSize) != percentage {
				percentage = int64(100 - (remaining*100)/fileSize)
				fmt.Print("\033[2K\r")
				printLine(fileName, float64(percentage))
			}
		}

		err = os.Rename(tmpName, fileName)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("\n")
		//fmt.Printf("Done: %v\n", fullName)
	}
	fmt.Println("Finished receiving")
	return 0
}
