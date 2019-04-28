package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_remarkLength = 200
)

func adminSubject(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sub, err := rpSvc.Subject(c, oid, int32(tp))
	if err != nil {
		log.Error("rpSvr.AdminGetSubjectState(oid%d,tp,%d)error(%v)", oid, int32(tp))
		c.JSON(nil, err)
		return
	}
	c.JSON(sub, nil)
	return
}

// adminSubjectState modify subject state
func adminSubjectState(c *bm.Context) {
	v := new(struct {
		Oid    []int64 `form:"oid,split" validate:"required"`
		Type   int32   `form:"type" validate:"required"`
		State  int32   `form:"state"`
		Remark string  `form:"remark"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	var adid int64
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	fails, err := rpSvc.ModifySubState(c, adid, adName, v.Oid, v.Type, v.State, v.Remark)
	c.JSON(fails, err)
	return
}

// SubLogSearch returns all subjects in recent 3 months by default,
// accept start time, end time, page, pagesize, order, sort as parameters.
func SubLogSearch(c *bm.Context) {
	v := new(struct {
		Oid       int64  `form:"oid"`
		Type      int32  `form:"type"`
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
		Page      int64  `form:"pn"`
		PageSize  int64  `form:"ps"`
		Order     string `form:"order"`
		Sort      string `form:"sort"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	if v.EndTime == "" {
		v.EndTime = time.Now().Format(model.DateFormat)
	}
	// 默认只展示3个月内的数据
	if v.StartTime == "" {
		v.StartTime = time.Now().AddDate(0, -3, 0).Format(model.DateFormat)
	}
	sp := model.LogSearchParam{
		Oid:       v.Oid,
		Type:      v.Type,
		CtimeFrom: v.StartTime,
		CtimeTo:   v.EndTime,
		Pn:        v.Page,
		Ps:        v.PageSize,
		Order:     v.Order,
		Sort:      v.Sort,
	}
	data, err := rpSvc.SubjectLog(c, sp)
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)
	return
}

// SubFreeze freeze or unfreeze the comments.
func SubFreeze(c *bm.Context) {
	v := new(struct {
		Oid    []int64 `form:"oid,split"`
		Type   int32   `form:"type"`
		Freeze int32   `form:"freeze"`
		Remark string  `form:"remark"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	var adid int64
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	fails, err := rpSvc.FreezeSub(c, adid, adName, v.Oid, v.Type, v.Freeze, v.Remark)
	c.JSON(fails, err)
	return
}
