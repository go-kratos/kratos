package http

import (
	"fmt"
	"go-common/app/admin/main/videoup/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// searchVideo search video entrance
func searchVideo(c *bm.Context) {
	var (
		err error
	)
	p := &search.VideoParams{}

	if err = c.Bind(p); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if p.Action == "trash" {
		authSrc.Permit("VIDEO_TRASH")(c)
	} else if p.Action == "list" {
		authSrc.Permit("VIDEO_LIST")(c)
	}
	if c.IsAborted() {
		return
	}
	if p.MonitorList != "" { //如果是监控列表，则去掉其他筛选条件（监控列表没有接搜索）
		p = &search.VideoParams{
			MonitorList: p.MonitorList,
			Pn:          p.Pn,
			Ps:          p.Ps,
		}
	}
	videoData, err := vdaSvc.SearchVideo(c, p)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(videoData, nil)
}

// searchCopyright search copyright entrance
func searchCopyright(c *bm.Context) {
	var (
		err error
		p   struct {
			Kw string `form:"text"`
		}
	)
	if err = c.Bind(&p); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	copyright, err := vdaSvc.SearchCopyright(c, p.Kw)
	if err != nil {
		log.Error("searchCopyright(%v) error(%v)", p, err)
		c.JSONMap(map[string]interface{}{
			"message": fmt.Sprintf("版权接口请求失败error(%v)。找谷安-Kwan", err),
		}, ecode.RequestErr)
		return
	}
	c.JSON(copyright.Result, nil)
}

func searchArchive(c *bm.Context) {
	var (
		err error
	)
	p := &search.ArchiveParams{}

	if err = c.Bind(p); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	archiveData, err := vdaSvc.SearchArchive(c, p)
	if c.IsAborted() {
		return
	}
	if err != nil {
		log.Error("searchArchive(%v) error(%v)", p, err)
		c.JSONMap(map[string]interface{}{
			"message": fmt.Sprintf("搜索接口请求失败error(%v)。", err),
		}, ecode.RequestErr)
		return
	}
	c.JSON(archiveData, nil)
}
