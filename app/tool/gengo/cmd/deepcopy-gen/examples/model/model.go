package model

import (
	"go-common/library/time"
)

// MemberBase is
type MemberBase struct {
	Mid      int64
	Name     string
	Sex      int64
	Face     string
	Sign     string
	Rank     int64
	Birthday time.Time
}

// Names is
type Names struct {
	Names map[int64]string
}

// Mids is
type Mids struct {
	Mids []int64
}
