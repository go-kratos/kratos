package model

import "go-common/library/time"

// ArgID .
type ArgID struct {
	ID int64 `form:"id" validate:"required"`
}

// ArgPointHistory .
type ArgPointHistory struct {
	Mid             int64 `form:"id"`
	ChangeType      int64 `form:"change_type"`
	StartChangeTime int64 `form:"start_change_time"`
	EndChangeTime   int64 `form:"end_change_time"`
	BatchID         int64 `form:"batch_id"`
	RelationID      int64 `form:"relation_id"`
}

// ArgIDExtra .
type ArgIDExtra struct {
	ID       int64 `form:"id" validate:"required"`
	Status   int8  `form:"status"  validate:"required"`
	Operator string
}

// ArgPage .
type ArgPage struct {
	Ps     int `form:"ps"`
	Pn     int `form:"pn"`
	Status int `form:"status"`
}

// ArgPoolID .
type ArgPoolID struct {
	PoolID int `form:"pool_id" validate:"required"`
}

// ArgReSource .
type ArgReSource struct {
	ID        int       `form:"id"`
	Increment int       `form:"increment"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
}

// ArgAddPrivilege add privilege.
type ArgAddPrivilege struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Title    string `json:"title" form:"title" validate:"required"`
	Explain  string `json:"explain" form:"explain" validate:"required"`
	Type     int8   `json:"type" form:"type"`
	LangType int8   `json:"lang_type" form:"lang_type"`
	Operator string `json:"operator"`
	WebLink  string `json:"web_link" form:"web_link"`
	AppLink  string `json:"app_link" form:"app_link"`
}

// ArgUpdatePrivilege update privilege.
type ArgUpdatePrivilege struct {
	ID       int64  `form:"id" validate:"required"`
	Name     string `json:"name" form:"name" validate:"required"`
	Title    string `json:"title" form:"title" validate:"required"`
	Explain  string `json:"explain" form:"explain" validate:"required"`
	Type     int8   `json:"type" form:"type"`
	Operator string `json:"operator"`
	WebLink  string `json:"web_link" form:"web_link"`
	AppLink  string `json:"app_link" form:"app_link"`
}

// ArgImage arg image.
type ArgImage struct {
	IconFileType     string
	IconBody         []byte
	IconGrayFileType string
	IconGrayBody     []byte
	WebImageFileType string
	WebImageBody     []byte
	AppImageFileType string
	AppImageBody     []byte
}

// ArgStatePrivilege def.
type ArgStatePrivilege struct {
	ID     int64 `form:"id" validate:"required"`
	Status int8  `form:"state"`
}

// ArgPivilegeID def.
type ArgPivilegeID struct {
	ID int64 `form:"id" validate:"required"`
}

// ArgOrder def.
type ArgOrder struct {
	AID int64 `form:"aid" validate:"required"`
	BID int64 `form:"bid" validate:"required"`
}

// ArgAddJointly arg add jointly.
type ArgAddJointly struct {
	Title     string `form:"title" validate:"required"`
	Content   string `form:"content"`
	StartTime int64  `form:"start_time" validate:"required"`
	EndTime   int64  `form:"end_time" validate:"required"`
	Link      string `form:"link" validate:"required"`
	IsHot     int8   `form:"is_hot" `
	Operator  string
}

// ArgModifyJointly arg modify jointly.
type ArgModifyJointly struct {
	ID        int64  `form:"id" validate:"required"`
	Title     string `form:"title" validate:"required"`
	Content   string `form:"content" validate:"required"`
	Link      string `form:"link" validate:"required"`
	IsHot     int8   `form:"is_hot" `
	StartTime int64  `form:"start_time" validate:"required"`
	EndTime   int64  `form:"end_time" validate:"required"`
	Operator  string
}

// ArgQueryJointly query jointly params .
type ArgQueryJointly struct {
	State int8 `form:"state" `
}

// ArgJointlyID .
type ArgJointlyID struct {
	ID int64 `form:"id" validate:"required"`
}

//ArgPayOrder qeury order.
type ArgPayOrder struct {
	Mid     int64  `form:"mid"`
	OrderNo string `form:"order_no"`
	Status  int8   `form:"status"`
	PN      int    `form:"pn" default:"1"`
	PS      int    `form:"ps" default:"20"`
}
