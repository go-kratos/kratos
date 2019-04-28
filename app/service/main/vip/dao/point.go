package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_allPointExchangePriceSQL = "SELECT origin_point,current_point,month,promotion_tip,promotion_color FROM vip_point_exchange_price"
	_pointExchangePriceSQL    = "SELECT origin_point,current_point,month,promotion_tip,promotion_color FROM vip_point_exchange_price WHERE month = ?"
)

//AllPointExchangePrice .
func (d *Dao) AllPointExchangePrice(c context.Context) (pe []*model.PointExchangePrice, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _allPointExchangePriceSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PointExchangePrice)
		if err = rows.Scan(&r.OriginPoint, &r.CurrentPoint, &r.Month, &r.PromotionTip, &r.PromotionColor); err != nil {
			pe = nil
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
		pe = append(pe, r)
	}

	return
}

//PointExchangePrice def.
func (d *Dao) PointExchangePrice(c context.Context, month int16) (pe *model.PointExchangePrice, err error) {
	row := d.db.QueryRow(c, _pointExchangePriceSQL, month)
	pe = new(model.PointExchangePrice)
	if err = row.Scan(&pe.OriginPoint, &pe.CurrentPoint, &pe.Month, &pe.PromotionTip, &pe.PromotionColor); err != nil {
		if err == sql.ErrNoRows {
			pe = nil
			err = nil
		} else {
			err = errors.Wrapf(err, "scan error")
		}
	}
	return
}
