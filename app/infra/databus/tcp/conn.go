package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type conn struct {
	// conn
	conn net.Conn
	// Read
	readTimeout time.Duration
	br          *bufio.Reader
	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer
	// Scratch space for formatting argument length.
	// '*' or '$', length, "\r\n"
	lenScratch [32]byte
	// Scratch space for formatting integers and floats.
	numScratch [40]byte
}

// newConn returns a new connection for the given net connection.
func newConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) *conn {
	return &conn{
		conn:         netConn,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		br:           bufio.NewReaderSize(netConn, _readBufSize),
		bw:           bufio.NewWriterSize(netConn, _writeBufSize),
	}
}

// Read read data from connection
func (c *conn) Read() (cmd string, args [][]byte, err error) {
	var (
		ln, cn int
		bs     []byte
	)
	if c.readTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	// start read
	if bs, err = c.readLine(); err != nil {
		return
	}
	if len(bs) < 2 {
		err = fmt.Errorf("read error data(%s) from connection", bs)
		return
	}
	// maybe a cmd that without any params is received,such as: QUIT
	if strings.ToLower(string(bs)) == _quit {
		cmd = _quit
		return
	}
	// get param number
	if ln, err = parseLen(bs[1:]); err != nil {
		return
	}
	args = make([][]byte, 0, ln-1)
	for i := 0; i < ln; i++ {
		if cn, err = c.readLen(_protoBulk); err != nil {
			return
		}
		if bs, err = c.readData(cn); err != nil {
			return
		}
		if i == 0 {
			cmd = strings.ToLower(string(bs))
			continue
		}
		args = append(args, bs)
	}
	return
}

// WriteError write error to connection and close connection
func (c *conn) WriteError(err error) {
	if c.writeTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err = c.Write(proto{prefix: _protoErr, message: err.Error()}); err != nil {
		c.Close()
		return
	}
	c.Flush()
	c.Close()
}

// Write write data to connection
func (c *conn) Write(p proto) (err error) {
	if c.writeTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	// start write
	switch p.prefix {
	case _protoStr:
		err = c.writeStatus(p.message)
	case _protoErr:
		err = c.writeError(p.message)
	case _protoInt:
		err = c.writeInt64(int64(p.integer))
	case _protoBulk:
		// c.writeString(p.message)
		err = c.writeBytes([]byte(p.message))
	case _protoArray:
		err = c.writeLen(p.prefix, p.integer)
	}
	return
}

// Flush flush connection
func (c *conn) Flush() error {
	if c.writeTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	return c.bw.Flush()
}

// Close close connection
func (c *conn) Close() error {
	return c.conn.Close()
}

// parseLen parses bulk string and array lengths.
func parseLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errors.New("malformed length")
	}
	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}
	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, errors.New("illegal bytes in length")
		}
		n += int(b - '0')
	}
	return n, nil
}

func (c *conn) readLine() ([]byte, error) {
	p, err := c.br.ReadBytes('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("long response line")
	}
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, errors.New("bad response line terminator")
	}
	return p[:i], nil
}

func (c *conn) readLen(prefix byte) (int, error) {
	ls, err := c.readLine()
	if err != nil {
		return 0, err
	}
	if len(ls) < 2 {
		return 0, errors.New("illegal bytes in length")
	}
	if ls[0] != prefix {
		return 0, errors.New("illegal bytes in length")
	}
	return parseLen(ls[1:])
}

func (c *conn) readData(n int) ([]byte, error) {
	if n > _maxValueSize {
		return nil, errors.New("exceeding max value limit")
	}
	buf := make([]byte, n+2)
	r, err := io.ReadFull(c.br, buf)
	if err != nil {
		return nil, err
	}
	if n != r-2 {
		return nil, errors.New("invalid bytes in len")
	}
	return buf[:n], err
}

func (c *conn) writeLen(prefix byte, n int) error {
	c.lenScratch[len(c.lenScratch)-1] = '\n'
	c.lenScratch[len(c.lenScratch)-2] = '\r'
	i := len(c.lenScratch) - 3
	for {
		c.lenScratch[i] = byte('0' + n%10)
		i--
		n = n / 10
		if n == 0 {
			break
		}
	}
	c.lenScratch[i] = prefix
	_, err := c.bw.Write(c.lenScratch[i:])
	return err
}

func (c *conn) writeStatus(s string) (err error) {
	c.bw.WriteByte(_protoStr)
	c.bw.WriteString(s)
	_, err = c.bw.WriteString("\r\n")
	return
}

func (c *conn) writeError(s string) (err error) {
	c.bw.WriteByte(_protoErr)
	c.bw.WriteString(s)
	_, err = c.bw.WriteString("\r\n")
	return
}

func (c *conn) writeInt64(n int64) (err error) {
	c.bw.WriteByte(_protoInt)
	c.bw.Write(strconv.AppendInt(c.numScratch[:0], n, 10))
	_, err = c.bw.WriteString("\r\n")
	return
}

func (c *conn) writeString(s string) (err error) {
	c.writeLen(_protoBulk, len(s))
	c.bw.WriteString(s)
	_, err = c.bw.WriteString("\r\n")
	return
}

func (c *conn) writeBytes(s []byte) (err error) {
	if len(s) == 0 {
		c.bw.WriteByte('$')
		c.bw.Write(_nullBulk)
	} else {
		c.writeLen(_protoBulk, len(s))
		c.bw.Write(s)
	}
	_, err = c.bw.WriteString("\r\n")
	return
}

func (c *conn) writeStrings(ss []string) (err error) {
	c.writeLen(_protoArray, len(ss))
	for _, s := range ss {
		if err = c.writeString(s); err != nil {
			return
		}
	}
	return
}
