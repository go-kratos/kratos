package http

import (
	"go-common/app/admin/main/credit/model/blocked"
	bm "go-common/library/net/http/blademaster"
)

func httpData(c *bm.Context, data interface{}, pager *blocked.Pager) {
	res := make(map[string]interface{})
	if data == nil {
		data = struct{}{}
	}
	if pager == nil {
		pager = &blocked.Pager{}
	}
	res["data"] = data
	res["pager"] = &blocked.Pager{
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
