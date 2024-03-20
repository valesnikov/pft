package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func getFiles(destDir string, conn io.ReadWriteCloser) error {
	defer conn.Close()
	err := checkHeaders(RCV_HEADER, conn)
	if err != nil {
		return err
	}
	recvBuf := make([]byte, TransmissionBufferSize)
	fmt.Println("Started receiving")

	for {
		done, err := func() (bool, error) {
			header, err := ReadFileHeader(conn)
			if err != nil {
				return false, err
			}

			if header.Name == "" {
				return true, nil //all files received
			}

			fileName := path.Join(destDir, header.Name)
			tmpName := fileName + ".pft_tmp"

			dir, _ := path.Split(fileName)
			if dir != "" {
				os.MkdirAll(dir, 0777)
			}

			file, err := os.Create(tmpName)
			if err != nil {
				return false, err
			}
			defer os.Remove(tmpName)

			remaining := header.Size
			percentage := int64(-1)

			for remaining > 0 {
				var msg_size int = TransmissionBufferSize
				if remaining < uint64(TransmissionBufferSize) {
					msg_size = int(remaining)
				}

				nRead, err := conn.Read(recvBuf[:msg_size])
				if err != nil {
					file.Close()
					return false, err
				}

				_, err = file.Write(recvBuf[:nRead])
				if err != nil {
					file.Close()
					return false, err
				}

				remaining -= uint64(nRead)
				if int64(100-(remaining*100)/header.Size) != percentage {
					percentage = int64(100 - (remaining*100)/header.Size)
					fmt.Print("\033[2K\r")
					printLine(fileName, float64(percentage))
				}
			}
			fmt.Println("")
			
			hash, err := getFileHash(file)
			if err != nil {
				
				return false, err
			}

			file.Close()

			if hash == header.Hash {
				err = os.Rename(tmpName, fileName)
			} else {
				fmt.Printf("failed to receive: %s", fileName)
				return false, nil
			}

			if err != nil {
				fmt.Println(err)
			}
			

			return false, nil
		}()

		if err != nil {
			return err
		}
		if done {
			break
		}
	}
	fmt.Println("Finished receiving")
	return nil
}
