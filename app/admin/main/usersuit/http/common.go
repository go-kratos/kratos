package http

import (
	"go-common/app/admin/main/usersuit/model"
	bm "go-common/library/net/http/blademaster"
)

func httpData(c *bm.Context, data interface{}, pager *model.Pager) {
	res := make(map[string]interface{})
	if data == nil {
		data = struct{}{}
	}
	if pager == nil {
		pager = &model.Pager{}
	}
	res["data"] = data
	res["pager"] = &model.Pager{
		Total: pager.Total,
		PN:    pager.PN,
		PS:    pager.PS,
		Order: pager.Order,
		Sort:  pager.Sort,
	}
	c.JSONMap(res, nil)
}

func httpCode(c *bm.Context, err error) {
	c.JSON(nil, err)
}
