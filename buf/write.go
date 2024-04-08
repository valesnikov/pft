package buf

import "runtime"

func (b *Buf) write(p []byte) int {
	n := 0
	if b.w > b.r || (b.w == b.r && b.len == 0) {
		n += copy(b.d[b.w:], p[n:])
		b.w = (b.w + n) % len(b.d)
		b.len += n
		if n == len(p) {
			return n
		}
	}
	num := copy(b.d[b.w: b.r], p[n:])
	b.w = (b.w + num) % len(b.d)
	n += num
	b.len += num
	return n
}

func (b *Buf) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	n = 0
	for n != len(p) {
	RECHECK:
		if b.err != nil {
			return 0, b.err
		}
		if b.len == len(b.d) {
			b.m.Unlock()
			for b.len == len(b.d) && b.err == nil{
				runtime.Gosched()
			}
			b.m.Lock()
			goto RECHECK
		}
		n += b.write(p[n:])
	}
	return n, b.err
}

// might write less than len(p)
func (b *Buf) WriteNonBlocking(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	return b.write(p), nil
}
