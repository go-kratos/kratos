package model

import (
	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/libs/bytes"
	"go-common/app/service/main/broadcast/libs/encoding/binary"
	"go-common/app/service/main/broadcast/libs/websocket"
)

const (
	maxBodySizeV1 = int32(1 << 10)
	// size
	packSizeV1      = 4
	headerSizeV1    = 2
	verSizeV1       = 2
	operationSizeV1 = 4
	seqIDSizeV1     = 4
	heartbeatSizeV1 = 4
	rawHeaderSizeV1 = packSizeV1 + headerSizeV1 + verSizeV1 + operationSizeV1 + seqIDSizeV1
	maxPackSizeV1   = maxBodySizeV1 + int32(rawHeaderSizeV1)
	// offset
	packOffsetV1      = 0
	headerOffsetV1    = packOffsetV1 + packSizeV1
	verOffsetV1       = headerOffsetV1 + headerSizeV1
	operationOffsetV1 = verOffsetV1 + verSizeV1
	seqIDOffsetV1     = operationOffsetV1 + operationSizeV1
	heartbeatOffsetV1 = seqIDOffsetV1 + seqIDSizeV1
)

// WriteToV1 .
func (p *Proto) WriteToV1(b *bytes.Writer) {
	var (
		packLen = rawHeaderSizeV1 + int32(len(p.Body))
		buf     = b.Peek(rawHeaderSizeV1)
	)
	binary.BigEndian.PutInt32(buf[packOffsetV1:], packLen)
	binary.BigEndian.PutInt16(buf[headerOffsetV1:], int16(rawHeaderSizeV1))
	binary.BigEndian.PutInt16(buf[verOffsetV1:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[operationOffsetV1:], p.Operation)
	binary.BigEndian.PutInt32(buf[seqIDOffsetV1:], p.SeqId)
	if p.Body != nil {
		b.Write(p.Body)
	}
}

// ReadTCPV1 .
func (p *Proto) ReadTCPV1(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(rawHeaderSizeV1); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[packOffsetV1:headerOffsetV1])
	headerLen = binary.BigEndian.Int16(buf[headerOffsetV1:verOffsetV1])
	p.Ver = int32(binary.BigEndian.Int16(buf[verOffsetV1:operationOffsetV1]))
	p.Operation = binary.BigEndian.Int32(buf[operationOffsetV1:seqIDOffsetV1])
	p.SeqId = binary.BigEndian.Int32(buf[seqIDOffsetV1:])
	if packLen > maxPackSizeV1 {
		return ErrProtoPackLen
	}
	if headerLen != rawHeaderSizeV1 {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}

// WriteTCPV1 .
func (p *Proto) WriteTCPV1(wr *bufio.Writer) (err error) {
	var (
		buf     []byte
		packLen int32
	)
	if p.Operation == OpRaw {
		// write without buffer, job concact proto into raw buffer
		_, err = wr.WriteRaw(p.Body)
		return
	}
	packLen = rawHeaderSizeV1 + int32(len(p.Body))
	if buf, err = wr.Peek(rawHeaderSizeV1); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[packOffsetV1:], packLen)
	binary.BigEndian.PutInt16(buf[headerOffsetV1:], int16(rawHeaderSizeV1))
	binary.BigEndian.PutInt16(buf[verOffsetV1:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[operationOffsetV1:], p.Operation)
	binary.BigEndian.PutInt32(buf[seqIDOffsetV1:], p.SeqId)
	if p.Body != nil {
		_, err = wr.Write(p.Body)
	}
	return
}

// WriteTCPHeartV1 .
func (p *Proto) WriteTCPHeartV1(wr *bufio.Writer, online int32) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = rawHeaderSizeV1 + heartbeatSizeV1
	if buf, err = wr.Peek(packLen); err != nil {
		return
	}
	// header
	binary.BigEndian.PutInt32(buf[packOffsetV1:], int32(packLen))
	binary.BigEndian.PutInt16(buf[headerOffsetV1:], int16(rawHeaderSizeV1))
	binary.BigEndian.PutInt16(buf[verOffsetV1:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[operationOffsetV1:], p.Operation)
	binary.BigEndian.PutInt32(buf[seqIDOffsetV1:], p.SeqId)
	// body
	binary.BigEndian.PutInt32(buf[heartbeatOffsetV1:], online)
	return
}

// ReadWebsocketV1 .
func (p *Proto) ReadWebsocketV1(ws *websocket.Conn) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) < rawHeaderSizeV1 {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[packOffsetV1:headerOffsetV1])
	headerLen = binary.BigEndian.Int16(buf[headerOffsetV1:verOffsetV1])
	p.Ver = int32(binary.BigEndian.Int16(buf[verOffsetV1:operationOffsetV1]))
	p.Operation = binary.BigEndian.Int32(buf[operationOffsetV1:seqIDOffsetV1])
	p.SeqId = binary.BigEndian.Int32(buf[seqIDOffsetV1:])
	if packLen > maxPackSizeV1 {
		return ErrProtoPackLen
	}
	if headerLen != rawHeaderSizeV1 {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

// WriteWebsocketV1 .
func (p *Proto) WriteWebsocketV1(ws *websocket.Conn) (err error) {
	var (
		buf     []byte
		packLen int
	)
	if p.Operation == OpRaw {
		err = ws.WriteMessage(websocket.BinaryMessage, p.Body)
		return
	}
	packLen = rawHeaderSizeV1 + len(p.Body)
	if err = ws.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = ws.Peek(rawHeaderSizeV1); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[packOffsetV1:], int32(packLen))
	binary.BigEndian.PutInt16(buf[headerOffsetV1:], int16(rawHeaderSizeV1))
	binary.BigEndian.PutInt16(buf[verOffsetV1:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[operationOffsetV1:], p.Operation)
	binary.BigEndian.PutInt32(buf[seqIDOffsetV1:], p.SeqId)
	if p.Body != nil {
		err = ws.WriteBody(p.Body)
	}
	return
}

// WriteWebsocketHeartV1 .
func (p *Proto) WriteWebsocketHeartV1(wr *websocket.Conn, online int32) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = rawHeaderSizeV1 + heartbeatSizeV1
	// websocket header
	if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = wr.Peek(packLen); err != nil {
		return
	}
	// proto header
	binary.BigEndian.PutInt32(buf[packOffsetV1:], int32(packLen))
	binary.BigEndian.PutInt16(buf[headerOffsetV1:], int16(rawHeaderSizeV1))
	binary.BigEndian.PutInt16(buf[verOffsetV1:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[operationOffsetV1:], p.Operation)
	binary.BigEndian.PutInt32(buf[seqIDOffsetV1:], p.SeqId)
	// proto body
	binary.BigEndian.PutInt32(buf[heartbeatOffsetV1:], online)
	return
}
