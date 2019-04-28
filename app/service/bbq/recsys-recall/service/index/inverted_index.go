package index

import (
	"encoding/binary"
	"errors"
)

const (
	_invertedIndexHeadSize  = 13
	_invertedIndexBlockSize = 8
	_magicNumber            = 0xffffffffdeadbeef
)

// InvertedIndexHead .
type InvertedIndexHead struct {
	Magic      uint64
	SourceType byte
	BodyLength uint32
}

// InvertedIndex .
type InvertedIndex struct {
	head *InvertedIndexHead
	Data []uint64
}

func (h *InvertedIndexHead) load(b []byte) error {
	if len(b) != _invertedIndexHeadSize {
		return errors.New("invalid head length")
	}
	h.Magic = binary.BigEndian.Uint64(b[:8])
	if h.Magic != _magicNumber {
		return errors.New("invalid head")
	}
	h.SourceType = b[8]
	h.BodyLength = binary.BigEndian.Uint32(b[9:])
	return nil
}

func (h *InvertedIndexHead) serialize() []byte {
	b := make([]byte, _invertedIndexHeadSize)

	b[0] = byte(h.Magic >> 56)
	b[1] = byte(h.Magic >> 48)
	b[2] = byte(h.Magic >> 40)
	b[3] = byte(h.Magic >> 32)
	b[4] = byte(h.Magic >> 24)
	b[5] = byte(h.Magic >> 16)
	b[6] = byte(h.Magic >> 8)
	b[7] = byte(h.Magic)

	b[8] = h.SourceType

	b[9] = byte(h.BodyLength >> 24)
	b[10] = byte(h.BodyLength >> 16)
	b[11] = byte(h.BodyLength >> 8)
	b[12] = byte(h.BodyLength)

	return b
}

// Load 反序列化
func (ii *InvertedIndex) Load(b []byte) error {
	if len(b) < _invertedIndexHeadSize {
		return errors.New("invalid data")
	}
	ii.head = new(InvertedIndexHead)
	if err := ii.head.load(b[:_invertedIndexHeadSize]); err != nil {
		return err
	}

	offset := _invertedIndexHeadSize
	blocks := int(ii.head.BodyLength)
	ii.Data = make([]uint64, blocks)

	for i := 0; i < blocks; i++ {
		cursor := offset + (i * _invertedIndexBlockSize)
		tmp := b[cursor : cursor+_invertedIndexBlockSize]
		if len(tmp) != _invertedIndexBlockSize {
			return errors.New("invalid item length")
		}
		ii.Data[i] = binary.BigEndian.Uint64(tmp)
	}

	return nil
}

// Serialize 序列化倒排索引
func (ii *InvertedIndex) Serialize() []byte {
	ii.head = &InvertedIndexHead{
		Magic:      _magicNumber,
		SourceType: byte(1),
		BodyLength: uint32(len(ii.Data)),
	}
	totalLen := _invertedIndexHeadSize + (len(ii.Data) * _invertedIndexBlockSize)
	b := make([]byte, totalLen)

	// head
	hb := ii.head.serialize()
	copy(b, hb)

	// body
	offset := _invertedIndexHeadSize
	for i := 0; i < len(ii.Data); i++ {
		b[offset+0] = byte(ii.Data[i] >> 56)
		b[offset+1] = byte(ii.Data[i] >> 48)
		b[offset+2] = byte(ii.Data[i] >> 40)
		b[offset+3] = byte(ii.Data[i] >> 32)
		b[offset+4] = byte(ii.Data[i] >> 24)
		b[offset+5] = byte(ii.Data[i] >> 16)
		b[offset+6] = byte(ii.Data[i] >> 8)
		b[offset+7] = byte(ii.Data[i])

		offset = offset + _invertedIndexBlockSize
	}

	return b
}
