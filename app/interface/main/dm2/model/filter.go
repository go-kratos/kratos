package model

import (
	"go-common/library/time"
)

// all const variable used in dm filter
const (
	FilterUnActive int8 = 0
	FilterActive   int8 = 1

	FilterTypeText   int8 = 0 //  文本类型
	FilterTypeRegex  int8 = 1 //  正则类型
	FilterTypeID     int8 = 2 //  用户ID类型
	FilterTypeBottom int8 = 4
	FilterTypeTop    int8 = 5
	FilterTypeRev    int8 = 6

	FilterLenText  = 50  // 文本类型最大字符长度
	FilterLenRegex = 200 // 正则类型最大字符长度

	FilterMaxUserText = 1000 // 用户关键字最大条数
	FilterMaxUserReg  = 100  // 用户正则最大条数
	FilterMaxUserID   = 1000 // 用户黑名单最大条数
	FilterMaxUpText   = 500  // up主关键字最大条数
	FilterMaxUpReg    = 100  // up主正则最大条数
	FilterMaxUpID     = 1000 // up主黑名单最大条数

	FilterNotExist = -1000000

	FilterContent = "10000,20000,25000,30000"
)

// UserFilter define a new struct, consistent with table "dm_filter_user_%"
type UserFilter struct {
	ID      int64     `json:"id"`
	Mid     int64     `json:"mid"`
	Type    int8      `json:"type"`
	Filter  string    `json:"filter"`
	Comment string    `json:"comment"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

// UpFilter filter of upper
type UpFilter struct {
	ID      int64     `json:"id"`
	Mid     int64     `json:"mid"`
	Type    int8      `json:"type"`
	Filter  string    `json:"filter"`
	Active  int8      `json:"active"`
	Comment string    `json:"comment"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// GlobalFilter define a new struct, consistent with table "dm_sys_filter"
type GlobalFilter struct {
	ID     int64     `json:"id"`
	Type   int8      `json:"type"`
	Filter string    `json:"filter"`
	Ctime  time.Time `json:"-"`
	Mtime  time.Time `json:"-"`
}
