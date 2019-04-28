package model

import (
	"errors"
	"sync"

	"go-common/library/time"
)

// ErrBusinessNotReady business is not ready.
var ErrBusinessNotReady = errors.New("error buiness is not ready")

// ArgBusiness rpc arg
type ArgBusiness struct {
	BusinessID int64
	Token      string
}

// Key rpc sharding key.
func (a *ArgBusiness) Key() int64 {
	return a.BusinessID
}

// Business business seq struct
type Business struct {
	ID            int64
	CurSeq        int64
	MaxSeq        int64
	Step          int64
	Perch         int64
	BenchTime     int64
	LastTimestamp int64
	Token         string
	CTime         time.Time
	MTime         time.Time
	Mutex         sync.Mutex
}

// SeqVersion seq-server version
type SeqVersion struct {
	IDC     int64
	SvrNum  int64
	SvrTime int64
}
