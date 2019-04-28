package model

// const msg
const (
	MsgTypeCustom = int8(1)
)

// SysMsg msg struct
type SysMsg struct {
	IsMsg    bool
	Type     int8
	MID      int64
	Title    string
	Content  string
	RemoteIP string
}

// MsgInfo get msg info
func MsgInfo(msg *SysMsg) (title, content, ip string) {
	switch msg.Type {
	case MsgTypeCustom:
		title = msg.Title
		content = msg.Content
		ip = msg.RemoteIP
	}
	return
}
