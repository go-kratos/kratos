package http

import bm "go-common/library/net/http/blademaster"

func timemachine2018(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	loginMid := midStr.(int64)
	c.JSON(tmSvc.Timemachine2018(c, loginMid, v.Mid))
}

func timemachine2018Raw(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	loginMid := midStr.(int64)
	c.JSON(tmSvc.Timemachine2018Raw(c, loginMid, v.Mid))
}

func timemachine2018Cache(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	loginMid := midStr.(int64)
	c.JSON(tmSvc.Timemachine2018Cache(c, loginMid, v.Mid))
}

//func startTmProc(c *bm.Context) {
//	c.JSON(nil, tmSvc.StartTmproc(c))
//}
//
//func stopTmProc(c *bm.Context) {
//	c.JSON(nil, tmSvc.StopTmproc(c))
//}
