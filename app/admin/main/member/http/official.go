package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"go-common/app/admin/main/member/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func officials(ctx *bm.Context) {
	arg := &model.ArgOfficial{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	os, total, err := svc.Officials(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"officials": os,
		"page": map[string]int{
			"num":   arg.Pn,
			"size":  arg.Ps,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func officialsExcel(ctx *bm.Context) {
	arg := &model.ArgOfficial{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	os, _, err := svc.Officials(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}

	data := make([][]string, 0, len(os))
	data = append(data, []string{"完成认证时间", "用户ID", "昵称", "认证类型", "认证称号", "称号后缀"})
	for _, of := range os {
		fields := []string{
			of.CTime.Time().Format("2006-01-02 15:04:05"),
			strconv.FormatInt(of.Mid, 10),
			of.Name,
			model.OfficialRoleName(of.Role),
			of.Title,
			of.Desc,
		}
		data = append(data, fields)
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			ctx.JSON(nil, err)
			return
		}
	}
	w.Flush()
	res := buf.Bytes()
	ctx.Writer.Header().Set("Content-Type", "application/csv")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", "官方认证名单"))
	ctx.Writer.Write([]byte("\xEF\xBB\xBF")) // 写 UTF-8 的 BOM 头，别删，删了客服就找上门了
	ctx.Writer.Write(res)
}

func officialDocs(ctx *bm.Context) {
	arg := &model.ArgOfficialDoc{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	ods, total, err := svc.OfficialDocs(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"officials": ods,
		"page": map[string]int{
			"num":   arg.Pn,
			"size":  arg.Ps,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func officialDocsExcel(ctx *bm.Context) {
	arg := &model.ArgOfficialDoc{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	ods, _, err := svc.OfficialDocs(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	data := make([][]string, 0, len(ods))
	data = append(data, []string{"申请时间", "认证类型", "用户ID", "用户昵称", "审核状态", "认证称号", "称号后缀", "操作人"})
	for _, od := range ods {
		fields := []string{
			od.CTime.Time().Format("2006-01-02 15:04:05"),
			model.OfficialRoleName(od.Role),
			strconv.FormatInt(od.Mid, 10),
			od.Name,
			model.OfficialStateName(od.State),
			od.Title,
			od.Desc,
			od.Uname,
		}
		data = append(data, fields)
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			ctx.JSON(nil, err)
			return
		}
	}
	w.Flush()
	res := buf.Bytes()
	ctx.Writer.Header().Set("Content-Type", "application/csv")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", "官方认证审核记录"))
	ctx.Writer.Write([]byte("\xEF\xBB\xBF")) // 写 UTF-8 的 BOM 头，别删，删了客服就找上门了
	ctx.Writer.Write(res)
}

func officialDoc(ctx *bm.Context) {
	arg := &model.ArgMid{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	od, logs, block, spys, rn, sameCreditCodeMids, err := svc.OfficialDoc(ctx, arg.Mid)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"official":              od,
		"logs":                  logs.Result,
		"block":                 block,
		"spys":                  spys,
		"realname":              rn,
		"same_credit_code_mids": sameCreditCodeMids,
	}
	ctx.JSON(res, nil)
}

func officialDocAudit(ctx *bm.Context) {
	arg := &model.ArgOfficialAudit{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.State == model.OfficialStateNoPass && arg.Reason == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if arg.Source != "" {
		arg.Reason = fmt.Sprintf("%s（来源：%s）", arg.Reason, arg.Source)
	}
	ctx.JSON(nil, svc.OfficialDocAudit(ctx, arg))
}

func officialDocEdit(ctx *bm.Context) {
	arg := &model.ArgOfficialEdit{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.OfficialDocEdit(ctx, arg))
}

func officialDocSubmit(ctx *bm.Context) {
	arg := &model.ArgOfficialSubmit{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.OfficialDocSubmit(ctx, arg))
}
