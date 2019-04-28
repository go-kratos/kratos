package http

import (
	pb "go-common/app/service/main/sms/api"
	bm "go-common/library/net/http/blademaster"
)

func addTemplate(ctx *bm.Context) {
	req := new(pb.AddTemplateReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(svc.AddTemplate(ctx, req))
}

func updateTemplate(ctx *bm.Context) {
	req := new(pb.UpdateTemplateReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(svc.UpdateTemplate(ctx, req))
}

func templateList(ctx *bm.Context) {
	req := new(pb.TemplateListReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	res, err := svc.TemplateList(ctx, req)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	pager := struct {
		Pn    int32 `json:"page"`
		Ps    int32 `json:"pagesize"`
		Total int32 `json:"total"`
	}{
		Pn:    req.Pn,
		Ps:    req.Ps,
		Total: res.Total,
	}
	data := map[string]interface{}{
		"data":  res.List,
		"pager": pager,
	}
	ctx.JSONMap(data, nil)
}
