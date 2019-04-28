package param

// GroupParam describe the group param
type GroupParam struct {
	Oid       int64  `form:"oid" json:"oid" validate:"required,min=1"`
	Business  int8   `form:"business" json:"business" validate:"required,min=1"`
	Rid       int8   `form:"rid" json:"rid"`
	State     int8   `form:"state" json:"state"`
	Tid       int64  `form:"tid" json:"tid" validate:"required,min=1"`
	Note      string `form:"note" json:"note"`
	AdminID   int64  `json:"adminid"`
	AdminName string `json:"admin_name"`
}

// GroupRoleSetParam .
type GroupRoleSetParam struct {
	GID       []int64 `form:"id,split" json:"id" validate:"required"`
	AdminID   int64   `json:"admin_id"`
	AdminName string  `json:"admin_name"`
	BID       int8    `form:"bid" json:"bid" validate:"required,min=1"`
	RID       int8    `json:"rid"`
	TID       int64   `form:"tid" json:"tid" validate:"min=-1"`
	Note      string  `form:"note" json:"note"`
}

// GroupResParam describe the set result request params to a group row
type GroupResParam struct {
	Oid         int64       `form:"oid" json:"oid" validate:"required,min=1"`
	Business    int8        `form:"business" json:"business" validate:"required,min=1"`
	State       int8        `form:"state" json:"state" validate:"required,min=1,max=14"`
	AdminID     int64       `json:"adminid"`
	AdminName   string      `json:"admin_name"`
	Reason      string      `form:"reason" json:"reason"`
	ISDisplay   bool        `form:"is_display" json:"is_display"`
	IsMessage   bool        `form:"is_message" json:"is_message"`
	ReviewState int         `form:"review_state" json:"review_state"`
	Extra       *GroupExtra `json:"extra"`
}

// BatchGroupResParam describe the set result request params to a set of groups
type BatchGroupResParam struct {
	GID         []int64     `form:"id,split" json:"id"`
	Oids        []int64     `form:"oids,split" json:"oids" validate:"required,gt=0"`
	Business    int8        `form:"business" json:"business" validate:"required,min=1"`
	Role        int8        `form:"role" json:"role"`
	State       int8        `form:"state" json:"state" validate:"required,min=1,max=14"`
	AdminID     int64       `json:"adminid"`
	AdminName   string      `json:"admin_name"`
	Reason      string      `form:"reason" json:"reason"`
	ISDisplay   bool        `form:"is_display" json:"is_display"`
	IsMessage   bool        `form:"is_message" json:"is_message"`
	ReviewState int         `form:"review_state" json:"review_state"`
	Extra       *GroupExtra `json:"extra"`
}

// GroupExtra .
type GroupExtra struct {
	ISDisplay   bool `form:"is_display" json:"is_display"`
	IsMessage   bool `form:"is_message" json:"is_message"`
	ReviewState int  `form:"review_state" json:"review_state"`
}

// GroupListParamV3 .
type GroupListParamV3 struct {
	Business     int8     `form:"business" validate:"required"`
	Oid          []string `form:"oid,split"`
	Rid          []int8   `form:"rid,split"` //role
	Fid          []int64  `form:"fid,split"` //flow
	Eid          []int64  `form:"eid,split"`
	Mid          []int64  `form:"mid,split"`        // workflow_business mid
	ReportMid    []int64  `form:"report_mid,split"` // workflow_chall mid
	FirstUserTid []int64  `form:"first_user_tid"`
	State        []int8   `form:"state,split"`
	Tid          []int64  `form:"tid,split"`
	Round        []int64  `form:"round,split"`
	TypeID       []int64  `form:"typeid,split"`
	KWPriority   bool     `form:"kw_priority"`
	KW           []string `form:"kw,split"`
	KWField      []string `form:"kw_field,split"`
	Order        string   `form:"order" default:"lasttime"`
	Sort         string   `form:"sort" default:"desc"`
	PN           int64    `form:"pn" default:"1"`
	PS           int64    `form:"ps" default:"50"`
	AdminName    []string `form:"admin_name,split"`
	CTimeFrom    string   `form:"ctime_from"`
	CTimeTo      string   `form:"ctime_to"`
}

// GroupStateSetParam .
type GroupStateSetParam struct {
	ID            []int64 `form:"id,split" json:"id" validate:"required"`
	Business      int8    `form:"business" json:"business" validate:"required"`
	State         int8    `form:"state" json:"state" validate:"required"`
	Tid           int64   `form:"tid" json:"tid"` //处理理由tag_id
	Rid           int8    `form:"rid" json:"rid" validate:"required"`
	Reason        string  `form:"reason" json:"reason"`
	IsDisplay     bool    `form:"is_display" json:"is_display"`
	IsMessage     bool    `form:"is_message" json:"is_message"`           //通知举报人
	IsMessageUper bool    `form:"is_message_uper" json:"is_message_uper"` //通知被举报人(up主)
	ReviewState   int     `form:"review_state" json:"review_state"`
	DecreaseMoral int64   `form:"decrease_moral" json:"decrease_moral" validate:"max=0"` //扣节操
	DisposeMode   int     `form:"dispose_mode" json:"dispose_mode" validate:"min=0"`     //内容处理方式,批量操作不支持处理内容
	BlockDay      int64   `form:"block_day" json:"block_day"`                            //封禁时间
	BlockReason   int8    `form:"block_reason" json:"block_reason"`                      //封禁理由
	AdminID       int64   `json:"admin_id"`
	AdminName     string  `json:"admin_name"`
}

// GroupStatePublicReferee .
type GroupStatePublicReferee struct {
	ID        []int64 `form:"id,split" json:"id" validate:"required"`
	Business  int8    `form:"business" json:"business" validate:"required"`
	AdminID   int64   `json:"admin_id"`
	AdminName string  `json:"admin_name"`
	State     int8    `json:"state"`
}

// UpExtraParam describe the request params to batch update group extra data
type UpExtraParam struct {
	Gids      []int64 `form:"gid,split" json:"gid" validate:"required,min=1"`
	Extra     string  `form:"extra" json:"extra" validate:"required"`
	AdminID   int64   `json:"admin_id"`
	AdminName string  `json:"admin_name"`
}

// GroupPendingParam .
type GroupPendingParam struct {
	Business int8   `form:"business" json:"business" validate:"required,min=1"`
	Rid      []int8 `form:"rid,split" json:"rid"`
}
