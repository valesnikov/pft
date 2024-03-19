package main

import (
	"bytes"
	crnd "crypto/rand"
	"math/rand"
	"os"
	"path"
	"testing"
)

func Test_NewFileHeader(t *testing.T) {
	dirIn, err := os.MkdirTemp(".", "test_headers")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirIn)

	const fileNum = 100
	rightHeaders := [fileNum]FileHeader{}

	const maxSizeB = 1024 * 1024 * 10 //10Mb
	inFileNames := [fileNum]string{}

	for i := 0; i < fileNum; i++ {
		file, err := os.CreateTemp(dirIn, "rndfile")
		if err != nil {
			panic(err)
		}
		b := make([]byte, rand.Int()%maxSizeB)
		_, err = crnd.Read(b)
		if err != nil {
			panic(err)
		}

		_, err = file.Write(b)
		if err != nil {
			panic(err)
		}

		inFileNames[i] = file.Name()

		fileNameBase := path.Base(file.Name())

		hash, err := getFileHash(file)
		if err != nil {
			panic(err)
		}
		file.Close()
		rightHeaders[i] = FileHeader{
			Size:     uint64(len(b)),
			Hash:     hash,
			Name:     fileNameBase,
		}
	}

	for i := 0; i < fileNum; i++ {
		header, err := MakeFileHeader(inFileNames[i])
		if err != nil {
			panic(err)
		}
		if header != rightHeaders[i] {
			t.Errorf("The files headers are different")
		}

		deserHeader, err := ReadFileHeader(bytes.NewReader(header.Serialize()))
		if err != nil {
			panic(err)
		}

		if deserHeader != header {
			t.Errorf("The files headers are different after serialize and deserialize")
		}
	}
}
