package model

import (
	"go-common/library/time"
)

const (
	// DESC 降序排列
	DESC = int32(0)
	// ASC 升序排列
	ASC = int32(1)

	// VerifyUnknown tag审核状态未知
	VerifyUnknown = int32(-1)
	// VerifyNone tag无需审核
	VerifyNone = int32(0)
	// VerifyWait tag待审核
	VerifyWait = int32(1)
	// VerifyDone tag已审核
	VerifyDone = int32(2)

	// StateUnknown tag状态未知
	StateUnknown = int32(-1)
	// StateNormal 正常
	StateNormal = int32(0)
	// StateDel 删除
	StateDel = int32(1)
	// StateShield 屏蔽
	StateShield = int32(2)

	// DefaultPageNum Default Page Num.
	DefaultPageNum = int32(1)
	// DefaultPagesize default page size.
	DefaultPagesize = int32(50)
	// DefaultSearchNum default tag search num.
	DefaultSearchNum = int32(10)

	// DefaultOrder DefaultOrder.
	DefaultOrder = "ctime"
	// DefaultSort DefaultSort.
	DefaultSort = "DESC"

	// TNameMaxLen tag name length .
	TNameMaxLen = 32

	// TagHot hot tag
	TagHot = int32(0)
	// TagSubmission submission tag
	TagSubmission = int32(1)

	// OperateAdd 操作增加
	OperateAdd = int32(0)
	// OperateDel 操作删除
	OperateDel = int32(1)

	//QuerryByResInfo 按照资源信息查询
	QuerryByResInfo = int32(0)
	// QuerryByLimitState 按照限制状态查询
	QuerryByLimitState = int32(1)

	// ResLimitNone 资源限制状态：无
	ResLimitNone = int32(0)

	// TypeUnknow TypeUnknow.
	TypeUnknow = int32(-1)
	// TypeUser 用户tag
	TypeUser = int32(0)
	// TypeUper up主tag
	TypeUper = int32(1)
	// TypeBiliClass 官方分类tag.
	TypeBiliClass = int32(2)
	// TypeBiliContent 官方内容tag.
	TypeBiliContent = int32(3)
	// TypeBiliActivity 官方活动tag.
	TypeBiliActivity = int32(4)

	// TagOperateYes TagOperateYes.
	TagOperateYes = int32(1)
	// TagOperateNO TagOperateNO.
	TagOperateNO = int32(0)
)

// Tag Tag.
type Tag struct {
	ID           int64     `json:"id"`            //主键
	Type         int32     `json:"type"`          //类型 0-用户tag 1-up主tag 2-官方分类tag 3-官方内容tag
	Name         string    `json:"name"`          //tag name
	Cover        string    `json:"cover"`         //封面地址
	Content      string    `json:"content"`       //tag简介
	Verify       int32     `json:"verify"`        //状态 0-无需审核 1-待审核 2-已审核
	Attr         int32     `json:"attr"`          // 属性： 0:锁定 1:锁定
	State        int32     `json:"state"`         //状态 0-正常 1-删除 2-屏蔽
	CTime        time.Time `json:"ctime"`         //创建时间
	MTime        time.Time `json:"mtime"`         //最后修改时间
	HeadCover    string    `json:"head_cover"`    //头图
	ShortContent string    `json:"short_content"` // 短评
}

// TagCount tag count.
type TagCount struct {
	Tid  int64 `json:"tag_id"`
	Bind int64 `json:"bind_count"`
	Sub  int64 `json:"sub_count"`
}

// MngSearchTag MngSearchTag.
type MngSearchTag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    int32  `json:"tag_type"` //类型 0-用户tag 1-up主tag 2-官方分类tag 3-官方内容tag
	Use     int64  `json:"use_count"`
	Atten   int64  `json:"atten_count"`
	State   int32  `json:"state"`        //状态 0-正常 1-删除 2-屏蔽
	Verify  int32  `json:"verify_state"` //状态 0-无需审核 1-待审核 2-已审核
	CTime   string `json:"ctime"`        //TODO  string ===> time.Time
	MTime   string `json:"mtime"`        //TODO  string ===> time.Time
}

// UpdateESearchTag UpdateESearchTag.
type UpdateESearchTag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	Type    *int32 `json:"tag_type,omitempty"` //类型 0-用户tag 1-up主tag 2-官方分类tag 3-官方内容tag
	Use     int64  `json:"use_count,omitempty"`
	Atten   int64  `json:"atten_count,omitempty"`
	State   *int32 `json:"state,omitempty"`        //状态 0-正常 1-删除 2-屏蔽
	Verify  *int32 `json:"verify_state,omitempty"` //状态 0-无需审核 1-待审核 2-已审核
	CTime   string `json:"ctime,omitempty"`        //TODO  string ===> time.Time
	MTime   string `json:"mtime,omitempty"`        //TODO  string ===> time.Time
}

// Page pagination
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// MngSearchTagList .
type MngSearchTagList struct {
	Result []*MngSearchTag `json:"result"`
	Page   *Page           `json:"page"`
}

// TagInfo tag info.
type TagInfo struct {
	ID      int64     `json:"id"`      //主键
	Type    int32     `json:"type"`    //类型 0-用户tag 1-up主tag 2-官方分类tag 3-官方内容tag
	Name    string    `json:"name"`    //tag name
	Cover   string    `json:"cover"`   //封面地址
	Content string    `json:"content"` //tag简介
	Verify  int32     `json:"verify"`  //状态 0-无需审核 1-待审核 2-已审核
	Attr    int32     `json:"attr"`    // 属性： 0:锁定 1:锁定
	State   int32     `json:"state"`   //状态 0-正常 1-删除 2-屏蔽
	Bind    int64     `json:"bind_count"`
	Sub     int64     `json:"sub_count"`
	CTime   time.Time `json:"ctime"` //创建时间
	MTime   time.Time `json:"mtime"` //最后修改时间
}

// ResTagLog Resource log.
type ResTagLog struct {
	ID     int64     `json:"id"`
	Oid    int64     `json:"oid"`
	Mid    int64     `json:"mid"`
	Tid    int64     `json:"tid"`
	Tname  string    `json:"tname"`
	Typ    int8      `json:"type"`
	Role   int8      `json:"role"`
	Action int8      `json:"action"`
	Remark string    `json:"remark"`
	State  int8      `json:"state"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"mtime"`
}

// DependServiceHost tag
type DependServiceHost struct {
	// MngSearchHost  string
	PlatformHost   string
	HotTagHost     string
	ArchiveHotHost string
	AccountHost    string
	MessageHost    string
}

// UpdateTag UpdateTag.
type UpdateTag struct {
	ID    int64 `json:"id"`
	State int32 `json:"state"`
	VSate int32 `json:"verify_state"`
}
