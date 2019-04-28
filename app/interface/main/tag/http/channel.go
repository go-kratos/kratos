package http

import (
	"strconv"

	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func channelCategory(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			From int32 `form:"from"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(svr.ChannelCategories(c, &model.ArgChannelCategories{From: param.From, RealIP: metadata.String(c, metadata.RemoteIP)}))
}

func channelRule(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid int64 `form:"tid" validate:"required,gte=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(svr.ChannelRule(c, param.Tid))
}

func channeList(c *bm.Context) {
	var (
		err   error
		mid   int64
		param = new(struct {
			ID   int64 `form:"id" validate:"required,gte=0"`
			From int32 `form:"from"`
		})
	)
	midStr := c.Request.Form.Get("mid")
	if err = c.Bind(param); err != nil {
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		err = nil
	}
	c.JSON(svr.ChanneList(c, mid, param.ID, param.From))
}

func channelRecommand(c *bm.Context) {
	var (
		err   error
		mid   int64
		param = new(struct {
			From int32 `form:"from"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	midStr := c.Request.Form.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		err = nil
	}
	c.JSON(svr.RecommandChannel(c, mid, param.From))
}

func channelDiscover(c *bm.Context) {
	var (
		err   error
		mid   int64
		param = new(struct {
			From int32 `form:"from"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	midStr := c.Request.Form.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		err = nil
	}
	c.JSON(svr.DiscoveryChannel(c, mid, param.From))
}

func channelResource(c *bm.Context) {
	var (
		err   error
		param = new(model.ArgChannelResource)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Tid == 0 && param.Name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.From != model.ChannelFromH5 {
		param.From = model.ChannelFromApp
	}
	param.RealIP = metadata.String(c, metadata.RemoteIP)
	c.JSON(svr.ChannelResources(c, param))
}

func resourceCheckBack(c *bm.Context) {
	var (
		err  error
		oids []int64
		tp   int
	)
	param := c.Request.Form
	oidStr := param.Get("oids")
	tpStr := param.Get("type")
	if oids, err = xstr.SplitInts(oidStr); err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(oids) > model.ResMaxNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.Atoi(tpStr); err != nil || tp <= 0 {
		log.Error("strconv.Atoi(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.ResChannelCheckBack(c, oids, int32(tp)))
}

func resourceInfos(c *bm.Context) {
	var (
		err   error
		param = new(model.ReqChannelResourceInfos)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	lenIDs := len(param.IDs)
	lenOids := len(param.Oids)
	lenTids := len(param.Tids)
	if (lenOids != lenIDs) || (lenOids != lenTids) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.ResChannelInfos(c, param))
}

func channelSubCustoms(c *bm.Context) {
	var (
		err           error
		mid, vmid     int64
		tp, order     int
		ps, pn, total int
		params        = c.Request.Form
	)
	tpStr := params.Get("type")
	orderStr := params.Get("order")
	vmidStr := params.Get("vmid")
	psStr := params.Get("ps")
	pnStr := params.Get("pn")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	vmid, _ = strconv.ParseInt(vmidStr, 10, 64)
	if vmid > 0 {
		mid = vmid
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.Atoi(tpStr); err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > model.SubTagMaxNum {
		ps = model.SubTagMaxNum
	}
	order, _ = strconv.Atoi(orderStr)
	if order != model.SortOrderASC {
		order = model.SortOrderDESC
	}
	cst, sst, total, err := svr.CustomSubTags(c, mid, order, tp, ps, pn)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 3)
	data["page"] = map[string]int{
		"num":   pn,
		"size":  ps,
		"total": total,
	}
	data["custom"] = cst
	data["standard"] = sst
	c.JSON(data, nil)

}

func upChannelSubCustoms(c *bm.Context) {
	var (
		err    error
		mid    int64
		tp     int
		tids   []int64
		params = c.Request.Form
	)
	tidsStr := params.Get("tids")
	tpStr := params.Get("type")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tidsStr != "" {
		if tids, err = xstr.SplitInts(tidsStr); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", tidsStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tp, err = strconv.Atoi(tpStr); err != nil || tp <= 0 {
		log.Error(" strconv.Atoi(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(tids) > model.MaxChannelSortNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.UpCustomSubChannels(c, mid, tids, tp))
}

func channelSquare(c *bm.Context) {
	var (
		err   error
		param = new(model.ReqChannelSquare)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(svr.ChannelSquare(c, param))
}

func channelDetail(c *bm.Context) {
	var (
		err   error
		mid   int64
		param = new(model.ReqChannelDetail)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.TName == "" && param.Tid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
		if mid > 0 {
			param.Mid = mid
		}
	}
	c.JSON(svr.ChannelDetail(c, param))
}
