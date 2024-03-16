package main

import (
	"fmt"
	"net"
	"os"
)

func sendFiles(files []string, conn net.Conn) error {
	defer conn.Close()
	err := checkHeaders(SND_HEADER, conn)
	if err != nil {
		return err
	}
	sendBuf := make([]byte, BUFSIZE)

	for _, filepath := range files {
		err := func() error {

			header, err := MakeFileHeader(filepath)
			if err != nil {
				return err
			}

			_, err = conn.Write(header.Serialize())
			if err != nil {
				return err
			}

			remaining := int64(header.Size)
			percentage := int64(-1)

			file, err := os.Open(filepath)
			if err != nil {
				return err
			}
			defer file.Close()

			for remaining > 0 {
				var msg_size int = BUFSIZE
				if remaining < int64(BUFSIZE) {
					msg_size = int(remaining)
				}

				nRead, err := file.Read(sendBuf[:msg_size])
				if err != nil {
					return err
				}
				if nRead < msg_size {
					msg_size = nRead
				}

				_, err = conn.Write(sendBuf[:msg_size])
				if err != nil {
					return err
				}

				remaining -= int64(msg_size)
				if 100-(remaining*100)/int64(header.Size) != percentage {
					percentage = 100 - (remaining*100)/int64(header.Size)
					fmt.Print("\033[2K\r")
					printLine(filepath, float64(percentage))
				}
			}
			fmt.Println("")
			return nil
		}()

		if err != nil {
			return err
		}
	}
	nullHeader := FileHeader{ //after last file
		NameSize: 0,
		Size:     0,
		Hash:     0,
	}
	_, err = conn.Write(nullHeader.Serialize())
	if err != nil {
		return err
	}

	return nil
}
