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
	8 byte - hash
	{filename size} byte - filename
	{file size} byte - file
}
*/

/*
HEADER_SIZE byte - header
8 byte - num connections
8 byte -
*/

const HEADER_SIZE = 8

var SND_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x73, 0x30, 0x30, 0x32, 0x0a} //pfts002\n
var RCV_HEADER = [HEADER_SIZE]byte{0x70, 0x66, 0x74, 0x72, 0x30, 0x30, 0x32, 0x0a} //pftr002\n

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
