package buf

import (
	"io"
	"runtime"
	"sync"
)

const defaultSize = 16 * 1024

type Buf struct {
	m    sync.Mutex
	d    []byte
	r, w int
	len  int
	err  error
}

func New(size int) *Buf {
	return &Buf{
		d:   make([]byte, size),
		r:   0,
		w:   0,
		len: 0,
	}
}

func (b *Buf) Close() error {
	b.m.Lock()
	b.err = io.EOF
	b.m.Unlock()

	for b.len != 0 {
		runtime.Gosched()
	}
	return nil
}

func (b *Buf) CloseWithError(err error) error {
	b.m.Lock()
	b.err = err
	if err == nil {
		b.err = io.EOF
	}
	b.m.Unlock()
	
	for b.len != 0 {
		runtime.Gosched()
	}
	return nil
}

func NewReader(r io.Reader) io.Reader {
	b := New(defaultSize)
	go func() {
		_, err := io.Copy(b, r)
		if err != nil {
			b.CloseWithError(err)
		} else {
			b.Close()
		}
	}()
	return b
}

func NewWriter(w io.Writer) io.WriteCloser {
	b := New(defaultSize)
	go func() {
		_, err := io.Copy(w, b)
		if err != nil {
			b.CloseWithError(err)
		} else {
			b.Close()
		}
	}()
	return b
}