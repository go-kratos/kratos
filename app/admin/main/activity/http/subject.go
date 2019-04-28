package http

import (
	"go-common/app/admin/main/activity/model"
	bm "go-common/library/net/http/blademaster"
)

func listInfosAll(c *bm.Context) {
	arg := new(model.ListSub)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.SubjectList(c, arg))
}

func videoList(c *bm.Context) {
	c.JSON(actSrv.VideoList(c))
}

func addActSubject(c *bm.Context) {
	arg := new(model.AddList)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.AddActSubject(c, arg))
}

func updateInfoAll(c *bm.Context) {
	type upStr struct {
		model.AddList
		Sid int64 `form:"sid" validate:"min=1"`
	}
	arg := new(upStr)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.UpActSubject(c, &arg.AddList, arg.Sid))
}

func subPro(c *bm.Context) {
	type subStr struct {
		Sid int64 `form:"sid" validate:"min=1"`
	}
	arg := new(subStr)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.SubProtocol(c, arg.Sid))
}

func timeConf(c *bm.Context) {
	type subStr struct {
		Sid int64 `form:"sid" validate:"required"`
	}
	arg := new(subStr)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.TimeConf(c, arg.Sid))
}

func article(c *bm.Context) {
	type subStr struct {
		Aids []int64 `form:"aids,split" validate:"min=1,required"`
	}
	arg := new(subStr)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.GetArticleMetas(c, arg.Aids))
}
