package model

// const .
const (
	IdentifyOk     = 0
	IdentifyNoInfo = 1

	PhoneOk    = 0
	Phone17x   = 1
	PhoneEmpty = 2
)

// IdentifyInfo .
type IdentifyInfo struct {
	Identify int8 `json:"'identify'"`
	Phone    int8 `json:"phone"`
}
