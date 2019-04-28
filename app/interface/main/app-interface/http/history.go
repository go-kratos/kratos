package http

import (
	"strconv"
	"strings"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/history"
	hismodle "go-common/app/interface/main/history/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

var (
	busMap = map[string][]string{
		"":        {},
		"all":     {},
		"archive": {"archive", "pgc"},
		"article": {"article", "article-list"},
		"live":    {"live"},
	}
)

//history list
func historyList(c *bm.Context) {
	param := &history.HisParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = 1
	}
	if param.Ps > 20 || param.Ps <= 0 {
		param.Ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(historySvr.List(c, param.Mid, param.Build, param.Pn, param.Ps, param.Platform, plat))
}

// shortAll get shorturl list
func live(c *bm.Context) {
	param := &history.LiveParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	roomIDs, err := xstr.SplitInts(param.RoomIDs)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(historySvr.Live(c, roomIDs))
}

//history list
func liveList(c *bm.Context) {
	param := &history.HisParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = 1
	}
	if param.Ps > 20 || param.Ps <= 0 {
		param.Ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(historySvr.LiveList(c, param.Mid, param.Build, param.Pn, param.Ps, param.Platform, plat))
}

//history cursor
func historyCursor(c *bm.Context) {
	param := &history.HisParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if param.Ps > 20 || param.Ps <= 0 {
		param.Ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	businesses, ok := busMap[param.Business]
	if !ok {
		log.Error("historyCursor invalid business(%s)", param.Business)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(historySvr.Cursor(c, param.Mid, param.Build, param.Max, param.Ps, param.Platform, param.MaxTP, plat, businesses))
}

//history del
func historyDel(c *bm.Context) {
	param := &history.DelParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	var hisRes []*hismodle.Resource
	for _, boid := range param.Boids {
		bo := strings.Split(boid, "_")
		if len(bo) != 2 {
			log.Error("historyDel invalid param(%+v)", param)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		oid, _ := strconv.ParseInt(bo[1], 10, 0)
		if oid == 0 {
			log.Error("historyDel invalid param(%+v)", param)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		hisRes = append(hisRes, &hismodle.Resource{
			Oid:      oid,
			Business: bo[0],
		})
	}
	c.JSON(nil, historySvr.Del(c, param.Mid, hisRes))
}

//history clear
func historyClear(c *bm.Context) {
	param := &history.HisParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	businesses, ok := busMap[param.Business]
	if !ok {
		log.Error("historyCursor invalid business(%s)", param.Business)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, historySvr.Clear(c, param.Mid, businesses))
}
