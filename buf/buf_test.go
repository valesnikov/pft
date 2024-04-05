package buf

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

func Test_FullWR(t *testing.T) {
	const max_size = 4096 //2^n
	data := [max_size]byte{}
	test_data := [max_size]byte{}
	for i := 1; i <= max_size; i *= 2 {
		b := New(i)

		d := data[:i]
		rand.Read(d)

		n, err := b.Write(d)
		if err != nil {
			t.Errorf("buf.Write(d), len(d)=%v, error = %v\n", len(d), err)
		}
		if n != len(d) {
			t.Errorf("buf.Write(d), len(d)=%v, n = %v\n", len(d), n)
		}

		td := test_data[:i]
		n, err = b.Read(td)
		if err != nil {
			t.Errorf("buf.Read(td), len(td)=%v, error = %v\n", len(td), err)
		}
		if n != len(d) {
			t.Errorf("buf.Read(td), len(td)=%v, n = %v\n", len(td), n)
		}
		if !bytes.Equal(d, td) {
			t.Errorf("data dont match i=%v\n", i)
		}
	}
}

func Test_Overflow(t *testing.T) {
	const buf_size = 13 //2^n
	const data_size = 409700

	d := make([]byte, data_size)
	rand.Read(d)
	b := New(buf_size)

	go func() {
		n, err := b.Write(d)
		if err != nil {
			t.Errorf("buf.Write(d), len(d)=%v, error = %v\n", len(d), err)
		}
		if n != len(d) {
			t.Errorf("buf.Write(d), len(d)=%v, n = %v\n", len(d), n)
		}
	}()

	td := make([]byte, data_size)
	n, err := io.ReadFull(b, td)
	if err != nil {
		t.Errorf("buf.Read(td), len(td)=%v, error = %v\n", len(td), err)
	}
	if n != len(d) {
		t.Errorf("buf.Read(td), len(td)=%v, n = %v\n", len(td), n)
	}
	if !bytes.Equal(d, td) {
		t.Errorf("data dont match\n")
	}
}

func Test_Close(t *testing.T) {
	const size = 16
	b := New(size)
	d := make([]byte, size)
	rand.Read(d)
	n, err := b.Write(d)
	if err != nil {
		t.Errorf("buf.Write(d), len(d)=%v, error = %v\n", len(d), err)
	}
	if n != len(d) {
		t.Errorf("buf.Write(d), len(d)=%v, n = %v\n", len(d), n)
	}

	b.Close()
	td := make([]byte, size)
	n, err = b.Read(td)
	if err != nil {
		t.Errorf("buf.Read(td), len(td)=%v, error = %v\n", len(td), err)
	}
	if n != len(d) {
		t.Errorf("buf.Read(td), len(td)=%v, n = %v\n", len(td), n)
	}
	if !bytes.Equal(d, td) {
		t.Errorf("data dont match size=%v\n", size)
	}

	_, err = b.Read(td)
	if err != io.EOF {
		t.Errorf("Unexpected err = \"%v\", want err = io.EOF\n", err)
	}
}