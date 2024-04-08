package buf

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
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

func Test_Reader(t *testing.T) {
	const size = 4096
	d := make([]byte, size)
	rand.Read(d)

	r := NewReader(bytes.NewReader(d))

	td, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("error = %v\n", err)
	}
	if !bytes.Equal(d, td) {
		t.Errorf("data dont match size=%v\n", size)
	}
}

func Test_Writer(t *testing.T) {
	const size = 4096
	d := make([]byte, size)
	rand.Read(d)

	bb := bytes.NewBuffer([]byte{})
	w := NewWriter(bb)

	_, err := w.Write(d)
	if err != nil {
		t.Errorf("error = %v\n", err)
	}
	w.Close()
	td, err := io.ReadAll(bb)
	if err != nil {
		t.Errorf("error = %v\n", err)
	}
	if !bytes.Equal(d, td) {
		t.Errorf("data dont match size=%v\n", size)
		fmt.Println(d, td)
	}
}

func Test_SmallBlock(t *testing.T) {
	const bsize = 137
	const dsize = 17
	b := New(bsize)
	d := make([]byte, dsize)
	td := make([]byte, dsize)
	for i := 0; i < 1000; i++ {
		rand.Read(d)
		n, err := b.Write(d)
		if err != nil {
			t.Errorf("error = %v\n", err)
		}
		if n != len(d) {
			t.Errorf("buf.Write(d), len(d)=%v, n = %v\n", len(d), n)
		}
		n, err = b.Read(td)
		if err != nil {
			t.Errorf("error = %v\n", err)
		}
		if n != len(td) {
			t.Errorf("buf.Read(td), len(td)=%v, n = %v\n", len(td), n)
		}
		if !bytes.Equal(d, td) {
			t.Errorf("data dont match\n",)
			fmt.Println(d, td)
		}
	}

}
