package model

import (
	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

// SimpleArc is the simple struct of archive
type SimpleArc struct {
	ID      int   `gorm:"column:id"`
	AID     int64 `gorm:"column:aid"`
	MID     int   `gorm:"column:mid"`
	TypeID  int32 `gorm:"column:typeid"`
	Title   string
	Content string
	Cover   string
	Deleted int
	Result  int
	Valid   int
	Mtime   time.Time
	Pubtime time.Time
}

// Archive archive def. corresponding to our table structure
type Archive struct {
	ID         int       `gorm:"column:id" json:"id"`
	AID        int64     `gorm:"column:aid" json:"aid"`
	MID        int       `gorm:"column:mid" json:"mid"`
	TypeID     int32     `gorm:"column:typeid" json:"typeid"`
	Videos     int       `gorm:"column:videos" json:"videos"`
	Title      string    `gorm:"column:title" json:"title"`
	Cover      string    `gorm:"column:cover" json:"cover"`
	Content    string    `gorm:"column:content" json:"content"`
	Duration   int       `gorm:"column:duration" json:"duration"`
	Copyright  int       `gorm:"column:copyright" json:"copyright"`
	Pubtime    time.Time `gorm:"column:pubtime" json:"pubtime"`
	InjectTime time.Time `gorm:"column:inject_time" json:"inject_time"`
	Ctime      time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime      time.Time `gorm:"column:mtime" json:"mtime"`
	State      int       `gorm:"column:state" json:"state"`
	Manual     int       `gorm:"column:manual" json:"manual"`
	Valid      uint8     `gorm:"column:valid" json:"valid"`
	Submit     uint8     `gorm:"column:submit" json:"submit"`
	Retry      int       `gorm:"column:retry" json:"retry"`
	Result     uint8     `gorm:"column:result" json:"result"`
	Deleted    uint8     `gorm:"column:deleted" json:"deleted"`
	Reason     string    `gorm:"column:reason" json:"reason"`
}

// ArcPager is the result and page of archive query.
type ArcPager struct {
	Items []*ArcList `json:"items"`
	Page  *Page      `json:"page"`
}

// ArcListParam is archive list request params
type ArcListParam struct {
	ID     string `form:"id" json:"id"`
	Title  string `form:"title" json:"title"`
	CID    string `form:"cid" json:"cid"`
	Typeid int32  `form:"typeid" json:"typeid"`
	Valid  string `form:"valid" json:"valid"`
	Pid    int32  `form:"pid" json:"-"`
	Order  int    `form:"order" json:"order" default:"2"`
	Mid    int64  `form:"mid" json:"mid"`
	UpName string `form:"up_name"`
	PageCfg
}

// AddResp is for the response for adding archives/uppers
type AddResp struct {
	Succ     []int64 `json:"succ"`     // successfully added ids
	Exist    []int64 `json:"exist"`    // the ids already exist in our DB
	Invalids []int64 `json:"invalids"` // the invalid ids ( not exist in archives/uppers )
}

// ArcType arctype
type ArcType struct {
	ID   int16  `json:"id"`
	Pid  int16  `json:"pid"`
	Name string `json:"name"`
}

// ArcDB is the archive query result
type ArcDB struct {
	ArcCore
	Pubdate time.Time `gorm:"column:pubtime"`
}

// ArcCore is the archive core struct
type ArcCore struct {
	ID      string    `json:"id"`
	CID     string    `json:"cid" gorm:"column:aid"`
	TypeID  int32     `json:"typeid" gorm:"column:typeid"`
	Title   string    `json:"title"`
	Valid   string    `json:"valid" gorm:"column:valid"`
	Mtime   time.Time `json:"mtime"`
	Content string    `json:"content"`
	Cover   string    `json:"cover"`
	MID     int64     `json:"mid" gorm:"column:mid"`
}

// ArcList def.
type ArcList struct {
	ArcCore
	PTypeID int32  `json:"parent_typeid"`
	Pubdate string `json:"pubdate"`
	UpName  string `json:"up_name"`
}

// ToList def.
func (v *ArcDB) ToList(pid int32) (res *ArcList) {
	return &ArcList{
		ArcCore: v.ArcCore,
		PTypeID: pid,
		Pubdate: v.Pubdate.Time().Format(mysql.TimeFormat),
	}
}

// UgcType ugc archive category typelist
type UgcType struct {
	ID       int32      `json:"id"`
	Name     string     `json:"name"`
	Children []UgcCType `json:"children"`
}

// UgcCType ugc archive children category type
type UgcCType struct {
	Pid  int32  `json:"pid"`
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Category is for getting pid and name from archive category
type Category struct {
	Pid, Name string
}

// AvailTps structure in memory
type AvailTps struct {
	PassedTps []UgcType
	AllTps    []UgcType
}

// TableName ugc_archive
func (v ArcDB) TableName() string {
	return "ugc_archive"
}

// TableName ugc_archive
func (a SimpleArc) TableName() string {
	return "ugc_archive"
}

// TableName ugc_archive
func (a Archive) TableName() string {
	return "ugc_archive"
}
