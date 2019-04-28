package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

var (
	_username = "shaozhenyu"
)

func autoBreach(c *bm.Context) {
	var err error
	v := new(struct {
		Type   int     `form:"type"`
		AIDs   []int64 `form:"aids,split" validate:"required"`
		MID    int64   `form:"mid" validate:"required"`
		Reason string  `form:"reason" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}

	err = incomeSvr.ArchiveBreach(c, v.Type, v.AIDs, v.MID, v.Reason, _username)
	if err != nil {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func autoDismiss(c *bm.Context) {
	var err error
	v := new(struct {
		Type   int    `form:"type"`
		MID    int64  `form:"mid" validate:"required"`
		Reason string `form:"reason" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	err = svr.AutoDismiss(c, _username, v.Type, v.MID, v.Reason)
	if err != nil {
		log.Error("growup svr.Dismiss error(%v)", err)
	}
	c.JSON(nil, err)
}

func autoForbid(c *bm.Context) {
	var err error
	v := new(struct {
		Type   int    `form:"type"`
		MID    int64  `form:"mid" validate:"required"`
		Reason string `form:"reason" validate:"required"`
		Days   int    `form:"days"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	err = svr.AutoForbid(c, _username, v.Type, v.MID, v.Reason, v.Days, v.Days*86400)
	if err != nil {
		log.Error("growup svr.Forbid error(%v)", err)
	}
	c.JSON(nil, err)
}
