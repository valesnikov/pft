package main

import (
	"bytes"
	"io"
	"testing"
)

type testHeaderConn struct{ msg []byte }

func (c testHeaderConn) Read(p []byte) (n int, err error) {
	return bytes.NewReader(c.msg).Read(p)
}
func (c testHeaderConn) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var WRONG_HEADER = [HEADER_SIZE]byte{0x71, 0x65, 0x75, 0x72, 0x31, 0x29, 0x32, 0x09}

func Test_checkHeaders(t *testing.T) {
	type args struct {
		header [HEADER_SIZE]byte
		conn   io.ReadWriter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"send-receive", args{SND_HEADER, testHeaderConn{RCV_HEADER[:]}}, false},
		{"receive-send", args{RCV_HEADER, testHeaderConn{SND_HEADER[:]}}, false},
		{"send-send", args{SND_HEADER, testHeaderConn{SND_HEADER[:]}}, true},
		{"receive-receive", args{RCV_HEADER, testHeaderConn{RCV_HEADER[:]}}, true},

		{"wrong-receive", args{WRONG_HEADER, testHeaderConn{RCV_HEADER[:]}}, true},
		{"receive-wrong", args{RCV_HEADER, testHeaderConn{WRONG_HEADER[:]}}, true},
		{"wrong-send", args{WRONG_HEADER, testHeaderConn{SND_HEADER[:]}}, true},
		{"send-wrong", args{SND_HEADER, testHeaderConn{WRONG_HEADER[:]}}, true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkHeaders(tt.args.header, tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("checkHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bufSizeToNum(t *testing.T) {
	type args struct {
		size string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{"1"}, 1, false},
		{"2", args{"92345678"}, 92345678, false},
		{"3", args{"2K"}, 2 * 1024, false},
		{"4", args{"642K"}, 642 * 1024, false},
		{"5", args{"3M"}, 3 * 1024 * 1024, false},
		{"6", args{"54M"}, 54 * 1024 * 1024, false},
		{"7", args{"2G"}, 2 * 1024 * 1024 * 1024, false},
		{"8", args{"1T"}, 0, true},
		{"9", args{"shue"}, 0, true},
		{"10", args{"1D"}, 0, true},
		{"11", args{"-212"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bufSizeToNum(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("bufSizeToNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bufSizeToNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
