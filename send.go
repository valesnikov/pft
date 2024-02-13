package main

import (
	"net"
	"fmt"
	"os"
	"encoding/binary"
)

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
				fmt.Print("\033[2K\r")
				fmt.Print(progressBar(float64(percentage), 10, 0))
				fmt.Printf("- %v", fileName)
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
