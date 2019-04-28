package archive

import (
	"go-common/app/admin/main/videoup/model/message"
	"go-common/library/time"
)

const (
	// ActionVideoSubmit 视频提交
	ActionVideoSubmit = "videoSubmit"
	// ActionArchiveSubmit 稿件提交
	ActionArchiveSubmit = "archiveSubmit"
	// ActionArchiveSecondRound 无稿件信息修改的补发二审消息
	ActionArchiveSecondRound = "archiveSecondRound"
	// ActionArchiveAttr 稿件属性修改
	ActionArchiveAttr = "archiveAttr"
	// ActionArchiveTypeID 稿件分区修改
	ActionArchiveTypeID = "archiveTypeID"
	//ActionArchiveTag 保存稿件标签
	ActionArchiveTag = "archiveTag"
	//ActionArchiveTagRecheck 保存稿件标签，并频道回查
	ActionArchiveTagRecheck = "archiveTagRecheck"
	//FromListChannelReview 提交来源列表 频道回查
	FromListChannelReview = "channel_review"
)

// VideoParam video struct
type VideoParam struct { // TODO: batch param
	ID            int64                `json:"id"`
	Filename      string               `json:"filename"`
	Aid           int64                `json:"aid"`
	Mid           int64                `json:"mid"`
	RegionID      int16                `json:"region_id"`
	VideoDesign   *message.VideoDesign `json:"video_design,omitempty"`
	Status        int16                `json:"status"`
	CTime         time.Time            `json:"ctime"`
	Cid           int64                `json:"cid,omitempty"`
	Title         string               `json:"title,omitempty"`
	Desc          string               `json:"desc,omitempty"`
	Index         int                  `json:"index,omitempty"`
	SrcType       string               `json:"src_type,omitempty"`
	Playurl       string               `json:"playurl,omitempty"`
	FailCode      int8                 `json:"failinfo,omitempty"`
	Duration      int64                `json:"duration,omitempty"`
	XcodeState    int8                 `json:"xcode_state,omitempty"`
	Attribute     int32                `json:"attribute,omitempty"`
	Filesize      int64                `json:"filesize,omitempty"`
	WebLink       string               `json:"weblink,omitempty"`
	Resolutions   string               `json:"resolutions,omitempty"`
	Encoding      int8                 `json:"encoding"`
	EncodePurpose string               `json:"encode_purpose,omitempty"`
	UID           int64                `json:"uid,omitempty"`
	TaskID        int64                `json:"task_id,omitempty"`
	Oname         string               `json:"oname,omitempty"`
	TagID         int64                `json:"tag_id,omitempty"`
	Reason        string               `json:"reason,omitempty"`
	ReasonID      int64                `json:"reject_reason_id,omitempty"`
	Note          string               `json:"note,omitempty"`
	Attrs         *AttrParam           `json:"attrs,omitempty"`
}

// ArcParam sencond round param
type ArcParam struct {
	Aid           int64       `json:"id"`
	Mid           int64       `json:"mid"`
	UID           int64       `json:"uid"`
	UName         string      `json:"uname"`
	CanCelMission bool        `json:"cancel_mission"`
	Cover         string      `json:"cover"`
	Source        string      `json:"source"`
	URL           string      `json:"redirecturl"`
	Forward       int64       `json:"forward"`
	PTime         time.Time   `json:"pubtime"`
	DTime         time.Time   `json:"delaytime"`
	CTime         time.Time   `json:"ctime"`
	Delay         bool        `json:"delay"`
	Tag           string      `json:"tag,omitempty"`
	IsUpBind      bool        `json:"is_up_bind"`
	SyncHiddenTag bool        `json:"sync_hidden_tag"`
	Copyright     int8        `json:"copyright"`
	FlagCopyright bool        `json:"flag_copyright"`
	Access        int16       `json:"access"`
	State         int8        `json:"state"`
	Round         int8        `json:"round"`
	Title         string      `json:"title,omitempty"`
	TypeID        int16       `json:"typeid"`
	Content       string      `json:"content"`
	Note          string      `json:"note"`
	Attrs         *AttrParam  `json:"attrs,omitempty"`
	Forbid        *ForbidAttr `json:"forbid"`
	Author        string      `json:"author"`
	RejectReason  string      `json:"reject_reason"`
	ReasonID      int64       `json:"reason_id"`
	ChangeDelay   bool        `json:"change_delay"`
	Notify        bool        `json:"notify"`
	NoEmail       bool        `json:"no_email"`
	ForceSync     bool        `json:"force_sync"`
	OnFlowID      int64       `json:"on_flow_id"`
	Dynamic       string      `json:"dynamic"`
	Porder
	UpNote string `json:"highrisk_note"`
	// AdminChange
	AdminChange   bool             `json:"admin_change"`
	PolicyID      int64            `json:"policy_id"`
	ApplyUID      int64            `json:"apply_uid"`
	FromList      string           `json:"from_list"`
	FlowAttribute map[string]int32 `json:"flow_attribute"`
}

//Porder table
type Porder struct {
	IndustryID   int64  `json:"industry_id"`
	BrandID      int64  `json:"brand_id"`
	BrandName    string `json:"brand_name"`
	Official     int8   `json:"official"`
	ShowType     string `json:"show_type"`
	ShowFront    int8   `json:"show_front"`
	Advertiser   string `json:"advertiser"`
	Agent        string `json:"agent"`
	GroupID      int64  `json:"group_id"`
	State        int8   `json:"state"`
	PorderAction string `json:"porder_action"`
}

//PorderConfig table
type PorderConfig struct {
	ID    int64  `json:"id"`
	Type  int8   `json:"type"`
	Name  string `json:"name"`
	State int8   `json:"state"`
	Code  string `json:"code"`
	Rank  int8   `json:"rank"`
}

// AttrParam bit
type AttrParam struct {
	NoRank      int32 `json:"no_rank,omitempty"`       // 0
	NoDynamic   int32 `json:"no_dynamic,omitempty"`    // 1
	NoWeb       int32 `json:"no_web,omitempty"`        // 2
	NoMobile    int32 `json:"no_mobile,omitempty"`     // 3
	NoSearch    int32 `json:"no_search,omitempty"`     // 4
	OverseaLock int32 `json:"oversea_block,omitempty"` // 5
	NoRecommend int32 `json:"no_recommend,omitempty"`  // 6
	NoReprint   int32 `json:"no_reprint,omitempty"`    // 7
	HasHD5      int32 `json:"is_hd,omitempty"`         // 8
	IsPGC       int32 `json:"is_pgc,omitempty"`        // 9
	AllowBp     int32 `json:"allow_bp,omitempty"`      // 10
	IsBangumi   int32 `json:"is_bangumi,omitempty"`    // 11
	IsPorder    int32 `json:"is_porder,omitempty"`     // 12
	LimitArea   int32 `json:"limit_area,omitempty"`    // 13
	AllowTag    int32 `json:"allow_tag,omitempty"`     // 14
	JumpURL     int32 `json:"is_jumpurl,omitempty"`    // 16
	IsMovie     int32 `json:"is_movie,omitempty"`      // 17
	BadgePay    int32 `json:"is_pay,omitempty"`        // 18
	PushBlog    int32 `json:"push_blog,omitempty"`     // 20
	ParentMode  int32 `json:"parent_mode,omitempty"`   // 21
	UGCPay      int32 `json:"ugcpay,omitempty"`        // 22
}

// IndexParam index_order.
type IndexParam struct {
	Aid       int64 `json:"aid"`
	ListOrder []*struct {
		ID    int64 `json:"id"`
		Index int   `json:"index"`
	} `json:"list_order"`
}

// MultSyncParam bath sync.
type MultSyncParam struct {
	Action     string      `json:"action"`
	VideoParam *VideoParam `json:"videoParam,omitempty"`
	ArcParam   *ArcParam   `json:"archiveParam,omitempty"`
}

// SyncAction sync action.
type SyncAction struct {
	Action string `json:"action"`
}

//TagParam update archive tag
type TagParam struct {
	AID               int64  `form:"aid"  validate:"required"`
	Tags              string `form:"tags"`
	FromChannelReview string `form:"channel_review"`
}

//BatchTagParam update batch archives' tag
type BatchTagParam struct {
	AIDs          []int64 `form:"aids,split" validate:"gt=0,dive,gt=0"`
	Action        string  `form:"action"`
	Tags          string  `form:"tags"`
	Note          string  `form:"note"`
	IsUpBind      bool    `form:"is_up_bind"`
	SyncHiddenTag bool    `form:"sync_hidden_tag"`
	FromList      string  `form:"from_list"`
}

//ChannelReviewInfo  频道回查检查
type ChannelReviewInfo struct {
	AID            int64
	ChannelIDs     string
	NeedReview     bool
	CanOperRecheck bool
}
