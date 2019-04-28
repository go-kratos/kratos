package http

import (
	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func groupList(c *bm.Context) {
	var params = &model.GroupListParams{}
	if err := c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(mfSvc.GroupList(c, params))
}

func groupAdd(c *bm.Context) {
	var params = &model.Group{}

	if err := c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(mfSvc.AddGroup(c, params))
}

func groupUpdate(c *bm.Context) {
	var (
		params = &model.Group{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	err = mfSvc.UpdateGroup(c, params)
	c.JSON(nil, err)
}

func groupDelete(c *bm.Context) {
	var (
		params = &model.Target{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	err = mfSvc.DeleteGroup(c, params.ID)
	c.JSON(nil, err)
}

func targetList(c *bm.Context) {
	var (
		params = &model.Target{}
		pn, ps int
		err    error
		form   = c.Request.Form
		sort   string
		draw   int
		res    *model.Targets
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	pnStr := form.Get("page")
	psStr := form.Get("pagesize")
	drawStr := form.Get("draw")
	sort = form.Get("sort")
	if sort == "" {
		sort = "ctime,0"
	}
	if drawStr != "" {
		if draw, err = strconv.Atoi(drawStr); err != nil {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		pnStr = form.Get("start")
		psStr = form.Get("length")
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if ps == 0 {
		ps = 20
	}
	if draw > 0 {
		pn = (pn + ps) / ps
	}
	if pn == 0 {
		pn = 1
	}
	res, err = mfSvc.TargetList(c, params, pn, ps, sort)
	if draw > 0 {
		res.Draw = draw
	}
	c.JSON(res, err)
}

func targetAdd(c *bm.Context) {
	var params = &model.Target{}
	if err := c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(mfSvc.AddTarget(c, params))
}

func targetUpdate(c *bm.Context) {
	var (
		params = &model.Target{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	err = mfSvc.UpdateTarget(c, params)
	c.JSON(nil, err)
}

func targetSync(c *bm.Context) {
	var (
		params = &model.Target{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	if params.ID == 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	err = mfSvc.TargetSync(c, params.ID, params.State)
	c.JSON(nil, err)
}

func productAdd(c *bm.Context) {
	var (
		params = &model.Product{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(mfSvc.AddProduct(c, params))
}

func productUpdate(c *bm.Context) {
	var (
		params = &model.Product{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	err = mfSvc.UpdateProduct(c, params)
	c.JSON(nil, err)
}

func productDelete(c *bm.Context) {
	var (
		params = &model.Product{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	if params.ID == 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	err = mfSvc.DeleteProduct(c, params.ID)
	c.JSON(nil, err)
}

func productList(c *bm.Context) {
	c.JSON(mfSvc.AllProducts(c))
}

func collect(c *bm.Context) {
	var (
		params = &monitor.Log{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	midInter, _ := c.Get("mid")
	if midInter != nil {
		params.Mid = strconv.FormatInt(midInter.(int64), 10)
	}
	params.IP = metadata.String(c, metadata.RemoteIP)
	params.Buvid = c.Request.Header.Get("Buvid")
	params.UserAgent = c.Request.Header.Get("User-Agent")
	mfSvc.Collect(c, params)
	c.JSON(nil, nil)
}
