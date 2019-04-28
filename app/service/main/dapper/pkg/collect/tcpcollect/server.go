package tcpcollect

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/model"
	"go-common/app/service/main/dapper/pkg/process"
	"go-common/library/log"
	protogen "go-common/library/net/trace/proto"
	"go-common/library/stat/counter"
	"go-common/library/stat/prom"
)

var (
	collectCount    = prom.New().WithCounter("dapper_collect_count", []string{"remote_host"})
	collectErrCount = prom.New().WithCounter("dapper_collect_err_count", []string{"remote_host"})
)

const (
	_magicSize  = 2
	_headerSize = 6
)

var (
	_magicBuf  = []byte{0xAC, 0xBE}
	_separator = []byte("\001")
)

// ClientStatus agent client status
type ClientStatus struct {
	Addr         string
	Counter      counter.Counter
	ErrorCounter counter.Counter
	UpTime       int64
}

func (c *ClientStatus) incr(iserr bool) {
	if iserr {
		collectErrCount.Incr(c.ClientHost())
	}
	collectCount.Incr(c.ClientHost())
	c.Counter.Add(1)
}

// ClientHost extract from client addr
func (c *ClientStatus) ClientHost() string {
	host, _, _ := net.SplitHostPort(c.Addr)
	return host
}

// TCPCollect tcp server.
type TCPCollect struct {
	cfg       *conf.Collect
	lis       net.Listener
	clientMap map[string]*ClientStatus
	rmx       sync.RWMutex
	ps        []process.Processer
}

// New tcp server.
func New(cfg *conf.Collect) *TCPCollect {
	svr := &TCPCollect{
		cfg:       cfg,
		clientMap: make(map[string]*ClientStatus),
	}
	return svr
}

// RegisterProcess implement process.Processer
func (s *TCPCollect) RegisterProcess(p process.Processer) {
	s.ps = append(s.ps, p)
}

func (s *TCPCollect) addClient(cs *ClientStatus) {
	s.rmx.Lock()
	defer s.rmx.Unlock()
	s.clientMap[cs.Addr] = cs
}

func (s *TCPCollect) removeClient(cs *ClientStatus) {
	s.rmx.Lock()
	defer s.rmx.Unlock()
	delete(s.clientMap, cs.Addr)
}

// ClientStatus ClientStatus
func (s *TCPCollect) ClientStatus() []*ClientStatus {
	s.rmx.RLock()
	defer s.rmx.RUnlock()
	css := make([]*ClientStatus, 0, len(s.clientMap))
	for _, cs := range s.clientMap {
		css = append(css, cs)
	}
	return css
}

// Start tcp server.
func (s *TCPCollect) Start() error {
	var err error
	if s.lis, err = net.Listen(s.cfg.Network, s.cfg.Addr); err != nil {
		return err
	}
	go func() {
		for {
			conn, err := s.lis.Accept()
			if err != nil {
				if netE, ok := err.(net.Error); ok && netE.Temporary() {
					log.Error("l.Accept() error(%v)", err)
					time.Sleep(time.Second)
					continue
				}
				return
			}
			go s.serveConn(conn)
		}
	}()
	log.Info("tcp server start addr:%s@%s", s.cfg.Network, s.cfg.Addr)
	return nil
}

// Close tcp server.
func (s *TCPCollect) Close() error {
	return s.lis.Close()
}

func (s *TCPCollect) serveConn(conn net.Conn) {
	log.Info("serverConn remoteIP:%s", conn.RemoteAddr().String())
	cs := &ClientStatus{
		Addr:         conn.RemoteAddr().String(),
		Counter:      counter.NewRolling(time.Second, 100),
		ErrorCounter: counter.NewGauge(),
		UpTime:       time.Now().Unix(),
	}
	s.addClient(cs)
	defer conn.Close()
	defer s.removeClient(cs)
	rd := bufio.NewReaderSize(conn, 65536)
	for {
		buf, err := s.tailPacket(rd)
		if err != nil {
			log.Error("s.tailPacket() remoteIP:%s error(%v)", conn.RemoteAddr().String(), err)
			cs.incr(true)
			return
		}
		if len(buf) == 0 {
			log.Error("s.tailPacket() is empty")
			cs.incr(true)
			continue
		}
		data := buf
		fields := bytes.Split(buf, _separator)
		if len(fields) >= 16 {
			if data, err = s.legacySpan(fields[2:]); err != nil {
				log.Error("convert legacy span error: %s", err)
				continue
			}
		}
		protoSpan := new(protogen.Span)
		if err = proto.Unmarshal(data, protoSpan); err != nil {
			log.Error("unmarshal data %s error: %s", err, data)
			continue
		}
		for _, p := range s.ps {
			if pe := p.Process(context.Background(), (*model.ProtoSpan)(protoSpan)); pe != nil {
				log.Error("process span %s error: %s", protoSpan, err)
			}
		}
		cs.incr(err != nil)
	}
}

func (s *TCPCollect) tailPacket(rr *bufio.Reader) (res []byte, err error) {
	var buf []byte
	// peek magic
	for {
		if buf, err = rr.Peek(_magicSize); err != nil {
			return
		}
		if bytes.Equal(buf, _magicBuf) {
			break
		}
		rr.Discard(1)
	}
	// peek length
	if buf, err = rr.Peek(_headerSize); err != nil {
		return
	}
	// peek body
	packetLen := int(binary.BigEndian.Uint32(buf[_magicSize:_headerSize]))
	if buf, err = rr.Peek(_headerSize + packetLen); err != nil {
		return
	}
	res = buf[_headerSize+_magicSize:]
	rr.Discard(packetLen + _headerSize)
	return
}

// startTime/endTime/traceID/spanID/parentID/event/level/class/sample/address/family/title/comment/caller/error
func (s *TCPCollect) legacySpan(fields [][]byte) ([]byte, error) {
	startAt, _ := strconv.ParseInt(string(fields[0]), 10, 64)
	finishAt, _ := strconv.ParseInt(string(fields[1]), 10, 64)
	traceID, _ := strconv.ParseUint(string(fields[2]), 10, 64)
	spanID, _ := strconv.ParseUint(string(fields[3]), 10, 64)
	parentID, _ := strconv.ParseUint(string(fields[4]), 10, 64)
	event, _ := strconv.Atoi(string(fields[5]))
	start := 8
	if len(fields) == 14 {
		start = 7
	}
	address := string(fields[start+1])
	family := string(fields[start+2])
	title := string(fields[start+3])
	comment := string(fields[start+4])
	caller := string(fields[start+5])
	errMsg := string(fields[start+6])

	span := &protogen.Span{Version: 2}
	span.ServiceName = family
	span.OperationName = title
	span.Caller = caller
	span.TraceId = traceID
	span.SpanId = spanID
	span.ParentId = parentID
	span.StartTime = &timestamp.Timestamp{
		Seconds: startAt / int64(time.Second),
		Nanos:   int32(startAt % int64(time.Second)),
	}
	d := finishAt - startAt
	span.Duration = &duration.Duration{
		Seconds: d / int64(time.Second),
		Nanos:   int32(d % int64(time.Second)),
	}
	if event == 0 {
		span.Tags = append(span.Tags, &protogen.Tag{Key: "span.kind", Kind: protogen.Tag_STRING, Value: []byte("client")})
	} else {
		span.Tags = append(span.Tags, &protogen.Tag{Key: "span.kind", Kind: protogen.Tag_STRING, Value: []byte("server")})
	}
	span.Tags = append(span.Tags, &protogen.Tag{Key: "legacy.address", Kind: protogen.Tag_STRING, Value: []byte(address)})
	span.Tags = append(span.Tags, &protogen.Tag{Key: "legacy.comment", Kind: protogen.Tag_STRING, Value: []byte(comment)})
	if errMsg != "" {
		span.Logs = append(span.Logs, &protogen.Log{Key: "legacy.error", Fields: []*protogen.Field{&protogen.Field{Key: "error", Value: []byte(errMsg)}}})
	}
	return proto.Marshal(span)
}
