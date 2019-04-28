package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func panelInfo(c *bm.Context) {
	arg := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(tvSrv.PanelInfo(arg.ID))

}

func panelStatus(c *bm.Context) {
	args := new(struct {
		ID     int64 `form:"id" validate:"required"`
		Status int64 `form:"status" default:"-1"`
	})
	if err := c.Bind(args); err != nil {
		return
	}
	if args.Status == -1 {
		renderErrMsg(c, ecode.RequestErr.Code(), "状态不能为空")
		return
	}

	c.JSON(nil, tvSrv.PanelStatus(args.ID, args.Status))

}

func savePanel(c *bm.Context) {
	panel := new(model.TvPriceConfig)

	if err := c.Bind(panel); err != nil {
		return
	}

	if panel.PID == 0 && panel.Price == 0 {
		renderErrMsg(c, ecode.RequestErr.Code(), "价格不能为空")
		return
	}

	if panel.PID == 0 && panel.Month == 0 {
		renderErrMsg(c, ecode.RequestErr.Code(), "月份不能为空")
		return
	}

	if panel.PID != 0 && panel.Stime >= panel.Etime {
		renderErrMsg(c, ecode.RequestErr.Code(), "活动开始时间不能晚于结束时间")
		return
	}

	c.JSON(nil, tvSrv.SavePanel(c, panel))

}

func panelList(c *bm.Context) {
	args := new(struct {
		Platform int64 `form:"platform"`
		Month    int64 `form:"month"`
		SubType  int64 `form:"sub_type" default:"-1"`
		SuitType int64 `form:"suit_type" default:"-1"`
	})

	if err := c.Bind(args); err != nil {
		return
	}

	c.JSON(tvSrv.PanelList(args.Platform, args.Month, args.SubType, args.SuitType))

}
