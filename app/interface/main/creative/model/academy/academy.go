package academy

import (
	"go-common/library/time"

	mdlArt "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
)

const (
	_ = iota
	//Course 教程级别
	Course
	//Operation 运营标签
	Operation
	//Classify  分类标签
	Classify
	//ArticleClass 专栏分类
	ArticleClass
	//H5 移动端tag分类
	H5
	//RecommendTag 推荐理由
	RecommendTag

	//BusinessForAll //所有类型稿件
	BusinessForAll = 0
	//BusinessForArchive //视频稿件
	BusinessForArchive = 1
	//BusinessForArticle //专栏稿件
	BusinessForArticle = 2
)

//H5Conf for h5 conf.
type H5Conf struct {
	//OfficialID 官方推荐id
	OfficialID int64
	//EditorChoiceID 编辑精选id
	EditorChoiceID int64
	//NewbCourseID 新人课程id
	NewbCourseID int64
	//ResourceID 资源管理位id 防重复
	ResourceID int64
}

//TagClassMap for tag type map.
func TagClassMap(ty int) (s string) {
	switch ty {
	case Course:
		s = "course_level"
	case Operation:
		s = "operation_tag"
	case Classify:
		s = "classify_tag"
	case ArticleClass:
		s = "article_class"
	case H5:
		s = "h5"
	case RecommendTag:
		s = "recommend_tag"
	}
	return
}

//Tag for academy tag.
type Tag struct {
	ID       int64     `json:"id"`
	ParentID int64     `json:"parent_id"`
	Type     int8      `json:"type"`
	State    int8      `json:"-"`
	Business int8      `json:"-"`
	Name     string    `json:"name"`
	Desc     string    `json:"-"`
	CTime    time.Time `json:"-"`
	MTime    time.Time `json:"-"`
	Children []*Tag    `json:"children,omitempty"`
}

//TagClassify map tag type name.
func TagClassify() map[int]string {
	return map[int]string{
		Course:       "教程级别",
		Operation:    "运营标签",
		Classify:     "分类标签",
		ArticleClass: "专栏分类",
	}
}

//Archive for academy achive & article.
type Archive struct {
	ID       int64     `json:"id"`
	OID      int64     `json:"oid"`
	State    int8      `json:"-"`
	Business int       `json:"business"`
	CTime    time.Time `json:"-"`
	MTime    time.Time `json:"-"`
	TIDs     []int64   `json:"-"`
}

//ArchiveTag for academy achive & tag relation .
type ArchiveTag struct {
	ID    int64     `json:"id"`
	OID   int64     `json:"oid"`
	TID   int64     `json:"tid"`
	State int8      `json:"-"`
	CTime time.Time `json:"-"`
	MTime time.Time `json:"-"`
}

//ArchiveMeta for archive meta.
type ArchiveMeta struct {
	OID            int64             `json:"oid"`
	MID            int64             `json:"mid"`
	State          int32             `json:"state"`
	Forbid         int8              `json:"forbid"`
	Cover          string            `json:"cover"`
	Type           string            `json:"type"`
	Title          string            `json:"title"`
	HighLightTitle string            `json:"highlight_title"`
	UName          string            `json:"uname"`
	Face           string            `json:"face"`
	Comment        string            `json:"comment"`
	CTime          time.Time         `json:"-"`
	MTime          time.Time         `json:"-"`
	Tags           map[string][]*Tag `json:"tags"`
	Duration       int64             `json:"duration"`
	ArcStat        *api.Stat         `json:"arc_stat,omitempty"`
	ArtStat        *mdlArt.Stats     `json:"art_stat,omitempty"`
	Business       int               `json:"business"`
	Rights         api.Rights        `json:"rights,omitempty"`
}

//ArchiveList for archive list.
type ArchiveList struct {
	Items []*ArchiveMeta `json:"items"`
	Page  *ArchivePage   `json:"page"`
}

//ArchivePage for archive pagination.
type ArchivePage struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

//FeedBack for user advise.
type FeedBack struct {
	// MID      int64     `json:"mid"` //TODO
	Category string    `json:"category"`
	Course   string    `json:"course"`
	Suggest  string    `json:"suggest"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// EsParam for es page.
type EsParam struct {
	OID      int64
	Tid      []int64
	TidsMap  map[int][]int64
	Business int
	Pn       int
	Ps       int
	Keyword  string
	Order    string
	IP       string
	Seed     int64 //支持h5随机推荐
	Duration int   //支持h5时长筛选
}

// EsPage for es page.
type EsPage struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// EsArc for search archive.
type EsArc struct {
	OID      int64    `json:"oid"`
	TID      []int64  `json:"tid"`
	Business int      `json:"business"`
	Title    []string `json:"title"` //highlight
}

// SearchResult archive list from search.
type SearchResult struct {
	Page   *EsPage  `json:"page"`
	Result []*EsArc `json:"result"`
}

//LinkTag for link tag.
type LinkTag struct {
	ID     int64 `json:"id"`
	TID    int64 `json:"tid"`
	LinkID int64 `json:"link_id"`
}

//RecArchive for archive.
type RecArchive struct {
	OID      int64             `json:"oid"`
	MID      int64             `json:"mid"`
	Cover    string            `json:"cover"`
	Title    string            `json:"title"`
	Business int               `json:"business,omitempty"` //只针对标签课程
	Duration int64             `json:"duration,omitempty"`
	ArcStat  *api.Stat         `json:"arc_stat,omitempty"`
	ArtStat  *mdlArt.Stats     `json:"art_stat,omitempty"`
	Tags     map[string][]*Tag `json:"tags,omitempty"`
}

//RecArcList for recommend archive list.
type RecArcList struct {
	Items []*RecArchive `json:"items"`
	Name  string        `json:"name"`
	TID   int64         `json:"tid"`
}

//RecConf for tag conf.
type RecConf struct {
	TIDs []int64
	PID  int64
}

//KV key for tag ids val for type ids
type KV struct {
	Key []int64 `json:"key"`
	Val []int64 `json:"val"`
}

//CourseRec for course rec
type CourseRec struct {
	ID    int64 `json:"id"`
	Rank  int64 `json:"rank"`
	Shoot *KV   `json:"shoot"`
	Scene *KV   `json:"scene"`
	Edit  *KV   `json:"edit"`
	Mmd   *KV   `json:"mmd"`
	Sing  *KV   `json:"sing"`
	Bang  *KV   `json:"bang"`
	Other *KV   `json:"other"`
}

//Drawn for Drawn rec
type Drawn struct {
	ID         int64 `json:"id"`
	Rank       int64 `json:"rank"`
	MobilePlan *KV   `json:"mobile_plan"`
	ScreenPlan *KV   `json:"screen_plan"`
	RecordPlan *KV   `json:"record_plan"`
	Other      *KV   `json:"other"`
}

//Video for Video rec
type Video struct {
	ID          int64 `json:"id"`
	Rank        int64 `json:"rank"`
	MobileMake  *KV   `json:"mobile_make"`
	AudioEdit   *KV   `json:"audio_edit"`
	EditCompose *KV   `json:"edit_compose"`
	Other       *KV   `json:"other"`
}

//Person for person rec
type Person struct {
	ID    int64 `json:"id"`
	Rank  int64 `json:"rank"`
	Other *KV   `json:"other"`
}

//Recommend for all type
type Recommend struct {
	Course *CourseRec `json:"course"`
	Drawn  *Drawn     `json:"drawn"`
	Video  *Video     `json:"video"`
	Person *Person    `json:"person"`
}
