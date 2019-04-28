package model

import "go-common/library/time"

// VipResourcePool vip_resource_pool table
type VipResourcePool struct {
	ID             int       `json:"id"`
	PoolName       string    `json:"pool_name"`
	BusinessID     int       `json:"business_id"`
	BusinessName   string    `json:"business_name"`
	Reason         string    `json:"reason"`
	CodeExpireTime time.Time `json:"code_expire_time"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Contacts       string    `json:"contacts"`
	ContactsNumber string    `json:"contacts_number"`
	Ctime          time.Time `json:"-"`
	Mtime          time.Time `json:"-"`
}

// VipResourceBatch vip_resource_batch table
type VipResourceBatch struct {
	ID             int       `json:"id"`
	PoolID         int       `json:"pool_id"`
	Unit           int       `json:"unit"`
	Count          int       `json:"count"`
	Ver            int       `json:"ver"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	SurplusCount   int       `json:"surplus_count"`
	CodeUseCount   int       `json:"code_use_count"`
	DirectUseCount int       `json:"direct_use_count"`
	Ctime          time.Time `json:"-"`
	Mtime          time.Time `json:"-"`
}

//VipBusinessInfo vip_business_info table
type VipBusinessInfo struct {
	ID             int       `json:"id" form:"id"`
	BusinessName   string    `json:"business_name" form:"business_name"`
	BusinessType   int       `json:"business_type" form:"business_type"`
	Status         int       `json:"status" form:"status"`
	AppKey         string    `json:"app_key" form:"app_key"`
	Secret         string    `json:"-" form:"secret"`
	Contacts       string    `json:"contacts" form:"contacts"`
	ContactsNumber string    `json:"contacts_number" form:"contacts_number"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
}

//VipChangeHistory vip_change_history table
type VipChangeHistory struct {
	ID         int       `json:"id"`
	Mid        int64     `json:"mid"`
	ChangeType int       `json:"change_type"`
	ChangeTime time.Time `json:"change_time"`
	Days       int       `json:"days"`
	OperatorID string    `json:"operator_id"`
	RelationID string    `json:"relation_id"`
	BatchID    int       `json:"batch_id"`
	Remark     string    `json:"remark"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

//VipUserInfo vip_user_info table
type VipUserInfo struct {
	ID                   int       `json:"id"`
	Mid                  int       `json:"mid"`
	VipType              int       `json:"vipType"`
	VipPayType           int       `json:"vipPayType"`
	VipStatus            int       `json:"vipStatus"`
	VipStartTime         time.Time `json:"vipStartTime"`
	VipOverdueTime       time.Time `json:"vipOverdueTime"`
	AnnualVipOverdueTime time.Time `json:"annualVipOverdueTime"`
	VipRecentTime        time.Time `json:"vipRecentTime"`
	Ctime                time.Time `json:"ctime"`
	Mtime                time.Time `json:"mtime"`
}

//VipChangeBo vip change
type VipChangeBo struct {
	Mid        int
	ChangeType int
	ChangeTime time.Time
	RelationID string
	Remark     string
	Days       int
	BatchID    int
	OperatorID string
}

//HandlerVip vip handler
type HandlerVip struct {
	OldVipUser *VipUserInfo
	VipUser    *VipUserInfo
	HistoryID  int
	Days       int
	Mid        int
}

//BcoinSendBo bcoinSendBo
type BcoinSendBo struct {
	Amount     int
	DayOfMonth int
	DueDate    time.Time
}

//VipBcoinSalary vip_bcoin_salary table
type VipBcoinSalary struct {
	ID            int       `json:"id"`
	Mid           int       `json:"mid"`
	Status        int       `json:"status"`
	GiveNowStatus int       `json:"giveNowStatus"`
	Month         time.Time `json:"month"`
	Amount        int       `json:"amount"`
	Memo          string    `json:"memo"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

//VipConfig vipConfig
type VipConfig struct {
	ID           int       `json:"id"`
	ConfigKey    string    `json:"configKey"`
	Name         string    `json:"name"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	OperatorID   int       `json:"operatorId"`
	OperatorName string    `json:"operatorName"`
	Mtime        time.Time `json:"mtime"`
}

// VipAppVersion app version.
type VipAppVersion struct {
	ID         int64  `json:"id"`
	PlatformID int8   `json:"platform_id"`
	Version    string `json:"version"`
	Tip        string `json:"tip"`
	Operator   string `json:"operator"`
	Link       string `json:"link"`
}

// VipPrivilege .
type VipPrivilege struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	OrderNum   int    `json:"order_num"`
	Remark     string `json:"remark"`
	BgColor    string `json:"bg_color"`
	Type       int    `json:"type"`
	Deleted    int    `json:"deleted"`
	PcLink     string `json:"pc_link"`
	MobileLink string `json:"mobile_link"`
}

// VipPrivilegeMapping  .
type VipPrivilegeMapping struct {
	ID          int    `json:"id"`
	PrivilegeID int    `json:"privilege_id"`
	PlatformID  int    `json:"platform_id"`
	Icon        string `json:"icon"`
	Status      int    `json:"status"`
	Operator    string `json:"operator"`
}

// const vip enum value
const (
	//ChangeType
	ChangeTypePointExhchange  = 1 // 积分兑换
	ChangeTypeRechange        = 2 //充值开通
	ChangeTypeSystem          = 3 // 系统发放
	ChangeTypeActiveGive      = 4 //活动赠送
	ChangeTypeRepeatDeduction = 5 //重复领取扣除
	ChangeTypeSystemDrawback  = 7 //系统扣减

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
)

// const .
const (
	NOTUSER = iota + 1
	USED
	FROZEN
)
