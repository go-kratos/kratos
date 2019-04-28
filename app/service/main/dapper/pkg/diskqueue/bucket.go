package diskqueue

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	_blockByte int32 = 512
	_lenByte   int32 = 2
	_dataByte        = _blockByte - _lenByte
)

var errBucketFull = errors.New("bucket is full or not enough")
var fullHeader = []byte{1, 254}
var nextHeader = []byte{1, 255}
var magicHeader = []byte{'D', 'Q'}

type memBucketPool struct {
	cap  int32
	pool sync.Pool
}

func newMemBucketPool(bucketByte int32) *memBucketPool {
	return &memBucketPool{
		pool: sync.Pool{New: func() interface{} {
			return make([]byte, bucketByte)
		}},
		cap: bucketByte / _blockByte,
	}
}

func (m *memBucketPool) new() *memBucket {
	data := m.pool.Get().([]byte)
	return &memBucket{data: data, cap: m.cap}
}

func (m *memBucketPool) free(bucket *memBucket) {
	m.pool.Put(bucket.data)
}

type memBucket struct {
	sync.Mutex
	cap     int32
	readAt  int32
	writeAt int32
	data    []byte
}

func (m *memBucket) push(p []byte) error {
	m.Lock()
	defer m.Unlock()
	length := int32(len(p))
	if length > _dataByte*(m.cap-m.writeAt) {
		return errBucketFull
	}
	// if p length < blockbyte write it direct
	if length < _dataByte {
		ds := m.writeAt * _blockByte
		binary.BigEndian.PutUint16(m.data[ds:], uint16(length))
		copy(m.data[ds+_lenByte:], p)
		m.writeAt++
		return nil
	}
	// loop write block
	blocks := length / _dataByte
	re := length % _dataByte
	var i int32
	for i = 0; i < blocks-1; i++ {
		ds := m.writeAt * _blockByte
		copy(m.data[ds:], nextHeader)
		ps := i * _dataByte
		copy(m.data[ds+_lenByte:], p[ps:ps+_dataByte])
		m.writeAt++
	}
	var nh []byte
	if re == 0 {
		nh = fullHeader
	} else {
		nh = nextHeader
	}
	ds := m.writeAt * _blockByte
	copy(m.data[ds:], nh)
	ps := (blocks - 1) * _dataByte
	copy(m.data[ds+_lenByte:], p[ps:ps+_dataByte])
	m.writeAt++
	if re != 0 {
		ds := m.writeAt * _blockByte
		binary.BigEndian.PutUint16(m.data[ds:], uint16(re))
		copy(m.data[ds+_lenByte:], p[blocks*_dataByte:])
		m.writeAt++
	}
	return nil
}

func (m *memBucket) pop() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	if m.readAt >= m.writeAt {
		return nil, io.EOF
	}
	ret := make([]byte, 0, _blockByte)
	for m.readAt < m.writeAt {
		ds := m.readAt * _blockByte
		m.readAt++
		l := int32(binary.BigEndian.Uint16(m.data[ds : ds+_lenByte]))
		if l <= _dataByte {
			ret = append(ret, m.data[ds+_lenByte:ds+_lenByte+l]...)
			break
		}
		ret = append(ret, m.data[ds+_lenByte:ds+_blockByte]...)
	}
	return ret, nil
}

func (m *memBucket) dump(w io.Writer) (int, error) {
	header := make([]byte, 10)
	copy(header, magicHeader)
	binary.BigEndian.PutUint32(header[2:6], uint32(m.readAt))
	binary.BigEndian.PutUint32(header[6:10], uint32(m.writeAt))
	n1, err := w.Write(header)
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(m.data[:m.writeAt*_blockByte])
	return n1 + n2, err
}

func newFileBucket(fpath string) (*fileBucket, error) {
	fp, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	header := make([]byte, 10)
	n, err := fp.Read(header)
	if err != nil {
		return nil, err
	}
	if n != 10 {
		return nil, fmt.Errorf("expect read 10 byte header get: %d", n)
	}
	if !bytes.Equal(header[:2], magicHeader) {
		return nil, fmt.Errorf("invalid magic %s", header[:2])
	}
	readAt := int32(binary.BigEndian.Uint32(header[2:6]))
	writeAt := int32(binary.BigEndian.Uint32(header[6:10]))
	if _, err = fp.Seek(int64(readAt*_blockByte), os.SEEK_CUR); err != nil {
		return nil, err
	}
	return &fileBucket{
		fp:      fp,
		readAt:  readAt,
		writeAt: writeAt,
		bufRd:   bufio.NewReader(fp),
	}, nil
}

type fileBucket struct {
	sync.Mutex
	fp      *os.File
	readAt  int32
	writeAt int32
	bufRd   *bufio.Reader
}

func (f *fileBucket) pop() ([]byte, error) {
	f.Lock()
	defer f.Unlock()
	if f.readAt >= f.writeAt {
		return nil, io.EOF
	}
	ret := make([]byte, 0, _blockByte)
	block := make([]byte, _blockByte)
	for f.readAt < f.writeAt {
		n, err := f.bufRd.Read(block)
		if err != nil {
			return nil, err
		}
		if int32(n) != _blockByte {
			return nil, fmt.Errorf("expect read %d byte data get %d", _blockByte, n)
		}
		l := int32(binary.BigEndian.Uint16(block[:2]))
		if l <= _dataByte {
			ret = append(ret, block[2:2+l]...)
			break
		}
		ret = append(ret, block[2:_blockByte]...)
	}
	return ret, nil
}

func (f *fileBucket) close() error {
	return f.fp.Close()
}
