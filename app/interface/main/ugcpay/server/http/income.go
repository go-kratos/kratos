package http

import (
	"strconv"

	api "go-common/app/interface/main/ugcpay/api/http"
	"go-common/app/interface/main/ugcpay/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func incomeAssetOverview(ctx *bm.Context) {
	var (
		err    error
		resp   *api.RespIncomeAssetOverview
		inc    *model.IncomeAssetOverview
		mid, _ = ctx.Get("mid")
	)
	if inc, err = srv.IncomeAssetOverview(ctx, mid.(int64)); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp = &api.RespIncomeAssetOverview{
		Total:         inc.Total,
		TotalBuyTimes: inc.TotalBuyTimes,
		MonthNew:      inc.MonthNew,
		DayNew:        inc.DayNew,
	}
	ctx.JSON(resp, err)
}

func incomeAssetMonthly(ctx *bm.Context) {
	var (
		err    error
		arg    = &api.ArgIncomeAssetList{}
		resp   = &api.RespIncomeAssetList{List: make([]*api.RespIncomeAsset, 0)}
		inc    *model.IncomeAssetMonthly
		page   *model.Page
		ver    int64
		mid, _ = ctx.Get("mid")
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.PS == 0 {
		arg.PS = 20
	}
	if arg.PN == 0 {
		arg.PN = 1
	}
	if arg.Ver != "" {
		if ver, err = strconv.ParseInt(arg.Ver, 10, 64); err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		// ver=0 代表总计
		ver = 0
	}

	if inc, page, err = srv.IncomeAssetList(ctx, mid.(int64), ver, arg.PS, arg.PN); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.Page = api.RespPage{
		Num:   page.Num,
		Size:  page.Size,
		Total: page.Total,
	}
	for _, i := range inc.List {
		resp.List = append(resp.List, &api.RespIncomeAsset{
			OID:           i.OID,
			OType:         i.OType,
			Title:         i.Title,
			Currency:      i.Currency,
			Price:         i.Price,
			TotalBuyTimes: i.TotalBuyTimes,
			NewBuyTimes:   i.NewBuyTimes,
			TotalErrTimes: i.TotalErrTimes,
			NewErrTimes:   i.NewErrTimes,
		})
	}
	ctx.JSON(resp, err)
}
