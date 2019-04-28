package order

import (
	"context"
	"go-common/app/interface/main/creative/model/order"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"net/url"
	"strconv"
)

// ExecuteOrders orders.
func (d *Dao) ExecuteOrders(c context.Context, mid int64, ip string) (orders []*order.Order, err error) {
	params := url.Values{}
	params.Set("up_mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int            `json:"code"`
		Orders []*order.Order `json:"data"`
	}
	if err = d.chaodian.Get(c, d.executeOrdersURI, ip, params, &res); err != nil {
		log.Error("chaodian url(%s) response(%v) error(%v)", d.executeOrdersURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.executeOrdersURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.executeOrdersURI, res)
		orders = nil
		return
	}
	orders = res.Orders
	return
}

// Ups ups
func (d *Dao) Ups(c context.Context) (ups map[int64]int64, err error) {
	params := url.Values{}
	var res struct {
		Code int     `json:"code"`
		Ups  []int64 `json:"data"`
	}
	if err = d.chaodian.Get(c, d.upsURI, "", params, &res); err != nil {
		log.Error("chaodian url(%s) response(%v) error(%v)", d.upsURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.upsURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.upsURI, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	ups = make(map[int64]int64)
	for _, v := range res.Ups {
		ups[v] = v
	}
	return
}

// OrderByAid order by aid.
func (d *Dao) OrderByAid(c context.Context, aid int64) (orderID int64, orderName string, gameBaseID int64, err error) {
	params := url.Values{}
	params.Set("av_id", strconv.FormatInt(aid, 10))
	var res struct {
		Code  int         `json:"code"`
		Order order.Order `json:"data"`
	}
	if err = d.chaodian.Get(c, d.getOrderByAidURI, "", params, &res); err != nil {
		log.Error("chaodian url(%s) response(%v) error(%v)", d.getOrderByAidURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("chaodian (%s)", d.getOrderByAidURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.getOrderByAidURI, res)
		orderID = 0
		return
	}
	orderID = res.Order.ExeOdID
	orderName = res.Order.BzOdName
	gameBaseID = res.Order.GameBaseID
	log.Info("chaodian GetOrderByAid Res OrderInfo (%+v)", res.Order)
	return
}

// Unbind unbind order id.
func (d *Dao) Unbind(c context.Context, mid, aid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("status", "-1")
	params.Set("av_id", strconv.FormatInt(aid, 10))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.chaodian.Post(c, d.archiveStatusURI, ip, params, &res); err != nil {
		log.Error("chaodian d.chaodian.POST uri(%s) aid(%d) mid(%d) error(%v)", d.archiveStatusURI+"?"+params.Encode(), mid, aid, err)
		return
	}
	log.Info("chaodian Unbind url with params: (%s)", d.archiveStatusURI+"?"+params.Encode())
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("chaodian unbind uri(%s) aid(%d) mid(%d) res.code(%d) error(%v)", d.archiveStatusURI+"?"+params.Encode(), mid, aid, res.Code, err)
	}
	return
}

// Oasis for orders.
func (d *Dao) Oasis(c context.Context, mid int64, ip string) (oa *order.Oasis, err error) {
	params := url.Values{}
	params.Set("up_id", strconv.FormatInt(mid, 10))
	var res struct {
		Code int          `json:"code"`
		Data *order.Oasis `json:"data"`
	}
	if err = d.chaodian.Get(c, d.oasisURI, ip, params, &res); err != nil {
		log.Error("chaodian up_execute_order_statics url(%s) response(%v) error(%v)", d.oasisURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("chaodian up_execute_order_statics mid(%d) url(%s) res(%+v)", mid, d.oasisURI+"?"+params.Encode(), res)
	if res.Code != 0 {
		log.Error("chaodian up_execute_order_statics url(%s) res(%v)", d.oasisURI, res)
		oa = nil
		return
	}
	oa = res.Data
	return
}

// LaunchTime publish time from order id.
func (d *Dao) LaunchTime(c context.Context, orderID int64, ip string) (beginDate xtime.Time, err error) {
	params := url.Values{}
	params.Set("execute_order_id", strconv.FormatInt(orderID, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			BeginDate xtime.Time `json:"begin_date"`
		} `json:"data"`
	}
	if err = d.chaodian.Get(c, d.launchTimeURI, "", params, &res); err != nil {
		log.Error("chaodian LaunchTime url(%s) response(%v) error(%v)", d.launchTimeURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("chaodian LaunchTime url(%s)", d.launchTimeURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian LaunchTime url(%s) res(%v)", d.launchTimeURI, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	beginDate = res.Data.BeginDate
	return
}
