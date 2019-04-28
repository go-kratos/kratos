package infoc

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
	"strconv"
	"sync"
	"time"

	"go-common/library/log"
)

var (
	_infoc2Magic = []byte{172, 190} // NOTE: magic 0xAC0xBE
	_infoc2Type  = []byte{0, 0}     // NOTE: type 0

	_infocTimeout = 500 * time.Millisecond
)

// Config is infoc config.
type Config struct {
	TaskID string
	// udp or tcp
	Proto string
	Addr  string
	// reporter
	ReporterAddr string
}

// Infoc infoc struct.
type Infoc struct {
	c      *Config
	header []byte
	// udp or tcp
	conn net.Conn
	lock sync.Mutex
	// reporter
	reporter *reporter
}

// New new infoc2 logger.
func New(c *Config) (i *Infoc) {
	i = &Infoc{
		c:      c,
		header: []byte(c.TaskID),
	}
	var err error
	if i.conn, err = net.Dial(i.c.Proto, i.c.Addr); err != nil {
		log.Error("infoc net dial error(%v)", err)
	}
	if c.ReporterAddr != "" {
		i.reporter = newReporter(c.TaskID, c.ReporterAddr)
		go i.reporter.reportproc()
	}
	return
}

// Rows the affected by binlog enent.
func (i *Infoc) Rows(rows int64) {
	if i.reporter != nil {
		i.reporter.receiveIncr(rows)
	}
}

// Send send message.
func (i *Infoc) Send(ctx context.Context, key string, v interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(v); err != nil {
		log.Error("json.Marshal(%v) error(%v)", v, err)
		return
	}
	var (
		res bytes.Buffer
		buf bytes.Buffer
	)
	res.Write(_infoc2Magic)
	// type and body buf, for calc length.
	buf.Write(_infoc2Type)
	buf.Write(i.header)
	buf.WriteString(strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	// // append first arg
	if _, err = buf.WriteString(string(b)); err != nil {
		return
	}
	// put length
	var ls [4]byte
	binary.BigEndian.PutUint32(ls[:], uint32(buf.Len()))
	res.Write(ls[:])       // NOTE: write length
	res.Write(buf.Bytes()) // NOTE：write type and body
	// write
	if err = i.write(res.Bytes()); err != nil {
		log.Error("infoc write error(%v)", err)
		return
	}
	if i.reporter != nil {
		i.reporter.sendIncr(1)
	}
	return
}

// write write data into connection.
func (i *Infoc) write(bs []byte) (err error) {
	defer func() {
		if err != nil {
			if i.conn != nil {
				i.conn.Close()
			}
			i.conn = nil
		}
		i.lock.Unlock()
	}()
	i.lock.Lock()
	// connection and write
	if i.conn == nil {
		if i.conn, err = net.DialTimeout(i.c.Proto, i.c.Addr, _infocTimeout); err != nil {
			log.Error("infoc net dial error(%v)", err)
			return
		}
	}
	if i.c.Proto == "tcp" {
		i.conn.SetDeadline(time.Now().Add(_infocTimeout))
	}
	if _, err = i.conn.Write(bs); err != nil {
		log.Error("infoc net write error(%v)", err)
	}
	return
}

// Flush flush reporter count.
func (i *Infoc) Flush() {
	if i.reporter != nil {
		i.reporter.flush()
	}
}

// Close close resource.
func (i *Infoc) Close() {
	if i.conn != nil {
		i.conn.Close()
	}
}
