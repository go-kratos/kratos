package deliver

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"
)

func TestDeliver(t *testing.T) {
	buf := &bytes.Buffer{}
	lis, err := net.Listen("tcp", "127.0.0.1:12233")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(buf, conn)
	}()
	data := []byte("hello world")
	readed := make(chan bool, 1)
	d, err := New([]string{"127.0.0.1:12233"}, func() ([]byte, error) {
		readed <- true
		return data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)
	if !bytes.Equal(buf.Bytes()[0:2], _magicBuf) {
		t.Error("invalid data, wrong magic header")
	}
	if int(binary.BigEndian.Uint32(buf.Bytes()[2:6])) != len(data) {
		t.Error("wrong data length")
	}
	if !bytes.Equal(buf.Bytes()[6:], data) {
		t.Errorf("invalid content %s", buf.Bytes()[6:])
	}
	d.Close()
}
