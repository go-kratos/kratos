package statsd

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"go-common/library/log"
)

const (
	_quit = ""
	_size = 1400
)

// Config statsd config.
type Config struct {
	Project  string
	Addr     string
	ChanSize int
}

// Statsd statsd struct.
type Statsd struct {
	// project.hostname.api
	// Make sure no '/' in the api.
	c        *Config
	business string
	r        *rand.Rand
	stats    chan string
}

// New new a statsd struct.
func New(c *Config) (s *Statsd) {
	s = new(Statsd)
	s.c = c
	s.business = fmt.Sprintf("%s", c.Project)
	// init rand
	s.r = rand.New(rand.NewSource(time.Now().Unix()))
	// init stat channel
	s.stats = make(chan string, c.ChanSize)
	go s.writeproc()
	return
}

// send data to udp statsd daemon
func (s *Statsd) send(data string, rate float32) {
	if rate < 1 && s.r != nil {
		if s.r.Float32() < rate {
			return
		}
	}
	select {
	case s.stats <- data:
	default:
		log.Warn("Statsd stat channel is full")
	}
}

// writeproc write data into connection.
func (s *Statsd) writeproc() {
	var (
		err  error
		l    int
		stat string
		conn net.Conn
		buf  bytes.Buffer
		tick = time.Tick(1 * time.Second)
	)
	for {
		select {
		case stat = <-s.stats:
			if stat == _quit {
				if conn != nil {
					conn.Close()
				}
				return
			}
		case <-tick:
			if l = buf.Len(); l > 0 {
				conn.Write(buf.Bytes()[:l-1])
				buf.Reset()
			}
			continue
		}
		if conn == nil {
			if conn, err = net.Dial("udp", s.c.Addr); err != nil {
				log.Error("net.Dial('udp', %s) error(%v)", s.c.Addr, err)
				time.Sleep(time.Second)
				continue
			}
		}
		if l = buf.Len(); l+len(stat) >= _size {
			conn.Write(buf.Bytes()[:l-1])
			buf.Reset()
		}
		buf.WriteString(stat)
		buf.WriteByte('\n')
	}
}

// Close close the connection.
func (s *Statsd) Close() {
	s.stats <- _quit
}

// Timing log timing information (in milliseconds) without sampling
func (s *Statsd) Timing(name string, time int64, extra ...string) {
	val := formatTiming(s.business, name, time, extra...)
	s.send(val, 1)
}

// Incr increments one stat counter without sampling
func (s *Statsd) Incr(name string, extra ...string) {
	val := formatIncr(s.business, name, extra...)
	s.send(val, 1)
}

// State set state
func (s *Statsd) State(stat string, val int64, extra ...string) {
	return
}

func formatIncr(business, name string, extra ...string) string {
	ss := []string{business, name}
	ss = append(ss, extra...)
	return strings.Join(ss, ".") + ":1|c"
}

func formatTiming(business, name string, time int64, extra ...string) string {
	ss := []string{business, name}
	ss = append(ss, extra...)
	return strings.Join(ss, ".") + fmt.Sprintf(":%d|ms", time)
}
