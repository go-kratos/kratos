package egg

import (
	"go-common/app/admin/main/feed/model/common"
	"go-common/library/time"
)

var (
	//NotDelete egg not deleted
	NotDelete uint8
	//Delete egg deleted
	Delete uint8 = 1
	//Publish egg publish
	Publish uint8 = 1
	//NotPublish egg not publish
	NotPublish uint8
	//Business log businessID
	Business = 201
)

//Obj add egg object
type Obj struct {
	Query     []string  `json:"query" form:"query,split" validate:"required"`
	Stime     time.Time `json:"stime" form:"stime"`
	Etime     time.Time `json:"etime" form:"etime"`
	ShowCount int       `json:"show_count" form:"show_count" validate:"required"`
	Plat      string    `json:"plat" form:"plat" validate:"required"`
}

//ObjUpdate add egg object
type ObjUpdate struct {
	ID        uint      `form:"id" validate:"required"`
	Query     []string  `json:"query" form:"query,split" validate:"required"`
	Stime     time.Time `json:"stime" form:"stime"`
	Etime     time.Time `json:"etime" form:"etime"`
	ShowCount int       `json:"show_count" form:"show_count" validate:"required"`
	Plat      string    `json:"plat" form:"plat" validate:"required"`
}

//Plat egg plat
type Plat struct {
	EggID      uint   `json:"egg_id"`
	Plat       uint8  `json:"plat"`
	Conditions string `json:"conditions"`
	Build      string `json:"build"`
	URL        string `json:"url"`
	Md5        string `json:"md5"`
	Size       uint   `json:"size"`
	Deleted    uint8  `json:"deleted"`
}

//Query egg query
type Query struct {
	EggID   uint
	Word    string
	STime   time.Time
	ETime   time.Time
	Deleted uint8
}

//Egg egg
type Egg struct {
	ID        uint
	Stime     time.Time
	Etime     time.Time
	ShowCount int
	UID       int64 `gorm:"column:uid"`
	Publish   uint8
	Person    string
	Delete    uint8
}

//IndexParam Index egg index param
type IndexParam struct {
	ID     string `json:"id" form:"id"`              // ID
	Stime  string `json:"stime" form:"stime"`        // 开始时间
	Etime  string `json:"etime" form:"etime"`        // 结束时间
	Person string `json:"person" form:"person"`      // 创建人
	Word   string `json:"word" form:"word"`          // 关键词
	Ps     int    `json:"ps" form:"ps" default:"20"` // 分页大小
	Pn     int    `json:"pn" form:"pn" default:"1"`  // 第几个分页
}

//Index egg index
type Index struct {
	ID        uint      `json:"id"`
	Words     string    `json:"words"`
	Stime     time.Time `json:"stime"`
	Etime     time.Time `json:"etime"`
	Plat      []Plat    `json:"plat"`
	ShowCount int       `json:"show_count"`
	Publish   uint8     `json:"publish"`
	Person    string    `json:"person"`
}

//IndexPager return values
type IndexPager struct {
	Item []*Index    `json:"item"`
	Page common.Page `json:"page"`
}

//SearchEgg for searching
type SearchEgg struct {
	ID    uint           `json:"id"`
	Words []string       `json:"query_list"`
	Stime time.Time      `json:"stime"`
	Etime time.Time      `json:"etime"`
	Plat  map[uint8]Plat `json:"plat"`
	//Plat      []Plat         `json:"plat"`
	ShowCount int   `json:"show_count"`
	Publish   uint8 `json:"publish"`
}

//SearchEggWeb for searching
type SearchEggWeb struct {
	ID    uint             `json:"id"`
	Words []string         `json:"query_list"`
	Stime time.Time        `json:"stime"`
	Etime time.Time        `json:"etime"`
	Plat  map[uint8][]Plat `json:"plat"`
	//Plat      []Plat         `json:"plat"`
	ShowCount int   `json:"show_count"`
	Publish   uint8 `json:"publish"`
}

// TableName Egg
func (a SearchEggWeb) TableName() string {
	return "egg"
}

// TableName Egg
func (a Egg) TableName() string {
	return "egg"
}

// TableName Egg plat
func (a Plat) TableName() string {
	return "egg_plat"
}

// TableName Egg query
func (a Query) TableName() string {
	return "egg_query"
}

// TableName Egg
func (a Index) TableName() string {
	return "egg"
}

// TableName Egg
func (a SearchEgg) TableName() string {
	return "egg"
}
