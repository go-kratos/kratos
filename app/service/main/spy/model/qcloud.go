package model

// Reg Pro args.
const (
	AccountType = 4
)

// Tel level .
const (
	Nomal int8 = iota
	LevelOne
	LevelTwo
	LevelThree
	LevelFour
)

// QcloudRegProResp def.
type QcloudRegProResp struct {
	Code     int    `json:"code"`
	CodeDesc string `json:"codeDesc"`
	Message  string `json:"message"`
	Level    int8   `json:"level"`
}
