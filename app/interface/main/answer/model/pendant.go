package model

// const .
const (
	PendantNotGet = 0
	PendantGet    = 1
)

// Pendant .
type Pendant struct {
	Pid  int    `json:"pid"`
	Name string `json:"name"`
}

// ReqPendant .
type ReqPendant struct {
	HID int64 `form:"hid" validate:"required"`
	MID int64
}
