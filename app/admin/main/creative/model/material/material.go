package material

import (
	xtime "go-common/library/time"
)

// consts .
const (
	StateDelete = 2
	StateOff    = 1
	StateOn     = 0
	//注意 因为历史原因  bgm 和其他素材没能在bilibili_creative.material一个表集中管理  针对素材类型 为bgm保留了type=3
	//字幕库
	TypeSubTitle = int8(0)
	//字体库
	TypeFont = int8(1)
	//滤镜库
	TypeFilter = int8(2)
	//bgm库
	TypeBGM = int8(3)
	//热词
	TypeHotWord = int8(4)
	//拍摄贴纸  ext 新增贴纸类型    默认为0 普通贴纸，存储格式是bitmask参考属性位   0普通   1人脸   2手势   3画面效果  （不是自然数顺序 服务端不校验）
	TypeSticks = int8(5)
	//贴纸Icon
	TypeSticksIcon = int8(6)
	//投稿贴纸
	TypeCreativeSticks = int8(7)
	//投稿转场
	TypeCreativeTransition = int8(8)
	//合拍库
	TypeCooperate = int8(9)
	//主题库
	TypeTheme = int8(10)
)

var (
	_materialtype = map[int8]string{
		TypeSubTitle:           "字幕库",
		TypeFont:               "字体库",
		TypeFilter:             "滤镜库",
		TypeBGM:                "bgm库",
		TypeHotWord:            "热词",
		TypeSticks:             "贴纸",
		TypeSticksIcon:         "贴纸Icon",
		TypeCreativeSticks:     "投稿贴纸",
		TypeCreativeTransition: "投稿转场",
		TypeCooperate:          "合拍库",
		TypeTheme:              "主题库",
	}
)

// InMaterialType in correct materialtype.
func InMaterialType(cate int8) (ok bool) {
	_, ok = _materialtype[cate]
	return
}

// Material model is the model for Material
type Material struct {
	ID            int64      `json:"id" form:"id" gorm:"column:id"`
	UID           int64      `json:"uid" form:"id" gorm:"column:uid"`
	Name          string     `json:"name" form:"name" gorm:"column:name"`
	Extra         string     `json:"extra" form:"extra" gorm:"column:extra"`
	Rank          int        `json:"rank" form:"rank" gorm:"column:rank"`
	Type          int8       `json:"type" form:"type" gorm:"column:type"`
	Platform      int        `json:"platform" form:"platform" gorm:"column:platform"`
	Build         string     `json:"build" form:"build" gorm:"column:build"`
	State         int8       `json:"state" form:"state" gorm:"column:state"`
	CategoryID    int64      `json:"category_id"  gorm:"-"`
	CategoryIndex int64      `json:"category_index" gorm:"-"`
	CategoryName  string     `json:"category_name" gorm:"-"`
	CTime         xtime.Time `json:"ctime" form:"ctime" gorm:"column:ctime"`
	MTime         xtime.Time `json:"mtime" form:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Material) TableName() string {
	return "material"
}

// Result def.
type Result struct {
	Items []*Material `json:"items"`
	Pager *Pager      `json:"pager"`
}

// Pager Pager def.
type Pager struct {
	Num   int   `json:"num"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// Param is used to parse user request
type Param struct {
	ID            int64  `form:"id" gorm:"column:id" json:"id"`
	Name          string `form:"name" gorm:"column:name" json:"name"`
	Extra         string `form:"extra" gorm:"column:extra" json:"extra"`
	Rank          int    `form:"rank" gorm:"column:rank" json:"rank"`
	Type          int8   `form:"type" gorm:"column:type" json:"type"`
	Cover         string `form:"cover" json:"cover"`
	Platform      int    `form:"platform" json:"platform"`
	Build         string `form:"build" json:"build"`
	DownloadURL   string `form:"download_url" json:"download_url"`
	ExtraURL      string `form:"extra_url" json:"extra_url"`
	ExtraField    string `form:"extra_field" json:"extra_field"`
	Max           int8   `form:"max" json:"max"`
	CategoryID    int64  `form:"category_id" json:"category_id"`
	CategoryIndex int64  `form:"category_index" json:"category_index"`
	SubType       int8   `form:"sub_type"  json:"sub_type"`
	Style         int8   `form:"style"  json:"style"`
	Tip           string `form:"tip"  json:"tip"`
	WhilteList    int8   `form:"white_list"  json:"white_list"`
	MaterialAID   int64  `form:"material_aid"  json:"material_aid"`
	MaterialCID   int64  `form:"material_cid"  json:"material_cid"`
	DemoAID       int64  `form:"demo_aid"  json:"demo_aid"`
	DemoCID       int64  `form:"demo_cid"  json:"demo_cid"`
	MissionID     int64  `form:"mission_id"  json:"mission_id"`
	FilterType    int8   `form:"filter_type"  json:"filter_type"`
}

// TableName is used to identify table name in gorm
func (Param) TableName() string {
	return "material"
}
