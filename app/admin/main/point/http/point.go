package http

import (
	"go-common/app/admin/main/point/model"
	pointmol "go-common/app/service/main/point/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func pointConfList(c *bm.Context) {
	var (
		err error
		res []*model.PointConf
	)
	if res, err = svc.PointConfList(c); err != nil {
		return
	}
	c.JSON(&model.PageInfo{Count: len(res), Item: res}, nil)
}

func pointConfInfo(c *bm.Context) {
	var (
		err error
		res *model.PointConf
	)
	arg := &model.ArgID{}
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = svc.PointCoinInfo(c, arg.ID); err != nil {
		log.Error("svc.PointCoinInfo(%d), err(%+v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func pointConfAdd(c *bm.Context) {
	var (
		err error
	)
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	pc := &model.PointConf{}
	if err = c.Bind(pc); err != nil {
		return
	}
	pc.Operator = opI.(string)
	if _, err = svc.PointCoinAdd(c, pc); err != nil {
		log.Error("svc.PointCoinAdd(%+v), err(%+v)", pc, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func pointConfEdit(c *bm.Context) {
	var (
		err error
	)
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	pc := &model.PointConf{}
	if err = c.Bind(pc); err != nil {
		return
	} else if pc.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pc.Operator = opI.(string)
	if err = svc.PointCoinEdit(c, pc); err != nil {
		log.Error("svc.PointCoinEdit(%+v), err(%+v)", pc, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func pointHistory(c *bm.Context) {
	var err error
	arg := &model.ArgPointHistory{}
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(svc.PointHistory(c, arg))
}

func pointUserAdd(c *bm.Context) {
	var (
		err error
	)
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	req := new(model.ArgPoint)
	if err = c.Bind(req); err != nil {
		log.Error("point add bind %+v", err)
		return
	}
	if req.Point <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg := new(pointmol.ArgPoint)
	arg.Mid = req.Mid
	arg.Point = req.Point
	arg.Remark = req.Remark
	arg.Operator = opI.(string)
	arg.ChangeType = model.PointSystem
	if err = svc.PointAdd(c, arg); err != nil {
		log.Error("point add(%+v) faild(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
