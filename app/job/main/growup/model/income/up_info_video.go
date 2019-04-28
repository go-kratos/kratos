package income

import (
	"go-common/library/time"
)

// Signed signed up
type Signed struct {
	MID          int64
	AccountState int
	SignedAt     time.Time
	QuitAt       time.Time
	IsDeleted    int
}
