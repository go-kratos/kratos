package model

// 一些常量
const (
	MaxBatchLen             = 10
	MaxTopicNameLen         = 64
	MaxTopicDescLen         = 256
	MaxSvTopicNum           = 15
	MaxTopicVideoLen        = 10
	MaxTopicLen             = 10
	TopicVideoSize          = 10
	DiscoveryTopicVideoSize = 6
	DiscoveryTopicSize      = 3
	CmsTopicSize            = 10
	MaxDiscoveryTopicPage   = 300
	MaxTopicVideoOffset     = 1000
	MaxStickTopicNum        = 10
	MaxStickTopicVideoNum   = 6
)

// Topic状态
const (
	TopicStateAvailable   = 0
	TopicStateUnavailable = 1
)

// redis key format
const (
	RedisStickTopicKey      = "stick:topic"
	ReidsStickTopicVideoKey = "stick:topic:video:%d"
)

// CursorValue 发现页下/话题详情页下的cursor
type CursorValue struct {
	// ！！！注意：这里的offset=db_offset+1
	Offset    int `json:"offset"`     // 默认值为0，从1开始，parseCursor中设置
	StickRank int `json:"stick_rank"` // 默认值为0，从1开始
}
