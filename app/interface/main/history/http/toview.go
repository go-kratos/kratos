package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

// toView return the user toview list
func toView(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(Page)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Pn < 1 {
		v.Pn = 1
	}
	if v.Ps > cnf.History.Max || v.Ps <= 0 {
		v.Ps = cnf.History.Max
	}
	list, count, err := hisSvc.ToView(c, mid, v.Pn, v.Ps, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"list":  list,
		"count": count,
	}
	c.JSON(data, nil)
}

// toView return the user toview list
func webToView(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(Page)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Pn < 1 {
		v.Pn = 1
	}
	if v.Ps > cnf.History.Max || v.Ps <= 0 {
		v.Ps = cnf.History.Max
	}
	list, count, err := hisSvc.WebToView(c, mid, v.Pn, v.Ps, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"list":  list,
		"count": count,
	}
	c.JSON(data, nil)
}

// delToView delete the user video of toview.
func delToView(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			Aids   []int64 `form:"aid,split"`
			Viewed bool    `form:"viewed"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, hisSvc.DelToView(c, mid, v.Aids, v.Viewed, metadata.String(c, metadata.RemoteIP)))
}

// addToView add video to the user toview list.
func addToView(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			Aid int64 `form:"aid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if collector != nil {
		collector.InfoAntiCheat2(c, "", strconv.FormatInt(v.Aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(v.Aid, 10), infoc.ItemTypeAv, infoc.ActionToView, strconv.FormatInt(v.Aid, 10))
	}
	c.JSON(nil, hisSvc.AddToView(c, mid, v.Aid, metadata.String(c, metadata.RemoteIP)))
}

// addToViews add videos to the user toview list.
func addMultiToView(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			Aids []int64 `form:"aids,split"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if collector != nil {
		collector.InfoAntiCheat2(c, "", xstr.JoinInts(v.Aids), strconv.FormatInt(mid, 10), xstr.JoinInts(v.Aids), infoc.ItemTypeAv, infoc.ActionToView, xstr.JoinInts(v.Aids))
	}
	c.JSON(nil, hisSvc.AddMultiToView(c, mid, v.Aids, metadata.String(c, metadata.RemoteIP)))
}

func managerToView(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mid int64 `form:"mid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	list, err := hisSvc.ManagerToView(c, v.Mid, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"list": list,
	}
	c.JSON(data, nil)
}

// remainingToView get the quantity of user's remaining toview.
func remainingToView(c *bm.Context) {
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	remaining, err := hisSvc.RemainingToView(c, mid, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"count": remaining,
	}
	c.JSON(data, nil)
}

// clearToView clear the user toview list.
func clearToView(c *bm.Context) {
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, hisSvc.ClearToView(c, mid))
}
