package buf

import "runtime"

func (b *Buf) read(p []byte) int {
	n := 0
	if b.r > b.w || (b.r == b.w && b.len == len(b.d)) {
		n += copy(p[n:], b.d[b.r:])
		b.r = (b.r + n) % len(b.d)
		b.len -= n
		if n == len(p) {
			return n
		}
	}
	num := copy(p[n:], b.d[b.r: b.w])
	b.r += num
	n += num
	b.len -= num
	return n
}

func (b *Buf) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()

RECHECK:
	if b.len == 0 {
		if b.err != nil {
			return 0, b.err
		}
		b.m.Unlock()
		for b.len == 0 {
			runtime.Gosched()
		}
		b.m.Lock()
		goto RECHECK
	}

	return b.read(p), nil
}

func (b *Buf) ReadNonBlocking(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.read(p), nil
}
