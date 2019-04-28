package model

import (
	"go-common/library/time"
)

//VipBusinessInfo vip_business_info table
type VipBusinessInfo struct {
	ID             int64     `json:"id"`
	BusinessName   string    `json:"businessName"`
	BusinessType   int8      `json:"businessType"`
	Status         int8      `json:"status"`
	AppKey         string    `json:"appKey"`
	Secret         string    `json:"secret"`
	Contacts       string    `json:"contacts"`
	ContactsNumber string    `json:"contactsNumber"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
}

// VipResourcePool vip_resource_pool table
type VipResourcePool struct {
	ID             int64     `json:"id"`
	PoolName       string    `json:"poolName"`
	BusinessID     int64     `json:"businessId"`
	BusinessName   string    `json:"businessName"`
	Reason         string    `json:"reason"`
	CodeExpireTime time.Time `json:"codeExpireTime"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Contacts       string    `json:"contacts"`
	ContactsNumber string    `json:"contactsNumber"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
}

// VipResourceBatch vip_resource_batch table
type VipResourceBatch struct {
	ID             int64     `json:"id"`
	PoolID         int64     `json:"poolId"`
	Unit           int64     `json:"unit"`
	Count          int64     `json:"count"`
	Ver            int64     `json:"ver"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	SurplusCount   int64     `json:"surplusCount"`
	CodeUseCount   int64     `json:"codeUseCount"`
	DirectUseCount int64     `json:"directUseCount"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
}

//VipResourceCode vip resource code.
type VipResourceCode struct {
	ID          int64     `json:"id"`
	BatchCodeID int64     `json:"batch_code_id"`
	Status      int8      `json:"status"`
	Days        int32     `json:"days"`
	RelationID  string    `json:"relation_id"`
	Code        string    `json:"code"`
	Mid         int64     `json:"mid"`
	UseTime     time.Time `json:"use_time"`
}

//VipResourceBatchCode vip resource batchcode.
type VipResourceBatchCode struct {
	ID           int64     `json:"id"`
	BusinessID   int64     `json:"business_id"`
	PoolID       int64     `json:"pool_id"`
	Status       int8      `json:"status"`
	Type         int8      `json:"type"`
	MaxCount     int64     `json:"max_count"`
	LimitDay     int64     `json:"limit_day"`
	BatchName    string    `json:"batch_name"`
	Reason       string    `json:"reason"`
	Unit         int32     `json:"unit"`
	Count        int32     `json:"count"`
	SurplusCount int32     `json:"surplus_count"`
	Price        float64   `json:"price"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

//VipActiveShow vip active show.
type VipActiveShow struct {
	ID            int64  `json:"id"`
	ProductName   string `json:"product_name"`
	ProductPic    string `json:"product_pic"`
	ProductDetail string `json:"product_detail"`
	RelationID    string `json:"relation_id"`
	BusID         string `json:"bus_id"`
	UseType       string `json:"use_type"`
	Type          int16  `json:"type"`
}

//CodeInfoResp code info Response
type CodeInfoResp struct {
	ID       int64     `json:"id"`
	UserTime time.Time `json:"user_time"`
	Code     string    `json:"code"`
}
