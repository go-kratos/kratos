package trace

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newServer(w io.Writer, network, address string) (func() error, error) {
	lis, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			lis.Close()
			log.Fatal(err)
		}
		io.Copy(w, conn)
		conn.Close()
		done <- struct{}{}
	}()
	return func() error {
		<-done
		return lis.Close()
	}, nil
}

func TestReportTCP(t *testing.T) {
	buf := &bytes.Buffer{}
	cancel, err := newServer(buf, "tcp", "127.0.0.1:6077")
	if err != nil {
		t.Fatal(err)
	}
	report := newReport("tcp", "127.0.0.1:6077", 0, 0).(*connReport)
	data := []byte("hello, world")
	report.writePackage(data)
	if err := report.Close(); err != nil {
		t.Error(err)
	}
	cancel()
	assert.Equal(t, data, buf.Bytes(), "receive data")
}

func newUnixgramServer(w io.Writer, address string) (func() error, error) {
	conn, err := net.ListenPacket("unixgram", address)
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		p := make([]byte, 4096)
		n, _, err := conn.ReadFrom(p)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(p[:n])
		done <- struct{}{}
	}()
	return func() error {
		<-done
		return conn.Close()
	}, nil
}

func TestReportUnixgram(t *testing.T) {
	os.Remove("/tmp/trace.sock")
	buf := &bytes.Buffer{}
	cancel, err := newUnixgramServer(buf, "/tmp/trace.sock")
	if err != nil {
		t.Fatal(err)
	}
	report := newReport("unixgram", "/tmp/trace.sock", 0, 0).(*connReport)
	data := []byte("hello, world")
	report.writePackage(data)
	if err := report.Close(); err != nil {
		t.Error(err)
	}
	cancel()
	assert.Equal(t, data, buf.Bytes(), "receive data")
}
