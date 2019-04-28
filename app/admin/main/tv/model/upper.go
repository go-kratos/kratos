package model

import "go-common/library/time"

// Upper corresponds to the structure of upper in our DB
type Upper struct {
	ID      int       `json:"id"`
	MID     int64     `json:"mid" gorm:"column:mid"`
	State   int       `json:"state"`
	Toinit  int       `json:"toinit"`
	Retry   int       `json:"retry"`
	Deleted int       `json:"deleted"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// UpperR corresponds to the structure of upper to show in front-end
type UpperR struct {
	MID   int64  `json:"mid"`
	State int    `json:"state"`
	Name  string `json:"name"`
	Ctime string `json:"ctime"`
	Mtime string `json:"mtime"`
}

// UpperPager def.
type UpperPager struct {
	Items []*UpperR `json:"items"`
	Page  *Page     `json:"page"`
}

// TableName ugc_uploader
func (a Upper) TableName() string {
	return "ugc_uploader"
}

// ImportResp is for the response for import uppers' videos
type ImportResp struct {
	NotExist []int64 `json:"not_exist"` // not existing uppers
	Succ     []int64 `json:"succ"`      // succesffuly updated ids
}

// ReqUpCms is the request structure of upcmsList
type ReqUpCms struct {
	Order int    `form:"order" validate:"required,min=3,max=4" default:"3"` // 3 = mtime Desc, 4 = mtime Asc
	Pn    int    `form:"pn" default:"1"`
	Name  string `form:"name"`
	MID   int64  `form:"mid"`
	Valid string `form:"valid"` // 0 = offline, 1 = online
}

// CmsUpper corresponds to the structure of upper for CMS in our DB
type CmsUpper struct {
	MID      int64     `json:"mid" gorm:"column:mid"`
	Mtime    time.Time `json:"-"`
	MtimeStr string    `json:"mtime" gorm:"-"`
	CmsName  string    `json:"cms_name"`
	OriName  string    `json:"ori_name"`
	CmsFace  string    `json:"cms_face"`
	Valid    int       `json:"valid"`
}

// ReqUpEdit is the request of up edit function
type ReqUpEdit struct {
	MID  int64  `form:"mid" validate:"required"`
	Name string `form:"name" validate:"required"`
	Face string `form:"face" validate:"required"`
}

// TableName ugc_uploader
func (a CmsUpper) TableName() string {
	return "ugc_uploader"
}

// CmsUpperPager is cms upper pager
type CmsUpperPager struct {
	Items []*CmsUpper `json:"items"`
	Page  *Page       `json:"page"`
}

// RespUpAudit is the response of up audit function
type RespUpAudit struct {
	Succ    []int64 `json:"succ"`
	Invalid []int64 `json:"invalid"`
}

// UpMC is upper info in MC
type UpMC struct {
	ID      int
	MID     int64 `gorm:"column:mid"`
	Toinit  int
	Submit  int    // 1=need report
	OriName string `gorm:"column:ori_name"` // original name
	CMSName string `gorm:"column:cms_name"` // cms intervened name
	OriFace string `gorm:"column:ori_face"` // original face
	CMSFace string `gorm:"column:cms_face"` // cms intervened face
	Valid   int    // auth info: 1=online,0=hidden
	Deleted int
}

// TableName ugc_uploader
func (a UpMC) TableName() string {
	return "ugc_uploader"
}
