package client

import (
	"encoding/binary"

	"github.com/juju/errors"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/hack"
)

func (c *Conn) readUntilEOF() (err error) {
	var data []byte

	for {
		data, err = c.ReadPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			return
		}
	}
	return
}

func (c *Conn) isEOFPacket(data []byte) bool {
	return data[0] == EOF_HEADER && len(data) <= 5
}

func (c *Conn) handleOKPacket(data []byte) (*Result, error) {
	var n int
	var pos int = 1

	r := new(Result)

	r.AffectedRows, _, n = LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2
	}

	//new ok package will check CLIENT_SESSION_TRACK too, but I don't support it now.

	//skip info
	return r, nil
}

func (c *Conn) handleErrorPacket(data []byte) error {
	e := new(MyError)

	var pos int = 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		e.State = hack.String(data[pos : pos+5])
		pos += 5
	}

	e.Message = hack.String(data[pos:])

	return e
}

func (c *Conn) readOK() (*Result, error) {
	data, err := c.ReadPacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if data[0] == OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		return nil, errors.New("invalid ok packet")
	}
}

func (c *Conn) readResult(binary bool) (*Result, error) {
	data, err := c.ReadPacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if data[0] == OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else if data[0] == LocalInFile_HEADER {
		return nil, ErrMalformPacket
	}

	return c.readResultset(data, binary)
}

func (c *Conn) readResultset(data []byte, binary bool) (*Result, error) {
	result := &Result{
		Status:       0,
		InsertId:     0,
		AffectedRows: 0,

		Resultset: &Resultset{},
	}

	// column count
	count, _, n := LengthEncodedInt(data)

	if n-len(data) != 0 {
		return nil, ErrMalformPacket
	}

	result.Fields = make([]*Field, count)
	result.FieldNames = make(map[string]int, count)

	if err := c.readResultColumns(result); err != nil {
		return nil, errors.Trace(err)
	}

	if err := c.readResultRows(result, binary); err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}

func (c *Conn) readResultColumns(result *Result) (err error) {
	var i int = 0
	var data []byte

	for {
		data, err = c.ReadPacket()
		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			if i != len(result.Fields) {
				err = ErrMalformPacket
			}

			return
		}

		result.Fields[i], err = FieldData(data).Parse()
		if err != nil {
			return
		}

		result.FieldNames[hack.String(result.Fields[i].Name)] = i

		i++
	}
}

func (c *Conn) readResultRows(result *Result, isBinary bool) (err error) {
	var data []byte

	for {
		data, err = c.ReadPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		result.RowDatas = append(result.RowDatas, data)
	}

	result.Values = make([][]interface{}, len(result.RowDatas))

	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary)

		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}
