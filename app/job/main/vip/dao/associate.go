package dao

import (
	"context"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_notGrantActOrder = "SELECT id,mid,order_no,product_id,months,panel_type,associate_state,ctime,mtime FROM vip_order_activity_record WHERE associate_state = 0 AND panel_type = ? limit ?;"
)

// NotGrantActOrders not grant activity order.
func (d *Dao) NotGrantActOrders(c context.Context, panelType string, limit int) (res []*model.VipOrderActivityRecord, err error) {
	var rows *sql.Rows
	if rows, err = d.oldDb.Query(c, _notGrantActOrder, panelType, limit); err != nil {
		err = errors.Wrapf(err, "dao associate not grants query (%s,%d)", panelType, limit)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.VipOrderActivityRecord)
		if err = rows.Scan(&r.ID, &r.Mid, &r.OrderNo, &r.ProductID, &r.Months, &r.PanelType, &r.AssociateState, &r.Ctime, &r.Mtime); err != nil {
			err = errors.Wrapf(err, "dao associate not grants scan (%s,%d)", panelType, limit)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
