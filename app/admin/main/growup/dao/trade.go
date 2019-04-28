package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

// TradeDao . r/w
type TradeDao interface {
	AddGoods(c context.Context, fields string, values string) (int64, error)
	UpdateGoods(c context.Context, set string, where string, IDs []int64) (int64, error)
	GoodsCount(c context.Context, where string) (int64, error)
	GoodsList(c context.Context, where string, from, limit int) ([]*model.GoodsInfo, error)
	OrderList(c context.Context, where string, from, limit int) ([]*model.OrderInfo, error)
	OrderCount(c context.Context, where string) (int64, error)
}

const (
	_addGoodsSQL    = "INSERT INTO creative_goods(%s) VALUES %s"
	_updateGoodsSQL = "UPDATE creative_goods SET %s WHERE id in (%s)"
	_listGoodsSQL   = "SELECT id,goods_type,ex_product_id,ex_resource_id,is_display,display_on_time,discount FROM creative_goods WHERE %s ORDER By id DESC LIMIT ?,?"
	_listOrderSQL   = "SELECT mid, order_no, order_time, goods_id, goods_type, goods_name, goods_price, goods_cost FROM creative_order WHERE %s ORDER BY order_time DESC LIMIT ?,?"
	_countOrderSQL  = "SELECT count(id) FROM creative_order WHERE %s"
	_countGoodsSQL  = "SELECT count(id) FROM creative_goods WHERE %s"
)

// AddGoods .
func (d *Dao) AddGoods(c context.Context, fields string, values string) (int64, error) {
	rows, err := d.rddb.Exec(c, fmt.Sprintf(_addGoodsSQL, fields, values))
	if err != nil {
		log.Error("d.AddGoods db.Exec err(%v)", err)
		return 0, err
	}
	return rows.RowsAffected()
}

// UpdateGoods .
func (d *Dao) UpdateGoods(c context.Context, set string, where string, IDs []int64) (int64, error) {
	s := _updateGoodsSQL
	if where != "" {
		s = s + " AND " + where
	}
	rows, err := d.rddb.Exec(c, fmt.Sprintf(s, set, xstr.JoinInts(IDs)))
	if err != nil {
		log.Error("d.UpdateGoods db.Exec err(%v)", err)
		return 0, err
	}
	return rows.RowsAffected()
}

// GoodsList .
func (d *Dao) GoodsList(c context.Context, where string, from, limit int) (res []*model.GoodsInfo, err error) {
	res = make([]*model.GoodsInfo, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_listGoodsSQL, where), from, limit)
	if err != nil {
		log.Error("d.GoodsList db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := new(model.GoodsInfo)
		if err = rows.Scan(&m.ID, &m.GoodsType, &m.ProductID, &m.ResourceID, &m.IsDisplay, &m.DisplayOnTime, &m.Discount); err != nil {
			log.Error("d.GoodsList rows.Scan err(%v)", err)
			return nil, err
		}
		res = append(res, m)
	}
	err = rows.Err()
	return
}

// GoodsCount .
func (d *Dao) GoodsCount(c context.Context, where string) (total int64, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_countGoodsSQL, where)).Scan(&total)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// OrderCount .
func (d *Dao) OrderCount(c context.Context, where string) (total int64, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_countOrderSQL, where)).Scan(&total)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// OrderList .
func (d *Dao) OrderList(c context.Context, where string, from, limit int) (res []*model.OrderInfo, err error) {
	res = make([]*model.OrderInfo, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_listOrderSQL, where), from, limit)
	if err != nil {
		log.Error("d.OrderList db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := new(model.OrderInfo)
		if err = rows.Scan(&v.MID, &v.OrderNo, &v.OrderTime, &v.GoodsID, &v.GoodsType, &v.GoodsName, &v.GoodsPrice, &v.GoodsCost); err != nil {
			log.Error("d.OrderList rows.Scan err(%v)", err)
			return
		}
		res = append(res, v)
	}
	err = rows.Err()
	return
}
