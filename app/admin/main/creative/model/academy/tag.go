package academy

const (
	_ = iota
	//Course 教程级别
	Course
	//Operation 运营标签
	Operation
	//Classify 分类标签
	Classify
	//ArticleClass 专栏分类
	ArticleClass
	//H5 手机端分类标签
	H5
	//Recommend 理由标签
	Recommend
)

const (
	//StateUnBlock 解冻状态
	StateUnBlock = 0
	//StateBlock 冻结状态
	StateBlock = -1
)

//TableName get table name
func (t *Tag) TableName() string {
	return "academy_tag"
}

//Tag for academy tag.
type Tag struct {
	ID       int64  `gorm:"column:id"`
	ParentID int64  `gorm:"column:parent_id"`
	Type     int8   `gorm:"column:type"`
	State    int8   `gorm:"column:state"`
	Business int8   `gorm:"column:business"`
	Name     string `gorm:"column:name"`
	Desc     string `gorm:"column:desc"`
	CTime    string `gorm:"column:ctime"`
	MTime    string `gorm:"column:mtime"`
	Rank     int64  `gorm:"column:rank"`
	Children []*Tag `json:"children,omitempty"`
}

//TagMeta for academy tag reuslt.
type TagMeta struct {
	ID       int64      `json:"id"`
	ParentID int64      `json:"parent_id"`
	Type     int8       `json:"type"`
	State    int8       `json:"state"`
	Business int8       `json:"business"`
	Count    int        `json:"count"`
	Name     string     `json:"name"`
	Desc     string     `json:"desc"`
	Children []*TagMeta `json:"children,omitempty"`
	Rank     int64      `gorm:"column:rank"`
	LinkID   []int64    `json:"link_id,omitempty"`
}

//TagClass for tag type name map.
func TagClass() map[int]string {
	return map[int]string{
		Course:       "教程级别",
		Operation:    "运营标签",
		Classify:     "分类标签",
		ArticleClass: "专栏分类",
		H5:           "手机端分类标签",
		Recommend:    "推荐理由",
	}
}
