package model

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"time"
)

//GetSpecialArg arg
type GetSpecialArg struct {
	GroupID   int    `form:"group_id"`
	UID       int    `form:"uid"`
	FromTime  string `form:"from_time"`  // "2006-01-02 15:04:05"
	ToTime    string `form:"to_time"`    // "2006-01-02 15:04:05"
	Order     string `form:"order" `     // 根据mtime排序，默认从升序, 取值，asc/desc
	Export    string `form:"export"`     // csv
	Charset   string `form:"charset"`    // 导出编码格式，默认gbk
	Pn        uint   `form:"pn"`         // 页码，默认1
	Ps        uint   `form:"ps"`         // 每页数量，默认20
	Mids      string `form:"mids"`       // 用户ids,以,分隔
	AdminName string `form:"admin_name"` // 管理员昵称
}

//UpSpecialWithName arg with name
type UpSpecialWithName struct {
	UpSpecial
	UName     string `json:"uname"`
	AdminName string `json:"admin_name"`
}

//GetSpecialByMidArg arg
type GetSpecialByMidArg struct {
	Mid int `form:"mid" validate:"required"`
}

//Copy copy
func (u *UpSpecialWithName) Copy(special *UpSpecial) {
	u.UpSpecial = *special
}

//GetTitleFields get title fields
func (u *UpSpecialWithName) GetTitleFields() []string {
	return []string{
		"配置时间",
		"MID",
		"昵称",
		"所属用户组",
		"描述",
		"配置人",
		"用户组Id",
	}
}

//ToStringFields get to string fields
func (u *UpSpecialWithName) ToStringFields() []string {
	var fields []string
	var mtime = time.Unix(int64(u.MTime), 0)
	fields = append(fields, fmt.Sprintf("%v", mtime.Format(mysql.TimeFormat)))
	fields = append(fields, fmt.Sprintf("%v", u.Mid))
	fields = append(fields, fmt.Sprintf("%v", u.UName))
	fields = append(fields, fmt.Sprintf("%v", u.GroupName))
	fields = append(fields, fmt.Sprintf("%v", u.Note))
	fields = append(fields, fmt.Sprintf("%v", u.AdminName))
	fields = append(fields, fmt.Sprintf("%v", u.GroupID))
	return fields
}

// UpsPage CategoryPager def.
type UpsPage struct {
	Items []*UpSpecialWithName `json:"items"`
	Pager *Pager               `json:"page"`
}

// Pager Common Pager def.
type Pager struct {
	Num   uint `json:"num"`
	Size  uint `json:"size"`
	Total int  `json:"total"`
}
