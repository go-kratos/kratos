package model

import (
	"go-common/library/time"
)

//allowance origin
const (
	AllowanceNone = iota
	AllowanceSystemAdmin
	AllowanceBusinessReceive
)

// blance change type
const (
	VipSalary int64 = iota + 1
	SystemAdminSalary
	Consume
	ConsumeFaildBack
)

// coupon type
const (
	CouponVideo = iota + 1
	CouponCartoon
	CouponAllowance
	CouponAllowanceCode
)

// coupon state.
const (
	NotUsed = iota
	InUse
	Used
	Expire
	Block
)

// allowance explain
const (
	NoLimitExplain = "不限定"
	ScopeFmt       = "仅限%s端使用"
)

// batch state
const (
	BatchStateNormal int8 = iota
	BatchStateBlock
)

// batch origin
const (
	AdminSalaryOrigin int64 = iota + 1
)

// allowance change type
const (
	AllowanceSalary int8 = iota + 1
	AllowanceConsume
	AllowanceCancel
	AllowanceConsumeSuccess
	AllowanceConsumeFaild
	AllowanceBlock
	AllowanceUnBlock
)

// coupon_batch_info表 product_limit_renewal字段.
const (
	ProdLimRenewalAll int8 = iota
	ProdLimRenewalAuto
	ProdLimRenewalNotAuto
)

// coupon_batch_info表 product_limit_renewal字段.
const (
	None           int8 = 0
	ProdLimMonth1       = 1
	ProdLimMonth3       = 3
	ProdLimMonth12      = 12
)

// ProdLimit .
var (
	ProdLimMonthMap   = map[int8]string{None: "", ProdLimMonth1: "1月", ProdLimMonth3: "3月", ProdLimMonth12: "12月"}
	ProdLimRenewalMap = map[int8]string{ProdLimRenewalAll: "", ProdLimRenewalAuto: "自动续期", ProdLimRenewalNotAuto: "非自动续期"}
)

// PageInfo common page info.
type PageInfo struct {
	Count       int         `json:"count"`
	CurrentPage int         `json:"currentPage,omitempty"`
	Item        interface{} `json:"item"`
}

// CouponBatchInfo info.
type CouponBatchInfo struct {
	ID             int64     `json:"id"`
	AppID          int64     `json:"app_id"`
	Name           string    `json:"name"`
	BatchToken     string    `json:"batch_token"`
	MaxCount       int64     `json:"max_count"`
	CurrentCount   int64     `json:"current_count"`
	StartTime      int64     `json:"start_time"`
	ExpireTime     int64     `json:"expire_time"`
	ExpireDay      int64     `json:"expire_day"`
	Ver            int64     `json:"ver"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
	Operator       string    `json:"operator"`
	LimitCount     int64     `json:"limit_count"`
	FullAmount     float64   `json:"full_amount"`
	Amount         float64   `json:"amount"`
	State          int8      `json:"state"`
	CouponType     int8      `json:"coupon_type"`
	PlatformLimit  string    `json:"platform_limit"`
	ProdLimMonth   int8      `json:"product_limit_month"`
	ProdLimRenewal int8      `json:"product_limit_Renewal"`
}

// ArgBatchInfo arg.
type ArgBatchInfo struct {
	AppID      int64  `form:"app_id" validate:"required,min=1,gte=1"`
	Name       string `form:"name" validate:"required"`
	MaxCount   int64  `form:"max_count" validate:"required,min=1,gte=1"`
	LimitCount int64  `form:"limit_count"`
	StartTime  int64  `form:"start_time" validate:"required,min=1,gte=1"`
	ExpireTime int64  `form:"end_time" validate:"required,min=1,gte=1"`
}

// ArgAllowanceBatchInfo allowance arg.
type ArgAllowanceBatchInfo struct {
	AppID          int64   `form:"app_id" validate:"required,min=1,gte=1"`
	Name           string  `form:"name" validate:"required"`
	MaxCount       int64   `form:"max_count"`
	LimitCount     int64   `form:"limit_count"`
	StartTime      int64   `form:"start_time"`
	ExpireTime     int64   `form:"end_time"`
	ExpireDay      int64   `form:"expire_day" default:"-1"`
	Amount         float64 `form:"amount" validate:"required,min=1,gte=1"`
	FullAmount     float64 `form:"full_amount" validate:"required,min=1,gte=1"`
	PlatformLimit  []int64 `form:"platform_limit,split"`
	ProdLimMonth   int8    `form:"product_limit_month"`
	ProdLimRenewal int8    `form:"product_limit_Renewal" validate:"gte=0,lte=2"`
}

// ArgAllowanceBatchInfoModify allowance modify arg.
type ArgAllowanceBatchInfoModify struct {
	ID             int64   `form:"id" validate:"required,min=1,gte=1"`
	AppID          int64   `form:"app_id" validate:"required,min=1,gte=1"`
	Name           string  `form:"name" validate:"required"`
	MaxCount       int64   `form:"max_count" `
	LimitCount     int64   `form:"limit_count"`
	PlatformLimit  []int64 `form:"platform_limit,split"`
	ProdLimMonth   int8    `form:"product_limit_month" validate:"gte=0"`
	ProdLimRenewal int8    `form:"product_limit_Renewal" validate:"gte=0,lte=2"`
}

// ArgAllowance arg.
type ArgAllowance struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// ArgAllowanceInfo arg.
type ArgAllowanceInfo struct {
	BatchToken string `form:"batch_token" validate:"required"`
}

// ArgAllowanceSalary allowance salary arg.
type ArgAllowanceSalary struct {
	Mids       []int64 `form:"mids,split"`
	BatchToken string  `form:"batch_token" validate:"required"`
	MsgType    string  `form:"msg_type" default:"vip"`
}

// ArgAllowanceState arg.
type ArgAllowanceState struct {
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	CouponToken string `form:"coupon_token" validate:"required"`
}

// ArgBatchList arg.
type ArgBatchList struct {
	AppID int64 `form:"app_id"`
	Type  int8  `form:"type" default:"3"`
}

// ArgSalaryCoupon salary coupon.
type ArgSalaryCoupon struct {
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	CouponType  int64  `form:"coupon_type" validate:"required,min=1,gte=1"`
	Count       int    `form:"count" validate:"required,min=1,gte=1"`
	BranchToken string `form:"branch_token" validate:"required"`
}

// ArgUploadFile upload file arg.
type ArgUploadFile struct {
	FileURL string `form:"url" validate:"required"`
}

// CouponBatchResp resp.
type CouponBatchResp struct {
	ID                  int64   `json:"id"`
	AppID               int64   `json:"app_id"`
	AppName             string  `json:"app_name"`
	Name                string  `json:"name"`
	BatchToken          string  `json:"batch_token"`
	MaxCount            int64   `json:"max_count"`
	CurrentCount        int64   `json:"current_count"`
	StartTime           int64   `json:"start_time"`
	ExpireTime          int64   `json:"expire_time"`
	ExpireDay           int64   `json:"expire_day"`
	Operator            string  `json:"operator"`
	LimitCount          int64   `json:"limit_count"`
	ProductLimitExplain string  `json:"product_limit_explain"`
	PlatfromLimit       []int64 `json:"platform_limit"`
	UseLimitExplain     string  `json:"use_limit_explain"`
	State               int8    `json:"state"`
	Amount              float64 `json:"amount"`
	FullAmount          float64 `json:"full_amount"`
	ProdLimMonth        int8    `json:"product_limit_month"`
	ProdLimRenewal      int8    `json:"product_limit_Renewal"`
}

// AppInfo app info.
type AppInfo struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Appkey    string    `json:"appkey"`
	NotifyURL string    `json:"notify_url"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// CouponResp def.
type CouponResp struct {
	Token     string `json:"token"`
	Mid       int64  `json:"mid"`
	GrantTime int64  `json:"grant_time"`
	UseTime   int64  `json:"use_time"`
	State     int8   `json:"state"`
	Remark    int8   `json:"remark"`
}

// CouponAllowanceInfo coupon allowance info.
type CouponAllowanceInfo struct {
	ID          int64     `json:"id"`
	CouponToken string    `json:"coupon_token"`
	Mid         int64     `json:"mid"`
	State       int32     `json:"state"`
	StartTime   int64     `json:"start_time"`
	ExpireTime  int64     `json:"expire_time"`
	Origin      int64     `json:"origin"`
	OrderNO     string    `json:"order_no"`
	Ver         int64     `json:"ver"`
	Remark      string    `json:"remark"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
	BatchToken  string    `json:"batch_token"`
	Amount      float64   `json:"amount"`
	FullAmount  float64   `json:"full_amount"`
	AppID       int64     `json:"app_id"`
}

// CouponAllowanceChangeLog coupon allowance change log.
type CouponAllowanceChangeLog struct {
	ID          int64     `json:"-"`
	CouponToken string    `json:"coupon_token"`
	OrderNO     string    `json:"order_no"`
	Mid         int64     `json:"mid"`
	State       int8      `json:"state"`
	ChangeType  int8      `json:"change_type"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// ProdLimExplainFmt .
func (c *CouponBatchResp) ProdLimExplainFmt(prodLimMonth, prodLimRenewal int8) {
	if prodLimMonth == None && prodLimRenewal == None {
		c.ProductLimitExplain = NoLimitExplain
	}
	pstr := ""
	if limm, ok := ProdLimMonthMap[prodLimMonth]; ok {
		pstr += limm
	}
	if limr, ok := ProdLimRenewalMap[prodLimRenewal]; ok {
		pstr += "、" + limr
	}
	c.ProductLimitExplain = pstr
}

//Sizer .
type Sizer interface {
	Size() int64
}

//ArgCouponViewBatch .
type ArgCouponViewBatch struct {
	ID           int64  `form:"id"`
	Name         string `form:"name" validate:"required"`
	AppID        int64  `form:"app_id" validate:"required,min=1"`
	MaxCount     int64  `form:"max_count" default:"-1"`
	CurrentCount int64  `form:"current_count"`
	LimitCount   int64  `form:"limit_count" default:"-1"`
	StartTime    int64  `form:"start_time" validate:"required,min=1"`
	ExpireTime   int64  `form:"end_time" validate:"required,min=1"`
	Operator     string `form:"operator"`
	Ver          int64
	BatchToken   string
	CouponType   int8
}

//ArgSearchCouponView .
type ArgSearchCouponView struct {
	PN          int    `form:"pn" default:"1"`
	PS          int    `form:"ps" default:"20"`
	Mid         int64  `form:"mid" validate:"required"`
	CouponToken string `form:"coupon_token"`
	AppID       int64  `form:"app_id"`
	BatchToken  string `form:"batch_token"`
	BatchTokens []string
}

//CouponInfo .
type CouponInfo struct {
	CouponToken string    `json:"coupon_token"`
	Mid         int64     `json:"mid"`
	State       int8      `json:"state"`
	StartTime   int64     `json:"start_time"`
	ExpireTime  int64     `json:"expire_time"`
	Origin      int8      `json:"origin"`
	CouponType  int8      `json:"coupon_type"`
	OrderNo     string    `json:"order_no"`
	OID         int32     `json:"oid"`
	Remark      string    `json:"remark"`
	UseVer      int64     `json:"use_ver"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
	BatchToken  string    `json:"batch_token"`
	Title       string    `json:"title"`
	BatchName   string    `json:"batch_name"`
}

//PGCInfoResq .
type PGCInfoResq struct {
	Title string `json:"title"`
}

//CouponChangeLog .
type CouponChangeLog struct {
	CouponToken string `json:"coupon_token"`
	Mid         int64  `json:"mid"`
	State       int8   `json:"state"`
}

// ArgBatchSalaryCoupon batch salary coupon.
type ArgBatchSalaryCoupon struct {
	FileURL     string `form:"file_url" validate:"required"`
	Count       int64  `form:"count" validate:"required,min=1,gte=1"`
	BranchToken string `form:"branch_token" validate:"required"`
	SliceSize   int    `form:"slice_size" default:"100" validate:"min=100,max=10000"`
}

// ArgCouponCode coupon code.
type ArgCouponCode struct {
	ID          int64  `form:"id"`
	BatchToken  string `form:"batch_token"`
	State       int32  `form:"state"`
	Code        string `form:"code"`
	Mid         int64  `form:"mid"`
	CouponType  int32  `form:"coupon_type"`
	CouponToken string `form:"coupon_token"`
	Pn          int    `form:"pn"`
	Ps          int    `form:"ps"`
}

// CouponCode coupon code.
type CouponCode struct {
	ID          int64     `json:"id"`
	BatchToken  string    `json:"batch_token"`
	State       int32     `json:"state"`
	Code        string    `json:"code"`
	Mid         int64     `json:"mid"`
	CouponType  int32     `json:"coupon_type"`
	CouponToken string    `json:"coupon_token"`
	Ver         int64     `json:"ver"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// CodePage code page.
type CodePage struct {
	Count    int64         `json:"count"`
	CodeList []*CouponCode `json:"code_list"`
}

// coupon code state.
const (
	CodeStateNotUse = iota + 1
	CodeStateUsed
	CodeStateBlock
	CodeStateExpire
)

// batch code max count.
const (
	BatchCodeMaxCount = 50000
	BatchAddCodeSlice = 100
)

// code batch state.
const (
	CodeBatchUsable = iota
	CodeBatchBlock
	CodeBatchExpire
)
