package http

import (
	"net/http"

	bm "go-common/library/net/http/blademaster"
)

// get user get a image.
func get(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Token string `form:"token"  validate:"required"`
			Bid   string `form:"bid"  validate:"required"`
		})
		img []byte
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if img, err = svr.CaptchaImg(c, v.Token, v.Bid); err != nil {
		c.JSON(nil, err)
		return
	}
	code := http.StatusOK
	c.Render(code, Image{
		Body: img,
	})
}

// token third business get token.
func token(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Bid string `form:"bid"  validate:"required"`
		})
		token, url string
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if url, token, err = svr.Token(c, v.Bid); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 1)
	data["data"] = map[string]string{
		"token": token,
		"url":   url,
	}
	c.JSONMap(data, nil)
}

// verify third business verify.
func verify(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Token string `form:"token"  validate:"required"`
			Code  string `form:"code"  validate:"required"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	err = svr.VerifyCaptcha(c, v.Token, v.Code)
	c.JSON(nil, err)
}
