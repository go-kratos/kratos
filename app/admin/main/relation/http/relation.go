package http

import (
	"time"

	"go-common/app/admin/main/relation/model"
	bm "go-common/library/net/http/blademaster"
)

func followers(c *bm.Context) {
	params := &model.FollowersParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	if params.PN <= 0 {
		params.PN = 1
	}
	if params.PS <= 0 {
		params.PS = 50
	}
	params.Order = "mtime"
	if params.Sort != "desc" {
		params.Sort = "asc"
	}

	c.JSON(svc.Followers(c, params))
}

func followings(c *bm.Context) {
	params := &model.FollowingsParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	if params.PN <= 0 {
		params.PN = 1
	}
	if params.PS <= 0 {
		params.PS = 50
	}
	params.Order = "mtime"
	if params.Sort != "desc" {
		params.Sort = "asc"
	}

	c.JSON(svc.Followings(c, params))
}

func logs(c *bm.Context) {
	params := &model.LogsParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	now := time.Now()
	from := time.Unix(0, 0)
	c.JSON(svc.RelationLog(c, params.Mid, params.Fid, from, now))
}

func stat(ctx *bm.Context) {
	params := &model.ArgMid{}
	if err := ctx.Bind(params); err != nil {
		return
	}
	ctx.JSON(svc.Stat(ctx, params))
}

func stats(ctx *bm.Context) {
	params := &model.ArgMids{}
	if err := ctx.Bind(params); err != nil {
		return
	}
	ctx.JSON(svc.Stats(ctx, params))
}
