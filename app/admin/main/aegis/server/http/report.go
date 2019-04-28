package http

import (
	"net/http"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func commonPre(c *bm.Context) (res map[string][]*model.ReportFlowItem, form map[string][24]*model.ReportFlowItem, err error) {
	opt := new(model.OptReport)
	if err = c.Bind(opt); err != nil {
		return
	}

	if opt.UName == "debug" {
		opt.Type = model.TypeTotal
		opt.UName = "total"
	}

	return srv.ReportTaskflow(c, opt)
}

func taskflow(c *bm.Context) {
	res, form, err := commonPre(c)
	c.JSONMap(map[string]interface{}{
		"data": res,
		"form": form,
	}, err)
}

func taskflowCSV(c *bm.Context) {
	opt := new(model.OptReport)
	if err := c.Bind(opt); err != nil {
		return
	}
	res, _, err := commonPre(c)
	if err != nil || len(res) == 0 {
		c.JSON(nil, err)
	}

	csv, err := formatReport(res)
	if err != nil {
		log.Error("formatReport error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: FormatCSV(csv),
		Title:   "任务吞吐量报表",
	})
}

func tasksubmitpre(c *bm.Context, pager bool) (res *model.ReportSubmitRes, p *model.DayPager, err error) {
	pm := new(model.OptReportSubmit)
	if err = c.Bind(pm); err != nil {
		return
	}
	if pm.FlowID <= 0 || pm.Bt > pm.Et || (pager && (pm.Pn > pm.Et || pm.Pn < pm.Bt)) {
		err = ecode.RequestErr
		return
	}
	//默认展示近期7天内的数据
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if pm.Bt == "" && pm.Et == "" {
		pm.Bt = time.Now().Add(-8 * 24 * time.Hour).Format("2006-01-02")
		pm.Et = yesterday
	} else if pm.Et == "" { //最近日期
		pm.Et = yesterday
	} else if pm.Bt == "" { //最远日期
		pm.Bt = "2019-01-01"
	}
	//分页的话，默认获取第一页内容（按照日期降序排列，则为最后一日的内容)
	bt := pm.Bt
	et := pm.Et
	if pager {
		if pm.Pn == "" {
			pm.Pn = et
		}
		pm.Et = pm.Pn
		pm.Bt = pm.Pn
	}

	if res, err = srv.ReportTaskSubmit(c, pm); err != nil || res == nil {
		return
	}

	//按照日期降序排列
	if pager {
		p = &model.DayPager{
			Pn:      pm.Pn,
			IsLast:  bt == pm.Pn,
			IsFirst: et == pm.Pn,
		}
	}
	return
}

func taskSubmit(c *bm.Context) {
	res, pager, err := tasksubmitpre(c, true)
	c.JSONMap(map[string]interface{}{
		"data":  res,
		"pager": pager,
	}, err)
}

func taskSubmitCSV(c *bm.Context) {
	res, _, err := tasksubmitpre(c, false)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	csv, err := formatReportTaskSubmit(res)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.Render(http.StatusOK, CSV{
		Content: FormatCSV(csv),
		Title:   "任务操作报表",
	})
}
