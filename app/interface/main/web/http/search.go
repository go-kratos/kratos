package http

import (
	"strconv"

	"go-common/app/interface/main/web/model"
	bm "go-common/library/net/http/blademaster"
)

func searchAll(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	v := new(model.SearchAllArg)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn <= 0 {
		v.Pn = 1
	}
	singleColumnStr := c.Request.Form.Get("single_column")
	if v.SingleColumn, err = strconv.Atoi(singleColumnStr); err != nil {
		v.SingleColumn = -1
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.SearchAll(c, mid, v, buvid, c.Request.Header.Get("User-Agent"), ""))
}

func searchByType(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	v := new(model.SearchTypeArg)
	if err = c.Bind(v); err != nil {
		return
	}
	singleColumnStr := c.Request.Form.Get("single_column")
	if v.SingleColumn, err = strconv.Atoi(singleColumnStr); err != nil {
		v.SingleColumn = -1
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.SearchByType(c, mid, v, buvid, c.Request.Header.Get("User-Agent")))
}

func searchRec(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	v := new(struct {
		Pn         int    `form:"page" default:"1" validate:"min=1"`
		Ps         int    `form:"pagesize" default:"5" validate:"min=1"`
		Keyword    string `form:"keyword"`
		FromSource string `form:"from_source"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.SearchRec(c, mid, v.Pn, v.Ps, v.Keyword, v.FromSource, buvid, c.Request.Header.Get("User-Agent")))
}

func searchDefault(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	v := new(struct {
		FromSource string `form:"from_source"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.SearchDefault(c, mid, v.FromSource, buvid, c.Request.Header.Get("User-Agent")))
}

func upRec(c *bm.Context) {
	var buvid string
	v := new(model.SearchUpRecArg)
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	c.JSON(webSvc.UpRec(c, mid, v, buvid))
}

func searchEgg(c *bm.Context) {
	v := new(struct {
		EggID int64 `form:"egg_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(webSvc.SearchEgg(c, v.EggID))
}
