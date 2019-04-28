package model

import (
	xtime "go-common/library/time"
)

// Report APP report info.
type Report struct {
	ID           int64  `json:"id"`
	APPID        int64  `json:"app_id"` // application
	PlatformID   int    `json:"platform_id"`
	Mid          int64  `json:"mid"`
	Buvid        string `json:"buvid"`
	Build        int    `json:"build"`
	TimeZone     int    `json:"time_zone"`
	NotifySwitch int    `json:"notify_switch"` // global notification switch
	DeviceToken  string `json:"device_token"`
	DeviceBrand  string `json:"device_brand"`
	DeviceModel  string `json:"device_model"`
	OSVersion    string `json:"os_version"`
	Extra        string `json:"extra"`
	Dtime        int64  `json:"dtime"`
}

// Task push task info.
type Task struct {
	ID          string         `json:"id"`   // task id
	Job         int64          `json:"job"`  // 多个子任务拥有同一个 job name
	Type        int            `json:"type"` // 任务类型 1:后台全量 2:后台批量 3:业务推送
	APPID       int64          `json:"app_id"`
	BusinessID  int64          `json:"business_id"`
	PlatformID  int            `json:"platform_id"`
	Platform    []int          `json:"platform"`
	Title       string         `json:"title"`
	Summary     string         `json:"summary"`
	LinkType    int8           `json:"link_type"`
	LinkValue   string         `json:"link_value"`
	Build       map[int]*Build `json:"build"`
	Sound       int            `json:"sound"`
	Vibration   int            `json:"vibration"`
	PassThrough int            `json:"pass_through"`
	MidFile     string         `json:"mid_file"`
	Progress    *Progress      `json:"progress"`
	PushTime    xtime.Time     `json:"push_time"`
	ExpireTime  xtime.Time     `json:"expire_time"`
	Status      int8           `json:"status"`
	Group       string         `json:"group"`
	ImageURL    string         `json:"image_url"`
	Extra       *TaskExtra     `json:"extra"`
}

// TaskExtra task extra.
type TaskExtra struct {
	Group    string `json:"group"`
	Filename string `json:"filename,omitempty"` // 任务文件的名称(前端展示用)
}

// Build version limit.
type Build struct {
	Build     int    `json:"build"`
	Condition string `json:"condition"`
}

// Progress task push progress.
type Progress struct {
	// total indicators
	Status           int8  `json:"st"`            // 任务状态
	MidTotal         int64 `json:"mid_total"`     // 任务接收到的mid总数
	MidValid         int64 `json:"mid_valid"`     // 能查到token的mid数
	MidMissed        int64 `json:"mm"`            // mid_missed 查不到token的mid数
	MidMissedSuccess int64 `json:"mms"`           // mid_missed_success 无效mid补偿推送成功的
	MidMissedFailed  int64 `json:"mmf"`           // mid_missed_failed 无效mid补偿推送失败的
	TokenTotal       int64 `json:"token_total"`   // 一共要推送的token数
	TokenValid       int64 `json:"token_valid"`   // 有效token
	TokenDelay       int64 `json:"token_delay"`   // 延迟推送的token
	TokenSuccess     int64 `json:"token_success"` // 推送成功的
	TokenFailed      int64 `json:"token_failed"`  // 推送失败的

	// brand indicators
	Brands map[int]int64 `json:"brands"` // 各品牌统计数据

	// server indicators
	RetryTimes int64      `json:"retry"` // 重试次数
	BeginTime  xtime.Time `json:"btime"` // 开始时间
	PushTime   xtime.Time `json:"ptime"` // 开始推送时间
	EndTime    xtime.Time `json:"etime"` // 结束时间
}

// APP appication
type APP struct {
	ID            int64
	Name          string
	PushLimitUser int
}

// Business business
type Business struct {
	ID            int64              `json:"id"`
	APPID         int64              `json:"app_id"`
	Name          string             `json:"name"`
	Desc          string             `json:"desc"`
	Token         string             `json:"token"`
	Sound         int                `json:"sound"`
	Vibration     int                `json:"vibration"`
	ReceiveSwitch int                `json:"receive_switch"`
	PushSwitch    int                `json:"push_switch"`
	SilentTime    BusinessSilentTime `json:"silent_time"`
	PushLimitUser int                `json:"push_limit_user"`
	Whitelist     int                `json:"whitelist"`
}

// BusinessSilentTime .
type BusinessSilentTime struct {
	BeginHour, EndHour     int
	BeginMinute, EndMinute int
}

// PushInfo push message.
type PushInfo struct {
	Job         int64
	TaskID      string
	APPID       int64
	Title       string
	Summary     string
	LinkType    int8
	LinkValue   string
	PushTime    xtime.Time
	ExpireTime  xtime.Time
	PassThrough int
	Sound       int
	Vibration   int
	ImageURL    string
}

// PushItem push item.
type PushItem struct {
	Platform  int
	Token     string
	Mid       int64
	Sound     int
	Vibration int
	Build     int
}

// PushChanItem push channel item.
type PushChanItem struct {
	Info *PushInfo
	Item *PushItem
}

// PushChanItems push channel item.
type PushChanItems struct {
	Info  *PushInfo
	Items []*PushItem
}

// Auth cert or auth info.
type Auth struct {
	APPID      int64
	PlatformID int
	Name       string // 第三方名称 for android，例如 小米
	Key        string // android的包名 或 iOS的 cert key
	Value      string // android的 auth 或 iOS的 cert value
	BundleID   string // just for iOS
}

// HTTPResponse http response.
type HTTPResponse struct {
	Code int
	Msg  string
}

// Callback push callback.
type Callback struct {
	Task     string
	APP      int64
	Platform int
	Mid      int64
	Pid      int // mobi_app ID
	Token    string
	Buvid    string
	Click    uint8 // 是否被点击
	Brand    int
	Extra    *CallbackExtra
}

// CallbackExtra .
type CallbackExtra struct {
	Status  int `json:"st"`
	Channel int `json:"chan"`
}
