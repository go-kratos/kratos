package http

import (
	"time"

	"go-common/app/admin/main/tag/conf"
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func reportList(c *bm.Context) {
	var (
		err     error
		total   int64
		rptInfo []*model.ReportInfo
		param   = new(model.ParamReportList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Oid != 0 && param.Type <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Audit == model.AuditFirst {
		if param.State > model.AuditNotHanleFirst && param.State < model.AuditHanledFirst {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else if param.Audit == model.AuditSecond {
		if param.State <= model.AuditNotHanleFirst || param.State >= model.AuditHanledFirst {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if len(param.STime) == 0 {
		timeNow := time.Now()
		timeThree := timeNow.AddDate(0, conf.Conf.Tag.SelectTime, 0)
		param.STime = timeThree.Format("2006-01-02 15:04:05")
		param.ETime = timeNow.Format("2006-01-02 15:04:05")
	}
	if rptInfo, total, err = svc.ReportList(c, param); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	data["list"] = rptInfo
	c.JSON(data, nil)
}

func reportInfo(c *bm.Context) {
	var (
		err   error
		rpt   []*model.ReportDetail
		param = new(struct {
			ID int64 `form:"id" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if rpt, err = svc.ReportInfo(c, param.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rpt, nil)
}

func reportHandle(c *bm.Context) {
	var (
		err      error
		uid      int64
		username string
		param    = new(model.ParamReportHandle)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	uid, username = managerInfo(c)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	svc.ReportHandle(c, username, uid, param.ID, param.Audit, param.Action)
	c.JSON(nil, nil)
}

func reportState(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamReportState)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.ReportState(c, param.ID, param.State))
}

func reportIgnore(c *bm.Context) {
	var (
		err error
		// uid      int64
		username string
		param    = new(model.ParamReport)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, username = managerInfo(c)
	if username == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.ReportIgnore(c, username, param.Audit, param.IDs))
}

func reportDelete(c *bm.Context) {
	var (
		err      error
		uid      int64
		username string
		param    = new(model.ParamReport)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	uid, username = managerInfo(c)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.ReportDelete(c, username, uid, param.IDs, param.Audit))
}

func reportPunish(c *bm.Context) {
	var (
		err   error
		uname string
		param = new(model.ParamReportPunish)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, uname = managerInfo(c)
	if uname == "" {
		log.Error("could not get login username.")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.ReasonType == 0 && param.Moral == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Moral != 0 && param.Reason == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.ReportPunish(c, uname, param.Remark, param.Note, param.ID, param.Reason, param.Moral, param.Notify, param.ReasonType, param.BlockTimeLength))
}

func reportLogList(c *bm.Context) {
	var (
		err   error
		total int64
		logs  []*model.ReportLog
		param = new(model.ParamReportLog)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if logs, total, err = svc.ReportLogList(c, param.Oid, param.Tid, param.Mid, param.Rid, param.Type, param.Pn, param.Ps, param.HandleType, param.STime, param.ETime, param.Username); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["log"] = logs
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	c.JSON(data, nil)
}

func reportLogInfo(c *bm.Context) {
	var (
		err   error
		logs  []*model.ReportLog
		param = new(struct {
			ID int64 `form:"id"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if logs, err = svc.ReportLogInfo(c, param.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(logs, nil)
}
