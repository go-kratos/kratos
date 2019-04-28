package tcp

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type connMock struct {
}

func (c *connMock) Read(b []byte) (n int, err error) {
	buf := []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n")
	copy(b, buf)
	return len(buf), nil
}
func (c *connMock) Write(b []byte) (n int, err error) {
	t := []byte{_protoStr}
	t = append(t, []byte(_ok)...)
	t = append(t, []byte("\r\n")...)
	if !bytes.Equal(b, t) {
		return 0, fmt.Errorf("%s not equal %s", b, t)
	}
	return len(b), nil
}
func (c *connMock) Close() error {
	return nil
}
func (c *connMock) LocalAddr() net.Addr {
	return nil
}
func (c *connMock) RemoteAddr() net.Addr {
	return nil
}
func (c *connMock) SetDeadline(t time.Time) error {
	return nil
}
func (c *connMock) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *connMock) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestConn(t *testing.T) {
	Convey("test conn:", t, func() {
		connMock := &connMock{}
		conn := newConn(connMock, time.Second, time.Second)
		// write
		p := proto{prefix: _protoStr, message: _ok}
		err := conn.Write(p)
		So(err, ShouldBeNil)
		err = conn.Flush()
		So(err, ShouldBeNil)
		// read
		cmd, args, err := conn.Read()
		So(err, ShouldBeNil)
		So(cmd, ShouldEqual, "set")
		So(string(args[0]), ShouldEqual, "mykey")
		So(string(args[1]), ShouldEqual, "myvalue")
	})
}
