package model

// ArgMid is rpc mid params.
type ArgMid struct {
	Mid    int64
	RealIP string
}

// ArgVote is rpc vote params.
type ArgVote struct {
	Mid     int64
	Cid     int64 `form:"cid"  validate:"required"`
	RealIP  string
	Vote    int8    `form:"vote" validate:"min=1,max=4"`
	Attr    int8    `form:"attr" validate:"min=0,max=1" default:"0"`
	Content string  `form:"content"`
	Likes   []int64 `form:"likes,split"  validate:"min=0,max=20"`
	Hates   []int64 `form:"hates,split" validate:"min=0,max=20"`
	AType   int8    `form:"apply_type" default:"0"`
	AReason int8    `form:"apply_reason" default:"0"`
}

// ArgMidCid is rpc mid and cid params.
type ArgMidCid struct {
	Mid, Cid int64
	RealIP   string
}

// ArgCid is rpc cid params.
type ArgCid struct {
	Cid    int64 `form:"cid"`
	RealIP string
}

// ArgCaseList is rpc case list params.
type ArgCaseList struct {
	Mid    int64
	RealIP string
	Pn     int64
	Ps     int64
}

// ArgSetQs is rpc set question params.
type ArgSetQs struct {
	ID     int64
	Ans    int64
	Status int64
}

// ArgAns is rpc answer params.
type ArgAns struct {
	Mid    int64
	RealIP string
	Refer  string
	UA     string
	Buvid  string
	Ans    *LabourAns
}

// ArgOpinion is rpc opinion arg.
type ArgOpinion struct {
	Cid   int64 `form:"cid"   validate:"required"`
	PN    int64 `form:"pn" default:"1"`
	PS    int64 `form:"ps" validate:"min=0,max=10" default:"10"`
	IP    string
	Otype int8 `form:"otype" validate:"min=1,max=2" default:"1"`
}

// ArgID id.
type ArgID struct {
	ID int64
}

// ArgBlocked struct
type ArgBlocked struct {
	Otype int64
	Btype int64
	PS    int64
	PN    int64
}

// ArgAnnounce struct
type ArgAnnounce struct {
	Type int8
	PS   int64
	PN   int64
}
