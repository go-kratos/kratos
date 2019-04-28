package order

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/videoup/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

const (
	_executeOrders = "/api/open_api/v2/execute_orders"
	_ups           = "/api/open_api/v2/ups"
	_launchtime    = "/api/open_api/v2/execute_orders/launch_time"
	_useExeOrder   = "/api/open_api/v2/execute_orders/use"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *bm.Client
	// uri
	executeOrdersURI string
	upsURI           string
	launchTimeURI    string
	useExeOrderURI   string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                c,
		client:           bm.NewClient(c.HTTPClient.Chaodian),
		executeOrdersURI: c.Host.Chaodian + _executeOrders,
		upsURI:           c.Host.Chaodian + _ups,
		launchTimeURI:    c.Host.Chaodian + _launchtime,
		useExeOrderURI:   c.Host.Chaodian + _useExeOrder,
	}
	return
}

// ExecuteOrders execute order ids.
func (d *Dao) ExecuteOrders(c context.Context, mid int64, ip string) (orderIds map[int64]int64, err error) {
	params := url.Values{}
	params.Set("up_mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int `json:"code"`
		Orders []*struct {
			ExeOdID int64 `json:"execute_order_id"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.executeOrdersURI, ip, params, &res); err != nil {
		log.Error("chaodian url(%s) response(%+v) error(%v)", d.executeOrdersURI+"?"+params.Encode(), res, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.executeOrdersURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.executeOrdersURI, res)
		err = ecode.VideoupOrderAPIErr
		return
	}
	orderIds = make(map[int64]int64)
	for _, v := range res.Orders {
		orderIds[v.ExeOdID] = v.ExeOdID
	}
	return
}

// Ups order ups.
func (d *Dao) Ups(c context.Context) (ups map[int64]int64, err error) {
	params := url.Values{}
	var res struct {
		Code int     `json:"code"`
		Ups  []int64 `json:"data"`
	}
	if err = d.client.Get(c, d.upsURI, "", params, &res); err != nil {
		log.Error("chaodian url(%s) response(%+v) error(%v)", d.upsURI+"?"+params.Encode(), res, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.upsURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.upsURI, res)
		err = ecode.VideoupOrderAPIErr
		return
	}
	ups = make(map[int64]int64)
	for _, v := range res.Ups {
		ups[v] = v
	}
	return
}

// BindOrder bind order with up.
func (d *Dao) BindOrder(c context.Context, mid, aid, orderID int64, ip string) (err error) {
	params := url.Values{}
	params.Set("execute_order_id", strconv.FormatInt(orderID, 10))
	params.Set("av_id", strconv.FormatInt(aid, 10))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Post(c, d.useExeOrderURI, ip, params, &res); err != nil {
		log.Error("d.client.Do uri(%s) aid(%d) mid(%d) orderID(%d) error(%v)", d.useExeOrderURI+"?"+params.Encode(), mid, aid, orderID, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.useExeOrderURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("d.client.Do uri(%s) aid(%d) mid(%d) orderID(%d) res.code(%d) error(%v)", d.useExeOrderURI+"?"+params.Encode(), mid, aid, orderID, res.Code, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	return
}

// PubTime publish time from order id.
func (d *Dao) PubTime(c context.Context, mid, orderID int64, ip string) (ptime xtime.Time, err error) {
	params := url.Values{}
	params.Set("execute_order_id", strconv.FormatInt(orderID, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			BeginDate xtime.Time `json:"begin_date"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.launchTimeURI, "", params, &res); err != nil {
		log.Error("chaodian url(%s) response(%+v) error(%v)", d.launchTimeURI+"?"+params.Encode(), res, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	log.Info("chaodian url(%s)", d.launchTimeURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("chaodian url(%s) res(%v)", d.launchTimeURI, res)
		err = ecode.VideoupOrderAPIErr
		return
	}
	ptime = res.Data.BeginDate
	return
}
