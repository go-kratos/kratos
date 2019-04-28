package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// byMid get vipinfo by mid.
func byMid(c *bm.Context) {
	var (
		res *model.VipInfoResp
		err error
	)
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.ByMid(c, arg.Mid); err != nil {
		log.Error("vipSvc.ByMid(%d) err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func vipInfos(c *bm.Context) {
	var (
		vMap map[int64]*model.VipInfoResp
		err  error
	)
	arg := new(struct {
		Mids []int64 `form:"mids,split" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err %+v", err)
		return
	}
	if vMap, err = vipSvc.VipInfos(c, arg.Mids); err != nil {
		log.Error("vipSvc.VipInfos(%v)  err(%+v)", arg.Mids, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(vMap, nil)
}

func vipHistory(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgChangeHistory)
	if c.Bind(arg); err != nil {
		log.Error("vipHistory(%d) err(%+v)", arg.Mid, err)
		return
	}

	vh, count, err := vipSvc.History(c, arg)
	rel := make(map[string]interface{})
	rel["data"] = vh
	rel["total"] = count

	c.JSON(rel, err)
}

func vipH5History(c *bm.Context) {
	var (
		err error
	)

	arg := new(model.ArgChangeHistory)
	if c.Bind(arg); err != nil {
		log.Error("vipH5History(%d) err(%+v)", arg.Mid, err)
		return
	}

	c.JSON(vipSvc.H5History(c, arg))
}

// vipInfo (for old service).
func vipInfo(c *bm.Context) {
	var (
		res *model.VipInfoBoResp
		err error
	)
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.VipInfoBo(c, arg.Mid); err != nil {
		log.Error("vipSvc.VipInfo(%d) err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
