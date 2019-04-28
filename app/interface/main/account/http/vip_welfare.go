package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func welfareList(c *bm.Context) {
	arg := new(struct {
		Tid       int64 `form:"tid"`
		Recommend int64 `form:"recommend"`
		Pn        int64 `form:"pn"`
		Ps        int64 `form:"ps"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	c.JSON(vipSvc.WelfareList(c, arg.Tid, arg.Recommend, arg.Pn, arg.Ps))
}

func welfareTypeList(c *bm.Context) {
	c.JSON(vipSvc.WelfareTypeList(c))
}

func welfareInfo(c *bm.Context) {
	userId := int64(0)
	mid, exists := c.Get("mid")
	if exists {
		userId = mid.(int64)
	}
	arg := new(struct {
		Wid int64 `form:"id"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	c.JSON(vipSvc.WelfareInfo(c, arg.Wid, userId))
}

func receiveWelfare(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(struct {
		Wid int64 `form:"id"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	c.JSON(vipSvc.WelfareReceive(c, arg.Wid, mid.(int64)))
}

func myWelfare(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(vipSvc.MyWelfare(c, mid.(int64)))
}
