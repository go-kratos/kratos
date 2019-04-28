package mcnmodel

import (
	"sort"

	mcnadminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/datamodel"
	arcgrpc "go-common/app/service/main/archive/api"

	"go-common/library/time"
)

// McnGetStateReply .
type McnGetStateReply struct {
	State        int8   `json:"state"`
	RejectReason string `json:"reject_reason"`
}

// CommonReply reply with nothing, base struct
type CommonReply struct{}

//Sorter can sort
type Sorter interface {
	Sort()
}

// McnGetBindReply .
type McnGetBindReply struct {
	BindID        int64                 `json:"bind_id" gorm:"column:bind_id"`
	CompanyName   string                `json:"company_name" gorm:"column:company_name"`
	McnMid        uint32                `json:"mcn_mid" gorm:"column:mcn_mid"`
	McnName       string                `json:"mcn_name"`
	UpAuthLink    string                `json:"up_auth_link" gorm:"column:up_auth_link"`
	OldPerm       uint32                `json:"-" gorm:"column:old_permission"`
	NewPerm       uint32                `json:"-" gorm:"column:new_permission"`
	OldPermission mcnadminmodel.Permits `json:"old_permission"`
	NewPermission mcnadminmodel.Permits `json:"new_permission"`
}

//Finish call this before send
func (s *McnGetBindReply) Finish() {
	s.OldPermission.SetAttrPermitVal(s.OldPerm)
	s.NewPermission.SetAttrPermitVal(s.NewPerm)
}

// McnGetDataSummaryReply .
type McnGetDataSummaryReply struct {
	datamodel.McnStatisticBaseInfo2
	IsNull                 bool      `json:"is_null"` // 缓存用来标记
	UpCount                int64     `json:"up_count"`
	UpCountDiff            int64     `json:"-"`
	FansCountAccumulate    int64     `json:"fans_count_accumulate"`
	FansCountDiff          int64     `json:"fans_count_diff"`
	PlayCountAccumulate    int64     `json:"play_count_accumulate"`
	PlayCountDiff          int64     `json:"play_count_diff"`
	ArchiveCountAccumulate int64     `json:"archive_count_accumulate"`
	ArchiveCountDiff       int64     `json:"archive_count_diff"`
	GenerateDate           time.Time `json:"generate_date"`
}

// CopyFrom .
func (m *McnGetDataSummaryReply) CopyFrom(v *McnDataSummary) {
	if v == nil {
		return
	}
	m.UpCount = int64(v.UpCount)
	m.FansCountAccumulate = int64(v.FansCountAccumulate)
	m.PlayCountAccumulate = int64(v.PlayCountAccumulate)
	m.ArchiveCountAccumulate = int64(v.ArchiveCountAccumulate)
	m.GenerateDate = v.GenerateDate
}

//CopyFromDmConMcnArchiveD .
func (m *McnGetDataSummaryReply) CopyFromDmConMcnArchiveD(v *datamodel.DmConMcnArchiveD) {
	if v == nil {
		return
	}
	m.UpCount = v.UpAll
	// m.UpCountDiff 不再使用了。
	m.ArchiveCountAccumulate = v.ArchiveAll
	m.ArchiveCountDiff = v.ArchiveInc
	m.PlayCountAccumulate = v.PlayAll
	m.PlayCountDiff = v.PlayInc
	m.FansCountAccumulate = v.FansAll
	m.FansCountDiff = v.FansInc
	m.McnStatisticBaseInfo2 = v.McnStatisticBaseInfo2
	m.GenerateDate = v.LogDate.Time()
}

// CalcDiff .
func (m *McnGetDataSummaryReply) CalcDiff(lastDay *McnDataSummary) {
	if lastDay == nil {
		return
	}
	m.UpCountDiff = m.UpCount - int64(lastDay.UpCount)
	m.FansCountDiff = m.FansCountAccumulate - int64(lastDay.FansCountAccumulate)
	m.PlayCountDiff = m.PlayCountAccumulate - int64(lastDay.PlayCountAccumulate)
	m.ArchiveCountDiff = m.ArchiveCountAccumulate - int64(lastDay.ArchiveCountAccumulate)
}

// McnUpDataInfo mcn data
type McnUpDataInfo struct {
	McnDataUp
	ActiveTid        int16                 `json:"active_tid"`
	TidName          string                `json:"tid_name"`
	BeginDate        time.Time             `json:"begin_date" gorm:"column:begin_date"`
	EndDate          time.Time             `json:"end_date" gorm:"column:end_date"`
	State            int8                  `json:"state" gorm:"column:state"`
	Name             string                `json:"name"`
	Permission       uint32                `gorm:"column:permission" json:"-"`
	PublicationPrice int64                 `gorm:"column:publication_price" json:"publication_price"` // 单位：1/1000 元
	Permits          mcnadminmodel.Permits `json:"permits"`
}

// HideData if state is not right, hide the data
func (m *McnUpDataInfo) HideData(hideDate bool) {
	m.FansIncreaseMonth = 0
	m.PlayCount = 0
	m.ArchiveCount = 0
	m.FansIncreaseAccumulate = 0
	m.FansCount = 0
	m.FansCountActive = 0
	if hideDate {
		m.BeginDate = 0
		m.EndDate = 0
	}
}

//Finish call finish before send out
func (m *McnUpDataInfo) Finish() {
	m.Permits.SetAttrPermitVal(m.Permission)
}

// McnBindUpApplyReply .
type McnBindUpApplyReply struct {
	BindID int64 `json:"bind_id"`
}

// McnGetUpListReply list
type McnGetUpListReply struct {
	Result []*McnUpDataInfo `json:"result"`
	model.PageResult
}

//Finish call finish before send out
func (m *McnGetUpListReply) Finish() {
	for _, v := range m.Result {
		v.Finish()
	}
}

// McnExistReply exist
type McnExistReply struct {
	Exist int `json:"exist"`
}

// McnGetAccountReply reply
type McnGetAccountReply struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
}

// McnGetMcnOldInfoReply req
type McnGetMcnOldInfoReply struct {
	CompanyName      string `json:"company_name"`
	CompanyLicenseID string `json:"company_license_id"`
	ContactName      string `json:"contact_name"`
	ContactTitle     string `json:"contact_title"`
	ContactIdcard    string `json:"contact_idcard"`
	ContactPhone     string `json:"contact_phone"`
}

// Copy .
func (m *McnGetMcnOldInfoReply) Copy(v *McnSign) {
	if v == nil {
		return
	}
	m.CompanyName = v.CompanyName
	m.CompanyLicenseID = v.CompanyLicenseID
	m.ContactName = v.ContactName
	m.ContactTitle = v.ContactTitle
	m.ContactIdcard = v.ContactIdcard
	m.ContactPhone = v.ContactPhone
}

// RankDataInterface 用来取排行用的数据
type RankDataInterface interface {
	GetTid() int16
	GetDataType() DataType
	GetValue() int64
}

// RankDataBase 基本排行信息
type RankDataBase struct {
	Tid      int16    `json:"tid"`
	DataType DataType `json:"data_type"`
}

// GetTid .
func (r *RankDataBase) GetTid() int16 {
	return r.Tid
}

// GetDataType .
func (r *RankDataBase) GetDataType() DataType {
	return r.DataType
}

// RankUpFansInfo reply info
type RankUpFansInfo struct {
	RankDataBase
	Mid            int64  `json:"mid"`
	UpFaceLink     string `json:"up_face_link"`
	Name           string `json:"name"`
	FansIncrease   int64  `json:"fans_increase"`
	FansAccumulate int64  `json:"fans_accumulate"`
	TidName        string `json:"tid_name"`
}

// GetValue .
func (r *RankUpFansInfo) GetValue() int64 {
	return r.FansIncrease
}

// Copy copy from rank up .
func (r *RankUpFansInfo) Copy(v *McnRankUpFan) {
	if v == nil {
		return
	}
	r.DataType = v.DataType
	r.Tid = v.ActiveTid
	r.Mid = v.UpMid
	r.FansIncrease = v.Value1
	r.FansAccumulate = v.Value2
}

// McnGetRankUpFansReply reply
type McnGetRankUpFansReply struct {
	Result   []RankDataInterface `json:"result"` // 按顺序进行排名
	TypeList []*TidnameInfo      `json:"type_list"`
}

// RankArchiveLikeInfo archive like rank info
type RankArchiveLikeInfo struct {
	RankDataBase
	ArchiveID       int64          `json:"archive_id"` // 稿件ID
	ArchiveTitle    string         `json:"archive_title"`
	Pic             string         `json:"pic"` // 封面
	TidName         string         `json:"tid_name"`
	LikesIncrease   int64          `json:"likes_increase"`
	LikesAccumulate int64          `json:"likes_accumulate"`
	PlayIncrease    int64          `json:"play_increase"`
	Ctime           time.Time      `json:"ctime"`
	Author          arcgrpc.Author `json:"author"` // up主信息
	Stat            arcgrpc.Stat   `json:"stat"`   // 统计信息
}

// GetValue .
func (r *RankArchiveLikeInfo) GetValue() int64 {
	return r.LikesIncrease
}

// CopyFromDB .
func (r *RankArchiveLikeInfo) CopyFromDB(v *McnRankArchiveLike) {
	if v == nil {
		return
	}

	r.ArchiveID = v.ArchiveID
	r.LikesIncrease = v.LikeCount
	r.Tid = v.Tid
	r.DataType = v.DataType
	r.PlayIncrease = v.PlayIncr
}

// CopyFromArchive copy from archive info from archive service
func (r *RankArchiveLikeInfo) CopyFromArchive(v *arcgrpc.Arc) {
	if v == nil {
		return
	}
	r.ArchiveTitle = v.Title
	r.Pic = v.Pic
	r.Ctime = v.Ctime
	r.TidName = v.TypeName
	r.Author = v.Author
	r.Stat = v.Stat
	r.Tid = int16(v.TypeID)
	r.LikesAccumulate = int64(v.Stat.Like)
}

// TidnameInfo tid name
type TidnameInfo struct {
	Tid  int16  `json:"tid"`
	Name string `json:"name"`
}

// McnGetRecommendPoolInfo recomend info
type McnGetRecommendPoolInfo struct {
	UpMid                  int64  `json:"up_mid"`
	FansCount              int64  `json:"fans_count"`
	FansCountIncreaseMonth int64  `json:"fans_count_increase_month"`
	ArchiveCount           int64  `json:"archive_count"`
	ActiveTid              int16  `json:"active_tid"`
	UpName                 string `json:"up_name"`
	TidName                string `json:"tid_name"`
}

// Copy copy from db
func (m *McnGetRecommendPoolInfo) Copy(v *McnUpRecommendPool) {
	if v == nil {
		return
	}
	m.UpMid = v.UpMid
	m.FansCount = v.FansCount
	m.FansCountIncreaseMonth = v.FansCountIncreaseMonth
	m.ArchiveCount = v.ArchiveCount
	m.ActiveTid = v.ActiveTid
}

// McnGetRecommendPoolReply result
type McnGetRecommendPoolReply struct {
	model.PageResult
	Result []*McnGetRecommendPoolInfo `json:"result"`
}

// McnGetRecommendPoolTidListReply result
type McnGetRecommendPoolTidListReply struct {
	Result []*TidnameInfo `json:"result"`
}

// --------------------------- 三期需求

// McnGetIndexIncReply 播放/弹幕/评论/分享/硬币/收藏/点赞数每日增量
type McnGetIndexIncReply struct {
	Result []*datamodel.DmConMcnIndexIncD `json:"result"`
}

//McnGetIndexSourceReply 播放/弹幕/评论/分享/硬币/收藏/点赞来源分区
type McnGetIndexSourceReply struct {
	Result []*datamodel.DmConMcnIndexSourceD `json:"result"`
}

//Sort sort
func (s *McnGetIndexSourceReply) Sort() {
	sort.Slice(s.Result, func(i, j int) bool {
		return s.Result[i].Value > s.Result[j].Value
	})
}

//McnGetPlaySourceReply #mcn稿件播放来源占比
type McnGetPlaySourceReply struct {
	datamodel.DmConMcnPlaySourceD
}

// McnGetMcnFansReply mcn粉丝基本数据
type McnGetMcnFansReply struct {
	datamodel.DmConMcnFansD
}

//McnGetMcnFansIncReply mcn粉丝按天增量
type McnGetMcnFansIncReply struct {
	Result []*datamodel.DmConMcnFansIncD `json:"result"`
}

//McnGetMcnFansDecReply mcn粉丝按天取关数
type McnGetMcnFansDecReply struct {
	Result []*datamodel.DmConMcnFansDecD `json:"result"`
}

//McnGetMcnFansAttentionWayReply mcn粉丝关注渠道
type McnGetMcnFansAttentionWayReply struct {
	datamodel.DmConMcnFansAttentionWayD
}

// McnGetBaseFansAttrReply # mcn 性别占比 + 观众年龄 + 观看途径
type McnGetBaseFansAttrReply struct {
	FansSex     *datamodel.DmConMcnFansSexW     `json:"fans_sex"`
	FansAge     *datamodel.DmConMcnFansAgeW     `json:"fans_age"`
	FansPlayWay *datamodel.DmConMcnFansPlayWayW `json:"fans_play_way"`
}

// McnGetFansAreaReply 游客/粉丝地区分布
type McnGetFansAreaReply struct {
	Result []*datamodel.DmConMcnFansAreaW `json:"result"`
}

//Sort sort.
func (s *McnGetFansAreaReply) Sort() {
	sort.Slice(s.Result, func(i, j int) bool {
		return s.Result[i].User > s.Result[j].User
	})
}

// McnGetFansTypeReply 游客/粉丝倾向分布
type McnGetFansTypeReply struct {
	Result []*datamodel.DmConMcnFansTypeW `json:"result"`
}

//Sort sort.
func (s *McnGetFansTypeReply) Sort() {
	sort.Slice(s.Result, func(i, j int) bool {
		return s.Result[i].Play > s.Result[j].Play
	})
}

// McnGetFansTagReply 游客/粉丝标签地图分布
type McnGetFansTagReply struct {
	Result []*datamodel.DmConMcnFansTagW `json:"result"`
}

// ----------------------------- 4期需求

//McnChangePermitReply .
type McnChangePermitReply = McnBindUpApplyReply

//McnGetUpPermitReply 4.2
type McnGetUpPermitReply struct {
	Old          *mcnadminmodel.Permits `json:"old"`
	New          *mcnadminmodel.Permits `json:"new"`
	ContractLink string                 `json:"contract_link"`
}

// McnPublicationPriceChangeReply 4.4
type McnPublicationPriceChangeReply = CommonReply

// McnBaseInfoReply .
type McnBaseInfoReply struct {
	ID                 int64                  `json:"id"`
	McnMid             int64                  `json:"mcn_mid"`
	CompanyName        string                 `json:"company_name"`
	CompanyLicenseID   string                 `json:"company_license_id"`
	CompanyLicenseLink string                 `json:"company_license_link"`
	ContractLink       string                 `json:"contract_link"`
	ContactName        string                 `json:"contact_name"`
	ContactTitle       string                 `json:"contact_title"`
	ContactIdcard      string                 `json:"contact_idcard"`
	ContactPhone       string                 `json:"contact_phone"`
	BeginDate          time.Time              `json:"begin_date"`
	EndDate            time.Time              `json:"end_date"`
	RejectReason       string                 `json:"reject_reason"`
	RejectTime         time.Time              `json:"reject_time"`
	State              model.MCNSignState     `json:"state"`
	Permission         uint32                 `json:"permission"`
	Ctime              time.Time              `json:"ctime"`
	Mtime              time.Time              `json:"mtime"`
	SignPayInfo        []*SignPayInfoReply    `json:"sign_pay_info"`
	Permits            *mcnadminmodel.Permits `json:"permits"` // 权限集合
}

// SignPayInfoReply  .
type SignPayInfoReply struct {
	State    mcnadminmodel.MCNPayState `json:"-"`
	DueDate  time.Time                 `json:"due_date"`
	PayValue int64                     `json:"pay_value"` // thousand bit
}

// AttrPermitVal get Permission all.
func (m *McnBaseInfoReply) AttrPermitVal() {
	m.Permits = &mcnadminmodel.Permits{}
	m.Permits.SetAttrPermitVal(m.Permission)
}

// CopyFromMcnInfo .
func (m *McnBaseInfoReply) CopyFromMcnInfo(v *McnSign) {
	if v == nil {
		return
	}
	m.ID = v.ID
	m.McnMid = v.McnMid
	m.CompanyName = v.CompanyName
	m.CompanyLicenseID = v.CompanyLicenseID
	m.CompanyLicenseLink = v.CompanyLicenseLink
	m.ContractLink = v.ContractLink
	m.ContactName = v.ContactName
	m.ContactTitle = v.ContactTitle
	m.ContactIdcard = v.ContactIdcard
	m.ContactPhone = v.ContactPhone
	m.BeginDate = v.BeginDate
	m.EndDate = v.EndDate
	m.RejectReason = v.RejectReason
	m.RejectTime = v.RejectTime
	m.State = v.State
	m.Permission = v.Permission
	m.Ctime = v.Ctime
	m.Mtime = v.Mtime
	m.AttrPermitVal()
}
