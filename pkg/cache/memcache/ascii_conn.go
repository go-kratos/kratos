package memcache

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	pkgerr "github.com/pkg/errors"
)

var (
	crlf                   = []byte("\r\n")
	space                  = []byte(" ")
	replyOK                = []byte("OK\r\n")
	replyStored            = []byte("STORED\r\n")
	replyNotStored         = []byte("NOT_STORED\r\n")
	replyExists            = []byte("EXISTS\r\n")
	replyNotFound          = []byte("NOT_FOUND\r\n")
	replyDeleted           = []byte("DELETED\r\n")
	replyEnd               = []byte("END\r\n")
	replyTouched           = []byte("TOUCHED\r\n")
	replyClientErrorPrefix = []byte("CLIENT_ERROR ")
	replyServerErrorPrefix = []byte("SERVER_ERROR ")
)

var _ protocolConn = &asiiConn{}

// asiiConn is the low-level implementation of Conn
type asiiConn struct {
	err  error
	conn net.Conn
	// Read & Write
	readTimeout  time.Duration
	writeTimeout time.Duration
	rw           *bufio.ReadWriter
}

func replyToError(line []byte) error {
	switch {
	case bytes.Equal(line, replyStored):
		return nil
	case bytes.Equal(line, replyOK):
		return nil
	case bytes.Equal(line, replyDeleted):
		return nil
	case bytes.Equal(line, replyTouched):
		return nil
	case bytes.Equal(line, replyNotStored):
		return ErrNotStored
	case bytes.Equal(line, replyExists):
		return ErrCASConflict
	case bytes.Equal(line, replyNotFound):
		return ErrNotFound
	case bytes.Equal(line, replyNotStored):
		return ErrNotStored
	case bytes.Equal(line, replyExists):
		return ErrCASConflict
	}
	return pkgerr.WithStack(protocolError(string(line)))
}

func (c *asiiConn) Populate(ctx context.Context, cmd string, key string, flags uint32, expiration int32, cas uint64, data []byte) error {
	var err error
	c.conn.SetWriteDeadline(shrinkDeadline(ctx, c.writeTimeout))
	// <command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
	if cmd == "cas" {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d %d\r\n", cmd, key, flags, expiration, len(data), cas)
	} else {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d\r\n", cmd, key, flags, expiration, len(data))
	}
	if err != nil {
		return c.fatal(err)
	}
	c.rw.Write(data)
	c.rw.Write(crlf)
	if err = c.rw.Flush(); err != nil {
		return c.fatal(err)
	}
	c.conn.SetReadDeadline(shrinkDeadline(ctx, c.readTimeout))
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return c.fatal(err)
	}
	return replyToError(line)
}

// newConn returns a new memcache connection for the given net connection.
func newASCIIConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) (protocolConn, error) {
	if writeTimeout <= 0 || readTimeout <= 0 {
		return nil, pkgerr.Errorf("readTimeout writeTimeout can't be zero")
	}
	c := &asiiConn{
		conn: netConn,
		rw: bufio.NewReadWriter(bufio.NewReader(netConn),
			bufio.NewWriter(netConn)),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
	return c, nil
}

func (c *asiiConn) Close() error {
	if c.err == nil {
		c.err = pkgerr.New("memcache: closed")
	}
	return c.conn.Close()
}

func (c *asiiConn) fatal(err error) error {
	if c.err == nil {
		c.err = pkgerr.WithStack(err)
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	}
	return c.err
}

func (c *asiiConn) Err() error {
	return c.err
}

func (c *asiiConn) Get(ctx context.Context, key string) (result *Item, err error) {
	c.conn.SetWriteDeadline(shrinkDeadline(ctx, c.writeTimeout))
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", key); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.parseGetReply(ctx, func(it *Item) {
		result = it
	}); err != nil {
		return
	}
	if result == nil {
		return nil, ErrNotFound
	}
	return
}

func (c *asiiConn) GetMulti(ctx context.Context, keys ...string) (map[string]*Item, error) {
	var err error
	c.conn.SetWriteDeadline(shrinkDeadline(ctx, c.writeTimeout))
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", strings.Join(keys, " ")); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	results := make(map[string]*Item, len(keys))
	if err = c.parseGetReply(ctx, func(it *Item) {
		results[it.Key] = it
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *asiiConn) parseGetReply(ctx context.Context, f func(*Item)) error {
	c.conn.SetReadDeadline(shrinkDeadline(ctx, c.readTimeout))
	for {
		line, err := c.rw.ReadSlice('\n')
		if err != nil {
			return c.fatal(err)
		}
		if bytes.Equal(line, replyEnd) {
			return nil
		}
		if bytes.HasPrefix(line, replyServerErrorPrefix) {
			errMsg := line[len(replyServerErrorPrefix):]
			return c.fatal(protocolError(errMsg))
		}
		it := new(Item)
		size, err := scanGetReply(line, it)
		if err != nil {
			return c.fatal(err)
		}
		it.Value = make([]byte, size+2)
		if _, err = io.ReadFull(c.rw, it.Value); err != nil {
			return c.fatal(err)
		}
		if !bytes.HasSuffix(it.Value, crlf) {
			return c.fatal(protocolError("corrupt get reply, no except CRLF"))
		}
		it.Value = it.Value[:size]
		f(it)
	}
}

func scanGetReply(line []byte, item *Item) (size int, err error) {
	pattern := "VALUE %s %d %d %d\r\n"
	dest := []interface{}{&item.Key, &item.Flags, &size, &item.cas}
	if bytes.Count(line, space) == 3 {
		pattern = "VALUE %s %d %d\r\n"
		dest = dest[:3]
	}
	n, err := fmt.Sscanf(string(line), pattern, dest...)
	if err != nil || n != len(dest) {
		return -1, fmt.Errorf("memcache: unexpected line in get response: %q", line)
	}
	return size, nil
}

func (c *asiiConn) Touch(ctx context.Context, key string, expire int32) error {
	line, err := c.writeReadLine(ctx, "touch %s %d\r\n", key, expire)
	if err != nil {
		return err
	}
	return replyToError(line)
}

func (c *asiiConn) IncrDecr(ctx context.Context, cmd, key string, delta uint64) (uint64, error) {
	line, err := c.writeReadLine(ctx, "%s %s %d\r\n", cmd, key, delta)
	if err != nil {
		return 0, err
	}
	switch {
	case bytes.Equal(line, replyNotFound):
		return 0, ErrNotFound
	case bytes.HasPrefix(line, replyClientErrorPrefix):
		errMsg := line[len(replyClientErrorPrefix):]
		return 0, pkgerr.WithStack(protocolError(errMsg))
	}
	val, err := strconv.ParseUint(string(line[:len(line)-2]), 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *asiiConn) Delete(ctx context.Context, key string) error {
	line, err := c.writeReadLine(ctx, "delete %s\r\n", key)
	if err != nil {
		return err
	}
	return replyToError(line)
}

func (c *asiiConn) writeReadLine(ctx context.Context, format string, args ...interface{}) ([]byte, error) {
	var err error
	c.conn.SetWriteDeadline(shrinkDeadline(ctx, c.writeTimeout))
	_, err = fmt.Fprintf(c.rw, format, args...)
	if err != nil {
		return nil, c.fatal(pkgerr.WithStack(err))
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(pkgerr.WithStack(err))
	}
	c.conn.SetReadDeadline(shrinkDeadline(ctx, c.readTimeout))
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return line, c.fatal(pkgerr.WithStack(err))
	}
	return line, nil
}
