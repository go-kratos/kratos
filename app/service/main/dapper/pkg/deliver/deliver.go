package deliver

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"go-common/library/log"
)

var (
	_magicBuf = []byte{0xAC, 0xBE}
	_bufpool  sync.Pool
)

func init() {
	rand.Seed(time.Now().UnixNano())
	_bufpool = sync.Pool{New: func() interface{} {
		return make([]byte, 0, 4096)
	}}
}

func freeBuf(buf []byte) {
	buf = buf[:0]
	_bufpool.Put(buf)
}

func getBuf() []byte {
	return _bufpool.Get().([]byte)
}

// Deliver deliver span to dapper-service through tcp
type Deliver struct {
	servers []string
	readFn  func() ([]byte, error)
	conn    *net.TCPConn
	dataCh  chan []byte
	closeCh chan struct{}
	closed  bool
}

// New Deliver
func New(servers []string, readFn func() ([]byte, error)) (*Deliver, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server provide")
	}
	d := &Deliver{
		servers: servers,
		readFn:  readFn,
		closeCh: make(chan struct{}, 1),
		dataCh:  make(chan []byte),
	}
	return d, d.start()
}

func (d *Deliver) start() error {
	if err := d.dial(); err != nil {
		return err
	}
	go d.fetch()
	go d.loop()
	return nil
}

func (d *Deliver) fetch() {
	for {
		if d.closed {
			return
		}
		data, err := d.readFn()
		if err != nil {
			log.Error("deliver read data error: %s", err)
			continue
		}
		d.dataCh <- data
	}
}

func (d *Deliver) loop() {
	for {
		select {
		case <-d.closeCh:
			return
		case data := <-d.dataCh:
			data = warpData(data)
		send:
			_, err := d.conn.Write(data)
			if err == nil {
				freeBuf(data)
				continue
			}
			d.reDial()
			goto send
		}
	}
}

// Close deliver
func (d *Deliver) Close() error {
	if d.closed {
		return fmt.Errorf("already closed")
	}
	d.closed = true
	d.closeCh <- struct{}{}
	timer := time.NewTimer(50 * time.Millisecond)
	select {
	case data := <-d.dataCh:
		// write last data to conn
		_, err := d.conn.Write(data)
		return fmt.Errorf("write last data error: %s", err)
	case <-timer.C:
		return nil
	}
	return nil
}

func (d *Deliver) reDial() {
	if d.conn != nil {
		d.conn.Close()
	}
	for {
		if err := d.dial(); err != nil {
			log.Error("redial error: %s, retry after second", err)
			time.Sleep(time.Second)
		}
		break
	}
}

func (d *Deliver) dial() error {
	server := chioceServer(d.servers)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return fmt.Errorf("dial tcp://%s error: %s", server, err)
	}
	d.conn = conn.(*net.TCPConn)
	d.conn.SetKeepAlive(true)
	return nil
}

func chioceServer(servers []string) string {
	return servers[rand.Intn(len(servers))]
}

func warpData(data []byte) []byte {
	buf := getBuf()
	buf = append(buf, _magicBuf...)
	buf = append(buf, []byte{0, 0, 0, 0, 0, 0}...)
	binary.BigEndian.PutUint32(buf[2:6], uint32(len(data)+2))
	buf = append(buf, data...)
	return buf
}
