package http

import (
	searchMdl "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_orderAll  = "totalrank"
	_orderDate = "pubdate"
)

// searchTypes returns the ugc types
func searchTypes(c *bm.Context) {
	c.JSON(tvSvc.SearchTypes())
}

// searchResult returns the result of search
func searchResult(c *bm.Context) {
	v := new(searchMdl.ReqSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Order != "" && v.Order != _orderAll && v.Order != _orderDate {
		c.JSON(nil, ecode.TvSearchOrderWrong)
	}
	c.JSON(secSvc.SearchRes(c, v))
}

// searchSug search suggestion words
func searchSug(c *bm.Context) {
	v := new(searchMdl.ReqSug)
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := secSvc.SearchSug(c, v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"result":  data.Result,
		"message": "success",
		"cost":    data.Cost,
		"stoken":  data.Stoken,
	}, nil)
}

func pgcIdx(c *bm.Context) {
	v := new(searchMdl.ReqPgcIdx)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSvc.EsPgcIdx(c, v))
}

func ugcIdx(c *bm.Context) {
	v := new(searchMdl.ReqUgcIdx)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSvc.EsUgcIdx(c, v))
}
