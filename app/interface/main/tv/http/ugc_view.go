package http

import (
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// viewPage view page handler.
func view(c *bm.Context) {
	var (
		v = new(struct {
			AID       int64  `form:"aid" validate:"required,min=1"`
			AccessKey string `form:"access_key"`
		})
		mid int64
	)
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// view
	now := time.Now()
	view, isok, errMsg, err := viewSvc.View(c, mid, v.AID, v.AccessKey, metadata.String(c, metadata.RemoteIP), now)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// err msg logic
	if !isok {
		c.JSONMap(map[string]interface{}{
			"data":    isok,
			"message": errMsg,
		}, ecode.CopyrightLimit)
		return
	}
	c.JSON(view, nil)
}
