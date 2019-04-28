package mcnmodel

import (
	mcnadminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/tool/validate"
	"go-common/library/time"
)

//CookieMidInterface cookie set interface, set mid from cookie to arg
type CookieMidInterface interface {
	SetMid(midFromCookie int64)
}

//McnCommonReq common mcn
type McnCommonReq struct {
	McnCheatReq
	SignID int64 `form:"sign_id"`
	McnMid int64
}

//CheatInterface cheat interface
type CheatInterface interface {
	// Cheat return true if cheated, false if not cheated
	Cheat() bool
}

//McnCheatReq cheat
type McnCheatReq struct {
	TMcnMid int64 `form:"t_mcn_mid"`
}

//Cheat .
func (m *McnCommonReq) Cheat() bool {
	if m.TMcnMid == 0 {
		return false
	}
	m.SetMid(m.TMcnMid)
	return true
}

//SetMid set mid
func (m *McnCommonReq) SetMid(midFromCookie int64) {
	m.McnMid = midFromCookie
}

//UpCommonReq common up
type UpCommonReq struct {
	UpMid int64
}

//SetMid set mid
func (m *UpCommonReq) SetMid(midFromCookie int64) {
	m.UpMid = midFromCookie
}

//GetStateReq get state
type GetStateReq struct {
	McnCommonReq
}

//McnApplyReq apply req
type McnApplyReq struct {
	McnCommonReq
	CompanyName        string `form:"company_name"`
	CompanyLicenseID   string `form:"company_license_id"`
	ContactName        string `form:"contact_name"`
	ContactTitle       string `form:"contact_title"`
	ContactIdcard      string `form:"contact_idcard" validate:"idcheck"`
	ContactPhone       string `form:"contact_phone" validate:"phonecheck"`
	CompanyLicenseLink string `form:"company_license_link" validate:"httpcheck"`
	ContractLink       string `form:"contract_link" validate:"httpcheck"`
}

//CopyTo .
func (m *McnApplyReq) CopyTo(v *McnSign) {
	if v == nil {
		return
	}
	v.McnMid = m.McnMid
	v.CompanyName = m.CompanyName
	v.CompanyLicenseID = m.CompanyLicenseID
	v.ContactName = m.ContactName
	v.ContactTitle = m.ContactTitle
	v.ContactIdcard = m.ContactIdcard
	v.ContactPhone = m.ContactPhone
	v.CompanyLicenseLink = m.CompanyLicenseLink
	v.ContractLink = m.ContractLink
}

//McnBindUpApplyReq .
type McnBindUpApplyReq struct {
	McnCommonReq
	UpMid        int64     `form:"up_mid"`
	BeginDate    time.Time `form:"begin_date"`
	EndDate      time.Time `form:"end_date"`
	ContractLink string    `form:"contract_link"` // 手动检查http格式
	UpAuthLink   string    `form:"up_auth_link"`  // 手动检查http格式
	UpType       int8      `form:"up_type"`       // 用户类型，0为站内，1为站外
	SiteLink     string    `form:"site_link"`     //up主站外账号链接, 如果up type为1，该项必填
	mcnadminmodel.Permits
	PublicationPrice int64 `form:"publication_price"` // 单位：1/1000 元
}

//IsSiteInfoOk 检查站外up主信息是否OK，如果不是站外Up主，直接返回ok
func (m *McnBindUpApplyReq) IsSiteInfoOk() bool {
	if m.UpType == 0 {
		return true
	}
	return validate.RegHTTPCheck.MatchString(m.SiteLink)
}

//CopyTo .
func (m *McnBindUpApplyReq) CopyTo(v *McnUp) {
	v.UpMid = m.UpMid
	v.McnMid = m.McnMid
	v.BeginDate = m.BeginDate
	v.EndDate = m.EndDate
	v.ContractLink = m.ContractLink
	v.UpAuthLink = m.UpAuthLink
	v.UpType = m.UpType
	v.SiteLink = m.SiteLink
	v.Permission = uint32(m.GetAttrPermitVal())
	v.PublicationPrice = m.PublicationPrice
}

//McnUpConfirmReq .
type McnUpConfirmReq struct {
	UpCommonReq
	BindID int64 `form:"bind_id"`
	Choice bool  `form:"choice"`
}

//McnUpGetBindReq .
type McnUpGetBindReq struct {
	UpCommonReq
	BindID int64 `form:"bind_id"`
}

// McnGetDataSummaryReq req
type McnGetDataSummaryReq = McnCommonReq

//McnGetUpListReq req
type McnGetUpListReq struct {
	McnCommonReq
	UpMid int64 `form:"up_mid"`
	model.PageArg
}

//McnGetAccountReq req
type McnGetAccountReq struct {
	Mid int64 `form:"mid"`
}

// McnGetMcnOldInfoReq req
type McnGetMcnOldInfoReq struct {
	McnCommonReq
}

// McnGetRankReq req to 获取排行
type McnGetRankReq struct {
	McnCommonReq
	Tid      int16    `form:"tid"` // 分区
	DataType DataType `form:"data_type"`
}

// McnGetRecommendPoolReq get recommend pool
type McnGetRecommendPoolReq struct {
	McnCommonReq
	model.PageArg
	Tid        int16  `form:"tid"`
	OrderField string `form:"order_field"`
	Sort       string `form:"sort"`
}

// McnGetRecommendPoolTidListReq common req
type McnGetRecommendPoolTidListReq = McnCommonReq

// ------inner request

// McnGetRankAPIReq req to 获取排行
type McnGetRankAPIReq struct {
	SignID   int64    `form:"sign_id"`
	Tid      int16    `form:"tid"` // 分区
	DataType DataType `form:"data_type"`
}

// 播放/弹幕/评论/分享/硬币/收藏/点赞数
const (
	ActionTypePlay  = "play"  //播放
	ActionTypeDanmu = "danmu" //弹幕
	ActionTypeReply = "reply" //评论
	ActionTypeShare = "share" //分享
	ActionTypeCoin  = "coin"  //硬币
	ActionTypeFav   = "fav"   //收藏
	ActionTypeLike  = "like"  //点赞数
)

const (
	// UserTypeGuest .
	UserTypeGuest = "guest" // 游客
	// UserTypeFans .
	UserTypeFans = "fans" // 粉丝
)

//McnGetIndexIncReq 增量趋势
type McnGetIndexIncReq struct {
	McnCommonReq
	Type string `form:"type"`
}

//McnGetIndexSourceReq 来源分区
type McnGetIndexSourceReq = McnGetIndexIncReq

//McnGetPlaySourceReq 稿件播放来源占比
type McnGetPlaySourceReq struct {
	McnCommonReq
}

//McnGetMcnFansReq mcn
type McnGetMcnFansReq = McnCommonReq

//McnGetMcnFansIncReq mcn粉丝按天关注数
type McnGetMcnFansIncReq = McnCommonReq

//McnGetMcnFansDecReq mcn粉丝按天取关数
type McnGetMcnFansDecReq = McnCommonReq

//McnGetMcnFansAttentionWayReq mcn粉丝关注渠道
type McnGetMcnFansAttentionWayReq = McnCommonReq

// McnGetBaseFansAttrReq  mcn 游客和粉丝基本数据请求
type McnGetBaseFansAttrReq struct {
	McnCommonReq
	UserType string `form:"user_type"`
}

// McnGetFansAreaReq mcn 地区分布请求
type McnGetFansAreaReq = McnGetBaseFansAttrReq

// McnGetFansTypeReq  mcn  游客/粉丝倾向分布请求
type McnGetFansTypeReq = McnGetBaseFansAttrReq

// McnGetFansTagReq  mcn  游客/粉丝标签地图分布请求
type McnGetFansTagReq = McnGetBaseFansAttrReq

//McnChangePermitReq change permit
type McnChangePermitReq struct {
	McnCommonReq
	UpMid int64 `form:"up_mid"`
	mcnadminmodel.Permits
	UpAuthLink string `form:"up_auth_link" validate:"httpcheck"`
}

//McnPublicationPriceChangeReq change publication price
type McnPublicationPriceChangeReq struct {
	McnCommonReq
	Price int64 `form:"price"`
	UpMid int64 `form:"up_mid"`
}
