package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(v1.ImgUploadRequest)
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(srv.Upload(c, mid.(int64), arg.Type))
}

func perUpload(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	req := new(v1.PreUploadRequest)
	if err := c.Bind(req); err != nil {
		log.Errorw(c, "event", "bind param err")
		return
	}
	c.JSON(srv.PreUpload(c, req, mid.(int64)))
}

func callBack(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	req := new(v1.CallBackRequest)
	if err := c.Bind(req); err != nil {
		log.Errorw(c, "event", "bind param err", err, "err")
		return
	}
	c.JSON(srv.CallBack(c, req, mid.(int64)))
}

func uploadCheck(c *bm.Context) {
	tmp, exists := c.Get("mid")
	if !exists || tmp.(int64) == 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}

	c.JSON(srv.VideoUploadCheck(c, tmp.(int64)))
}
func homeimg(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	req := new(v1.HomeImgRequest)
	if err := c.Bind(req); err != nil {
		log.Errorw(c, "event", "bind param err", err, "err")
		return
	}
	c.JSON(srv.HomeImg(c, req, mid.(int64)))
}
