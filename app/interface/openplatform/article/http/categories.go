package http

import (
	artmdl "go-common/app/interface/openplatform/article/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func categories(c *bm.Context) {
	data, err := artSrv.ListCategories(c, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if data == nil {
		data = artmdl.Categories{}
	}
	c.JSON(data, nil)
}
