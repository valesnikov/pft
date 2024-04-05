package buf

import (
	"io"
	"sync"
)

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
	b.err = io.EOF
	return nil
}

func (b *Buf) CloseWithError(err error) error {
	b.err = err
	if err == nil {
		b.err = io.EOF
	}
	return nil
}
