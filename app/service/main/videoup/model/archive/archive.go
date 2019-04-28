package archive

import (
	"go-common/library/time"
)

// Archive is archive model.
type Archive struct {
	Aid          int64         `json:"aid"`
	Mid          int64         `json:"mid"`
	TypeID       int16         `json:"tid"`
	HumanRank    int           `json:"-"`
	Title        string        `json:"title"`
	Author       string        `json:"-"`
	Cover        string        `json:"cover"`
	RejectReason string        `json:"reject_reason"`
	Tag          string        `json:"tag"`
	Duration     int64         `json:"duration"`
	Copyright    int8          `json:"copyright"`
	Desc         string        `json:"desc"`
	MissionID    int64         `json:"mission_id"`
	Round        int8          `json:"-"`
	Forward      int64         `json:"-"`
	Attribute    int32         `json:"attribute"`
	Access       int16         `json:"-"`
	State        int8          `json:"state"`
	Source       string        `json:"source"`
	NoReprint    int32         `json:"no_reprint"`
	UGCPay       int32         `json:"ugcpay"`
	OrderID      int64         `json:"order_id"`
	UpFrom       int8          `json:"up_from"`
	Dynamic      string        `json:"dynamic"`
	DescFormatID int64         `json:"desc_format_id"`
	Porder       *Porder       `json:"porder"`
	Staffs       []*StaffApply `json:"staffs"`
	POI          *PoiObj       `json:"poi_object"`
	Vote         *Vote         `json:"vote"`
	DTime        time.Time     `json:"dtime"`
	PTime        time.Time     `json:"ptime"`
	CTime        time.Time     `json:"ctime"`
	MTime        time.Time     `json:"mtime"`
}

// AttrSet set attribute.
func (a *Archive) AttrSet(v int32, bit uint) {
	a.Attribute = a.Attribute&(^(1 << bit)) | (v << bit)
}

// AttrVal get attribute.
func (a *Archive) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// NotAllowUp check archive is or not allow update state.
func (a *Archive) NotAllowUp() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidLock || a.State == StateForbidPolice
}

// SimpleArchive str
type SimpleArchive struct {
	Aid    int64    `json:"aid"`
	Title  string   `json:"title"`
	Mid    int64    `json:"mid"`
	Videos []*Video `json:"videos,omitempty"`
}

// Addit str
type Addit struct {
	Aid           int64  `json:"aid"`
	MissionID     int64  `json:"mission_id"`
	UpFrom        int8   `json:"up_from"`
	FromIP        int64  `json:"from_ip"`
	IPv6          []byte `json:"ipv6"`
	Source        string `json:"source"`
	OrderID       int64  `json:"order_id"`
	RecheckReason string `json:"recheck_reason"`
	RedirectURL   string `json:"redirect_url"`
	FlowID        int64  `json:"flow_id"`
	Advertiser    string `json:"advertiser"`
	FlowRemark    string `json:"flow_remark"`
	DescFormatID  int64  `json:"desc_format_id"`
	Desc          string `json:"desc"`
	Dynamic       string `json:"dynamic"`
}

// Delay str
type Delay struct {
	Aid   int64
	State int8
	DTime time.Time
}

// Type info
type Type struct {
	ID   int16  `json:"id"`
	PID  int16  `json:"pid"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

// Alert str
type Alert struct {
	Key   string
	Value int64
	Limit int64
}

// Up str
type Up struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	GroupName string    `json:"group_name" `
	GroupTag  string    `json:"group_tag"`
	Mid       int64     `json:"mid"`
	Note      string    `json:"note"`
	CTime     time.Time `json:"ctime"`
}
