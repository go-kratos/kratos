package http

import (
	"net/http"
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func checkgroup() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		uid, _ := getUIDName(ctx)
		role, err := srv.CheckGroup(ctx, uid)
		if err != nil || role == 0 {
			data := map[string]interface{}{
				"code":    ecode.RequestErr,
				"message": "权限错误",
			}
			ctx.Render(http.StatusOK, render.MapJSON(data))
			ctx.Abort()
			return
		}
	}
}

// 校验任务操作权限
func checkowner() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		tidS := ctx.Request.Form.Get("task_id")
		tid, err := strconv.Atoi(tidS)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			ctx.Abort()
			return
		}

		uid, _ := getUIDName(ctx)
		if err := srv.CheckOwner(ctx, int64(tid), uid); err != nil {
			data := map[string]interface{}{
				"code":    ecode.RequestErr,
				"message": err.Error(),
			}
			ctx.Render(http.StatusOK, render.MapJSON(data))
			ctx.Abort()
			return
		}
	}
}
