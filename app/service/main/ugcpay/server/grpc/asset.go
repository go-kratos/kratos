package grpc

import (
	"context"
	"fmt"

	"go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/app/service/main/ugcpay/model"
	"go-common/app/service/main/ugcpay/service"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// New Identify warden rpc server
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	w.Use(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if resp, err = handler(ctx, req); err == nil {
			log.Infov(ctx,
				log.KV("path", info.FullMethod),
				log.KV("caller", metadata.String(ctx, metadata.Caller)),
				log.KV("remote_ip", metadata.String(ctx, metadata.RemoteIP)),
				log.KV("args", fmt.Sprintf("%s", req)),
				log.KV("retVal", fmt.Sprintf("%s", resp)))
		}
		return
	})
	v1.RegisterUGCPayServer(w.Server(), &UGCPayServer{s})
	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

// UGCPayServer .
type UGCPayServer struct {
	svr *service.Service
}

var _ v1.UGCPayServer = &UGCPayServer{}

// AssetRegister .
func (u *UGCPayServer) AssetRegister(ctx context.Context, req *v1.AssetRegisterReq) (*v1.EmptyStruct, error) {
	err := u.svr.AssetRegister(ctx, req.Mid, req.Oid, req.Otype, req.Currency, req.Price)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyStruct{}, nil
}

// AssetQuery .
func (u *UGCPayServer) AssetQuery(ctx context.Context, req *v1.AssetQueryReq) (*v1.AssetQueryResp, error) {
	res, pp, err := u.svr.AssetQuery(ctx, req.Oid, req.Otype, req.Currency)
	if err != nil {
		return nil, err
	}
	return &v1.AssetQueryResp{
		Price:         res.Price,
		PlatformPrice: pp,
	}, nil
}

// AssetRelation .
func (u *UGCPayServer) AssetRelation(ctx context.Context, req *v1.AssetRelationReq) (*v1.AssetRelationResp, error) {
	res, err := u.svr.AssetRelation(ctx, req.Mid, req.Oid, req.Otype)
	if err != nil {
		return nil, err
	}
	return &v1.AssetRelationResp{
		State: res,
	}, nil
}

// AssetRelationDetail .
func (u *UGCPayServer) AssetRelationDetail(ctx context.Context, req *v1.AssetRelationDetailReq) (*v1.AssetRelationDetailResp, error) {
	state, err := u.svr.AssetRelation(ctx, req.Mid, req.Oid, req.Otype)
	if err != nil {
		return nil, err
	}
	res, pp, err := u.svr.AssetQuery(ctx, req.Oid, req.Otype, req.Currency)
	if err != nil {
		return nil, err
	}
	return &v1.AssetRelationDetailResp{
		RelationState:      state,
		AssetPrice:         res.Price,
		AssetPlatformPrice: pp,
	}, nil
}

// TradeCreate .
func (u *UGCPayServer) TradeCreate(ctx context.Context, req *v1.TradeCreateReq) (*v1.TradeCreateResp, error) {
	orderID, payData, err := u.svr.TradeCreate(ctx, req.Platform, req.Mid, req.Oid, req.Otype, req.Currency)
	if err != nil {
		return nil, err
	}
	return &v1.TradeCreateResp{
		OrderId: orderID,
		PayData: payData,
	}, nil
}

// TradeQuery .
func (u *UGCPayServer) TradeQuery(ctx context.Context, req *v1.TradeOrderReq) (*v1.TradeOrderResp, error) {
	order, err := u.svr.TradeQuery(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	order.State = order.ReturnState()
	return &v1.TradeOrderResp{
		OrderId:  order.OrderID,
		Mid:      order.MID,
		Biz:      order.Biz,
		Platform: order.Platform,
		Oid:      order.OID,
		Otype:    order.OType,
		Fee:      order.Fee,
		Currency: order.Currency,
		PayId:    order.PayID,
		State:    order.State,
		Reason:   order.PayReason,
	}, nil
}

// TradeCancel .
func (u *UGCPayServer) TradeCancel(ctx context.Context, req *v1.TradeOrderReq) (*v1.EmptyStruct, error) {
	err := u.svr.TradeCancel(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyStruct{}, nil
}

// TradeConfirm .
func (u *UGCPayServer) TradeConfirm(ctx context.Context, req *v1.TradeOrderReq) (*v1.TradeOrderResp, error) {
	order, err := u.svr.TradeConfirm(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	order.State = order.ReturnState()
	return &v1.TradeOrderResp{
		OrderId:  order.OrderID,
		Mid:      order.MID,
		Biz:      order.Biz,
		Platform: order.Platform,
		Oid:      order.OID,
		Otype:    order.OType,
		Fee:      order.Fee,
		Currency: order.Currency,
		PayId:    order.PayID,
		State:    order.State,
		Reason:   order.PayReason,
	}, nil
}

// TradeRefund .
func (u *UGCPayServer) TradeRefund(ctx context.Context, req *v1.TradeOrderReq) (*v1.EmptyStruct, error) {
	err := u.svr.TradeRefund(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyStruct{}, nil
}

// IncomeUserAssetOverview .
func (u *UGCPayServer) IncomeUserAssetOverview(ctx context.Context, req *v1.IncomeUserAssetOverviewReq) (*v1.IncomeUserAssetOverviewResp, error) {
	overview, monthReady, newDailyBill, err := u.svr.IncomeUserAssetOverview(ctx, req.Mid, "bp")
	if err != nil {
		return nil, err
	}
	// 总计收入
	resp := &v1.IncomeUserAssetOverviewResp{}
	if overview != nil {
		resp.Total = overview.TotalIn - overview.TotalOut
		resp.TotalBuyTimes = overview.PaySuccess - overview.PayError
	}
	// 本月新增收入
	resp.MonthNew = monthReady
	// 前日新增收入
	if newDailyBill != nil {
		resp.DayNew = newDailyBill.In - newDailyBill.Out
	}
	log.Info("IncomeUserAssetOverview grpc resp: %+v", resp)
	return resp, nil
}

// IncomeUserAssetList .
func (u *UGCPayServer) IncomeUserAssetList(ctx context.Context, req *v1.IncomeUserAssetListReq) (resp *v1.IncomeUserAssetListResp, err error) {
	// 月度收入
	if req.Ver > 0 {
		return u.incomeUserAssetListByVer(ctx, req.Mid, "bp", req.Ver, req.Pn, req.Ps)
	}
	// 总计收入
	return u.incomeUserAssetListByAll(ctx, req.Mid, "bp", req.Pn, req.Ps)
}

func (u *UGCPayServer) incomeUserAssetListByAll(ctx context.Context, mid int64, currency string, pn, ps int64) (resp *v1.IncomeUserAssetListResp, err error) {
	var (
		allList *model.AggrIncomeUserAssetList
		allPage *model.Page
	)
	// 获得总计收入
	if allList, allPage, err = u.svr.IncomeUserAssetList(ctx, mid, currency, 0, pn, ps); err != nil {
		return
	}
	// 查询 asset price
	assetPriceMap := make(map[string]int64) // assetPriceMap map[assetKey]price
	for _, a := range allList.Assets {
		as, _, err := u.svr.AssetQuery(ctx, a.OID, a.OType, a.Currency)
		if err != nil {
			log.Error("IncomeUserAssetList found invalid asset oid: %d, otype: %s, err: %+v", a.OID, a.OType, err)
			err = nil
			continue
		}
		assetPriceMap[assetKey(a.OID, a.OType, a.Currency)] = as.Price
	}
	// 写入返回值
	resp = &v1.IncomeUserAssetListResp{
		Page: &v1.Page{
			Num:   allPage.Num,
			Size_: allPage.Size,
			Total: allPage.Total,
		},
	}
	for _, a := range allList.Assets {
		asset := &v1.IncomeUserAsset{
			Oid:           a.OID,
			Otype:         a.OType,
			Currency:      a.Currency,
			Price:         assetPriceMap[assetKey(a.OID, a.OType, a.Currency)],
			TotalBuyTimes: a.PaySuccess,
			NewBuyTimes:   0,
			TotalErrTimes: a.PayError,
			NewErrTimes:   0,
		}
		resp.List = append(resp.List, asset)
	}
	return resp, nil
}

func (u *UGCPayServer) incomeUserAssetListByVer(ctx context.Context, mid int64, currency string, ver int64, pn, ps int64) (resp *v1.IncomeUserAssetListResp, err error) {
	var (
		monthList *model.AggrIncomeUserAssetList
		monthPage *model.Page
	)
	// 获得月份收入
	if monthList, monthPage, err = u.svr.IncomeUserAssetList(ctx, mid, currency, ver, pn, ps); err != nil {
		return
	}
	// 查询 asset price
	assetPriceMap := make(map[string]int64) // assetPriceMap map[assetKey]price
	for _, a := range monthList.Assets {
		as, _, err := u.svr.AssetQuery(ctx, a.OID, a.OType, a.Currency)
		if err != nil {
			log.Error("IncomeUserAssetList found invalid asset oid: %d, otype: %s, err: %+v", a.OID, a.OType, err)
			err = nil
			continue
		}
		assetPriceMap[assetKey(a.OID, a.OType, a.Currency)] = as.Price
	}
	// 写入返回值
	resp = &v1.IncomeUserAssetListResp{
		Page: &v1.Page{
			Num:   monthPage.Num,
			Size_: monthPage.Size,
			Total: monthPage.Total,
		},
	}
	for _, a := range monthList.Assets {
		var (
			asset = &v1.IncomeUserAsset{
				Oid:           a.OID,
				Otype:         a.OType,
				Currency:      a.Currency,
				Price:         0,
				TotalBuyTimes: 0,
				NewBuyTimes:   a.PaySuccess,
				TotalErrTimes: 0,
				NewErrTimes:   a.PayError,
			}
		)
		asset.Price = assetPriceMap[assetKey(a.OID, a.OType, a.Currency)]
		allAsset, err := u.svr.IncomeUserAsset(ctx, mid, a.OID, a.OType, a.Currency, 0)
		if err != nil {
			log.Error("u.svr.IncomeUserAsset mid: %d, oid: %d, otype: %s, currency: %s, err: %+v", mid, a.OID, a.OType, a.Currency, err)
			err = nil
			continue
		}
		if allAsset == nil {
			log.Error("u.svr.IncomeUserAsset got nil asset, mid: %d, oid: %d, otype: %s, currency: %s", mid, a.OID, a.OType, a.Currency)
			err = nil
			continue
		}
		asset.TotalBuyTimes = allAsset.PaySuccess
		asset.TotalErrTimes = allAsset.PayError
		resp.List = append(resp.List, asset)
	}
	return
}

func assetKey(oid int64, otype, currency string) string {
	return fmt.Sprintf("%d_%s_%s", oid, otype, currency)
}
