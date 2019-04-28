package model

import "go-common/library/time"

//ResoucePoolBo  pool bo
type ResoucePoolBo struct {
	PN             int       `form:"pn" default:"1"`
	PS             int       `form:"ps" default:"20"`
	ID             int       `form:"pool_id"`
	PoolName       string    `form:"pool_name"`
	BusinessID     int       `form:"biz_id"`
	StartTime      time.Time `form:"start_time"`
	EndTime        time.Time `form:"end_time"`
	BatchID        int       `form:"batch_id"`
	Reason         string    `form:"reason"`
	CodeExpireTime time.Time `form:"code_expire_time"`
	Contacts       string    `form:"contacts"`
	ContactsNumber string    `form:"contacts_number"`
}

//ResouceBatchBo resouce batch bo
type ResouceBatchBo struct {
	ID             int       `form:"id"`
	PoolID         int       `form:"pool_id"`
	Unit           int       `form:"unit"`
	Count          int       `form:"count"`
	StartTime      time.Time `form:"start_time"`
	EndTime        time.Time `form:"end_time"`
	SurplusCount   int       `form:"surplus_count"`
	CodeUseCount   int       `form:"code_use_count"`
	DirectUseCount int       `form:"direct_use_count"`
}

//ResouceBatchVo resouce batch vo
type ResouceBatchVo struct {
	VipResourceBatch
	PoolName string `json:"pool_name"`
}

//ArgPrivilege .
type ArgPrivilege struct {
	PrivilegeID int                   `form:"privilege_id"`
	Name        string                `form:"name"`
	Remark      string                `form:"remark"`
	PcLink      string                `form:"pc_link"`
	H5Link      string                `form:"h5_link"`
	BgColor     string                `form:"bg_color"`
	Type        int                   `form:"type"`
	Mapping     []ArgPrivilegeMapping `form:"platforms"`
}

//ArgPrivilegeMapping .
type ArgPrivilegeMapping struct {
	Status   int    `form:"status"`
	Platform int    `form:"platform"`
	Icon     string `form:"icon"`
}

// ArgCode .
type ArgCode struct {
	ID           int64     `form:"id"`
	Code         string    `form:"code"`
	Mid          int64     `form:"mid"`
	BusinessID   int64     `form:"business_id"`
	PoolID       int64     `form:"pool_id"`
	BatchCodeID  int64     `form:"batch_code_id"`
	Status       int8      `form:"status"`
	UseStartTime time.Time `form:"use_start_time"`
	UseEndTime   time.Time `form:"use_end_time"`
	BatchCodeIDs []int64   `form:"batch_code_ids"`
}

// ArgBatchCode .
type ArgBatchCode struct {
	ID         int64     `form:"id"`
	BusinessID int64     `form:"business_id"`
	PoolID     int64     `form:"pool_id"`
	Name       string    `form:"name"`
	Status     int8      `form:"status"`
	StartTime  time.Time `form:"start_time"`
	EndTime    time.Time `form:"end_time"`
}

// ArgPushData .
type ArgPushData struct {
	ProgressStatus int8 `form:"progress_status"`
	Status         int8 `form:"status"`
	PN             int  `form:"pn" default:"1"`
	PS             int  `form:"ps" default:"20"`
}
