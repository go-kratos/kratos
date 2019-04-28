package academy

const (
	//StateRemove 移除状态
	StateRemove = -1
	//StateNormal 正常状态
	StateNormal = 0
	//BusinessForArchvie 稿件
	BusinessForArchvie = 1
	//BusinessForArticle 专栏
	BusinessForArticle = 2
	//LogClientAcademy 日志服务类型
	LogClientAcademy = 181
	//DefaultState check search archive state
	DefaultState = 2018
)

//TableName get table name
func (a *Archive) TableName() string {
	return "academy_archive"
}

//Archive for academy achive & article.
type Archive struct {
	ID       int64  `gorm:"column:id"`
	OID      int64  `gorm:"column:oid"`
	Title    string `gorm:"column:title"`
	State    int8   `gorm:"column:state"`
	Business int8   `gorm:"column:business"`
	CTime    string `gorm:"column:ctime"`
	MTime    string `gorm:"column:mtime"`
	Comment  string `gorm:"column:comment"`
	Hot      int64  `gorm:"column:hot"`
}

//TableName get table name
func (at *ArchiveTag) TableName() string {
	return "academy_archive_tag"
}

//ArchiveTag for academy achive & tag relation .
type ArchiveTag struct {
	ID       int64  `gorm:"column:id"`
	OID      int64  `gorm:"column:oid"`
	TID      int64  `gorm:"column:tid"`
	State    int8   `gorm:"column:state"`
	Business int8   `gorm:"column:business"`
	CTime    string `gorm:"column:ctime"`
	MTime    string `gorm:"column:mtime"`
}

//ArchiveOrigin for archive list.
type ArchiveOrigin struct {
	OID      int64
	TIDs     []int64
	Comment  string
	Business int8
}

//ArchiveCount get archive count by tid.
type ArchiveCount struct {
	TID   int64 `gorm:"column:tid"`
	Count int   `gorm:"column:count"` //当前tag关联的稿件量
}

//ArchiveMeta for archive meta.
type ArchiveMeta struct {
	OID     int64              `json:"oid"`
	State   int32              `json:"state"`
	Forbid  int8               `json:"forbid"`
	Cover   string             `json:"cover"`
	Type    string             `json:"type"`
	Title   string             `json:"title"`
	UName   string             `json:"uname"`
	Comment string             `json:"comment"`
	CTime   int64              `json:"ctime"`
	MTime   int64              `json:"mtime"`
	Tags    map[int][]*TagMeta `json:"tags"`
	Hot     int64              `json:"hot"`
}

//ArchiveTags for archive tag relation.
type ArchiveTags struct {
	ID       int64 `gorm:"column:id"`
	TID      int64 `gorm:"column:tid"`
	OID      int64 `gorm:"column:oid"`
	Type     int8  `gorm:"column:type"`
	Business int8  `gorm:"column:business"`
}

//Archives for archive list
type Archives struct {
	Pager *Pager         `json:"pager"`
	Items []*ArchiveMeta `json:"items"`
}

// Pager Pager def.
type Pager struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// LogParam for manager.
type LogParam struct {
	UID    int64  `json:"uid"`
	UName  string `json:"uname"`
	Action string `json:"action"`
	TID    int64  `json:"tid"`
	OIDs   string `json:"oids"`
	OName  string `json:"oname"`
	OState int8   `json:"ostate"`
}

// EsParam for es param.
type EsParam struct {
	OID       int64
	Business  int8
	Keyword   string
	Uname     string
	TID       []int64
	Copyright int
	State     int
	Pn        int
	Ps        int
	IP        string
	TidsMap   map[int][]int64
}

// EsPage for es page.
type EsPage struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// EsArc for search archive.
type EsArc struct {
	OID int64   `json:"oid"`
	TID []int64 `json:"tid"`
}

// SearchResult archive list from search.
type SearchResult struct {
	Page   *EsPage  `json:"page"`
	Result []*EsArc `json:"result"`
}
