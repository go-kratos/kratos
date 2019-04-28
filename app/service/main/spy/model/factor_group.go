package model

import (
	"go-common/library/time"
)

// FactorGroup is.
type FactorGroup struct {
	ID    int64
	Name  string //风险因子组名
	CTime time.Time
}
