package http

import (
	"strconv"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func addFavTopic(c *bm.Context) {
	params := c.Request.Form
	midIfc, _ := c.Get("mid")
	tpStr := params.Get("tpid")
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%s)", tpStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.AddFavTopic(c, midIfc.(int64), tp, c.Request.Header.Get("Cookie"), params.Get("access_key"))
	c.JSON(nil, err)
}

func delFavTopic(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	tpStr := params.Get("tpid")
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%s)", tpStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.DelFavTopic(c, mid.(int64), tp)
	c.JSON(nil, err)
}

// isTopicFavouried determine topic whether or not favouried by mid
func isTopicFavoured(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	tpStr := params.Get("tpid")
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%s)", tpStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	faved, err := favSvc.IsTopicFavoured(c, mid.(int64), tp)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{"favoured": faved}
	c.JSON(data, nil)
}

func favTopics(c *bm.Context) {
	var appInfo *model.AppInfo
	params := c.Request.Form
	mid, _ := c.Get("mid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps >= conf.Conf.Fav.MaxPagesize || ps <= 0 {
		ps = conf.Conf.Fav.MaxPagesize
	}
	platformStr := params.Get("platform")
	buildStr := params.Get("build")
	mobiAppStr := params.Get("mobi_app")
	deviceStr := params.Get("device")
	if platformStr != "" && buildStr != "" {
		appInfo = &model.AppInfo{
			Platform: platformStr,
			Build:    buildStr,
			MobiApp:  mobiAppStr,
			Device:   deviceStr,
		}
	}
	data, err := favSvc.FavTopics(c, mid.(int64), pn, ps, appInfo)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
