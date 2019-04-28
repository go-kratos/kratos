package http

import (
	"strconv"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strings"
)

func topArc(c *bm.Context) {
	var (
		mid, vmid int64
		err       error
	)
	vmidStr := c.Request.Form.Get("vmid")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.TopArc(c, mid, vmid))
}

func setTopArc(c *bm.Context) {
	v := new(struct {
		Aid    int64  `form:"aid" validate:"min=1"`
		Reason string `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	reason := strings.TrimSpace(v.Reason)
	if len([]rune(reason)) > conf.Conf.Rule.MaxTopReasonLen {
		c.JSON(nil, ecode.TopReasonLong)
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.SetTopArc(c, mid, v.Aid, reason))
}

func cancelTopArc(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.DelTopArc(c, mid))
}

func masterpiece(c *bm.Context) {
	var (
		mid, vmid int64
		err       error
	)
	vmidStr := c.Request.Form.Get("vmid")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.Masterpiece(c, mid, vmid))
}

func addMasterpiece(c *bm.Context) {
	v := new(struct {
		Aid    int64  `form:"aid" validate:"min=1"`
		Reason string `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	reason := strings.TrimSpace(v.Reason)
	if len([]rune(reason)) > conf.Conf.Rule.MaxMpReasonLen {
		c.JSON(nil, ecode.TopReasonLong)
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.AddMasterpiece(c, mid, v.Aid, reason))
}

func editMasterpiece(c *bm.Context) {
	v := new(struct {
		Aid    int64  `form:"aid" validate:"min=1"`
		PreAid int64  `form:"pre_aid" validate:"min=1"`
		Reason string `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	reason := strings.TrimSpace(v.Reason)
	if len([]rune(reason)) > conf.Conf.Rule.MaxMpReasonLen {
		c.JSON(nil, ecode.TopReasonLong)
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.EditMasterpiece(c, mid, v.Aid, v.PreAid, reason))
}

func cancelMasterpiece(c *bm.Context) {
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.CancelMasterpiece(c, mid, v.Aid))
}

func arcSearch(c *bm.Context) {
	var (
		res   *model.SearchRes
		total int
		mid   int64
		err   error
		v     = new(model.SearchArg)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.CheckType != "" {
		if _, ok := model.ArcCheckType[v.CheckType]; !ok {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if v.CheckID <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if res, total, err = spcSvc.ArcSearch(c, mid, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"pn":    v.Pn,
		"ps":    v.Ps,
		"count": total,
	}
	data["page"] = page
	data["list"] = res
	c.JSON(data, nil)
}

func arcList(c *bm.Context) {
	var (
		rs  *model.UpArc
		err error
	)
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
		Pn  int32 `form:"pn" default:"1" validate:"min=1"`
		Ps  int32 `form:"ps" default:"20" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if rs, err = spcSvc.UpArcs(c, v.Mid, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int64{
		"pn":    int64(v.Pn),
		"ps":    int64(v.Ps),
		"count": rs.Count,
	}
	data["page"] = page
	data["archives"] = rs.List
	c.JSON(data, nil)
}
