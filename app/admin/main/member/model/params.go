package model

import (
	"go-common/app/admin/main/member/model/block"
	xtime "go-common/library/time"
)

// ArgMid is.
type ArgMid struct {
	Mid int64 `form:"mid" validate:"min=1,required"`
}

// ArgMids is.
type ArgMids struct {
	Mid        []int64 `form:"mid,split" validate:"dive,gt=0"`
	Operator   string  `form:"operator"`
	OperatorID int64   `form:"operator_id"`
}

// ArgExpSet is.
type ArgExpSet struct {
	Mid        int64   `form:"mid" validate:"min=1"`
	Exp        float64 `form:"exp" validate:"required"`
	Reason     string  `form:"reason" validate:"required"`
	Operator   string  `form:"operator"`
	OperatorID int64   `form:"operator_id"`
	IP         string  `form:"ip"`
}

// ArgMoralSet is.
type ArgMoralSet struct {
	Mid        int64   `form:"mid" validate:"min=1"`
	Moral      float64 `form:"moral" validate:"required"`
	Reason     string  `form:"reason" validate:"required"`
	Operator   string  `form:"operator"`
	OperatorID int64   `form:"operator_id"`
	IP         string  `form:"ip"`
}

// ArgRankSet is.
type ArgRankSet struct {
	Mid        int64  `form:"mid" validate:"min=1"`
	Rank       int64  `form:"rank" validate:"required"`
	Reason     string `form:"reason" validate:"required"`
	Operator   string `form:"operator"`
	OperatorID int64  `form:"operator_id"`
	IP         string `form:"ip"`
}

// ArgCoinSet is.
type ArgCoinSet struct {
	Mid        int64   `form:"mid" validate:"min=1"`
	Coins      float64 `form:"coins" validate:"required"`
	Reason     string  `form:"reason" validate:"required"`
	Operator   string  `form:"operator"`
	OperatorID int64   `form:"operator_id"`
	IP         string  `form:"ip"`
}

// ArgAdditRemarkSet is.
type ArgAdditRemarkSet struct {
	Mid    int64  `form:"mid" validate:"min=1"`
	Remark string `form:"remark"`
}

// ArgBaseReview is.
type ArgBaseReview struct {
	Mid      []int64 `form:"mid,split"`
	StartMid int64   `form:"start_mid" validate:"min=0"`
	EndMid   int64   `form:"end_mid" validate:"min=0"`
}

// Mids mid list.
func (amr *ArgBaseReview) Mids() []int64 {
	mids := amr.Mid
	for i := amr.StartMid; i <= amr.EndMid; i++ {
		mids = append(mids, i)
	}
	return mids
}

// ArgList is.
type ArgList struct {
	Mid     int64  `form:"mid"`
	Keyword string `form:"keyword"`
	PN      int64  `form:"pn"`
	PS      int64  `form:"ps"`
}

// ArgOfficial is.
type ArgOfficial struct {
	Mid   int64      `form:"mid"`
	Role  []int64    `form:"role,split"`
	STime xtime.Time `form:"stime"`
	ETime xtime.Time `form:"etime"`
	Pn    int        `form:"pn"`
	Ps    int        `form:"ps"`
}

// ArgOfficialDoc is.
type ArgOfficialDoc struct {
	Mid   int64      `form:"mid"`
	Role  []int64    `form:"role,split"`
	State []int64    `form:"state,split"`
	STime xtime.Time `form:"stime"`
	ETime xtime.Time `form:"etime"`
	Uname string     `form:"uname"`
	Pn    int        `form:"pn"`
	Ps    int        `form:"ps"`
}

// ArgOfficialAudit is.
type ArgOfficialAudit struct {
	Mid        int64  `form:"mid" validate:"min=1"`
	State      int8   `form:"state" validate:"min=1"`
	UID        int64  `form:"uid" validate:"min=1"`
	Uname      string `form:"uname" validate:"min=1"`
	Reason     string `form:"reason"`
	Source     string `form:"source"`
	IsInternal bool   `form:"is_internal"`
}

// ArgOfficialEdit is.
type ArgOfficialEdit struct {
	Mid   int64  `form:"mid" validate:"min=1,required"`
	Role  int8   `form:"role" validate:"min=0"`
	Name  string `form:"name" validate:"gt=1,required"`
	Title string `form:"title" validate:"gt=1,required"`
	Desc  string `form:"desc"`

	// extra
	Telephone         string `form:"telephone"`
	Email             string `form:"email"`
	Address           string `form:"address"`
	Supplement        string `form:"supplement"`
	Company           string `form:"company"`
	Operator          string `form:"operator"`
	CreditCode        string `form:"credit_code"`
	Organization      string `form:"organization"`
	OrganizationType  string `form:"organization_type"`
	BusinessLicense   string `form:"business_license"`
	BusinessScale     string `form:"business_scale"`
	BusinessLevel     string `form:"business_level"`
	BusinessAuth      string `form:"business_auth"`
	OfficalSite       string `form:"official_site"`
	RegisteredCapital string `form:"registered_capital"`

	SendMessage    bool   `form:"send_msg"`
	MessageTitle   string `form:"msg_title"`
	MessageContent string `form:"msg_content"`

	UID   int64  `form:"uid" validate:"min=1"`
	Uname string `form:"uname" validate:"min=1"`

	IsInternal bool `form:"is_internal"`
}

// ArgOfficialSubmit arg submit official doc
type ArgOfficialSubmit struct {
	Mid   int64  `form:"mid"`
	Name  string `form:"name"`
	Role  int8   `form:"role"`
	Title string `form:"title"`
	Desc  string `form:"desc"`

	// extra
	Realname          int8   `form:"realname"`
	Operator          string `form:"operator"`
	Telephone         string `form:"telephone"`
	Email             string `form:"email"`
	Address           string `form:"address"`
	Company           string `form:"company"`
	CreditCode        string `form:"credit_code"`       // 社会信用代码
	Organization      string `form:"organization"`      // 政府或组织名称
	OrganizationType  string `form:"organization_type"` // 组织或机构类型
	BusinessLicense   string `form:"business_license"`  // 企业营业执照
	BusinessScale     string `form:"business_scale"`    // 企业规模
	BusinessLevel     string `form:"business_level"`    // 企业登记
	BusinessAuth      string `form:"business_auth"`     // 企业授权函
	Supplement        string `form:"supplement"`        // 其他补充材料
	Professional      string `form:"professional"`      // 专业资质
	Identification    string `form:"identification"`    // 身份证明
	OfficalSite       string `form:"official_site"`
	RegisteredCapital string `form:"registered_capital"`

	UID   int64  `form:"uid"`
	Uname string `form:"uname"`

	IsInternal   bool   `form:"is_internal"`
	SubmitSource string `form:"submit_source"`
}

// ArgFaceHistory is.
type ArgFaceHistory struct {
	Mid      int64      `form:"mid"`
	Operator string     `form:"operator"`
	Status   []int8     `form:"status,split"`
	STime    xtime.Time `form:"stime" validate:"min=0"`
	ETime    xtime.Time `form:"etime" validate:"min=0"`

	PS int `form:"ps" validate:"min=0,max=50"`
	PN int `form:"pn" validate:"min=0"`
}

// ArgMonitor is.
type ArgMonitor struct {
	Mid int64 `form:"mid"`
	Pn  int   `form:"pn"`
	Ps  int   `form:"ps"`
}

// ArgAddMonitor is.
type ArgAddMonitor struct {
	Mid        int64  `form:"mid" validate:"min=1,required"`
	Operator   string `form:"operator"`
	OperatorID int64  `form:"operator_id"`
	Remark     string `form:"remark"`
}

// ArgDelMonitor is.
type ArgDelMonitor struct {
	Mid        int64  `form:"mid" validate:"min=1,required"`
	Operator   string `form:"operator"`
	OperatorID int64  `form:"operator_id"`
	Remark     string `form:"remark"`
}

// ArgReviewList is.
type ArgReviewList struct {
	Mid       int64      `form:"mid"`
	Property  []int8     `form:"property,split"`
	Operator  string     `form:"operator"`
	State     []int8     `form:"state,split"`
	IsDesc    bool       `form:"is_desc"`
	IsMonitor bool       `form:"is_monitor"`
	ForceDB   bool       `form:"force_db"`
	STime     xtime.Time `form:"stime" validate:"min=0"`
	ETime     xtime.Time `form:"etime" validate:"min=0"`
	Ps        int        `form:"ps" validate:"min=0,max=50"`
	Pn        int        `form:"pn" validate:"min=0"`
}

// ArgReviewAudit is.
type ArgReviewAudit struct {
	ID         []int64 `form:"id,split" validate:"dive,gt=0"`
	State      int8    `form:"state" validate:"min=1"`
	Operator   string  `form:"operator"`
	OperatorID int64   `form:"operator_id"`
	Remark     string  `form:"remark"`
	BlockUser  bool    `form:"block_user"`
	//for block
	ArgBatchBlock
}

// ArgBatchBlock .
type ArgBatchBlock struct {
	Source   block.BlockMgrSource `form:"block_source"`
	Area     block.BlockArea      `form:"block_area"`
	Reason   string               `form:"block_reason"`
	Comment  string               `form:"block_comment"`
	Action   block.BlockAction    `form:"block_action"`
	Duration int64                `form:"block_duration"` // 单位：天
	Notify   bool                 `form:"block_notify"`
}

// Validate .
func (p *ArgBatchBlock) Validate() bool {
	// p.MIDs = intsSet(p.MIDs)
	// if len(p.MIDs) == 0 || len(p.MIDs) > 200 {
	// 	return false
	// }
	// if p.AdminID <= 0 {
	// 	return false
	// }
	// if p.AdminName == "" {
	// 	return false
	// }
	if p.Source != block.BlockMgrSourceSys && p.Source != block.BlockMgrSourceCredit {
		return false
	}
	if !p.Area.Contain() {
		return false
	}
	if p.Comment == "" {
		return false
	}
	if p.Action != block.BlockActionForever && p.Action != block.BlockActionLimit {
		return false
	}
	if p.Action == block.BlockActionLimit {
		if p.Duration <= 0 {
			return false
		}
	}
	return true
}

// ArgReview is.
type ArgReview struct {
	ID int64 `form:"id" validate:"min=1"`
}

// ArgPubExpMsg is.
type ArgPubExpMsg struct {
	Event string `form:"event" validate:"min=1,required"`
	Mid   int64  `form:"mid" validate:"min=1,required"`
	IP    string `form:"ip"`
	Ts    int64  `form:"ts"`
}

// Mode is.
func (a *ArgFaceHistory) Mode() string {
	if a.Mid > 0 && a.Operator != "" {
		return "op"
	}
	if a.Mid > 0 {
		return "mid"
	}
	return "op"
}

// ArgBatchFormal is
type ArgBatchFormal struct {
	FileData   []byte
	Operator   string `form:"operator"`
	OperatorID int64  `form:"operator_id"`
}

// ArgRealnameSubmit is
type ArgRealnameSubmit struct {
	Mid             int64  `form:"mid" validate:"required"`
	Realname        string `form:"realname" validate:"required"`
	CardType        int8   `form:"card_type"`
	CardNum         string `form:"card_num" validate:"required"`
	Country         int16  `form:"country"`
	FrontImageToken string `form:"front_image_token" validate:"required"`
	BackImageToken  string `form:"back_image_token" validate:"required"`
	HandImageToken  string `form:"hand_image_token"`

	Operator   string `form:"operator" validate:"required"`
	OperatorID int64  `form:"operator_id" validate:"required"`
	Remark     string `form:"remark" validate:"required"`
}
