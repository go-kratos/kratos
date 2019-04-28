package model

import (
	"fmt"
)

var (
	_upTpl   = `由#{%s}{"http://space.bilibili.com/%d"}在您的稿件《#{%s}{"http://www.bilibili.com/video/av%d"}》下创作的%s语言版本的字幕，已被管理员流转到%s。《#{查看详情}{"https://member.bilibili.com/v2#/zimu/my-zimu/zimu-editor?aid=%d&cid=%d&subtitleId=%d"}》）`
	_userTpl = `您在稿件《#{%s}{"http://www.bilibili.com/video/av%d"}》下创作的%s语言版本的字幕，已被管理员流转到%s状态。《#{查看详情}{"https://member.bilibili.com/v2#/zimu/my-zimu/zimu-editor?aid=%d&cid=%d&subtitleId=%d"}》）`

	// StatusContent .
	StatusContent = map[uint8]string{
		uint8(SubtitleStatusDraft):         "草稿",
		uint8(SubtitleStatusToAudit):       "待审核",
		uint8(SubtitleStatusAuditBack):     "驳回",
		uint8(SubtitleStatusRemove):        "删除",
		uint8(SubtitleStatusPublish):       "发布",
		uint8(SubtitleStatusManagerBack):   "管理员驳回",
		uint8(SubtitleStatusManagerRemove): "管理员删除",
	}
)

// WorkFlowSubtitleDispose .
const (
	WorkFlowSubtitleDisposeManagerBack   = int32(2) // 管理员回退
	WorkFlowSubtitleDisposeManagerDelete = int32(3) // 管理员删除
)

// SubtitleStatus .
type SubtitleStatus uint8

// SubtitleStatus
const (
	SubtitleStatusUnknown SubtitleStatus = iota
	SubtitleStatusDraft
	SubtitleStatusToAudit
	SubtitleStatusAuditBack
	SubtitleStatusRemove
	SubtitleStatusPublish
	SubtitleStatusCheckToAudit
	SubtitleStatusCheckPublish
	SubtitleStatusManagerBack
	SubtitleStatusManagerRemove
)

// WorkFlowSubtitleArg .
type WorkFlowSubtitleArg struct {
	Object  *WorkFlowObject   `json:"object"`
	Targets []*WorkFlowTarget `json:"targets"`
}

// WorkFlowObject .
type WorkFlowObject struct {
	Ids         []int64 `json:"id"`
	Business    int64   `json:"business"`
	Reason      string  `json:"reason"`
	DisposeMode int32   `json:"dispose_mode"`
}

// WorkFlowTarget .
type WorkFlowTarget struct {
	ID  int64 `json:"id"`
	Eid int64 `json:"eid"`
	Oid int64 `json:"oid"`
}

// SearchSubtitleResult result from search
type SearchSubtitleResult struct {
	Page   *Page
	Result []*struct {
		Oid int64 `json:"oid"`
		ID  int64 `json:"id"`
	} `json:"result"`
}

// SubtitleList .
type SubtitleList struct {
	Page      *Page             `json:"page"`
	Subtitles []*SearchSubtitle `json:"subtitles"`
}

// SearchSubtitle .
type SearchSubtitle struct {
	ID          int64  `json:"subtitle_id"`
	Oid         int64  `json:"oid"`
	Aid         int64  `json:"aid"`
	ArchiveName string `json:"archive_name"`
	VideoName   string `json:"video_name"`
	AuthorID    int64  `json:"author_id"`
	Lan         string `json:"lan"`
	LanDoc      string `json:"lan_doc"`
	Status      uint8  `json:"status"`
	IsSign      bool   `json:"is_sign"`
	IsLock      bool   `json:"is_lock"`
	SubtitleURL string `json:"subtitle_url"`
	Mtime       int64  `json:"mtime"`
}

// SubtitlePub .
type SubtitlePub struct {
	Oid        int64
	Type       int32
	Lan        uint8
	SubtitleID int64
	IsDelete   bool
}

// Subtitle .
type Subtitle struct {
	ID          int64          `json:"id"`
	Oid         int64          `json:"oid"`
	Type        int32          `json:"type"`
	Lan         uint8          `json:"lan"`
	Aid         int64          `json:"aid"`
	Mid         int64          `json:"mid"`
	AuthorID    int64          `json:"author_id"`
	UpMid       int64          `json:"up_mid"`
	IsSign      bool           `json:"is_sign"`
	IsLock      bool           `json:"is_lock"`
	Status      SubtitleStatus `json:"status"`
	PubTime     int64          `json:"-"`
	SubtitleURL string         `json:"subtitle_url"`
	Mtime       int64          `json:"mtime"`
}

// SubtitleSearchArg .
type SubtitleSearchArg struct {
	Aid      int64
	Oid      int64
	Mid      int64
	UpperMid int64
	Status   uint8
	Lan      uint8
	Ps       int32
	Pn       int32
}

// SubtitleArg .
type SubtitleArg struct {
	Aid      int64  `form:"aid"`
	Oid      int64  `form:"oid"`
	Mid      int64  `form:"mid"`
	UpperMid int64  `form:"up_mid"`
	Status   uint8  `form:"status"`
	Lan      string `form:"lan"`
	Ps       int32  `form:"ps" validate:"required"`
	Pn       int32  `form:"pn" validate:"required"`
}

// EditSubtitleArg .
type EditSubtitleArg struct {
	Oid          int64 `form:"oid" validate:"required"`
	SubtileID    int64 `form:"subtitle_id" validate:"required"`
	Status       uint8 `form:"status"`
	NotifyAuthor bool  `form:"notify_author"`
	NotifyUpper  bool  `form:"notify_upper"`
}

// SubtitleContext .
type SubtitleContext struct {
	// UpdateStatus  bool // 更新状态  无事务
	DraftCache    bool // 删除草稿缓存
	SubtitleCache bool // 删除字幕个体缓存 （必须）
	RebuildPub    bool // 触发到发布状态
	CheckHasDraft bool // 处罚到草稿状态
}

// Build .
func (sc *SubtitleContext) Build(origin, dst SubtitleStatus) {
	sc.SubtitleCache = true
	switch origin {
	case SubtitleStatusDraft, SubtitleStatusToAudit:
		sc.DraftCache = true
	case SubtitleStatusPublish:
		sc.RebuildPub = true
	}
	switch dst {
	case SubtitleStatusDraft, SubtitleStatusToAudit:
		sc.CheckHasDraft = true
	case SubtitleStatusPublish:
		sc.RebuildPub = true
	}
}

// SubtitleLans .
type SubtitleLans []*SubtitleLan

// SubtitleLan .
type SubtitleLan struct {
	Code  int64  `json:"code"`
	Lan   string `json:"lan"`
	DocZh string `json:"doc_zh"`
	DocEn string `json:"-"`
}

// GetByLan .
func (ss SubtitleLans) GetByLan(lan string) (code int64) {
	for _, s := range ss {
		if s.Lan == lan {
			return s.Code
		}
	}
	return 0
}

// GetByID .
func (ss SubtitleLans) GetByID(lanID int64) (lan string, doc string) {
	for _, s := range ss {
		if s.Code == lanID {
			return s.Lan, s.DocZh
		}
	}
	return
}

// NotifySubtitleUser .
type NotifySubtitleUser struct {
	Mid         int64
	Aid         int64
	Oid         int64
	SubtitleID  int64
	ArchiveName string
	LanDoc      string
	Status      string
}

// NotifySubtitleUp .
type NotifySubtitleUp struct {
	Mid         int64
	AuthorID    int64
	AuthorName  string
	Aid         int64
	Oid         int64
	SubtitleID  int64
	ArchiveName string
	LanDoc      string
	Status      string
}

// Msg .
func (ns *NotifySubtitleUp) Msg() string {
	return fmt.Sprintf(_upTpl, ns.AuthorName, ns.AuthorID, ns.ArchiveName, ns.Aid, ns.LanDoc, ns.Status, ns.Aid, ns.Oid, ns.SubtitleID)
}

// Msg .
func (ns *NotifySubtitleUser) Msg() string {
	return fmt.Sprintf(_userTpl, ns.ArchiveName, ns.Aid, ns.LanDoc, ns.Status, ns.Aid, ns.Oid, ns.SubtitleID)
}
