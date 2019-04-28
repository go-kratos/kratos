package http

import (
	"strings"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func archiveSearch(c *bm.Context) {
	form := c.Request.Form
	appidStr := form.Get("appid")
	switch appidStr {
	case "archive_check":
		archiveCheck(c)
	case "video":
		video(c)
	case "task_qa":
		taskQa(c)
	case "archive_commerce":
		archiveCommerce(c)
	default:
		c.JSON(nil, ecode.RequestErr)
		return
	}
}

func archiveCheck(c *bm.Context) {
	var (
		err error
		sp  = &model.ArchiveCheckParams{
			Bsp: &model.BasicSearchParams{},
		}
		form = c.Request.Form
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	orderStr := form.Get("order")
	if len(sp.Bsp.Order) == 1 && sp.Bsp.Order[0] == "" {
		sp.Bsp.Order = nil
	}
	if orderStr != "" {
		sp.Bsp.Order = strings.Split(orderStr, ",")
	}
	kwFieldsStr := form.Get("kw_fields")
	if kwFieldsStr == "" {
		sp.Bsp.KwFields = []string{"title", "content", "tag", "author"}
	} else {
		sp.Bsp.KwFields = strings.Split(kwFieldsStr, ",")
	}
	res, err := svr.ArchiveCheck(c, sp)
	if err != nil {
		log.Error("svr.ArchiveCheck(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func video(c *bm.Context) {
	var (
		err error
		sp  = &model.VideoParams{
			Bsp: &model.BasicSearchParams{},
		}
		form = c.Request.Form
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if form.Get("order_type") == "" {
		sp.OrderType = -1
	}
	if form.Get("kw_fields") == "" {
		sp.Bsp.KwFields = []string{"arc_title", "arc_author"}
	}
	res, err := svr.Video(c, sp)
	if err != nil {
		log.Error("svr.ArchiveCheck(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func taskQa(c *bm.Context) {
	var (
		err error
		sp  = &model.TaskQa{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.TaskQa(c, sp)
	if err != nil {
		log.Error("svr.TaskQa(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func archiveCommerce(c *bm.Context) {
	var (
		err error
		sp  = &model.ArchiveCommerce{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.ArchiveCommerce(c, sp)
	if err != nil {
		log.Error("svr.ArchiveCommerce(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}
