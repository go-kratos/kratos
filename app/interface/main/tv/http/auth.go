package http

import (
	bm "go-common/library/net/http/blademaster"
)

func loadEP(c *bm.Context) {
	// bind
	v := new(struct {
		EPID  int64 `form:"epid" validate:"required,min=1"`
		Build int64 `form:"build"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	// take cache and distinguish the error msg
	if isok, errMsg, err := tvSvc.EpMsg(v.EPID, v.Build); err != nil {
		c.JSON(nil, err)
	} else {
		c.JSONMap(map[string]interface{}{
			"data":    isok,
			"message": errMsg,
		}, nil)
	}
}

func loadVideo(c *bm.Context) {
	// bind
	v := new(struct {
		CID int64 `form:"cid" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	// take cache and distinguish the error msg
	if isok, errMsg, err := viewSvc.VideoMsg(c, v.CID); err != nil {
		c.JSON(nil, err)
	} else {
		c.JSONMap(map[string]interface{}{
			"data":    isok,
			"message": errMsg,
		}, nil)
	}
}
