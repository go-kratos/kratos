package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/app/interface/main/web-feed/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func feed(c *bm.Context) {
	var (
		params = c.Request.Form
		pn, ps int
		mid    int64
		feeds  []*model.Feed
		err    error
	)
	pnStr := params.Get("pn")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	psStr := params.Get("ps")
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Feed.DefaultSize
	} else if ps > conf.Conf.Feed.MaxSize {
		ps = conf.Conf.Feed.MaxSize
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if feeds, err = feedSrv.Feed(c, mid, pn, ps); err != nil {
		log.Error("feedSrv.Feed(%d,%d,%d) error(%v)", mid, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(feeds, err)
}

func feedUnread(c *bm.Context) {
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	count, err := feedSrv.UnreadCount(c, mid)
	if err != nil {
		log.Error("feedSrv.UnreadCount(%d,%v) error(%v)", mid, time.Now(), err)
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Count int `json:"count"`
	}{Count: count}, nil)
}

func articleFeedUnread(c *bm.Context) {
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	count, err := feedSrv.ArticleUnreadCount(c, mid)
	if err != nil {
		log.Error("feedSrv.ArticleUnreadCount(%d,%v) error(%v)", mid, time.Now(), err)
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Count int `json:"count"`
	}{Count: count}, nil)
}

func articleFeed(c *bm.Context) {
	var (
		params = c.Request.Form
		pn, ps int
		feeds  []*artmdl.Meta
		mid    int64
		err    error
	)
	pnStr := params.Get("pn")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	psStr := params.Get("ps")
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Feed.DefaultSize
	} else if ps > conf.Conf.Feed.MaxSize {
		ps = conf.Conf.Feed.MaxSize
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if feeds, err = feedSrv.ArticleFeed(c, mid, pn, ps); err != nil {
		log.Error("feedSrv.ArticleFeed(%d,%d,%d) error(%v)", mid, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(feeds, nil)
}
