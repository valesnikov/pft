package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash/v2"
	"github.com/klauspost/compress/zstd"
)

func sendFiles(names []string, conn io.Writer, bufSize int) error {
	zconn, err := zstd.NewWriter(conn, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		return err
	}
	defer zconn.Close()

	filesOpen, filesNames, err := prepareFileNames(names)
	if err != nil {
		return err
	}

	sendBuf := make([]byte, bufSize)

	for i, filepath := range filesOpen {
		err := func() error {

			header, err := MakeFileHeader(filepath)
			if err != nil {
				return err
			}
			header.Name = filesNames[i]

			_, err = zconn.Write(header.Serialize())
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

			hashWriter := xxhash.New()
			writer := io.MultiWriter(hashWriter, zconn)

			if remaining == 0 {
				printLine(filepath, 100)
			}

			for remaining > 0 {
				var msg_size int = bufSize
				if remaining < int64(bufSize) {
					msg_size = int(remaining)
				}

				nRead, err := file.Read(sendBuf[:msg_size])
				if err != nil {
					return err
				}
				if nRead < msg_size {
					msg_size = nRead
				}

				_, err = writer.Write(sendBuf[:msg_size])
				if err != nil {
					return err
				}

				remaining -= int64(msg_size)
				if 100-(remaining*100)/int64(header.Size) != percentage {
					percentage = 100 - (remaining*100)/int64(header.Size)
					cleanLine()
					printLine(filepath, float64(percentage))
				}
			}

			hashBuf := [8]byte{}
			binary.BigEndian.PutUint64(hashBuf[:], hashWriter.Sum64())

			_, err = zconn.Write(hashBuf[:])
			if err != nil {
				return err
			}

			fmt.Println("")
			return nil
		}()

		if err != nil {
			return err
		}
	}
	nullHeader := FileHeader{ //after last file
		Size: 0,
		Name: "",
	}
	_, err = zconn.Write(nullHeader.Serialize())
	if err != nil {
		return err
	}

	return nil
}
