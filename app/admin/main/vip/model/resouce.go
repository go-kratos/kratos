package model

import "go-common/library/time"

// ResourceCode .
type ResourceCode struct {
	ID          int64     `json:"id"`
	BatchCodeID int64     `json:"batch_code_id"`
	Status      int8      `json:"status"`
	Code        string    `json:"code"`
	Mid         int64     `json:"mid"`
	UseTime     time.Time `json:"use_time"`
	Ctime       time.Time `json:"ctime"`
}

// BatchCode .
type BatchCode struct {
	ID             int64     `json:"id" form:"id"`
	BusinessID     int64     `json:"business_id" form:"business_id" validate:"required"`
	PoolID         int64     `json:"pool_id" form:"pool_id" validate:"required"`
	Status         int8      `json:"status" `
	Type           int8      `json:"type" form:"type"`
	BatchName      string    `json:"batch_name" form:"batch_name" validate:"required"`
	MaxCount       int64     `json:"max_count" form:"max_count"`
	LimitDay       int64     `json:"limit_day" form:"limit_day" validate:"max=10000,min=-1"`
	Reason         string    `json:"reason" form:"reason" validate:"required"`
	Unit           int32     `json:"unit" form:"unit" validate:"required"`
	Count          int64     `json:"count" form:"count" validate:"required"`
	SurplusCount   int64     `json:"surplus_count"`
	Price          float64   `json:"price" form:"price" validate:"required"`
	StartTime      time.Time `json:"start_time" form:"start_time" validate:"required"`
	EndTime        time.Time `json:"end_time" form:"end_time" validate:"required"`
	Contacts       string    `json:"contacts" form:"contacts"`
	ContactsNumber string    `json:"contacts_number" form:"contacts_number"`
	Operator       string    `json:"operator"`
	Ctime          time.Time `json:"ctime"`
}

// CodeVo .
type CodeVo struct {
	ResourceCode
	BatchName   string    `json:"batch_name"`
	BatchStatus int8      `json:"batch_status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Unit        int32     `json:"unit"`
}
