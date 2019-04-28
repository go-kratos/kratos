package model

import (
	"time"

	xtime "go-common/library/time"
)

// ArgID .
type ArgID struct {
	ID     int64
	Mid    int64
	RealIP string
}

// ArgName .
type ArgName struct {
	Name   string
	Mid    int64
	RealIP string
}

// ArgIDs .
type ArgIDs struct {
	IDs    []int64
	Mid    int64
	RealIP string
}

// ArgNames .
type ArgNames struct {
	Names  []string
	Mid    int64
	RealIP string
}

// ArgAid .
type ArgAid struct {
	Aid    int64
	Mid    int64
	RealIP string
}

// ArgSub .
type ArgSub struct {
	Mid    int64
	Vmid   int64
	Pn     int
	Ps     int
	Order  int
	RealIP string
}

// ArgAddSub .
type ArgAddSub struct {
	Mid    int64
	Tids   []int64
	Now    time.Time
	RealIP string
}

// ArgCancelSub .
type ArgCancelSub struct {
	Mid    int64
	Tid    int64
	Now    time.Time
	RealIP string
}

// ArgUpdateCustomSort .
type ArgUpdateCustomSort struct {
	Tids   string
	Mid    int64
	Type   int
	RealIP string
}

// ArgCustomSort .
type ArgCustomSort struct {
	Mid    int64
	Type   int
	Order  int
	Pn     int
	Ps     int
	RealIP string
}

const (
	// UpRole .
	UpRole = iota // up主角色
	// UserRole .
	UserRole // 普通用户
	// AdminRole .
	AdminRole // 管理员
)

const (
	// PicResType .
	PicResType = iota + 1 // 图文资源
)

// ArgBind .
type ArgBind struct {
	Oid    int64
	Mid    int64
	Type   int8
	Names  []string
	RealIP string
}

// ArgUserAdd .
type ArgUserAdd struct {
	Oid    int64
	Mid    int64
	Type   int8
	Name   string
	Role   int8
	RealIP string
}

// ArgUserDel .
type ArgUserDel struct {
	Oid    int64
	Tid    int64
	Type   int8
	Mid    int64
	Role   int8
	RealIP string
}

// ArgResTags .
type ArgResTags struct {
	Oids   []int64
	Type   int8
	Mid    int64
	RealIP string
}

// ArgChannelResource ArgChannelResource.
type ArgChannelResource struct {
	Tid        int64  `form:"tid"`
	Mid        int64  `form:"mid"`
	Plat       int32  `form:"plat"`
	LoginEvent int32  `form:"login_event"`
	RequestCNT int32  `form:"request_cnt"`
	DisplayID  int32  `form:"display_id"`
	From       int32  `form:"from"`
	Type       int32  `form:"type"`
	Build      int32  `form:"build"`
	Name       string `form:"tname"`
	Buvid      string `form:"buvid"`
	Channel    int32
	RealIP     string
}

// ArgChanneList arg channel list.
type ArgChanneList struct {
	ID     int64
	Mid    int64
	From   int32
	RealIP string
}

// ArgChannelCategories arg channel categories.
type ArgChannelCategories struct {
	From   int32
	RealIP string
}

// ArgDiscoverChanneList .
type ArgDiscoverChanneList struct {
	Mid    int64
	From   int32
	RealIP string
}

// ArgRecommandChannel .
type ArgRecommandChannel struct {
	Mid    int64
	From   int32
	RealIP string
}

// ArgResChannelCheck .
type ArgResChannelCheck struct {
	Oids   []int64
	Type   int32
	RealIP string
}

// ArgResChannel .
type ArgResChannel struct {
	Oids   []int64
	Type   int32
	Mng    int32
	Mid    int64
	RealIP string
}

const (
	// UserTag tag Type
	UserTag = int8(0) // 普通tag
	// UpTag tag Type
	UpTag = int8(1) // up主tag
	// OfficailClassifyTag tag Type
	OfficailClassifyTag = int8(2) // 官方-分类tag
	// OfficailContentTag tag Type
	OfficailContentTag = int8(3) // 官方-内容tag
	// OfficailActiveTag tag Type
	OfficailActiveTag = int8(4) // 官方-活动tag
)

// Tag .
type Tag struct {
	ID           int64      `json:"tag_id"`
	Name         string     `json:"tag_name"`
	Cover        string     `json:"cover"`
	HeadCover    string     `json:"head_cover"`
	Content      string     `json:"content"`
	ShortContent string     `json:"short_content"`
	Type         int8       `json:"type"`
	State        int8       `json:"state"`
	CTime        xtime.Time `json:"ctime"`
	MTime        xtime.Time `json:"-"`
	// tag count
	Count struct {
		View  int `json:"view"`
		Use   int `json:"use"`
		Atten int `json:"atten"`
	} `json:"count"`
	// subscriber
	IsAtten int8 `json:"is_atten"`
	// archive_tag
	Role      int8  `json:"-"`
	Likes     int64 `json:"likes"`
	Hates     int64 `json:"hates"`
	Attribute int8  `json:"attribute"`
	Liked     int8  `json:"liked"`
	Hated     int8  `json:"hated"`
}

// Sub .
type Sub struct {
	Tags  []*Tag `json:"tags"`
	Total int    `json:"total"`
}

// Tags .
type Tags []*Tag

func (t Tags) Len() int { return len(t) }

func (t Tags) Less(i, j int) bool {
	if t[i].Likes > t[j].Likes {
		return true
	} else if t[i].Likes == t[j].Likes {
		if t[i].Type > t[j].Type {
			return true
		} else if t[i].Type == t[j].Type {
			if (t[i].Role == 0 && t[j].Role != 0) || (t[i].Role == 2 && t[j].Role == 1) {
				return true
			}
			if t[i].Role == t[j].Role && t[i].Hates < t[j].Hates {
				return true
			}
			if t[i].Hates == t[j].Hates && t[i].CTime < t[j].CTime {
				return true
			}
		}
	}
	return false
}

func (t Tags) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
