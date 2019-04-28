package academy

//TableName get table name
func (o *Occupation) TableName() string {
	return "academy_occupation"
}

//TableName get table name
func (s *Skill) TableName() string {
	return "academy_skill"
}

//TableName get table name
func (s *Software) TableName() string {
	return "academy_software"
}

//TableName get table name
func (s *ArcSkill) TableName() string {
	return "academy_arc_skill"
}

//TableName get table name
func (s *TagLink) TableName() string {
	return "academy_tag_link"
}

//Occupation for academy occupation.
type Occupation struct {
	ID           int64    `gorm:"column:id"    form:"id" json:"id"`
	Rank         int64    `gorm:"column:rank"  form:"rank" json:"rank"`
	State        int      `gorm:"column:state" form:"state" json:"state"`
	Name         string   `gorm:"column:name"  form:"name" json:"name"`
	Desc         string   `gorm:"column:desc"  form:"desc" json:"desc"`
	Logo         string   `gorm:"column:logo"  form:"logo" json:"logo"`
	MainStep     string   `gorm:"column:main_step"     form:"main_step" json:"main_step"`
	MainSoftware string   `gorm:"column:main_software" form:"main_software" json:"main_software"`
	CTime        string   `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime        string   `gorm:"column:mtime" form:"mtime" json:"-"`
	Skill        []*Skill `gorm:"-" form:"-" json:"skill"`
	Count        int      `gorm:"-" form:"-" json:"count"`
}

//Skill for academy skill.
type Skill struct {
	ID       int64       `gorm:"column:id"    form:"id" json:"id"`
	OID      int64       `gorm:"column:oid"   form:"oid" json:"oid"`
	State    int         `gorm:"column:state" form:"state" json:"state"`
	Name     string      `gorm:"column:name"  form:"name" json:"name"`
	Desc     string      `gorm:"column:desc"  form:"desc" json:"desc"`
	CTime    string      `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime    string      `gorm:"column:mtime" form:"mtime" json:"-"`
	Software []*Software `gorm:"-" form:"-" json:"software"`
	Count    int         `gorm:"-" form:"-" json:"count"`
}

//Software for academy software.
type Software struct {
	ID    int64  `gorm:"column:id"    form:"id" json:"id"`
	SkID  int64  `gorm:"column:skid"  form:"skid" json:"skid"`
	State int    `gorm:"column:state" form:"state" json:"state"`
	Name  string `gorm:"column:name"  form:"name" json:"name"`
	Desc  string `gorm:"column:desc"  form:"desc" json:"desc"`
	CTime string `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime string `gorm:"column:mtime" form:"mtime" json:"-"`
	Count int    `gorm:"-" form:"-" json:"count"`
}

//ArcSkill for academy archive relation to occupation & skill & software.
type ArcSkill struct {
	ID    int64  `gorm:"column:id"   form:"id" json:"id"`
	AID   int64  `gorm:"column:aid"  form:"aid" json:"aid"`
	PID   int64  `gorm:"column:pid"  form:"pid" json:"pid"`
	SkID  int64  `gorm:"column:skid" form:"skid" json:"skid"`
	SID   int64  `gorm:"column:sid"  form:"sid" json:"sid"`
	Type  int    `gorm:"column:type" form:"type" json:"type"`
	State int    `gorm:"column:state" form:"state" json:"state"`
	CTime string `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime string `gorm:"column:mtime" form:"mtime" json:"-"`
	Title string `gorm:"-" form:"title" json:"title"`
	Pic   string `gorm:"-" form:"pic" json:"pic"`
	Pn    int    `gorm:"-" form:"pn" json:"-"`
	Ps    int    `gorm:"-" form:"ps" json:"-"`
}

//ArcSkills for archive skill list
type ArcSkills struct {
	Pager *Pager      `json:"pager"`
	Items []*ArcSkill `json:"items"`
}

//TagLink for academy h5 tag relation to web tags.
type TagLink struct {
	ID     int64  `gorm:"column:id"  form:"id" json:"id"`
	TID    int64  `gorm:"column:tid" form:"tid" json:"tid"`
	LinkID int64  `gorm:"column:link_id" form:"link_id" json:"link_id"`
	CTime  string `gorm:"column:ctime"   form:"ctime"  json:"-"`
	MTime  string `gorm:"column:mtime"   form:"mtime" json:"-"`
}

//TableName get table name
func (sk *SearchKeywords) TableName() string {
	return "academy_search_keywords"
}

//SearchKeywords for academy h5 search keywords.
type SearchKeywords struct {
	ID       int64             `gorm:"column:id"  form:"id" json:"id"`
	Rank     int64             `gorm:"column:rank" form:"rank" json:"rank"`
	ParentID int64             `gorm:"column:parent_id" form:"parent_id" json:"parent_id"`
	State    int8              `gorm:"column:state" form:"state" json:"state"`
	Name     string            `gorm:"column:name"  form:"name" json:"name"`
	Comment  string            `gorm:"column:comment" form:"comment" json:"comment"`
	CTime    string            `gorm:"column:ctime"   form:"ctime"  json:"-"`
	MTime    string            `gorm:"column:mtime"   form:"mtime" json:"-"`
	Count    int               `gorm:"-" form:"-" json:"count,omitempty"`
	Children []*SearchKeywords `json:"children,omitempty"`
}
