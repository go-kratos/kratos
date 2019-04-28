package udpcollect

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestUDPCollect(t *testing.T) {
	count := 0
	data := []byte("hello world")
	collect, err := New("unixgram:///tmp/test.sock", 2, func(p []byte) error {
		count++
		if !bytes.Equal(p, data) {
			t.Errorf("invalid p: %s", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := collect.Start(); err != nil {
		t.Fatal(err)
	}
	conn, err := net.DialTimeout("unixgram", "/tmp/test.sock", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		conn.Write(data)
	}
	time.Sleep(time.Second)
	collect.Close()
	if count != 20 {
		t.Errorf("wrong get %d != 20", count)
	}
}
