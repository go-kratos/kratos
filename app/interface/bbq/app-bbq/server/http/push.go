package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func pushRegister(c *bm.Context) {
	args := &v1.PushRegisterRequest{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	var mid int64
	if tmp, ok := c.Get("mid"); ok {
		mid = tmp.(int64)
	}
	buvid := c.Request.Header.Get("Buvid")
	req := &notice.UserPushDev{
		Mid:        mid,
		RegisterId: args.RegID,
		Buvid:      buvid,
		Platform:   1,
		Sdk:        1,
	}
	if args.Platform == "ios" {
		req.Platform = 2
	}
	c.JSON(srv.PushRegister(c, req))
	// 埋点
	uiLog(c, model.ActionPushRegister, args)
}

func pushLogout(c *bm.Context) {
	args := &notice.UserPushDev{}
	var mid int64
	if tmp, ok := c.Get("mid"); ok {
		mid = tmp.(int64)
	}
	args.Mid = mid
	buvid := c.Request.Header.Get("Buvid")
	args.Buvid = buvid
	c.JSON(srv.PushLogout(c, args))
}

func pushCallback(c *bm.Context) {
	args := &v1.PushCallbackRequest{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	var mid int64
	if tmp, ok := c.Get("mid"); ok {
		mid = tmp.(int64)
	}
	buvid := c.Request.Header.Get("Buvid")

	c.JSON(srv.PushCallback(c, args, mid, buvid))
	// 埋点
	uiLog(c, model.ActionPushCallback, args)
}
