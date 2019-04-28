package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go-common/library/log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LancerMinimumLength   = 6
	LancerLengthPartBegin = 2
	LancerLengthPartEnd   = 6
)

type LancerLogStream struct {
	conn        net.Conn
	addr        string
	pool        sync.Pool
	dataChannel chan *LancerData
	quit        chan struct{}
	wg          sync.WaitGroup
	splitter    string
	replacer    string
	timeout     time.Duration
}

type LancerData struct {
	bytes.Buffer
	logid    string
	isAppend bool
	lancer   *LancerLogStream
}

func NewLancerLogStream(address string, capacity int, timeout time.Duration) *LancerLogStream {
	c, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Error("[Lancer]Dial lancer %s error:%+v", address, err)
		c = nil
	}
	lancer := &LancerLogStream{
		conn: c,
		addr: address,
		pool: sync.Pool{
			New: func() interface{} {
				return new(LancerData)
			},
		},
		dataChannel: make(chan *LancerData, capacity),
		quit:        make(chan struct{}),
		splitter:    "\u0001",
		replacer:    "|",
		timeout:     timeout,
	}
	lancer.wg.Add(1)
	go lancer.processor()
	return lancer
}

func (lancer *LancerLogStream) Close() {
	close(lancer.dataChannel)
	close(lancer.quit)
	lancer.wg.Wait()
}

func (lancer *LancerLogStream) processor() {
	defer func() {
		if lancer.conn != nil {
			lancer.conn.Close()
		}
		lancer.wg.Done()
	}()
	var lastFail *LancerData
PROCESSOR:
	for {
		select {
		case <-lancer.quit:
			return
		default:
		}
		if lastFail != nil {
			if err := lancer.write(lastFail.Bytes()); err != nil {
				runtime.Gosched()
				continue PROCESSOR
			}
			lancer.pool.Put(lastFail)
			lastFail = nil
		}
		for b := range lancer.dataChannel {
			if err := lancer.write(b.Bytes()); err != nil {
				lastFail = b
				runtime.Gosched()
				continue PROCESSOR
			}
			lancer.pool.Put(b)
		}
		return
	}
}

func (lancer *LancerLogStream) write(b []byte) error {
	if lancer.conn == nil {
		c, err := net.DialTimeout("tcp", lancer.addr, lancer.timeout)
		if err != nil {
			log.Error("[Lancer]Dial %s error:%+v", lancer.addr, err)
			return err
		}
		lancer.conn = c
	}
	_, err := lancer.conn.Write(b)
	if err != nil {
		log.Error("[Lancer]Conn write error:%+v", err)
		lancer.conn.Close()
		lancer.conn = nil
	}
	return err
}

func (lancer *LancerLogStream) NewLancerData(logid string, token string) *LancerData {
	ld := lancer.pool.Get().(*LancerData)
	ld.Reset()
	ld.lancer = lancer
	ld.isAppend = false
	ld.logid = logid
	ld.Write([]byte{0xAC, 0xBE})
	ld.Write([]byte{0, 0, 0, 0})
	ld.Write([]byte{0, 1})
	header := fmt.Sprintf("logId=%s&timestamp=%d&token=%s&version=1.1", logid,
		time.Now().UnixNano()/int64(time.Millisecond), token)
	headerLength := uint16(len(header))
	ld.Write([]byte{byte(headerLength >> 8), byte(headerLength)})
	ld.Write([]byte(header))
	return ld
}

func (ld *LancerData) PutString(v string) *LancerData {
	ld.splitter()
	ld.WriteString(strings.Replace(v, ld.lancer.splitter, ld.lancer.replacer, -1))
	return ld
}

func (ld *LancerData) PutTimestamp(v time.Time) *LancerData {
	return ld.PutInt(v.Unix())
}

func (ld *LancerData) PutUint(v uint64) *LancerData {
	ld.splitter()
	ld.WriteString(strconv.FormatUint(v, 10))
	return ld
}

func (ld *LancerData) PutInt(v int64) *LancerData {
	ld.splitter()
	ld.WriteString(strconv.FormatInt(v, 10))
	return ld
}

func (ld *LancerData) PutFloat(v float64) *LancerData {
	ld.splitter()
	ld.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	return ld
}

func (ld *LancerData) PutBool(v bool) *LancerData {
	ld.splitter()
	ld.WriteString(strconv.FormatBool(v))
	return ld
}

func (ld *LancerData) splitter() {
	if ld.isAppend {
		ld.WriteString(ld.lancer.splitter)
	}
	ld.isAppend = true
}

func (ld *LancerData) Commit() error {
	if ld.Len() < LancerMinimumLength {
		return errors.New("protocol error")
	}
	l := uint32(ld.Len()) - LancerMinimumLength
	binary.BigEndian.PutUint32(ld.Bytes()[LancerLengthPartBegin:LancerLengthPartEnd], l)
	select {
	case ld.lancer.dataChannel <- ld:
		return nil
	default:
		ld.lancer.pool.Put(ld)
		return errors.New("lancer channel is full")
	}
}
