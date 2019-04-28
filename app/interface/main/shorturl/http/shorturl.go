package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/shorturl/conf"
	"go-common/app/interface/main/shorturl/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// add short url from long url.
func add(c *bm.Context) {
	param := &model.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	// check args
	uri := strings.TrimSpace(param.Uri)
	if uri == "" {
		log.Error("add short url args empty long(%s)", uri)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	short, err := suSvr.Add(c, param.Mid, uri)
	if err != nil {
		log.Error("suSvr.Add error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := map[string]string{
		"url": conf.Conf.Host.Default + short,
	}
	c.JSON(data, nil)
}

// jump redirect short url to long url.
func jump(c *bm.Context) {
	// check path
	if len(c.Request.URL.Path) == 0 || c.Request.URL.Path == "/" || c.Request.URL.Path == "/favicon.ico" || strings.HasPrefix(c.Request.URL.Path, "/x/") {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	su, err := suSvr.ShortCache(c, c.Request.URL.Path[1:])
	if err != nil {
		log.Error("suSvr.Get url(%v) error(%v)", c.Request.URL.Path[1:], err)
		c.JSON(nil, err)
		return
	}
	if su == nil || su.Long == "" || su.State == model.StateDelted {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if !strings.HasPrefix(su.Long, "http://") && !strings.HasPrefix(su.Long, "https://") {
		su.Long = "http://" + su.Long
		return
	}
	// redirect
	http.Redirect(c.Writer, c.Request, su.Long, http.StatusFound)
}

// shortAll get shorturl list
func shortAll(c *bm.Context) {
	param := &model.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	pn, err := strconv.Atoi(param.Pn)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(param.Ps)
	if err != nil || ps > 20 || ps <= 0 {
		ps = 20
	}
	long := strings.TrimSpace(param.Uri)
	data, err := suSvr.ShortLimit(c, pn, ps, param.Mid, long)
	if err != nil {
		log.Error("suSvr.ShortLimit error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, su := range data {
		su.Short = conf.Conf.Host.Default + su.Short
	}
	c.JSONMap(map[string]interface{}{
		"data": data,
		"size": 2233,
	}, nil)
}

// shortState set state
func shortUpdate(c *bm.Context) {
	param := &model.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	uri := strings.TrimSpace(param.Uri)
	if uri == "" {
		log.Error("add short url args empty long(%s)", param.Uri)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Mid <= 0 {
		log.Error("mid less than 0 error(%v)", param.Mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := suSvr.ShortUpdate(context.TODO(), param.ID, param.Mid, uri)
	if err != nil {
		log.Error("suSvr.ShortUpdate error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// shortState set state
func shortDel(c *bm.Context) {
	param := &model.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	if param.Mid <= 0 {
		log.Error("mid less than 0 error(%v)", param.Mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := suSvr.ShortDel(c, param.ID, param.Mid, time.Now())
	if err != nil {
		log.Error("suSvr.ShortState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// shortById by id
func shortByID(c *bm.Context) {
	param := &model.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	data, err := suSvr.ShortByID(c, param.ID)
	if err != nil {
		log.Error("suSvr.ShortState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func ping(c *bm.Context) {
	if err := suSvr.Ping(c); err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		log.Error("shorturl service ping error(%v)", err)
	}
}
