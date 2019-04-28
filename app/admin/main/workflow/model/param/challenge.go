package param

import (
	"net/url"
	"strconv"

	"go-common/library/xstr"
)

// ChallengeListCommonParam .
type ChallengeListCommonParam struct {
	Business           int8     `form:"business" validate:"required"`
	IDs                []int64  `form:"ids,split"`
	Oids               []string `form:"oids,split"`
	ObjectMids         []int64  `form:"object_mids,split"`
	Mids               []int64  `form:"mids,split"`
	Gids               []int64  `form:"gids,split"`
	States             []int64  `form:"states,split"`
	TypeIDs            []int64  `form:"typeids,split"`
	Tids               []int64  `form:"tids,split" validate:"dive,gt=0"`
	Rounds             []int64  `form:"rounds,split"`
	AssigneeAdminIDs   []int64  `form:"assignee_adminids,split"`
	AssigneeAdminNames []string `form:"assignee_adminnames,split"`
	AdminIDs           []int64  `form:"adminids,split"`
	BusinessStates     []int64  `form:"business_states,split"`
	DispatchStates     []int64  `form:"dispatch_states,split"`
	Title              string   `form:"title"`
	Content            string   `form:"content"`
	AdminReply         string   `form:"admin_reply"`
	UserReply          string   `form:"user_reply"`
	CTimeFrom          string   `form:"ctime_from"`
	CTimeTo            string   `form:"ctime_to"`
	Order              string   `form:"order" default:"ctime"`
	Sort               string   `form:"sort_order" default:"desc"`
	PS                 int      `form:"ps" default:"50"`
	PN                 int      `form:"pn" default:"1"`
}

// ChallengeListV3Param .
type ChallengeListV3Param struct {
	Business           int8     `form:"business" validate:"required"`
	IDs                []int64  `form:"cid,split"`
	Oids               []string `form:"oid,split"`
	Mids               []int64  `form:"mid,split"`
	Gids               []int64  `form:"gid,split"`
	States             []int64  `form:"state,split"`
	TypeIDs            []int64  `form:"typeid,split"`
	Tids               []int64  `form:"tid,split"`
	Roles              []int64  `form:"role,split"`
	AssigneeAdminIDs   []int64  `form:"assignee_adminid,split"`
	AssigneeAdminNames []string `form:"assignee_admin_name,split"`
	AdminIDs           []int64  `form:"adminid,split"`
	AdminNames         []string `form:"admin_name,split"`
	BusinessStates     []int64  `form:"business_state,split"`
	KW                 []string `form:"kw,split"`
	KWField            []string `form:"kw_field,split"`
	CTimeFrom          string   `form:"ctime_from"`
	CTimeTo            string   `form:"ctime_to"`
	Order              string   `form:"order" default:"id"`
	Sort               string   `form:"sort" default:"desc"`
	PS                 int      `form:"ps" default:"50"`
	PN                 int      `form:"pn" default:"1"`
}

// ChallRstParam describe the reset request params to a challenge row
type ChallRstParam struct {
	Cid       int64  `form:"cid" json:"cid" validate:"required,min=1"`
	State     int8   `form:"state" json:"state" validate:"min=0,max=14"`
	AdminID   int64  `json:"adminid"`
	AdminName string `json:"admin_name"`
	Reason    string `form:"reason" json:"reason"`
	Business  int8   `form:"business" json:"business"`
}

// ChallUpParam describe the update request params to a challenge row
type ChallUpParam struct {
	Cid       int64  `form:"cid" json:"cid" validate:"required,min=1"`
	Tid       int64  `form:"tid" json:"tid"`
	Note      string `form:"note" json:"note"`
	AdminID   int64  `form:"adminid" json:"adminid"`
	AdminName string `json:"admin_name"`
	Business  int8   `form:"business" json:"business"`
	Role      int8   `form:"role" json:"role"`
}

// ChallResParam describe the set result request params to a challenge row
type ChallResParam struct {
	Cid       int64  `json:"cid" form:"cid" validate:"required,min=1"`
	State     int8   `json:"state" form:"state" validate:"min=0,max=14"`
	Reason    string `json:"reason" form:"reason"`
	AdminID   int64  `json:"adminid" form:"adminid"`
	AdminName string `json:"admin_name"`
}

// BatchChallResParam describe the set result request params to a set of challenges
type BatchChallResParam struct {
	Cids      []int64 `json:"cids" form:"cid,split" validate:"required,gt=0"`
	State     int8    `json:"state" form:"state" validate:"min=0,max=14"`
	Business  int8    `form:"business" json:"business"`
	Role      int8    `form:"role" json:"role"`
	AdminID   int64   `json:"adminid"`
	AdminName string  `json:"admin_name"`
	Reason    string  `json:"reason" form:"reason"`
}

// ChallSetParamV3 .
type ChallSetParamV3 struct {
	ID      []int64 `json:"id" form:"id,split" validate:"required,gt=0"`
	State   int8    `json:"state" form:"state" validate:"min=0,max=14"`
	AdminID int64   `json:"adminid"`
	Reason  string  `json:"reason" form:"reason"`
}

// BatchChallBusStateParam .
type BatchChallBusStateParam struct {
	Cids              []int64 `form:"cid,split" json:"cid" validate:"required,gt=0"`
	AssigneeAdminID   int64   `json:"assignee_admin_id"`
	AssigneeAdminName string  `json:"assignee_admin_name"`
	Business          int8    `form:"business"`
	Role              int8    `form:"role"`
	BusState          int8    `form:"business_state" json:"business_state" validate:"min=0,max=14"`
}

// EventParam is used to parse user request
type EventParam struct {
	Cid         int64  `json:"cid" form:"cid" validate:"required,min=1"`
	AdminID     int64  `json:"adminid" form:"adminid"`
	AdminName   string `json:"admin_name"`
	Content     string `json:"content" form:"content"`
	Attachments string `json:"attachments" form:"attachments"`
	Event       int8   `json:"event" form:"event" validate:"required,min=1"`
}

// BatchEventParam is used to parse user request
type BatchEventParam struct {
	Cids        []int64 `json:"cids,split" form:"cids,split" validate:"required,dive,gt=0"`
	AdminID     int64   `json:"adminid" form:"adminid"`
	AdminName   string  `json:"admin_name"`
	Content     string  `json:"content" form:"content"`
	Attachments string  `json:"attachments" form:"attachments"`
	Event       int8    `json:"event" form:"event" validate:"required,min=1"`
}

// ChallExtraParam describe the request params to update challenge extra data
type ChallExtraParam struct {
	Cid       int64                  `json:"cid" validate:"required,min=1"`
	AdminID   int64                  `json:"adminid" validate:"required,min=1"`
	AdminName string                 `json:"admin_name"`
	Extra     map[string]interface{} `json:"extra" validate:"required"`
}

// ChallExtraParamV3 .
type ChallExtraParamV3 struct {
	Cids      []int64 `json:"cid" form:"cid,split" validate:"required,dive,gt=0"`
	AdminID   int64   `json:"adminid"`
	AdminName string  `json:"admin_name"`
	Extra     string  `json:"extra" form:"extra" validate:"required"`
}

// BatchChallExtraParam describe the request params to batch update challenges extra data
type BatchChallExtraParam struct {
	Cids      []int64                `json:"cid" form:"cid" validate:"required,min=1"`
	Business  int8                   `json:"business" form:"business"`
	AdminID   int64                  `json:"adminid" validate:"required,min=1"`
	AdminName string                 `json:"admin_name"`
	Extra     map[string]interface{} `json:"extra" form:"extra" validate:"required"`
}

// BusChallsBusStateParam describe the request params to update business state of challenges in business
type BusChallsBusStateParam struct {
	Business     int8                   `json:"business" validate:"required,min=1"`
	Oid          int64                  `json:"oid" validate:"required,min=1"`
	AdminID      int64                  `json:"adminid" validate:"required,min=1"`
	BusState     int8                   `json:"business_state" validate:"min=0,max=14"`
	PreBusStates []int8                 `json:"pre_business_states" validate:"dive,gt=-1"`
	Extra        map[string]interface{} `json:"extra"`
}

// ValidComponent will verify the component field is valid
func (e *EventParam) ValidComponent() bool {
	if e.Cid > 0 &&
		e.AdminID > 0 &&
		e.Content != "" &&
		e.Event > 0 {
		return true
	}

	return false
}

// ValidComponent will verify the component field is valid
func (be *BatchEventParam) ValidComponent() bool {
	if len(be.Cids) > 0 &&
		be.AdminID > 0 &&
		be.Content != "" &&
		be.Event > 0 {
		return true
	}

	return false
}

// MessageParam is the model to send message to end user
type MessageParam struct {
	Type     string
	Source   int8
	DataType int8
	MC       string
	Title    string
	Context  string
	MidList  []int64
}

// Query method will serialize all conditions into a url.Values struct
func (mp *MessageParam) Query() (uv url.Values) {
	uv = url.Values{}

	uv.Set("type", mp.Type)
	uv.Set("source", strconv.Itoa(int(mp.Source)))
	uv.Set("data_type", strconv.Itoa(int(mp.DataType)))
	uv.Set("mc", mp.MC)
	uv.Set("title", mp.Title)
	uv.Set("context", mp.Context)
	uv.Set("mid_list", xstr.JoinInts(mp.MidList))

	return uv
}
