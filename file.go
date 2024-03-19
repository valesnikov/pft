package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"
	"github.com/cespare/xxhash/v2"
)

var ErrWrongHash = errors.New("check hash: hashs dont match")

type FileHeader struct {
	Size uint64
	Hash uint64
	Name string
}

func MakeFileHeader(filepath string) (FileHeader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return FileHeader{}, err
	}
	defer file.Close()

	hash, err := getFileHash(file)
	if err != nil {
		return FileHeader{}, err
	}

	fStat, err := file.Stat()
	if err != nil {
		return FileHeader{}, err
	}

	fileName := path.Base(file.Name())

	return FileHeader{
		Size: uint64(fStat.Size()),
		Hash: hash,
		Name: fileName,
	}, nil
}

func (fh *FileHeader) Serialize() []byte {
	buf := make([]byte, 24+len(fh.Name))
	binary.BigEndian.PutUint64(buf[0:8], uint64(len(fh.Name)))
	binary.BigEndian.PutUint64(buf[8:16], fh.Size)
	binary.BigEndian.PutUint64(buf[16:24], fh.Hash)
	copy(buf[24:], []byte(fh.Name))
	return buf
}

func ReadFileHeader(reader io.Reader) (FileHeader, error) {
	sizesBuf := [24]byte{}
	_, err := io.ReadFull(reader, sizesBuf[:])
	if err != nil {
		return FileHeader{}, err
	}
	nameSize := binary.BigEndian.Uint64(sizesBuf[0:8])
	fileSize := binary.BigEndian.Uint64(sizesBuf[8:16])
	fileHash := binary.BigEndian.Uint64(sizesBuf[16:24])
	if nameSize == 0 || fileSize == 0 {
		return FileHeader{
			Size: 0,
			Hash: 0,
			Name: "",
		}, nil
	}

	nameBuf := make([]byte, nameSize)
	_, err = io.ReadFull(reader, nameBuf)
	if err != nil {
		return FileHeader{}, err
	}
	return FileHeader{
		Size:     fileSize,
		Hash:     fileHash,
		Name:     string(nameBuf),
	}, nil
}

func getFileHash(file io.ReadSeeker) (uint64, error) {
	hash := xxhash.New()
	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	file.Seek(0, io.SeekStart)
	defer file.Seek(pos, io.SeekStart)

	_, err = io.Copy(hash, file)
	if err != nil {
		return 0, err
	}
	return hash.Sum64(), nil
}
