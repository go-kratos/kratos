package model

import (
	"errors"

	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/libs/bytes"
	"go-common/app/service/main/broadcast/libs/encoding/binary"
	"go-common/app/service/main/broadcast/libs/websocket"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize        = 4
	_headerSize      = 2
	_verSize         = 2
	_operationSize   = 4
	_seqIDSize       = 4
	_compressSize    = 1
	_contentTypeSize = 1
	_rawHeaderSize   = _packSize + _headerSize + _verSize + _operationSize + _seqIDSize + _compressSize + _contentTypeSize
	_maxPackSize     = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset        = 0
	_headerOffset      = _packOffset + _packSize
	_verOffset         = _headerOffset + _headerSize
	_operationOffset   = _verOffset + _verSize
	_seqIDOffset       = _operationOffset + _operationSize
	_compressOffset    = _seqIDOffset + _seqIDSize
	_contentTypeOffset = _compressOffset + _compressSize
)

var (
	emptyJSONBody = []byte("{}")

	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Operation: OpProtoReady}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Operation: OpProtoFinish}
)

// WriteTo write a proto to bytes writer.
func (p *Proto) WriteTo(b *bytes.Writer) {
	var (
		packLen = _rawHeaderSize + int32(len(p.Body))
		buf     = b.Peek(_rawHeaderSize)
	)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_operationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[_seqIDOffset:], p.SeqId)
	binary.BigEndian.PutInt8(buf[_compressOffset:], int8(p.Compress))
	binary.BigEndian.PutInt8(buf[_contentTypeOffset:], int8(p.ContentType))
	if p.Body != nil {
		b.Write(p.Body)
	}
}

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_operationOffset]))
	p.Operation = binary.BigEndian.Int32(buf[_operationOffset:_seqIDOffset])
	p.SeqId = binary.BigEndian.Int32(buf[_seqIDOffset:_compressOffset])
	p.Compress = int32(binary.BigEndian.Int8(buf[_compressOffset:_contentTypeOffset]))
	p.ContentType = int32(binary.BigEndian.Int8(buf[_contentTypeOffset:]))
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}

// WriteTCP write a proto to TCP writer.
func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	var (
		buf     []byte
		packLen int32
	)
	if p.Operation == OpRaw {
		// write without buffer, job concact proto into raw buffer
		_, err = wr.WriteRaw(p.Body)
		return
	}
	packLen = _rawHeaderSize + int32(len(p.Body))
	if buf, err = wr.Peek(_rawHeaderSize); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_operationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[_seqIDOffset:], p.SeqId)
	binary.BigEndian.PutInt8(buf[_compressOffset:], int8(p.Compress))
	binary.BigEndian.PutInt8(buf[_contentTypeOffset:], int8(p.ContentType))
	if p.Body != nil {
		_, err = wr.Write(p.Body)
	}
	return
}

// ReadWebsocket read a proto from websocket connection.
func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_operationOffset]))
	p.Operation = binary.BigEndian.Int32(buf[_operationOffset:_seqIDOffset])
	p.SeqId = binary.BigEndian.Int32(buf[_seqIDOffset:_compressOffset])
	p.Compress = int32(binary.BigEndian.Int8(buf[_compressOffset:_contentTypeOffset]))
	p.ContentType = int32(binary.BigEndian.Int8(buf[_contentTypeOffset:]))
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Proto) WriteWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf     []byte
		packLen int
	)
	// NOTE: 通过 OpRaw = 9 为ws批量消息处理
	//	if p.Operation == OpRaw {
	//		err = ws.WriteMessage(websocket.BinaryMessage, p.Body)
	//		return
	//	}
	packLen = _rawHeaderSize + len(p.Body)
	if err = ws.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = ws.Peek(_rawHeaderSize); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_operationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[_seqIDOffset:], p.SeqId)
	binary.BigEndian.PutInt8(buf[_compressOffset:], int8(p.Compress))
	binary.BigEndian.PutInt8(buf[_contentTypeOffset:], int8(p.ContentType))
	if p.Body != nil {
		err = ws.WriteBody(p.Body)
	}
	return
}

// WriteWebsocketHeart write a heartbeat proto to websocket connnection.
func (p *Proto) WriteWebsocketHeart(wr *websocket.Conn) (err error) {
	var (
		buf     []byte
		packLen int
	)
	if len(p.Body) == 0 {
		p.Body = emptyJSONBody
	}
	packLen = _rawHeaderSize + len(p.Body)
	// websocket header
	if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = wr.Peek(_rawHeaderSize); err != nil {
		return
	}
	// proto header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_operationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[_seqIDOffset:], p.SeqId)
	binary.BigEndian.PutInt8(buf[_compressOffset:], int8(p.Compress))
	binary.BigEndian.PutInt8(buf[_contentTypeOffset:], int8(p.ContentType))
	// proto body
	if p.Body != nil {
		err = wr.WriteBody(p.Body)
	}
	return
}

// WriteTCPHeart write a heartbeat proto to TCP writer.
func (p *Proto) WriteTCPHeart(wr *bufio.Writer) (err error) {
	var (
		buf     []byte
		packLen int32
	)
	if len(p.Body) == 0 {
		p.Body = emptyJSONBody
	}
	packLen = _rawHeaderSize + int32(len(p.Body))
	if buf, err = wr.Peek(_rawHeaderSize); err != nil {
		return
	}
	// header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_operationOffset:], p.Operation)
	binary.BigEndian.PutInt32(buf[_seqIDOffset:], p.SeqId)
	binary.BigEndian.PutInt8(buf[_compressOffset:], int8(p.Compress))
	binary.BigEndian.PutInt8(buf[_contentTypeOffset:], int8(p.ContentType))
	// body
	if p.Body != nil {
		_, err = wr.Write(p.Body)
	}
	return
}
