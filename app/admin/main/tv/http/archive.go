package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

//arcOnline archive online
func arcOnline(c *bm.Context) {
	arcAction(c, 1)
}

func arcHidden(c *bm.Context) {
	arcAction(c, 2)
}

func arcAction(c *bm.Context, action int) {
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
	if err := tvSrv.ArcAction(param.IDs, action); err != nil {
		res["message"] = "更新数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON("成功", nil)
}

// archive list repository
func arcList(c *bm.Context) {
	var (
		res   = make(map[string]interface{})
		param = new(model.ArcListParam)
	)
	if err := c.Bind(param); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pager, err := tvSrv.ArchiveList(c, param); err != nil {
		res["message"] = "获取数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
	} else {
		c.JSON(pager, nil)
	}
}

//arcCategory archive category
func arcCategory(c *bm.Context) {
	c.JSON(tvSrv.GetTps(c, true))
}

// auditCategory gets audit consult used categorys
func auditCategory(c *bm.Context) {
	c.JSON(tvSrv.GetTps(c, false))
}

//arcTypeRPC get archive type from rpc
func arcTypeRPC(c *bm.Context) {
	c.JSON(tvSrv.ArcTypes, nil)
}

func arcUpdate(c *bm.Context) {
	param := new(struct {
		ID      int64  `form:"id" validate:"required"`
		Cover   string `form:"cover" validate:"required"`
		Content string `form:"content" validate:"required"`
		Title   string `form:"title" validate:"required"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.ArcUpdate(param.ID, param.Cover, param.Content, param.Title))
}

func unShelve(c *bm.Context) {
	var (
		username string
		param    = new(model.ReqUnshelve)
	)
	if err := c.Bind(param); err != nil {
		return
	}
	if un, ok := c.Get("username"); ok {
		username = un.(string)
	} else {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	c.JSON(tvSrv.Unshelve(c, param, username))
}
