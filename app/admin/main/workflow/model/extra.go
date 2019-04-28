package model

import (
	"net/url"

	xtime "go-common/library/time"
)

// ChannelInfoResult .
type ChannelInfoResult struct {
	*CommonResponse
	Data map[int64]*ArchiveChannelData
}

// ArchiveChannelData .
type ArchiveChannelData struct {
	Tids       []int64            `json:"tids"`
	ChannelMap map[int64]*Channel `json:"channels"`
	Archive    *Archive           `json:"-"`
}

// Channel .
type Channel struct {
	//Tid      int64    `json:"tid"`
	TName string `json:"tname"`
	//HitRule  string   `json:"hit_rule"`
	//HitNames []string `json:"hit_names"`
	HitRules []string `json:"hit_rules"`
	Archive  *Archive `json:"archive"`
}

// Review review struct
// http://info.bilibili.co/pages/viewpage.action?pageId=5417343
type Review struct {
	ReviewID  int64      `json:"review_id"`
	Type      int8       `json:"type"`
	Mid       int64      `json:"mid"`
	MediaName string     `json:"media_name"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Status    int8       `json:"status"`
	IsDelete  int8       `json:"is_delete"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"mtime"`
	TypeName  string     `json:"type_name"`
	Author    string     `json:"author"`
	Composer  *Account   `json:"-"`
}

// MultiReviewInfoResult response of multi review info api
type MultiReviewInfoResult struct {
	Code    int               `json:"code"`
	Result  map[int64]*Review `json:"result"`
	Message string            `json:"message"`
	TTL     int64             `json:"ttl"`
}

// SpecialTag .
type SpecialTag struct {
	MID       int64  `json:"mid"`
	GroupTag  string `json:"group_tag"`
	BGColor   string `json:"bg_color"`
	FontColor string `json:"font_color"`
	GroupID   int    `json:"group_id"`
}

// UperSpecial .
// http://info.bilibili.co/pages/viewpage.action?pageId=8479274
type UperSpecial struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int64  `json:"ttl"`
	Data    struct {
		Page  Page          `json:"page"`
		Items []*SpecialTag `json:"items"`
	} `json:"data"`
}

// Account account info struct
// http://info.bilibili.co/pages/viewpage.action?pageId=8465956#V3HTTP接口-用户全量信息含计数ProfileStat
// ignore unused field
type Account struct {
	Mid        int64         `json:"mid"`
	Name       string        `json:"name"`
	Rank       int32         `json:"rank"`
	Follower   int64         `json:"follower"`
	Official   *Official     `json:"official"`
	SpecialTag []*SpecialTag `json:"special_tag"`
}

// Official .
type Official struct {
	Role int8 `json:"role"`
}

// Relation .
// http://info.bilibili.co/pages/viewpage.action?pageId=8742464
type Relation struct {
	Follower int `json:"follower"`
}

// AccountInfosResult response of account info api
type AccountInfosResult struct {
	Code    int      `json:"code"`
	Data    *Account `json:"data"`
	Message string   `json:"message"`
	TTL     int64    `json:"ttl"`
}

// ArchiveComplainGroup is the model for archvie appeal group view
type ArchiveComplainGroup struct {
	Group
	Archive       *Archive      `json:"archive"`
	ChallengeTags ChallTagSlice `json:"challenge_tags"`
}

// ArchiveAppealGroup is the model for archvie appeal group view
type ArchiveAppealGroup struct {
	ArchiveComplainGroup

	// Appeal Related
	Cid               int64    `json:"cid"`
	AppealReason      string   `json:"appeal_reason"`
	AppealAttachments []string `json:"appeal_attachments"`
}

// ArchiveListPage is the model for archive list result
type ArchiveListPage struct {
	Items interface{} `json:"items"`
	Page  *Page       `json:"page"`
}

// Archive is archive info struct
// http://info.bilibili.co/pages/viewpage.action?pageId=3686646
type Archive struct {
	Author       string   `json:"author"`
	State        int32    `json:"state"`
	Mid          int64    `json:"mid"`
	TypeID       int32    `json:"tid"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	RejectReason string   `json:"reject_reason"`
	Composer     *Account `json:"-"`

	OPName    string `json:"op_name"`
	OPContent string `json:"op_content"`
	OPRemark  string `json:"op_remark"`
}

// ArchiveResult is response of archive info api
type ArchiveResult struct {
	*Archive `json:"archive"`
}

// MultiArchiveInfoResult is response of multi archive info api
type MultiArchiveInfoResult struct {
	Code    int                      `json:"code"`
	Data    map[int64]*ArchiveResult `json:"data"`
	Message string                   `json:"message"`
	TTL     int64                    `json:"ttl"`
}

// ArchiveType archive type struct
type ArchiveType struct {
	ID          int64  `json:"id"`
	PID         int64  `json:"pid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ArchiveTypeResult response of archive type api
type ArchiveTypeResult struct {
	Code    int                    `json:"code"`
	Data    map[int32]*ArchiveType `json:"data"`
	Message string                 `json:"message"`
	TTL     int64                  `json:"ttl"`
}

// FixAttachments will fix attachments url as user friendly
// ignore https case
// FIXME: this should be removed after attachment url is be normed
func (ag *ArchiveAppealGroup) FixAttachments() {
	if ag.AppealAttachments == nil {
		return
	}
	fixed := make([]string, 0, len(ag.AppealAttachments))
	for _, a := range ag.AppealAttachments {
		u, err := url.Parse(a)
		if err != nil {
			continue
		}
		u.Scheme = "http"
		fixed = append(fixed, u.String())
	}
	ag.AppealAttachments = fixed
}

// CreditMeta .
type CreditMeta struct {
	*Business
	Producer *Account `json:"producer"`
}

// BlockInfoResp .
type BlockInfoResp struct {
	*CommonResponse
	Data *struct {
		Mid        int64 `json:"mid"`
		Status     int64 `json:"status"`
		StartTime  int64 `json:"start_time"`
		EndTime    int64 `json:"end_time"`
		BlockedSum int64 `json:"blocked_sum"`
	} `json:"data"`
}

// BlockNumResp .
type BlockNumResp struct {
	*CommonResponse
	Data *struct {
		BlockedSum int64 `json:"blockedSum"`
	} `json:"data"`
}

// BlockCaseAdd .
// http://info.bilibili.co/pages/viewpage.action?pageId=5417571
type BlockCaseAdd struct {
	Aid           int64  `json:"aid"`
	Mid           int64  `json:"mid"`
	Operator      string `json:"operator"`
	OperID        int64  `json:"oper_id"`
	OriginType    int64  `json:"origin_type"`
	OriginContent string `json:"origin_content"`
	OriginTitle   string `json:"origin_title"`
	OriginURL     string `json:"origin_url"`
	ReasonType    int64  `json:"reason_type"`
	Oid           int64  `json:"oid"`
	RpID          int64  `json:"rp_id"`
	TagID         int64  `json:"tag_id"`
	Type          int64  `json:"type"`
	Page          int64  `json:"page"`
	BusinessTime  int64  `json:"business_time"`
}
