package model

// MSGType .
type MSGType uint8

// const .
const (
	MSGTypeBlock MSGType = iota + 1
	MSGTypeBlockRemove
)

// MSG .
type MSG struct {
	Code    string
	Title   string
	Content string
}
