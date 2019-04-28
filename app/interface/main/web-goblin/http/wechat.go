package http

import (
	"net/http"

	bm "go-common/library/net/http/blademaster"
)

func qrcode(c *bm.Context) {
	v := new(struct {
		JSON string `form:"json" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	qrcode, err := srvWechat.Qrcode(c, v.JSON)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Bytes(http.StatusOK, "", qrcode)
}
