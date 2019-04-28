package http

import (
	"strconv"

	"go-common/app/interface/main/web/conf"
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// tagAids gets aids via tag
func tagAids(c *bm.Context) {
	var (
		err           error
		tid           int64
		pn, ps, total int
		arcs          []*v1.Arc
		params        = c.Request.Form
		tidStr        = params.Get("tid") // tag id
		pnStr         = params.Get("pn")
		psStr         = params.Get("ps")
	)
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Tag.MaxSize {
		ps = conf.Conf.Tag.MaxSize
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if total, arcs, err = webSvc.TagAids(c, tid, pn, ps); err != nil {
		log.Error("webSvc.TagAids(%d, %d, %d) error(%v)", tid, pn, ps, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"data":  arcs,
		"total": total,
	}, nil)
}

func tagDetail(c *bm.Context) {
	v := new(struct {
		TagID int64 `form:"tag_id" validate:"min=1"`
		Ps    int   `form:"ps" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(webSvc.TagDetail(c, v.TagID, v.Ps))
}
