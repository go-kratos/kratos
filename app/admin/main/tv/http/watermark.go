package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func waterMarklist(c *bm.Context) {
	var (
		err    error
		res    = map[string]interface{}{}
		pagers *model.WaterMarkListPager
	)
	param := new(model.WaterMarkListParam)
	if err = c.Bind(param); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pagers, err = tvSrv.WaterMarkist(c, param); err != nil {
		res["message"] = "获取数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}

	c.JSON(pagers, nil)
}

func waterMarkAdd(c *bm.Context) {
	var (
		err    error
		res    = map[string]interface{}{}
		addRes *model.AddEpIDResp
	)
	param := new(struct {
		IDs  []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
		Type string  `form:"type" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Type == "seasonid" {
		if addRes, err = tvSrv.AddSeasonID(c, param.IDs); err != nil {
			res["message"] = "添加数据失败!" + err.Error()
			c.JSONMap(res, ecode.RequestErr)
			return
		}
	} else if param.Type == "epid" {
		if addRes, err = tvSrv.AddEpID(c, param.IDs); err != nil {
			res["message"] = "添加数据失败!" + err.Error()
			c.JSONMap(res, ecode.RequestErr)
			return
		}
	}
	c.JSON(addRes, nil)
}

func waterMarkDelete(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if err = tvSrv.DeleteWatermark(param.IDs); err != nil {
		res["message"] = "删除数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func transList(c *bm.Context) {
	var param = new(model.TransReq)
	if err := c.Bind(param); err != nil {
		return
	}
	if data, err := tvSrv.TransList(c, param); err == nil || err == ecode.NothingFound {
		c.JSON(data, nil)
	} else {
		c.JSON(data, err)
	}
}
