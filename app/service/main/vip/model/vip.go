package model

import (
	"go-common/library/time"
)

//VipChangeHistory vip_change_history table
type VipChangeHistory struct {
	ID          int64     `json:"id"`
	Mid         int64     `json:"mid"`
	ChangeType  int8      `json:"changeType"`
	ChangeTime  time.Time `json:"changeTime"`
	Days        int64     `json:"days"`
	Month       int16     `json:"month"`
	OperatorID  string    `json:"operatorId"`
	RelationID  string    `json:"relationId"`
	BatchID     int64     `json:"batchId"`
	Remark      string    `json:"remark"`
	Ctime       time.Time `json:"ctime"`
	BatchCodeID int64     `json:"batchCodeId"`
}

//VipAppInfo vip app info
type VipAppInfo struct {
	ID       int64     `json:"id"`
	Type     int8      `json:"type"`
	Name     string    `json:"name"`
	AppKey   string    `json:"appKey"`
	PurgeURL string    `json:"purgeUrl"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

//VipBcoinSalary vip_bcoin_salary table
type VipBcoinSalary struct {
	ID            int64     `json:"id"`
	Mid           int64     `json:"mid"`
	Status        int8      `json:"status"`
	GiveNowStatus int8      `json:"give_now_status"`
	Month         time.Time `json:"month"`
	PayDay        time.Time `json:"payday"`
	Amount        int64     `json:"amount"`
	Memo          string    `json:"memo"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

//VipConfig vipConfig
type VipConfig struct {
	ID           int64     `json:"id"`
	ConfigKey    string    `json:"configKey"`
	Name         string    `json:"name"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	OperatorID   int64     `json:"operatorId"`
	OperatorName string    `json:"operatorName"`
	Mtime        time.Time `json:"mtime"`
}

//VipChangeBo vip change
type VipChangeBo struct {
	Mid         int64
	ChangeType  int8
	ChangeTime  time.Time
	RelationID  string
	Remark      string
	Days        int64
	Months      int16
	BatchID     int64
	BatchCodeID int64
	OperatorID  string
}

//HandlerVip vip handler
type HandlerVip struct {
	OldVipUser *VipInfoDB
	VipUser    *VipInfoDB
	HistoryID  int64
	Days       int64
	Months     int16
	Mid        int64
	ToMid      int64
}

//OldHandlerVip old vip handler
type OldHandlerVip struct {
	OldVipUser *VipUserInfo
	VipUser    *VipUserInfo
	HistoryID  int64
	Days       int64
	Months     int16
	Mid        int64
	ToMid      int64
}

//BcoinSendBo bcoinSendBo
type BcoinSendBo struct {
	Amount     int64
	DayOfMonth int64
	DueDate    time.Time
}

//VipBo vipBo
type VipBo struct {
	Mid       int64 `json:"mid"`
	VipStatus int8  `json:"vipStatus"`
	VipType   int8  `json:"vipType"`
}

//VipListVo vipListVo
type VipListVo struct {
	VipList []*VipBo `json:"vipList"`
	ID      int64    `json:"id"`
}

// VipInfoResp vipinfo resp.
type VipInfoResp struct {
	Mid            int64  `json:"mid"`
	VipType        int8   `json:"vip_type"`
	PayType        int8   `json:"pay_type"`
	PayChannelID   int32  `json:"pay_channel_id"`
	VipStatus      int32  `json:"vip_status"`
	VipTotalMsec   int64  `json:"vip_total_sec"`
	VipHoldMsec    int64  `json:"vip_hold_sec"`
	VipDueMsec     int64  `json:"vip_due_sec"`
	VipSurplusMsec int64  `json:"vip_surplus_sec"`
	DueRemark      string `json:"due_remark"`
	VipDueDate     int64  `json:"vip_due_date"`
	VipRecentTime  int64  `json:"vip_recent_time"`
	AutoRenewed    int32  `json:"auto_renewed"`
}

//VipChangeHistoryVo .
type VipChangeHistoryVo struct {
	ID            string           `json:"id"`
	ChangeType    int8             `json:"change_type"`
	ChangeTypeStr string           `json:"change_type_str"`
	ChangeTime    int64            `json:"change_time"`
	Month         int16            `json:"month"`
	OpenRemark    string           `json:"open_remark"`
	Days          int32            `json:"days"`
	Remark        string           `json:"remark"`
	Actives       []*VipActiveShow `json:"actives"`
}

//Eunm vip enum value.
const (
	//ChangeType
	ChangeTypePointExhchange  = 1 // 积分兑换
	ChangeTypeRechange        = 2 //充值开通
	ChangeTypeSystem          = 3 // 系统发放
	ChangeTypeActiveGive      = 4 //活动赠送
	ChangeTypeRepeatDeduction = 5 //重复领取扣除

	VipDaysMonth = 31
	VipDaysYear  = 366

	NotVip    = 0 //非大会员
	Vip       = 1 //月度大会员
	AnnualVip = 2 //年度会员

	VipStatusOverTime    = 0 //过期
	VipStatusNotOverTime = 1 //未过期
	VipStatusFrozen      = 2 //冻结
	VipStatusBan         = 3 //封禁

	VipAppUser  = 1 //大会员对接业务方user缓存
	VipAppPoint = 2 //大会员对接业务方积分缓存

	VipChangeFrozen   = -1 //冻结
	VipChangeUnFrozen = 0  //解冻
	VipChangeOpen     = 1  //开通
	VipChangeModify   = 2  //变更

	VipBusinessStatusOpen  = 0 //有效
	VipBusinessStatusClose = 1 //无效

	VipOpenMsgTitle     = "大会员开通成功"
	VipSystemNotify     = 4
	VipOpenMsg          = "恭喜您已开通大会员服务%s！"
	VipOpenKMsg         = "恭喜您已续期大会员服务%s！"
	VipBcoinGiveContext = "尊敬的年度大会员，您本月%dB币到账啦！请您随意挥霍，注意会在次月%d日清零哦！"
	VipBcoinGiveTitle   = "B币到账通知"

	VipOpenMsgCode      = "10_1_1"
	VipBcoinGiveMsgCode = "10_99_2"
	VipCustomizeMsgCode = "10_99_1"

	AnnualVipBcoinDay              = "annual_vip_bcoin_day"                //年费VIPB券发放每月第几天
	AnnualVipBcoinCouponMoney      = "annual_vip_bcoin_coupon_money"       //年费VIP返回B券金额
	AnnualVipBcoinCouponActivityID = "annual_vip_bcoin_coupon_activity_id" //年费VIP返B券活动ID

)

// vip AccessStatus.
const (
	WebHadAccess int32 = iota
)

//vip renew type
const (
	NomalVip = iota
	AuoRenewVip
)
