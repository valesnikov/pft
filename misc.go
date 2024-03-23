package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

/*
len(HEADER) byte - header
while {filename size} != 0 {
	8 byte - filename size
	8 byte - file size
	{filename size} byte - filename
	{file size} byte - file
	8 byte - hash
}
*/

const HEADER_SIZE = 8

var SND_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x73, 0x30, 0x30, 0x33, 0x0a} //pfts003\n
var RCV_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x72, 0x30, 0x30, 0x33, 0x0a} //pftr003\n

var headerTemplate = regexp.MustCompile(`^pft[rs]\d{3}\n`)

var ErrWrongHeader = errors.New("check headers: the header of the second party is not correct")
var ErrEqHeader = errors.New("check headers: the second party is also the sender/receiver")
var ErrOldHeader = errors.New("check headers: the second party has an old incompatible version")
var ErrNewHeader = errors.New("check headers: the second party has a new incompatible version")
var ErrNoHeader = errors.New("check headers: failed to receive or send the header")

func checkHeaders(header [HEADER_SIZE]byte, conn io.ReadWriter) error {
	var hdr = [HEADER_SIZE]byte{}

	if header == SND_HEADER {
		_, err := conn.Write(header[:]) //send header
		if err != nil {
			fmt.Println(err)
			return ErrNoHeader
		}
		_, err = io.ReadFull(conn, hdr[:])
		if err != nil {
			fmt.Println(err)
			return ErrNoHeader
		}
	} else if header == RCV_HEADER {
		_, err := io.ReadFull(conn, hdr[:])
		if err != nil {
			fmt.Println(err)
			return ErrNoHeader
		}

		_, err = conn.Write(header[:]) //send header
		if err != nil {
			fmt.Println(err)
			return ErrNoHeader
		}
	}

	if len(headerTemplate.Find(hdr[:])) == 0 {
		return ErrWrongHeader
	} //correct format

	thisMode := header[3]
	mode := hdr[3]
	if mode == thisMode {
		return ErrEqHeader
	}
	var version, thisVersion int
	fmt.Sscan(string(hdr[4:7]), &version)
	fmt.Sscan(string(header[4:7]), &thisVersion)

	if version > thisVersion {
		return ErrNewHeader
	} else if version < thisVersion {
		return ErrOldHeader
	}

	return nil
}

/*
/a/b/c.d to  c.d
/e/g/f 	 to  f/a.c, f/b.d, f/g.n
*/
func halalizeFileName(names []string) (forOpen, forSend []string, err error) {
	forOpen = make([]string, 0, len(names))
	forSend = make([]string, 0, len(names))

	var addEntry func(string, string) error
	addEntry = func(fullPath, name string) error {
		fi, err := os.Stat(fullPath)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			forOpen = append(forOpen, fullPath)
			forSend = append(forSend, name)
		} else {
			entries, err := os.ReadDir(fullPath)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				err := addEntry(path.Join(fullPath, entry.Name()), path.Join(name, entry.Name()))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, fullPath := range names {
		fi, err := os.Stat(fullPath)
		if err != nil {
			return nil, nil, err
		}

		err = addEntry(fullPath, fi.Name())
		if err != nil {
			return nil, nil, err
		}
	}

	return forOpen, forSend, nil
}

var bufSizeTemplate = regexp.MustCompile(`^\d+[KMG]?$`)
var ErrBSWrongFormat = errors.New("buf size: wrong format, use {num}[K/M/G]")
var ErrBSLarge = errors.New("buf size: size too large ")

func bufSizeToNum(size string) (int, error) {
	if len(bufSizeTemplate.FindString(size)) == 0 {
		return 0, ErrBSWrongFormat
	} //correct format
	res := 0
	mul := 'B'
	fmt.Sscanf(size, "%d%c", &res, &mul)
	//fmt.Printf("%d-%c\n", res, mul)
	switch mul {
	case 'B':
		res *= 1
	case 'K':
		res *= 1024
	case 'M':
		res *= 1024*1024
	case 'G':
		res *= 1024*1024*1024
	}

	return res, nil
}
