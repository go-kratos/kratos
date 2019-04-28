package memcache

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	pkgerr "github.com/pkg/errors"
)

var (
	crlf                   = []byte("\r\n")
	spaceStr               = string(" ")
	replyOK                = []byte("OK\r\n")
	replyStored            = []byte("STORED\r\n")
	replyNotStored         = []byte("NOT_STORED\r\n")
	replyExists            = []byte("EXISTS\r\n")
	replyNotFound          = []byte("NOT_FOUND\r\n")
	replyDeleted           = []byte("DELETED\r\n")
	replyEnd               = []byte("END\r\n")
	replyTouched           = []byte("TOUCHED\r\n")
	replyValueStr          = "VALUE"
	replyClientErrorPrefix = []byte("CLIENT_ERROR ")
	replyServerErrorPrefix = []byte("SERVER_ERROR ")
)

const (
	_encodeBuf = 4096 // 4kb
	// 1024*1024 - 1, set error???
	_largeValue = 1000 * 1000 // 1MB
)

type reader struct {
	io.Reader
}

func (r *reader) Reset(rd io.Reader) {
	r.Reader = rd
}

// conn is the low-level implementation of Conn
type conn struct {
	// Shared
	mu   sync.Mutex
	err  error
	conn net.Conn
	// Read & Write
	readTimeout  time.Duration
	writeTimeout time.Duration
	rw           *bufio.ReadWriter
	// Item Reader
	ir bytes.Reader
	// Compress
	gr gzip.Reader
	gw *gzip.Writer
	cb bytes.Buffer
	// Encoding
	edb bytes.Buffer
	// json
	jr reader
	jd *json.Decoder
	je *json.Encoder
	// protobuffer
	ped *proto.Buffer
}

// DialOption specifies an option for dialing a Memcache server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	dial         func(network, addr string) (net.Conn, error)
}

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = d
	}}
}

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = d
	}}
}

// DialConnectTimeout specifies the timeout for connecting to the Memcache server.
func DialConnectTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		dialer := net.Dialer{Timeout: d}
		do.dial = dialer.Dial
	}}
}

// DialNetDial specifies a custom dial function for creating TCP
// connections. If this option is left out, then net.Dial is
// used. DialNetDial overrides DialConnectTimeout.
func DialNetDial(dial func(network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dial = dial
	}}
}

// Dial connects to the Memcache server at the given network and
// address using the specified options.
func Dial(network, address string, options ...DialOption) (Conn, error) {
	do := dialOptions{
		dial: net.Dial,
	}
	for _, option := range options {
		option.f(&do)
	}
	netConn, err := do.dial(network, address)
	if err != nil {
		return nil, pkgerr.WithStack(err)
	}
	return NewConn(netConn, do.readTimeout, do.writeTimeout), nil
}

// NewConn returns a new memcache connection for the given net connection.
func NewConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) Conn {
	if writeTimeout <= 0 || readTimeout <= 0 {
		panic("must config memcache timeout")
	}
	c := &conn{
		conn: netConn,
		rw: bufio.NewReadWriter(bufio.NewReader(netConn),
			bufio.NewWriter(netConn)),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
	c.jd = json.NewDecoder(&c.jr)
	c.je = json.NewEncoder(&c.edb)
	c.gw = gzip.NewWriter(&c.cb)
	c.edb.Grow(_encodeBuf)
	// NOTE reuse bytes.Buffer internal buf
	// DON'T concurrency call Scan
	c.ped = proto.NewBuffer(c.edb.Bytes())
	return c
}

func (c *conn) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = pkgerr.New("memcache: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = pkgerr.WithStack(err)
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	}
	c.mu.Unlock()
	return c.err
}

func (c *conn) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *conn) Add(item *Item) error {
	return c.populate("add", item)
}

func (c *conn) Set(item *Item) error {
	return c.populate("set", item)
}

func (c *conn) Replace(item *Item) error {
	return c.populate("replace", item)
}

func (c *conn) CompareAndSwap(item *Item) error {
	return c.populate("cas", item)
}

func (c *conn) populate(cmd string, item *Item) (err error) {
	if !legalKey(item.Key) {
		return pkgerr.WithStack(ErrMalformedKey)
	}
	var res []byte
	if res, err = c.encode(item); err != nil {
		return
	}
	l := len(res)
	count := l/(_largeValue) + 1
	if count == 1 {
		item.Value = res
		return c.populateOne(cmd, item)
	}
	nItem := &Item{
		Key:        item.Key,
		Value:      []byte(strconv.Itoa(l)),
		Expiration: item.Expiration,
		Flags:      item.Flags | flagLargeValue,
	}
	err = c.populateOne(cmd, nItem)
	if err != nil {
		return
	}
	k := item.Key
	nItem.Flags = item.Flags
	for i := 1; i <= count; i++ {
		if i == count {
			nItem.Value = res[_largeValue*(count-1):]
		} else {
			nItem.Value = res[_largeValue*(i-1) : _largeValue*i]
		}
		nItem.Key = fmt.Sprintf("%s%d", k, i)
		if err = c.populateOne(cmd, nItem); err != nil {
			return
		}
	}
	return
}

func (c *conn) populateOne(cmd string, item *Item) (err error) {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	// <command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
	if cmd == "cas" {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d %d\r\n",
			cmd, item.Key, item.Flags, item.Expiration, len(item.Value), item.cas)
	} else {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d\r\n",
			cmd, item.Key, item.Flags, item.Expiration, len(item.Value))
	}
	if err != nil {
		return c.fatal(err)
	}
	c.rw.Write(item.Value)
	c.rw.Write(crlf)
	if err = c.rw.Flush(); err != nil {
		return c.fatal(err)
	}
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return c.fatal(err)
	}
	switch {
	case bytes.Equal(line, replyStored):
		return nil
	case bytes.Equal(line, replyNotStored):
		return ErrNotStored
	case bytes.Equal(line, replyExists):
		return ErrCASConflict
	case bytes.Equal(line, replyNotFound):
		return ErrNotFound
	}
	return pkgerr.WithStack(protocolError(string(line)))
}

func (c *conn) Get(key string) (r *Item, err error) {
	if !legalKey(key) {
		return nil, pkgerr.WithStack(ErrMalformedKey)
	}
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", key); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.parseGetReply(func(it *Item) {
		r = it
	}); err != nil {
		return
	}
	if r == nil {
		err = ErrNotFound
		return
	}
	if r.Flags&flagLargeValue != flagLargeValue {
		return
	}
	if r, err = c.getLargeValue(r); err != nil {
		return
	}
	return
}

func (c *conn) GetMulti(keys []string) (res map[string]*Item, err error) {
	for _, key := range keys {
		if !legalKey(key) {
			return nil, pkgerr.WithStack(ErrMalformedKey)
		}
	}
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", strings.Join(keys, " ")); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	res = make(map[string]*Item, len(keys))
	if err = c.parseGetReply(func(it *Item) {
		res[it.Key] = it
	}); err != nil {
		return
	}
	for k, v := range res {
		if v.Flags&flagLargeValue != flagLargeValue {
			continue
		}
		r, err := c.getLargeValue(v)
		if err != nil {
			return res, err
		}
		res[k] = r
	}
	return
}

func (c *conn) getMulti(keys []string) (res map[string]*Item, err error) {
	for _, key := range keys {
		if !legalKey(key) {
			return nil, pkgerr.WithStack(ErrMalformedKey)
		}
	}
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", strings.Join(keys, " ")); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	res = make(map[string]*Item, len(keys))
	err = c.parseGetReply(func(it *Item) {
		res[it.Key] = it
	})
	return
}

func (c *conn) getLargeValue(it *Item) (r *Item, err error) {
	l, err := strconv.Atoi(string(it.Value))
	if err != nil {
		return
	}
	count := l/_largeValue + 1
	keys := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		keys = append(keys, fmt.Sprintf("%s%d", it.Key, i))
	}
	items, err := c.getMulti(keys)
	if err != nil {
		return
	}
	if len(items) < count {
		err = ErrNotFound
		return
	}
	v := make([]byte, 0, l)
	for _, k := range keys {
		if items[k] == nil || items[k].Value == nil {
			err = ErrNotFound
			return
		}
		v = append(v, items[k].Value...)
	}
	it.Value = v
	it.Flags = it.Flags ^ flagLargeValue
	r = it
	return
}

func (c *conn) parseGetReply(f func(*Item)) error {
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
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
	if !bytes.HasSuffix(line, crlf) {
		return 0, protocolError("corrupt get reply, no except CRLF")
	}
	// VALUE <key> <flags> <bytes> [<cas unique>]
	chunks := strings.Split(string(line[:len(line)-2]), spaceStr)
	if len(chunks) < 4 {
		return 0, protocolError("corrupt get reply")
	}
	if chunks[0] != replyValueStr {
		return 0, protocolError("corrupt get reply, no except VALUE")
	}
	item.Key = chunks[1]
	flags64, err := strconv.ParseUint(chunks[2], 10, 32)
	if err != nil {
		return 0, err
	}
	item.Flags = uint32(flags64)
	if size, err = strconv.Atoi(chunks[3]); err != nil {
		return
	}
	if len(chunks) > 4 {
		item.cas, err = strconv.ParseUint(chunks[4], 10, 64)
	}
	return
}

func (c *conn) Touch(key string, expire int32) (err error) {
	if !legalKey(key) {
		return pkgerr.WithStack(ErrMalformedKey)
	}
	line, err := c.writeReadLine("touch %s %d\r\n", key, expire)
	if err != nil {
		return err
	}
	switch {
	case bytes.Equal(line, replyTouched):
		return nil
	case bytes.Equal(line, replyNotFound):
		return ErrNotFound
	default:
		return pkgerr.WithStack(protocolError(string(line)))
	}
}

func (c *conn) Increment(key string, delta uint64) (uint64, error) {
	return c.incrDecr("incr", key, delta)
}

func (c *conn) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return c.incrDecr("decr", key, delta)
}

func (c *conn) incrDecr(cmd, key string, delta uint64) (uint64, error) {
	if !legalKey(key) {
		return 0, pkgerr.WithStack(ErrMalformedKey)
	}
	line, err := c.writeReadLine("%s %s %d\r\n", cmd, key, delta)
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

func (c *conn) Delete(key string) (err error) {
	if !legalKey(key) {
		return pkgerr.WithStack(ErrMalformedKey)
	}
	line, err := c.writeReadLine("delete %s\r\n", key)
	if err != nil {
		return err
	}
	switch {
	case bytes.Equal(line, replyOK):
		return nil
	case bytes.Equal(line, replyDeleted):
		return nil
	case bytes.Equal(line, replyNotStored):
		return ErrNotStored
	case bytes.Equal(line, replyExists):
		return ErrCASConflict
	case bytes.Equal(line, replyNotFound):
		return ErrNotFound
	}
	return pkgerr.WithStack(protocolError(string(line)))
}

func (c *conn) writeReadLine(format string, args ...interface{}) ([]byte, error) {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	_, err := fmt.Fprintf(c.rw, format, args...)
	if err != nil {
		return nil, c.fatal(pkgerr.WithStack(err))
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(pkgerr.WithStack(err))
	}
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return line, c.fatal(pkgerr.WithStack(err))
	}
	return line, nil
}

func (c *conn) Scan(item *Item, v interface{}) (err error) {
	c.ir.Reset(item.Value)
	if item.Flags&FlagGzip == FlagGzip {
		if err = c.gr.Reset(&c.ir); err != nil {
			return
		}
		if err = c.decode(&c.gr, item, v); err != nil {
			err = pkgerr.WithStack(err)
			return
		}
		err = c.gr.Close()
	} else {
		err = c.decode(&c.ir, item, v)
	}
	err = pkgerr.WithStack(err)
	return
}

func (c *conn) WithContext(ctx context.Context) Conn {
	// FIXME: implement WithContext
	return c
}

func (c *conn) encode(item *Item) (data []byte, err error) {
	if (item.Flags | _flagEncoding) == _flagEncoding {
		if item.Value == nil {
			return nil, ErrItem
		}
	} else if item.Object == nil {
		return nil, ErrItem
	}
	// encoding
	switch {
	case item.Flags&FlagGOB == FlagGOB:
		c.edb.Reset()
		if err = gob.NewEncoder(&c.edb).Encode(item.Object); err != nil {
			return
		}
		data = c.edb.Bytes()
	case item.Flags&FlagProtobuf == FlagProtobuf:
		c.edb.Reset()
		c.ped.SetBuf(c.edb.Bytes())
		pb, ok := item.Object.(proto.Message)
		if !ok {
			err = ErrItemObject
			return
		}
		if err = c.ped.Marshal(pb); err != nil {
			return
		}
		data = c.ped.Bytes()
	case item.Flags&FlagJSON == FlagJSON:
		c.edb.Reset()
		if err = c.je.Encode(item.Object); err != nil {
			return
		}
		data = c.edb.Bytes()
	default:
		data = item.Value
	}
	// compress
	if item.Flags&FlagGzip == FlagGzip {
		c.cb.Reset()
		c.gw.Reset(&c.cb)
		if _, err = c.gw.Write(data); err != nil {
			return
		}
		if err = c.gw.Close(); err != nil {
			return
		}
		data = c.cb.Bytes()
	}
	if len(data) > 8000000 {
		err = ErrValueSize
	}
	return
}

func (c *conn) decode(rd io.Reader, item *Item, v interface{}) (err error) {
	var data []byte
	switch {
	case item.Flags&FlagGOB == FlagGOB:
		err = gob.NewDecoder(rd).Decode(v)
	case item.Flags&FlagJSON == FlagJSON:
		c.jr.Reset(rd)
		err = c.jd.Decode(v)
	default:
		data = item.Value
		if item.Flags&FlagGzip == FlagGzip {
			c.edb.Reset()
			if _, err = io.Copy(&c.edb, rd); err != nil {
				return
			}
			data = c.edb.Bytes()
		}
		if item.Flags&FlagProtobuf == FlagProtobuf {
			m, ok := v.(proto.Message)
			if !ok {
				err = ErrItemObject
				return
			}
			c.ped.SetBuf(data)
			err = c.ped.Unmarshal(m)
		} else {
			switch v.(type) {
			case *[]byte:
				d := v.(*[]byte)
				*d = data
			case *string:
				d := v.(*string)
				*d = string(data)
			case interface{}:
				err = json.Unmarshal(data, v)
			}
		}
	}
	return
}

func legalKey(key string) bool {
	if len(key) > 250 || len(key) == 0 {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] <= ' ' || key[i] == 0x7f {
			return false
		}
	}
	return true
}
