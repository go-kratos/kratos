package http

import (
	"encoding/json"

	bm "go-common/library/net/http/blademaster"
)

func sudo(ctx *bm.Context) {
	cmd := ctx.Request.Form.Get("cmd")
	if cmd == "" {
		ctx.AbortWithStatus(400)
		return
	}
	ctx.Set("command", cmd)
}

func notityPurgeCache(ctx *bm.Context) {
	cmd, ok := ctx.Get("command")
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}
	plain, ok := cmd.(string)
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}

	var param struct {
		Mid    int64  `json:"mid"`
		Action string `json:"action"`
	}
	if err := json.Unmarshal([]byte(plain), &param); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	if param.Mid <= 0 {
		ctx.AbortWithStatus(400)
		return
	}
	if param.Action == "" {
		param.Action = "updateByAdmin"
	}
	ctx.JSON(nil, memberSvc.NotityPurgeCache(ctx, param.Mid, param.Action))
}
