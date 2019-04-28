package block

import (
	"fmt"

	model "go-common/app/service/main/member/model/block"
	"go-common/library/log"
)

const (
	_creditLimit   = `抱歉，你的账号因“%s%s”，现已进行封禁%d天处理，账号解封需要满足以下两个条件:1.账号封禁时间已满。2.完成解封答题（ #{点击进入解封答题}{"http://www.bilibili.com/blackroom/releaseexame.html"} ）全部完成后解封。封禁期间将无法投稿、发送及回复消息，无法发布评论、弹幕，无法对他人评论进行回复、赞踩操作，无法进行投币、编辑标签、添加关注、添加收藏操作。解封后恢复正常，还请遵守社区规范，共同维护良好的社区氛围！`
	_creditForever = `抱歉，你的账号因“%s%s”，现已进行永久封禁处理。封禁期间将无法投稿、发送及回复消息，无法发布评论、弹幕，无法对他人评论进行回复、赞踩操作，无法进行投币、编辑标签、添加关注、添加收藏操作。解封后恢复正常，还请遵守社区规范，共同维护良好的社区氛围！`

	_sysLimit   = `抱歉，你的账号因“%s”，现已进行封禁%d天处理，账号解封需要满足以下两个条件:1.账号封禁时间已满。2.完成解封答题（ #{点击进入解封答题}{"http://www.bilibili.com/blackroom/releaseexame.html"} ）全部完成后解封。封禁期间将无法投稿、发送及回复消息，无法发布评论、弹幕，无法对他人评论进行回复、赞踩操作，无法进行投币、编辑标签、添加关注、添加收藏操作。解封后恢复正常，还请遵守社区规范，共同维护良好的社区氛围！`
	_sysForever = `抱歉，你的账号因“%s”，现已进行永久封禁处理。封禁期间将无法投稿、发送及回复消息，无法发布评论、弹幕，无法对他人评论进行回复、赞踩操作，无法进行投币、编辑标签、添加关注、添加收藏操作。解封后恢复正常，还请遵守社区规范，共同维护良好的社区氛围！`

	_remove = `你的账号已经解除封禁，封禁期间禁止使用的各项社区功能已经恢复。请遵守社区规范，共同维护良好的社区氛围。`
)

// MSGInfo get msg info
func (s *Service) MSGInfo(source model.BlockSource, action model.BlockAction, area model.BlockArea, reason string, days int64) (code string, title, content string) {
	// 小黑屋封禁
	if source == model.BlockSourceBlackHouse {
		areaStr := area.String()
		if areaStr != "" {
			areaStr = fmt.Sprintf("在%s中", areaStr)
		}
		if action == model.BlockActionLimit {
			code = "2_3_2"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_creditLimit, areaStr, s.convertReason(reason), days)
			return
		}
		if action == model.BlockActionForever {
			code = "2_3_3"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_creditForever, areaStr, s.convertReason(reason))
			return
		}
	}
	// B+小黑屋封禁
	if source == model.BlockSourceBplus {
		if action == model.BlockActionLimit {
			code = "2_3_2"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_sysLimit, reason, days)
			return
		}
		if action == model.BlockActionForever {
			code = "2_3_3"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_sysForever, reason)
			return
		}
	}
	// 系统封禁
	if source == model.BlockSourceSys {
		if action == model.BlockActionLimit {
			code = "2_3_4"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_sysLimit, s.convertReason(reason), days)
			return
		}
		if action == model.BlockActionForever {
			code = "2_3_5"
			title = "账号违规处理通知"
			content = fmt.Sprintf(_sysForever, s.convertReason(reason))
			return
		}
	}
	if action == model.BlockActionAdminRemove || action == model.BlockActionSelfRemove {
		code = "2_3_6"
		title = "账号封禁解除通知"
		content = _remove
		return
	}
	log.Error("s.MSGInfo unkown source[%v] action[%v] area[%v] reason[%s] days[%d]", source, action, area, reason, days)
	return
}

func (s *Service) convertReason(reason string) string {
	switch reason {
	case "账号资料相关违规":
		return "账号资料违规"
	case "作品投稿违规":
		return "作品投稿违规"
	case "异常注册账号":
		return "异常注册"
	case "异常答题账号":
		return "异常答题"
	case "异常数据行为":
		return "异常数据行为"
	case "发布违规信息":
		return "发布违规信息"
	case "其他自动封禁", "手动封禁":
		return "违反社区规则"
	default:
		return reason
	}
}
