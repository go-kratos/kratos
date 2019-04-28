package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/juju/errors"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/packet"
)

type Conn struct {
	*packet.Conn

	user      string
	password  string
	db        string
	TLSConfig *tls.Config

	capability uint32

	status uint16

	charset string

	salt []byte

	connectionID uint32
}

func getNetProto(addr string) string {
	proto := "tcp"
	if strings.Contains(addr, "/") {
		proto = "unix"
	}
	return proto
}

// Connect to a MySQL server, addr can be ip:port, or a unix socket domain like /var/sock.
// Accepts a series of configuration functions as a variadic argument.
func Connect(addr string, user string, password string, dbName string, options ...func(*Conn)) (*Conn, error) {
	proto := getNetProto(addr)

	c := new(Conn)

	var err error
	conn, err := net.DialTimeout(proto, addr, 10*time.Second)
	if err != nil {
		return nil, errors.Trace(err)
	}

	c.Conn = packet.NewConn(conn)
	c.user = user
	c.password = password
	c.db = dbName

	//use default charset here, utf-8
	c.charset = DEFAULT_CHARSET

	// Apply configuration functions.
	for i := range options {
		options[i](c)
	}

	if err = c.handshake(); err != nil {
		return nil, errors.Trace(err)
	}

	return c, nil
}

func (c *Conn) handshake() error {
	var err error
	if err = c.readInitialHandshake(); err != nil {
		c.Close()
		return errors.Trace(err)
	}

	if err := c.writeAuthHandshake(); err != nil {
		c.Close()

		return errors.Trace(err)
	}

	if _, err := c.readOK(); err != nil {
		c.Close()
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) Ping() error {
	if err := c.writeCommand(COM_PING); err != nil {
		return errors.Trace(err)
	}

	if _, err := c.readOK(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) UseDB(dbName string) error {
	if c.db == dbName {
		return nil
	}

	if err := c.writeCommandStr(COM_INIT_DB, dbName); err != nil {
		return errors.Trace(err)
	}

	if _, err := c.readOK(); err != nil {
		return errors.Trace(err)
	}

	c.db = dbName
	return nil
}

func (c *Conn) GetDB() string {
	return c.db
}

func (c *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	if len(args) == 0 {
		return c.exec(command)
	} else {
		if s, err := c.Prepare(command); err != nil {
			return nil, errors.Trace(err)
		} else {
			var r *Result
			r, err = s.Execute(args...)
			s.Close()
			return r, err
		}
	}
}

func (c *Conn) Begin() error {
	_, err := c.exec("BEGIN")
	return errors.Trace(err)
}

func (c *Conn) Commit() error {
	_, err := c.exec("COMMIT")
	return errors.Trace(err)
}

func (c *Conn) Rollback() error {
	_, err := c.exec("ROLLBACK")
	return errors.Trace(err)
}

func (c *Conn) SetCharset(charset string) error {
	if c.charset == charset {
		return nil
	}

	if _, err := c.exec(fmt.Sprintf("SET NAMES %s", charset)); err != nil {
		return errors.Trace(err)
	} else {
		c.charset = charset
		return nil
	}
}

func (c *Conn) FieldList(table string, wildcard string) ([]*Field, error) {
	if err := c.writeCommandStrStr(COM_FIELD_LIST, table, wildcard); err != nil {
		return nil, errors.Trace(err)
	}

	data, err := c.ReadPacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	fs := make([]*Field, 0, 4)
	var f *Field
	if data[0] == ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		for {
			if data, err = c.ReadPacket(); err != nil {
				return nil, errors.Trace(err)
			}

			// EOF Packet
			if c.isEOFPacket(data) {
				return fs, nil
			}

			if f, err = FieldData(data).Parse(); err != nil {
				return nil, errors.Trace(err)
			}
			fs = append(fs, f)
		}
	}
	return nil, fmt.Errorf("field list error")
}

func (c *Conn) SetAutoCommit() error {
	if !c.IsAutoCommit() {
		if _, err := c.exec("SET AUTOCOMMIT = 1"); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (c *Conn) IsAutoCommit() bool {
	return c.status&SERVER_STATUS_AUTOCOMMIT > 0
}

func (c *Conn) IsInTransaction() bool {
	return c.status&SERVER_STATUS_IN_TRANS > 0
}

func (c *Conn) GetCharset() string {
	return c.charset
}

func (c *Conn) GetConnectionID() uint32 {
	return c.connectionID
}

func (c *Conn) HandleOKPacket(data []byte) *Result {
	r, _ := c.handleOKPacket(data)
	return r
}

func (c *Conn) HandleErrorPacket(data []byte) error {
	return c.handleErrorPacket(data)
}

func (c *Conn) ReadOKPacket() (*Result, error) {
	return c.readOK()
}

func (c *Conn) exec(query string) (*Result, error) {
	if err := c.writeCommandStr(COM_QUERY, query); err != nil {
		return nil, errors.Trace(err)
	}

	return c.readResult(false)
}
