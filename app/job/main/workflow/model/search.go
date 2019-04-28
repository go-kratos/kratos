package model

// ChallSearchCommonRes .
type ChallSearchCommonRes struct {
	Page   *page                    `json:"page"`
	Result []*ChallSearchCommonData `json:"result"`
}

// ChallSearchCommonData .
type ChallSearchCommonData struct {
	ID int64 `json:"id"`
}

// Page .
type page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// GroupSearchCommonCond is the common condition model to send group search request
type GroupSearchCommonCond struct {
	Fields    []string
	Business  int8
	IDs       []int64
	Oids      []string
	Tids      []int64
	States    []int8
	Mids      []int64
	Rounds    []int64
	TypeIDs   []int64
	FID       []int64
	RID       []int8
	EID       []int64
	TagRounds []int64

	ReportMID []int64 // report_mid
	AuthorMID []int64 // mid

	KW       []string
	KWFields []string

	CTimeFrom string
	CTimeTo   string

	PN    int64
	PS    int64
	Order string
	Sort  string
}

// AppealSearchCond .
type AppealSearchCond struct {
	Fields        []string
	IDs           []int64
	Rids          []int32
	Tids          []int64
	Bid           []int
	Mids          []int64
	Oids          []int64
	AuditState    []int8
	TransferState []int8
	AssignState   []int8
	Weight        int64
	Degree        []int8
	AuditAdmin    []int32
	TransferAdmin []int32
	TypeIDs       []int64 // workflow_business table
	KW            []string
	KWFields      []string
	DTimeFrom     string
	DTimeTo       string
	TTimeFrom     string
	TTimeTo       string
	CTimeFrom     string
	CTimeTo       string
	MTimeFrom     string
	MTimeTo       string
	PN            int
	PS            int
	Order         string
	Sort          string
}

// AppealSearchRes .
type AppealSearchRes struct {
	Page   *page               `json:"page"`
	Result []*AppealSearchData `json:"result"`
}

// AppealSearchData .
type AppealSearchData struct {
	ID            int64 `json:"id"`
	Bid           int32 `json:"bid"`
	Tid           int32 `json:"tid"`
	Mid           int64 `json:"mid"`
	Oid           int64 `json:"oid"`
	AuditState    int8  `json:"audit_state"`
	TransferState int8  `json:"transfer_state"`
	AssignState   int8  `json:"assign_state"`
	TransferAdmin int   `json:"transfer_adminid"`
	Weight        int64 `json:"weight"`
}
