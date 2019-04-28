package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/growup/model"
	vip "go-common/app/service/main/vip/model"

	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_panelType = "incentive"
)

const (
	_goodsInfoSQL       = "SELECT ex_product_id,ex_resource_id,display_on_time,goods_type,discount FROM creative_goods WHERE is_display = ?"
	_goodsResourceIDSQL = "SELECT ex_resource_id,discount,goods_type FROM creative_goods WHERE ex_product_id = ? AND goods_type = ? AND is_display = 2"

	_goodsOrderSQL      = "SELECT order_time,goods_name,goods_price,goods_cost FROM creative_order WHERE mid = ? ORDER BY order_time DESC LIMIT ?,?"
	_goodsOrderCountSQL = "SELECT count(*) FROM creative_order WHERE mid = ?"

	// insert
	_txInGoodsOrderSQL = "INSERT INTO creative_order(mid,order_no,order_time,goods_type,goods_id,goods_name,goods_price,goods_cost) VALUES(?,?,?,?,?,?,?,?)"
)

// GetGoodsByProductID get goods by product id
func (d *Dao) GetGoodsByProductID(c context.Context, productID string, goodsType int) (goods *model.GoodsInfo, err error) {
	goods = &model.GoodsInfo{}
	err = d.db.QueryRow(c, _goodsResourceIDSQL, productID, goodsType).Scan(&goods.ResourceID, &goods.Discount, &goods.GoodsType)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// GetGoodsOrderCount get order count
func (d *Dao) GetGoodsOrderCount(c context.Context, mid int64) (count int, err error) {
	err = d.db.QueryRow(c, _goodsOrderCountSQL, mid).Scan(&count)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// GetGoodsOrders get orders
func (d *Dao) GetGoodsOrders(c context.Context, mid int64, start, end int) (orders []*model.GoodsOrder, err error) {
	orders = make([]*model.GoodsOrder, 0)
	rows, err := d.db.Query(c, _goodsOrderSQL, mid, start, end)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		g := &model.GoodsOrder{}
		err = rows.Scan(&g.OrderTime, &g.GoodsName, &g.GoodsPrice, &g.GoodsCost)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		orders = append(orders, g)
	}
	err = rows.Err()
	return
}

// GetDisplayGoods get up from up_activity
func (d *Dao) GetDisplayGoods(c context.Context, isDisplay int) (goods []*model.GoodsInfo, err error) {
	goods = make([]*model.GoodsInfo, 0)
	rows, err := d.db.Query(c, _goodsInfoSQL, isDisplay)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		g := &model.GoodsInfo{}
		err = rows.Scan(&g.ProductID, &g.ResourceID, &g.DisplayOnTime, &g.GoodsType, &g.Discount)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		goods = append(goods, g)
	}
	err = rows.Err()
	return
}

// TxInsertGoodsOrder insert goods order
func (d *Dao) TxInsertGoodsOrder(tx *xsql.Tx, o *model.GoodsOrder) (rows int64, err error) {
	res, err := tx.Exec(_txInGoodsOrderSQL, o.MID, o.OrderNo, o.OrderTime, o.GoodsType, o.GoodsID, o.GoodsName, o.GoodsPrice, o.GoodsCost)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _txInGoodsOrderSQL, err)
		return
	}
	return res.RowsAffected()
}

// ListVipProducts list vip products
func (d *Dao) ListVipProducts(c context.Context, mid int64) (r map[string]*model.GoodsInfo, err error) {
	r = make(map[string]*model.GoodsInfo)
	res, err := d.vip.VipPanelInfo5(c, &vip.ArgPanel{PanelType: _panelType, Mid: mid})
	if err != nil {
		log.Error("vipRPC.VipPanelInfo5 err(%v)", err)
		return
	}
	if _, err = json.Marshal(res.Vps); err != nil {
		log.Error("json.Marshal err(%v)", err)
		return
	}
	for _, v := range res.Vps {
		m := new(model.GoodsInfo)
		m.ProductID = v.PdID
		m.ProductName = v.PdName
		// 大会员实时价格 = 激励兑换实时成本价; 单位元转换为单位分
		m.OriginPrice = int64(Round(Mul(v.DPrice, float64(100)), 0))
		m.Month = v.Month
		r[v.PdID] = m
	}
	return
}

// ExchangeBigVIP exchange big vip
func (d *Dao) ExchangeBigVIP(c context.Context, mid, resourceID, orderNo int64, remark string) (err error) {
	params := url.Values{}
	params.Set("batchId", strconv.FormatInt(resourceID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("orderNo", strconv.FormatInt(orderNo, 10))
	params.Set("remark", remark)

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	url := d.c.Host.VipURI
	if err = d.httpRead.Post(c, url, "", params, &res); err != nil {
		log.Error("d.client.Post url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("ExchangeBigVIP code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), res.Message)
		err = fmt.Errorf("ExchangeBigVIP error(%s)", res.Message)
	}
	return
}
