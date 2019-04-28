package http

import (
	"strconv"

	"go-common/app/interface/main/app-show/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func dailyID(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	dailyIDStr := params.Get("daily_id")
	dailyID, err := strconv.Atoi(dailyIDStr)
	if err != nil {
		log.Error("dailyID(%s) error(%v)", dailyIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > 60 || ps <= 0 {
		ps = 60
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := dailySvc.Daily(c, plat, build, dailyID, pn, ps)
	returnJSON(c, data, nil)
}

func columnList(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	buildStr := params.Get("build")
	categoryIDStr := params.Get("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Error("categoryID(%s) error(%v)", categoryIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := dailySvc.ColumnList(plat, build, categoryID)
	returnJSON(c, data, nil)
}

func category(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	buildStr := params.Get("build")
	categoryIDStr := params.Get("category_id")
	columnIDStr := params.Get("column_id")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Error("categoryID(%s) error(%v)", categoryIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	columnID, _ := strconv.Atoi(columnIDStr)
	plat := model.Plat(mobiApp, device)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > 60 || ps <= 0 {
		ps = 60
	}
	data := dailySvc.Category(plat, build, categoryID, columnID, pn, ps)
	returnJSON(c, data, nil)
}
