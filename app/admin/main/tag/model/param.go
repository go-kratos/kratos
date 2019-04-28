package model

// ParamPage ParamPage.
type ParamPage struct {
	Pn int32 `form:"pn" default:"1"`
	Ps int32 `form:"ps" default:"50"`
}

// ESTag ESTag.
type ESTag struct {
	IDs     []int64 `form:"ids,split"`
	TagType int32   `form:"tag_type" default:"-1"`
	State   int32   `form:"tag_state" default:"-1"`
	Vstate  int32   `form:"verify_state" default:"-1"`
	Sort    string  `form:"sort" default:"desc"`
	Order   string  `form:"order" default:"ctime"`
	Keyword string  `form:"keyword"`
	ParamPage
}

// ParamTagEdit param tag edit.
type ParamTagEdit struct {
	TName   string `form:"name" validate:"required"`
	Content string `form:"content"`
	TP      int32  `form:"type" validate:"gte=0"`
	Tid     int64  `form:"tid"`
}

// ParamTagState param tag state.
type ParamTagState struct {
	Tid   int64 `form:"tid" validate:"required,gt=0"`
	State int32 `form:"state" validate:"gte=0,lte=2"`
}

// ParamRelationList RelationList.
type ParamRelationList struct {
	Type    int32  `form:"type" validate:"gte=0,lte=1"`
	TName   string `form:"tname"`
	Oid     int64  `form:"oid"`
	OidType int32  `form:"oid_type"`
	ParamPage
}

// ParamRelationAdd ParamRelationAdd.
type ParamRelationAdd struct {
	TName string `form:"tname" validate:"required"`
	Oid   int64  `form:"oid" validate:"required,gt=0"`
	Type  int32  `form:"type" validate:"required,gt=0"`
}

// ParamRelation ParamRelation.
type ParamRelation struct {
	Tid  int64 `form:"tid" validate:"required,gt=0"`
	Oid  int64 `form:"oid" validate:"required,gt=0"`
	Type int32 `form:"type" validate:"required,gt=0"`
}

// ParamHotList ParamHotList.
type ParamHotList struct {
	Type int32 `form:"type"  validate:"gte=0"`
	Rid  int64 `form:"rid"  validate:"required,gt=0"`
	Prid int64 `form:"prid"  validate:"required,gt=0"`
}

// ParamResLogList ParamResLogList.
type ParamResLogList struct {
	Oid    int64 `form:"oid" validate:"required,gt=0"`
	TP     int32 `form:"type" validate:"required,gt=0"`
	Role   int32 `form:"role" validate:"lte=2"`
	Action int32 `form:"action" validate:"lte=4"`
	ParamPage
}

// ParamResLogState ParamResLogState.
type ParamResLogState struct {
	ID    int64 `form:"id" validate:"required,gt=0"`
	Oid   int64 `form:"oid" validate:"required,gt=0"`
	TP    int32 `form:"type" validate:"required,gt=0"`
	State int32 `form:"state" validate:"lte=1,gte=0"`
}

// ParamResLimit ParamResLimit.
type ParamResLimit struct {
	Oid        int64 `form:"oid"`
	OidType    int32 `form:"oid_type"`
	TP         int32 `form:"type" validate:"gte=0,lte=1"`
	LimitState int32 `form:"state"`
	ParamPage
}

// ParamResLimitState ParamResLimitState.
type ParamResLimitState struct {
	Oid     int64 `form:"oid" validate:"required,gt=0"`
	Type    int32 `form:"type" validate:"required,gt=0"`
	Operate int32 `form:"operate" validate:"gte=0"`
}

// ParamSynonymList ParamSynonymList.
type ParamSynonymList struct {
	Keyword string `form:"keyword"`
	ParamPage
}

// ParamSynonymEdit ParamSynonymEdit.
type ParamSynonymEdit struct {
	TName  string  `form:"tname" validate:"required"`
	Adverb []int64 `form:"adverb,split" validate:"required,min=1,dive,gt=0"`
}

// ParamSynonymDel ParamSynonymDel.
type ParamSynonymDel struct {
	Tid    int64   `form:"tid" validate:"required,gt=0"`
	Adverb []int64 `form:"adverb,split"`
}

// ParamSynonymExist ParamSynonymExist.
type ParamSynonymExist struct {
	TName  string `form:"tname" validate:"required"`
	Adverb string `form:"adverb" validate:"required"`
}

// ParamReportList ParamReportList.
type ParamReportList struct {
	Audit  int32   `form:"audit" validate:"required,gt=0,lte=2"` //1:一审; 2:二审
	State  int32   `form:"state" validate:"gte=0,lte=4"`         //处理结果 0:一审未处理，1:二审已处理;2:二审不处理；3:二审未处理; 4:一审已处理
	Reason int32   `form:"reason"`                               //举报原因
	Mid    int64   `form:"mid"`                                  //mid：被举报人  rptMid：举报人
	RptMid int64   `form:"rpt_mid"`
	STime  string  `form:"stime"`
	ETime  string  `form:"etime"`
	Rid    []int64 `form:"rid,split"`
	Oid    int64   `form:"oid"`
	Type   int32   `form:"type"`
	TName  string  `form:"tname"` // tag名称.
	Tids   []int64 `form:"tids"`  // tag id.
	ParamPage
}

// ParamReportHandle ParamReportHandle.
type ParamReportHandle struct {
	ID     int64 `form:"id" validate:"required,gt=0"`
	Audit  int32 `form:"audit" validate:"required,gte=1,lte=2"` //1:一审; 2:二审
	Action int32 `form:"action" validate:"gte=0,lte=1"`
}

// ParamReportState ParamReportState.
type ParamReportState struct {
	ID    int64 `form:"id" validate:"required,gt=0"`
	State int32 `form:"state" validate:"gte=0,lte=4"` //处理结果 0:一审未处理，1:二审已处理;2:二审不处理；3:二审未处理; 4:一审已处理
}

// ParamReportPunish ParamReportPunish.
type ParamReportPunish struct {
	ID              []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	Audit           int32   `form:"audit" validate:"required,gte=1,lte=2"` //1:一审; 2:二审
	ReasonType      int32   `form:"reasontype"`
	Moral           int32   `form:"moral"`
	BlockTimeLength int32   `form:"blocktimelength"`
	Notify          int32   `form:"is_notify"`
	Reason          int32   `form:"reason"`
	Remark          string  `form:"remark"`
	Note            string  `form:"note"`
}

// ParamReportLog ParamReportLog.
type ParamReportLog struct {
	Oid        int64   `form:"oid"`
	Type       int32   `form:"type"`
	Tid        int64   `form:"tid"`
	Rid        int64   `form:"rid"`
	Mid        int64   `form:"mid"`
	STime      string  `form:"stime"`
	ETime      string  `form:"etime"`
	Username   string  `form:"username"`
	HandleType []int64 `form:"handle_type,split"`
	ParamPage
}

// ParamReport ParamReport.
type ParamReport struct {
	IDs   []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	Audit int32   `form:"audit" validate:"required,gte=1,lte=2"` //1:一审; 2:二审
}

// ParamChanneList ParamChanneList.
type ParamChanneList struct {
	ParamPage
	IDs       []int64 `form:"ids,split"`
	Type      int32   `form:"type"`
	Operator  string  `form:"operator"`
	Sort      string  `form:"sort"`
	Order     string  `form:"order"`
	State     int32   `form:"state"`
	STime     string  `form:"stime"`
	ETime     string  `form:"etime"`
	INTShield int32   `form:"int_shield" default:"-1"`
}
