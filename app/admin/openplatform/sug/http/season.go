package http

import (
	"go-common/app/admin/openplatform/sug/model"
	bm "go-common/library/net/http/blademaster"
)

func matchOperate(c *bm.Context) {
	v := &model.MatchOperate{}
	c.Bind(v)
	c.JSON(srv.MatchOperate(c, v))
}

func search(c *bm.Context) {
	v := &model.Search{}
	c.Bind(v)
	c.JSON(srv.Search(c, v))
}

func sourceSearch(c *bm.Context) {
	v := &model.SourceSearch{}
	c.Bind(v)
	c.JSON(srv.SourceSearch(c, v))
}
