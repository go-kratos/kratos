package infoc

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/netutil"
	xtime "go-common/library/time"
)

const (
	_infocSpliter  = "\001"
	_infocReplacer = "|"
	_infocLenStart = 2
	_infocLenEnd   = 6
	_protocolLen   = 6
	_infocTimeout  = 50 * time.Millisecond
)

var (
	_infocMagic     = []byte{172, 190}   // NOTE: magic 0xAC0xBE
	_infocHeaderLen = []byte{0, 0, 0, 0} // NOTE: body len placeholder
	_infocType      = []byte{0, 0}       // NOTE: type 0
	_maxRetry       = 10
	// ErrFull error chan buffer full.
	ErrFull = errors.New("infoc: chan buffer full")
)

var (
	// ClientWeb ...
	ClientWeb     = "web"
	ClientIphone  = "iphone"
	ClientIpad    = "ipad"
	ClientAndroid = "android"

	ItemTypeAv       = "av"
	ItemTypeBangumi  = "bangumi"
	ItemTypeLive     = "live"
	ItemTypeTopic    = "topic"
	ItemTypeRank     = "rank"
	ItemTypeActivity = "activity"
	ItemTypeTag      = "tag"
	ItemTypeAD       = "ad"
	ItemTypeLV       = "lv"

	ActionClick     = "click"
	ActionPlay      = "play"
	ActionFav       = "fav"
	ActionCoin      = "coin"
	ActionDM        = "dm"
	ActionToView    = "toview"
	ActionShare     = "share"
	ActionSpace     = "space"
	Actionfollow    = "follow"
	ActionHeartbeat = "heartbeat"
	ActionAnswer    = "answer"
)

// Config is infoc config.
type Config struct {
	TaskID string
	// udp or tcp
	Proto        string
	Addr         string
	ChanSize     int
	DialTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// Infoc infoc struct.
type Infoc struct {
	c            *Config
	header       []byte
	msgs         chan *bytes.Buffer
	dialTimeout  time.Duration
	writeTimeout time.Duration
	pool         sync.Pool
	waiter       sync.WaitGroup
}

// New new infoc logger.
func New(c *Config) (i *Infoc) {
	i = &Infoc{
		c:            c,
		header:       []byte(c.TaskID),
		msgs:         make(chan *bytes.Buffer, c.ChanSize),
		dialTimeout:  time.Duration(c.DialTimeout),
		writeTimeout: time.Duration(c.WriteTimeout),
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
	if i.dialTimeout == 0 {
		i.dialTimeout = _infocTimeout
	}
	if i.writeTimeout == 0 {
		i.writeTimeout = _infocTimeout
	}
	i.waiter.Add(1)
	go i.writeproc()
	return
}

func (i *Infoc) buf() *bytes.Buffer {
	return i.pool.Get().(*bytes.Buffer)
}

func (i *Infoc) freeBuf(buf *bytes.Buffer) {
	buf.Reset()
	i.pool.Put(buf)
}

// Info record log to file.
func (i *Infoc) Info(args ...interface{}) (err error) {
	err, res := i.info(args...)
	if err != nil {
		return
	}
	select {
	case i.msgs <- res:
	default:
		i.freeBuf(res)
		err = ErrFull
	}
	return
}

// Infov support filter mirror request
func (i *Infoc) Infov(ctx context.Context, args ...interface{}) (err error) {
	if metadata.Bool(ctx, metadata.Mirror) {
		return
	}

	return i.Info(args...)
}

func getValue(i interface{}) (s string) {
	switch v := i.(type) {
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case string:
		s = v
	case bool:
		s = strconv.FormatBool(v)
	default:
		s = fmt.Sprint(i)
	}
	return
}

// Close close the connection.
func (i *Infoc) Close() error {
	i.msgs <- nil
	i.waiter.Wait()
	return nil
}

// writeproc write data into connection.
func (i *Infoc) writeproc() {
	var (
		msg  *bytes.Buffer
		conn net.Conn
		err  error
	)
	bc := netutil.BackoffConfig{
		MaxDelay:  15 * time.Second,
		BaseDelay: 1.0 * time.Second,
		Factor:    1.6,
		Jitter:    0.2,
	}
	for {
		if msg = <-i.msgs; msg == nil {
			break // quit infoc writeproc
		}
		var j int
		for j = 0; j < _maxRetry; j++ {
			if conn == nil || err != nil {
				if conn, err = net.DialTimeout(i.c.Proto, i.c.Addr, i.dialTimeout); err != nil {
					log.Error("infoc net dial error(%v)", err)
					time.Sleep(bc.Backoff(j))
					continue
				}
			}
			if i.writeTimeout != 0 {
				conn.SetWriteDeadline(time.Now().Add(i.writeTimeout))
			}
			if _, err = conn.Write(msg.Bytes()); err != nil {
				log.Error("infoc conn write error(%v)", err)
				conn.Close()
				time.Sleep(bc.Backoff(j))
				continue
			}
			break
		}
		if j == _maxRetry {
			log.Error("infoc reached max retry times")
		}
		i.freeBuf(msg)
	}
	i.waiter.Done()
	if conn != nil && err == nil {
		conn.Close()
	}
}

func (i *Infoc) info(args ...interface{}) (err error, buf *bytes.Buffer) {
	if len(args) == 0 {
		return nil, nil
	}
	res := i.buf()
	res.Write(_infocMagic)     // type and body buf, for calc length.
	res.Write(_infocHeaderLen) // placeholder
	res.Write(_infocType)
	res.Write(i.header)
	res.WriteString(strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	// // append first arg
	_, err = res.WriteString(getValue(args[0]))
	for _, arg := range args[1:] {
		// append ,arg
		res.WriteString(_infocSpliter)
		_, err = res.WriteString(strings.Replace(getValue(arg), _infocSpliter, _infocReplacer, -1))
	}
	if err != nil {
		i.freeBuf(res)
		return
	}
	bs := res.Bytes()
	binary.BigEndian.PutUint32(bs[_infocLenStart:_infocLenEnd], uint32(res.Len()-_protocolLen))
	buf = res
	return
}
