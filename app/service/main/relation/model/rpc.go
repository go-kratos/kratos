package model

// ArgMid mid
type ArgMid struct {
	Mid    int64 `json:"mid" form:"mid" validate:"required" params:"mid;Required;Min(1)"`
	RealIP string
}

// ArgSameFollowing is
type ArgSameFollowing struct {
	Mid1 int64 `json:"mid1" form:"mid1" validate:"required" params:"mid1;Required;Min(1)"`
	Mid2 int64 `json:"mid2" form:"mid2" validate:"required" params:"mid2;Required;Min(1)"`
}

// ArgMids mids
type ArgMids struct {
	Mids   []int64
	RealIP string
}

// ArgFollowing following
type ArgFollowing struct {
	Mid    int64
	Fid    int64 `json:"fid" form:"fid" validate:"required" params:"fid;Required;Min(1)"`
	Source uint8
	RealIP string
	Action int8
	Infoc  map[string]string
}

// ArgRelation relation
type ArgRelation struct {
	Mid, Fid int64
	RealIP   string
}

// ArgRelations relations
type ArgRelations struct {
	Mid    int64
	Fids   []int64
	RealIP string
}

// ArgTag tag
type ArgTag struct {
	Mid    int64
	Tag    string
	RealIP string
}

// ArgTagId tag id
type ArgTagId struct {
	Mid    int64
	TagId  int64
	RealIP string
}

// ArgTagDel tag del
type ArgTagDel struct {
	Mid    int64
	TagId  int64
	RealIP string
}

// ArgTagUpdate tag update
type ArgTagUpdate struct {
	Mid    int64
	TagId  int64
	New    string
	RealIP string
}

// ArgTagsMoveUsers tags move users
type ArgTagsMoveUsers struct {
	Mid         int64
	BeforeID    int64
	AfterTagIds string
	Fids        string
	RealIP      string
}

// ArgPrompt rpc promt arg.
type ArgPrompt struct {
	Mid   int64 `form:"mid" params:"mid"`
	Fid   int64 `form:"fid" validate:"required" params:"fid;Required;Min(1)"`
	Btype int8  `form:"btype" validate:"required,min=1" params:"btype;Required;Min(1)"`
}

// ArgAchieveGet is
type ArgAchieveGet struct {
	Award string `form:"award" validate:"required"`
	Mid   int64  `form:"mid" validate:"required"`
}

// ArgAchieve is
type ArgAchieve struct {
	AwardToken string `form:"award_token" validate:"required"`
}

// FollowerNotifySetting show the follower-notify setting state
type FollowerNotifySetting struct {
	Mid     int64 `json:"mid"`
	Enabled bool  `json:"enabled"`
}
