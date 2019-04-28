package http

import (
	"fmt"
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_jsonErr = "Incorrect Json Format"
)

func intervsRank(c *bm.Context) {
	var (
		form    = new(model.RankListReq)
		request = new(model.IntervListReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	request.FromRank(form)
	c.JSON(tvSrv.Intervs(request))
}

func intervsMod(c *bm.Context) {
	var (
		form    = new(model.ModListReq)
		request = new(model.IntervListReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	request.FromMod(form)
	c.JSON(tvSrv.Intervs(request))
}

func intervsIndex(c *bm.Context) {
	var (
		form    = new(model.IdxListReq)
		request = new(model.IntervListReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	request.FromIndex(form)
	c.JSON(tvSrv.Intervs(request))
}

func rankPub(c *bm.Context) {
	var (
		form = new(model.RankPubReq)
		req  = new(model.IntervPubReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	if err := req.FromRank(form); err != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), _jsonErr)
		return
	}
	intervPublish(c, req)
}

func indexPub(c *bm.Context) {
	var (
		form = new(model.IdxPubReq)
		req  = new(model.IntervPubReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	if err := req.FromIndex(form); err != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), _jsonErr)
		return
	}
	intervPublish(c, req)
}

func modPub(c *bm.Context) {
	var (
		form = new(model.ModPubReq)
		req  = new(model.IntervPubReq)
	)
	if err := c.Bind(form); err != nil {
		return
	}
	if err := req.FromMod(form); err != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), _jsonErr)
		return
	}
	intervPublish(c, req)
}

// combine the alert msg for too many interventions and cut the slice
func alertMsg(items []*model.SimpleRank, nbLimit int) (msg string, restItems []*model.SimpleRank) {
	var (
		length = len(items)
	)
	if length <= nbLimit {
		return "", items
	}
	msg = "以下内容因超量未发布干预："
	for i := nbLimit; i < length; i++ {
		if i+1 == length {
			msg = msg + fmt.Sprintf("id%d", items[i].ContID)
			continue
		}
		msg = msg + fmt.Sprintf("id%d,", items[i].ContID)
	}
	return msg, items[:nbLimit]
}

func intervPublish(c *bm.Context, req *model.IntervPubReq) {
	var (
		err       error
		invalid   *model.RankError
		alertInfo string // used when too many interventions published
	)
	alertInfo, req.Items = alertMsg(req.Items, tvSrv.IntervLimit) // too many intervention treatment
	invalid, err = tvSrv.RefreshIntervs(req)
	if err != nil {
		log.Error("RefreshIntervs Error %v", err)
		c.JSON(nil, err)
		return
	}
	if invalid != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), fmt.Sprintf("发布失败，以下内容状态错误：id%d", invalid.SeasonID))
		return
	}
	if alertInfo != "" {
		renderErrMsg(c, ecode.OK.Code(), alertInfo)
		return
	}
	renderErrMsg(c, ecode.OK.Code(), "发布成功")
}
