package model

import "go-common/app/interface/main/dm2/model/oplog"

// ArgOids arguments passed by rpc client
type ArgOids struct {
	Type   int32
	Plat   int8
	Oids   []int64
	RealIP string
}

// ArgEditDMState arguments passed by rpc client
type ArgEditDMState struct {
	Type         int32
	Oid          int64
	Mid          int64
	State        int32
	Dmids        []int64
	Source       oplog.Source
	OperatorType oplog.OperatorType
}

// ArgEditDMPool arguments passed by rpc client
type ArgEditDMPool struct {
	Type         int32
	Oid          int64
	Mid          int64
	Pool         int32
	Dmids        []int64
	Source       oplog.Source
	OperatorType oplog.OperatorType
}

// ArgEditDMAttr edit dm attr,bit: AttrXXX defined in model,value:AttrYes/AttrNo
type ArgEditDMAttr struct {
	Type         int32
	Oid          int64
	Mid          int64
	Bit          uint
	Value        int32
	Dmids        []int64
	Source       oplog.Source
	OperatorType oplog.OperatorType
}

// ArgAdvance arguments passed by rpc client
type ArgAdvance struct {
	Mid  int64
	Cid  int64
	Mode string
}

// ArgMid arguments passed by rpc client
type ArgMid struct {
	Mid int64
}

// ArgUpAdvance arguments passed by rpc client
type ArgUpAdvance struct {
	Mid int64
	ID  int64
}

// ArgAddUserFilters add user filters
type ArgAddUserFilters struct {
	Type    int8
	Mid     int64
	Filters []string
	Comment string
}

// ArgDelUserFilters delete user filter
type ArgDelUserFilters struct {
	Mid int64
	IDs []int64
}

// ArgAddUpFilters add up filters
type ArgAddUpFilters struct {
	Type    int8
	Mid     int64
	Filters []string
}

// ArgUpFilters up filters
type ArgUpFilters struct {
	Mid int64
}

// ArgEditUpFilters edit up filter
type ArgEditUpFilters struct {
	Type    int8
	Mid     int64
	Active  int8
	Filters []string
}

// ArgAddGlobalFilter add global filter
type ArgAddGlobalFilter struct {
	Type   int8
	Filter string
}

// ArgGlobalFilters get global filter
type ArgGlobalFilters struct {
}

// ArgDelGlobalFilters delete global filters
type ArgDelGlobalFilters struct {
	IDs []int64
}

// ArgBanUsers ban users
type ArgBanUsers struct {
	Mid   int64
	Oid   int64
	DMIDs []int64
}

// ArgCancelBanUsers cancel banned users
type ArgCancelBanUsers struct {
	Mid     int64
	Aid     int64
	Filters []string
}

// ArgCid arg cid
type ArgCid struct {
	Cid int64
}

// ArgSubtitleGet .
type ArgSubtitleGet struct {
	Aid  int64
	Oid  int64
	Type int32
}

// ArgSubtitleAllowSubmit .
type ArgSubtitleAllowSubmit struct {
	Aid         int64
	AllowSubmit bool
	Lan         string
}

// SubtitleSubjectReply .
type SubtitleSubjectReply struct {
	AllowSubmit bool
	Lan         string
	LanDoc      string
}

// ArgArchiveID .
type ArgArchiveID struct {
	Aid int64
}

// ArgMask .
type ArgMask struct {
	Cid  int64
	Plat int8
}
