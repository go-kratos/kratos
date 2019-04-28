package block

import (
	model "go-common/app/admin/main/member/model/block"
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

func history(c *bm.Context) {
	var (
		err error
		v   = &model.ParamHistory{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(svc.History(c, v.MID, v.PS, v.PN, v.Desc))
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
