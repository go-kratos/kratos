package codec

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
)

// Decoder is
type Decoder struct {
	mutex sync.Mutex // each item must be received atomically
	r     io.Reader  // source of the data
	buf   decBuffer  // buffer for more efficient i/o from r
	err   error
}

// NewDecoder is
func NewDecoder(r io.Reader) *Decoder {
	dec := new(Decoder)
	// We use the ability to read bytes as a plausible surrogate for buffering.
	if _, ok := r.(io.ByteReader); !ok {
		r = bufio.NewReader(r)
	}
	dec.r = r
	return dec
}

func catchError(err *error) {
	if e := recover(); e != nil {
		pe, ok := e.(error)
		if !ok {
			panic(e)
		}
		if !strings.HasPrefix(pe.Error(), "proto") {
			panic(pe)
		}
		*err = pe
	}
}

// Decode is
func (dec *Decoder) Decode(e proto.Message, size int) error {
	return dec.DecodeValue(e, size)
}

// DecodeValue is
func (dec *Decoder) DecodeValue(v proto.Message, size int) error {
	// Make sure we're single-threaded through here.
	dec.mutex.Lock()
	defer dec.mutex.Unlock()

	dec.buf.Reset() // In case data lingers from previous invocation.
	dec.err = nil
	dec.decodeValue(v, size)
	return dec.err
}

func (dec *Decoder) decodeValue(value proto.Message, size int) {
	defer catchError(&dec.err)
	if err := dec.readChunk(size); err != nil {
		dec.err = err
		return
	}
	if err := proto.Unmarshal(dec.buf.Bytes(), value); err != nil {
		dec.err = err
		return
	}
	return
}

func (dec *Decoder) readChunk(chunkSize int) error {
	if dec.buf.Len() != 0 {
		// The buffer should always be empty now.
		panic("non-empty decoder buffer")
	}
	// Read the data
	dec.buf.Size(chunkSize)
	if _, err := io.ReadFull(dec.r, dec.buf.Bytes()); err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	return nil
}

type decBuffer struct {
	data   []byte
	offset int // Read offset.
}

func (d *decBuffer) Read(p []byte) (int, error) {
	n := copy(p, d.data[d.offset:])
	if n == 0 && len(p) != 0 {
		return 0, io.EOF
	}
	d.offset += n
	return n, nil
}

func (d *decBuffer) Drop(n int) {
	if n > d.Len() {
		panic("drop")
	}
	d.offset += n
}

// Size grows the buffer to exactly n bytes, so d.Bytes() will
// return a slice of length n. Existing data is first discarded.
func (d *decBuffer) Size(n int) {
	d.Reset()
	if cap(d.data) < n {
		d.data = make([]byte, n)
	} else {
		d.data = d.data[0:n]
	}
}

func (d *decBuffer) ReadByte() (byte, error) {
	if d.offset >= len(d.data) {
		return 0, io.EOF
	}
	c := d.data[d.offset]
	d.offset++
	return c, nil
}

func (d *decBuffer) Len() int {
	return len(d.data) - d.offset
}

func (d *decBuffer) Bytes() []byte {
	return d.data[d.offset:]
}

func (d *decBuffer) Reset() {
	d.data = d.data[0:0]
	d.offset = 0
}
