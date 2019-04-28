package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upRating(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		log.Error("error get mid")
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	rating, err := svc.UpRating(c, mid)
	if err != nil {
		log.Error("svc.UpRating mid(%v) err(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rating, err)
}
