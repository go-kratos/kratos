package dao

import (
	"context"

	archive "go-common/app/service/main/archive/api"
	ugcpay "go-common/app/service/main/ugcpay/api/grpc/v1"

	"go-common/app/interface/main/ugcpay/model"
)

// TradeCreate create trade order for mid
func (d *Dao) TradeCreate(ctx context.Context, platform string, mid int64, oid int64, otype string, currency string) (orderID string, payData string, err error) {
	var (
		req = &ugcpay.TradeCreateReq{
			Platform: platform,
			Mid:      mid,
			Oid:      oid,
			Otype:    otype,
			Currency: currency,
		}
		reply *ugcpay.TradeCreateResp
	)
	if reply, err = d.ugcpayAPI.TradeCreate(ctx, req); err != nil {
		return
	}
	orderID = reply.OrderId
	payData = reply.PayData
	return
}

// TradeQuery query trade order by orderID
func (d *Dao) TradeQuery(ctx context.Context, orderID string) (order *model.TradeOrder, err error) {
	var (
		req = &ugcpay.TradeOrderReq{
			Id: orderID,
		}
		reply *ugcpay.TradeOrderResp
	)
	if reply, err = d.ugcpayAPI.TradeQuery(ctx, req); err != nil {
		return
	}
	order = &model.TradeOrder{
		OrderID:  reply.OrderId,
		MID:      reply.Mid,
		Biz:      reply.Biz,
		Platform: reply.Platform,
		OID:      reply.Oid,
		OType:    reply.Otype,
		Fee:      reply.Fee,
		Currency: reply.Currency,
		PayID:    reply.PayId,
		State:    reply.State,
		Reason:   reply.Reason,
	}
	return
}

// TradeConfirm confirm trade order by orderID
func (d *Dao) TradeConfirm(ctx context.Context, orderID string) (order *model.TradeOrder, err error) {
	var (
		req = &ugcpay.TradeOrderReq{
			Id: orderID,
		}
		reply *ugcpay.TradeOrderResp
	)
	if reply, err = d.ugcpayAPI.TradeConfirm(ctx, req); err != nil {
		return
	}
	order = &model.TradeOrder{
		OrderID:  reply.OrderId,
		MID:      reply.Mid,
		Biz:      reply.Biz,
		Platform: reply.Platform,
		OID:      reply.Oid,
		OType:    reply.Otype,
		Fee:      reply.Fee,
		Currency: reply.Currency,
		PayID:    reply.PayId,
		State:    reply.State,
		Reason:   reply.Reason,
	}
	return
}

// TradeCancel cancel trade order by orderID
func (d *Dao) TradeCancel(ctx context.Context, orderID string) (err error) {
	var (
		req = &ugcpay.TradeOrderReq{
			Id: orderID,
		}
	)
	if _, err = d.ugcpayAPI.TradeCancel(ctx, req); err != nil {
		return
	}
	return
}

// Income

// IncomeAssetOverview .
func (d *Dao) IncomeAssetOverview(ctx context.Context, mid int64) (inc *model.IncomeAssetOverview, err error) {
	var (
		req = &ugcpay.IncomeUserAssetOverviewReq{
			Mid: mid,
		}
		reply *ugcpay.IncomeUserAssetOverviewResp
	)
	if reply, err = d.ugcpayAPI.IncomeUserAssetOverview(ctx, req); err != nil {
		return
	}
	inc = &model.IncomeAssetOverview{
		Total:         reply.Total,
		TotalBuyTimes: reply.TotalBuyTimes,
		MonthNew:      reply.MonthNew,
		DayNew:        reply.DayNew,
	}
	return
}

// IncomeUserAssetList .
func (d *Dao) IncomeUserAssetList(ctx context.Context, mid int64, ver int64, ps, pn int64) (inc *model.IncomeAssetMonthly, err error) {
	var (
		req = &ugcpay.IncomeUserAssetListReq{
			Mid: mid,
			Ver: ver,
			Ps:  ps,
			Pn:  pn,
		}
		reply *ugcpay.IncomeUserAssetListResp
	)
	if reply, err = d.ugcpayAPI.IncomeUserAssetList(ctx, req); err != nil {
		return
	}
	inc = &model.IncomeAssetMonthly{
		List: make([]*model.IncomeAssetMonthlyByContent, 0),
		Page: &model.Page{
			Num:   reply.Page.Num,
			Size:  reply.Page.Size_,
			Total: reply.Page.Total,
		},
	}
	for _, c := range reply.List {
		inc.List = append(inc.List, &model.IncomeAssetMonthlyByContent{
			OID:           c.Oid,
			OType:         c.Otype,
			Currency:      c.Currency,
			Price:         c.Price,
			TotalBuyTimes: c.TotalBuyTimes,
			NewBuyTimes:   c.NewBuyTimes,
			TotalErrTimes: c.TotalErrTimes,
			NewErrTimes:   c.NewErrTimes,
		})
	}
	return
}

// archive

// ArchiveTitles 通过 aid list 获取稿件标题
func (d *Dao) ArchiveTitles(ctx context.Context, aids []int64) (arcTitles map[int64]string, err error) {
	var (
		req = &archive.ArcsRequest{
			Aids: aids,
		}
		reply *archive.ArcsReply
	)
	if reply, err = d.archiveAPI.Arcs(ctx, req); err != nil {
		return
	}
	arcTitles = make(map[int64]string)
	for aid, a := range reply.Arcs {
		arcTitles[aid] = a.Title
	}
	return
}
