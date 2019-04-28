package sock

import (
	"bufio"
	"runtime"
	"errors"
	"net"
	"os"
	"time"
	"io"
	"path/filepath"
	"fmt"
	"context"
	"strings"

	"go-common/library/log"
	"go-common/app/service/ops/log-agent/input"
	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
	"go-common/app/service/ops/log-agent/pkg/lancermonitor"
	"go-common/app/service/ops/log-agent/pkg/lancerroute"

	"github.com/BurntSushi/toml"
)

const (
	_logIdLen           = 6
	_logLancerHeaderLen = 19
	_logSeparator       = byte('\u0001')
)

var (
	// ErrInvalidAddr invalid address.
	ErrInvalidAddr = errors.New("invalid address")
	// logMagic log magic.
	logMagic = []byte{0xAC, 0xBE}
	local, _ = time.LoadLocation("Local")
)

func init() {
	err := input.Register("sock", NewSock)
	if err != nil {
		panic(err)
	}
}

type Sock struct {
	c        *Config
	output   chan<- *event.ProcessorEvent
	readChan chan *event.ProcessorEvent
	ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
}

func NewSock(ctx context.Context, config interface{}, output chan<- *event.ProcessorEvent) (input.Input, error) {
	sock := new(Sock)
	if c, ok := config.(*Config); !ok {
		return nil, fmt.Errorf("Error config for Sock Input")
	} else {
		if err := c.ConfigValidate(); err != nil {
			return nil, err
		}
		sock.c = c
	}
	sock.output = output
	sock.ctx, sock.cancel = context.WithCancel(ctx)
	sock.readChan = make(chan *event.ProcessorEvent, sock.c.ReadChanSize)
	sock.closed = false

	return sock, nil
}

func DecodeConfig(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	c = new(Config)
	if err = md.PrimitiveDecode(primValue, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Sock) Run() (err error) {
	s.listen()
	go s.writetoProcessor()
	return nil
}

func (s *Sock) Stop() {
	s.closed = true
	s.cancel()
}

func (s *Sock) Ctx() context.Context {
	return s.ctx
}

//listen listen unix socket to buffer
func (s *Sock) listen() {
	// SOCK_DGRAM
	if s.c.UdpAddr != "" {
		if flag, _ := pathExists(s.c.UdpAddr); flag {
			os.Remove(s.c.UdpAddr)
		}
		if flag, _ := pathExists(filepath.Dir(s.c.UdpAddr)); !flag {
			os.Mkdir(filepath.Dir(s.c.UdpAddr), os.ModePerm)
		}
		udpconn, err := net.ListenPacket("unixgram", s.c.UdpAddr)
		if err != nil {
			panic(err)
		}
		if err := os.Chmod(s.c.UdpAddr, 0777); err != nil {
			panic(err)
		}
		log.Info("start listen: %s", s.c.UdpAddr)
		for i := 0; i < runtime.NumCPU(); i++ {
			go s.udpread(udpconn)
		}
	}
	// SOCK_SEQPACKET
	if s.c.TcpAddr != "" {
		if flag, _ := pathExists(s.c.TcpAddr); flag {
			os.Remove(s.c.TcpAddr)
		}
		if flag, _ := pathExists(filepath.Dir(s.c.TcpAddr)); !flag {
			os.Mkdir(filepath.Dir(s.c.TcpAddr), os.ModePerm)
		}
		tcplistener, err := net.Listen("unixpacket", s.c.TcpAddr)
		if err != nil {
			panic(err)
		}
		if err := os.Chmod(s.c.TcpAddr, 0777); err != nil {
			panic(err)
		}
		log.Info("start listen: %s", s.c.TcpAddr)
		go s.tcpread(tcplistener)
	}
}

func (s *Sock) writeToReadChan(e *event.ProcessorEvent) {
	lancermonitor.IncreaseLogCount("agent.receive.count", e.LogId)
	select {
	case s.readChan <- e:
	default:
		event.PutEvent(e)
		flowmonitor.Fm.AddEvent(e, "log-agent.input.sock", "ERROR", "read chan full")
		log.Warn("sock read chan full, discard log")
	}
}

// tcpread accept tcp connection
func (s *Sock) tcpread(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleTcpConn(conn)
	}
}

// process single tcp connection
func (s *Sock) handleTcpConn(conn net.Conn) {
	defer conn.Close()
	rd := bufio.NewReaderSize(conn, s.c.TcpBatchMaxBytes)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(s.c.TcpReadTimeout)))
		b, err := rd.ReadSlice(_logSeparator)

		if err == nil && len(b) <= _logLancerHeaderLen {
			continue
		}

		if err == nil {
			e := s.preproccess(b[:len(b)-1])
			if e != nil {
				s.writeToReadChan(e)
				flowmonitor.Fm.AddEvent(e, "log-agent.input.sock", "OK", "received")
				continue
			}
		}

		// conn closed and return EOF
		if err == io.EOF {
			e := s.preproccess(b)
			if e != nil {
				s.writeToReadChan(e)
				flowmonitor.Fm.AddEvent(e, "log-agent.input.sock", "OK", "received")
				continue
			}
			log.Info("get EOF from conn, close conn")
			return
		}

		log.Error("read from tcp conn error(%v). close conn", err)
		return
	}
}

//read read from unix socket conn
func (s *Sock) udpread(conn net.PacketConn) {
	b := make([]byte, s.c.UdpPacketMaxSize)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(s.c.UdpReadTimeout)))
		l, _, err := conn.ReadFrom(b)

		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			log.Error("conn.ReadFrom() error(%v)", err)
			continue
		}

		e := s.preproccess(b[:l])

		if e != nil {
			flowmonitor.Fm.AddEvent(e, "log-agent.input.sock", "OK", "received")
			s.writeToReadChan(e)
		}
	}
}

func (s *Sock) preproccess(b []byte) *event.ProcessorEvent {
	if len(b) <= _logLancerHeaderLen {
		return nil
	}
	e := event.GetEvent()
	e.LogId = string(b[:_logIdLen])
	e.Destination = lancerroute.GetLancerByLogid(e.LogId)
	e.Write(b[_logLancerHeaderLen:])
	e.Source = "sock"
	flowmonitor.Fm.AddEvent(e, "log-agent.input.sock", "OK", "received")
	return e
}

// pathExists judge if the file exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Sock) writetoProcessor() {
	for {
		select {
		case e := <-s.readChan:
			s.output <- e
		case <-s.ctx.Done():
			return
		}
	}
}
