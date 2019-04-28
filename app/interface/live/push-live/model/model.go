package model

// ApPushTask struct of table link_push.ap_push_task
type ApPushTask struct {
	ID         int64  `json:"id"`
	Type       int    `json:"type"`
	TargetID   int64  `json:"target_id"`
	AlertTitle string `json:"alert_title"`
	AlertBody  string `json:"alert_body"`
	MidSource  int    `json:"mid_source"`
	LinkType   int    `json:"link_type"`
	LinkValue  string `json:"link_value"`
	Total      int    `json:"total"`
	ExpireTime int    `json:"expire_time"`
	Group      string
}

// StartLiveMessage StartLiveNotify-T message
type StartLiveMessage struct {
	TargetID   int64  `json:"target_id"`
	Uname      string `json:"uname"`
	LinkValue  string `json:"link_value"`
	ExpireTime int    `json:"expire_time"`
	RoomTitle  string `json:"room_title"`
}

// LiveCommonMessage LivePushCommon-T message
type LiveCommonMessage struct {
	Topic      string                   `json:"topic"`
	MsgID      string                   `json:"msg_id"`
	MsgKey     string                   `json:"msg_key"`
	MsgContent LiveCommonMessageContent `json:"msg_content"`
}

// LiveCommonMessageContent LivePushCommon-T message.msg_content
type LiveCommonMessageContent struct {
	Business   int    `json:"business"`
	Group      string `json:"group"`
	Mids       string `json:"mids"`
	AlertTitle string `json:"alert_title"`
	AlertBody  string `json:"alert_body"`
	LinkValue  string `json:"link_value"`
	LinkType   int    `json:"link_type"`
	ExpireTime int    `json:"expire_time"`
}

// 直播开关DB相关配置
const (
	LivePushType        = 1001
	LivePushSwitchOn    = 1
	LivePushConfigOn    = 1
	PushIntervalKey     = "push_interval"
	PushIntervalDefault = 1800
)

/**
 * 推送类型
 * 注意：这里复用这个常量定义，1 & 2是hbase中关系链的类型，但是RelationAll=3不是，这里只是个业务概念上的类型
 * 	表示取所有关注数据
 */
const (
	// RelationAttention 关注
	RelationAttention = iota + 1
	// RelationSpecial 特别关注
	RelationSpecial
	// RelationAll 关注+特别关注
	RelationAll
)

// 推送后台策略,DB中的记录
const (
	StrategySwitch        = "Switch"           //开启推送开关
	StrategySpecial       = "Special"          //特别关注
	StrategyFans          = "Fans"             //关注
	StrategySwitchSpecial = "SwitchAndSpecial" //开启开关且特别关注
)

// 推送任务标记mid来源，组合来源则取交则可
const (
	TaskSourceSwitch    = 1
	TaskSourceSpecial   = 2
	TaskSourceFans      = 4
	TaskSourceSwitchSpe = 8
)

// 开播提醒消息的group信息
const (
	AttentionGroup      = "follow"               // 关注
	SpecialGroup        = "sfollow"              // 特别关注
	ActivityAppointment = "activity_appointment" // 预约
)

// 业务business配置
const (
	StartLiveBusiness = 1
	ActivityBusiness  = 111
)
