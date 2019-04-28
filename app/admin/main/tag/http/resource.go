package http

import (
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func resourceLogList(c *bm.Context) {
	var (
		err   error
		total int64
		res   []*model.ResTagLog
		param = new(model.ParamResLogList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Role < model.ResRoleALL {
		param.Role = model.ResRoleALL
	}
	if param.Action < model.ResTagALL {
		param.Action = model.ResTagALL
	}
	if param.Pn < 1 {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if res, total, err = svc.ResourceLogs(c, param.Oid, param.TP, param.Role, param.Action, param.Pn, param.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["logs"] = res
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	c.JSON(data, nil)
}

func resourceLogState(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamResLogState)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.UpdateResLogState(c, param.ID, param.Oid, param.TP, param.State))
}

func resourceLimit(c *bm.Context) {
	var (
		err      error
		total    int64
		resource []*model.LimitRes
		limitRes *model.LimitRes
		param    = new(model.ParamResLimit)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.TP == model.QuerryByResInfo {
		if param.Oid <= 0 || param.OidType <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if limitRes, err = svc.ResourceByOid(c, param.Oid, param.OidType); err != nil {
			c.JSON(nil, err)
			return
		}
		total = 1
		param.Ps = 1
		param.Pn = 1
		resource = append(resource, limitRes)
	} else if param.TP == model.QuerryByLimitState {
		if param.LimitState < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if param.Pn <= 0 {
			param.Pn = model.DefaultPageNum
		}
		if param.Ps <= 0 {
			param.Ps = model.DefaultPagesize
		}
		if resource, total, err = svc.ResByOperate(c, param.LimitState, param.Pn, param.Ps); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	data := make(map[string]interface{}, 2)
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	data["resource"] = resource
	c.JSON(data, nil)
}

func resourceLimitState(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamResLimitState)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.UpdateResLimitState(c, param.Oid, param.Type, param.Operate))
}
