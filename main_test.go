package main

import (
	"bytes"
	crnd "crypto/rand"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"
)

func Test_SendAndReceive(t *testing.T) {
	dirIn, err := os.MkdirTemp(".", "test_in")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirIn)

	dirOut, err := os.MkdirTemp(".", "test_out")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirOut)

	const fileNum = 100
	const maxSizeB = 1024 * 1024 * 10 //10Mb

	inFileNames := [fileNum]string{}
	outFileNames := [fileNum]string{}

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
		outFileNames[i] = file.Name()
		file.Close()
	}

	outConn, inConn := net.Pipe()
	fmt.Println("start testing")

	go sendFiles(inFileNames[:], inConn)
	getFiles(dirOut, outConn)

	for i := 0; i < fileNum; i++ {
		f1, err := os.ReadFile(inFileNames[i])
		if err != nil {
			panic(err)
		}
		f2, err := os.ReadFile(outFileNames[i])
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(f1, f2) {
			t.Errorf("The files are different after the transfer")
		}
	}

}