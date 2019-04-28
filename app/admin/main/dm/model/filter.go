package model

import (
	"go-common/library/time"
)

// all const variable used in dm filter
const (
	FilterUnActive int8 = 0
	FilterActive   int8 = 1

	FilterTypeAll   int8 = -1 //  所有类型
	FilterTypeText  int8 = 0  //  文本类型
	FilterTypeRegex int8 = 1  //  正则类型
	FilterTypeID    int8 = 2  //  用户ID类型

	FilterMaxUpText = 500  // up主关键字最大条数
	FilterMaxUpReg  = 100  // up主正则最大条数
	FilterMaxUpID   = 1000 // up主黑名单最大条数
)

// UpFilter define a new struct, consistent with table "dm_filter_up_%"
type UpFilter struct {
	ID     int64     `json:"id"`
	Oid    int64     `json:"oid"`
	Type   int8      `json:"type"`
	Filter string    `json:"filter"`
	Ctime  time.Time `json:"ctime"`
}

// UpFilterRes return UpFilters and PageInfo
type UpFilterRes struct {
	Result []*UpFilter `json:"result"`
	Page   *PageInfo   `json:"page"`
}
