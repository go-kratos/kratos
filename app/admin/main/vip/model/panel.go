package model

import "go-common/library/time"

// VipPriceConfigPlat vip价格面版平台 1. 其他平台 2.IOS平台 3.IOS的HD平台
type VipPriceConfigPlat int64

// VipPriceConfigStatus vip价格面版配置状态 0. 有效 1. 失效 2.待生效
type VipPriceConfigStatus int8

const (
	// VipPriceConfigStatusON 有效
	VipPriceConfigStatusON VipPriceConfigStatus = iota
	// VipPriceConfigStatusOFF 失效
	VipPriceConfigStatusOFF
	// VipPriceConfigStatusFuture 待生效
	VipPriceConfigStatusFuture
)

// vip_price_config suit_type
const (
	AllUser int8 = iota
	OldVIP
	NewVIP
	OldSubVIP
	NewSubVIP
	OldPackVIP
	NewPackVIP
)

// const .
const (
	DefualtZeroTimeFromDB = 0
	TimeFormatDay         = "2006-01-02 15:04:05"
	DefulatTimeFromDB     = "1970-01-01 08:00:00"
)

// VipPriceConfig struct .
type VipPriceConfig struct {
	ID          int64                `gorm:"column:id" json:"id"`
	Plat        VipPriceConfigPlat   `gorm:"column:platform" json:"platform"`
	PdName      string               `gorm:"column:product_name" json:"product_name"`
	PdID        string               `gorm:"column:product_id" json:"product_id"`
	SuitType    int8                 `gorm:"column:suit_type" json:"suit_type"`
	Month       int16                `gorm:"column:month" json:"month"`
	SubType     int8                 `gorm:"column:suit_type" json:"sub_type"`
	OPrice      float64              `gorm:"column:original_price" json:"original_price"`
	NPrice      float64              `json:"now_price"`
	Selected    int8                 `gorm:"column:selected" json:"selected"`
	Remark      string               `gorm:"column:remark" json:"remark"`
	Status      VipPriceConfigStatus `gorm:"column:status" json:"status"`
	Operator    string               `gorm:"column:operator" json:"operator"`
	OpID        int64                `gorm:"column:oper_id" json:"oper_id"`
	Superscript string               `gorm:"column:superscript" json:"superscript"`
	CTime       time.Time            `gorm:"column:ctime" json:"ctime"`
	MTime       time.Time            `gorm:"column:mtime" json:"mtime"`
	StartBuild  int64                `json:"start_build"`
	EndBuild    int64                `json:"end_build"`
}

// VipPriceConfigV2 struct .
type VipPriceConfigV2 struct {
	ID            int64     `gorm:"column:id" json:"id" form:"id"`
	Platform      int64     `gorm:"column:platform" json:"platform" form:"platform"`
	ProductName   string    `gorm:"column:product_name" json:"product_name" form:"product_name"`
	ProductID     string    `gorm:"column:product_id" json:"product_id" form:"product_id"`
	SuitType      int64     `gorm:"column:suit_type" json:"suit_type" form:"suit_type"`
	Month         int64     `gorm:"column:month" json:"month" form:"month"`
	SubType       int64     `gorm:"column:sub_type" json:"sub_type" form:"sub_type"`
	OriginalPrice float64   `gorm:"column:original_price" json:"original_price" form:"original_price"`
	Selected      int8      `gorm:"column:selected" json:"selected" form:"selected"`
	Remark        string    `gorm:"column:remark" json:"remark" form:"remark"`
	Status        int8      `gorm:"column:status" json:"status" form:"status"`
	Operator      string    `gorm:"column:operator" json:"operator" form:"operator"`
	OperID        int64     `gorm:"column:oper_id" json:"oper_id" form:"oper_id"`
	Ctime         time.Time `gorm:"column:ctime" json:"ctime" form:"ctime"`
	Mtime         time.Time `gorm:"column:mtime" json:"mtime" form:"mtime"`
	Superscript   string    `gorm:"column:superscript" json:"superscript" form:"superscript"`
	StartBuild    int64     `json:"start_build"`
	EndBuild      int64     `json:"end_build"`
}

// VipDPriceConfig price discount config.
type VipDPriceConfig struct {
	DisID      int64     `json:"discount_id"`
	ID         int64     `json:"vpc_id"`
	PdID       string    `json:"product_id"`
	DPrice     float64   `json:"discount_price"`
	STime      time.Time `json:"stime"`
	ETime      time.Time `json:"etime"`
	Remark     string    `json:"remark"`
	Operator   string    `json:"operator"`
	OpID       int64     `json:"oper_id"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
	FirstPrice float64   `json:"first_price"`
}

// VipPriceDiscountConfigV2 table vip_price_discount_config_v2 struct .
type VipPriceDiscountConfigV2 struct {
	ID            int64     `gorm:"column:id" json:"id" form:"id"`
	VpcID         int64     `gorm:"column:vpc_id" json:"vpc_id" form:"vpc_id"`
	ProductID     string    `gorm:"column:product_id" json:"product_id" form:"product_id"`
	DiscountPrice float64   `gorm:"column:discount_price" json:"discount_price" form:"discount_price"`
	Stime         time.Time `gorm:"column:stime" json:"stime" form:"stime"`
	Etime         time.Time `gorm:"column:etime" json:"etime" form:"etime"`
	Remark        string    `gorm:"column:remark" json:"remark" form:"remark"`
	Operator      string    `gorm:"column:operator" json:"operator" form:"operator"`
	OperID        int64     `gorm:"column:oper_id" json:"oper_id" form:"oper_id"`
	Ctime         time.Time `gorm:"column:ctime" json:"ctime" form:"ctime"`
	Mtime         time.Time `gorm:"column:mtime" json:"mtime" form:"mtime"`
}

// ArgAddOrUpVipPrice .
type ArgAddOrUpVipPrice struct {
	ID          int64              `form:"id"`
	Plat        VipPriceConfigPlat `form:"platform" validate:"required"`
	PdName      string             `form:"product_name" validate:"required"`
	PdID        string             `form:"product_id" validate:"required"`
	Month       int16              `form:"month" validate:"required"`
	SubType     int8               `form:"sub_type"`
	SuitType    int8               `form:"suit_type"`
	OPrice      float64            `form:"original_price" validate:"required"`
	Remark      string             `form:"remark"`
	Operator    string             `form:"operator" validate:"required"`
	OpID        int64              `form:"oper_id" validate:"required"`
	Selected    int8               `form:"selected"`
	Superscript string             `form:"superscript"`
	StartBuild  int64              `form:"start_build"`
	EndBuild    int64              `form:"end_build"`
}

// ArgAddOrUpVipDPrice .
type ArgAddOrUpVipDPrice struct {
	DisID      int64     `form:"discount_id"`
	ID         int64     `form:"vpc_id" validate:"required"`
	PdID       string    `form:"product_id" validate:"required"`
	DPrice     float64   `form:"discount_price"`
	STime      time.Time `form:"stime" validate:"required"`
	ETime      time.Time `form:"etime"`
	Remark     string    `form:"remark"`
	Operator   string    `form:"operator" validate:"required"`
	OpID       int64     `form:"oper_id" validate:"required"`
	FirstPrice float64   `form:"first_price"`
}

// CheckProductID .
// func (vpc *VipPriceConfig) CheckProductID(arg *ArgAddOrUpVipDPrice) bool {
// 	return (vpc.Plat == PlatVipPriceConfigIOS || vpc.Plat == PlatVipPriceConfigIOSHD || vpc.Plat == PlatVipPriceConfigIphoneB) && arg.PdID == ""
// }

// ExistPlat .
// func (aavp *ArgAddOrUpVipPrice) ExistPlat() bool {
// 	return aavp.Plat == PlatVipPriceConfigOther ||
// 		aavp.Plat == PlatVipPriceConfigIOS ||
// 		aavp.Plat == PlatVipPriceConfigIOSHD ||
// 		aavp.Plat == PlatVipPriceConfigFriendsGift ||
// 		aavp.Plat == PlatVipPriceConfigInternational ||
// 		aavp.Plat == PlatVipPriceConfigIphoneB ||
// 		aavp.Plat == PlatVipPriceConfigCheck
// }

// ArgVipPriceID .
type ArgVipPriceID struct {
	ID int64 `form:"id" validate:"required"`
}

// ArgVipDPriceID .
type ArgVipDPriceID struct {
	DisID int64 `form:"discount_id" validate:"required"`
}

// ArgVipPrice .
type ArgVipPrice struct {
	Plat     VipPriceConfigPlat `form:"platform" default:"-1"`
	Month    int16              `form:"month" default:"-1"`
	SubType  int8               `form:"sub_type" default:"-1"`
	SuitType int8               `form:"suit_type" default:"-1"`
}
