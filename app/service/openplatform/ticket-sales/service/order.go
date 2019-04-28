package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	rpc "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/dao"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	oca "go-common/app/service/openplatform/ticket-sales/model/order_checker/account"
	oci "go-common/app/service/openplatform/ticket-sales/model/order_checker/item"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
)

//OrderChecker 下单前置检查器，每个检查器需要有个Check方法，接收下单请求做为参数，返回和下单请求数一样多且顺序也一样的ecode.Codes，errcode不为0的在后续流程中拦截
type OrderChecker interface {
	Check(context.Context, *rpc.CreateOrdersRequest) ([]ecode.Codes, error)
}

//OrderLocker 下单前置锁
type OrderLocker interface {
	Lock(context.Context, *rpc.CreateOrdersRequest) ([]ecode.Codes, error)
}

//OrderProcesser 下单处理器
type OrderProcesser struct {
}

//useChecker 使用检查器
func (p *OrderProcesser) useChecker(ctx context.Context, req *rpc.CreateOrdersRequest, rules ...OrderChecker) (ee []ecode.Codes, err error) {
	if len(rules) == 0 {
		return
	}
	ee, err = rules[0].Check(ctx, req)
	if err != nil {
		return
	}
	l := len(req.Orders)
	if len(ee) != l {
		return nil, errors.New("ecode length must be equal with order length")
	}
	if len(rules) == 1 {
		return
	}
	l1 := 0
	for i := 0; i < l; i++ {
		if ee[i] == nil {
			l1++
		}
	}
	req1 := &rpc.CreateOrdersRequest{
		Orders: make([]*rpc.CreateOrderRequest, l1),
	}
	//m key:成功的ecode在原数组中的下标, val:去掉失败记录后的下标
	m := make(map[int]int, l1)
	j := 0
	for i := 0; i < l; i++ {
		if ee[i] == nil {
			req1.Orders[j] = req.Orders[i]
			m[i] = j
			j++
		}
	}
	ec1, err := p.useChecker(ctx, req1, rules[1:]...)
	if err != nil {
		return
	}
	if ec1 != nil {
		for k, v := range m {
			if ec1[v] != nil {
				ee[k] = ec1[v]
			}
		}
	}
	return
}

//ListOrders 订单列表
func (s *Service) ListOrders(ctx context.Context, req *rpc.ListOrdersRequest) (res *rpc.ListOrdersResponse, err error) {
	query := (*model.OrderMainQuerier)(req)
	cnt, err := s.dao.OrderCount(ctx, query)
	if err != nil {
		return
	}
	res = &rpc.ListOrdersResponse{}
	res.Count = cnt
	if query.Limit == 0 {
		query.Limit = dao.DefaultOrderPSize
	}
	if query.OrderBy == "" {
		query.OrderBy = dao.DefaultOrderOrderBy
	}
	list, err := s.dao.Orders(ctx, query)
	if err != nil {
		return
	}
	ll := len(list)
	if res.Count == 0 && ll > 0 {
		res.Count = int64(ll)
	}
	oids := make([]int64, ll)
	for k, v := range list {
		oids[k] = v.OrderID
	}
	dts, err := s.dao.OrderDetails(ctx, oids)
	if err != nil {
		return
	}
	skus, err := s.dao.OrderSKUs(ctx, oids)
	if err != nil {
		return
	}
	chs, err := s.dao.OrderPayCharges(ctx, oids)
	if err != nil {
		return
	}
	res.List = make([]*rpc.OrderResponse, ll)
	for k, v := range list {
		o := &rpc.OrderResponse{}
		o.OrderID = v.OrderID
		o.UID = v.UID
		o.OrderType = v.OrderType
		o.ItemID = v.ItemID
		o.ItemInfo = v.ItemInfo
		o.Count = v.Count
		o.TotalMoney = v.TotalMoney
		o.PayMoney = v.PayMoney
		o.ExpressFee = v.ExpressFee
		o.Status = v.Status
		o.SubStatus = v.SubStatus
		o.RefundStatus = v.RefundStatus
		o.CTime = v.CTime
		o.MTime = v.MTime
		o.Source = v.Source
		o.IsDeleted = v.IsDeleted
		if dt, ok := dts[v.OrderID]; ok {
			o.Detail = &rpc.OrderResponseMore{}
			o.Detail.Coupon = dt.Coupon
			o.Detail.Deliver = dt.DeliverDetail
			o.Detail.Extra = dt.Detail
			o.Detail.ExpressCO = dt.ExpressCO
			o.Detail.ExpressNO = dt.ExpressNO
			o.Detail.ExpressType = dt.ExpressType
			o.Detail.Remark = dt.Remark
			o.Detail.DeviceType = dt.DeviceType
			o.Detail.IP = 0
			o.Detail.MSource = dt.MSource
			o.Detail.Buyers = []*_type.OrderBuyer{
				{
					Name:       dt.Buyer,
					Tel:        dt.Tel,
					PersonalID: dt.PersonalID,
				},
			}
		}
		if ss, ok := skus[v.OrderID]; ok {
			o.SKUs = make([]*_type.OrderSKU, len(ss))
			for k, v := range ss {
				o.SKUs[k] = (*_type.OrderSKU)(v)
			}
		}
		if ch, ok := chs[v.OrderID]; ok {
			o.PayCharge = (*_type.OrderPayCharge)(ch)
			o.PayCharge.PayTime = v.PayTime
		}
		res.List[k] = o
	}
	return
}

//CreateOrders 创建订单
func (s *Service) CreateOrders(ctx context.Context, req *rpc.CreateOrdersRequest) (res *rpc.CreateOrdersResponse, err error) {
	p := &OrderProcesser{}
	ac := oca.New(s.dao)
	ic := oci.New(s.dao, ac)
	ee, err := p.useChecker(ctx, req, ac, ic)
	if err != nil {
		return
	}
	res = &rpc.CreateOrdersResponse{
		Result: make([]*rpc.CreateOrderResult, len(req.Orders)),
	}
	var n, i, j int
	var skuN int64
	for k, v := range req.Orders {
		if ee[k] == nil {
			n++
			for _, s := range v.SKUs {
				skuN += s.Count
			}
		} else {
			res.Result[k] = &rpc.CreateOrderResult{
				Code:    ee[k].Code(),
				Message: ee[k].Message(),
			}
		}
	}
	if n == 0 {
		return
	}
	oids, err := s.dao.OrderID(ctx, n)
	if err != nil {
		return
	}
	orders := make([]*model.OrderMain, n)
	dts := make([]*model.OrderDetail, n)
	skus := make([]*model.OrderSKU, skuN)
	if err != nil {
		return
	}
	for k, v := range req.Orders {
		if ee[k] != nil {
			continue
		}
		item := ic.ItemInfos.BaseInfo[v.ProjectID]
		opt := ic.ItemInfos.BillOpt[v.ProjectID]
		scr := item.Screen[v.ScreenID]
		orders[i] = &model.OrderMain{
			OrderID:   oids[i],
			UID:       fmt.Sprintf("%d", v.UID),
			OrderType: v.OrderType,
			ItemID:    v.ProjectID,
			ItemInfo: &_type.OrderItemInfo{
				Name:           item.Name,
				Img:            item.Img.First,
				ScreenID:       v.ScreenID,
				ScreenType:     int16(scr.Type),
				DeliverType:    int16(scr.DeliveryType),
				ExpressFee:     int64(opt.ExpTip),
				VIPExpressFree: int16(opt.VipExpFree),
				VerID:          0, //todo 补充ver_id
			},
			TotalMoney: ic.Prices[k].Total,
			PayMoney:   ic.Prices[k].Pay,
			ExpressFee: ic.Prices[k].ExpFee,
			Source:     v.Source,
			IsDeleted:  v.IsDeleted,
		}
		//非大会员不免邮
		for _, s := range v.SKUs {
			tk := scr.Ticket[int64(s.SKUID)]
			orders[i].Count += s.Count
			skus[j] = &model.OrderSKU{
				OrderID:     orders[i].OrderID,
				SKUID:       s.SKUID,
				Count:       s.Count,
				OriginPrice: int64(tk.PriceList.OriPrice),
				Price:       int64(tk.PriceList.Price),
				TicketType:  int16(tk.Type),
			}
			j++
		}
		if orders[i].PayMoney == 0 {
			orders[i].Status = consts.OrderStatusPaid
			orders[i].SubStatus = consts.SubStatusPaid
		} else {
			orders[i].Status = consts.OrderStatusUnpaid
			orders[i].SubStatus = consts.SubStatusUnpaid
		}
		dts[i] = &model.OrderDetail{
			OrderID:       orders[i].OrderID,
			IP:            net.ParseIP(metadata.String(ctx, metadata.RemoteIP)),
			DeviceType:    v.DeviceType,
			DeliverDetail: v.DeliverDetail,
		}
		if len(v.Buyers) > 0 {
			dts[i].Buyer = v.Buyers[0].Name
			dts[i].Tel = v.Buyers[0].Tel
			dts[i].PersonalID = v.Buyers[0].PersonalID
		}
		i++
	}
	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return
	}
	s.dao.TxInsertOrders(tx, orders)
	s.dao.TxInsertOrderDetails(tx, dts)
	s.dao.TxInsertOrderSKUs(tx, skus)
	tx.Rollback()
	return
}

//ListOrderLogs 日志列表
func (s *Service) ListOrderLogs(ctx context.Context, req *rpc.ListOrderLogRequest) (res *rpc.ListOrderLogResponse, err error) {
	res = new(rpc.ListOrderLogResponse)
	orderLogs, err := s.dao.GetOrderLogList(ctx, req.OrderID, req.Offset, req.Limit, req.OrderBy)
	if err != nil {
		return
	}
	res.List = append(res.List, orderLogs...)
	cnt, err := s.dao.GetOrderLogCnt(ctx, req.OrderID)
	if err != nil {
		return
	}
	res.Cnt = cnt
	return
}

//AddOrderLogs 日志添加
func (s *Service) AddOrderLogs(ctx context.Context, req *rpc.AddOrderLogRequest) (res *rpc.AddOrderLogResponse, err error) {
	res = new(rpc.AddOrderLogResponse)
	insertID, err := s.dao.AddOrderLog(ctx, req.Data)
	if err != nil {
		return
	}
	res.Id = insertID
	return
}

//GetSettleOrders 获取结算订单
func (s *Service) GetSettleOrders(ctx context.Context, dt string, ref bool, ext string, size int) (res *model.SettleOrders, err error) {
	t, err := time.Parse("2006-01-02", dt)
	if err != nil {
		return
	}
	var id int64
	if ext != "" {
		id, err = strconv.ParseInt(ext, 10, 64)
		if err != nil {
			return
		}
	}
	var offset int64
	if ref {
		res, offset, err = s.dao.RawGetSettleRefunds(ctx, t, t.Add(time.Hour*24), id, size)
	} else {
		res, offset, err = s.dao.RawGetSettleOrders(ctx, t, t.Add(time.Hour*24), id, size)
	}
	if offset > 0 {
		res.ExtParams = strconv.FormatInt(offset, 10)
	}
	return
}

//RepushSettleOrders 重推结算订单
func (s *Service) RepushSettleOrders(ctx context.Context, req interface{}) error {
	return s.dao.DatabusPub(ctx, "settle_repush", req)
}
