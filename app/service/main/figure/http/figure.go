package http

import (
	"strconv"

	"go-common/app/service/main/figure/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func figureInfo(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form

		midStr = params.Get("mid")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%+v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.FigureWithRank(c, mid))
}

func figureInfos(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchInfo{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if len(v.MIDs) > 50 {
		err = ecode.RequestErr
		return
	}
	var (
		frs []*model.FigureWithRank
	)
	frs, err = svc.BatchFigureWithRank(c, v.MIDs)
	if frs == nil {
		frs = make([]*model.FigureWithRank, 0)
	}
	c.JSON(frs, err)
}
