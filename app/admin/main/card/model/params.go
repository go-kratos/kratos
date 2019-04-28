package model

// ArgQueryGroup query group arg.
type ArgQueryGroup struct {
	GroupID int64 `form:"group_id"`
	State   int8  `form:"state" default:"-1"`
}

// ArgQueryCards query cards arg.
type ArgQueryCards struct {
	GroupID int64 `form:"group_id"`
}

// ArgState update state.
type ArgState struct {
	ID    int64 `form:"id" validate:"required,min=1,gte=1"`
	State int8  `form:"state"`
}

// ArgID arg id.
type ArgID struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// ArgIds ids arg.
type ArgIds struct {
	Ids []int64 `form:"ids,split" validate:"min=1,max=50"`
}

// AddGroup add group arg.
type AddGroup struct {
	Name     string `form:"name" validate:"required" gorm:"column:name"`
	State    int8   `form:"state" gorm:"column:state"`
	Operator string `gorm:"column:operator"`
	OrderNum int64  `gorm:"column:order_num"`
}

// UpdateGroup update group arg.
type UpdateGroup struct {
	Name     string `form:"name" validate:"required" gorm:"column:name"`
	State    int8   `form:"state" gorm:"column:state"`
	Operator string `gorm:"column:operator"`
	ID       int64  `form:"id" validate:"required,min=1,gte=1"`
}

// AddCard add card arg.
type AddCard struct {
	Name            string `json:"name" gorm:"column:name" form:"name" validate:"required"`
	State           int32  `json:"state" gorm:"column:state" form:"state" `
	IsHot           int32  `json:"is_hot" gorm:"column:is_hot" form:"is_hot"`
	CardURL         string `json:"card_url" gorm:"column:card_url"`
	BigCradURL      string `json:"big_crad_url" gorm:"column:big_crad_url"`
	CardType        int32  `json:"card_type" gorm:"column:card_type" form:"card_type"`
	OrderNum        int64  `json:"order_num" gorm:"column:order_num"`
	Operator        string `json:"operator" gorm:"column:operator"`
	GroupID         int64  `json:"group_id" gorm:"column:group_id" form:"group_id" validate:"required"`
	CardFileType    string `gorm:"-"`
	CardBody        []byte `gorm:"-"`
	BigCardFileType string `gorm:"-"`
	BigCardBody     []byte `gorm:"-"`
}

// UpdateCard update card info.
type UpdateCard struct {
	ID              int64  `form:"id" validate:"required,min=1,gte=1"`
	Name            string `json:"name" gorm:"column:name" form:"name" validate:"required"`
	State           int32  `json:"state" gorm:"column:state" form:"state" `
	IsHot           int32  `json:"is_hot" gorm:"column:is_hot" form:"is_hot"`
	CardURL         string `json:"card_url" gorm:"column:card_url"`
	BigCradURL      string `json:"big_crad_url" gorm:"column:big_crad_url"`
	Operator        string `json:"operator" gorm:"column:operator"`
	CardFileType    string `gorm:"-"`
	CardBody        []byte `gorm:"-"`
	BigCardFileType string `gorm:"-"`
	BigCardBody     []byte `gorm:"-"`
}
