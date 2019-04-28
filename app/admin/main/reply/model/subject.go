package model

import (
	"encoding/json"
	"fmt"
	"go-common/library/ecode"
	"go-common/library/time"
)

const (
	// SubTypeArchive 稿件
	SubTypeArchive int32 = 1
	// SubTypeTopic 话题
	SubTypeTopic = 2
	// SubTypeDrawyoo 画站 NOTE Deprecated
	SubTypeDrawyoo = 3
	// SubTypeActivity 活动
	SubTypeActivity = 4
	// SubTypeLiveVideo 直播小视频
	SubTypeLiveVideo = 5
	// SubTypeForbiden 封禁信息
	SubTypeForbiden = 6
	// SubTypeNotice 公告信息
	SubTypeNotice = 7
	// SubTypeLiveAct 直播活动
	SubTypeLiveAct = 8
	// SubTypeActArc 活动稿件
	SubTypeActArc = 9
	// SubTypeLiveNotice 直播公告
	SubTypeLiveNotice = 10
	// SubTypeLivePicture 直播有文画
	SubTypeLivePicture = 11
	// SubTypeArticle 文章
	SubTypeArticle = 12
	// SubTypeTicket 票务
	SubTypeTicket = 13
	// SubTypeMusic 音乐
	SubTypeMusic = 14
	// SubTypeCredit 风纪委
	SubTypeCredit = 15
	// SubTypePgcCmt 点评
	SubTypePgcCmt = 16
	// SubTypeDynamic 动态
	SubTypeDynamic = 17
	// SubTypePlaylist 播单
	SubTypePlaylist = 18
	// SubTypeMusicList 音乐播单
	SubTypeMusicList = 19
	// SubTypeComicSeason 漫画部评论
	SubTypeComicSeason = 20
	// SubTypeComicEpisode 漫画话评论
	SubTypeComicEpisode = 21
	// SubStateNormal 正常
	SubStateNormal int32 = 0
	// SubStateForbid 关闭
	SubStateForbid = 1
	// SubStateMonitor 监控中 NOTE Deprecated to use attr 2
	SubStateMonitor = 2
	// SubAttrTopAdmin 管理员置顶
	SubAttrTopAdmin uint32 = 0
	// SubAttrTopUpper up主置顶
	SubAttrTopUpper = 1
	// SubAttrMonitor 监控中
	SubAttrMonitor = 2
	// SubAttrConfig 存在配置
	SubAttrConfig = 3
	// SubAttrAudit 先审后发
	SubAttrAudit = 4
	// SubAttrFrozen 冻结评论
	SubAttrFrozen = 5

	// SubAttrFolded  有折叠的根评论
	SubAttrFolded = 7

	// SujectAllow subject_allow
	SujectAllow = "subject_allow"
	// SujectForbid SujectForbid
	SujectForbid = "subject_forbid"
	// SujectFrozen SujectFrozen
	SujectFrozen = "subject_frozen"
	// SujectUnfrozenAllow SujectUnfrozenAllow
	SujectUnfrozenAllow = "subject_unfrozen_allow"
	// SujectUnfrozenForbid SujectUnfrozenForbid
	SujectUnfrozenForbid = "subject_unfrozen_forbid"
)

// Subject reply subject.
type Subject struct {
	ID     int64     `json:"-"`
	Oid    int64     `json:"oid"`
	Type   int32     `json:"type"`
	Mid    int64     `json:"mid"`
	Count  int32     `json:"count"`
	RCount int32     `json:"rcount"`
	ACount int32     `json:"acount"`
	State  int32     `json:"state"`
	Attr   uint32    `json:"attr"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"-"`
	Meta   string    `json:"meta"`
}

// HasFolded ...
func (s *Subject) HasFolded() bool {
	return s.AttrVal(AttrFolded) == AttrYes
}

// MarkHasFolded ...
func (s *Subject) MarkHasFolded() {
	s.AttrSet(AttrYes, AttrFolded)
}

// UnmarkHasFolded ...
func (s *Subject) UnmarkHasFolded() {
	s.AttrSet(AttrNo, AttrFolded)
}

// SubjectMeta SubjectMeta
type SubjectMeta struct {
	AdminTop int64 `json:"atop,omitempty"`
	UpperTop int64 `json:"utop,omitempty"`
}

// TopSet TopSet
func (s *Subject) TopSet(top int64, typ uint, act uint32) (err error) {
	var meta SubjectMeta
	if s.Meta != "" {
		err = json.Unmarshal([]byte(s.Meta), &meta)
		if err != nil {
			return
		}
	}
	if act == 1 {
		if typ == 0 {
			if meta.AdminTop == top {
				return fmt.Errorf("subject already have the same top")
			}
			meta.AdminTop = top
		} else {
			if meta.UpperTop == top {
				return fmt.Errorf("subject already have the same top")
			}
			meta.UpperTop = top
		}
	} else {
		if typ == 0 {
			meta.AdminTop = 0
		} else {
			meta.UpperTop = 0
		}
	}
	var content []byte
	content, err = json.Marshal(meta)
	if err != nil {
		return err
	}
	s.Meta = string(content)
	return
}

// AttrVal return val of subject'attr.
func (s *Subject) AttrVal(bit uint32) uint32 {
	return (s.Attr >> bit) & uint32(1)
}

// AttrSet set val of subject'attr.
func (s *Subject) AttrSet(v uint32, bit uint32) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}

// CheckSubState check subject state is ok.
func CheckSubState(state int32) (err error) {
	if state < SubStateNormal || state > SubStateForbid {
		err = ecode.ReplyIllegalSubState
	}
	return
}

// SubjectLog info
type SubjectLog struct {
	AdminID   int64  `json:"adminid"`
	AdminName string `json:"admin_name"`
	Oid       string `json:"oid"`
	Type      int32  `json:"type"`
	Remark    string `json:"remark"`
	Action    string `json:"action"`
	CTime     string `json:"ctime"`
	State     int32  `json:"state"`
}

// SubjectLogRes format
type SubjectLogRes struct {
	Logs  []*SubjectLog `json:"logs"`
	Page  Page          `json:"page"`
	Order string        `json:"order"`
	Sort  string        `json:"sort"`
}
