package dhh

import "time"

// DHHDB 大航海信息DB层
type DHHDB struct {
	ID            int64
	UID           int64
	TargetId      int64
	PrivilegeType int64
	StartTime     time.Time
	ExpiredTime   time.Time
	Ctime         time.Time
	Utime         time.Time
}

// DHHDBTime 大航海信息验DB层(时间转换)
type DHHDBTime struct {
	ID            int64
	UID           int64
	TargetId      int64
	PrivilegeType int64
	StartTime     string
	ExpiredTime   string
	Ctime         string
	Utime         string
}

// ModelDHHList dhh list
type ModelDHHList struct {
	Data []DHHDB
}

// ModelExpLog 行为日志上报结构
type ModelExpLog struct {
	Mid  int64
	Uexp int64
	Rexp int64
	Ts   int64
	// 业务来源
	ReqBizDesc string
	Buvid      string
	// 具体描述
	Content map[string]string
}

// DaHangHaiRedis 等级基础结构
type DaHangHaiRedis struct {
	Id            int64  `json:"id"`
	Uid           int64  `json:"uid"`
	TargetId      int64  `json:"target_id"`
	PrivilegeType int64  `json:"privilege_type"`
	StartTime     string `json:"start_time"`
	ExpiredTime   string `json:"expired_time"`
	Ctime         string `json:"ctime"`
	Utime         string `json:"utime"`
}

// DaHangHaiRedis2 等级基础结构
type DaHangHaiRedis2 struct {
	Id            string `json:"id"`
	Uid           string `json:"uid"`
	TargetId      string `json:"target_id"`
	PrivilegeType string `json:"privilege_type"`
	StartTime     string `json:"start_time"`
	ExpiredTime   string `json:"expired_time"`
	Ctime         string `json:"ctime"`
	Utime         string `json:"utime"`
}
