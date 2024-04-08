package main

import (
	"encoding/binary"
	"io"
	"os"
	"path"
)

type FileHeader struct {
	Size uint64
	Name string
}

func MakeFileHeader(filepath string) (FileHeader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return FileHeader{}, err
	}
	defer file.Close()

	fStat, err := file.Stat()
	if err != nil {
		return FileHeader{}, err
	}

	fileName := path.Base(file.Name())

	return FileHeader{
		Size: uint64(fStat.Size()),
		Name: fileName,
	}, nil
}

func (fh FileHeader) Serialize() []byte {
	buf := make([]byte, 16+len(fh.Name))
	binary.BigEndian.PutUint64(buf[0:8], uint64(len(fh.Name)))
	binary.BigEndian.PutUint64(buf[8:16], fh.Size)
	copy(buf[16:], []byte(fh.Name))
	return buf
}

func ReadFileHeader(reader io.Reader) (FileHeader, error) {
	sizesBuf := [16]byte{}
	_, err := io.ReadFull(reader, sizesBuf[:])
	if err != nil {
		return FileHeader{}, err
	}
	nameSize := binary.BigEndian.Uint64(sizesBuf[0:8])
	fileSize := binary.BigEndian.Uint64(sizesBuf[8:16])
	if nameSize == 0 {
		return FileHeader{
			Size: fileSize,
			Name: "",
		}, nil
	}

	nameBuf := make([]byte, nameSize)
	_, err = io.ReadFull(reader, nameBuf)
	if err != nil {
		return FileHeader{}, err
	}
	return FileHeader{
		Size: fileSize,
		Name: string(nameBuf),
	}, nil
}
