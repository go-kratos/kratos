package model

import (
	"encoding/csv"
	"fmt"
	"time"

	dtmdl "go-common/app/interface/main/mcn/model/datamodel"
	xtime "go-common/library/time"
)

// Permits .
type Permits struct {
	BasePermission   uint8 `form:"base_permission" json:"base_permission" validate:"min=0,max=1"`     // 基础权限
	DataPermission   uint8 `form:"data_permission" json:"data_permission" validate:"min=0,max=1"`     // 数据权限
	RecPermission    uint8 `form:"rec_permission" json:"rec_permission" validate:"min=0,max=1"`       // 推荐权限
	DepartPermission uint8 `form:"depart_permission" json:"depart_permission" validate:"min=0,max=1"` // 起飞权限
}

// SetAttrPermitVal set struct from permission
func (p *Permits) SetAttrPermitVal(val uint32) {
	p.BasePermission = AttrVal(val, uint(AttrBasePermitBit))
	p.DataPermission = AttrVal(val, uint(AttrDataPermitBit))
	p.RecPermission = AttrVal(val, uint(AttrRecPermitBit))
	p.DepartPermission = AttrVal(val, uint(AttrDepartPermitBit))
}

// GetAttrPermitVal .
func (p *Permits) GetAttrPermitVal() (permission uint32) {
	permission = AttrSet(permission, p.BasePermission, uint(AttrBasePermitBit))
	permission = AttrSet(permission, p.DataPermission, uint(AttrDataPermitBit))
	permission = AttrSet(permission, p.RecPermission, uint(AttrRecPermitBit))
	permission = AttrSet(permission, p.DepartPermission, uint(AttrDepartPermitBit))
	return
}

// AttrSet set Permission.
func AttrSet(dest uint32, bitValue uint8, bit uint) (res uint32) {
	res = dest&(^(1 << bit)) | (uint32(bitValue) << bit)
	return
}

// AttrVal get Permission.
func AttrVal(v uint32, bit uint) uint8 {
	return uint8((v >> bit) & 1)
}

// MCNSignInfoReply .
type MCNSignInfoReply struct {
	SignID             int64               `json:"sign_id"`
	McnMid             int64               `json:"mcn_mid"`
	McnName            string              `json:"mcn_name"`
	CompanyName        string              `json:"company_name"`
	CompanyLicenseID   string              `json:"company_license_id"`
	CompanyLicenseLink string              `json:"company_license_link"`
	ContractLink       string              `json:"contract_link"`
	ContactName        string              `json:"contact_name"`
	ContactTitle       string              `json:"contact_title"`
	ContactPhone       string              `json:"contact_phone"`
	ContactIdcard      string              `json:"contact_idcard"`
	BeginDate          xtime.Time          `json:"begin_date"`
	EndDate            xtime.Time          `json:"end_date"`
	State              MCNSignState        `json:"state"`
	RejectTime         xtime.Time          `json:"reject_time"`
	RejectReason       string              `json:"reject_reason"`
	Ctime              xtime.Time          `json:"ctime"`
	Mtime              xtime.Time          `json:"mtime"`
	SignPayInfo        []*SignPayInfoReply `json:"sign_pay_info"`
	Permission         uint32              `json:"permission"`
	Permits            *Permits            `json:"permits"` // 权限集合
}

// AttrPermitVal get Permission all.
func (n *MCNSignInfoReply) AttrPermitVal() {
	n.Permits = &Permits{}
	n.Permits.SetAttrPermitVal(n.Permission)
}

// MCNSignListReply .
type MCNSignListReply struct {
	List []*MCNSignInfoReply `json:"result"`
	PageResult
}

// SignPayInfoReply  .
type SignPayInfoReply struct {
	SignPayID int64       `json:"sign_pay_id,omitempty"`
	McnMid    int64       `json:"mcn_mid"`
	SignID    int64       `json:"sign_id,omitempty"`
	State     MCNPayState `json:"state"`
	DueDate   xtime.Time  `json:"due_date"`
	PayValue  int64       `json:"pay_value"` // thousand bit
}

// MCNUPInfoReply .
type MCNUPInfoReply struct {
	SignUpID               int64      `json:"sign_up_id"`
	SignID                 int64      `json:"sign_id"`
	McnMid                 int64      `json:"mcn_mid"`
	UpMid                  int64      `json:"up_mid"`
	BeginDate              xtime.Time `json:"begin_date"`
	EndDate                xtime.Time `json:"end_date"`
	ContractLink           string     `json:"contract_link"`
	UpAuthLink             string     `json:"up_auth_link"`
	RejectTime             xtime.Time `json:"reject_time"`
	RejectReason           string     `json:"reject_reason"`
	State                  MCNUPState `json:"state"`
	StateChangeTime        xtime.Time `json:"state_change_time"`
	Ctime                  xtime.Time `json:"ctime"`
	Mtime                  xtime.Time `json:"mtime"`
	UpName                 string     `json:"up_name"`
	McnName                string     `json:"mcn_name"`
	ActiveTid              int16      `json:"active_tid"`
	TpName                 string     `json:"type_name"`
	FansCount              int64      `json:"fans_count"`
	FansCountActive        int64      `json:"fans_count_active"`
	FansIncreaseAccumulate int64      `json:"fans_increase_accumulate"`
	ArchiveCount           int64      `json:"archive_count"`
	PlayCount              int64      `json:"play_count"`
	UPType                 int8       `json:"up_type"`
	SiteLink               string     `json:"site_link"`
	ConfirmTime            xtime.Time `json:"confirm_time"`
	PubPrice               int64      `json:"publication_price"`
	Permission             uint32     `json:"permission"`
	Permits                *Permits   `json:"permits"` // 权限集合
}

// AttrPermitVal get Permission all.
func (n *MCNUPInfoReply) AttrPermitVal() {
	n.Permits = &Permits{}
	n.Permits.SetAttrPermitVal(n.Permission)
}

// MCNUPReviewListReply .
type MCNUPReviewListReply struct {
	List []*MCNUPInfoReply `json:"result"`
	PageResult
}

// UpBaseInfo .
type UpBaseInfo struct {
	Mid                    int64 `json:"mid"`
	FansCount              int64 `json:"fans_count"`
	ActiveTid              int16 `json:"active_tid"`
	ArticleCountAccumulate int64 `json:"article_count_accumulate"`
}

// UpPlayInfo .
type UpPlayInfo struct {
	Mid                 int64 `json:"mid"`
	ArticleCount        int64 `json:"article_count"`
	PlayCountAccumulate int64 `json:"play_count_accumulate"`
	PlayCountAverage    int64 `json:"play_count_average"`
}

// MCNListReply struct .
type MCNListReply struct {
	List []*MCNListOne `json:"result"`
	PageResult
}

// MCNListOne struct .
type MCNListOne struct {
	ID                        int64               `json:"id"`
	MCNMID                    int64               `json:"mcn_mid"`
	MCNName                   string              `json:"mcn_name"`
	UPCount                   int64               `json:"up_count"`
	FansCountAccumulate       int64               `json:"fans_count_accumulate"`
	FansCountOnlineAccumulate int64               `json:"fans_count_online_accumulate"`
	FansCountRealAccumulate   int64               `json:"fans_count_real_accumulate"`
	FansCountCheatAccumulate  int64               `json:"fans_count_cheat_accumulate"`
	GenerateDate              xtime.Time          `json:"generate_date"`
	BeginDate                 xtime.Time          `json:"begin_date"`
	EndDate                   xtime.Time          `json:"end_date"`
	State                     MCNSignState        `json:"state"`
	PayInfos                  []*SignPayInfoReply `json:"pay_infos"`
	Permission                uint32              `json:"permission"`
	Permits                   *Permits            `json:"permits"` // 权限集合
}

// AttrPermitVal get Permission all.
func (n *MCNListOne) AttrPermitVal() {
	n.Permits = &Permits{}
	n.Permits.SetAttrPermitVal(n.Permission)
}

// MCNInfoReply struct .
type MCNInfoReply struct {
	MCNSign
	UPCount                   int64 `json:"up_count"`
	ArchiveCountAccumulate    int64 `json:"archive_count_accumulate"`
	PlayCountAccumulate       int64 `json:"play_count_accumulate"`
	FansCountAccumulate       int64 `json:"fans_count_accumulate"`
	FansCountOnline           int64 `json:"fans_count_online"`
	FansCountReal             int64 `json:"fans_count_real"`
	FansCountCheat            int64 `json:"fans_count_cheat"`
	FansCountRealAccumulate   int64 `json:"fans_count_real_accumulate"`
	FansCountOnlineAccumulate int64 `json:"fans_count_online_accumulate"`
}

// MCNUPListReply struct .
type MCNUPListReply struct {
	List []*MCNUPInfoReply `json:"result"`
	PageResult
}

// MCNCheatReply struct .
type MCNCheatReply struct {
	SignID                          int64  `json:"sign_id"`
	MCNMID                          int64  `json:"mcn_mid"`
	MCNName                         string `json:"mcn_name"`
	UpMID                           int64  `json:"up_mid"`
	UpName                          string `json:"up_name"`
	FansCountCheatAccumulate        int64  `json:"fans_count_cheat_accumulate"`
	FansCountCheatIncreaseDay       int64  `json:"fans_count_cheat_increase_day"`
	FansCountReal                   int64  `json:"fans_count_real"`
	FansCountCheatCleanedAccumulate int64  `json:"fans_count_cheat_cleaned_accumulate"`
}

// MCNCheatListReply struct.
type MCNCheatListReply struct {
	List []*MCNCheatReply `json:"result"`
	PageResult
}

// MCNCheatUPReply struct .
type MCNCheatUPReply struct {
	GenerateDate                    xtime.Time `json:"generate_date"`
	FansCountCheatIncreaseDay       int64      `json:"fans_count_cheat_increase_day"`
	MCNMID                          int64      `json:"mcn_mid"`
	MCNName                         string     `json:"mcn_name"`
	SignID                          int64      `json:"sign_id"`
	FansCountCheatAccumulate        int64      `json:"fans_count_cheat_accumulate"`
	FansCountCheatCleanedAccumulate int64      `json:"fans_count_cheat_cleaned_accumulate"`
	FansCountReal                   int64      `json:"fans_count_real"`
}

// MCNCheatUPListReply struct .
type MCNCheatUPListReply struct {
	List []*MCNCheatUPReply `json:"result"`
	PageResult
}

// MCNImportUPInfoReply struct .
type MCNImportUPInfoReply struct {
	ID                   int64  `json:"id"`
	MCNMID               int64  `json:"mcn_mid"`
	SignID               int64  `json:"sign_id"`
	UpMID                int64  `json:"up_mid"`
	UpName               string `json:"up_name"`
	StandardFansDate     int64  `json:"standard_fans_date"`
	StandardArchiveCount int64  `json:"standard_archive_count"`
	StandardFansCount    int64  `json:"standard_fans_count"`
	IsReward             int8   `json:"is_reward"`
	JoinTime             int32  `json:"join_time"`
}

// MCNIncreaseReply struct .
type MCNIncreaseReply struct {
	ID                        int64      `json:"id"`
	SignID                    int64      `json:"sign_id"`
	DataType                  int8       `json:"data_type"`
	ActiveTID                 int64      `json:"active_tid"`
	GenerateDate              xtime.Time `json:"generate_date"`
	UPCount                   int64      `json:"up_count"`
	FansCountOnlineAccumulate int64      `json:"fans_count_online_accumulate"`
	FansCountRealAccumulate   int64      `json:"fans_count_real_accumulate"`
	FansCountCheatAccumulate  int64      `json:"fans_count_cheat_accumulate"`
	FansCountIncreaseDay      int64      `json:"fans_count_increase_day"`
	ArchiveCountAccumulate    int64      `json:"archive_count_accumulate"`
	ArchiveCountDay           int64      `json:"archive_count_day"`
	PlayCountAccumulate       int64      `json:"play_count_accumulate"`
	PlayCountIncreaseDay      int64      `json:"play_count_increase_day"`
	FansCountAccumulate       int64      `json:"fans_count_accumulate"`
}

// MCNIncreaseListReply struct .
type MCNIncreaseListReply struct {
	List []*MCNIncreaseReply `json:"result"`
	PageResult
}

//GetFileName get file name
func (q *MCNListReply) GetFileName() string {
	return fmt.Sprintf("%s_%s.csv", "MCN列表", time.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (q *MCNListReply) ToCsv(writer *csv.Writer) {
	var title = []string{
		"ID",
		"MCN_ID",
		"MCN_昵称",
		"签约UP主数",
		"累计粉丝数",
		"累计线上涨粉数",
		"累计实际粉丝数",
		"累计作弊粉丝数",
		"签约周期",
		"付款周期",
		"账号状态",
	}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.List {
		var record []string
		var payString string
		if len(v.PayInfos) > 0 {
			for _, pv := range v.PayInfos {
				payString += fmt.Sprintf("%s-%d-%s ", pv.DueDate.Time().Format(TimeFormatDay), pv.PayValue/1000, pv.State.String())
			}
		}
		record = append(record,
			intFormat(v.ID),
			intFormat(v.MCNMID),
			v.MCNName,
			intFormat(v.UPCount),
			intFormat(v.FansCountAccumulate),
			intFormat(v.FansCountOnlineAccumulate),
			intFormat(v.FansCountRealAccumulate),
			intFormat(v.FansCountCheatAccumulate),
			fmt.Sprintf("%s-%s", v.BeginDate.Time().Format(TimeFormatDay), v.EndDate.Time().Format(TimeFormatDay)),
			payString,
			v.State.String(),
		)
		writer.Write(record)
	}
}

//GetFileName get file name
func (q *MCNUPListReply) GetFileName() string {
	return fmt.Sprintf("%s_%s.csv", "MCN UP主列表", time.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (q *MCNUPListReply) ToCsv(writer *csv.Writer) {
	var title = []string{
		"ID",
		"UP主UID",
		"UP主昵称",
		"粉丝总量",
		"活跃粉丝量",
		"粉数增长量",
		"稿件量",
		"播放量",
		"分区",
		"账号状态",
		"签约周期",
	}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.List {
		var record []string
		record = append(record,
			intFormat(v.SignUpID),
			intFormat(v.UpMid),
			v.UpName,
			intFormat(v.FansCount),
			intFormat(v.FansCountActive),
			intFormat(v.FansIncreaseAccumulate),
			intFormat(v.ArchiveCount),
			intFormat(v.PlayCount),
			v.TpName,
			v.State.String(),
			fmt.Sprintf("%s-%s", v.BeginDate.Time().Format(TimeFormatDay), v.EndDate.Time().Format(TimeFormatDay)),
		)

		writer.Write(record)
	}
}

// McnUpRecommendPool .
type McnUpRecommendPool struct {
	ID                     int64                `json:"id"`
	UpMid                  int64                `json:"up_mid"`
	UpName                 string               `json:"up_name"`
	FansCount              int64                `json:"fans_count"`
	FansCountIncreaseMonth int64                `json:"fans_count_increase_month"`
	ArchiveCount           int64                `json:"archive_count"`
	PlayCountAccumulate    int64                `json:"play_count_accumulate"`
	PlayCountAverage       int64                `json:"play_count_average"`
	ActiveTid              int16                `json:"active_tid"`
	TpName                 string               `json:"type_name"`
	LastArchiveTime        xtime.Time           `json:"last_archive_time"`
	State                  MCNUPRecommendState  `json:"state"`
	Source                 MCNUPRecommendSource `json:"source"`
	GenerateTime           xtime.Time           `json:"generate_time"`
	Ctime                  xtime.Time           `json:"ctime"`
	Mtime                  xtime.Time           `json:"mtime"`
}

// McnUpRecommendListReply struct .
type McnUpRecommendListReply struct {
	List []*McnUpRecommendPool `json:"result"`
	PageResult
}

//GetFileName get file name
func (list *McnUpRecommendListReply) GetFileName() string {
	return fmt.Sprintf("%s_%s.csv", "MCN推荐池列表", time.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (list *McnUpRecommendListReply) ToCsv(writer *csv.Writer) {
	var title = []string{
		"UP主UID",
		"up主昵称",
		"粉丝量",
		"本月粉丝增长量",
		"累积播放量",
		"稿均播放量",
		"分区",
		"最近投稿时间",
		"来源",
		"推荐池状态",
		"数据更新时间",
	}
	writer.Write(title)
	if list == nil {
		return
	}
	for _, v := range list.List {
		var record []string
		record = append(record,
			intFormat(v.UpMid),
			v.UpName,
			intFormat(v.FansCount),
			intFormat(v.FansCountIncreaseMonth),
			intFormat(v.PlayCountAccumulate),
			intFormat(v.PlayCountAverage),
			v.TpName,
			v.LastArchiveTime.Time().Format(TimeFormatSec),
			v.Source.String(),
			v.State.String(),
			v.GenerateTime.Time().Format(TimeFormatSec),
		)
		writer.Write(record)
	}
}

// McnGetRankUpFansReply reply
type McnGetRankUpFansReply struct {
	Result   []*RankArchiveLikeInfo `json:"result"` // 按顺序进行排名
	TypeList []*TidnameInfo         `json:"type_list"`
}

// GetFileName get file name
func (list *McnGetRankUpFansReply) GetFileName() string {
	return fmt.Sprintf("%s_%s.csv", "top稿件列表", time.Now().Format(dateTimeFmt))
}

// ToCsv to buffer
func (list *McnGetRankUpFansReply) ToCsv(writer *csv.Writer) {
	var title = []string{
		"稿件ID",
		"稿件标题",
		"UP主UID",
		"UP主昵称",
		"新增点赞数",
		"累积点赞数",
		"新增播放数",
		"累积播放数",
		"分区",
		"上传日期",
	}
	writer.Write(title)
	if list == nil {
		return
	}
	for _, v := range list.Result {
		var record []string
		record = append(record,
			intFormat(v.ArchiveID),
			v.ArchiveTitle,
			intFormat(v.Author.Mid),
			v.Author.Name,
			intFormat(v.LikesIncrease),
			intFormat(v.LikesAccumulate),
			intFormat(v.PlayIncrease),
			intFormat(v.PlayAccumulate),
			v.TidName,
			v.Ctime.Time().Format(TimeFormatSec),
		)
		writer.Write(record)
	}
}

// McnGetMcnFansReply reply 粉丝分析.
type McnGetMcnFansReply struct {
	FansOverview *dtmdl.DmConMcnFansD        `json:"fans_overview"` // 粉丝概况
	FansSex      *dtmdl.DmConMcnFansSexW     `json:"fans_sex"`      // 粉丝性别
	FansAge      *dtmdl.DmConMcnFansAgeW     `json:"fans_age"`      // 粉丝年龄
	FansPlayWay  *dtmdl.DmConMcnFansPlayWayW `json:"fans_play_way"` // 粉丝观看途径
	FansArea     []*dtmdl.DmConMcnFansAreaW  `json:"fans_area"`     // 粉丝地区分布
	FansType     []*dtmdl.DmConMcnFansTypeW  `json:"fans_type"`     // 粉丝倾向分布
	FansTag      []*dtmdl.DmConMcnFansTagW   `json:"fans_tag"`      // 粉丝标签地图分布
}

// McnUpPermissionApply .
type McnUpPermissionApply struct {
	ID            int64                `json:"id"`
	McnMid        int64                `json:"mcn_mid"`
	UpMid         int64                `json:"up_mid"`
	McnName       string               `json:"mcn_name"`
	UpName        string               `json:"up_name"`
	SignID        int64                `json:"sign_id"`
	FansCount     int64                `json:"fans_count"`
	UpAuthLink    string               `json:"up_auth_link"`
	ActiveTID     int16                `json:"active_tid"`
	TypeName      string               `json:"type_name"`
	RejectReason  string               `json:"reject_reason"`
	RejectTime    xtime.Time           `json:"reject_time"`
	State         MCNUPPermissionState `json:"state"`
	Ctime         xtime.Time           `json:"ctime"`
	Mtime         xtime.Time           `json:"mtime"`
	AdminID       int64                `json:"admin_id"`
	AdminName     string               `json:"admin_name"`
	OldPermits    *Permits             `json:"old_permits"`
	NewPermits    *Permits             `json:"new_permits"`
	NewPermission uint32               `json:"-"`
	OldPermission uint32               `json:"-"`
}

// AttrPermitVal get Permission all.
func (n *McnUpPermissionApply) AttrPermitVal() {
	n.OldPermits, n.NewPermits = &Permits{}, &Permits{}
	n.OldPermits.SetAttrPermitVal(n.OldPermission)
	n.NewPermits.SetAttrPermitVal(n.NewPermission)
}

// McnUpPermitApplyListReply struct .
type McnUpPermitApplyListReply struct {
	List []*McnUpPermissionApply `json:"result"`
	PageResult
}
