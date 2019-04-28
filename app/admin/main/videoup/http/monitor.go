package http

import (
	"encoding/json"
	"go-common/app/admin/main/videoup/model/monitor"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// monitorRuleResult 获取监控规则的监控结果
func monitorRuleResult(c *bm.Context) {
	var (
		err error
		res []*monitor.RuleResultData
		p   = &monitor.RuleResultParams{}
	)
	if err = c.Bind(p); err != nil {
		return
	}
	if p == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Type == 0 || p.Business == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if res, err = vdaSvc.MonitorResult(c, p); err != nil {
		log.Error("vdaSvc.MonitorResult(%v) error(%v)", p, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, nil)
}

// monitorRuleUpdate 更新/添加监控规则
func monitorRuleUpdate(c *bm.Context) {
	var (
		err error
		p   = new(struct {
			Rule string `form:"rule" validate:"required"`
		})
		rule = &monitor.Rule{}
	)
	if err = c.Bind(p); err != nil {
		log.Error("c.Bind(%v) error(%v)", p, err)
		return
	}
	if p == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(p.Rule), rule); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", p.Rule, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if rule.Type == 0 || rule.Business == 0 || rule.RuleConf == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	UID, _ := getUIDName(c)
	rule.UID = UID
	if err = vdaSvc.MonitorRuleUpdate(c, rule); err != nil {
		log.Error("vdaSvc.MonitorRuleUpdate(%v) error(%v)", rule, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func monitorRuleResultOids(c *bm.Context) {
	var (
		err error
		p   = new(struct {
			Type     int8  `form:"type" validate:"required"`
			Business int8  `form:"business" validate:"required"`
			ID       int64 `form:"id" validate:"required"`
		})
		total  int
		oidMap map[int64]int
	)
	if err = c.Bind(p); err != nil {
		log.Error("c.Bind(%v) error(%v)", p, err)
		return
	}
	if p.Type == 0 || p.Business == 0 || p.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if total, oidMap, _, err = vdaSvc.MoniStayOids(c, p.Type, p.Business, p.ID); err != nil {
		log.Error("vdaSvc.MoniStatsOids(%v) error(%v)", p, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	r := new(struct {
		Total int           `json:"total"`
		Oids  map[int64]int `json:"oids"`
	})
	r.Total = total
	r.Oids = oidMap
	c.JSON(r, nil)
}

func monitorNotify(c *bm.Context) {
	var (
		err error
		res []*monitor.RuleResultData
	)
	if res, err = vdaSvc.MonitorNotifyResult(c); err != nil {
		log.Error("vdaSvc.MonitorNotifyResult() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, nil)
}
