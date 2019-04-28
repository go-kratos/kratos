package http

import (
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func actionCount(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			RpID    int64 `form:"rpid" validate:"required"`
			Oid     int64 `form:"oid" validate:"required"`
			Type    int32 `form:"type" validate:"required"`
			AdminID int64 `form:"admin_id"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	like, hate, err := rpSvc.ActionCount(c, v.RpID, v.Oid, v.AdminID, v.Type)
	if err != nil {
		log.Warn("svc.ActionInfo(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"like": like,
		"hate": hate,
	}
	c.JSONMap(data, nil)
}

func actionUpdate(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			RpID    int64  `form:"rpid" validate:"required"`
			Oid     int64  `form:"oid" validate:"required"`
			Type    int32  `form:"type" validate:"required"`
			Action  int32  `form:"action" validate:"required,min=1,max=2"`
			Count   int32  `form:"count"`
			AdminID int64  `form:"admin_id"`
			Remark  string `form:"remark"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	switch v.Action {
	case model.ActionLike:
		if v.Count < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		err = rpSvc.UpActionLike(c, v.RpID, v.Oid, adid, v.Type, v.Count, v.Remark)
	case model.ActionHate:
		err = rpSvc.UpActionHate(c, v.RpID, v.Oid, adid, v.Type, v.Count, v.Remark)
	default:
		err = ecode.RequestErr
	}
	if err != nil {
		log.Warn("rpSvc.ActionUpdate(%v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
