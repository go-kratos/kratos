package model

// ArgPointAdd .
type ArgPointAdd struct {
	Mid        int64   `json:"mid" form:"mid" validate:"required,min=1,gte=1"`
	ChangeType int     `json:"change_type" form:"change_type" validate:"required"`
	RelationID string  `json:"relation_id" form:"relation_id"`
	Bcoin      float64 `json:"bcoin" form:"bcoin" validate:"required"`
	Remark     string  `json:"remark" form:"remark"`
	OrderID    string  `json:"order_id" form:"order_id" validate:"required"`
}

// ArgMid .
type ArgMid struct {
	Mid int64 `form:"mid" validate:"required,min=1,gte=1"`
}

//ArgPointHistory .
type ArgPointHistory struct {
	Cursor int `form:"cursor"`
	PS     int `form:"ps"`
	PN     int `form:"pn"`
}

//ArgOldPointHistory .
type ArgOldPointHistory struct {
	Mid int64 `form:"mid"`
	PS  int   `form:"ps"`
	PN  int   `form:"pn"`
}

// ArgConfig biz config.
type ArgConfig struct {
	Mid        int64   `form:"mid" validate:"required,min=1,gte=1"`
	Bp         float64 `form:"bp"`
	ChangeType int8    `form:"change_type"`
}

// ArgPointConsume .
type ArgPointConsume struct {
	Mid        int64  `form:"mid" validate:"required,min=1,gte=1"`
	ChangeType int64  `form:"change_type" validate:"required,min=1,gte=1"`
	RelationID string `form:"relation_id"`
	Point      int64  `form:"point"`
	Remark     string `form:"remark"`
}

// ArgPoint .
type ArgPoint struct {
	Mid        int64  `form:"mid" validate:"required,min=1,gte=1"`
	ChangeType int64  `form:"change_type" validate:"required,min=1,gte=1"`
	Point      int64  `form:"point"`
	Remark     string `form:"remark"`
	Operator   string `form:"operator"`
}
