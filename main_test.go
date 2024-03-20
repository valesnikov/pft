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

	//sendFiles(inFileNames[:], inConn)
	errChn := make(chan error)

	go func() {
		err := sendFiles(inFileNames[:], inConn)
		errChn <- err
	}()

	rErr := getFiles(dirOut, outConn)
	if rErr != nil {
		t.Error(rErr)
		return
	}
	sErr := <-errChn
	if sErr != nil {
		t.Error(sErr)
		return
	}

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

func Test_Archive(t *testing.T) {
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

	fileW, err := os.CreateTemp(".", "archive")
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileW.Name())

	fmt.Println("start testing")

	sErr := sendFiles(inFileNames[:], newPftWriter(fileW))
	if sErr != nil {
		t.Error(sErr)
		return
	}

	fileR, err := os.Open(fileW.Name())
	if err != nil {
		panic(err)
	}

	rErr := getFiles(dirOut, newPftReader(fileR))
	if rErr != nil {
		t.Error(rErr)
		return
	}

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
