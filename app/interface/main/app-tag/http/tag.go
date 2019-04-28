package http

import (
	"net/http"
	"strconv"
	"time"

	"go-common/app/interface/main/app-tag/model"
	"go-common/app/interface/main/app-tag/model/region"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

var (
	_emptyShowItems = []*region.ShowItem{}
)

const (
	_headerBuvid = "Buvid"
)

func tagDetail(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		header = c.Request.Header
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	pullStr := params.Get("pull")
	ctimeStr := params.Get("ctime")
	tidStr := params.Get("tag_id")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	// check params
	// tag id
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil || tid < 1 {
		log.Error("tid(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// pull default
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		pull = true
	}
	// ctime
	ctime, err := strconv.ParseInt(ctimeStr, 10, 64)
	if err != nil || ctime < 0 {
		ctime = 0
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get(_headerBuvid)
	data := map[string]interface{}{}
	tag, similar, err := tagSvr.TagDetail(c, plat, mid, tid, time.Now())
	if err != nil {
		log.Error("tagSvr.TagDetail(%d,%d,%d) error(%v)", plat, mid, tid, err)
		c.JSON(nil, err)
		return
	}
	is, ctop, cbottom, err := tagSvr.TagDefault(c, plat, tid, mid, ctime, pull, build, buvid, time.Now())
	if err != nil {
		log.Error("tagSvr.TagDefault(%d,%d,%d,%t,%d) error(%v)", plat, mid, tid, pull, ctime, err)
	}
	if len(is) == 0 {
		if is, err = tagSvr.TagNew(c, plat, tid, 1, 20, time.Now()); err != nil {
			log.Error("agSvr.TagNew(%d,%d,%d) error(%v)", tid, 1, 20, err)
		}
	}
	if len(is) != 0 {
		data["item"] = is
	} else {
		data["item"] = []struct{}{}
	}
	data["tag"] = tag
	data["similar_tag"] = similar
	data["ctop"] = ctop
	data["cbottom"] = cbottom
	returnJSON(c, data, nil)
}

func tagDefault(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		header = c.Request.Header
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	pullStr := params.Get("pull")
	ctimeStr := params.Get("ctime")
	tidStr := params.Get("tag_id")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	buvid := header.Get(_headerBuvid)
	// check params
	// tag id
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil || tid < 1 {
		log.Error("tid(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// pull default
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		pull = true
	}
	// ctime
	ctime, err := strconv.ParseInt(ctimeStr, 10, 64)
	if err != nil || ctime < 0 {
		ctime = 0
	}
	plat := model.Plat(mobiApp, device)
	data := map[string]interface{}{}
	is, ctop, cbottom, err := tagSvr.TagDefault(c, plat, tid, mid, ctime, pull, build, buvid, time.Now())
	if err != nil {
		log.Error("tagSvr.TagDefault(%d,%d,%d,%t,%d) error(%v)", plat, mid, tid, pull, ctime, err)
		data["item"] = []struct{}{}
		return
	}
	data["item"] = is
	if pull {
		data["ctime"] = ctop
	} else {
		data["ctime"] = cbottom
	}
	returnJSON(c, data, nil)
}

func tagNew(c *bm.Context) {
	var (
		params = c.Request.Form
	)
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	tidStr := params.Get("tag_id")
	pnStr := params.Get("pn")
	// psStr := params.Get("ps")
	// check params
	// tag id
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil || tid < 1 {
		log.Error("tid(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check page
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps := 20
	plat := model.Plat(mobiApp, device)
	data := map[string]interface{}{}
	is, err := tagSvr.TagNew(c, plat, tid, pn, ps, time.Now())
	if err != nil {
		log.Error("agSvr.TagNew(%d,%d,%d) error(%v)", tid, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	data["item"] = is
	returnJSON(c, data, nil)
}

// tagDynamic tag Dynamic
func tagDynamic(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	ridStr := params.Get("rid")
	reidStr := params.Get("reid")
	tagIDStr := params.Get("tag_id")
	tagName := params.Get("tag_name")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, _ := strconv.Atoi(ridStr)
	reid, _ := strconv.Atoi(reidStr)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	buvid := header.Get(_headerBuvid)
	plat := model.Plat(mobiApp, device)
	var tagID int64
	if tagIDStr != "" {
		if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tagID == 0 {
		if tagID, err = tagSvr.TagIDByName(c, tagName); err != nil {
			log.Error("tagName(%s) error(%v)", tagName, err)
			c.JSON(nil, err)
			return
		}
	}
	data := tagSvr.TagDynamic(c, plat, build, rid, reid, tagID, mid, mobiApp, buvid, device, time.Now(), true)
	returnJSON(c, data, nil)
}

// tagDynamicList
func tagDynamicList(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)

	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	reidStr := params.Get("reid")
	ridStr := params.Get("rid")
	ctimeStr := params.Get("ctime")
	buildStr := params.Get("build")
	tagIDStr := params.Get("tag_id")
	tagName := params.Get("tag_name")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	reid, _ := strconv.Atoi(reidStr)
	rid, _ := strconv.Atoi(ridStr)
	pullStr := params.Get("pull")
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		log.Error("pullStr(%s) error(%v)", pullStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// ctime
	ctime, err := strconv.ParseInt(ctimeStr, 10, 64)
	if err != nil || ctime < 0 {
		ctime = 0
	}
	var tagID int64
	if tagIDStr != "" {
		if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tagID == 0 {
		if tagID, err = tagSvr.TagIDByName(c, tagName); err != nil {
			log.Error("tagName(%s) error(%v)", tagName, err)
			c.JSON(nil, err)
			return
		}
	} else {
		tagName, _ = tagSvr.TagNameByID(c, tagID)
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	data := tagSvr.TagDynamicList(c, plat, build, rid, reid, tagID, pull, ctime, mid, mobiApp, buvid, device, tagName, time.Now())
	returnJSON(c, data, nil)
}

func tagRankList(c *bm.Context) {
	params := c.Request.Form
	reidStr := params.Get("reid")
	pnStr := params.Get("pn")
	// psStr := params.Get("ps")
	ps := 20
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	tagIDStr := params.Get("tag_id")
	tagName := params.Get("tag_name")
	orderStr := params.Get("order")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	reid, _ := strconv.Atoi(reidStr)
	var (
		tagID int64
		err   error
	)
	if tagIDStr != "" {
		if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tagID == 0 {
		if tagID, err = tagSvr.TagIDByName(c, tagName); err != nil {
			log.Error("tagName(%s) error(%v)", tagName, err)
			c.JSON(nil, err)
			return
		}
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	if pn*ps > 400 {
		returnJSON(c, _emptyShowItems, nil)
		return
	}
	order := ""
	switch orderStr {
	case "new":
		order = "new"
	}
	plat := model.Plat(mobiApp, device)
	data := tagSvr.TagRankList(c, plat, build, reid, tagID, pn, ps, order, mobiApp)
	returnJSON(c, data, nil)
}

func tagTab(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("rid")
	tagIDStr := params.Get("tag_id")
	tagName := params.Get("tag_name")
	rid, _ := strconv.Atoi(ridStr)
	var (
		tagID int64
		err   error
	)
	if tagIDStr != "" {
		if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tagID == 0 {
		if tagID, err = tagSvr.TagIDByName(c, tagName); err != nil {
			log.Error("tagName(%s) error(%v)", tagName, err)
			c.JSON(nil, err)
			return
		}
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data := tagSvr.TagTab(c, rid, tagID, mid, time.Now())
	returnJSON(c, data, nil)
}

// tagDynamic tag Dynamic
func tagDynamicIndex(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	ridStr := params.Get("rid")
	reidStr := params.Get("reid")
	tagIDStr := params.Get("tag_id")
	tagName := params.Get("tag_name")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	rid, _ := strconv.Atoi(ridStr)
	reid, _ := strconv.Atoi(reidStr)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	var tagID int64
	if tagIDStr != "" {
		if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tagID == 0 {
		if tagID, err = tagSvr.TagIDByName(c, tagName); err != nil {
			log.Error("tagName(%s) error(%v)", tagName, err)
			c.JSON(nil, err)
			return
		}
	}
	data := tagSvr.TagDynamic(c, plat, build, rid, reid, tagID, mid, mobiApp, buvid, device, time.Now(), false)
	returnJSON(c, data, nil)
}

//returnJSON return json no message
func returnJSON(c *bm.Context, data interface{}, err error) {
	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)
	c.Render(code, render.JSON{
		Code:    bcode.Code(),
		Message: "",
		Data:    data,
	})
}
