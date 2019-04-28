package blocked

// ArgBlockedSearch param struct.
type ArgBlockedSearch struct {
	Keyword       string `form:"keyword" default:"-"`
	UID           int64  `form:"mid" default:"-100"`
	OPID          int64  `form:"op_id" default:"-100"`
	OriginType    int8   `form:"origin_type" default:"-100"`
	BlockedType   int8   `form:"blocked_type" default:"-100"`
	PublishStatus int8   `form:"publish_status" default:"-100"`
	Start         string `form:"start"`
	End           string `form:"end"`
	PN            int    `form:"pn" default:"1"`
	PS            int    `form:"ps" default:"50"`
	Order         string `form:"order" default:"id"`
	Sort          string `form:"sort" default:"desc"`
}

// ArgPublishSearch arg publish search
type ArgPublishSearch struct {
	Keyword  string `form:"keyword" default:"-"`
	Type     int8   `form:"type" default:"-100"`
	ShowFrom string `form:"start"`
	ShowTo   string `form:"end"`
	Order    string `form:"order" default:"id"`
	Sort     string `form:"sort" default:"desc"`
	PN       int    `form:"pn" default:"1"`
	PS       int    `form:"ps" default:"50"`
}

// ArgCaseSearch struct
type ArgCaseSearch struct {
	Keyword    string `form:"keyword" default:"-"`
	OriginType int8   `form:"origin_type" default:"-100"`
	Status     int8   `form:"status" default:"-100"`
	CaseType   int8   `form:"case_type" default:"-100"`
	UID        int64  `form:"uid" default:"-100"`
	OPID       int64  `form:"op_id" default:"-100"`
	TimeFrom   string `form:"start"`
	TimeTo     string `form:"end"`
	Order      string `form:"order" default:"id"`
	Sort       string `form:"sort" default:"desc"`
	PN         int    `form:"pn" default:"1"`
	PS         int    `form:"ps" default:"50"`
}

// ArgJurySearch struct
type ArgJurySearch struct {
	UID         int64  `form:"mid" default:"-100"`
	Status      int8   `form:"status" default:"-100"`
	Black       int8   `form:"type" default:"-100"`
	ExpiredFrom string `form:"start"`
	ExpiredTo   string `form:"end"`
	Order       string `form:"order" default:"id"`
	Sort        string `form:"sort" default:"desc"`
	PN          int    `form:"pn" default:"1"`
	PS          int    `form:"ps" default:"50"`
}

// ArgAddJurys struct
type ArgAddJurys struct {
	MIDs []int64 `form:"mids,split" validate:"required"`
	OPID int64   `form:"op_id" validate:"required"`
	Day  int     `form:"day" validate:"required"`
	Send int8    `form:"send" validate:"min=0,max=1"`
}

// ArgOpinionSearch struct
type ArgOpinionSearch struct {
	UID   int64  `form:"mid" default:"-100"`
	CID   int64  `form:"cid" default:"-100"`
	Vote  int    `form:"vote" default:"-100"`
	State int8   `form:"state" default:"-100"`
	Order string `form:"order" default:"id"`
	Sort  string `form:"sort" default:"desc"`
	PN    int    `form:"pn" default:"1"`
	PS    int    `form:"ps" default:"50"`
}

// ArgKpiPointSearch param struct.
type ArgKpiPointSearch struct {
	UID   int64  `form:"uid"  default:"-100"`
	Start string `form:"start" default:"-"`
	End   string `form:"end" default:"-"`
	Order string `form:"order" default:"id"`
	Sort  string `form:"sort" default:"desc"`
	PN    int    `form:"pn" default:"1"`
	PS    int    `form:"ps" default:"50"`
}

// ArgKpiSearch param struct.
type ArgKpiSearch struct {
	UID   int64  `form:"uid"  default:"0"`
	Start string `form:"start"`
	End   string `form:"end"`
	PN    int    `form:"pn" default:"1"`
	PS    int    `form:"ps" default:"20"`
}

// ArgPublish param struct.
type ArgPublish struct {
	ID            int64  `form:"id"`
	OID           int64  `form:"op_id" validate:"required"`
	PType         int8   `form:"publish_type" validate:"min=1,max=4"`
	PublishStatus int8   `form:"publish_status" validate:"min=0,max=1"`
	StickStatus   int8   `form:"stick_status"  validate:"min=0,max=1"`
	SubTitle      string `form:"sub_title"`
	Title         string `form:"title"`
	URL           string `form:"url"`
	Content       string `form:"content"`
	ShowTime      string `form:"show_time"`
}

// ArgCase param struct.
type ArgCase struct {
	ID            int64  `form:"id"`
	UID           int64  `form:"uid" validate:"required"`
	Otype         int8   `form:"origin_type" validate:"min=0,max=20"`
	ReasonType    int8   `form:"reason_type" validate:"min=0,max=40"`
	PunishResult  int8   `form:"punish_result" validate:"min=0,max=10"`
	BlockedDays   int    `form:"blocked_days"`
	OriginTitle   string `form:"origin_title" validate:"required"`
	OriginURL     string `form:"origin_url" validate:"required"`
	OriginContent string `form:"origin_content"`
	RelationID    string `form:"relation_id"`
	OID           int64  `form:"op_id" validate:"required"`
}

// ArgUpStatus param struct
type ArgUpStatus struct {
	IDS    []int64 `form:"ids,split" validate:"min=1,max=100"`
	OID    int64   `form:"op_id"  validate:"required"`
	Status int8    `form:"status"`
}

// ArgUpInfo param struct
type ArgUpInfo struct {
	ID      int64  `form:"id" validate:"required"`
	OID     int64  `form:"op_id" validate:"required"`
	Status  int8   `form:"status" validate:"min=0,max=1"`
	Content string `form:"content"`
}

// ArgCaseConf param struct
type ArgCaseConf struct {
	CaseGiveHours      int `form:"case_give_hours"  default:"0"`
	CaseCheckHours     int `form:"case_check_hours"  default:"0"`
	JuryVoteRadio      int `form:"jury_vote_radio"  default:"0"`
	CaseJudgeRadio     int `form:"case_judge_radio"  default:"0"`
	CaseVoteMin        int `form:"case_vote_min"  default:"0"`
	CaseObtainMax      int `form:"case_obtain_max"  default:"0"`
	CaseVoteMax        int `form:"case_vote_max"  default:"0"`
	JuryApplyMax       int `form:"jury_apply_max"  default:"0"`
	CaseLoadMax        int `form:"case_load_max" default:"0"`
	CaseLoadSwitch     int `form:"case_load_switch" default:"0"`
	CaseVoteMaxPercent int `form:"case_vote_max_percent" default:"0"`
	OID                int `form:"op_id"  validate:"required"`
}

// ArgAutoCaseConf param struct.
type ArgAutoCaseConf struct {
	ID          int64   `form:"id"`
	Platform    int8    `form:"platform"  validate:"required"`
	Reasons     []int64 `form:"reasons,split"`
	ReportScore int     `form:"report_score" default:"0"`
	Likes       int     `form:"likes" default:"0"`
	OID         int64   `form:"op_id"  validate:"required"`
}

// Pager param struct.
type Pager struct {
	Total int    `json:"total"`
	PN    int    `json:"page"`
	PS    int    `json:"pagesize"`
	Order string `json:"order"`
	Sort  string `json:"sort"`
}

// ArgVoteNum param struct.
type ArgVoteNum struct {
	OID   int64 `form:"op_id"   validate:"required"`
	RateS int8  `form:"rate_s" default:"1"`
	RateA int8  `form:"rate_a" default:"1"`
	RateB int8  `form:"rate_b" default:"1"`
	RateC int8  `form:"rate_c" default:"1"`
	RateD int8  `form:"rate_d" default:"1"`
}
