package search

import (
	"go-common/app/admin/main/workflow/model"
)

// search appid
const (
	GroupSrhComID = "workflow_group_common"
	ChallSrhComID = "workflow_chall_common"
)

// GroupSearchCommonCond is the common condition model to send group search request
type GroupSearchCommonCond struct {
	Fields       []string
	Business     int8
	IDs          []int64
	Oids         []string
	Tids         []int64
	States       []int8
	Mids         []int64
	Rounds       []int64
	TypeIDs      []int64
	FID          []int64
	RID          []int8
	EID          []int64
	TagRounds    []int64
	FirstUserTid []int64

	ReportMID []int64 // report_mid
	AuthorMID []int64 // mid

	KWPriority bool
	KW         []string
	KWFields   []string

	CTimeFrom string
	CTimeTo   string

	PN    int64
	PS    int64
	Order string
	Sort  string
}

// GroupSearchCommonData .
type GroupSearchCommonData struct {
	ID           int64 `json:"id"`
	Oid          int64 `json:"oid"`
	Mid          int64 `json:"mid"`
	TypeID       int64 `json:"typeid"`
	Eid          int64 `json:"eid"`
	FirstUserTid int64 `json:"first_user_tid"`
}

// GroupSearchCommonResp .
type GroupSearchCommonResp struct {
	Page   *model.Page              `json:"page"`
	Result []*GroupSearchCommonData `json:"result"`
}

// ChallSearchCommonCond is the common condition model to send challenge search request
type ChallSearchCommonCond struct {
	// Using int64 directly
	Fields             []string
	Business           int8
	IDs                []int64
	Gids               []int64
	Oids               []string
	Tids               []int64
	Mids               []int64
	ObjectMids         []int64
	Rounds             []int64
	TypeIDs            []int64
	AssigneeAdminIDs   []int64
	AssigneeAdminNames []string
	AdminIDs           []int64
	States             []int64
	BusinessStates     []int64

	CTimeFrom string
	CTimeTo   string
	KW        []string
	KWFields  []string
	Distinct  []string

	PN    int
	PS    int
	Order string
	Sort  string
}

// FormatState transform add state in queue into search cond
func (csc *ChallSearchCommonCond) FormatState() {
	for _, busState := range csc.BusinessStates {
		if busState == model.QueueBusinessStateBefore {
			csc.BusinessStates = append(csc.BusinessStates, model.QueueBusinessState)
		}
	}

	for _, st := range csc.States {
		if st == model.QueueStateBefore {
			csc.States = append(csc.States, model.QueueState)
		}
	}
}

// ChallSearchCommonData .
type ChallSearchCommonData struct {
	ID       int64       `json:"id"`
	Oid      int64       `json:"oid"`
	Mid      int64       `json:"mid"`
	Gid      int64       `json:"gid"`
	Tid      int64       `json:"tid"`
	CountTid int64       `json:"count_tid"`
	State    interface{} `json:"state"` //兼容 int string
	Title    string      `json:"title"`
	Business int8        `json:"business"`
	TypeID   int64       `json:"typeid"`
	CTime    string      `json:"ctime"`
}

// ChallSearchCommonResp .
type ChallSearchCommonResp struct {
	Page   *model.Page              `json:"page"`
	Result []*ChallSearchCommonData `json:"result"`
}

// ChallReleaseUpSearchCond .
type ChallReleaseUpSearchCond struct {
	Cids            []int64 `json:"id"`
	AssigneeAdminID int64   `json:"assignee_adminid"`
	BusinessState   int64   `json:"business_state"`
}

// ChallReleaseUpSearchCondItem .
type ChallReleaseUpSearchCondItem struct {
	Cid             int64 `json:"id"`
	AssigneeAdminID int64 `json:"assignee_adminid"`
	BusinessState   int64 `json:"business_state"`
}

// ChallUpSearchResult .
type ChallUpSearchResult struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`
}
