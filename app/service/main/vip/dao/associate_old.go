package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_activityOrderSQL       = "SELECT id,mid,order_no,product_id,months,panel_type,associate_state,ctime,mtime FROM vip_order_activity_record WHERE order_no = ?;"
	_updateActivityStateSQL = "UPDATE vip_order_activity_record SET associate_state =? WHERE order_no = ?;"
	_productBuyCountSQL     = "SELECT current_count as count FROM vip_product_pay_record WHERE mid = ? AND months=? AND panel_type = ?;"
)

//ActivityOrder get activity order by order_no.
func (d *Dao) ActivityOrder(c context.Context, orderNO string) (res *model.VipOrderActivityRecord, err error) {
	res = new(model.VipOrderActivityRecord)
	if err = d.olddb.QueryRow(c, _activityOrderSQL, orderNO).
		Scan(&res.ID, &res.Mid, &res.OrderNO, &res.ProductID, &res.Months, &res.PanelType, &res.AssociateState, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao activity order(%s)", orderNO)
	}
	return
}

// UpdateActivityState update act vip grant state.
func (d *Dao) UpdateActivityState(c context.Context, state int8, orderNO string) (aff int64, err error) {
	var res xsql.Result
	if res, err = d.olddb.Exec(c, _updateActivityStateSQL, state, orderNO); err != nil {
		err = errors.Wrapf(err, "dao update associate state(%d,%s)", state, orderNO)
		return
	}
	return res.RowsAffected()
}

// CountProductBuy get user by product count.
func (d *Dao) CountProductBuy(c context.Context, mid int64, months int32, panelType string) (count int64, err error) {
	row := d.olddb.QueryRow(c, _productBuyCountSQL, mid, months, panelType)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao update associate state(%d,%d,%s)", mid, months, panelType)
	}
	return
}
