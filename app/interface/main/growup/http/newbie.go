package http

import (
	"go-common/app/interface/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/net/http/blademaster"
)

func upNewbieLetter(c *blademaster.Context) {
	iMid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}

	req := new(model.NewbieLetterReq)
	if err := c.Bind(req); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	req.Mid = iMid.(int64)
	res, err := newbieSvr.Letter(c, req)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
