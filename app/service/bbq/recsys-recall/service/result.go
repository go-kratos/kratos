package service

import (
	"encoding/binary"
	"errors"
	"go-common/app/service/bbq/recsys-recall/model"
)

// Result .
type Result struct {
	TotalHit    int32
	FilterCount int32
	FinalCount  int32
	Tuples      []*model.Tuple
}

// RecallResult .
type RecallResult struct {
	Result
	Tag      string
	Name     string
	Priority int32
}

// ToBytes .
func (rr *Result) ToBytes() *[]byte {
	totalLen := 12 + 12*len(rr.Tuples)
	b := make([]byte, totalLen)

	// total hit
	offset := 0
	b[offset] = byte(rr.TotalHit)
	b[offset+1] = byte(rr.TotalHit >> 8)
	b[offset+2] = byte(rr.TotalHit >> 16)
	b[offset+3] = byte(rr.TotalHit >> 24)

	// filter
	offset += 4
	b[offset] = byte(rr.FilterCount)
	b[offset+1] = byte(rr.FilterCount >> 8)
	b[offset+2] = byte(rr.FilterCount >> 16)
	b[offset+3] = byte(rr.FilterCount >> 24)

	// final
	offset += 4
	b[offset] = byte(rr.FinalCount)
	b[offset+1] = byte(rr.FinalCount >> 8)
	b[offset+2] = byte(rr.FinalCount >> 16)
	b[offset+3] = byte(rr.FinalCount >> 24)

	// tuples
	offset += 4
	for _, v := range rr.Tuples {
		tuple := v.ToBytes()
		for i := range tuple {
			b[offset+i] = tuple[i]
		}
		offset += len(tuple)
	}

	return &b
}

func parseResult(raw *[]byte) (*Result, error) {
	if len(*raw) <= 0 {
		return nil, errors.New("parse recall result invalid length")
	}

	// total hit
	offset := 0
	totalHit := binary.LittleEndian.Uint32((*raw)[offset : offset+4])

	// filter
	offset = offset + 4
	filterCount := binary.LittleEndian.Uint32((*raw)[offset : offset+4])

	// final
	offset = offset + 4
	finalCount := binary.LittleEndian.Uint32((*raw)[offset : offset+4])

	// tuple
	offset = offset + 4
	nblocks := (len((*raw)) - offset) / model.TupleSize()
	tuples := make([]*model.Tuple, nblocks)
	for i := 0; i < nblocks; i++ {
		tuples[i] = model.ParseTuple((*raw)[offset : offset+model.TupleSize()])
		offset += model.TupleSize()
	}
	return &Result{
		TotalHit:    int32(totalHit),
		FilterCount: int32(filterCount),
		FinalCount:  int32(finalCount),
		Tuples:      tuples,
	}, nil
}
