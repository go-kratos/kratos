package http

import (
	http "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func shareURL(c *bm.Context) {
	var device *bm.Device
	if dev, _ := c.Get("device"); dev != nil {
		device = dev.(*bm.Device)
	}
	mid := int64(0)
	if v, _ := c.Get("mid"); v != nil {
		mid = v.(int64)
	}

	arg := &http.ShareRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(srv.GetShareURL(c, mid, device, arg))
}

func shareCallback(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	if mid == nil {
		mid = int64(0)
	}
	arg := &http.ShareCallbackRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := srv.ShareCallback(c, mid.(int64), dev.(*bm.Device), arg)
	c.JSON(resp, err)

	// 埋点
	if err != nil {
		return
	}
	ext := struct {
		Svid    int64
		Channel int32
	}{
		Svid:    arg.Svid,
		Channel: arg.Channel,
	}
	uiLog(c, model.ActionShare, ext)
}
