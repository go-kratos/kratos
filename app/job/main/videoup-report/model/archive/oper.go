package archive

import "fmt"

const (
	//OperTypeNoChannel oper type
	OperTypeNoChannel = int8(1)
	// OperStyleOne 操作展示类型1：[%s]从[%v]设为[%v]
	OperStyleOne = int8(1)
	// OperStyleTwo 操作展示类型2：[%s]%v:%v
	OperStyleTwo = int8(2)
)

var (
	//FlowOperType type
	FlowOperType = map[int64]int8{
		FLowGroupIDChannel: OperTypeNoChannel,
	}
	_operType = map[int8]string{
		OperTypeNoChannel: "频道禁止",
	}
)

// Operformat oper format.
func Operformat(tagID int8, old, new interface{}, style int8) (cont string) {
	var template string
	switch style {
	case OperStyleOne:
		template = "[%s]从[%v]设为[%v]"
	case OperStyleTwo:
		template = "[%s]%v:%v"
	}
	cont = fmt.Sprintf(template, _operType[tagID], old, new)
	return
}

//VideoOper 视频审核记录结构
type VideoOper struct {
	AID       int64
	UID       int64
	VID       int64
	Status    int
	Content   string
	Attribute int32
	LastID    int64
	Remark    string
}
