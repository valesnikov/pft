package main

import (
	"fmt"
	"net"
	"os"
	"path"
)

func getFiles(destDir string, conn net.Conn) error {
	defer conn.Close()
	err := checkHeaders(RCV_HEADER, conn)
	if err != nil {
		return err
	}
	recvBuf := make([]byte, BUFSIZE)
	fmt.Println("Started receiving")

	for {
		done, err := func() (bool, error) {
			header, err := ReadFileHeader(conn)
			if err != nil {
				return false, err
			}
			if header.Size == 0 || header.NameSize == 0 {
				return true, nil //all files received
			}

			fullName := header.Name
			fileName := path.Join(destDir, path.Base(fullName))
			tmpName := fileName + ".pft_tmp"

			file, err := os.Create(tmpName)
			if err != nil {
				return false, err
			}
			defer os.Remove(tmpName)

			remaining := header.Size
			percentage := int64(-1)

			for remaining > 0 {
				var msg_size int = BUFSIZE
				if remaining < uint64(BUFSIZE) {
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

			file.Close()

			hash, err := getFileHash(tmpName)
			if err != nil {
				return false, err
			}

			if hash == header.Hash {
				err = os.Rename(tmpName, fileName)
			} else {
				return false, ErrWrongHash
			}

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("")

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
