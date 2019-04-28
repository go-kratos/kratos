package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

//queryContactLogs queryContactLogs
func queryContactLogs(c *bm.Context) {
	v := &model.QueryContactLogRequest{}
	if err := c.Bind(v); err != nil {
		return
	}

	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(srv.QueryContactLogs(c, v))
}

//queryContactLogs queryContactLogs
func queryRedisdata(c *bm.Context) {
	v := &model.QueryContactLogRequest{}
	if err := c.Bind(v); err != nil {
		return
	}

	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(srv.QueryContactLogs(c, v))
}

// SyncContacts sync the wechat contacts 更新企业微信列表（用户信息和saga信息）
func queryContacts(ctx *bm.Context) {
	var (
		req = &model.EmptyReq{}
		v   = &model.Pagination{}
		err error
	)
	if err = ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if err = ctx.Bind(v); err != nil {
		return
	}
	ctx.JSON(srv.QueryContacts(ctx, v))
}

func createWechat(ctx *bm.Context) {
	var (
		username string
		err      error
	)
	req := &model.CreateChatReq{}
	if err = ctx.BindWith(req, binding.JSON); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if username, err = getUsername(ctx); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.CreateWechat(ctx, req, username))
}

func queryWechatCreateLog(ctx *bm.Context) {
	var (
		req  = &model.Pagination{}
		err  error
		user string
	)
	if err = ctx.Bind(req); err != nil {
		return
	}
	if err = req.Verify(); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if user, err = getUsername(ctx); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.QueryWechatCreateLog(ctx, req, user))
}

func getWechat(ctx *bm.Context) {
	ctx.JSON(srv.WechatParams(ctx, ctx.Request.Form.Get("chatid")))
}

func sendGroupWechat(ctx *bm.Context) {
	req := &model.SendChatReq{}
	if err := ctx.BindWith(req, binding.JSON); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.SendGroupWechat(ctx, req))
}

func sendWechat(ctx *bm.Context) {
	req := &model.SendMessageReq{}
	if err := ctx.BindWith(req, binding.JSON); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.SendWechat(ctx, req))
}

func updateWechat(ctx *bm.Context) {
	req := &model.UpdateChatReq{}
	if err := ctx.BindWith(req, binding.JSON); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.UpdateWechat(ctx, req))
}

func syncWechatContacts(ctx *bm.Context) {
	ctx.JSON(srv.SyncWechatContacts(ctx))
}
