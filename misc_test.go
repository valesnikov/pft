package main

import (
	"io"
	"bytes"
	"testing"
)

type TestHeaderConn struct{msg []byte}
func (c TestHeaderConn) Read(p []byte) (n int, err error) {
	return bytes.NewReader(c.msg).Read(p)
}
func (c TestHeaderConn) Write(p []byte) (n int, err error) {
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
		{"send-receive", args{SND_HEADER, TestHeaderConn{RCV_HEADER[:]}}, false},
		{"receive-send", args{RCV_HEADER, TestHeaderConn{SND_HEADER[:]}}, false},
		{"send-send", args{SND_HEADER, TestHeaderConn{SND_HEADER[:]}}, true},
		{"receive-receive", args{RCV_HEADER, TestHeaderConn{RCV_HEADER[:]}}, true},
		
		{"wrong-receive", args{WRONG_HEADER, TestHeaderConn{RCV_HEADER[:]}}, true},
		{"receive-wrong", args{RCV_HEADER, TestHeaderConn{WRONG_HEADER[:]}}, true},
		{"wrong-send", args{WRONG_HEADER, TestHeaderConn{SND_HEADER[:]}}, true},
		{"send-wrong", args{SND_HEADER, TestHeaderConn{WRONG_HEADER[:]}}, true},

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
