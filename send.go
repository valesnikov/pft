package main

import (
	"fmt"
	"net"
	"os"
)

func sendFiles(names []string, conn net.Conn) error {
	defer conn.Close()
	filesOpen, filesNames, err := halalizeFileName(names)
	if err != nil {
		return err
	}

	err = checkHeaders(SND_HEADER, conn)
	if err != nil {
		return err
	}
	sendBuf := make([]byte, TransmissionBufferSize)

	for i, filepath := range filesOpen {
		err := func() error {

			header, err := MakeFileHeader(filepath)
			if err != nil {
				return err
			}
			header.Name = filesNames[i]

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
				var msg_size int = TransmissionBufferSize
				if remaining < int64(TransmissionBufferSize) {
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
		Size:     0,
		Hash:     0,
		Name: 	 "",
	}
	_, err = conn.Write(nullHeader.Serialize())
	if err != nil {
		return err
	}

	return nil
}
