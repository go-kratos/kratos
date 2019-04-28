package model

// IdentifyStatus .
type IdentifyStatus int8

// IdentifyStatus enum
const (
	IdentifyNotOK IdentifyStatus = iota
	IdentifyOK
	// ApiIdentifyOk identify ok
	APIIdentifyOk = 0
	// ApiIdentifyNoInfo no identify info
	APIIdentifyNoInfo = 1
)

// Identification .
type Identification struct {
	Identification IdentifyStatus `json:"identification"`
}

// IdentifyInfo identify info.
type IdentifyInfo struct {
	Identify int8 `json:"'identify'"`
	Phone    int8 `json:"phone"`
}

// IdentifyApply identify apply info.
type IdentifyApply struct {
	Identify int8 `json:"'identify'"`
	Applied  bool `json:"applied"`
}
