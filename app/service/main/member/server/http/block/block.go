package block

import (
	"time"

	model "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func info(c *bm.Context) {
	var (
		err error
		v   = &model.ParamInfo{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	var infos []*model.BlockInfo
	if infos, err = svc.Infos(c, []int64{v.MID}); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(infos) != 1 {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(infos[0], nil)
}

func batchInfo(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchInfo{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(svc.Infos(c, v.MIDs))
}

func batchDetail(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchDetail{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(svc.UserDetails(c, v.MIDs))
}

func block(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBlock{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	duration := time.Duration(v.Duration) * time.Second
	c.JSON(nil, svc.Block(c, []int64{v.MID}, v.Source, v.Area, v.Action, v.StartTime, duration, v.Operator, v.Reason, v.Comment, v.Notify))
}

func batchBlock(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchBlock{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	duration := time.Duration(v.Duration) * time.Second
	c.JSON(nil, svc.Block(c, v.MIDs, v.Source, v.Area, v.Action, v.StartTime, duration, v.Operator, v.Reason, v.Comment, v.Notify))
}

func remove(c *bm.Context) {
	var (
		err error
		v   = &model.ParamRemove{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(nil, svc.Remove(c, []int64{v.MID}, v.Source, model.BlockAreaNone, v.Operator, v.Reason, v.Comment, v.Notify))
}

func batchRemove(c *bm.Context) {
	var (
		err error
		v   = &model.ParamBatchRemove{}
	)
	if err = bind(c, v); err != nil {
		return
	}
	c.JSON(nil, svc.Remove(c, v.MIDs, v.Source, model.BlockAreaNone, v.Operator, v.Reason, v.Comment, v.Notify))
}
