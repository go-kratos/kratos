package model

import (
	"fmt"

	"go-common/library/time"
)

const (
	//ModulesNotDelete module not delete
	ModulesNotDelete = 0
	//ModulesDelete module delete
	ModulesDelete = 1
	//ModulesValid module is valid
	ModulesValid = 1
	//ModulesPublishYes module is publish status in MC
	ModulesPublishYes = 1
	//ModulesPublishNo module is not publish status in MC
	ModulesPublishNo = 0
	//PageMain 主页
	PageMain = 0
	//PageJP 追番
	PageJP = 1
	//PageMovie 电影
	PageMovie = 2
	//PageDocumentary 纪录片
	PageDocumentary = 3
	//PageCN 国创
	PageCN = 4
	//PageSoapopera 电视剧
	PageSoapopera = 5
	//TypeSevenFocus 首页七格焦点图
	TypeSevenFocus = 1
	//TypeFiveFocus 5格焦点
	TypeFiveFocus = 2
	//TypeSixFocus 6格焦点
	TypeSixFocus = 3
	//TypeVertListFirst 竖图1列表
	TypeVertListFirst = 4
	//TypeVertListSecond 竖图2列表
	TypeVertListSecond = 5
	//TypeHorizList 横图列表
	TypeHorizList = 6
	//TypeZhuiFan 追番模块
	TypeZhuiFan = 7
)

// Modules is use for Modular
type Modules struct {
	ID       uint64 `json:"id"`
	PageID   string `json:"page_id" form:"page_id" validate:"required"`
	Flexible string `json:"flexible" form:"flexible" validate:"required"`
	Icon     string `json:"icon" form:"icon"`
	Title    string `json:"title" form:"title" validate:"required"`
	Capacity uint64 `json:"capacity" form:"capacity" validate:"required"`
	More     string `json:"more" form:"more" validate:"required"`
	Order    uint8  `json:"order"`
	Moretype string `json:"moretype" form:"moretype"`
	Morepage int64  `json:"morepage" form:"morepage"`
	Deleted  uint8  `json:"-"`
	Valid    uint8  `json:"valid"`
	ModCore
}

// ModulesAddParam is use for Modular add param
type ModulesAddParam struct {
	ID       uint64 `form:"id" validate:"required"`
	PageID   string `form:"page_id" validate:"required"`
	Flexible string `form:"flexible" validate:"required"`
	Icon     string `form:"icon"`
	Title    string `form:"title" validate:"required"`
	Capacity uint64 `form:"capacity" validate:"required"`
	More     string `form:"more" validate:"required"`
	Moretype string `json:"moretype" form:"moretype"`
	Morepage int64  `json:"morepage" form:"morepage"`
	Order    uint8
	ModCore
}

// ModCore def.
type ModCore struct {
	Type    string `json:"type" form:"type" validate:"required"`
	Source  string `json:"source" form:"source" validate:"required"`
	SrcType int    `json:"src_type" form:"src_type" validate:"required"`
}

//ModPub is used for store publish status
type ModPub struct {
	Time  string
	State uint8
}

//ModulesList is used for function module list
type ModulesList struct {
	Items    []*Modules `json:"items"`
	PubState uint8      `json:"pubstate"`
	PubTime  string     `json:"pubtime"`
}

// TableName tv modules
func (a Modules) TableName() string {
	return "tv_modules"
}

//CommonCat , PGC types or ugc second level types
type CommonCat struct {
	ID   int32  `json:"id"`
	PID  int32  `json:"pid"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

//ParentCat : ugc first level types
type ParentCat struct {
	ID       int32        `json:"id"`
	Name     string       `json:"name"`
	Type     int          `json:"type"`
	Children []*CommonCat `json:"children,omitempty"`
}

//SupCats : support category map
type SupCats struct {
	UgcMap map[int32]int
	PgcMap map[int32]int
}

// AbnorCids is the export format for abnormal cids
type AbnorCids struct {
	CID        int64  `json:"cid"`
	VideoTitle string `json:"video_title"`
	CTime      string `json:"ctime"`
	AID        int64  `json:"aid"`
	ArcTitle   string `json:"arc_title"`
	PubTime    string `json:"pub_time"`
}

// Export transforms the structure to export csv data
func (v *AbnorCids) Export() (res []string) {
	res = append(res, fmt.Sprintf("%d", v.CID))
	res = append(res, v.VideoTitle)
	res = append(res, v.CTime)
	res = append(res, fmt.Sprintf("%d", v.AID))
	res = append(res, v.ArcTitle)
	res = append(res, v.PubTime)
	return
}

// AbnorVideo def.
type AbnorVideo struct {
	CID        int64
	VideoTitle string
	CTime      time.Time
	AID        int64
}

// ToCids transforms the archive & video to abnormal cid export structure
func (v *AbnorVideo) ToCids(arc *Archive) *AbnorCids {
	return &AbnorCids{
		CID:        v.CID,
		VideoTitle: v.VideoTitle,
		CTime:      v.CTime.Time().Format("2006-01-02 15:04:05"),
		AID:        v.AID,
		ArcTitle:   arc.Title,
		PubTime:    arc.Pubtime.Time().Format("2006-01-02 15:04:05"),
	}
}
