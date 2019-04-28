package service

import (
	"fmt"

	"go-common/app/admin/main/block/conf"
	"go-common/app/admin/main/block/model"
	"go-common/library/log"
)

// MSGInfo get msg info
func (s *Service) MSGInfo(source model.BlockSource, action model.BlockAction, area model.BlockArea, reason string, days int64) (code string, title, content string) {
	if source == model.BlockSourceBlackHouse {
		areaStr := area.String()
		if areaStr != "" {
			areaStr = fmt.Sprintf("在%s中", areaStr)
		}
		if action == model.BlockActionLimit {
			code = conf.Conf.Property.MSG.BlackHouseLimit.Code
			title = conf.Conf.Property.MSG.BlackHouseLimit.Title
			content = fmt.Sprintf(conf.Conf.Property.MSG.BlackHouseLimit.Content, areaStr, s.convertReason(reason), days)
			return
		}
		if action == model.BlockActionForever {
			code = conf.Conf.Property.MSG.BlackHouseForever.Code
			title = conf.Conf.Property.MSG.BlackHouseForever.Title
			content = fmt.Sprintf(conf.Conf.Property.MSG.BlackHouseForever.Content, areaStr, s.convertReason(reason))
			return
		}
	}
	if source == model.BlockSourceSys {
		if action == model.BlockActionLimit {
			code = conf.Conf.Property.MSG.SysLimit.Code
			title = conf.Conf.Property.MSG.SysLimit.Title
			content = fmt.Sprintf(conf.Conf.Property.MSG.SysLimit.Content, s.convertReason(reason), days)
			return
		}
		if action == model.BlockActionForever {
			code = conf.Conf.Property.MSG.SysForever.Code
			title = conf.Conf.Property.MSG.SysForever.Title
			content = fmt.Sprintf(conf.Conf.Property.MSG.SysForever.Content, s.convertReason(reason))
			return
		}
	}
	if action == model.BlockActionAdminRemove || action == model.BlockActionSelfRemove {
		code = conf.Conf.Property.MSG.BlockRemove.Code
		title = conf.Conf.Property.MSG.BlockRemove.Title
		content = conf.Conf.Property.MSG.BlockRemove.Content
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
