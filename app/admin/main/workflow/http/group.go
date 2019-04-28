package http

import (
	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func setGroupStateV3(ctx *bm.Context) {
	gssp := &param.GroupStateSetParam{}
	if err := ctx.BindWith(gssp, binding.FormPost); err != nil {
		return
	}
	gssp.AdminID, gssp.AdminName = adminInfo(ctx)

	// check ban account operate    账号封禁不支持批处理
	if len(gssp.ID) > 1 && gssp.BlockDay != 0 {
		ctx.JSON(nil, ecode.WkfBanNotSupportBatchOperate)
		return
	}

	ctx.JSON(nil, wkfSvc.SetGroupState(ctx, gssp))
}

func groupListV3(ctx *bm.Context) {
	v := new(param.GroupListParamV3)
	if err := ctx.Bind(v); err != nil {
		return
	}
	gscc := &search.GroupSearchCommonCond{
		Fields:       []string{"id", "oid", "typeid", "mid", "eid", "report_mid", "title", "first_user_tid"},
		Business:     v.Business,
		Oids:         v.Oid,
		Mids:         v.Mid,
		States:       v.State,
		TypeIDs:      v.TypeID,
		Rounds:       v.Round,
		RID:          v.Rid,
		FID:          v.Fid,
		EID:          v.Eid,
		Tids:         v.Tid,
		FirstUserTid: v.FirstUserTid,
		Order:        v.Order,
		Sort:         v.Sort,
		PN:           v.PN,
		PS:           v.PS,
		KWPriority:   v.KWPriority,
		KW:           v.KW,
		KWFields:     v.KWField,
		CTimeFrom:    v.CTimeFrom,
		CTimeTo:      v.CTimeTo,
		ReportMID:    v.ReportMid,
	}
	ctx.JSON(wkfSvc.GroupListV3(ctx, gscc))
}

func setGroupRole(ctx *bm.Context) {
	grsp := &param.GroupRoleSetParam{}
	if err := ctx.BindWith(grsp, binding.FormPost); err != nil {
		return
	}
	grsp.AdminID, grsp.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpGroupRole(ctx, grsp))
}

func upGroupExtra(ctx *bm.Context) {
	uep := &param.UpExtraParam{}
	if err := ctx.BindWith(uep, binding.Form); err != nil {
		return
	}
	uep.AdminID, uep.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpGroupExtra(ctx, uep))
}

func setPublicReferee(ctx *bm.Context) {
	gspr := &param.GroupStatePublicReferee{}
	if err := ctx.BindWith(gspr, binding.FormPost); err != nil {
		return
	}
	// if bid support public judge
	if gspr.Business != model.CommentComplain {
		ctx.JSON(nil, ecode.WkfBidNotSupportPublicReferee)
	}
	gspr.AdminID, gspr.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.SetPublicReferee(ctx, gspr))
}

func countPendingGroup(ctx *bm.Context) {
	gpp := &param.GroupPendingParam{}
	if err := ctx.Bind(gpp); err != nil {
		return
	}
	gscc := &search.GroupSearchCommonCond{
		Fields:   []string{"id"},
		Business: gpp.Business,
		RID:      gpp.Rid,
		States:   []int8{model.Pending},
		PS:       1,
		PN:       1,
		Order:    "id",
		Sort:     "desc",
	}
	ctx.JSON(wkfSvc.GroupPendingCount(ctx, gscc))
}
