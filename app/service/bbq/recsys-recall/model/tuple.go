package model

import (
	"encoding/binary"
	"math"
)

const (
	_tupleSize = 12
)

// Tuple .
type Tuple struct {
	Svid  uint64
	Score float32
}

// PriorityTuple .
type PriorityTuple struct {
	Tuple
	Tag      string
	Name     string
	Priority int32
}

// ToBytes .
func (t *Tuple) ToBytes() []byte {
	b := make([]byte, 12)

	b[0] = byte(t.Svid)
	b[1] = byte(t.Svid >> 8)
	b[2] = byte(t.Svid >> 16)
	b[3] = byte(t.Svid >> 24)
	b[4] = byte(t.Svid >> 32)
	b[5] = byte(t.Svid >> 40)
	b[6] = byte(t.Svid >> 48)
	b[7] = byte(t.Svid >> 56)
	// score
	score := math.Float32bits(t.Score)
	b[8] = byte(score)
	b[9] = byte(score >> 8)
	b[10] = byte(score >> 16)
	b[11] = byte(score >> 24)

	return b
}

// ParseTuple .
func ParseTuple(b []byte) *Tuple {
	svid := binary.LittleEndian.Uint64(b[:8])
	score := math.Float32frombits(binary.LittleEndian.Uint32(b[8:12]))
	return &Tuple{
		Svid:  svid,
		Score: score,
	}
}

// TupleSize .
func TupleSize() int {
	return _tupleSize
}
