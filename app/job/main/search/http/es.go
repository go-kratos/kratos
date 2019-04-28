package http

import (
	"strconv"

	"go-common/app/job/main/search/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func action(c *bm.Context) {
	var (
		params           = c.Request.Form
		recoverID        int64
		writeEntityIndex bool
	)
	appid := params.Get("appid")
	action := params.Get("action")
	if appid == "" || action == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !model.ExistsAction[action] {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Get("recover_id") != "" {
		if rid, err := strconv.ParseInt(params.Get("recover_id"), 10, 64); err == nil {
			recoverID = rid
		}
	}
	if params.Get("entity") == "1" {
		writeEntityIndex = true
	} else {
		writeEntityIndex = false
	}
	c.JSON(svr.HTTPAction(ctx, appid, action, recoverID, writeEntityIndex))
}

func stat(c *bm.Context) {
	var (
		params = c.Request.Form
	)
	appid := params.Get("appid")
	if appid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Stat(ctx, appid))
}
