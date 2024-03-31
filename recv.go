package main

import (
	"encoding/binary"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"io"
	"os"
	"path"
)

func getFiles(destDir string, conn io.ReadWriter, bufSize int) error {
	err := checkHeaders(RCV_HEADER, conn)
	if err != nil {
		return err
	}
	recvBuf := make([]byte, bufSize)
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
				os.MkdirAll(dir, 0755)
			}

			file, err := os.Create(tmpName)
			if err != nil {
				return false, err
			}
			defer os.Remove(tmpName)

			remaining := header.Size
			percentage := int64(-1)

			hashWriter := xxhash.New()
			writer := io.MultiWriter(hashWriter, file)

			for remaining > 0 {
				var msg_size int = bufSize
				if remaining < uint64(bufSize) {
					msg_size = int(remaining)
				}

				nRead, err := conn.Read(recvBuf[:msg_size])
				if err != nil {
					file.Close()
					return false, err
				}

				_, err = writer.Write(recvBuf[:nRead])
				if err != nil {
					file.Close()
					return false, err
				}

				remaining -= uint64(nRead)
				if int64(100-(remaining*100)/header.Size) != percentage {
					percentage = int64(100 - (remaining*100)/header.Size)
					cleanLine()
					printLine(fileName, float64(percentage))
				}
			}
			file.Close()
			fmt.Println("")

			hash := hashWriter.Sum64()
			hashBuf := [8]byte{}
			_, err = io.ReadFull(conn, hashBuf[:])
			if err != nil {
				return false, err
			}
			testHash := binary.BigEndian.Uint64(hashBuf[:])

			if hash == testHash {
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
