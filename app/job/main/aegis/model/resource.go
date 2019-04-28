package model

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	xtime "go-common/library/time"
)

// Resource .
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

// TableName .
func (t *Resource) TableName() string {
	return "resource"
}

//AddOption add option
type AddOption struct {
	Resource
	State int   `form:"state" json:"state"`
	NetID int64 `form:"net_id" json:"net_id"`
}

// ToQueryURI convert field to uri.
func (opt AddOption) ToQueryURI() url.Values {
	var params = url.Values{}
	params.Add("business_id", strconv.Itoa(int(opt.BusinessID)))
	params.Add("net_id", strconv.Itoa(int(opt.NetID)))
	params.Add("oid", opt.OID)
	params.Add("mid", strconv.Itoa(int(opt.MID)))
	params.Add("content", opt.Content)
	params.Add("extra1", strconv.Itoa(int(opt.Extra1)))
	params.Add("extra2", strconv.Itoa(int(opt.Extra2)))
	params.Add("extra3", strconv.Itoa(int(opt.Extra3)))
	params.Add("extra4", strconv.Itoa(int(opt.Extra4)))
	params.Add("extra5", strconv.Itoa(int(opt.Extra5)))
	params.Add("extra5", strconv.Itoa(int(opt.Extra6)))
	params.Add("extra1s", opt.Extra1s)
	params.Add("extra2s", opt.Extra2s)
	params.Add("extra3s", opt.Extra3s)
	params.Add("extra4s", opt.Extra4s)
	params.Add("extratime1", opt.ExtraTime1)
	params.Add("octime", opt.OCtime)
	params.Add("ptime", opt.Ptime)
	params.Add("metadata", opt.MetaData)
	return params
}

//UpdateOption update option
type UpdateOption struct {
	BusinessID int64                  `json:"business_id"`
	NetID      int64                  `json:"net_id"`
	OID        string                 `json:"oid"`
	Update     map[string]interface{} `json:"update"`
}

//ToQueryURI convert field to uri.
func (opt UpdateOption) ToQueryURI() url.Values {
	var params = url.Values{}
	params.Add("business_id", strconv.Itoa(int(opt.BusinessID)))
	params.Add("net_id", strconv.Itoa(int(opt.NetID)))
	params.Add("oid", opt.OID)

	if bs, err := json.Marshal(opt.Update); err == nil && len(bs) > 0 {
		params.Add("update", string(bs))
	}
	return params
}

//CancelOption .
type CancelOption struct {
	BusinessID int64    `json:"business_id"`
	Oids       []string `json:"oids"`
	Reason     string   `json:"reason"`
}

// ToQueryURI convert field to uri.
func (opt CancelOption) ToQueryURI() url.Values {
	var params = url.Values{}
	params.Add("business_id", strconv.Itoa(int(opt.BusinessID)))
	params.Add("oids", strings.Join(opt.Oids, ","))
	params.Add("reason", opt.Reason)

	return params
}
