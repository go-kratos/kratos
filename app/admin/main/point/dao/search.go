package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/point/model"
	"go-common/library/database/elastic"

	"github.com/pkg/errors"
)

// PointHistory get point change history from es.
func (d *Dao) PointHistory(c context.Context, arg *model.ArgPointHistory) (res *model.SearchData, err error) {
	var changeTimeFrom, changeTimeTo string
	req := d.es.NewRequest(_searchBussinss).Index(_searchBussinss).Pn(int(arg.PN)).Ps(int(arg.PS))
	req.Fields("id", "mid", "order_id", "relation_id", "point_balance", "change_time", "change_type", "remark", "operator")
	if arg.ChangeType != 0 {
		req.WhereEq("change_type", arg.ChangeType)
	}
	if arg.Mid > 0 {
		req.WhereEq("mid", arg.Mid)
	}
	if arg.StartChangeTime != 0 {
		changeTimeFrom = time.Unix(arg.StartChangeTime, 0).Format("2006-01-02 15:04:05")
	}
	if arg.EndChangeTime != 0 {
		changeTimeTo = time.Unix(arg.EndChangeTime, 0).Format("2006-01-02 15:04:05")
	}
	req.WhereRange("change_time", changeTimeFrom, changeTimeTo, elastic.RangeScopeLcRc)
	req.Order("change_time", "desc")
	res = &model.SearchData{}
	if err = req.Scan(c, &res); err != nil {
		err = errors.WithStack(err)
	}
	return
}
