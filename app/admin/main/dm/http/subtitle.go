package http

import (
	"go-common/app/admin/main/dm/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// subtitleList 字幕后台管理搜索
func subtitleList(c *bm.Context) {
	var (
		v = new(model.SubtitleArg)
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SubtitleList(c, v))
}

// subtitleEdit 字幕操作
func subtitleEdit(c *bm.Context) {
	var (
		v = new(model.EditSubtitleArg)
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.EditSubtitle(c, v))
}

// workflow 回调函数
func subtitleEditCallback(c *bm.Context) {
	var (
		v = new(model.WorkFlowSubtitleArg)
	)
	if err := c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, dmSvc.WorkFlowEditSubtitle(c, v))
}

// subtitleStatusList 字幕状态列表 给举报使用
func subtitleStatusList(c *bm.Context) {
	c.JSON(dmSvc.SubtitleStatusList(c))
}

func subtitleLanList(c *bm.Context) {
	c.JSON(dmSvc.SubtitleLanList(c))
}

// subtitleSwitch 字幕开关
func subtitleSwitch(c *bm.Context) {
	var (
		v = new(struct {
			Aid    int64 `form:"aid" validate:"required"`
			Allow  bool  `form:"allow"`
			Closed bool  `form:"closed"`
		})
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitleSwitch(c, v.Aid, v.Allow, v.Closed))
}
