package memcache

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/gogo/protobuf/proto"
)

type reader struct {
	io.Reader
}

func (r *reader) Reset(rd io.Reader) {
	r.Reader = rd
}

const _encodeBuf = 4096 // 4kb

type encodeDecode struct {
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

func newEncodeDecoder() *encodeDecode {
	ed := &encodeDecode{}
	ed.jd = json.NewDecoder(&ed.jr)
	ed.je = json.NewEncoder(&ed.edb)
	ed.gw = gzip.NewWriter(&ed.cb)
	ed.edb.Grow(_encodeBuf)
	// NOTE reuse bytes.Buffer internal buf
	// DON'T concurrency call Scan
	ed.ped = proto.NewBuffer(ed.edb.Bytes())
	return ed
}

func (ed *encodeDecode) encode(item *Item) (data []byte, err error) {
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
		ed.edb.Reset()
		if err = gob.NewEncoder(&ed.edb).Encode(item.Object); err != nil {
			return
		}
		data = ed.edb.Bytes()
	case item.Flags&FlagProtobuf == FlagProtobuf:
		ed.edb.Reset()
		ed.ped.SetBuf(ed.edb.Bytes())
		pb, ok := item.Object.(proto.Message)
		if !ok {
			err = ErrItemObject
			return
		}
		if err = ed.ped.Marshal(pb); err != nil {
			return
		}
		data = ed.ped.Bytes()
	case item.Flags&FlagJSON == FlagJSON:
		ed.edb.Reset()
		if err = ed.je.Encode(item.Object); err != nil {
			return
		}
		data = ed.edb.Bytes()
	default:
		data = item.Value
	}
	// compress
	if item.Flags&FlagGzip == FlagGzip {
		ed.cb.Reset()
		ed.gw.Reset(&ed.cb)
		if _, err = ed.gw.Write(data); err != nil {
			return
		}
		if err = ed.gw.Close(); err != nil {
			return
		}
		data = ed.cb.Bytes()
	}
	if len(data) > 8000000 {
		err = ErrValueSize
	}
	return
}

func (ed *encodeDecode) decode(item *Item, v interface{}) (err error) {
	var (
		data []byte
		rd   io.Reader
	)
	ed.ir.Reset(item.Value)
	rd = &ed.ir
	if item.Flags&FlagGzip == FlagGzip {
		rd = &ed.gr
		if err = ed.gr.Reset(&ed.ir); err != nil {
			return
		}
		defer func() {
			if e := ed.gr.Close(); e != nil {
				err = e
			}
		}()
	}
	switch {
	case item.Flags&FlagGOB == FlagGOB:
		err = gob.NewDecoder(rd).Decode(v)
	case item.Flags&FlagJSON == FlagJSON:
		ed.jr.Reset(rd)
		err = ed.jd.Decode(v)
	default:
		data = item.Value
		if item.Flags&FlagGzip == FlagGzip {
			ed.edb.Reset()
			if _, err = io.Copy(&ed.edb, rd); err != nil {
				return
			}
			data = ed.edb.Bytes()
		}
		if item.Flags&FlagProtobuf == FlagProtobuf {
			m, ok := v.(proto.Message)
			if !ok {
				err = ErrItemObject
				return
			}
			ed.ped.SetBuf(data)
			err = ed.ped.Unmarshal(m)
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
