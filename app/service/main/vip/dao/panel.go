package dao

import (
	"context"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_vipPayOrderSuccsSQL            = "SELECT id,order_type,buy_months FROM vip_pay_order WHERE mid = ? and status= 2"
	_vipPriceConfigsSQL             = "SELECT id,platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,ctime,mtime,superscript,start_build,end_build FROM vip_price_config_v2 WHERE status = 0"
	_vipPriceDiscountConfigsSQL     = "SELECT vpc_id,product_id,discount_price,stime,etime,remark,ctime,mtime,first_price FROM vip_price_discount_config_v2 WHERE stime <= ? AND ((etime > ? AND  etime <> '1970-01-01 08:00:00') OR (etime = '1970-01-01 08:00:00'))"
	_vipPriceDiscountByProductIDSQL = "SELECT vpc_id,product_id,discount_price,stime,etime,remark,ctime,mtime FROM vip_price_discount_config_v2 WHERE product_id = ? ORDER BY mtime DESC"
	_vipPriceByProductIDSQL         = "SELECT id,platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,ctime,mtime,superscript FROM vip_price_config_v2 WHERE status = 0 AND product_id = ? ORDER BY mtime DESC LIMIT 1"
	_vipPriceByIDSQL                = "SELECT id,platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,ctime,mtime,superscript,start_build,end_build FROM vip_price_config_v2 WHERE id = ?"
)

// VipPayOrderSuccs get succ of vip pay orders.
func (d *Dao) VipPayOrderSuccs(c context.Context, mid int64) (mpo map[string]struct{}, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _vipPayOrderSuccsSQL, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	mpo = make(map[string]struct{})
	for rows.Next() {
		po := new(model.PayOrder)
		if err = rows.Scan(&po.ID, &po.OrderType, &po.BuyMonths); err != nil {
			if err != sql.ErrNoRows {
				err = errors.WithStack(err)
				return
			}
			mpo = nil
			err = nil
			return
		}
		mpo[po.DoPayOrderTypeKey()] = struct{}{}
	}
	err = rows.Err()
	return
}

// VipPriceConfigs get vip price configs.
func (d *Dao) VipPriceConfigs(c context.Context) (vpcs []*model.VipPriceConfig, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _vipPriceConfigsSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vpc := new(model.VipPriceConfig)
		if err = rows.Scan(&vpc.ID, &vpc.Plat, &vpc.PdName, &vpc.PdID, &vpc.SuitType, &vpc.Month, &vpc.SubType, &vpc.OPrice,
			&vpc.Selected, &vpc.Remark, &vpc.CTime, &vpc.MTime, &vpc.Superscript, &vpc.StartBuild, &vpc.EndBuild); err != nil {
			if err != sql.ErrNoRows {
				err = errors.WithStack(err)
				return
			}
			vpcs = nil
			err = nil
			return
		}
		vpcs = append(vpcs, vpc)
	}
	err = rows.Err()
	return
}

// VipPriceDiscountConfigs get vip price discount configs.
func (d *Dao) VipPriceDiscountConfigs(c context.Context) (mvp map[int64]*model.VipDPriceConfig, err error) {
	var (
		rows *sql.Rows
		now  = time.Now()
	)
	mvp = make(map[int64]*model.VipDPriceConfig)
	if rows, err = d.db.Query(c, _vipPriceDiscountConfigsSQL, now, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vpc := new(model.VipDPriceConfig)
		if err = rows.Scan(&vpc.ID, &vpc.PdID, &vpc.DPrice, &vpc.STime, &vpc.ETime, &vpc.Remark, &vpc.CTime, &vpc.MTime, &vpc.FirstPrice); err != nil {
			if err != sql.ErrNoRows {
				err = errors.WithStack(err)
				return
			}
			mvp = nil
			err = nil
			return
		}
		mvp[vpc.ID] = vpc
	}
	err = rows.Err()
	return
}

//VipPriceDiscountByProductID select vip price discount by product id.
func (d *Dao) VipPriceDiscountByProductID(c context.Context, productID string) (vpc []*model.VipDPriceConfig, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _vipPriceDiscountByProductIDSQL, productID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := new(model.VipDPriceConfig)
		if err = rows.Scan(&v.ID, &v.PdID, &v.DPrice, &v.STime, &v.ETime, &v.Remark, &v.CTime, &v.MTime); err != nil {
			err = errors.WithStack(err)
			vpc = nil
			err = nil
			return
		}
		vpc = append(vpc, v)
	}
	err = rows.Err()
	return
}

//VipPriceByProductID select vip price by product id.
func (d *Dao) VipPriceByProductID(c context.Context, productID string) (vpc *model.VipPriceConfig, err error) {
	var row = d.db.QueryRow(c, _vipPriceByProductIDSQL, productID)
	vpc = new(model.VipPriceConfig)
	if err = row.Scan(&vpc.ID, &vpc.Plat, &vpc.PdName, &vpc.PdID, &vpc.SuitType, &vpc.Month, &vpc.SubType, &vpc.OPrice,
		&vpc.Selected, &vpc.Remark, &vpc.CTime, &vpc.MTime, &vpc.Superscript); err != nil {
		if err == sql.ErrNoRows {
			vpc = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}

//VipPriceByID vip price by id.
func (d *Dao) VipPriceByID(c context.Context, id int64) (vpc *model.VipPriceConfig, err error) {
	var row = d.db.QueryRow(c, _vipPriceByIDSQL, id)
	vpc = new(model.VipPriceConfig)
	if err = row.Scan(&vpc.ID, &vpc.Plat, &vpc.PdName, &vpc.PdID, &vpc.SuitType, &vpc.Month, &vpc.SubType, &vpc.OPrice,
		&vpc.Selected, &vpc.Remark, &vpc.CTime, &vpc.MTime, &vpc.Superscript, &vpc.StartBuild, &vpc.EndBuild); err != nil {
		if err == sql.ErrNoRows {
			vpc = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			d.errProm.Incr("row_scan_db")
		}
	}
	return
}
