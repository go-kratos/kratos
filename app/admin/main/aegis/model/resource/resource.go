package resource

import (
	"encoding/json"
	"go-common/library/log"
	xtime "go-common/library/time"
)

//Resource .
type Resource struct {
	ID         int64      `json:"id" gorm:"primary_key" form:"id"`
	BusinessID int64      `json:"business_id" gorm:"column:business_id" form:"business_id"`
	OID        string     `json:"oid" gorm:"column:oid" form:"oid"`
	MID        int64      `json:"mid" gorm:"column:mid" form:"mid"`
	Content    string     `json:"content" gorm:"column:content" form:"content"`
	Extra1     int64      `json:"extra1" gorm:"column:extra1" form:"extra1"`
	Extra2     int64      `json:"extra2" gorm:"column:extra2" form:"extra2"`
	Extra3     int64      `json:"extra3" gorm:"column:extra3" form:"extra3"`
	Extra4     int64      `json:"extra4" gorm:"column:extra4" form:"extra4"`
	Extra1s    string     `json:"extra1s" gorm:"column:extra1s" form:"extra1s"`
	Extra2s    string     `json:"extra2s" gorm:"column:extra2s" form:"extra2s"`
	MetaData   string     `json:"metadata" gorm:"column:metadata" form:"metadata"`
	Ctime      xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime      xtime.Time `json:"mtime" gorm:"column:mtime"`
	Extra5     int64      `json:"extra5" gorm:"column:extra5" form:"extra5"`
	Extra6     int64      `json:"extra6" gorm:"column:extra6" form:"extra6"`
	Extra3s    string     `json:"extra3s" gorm:"column:extra3s" form:"extra3s"`
	Extra4s    string     `json:"extra4s" gorm:"column:extra4s" form:"extra4s"`
	ExtraTime1 string     `json:"extratime1" gorm:"column:extratime1" form:"extratime1"`
	OCtime     string     `json:"octime" gorm:"column:octime" form:"octime"`
	Ptime      string     `json:"ptime" gorm:"column:ptime" form:"ptime"`
}

// Result .
type Result struct {
	ID            int64           `json:"id" gorm:"primary_key" form:"id"`
	RID           int64           `json:"rid" gorm:"column:rid" form:"rid"`
	Attribute     int64           `json:"attribute" gorm:"column:attribute" form:"attribute" default:"-1"`
	Note          string          `json:"note" gorm:"column:note" form:"note" submit:"string"`
	RejectReason  string          `json:"reject_reason" gorm:"column:reject_reason" form:"reject_reason" submit:"string"`
	ReasonID      int64           `json:"reason_id" gorm:"column:reason_id" form:"reason_id" default:"0" submit:"int"`
	State         int             `json:"state" gorm:"column:state" form:"state"`
	PubTime       xtime.Time      `json:"pubtime" gorm:"column:pubtime"`
	DelTime       xtime.Time      `json:"deltime" gorm:"column:deltime"`
	Ctime         xtime.Time      `json:"ctime" gorm:"column:ctime"`
	Mtime         xtime.Time      `json:"mtime" gorm:"column:mtime"`
	AttributeList map[string]int8 `json:"attribute_list" gorm:"-" submit:"map"`
}

// AttrParse 属性值解析为属性展开结果
func (r *Result) AttrParse(cfg map[string]uint) {
	r.AttributeList = make(map[string]int8)
	for name, bit := range cfg {
		r.AttributeList[name] = int8((r.Attribute >> bit) & int64(1))
	}
}

// AttrSet 展开结果计算回属性值
func (r *Result) AttrSet(cfg map[string]uint) {
	var attr int64
	for name, bit := range cfg {
		if val, ok := r.AttributeList[name]; ok {
			attr += int64(val) << bit
		}
	}
	r.Attribute = attr
}

// MetaData 资源扩展数据项目
type MetaData struct {
	Name   string      `json:"name"`
	CNDesc string      `json:"cndesc"` // 中文描述
	Value  interface{} `json:"value"`
}

// Res .
type Res struct {
	ID            int64                  `json:"id" gorm:"primary_key" form:"id"`
	BusinessID    int64                  `json:"business_id" gorm:"column:business_id" form:"business_id"`
	OID           string                 `json:"oid" gorm:"column:oid" form:"oid"`
	MID           int64                  `json:"mid" gorm:"column:mid" form:"mid"`
	Content       string                 `json:"content" gorm:"column:content" form:"content"`
	Extra1        int64                  `json:"extra1" gorm:"column:extra1" form:"extra1"`
	Extra2        int64                  `json:"extra2" gorm:"column:extra2" form:"extra2"`
	Extra3        int64                  `json:"extra3" gorm:"column:extra3" form:"extra3"`
	Extra4        int64                  `json:"extra4" gorm:"column:extra4" form:"extra4"`
	Extra1s       string                 `json:"extra1s,omitempty" gorm:"column:extra1s" form:"extra1s"`
	Extra2s       string                 `json:"extra2s,omitempty" gorm:"column:extra2s" form:"extra2s"`
	MetaData      string                 `json:"metadata,omitempty" gorm:"column:metadata" form:"metadata"`
	Attribute     int64                  `json:"attribute" gorm:"column:attribute" form:"attribute"`
	Note          string                 `json:"note,omitempty" gorm:"column:note" form:"note"`
	RejectReason  string                 `json:"reject_reason,omitempty" gorm:"column:reject_reason" form:"reject_reason"`
	ReasonID      int64                  `json:"reason_id,omitempty" gorm:"column:reason_id" form:"reason_id"`
	State         int64                  `json:"state" gorm:"column:state" form:"state"`
	Pubtime       xtime.Time             `json:"pubtime,omitempty" gorm:"column:pubtime"`
	Deltime       xtime.Time             `json:"deltime,omitempty" gorm:"column:deltime"`
	Ctime         string                 `json:"ctime"`
	Mtime         xtime.Time             `json:"mtime" gorm:"column:mtime"`
	AttributeList map[string]int8        `json:"attribute_list,omitempty"`
	Metas         map[string]interface{} `json:"metas"`

	Extra5     int64  `json:"extra5"`
	Extra6     int64  `json:"extra6"`
	Extra3s    string `json:"extra3s,omitempty"`
	Extra4s    string `json:"extra4s,omitempty"`
	ExtraTime1 string `json:"extratime1,omitempty"`
	OCtime     string `json:"octime,omitempty"`
	Ptime      string `json:"ptime,omitempty"`
}

// AttrParse 属性值解析为属性展开结果
func (r *Res) AttrParse(cfg map[string]uint) {
	r.AttributeList = make(map[string]int8)
	for name, bit := range cfg {
		r.AttributeList[name] = int8((r.Attribute >> bit) & int64(1))
	}
}

//MetaParse .
func (r *Res) MetaParse() {
	if len(r.MetaData) > 0 {
		meta := make(map[string]interface{})
		if err := json.Unmarshal([]byte(r.MetaData), &meta); err != nil {
			log.Error("MetaParse error(%v)", err)
		}
		r.Metas = meta
	}
}

// Args .
type Args struct {
	RID        int64                  `json:"id" form:"rid"`
	BusinessID int64                  `json:"business_id" form:"business_id"`
	OID        int64                  `json:"oid"  form:"oid"`
	Changes    map[string]interface{} `json:"changes"`
}

// TableName .
func (r *Resource) TableName() string {
	return "resource"
}

//TableName .
func (r *Result) TableName() string {
	return "resource_result"
}
