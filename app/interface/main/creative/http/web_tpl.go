package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func webTemplates(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
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

func webAddTpl(c *bm.Context) {
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
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	typeid, _ := strconv.Atoi(typeidStr)
	if typeid < 0 {
		typeid = 0
	}
	if name == "" {
		log.Error("name can not be empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if copyright == "" {
		copyright = "Original"
	}
	if err := tplSvc.AddTemplate(c, mid, int16(typeid), copyright, name, title, tag, content, time.Now()); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webUpdateTpl(c *bm.Context) {
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
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typeid, _ := strconv.Atoi(typeidStr)
	if typeid < 0 {
		typeid = 0
	}
	if name == "" {
		log.Error("name is empty error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if copyright == "" {
		copyright = "Original"
	}
	// update
	c.JSON(nil, tplSvc.UpdateTemplate(c, id, mid, int16(typeid), copyright, name, title, tag, content, time.Now()))
}

func webDelTpl(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("tid")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
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
