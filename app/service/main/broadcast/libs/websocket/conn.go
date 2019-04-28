package websocket

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"go-common/app/service/main/broadcast/libs/bufio"
)

const (
	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	finBit  = 1 << 7
	rsv1Bit = 1 << 6
	rsv2Bit = 1 << 5
	rsv3Bit = 1 << 4
	opBit   = 0x0f

	// Frame header byte 1 bits from Section 5.2 of RFC 6455
	maskBit = 1 << 7
	lenBit  = 0x7f

	continuationFrame        = 0
	continuationFrameMaxRead = 100
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

var (
	// ErrMessageClose close control message
	ErrMessageClose = errors.New("close control message")
	// ErrMessageMaxRead continuation frrame max read
	ErrMessageMaxRead = errors.New("continuation frame max read")
)

// Conn represents a WebSocket connection.
type Conn struct {
	rwc io.ReadWriteCloser
	r   *bufio.Reader
	w   *bufio.Writer
}

// new connection
func newConn(rwc io.ReadWriteCloser, r *bufio.Reader, w *bufio.Writer) *Conn {
	return &Conn{rwc: rwc, r: r, w: w}
}

// WriteMessage write a message by type.
func (c *Conn) WriteMessage(msgType int, msg []byte) (err error) {
	if err = c.WriteHeader(msgType, len(msg)); err != nil {
		return
	}
	err = c.WriteBody(msg)
	return
}

// WriteHeader write header frame.
func (c *Conn) WriteHeader(msgType int, length int) (err error) {
	var h []byte
	if h, err = c.w.Peek(2); err != nil {
		return
	}
	// 1.First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	h[0] = 0
	h[0] |= finBit | byte(msgType)
	// 2.Second byte. Mask/Payload len(7bits)
	h[1] = 0
	switch {
	case length <= 125:
		// 7 bits
		h[1] |= byte(length)
	case length < 65536:
		// 16 bits
		h[1] |= 126
		if h, err = c.w.Peek(2); err != nil {
			return
		}
		binary.BigEndian.PutUint16(h, uint16(length))
	default:
		// 64 bits
		h[1] |= 127
		if h, err = c.w.Peek(8); err != nil {
			return
		}
		binary.BigEndian.PutUint64(h, uint64(length))
	}
	return
}

// WriteBody write a message body.
func (c *Conn) WriteBody(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.w.Write(b)
	}
	return
}

// Peek write peek.
func (c *Conn) Peek(n int) ([]byte, error) {
	return c.w.Peek(n)
}

// Flush flush writer buffer
func (c *Conn) Flush() error {
	return c.w.Flush()
}

// ReadMessage read a message.
func (c *Conn) ReadMessage() (op int, payload []byte, err error) {
	var (
		fin         bool
		finOp, n    int
		partPayload []byte
	)
	for {
		// read frame
		if fin, op, partPayload, err = c.readFrame(); err != nil {
			return
		}
		switch op {
		case BinaryMessage, TextMessage, continuationFrame:
			if fin && len(payload) == 0 {
				return op, partPayload, nil
			}
			// continuation frame
			payload = append(payload, partPayload...)
			if op != continuationFrame {
				finOp = op
			}
			// final frame
			if fin {
				op = finOp
				return
			}
		case PingMessage:
			// handler ping
			if err = c.WriteMessage(PongMessage, partPayload); err != nil {
				return
			}
		case PongMessage:
			// handler pong
		case CloseMessage:
			// handler close
			err = ErrMessageClose
			return
		default:
			err = fmt.Errorf("unknown control message, fin=%t, op=%d", fin, op)
			return
		}
		if n > continuationFrameMaxRead {
			err = ErrMessageMaxRead
			return
		}
		n++
	}
}

func (c *Conn) readFrame() (fin bool, op int, payload []byte, err error) {
	var (
		b          byte
		p          []byte
		mask       bool
		maskKey    []byte
		payloadLen int64
	)
	// 1.First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}
	// final frame
	fin = (b & finBit) != 0
	// rsv MUST be 0
	if rsv := b & (rsv1Bit | rsv2Bit | rsv3Bit); rsv != 0 {
		return false, 0, nil, fmt.Errorf("unexpected reserved bits rsv1=%d, rsv2=%d, rsv3=%d", b&rsv1Bit, b&rsv2Bit, b&rsv3Bit)
	}
	// op code
	op = int(b & opBit)
	// 2.Second byte. Mask/Payload len(7bits)
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}
	// is mask payload
	mask = (b & maskBit) != 0
	// payload length
	switch b & lenBit {
	case 126:
		// 16 bits
		if p, err = c.r.Pop(2); err != nil {
			return
		}
		payloadLen = int64(binary.BigEndian.Uint16(p))
	case 127:
		// 64 bits
		if p, err = c.r.Pop(8); err != nil {
			return
		}
		payloadLen = int64(binary.BigEndian.Uint64(p))
	default:
		// 7 bits
		payloadLen = int64(b & lenBit)
	}
	// read mask key
	if mask {
		maskKey, err = c.r.Pop(4)
		if err != nil {
			return
		}
	}
	// read payload
	if payloadLen > 0 {
		if payload, err = c.r.Pop(int(payloadLen)); err != nil {
			return
		}
		if mask {
			maskBytes(maskKey, 0, payload)
		}
	}
	return
}

// Close close the connection.
func (c *Conn) Close() error {
	return c.rwc.Close()
}

func maskBytes(key []byte, pos int, b []byte) int {
	for i := range b {
		b[i] ^= key[pos&3]
		pos++
	}
	return pos & 3
}
