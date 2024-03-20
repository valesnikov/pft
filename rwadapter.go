package main

import (
	"io"
)

type pftWriter struct { //io.ReadWriteCloser
	file   io.WriteCloser
	rcvPos int
	sndPos int
}

func newPftWriter(file io.WriteCloser) (fw *pftWriter) {
	fw = new(pftWriter)
	fw.file = file
	fw.rcvPos = 0
	return fw
}

func (w *pftWriter) Read(p []byte) (n int, err error) {
	if len(p) > len(RCV_HEADER[:]) {
		p = p[:len(RCV_HEADER[:])]
	}
	n = copy(p, RCV_HEADER[w.rcvPos:w.rcvPos+len(p)])
	w.rcvPos += len(p)
	return n, nil
}

func (w *pftWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *pftWriter) Close() error {
	return w.file.Close()
}

// -
// -
// -
// -
// -
// -
// -
// -
// -
// -
// -
// -

type pftReader struct { //io.ReadWriteCloser
	file      io.ReadCloser
	header    []byte
	headerPos int
}

func newPftReader(file io.ReadCloser) (fw *pftReader) {
	fw = new(pftReader)
	fw.file = file
	fw.header = SND_HEADER[:]
	fw.headerPos = 0
	return fw
}

func (w *pftReader) Read(p []byte) (n int, err error) {
	return w.file.Read(p)
}

func (w *pftReader) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *pftReader) Close() error {
	return w.file.Close()
}
