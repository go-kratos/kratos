package item

import (
	"context"
	"strconv"
	"strings"
	"time"

	acc "go-common/app/service/main/account/model"
	vip "go-common/app/service/main/vip/model"
	itm "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	rpc "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/dao"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/app/service/openplatform/ticket-sales/model/order_checker/account"
	"go-common/library/ecode"
)

//Checker 检查商品信息
type Checker struct {
	dao       *dao.Dao
	ItemInfos *itm.BillReply
	ac        *account.Checker
	Prices    []Prices
}

//Prices 输出订单商品价格
type Prices struct {
	Total  int64
	Pay    int64
	ExpFee int64
}

//New 新建一个检查类
func New(d *dao.Dao, ac *account.Checker) *Checker {
	return &Checker{
		dao: d,
		ac:  ac,
	}
}

//Check 检查商品信息
func (ic *Checker) Check(ctx context.Context, req *rpc.CreateOrdersRequest) (ee []ecode.Codes, err error) {
	itemIDs, scIDs, tkIDs := getItemIds(req.Orders)
	ic.ItemInfos, err = ic.dao.ItemBillInfo(ctx, itemIDs, scIDs, tkIDs)
	if err != nil {
		return
	}
	ee = make([]ecode.Codes, len(req.Orders))
	ic.Prices = make([]Prices, len(req.Orders))
begin:
	for k, v := range req.Orders {
		//检查订单类型
		if _, ok := consts.OrderTypes[v.OrderType]; !ok {
			ee[k] = ecode.TicketParamInvalid
			continue begin
		}
		//检查项目
		if e := ic.checkItem(v); e != nil {
			ee[k] = e
			continue begin
		}
		//检查购买人信息
		if e := ic.checkBuyer(v); e != nil {
			ee[k] = e
			continue begin
		}
		user := ic.ac.GetUser(v.UID)
		//检查票种
		cnt, skuIDs, e := ic.checkSKU(v, user)
		if e != nil {
			ee[k] = e
			continue begin
		}
		//检查剩余可购数
		boughtCnt, _ := ic.dao.RawBoughtCount(ctx, strconv.FormatInt(v.UID, 10), v.ProjectID, skuIDs)
		if e = ic.buyLimit(v, user, boughtCnt, cnt); e != nil {
			ee[k] = e
			continue begin
		}
		//检查商品价格
		ic.Prices[k], e = ic.checkPrice(v, user)
		if e != nil {
			ee[k] = e
			continue begin
		}
	}
	return
}

//getItemIds 依次返回订单请求中的商品ID、场次ID、票价ID
func getItemIds(orders []*rpc.CreateOrderRequest) ([]int64, []int64, []int64) {
	l := len(orders)
	itemIDsMp, scIDsMp, tkIDsMp := make(map[int64]bool, l), make(map[int64]bool, l), make(map[int64]bool, l)
	itemIDs, scIDs, tkIDs := make([]int64, l), make([]int64, l), make([]int64, l)
	var i1, i2, i3 int
	for _, o := range orders {
		if _, ok := itemIDsMp[o.ProjectID]; !ok {
			itemIDsMp[o.ProjectID] = true
			itemIDs[i1] = o.ProjectID
			i1++
		}
		if _, ok := scIDsMp[o.ScreenID]; !ok {
			scIDsMp[o.ScreenID] = true
			scIDs[i2] = o.ScreenID
			i2++
		}
		for _, sku := range o.SKUs {
			if _, ok := tkIDsMp[sku.SKUID]; !ok {
				tkIDsMp[sku.SKUID] = true
				tkIDs[i3] = sku.SKUID
				i3++
			}
		}
	}
	return itemIDs[:i1], scIDs[:i2], tkIDs[:i3]
}

func (ic *Checker) checkItem(o *rpc.CreateOrderRequest) ecode.Codes {
	//检查项目是否存在
	if _, ok := ic.ItemInfos.BaseInfo[o.ProjectID]; !ok {
		return ecode.TicketRecordLost
	}
	//todo 检查项目状态
	if _, ok := ic.ItemInfos.BillOpt[o.ProjectID]; !ok {
		ic.ItemInfos.BillOpt[o.ProjectID] = &itm.BillOpt{}
	}
	//检查场次是否存在
	scr, ok := ic.ItemInfos.BaseInfo[o.ProjectID].Screen[o.ScreenID]
	if !ok {
		return ecode.TicketRecordLost
	}
	//不选座场次必传sku
	if scr.PickSeat == consts.PickSeatNo && len(o.SKUs) == 0 {
		return ecode.TicketMissData
	}
	//检查必填配送信息
	if scr.DeliveryType == consts.DeliverTypeExpress && (o.DeliverDetail.Tel == "" || o.DeliverDetail.Addr == "" || o.DeliverDetail.Name == "") {
		return ecode.TicketMissData
	}
	return nil
}

func (ic *Checker) checkBuyer(o *rpc.CreateOrderRequest) ecode.Codes {
	sBuyer := ic.ItemInfos.BillOpt[o.ProjectID].BuyerInfo
	if sBuyer != "" && len(o.Buyers) == 0 {
		return ecode.TicketMissData
	}
	aBuyer := strings.Split(sBuyer, ",")
	for _, v := range aBuyer {
		f, _ := strconv.Atoi(v)
		if (f == consts.BuyerInfoTel && o.Buyers[0].Tel == "") ||
			(f == consts.BuyerInfoPerID && o.Buyers[0].PersonalID == "") {
			return ecode.TicketMissData
		}
	}
	return nil
}

func (ic *Checker) checkSKU(o *rpc.CreateOrderRequest, user *acc.Card) (cnt int64, skuIDs []int64, e ecode.Codes) {
	skuIDs = make([]int64, len(o.SKUs))
	now := time.Now().Unix()
	//检查每个sku可售态
	for k, sku := range o.SKUs {
		skuIDs[k] = sku.SKUID
		tk, ok := ic.ItemInfos.BaseInfo[o.ProjectID].Screen[o.ScreenID].Ticket[skuIDs[k]]
		if sku.Count <= 0 || !ok {
			e = ecode.TicketParamInvalid
			return
		}
		if tk.Time.SaleStime > now {
			e = ecode.TicketSaleNotStart
			return
		}
		if tk.Time.SaleEtime < now {
			e = ecode.TicketSaleEnd
			return
		}
		if user.Vip.Type < tk.BuyLimit {
			e = ecode.TicketNoPriv
			return
		}
		cnt += sku.Count
	}
	return
}

//buyLimit 检查订单可购买数，返回单用户可购买数与单笔订单可购买数
func (ic *Checker) buyLimit(o *rpc.CreateOrderRequest, user *acc.Card, boughtCnt int64, buyCnt int64) ecode.Codes {
	var lvLimits []*itm.BnlLevel
	var vipLimits map[int32]*itm.BnlLevel
	var oLimit, uLimit int64
	if opt, ok := ic.ItemInfos.BillOpt[o.ProjectID]; ok && opt.BuyLimit != nil {
		lvLimits = opt.BuyLimit.Level
		vipLimits = opt.BuyLimit.VIP
		oLimit = int64(opt.BuyLimit.Per)
	}
	uLevel, vipTyp := user.Level, user.Vip.Type
	if int(uLevel) >= len(lvLimits) {
		uLevel = 0
	}
	if uLimit == 0 {
		oLimit = consts.DefaultBuyNumLimit
	}
	if vipTyp == vip.NotVip || lvLimits[uLevel].ApplyToVip == 1 {
		uLimit = int64(lvLimits[uLevel].Max)
	} else {
		uLimit = int64(vipLimits[vipTyp].Max)
	}
	if buyCnt > oLimit {
		return ecode.TicketExceedLimit
	}
	if uLimit-boughtCnt-buyCnt <= 0 {
		return ecode.TicketExceedLimit
	}
	return nil
}

func (ic *Checker) checkPrice(o *rpc.CreateOrderRequest, user *acc.Card) (p Prices, e ecode.Codes) {
	opt := ic.ItemInfos.BillOpt[o.ProjectID]
	//非大会员不免邮
	if user == nil || user.Vip.Type != vip.AnnualVip || opt.VipExpFree == 0 {
		p.ExpFee = int64(opt.ExpTip)
		if p.ExpFee < 0 {
			p.ExpFee = 0
		}
	}
	p.Total = p.ExpFee
	for _, s := range o.SKUs {
		tk := ic.ItemInfos.BaseInfo[o.ProjectID].Screen[o.ScreenID].Ticket[int64(s.SKUID)]
		p.Total += int64(tk.PriceList.Price) * s.Count
	}
	p.Pay = p.Total
	if p.Pay != o.PayMoney {
		e = ecode.TicketPriceChanged
	}
	return
}
