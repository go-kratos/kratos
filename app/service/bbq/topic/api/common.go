package api

import (
	"context"
	"encoding/json"
	"go-common/library/log"
)

// Transform2Interface 转换成interface
func Transform2Interface(ctx context.Context, data []byte) (inter interface{}, err error) {
	err = json.Unmarshal(data, &inter)
	if err != nil {
		log.Errorw(ctx, "log", "transform to interface fail", "data", string(data))
		return
	}
	return
}

// 话题的状态
const (
	TopicStateAvailable   = 0
	TopicStateUnAvailable = 1

	TopicVideoStateAvailable   = 0
	TopicVideoStateUnAvailable = 1
)

// 话题热门类型的enum，用于TopicInfo->HotType字段
// 开始时使用了hot_type，但其实就是表示特殊的话题状态
const (
	TopicHotTypeHot     = 1 // 热门
	TopicHotTypeHistory = 2 // 历史，暂时只有客户端使用
	TopicHotTypeStick   = 4 // 置顶
)
