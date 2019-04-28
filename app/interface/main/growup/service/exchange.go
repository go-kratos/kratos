package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/interface/main/growup/dao"
	"go-common/app/interface/main/growup/model"

	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_display    = 2
	_bigVipType = 1
	_remark     = "激励兑换"
)

// GoodsState get goods new state
func (s *Service) GoodsState(c context.Context) (data interface{}, err error) {
	var newDisplay string
	goods, err := s.dao.GetDisplayGoods(c, _display)
	if err != nil {
		log.Error("s.dao.GetDisplayGoods error(%v)", err)
		return
	}
	if len(goods) > 0 {
		sort.Slice(goods, func(i, j int) bool {
			return goods[i].DisplayOnTime > goods[j].DisplayOnTime
		})
		newDisplay = fmt.Sprintf("%s+%d+%d", goods[0].ProductID, goods[0].GoodsType, goods[0].DisplayOnTime)
	}
	data = map[string]interface{}{
		"new_display":   newDisplay,
		"open_exchange": time.Now().Day() >= 8, // 1-7号不开放兑换
	}
	return
}

// GoodsShow show goods
func (s *Service) GoodsShow(c context.Context, mid int64) (goods []*model.GoodsInfo, err error) {
	goods, err = s.dao.GetDisplayGoods(c, _display)
	if err != nil {
		log.Error("s.dao.GetDisplayGoods error(%v)", err)
		return
	}
	vips, err := s.dao.ListVipProducts(c, mid)
	if err != nil {
		log.Error("s.dao.ListVipProducts error(%v)", err)
		return
	}

	// get vip price
	for _, g := range goods {
		if g.GoodsType != _bigVipType {
			continue
		}
		if v, ok := vips[g.ProductID]; ok {
			g.Month = v.Month
			g.ProductName = v.ProductName
			g.OriginPrice = v.OriginPrice
			g.CurrentPrice = int64(dao.Round(dao.Mul(float64(v.OriginPrice), float64(g.Discount)/float64(100)), 0))
		}
	}
	sort.Slice(goods, func(i, j int) bool {
		return goods[i].Month < goods[j].Month
	})
	return
}

// GoodsRecord get goods record
func (s *Service) GoodsRecord(c context.Context, mid int64, page, size int) (monthOrder map[string][]*model.GoodsOrder, count int, err error) {
	monthOrder = make(map[string][]*model.GoodsOrder)
	if page == 0 {
		page = 1
	}
	start, end := (page-1)*size, page*size

	count, err = s.dao.GetGoodsOrderCount(c, mid)
	if err != nil {
		log.Error("s.dao.GetGoodsRecordCount error(%v)", err)
		return
	}

	orders, err := s.dao.GetGoodsOrders(c, mid, start, end)
	if err != nil {
		log.Error("s.dao.GetGoodsOrders error(%v)", err)
		return
	}
	for i := 0; i < len(orders); i++ {
		month := orders[i].OrderTime.Time().Format("2006-01")
		if _, ok := monthOrder[month]; !ok {
			monthOrder[month] = make([]*model.GoodsOrder, 0)
		}
		monthOrder[month] = append(monthOrder[month], orders[i])
	}
	return
}

// GoodsBuy buy goods from creative
func (s *Service) GoodsBuy(c context.Context, mid int64, productID string, goodsType int, price int64) (err error) {
	date := time.Now()
	// 1-7号不开放兑换
	if date.Day() < 8 {
		err = ecode.GrowupGoodsTimeErr
		return
	}
	p, err := s.getProduct(c, mid, productID, goodsType)
	if err != nil {
		log.Error("s.getProduct error(%v)", err)
		return
	}
	p.CurrentPrice = int64(dao.Round(dao.Mul(float64(p.OriginPrice), float64(p.Discount)/float64(100)), 0))
	if price != p.CurrentPrice {
		err = ecode.GrowupPriceErr
		return
	}
	// check upwithdraw_income
	upAccount, totalUnwithdraw, err := s.getUpTotalUnwithdraw(c, mid)
	if err != nil {
		log.Error("s.getUpTotalUnwithdraw error(%v)", err)
		return
	}
	if totalUnwithdraw < p.CurrentPrice {
		err = ecode.GrowupPriceNotEnough
		return
	}
	uuid, err := s.sf.Generate()
	if err != nil {
		return
	}
	order := &model.GoodsOrder{
		MID:        mid,
		OrderNo:    fmt.Sprintf("DHY-%s-%d-%d", date.Format("20060102150405"), uuid, mid),
		OrderTime:  xtime.Time(date.Unix()),
		GoodsType:  p.GoodsType,
		GoodsID:    productID,
		GoodsName:  p.ProductName,
		GoodsPrice: p.CurrentPrice,
		GoodsCost:  p.OriginPrice,
	}
	// 清除up_summary redis缓存
	if err = s.dao.DelCacheKey(c, fmt.Sprintf("growup-up-summary:%d", mid)); err != nil {
		log.Error("s.dao.DelCacheKey error(%v)", err)
		return
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	// update up account
	rows, err := s.dao.TxUpdateUpAccountExchangeIncome(tx, mid, p.CurrentPrice, upAccount.Version)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.TxUpdateUpAccountExchangeIncome error(%v)", err)
		return
	}
	if rows != 1 {
		tx.Rollback()
		log.Error("s.dao.TxUpdateUpAccountExchangeIncome update rows(%d) != 1")
		err = ecode.GrowupBuyErr
		return
	}

	// generate order
	rows, err = s.dao.TxInsertGoodsOrder(tx, order)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.InsertGoodsOrder error(%v)", err)
		return
	}
	if rows != 1 {
		tx.Rollback()
		log.Error("s.dao.InsertGoodsOrder insert rows(%d) != 1", rows)
		err = ecode.GrowupBuyErr
		return
	}

	// use vip batch info
	if err = s.dao.ExchangeBigVIP(c, mid, p.ResourceID, uuid, _remark); err != nil {
		tx.Rollback()
		log.Error("s.dao.ExchangeBigVIP error(%v)", err)
		err = ecode.GrowupBuyErr
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

func (s *Service) getProduct(c context.Context, mid int64, productID string, goodsType int) (p *model.GoodsInfo, err error) {
	p, err = s.dao.GetGoodsByProductID(c, productID, goodsType)
	if err != nil {
		log.Error("s.dao.GetGoodsByProductID error(%v)", err)
		return
	}
	vips, err := s.dao.ListVipProducts(c, mid)
	if err != nil {
		log.Error("s.dao.ListVipProducts error(%v)", err)
		return
	}
	// check ResourceID
	if p.ResourceID == 0 || len(vips) == 0 {
		err = ecode.GrowupGoodsNotExist
		return
	}
	vip, ok := vips[productID]
	if !ok {
		err = ecode.GrowupGoodsNotExist
		return
	}
	p.OriginPrice = vip.OriginPrice
	p.ProductName = vip.ProductName
	return
}

func (s *Service) getUpTotalUnwithdraw(c context.Context, mid int64) (upAccount *model.UpAccount, unwithdraw int64, err error) {
	lastDay := time.Now().AddDate(0, 0, -1)
	upIncomes, err := s.dao.ListUpIncome(c, mid, "up_income", lastDay.Format(_layout), lastDay.Format(_layout))
	if err != nil {
		log.Error("s.dao.ListUpIncome error(%v)", err)
		return
	}
	var lastDayIncome int64
	for _, up := range upIncomes {
		if up.Date.Time().Format(_layout) == lastDay.Format(_layout) {
			lastDayIncome = up.Income
		}
	}

	upAccount, err = s.dao.ListUpAccount(c, mid)
	if err != nil {
		log.Error("s.dao.ListUpAccount error(%v)", err)
		return
	}
	if upAccount == nil {
		err = ecode.GrowupPriceNotEnough
		return
	}
	unwithdraw = upAccount.TotalUnwithdrawIncome - lastDayIncome
	return
}
