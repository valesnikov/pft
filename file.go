package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
)

type FileHeader struct {
	Size uint64
	Name string
	Flags uint32
}

func MakeFileHeader(filepath string, flags uint32) (FileHeader, error) {
	wrap_err := func (err error) error {return fmt.Errorf("make \"%s\" file header:\n%w", filepath, err)}

	file, err := os.Open(filepath)
	if err != nil {
		return FileHeader{}, wrap_err(err)
	}
	defer file.Close()

	fStat, err := file.Stat()
	if err != nil {
		return FileHeader{}, wrap_err(err)
	}

	fileName := path.Base(file.Name())

	return FileHeader{
		Size: uint64(fStat.Size()),
		Name: fileName,
		Flags: flags,
	}, nil
}

func (fh FileHeader) Serialize() []byte {
	buf := make([]byte, 16+len(fh.Name))
	binary.BigEndian.PutUint32(buf[0:4], fh.Flags)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(fh.Name)))
	binary.BigEndian.PutUint64(buf[8:16], fh.Size)
	copy(buf[16:], []byte(fh.Name))
	return buf
}

func ReadFileHeader(reader io.Reader) (FileHeader, error) {
	wrap_err := func (err error) error {return fmt.Errorf("read file header from \"%v\":\n%w", reader, err)}

	sizesBuf := [16]byte{}
	_, err := io.ReadFull(reader, sizesBuf[:])
	if err != nil {
		return FileHeader{}, wrap_err(err)
	}
	flags := binary.BigEndian.Uint32(sizesBuf[0:4])
	nameSize := binary.BigEndian.Uint32(sizesBuf[4:8])
	fileSize := binary.BigEndian.Uint64(sizesBuf[8:16])
	if nameSize == 0 {
		return FileHeader{
			Size: fileSize,
			Name: "",
			Flags: flags,
		}, nil
	}

	nameBuf := make([]byte, nameSize)
	_, err = io.ReadFull(reader, nameBuf)
	if err != nil {
		return FileHeader{}, wrap_err(err)
	}
	return FileHeader{
		Size: fileSize,
		Name: string(nameBuf),
		Flags: flags,
	}, nil
}
