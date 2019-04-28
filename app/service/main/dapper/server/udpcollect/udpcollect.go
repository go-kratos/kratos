package udpcollect

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"go-common/library/log"
)

const (
	_bufsize = 32 * 1024
)

// New UnixCollect
func New(addr string, workers int, writeFn func(p []byte) error) (*UDPCollect, error) {
	if workers == 0 {
		workers = 1
	}
	addrURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse addr error: %s", err)
	}
	return &UDPCollect{
		addr:    addrURL,
		writeFn: writeFn,
		workers: workers,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, _bufsize)
			},
		},
		readTimeout: 60 * time.Second,
	}, nil
}

// UDPCollect collect span data from unix socket
type UDPCollect struct {
	wg          sync.WaitGroup
	workers     int
	addr        *url.URL
	writeFn     func(p []byte) error
	readTimeout time.Duration
	pool        sync.Pool
	closed      bool
	pconn       net.PacketConn
}

// Start collector
func (u *UDPCollect) Start() error {
	var err error
	switch u.addr.Scheme {
	case "unixgram":
		u.pconn, err = listenUNIX(u.addr.Path)
	case "udp", "udp4", "udp6":
		u.pconn, err = listtenNet(u.addr.Scheme, u.addr.Host)
	default:
		return fmt.Errorf("unsupport network %s", u.addr.Scheme)
	}
	if err != nil {
		return fmt.Errorf("listen packet error: %s", err)
	}
	log.Info("dapper agent listen at: %s, workers: %d", u.addr, u.workers)
	u.wg.Add(u.workers)
	for i := 0; i < u.workers; i++ {
		go u.serve()
	}
	return nil
}

func listenUNIX(addr string) (net.PacketConn, error) {
	dirname := path.Dir(addr)
	info, err := os.Stat(dirname)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return nil, fmt.Errorf("create directory %s error: %s", dirname, err)
		}
	}
	if err == nil && !info.IsDir() {
		return nil, fmt.Errorf("%s is already exists and not a directory", dirname)
	}
	if _, err := os.Stat(addr); err == nil {
		// remove old socket file
		os.Remove(addr)
	}
	conn, err := net.ListenPacket("unixgram", addr)
	if err != nil {
		return nil, err
	}
	// make file permission to 666, so php can wirte span to this socket
	return conn, os.Chmod(addr, 0666)
}

func listtenNet(network, addr string) (net.PacketConn, error) {
	return net.ListenPacket(network, addr)
}
func (u *UDPCollect) serve() {
	defer u.wg.Done()
	for {
		if err := u.handler(u.pconn); err != nil {
			if strings.Contains(err.Error(), "closed") && u.closed {
				return
			}
			log.Error("handler PacketConn error: %s, retry after second", err)
			time.Sleep(time.Second)
		}
	}
}

func (u *UDPCollect) handler(pconn net.PacketConn) error {
	p := u.buffer()
	defer u.freeBuffer(p)
	pconn.SetReadDeadline(time.Now().Add(u.readTimeout))
	n, _, err := pconn.ReadFrom(p)
	if n > 0 {
		u.writeFn(p[:n])
	}
	if err == nil {
		return nil
	}
	if netErr, ok := err.(net.Error); ok {
		// ignore timeout and temporyary
		if netErr.Timeout() || netErr.Temporary() {
			return nil
		}
	}
	return err
}

func (u *UDPCollect) buffer() []byte {
	return u.pool.Get().([]byte)
}

func (u *UDPCollect) freeBuffer(p []byte) {
	u.pool.Put(p)
}

// Close udp collect
func (u *UDPCollect) Close() error {
	u.closed = true
	u.pconn.Close()
	// wait all workers exit
	u.wg.Wait()
	if u.addr.Scheme == "unixgram" {
		return os.Remove(u.addr.Path)
	}
	return nil
}
