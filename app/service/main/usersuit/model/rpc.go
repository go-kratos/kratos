package model

// ArgBuy buy
type ArgBuy struct {
	Mid int64
	Num int64
	IP  string
}

// ArgApply apply
type ArgApply struct {
	Mid    int64
	Code   string
	Cookie string
	IP     string
}

// ArgStat stat
type ArgStat struct {
	Mid int64
	IP  string
}

// ArgGenerate generator
type ArgGenerate struct {
	Mid       int64
	Num       int64
	ExpireDay int64
	IP        string
}

// ArgList list
type ArgList struct {
	Mid        int64
	Start, End int64
}

// ArgEquipment rpc pendant arg .
type ArgEquipment struct {
	Mid int64
	IP  string
}

// ArgEquipments rpc equipment arg .
type ArgEquipments struct {
	Mids []int64
	IP   string
}

// ArgEquip rpc equip arg.
type ArgEquip struct {
	Mid    int64
	Pid    int64
	Status int8
	IP     string
	Source int64
}

// ArgMid struct.
type ArgMid struct {
	Mid int64
}

// ArgMids struct.
type ArgMids struct {
	Mids []int64
}

// ArgMedalUserInfo struct.
type ArgMedalUserInfo struct {
	Mid    int64
	Cookie string
	IP     string
}

// ArgMedalInstall struct.
type ArgMedalInstall struct {
	Mid         int64 `form:"mid" validate:"gt=0,required"`
	Nid         int64 `form:"nid" validate:"gt=0,required"`
	IsActivated int8  `form:"isActivated"`
}

// ArgGrantByMids one pendant give to multiple users.
type ArgGrantByMids struct {
	BatchNo string
	Mids    []int64
	Pid     int64
	Expire  int64
}

// ArgGPMID .
type ArgGPMID struct {
	MID int64
	GID int64
}
