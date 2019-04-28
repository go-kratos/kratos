package http

import (
	"context"
	"net/http"
	"strconv"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/net/http/blademaster"
)

func recommendAdd(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.RecommendUpReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.RecommendUpReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.RecommendAdd(cont, arg.(*model.RecommendUpReq))
		},
		"recommendAdd")(c)
}

func recommendOP(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.RecommendStateOpReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.RecommendStateOpReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.RecommendOP(cont, arg.(*model.RecommendStateOpReq))
		},
		"recommendOP")(c)
}

func recommendList(c *blademaster.Context) {
	httpGetWriterByExport(
		new(model.MCNUPRecommendReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.RecommendList(cont, arg.(*model.MCNUPRecommendReq))
		},
		"recommendList")(c)
}
