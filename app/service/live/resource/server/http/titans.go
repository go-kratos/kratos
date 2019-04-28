package http

import (
	"context"
	"go-common/app/service/live/resource/api/http/v1"
	"go-common/library/ecode"
	"go-common/library/net/http/blademaster"
	"time"
)

func getNodes(c *blademaster.Context) {
	res := map[string]interface{}{}
	res["data"] = ""
	cookie := c.Request.Header.Get("Cookie")
	team := c.Request.FormValue("team")
	node := c.Request.FormValue("node")
	username, err := c.Request.Cookie("username")
	if err != nil || cookie == "" || username == nil {
		err = ecode.Error(1, "cookie未获取到")
		c.JSONMap(res, err)
		return
	}
	ctx, cancel := context.WithTimeout(c, 800*time.Millisecond)
	defer cancel()
	sRes, err := titansService.GetMyTreeApps(ctx, &v1.TreeAppsReq{
		Team: team,
		Node: node,
	}, cookie, username.Value)
	res["msg"] = ""
	res["message"] = ""
	res["data"] = sRes
	c.JSONMap(res, err)
}
