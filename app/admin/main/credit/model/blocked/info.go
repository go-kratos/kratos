package blocked

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/admin/main/credit/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// const info
const (
	// info publish_status
	StatusClose = int8(0) // 案件关闭状态
	StatusOpen  = int8(1) // 案件公开状态
	// info block_type
	PunishBlock = 0 // 系统封禁
	PunishJury  = 1 // 风纪仲裁
	// punish type.
	PunishTypeMoral        = int8(1) // 节操
	PunishTypeBlockTime    = int8(2) // 封禁
	PunishTypeBlockForever = int8(3) // 永久封禁
	// info status.
	BlockStateOpen  = int8(0) // 未解禁
	BlockStateClose = int8(1) // 已解禁
)

// var info
var (
	PStatusDesc = map[int8]string{
		StatusClose: "不公开",
		StatusOpen:  "公开",
	}
	BTypeDesc = map[int8]string{
		PunishBlock: "系统封禁",
		PunishJury:  "风纪仲裁",
	}
)

// Info is blocked_info model.
type Info struct {
	ID                  int64      `gorm:"column:id" json:"id"`
	UID                 int64      `gorm:"column:uid" json:"uid"`
	UName               string     `gorm:"column:uname" json:"uname"`
	Status              int8       `gorm:"column:status" json:"status"`
	OriginTitle         string     `gorm:"column:origin_title" json:"origin_title"`
	OriginURL           string     `gorm:"column:origin_url" json:"origin_url"`
	OriginContent       string     `gorm:"column:origin_content" json:"origin_content"`
	OriginContentModify string     `gorm:"column:origin_content_modify" json:"origin_content_modify"`
	OriginType          int8       `gorm:"column:origin_type" json:"origin_type"`
	BlockedDays         int        `gorm:"column:blocked_days" json:"blocked_days"`
	BlockedForever      int8       `gorm:"column:blocked_forever" json:"blocked_forever"`
	BlockedType         int8       `gorm:"column:blocked_type" json:"blocked_type"`
	BlockedRemark       string     `gorm:"column:blocked_remark" json:"blocked_remark"`
	CaseID              int64      `gorm:"column:case_id" json:"case_id"`
	MoralNum            int        `gorm:"column:moral_num" json:"moral_num"`
	ReasonType          int8       `gorm:"column:reason_type" json:"reason_type"`
	PublishStatus       int8       `gorm:"column:publish_status" json:"publish_status"`
	PunishType          int8       `gorm:"column:punish_type" json:"punish_type"`
	PunishTime          xtime.Time `gorm:"column:punish_time" json:"punish_time"`
	PublishTime         xtime.Time `gorm:"column:publish_time" json:"publish_time"`
	OperID              int64      `gorm:"column:oper_id" json:"oper_id"`
	CTime               xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime               xtime.Time `gorm:"column:mtime" json:"mtime"`
	PublishStatusDesc   string     `gorm:"-" json:"publish_status_desc"`
	OriginTypeDesc      string     `gorm:"-" json:"origin_type_desc"`
	BlockedTypeDesc     string     `gorm:"-" json:"blocked_type_desc"`
	BlockedDaysDesc     string     `gorm:"-" json:"blocked_days_desc"`
	ReasonTypeDesc      string     `gorm:"-" json:"reason_type_desc"`
	OPName              string     `gorm:"-" json:"oname"`
	OOPName             string     `gorm:"column:operator_name" json:"-"`
}

// InfoList is info list.
type InfoList struct {
	IDs  []int64
	List []*Info
}

// InfoDesc is Info_desc model.
type InfoDesc struct {
	ID                string `json:"id"`
	PunishTime        string `json:"punish_time"`
	OriginTypeDesc    string `json:"origin_type_desc"`
	ReasonTypeDesc    string `json:"reason_type_desc"`
	PublishStatusDesc string `json:"publish_status_desc"`
	BlockedTypeDesc   string `json:"blocked_type_desc"`
	OriginContent     string `json:"origin_content"`
	BlockedDaysDesc   string `json:"blocked_days_desc"`
	UName             string `json:"uname"`
	UID               string `json:"uid"`
	OPName            string `json:"oname"`
}

// TableName Info tablename
func (*Info) TableName() string {
	return "blocked_info"
}

// DealInfo deal with info data.
func DealInfo(infos []*Info) (data [][]string, err error) {
	var infoDescs []*InfoDesc
	for _, v := range infos {
		infoDesc := &InfoDesc{
			ID:                strconv.FormatInt(v.ID, 10),
			PunishTime:        v.PunishTime.Time().Format(model.TimeFormatSec),
			OriginTypeDesc:    v.OriginTypeDesc,
			PublishStatusDesc: v.PublishStatusDesc,
			BlockedTypeDesc:   v.BlockedTypeDesc,
			ReasonTypeDesc:    v.ReasonTypeDesc,
			OriginContent:     v.OriginContent,
			BlockedDaysDesc:   v.BlockedDaysDesc,
			UName:             v.UName,
			UID:               strconv.FormatInt(v.UID, 10),
			OPName:            v.OPName,
		}
		infoDescs = append(infoDescs, infoDesc)
	}
	infoMap, _ := json.Marshal(infoDescs)
	var objmap []map[string]string
	if err = json.Unmarshal(infoMap, &objmap); err != nil {
		log.Error("Unmarshal(%s) error(%v)", string(infoMap), err)
		return
	}
	data = append(data, []string{"ID", "惩罚时间", "类型", "状态", "封禁类型", "理由类型", "原文概要", "处罚结果", "用户", "用户ID", "操作者"})
	for _, v := range objmap {
		var fields []string
		fields = append(fields, v["id"])
		fields = append(fields, v["punish_time"])
		fields = append(fields, v["origin_type_desc"])
		fields = append(fields, v["publish_status_desc"])
		fields = append(fields, v["blocked_type_desc"])
		fields = append(fields, v["reason_type_desc"])
		fields = append(fields, v["origin_content"])
		fields = append(fields, v["blocked_days_desc"])
		fields = append(fields, v["uname"])
		fields = append(fields, v["uid"])
		fields = append(fields, v["oname"])
		data = append(data, fields)
	}
	return
}

// BDaysDesc is blocked_days_desc.
func BDaysDesc(bDays, moralNum int, pType, bForever int8) string {
	switch {
	case pType == PunishTypeMoral:
		return fmt.Sprintf(blockedDesc[BlockMoralNum], moralNum)
	case bDays == BlockForever && bForever == OnBlockedForever:
		return blockedDesc[BlockForever]
	case pType == PunishTypeBlockTime:
		if bDays == BlockThree || bDays == BlockSeven || bDays == BlockFifteen {
			return blockedDesc[bDays]
		}
		return strconv.Itoa(bDays) + blockedDesc[BlockCustom]
	}
	return ""
}
