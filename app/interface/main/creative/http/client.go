package http

import (
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/order"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
	"time"
)

func clientViewArc(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	av, err := arcSvc.View(c, mid, aid, ip, archive.PlatformWindows)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if av == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	av.Archive.Desc = archive.ShortDesc(av.Archive.Desc)
	arcElec, err := elecSvc.ArchiveState(c, aid, mid, ip)
	if err != nil {
		log.Error("archive(%d) error(%v)", mid, err)
	}
	c.JSON(map[string]interface{}{
		"archive":      av.Archive,
		"videos":       av.Videos,
		"archive_elec": arcElec,
	}, nil)
}

func clientDelArc(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// del
	c.JSON(nil, arcSvc.Del(c, mid, aid, ip))
}

func clientArchives(c *bm.Context) {

	params := c.Request.Form
	class := params.Get("class")
	order := params.Get("order")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check params
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 50 {
		ps = 10
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	arc, err := arcSvc.Archives(c, mid, int16(tid), keyword, order, class, ip, pn, ps, 0)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arc, nil)
}

func clientArchiveSearch(c *bm.Context) {
	params := c.Request.Form
	class := params.Get("class")
	order := params.Get("order")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check params
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 50 {
		ps = 10
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	arc, err := arcSvc.Search(c, mid, int16(tid), keyword, order, class, ip, pn, ps, 0)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arc, nil)
}

func clientPre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	mf, err := accSvc.MyInfo(c, mid, ip, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// commercial order
	orders := make([]*order.Order, 0)
	if arcSvc.AllowOrderUps(mid) {
		orders, _ = arcSvc.ExecuteOrders(c, mid, metadata.String(c, metadata.RemoteIP))
	}
	mf.Commercial = arcSvc.AllowCommercial(c, mid)
	c.JSON(map[string]interface{}{
		"typelist":   arcSvc.Types(c, "ch"),
		"activities": arcSvc.Activities(c),
		"myinfo":     mf,
		"orders":     orders,
	}, nil)
}

func clientTemplates(c *bm.Context) {

	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	tps, err := tplSvc.Templates(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tps, nil)
}

func clientTags(c *bm.Context) {
	params := c.Request.Form
	tidStr := params.Get("typeid")
	title := params.Get("title")
	filename := params.Get("filename")
	desc := params.Get("desc")
	cover := params.Get("cover")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check params
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	tags, err := dataSvc.Tags(c, mid, uint16(tid), title, filename, desc, cover, archive.TagPredictFromWindows)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// del
	c.JSON(tags, nil)
}

func clientAddTpl(c *bm.Context) {
	params := c.Request.Form
	typeidStr := params.Get("typeid")
	copyright := params.Get("arctype")
	name := params.Get("name")
	title := params.Get("title")
	tag := params.Get("keywords")
	content := params.Get("description")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	typeid, err := strconv.ParseInt(typeidStr, 10, 16)
	if err != nil || typeid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if name == "" {
		log.Error("name is empty error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// add
	err = tplSvc.AddTemplate(c, mid, int16(typeid), copyright, name, title, tag, content, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func clientUpdateTpl(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("tid")
	typeidStr := params.Get("typeid")
	copyright := params.Get("arctype")
	name := params.Get("name")
	title := params.Get("title")
	tag := params.Get("keywords")
	content := params.Get("description")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typeid, err := strconv.ParseInt(typeidStr, 10, 16)
	if err != nil || typeid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if name == "" {
		log.Error("name is empty error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// update
	c.JSON(nil, tplSvc.UpdateTemplate(c, id, mid, int16(typeid), copyright, name, title, tag, content, time.Now()))
}

func clientDelTpl(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("tid")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// del
	c.JSON(nil, tplSvc.DelTemplate(c, id, mid, time.Now()))
}
