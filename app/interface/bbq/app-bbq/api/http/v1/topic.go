package v1

import topic "go-common/app/service/bbq/topic/api"

// TopicVideo 话题视频的结构
type TopicVideo struct {
	*VideoResponse
	CursorValue string `json:"cursor_value"` // 透传给客户端，标记在列表中的位置
	HotType     int64  `json:"hot_type"`     // 热门类型，直接用topic给的数据
}

//TopicDetail 话题详情页，可作为详情页回包，也可作为发现页话题列表的item
type TopicDetail struct {
	HasMore   bool             `json:"has_more"`
	TopicInfo *topic.TopicInfo `json:"topic_info,omitempty"`
	List      []*TopicVideo    `json:"list,omitempty"`
}

// DiscoveryRes 发现页返回结构
type DiscoveryRes struct {
	BannerList []*Banner      `json:"banner_list"`
	HotWords   []string       `json:"hot_words"`
	TopicList  []*TopicDetail `json:"topic_list"`
	HasMore    bool           `json:"has_more"`
}

// Banner Banner结构
type Banner struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   int16  `json:"type"`
	Scheme string `json:"scheme"`
	PIC    string `json:"pic"`
}

//DiscoveryReq 发现页请求
type DiscoveryReq struct {
	Page int32 `form:"page"  validate:"gt=0,required"`
}

// TopicSearchReq 话题搜索请求
type TopicSearchReq struct {
	Page    int32  `form:"page"  validate:"gt=0,required"`
	Keyword string `from:"keyword"`
}

// TopicSearchResponse 话题搜索回包
type TopicSearchResponse struct {
	HasMore bool               `json:"has_more"`
	List    []*topic.TopicInfo `json:"list,omitempty"`
}
