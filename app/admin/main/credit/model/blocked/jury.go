package blocked

import (
	"encoding/json"
	"go-common/library/log"
	"strconv"

	"go-common/app/admin/main/credit/model"
	xtime "go-common/library/time"
)

// const jury
const (
	// judge status.
	JudgeTypeUndeal  = int8(0) // 未裁决
	JudgeTypeViolate = int8(1) // 违规
	JudgeTypeLegal   = int8(2) // 未违
	// judge invalid_reason
	JuryBlocked = int8(1)
	JuryExpire  = int8(2)
	JuryAdmin   = int8(3)
	// jury black or white
	JuryNormal = int8(0)
	JuryBlack  = int8(1)
	JuryWhite  = int8(2)
	// JuryStatus
	JuryStatusOn   = int8(1)
	JuryStatusDown = int8(2)
)

// var jury
var (
	JuryerStyle = map[int8]string{
		JuryNormal: "正常",
		JuryBlack:  "黑名单",
		JuryWhite:  "白名单",
	}
	JuryerStatus = map[int8]string{
		JuryStatusOn:   "有效",
		JuryStatusDown: "失效",
	}
)

// Jury blocked_jury model.
type Jury struct {
	ID            int64      `gorm:"column:id" json:"id"`
	UID           int64      `gorm:"column:mid" json:"uid"`
	OPID          int64      `gorm:"column:oper_id" json:"oper_id"`
	UName         string     `gorm:"-" json:"uname"`
	Status        int8       `gorm:"column:status" json:"status"`
	StatusDesc    string     `gorm:"-" json:"status_desc"`
	Expired       xtime.Time `gorm:"column:expired" json:"expired"`
	EffectDay     xtime.Time `gorm:"-" json:"effect_day"`
	InvalidReason int8       `gorm:"column:invalid_reason" json:"invalid_reason"`
	VoteTotal     int        `gorm:"column:vote_total" json:"vote_total"`
	VoteRight     int        `gorm:"column:vote_right" json:"vote_right"`
	Total         int        `gorm:"column:total" json:"total"`
	Black         int8       `gorm:"column:black" json:"black"`
	VoteRadio     string     `gorm:"-" json:"vote_radio"`
	BlackDesc     string     `gorm:"-" json:"black_desc"`
	Remark        string     `gorm:"column:remark" json:"remark"`
	CTime         xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime         xtime.Time `gorm:"column:mtime" json:"mtime"`
	OPName        string     `gorm:"-" json:"oname"`
}

// WebHook is work flow webhook .
type WebHook struct {
	Verb  string `json:"verb"`
	Actor struct {
		AdminID int64 `json:"admin_id"`
	} `json:"actor"`
	Object *struct {
		CID   int64   `json:"cid"`
		CIDs  []int64 `json:"cids"`
		State int     `json:"state"`
	} `json:"object"`
	Target *struct {
		CID      int64 `json:"cid"`
		OID      int64 `json:"oid"`
		Business int   `json:"business"`
		Mid      int64 `json:"mid"`
		Tid      int   `json:"tid"`
		State    int   `json:"state"`
	} `json:"target"`
}

// TableName jury tablename
func (*Jury) TableName() string {
	return "blocked_jury"
}

// JuryList is info list.
type JuryList struct {
	Count int
	PN    int
	PS    int
	Order string
	Sort  string
	IDs   []int64
	List  []*Jury
}

// JuryDesc struct
type JuryDesc struct {
	UID        string `json:"uid"`
	UName      string `json:"uname"`
	StatusDesc string `json:"status_desc"`
	BlackDesc  string `json:"black_desc"`
	VoteTotal  string `json:"vote_total"`
	VoteRadio  string `json:"vote_radio"`
	Expired    string `json:"expired"`
	Remark     string `json:"remark"`
	EffectDay  string `json:"effect_day"`
	OPName     string `json:"oname"`
}

// DealJury export data.
func DealJury(jurys []*Jury) (data [][]string, err error) {
	var jurysDesc []*JuryDesc
	if len(jurys) < 0 {
		return
	}
	for _, v := range jurys {
		juryDesc := &JuryDesc{
			UID:        strconv.FormatInt(v.UID, 10),
			UName:      v.UName,
			StatusDesc: v.StatusDesc,
			BlackDesc:  v.BlackDesc,
			VoteTotal:  strconv.FormatInt(int64(v.VoteTotal), 10),
			VoteRadio:  v.VoteRadio,
			Expired:    v.Expired.Time().Format(model.TimeFormatSec),
			Remark:     v.Remark,
			EffectDay:  v.EffectDay.Time().Format(model.TimeFormatSec),
			OPName:     v.OPName,
		}
		jurysDesc = append(jurysDesc, juryDesc)
	}
	jurysMap, _ := json.Marshal(jurysDesc)
	var objmap []map[string]string
	if err = json.Unmarshal(jurysMap, &objmap); err != nil {
		log.Error("Unmarshal(%s) error(%v)", string(jurysMap), err)
		return
	}
	data = append(data, []string{"UID", "昵称", "状态", "类型", "投票数", "投准率", "失效时间", "备注", "生效时间", "操作人"})
	for _, v := range objmap {
		var fields []string
		fields = append(fields, v["uid"])
		fields = append(fields, v["uname"])
		fields = append(fields, v["status_desc"])
		fields = append(fields, v["black_desc"])
		fields = append(fields, v["vote_total"])
		fields = append(fields, v["vote_radio"])
		fields = append(fields, v["expired"])
		fields = append(fields, v["remark"])
		fields = append(fields, v["effect_day"])
		fields = append(fields, v["oname"])
		data = append(data, fields)
	}
	return
}
