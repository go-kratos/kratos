package model

import (
	"go-common/library/time"
)

// ListParams query qa vide task list params
type ListParams struct {
	AuditStatus []int    `form:"auditStatus,split" validate:"omitempty,max=100"`
	TaskID      []int64  `form:"taskID,split" validate:"omitempty,max=100,dive,gt=0"`
	Keyword     []string `form:"keyword,split" validate:"omitempty,max=100"`
	UPGroup     []int64  `form:"upGroup,split" validate:"omitempty,max=100,dive,gt=0"`
	UID         []int64  `form:"uid,split" validate:"omitempty,max=100,dive,gt=0"`
	Limit       int      `form:"limit" validate:"omitempty,max=1000"`
	Seed        string   `form:"seed"`
	ArcTypeID   []int64  `form:"arcTypeid[]" validate:"omitempty,max=100,dive,gt=0"`
	TagID       []int64  `form:"tagID,split" validate:"omitempty,max=100,dive,gt=0"`
	State       int16    `form:"state"`
	CTimeFrom   string   `form:"ctimeFrom"`
	CTimeTo     string   `form:"ctimeTo"`
	FTimeFrom   string   `form:"ftimeFrom"`
	FTimeTo     string   `form:"ftimeTo"`
	FansFrom    int64    `form:"fansFrom"`
	FansTo      int64    `form:"fansTo"`
	Order       string   `form:"order" default:"id"`
	Sort        string   `form:"sort" default:"desc"`
	Ps          int      `form:"ps" default:"50" validate:"omitempty,gt=0,max=100"`
	Pn          int      `form:"pn" default:"1" validate:"omitempty,gt=0"`
}

//AddVideoParams add qa video task params
type AddVideoParams struct {
	OUID  int64  `json:"uid"`
	Oname string `json:"username"`
	VideoDetail
}

//QASubmitParams submit qa video task params
type QASubmitParams struct {
	ID           int64  `json:"id" form:"id" validate:"required,gt=0"`
	AuditStatus  int16  `json:"audit_status" form:"auditStatus"`
	Encoding     int32  `json:"encoding" form:"encoding"`
	Norank       int32  `json:"norank" form:"norank"`
	Nodynamic    int32  `json:"nodynamic" form:"nodynamic"`
	PushBlog     int32  `json:"push_blog" form:"push_blog"`
	Norecommend  int32  `json:"norecommend" form:"norecommend"`
	Nosearch     int32  `json:"nosearch" form:"nosearch"`
	OverseaBlock int32  `json:"oversea_block" form:"oversea_block"`
	TagID        int64  `json:"tag_id" form:"tagID" validate:"omitempty,gt=0"`
	ReasonID     int64  `json:"reason_id" form:"reasonID" validate:"omitempty,gt=0"`
	Reason       string `json:"reason" form:"reason"`
	Note         string `json:"note" form:"note"`
	QaTagID      int64  `json:"qa_tag_id" form:"qaTagid" validate:"required,gt=0"`
	QATag        string `json:"qa_tag" form:"qaTag" validate:"required"`
	QaNote       string `json:"qa_note" form:"qaNote"`
}

//TaskVideoItem qa vide task list item
type TaskVideoItem struct {
	ID          int64      `json:"id"`
	DetailID    int64      `json:"detail_id"`
	TaskID      int64      `json:"task_id"`
	TaskUTime   int64      `json:"task_utime"`
	CTime       string     `json:"ctime"`
	FTime       string     `json:"ftime"`
	AuditStatus int16      `json:"audit_status"`
	TagID       int64      `json:"audit_tagid"`
	MID         int64      `json:"mid"`
	UPName      string     `json:"up_name"`
	UPGroups    []int64    `json:"up_groups"`
	UPGroupList []*UPGroup `json:"up_group_list"`
	Fans        int64      `json:"fans"`
	ArcTitle    string     `json:"arc_title"`
	ArcTypeid   int64      `json:"arc_typeid"`
	UID         int64      `json:"uid"`
	User        *UserRole  `json:"user"`
	State       int16      `json:"state"`
	StateName   string     `json:"state_name"`
}

//Page page
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

//QAVideoList qa video task list return struct
type QAVideoList struct {
	Result []*TaskVideoItem `json:"result"`
	Page   Page             `json:"page"`
}

// VideoParam video struct
type VideoParam struct { // TODO: batch param
	ID            int64        `json:"id"`
	Filename      string       `json:"filename"`
	Aid           int64        `json:"aid"`
	Mid           int64        `json:"mid"`
	RegionID      int16        `json:"region_id"`
	VideoDesign   *VideoDesign `json:"video_design,omitempty"`
	Status        int16        `json:"status"`
	CTime         time.Time    `json:"ctime"`
	Cid           int64        `json:"cid,omitempty"`
	Title         string       `json:"title,omitempty"`
	Desc          string       `json:"desc,omitempty"`
	Index         int          `json:"index,omitempty"`
	SrcType       string       `json:"src_type,omitempty"`
	Playurl       string       `json:"playurl,omitempty"`
	FailCode      int8         `json:"failinfo,omitempty"`
	Duration      int64        `json:"duration,omitempty"`
	XcodeState    int8         `json:"xcode_state,omitempty"`
	Attribute     int32        `json:"attribute,omitempty"`
	Filesize      int64        `json:"filesize,omitempty"`
	WebLink       string       `json:"weblink,omitempty"`
	Resolutions   string       `json:"resolutions,omitempty"`
	Encoding      int8         `json:"encoding"`
	EncodePurpose string       `json:"encode_purpose,omitempty"`
	UID           int64        `json:"uid,omitempty"`
	TaskID        int64        `json:"task_id,omitempty"`
	Oname         string       `json:"oname,omitempty"`
	TagID         int64        `json:"tag_id,omitempty"`
	Reason        string       `json:"reason,omitempty"`
	ReasonID      int64        `json:"reject_reason_id,omitempty"`
	Note          string       `json:"note,omitempty"`
	Attrs         *AttrParam   `json:"attrs,omitempty"`
	Fans          int64        `json:"-"`
	CateID        int64        `json:"-"`
	UpFrom        int8         `json:"-"`
	TaskState     int8         `json:"-"`
	TypeID        int16        `json:"-"`
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
}

// VideoDesign mosaic and watermark
type VideoDesign struct {
	Mosaic    []*Mosaic    `json:"mosaic,omitempty"`
	WaterMark []*WaterMark `json:"watermark,omitempty"`
}

// Mosaic .
type Mosaic struct {
	X     int64 `json:"x" form:"mosaic[0][x]"`
	Y     int64 `json:"y" form:"mosaic[0][y]"`
	W     int64 `json:"w" form:"mosaic[0][w]"`
	H     int64 `json:"h" form:"mosaic[0][h]"`
	Start int64 `json:"start" form:"mosaic[0][start]"`
	End   int64 `json:"end" form:"mosaic[0][end]"`
}

// WaterMark .
type WaterMark struct {
	LOC   int8   `json:"loc,omitempty"`
	URL   string `json:"url,omitempty"`
	MD5   string `json:"md5,omitempty"`
	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
	X     int64  `json:"x,omitempty"`
	Y     int64  `json:"y,omitempty"`
}
