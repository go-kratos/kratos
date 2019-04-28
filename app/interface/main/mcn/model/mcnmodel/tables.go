package mcnmodel

import (
	"go-common/app/interface/main/mcn/model"
	"go-common/library/time"
)

// go get github.com/abice/go-enum
// go:generate go-enum -f=tables.go

// table names
const (
	//TableNameMcnUp up
	TableNameMcnUp = "mcn_up"
	//TableNameMcnSign sign
	TableNameMcnSign = "mcn_sign"
	//TableNameMcnDataSummary data summary
	TableNameMcnDataSummary = "mcn_data_summary"
	//TableNameMcnDataUp data up
	TableNameMcnDataUp = "mcn_data_up"
	//TableNameUpBaseInfo up base info
	TableNameUpBaseInfo = "up_base_info"
	//TableNameMcnRankUpFan rank for up fans
	TableNameMcnRankUpFan = "mcn_rank_up_fans"
	//TableNameMcnRankArchiveLike rank for archive likes
	TableNameMcnRankArchiveLike = "mcn_rank_archive_likes"
	//TableNameMcnUpRecommendPool up recomment pool
	TableNameMcnUpRecommendPool = "mcn_up_recommend_pool"
	TableMcnUpPermissionApply   = "mcn_up_permission_apply"
)

// DataType 数据类型，1累计，2昨日，3上周，4上月
/* ENUM(
Accumulate = 1
Day = 2
Week = 3
Month = 4
ActiveFans = 5
)*/
const (
	// DataTypeAccumulate is a DataType of type Accumulate
	DataTypeAccumulate DataType = 1
	// DataTypeDay is a DataType of type Day
	DataTypeDay DataType = 2
	// DataTypeWeek is a DataType of type Week
	DataTypeWeek DataType = 3
	// DataTypeMonth is a DataType of type Month
	DataTypeMonth DataType = 4
	// DataTypeActiveFans active fans
	DataTypeActiveFans DataType = 5
)

const (
	// McnDataTypeDay is a McnDataType of type Day
	McnDataTypeDay McnDataType = 1
	// McnDataTypeMonth is a McnDataType of type Month
	McnDataTypeMonth McnDataType = 2
)

// DataType .
type DataType int8

// McnUp up table
type McnUp struct {
	ID               int64            `json:"id" gorm:"column:id"`
	SignID           int64            `json:"sign_id" gorm:"column:sign_id"`
	McnMid           int64            `json:"mcn_mid" gorm:"column:mcn_mid"`
	UpMid            int64            `json:"up_mid" gorm:"column:up_mid"`
	BeginDate        time.Time        `json:"begin_date" gorm:"column:begin_date"`
	EndDate          time.Time        `json:"end_date" gorm:"column:end_date"`
	ContractLink     string           `json:"contract_link" gorm:"column:contract_link"`
	UpAuthLink       string           `json:"up_auth_link" gorm:"column:up_auth_link"`
	RejectReason     string           `json:"reject_reason" gorm:"column:reject_reason"`
	RejectTime       string           `json:"reject_time" gorm:"column:reject_time"`
	State            model.MCNUPState `json:"state" gorm:"column:state"`
	UpType           int8             `json:"up_type" gorm:"column:up_type"`     // 用户类型，0为站内，1为站外
	SiteLink         string           `json:"site_link" gorm:"column:site_link"` //up主站外账号链接, 如果up type为1，该项必填
	StateChangeTime  time.Time        `json:"state_change_time" gorm:"column:state_change_time"`
	ConfirmTime      time.Time        `json:"confirm_time" gorm:"column:confirm_time"`
	Ctime            time.Time        `json:"ctime" gorm:"column:ctime"`
	Mtime            time.Time        `json:"mtime" gorm:"column:mtime"`
	Permission       uint32           `gorm:"column:permission" json:"permission"`
	PublicationPrice int64            `gorm:"column:publication_price" json:"publication_price"` // 单位：1/1000 元
}

// TableName .
func (s *McnUp) TableName() string {
	return TableNameMcnUp
}

// IsBindable check if up canbe bind to other
func (s *McnUp) IsBindable() bool {
	return isUpBindable(s.State)
}

// IsBeingBindedWithMcn check this up is in the middle of being binded with mcn,
func (s *McnUp) IsBeingBindedWithMcn(mcn *McnSign) bool {
	if mcn == nil {
		return false
	}
	if s.SignID == mcn.ID &&
		(s.State == model.MCNUPStateOnReview || s.State == model.MCNUPStateNoAuthorize) {
		return true
	}
	return false
}

// McnSign mcn sign table
type McnSign struct {
	ID                 int64              `json:"id" gorm:"column:id"`
	McnMid             int64              `json:"mcn_mid" gorm:"column:mcn_mid"`
	CompanyName        string             `json:"company_name" gorm:"column:company_name"`
	CompanyLicenseID   string             `json:"company_license_id" gorm:"column:company_license_id"`
	CompanyLicenseLink string             `json:"company_license_link" gorm:"column:company_license_link"`
	ContractLink       string             `json:"contract_link" gorm:"column:contract_link"`
	ContactName        string             `json:"contact_name" gorm:"column:contact_name"`
	ContactTitle       string             `json:"contact_title" gorm:"column:contact_title"`
	ContactIdcard      string             `json:"contact_idcard" gorm:"column:contact_idcard"`
	ContactPhone       string             `json:"contact_phone" gorm:"column:contact_phone"`
	BeginDate          time.Time          `json:"begin_date" gorm:"column:begin_date"`
	EndDate            time.Time          `json:"end_date" gorm:"column:end_date"`
	RejectReason       string             `json:"reject_reason" gorm:"column:reject_reason"`
	RejectTime         time.Time          `json:"reject_time" gorm:"column:reject_time"`
	State              model.MCNSignState `json:"state" gorm:"column:state"`
	Ctime              time.Time          `json:"ctime" gorm:"column:ctime"`
	Mtime              time.Time          `json:"mtime" gorm:"column:mtime"`
	Permission         uint32             `json:"permission" gorm:"column:permission"`
}

// TableName table name
func (s *McnSign) TableName() string {
	return TableNameMcnSign
}

// McnDataType .
/* ENUM(
Day = 1
Month = 2
)*/
type McnDataType int8

// McnDataSummary table
type McnDataSummary struct {
	ID                       int64       `json:"id" gorm:"column:id"`
	McnMid                   int64       `json:"mcn_mid" gorm:"column:mcn_mid"`
	SignID                   int64       `json:"sign_id" gorm:"column:sign_id"`
	UpCount                  int64       `json:"up_count" gorm:"column:up_count"`
	FansCountAccumulate      int64       `json:"fans_count_accumulate" gorm:"column:fans_count_accumulate"`
	FansCountOnline          int64       `json:"fans_count_online" gorm:"column:fans_count_online"`
	FansCountReal            int64       `json:"fans_count_real" gorm:"column:fans_count_real"`
	FansCountCheatAccumulate int64       `json:"fans_count_cheat_accumulate" gorm:"column:fans_count_cheat_accumulate"`
	FansCountIncreaseDay     int64       `json:"fans_count_increase_day" gorm:"column:fans_count_increase_day"`
	PlayCountAccumulate      int64       `json:"play_count_accumulate" gorm:"column:play_count_accumulate"`
	PlayCountIncreaseDay     int64       `json:"play_count_increase_day" gorm:"column:play_count_increase_day"`
	ArchiveCountAccumulate   int64       `json:"archive_count_accumulate" gorm:"column:archive_count_accumulate"`
	ActiveTid                int64       `json:"active_tid" gorm:"column:active_tid"`
	GenerateDate             time.Time   `json:"generate_date" gorm:"column:generate_date"`
	DataType                 McnDataType `json:"data_type" gorm:"column:data_type"`
	Ctime                    time.Time   `json:"ctime" gorm:"column:ctime"`
	Mtime                    time.Time   `json:"mtime" gorm:"column:mtime"`
}

// TableName table name
func (s *McnDataSummary) TableName() string {
	return TableNameMcnDataSummary
}

// McnDataUp table name
type McnDataUp struct {
	ID                     int64     `json:"id" gorm:"column:id"`
	McnMid                 int64     `json:"mcn_mid" gorm:"column:mcn_mid"`
	SignID                 int64     `json:"sign_id" gorm:"column:sign_id"`
	UpMid                  int64     `json:"up_mid" gorm:"column:up_mid"`
	DataType               DataType  `json:"data_type" gorm:"column:data_type"`
	FansIncreaseAccumulate int32     `json:"fans_increase_accumulate" gorm:"column:fans_increase_accumulate"`
	ArchiveCount           int32     `json:"archive_count" gorm:"column:archive_count"`
	PlayCount              int64     `json:"play_count" gorm:"column:play_count"`
	FansIncreaseMonth      int64     `json:"fans_increase_month" gorm:"column:fans_increase_month"`
	FansCount              int64     `json:"fans_count" gorm:"column:fans_count"`
	FansCountActive        int64     `json:"fans_count_active" gorm:"column:fans_count_active"`
	GenerateDate           time.Time `json:"generate_date" gorm:"column:generate_date"`
	Ctime                  time.Time `json:"ctime" gorm:"column:ctime"`
	Mtime                  time.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName table name
func (s *McnDataUp) TableName() string {
	return TableNameMcnDataUp
}

// UpBaseInfo  struct
type UpBaseInfo struct {
	ID                     uint32 `gorm:"column:id"`
	Mid                    int64  `gorm:"column:mid"`
	ActiveTid              int64  `gorm:"column:active_tid"`
	ArticleCountAccumulate int    `gorm:"column:article_count_accumulate"`
	Activity               int    `gorm:"column:activity"`
	FansCount              int    `gorm:"column:fans_count"`
}

// TableName .
func (s *UpBaseInfo) TableName() string {
	return TableNameUpBaseInfo
}

// McnRankUpFan .
type McnRankUpFan struct {
	ID           int64     `json:"id" gorm:"column:id"`
	McnMid       int64     `json:"mcn_mid" gorm:"column:mcn_mid"`
	SignID       int64     `json:"sign_id" gorm:"column:sign_id"`
	UpMid        int64     `json:"up_mid" gorm:"column:up_mid"`
	Value1       int64     `json:"value1" gorm:"column:value1"`
	Value2       int64     `json:"value2" gorm:"column:value2"`
	ActiveTid    int16     `json:"active_tid" gorm:"column:active_tid"`
	DataType     DataType  `json:"data_type" gorm:"column:data_type"`
	GenerateDate time.Time `json:"generate_date" gorm:"column:generate_date"`
	Ctime        time.Time `json:"ctime" gorm:"column:ctime"`
	Mtime        time.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (s *McnRankUpFan) TableName() string {
	return TableNameMcnRankUpFan
}

// McnRankArchiveLike .
type McnRankArchiveLike struct {
	ID           int64     `json:"id" gorm:"column:id"`
	McnMid       int64     `json:"mcn_mid" gorm:"column:mcn_mid"`
	SignID       int64     `json:"sign_id" gorm:"column:sign_id"`
	UpMid        int64     `json:"up_mid" gorm:"column:up_mid"`
	ArchiveID    int64     `json:"archive_id" gorm:"column:avid"`
	LikeCount    int64     `json:"like_count" gorm:"column:like_count"`
	PlayIncr     int64     `json:"play_incr" gorm:"column:play_incr" `
	DataType     DataType  `json:"data_type" gorm:"column:data_type"`
	Tid          int16     `json:"tid" gorm:"column:tid"`
	GenerateDate time.Time `json:"generate_date" gorm:"column:generate_date"`
	Ctime        time.Time `json:"ctime" gorm:"column:ctime"`
	Mtime        time.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (s *McnRankArchiveLike) TableName() string {
	return TableNameMcnRankArchiveLike
}

// McnUpRecommendPool 推荐池 struct
type McnUpRecommendPool struct {
	ID                     int64     `gorm:"column:id" json:"id"`
	UpMid                  int64     `gorm:"column:up_mid" json:"up_mid"`
	FansCount              int64     `gorm:"column:fans_count" json:"fans_count"`
	FansCountIncreaseMonth int64     `gorm:"column:fans_count_increase_month" json:"fans_count_increase_month"`
	ArchiveCount           int64     `gorm:"column:archive_count" json:"archive_count"`
	PlayCountAccumulate    int64     `gorm:"column:play_count_accumulate" json:"play_count_accumulate"`
	PlayCountAverage       int64     `gorm:"column:play_count_average" json:"play_count_average"`
	ActiveTid              int16     `gorm:"column:active_tid" json:"active_tid"`
	LastArchiveTime        time.Time `gorm:"column:last_archive_time" json:"last_archive_time"`
	State                  uint8     `gorm:"column:state" json:"state"`
	Source                 int64     `gorm:"column:source" json:"source"`
	GenerateTime           time.Time `gorm:"column:generate_time" json:"generate_time"`
	Ctime                  time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime                  time.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName table name.
func (s *McnUpRecommendPool) TableName() string {
	return TableNameMcnUpRecommendPool
}

// MCNUPRecommendState 推荐状态
// MCNUPRecommendState .
type MCNUPRecommendState int8

// const .
const (
	// MCNUPRecommendStateUnKnown 未知状态
	MCNUPRecommendStateUnKnown MCNUPRecommendState = 0
	// MCNUPRecommendStateOff 未推荐
	MCNUPRecommendStateOff MCNUPRecommendState = 1
	// MCNUPRecommendStateOn  推荐中
	MCNUPRecommendStateOn MCNUPRecommendState = 2
	// MCNUPRecommendStateBan 禁止推荐
	MCNUPRecommendStateBan MCNUPRecommendState = 3
	// MCNUPRecommendStateDel 移除中
	MCNUPRecommendStateDel MCNUPRecommendState = 100
)

// MCNUPRecommendSource
// type MCNUPRecommendSource mcnadminmodel.MCNUPRecommendSource

//McnUpPermissionApply permission
type McnUpPermissionApply struct {
	ID            int64     `gorm:"column:id" json:"id"`
	McnMid        int64     `gorm:"column:mcn_mid" json:"mcn_mid"`
	UpMid         int64     `gorm:"column:up_mid" json:"up_mid"`
	SignID        int64     `gorm:"column:sign_id" json:"sign_id"`
	NewPermission uint32    `gorm:"column:new_permission" json:"new_permission"`
	OldPermission uint32    `gorm:"column:old_permission" json:"old_permission"`
	RejectReason  string    `gorm:"-" json:"reject_reason"`
	RejectTime    time.Time `gorm:"-" json:"reject_time"`
	State         int8      `gorm:"column:state" json:"state"`
	Ctime         time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime         time.Time `gorm:"column:mtime" json:"mtime"`
	AdminID       int64     `gorm:"-" json:"-"`
	AdminName     string    `gorm:"-" json:"-"`
	UpAuthLink    string    `gorm:"column:up_auth_link" json:"up_auth_link"`
}

// TableName table name.
func (s *McnUpPermissionApply) TableName() string {
	return TableMcnUpPermissionApply
}
