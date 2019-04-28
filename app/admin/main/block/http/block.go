package http

import (
	"go-common/app/admin/main/block/model"
	bm "go-common/library/net/http/blademaster"
)

func blockSearch(c *bm.Context) {
	var (
		err error
		v   = &model.ParamSearch{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(svc.Search(c, v.MIDs))
}

func blockHistory(c *bm.Context) {
	var (
		err error
		v   = &model.ParamHistory{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	var ret struct {
		Status  model.BlockStatus     `json:"status"`
		Total   int                   `json:"total"`
		History []*model.BlockHistory `json:"history"`
	}
	if ret.Status, ret.Total, ret.History, err = svc.History(c, v.MID, v.PS, v.PN); err != nil {
		c.JSON(nil, err)
		return
	}
	if ret.History == nil {
		ret.History = make([]*model.BlockHistory, 0)
	}
	c.JSON(ret, err)
}

func batchBlock(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchBlock{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(nil, svc.BatchBlock(c, v))
}

func batchRemove(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchRemove{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(nil, svc.BatchRemove(c, v))
}
