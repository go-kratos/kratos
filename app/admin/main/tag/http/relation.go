package http

import (
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func relationList(c *bm.Context) {
	var (
		err   error
		total int64
		param = new(model.ParamRelationList)
		res   []*model.Resource
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if param.Type == model.QueryTypeTName {
		if param.TName == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if total, res, err = svc.RelationListByTag(c, param.TName, param.Pn, param.Ps); err != nil {
			c.JSON(nil, err)
			return
		}
	} else if param.Type == model.QueryTypeOid {
		if param.Oid <= 0 || param.OidType < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if total, res, err = svc.RelationListByOid(c, param.Oid, param.OidType, param.Pn, param.Ps); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	data := make(map[string]interface{}, 2)
	data["relation"] = res
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	c.JSON(data, nil)
}

func relationAdd(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamRelationAdd)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	mid, _ := managerInfo(c)
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.RelationAdd(c, param.TName, param.Oid, mid, param.Type))
}
func relationLock(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamRelation)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.RelationLock(c, param.Tid, param.Oid, param.Type))
}

func relationUnLock(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamRelation)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.RelationUnLock(c, param.Tid, param.Oid, param.Type))
}

func relationDelete(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamRelation)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.RelationDelete(c, param.Tid, param.Oid, param.Type))
}

func regionArcRefresh(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Rid int64 `form:"rid" validate:"required,gt=0"`
			Tid int64 `form:"tid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.RegionTagArcsRefresh(c, param.Rid, param.Tid))
}
