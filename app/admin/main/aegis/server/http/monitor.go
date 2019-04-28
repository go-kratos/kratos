package http

import (
	"go-common/app/admin/main/aegis/model/monitor"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func monitorRuleResult(c *bm.Context) {
	var (
		err    error
		bidStr string
		bid    int64
		res    []*monitor.RuleResultData
	)
	bidStr = c.Request.Form.Get("bid")
	if bid, err = strconv.ParseInt(bidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if bid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if res, err = srv.MonitorBuzResult(c, bid); err != nil {
		log.Error("srv.MonitorResult(%d) error(%v)", bid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, nil)
}

func monitorResultOids(c *bm.Context) {
	var (
		err   error
		idStr string
		id    int64
		res   []struct {
			OID  int64 `json:"oid"`
			Time int   `json:"time"`
		}
		oidMap map[int64]int
	)
	idStr = c.Request.Form.Get("id")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidMap, err = srv.MonitorResultOids(c, id); err != nil {
		log.Error("srv.MonitorResultOids(%d) error(%v)", id, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for k, v := range oidMap {
		res = append(res, struct {
			OID  int64 `json:"oid"`
			Time int   `json:"time"`
		}{OID: k, Time: v})
	}
	c.JSON(res, nil)
}

func monitorRuleUpdate(c *bm.Context) {

}
