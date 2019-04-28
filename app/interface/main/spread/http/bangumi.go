package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func bangumiContent(c *bm.Context) {
	v := new(struct {
		Pn     int    `form:"pn" default:"1"`
		Ps     int    `form:"ps" validate:"min=1,max=1000" default:"20"`
		Type   int8   `form:"type" validate:"required"`
		Appkey string `form:"appkey" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := svc.BangumiContent(c, v.Pn, v.Ps, v.Type, v.Appkey)
	Render(c, data.Code, data.Message, data.Result, data.Total, err)
}

func bangumiOff(c *bm.Context) {
	v := new(struct {
		Pn        int    `form:"pn" default:"1"`
		Ps        int    `form:"ps" validate:"min=1,max=1000" default:"20"`
		Type      int8   `form:"type" validate:"required"`
		Timestamp int64  `form:"timestamp"`
		Appkey    string `form:"appkey" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	data, err := svc.BangumiOff(c, v.Pn, v.Ps, v.Type, v.Timestamp, v.Appkey)
	Render(c, data.Code, data.Message, data.Result, data.Total, err)
}
