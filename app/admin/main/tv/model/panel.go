package model

import "go-common/library/time"

// TvPriceConfig is tv vip pay order
type TvPriceConfig struct {
	ID          int64     `form:"id" json:"id"`
	PID         int64     `form:"pid" json:"pid" gorm:"column:pid"`
	Platform    int8      `form:"platform" json:"platform" validate:"required"`
	ProductName string    `form:"product_name" validate:"required" json:"product_name"`
	ProductID   string    `form:"product_id" validate:"required" json:"product_id"`
	SuitType    int8      `form:"suit_type" json:"suit_type" `
	Month       int64     `form:"month" json:"month"`
	SubType     int8      `form:"sub_type" json:"sub_type" `
	Price       int64     `form:"price" json:"price"`
	Selected    int8      `form:"selected" json:"selected"`
	Remark      string    `form:"remark" json:"remark"`
	Status      int8      `form:"status" json:"status"`
	Superscript string    `form:"superscript" json:"superscript"`
	Operator    string    `form:"operator" json:"operator"`
	OperId      int64     `form:"oper_id" json:"oper_id"`
	Stime       time.Time `form:"stime" json:"stime"`
	Etime       time.Time `form:"etime" json:"etime"`
	Mtime       time.Time `json:"mtime"`
}

// TvPriceConfigResp is used show panel info
type TvPriceConfigResp struct {
	ID          int64           `form:"id" json:"id"`
	PID         int64           `form:"pid" json:"pid" gorm:"column:pid"`
	ProductName string          `form:"product_name" json:"product_name"`
	ProductID   string          `form:"product_id" json:"product_id"`
	SuitType    int8            `form:"suit_type" json:"suit_type"`
	Month       int64           `form:"month" json:"month"`
	SubType     int8            `form:"sub_type" json:"sub_type"`
	Price       int64           `form:"price" json:"price"`
	OriginPrice int64           `form:"original_price" json:"original_price"`
	Selected    int8            `form:"selected" json:"selected"`
	Remark      string          `form:"remark" json:"remark"`
	Status      int8            `form:"status" json:"status"`
	Superscript string          `form:"superscript" json:"superscript"`
	Operator    string          `form:"operator" json:"operator"`
	OperId      int64           `form:"oper_id" json:"oper_id"`
	Ctime       time.Time       `json:"ctime"`
	Mtime       time.Time       `json:"mtime"`
	Items       []TvPriceConfig `json:"item"`
}

// TvPriceConfigListResp is used to list in TV panel list
type TvPriceConfigListResp struct {
	ID          int64     `form:"id" json:"id"`
	PID         int64     `form:"pid" json:"pid" gorm:"column:pid"`
	ProductName string    `form:"product_name" json:"product_name"`
	ProductID   string    `form:"product_id" json:"product_id"`
	SuitType    int8      `form:"suit_type" json:"suit_type"`
	Month       int64     `form:"month" json:"month"`
	SubType     int8      `form:"sub_type" json:"sub_type"`
	Price       int64     `form:"price" json:"price"`
	OriginPrice int64     `form:"original_price" json:"original_price"`
	Selected    int8      `form:"selected" json:"selected"`
	Status      int8      `form:"status" json:"status"`
	Operator    string    `form:"operator" json:"operator"`
	OperId      int64     `form:"oper_id" json:"oper_id"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// RemotePanel YST product res
type RemotePanel struct {
	Product []Product `json:"data"`
	Result  struct {
		ResultCode string `json:"result_code"`
		ResultMsg  string `json:"result_msg"`
	} `json:"result"`
}

// Product YST product
type Product struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	Title           string `json:"title"`
	Price           int64  `json:"price"`
	ComboPkgID      string `json:"combo_pkg_id"`
	ComboDes        string `json:"combo_des"`
	VideoType       string `json:"video_type"`
	VodType         string `json:"vod_type"`
	ProductDuration string `json:"product_duration"`
	Contract        string `json:"contract"`
	SuitType        int8   `json:"suit_type"`
}

// TableName tv_price_config
func (*TvPriceConfig) TableName() string {
	return "tv_price_config"
}

// TableName tv_price_config
func (*TvPriceConfigListResp) TableName() string {
	return "tv_price_config"
}
